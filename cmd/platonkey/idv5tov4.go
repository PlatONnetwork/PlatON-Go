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
update NodeID to enode.ID.
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
					utils.Fatalf("the NodeID can't be nil")
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
				fmt.Println("NodeID(V5)  : ", idpair.v5.String())
				fmt.Println("enode.ID(V4): ", idpair.v4.String())
				if i != len(outpairs)-1 {
					fmt.Println("---")
				}
			}
		}
		return nil
	},
}
