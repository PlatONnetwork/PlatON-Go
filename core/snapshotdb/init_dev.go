// +build test

package snapshotdb

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/robfig/cron"
	"math/rand"
	"os"
	"path"
	"time"
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
	logger.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	logger.Info("begin test")
	dbpath = path.Join(os.TempDir(), DBPath, fmt.Sprint(rand.Uint64()))
}

// New  create a new snapshotDB
func New() (DB, error) {
	p := path.Join(os.TempDir(), DBPath, fmt.Sprint(rand.Uint64()))
	logger.Info("begin newDB", "path", p)
	s, err := openFile(p, false)
	if err != nil {
		logger.Error(fmt.Sprint("open db file fail:", err))
		return nil, err
	}
	db, err := newDB(s)
	if err != nil {
		logger.Error(fmt.Sprint("new db fail:", err))
		return nil, err
	}
	db.corn = cron.New()
	if err := db.corn.AddFunc("@every 1s", dbInstance.schedule); err != nil {
		logger.Error(fmt.Sprint("new db fail", err))
		return nil, err
	}
	db.corn.Start()
	return db, nil
}
