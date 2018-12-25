package core

import "gopkg.in/urfave/cli.v1"

var (
	ConfigPathFlag = cli.StringFlag{
		Name:  "config",
		Usage: "config path",
	}
	PKFilePathFlag = cli.StringFlag{
		Name:  "pkfile",
		Value: "",
		Usage: "private key file path",
	}
	StabExecTimesFlag = cli.IntFlag{
		Name:  "times",
		Value: 1000,
		Usage: "execute times",
	}
	SendTxIntervalFlag = cli.IntFlag{
		Name:  "interval",
		Value: 10,
		Usage: "Time interval for sending transactions",
	}
	AccountSizeFlag = cli.IntFlag{
		Name:  "size",
		Value: 10,
		Usage: "account size",
	}
	TxJsonDataFlag = cli.StringFlag{
		Name:  "data",
		Usage: "transaction data",
	}
	ContractWasmFilePathFlag = cli.StringFlag{
		Name:  "code",
		Usage: "wasm file path",
	}
	ContractAddrFlag = cli.StringFlag{
		Name: "addr",

		Usage: "the contract address",
	}
	ContractFuncNameFlag = cli.StringFlag{
		Name:  "func",
		Usage: "function and param ,eg :set(1,\"a\")",
	}
	TransactionTypeFlag = cli.IntFlag{
		Name:  "type",
		Value: 2,
		Usage: "tx type ,default 2",
	}
	ContractAbiFilePathFlag = cli.StringFlag{
		Name:  "abi",
		Usage: "abi file path",
	}
	TransactionHashFlag = cli.StringFlag{
		Name:  "hash",
		Usage: "tx hash",
	}
	TxFromFlag = cli.StringFlag{
		Name:  "from",
		Usage: "transaction sender addr",
	}
	TxToFlag = cli.StringFlag{
		Name:  "to",
		Usage: "transaction acceptor addr",
	}
	TransferValueFlag = cli.StringFlag{
		Name:  "value",
		Value: "0xDE0B6B3A7640000", //one
		Usage: "transfer value",
	}

	deployCmdFlags = []cli.Flag{
		ContractWasmFilePathFlag,
		ContractAbiFilePathFlag,
		ConfigPathFlag,
	}
	invokeCmdFlags = []cli.Flag{
		ContractFuncNameFlag,
		ContractAbiFilePathFlag,
		ContractAddrFlag,
		ConfigPathFlag,
		TransactionTypeFlag,
	}

	sendTransactionCmdFlags = []cli.Flag{
		TxFromFlag,
		TxToFlag,
		TransferValueFlag,
		ConfigPathFlag,
	}
	sendRawTransactionCmdFlags = []cli.Flag{
		PKFilePathFlag,
		TxFromFlag,
		TxToFlag,
		TransferValueFlag,
		ConfigPathFlag,
	}
	getTxReceiptCmdFlags = []cli.Flag{
		TransactionHashFlag,
		ConfigPathFlag,
	}

	stabilityCmdFlags = []cli.Flag{
		PKFilePathFlag,
		StabExecTimesFlag,
		SendTxIntervalFlag,
		ConfigPathFlag,
	}
	stabPrepareCmdFlags = []cli.Flag{
		PKFilePathFlag,
		AccountSizeFlag,
		TransferValueFlag,
		ConfigPathFlag,
	}
)
