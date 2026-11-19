package core

import (
	"math/big"
	"os"
	"os/exec"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/trie"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/handler"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/stretchr/testify/require"
)

func TestHawkingGovernanceActivationBlock(t *testing.T) {
	// The reactor, VRF handler, snapshot DB and parallel executor are node
	// singletons. Give this replay its own node lifetime.
	const child = "PLATON_TEST_HAWKING_ACTIVATION"
	if os.Getenv(child) != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestHawkingGovernanceActivationBlock$")
		cmd.Env = append(os.Environ(), child+"=1")
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "%s", out)
		return
	}
	cfg := &params.ChainConfig{
		ChainID: big.NewInt(1337), PIP7ChainID: params.PrivatePIP7ChainID,
		GenesisVersion: params.FORKVERSION_1_5_0, PauliBlock: big.NewInt(0),
	}
	params.GetEc(params.DefaultMainNet)
	snapshotdb.SetDBPathWithNode(t.TempDir())
	db := snapshotdb.Instance()
	defer db.Close()
	_, err := gov.InitGenesisGovernParam(common.Hash{}, db, cfg.GenesisVersion)
	require.NoError(t, err)
	parentState, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	require.NoError(t, err)
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	contract := common.Address{7}
	parentState.SetBalance(crypto.PubkeyToAddress(key.PublicKey), big.NewInt(1000000))
	parentState.SetCode(contract, []byte{byte(vm.PUSH1), 42, byte(vm.PUSH0), byte(vm.SSTORE), byte(vm.STOP)})
	gdb := gov.NewGovDB(db)
	proposal := &gov.VersionProposal{ProposalID: common.Hash{5}, ProposalType: gov.Version, NewVersion: params.FORKVERSION_1_6_0, ActiveBlock: 1}
	require.NoError(t, gdb.AddActiveVersion(cfg.GenesisVersion, 0, parentState))
	require.NoError(t, gdb.SetProposal(proposal, parentState))
	require.NoError(t, gdb.SetTallyResult(gov.TallyResult{ProposalID: proposal.ProposalID, Status: gov.PreActive}, parentState))
	require.NoError(t, db.PutBaseDB(gov.KeyPreActiveProposal(), common.MustRlpEncode(proposal.ProposalID)))
	require.NoError(t, db.PutBaseDB(gov.KeyPreActiveVersion(), common.MustRlpEncode(proposal.NewVersion)))
	parent := &types.Header{Number: big.NewInt(0), Root: parentState.IntermediateRoot(true)}
	require.NoError(t, db.SetCurrent(parent.Hash(), *common.Big0, *common.Big0))
	header := &types.Header{ParentHash: parent.Hash(), Number: big.NewInt(1), BaseFee: big.NewInt(1), GasLimit: 1000000, Extra: make([]byte, 32+consensus.ExtraSeal)}
	active, certain := consensus.IsHawkingEnabled(cfg, header.Number, parent, func(common.Hash) (*state.StateDB, error) { return parentState, nil })
	require.True(t, active && certain)
	require.False(t, gov.HawkingActive(cfg, parent.Number, parentState))
	engine := consensus.NewFaker()
	processors := []Processor{NewStateProcessor(cfg, nil, engine), NewParallelStateProcessor(cfg, nil, engine)}
	// Header-only sync may defer the decision. Without activation in BeginBlock,
	// even an explicitly empty body list must be rejected by both processors.
	early := types.NewBlockWithHeader(header).WithWithdrawals(types.Withdrawals{})
	for _, p := range processors {
		_, _, _, err := p.Process(early, parentState.Copy(), vm.Config{})
		require.ErrorContains(t, err, "withdrawals before hawking")
	}

	vrf := handler.NewVrfHandler(parent.Nonce.Bytes())
	vrf.SetPrivateKey(key)
	proof, err := vrf.GenerateNonce(header.Number, parent.Hash())
	require.NoError(t, err)
	header.Nonce = types.EncodeNonce(proof)
	tx, err := types.SignTx(types.NewTransaction(0, contract, big.NewInt(0), 50000, big.NewInt(2), nil), types.MakeSigner(cfg, header.Number, true), key)
	require.NoError(t, err)
	block := types.NewBlockWithWithdrawals(header, types.Transactions{tx}, nil, types.Withdrawals{}, trie.NewStackTrie(nil))
	seal := block.Header()
	sig, err := crypto.Sign(seal.SealHash().Bytes(), key)
	require.NoError(t, err)
	copy(seal.Extra[32:], sig)
	block = block.WithSeal(seal)

	govPlugin := plugin.NewGovPlugin(db)
	govPlugin.SetChainConfig(cfg)
	bcr = &BlockChainReactor{vh: vrf, validatorMode: common.PPOS_VALIDATOR_MODE,
		beginRule: []int{params.GovernanceRule}, basePluginMap: map[int]plugin.BasePlugin{params.GovernanceRule: govPlugin}}
	NewExecutor(cfg, nil, vm.Config{})
	var serialState, serialReceipts common.Hash
	var serialGas uint64
	for i, p := range processors {
		st := parentState.Copy()
		receipts, _, gas, err := p.Process(block, st, vm.Config{})
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		require.Equal(t, uint64(types.ReceiptStatusSuccessful), receipts[0].Status)
		require.Equal(t, common.BigToHash(big.NewInt(42)), common.BytesToHash(st.GetState(contract, common.ZeroHash.Bytes())))
		require.Equal(t, gas, receipts[0].GasUsed)
		require.True(t, gov.NewGov(nil).Gte160VersionState(st))
		root, receiptsRoot := st.IntermediateRoot(true), types.DeriveSha(receipts, trie.NewStackTrie(nil))
		if i == 0 {
			serialState, serialReceipts = root, receiptsRoot
			serialGas = gas
		} else {
			require.Equal(t, serialGas, gas)
			require.Equal(t, serialState, root)
			require.Equal(t, serialReceipts, receiptsRoot)
		}
	}
}

