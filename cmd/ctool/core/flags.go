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
