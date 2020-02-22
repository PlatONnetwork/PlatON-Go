package plugin

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	FORKHASH    = "0xe01a99b5d742bca71acb656dd275a5d352afb58eb5f8839cf4563402c056837d"
	FORKNUM     = 521528
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
