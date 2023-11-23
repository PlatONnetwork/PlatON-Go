// Copyright 2021 The PlatON Network Authors
// This file is part of PlatON-Go.
//
// PlatON-Go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// PlatON-Go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PlatON-Go. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"github.com/PlatONnetwork/PlatON-Go/eth"

	"github.com/urfave/cli/v2"
)

var (
	AnalyzeStressTestCmd = &cli.Command{
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
