package plugin

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

func NewCollectDeclareVersionPlugin() *CollectDeclareVersionPlugin {
	cd := new(CollectDeclareVersionPlugin)
	cd.num = params.FORKNUM
	cd.version = params.FORKVERSION
	return cd
}

type CollectDeclareVersionPlugin struct {
	num     uint64
	version uint32
}

func (b *CollectDeclareVersionPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	if header.ParentHash.String() == params.FORKHASH && (header.Number.Uint64()-params.FORKNUM) == 1 {
		if err := gov.AddActiveVersion(b.version, header.Number.Uint64(), state); err != nil {
			return err
		}
		log.Debug("CollectDeclareVersionPlugin begin ClearProcessingProposals")
		if err := gov.ClearProcessingProposals(blockHash, state); err != nil {
			return err
		}
		gov.EnableCounter = true
		gov.NodeDeclaredVersionsCounter = make(map[discover.NodeID]uint32)
	}
	return nil
}

func (b *CollectDeclareVersionPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	if header.ParentHash.String() == params.FORKHASH && (header.Number.Uint64()-params.FORKNUM) == 1 {
		defer func() {
			gov.EnableCounter = false
		}()
		list, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())
		if err != nil {
			return err
		}
		size := 0
		for _, stakingNodeID := range list {
			if value, ok := gov.NodeDeclaredVersionsCounter[stakingNodeID]; ok {
				if value == b.version {
					size++
				}
			}
		}
		wantSize := (len(list)*2)/3 + 1
		log.Debug("CollectDeclareVersionPlugin begin ClearProcessingProposals,count size", "size", size, "want", wantSize)
		if size < wantSize {
			return fmt.Errorf("the block Collect Declare Version less than %v", wantSize)
		}
	}
	return nil
}

func (b *CollectDeclareVersionPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

func IsForkBlock(blockNumber uint64, parentHash string) bool {
	if blockNumber == params.FORKNUM+1 && parentHash == params.FORKHASH {
		return true
	}
	return false
}
