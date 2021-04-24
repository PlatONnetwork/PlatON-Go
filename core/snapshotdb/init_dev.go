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

// +build test

package snapshotdb

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

const (
	//DBPath path of db
	DBPath = "snapshotdb_test"
	//DBBasePath path of basedb
	DBBasePath  = "base"
	currentPath = "current"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	//		logger.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	logger.Info("begin test")
	dbpath = path.Join(os.TempDir(), DBPath, fmt.Sprint(rand.Uint64()))
	testChain := new(testchain)
	header := generateHeader(big.NewInt(1000000000), common.ZeroHash)
	testChain.h = append(testChain.h, header)
	SetDBBlockChain(testChain)
}
