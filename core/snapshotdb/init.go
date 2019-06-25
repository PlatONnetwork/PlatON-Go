package snapshotdb

import (
	"flag"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/robfig/cron"
	"os"
	"path"
)

const (
	DBPath      = "snapshotdb"
	DBBasePath  = "base"
	currentPath = "current"
	DBTestPath  = "snapshotdb_test"
)

func initDB() {
	var dbPath string
	if flag.Lookup("test.v") == nil {
		dbPath = DBPath
	} else {
		// for test
		log.Debug("run under go test")
		dbPath = path.Join(os.TempDir(), DBTestPath)
	}
	if findCurrent(dbPath) {
		db := new(snapshotDB)
		if err := db.recover(dbPath); err != nil {
			panic(err)
		}
		dbInstance = db
	} else {
		db, err := newDB(dbPath)
		if err != nil {
			panic(err)
		}
		dbInstance = db
	}
	dbInstance.corn = cron.New()
	if err := dbInstance.corn.AddFunc("@every 1s", dbInstance.schedule); err != nil {
		log.Error("corn AddFunc fail", err)
	}
	dbInstance.corn.Start()
}
