// +build !test

package snapshotdb

import (
	"fmt"
	"os"
)

const (
	//DBPath path of db
	DBPath = "snapshotdb"
	//DBBasePath path of basedb
	DBBasePath  = "base"
	currentPath = "current"
)

// New  create a new snapshotDB,will clear old snapshotDB data
func New(path string) (DB, error) {
	if err := os.RemoveAll(path); err != nil {
		return nil, err
	}
	s, err := openFile(path, false)
	if err != nil {
		logger.Error("open db file fail", "error", err, "path", dbpath)
		return nil, err
	}

	logger.Info("begin new", "path", path)
	db, err := newDB(s)
	if err != nil {
		logger.Error(fmt.Sprint("new db fail:", err))
		return nil, err
	}
	return db, nil
}
