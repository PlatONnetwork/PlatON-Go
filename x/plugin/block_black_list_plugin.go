package plugin

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	FORKHASH    = "0x7560a7552584fd6f5e992ad34721694a5253aa218138e3a0c0a93619b3a30805"
	FORKNUM     = 1000
	FORKVERSION = uint32(0<<16 | 9<<8 | 0)
)

var BlockBlackListERROR = fmt.Errorf("the block is exist in BlackList,hash:%v", FORKHASH)

type BlockBlackListPlugin struct {
	list []common.Hash
}

func NewBlockBlackListPlugin() *BlockBlackListPlugin {
	blackhash := common.HexToHash(FORKHASH)
	bl := new(BlockBlackListPlugin)
	bl.list = make([]common.Hash, 0)
	bl.list = append(bl.list, blackhash)
	return bl
}

func (b *BlockBlackListPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	for _, value := range b.list {
		if blockHash == value {
			return BlockBlackListERROR
		}
	}
	return nil
}

func (b *BlockBlackListPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	return nil
}

func (b *BlockBlackListPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}
