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

package consensus

import (
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

type EvidenceType uint8

type Evidence interface {
	//Verify(ecdsa.PublicKey) error
	Equal(Evidence) bool
	//return lowest number
	BlockNumber() uint64
	Epoch() uint64
	ViewNumber() uint64
	Hash() []byte
	//Address() common.NodeAddress
	NodeID() discover.NodeID
	BlsPubKey() *bls.PublicKey
	Validate() error
	Type() EvidenceType
	ValidateMsg() bool
}

type Evidences []Evidence

func (e Evidences) Len() int {
	return len(e)
}

type EvidencePool interface {
	//Deserialization of evidence
	//UnmarshalEvidence(data string) (Evidences, error)
	//Get current evidences
	Evidences() Evidences
	//Clear all evidences
	Clear(epoch uint64, viewNumber uint64)
	Close()
}
