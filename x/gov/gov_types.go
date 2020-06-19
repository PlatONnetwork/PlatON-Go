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
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

type TallyResult struct {
	ProposalID    common.Hash    `json:"proposalID"`
	Yeas          uint64         `json:"yeas"`
	Nays          uint64         `json:"nays"`
	Abstentions   uint64         `json:"abstentions"`
	AccuVerifiers uint64         `json:"accuVerifiers"`
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
	ParamVerifier ParamVerifier `json:"-"`
}
