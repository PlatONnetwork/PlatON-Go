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
	"fmt"

	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

type StakingDB struct {
	db snapshotdb.DB
}

func NewStakingDB() *StakingDB {
	return &StakingDB{
		db: snapshotdb.Instance(),
	}
}

func NewStakingDBWithDB(db snapshotdb.DB) *StakingDB {
	return &StakingDB{
		db: db,
	}
}

func (db *StakingDB) GetDB() snapshotdb.DB {
	return db.db
}

func (db *StakingDB) get(blockHash common.Hash, key []byte) ([]byte, error) {
	return db.db.Get(blockHash, key)
}

func (db *StakingDB) getFromCommitted(key []byte) ([]byte, error) {
	return db.db.GetFromCommittedBlock(key)
}

func (db *StakingDB) put(blockHash common.Hash, key, value []byte) error {
	return db.db.Put(blockHash, key, value)
}

func (db *StakingDB) del(blockHash common.Hash, key []byte) error {
	return db.db.Del(blockHash, key)
}

func (db *StakingDB) ranking(blockHash common.Hash, prefix []byte, ranges int) iterator.Iterator {
	return db.db.Ranking(blockHash, prefix, ranges)
}
func (db *StakingDB) Del(blockHash common.Hash, key []byte) error {
	return db.db.Del(blockHash, key)
}

func (db *StakingDB) GetLastKVHash(blockHash common.Hash) []byte {
	return db.db.GetLastKVHash(blockHash)
}

// about candidate ...

func (db *StakingDB) GetCandidateStore(blockHash common.Hash, addr common.NodeAddress) (*Candidate, error) {
	base, err := db.GetCanBaseStore(blockHash, addr)
	if nil != err {
		return nil, err
	}
	mutable, err := db.GetCanMutableStore(blockHash, addr)
	if nil != err {
		return nil, err
	}

	can := &Candidate{}
	can.CandidateBase = base
	can.CandidateMutable = mutable
	return can, nil
}

func (db *StakingDB) GetCandidateStoreByIrr(addr common.NodeAddress) (*Candidate, error) {
	base, err := db.GetCanBaseStoreByIrr(addr)
	if nil != err {
		return nil, err
	}
	mutable, err := db.GetCanMutableStoreByIrr(addr)
	if nil != err {
		return nil, err
	}

	can := &Candidate{}
	can.CandidateBase = base
	can.CandidateMutable = mutable
	return can, nil
}

func (db *StakingDB) GetCandidateStoreWithSuffix(blockHash common.Hash, suffix []byte) (*Candidate, error) {
	base, err := db.GetCanBaseStoreWithSuffix(blockHash, suffix)
	if nil != err {
		return nil, err
	}
	mutable, err := db.GetCanMutableStoreWithSuffix(blockHash, suffix)
	if nil != err {
		return nil, err
	}

	can := &Candidate{}
	can.CandidateBase = base
	can.CandidateMutable = mutable
	return can, nil
}

func (db *StakingDB) GetCandidateStoreByIrrWithSuffix(suffix []byte) (*Candidate, error) {
	base, err := db.GetCanBaseStoreByIrrWithSuffix(suffix)
	if nil != err {
		return nil, err
	}
	mutable, err := db.GetCanMutableStoreByIrrWithSuffix(suffix)
	if nil != err {
		return nil, err
	}

	can := &Candidate{}
	can.CandidateBase = base
	can.CandidateMutable = mutable
	return can, nil
}

func (db *StakingDB) SetCandidateStore(blockHash common.Hash, addr common.NodeAddress, can *Candidate) error {

	if err := db.SetCanBaseStore(blockHash, addr, can.CandidateBase); nil != err {
		return err
	}
	if err := db.SetCanMutableStore(blockHash, addr, can.CandidateMutable); nil != err {
		return err
	}
	return nil
}

func (db *StakingDB) DelCandidateStore(blockHash common.Hash, addr common.NodeAddress) error {
	if err := db.DelCanBaseStore(blockHash, addr); nil != err {
		return err
	}
	if err := db.DelCanMutableStore(blockHash, addr); nil != err {
		return err
	}
	return nil
}

// about canbase ...

