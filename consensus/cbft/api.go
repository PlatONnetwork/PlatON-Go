package cbft

// PublicBlockChainAPI provides an API to access the PlatON blockchain.
// It offers only methods that operate on public data that is freely available to anyone.

type CbftAPI interface {
	Status() string
	Evidences() string
}

type PublicConsensusAPI struct {
	engine CbftAPI
}

// NewPublicBlockChainAPI creates a new PlatON blockchain API.
func NewPublicConsensusAPI(engine CbftAPI) *PublicConsensusAPI {
	return &PublicConsensusAPI{engine: engine}
}

func (s *PublicConsensusAPI) ConsensusStatus() string {
	return s.engine.Status()
}

func (s *PublicConsensusAPI) Evidences() string {
	return s.engine.Evidences()
}
