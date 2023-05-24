// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"errors"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/state/snapshot"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"
	cli "gopkg.in/urfave/cli.v1"
)

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// emptyCode is the known hash of the empty EVM bytecode.
	emptyCode = crypto.Keccak256(nil)
)

var (
	snapshotCommand = cli.Command{
		Name:        "snapshot",
		Usage:       "A set of commands based on the snapshot",
		Category:    "MISCELLANEOUS COMMANDS",
		Description: "",
		Subcommands: []cli.Command{
			{
				Name:      "verify-state",
				Usage:     "Recalculate state hash based on the snapshot for verification",
				ArgsUsage: "<root>",
				Action:    utils.MigrateFlags(verifyState),
				Category:  "MISCELLANEOUS COMMANDS",
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.TestnetFlag,
				},
				Description: `
platon snapshot verify-state <state-root>
will traverse the whole accounts and storages set based on the specified
snapshot and recalculate the root hash of state for verification.
In other words, this command does the snapshot to trie conversion.
`,
			},
			{
				Name:      "traverse-state",
				Usage:     "Traverse the state with given root hash for verification",
				ArgsUsage: "<root>",
				Action:    utils.MigrateFlags(traverseState),
				Category:  "MISCELLANEOUS COMMANDS",
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.TestnetFlag,
				},
				Description: `
platon snapshot traverse-state <state-root>
will traverse the whole state from the given state root and will abort if any
referenced trie node or contract code is missing. This command can be used for
state integrity verification. The default checking target is the HEAD state.

It's also usable without snapshot enabled.
`,
			},
			{
				Name:      "traverse-rawstate",
				Usage:     "Traverse the state with given root hash for verification",
				ArgsUsage: "<root>",
				Action:    utils.MigrateFlags(traverseRawState),
				Category:  "MISCELLANEOUS COMMANDS",
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.TestnetFlag,
				},
				Description: `
platon snapshot traverse-rawstate <state-root>
will traverse the whole state from the given root and will abort if any referenced
trie node or contract code is missing. This command can be used for state integrity
verification. The default checking target is the HEAD state. It's basically identical
to traverse-state, but the check granularity is smaller. 

It's also usable without snapshot enabled.
`,
			},
		},
	}
)

func verifyState(ctx *cli.Context) error {
	stack, _ := makeConfigNode(ctx)
	defer stack.Close()

	chaindb := utils.MakeChainDatabase(ctx, stack, true)
	headBlock := rawdb.ReadHeadBlock(chaindb)
	if headBlock == nil {
		log.Error("Failed to load head block")
		return errors.New("no head block")
	}
	snaptree, err := snapshot.New(chaindb, trie.NewDatabase(chaindb), 256, headBlock.Root(), false, false, false)
	if err != nil {
		log.Error("Failed to open snapshot tree", "error", err)
		return err
	}
	if ctx.NArg() > 1 {
		log.Error("Too many arguments given")
		return errors.New("too many arguments")
	}
	var root = headBlock.Root()
	if ctx.NArg() == 1 {
		root, err = parseRoot(ctx.Args()[0])
		if err != nil {
			log.Error("Failed to resolve state root", "error", err)
			return err
		}
	}
	if err := snaptree.Verify(root); err != nil {
		log.Error("Failed to verfiy state", "error", err)
		return err
	}
	log.Info("Verified the state")
	return nil
}

