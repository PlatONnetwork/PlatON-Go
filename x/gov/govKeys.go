package gov

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	KeyDelimiter               = []byte(":")
	keyPrefixProposal          = []byte("Proposal")
	keyPrefixVote              = []byte("Vote")
	keyPrefixTallyResult       = []byte("TallyResult")
	keyPrefixVotingProposals   = []byte("VotingProposals")
	keyPrefixEndProposals      = []byte("EndProposals")
	//keyPrefixPreActiveProposal = []byte("PreActiveProposal")
	keyPrefixPreActiveVersion  = []byte("PreActiveVersion")
	keyPrefixActiveVersion     = []byte("ActiveVersion")
	keyPrefixVotedVerifiers    = []byte("VotedVerifiers")
	keyPrefixActiveNodes       = []byte("ActiveNodes")
	keyPrefixAccuVerifiers     = []byte("AccuVerifiers")
)

// 产生保存提案的key
func KeyProposal(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixProposal,
		proposalID.Bytes(),
	}, KeyDelimiter)

}

// 生成存储投票的key
func KeyVote(proposalID common.Hash, voter *discover.NodeID) []byte {
	return bytes.Join([][]byte{
		keyPrefixVote,
		proposalID.Bytes(),
		voter.Bytes(),
	}, KeyDelimiter)
}

// 生成投票结果的key
func KeyTallyResult(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixTallyResult,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

// 生成正在投票的提案列表的key
func KeyVotingProposals() []byte {
	return keyPrefixVotingProposals
}

// 生成投票结束的提案列表的key
func KeyEndProposalID() []byte {
	return keyPrefixEndProposals
}

// 生成已激活版本列表的key
func KeyActiveVersion() []byte {
	return keyPrefixActiveVersion
}

// 生成已激活版本列表的key
func KeyPreActiveVersion() []byte {
	return keyPrefixPreActiveVersion
}

// 生成已投票的验证人列表key
func KeyVotedVerifier(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixVotedVerifiers,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

// 生成升级节点列表的key
func KeyActiveNodes(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixActiveNodes,
		proposalID.Bytes(),
	}, KeyDelimiter)
}

// 生成提案投票期内累计不同验证人的key
func KeyAccuVerifiers(proposalID common.Hash) []byte {
	return bytes.Join([][]byte{
		keyPrefixAccuVerifiers,
		proposalID.Bytes(),
	}, KeyDelimiter)
}
