package main

import (
	"github.com/PlatONnetwork/PlatON-Go/cmd/ctool/core"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "platon transaction test util"
	app.Usage = "send a transaction to deploy or invoke contract"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cmd",
			Value: "",
			Usage: "cmd，eg:deploy,invoke,getTxReceipt",
		},

		//deploy param
		cli.StringFlag{
			Name:  "code, c",
			Value: "",
			Usage: "wasm file path",
		},

		//invoke contract param
		cli.StringFlag{
			Name:  "address, addr",
			Value: "",
			Usage: "the contract address",
		},
		cli.StringFlag{
			Name:  "func, f",
			Value: "",
			Usage: "function and param ,eg :set(1,\"a\")",
		},

		//invoke deploy param
		cli.StringFlag{
			Name:  "abi, a",
			Value: "",
			Usage: "abi file path",
		},
		cli.StringFlag{
			Name:  "config, cfg",
			Value: "",
			Usage: "config path",
		},

		//get tx receipt
		cli.StringFlag{
			Name:  "hash",
			Value: "",
			Usage: "tx hash",
		},
	}
	app.Action = func(c *cli.Context) error {
		cmd := c.String("cmd")

		if cmd == "" {
			fmt.Println("cmd can't be empty!")
			return nil
		}

		switch cmd {
		case "deploy":
			deploy(c)
		case "invoke":
			invoke(c)
		case "getTxReceipt":
			getTxReceipt(c)
		}

		return nil
	}

	app.Run(os.Args)
}

func deploy(c *cli.Context) {
	abiPath := c.String("abi")
	codePath := c.String("code")
	configPath := c.String("config")

	//param check
	if abiPath == "" {
		fmt.Println("abi can't be empty!")
		return
	}

	if codePath == "" {
		fmt.Println("code can't be empty!")
		return
	}
	if configPath == "" {
		dir, _ := os.Getwd()
		configPath = dir + "/config.json"
	}

	core.Deploy(abiPath, codePath, configPath)
}

func invoke(c *cli.Context) {
	addr := c.String("addr")
	abiPath := c.String("abi")
	funcParams := c.String("func")
	configPath := c.String("config")

	//param check
	if abiPath == "" {
		fmt.Println("abi can't be empty!")
		return
	}
	if addr == "" {
		fmt.Println("addr can't be empty!")
		return
	}
	if funcParams == "" {
		fmt.Println("func can't be empty!")
		return
	}
	if configPath == "" {
		dir, _ := os.Getwd()
		configPath = dir + "/config.json"
	}

	core.ContractInvoke(addr, abiPath, funcParams, configPath)
}

func getTxReceipt(c *cli.Context) {
	txHash := c.String("hash")
	configPath := c.String("config")
	if txHash == "" {
		fmt.Println("txHash can't be empty!")
		return
	}
	if configPath == "" {
		dir, _ := os.Getwd()
		configPath = dir + "/config.json"
	}

	core.GetTxReceipt(txHash, configPath)
}
