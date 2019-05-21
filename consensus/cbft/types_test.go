package cbft

import (
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
	ext := extList[50]
	assert.Equal(t, len(m.GetSubChainUnExecuted()), 99)
	assert.Equal(t, len(m.GetSubChainWithTwoThirdVotes(ext.block.Hash(), ext.block.NumberU64())), 50)
	assert.Equal(t, m.FindHighestConfirmedWithHeader().number, uint64(99))
	m.ClearParents(extList[2].block.Hash(), extList[2].block.NumberU64())
	assert.Equal(t, m.Len(), 98)
	m.ClearChildren(extList[2].block.Hash(), extList[2].block.NumberU64(), uint64(time.Now().UnixNano()))
	assert.Equal(t, m.Len(), 1)
	m.BaseBlock(extList[2].block.Hash(), extList[2].block.NumberU64())
	assert.Equal(t, extList[2].block.Hash(), m.head.block.Hash())
}

func TestSameNumberBlock(t *testing.T) {
	extList, m := newChain(4, 0)
	ext := newTestBlockExtExtra(uint64(3), extList[2].block.Hash(), []byte{0x01, 0x02}, 0)
	m.Add(ext.block.Hash(), ext.number, ext)
	t.Log(m.BlockString())
}
