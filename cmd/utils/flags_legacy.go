// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/eth/ethconfig"
	"strings"

	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/node"
)

var ShowDeprecated = cli.Command{
	Action:      showDeprecated,
	Name:        "show-deprecated-flags",
	Usage:       "Show flags that have been deprecated",
	ArgsUsage:   " ",
	Category:    "MISCELLANEOUS COMMANDS",
	Description: "Show flags that have been deprecated and will soon be removed",
}

var DeprecatedFlags = []cli.Flag{}

var (
	LegacyRPCEnabledFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Enable the HTTP-RPC server (deprecated, use --http)",
	}
	LegacyRPCListenAddrFlag = cli.StringFlag{
		Name:  "rpcaddr",
		Usage: "HTTP-RPC server listening interface (deprecated, use --http.addr)",
		Value: node.DefaultHTTPHost,
	}
	LegacyRPCPortFlag = cli.IntFlag{
		Name:  "rpcport",
		Usage: "HTTP-RPC server listening port (deprecated, use --http.port)",
		Value: node.DefaultHTTPPort,
	}
	LegacyRPCCORSDomainFlag = cli.StringFlag{
		Name:  "rpccorsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced) (deprecated, use --http.corsdomain)",
		Value: "",
	}
	LegacyRPCVirtualHostsFlag = cli.StringFlag{
		Name:  "rpcvhosts",
		Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (deprecated, use --http.vhosts)",
		Value: strings.Join(node.DefaultConfig.HTTPVirtualHosts, ","),
	}
	LegacyRPCApiFlag = cli.StringFlag{
		Name:  "rpcapi",
		Usage: "API's offered over the HTTP-RPC interface (deprecated, use --http.api)",
		Value: "",
	}
	LegacyWSListenAddrFlag = cli.StringFlag{
		Name:  "wsaddr",
		Usage: "WS-RPC server listening interface (deprecated, use --ws.addr)",
		Value: node.DefaultWSHost,
	}
	LegacyWSPortFlag = cli.IntFlag{
		Name:  "wsport",
		Usage: "WS-RPC server listening port (deprecated, use --ws.port)",
		Value: node.DefaultWSPort,
	}
	LegacyWSApiFlag = cli.StringFlag{
		Name:  "wsapi",
		Usage: "API's offered over the WS-RPC interface (deprecated, use --ws.api)",
		Value: "",
	}
	LegacyWSAllowedOriginsFlag = cli.StringFlag{
		Name:  "wsorigins",
		Usage: "Origins from which to accept websockets requests (deprecated, use --ws.origins)",
		Value: "",
	}
	LegacyGpoBlocksFlag = cli.IntFlag{
		Name:  "gpoblocks",
		Usage: "Number of recent blocks to check for gas prices (deprecated, use --gpo.blocks)",
		Value: ethconfig.Defaults.GPO.Blocks,
	}
	LegacyGpoPercentileFlag = cli.IntFlag{
		Name:  "gpopercentile",
		Usage: "Suggested gas price is the given percentile of a set of recent transaction gas prices (deprecated, use --gpo.percentile)",
		Value: ethconfig.Defaults.GPO.Percentile,
	}
)

// showDeprecated displays deprecated flags that will be soon removed from the codebase.
func showDeprecated(*cli.Context) {
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("The following flags are deprecated and will be removed in the future!")
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println()

	for _, flag := range DeprecatedFlags {
		fmt.Println(flag.String())
	}
}
