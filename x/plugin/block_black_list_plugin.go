package plugin

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	FORKHASH    = "0xec113de20b1becec357f7e74d8662ad904197581300445b1f024ce0fc28f55ea"
	FORKNUM     = 114
	FORKVERSION = uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)
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
