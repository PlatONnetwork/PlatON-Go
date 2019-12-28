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

// Stored the height of the block that was actually issuance
func StorageIncIssuanceNumber(hash common.Hash, snapshotDB snapshotdb.DB, incIssuanceNumber uint64) error {
	if err := snapshotDB.Put(hash, IncIssuanceNumberKey, common.Uint64ToBytes(incIssuanceNumber)); nil != err {
		log.Error("Failed to execute StorageIncIssuanceNumber function", "hash", hash.TerminalString(), "incIssuanceNumber", incIssuanceNumber, "err", err)
		return err
	}
	return nil
}

func LoadIncIssuanceNumber(hash common.Hash, snapshotDB snapshotdb.DB) (uint64, error) {
	incIssuanceNumberByte, err := snapshotDB.Get(hash, IncIssuanceNumberKey)
	if nil != err {
		if err == snapshotdb.ErrNotFound {
			return 0, nil
		}
		log.Error("Failed to execute LoadIncIssuanceNumber function", "hash", hash.TerminalString(), "key", string(IncIssuanceNumberKey), "err", err)
		return 0, err
	}
	return common.BytesToUint64(incIssuanceNumberByte), nil
}

// Determine whether the block height belongs to the last block at the end of the year according to the passed blockNumber
func IsYearEnd(hash common.Hash, blockNumber uint64) (bool, error) {
	number, err := LoadIncIssuanceNumber(hash, snapshotdb.Instance())
	if nil != err {
		return false, err
	}
	if number == blockNumber {
		return true, nil
	}
	return false, nil
}

// Store the expected time for increase issuance
func StorageIncIssuanceTime(hash common.Hash, snapshotDB snapshotdb.DB, incTime int64) error {
	if err := snapshotDB.Put(hash, IncIssuanceTimeKey, common.Int64ToBytes(incTime)); nil != err {
		log.Error("Failed to execute StorageIncIssuanceTime function", "hash", hash.TerminalString(), "key", string(IncIssuanceTimeKey),
			"value", incTime, "err", err)
		return err
	}
	return nil
}

func LoadIncIssuanceTime(hash common.Hash, snapshotDB snapshotdb.DB) (int64, error) {
	incTimeByte, err := snapshotDB.Get(hash, IncIssuanceTimeKey)
	if nil != err {
		if err != snapshotdb.ErrNotFound {
			log.Error("Failed to execute LoadIncIssuanceTime function", "hash", hash.TerminalString(), "key", string(IncIssuanceTimeKey), "err", err)
			return 0, err
		} else {
			return 0, nil
		}
	}
	return common.BytesToInt64(incTimeByte), nil
}
