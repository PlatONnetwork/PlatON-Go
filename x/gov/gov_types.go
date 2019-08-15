package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

type TallyResult struct {
	ProposalID    common.Hash    `json:"proposalID"`
	Yeas          uint16         `json:"yeas"`
	Nays          uint16         `json:"nays"`
	Abstentions   uint16         `json:"abstentions"`
	AccuVerifiers uint16         `json:"accuVerifiers"`
	Status        ProposalStatus `json:"status"`
	CanceledBy    common.Hash    `json:"canceledBy"`
}

type Vote struct {
	ProposalID common.Hash     `json:"proposalID"`
	VoteNodeID discover.NodeID `json:"voteNodeID"`
	VoteOption VoteOption      `json:"voteOption"`
}

type VoteValue struct {
	VoteNodeID discover.NodeID `json:"voteNodeID"`
	VoteOption VoteOption      `json:"voteOption"`
}

type ParamValue struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type ProgramVersionValue struct {
	ProgramVersion     uint32             `json:"ProgramVersion"`
	ProgramVersionSign common.VersionSign `json:"ProgramVersionSign"`
}

type ActiveVersionValue struct {
	ActiveVersion uint32 `json:"ActiveVersion"`
	ActiveBlock   uint64 `json:"ActiveBlock"`
}
