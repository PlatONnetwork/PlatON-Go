package cbft

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"math/big"
	"time"
)

var (
	chainConfig      = params.TestnetChainConfig
	testTxPoolConfig = core.DefaultTxPoolConfig
)

type testBackend struct {
	db     ethdb.Database
	chain  *core.BlockChain
	txpool *core.TxPool
	cache  *core.BlockChainCache
	cbft   *Cbft
	worker *mockWorker
}

type NodeData struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    common.Address
	nodeID     discover.NodeID
	index      int
}

type testValidator struct {
	owner  *NodeData
	neibor []*NodeData
}

type mockWorker struct {
	mux *event.TypeMux
}

func init() {
	chainConfig.Cbft.Period = 1
	chainConfig.Cbft.Epoch = 1
	chainConfig.Cbft.MaxLatency = 1
	chainConfig.Cbft.Duration = 10
	testTxPoolConfig.Journal = ""

}

func makeViewChangeVote(pri *ecdsa.PrivateKey, timestamp, blockNum uint64, blockHash common.Hash, proposalIndex uint32,
	proposalAddr common.Address, validatorIndex uint32, validatorAddr common.Address) *viewChangeVote {
	p := &viewChangeVote{
		Timestamp:      timestamp,
		BlockNum:       blockNum,
		BlockHash:      blockHash,
		ProposalIndex:  proposalIndex,
		ProposalAddr:   proposalAddr,
		ValidatorIndex: validatorIndex,
		ValidatorAddr:  validatorAddr,
	}
	if pri != nil {
		cb, _ := p.CannibalizeBytes()
		sign, _ := crypto.Sign(cb, pri)
		p.Signature.SetBytes(sign)
	}
	return p
}

func makePrepareVote(pri *ecdsa.PrivateKey, timestamp, blockNum uint64, blockHash common.Hash, validatorIndex uint32, validatorAddr common.Address) *prepareVote {
	p := &prepareVote{
		Timestamp:      timestamp,
		Number:         blockNum,
		Hash:           blockHash,
		ValidatorIndex: validatorIndex,
		ValidatorAddr:  validatorAddr,
	}

	if pri != nil {
		cb, _ := p.CannibalizeBytes()
		sign, _ := crypto.Sign(cb, pri)
		p.Signature.SetBytes(sign)
	}
	return p
}

func makeViewChange(pri *ecdsa.PrivateKey, timestamp, baseBlockNum uint64, baseBlockHash common.Hash, proposalIndex uint32, proposalAddr common.Address, prepareVote []*prepareVote) *viewChange {
	p := &viewChange{
		Timestamp:            timestamp,
		BaseBlockNum:         baseBlockNum,
		BaseBlockHash:        baseBlockHash,
		ProposalIndex:        proposalIndex,
		ProposalAddr:         proposalAddr,
		BaseBlockPrepareVote: prepareVote,
	}

	if pri != nil {
		cb, _ := p.CannibalizeBytes()
		sign, _ := crypto.Sign(cb, pri)
		p.Signature.SetBytes(sign)
	}
	return p
}

func createBlock(pri *ecdsa.PrivateKey, parent common.Hash, number uint64) *types.Block {

	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: parent,
	}

	sign, _ := crypto.Sign(header.SealHash().Bytes(), pri)
	header.Extra = make([]byte, 32+65)
	copy(header.Extra, sign)

	block := types.NewBlockWithHeader(header)
	return block
}

func nodeIndexNow(validators *testValidator, startTimestamp int64) *NodeData {
	now := time.Now().UnixNano() / 1e6

	distance := now - startTimestamp
	duration := chainConfig.Cbft.Duration * 1000
	total := int64(len(validators.neibor) + 1)

	index := distance % (duration * total) / duration
	//
	//if distance%(duration*total)%duration != 0 {
	//	index += 1
	//}

	if index == 0 {
		return validators.owner
	}
	return validators.neibor[index]
}

func CreateCBFT(path string, pri *ecdsa.PrivateKey) *Cbft {
	ctx := node.NewServiceContext(&node.Config{DataDir: path}, nil, new(event.TypeMux), nil)

	cbft := New(chainConfig.Cbft, ctx.EventMux, ctx)
	cbft.SetPrivateKey(pri)
	return cbft
}

func NewMockWorker(mux *event.TypeMux) *mockWorker {
	m := &mockWorker{mux: mux}
	go m.loop()
	return m
}

func (m *mockWorker) loop() {
	sub := m.mux.Subscribe(cbfttypes.CbftResult{})
	for {
		select {
		case <-sub.Chan():

		}
	}
}

func CreateBackend(engine *Cbft, nodes []discover.Node) *testBackend {
	var (
		db    = ethdb.NewMemDatabase()
		gspec = core.Genesis{
			Config: chainConfig,
			Alloc:  core.GenesisAlloc{},
		}
	)

	gspec.MustCommit(db)

	chain, _ := core.NewBlockChain(db, nil, gspec.Config, engine, vm.Config{}, nil)
	cache := core.NewBlockChainCache(chain)

	engine.SetBlockChainCache(cache)
	txpool := core.NewTxPool(testTxPoolConfig, chainConfig, core.NewTxPoolBlockChain(cache))

	engine.Start(chain, txpool, NewStaticAgency(nodes))

	return &testBackend{
		db:     db,
		chain:  chain,
		cache:  cache,
		txpool: txpool,
		worker: NewMockWorker(engine.eventMux),
	}
}

func createTestValidator(accounts []*ecdsa.PrivateKey) *testValidator {
	var validators testValidator
	for i, pri := range accounts {
		if i == 0 {
			validators.owner = &NodeData{
				privateKey: pri,
				publicKey:  &pri.PublicKey,
				address:    crypto.PubkeyToAddress(pri.PublicKey),
				nodeID:     discover.PubkeyID(&pri.PublicKey),
				index:      0,
			}
			continue
		}
		validators.neibor = append(validators.neibor, &NodeData{
			privateKey: pri,
			publicKey:  &pri.PublicKey,
			address:    crypto.PubkeyToAddress(pri.PublicKey),
			nodeID:     discover.PubkeyID(&pri.PublicKey),
			index:      1,
		})
	}
	return &validators
}

func (v *testValidator) Nodes() []discover.Node {
	var nodes []discover.Node
	nodes = append(nodes, discover.Node{ID: v.owner.nodeID})
	for _, n := range v.neibor {
		nodes = append(nodes, discover.Node{ID: n.nodeID})
	}
	return nodes
}

func randomCBFT(path string, i int) (*Cbft, *testBackend, *testValidator) {
	validators := createTestValidator(createAccount(i))
	engine := CreateCBFT(path, validators.owner.privateKey)
	backend := CreateBackend(engine, validators.Nodes())
	return engine, backend, validators
}
