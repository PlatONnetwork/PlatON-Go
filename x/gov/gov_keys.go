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
