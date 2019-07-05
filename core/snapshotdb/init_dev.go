// +build test

package snapshotdb

import (
	"github.com/PlatONnetwork/PlatON-Go/log"
	"os"
	"path"
)

const (
	//DBPath path of db
	DBPath = "snapshotdb_test"
	//DBBasePath path of basedb
	DBBasePath  = "base"
	currentPath = "current"
)

func init() {
	logger.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	logger.Info("begin test")
	dbpath = path.Join(os.TempDir(), DBPath)
}
