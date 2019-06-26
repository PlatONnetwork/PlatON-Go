// +build !test

package snapshotdb

import (
	"github.com/PlatONnetwork/PlatON-Go/node"
	"path"
)

const (
	//DBPath path of db
	DBPath = "snapshotdb"
	//DBBasePath path of basedb
	DBBasePath  = "base"
	currentPath = "current"
)

func init() {
	dbpath = path.Join(node.DefaultDataDir(), "platon", DBPath)
}
