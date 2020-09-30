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

package staking

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (
	CanBasePrefixStr           = "CanBase"
	CanMutablePrefixStr        = "CanMut"
	CanPowerPrefixStr          = "Power"
	UnStakeCountKeyStr         = "UnStakeCount"
	UnStakeItemKeyStr          = "UnStakeItem"
	DelegatePrefixStr          = "Del"
	EpochIndexKeyStr           = "EpochIndex"
	EpochValArrPrefixStr       = "EpochValArr"
	RoundIndexKeyStr           = "RoundIndex"
	RoundValArrPrefixStr       = "RoundValArr"
	AccountStakeRcPrefixStr    = "AccStakeRc"
	PPOSHASHStr                = "PPOSHASH"
	RoundValAddrArrPrefixStr   = "RoundValAddrArr"
	RoundAddrBoundaryPrefixStr = "RoundAddrBoundary"
)

var (
	CanBaseKeyPrefix        = []byte(CanBasePrefixStr)
	CanMutableKeyPrefix     = []byte(CanMutablePrefixStr)
	CanPowerKeyPrefix       = []byte(CanPowerPrefixStr)
	UnStakeCountKey         = []byte(UnStakeCountKeyStr)
	UnStakeItemKey          = []byte(UnStakeItemKeyStr)
	DelegateKeyPrefix       = []byte(DelegatePrefixStr)
	EpochIndexKey           = []byte(EpochIndexKeyStr)
	EpochValArrPrefix       = []byte(EpochValArrPrefixStr)
	RoundIndexKey           = []byte(RoundIndexKeyStr)
	RoundValArrPrefix       = []byte(RoundValArrPrefixStr)
	AccountStakeRcPrefix    = []byte(AccountStakeRcPrefixStr)
	PPOSHASHKey             = []byte(PPOSHASHStr)
	RoundValAddrArrPrefix   = []byte(RoundValAddrArrPrefixStr)
	RoundAddrBoundaryPrefix = []byte(RoundAddrBoundaryPrefixStr)

	b104Len = len(math.MaxBig104.Bytes())
)

// CanBase ...

func CanBaseKeyByAddr(addr common.NodeAddress) []byte {
	return append(CanBaseKeyPrefix, addr.Bytes()...)
}
func CanBaseKeyBySuffix(addr []byte) []byte {
	return append(CanBaseKeyPrefix, addr...)
}

func CanMutableKeyByAddr(addr common.NodeAddress) []byte {
	return append(CanMutableKeyPrefix, addr.Bytes()...)
}

func CanMutableKeyBySuffix(addr []byte) []byte {
	return append(CanMutableKeyPrefix, addr...)
}

// the candidate power key
func TallyPowerKey(programVersion uint32, shares *big.Int, stakeBlockNum uint64, stakeTxIndex uint32, nodeID discover.NodeID) []byte {

	// Only sort Major and Minor
	// eg. 1.1.x => 1.1.0
	subVersion := math.MaxInt32 - xutil.CalcVersion(programVersion)
	sortVersion := common.Uint32ToBytes(subVersion)

	priority := new(big.Int).Sub(math.MaxBig104, shares)
	zeros := make([]byte, b104Len)
	prio := append(zeros, priority.Bytes()...)

	id := nodeID.Bytes()

	num := common.Uint64ToBytes(stakeBlockNum)
	txIndex := common.Uint32ToBytes(stakeTxIndex)

	// some index of pivots
	indexPre := len(CanPowerKeyPrefix)
	indexVersion := indexPre + len(sortVersion)
	indexPrio := indexVersion + len(prio)
	indexNum := indexPrio + len(num)
	indexTxIndex := indexNum + len(txIndex)
	size := indexTxIndex + len(id)

	// construct key
	key := make([]byte, size)
	copy(key[:len(CanPowerKeyPrefix)], CanPowerKeyPrefix)
	copy(key[indexPre:indexVersion], sortVersion)
	copy(key[indexVersion:indexPrio], prio)
	copy(key[indexPrio:indexNum], num)
	copy(key[indexNum:indexTxIndex], txIndex)
	copy(key[indexTxIndex:], id)
	return key
}