func TestInnerAdminChangeAtGovernanceHead(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	stateDB := state.NewDatabase(db)
	st, err := state.New(common.Hash{}, stateDB, nil)
	require.NoError(t, err)
	stored := &params.ChainConfig{GenesisVersion: params.FORKVERSION_1_5_0}
	next := &params.ChainConfig{GenesisVersion: stored.GenesisVersion, Cbft: &params.CbftConfig{InnerAdmin: common.Address{1}}}
	var genesisHash common.Hash
	for i, version := range []uint32{params.FORKVERSION_1_5_0, params.FORKVERSION_1_6_0} {
		require.NoError(t, gov.NewGovDB(nil).AddActiveVersion(version, uint64(i), st))
		root, err := st.Commit(true)
		require.NoError(t, err)
		require.NoError(t, stateDB.TrieDB().Commit(root, false, false))
		head := &types.Header{ParentHash: genesisHash, Root: root, Number: big.NewInt(int64(i))}
		rawdb.WriteBlock(db, types.NewBlockWithHeader(head))
		rawdb.WriteCanonicalHash(db, head.Hash(), uint64(i))
		rawdb.WriteHeadHeaderHash(db, head.Hash())
		rawdb.WriteHeadBlockHash(db, head.Hash())
		if i == 0 {
			genesisHash = head.Hash()
			rawdb.WriteChainConfig(db, genesisHash, stored)
			cfg, err := ResolveInnerAdminConfig(db, stored, next)
			require.NoError(t, err)
			require.Same(t, next, cfg)
		} else {
			// Include the startup entry point: rejection must precede a config rewrite.
			_, _, err := SetupGenesisBlock(db, nil, &Genesis{Config: next})
			require.ErrorContains(t, err, "cannot change active Cbft.InnerAdmin")
			require.Equal(t, stored.InnerAdminAddress(), rawdb.ReadChainConfig(db, genesisHash).InnerAdminAddress())
		}
		st, err = state.New(root, stateDB, nil)
		require.NoError(t, err)
	}
}

func TestInnerAdminChangeWithHeaderOnlyHead(t *testing.T) {
	for _, tc := range []struct {
		name    string
		version uint32
	}{
		{"before activation", params.FORKVERSION_1_5_0},
		{"at activation", params.FORKVERSION_1_6_0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			gen := &Genesis{
				Config:        &params.ChainConfig{GenesisVersion: params.FORKVERSION_1_5_0, Cbft: &params.CbftConfig{InnerAdmin: common.Address{1}}},
				EconomicModel: params.GetEc(params.DefaultMainNet),
			}
			genesis, err := gen.Commit(db, nil)
			require.NoError(t, err)
			oldAdmin := gen.Config.InnerAdminAddress()
			gen.Config.Cbft.InnerAdmin = common.Address{2}

			stateDB := state.NewDatabase(db)
			st, err := state.New(genesis.Root(), stateDB, nil)
			require.NoError(t, err)
			require.NoError(t, gov.NewGovDB(nil).AddActiveVersion(tc.version, 1, st))
			root, err := st.Commit(true)
			require.NoError(t, err)
			block := types.NewBlockWithHeader(&types.Header{ParentHash: genesis.Hash(), Number: big.NewInt(1), Root: root})
			rawdb.WriteHeader(db, block.Header())
			rawdb.WriteCanonicalHash(db, block.Hash(), 1)
			rawdb.WriteHeadHeaderHash(db, block.Hash())
			_, err = state.New(root, state.NewDatabase(db), nil)
			require.Error(t, err)

			cfg, _, err := SetupGenesisBlock(db, nil, gen)
			require.NoError(t, err)
			require.Equal(t, oldAdmin, cfg.InnerAdminAddress())
			require.Equal(t, oldAdmin, rawdb.ReadChainConfig(db, genesis.Hash()).InnerAdminAddress())
			require.Equal(t, genesis.Hash(), rawdb.ReadHeadBlockHash(db))
			require.Equal(t, block.Hash(), rawdb.ReadHeadHeaderHash(db))
			wrongGenesis := *gen
			wrongGenesis.Timestamp++
			_, _, err = SetupGenesisBlock(db, nil, &wrongGenesis)
			require.ErrorAs(t, err, new(*GenesisMismatchError))

			require.NoError(t, stateDB.TrieDB().Commit(root, false, false))
			rawdb.WriteBlock(db, block)
			rawdb.WriteHeadBlockHash(db, block.Hash())
			cfg, _, err = SetupGenesisBlock(db, nil, gen)
			if tc.version == params.FORKVERSION_1_6_0 {
				require.ErrorContains(t, err, "cannot change active Cbft.InnerAdmin")
				require.Equal(t, oldAdmin, rawdb.ReadChainConfig(db, genesis.Hash()).InnerAdminAddress())
			} else {
				require.NoError(t, err)
				require.Equal(t, gen.Config.InnerAdminAddress(), cfg.InnerAdminAddress())
				require.Equal(t, cfg.InnerAdminAddress(), rawdb.ReadChainConfig(db, genesis.Hash()).InnerAdminAddress())
			}
		})
	}
}
