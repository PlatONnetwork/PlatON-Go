package ppos

import (
	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	RewardCmd = cli.Command{
		Name:  "reward",
		Usage: "use for reward",
		Subcommands: []cli.Command{
			getDelegateRewardCmd,
		},
	}
	getDelegateRewardCmd = cli.Command{
		Name:   "getDelegateReward",
		Usage:  "5100,query account not withdrawn commission rewards at each node,parameter:nodeList(can empty)",
		Action: getDelegateReward,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, nodeList, jsonFlag},
	}
	nodeList = cli.StringSliceFlag{
		Name:  "nodeList",
		Usage: "node list,may empty",
	}
)

func getDelegateReward(c *cli.Context) error {
	netCheck(c)
	nodeIDlist := c.StringSlice(nodeList.Name)
	idlist := make([]discover.NodeID, 0)
	for _, node := range nodeIDlist {
		nodeid, err := discover.HexID(node)
		if err != nil {
			return err
		}
		idlist = append(idlist, nodeid)
	}
	return query(c, 5100, idlist)
}
