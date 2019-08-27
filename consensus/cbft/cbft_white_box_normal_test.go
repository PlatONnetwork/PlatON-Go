package cbft

import (
	"strings"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/stretchr/testify/assert"
)

const (
	overIndexLimit       = "blockIndex higher than amount"
	existIndex           = "blockIndex already exists"
	firstBlockNotQCChild = "the first index block is not contiguous by local highestQC or highestLock"
	notExistPreIndex     = "previous index block not exists"
	diffPreIndexBlock    = "non contiguous index block"
	viewNumberTooLow     = "viewNumber too low"

	needChangeView = "need change view"
	needFetchBlock = "viewNumber higher than local"
)

func init() {
	bls.Init(bls.CurveFp254BNb)
}

func TestPrepareRules(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(1)
	node := MockNode(pk[0], sk[0], cbftnodes, 10000, 10)
	assert.Nil(t, node.Start())

	result := make(chan *types.Block, 1)

	parent := node.chain.Genesis()
	var blocks []*types.Block
	for i := 0; i < 6; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, node.engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		node.engine.OnSeal(block, result, nil)

		select {
		case b := <-result:
			assert.NotNil(t, b)
			parent = b
			blocks = append(blocks, b)
		}
	}
	assert.Equal(t, uint32(6), node.engine.state.NumViewBlocks())

	// test PrepareRules check index
	testPrepareRulesCheckIndex(t, node, blocks)
}

func testPrepareRulesCheckIndex(t *testing.T, node *TestCBFT, blocks []*types.Block) {
	// over index limit
	prepare := &protocols.PrepareBlock{
		Epoch:         node.engine.state.Epoch(),
		ViewNumber:    node.engine.state.ViewNumber(),
		Block:         blocks[0],
		BlockIndex:    10,
		ProposalIndex: 0,
	}
	err := node.engine.OnPrepareBlock("", prepare)
	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), overIndexLimit))
	}

	// exist index
	prepare.BlockIndex = uint32(len(blocks) - 1)
	err = node.engine.OnPrepareBlock("", prepare)
	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), existIndex))
	}

	// not exist previous index
	//prepare.BlockIndex = uint32(len(blocks) + 1)
	//err = node.engine.OnPrepareBlock("", prepare)
	//assert.NotNil(t, err)
	//if err != nil {
	//	assert.True(t, strings.HasPrefix(err.Error(), notExistPreIndex))
	//}

	// different previous index block
	prepare.BlockIndex = uint32(len(blocks))
	err = node.engine.OnPrepareBlock("", prepare)
	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), diffPreIndexBlock))
	}

	// firstBlock not qc child or lock child
	hash, number := node.engine.state.HighestQCBlock().Hash(), node.engine.state.HighestQCBlock().NumberU64()
	block, qc := node.engine.blockTree.FindBlockAndQC(hash, number)
	node.engine.changeView(node.engine.state.Epoch(), node.engine.state.ViewNumber()+1, block, qc, nil)
	b := NewBlock(blocks[len(blocks)-3].Hash(), blocks[len(blocks)-3].NumberU64()+1)
	prepare = &protocols.PrepareBlock{
		Epoch:         node.engine.state.Epoch(),
		ViewNumber:    node.engine.state.ViewNumber(),
		Block:         b,
		BlockIndex:    0,
		ProposalIndex: 0,
	}
	err = node.engine.OnPrepareBlock("", prepare)
	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), firstBlockNotQCChild))
	}

	// viewNumber too low
	prepare.ViewNumber = node.engine.state.ViewNumber() - 1
	err = node.engine.OnPrepareBlock("", prepare)
	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), viewNumberTooLow))
	}
}

func TestPrepareBehind(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(1)
	node := MockNode(pk[0], sk[0], cbftnodes, 10000, 10)
	assert.Nil(t, node.Start())

	result := make(chan *types.Block, 1)

	parent := node.chain.Genesis()
	var blocks []*types.Block
	for i := 0; i < 9; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, node.engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		node.engine.OnSeal(block, result, nil)

		select {
		case b := <-result:
			assert.NotNil(t, b)
			parent = b
			blocks = append(blocks, b)
		}
	}
	assert.Equal(t, uint32(9), node.engine.state.NumViewBlocks())

	// test arrive new view
	testPrepareNewView(t, node, blocks)
	// test need fetch block
	testPrepareFetch(t, node, blocks)
}

func testPrepareNewView(t *testing.T, node *TestCBFT, blocks []*types.Block) {
	// arrive new view
	b := NewBlock(blocks[len(blocks)-1].Hash(), blocks[len(blocks)-1].NumberU64()+1)
	vqc := &ctypes.ViewChangeQuorumCert{
		Epoch:       node.engine.state.Epoch(),
		ViewNumber:  node.engine.state.ViewNumber(),
		BlockHash:   blocks[len(blocks)-1].Hash(),
		BlockNumber: blocks[len(blocks)-1].NumberU64(),
	}
	prepare := &protocols.PrepareBlock{
		Epoch:         node.engine.state.Epoch(),
		ViewNumber:    node.engine.state.ViewNumber() + 1,
		Block:         b,
		BlockIndex:    0,
		ProposalIndex: 0,
		ViewChangeQC: &ctypes.ViewChangeQC{
			QCs: []*ctypes.ViewChangeQuorumCert{vqc},
		},
	}
	err := node.engine.safetyRules.PrepareBlockRules(prepare)
	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), needChangeView))
	}
}

func testPrepareFetch(t *testing.T, node *TestCBFT, blocks []*types.Block) {
	// need fetch block
	b := NewBlock(blocks[len(blocks)-1].Hash(), blocks[len(blocks)-1].NumberU64()+3)
	vqc := &ctypes.ViewChangeQuorumCert{
		Epoch:       node.engine.state.Epoch(),
		ViewNumber:  node.engine.state.ViewNumber(),
		BlockHash:   blocks[len(blocks)-1].Hash(),
		BlockNumber: blocks[len(blocks)-1].NumberU64(),
	}
	prepare := &protocols.PrepareBlock{
		Epoch:         node.engine.state.Epoch(),
		ViewNumber:    node.engine.state.ViewNumber() + 1,
		Block:         b,
		BlockIndex:    0,
		ProposalIndex: 0,
		ViewChangeQC: &ctypes.ViewChangeQC{
			QCs: []*ctypes.ViewChangeQuorumCert{vqc},
		},
	}
	err := node.engine.safetyRules.PrepareBlockRules(prepare)
	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), needFetchBlock))
	}
}
