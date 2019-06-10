package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

var (
	node1, _ = discover.HexID("f6b218929fa97cf77f1482c952698ec50448def9096110a611ffd1f712a73a4790bf9a05147345b3fe2d3316a4e6e5ba1da25fb5702742767f494fb8f6e09586")
)

func TestProcessingVote(t *testing.T) {
	pv := ProcessingVote{}

	pv.Clear()
	assert.Equal(t, 0, len(pv), "")
	assert.Equal(t, "[]", pv.String(), "")

	pv.Add(common.BytesToHash([]byte{1}), node1, &prepareVote{})
	assert.Equal(t, 1, len(pv), "")
	assert.Equal(t, "[hash:000000…000001,[[nodeId:f6b218929fa97cf7, vote:[Timestamp:0 Hash:000000…000000 Number:0 ValidatorAddr:0x0000000000000000000000000000000000000000 ValidatorIndex:0]]]]", pv.String())
	pv.Clear()
	assert.Equal(t, 0, len(pv), "")
}

func TestPendingVote(t *testing.T) {
	pv := PendingVote{}

	pv.Clear()
	assert.Equal(t, 0, len(pv), "")
	assert.Equal(t, "[]", pv.String(), "")

	pv.Add(common.BytesToHash([]byte{1}), &prepareVote{})
	assert.Equal(t, 1, len(pv), "")
	assert.Equal(t, "[[hash:000000…000001, vote:[Timestamp:0 Hash:000000…000000 Number:0 ValidatorAddr:0x0000000000000000000000000000000000000000 ValidatorIndex:0]]]", pv.String())

	pv.Clear()
	assert.Equal(t, 0, len(pv), "")
	assert.Equal(t, "[]", pv.String(), "")

}

func TestViewChangeVotes(t *testing.T) {
	v := ViewChangeVotes{}

	assert.Equal(t, "[]", v.String(), "")

	v = make(map[common.Address]*viewChangeVote)
	v[common.BytesToAddress([]byte{1})] = &viewChangeVote{}
	assert.Equal(t, "[[addr:0x0000000000000000000000000000000000000001, vote:[Timestamp:0 BlockNum:0 BlockHash:000000…000000 ValidatorIndex:0 ValidatorAddr:0x0000000000000000000000000000000000000000]]]",
		v.String())
}

func TestPendingBlock(t *testing.T) {
	pb := PendingBlock{}

	assert.Equal(t, 0, len(pb), "")
	assert.Equal(t, "[]", pb.String(), "")

	header := &types.Header{
		Number:     big.NewInt(int64(11111)),
		ParentHash: common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065749"),
		Extra:      []byte{},
	}
	block := types.NewBlockWithHeader(header)
	pb.Add(common.BytesToHash([]byte{1}), &prepareBlock{Block: block})
	assert.Equal(t, 1, len(pb), "")
	assert.Equal(t, "[[block:[Timestamp:0 Hash:5f688a…820372 Number:11111 ProposalAddr:0x0000000000000000000000000000000000000000, ProposalIndex:0]]]", pb.String())
	pb.Clear()
	assert.Equal(t, 0, len(pb), "")
}

func newTestBlockExt(number uint64, parent common.Hash, threshold int) *BlockExt {
	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: parent,
	}
	block := types.NewBlockWithHeader(header)
	return NewBlockExtByPeer(block, number, threshold)
}

func newTestBlockExtExtra(number uint64, parent common.Hash, extra []byte, threshold int) *BlockExt {
	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: parent,
		Extra:      extra,
	}
	block := types.NewBlockWithHeader(header)
	return NewBlockExtByPeer(block, number, threshold)
}

