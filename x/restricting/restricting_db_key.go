// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package restricting

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	RestrictingKeyPrefix         = []byte("RestrictInfo")
	RestrictRecordKeyPrefix      = []byte("RestrictRecord")
	InitialFoundationRestricting = []byte("InitialFoundationRestricting")
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
