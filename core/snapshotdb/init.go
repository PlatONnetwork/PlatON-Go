package snapshotdb

import (
	"flag"
	"fmt"
	"github.com/robfig/cron"
)

//todo 需要确认文件路径与操作系统
const (
	DBPath      = "./snapshotdb"
	DBBasePath  = "base"
	currentPath = "current"
)

func init() {
	if flag.Lookup("test.v") == nil {
		if findCurrent(DBPath) {
			db := new(snapshotDB)
			if err := db.recover(DBPath); err != nil {
				panic(err)
			}
			dbInstance = db
		} else {
			db, err := newDB(DBPath)
			if err != nil {
				panic(err)
			}
			dbInstance = db
		}
		go func() {
			c := cron.New()
			c.AddFunc("@every 1s", dbInstance.schedule)
			c.Start()
		}()
	} else {
		fmt.Println("run under go test")
	}

}
