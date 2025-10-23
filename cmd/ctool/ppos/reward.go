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
	"github.com/urfave/cli/v2"

	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
)

var (
	RewardCmd = &cli.Command{
		Name:  "reward",
		Usage: "use for reward",
		Subcommands: []*cli.Command{
			getDelegateRewardCmd,
		},
	}
	getDelegateRewardCmd = &cli.Command{
		Name:   "getDelegateReward",
		Usage:  "5100,query account not withdrawn commission rewards at each node,parameter:nodeList(can empty)",
		Before: netCheck,
		Action: getDelegateReward,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, nodeList, jsonFlag},
	}
	nodeList = &cli.StringSliceFlag{
		Name:  "nodeList",
		Usage: "node list,may empty",
	}
)

func getDelegateReward(c *cli.Context) error {
	nodeIDlist := c.StringSlice(nodeList.Name)
	idlist := make([]enode.IDv0, 0)
	for _, node := range nodeIDlist {
		nodeid, err := enode.HexIDv0(node)
		if err != nil {
			return err
		}
		idlist = append(idlist, nodeid)
	}
	return query(c, 5100, idlist)
}
