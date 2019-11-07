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