// traverseState is a helper function used for pruning verification.
// Basically it just iterates the trie, ensure all nodes and associated
// contract codes are present.
func traverseState(ctx *cli.Context) error {
	stack, _ := makeConfigNode(ctx)
	defer stack.Close()

	chaindb := utils.MakeChainDatabase(ctx, stack, true)
	headBlock := rawdb.ReadHeadBlock(chaindb)
	if headBlock == nil {
		log.Error("Failed to load head block")
		return errors.New("no head block")
	}
	if ctx.NArg() > 1 {
		log.Error("Too many arguments given")
		return errors.New("too many arguments")
	}
	var (
		root common.Hash
		err  error
	)
	if ctx.NArg() == 1 {
		root, err = parseRoot(ctx.Args()[0])
		if err != nil {
			log.Error("Failed to resolve state root", "error", err)
			return err
		}
		log.Info("Start traversing the state", "root", root)
	} else {
		root = headBlock.Root()
		log.Info("Start traversing the state", "root", root, "number", headBlock.NumberU64())
	}
	triedb := trie.NewDatabase(chaindb)
	t, err := trie.NewSecure(root, triedb)
	if err != nil {
		log.Error("Failed to open trie", "root", root, "error", err)
		return err
	}
	var (
		accounts   int
		slots      int
		codes      int
		lastReport time.Time
		start      = time.Now()
	)
	accIter := trie.NewIterator(t.NodeIterator(nil))
	for accIter.Next() {
		accounts += 1
		var acc state.Account
		if err := rlp.DecodeBytes(accIter.Value, &acc); err != nil {
			log.Error("Invalid account encountered during traversal", "error", err)
			return err
		}
		if acc.Root != emptyRoot {
			storageTrie, err := trie.NewSecure(acc.Root, triedb)
			if err != nil {
				log.Error("Failed to open storage trie", "root", acc.Root, "error", err)
				return err
			}
			storageIter := trie.NewIterator(storageTrie.NodeIterator(nil))
			for storageIter.Next() {
				slots += 1
			}
			if storageIter.Err != nil {
				log.Error("Failed to traverse storage trie", "root", acc.Root, "error", storageIter.Err)
				return storageIter.Err
			}
		}
		if !bytes.Equal(acc.CodeHash, emptyCode) {
			code := rawdb.ReadCode(chaindb, common.BytesToHash(acc.CodeHash))
			if len(code) == 0 {
				log.Error("Code is missing", "hash", common.BytesToHash(acc.CodeHash))
				return errors.New("missing code")
			}
			codes += 1
		}
		if time.Since(lastReport) > time.Second*8 {
			log.Info("Traversing state", "accounts", accounts, "slots", slots, "codes", codes, "elapsed", common.PrettyDuration(time.Since(start)))
			lastReport = time.Now()
		}
	}
	if accIter.Err != nil {
		log.Error("Failed to traverse state trie", "root", root, "error", accIter.Err)
		return accIter.Err
	}
	log.Info("State is complete", "accounts", accounts, "slots", slots, "codes", codes, "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

// traverseRawState is a helper function used for pruning verification.
// Basically it just iterates the trie, ensure all nodes and associated
// contract codes are present. It's basically identical to traverseState
// but it will check each trie node.
func traverseRawState(ctx *cli.Context) error {
	stack, _ := makeConfigNode(ctx)
	defer stack.Close()

	chaindb := utils.MakeChainDatabase(ctx, stack, true)
	headBlock := rawdb.ReadHeadBlock(chaindb)
	if headBlock == nil {
		log.Error("Failed to load head block")
		return errors.New("no head block")
	}
	if ctx.NArg() > 1 {
		log.Error("Too many arguments given")
		return errors.New("too many arguments")
	}
	var (
		root common.Hash
		err  error
	)
	if ctx.NArg() == 1 {
		root, err = parseRoot(ctx.Args()[0])
		if err != nil {
			log.Error("Failed to resolve state root", "error", err)
			return err
		}
		log.Info("Start traversing the state", "root", root)
	} else {
		root = headBlock.Root()
		log.Info("Start traversing the state", "root", root, "number", headBlock.NumberU64())
	}
	triedb := trie.NewDatabase(chaindb)
	t, err := trie.NewSecure(root, triedb)
	if err != nil {
		log.Error("Failed to open trie", "root", root, "error", err)
		return err
	}
	var (
		nodes      int
		accounts   int
		slots      int
		codes      int
		lastReport time.Time
		start      = time.Now()
	)
	accIter := t.NodeIterator(nil)
	for accIter.Next(true) {
		nodes += 1
		node := accIter.Hash()

		if node != (common.Hash{}) {
			// Check the present for non-empty hash node(embedded node doesn't
			// have their own hash).
			blob := rawdb.ReadTrieNode(chaindb, node)
			if len(blob) == 0 {
				log.Error("Missing trie node(account)", "hash", node)
				return errors.New("missing account")
			}
		}
		// If it's a leaf node, yes we are touching an account,
		// dig into the storage trie further.
		if accIter.Leaf() {
			accounts += 1
			var acc state.Account
			if err := rlp.DecodeBytes(accIter.LeafBlob(), &acc); err != nil {
				log.Error("Invalid account encountered during traversal", "error", err)
				return errors.New("invalid account")
			}
			if acc.Root != emptyRoot {
				storageTrie, err := trie.NewSecure(acc.Root, triedb)
				if err != nil {
					log.Error("Failed to open storage trie", "root", acc.Root, "error", err)
					return errors.New("missing storage trie")
				}
				storageIter := storageTrie.NodeIterator(nil)
				for storageIter.Next(true) {
					nodes += 1
					node := storageIter.Hash()

					// Check the present for non-empty hash node(embedded node doesn't
					// have their own hash).
					if node != (common.Hash{}) {
						blob := rawdb.ReadTrieNode(chaindb, node)
						if len(blob) == 0 {
							log.Error("Missing trie node(storage)", "hash", node)
							return errors.New("missing storage")
						}
					}
					// Bump the counter if it's leaf node.
					if storageIter.Leaf() {
						slots += 1
					}
				}
				if storageIter.Error() != nil {
					log.Error("Failed to traverse storage trie", "root", acc.Root, "error", storageIter.Error())
					return storageIter.Error()
				}
			}
			if !bytes.Equal(acc.CodeHash, emptyCode) {
				code := rawdb.ReadCode(chaindb, common.BytesToHash(acc.CodeHash))
				if len(code) == 0 {
					log.Error("Code is missing", "account", common.BytesToHash(accIter.LeafKey()))
					return errors.New("missing code")
				}
				codes += 1
			}
			if time.Since(lastReport) > time.Second*8 {
				log.Info("Traversing state", "nodes", nodes, "accounts", accounts, "slots", slots, "codes", codes, "elapsed", common.PrettyDuration(time.Since(start)))
				lastReport = time.Now()
			}
		}
	}
	if accIter.Error() != nil {
		log.Error("Failed to traverse state trie", "root", root, "error", accIter.Error())
		return accIter.Error()
	}
	log.Info("State is complete", "nodes", nodes, "accounts", accounts, "slots", slots, "codes", codes, "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

func parseRoot(input string) (common.Hash, error) {
	var h common.Hash
	if err := h.UnmarshalText([]byte(input)); err != nil {
		return h, err
	}
	return h, nil
}
