// Copyright 2015 The go-ethereum Authors
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
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"io"

	"os"
	"strconv"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

var (
	initCommand = cli.Command{
		Action:    utils.MigrateFlags(initGenesis),
		Name:      "init",
		Usage:     "Bootstrap and initialize a new genesis block",
		ArgsUsage: "<genesisPath>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The init command initializes a new genesis block and definition for the network.
This is a destructive action and changes the network in which you will be
participating.

It expects the genesis file as argument.`,
	}
	dumpGenesisCommand = cli.Command{
		Action:    utils.MigrateFlags(dumpGenesis),
		Name:      "dumpgenesis",
		Usage:     "Dumps genesis block JSON configuration to stdout",
		ArgsUsage: "",
		Flags: []cli.Flag{
			utils.DataDirFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The dumpgenesis command dumps the genesis block configuration in JSON format to stdout.`,
	}
	importPreimagesCommand = cli.Command{
		Action:    utils.MigrateFlags(importPreimages),
		Name:      "import-preimages",
		Usage:     "Import the preimage database from an RLP stream",
		ArgsUsage: "<datafile>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
	The import-preimages command imports hash preimages from an RLP encoded stream.`,
	}
	exportPreimagesCommand = cli.Command{
		Action:    utils.MigrateFlags(exportPreimages),
		Name:      "export-preimages",
		Usage:     "Export the preimage database into an RLP stream",
		ArgsUsage: "<dumpfile>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The export-preimages command export hash preimages to an RLP encoded stream`,
	}
	dumpCommand = cli.Command{
		Action:    utils.MigrateFlags(dump),
		Name:      "dump",
		Usage:     "Dump a specific block from storage",
		ArgsUsage: "[<blockHash> | <blockNum>]...",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The arguments are interpreted as block numbers or hashes.
Use "ethereum dump 0" to dump the genesis block.`,
	}
)

// initGenesis will initialise the given JSON format genesis file and writes it as
// the zero'd block (i.e. genesis) or will fail hard if it can't succeed.
func initGenesis(ctx *cli.Context) error {
	// Make sure we have a valid genesis JSON
	genesisPath := ctx.Args().First()
	if len(genesisPath) == 0 {
		utils.Fatalf("Must supply path to genesis JSON file")
	}

	genesis := new(core.Genesis)
	if err := genesis.InitGenesisAndSetEconomicConfig(genesisPath); err != nil {
		utils.Fatalf(err.Error())
	}

	// Open an initialise both full and light databases
	stack, _ := makeConfigNode(ctx)
	defer stack.Close()

	for _, name := range []string{"chaindata", "lightchaindata"} {
		chaindb, err := stack.OpenDatabase(name, 0, 0, "", false)
		if err != nil {
			utils.Fatalf("Failed to open database: %v", err)
		}
		var sdb snapshotdb.DB
		if name == "chaindata" {
			sdb, err = snapshotdb.Open(stack.ResolvePath(snapshotdb.DBPath), 0, 0, true)
			if err != nil {
				utils.Fatalf("Failed to open snapshotdb: %v", err)
			}

		}
		_, hash, err := core.SetupGenesisBlock(chaindb, sdb, genesis)
		if err != nil {
			utils.Fatalf("Failed to write genesis block: %v", err)
		}
		log.Info("Successfully wrote genesis state", "database", name, "hash", hash.Hex())
		if sdb != nil {
			if err := sdb.Close(); err != nil {
				utils.Fatalf("close base db fail: %v", err)
			}
		}
		chaindb.Close()
	}
	genesisFile, err := os.Create(stack.GenesisPath())
	if err != nil {
		utils.Fatalf("Failed create Genesis file: %v", err)
	}
	defer genesisFile.Close()

	file, err := os.Open(genesisPath)
	if err != nil {
		utils.Fatalf("Failed to read genesis file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(genesisFile, file); err != nil {
		utils.Fatalf("Failed Copy Genesis file: %v", err)
	}

	return nil
}

func dumpGenesis(ctx *cli.Context) error {
	genesis := utils.MakeGenesis(ctx)
	if genesis == nil {
		genesis = core.DefaultGenesisBlock()
	}
	if err := json.NewEncoder(os.Stdout).Encode(genesis); err != nil {
		utils.Fatalf("could not encode genesis")
	}
	return nil
}

// importPreimages imports preimage data from the specified file.
func importPreimages(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}
	stack, _ := makeFullNode(ctx)
	defer stack.Close()

	diskdb := utils.MakeChainDatabase(ctx, stack, false)

	start := time.Now()
	if err := utils.ImportPreimages(diskdb, ctx.Args().First()); err != nil {
		utils.Fatalf("Import error: %v\n", err)
	}
	fmt.Printf("Import done in %v\n", time.Since(start))
	return nil
}

// exportPreimages dumps the preimage data to specified json file in streaming way.
func exportPreimages(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}
	stack, _ := makeFullNode(ctx)
	defer stack.Close()

	diskdb := utils.MakeChainDatabase(ctx, stack, true)
	start := time.Now()

	if err := utils.ExportPreimages(diskdb, ctx.Args().First()); err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}
	fmt.Printf("Export done in %v\n", time.Since(start))
	return nil
}

func dump(ctx *cli.Context) error {
	stack, _ := makeFullNode(ctx)
	defer stack.Close()

	opts := &state.DumpConfig{
		OnlyWithAddresses: true,
		Max:               eth.AccountRangeMaxResults, // Sanity limit over RPC
	}

	chain, chainDb := utils.MakeChain(ctx, stack)
	for _, arg := range ctx.Args() {
		var block *types.Block
		if hashish(arg) {
			block = chain.GetBlockByHash(common.HexToHash(arg))
		} else {
			num, _ := strconv.Atoi(arg)
			block = chain.GetBlockByNumber(uint64(num))
		}
		if block == nil {
			fmt.Println("{}")
			utils.Fatalf("block not found")
		} else {
			state, err := state.New(block.Root(), state.NewDatabase(chainDb), nil)
			if err != nil {
				utils.Fatalf("could not create new state: %v", err)
			}
			fmt.Printf("%s\n", state.Dump(opts))
		}
	}
	chainDb.Close()
	return nil
}

// hashish returns true for strings that look like hashes.
func hashish(x string) bool {
	_, err := strconv.Atoi(x)
	return err != nil
}
