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
		Action: getVerifierList,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, jsonFlag},
	}
	getValidatorListCmd = cli.Command{
		Name:   "getValidatorList",
		Usage:  "1101,query the list of validators in the current consensus round",
		Action: getValidatorList,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, jsonFlag},
	}
	getCandidateListCmd = cli.Command{
		Name:   "getCandidateList",
		Usage:  "1102,Query the list of all real-time candidates",
		Action: getCandidateList,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, jsonFlag},
	}
	getRelatedListByDelAddrCmd = cli.Command{
		Name:   "getRelatedListByDelAddr",
		Usage:  "1103,Query the NodeID and pledge Id of the node entrusted by the current account address,parameter:add",
		Action: getRelatedListByDelAddr,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, addFlag, jsonFlag},
	}
	getDelegateInfoCmd = cli.Command{
		Name:   "getDelegateInfo",
		Usage:  "1104,Query the delegation information of the current single node,parameter:stakingBlock,address,nodeid",
		Action: getDelegateInfo,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, stakingBlockNumFlag, addFlag, nodeIdFlag, jsonFlag},
	}
	getCandidateInfoCmd = cli.Command{
		Name:   "getCandidateInfo",
		Usage:  "1105,Query the pledge information of the current node,parameter:nodeid",
		Action: getCandidateInfo,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, nodeIdFlag, jsonFlag},
	}
	getPackageRewardCmd = cli.Command{
		Name:   "getPackageReward",
		Usage:  "1200,query the block reward of the current settlement epoch",
		Action: getPackageReward,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, jsonFlag},
	}
	getStakingRewardCmd = cli.Command{
		Name:   "getStakingReward",
		Usage:  "1201,query the pledge reward of the current settlement epoch",
		Action: getStakingReward,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, jsonFlag},
	}
	getAvgPackTimeCmd = cli.Command{
		Name:   "getAvgPackTime",
		Usage:  "1202,average time to query packaged blocks",
		Action: getAvgPackTime,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, jsonFlag},
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
	netCheck(c)
	return query(c, 1100)
}

func getValidatorList(c *cli.Context) error {
	netCheck(c)
	return query(c, 1101)
}

func getCandidateList(c *cli.Context) error {
	netCheck(c)
	return query(c, 1102)
}

func getRelatedListByDelAddr(c *cli.Context) error {
	netCheck(c)
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
	netCheck(c)
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
	netCheck(c)
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
	netCheck(c)
	return query(c, 1200)
}

func getStakingReward(c *cli.Context) error {
	netCheck(c)
	return query(c, 1201)
}

func getAvgPackTime(c *cli.Context) error {
	netCheck(c)
	return query(c, 1202)
}
