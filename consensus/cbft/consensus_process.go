package cbft

//Perform security rule verification，store in blockTree, Whether to start synchronization
func (cbft *Cbft) OnPrepareBlock(msg *prepareBlock) error {
	return nil
}

//Perform security rule verification，store in blockTree, Whether to start synchronization
func (cbft *Cbft) OnPrepareVote(msg *prepareVote) error {
	return nil
}

//Perform security rule verification, view switching
func (cbft *Cbft) OnViewChange(msg *viewChange) error {
	return nil
}
