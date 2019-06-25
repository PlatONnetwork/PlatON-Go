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
	keyPrefixVotedVerifiers    = []byte("VotedVerifiers")
	keyPrefixActiveNodes       = []byte("ActiveNodes")
	keyPrefixAccuVerifiers     = []byte("AccuVerifiers")
)

// 提案的key
func KeyProposal(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixProposal,
		proposalID.Bytes(),
	}, KeyDelimiter)

}

// 投票的key
func KeyVote(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixVote,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

// 投票结果的key
func KeyTallyResult(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixTallyResult,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

// 正在投票的提案列表的key
func KeyVotingProposals() []byte {
	return keyPrefixVotingProposals
}

// 预生效提案ID的key
func KeyPreActiveProposalID() []byte {
	return keyPrefixEndProposals
}


// 所有操作均结束的提案列表的key
func KeyEndProposals() []byte {
	return keyPrefixEndProposals
}

// 生效版本的key
func KeyActiveVersion() []byte {
	return keyPrefixActiveVersion
}

// 预生效版本的key
func KeyPreActiveVersion() []byte {
	return keyPrefixPreActiveVersion
}

// 已投票的验证人列表key
func KeyVotedVerifier(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixVotedVerifiers,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

// 已升级节点列表的key
func KeyActiveNodes(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixActiveNodes,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

// 提案投票期内累积的不同验证人的key
func KeyAccuVerifiers(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixAccuVerifiers,
		proposalID.Bytes(),
	}, KeyDelimiter)
}
