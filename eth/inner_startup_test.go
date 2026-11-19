package eth

import (
	"math/big"
	"os"
	"os/exec"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/eth/ethconfig"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/stretchr/testify/require"
)

func TestInnerSignerStartup(t *testing.T) {
	// The blockchain reactor, executor and snapshot DB have process-wide lifetimes.
	const child = "PLATON_TEST_INNER_SIGNER_ROLE"
	role := os.Getenv(child)
	if role == "" {
		for _, role := range []string{"missing-key", "producer", "follower"} {
			t.Run(role, func(t *testing.T) {
				cmd := exec.Command(os.Args[0], "-test.run=^TestInnerSignerStartup$", "-test.timeout=45s")
				cmd.Env = append(os.Environ(), child+"="+role, "PLATON_TEST_INNER_DATADIR="+t.TempDir())
				out, err := cmd.CombinedOutput()
				require.NoError(t, err, "%s", out)
			})
		}
		return
	}
	stack, err := node.New(&node.Config{Name: "inner", DataDir: os.Getenv("PLATON_TEST_INNER_DATADIR"), P2P: p2p.Config{NoDiscovery: true, NoDial: true, MaxPeers: 1}})
	require.NoError(t, err)
	defer stack.Close()
	path := stack.ResolvePath(snapshotdb.DBPath)
	snapshotdb.SetDBPathWithNode(path)
	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	cfg := ethconfig.Defaults
	cfg.DatabaseCache, cfg.TrieDBCache, cfg.SnapshotCache = 16, 0, 0
	cfg.TrieCleanCacheJournal, cfg.TxPool.Journal = "", ""
	cfg.CbftConfig.BlsPriKey = &blsKey
	cfg.Genesis = core.DefaultGenesisBlock()
	validatorKey := stack.Config().NodeKey()
	if role == "follower" {
		validatorKey, err = crypto.GenerateKey()
		require.NoError(t, err)
	}
	adminKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	oldAdmin := crypto.PubkeyToAddress(adminKey.PublicKey)
	if role == "producer" {
		cfg.Miner.InnerAdminKeyFile = stack.ResolvePath("inner-admin.key")
		require.NoError(t, crypto.SaveECDSA(cfg.Miner.InnerAdminKeyFile, adminKey))
	}
	cfg.Genesis.Config = &params.ChainConfig{
		ChainID: big.NewInt(1337), PIP7ChainID: big.NewInt(1337), GenesisVersion: params.FORKVERSION_1_5_0,
		Cbft: &params.CbftConfig{Period: 20000, Amount: 10, InnerAdmin: oldAdmin, ValidatorMode: common.INNER_VALIDATOR_MODE, InitialNodes: []params.CbftNode{{
			Node: enode.NewV4(&validatorKey.PublicKey, nil, 0, 0), BlsPubKey: *blsKey.GetPublicKey(),
		}}},
	}

	// Simulate interrupted header sync with neither body nor state at the tip.
	db, err := stack.OpenDatabaseWithFreezer("chaindata", 16, 16, "", "", false)
	require.NoError(t, err)
	sdb, err := snapshotdb.Open(path, 16, 16, true, false)
	require.NoError(t, err)
	genesis, err := cfg.Genesis.Commit(db, sdb)
	require.NoError(t, err)
	header := &types.Header{ParentHash: genesis.Hash(), Number: big.NewInt(1), Root: common.Hash{1}}
	rawdb.WriteHeader(db, header)
	rawdb.WriteCanonicalHash(db, header.Hash(), 1)
	rawdb.WriteHeadHeaderHash(db, header.Hash())
	require.NoError(t, db.Close())
	require.NoError(t, sdb.Close())
	cfg.Genesis.Config.Cbft.InnerAdmin = common.Address{2}

	service, err := New(stack, &cfg)
	if role != "missing-key" {
		require.NoError(t, err)
		require.Equal(t, oldAdmin, service.BlockChain().Config().InnerAdminAddress())
		require.Equal(t, genesis.Hash(), service.BlockChain().CurrentBlock().Hash())
		require.Equal(t, header.Hash(), service.BlockChain().CurrentHeader().Hash())
		require.NoError(t, stack.Start())
		return
	}
	require.ErrorContains(t, err, "configure Miner.InnerAdminKeyFile")
	// Reopen every persistent store before closing the node or exiting the process.
	db, err = stack.OpenDatabaseWithFreezer("chaindata", 16, 16, "", "", false)
	require.NoError(t, err)
	require.Equal(t, oldAdmin, rawdb.ReadChainConfig(db, genesis.Hash()).InnerAdminAddress())
	require.Equal(t, header.Hash(), rawdb.ReadHeadHeaderHash(db))
	require.NoError(t, db.Close())
	sdb, err = snapshotdb.Open(path, 16, 16, true, false)
	require.NoError(t, err)
	require.NoError(t, sdb.Close())
	w, err := wal.NewWal(stack, "")
	require.NoError(t, err)
	w.Close()
	pool, err := evidence.NewEvidencePool(stack.ResolvePath, cfg.CbftConfig.EvidenceDir)
	require.NoError(t, err)
	pool.Close()
}
