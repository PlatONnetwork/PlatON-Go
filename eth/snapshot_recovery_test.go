package eth

import (
	"math/big"
	"os"
	"os/exec"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/consensus/misc"
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

func TestSnapshotRecoveryOnStartup(t *testing.T) {
	// The reactor, executor and SnapshotDB are process-wide singletons.
	const stageEnv = "PLATON_TEST_SNAPSHOT_STAGE"
	const modeEnv = "PLATON_TEST_SNAPSHOT_MODE"
	const datadirEnv = "PLATON_TEST_SNAPSHOT_DATADIR"
	stage := os.Getenv(stageEnv)
	if stage == "" {
		for _, mode := range []string{common.INNER_VALIDATOR_MODE, common.STATIC_VALIDATOR_MODE, "", common.PPOS_VALIDATOR_MODE} {
			name := mode
			if name == "" {
				name = "default"
			}
			t.Run(name, func(t *testing.T) {
				datadir := t.TempDir()
				for _, stage := range []string{"write", "restart"} {
					cmd := exec.Command(os.Args[0], "-test.run=^TestSnapshotRecoveryOnStartup$", "-test.timeout=45s")
					cmd.Env = append(os.Environ(), stageEnv+"="+stage, modeEnv+"="+mode, datadirEnv+"="+datadir)
					out, err := cmd.CombinedOutput()
					require.NoError(t, err, "%s: %s", stage, out)
				}
			})
		}
		return
	}

	mode := os.Getenv(modeEnv)
	stack, err := node.New(&node.Config{
		Name:    "snapshot-recovery",
		DataDir: os.Getenv(datadirEnv),
		P2P:     p2p.Config{NoDiscovery: true, NoDial: true, MaxPeers: 1},
	})
	require.NoError(t, err)
	defer stack.Close()
	snapshotdb.SetDBPathWithNode(stack.ResolvePath(snapshotdb.DBPath))

	cfg := ethconfig.Defaults
	cfg.DatabaseCache, cfg.TrieDBCache, cfg.SnapshotCache = 16, 0, 0
	cfg.TrieCleanCacheJournal, cfg.TxPool.Journal = "", ""
	cfg.NoPruning = true // Persist the post-block trie before reopening the database.
	var blsKey bls.SecretKey
	require.NoError(t, blsKey.SetLittleEndian([]byte{1}))
	cfg.CbftConfig.BlsPriKey = &blsKey
	nodeKey := stack.Config().NodeKey()
	cfg.Genesis = core.DefaultGenesisBlock()
	cfg.Genesis.Config = &params.ChainConfig{
		ChainID: big.NewInt(1337), PIP7ChainID: big.NewInt(1337), GenesisVersion: params.FORKVERSION_1_6_0, PauliBlock: big.NewInt(0),
		Cbft: &params.CbftConfig{
			Period: 20000, Amount: 10, ValidatorMode: mode,
			InitialNodes: []params.CbftNode{{Node: enode.NewV4(&nodeKey.PublicKey, nil, 0, 0), BlsPubKey: *blsKey.GetPublicKey()}},
		},
	}
	if mode == common.INNER_VALIDATOR_MODE {
		cfg.Genesis.Config.Cbft.InnerAdmin = crypto.PubkeyToAddress(nodeKey.PublicKey)
		cfg.Miner.InnerAdminKeyFile = stack.ResolvePath("inner-admin.key")
		if stage == "write" {
			require.NoError(t, crypto.SaveECDSA(cfg.Miner.InnerAdminKeyFile, nodeKey))
		}
	}

	var storedHead *types.Block
	if stage == "restart" {
		db, err := stack.OpenDatabaseWithFreezer("chaindata", 16, 16, "", "", false)
		require.NoError(t, err)
		storedHead = rawdb.ReadBlock(db, rawdb.ReadHeadBlockHash(db), 1)
		require.NoError(t, db.Close())
		require.NotNil(t, storedHead)
	}

	service, err := New(stack, &cfg)
	require.NoError(t, err)
	defer snapshotdb.Close()
	if stage == "write" {
		// Persist the executed block before the SnapshotDB commit.
		writeBlockBeforeSnapshotCommit(t, service, stack)
		require.Zero(t, snapshotdb.Instance().GetCurrent().GetHighest(true).Num.Uint64())
		return
	}

	head := service.BlockChain().CurrentBlock()
	require.Equal(t, storedHead.Hash(), head.Hash())
	require.Equal(t, storedHead.Root(), head.Root())
	_, err = service.BlockChain().StateAt(head.Root())
	require.NoError(t, err)
	highest := snapshotdb.Instance().GetCurrent().GetHighest(true)
	if mode == common.PPOS_VALIDATOR_MODE {
		require.Equal(t, head.NumberU64(), highest.Num.Uint64())
		require.Equal(t, head.Hash(), highest.Hash)
	} else {
		require.Zero(t, highest.Num.Uint64())
	}
}

func writeBlockBeforeSnapshotCommit(t *testing.T, service *Ethereum, stack *node.Node) {
	t.Helper()
	chain := service.BlockChain()
	parent := chain.CurrentBlock()
	header := &types.Header{
		ParentHash: parent.Hash(), Number: big.NewInt(1), Time: parent.Time() + 1,
		GasLimit: parent.GasLimit(), BaseFee: misc.CalcBaseFee(chain.Config(), parent.Header()),
	}
	require.NoError(t, service.Engine().Prepare(chain, header))
	reactor := core.GetReactorInstance()
	reactor.SetWorkerCoinBase(header, enode.PublicKeyToIDv0(&stack.Config().NodeKey().PublicKey))
	statedb, err := chain.StateAt(parent.Root())
	require.NoError(t, err)
	require.NoError(t, reactor.BeginBlocker(header, statedb))
	require.NoError(t, reactor.EndBlocker(header, statedb))
	block, err := service.Engine().Finalize(chain, header, statedb, nil, nil, nil)
	require.NoError(t, err)
	sealedHeader := block.Header()
	seal, err := crypto.Sign(sealedHeader.SealHash().Bytes(), stack.Config().NodeKey())
	require.NoError(t, err)
	copy(sealedHeader.Extra[len(sealedHeader.Extra)-consensus.ExtraSeal:], seal)
	block = block.WithSeal(sealedHeader)

	validators := utils.NewBitArray(1)
	validators.SetIndex(0, true)
	qc := &ctypes.QuorumCert{Epoch: state.DefaultEpoch, BlockHash: block.Hash(), BlockNumber: block.NumberU64(), ValidatorSet: validators}
	msg, err := qc.CannibalizeBytes()
	require.NoError(t, err)
	copy(qc.Signature[:], service.config.CbftConfig.BlsPriKey.Sign(string(msg)).Serialize())
	extra, err := ctypes.EncodeExtra(1, qc)
	require.NoError(t, err)
	block.SetExtraData(extra)
	require.NoError(t, chain.WriteBlockWithState(block, nil, nil, statedb, false, nil))
}
