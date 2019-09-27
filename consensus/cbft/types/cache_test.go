package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSyncCache(t *testing.T) {
	cache := NewSyncCache(2 * time.Millisecond)
	cache.Add(1)
	assert.False(t, cache.AddOrReplace(1))
	time.Sleep(3 * time.Millisecond)
	assert.True(t, cache.AddOrReplace(1))
	cache.Remove(1)
	assert.True(t, cache.AddOrReplace(1))
	cache.AddOrReplace(2)
	cache.AddOrReplace(3)
	assert.Equal(t, 3, cache.Len())
	cache.Purge()
	assert.Equal(t, 0, cache.Len())

}

func TestCSMsgPool(t *testing.T) {
	pool := NewCSMsgPool()

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareBlock(i, &MsgInfo{})
	}

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareVote(i, i+1, &MsgInfo{})
	}

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareQC(i, &MsgInfo{})
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareBlock(uint32(i)))
		assert.Equal(t, uint32(1), pool.getBlockMetric(i))
	}

	for i := uint32(10); i < 11; i++ {
		assert.Nil(t, pool.GetPrepareBlock(uint32(i)))
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareVote(i, i+1))
		assert.Equal(t, uint32(1), pool.getVoteMetric(i, i+1))
	}

	for i := uint32(10); i < 11; i++ {
		assert.Nil(t, pool.GetPrepareVote(i, i+1))
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareQC(uint32(i)))
		assert.Equal(t, uint32(1), pool.getQCMetric(i))
	}

	for i := uint32(10); i < 11; i++ {
		assert.Nil(t, pool.GetPrepareQC(uint32(i)))
	}

	//re-add
	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareBlock(i, &MsgInfo{})
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareBlock(uint32(i)))
		assert.Equal(t, uint32(2), pool.getBlockMetric(i))
	}
	assert.Equal(t, uint32(0), pool.getBlockMetric(11))

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareQC(i, &MsgInfo{})
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareQC(uint32(i)))
		assert.Equal(t, uint32(2), pool.getQCMetric(i))
	}
	assert.Equal(t, uint32(0), pool.getQCMetric(11))

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareVote(i, i+1, &MsgInfo{})
		pool.AddPrepareVote(i, i+2, &MsgInfo{})

	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareVote(i, i+1))
		assert.NotNil(t, pool.GetPrepareVote(i, i+2))
		assert.Equal(t, uint32(2), pool.getVoteMetric(i, i+1))
	}
	assert.Equal(t, uint32(0), pool.getVoteMetric(10, 11))

	pool.Purge()

	assert.Empty(t, pool.prepareBlocks)
	assert.Empty(t, pool.prepareVotes)
	assert.Empty(t, pool.prepareQC)

	assert.Empty(t, pool.qcMetric)
	assert.Empty(t, pool.blockMetric)
	assert.Empty(t, pool.voteMetric)
}
