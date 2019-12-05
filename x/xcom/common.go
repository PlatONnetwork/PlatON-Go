package xcom

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

// saves block average pack time (millisecond) to snapshot db.
func StorageAvgPackTime(hash common.Hash, snapshotDB snapshotdb.DB, avgPackTime uint64) error {
	if err := snapshotDB.Put(hash, AvgPackTimeKey, common.Uint64ToBytes(avgPackTime)); nil != err {
		log.Error("Failed to save block average pack time", "hash", hash.TerminalString(), "avgPackTime", avgPackTime, "err", err)
		return err
	}
	return nil
}

// gets block average pack time (millisecond) from snapshot db.
func LoadCurrentAvgPackTime() (uint64, error) {
	return LoadAvgPackTime(common.ZeroHash, snapshotdb.Instance())
}

// gets block average pack time (millisecond) from snapshot db.
func LoadAvgPackTime(hash common.Hash, snapshotDB snapshotdb.DB) (uint64, error) {
	avgPackTimeByte, err := snapshotDB.Get(hash, AvgPackTimeKey)

	if nil != err {
		if err == snapshotdb.ErrNotFound {
			return 0, nil
		}
		log.Error("Failed to load block average pack time", "hash", hash.TerminalString(), "key", string(AvgPackTimeKey), "err", err)
		return 0, err
	}
	return common.BytesToUint64(avgPackTimeByte), nil
}
