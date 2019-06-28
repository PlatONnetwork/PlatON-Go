// +build test

package snapshotdb

import (
	"flag"
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
	if flag.Lookup("test.bench") == nil {
		logger.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	} else {
		logger.SetHandler(log.DiscardHandler())
	}
	logger.Info("begin test")
	dbpath = path.Join(os.TempDir(), DBPath)
}
