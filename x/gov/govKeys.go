package gov

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	KeyDelimiter             = []byte(":")
	keyPrefixProposal        = []byte("Proposal")
	keyPrefixVote            = []byte("Vote")
	keyPrefixTallyResult     = []byte("TallyResult")
	keyPrefixVotingProposals = []byte("VotingProposals")
	keyPrefixEndProposals    = []byte("EndProposals")
	//keyPrefixPreActiveProposal = []byte("PreActiveProposal")
	keyPrefixPreActiveVersion = []byte("PreActiveVersion")
	keyPrefixActiveVersion    = []byte("ActiveVersion")
	keyPrefixDeclaredNodes    = []byte("DeclaredNodes")
	keyPrefixTotalVerifiers   = []byte("TotalVerifiers")
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
func KeyPreActiveProposals() []byte {
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

// 已升级节点列表的key
func KeyDeclaredNodes(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixDeclaredNodes,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

// 提案投票期内累积的不同验证人的key
func KeyTotalVerifiers(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixTotalVerifiers,
		proposalID.Bytes(),
	}, KeyDelimiter)
}
