package xcom

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	RestrictingKeyPrefix    = []byte("RestrictInfo")
	RestrictRecordKeyPrefix = []byte("RestrictRecord")
)

// RestrictingKey used for search restricting info. key: prefix + account
func GetRestrictingKey(account common.Address) []byte {
	return append(RestrictingKeyPrefix, account.Bytes()...)
}

// RestrictingKey used for search restricting entry info. key: prefix + account + blockNum
func GetReleaseAmountKey(account common.Address, blockNum uint64) []byte {
	release := append(account.Bytes(), common.Uint64ToBytes(blockNum)...)
	return append(RestrictingKeyPrefix, release...)
}

// ReleaseNumberKey used for search records at target release blockNumber. key: prefix + blockNum
func GetReleaseNumberKey(blockNum uint64) []byte {
	return append(RestrictRecordKeyPrefix, common.Uint64ToBytes(blockNum)...)
}

// ReleaseAccountKey used for search restricting account at target block index. key: prefix + blockNum + index
func GetReleaseAccountKey(blockNum uint64, index uint32) []byte {
	releaseIndex := append(common.Uint64ToBytes(blockNum), common.Uint32ToBytes(index)...)
	return append(RestrictRecordKeyPrefix, releaseIndex...)
}
