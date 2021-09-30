// Copyright 2021 The PlatON Network Authors
// This file is part of PlatON-Go.
//
// PlatON-Go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// PlatON-Go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PlatON-Go. If not, see <http://www.gnu.org/licenses/>.

package ppos

import (
	"errors"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	StakingCmd = cli.Command{
		Name:  "staking",
		Usage: "use for staking",
		Subcommands: []cli.Command{
			GetVerifierListCmd,
			getValidatorListCmd,
			getCandidateListCmd,
			getRelatedListByDelAddrCmd,
			getDelegateInfoCmd,
			getCandidateInfoCmd,
			getPackageRewardCmd,
			getStakingRewardCmd,
			getAvgPackTimeCmd,
		},
	}
	GetVerifierListCmd = cli.Command{
		Name:   "getVerifierList",
		Usage:  "1100,query the validator queue of the current settlement epoch",
		Before: netCheck,
		Action: getVerifierList,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, jsonFlag},
	}
	getValidatorListCmd = cli.Command{
		Name:   "getValidatorList",
		Usage:  "1101,query the list of validators in the current consensus round",
		Before: netCheck,
		Action: getValidatorList,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, jsonFlag},
	}
	getCandidateListCmd = cli.Command{
		Name:   "getCandidateList",
		Usage:  "1102,Query the list of all real-time candidates",
		Before: netCheck,
		Action: getCandidateList,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, jsonFlag},
	}
	getRelatedListByDelAddrCmd = cli.Command{
		Name:   "getRelatedListByDelAddr",
		Usage:  "1103,Query the NodeID and staking Id of the node entrusted by the current account address,parameter:add",
		Before: netCheck,
		Action: getRelatedListByDelAddr,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, addFlag, jsonFlag},
	}
	getDelegateInfoCmd = cli.Command{
		Name:   "getDelegateInfo",
		Usage:  "1104,Query the delegation information of the current single node,parameter:stakingBlock,address,nodeid",
		Before: netCheck,
		Action: getDelegateInfo,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, stakingBlockNumFlag, addFlag, nodeIdFlag, jsonFlag},
	}
	getCandidateInfoCmd = cli.Command{
		Name:   "getCandidateInfo",
		Usage:  "1105,Query the staking information of the current node,parameter:nodeid",
		Before: netCheck,
		Action: getCandidateInfo,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, nodeIdFlag, jsonFlag},
	}
	getPackageRewardCmd = cli.Command{
		Name:   "getPackageReward",
		Usage:  "1200,query the block reward of the current settlement epoch",
		Before: netCheck,
		Action: getPackageReward,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, jsonFlag},
	}
	getStakingRewardCmd = cli.Command{
		Name:   "getStakingReward",
		Usage:  "1201,query the staking reward of the current settlement epoch",
		Before: netCheck,
		Action: getStakingReward,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, jsonFlag},
	}
	getAvgPackTimeCmd = cli.Command{
		Name:   "getAvgPackTime",
		Usage:  "1202,average time to query packaged blocks",
		Before: netCheck,
		Action: getAvgPackTime,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, jsonFlag},
	}
	addFlag = cli.StringFlag{
		Name:  "address",
		Usage: "account address",
	}
	stakingBlockNumFlag = cli.Uint64Flag{
		Name:  "stakingBlock",
		Usage: "block height when staking is initiated",
	}
	nodeIdFlag = cli.StringFlag{
		Name:  "nodeid",
		Usage: "node id",
	}
)

func getVerifierList(c *cli.Context) error {
	return query(c, 1100)
}

func getValidatorList(c *cli.Context) error {
	return query(c, 1101)
}

func getCandidateList(c *cli.Context) error {
	return query(c, 1102)
}

func getRelatedListByDelAddr(c *cli.Context) error {
	addstring := c.String(addFlag.Name)
	if addstring == "" {
		return errors.New("The Del's account address is not set")
	}
	add, err := common.Bech32ToAddress(addstring)
	if err != nil {
		return err
	}
	return query(c, 1103, add)
}

func getDelegateInfo(c *cli.Context) error {
	addstring := c.String(addFlag.Name)
	if addstring == "" {
		return errors.New("The Del's account address is not set")
	}
	add, err := common.Bech32ToAddress(addstring)
	if err != nil {
		return err
	}
	nodeIDstring := c.String(nodeIdFlag.Name)
	if nodeIDstring == "" {
		return errors.New("The verifier's node ID is not set")
	}
	nodeid, err := discover.HexID(nodeIDstring)
	if err != nil {
		return err
	}
	stakingBlockNum := c.Uint64(stakingBlockNumFlag.Name)
	return query(c, 1104, stakingBlockNum, add, nodeid)
}

func getCandidateInfo(c *cli.Context) error {
	nodeIDstring := c.String(nodeIdFlag.Name)
	if nodeIDstring == "" {
		return errors.New("The verifier's node ID is not set")
	}
	nodeid, err := discover.HexID(nodeIDstring)
	if err != nil {
		return err
	}
	return query(c, 1105, nodeid)
}

func getPackageReward(c *cli.Context) error {
	return query(c, 1200)
}

func getStakingReward(c *cli.Context) error {
	return query(c, 1201)
}

func getAvgPackTime(c *cli.Context) error {
	return query(c, 1202)
}
