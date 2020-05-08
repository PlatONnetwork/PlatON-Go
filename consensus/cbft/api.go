package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
)

type Status struct {
	Tree      *types.BlockTree `json:"blockTree"`
	State     *state.ViewState `json:"state"`
	Validator bool             `json:"validator"`
}

// API defines an exposed API function interface.
type API interface {
	Status() *Status
	Evidences() string
	GetPrepareQC(number uint64) *types.QuorumCert
	GetSchnorrNIZKProve() (*bls.SchnorrProof, error)
}

// PublicConsensusAPI provides an API to access the PlatON blockchain.
// It offers only methods that operate on public data that
// is freely available to anyone.
type PublicConsensusAPI struct {
	engine API
}

// NewPublicConsensusAPI creates a new PlatON blockchain API.
func NewPublicConsensusAPI(engine API) *PublicConsensusAPI {
	return &PublicConsensusAPI{engine: engine}
}

// ConsensusStatus returns the status data of the consensus engine.
func (s *PublicConsensusAPI) ConsensusStatus() *Status {
	return s.engine.Status()
}

// Evidences returns the relevant data of the verification.
func (s *PublicConsensusAPI) Evidences() string {
	return s.engine.Evidences()
}

// GetPrepareQC returns the QC certificate corresponding to the blockNumber.
func (s *PublicConsensusAPI) GetPrepareQC(number uint64) *types.QuorumCert {
	return s.engine.GetPrepareQC(number)
}

func (s *PublicConsensusAPI) GetSchnorrNIZKProve() string {
	proof, err := s.engine.GetSchnorrNIZKProve()
	if nil != err {
		return err.Error()
	}
	proofByte, err := proof.MarshalText()
	if nil != err {
		return err.Error()
	}
	return string(proofByte)
}
