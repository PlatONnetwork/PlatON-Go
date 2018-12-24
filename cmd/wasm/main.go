package main

import (
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

var (
	app = utils.NewApp("", "the wasm command line interface")

	DebugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "output full trace logs",
	}
	MemProfileFlag = cli.StringFlag{
		Name:  "memprofile",
		Usage: "creates a memory profile at the given path",
	}
	TxTypeFlag = cli.Int64Flag{
		Name: "txtype",
		Usage: "The transaction type",
	}
	StatDumpFlag = cli.BoolFlag{
		Name:  "statdump",
		Usage: "displays stack and heap memory information",
	}
	CodeFlag = cli.StringFlag{
		Name:  "code",
		Usage: "WASM code",
	}
	CodeFileFlag = cli.StringFlag{
		Name:  "codefile",
		Usage: "File containing WASM code. If '-' is specified, code is read from stdin",
	}
	AbiFlag = cli.StringFlag{
		Name:  "abi",
		Usage: "WASM abi",
	}
	AbiFileFlag = cli.StringFlag{
		Name:  "abifile",
		Usage: "File containing WASM abi. If '-' is specified, abi is read from stdin",
	}
	GasFlag = cli.Uint64Flag{
		Name:  "gas",
		Usage: "gas limit for the wasm",
	}
	GasPriceFlag = utils.BigFlag{
		Name:  "gasPrice",
		Usage: "price set for the wasm",
	}
	ValueFlag = utils.BigFlag{
		Name:  "value",
		Usage: "value set for the wasm",
	}
	InputFlag = cli.StringFlag{
		Name:  "input",
		Usage: "input for the wasm",
	}
	VerbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "sets the verbosity level",
	}
	CreateFlag = cli.BoolFlag{
		Name:  "create",
		Usage: "indicates the action should be create rather than call",
	}
	GenesisFlag = cli.StringFlag{
		Name:  "prestate",
		Usage: "JSON file with prestate (genesis) config",
	}
	MachineFlag = cli.BoolFlag{
		Name:  "json",
		Usage: "output trace logs in machin readable format(json)",
	}
	SenderFlag = cli.StringFlag{
		Name:  "sender",
		Usage: "The transaction origin",
	}
	ReceiverFlag = cli.StringFlag{
		Name:  "receiver",
		Usage: "The transaction receiver (execution context)",
	}

)

func init() {
	app.Flags = []cli.Flag{
		CreateFlag,
		DebugFlag,
		VerbosityFlag,
		CodeFlag,
		CodeFileFlag,
		GasFlag,
		GasPriceFlag,
		ValueFlag,
		InputFlag,
		MemProfileFlag,
		StatDumpFlag,
		GenesisFlag,
		MachineFlag,
		SenderFlag,
		ReceiverFlag,
	}
	app.Commands = []cli.Command{
		runCommond,
		unittestCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
