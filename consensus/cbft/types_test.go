package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

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
