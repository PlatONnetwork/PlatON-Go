package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"sync"
)

type ProducerBlocks struct {
	author       discover.NodeID
	baseBlockNum uint64
	blocks       map[uint64]*types.Block
	lock         sync.Mutex
}

func NewProducerBlocks(author discover.NodeID, blockNum uint64) *ProducerBlocks {
	return &ProducerBlocks{
		author:       author,
		baseBlockNum: blockNum,
		blocks:       make(map[uint64]*types.Block),
	}
}
func (pb *ProducerBlocks) String() string {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	return fmt.Sprintf("author:%s base block:%d, total:%d", pb.author.String(), pb.baseBlockNum, len(pb.blocks))
}
func (pb *ProducerBlocks) SetAuthor(author discover.NodeID) {
	pb.lock.Lock()
	pb.author = author
	pb.blocks = make(map[uint64]*types.Block)
	pb.lock.Unlock()
}

func (pb *ProducerBlocks) AddBlock(block *types.Block) {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	if pb.baseBlockNum < block.NumberU64() {
		pb.blocks[block.NumberU64()] = block
	}
}

func (pb *ProducerBlocks) ExistBlock(block *types.Block) bool {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	if v, ok := pb.blocks[block.NumberU64()]; ok && v.Hash() == block.Hash() {
		return true
	}
	return false
}
func (pb *ProducerBlocks) Author() discover.NodeID {
	pb.lock.Lock()
	pb.lock.Unlock()
	return pb.author
}

func (pb *ProducerBlocks) MaxSequenceBlockNum() uint64 {
	pb.lock.Lock()
	pb.lock.Unlock()
	num := pb.baseBlockNum + 1
	for {
		_, ok := pb.blocks[num]
		if !ok {
			break
		}
		num++
	}
	return num - 1
}

func (pb *ProducerBlocks) MaxSequenceBlock() *types.Block {
	pb.lock.Lock()
	pb.lock.Unlock()
	var block *types.Block
	num := pb.baseBlockNum + 1
	for {
		b, ok := pb.blocks[num]
		if !ok {
			break
		}
		block = b
		num++
	}
	return block
}

func (pb *ProducerBlocks) Len() int {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	return len(pb.blocks)
}
