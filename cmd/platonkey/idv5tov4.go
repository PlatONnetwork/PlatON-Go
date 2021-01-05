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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discv5"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
)

var V5IDFileFlag = cli.StringFlag{
	Name:  "v5IDFile",
	Usage: "NodeID file like  [hex,hex...]",
}

type NodeIDPair struct {
	v4 enode.ID
	v5 discv5.NodeID
}
var commandV5ToV4 = cli.Command{
	Name:      "v5tov4",
	Usage:     "update NodeID to enode.ID",
	ArgsUsage: "[<node>]",
	Description: `
update Id to enode.ID.
`,
	Flags: []cli.Flag{
		jsonFlag,
	},
	Action: func(ctx *cli.Context) error {
		var nodes []string

		if ctx.IsSet(V5IDFileFlag.Name) {
			filePath := ctx.String(V5IDFileFlag.Name)
			nodeJson, err := ioutil.ReadFile(filePath)
			if err != nil {
				utils.Fatalf("Failed to read the nodefile at '%s': %v", filePath, err)
			}
			if err := json.Unmarshal(nodeJson, &nodes); err != nil {
				utils.Fatalf("Failed to json decode '%s': %v", filePath, err)
			}
		} else {
			for _, v5idStr := range ctx.Args() {
				if v5idStr == "" {
					utils.Fatalf("the Id can't be nil")
				}
				nodes = append(nodes, v5idStr)
			}
		}

		var outpairs []NodeIDPair
		for _, v5idstr := range nodes {
			v5id := discv5.MustHexID(v5idstr)
			v4id := enode.NodeIDToIDV4(v5id)
			out := NodeIDPair{
				v4: v4id,
				v5: v5id,
			}
			outpairs = append(outpairs, out)
		}

		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(outpairs)
		} else {
			for i, idpair := range outpairs {
				fmt.Println("Id(V5)  : ", idpair.v5.String())
				fmt.Println("enode.ID(V4): ", idpair.v4.String())
				if i != len(outpairs)-1 {
					fmt.Println("---")
				}
			}
		}
		return nil
	},
}
