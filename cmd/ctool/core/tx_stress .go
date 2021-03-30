package core

import (
	"github.com/PlatONnetwork/PlatON-Go/eth"

	"gopkg.in/urfave/cli.v1"
)

var (
	AnalyzeStressTestCmd = cli.Command{
		Name:   "analyzeStressTest",
		Usage:  "analyze the tx stress test source file to generate  result data",
		Action: analyzeStressTest,
		Flags:  txStressFlags,
	}
)

func analyzeStressTest(c *cli.Context) error {
	configPaths := c.StringSlice(TxStressSourceFilesPathFlag.Name)
	t := c.Int(TxStressStatisticTimeFlag.Name)
	output := c.String(TxStressOutPutFileFlag.Name)
	return eth.AnalyzeStressTest(configPaths, output, t)
}
