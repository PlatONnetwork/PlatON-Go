package cbft

import "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

//Perform security rule verification，store in blockTree, Whether to start synchronization
func (cbft *Cbft) OnPrepareBlock(msg *protocols.PrepareBlock) error {
	if err := cbft.safetyRules.PrepareBlockRules(msg); err != nil {
		if err.Fetch() {
			//todo fetch block
		}
	}
	cbft.state.AddPrepareBlock(msg)

	return nil
}

//Perform security rule verification，store in blockTree, Whether to start synchronization
func (cbft *Cbft) OnPrepareVote(msg *protocols.PrepareVote) error {
	if err := cbft.safetyRules.PrepareVoteRules(msg); err != nil {
		if err.Fetch() {
			//todo fetch block
		}
	}

	//todo parse pubkey as id
	cbft.state.AddPrepareVote("", msg)
	//todo new qc block
	return nil
}

//Perform security rule verification, view switching
func (cbft *Cbft) OnViewChange(msg *protocols.ViewChange) error {
	if err := cbft.safetyRules.ViewChangeRules(msg); err != nil {
		if err.Fetch() {
			//todo fetch block
		}
	}

	//todo parse pubkey as id
	cbft.state.AddViewChange("", msg)
	return nil
}
