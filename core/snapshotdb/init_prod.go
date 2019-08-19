// +build !test

package snapshotdb

import "errors"

const (
	//DBPath path of db
	DBPath = "snapshotdb"
	//DBBasePath path of basedb
	DBBasePath  = "base"
	currentPath = "current"
)

// New  create a new snapshotDB,will clear old snapshotDB data
func New() (DB, error) {
	if dbInstance != nil {
		if err := dbInstance.Clear(); err != nil {
			return nil, err
		}
	}
	if err := initDB(); err != nil {
		return nil, errors.New("init db fail:" + err.Error())
	}
	return dbInstance, nil
}
