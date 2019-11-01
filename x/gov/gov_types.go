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

type VoteInfo struct {
	ProposalID common.Hash     `json:"proposalID"`
	VoteNodeID discover.NodeID `json:"voteNodeID"`
	VoteOption VoteOption      `json:"voteOption"`
}

type VoteValue struct {
	VoteNodeID discover.NodeID `json:"voteNodeID"`
	VoteOption VoteOption      `json:"voteOption"`
}

type ActiveVersionValue struct {
	ActiveVersion uint32 `json:"ActiveVersion"`
	ActiveBlock   uint64 `json:"ActiveBlock"`
}
type ParamItem struct {
	Module string `json:"Module"`
	Name   string `json:"Name"`
	Desc   string `json:"Desc"`
}

type ParamValue struct {
	StaleValue  string `json:"StaleValue"`
	Value       string `json:"Value"`
	ActiveBlock uint64 `json:"ActiveBlock"`
}

type GovernParam struct {
	ParamItem     *ParamItem
	ParamValue    *ParamValue
	ParamVerifier ParamVerifier
}
