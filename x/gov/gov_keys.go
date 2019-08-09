package gov

import (
	"bytes"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	KeyDelimiter               = []byte(":")
	keyPrefixProposal          = []byte("Proposal")
	keyPrefixVote              = []byte("Vote")
	keyPrefixTallyResult       = []byte("TallyResult")
	keyPrefixVotingProposals   = []byte("VotingProposals")
	keyPrefixEndProposals      = []byte("EndProposals")
	keyPrefixPreActiveProposal = []byte("PreActiveProposal")
	keyPrefixPreActiveVersion  = []byte("PreActiveVersion")
	keyPrefixActiveVersion     = []byte("ActiveVersion")
	keyPrefixActiveNodes       = []byte("ActiveNodes")
	keyPrefixAccuVerifiers     = []byte("AccuVerifiers")
	keyPrefixParams            = []byte("Params")
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

func KeyPreActiveProposals() []byte {
	return keyPrefixPreActiveProposal
}

func KeyEndProposals() []byte {
	return keyPrefixEndProposals
}

func KeyActiveVersion() []byte {
	return keyPrefixActiveVersion
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

func KeyParams() []byte {
	return keyPrefixParams
}
