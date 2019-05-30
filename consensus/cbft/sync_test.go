package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducerBlocks(t *testing.T) {
	v := createTestValidator(createAccount(1)).owner
	v2 := createTestValidator(createAccount(1)).owner

	pb := NewProducerBlocks(v.nodeID, 10)

	pb.SetAuthor(v2.nodeID)

	assert.Equal(t, pb.Author(), v2.nodeID)
	blocks := createEmptyBlocks(v2.privateKey, common.BytesToHash(Rand32Bytes(32)), 10, 10)

	for _, b := range blocks {
		pb.AddBlock(b)
	}

	assert.Equal(t, 10, pb.Len())
	assert.True(t, pb.ExistBlock(blocks[0]))
	assert.Equal(t, blocks[9].Hash(), pb.MaxSequenceBlock().Hash())
	assert.Equal(t, blocks[9].NumberU64(), pb.MaxSequenceBlockNum())
}
