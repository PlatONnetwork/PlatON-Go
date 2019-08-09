package restricting

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	RestrictingKeyPrefix    = []byte("RestrictInfo")
	RestrictRecordKeyPrefix = []byte("RestrictRecord")
	EpochPrefix             = []byte("RestrictEpoch")
)

// RestrictingKey used for search restricting info. key: prefix + account
func GetRestrictingKey(account common.Address) []byte {
	return append(RestrictingKeyPrefix, account.Bytes()...)
}

// RestrictingKey used for search restricting entry info. key: prefix + epoch + account
func GetReleaseAmountKey(epoch uint64, account common.Address) []byte {
	release := append(common.Uint64ToBytes(epoch), account.Bytes()...)
	return append(RestrictingKeyPrefix, release...)
}

// ReleaseNumberKey used for search records at target epoch. key: prefix + epoch
func GetReleaseEpochKey(epoch uint64) []byte {
	return append(RestrictRecordKeyPrefix, common.Uint64ToBytes(epoch)...)
}

// ReleaseAccountKey used for search the index of the restricting account in the released account
// list at target epoch. key: prefix + epoch + index
func GetReleaseAccountKey(epoch uint64, index uint32) []byte {
	releaseIndex := append(common.Uint64ToBytes(epoch), common.Uint32ToBytes(index)...)
	return append(RestrictRecordKeyPrefix, releaseIndex...)
}

func GetLatestEpochKey() []byte {
	return append(EpochPrefix, []byte("latest")...)
}
