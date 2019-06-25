package snapshotdb

import (
	"flag"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/robfig/cron"
	"os"
	"path"
)

const (
	//DBPath path of db
	DBPath = "snapshotdb"
	//DBBasePath path of basedb
	DBBasePath  = "base"
	currentPath = "current"
	//DBTestPath path of testdb
	DBTestPath = "snapshotdb_test"
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
	s, err := openFile(dbPath, false)
	if err != nil {
		panic(err)
	}
	fds, err := s.List(TypeCurrent)
	if err != nil {
		panic(err)
	}
	if len(fds) > 0 {
		db := new(snapshotDB)
		if err := db.recover(s); err != nil {
			panic(err)
		}
		dbInstance = db
	} else {
		db, err := newDB(s)
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