func GetUnStakeCountKey(epoch uint64) []byte {
	return append(UnStakeCountKey, common.Uint64ToBytes(epoch)...)
}

func GetUnStakeItemKey(epoch, index uint64) []byte {

	epochByte := common.Uint64ToBytes(epoch)
	indexByte := common.Uint64ToBytes(index)

	markPre := len(UnStakeItemKey)
	markEpoch := markPre + len(epochByte)
	size := markEpoch + len(indexByte)

	key := make([]byte, size)
	copy(key[:markPre], UnStakeItemKey)
	copy(key[markPre:markEpoch], epochByte)
	copy(key[markEpoch:], indexByte)

	return key
}

func GetDelegateKey(delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber uint64) []byte {

	delAddrByte := delAddr.Bytes()
	nodeIdByte := nodeId.Bytes()
	stakeNumByte := common.Uint64ToBytes(stakeBlockNumber)

	markPre := len(DelegateKeyPrefix)
	markDelAddr := markPre + len(delAddrByte)
	markNodeId := markDelAddr + len(nodeIdByte)
	size := markNodeId + len(stakeNumByte)

	key := make([]byte, size)
	copy(key[:markPre], DelegateKeyPrefix)
	copy(key[markPre:markDelAddr], delAddrByte)
	copy(key[markDelAddr:markNodeId], nodeIdByte)
	copy(key[markNodeId:], stakeNumByte)

	return key
}

//notice this assume key must right
func DecodeDelegateKey(key []byte) (delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber uint64) {
	delegateKeyPrefixLength := len(DelegateKeyPrefix)
	delAddrLength := len(delAddr) + delegateKeyPrefixLength
	nodeIdLength := len(nodeId) + delAddrLength
	delAddr = common.BytesToAddress(key[delegateKeyPrefixLength:delAddrLength])
	nodeId = discover.MustBytesID(key[delAddrLength:nodeIdLength])
	stakeBlockNumber = common.BytesToUint64(key[nodeIdLength:])
	return
}

func GetDelegateKeyBySuffix(suffix []byte) []byte {
	return append(DelegateKeyPrefix, suffix...)
}

func GetEpochIndexKey() []byte {
	return EpochIndexKey
}

func GetEpochValArrKey(start, end uint64) []byte {
	startByte := common.Uint64ToBytes(start)
	endByte := common.Uint64ToBytes(end)

	markPre := len(EpochValArrPrefix)
	markStart := markPre + len(startByte)
	size := markStart + len(endByte)

	key := make([]byte, size)
	copy(key[:markPre], EpochValArrPrefix)
	copy(key[markPre:markStart], startByte)
	copy(key[markStart:], endByte)

	return key
}

func GetRoundIndexKey() []byte {
	return RoundIndexKey
}

func GetRoundValArrKey(start, end uint64) []byte {
	startByte := common.Uint64ToBytes(start)
	endByte := common.Uint64ToBytes(end)

	markPre := len(RoundValArrPrefix)
	markStart := markPre + len(startByte)
	size := markStart + len(endByte)

	key := make([]byte, size)
	copy(key[:markPre], RoundValArrPrefix)
	copy(key[markPre:markStart], startByte)
	copy(key[markStart:], endByte)

	return key
}

func GetAccountStakeRcKey(addr common.Address) []byte {
	return append(AccountStakeRcPrefix, addr.Bytes()...)
}

func GetPPOSHASHKey() []byte {
	return PPOSHASHKey
}

func GetRoundValAddrArrKey(round uint64) []byte {
	return append(RoundValAddrArrPrefix, common.Uint64ToBytes(round)...)
}

func GetRoundAddrBoundaryKey() []byte {
	return RoundAddrBoundaryPrefix
}
