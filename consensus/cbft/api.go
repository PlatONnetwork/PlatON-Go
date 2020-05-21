// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

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
