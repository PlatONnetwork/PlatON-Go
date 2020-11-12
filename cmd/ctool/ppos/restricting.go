package ppos

import (
	"errors"

	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	RestrictingCmd = cli.Command{
		Name:  "restricting",
		Usage: "use for restricting",
		Subcommands: []cli.Command{
			getRestrictingInfoCmd,
		},
	}
	getRestrictingInfoCmd = cli.Command{
		Name:   "getRestrictingInfo",
		Usage:  "4100,get restricting info,parameter:address",
		Action: getRestrictingInfo,
		Flags:  []cli.Flag{rpcUrlFlag, testNetFlag, addFlag, jsonFlag},
	}
)

func getRestrictingInfo(c *cli.Context) error {
	netCheck(c)
	addstring := c.String(addFlag.Name)
	if addstring == "" {
		return errors.New("The locked position release to the account account is not set")
	}
	add, err := common.Bech32ToAddress(addstring)
	if err != nil {
		return err
	}
	return query(c, 4100, add)
}
