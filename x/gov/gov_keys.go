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

package gov

import (
	"bytes"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	KeyDelimiter               = []byte(":")
	keyPrefixProposal          = []byte("PID")
	keyPrefixVote              = []byte("Vote")
	keyPrefixTallyResult       = []byte("Result")
	keyPrefixVotingProposals   = []byte("Votings")
	keyPrefixEndProposals      = []byte("Ends")
	keyPrefixPreActiveProposal = []byte("PreActPID")
	keyPrefixPreActiveVersion  = []byte("PreActVer")
	keyPrefixActiveVersions    = []byte("ActVers")
	keyPrefixActiveNodes       = []byte("ActNodes")
	keyPrefixAccuVerifiers     = []byte("AccuVoters")
	keyPrefixPIPIDs            = []byte("PIPIDs")
	keyPrefixParamItems        = []byte("ParamItems")
	keyPrefixParamValue        = []byte("ParamValue")
	keyGovernHASHKey           = []byte("GovernHASH")
)

func KeyProposal(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixProposal,
		proposalID.Bytes(),
	}, KeyDelimiter)

}

func KeyVote(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixVote,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

func KeyTallyResult(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixTallyResult,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

func KeyVotingProposals() []byte {
	return keyPrefixVotingProposals
}

func KeyPreActiveProposal() []byte {
	return keyPrefixPreActiveProposal
}

func KeyEndProposals() []byte {
	return keyPrefixEndProposals
}

func KeyActiveVersions() []byte {
	return keyPrefixActiveVersions
}

func KeyPreActiveVersion() []byte {
	return keyPrefixPreActiveVersion
}

func KeyActiveNodes(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixActiveNodes,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

func KeyAccuVerifier(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixAccuVerifiers,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

func KeyPIPIDs() []byte {
	return keyPrefixPIPIDs
}

func KeyParamItems() []byte {
	return keyPrefixParamItems
}
func KeyParamValue(module, name string) []byte {
	return bytes.Join([][]byte{
		keyPrefixParamValue,
		[]byte(module + "/" + name),
	}, KeyDelimiter)
}

func KeyGovernHASHKey() []byte {
	return keyGovernHASHKey
}