func newChain(len int, threshold int) ([]*BlockExt, *BlockExtMap) {
	parent := common.Hash{}
	ext := newTestBlockExt(0, parent, threshold)
	m := NewBlockExtMap(ext, threshold)
	extList := make([]*BlockExt, 0)
	extList = append(extList, ext)
	for i := 1; i < len; i++ {
		ext = newTestBlockExt(uint64(i), ext.block.Hash(), threshold)
		m.Add(ext.block.Hash(), ext.number, ext)
		extList = append(extList, ext)
	}
	return extList, m
}
func TestBlockExtMap(t *testing.T) {
	extList, m := newChain(100, 0)

	assert.Equal(t, 99, len(m.GetSubChainUnExecuted()))

	assert.Nil(t, m.FindHighestConfirmedWithHeader())

	unExeBlocks := m.GetSubChainUnExecuted()
	for _, e := range unExeBlocks {
		e.isExecuted = true
	}
	assert.Equal(t, uint64(99), m.FindHighestConfirmedWithHeader().number)
	m.ClearParents(extList[2].block.Hash(), extList[2].block.NumberU64())
	assert.Equal(t, 98, m.Len())
	m.ClearChildren(extList[2].block.Hash(), extList[2].block.NumberU64(), uint64(time.Now().UnixNano()))
	assert.Equal(t, 1, m.Len())
	m.BaseBlock(extList[2].block.Hash(), extList[2].block.NumberU64())
	assert.Equal(t, m.head.block.Hash(), extList[2].block.Hash())

}
func TestBlockExtMap_GetSubChainWithTwoThirdVotes(t *testing.T) {
	extList, m := newChain(100, 0)
	assert.Equal(t, 100, m.Total())

	for _, b := range extList {
		b.isExecuted = true
	}

	ext := extList[50]
	assert.Equal(t, 50, len(m.GetSubChainWithTwoThirdVotes(ext.block.Hash(), ext.block.NumberU64())))

}

func TestBlockExtMap_GetHasVoteWithoutBlock(t *testing.T) {
	extList, m := newChain(100, 0)
	assert.Equal(t, 100, m.Total())

	ext := m.FindHighestConfirmed(extList[0].block.Hash(), extList[0].number)
	assert.Nil(t, ext)

	unExeBlocks := m.GetSubChainUnExecuted()
	for _, e := range unExeBlocks {
		e.isExecuted = true
	}
	ext = m.FindHighestConfirmed(extList[0].block.Hash(), extList[0].number)
	assert.Equal(t, ext.number, extList[99].number)

	ext = m.FindHighestLogical(extList[0].block.Hash(), extList[0].number)
	assert.Equal(t, ext.number, extList[99].number)

	bs := m.findBlockByNumber(10, 20)
	assert.Len(t, bs, 11)

	bx := m.findBlockExtByNumber(10, 20)
	assert.Len(t, bx, 11)

	b := m.findBlockByHash(extList[30].block.Hash())
	assert.Equal(t, b.Hash(), extList[30].block.Hash())

	for _, b := range extList {
		b.block = nil
	}

	assert.Len(t, m.GetHasVoteWithoutBlock(100), 100)
}

func TestBlockExtMap_GetWithoutTwoThirdVotes(t *testing.T) {
	extList, m := newChain(100, 4)
	assert.Equal(t, 100, m.Total())

	for _, b := range extList {
		b.block = nil
	}

	assert.Len(t, m.GetWithoutTwoThirdVotes(100), 100)
}

func TestSameNumberBlock(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("catch panic.")
		}
	}()
	extList, m := newChain(4, 0)
	ext := newTestBlockExtExtra(uint64(3), extList[2].block.Hash(), []byte{0x01, 0x02}, 0)
	assert.Panics(t, func() { m.Add(ext.block.Hash(), ext.number, ext) })

	t.Log(m.BlockString())
}

func TestBlockExt(t *testing.T) {
	v := createTestValidator(createAccount(4))
	hash := common.BytesToHash(Rand32Bytes(32))

	// create view by validator 1
	view := makeViewChange(v.validator(1).privateKey, 1000*1+100, 0, hash, 1, v.validator(1).address, nil)
	ext := makeConfirmedBlock(v, hash, view, 1)[0]
	extSeal := NewBlockExtBySeal(ext.block, 20, 4)
	extSeal.Merge(ext)
	assert.Len(t, extSeal.Votes(), 0)
	extSeal = NewBlockExtBySeal(ext.block, 1, 4)

	extSeal.Merge(ext)

	assert.Len(t, extSeal.Votes(), 3)

	_, err := extSeal.PrepareBlock()

	assert.NotNil(t, err)

	assert.NotEmpty(t, extSeal.String())

}