func (db *StakingDB) GetCanBaseStore(blockHash common.Hash, addr common.NodeAddress) (*CandidateBase, error) {

	key := CanBaseKeyByAddr(addr)

	canByte, err := db.get(blockHash, key)

	if nil != err {
		return nil, err
	}

	var can CandidateBase
	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}

	return &can, nil
}

func (db *StakingDB) GetCanBaseStoreByIrr(addr common.NodeAddress) (*CandidateBase, error) {
	key := CanBaseKeyByAddr(addr)
	canByte, err := db.getFromCommitted(key)

	if nil != err {
		return nil, err
	}
	var can CandidateBase

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCanBaseStoreWithSuffix(blockHash common.Hash, suffix []byte) (*CandidateBase, error) {
	key := CanBaseKeyBySuffix(suffix)

	canByte, err := db.get(blockHash, key)

	if nil != err {
		return nil, err
	}
	var can CandidateBase

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCanBaseStoreByIrrWithSuffix(suffix []byte) (*CandidateBase, error) {
	key := CanBaseKeyBySuffix(suffix)
	canByte, err := db.getFromCommitted(key)

	if nil != err {
		return nil, err
	}
	var can CandidateBase

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) SetCanBaseStore(blockHash common.Hash, addr common.NodeAddress, can *CandidateBase) error {

	key := CanBaseKeyByAddr(addr)

	if val, err := rlp.EncodeToBytes(can); nil != err {
		return err
	} else {

		return db.put(blockHash, key, val)
	}
}

func (db *StakingDB) DelCanBaseStore(blockHash common.Hash, addr common.NodeAddress) error {
	key := CanBaseKeyByAddr(addr)
	return db.del(blockHash, key)
}

// about canmutable ...

func (db *StakingDB) GetCanMutableStore(blockHash common.Hash, addr common.NodeAddress) (*CandidateMutable, error) {

	key := CanMutableKeyByAddr(addr)

	canByte, err := db.get(blockHash, key)

	if nil != err {
		return nil, err
	}

	var can CandidateMutable
	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}

	return &can, nil
}

func (db *StakingDB) GetCanMutableStoreByIrr(addr common.NodeAddress) (*CandidateMutable, error) {
	key := CanMutableKeyByAddr(addr)
	canByte, err := db.getFromCommitted(key)

	if nil != err {
		return nil, err
	}
	var can CandidateMutable

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCanMutableStoreWithSuffix(blockHash common.Hash, suffix []byte) (*CandidateMutable, error) {
	key := CanMutableKeyBySuffix(suffix)

	canByte, err := db.get(blockHash, key)

	if nil != err {
		return nil, err
	}
	var can CandidateMutable

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCanMutableStoreByIrrWithSuffix(suffix []byte) (*CandidateMutable, error) {
	key := CanMutableKeyBySuffix(suffix)
	canByte, err := db.getFromCommitted(key)

	if nil != err {
		return nil, err
	}
	var can CandidateMutable

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) SetCanMutableStore(blockHash common.Hash, addr common.NodeAddress, can *CandidateMutable) error {

	key := CanMutableKeyByAddr(addr)

	if val, err := rlp.EncodeToBytes(can); nil != err {
		return err
	} else {

		return db.put(blockHash, key, val)
	}
}

func (db *StakingDB) DelCanMutableStore(blockHash common.Hash, addr common.NodeAddress) error {
	key := CanMutableKeyByAddr(addr)
	return db.del(blockHash, key)
}

// about candidate power ...

func (db *StakingDB) SetCanPowerStore(blockHash common.Hash, addr common.NodeAddress, can *Candidate) error {

	key := TallyPowerKey(can.ProgramVersion, can.Shares, can.StakingBlockNum, can.StakingTxIndex, can.NodeId)

	return db.put(blockHash, key, addr.Bytes())
}

func (db *StakingDB) DelCanPowerStore(blockHash common.Hash, can *Candidate) error {

	key := TallyPowerKey(can.ProgramVersion, can.Shares, can.StakingBlockNum, can.StakingTxIndex, can.NodeId)
	return db.del(blockHash, key)
}

// about UnStakeItem ...

func (db *StakingDB) AddUnStakeItemStore(blockHash common.Hash, epoch uint64, canAddr common.NodeAddress, stakeBlockNumber uint64, recovery bool) error {

	count_key := GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	var v uint64
	switch {
	case snapshotdb.NonDbNotFoundErr(err):
		return err
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	v++

	if err := db.put(blockHash, count_key, common.Uint64ToBytes(v)); nil != err {
		return err
	}
	item_key := GetUnStakeItemKey(epoch, v)

	unStakeItem := &UnStakeItem{
		NodeAddress:     canAddr,
		StakingBlockNum: stakeBlockNumber,
		Recovery:        recovery,
	}

	item, err := rlp.EncodeToBytes(unStakeItem)
	if nil != err {
		return err
	}

	return db.put(blockHash, item_key, item)
}

func (db *StakingDB) GetUnStakeCountStore(blockHash common.Hash, epoch uint64) (uint64, error) {
	count_key := GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return 0, err
	}
	return common.BytesToUint64(val), nil
}

func (db *StakingDB) GetUnStakeItemStore(blockHash common.Hash, epoch, index uint64) (*UnStakeItem, error) {
	item_key := GetUnStakeItemKey(epoch, index)
	itemByte, err := db.get(blockHash, item_key)
	if nil != err {
		return nil, err
	}

	var unStakeItem UnStakeItem
	if err := rlp.DecodeBytes(itemByte, &unStakeItem); nil != err {
		return nil, err
	}
	return &unStakeItem, nil
}

func (db *StakingDB) DelUnStakeCountStore(blockHash common.Hash, epoch uint64) error {
	count_key := GetUnStakeCountKey(epoch)

	return db.del(blockHash, count_key)
}

func (db *StakingDB) DelUnStakeItemStore(blockHash common.Hash, epoch, index uint64) error {
	item_key := GetUnStakeItemKey(epoch, index)

	return db.del(blockHash, item_key)
}

// about delegate ...

func (db *StakingDB) GetDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber uint64) (*Delegation, error) {

	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	delByte, err := db.get(blockHash, key)
	if nil != err {
		return nil, err
	}

	var del Delegation
	if err := rlp.DecodeBytes(delByte, &del); nil != err {
		return nil, err
	}

	return &del, nil
}

func (db *StakingDB) GetDelegateStoreByIrr(delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber uint64) (*Delegation, error) {
	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	delByte, err := db.getFromCommitted(key)
	if nil != err {
		return nil, err
	}

	var del Delegation
	if err := rlp.DecodeBytes(delByte, &del); nil != err {
		return nil, err
	}
	return &del, nil
}

func (db *StakingDB) GetDelegateStoreBySuffix(blockHash common.Hash, keySuffix []byte) (*Delegation, error) {
	key := GetDelegateKeyBySuffix(keySuffix)
	delByte, err := db.get(blockHash, key)
	if nil != err {
		return nil, err
	}

	var del Delegation
	if err := rlp.DecodeBytes(delByte, &del); nil != err {
		return nil, err
	}
	return &del, nil
}

type DelegationInfo struct {
	NodeID           discover.NodeID
	StakeBlockNumber uint64
	Delegation       *Delegation
}

type DelByDelegateEpoch []*DelegationInfo

func (d DelByDelegateEpoch) Len() int { return len(d) }
func (d DelByDelegateEpoch) Less(i, j int) bool {
	return d[i].Delegation.DelegateEpoch < d[j].Delegation.DelegateEpoch
}
func (d DelByDelegateEpoch) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func (db *StakingDB) GetDelegatesInfo(blockHash common.Hash, delAddr common.Address) ([]*DelegationInfo, error) {
	key := GetDelegateKeyBySuffix(delAddr.Bytes())
	itr := db.ranking(blockHash, key, 0)
	if itr.Error() != nil {
		return nil, itr.Error()
	}
	infos := make([]*DelegationInfo, 0)
	for itr.Next() {
		info := new(DelegationInfo)
		_, info.NodeID, info.StakeBlockNumber = DecodeDelegateKey(itr.Key())
		info.Delegation = new(Delegation)
		if err := rlp.DecodeBytes(itr.Value(), info.Delegation); err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func (db *StakingDB) SetDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	stakeBlockNumber uint64, del *Delegation) error {

	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	delByte, err := rlp.EncodeToBytes(del)
	if nil != err {
		return err
	}

	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) SetDelegateStoreBySuffix(blockHash common.Hash, suffix []byte, del *Delegation) error {
	key := GetDelegateKeyBySuffix(suffix)
	delByte, err := rlp.EncodeToBytes(del)
	if nil != err {
		return err
	}

	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) DelDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	stakeBlockNumber uint64) error {
	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	return db.del(blockHash, key)
}

func (db *StakingDB) DelDelegateStoreBySuffix(blockHash common.Hash, suffix []byte) error {
	key := GetDelegateKeyBySuffix(suffix)

	return db.del(blockHash, key)
}

// about epoch validates ...

func (db *StakingDB) SetEpochValIndex(blockHash common.Hash, indexArr ValArrIndexQueue) error {
	value, err := rlp.EncodeToBytes(indexArr)
	if nil != err {
		return err
	}

	return db.put(blockHash, GetEpochIndexKey(), value)
}

func (db *StakingDB) GetEpochValIndexByBlockHash(blockHash common.Hash) (ValArrIndexQueue, error) {
	val, err := db.get(blockHash, GetEpochIndexKey())
	if nil != err {
		return nil, err
	}
	var queue ValArrIndexQueue
	if err := rlp.DecodeBytes(val, &queue); nil != err {
		return nil, err
	}
	return queue, nil
}

func (db *StakingDB) GetEpochValIndexByIrr() (ValArrIndexQueue, error) {
	val, err := db.getFromCommitted(GetEpochIndexKey())
	if nil != err {
		return nil, err
	}
	var queue ValArrIndexQueue
	if err := rlp.DecodeBytes(val, &queue); nil != err {
		return nil, err
	}
	return queue, nil
}

func (db *StakingDB) SetEpochValList(blockHash common.Hash, start, end uint64, valArr ValidatorQueue) error {

	value, err := rlp.EncodeToBytes(valArr)
	if nil != err {
		return err
	}

	return db.put(blockHash, GetEpochValArrKey(start, end), value)
}

func (db *StakingDB) GetEpochValListByIrr(start, end uint64) (ValidatorQueue, error) {
	arrByte, err := db.getFromCommitted(GetEpochValArrKey(start, end))
	if nil != err {
		return nil, err
	}

	var arr ValidatorQueue
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) GetEpochValListByBlockHash(blockHash common.Hash, start, end uint64) (ValidatorQueue, error) {
	arrByte, err := db.get(blockHash, GetEpochValArrKey(start, end))
	if nil != err {
		return nil, err
	}

	var arr ValidatorQueue
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) DelEpochValListByBlockHash(blockHash common.Hash, start, end uint64) error {

	return db.del(blockHash, GetEpochValArrKey(start, end))
}

// about round validators

func (db *StakingDB) SetRoundValIndex(blockHash common.Hash, indexArr ValArrIndexQueue) error {
	value, err := rlp.EncodeToBytes(indexArr)
	if nil != err {
		return err
	}

	return db.put(blockHash, GetRoundIndexKey(), value)
}

func (db *StakingDB) GetRoundValIndexByBlockHash(blockHash common.Hash) (ValArrIndexQueue, error) {
	val, err := db.get(blockHash, GetRoundIndexKey())
	if nil != err {
		return nil, err
	}
	var queue ValArrIndexQueue
	if err := rlp.DecodeBytes(val, &queue); nil != err {
		return nil, err
	}
	return queue, nil
}

func (db *StakingDB) GetRoundValIndexByIrr() (ValArrIndexQueue, error) {
	val, err := db.getFromCommitted(GetRoundIndexKey())
	if nil != err {
		return nil, err
	}
	var queue ValArrIndexQueue
	if err := rlp.DecodeBytes(val, &queue); nil != err {
		return nil, err
	}
	return queue, nil
}

func (db *StakingDB) SetRoundValList(blockHash common.Hash, start, end uint64, valArr ValidatorQueue) error {

	value, err := rlp.EncodeToBytes(valArr)
	if nil != err {
		return err
	}

	return db.put(blockHash, GetRoundValArrKey(start, end), value)
}

func (db *StakingDB) GetRoundValListByIrr(start, end uint64) (ValidatorQueue, error) {
	arrByte, err := db.getFromCommitted(GetRoundValArrKey(start, end))
	if nil != err {
		return nil, err
	}

	var arr ValidatorQueue
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) GetRoundValListByBlockHash(blockHash common.Hash, start, end uint64) (ValidatorQueue, error) {
	arrByte, err := db.get(blockHash, GetRoundValArrKey(start, end))
	if nil != err {
		return nil, err
	}

	var arr ValidatorQueue
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) DelRoundValListByBlockHash(blockHash common.Hash, start, end uint64) error {

	return db.del(blockHash, GetRoundValArrKey(start, end))
}

// iterator ...

func (db *StakingDB) IteratorCandidatePowerByBlockHash(blockHash common.Hash, ranges int) iterator.Iterator {
	return db.ranking(blockHash, CanPowerKeyPrefix, ranges)
}

func (db *StakingDB) IteratorDelegateByBlockHashWithAddr(blockHash common.Hash, addr common.Address, ranges int) iterator.Iterator {
	prefix := append(DelegateKeyPrefix, addr.Bytes()...)
	return db.ranking(blockHash, prefix, ranges)
}

// about account staking reference count ...

func (db *StakingDB) AddAccountStakeRc(blockHash common.Hash, addr common.Address) error {
	key := GetAccountStakeRcKey(addr)
	val, err := db.get(blockHash, key)
	var v uint64
	switch {
	case snapshotdb.NonDbNotFoundErr(err):
		return err
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	v++

	return db.put(blockHash, key, common.Uint64ToBytes(v))
}

func (db *StakingDB) SubAccountStakeRc(blockHash common.Hash, addr common.Address) error {
	key := GetAccountStakeRcKey(addr)
	val, err := db.get(blockHash, key)
	var v uint64
	switch {
	case snapshotdb.NonDbNotFoundErr(err):
		return err
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	// Prevent large numbers from being directly called after the uint64 overflow
	if v == 0 {
		return nil
	}

	v--

	if v == 0 {

		return db.del(blockHash, key)
	} else {

		return db.put(blockHash, key, common.Uint64ToBytes(v))
	}
}

func (db *StakingDB) HasAccountStakeRc(blockHash common.Hash, addr common.Address) (bool, error) {
	key := GetAccountStakeRcKey(addr)
	val, err := db.get(blockHash, key)
	var v uint64
	switch {
	case snapshotdb.NonDbNotFoundErr(err):
		return false, err
	case snapshotdb.IsDbNotFoundErr(err):
		return false, nil
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	if v == 0 {
		return false, nil
	} else if v > 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("Account Stake Reference Count cannot be negative, account: %s", addr.String())
	}
}

// about round validator's addrs ...

func (db *StakingDB) StoreRoundValidatorAddrs(blockHash common.Hash, key []byte, arry []common.NodeAddress) error {
	value, err := rlp.EncodeToBytes(arry)
	if nil != err {
		return err
	}
	return db.put(blockHash, key, value)
}

func (db *StakingDB) DelRoundValidatorAddrs(blockHash common.Hash, key []byte) error {
	return db.del(blockHash, key)
}

func (db *StakingDB) LoadRoundValidatorAddrs(blockHash common.Hash, key []byte) ([]common.Address, error) {
	rlpValue, err := db.get(blockHash, key)
	if nil != err {
		return nil, err
	}
	var value []common.Address
	if err := rlp.DecodeBytes(rlpValue, &value); nil != err {
		return nil, err
	}
	return value, nil
}

func (db *StakingDB) SetRoundAddrBoundary(blockHash common.Hash, round uint64) error {
	return db.put(blockHash, GetRoundAddrBoundaryKey(), common.Uint64ToBytes(round))
}

func (db *StakingDB) GetRoundAddrBoundary(blockHash common.Hash) (uint64, error) {
	round, err := db.get(blockHash, GetRoundAddrBoundaryKey())
	if nil != err {
		return 0, err
	}
	return common.BytesToUint64(round), nil
}
