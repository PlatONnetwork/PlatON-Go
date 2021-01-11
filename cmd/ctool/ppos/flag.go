package ppos

import "gopkg.in/urfave/cli.v1"

var (
	rpcUrlFlag = cli.StringFlag{
		Name:  "rpcurl",
		Usage: "the rpc url",
	}

	jsonFlag = cli.BoolFlag{
		Name:  "json",
		Usage: "print raw transaction",
	}

	addressHRPFlag = cli.StringFlag{
		Name:  "addressHRP",
		Usage: "set address hrp",
	}
)
