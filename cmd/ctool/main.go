package main

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/cmd/ctool/core"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"gopkg.in/urfave/cli.v1"
	"os"
	"sort"
)

var (
	app = utils.NewApp("", "the wasm command line interface")
)

func init() {

	// Initialize the CLI app
	app.Commands = []cli.Command{
		core.DeployCmd,
		core.InvokeCmd,
		core.SendTransactionCmd,
		core.SendRawTransactionCmd,
		core.GetTxReceiptCmd,
		core.StabilityCmd,
		core.StabPrepareCmd,
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	app.After = func(ctx *cli.Context) error {
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
