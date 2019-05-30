package cbft

import (
	"crypto/ecdsa"
	"fmt"
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
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/deckarep/golang-set"
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
	owner     *NodeData
	neighbors []*NodeData
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

func makeConfirmedBlock(v *testValidator, root common.Hash, view *viewChange, num int) []*BlockExt {
	blocks := make([]*BlockExt, 0)
	parentHash := view.BaseBlockHash
	for i := uint64(1); i <= uint64(num); i++ {
		block := createBlockWithRootHash(v.validator(view.ProposalIndex).privateKey, parentHash, root, view.BaseBlockNum+i)
		ext := NewBlockExt(block, block.NumberU64(), v.len())
		ext.view = view
		for j := uint32(0); j < uint32(v.len()); j++ {
			if j != view.ProposalIndex {
				ext.prepareVotes.Add(makePrepareVote(v.validator(j).privateKey, view.Timestamp, block.NumberU64(), block.Hash(), j, v.validator(j).address))
				ext.viewChangeVotes = append(ext.viewChangeVotes, makeViewChangeVote(v.validator(j).privateKey, view.Timestamp, view.BaseBlockNum, view.BaseBlockHash, view.ProposalIndex, view.ProposalAddr, j, v.validator(j).address))
			}
		}
		extra := []byte{cbftVersion}
		bxBytes, _ := rlp.EncodeToBytes(ext.BlockExtra())
		extra = append(extra, bxBytes...)
		block.SetExtraData(extra)
		parentHash = block.Hash()
		blocks = append(blocks, ext)
	}
	return blocks
}

func createBlock(pri *ecdsa.PrivateKey, parent common.Hash, number uint64) *types.Block {

	header := &types.Header{
		Number:      big.NewInt(int64(number)),
		ParentHash:  parent,
		ReceiptHash: types.EmptyRootHash,
		TxHash:      types.EmptyRootHash,
	}

	sign, _ := crypto.Sign(header.SealHash().Bytes(), pri)
	header.Extra = make([]byte, 32+65)
	copy(header.Extra, sign)

	block := types.NewBlockWithHeader(header)
	return block
}

func createBlockWithRootHash(pri *ecdsa.PrivateKey, parent common.Hash, root common.Hash, number uint64) *types.Block {

	header := &types.Header{
		Number:      big.NewInt(int64(number)),
		ParentHash:  parent,
		ReceiptHash: types.EmptyRootHash,
		TxHash:      types.EmptyRootHash,
		Root:        root,
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
	total := int64(len(validators.neighbors) + 1)

	index := distance % (duration * total) / duration
	//
	//if distance%(duration*total)%duration != 0 {
	//	index += 1
	//}
	if index > int64(len(validators.neighbors)) {
		panic(fmt.Sprintf("now:%d, distance:%d, duration:%d, total:%d, index:%d", now, distance, duration, total, index))
	}
	if index == 0 {
		return validators.owner
	}
	return validators.neighbors[index-1]
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
		validators.neighbors = append(validators.neighbors, &NodeData{
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
	for _, n := range v.neighbors {
		nodes = append(nodes, discover.Node{ID: n.nodeID})
	}
	return nodes
}

func (v *testValidator) validator(index uint32) *NodeData {
	if index == 0 {
		return v.owner
	}
	return v.neighbors[index-1]
}
func (v *testValidator) len() int {
	return len(v.neighbors) + 1
}

func randomCBFT(path string, i int) (*Cbft, *testBackend, *testValidator) {
	validators := createTestValidator(createAccount(i))
	engine := CreateCBFT(path, validators.owner.privateKey)
	backend := CreateBackend(engine, validators.Nodes())
	return engine, backend, validators
}

func makeHandler(cbft *Cbft, pid string, msgHash common.Hash) (*handler) {
	handler := NewHandler(cbft)
	peerSets := newPeerSet()
	peer := &peer{
		id: pid,
		knownMessageHash: mapset.NewSet(),
	}
	peer.MarkMessageHash(msgHash)
	peerSets.Register(peer)
	handler.peers = peerSets
	return handler
}

func makeGetPrepareVote(blockNum uint64, blockHash common.Hash) *getPrepareVote {
	p := &getPrepareVote{
		Number:         blockNum,
		Hash:           blockHash,
		VoteBits: 		NewBitArray(32),
	}
	return p
}

func makePrepareVotes(pri *ecdsa.PrivateKey, timestamp, blockNum uint64, blockHash common.Hash, validatorIndex uint32, validatorAddr common.Address) *prepareVotes {
	pv := makePrepareVote(pri, timestamp, blockNum, blockHash, validatorIndex, validatorAddr)
	pvs := &prepareVotes{
		Hash: blockHash,
		Number: blockNum,
		Votes: []*prepareVote{ pv },
	}
	return pvs
}