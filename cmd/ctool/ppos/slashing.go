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

	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"

	"gopkg.in/urfave/cli.v1"
)

var (
	SlashingCmd = cli.Command{
		Name:  "slashing",
		Usage: "use for slashing",
		Subcommands: []cli.Command{
			checkDuplicateSignCmd,
			zeroProduceNodeListCmd,
		},
	}
	checkDuplicateSignCmd = cli.Command{
		Name:   "checkDuplicateSign",
		Usage:  "3001,query whether the node has been reported for too many signatures,parameter:duplicateSignType,nodeid,blockNum",
		Before: netCheck,
		Action: checkDuplicateSign,
		Flags: []cli.Flag{rpcUrlFlag, addressHRPFlag,
			cli.Uint64Flag{
				Name:  "duplicateSignType",
				Usage: "duplicateSign type,1：prepareBlock，2：prepareVote，3：viewChange",
			},
			nodeIdFlag,
			blockNumFlag, jsonFlag,
		},
	}
	zeroProduceNodeListCmd = cli.Command{
		Name:   "zeroProduceNodeList",
		Usage:  "3002,query the list of nodes with zero block",
		Before: netCheck,
		Action: zeroProduceNodeList,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, jsonFlag},
	}
	blockNumFlag = cli.Uint64Flag{
		Name:  "blockNum",
		Usage: "blockNum",
	}
)

func checkDuplicateSign(c *cli.Context) error {
	duplicateSignType := c.Uint64("duplicateSignType")

	nodeIDstring := c.String(nodeIdFlag.Name)
	if nodeIDstring == "" {
		return errors.New("The reported node ID is not set")
	}
	nodeid, err := enode.HexIDv0(nodeIDstring)
	if err != nil {
		return err
	}

	blockNum := c.Uint64(blockNumFlag.Name)

	return query(c, 3001, uint32(duplicateSignType), nodeid, blockNum)
}

func zeroProduceNodeList(c *cli.Context) error {
	return query(c, 3002)
}
