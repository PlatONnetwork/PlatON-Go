package cbft

import "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

//Perform security rule verification，store in blockTree, Whether to start synchronization
func (cbft *Cbft) OnPrepareBlock(msg *protocols.PrepareBlock) error {
	return nil
}

//Perform security rule verification，store in blockTree, Whether to start synchronization
func (cbft *Cbft) OnPrepareVote(msg *protocols.PrepareVote) error {
	return nil
}

//Perform security rule verification, view switching
func (cbft *Cbft) OnViewChange(msg *protocols.ViewChange) error {
	return nil
}
