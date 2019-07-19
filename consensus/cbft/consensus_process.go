package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

// Perform security rule verification，store in blockTree, Whether to start synchronization
func (cbft *Cbft) OnPrepareBlock(id string, msg *protocols.PrepareBlock) error {
	if err := cbft.safetyRules.PrepareBlockRules(msg); err != nil {
		if err.Fetch() {
			cbft.fetchBlock(id, msg.Block.Hash(), msg.Block.NumberU64())
			//todo fetch block
		}
	}
	cbft.state.AddPrepareBlock(msg)

	cbft.prepareBlockFetchRules(id, msg)
	return nil
}

// Perform security rule verification，store in blockTree, Whether to start synchronization
func (cbft *Cbft) OnPrepareVote(id string, msg *protocols.PrepareVote) error {
	if err := cbft.safetyRules.PrepareVoteRules(msg); err != nil {
		if err.Fetch() {
			//todo fetch block
		}
	}

	cbft.prepareVoteFetchRules(id, msg)
	//todo parse pubkey as id
	cbft.state.AddPrepareVote("", msg)
	//todo new qc block

	return nil
}

// Perform security rule verification, view switching
func (cbft *Cbft) OnViewChange(id string, msg *protocols.ViewChange) error {
	if err := cbft.safetyRules.ViewChangeRules(msg); err != nil {
		if err.Fetch() {
			cbft.fetchBlock(id, msg.BlockHash, msg.BlockNum)
		}
	}

	//todo parse pubkey as id
	cbft.state.AddViewChange("", msg)
	return nil
}

//Perform security rule verification, view switching
func (cbft *Cbft) OnInsertQCBlock(blocks []*types.Block, qcs []*ctypes.QuorumCert) error {
	//todo insert tree, update view
	return nil
}
