package staking

import (
	"encoding/hex"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type StakingDB struct {
	db snapshotdb.DB
}

func NewStakingDB() *StakingDB {
	return &StakingDB{
		db: snapshotdb.Instance(),
	}
}

func (db *StakingDB) get(blockHash common.Hash, key []byte) ([]byte, error) {
	return db.db.Get(blockHash, key)
}

func (db *StakingDB) getFromCommitted(key []byte) ([]byte, error) {
	return db.db.GetFromCommittedBlock(key)
}

func (db *StakingDB) put(blockHash common.Hash, key, value []byte) error {

	//// todo test
	//log.Debug("put", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key), "val", hex.EncodeToString(value))

	return db.db.Put(blockHash, key, value)
}

func (db *StakingDB) del(blockHash common.Hash, key []byte) error {

	//// todo test
	//log.Debug("del", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key))

	return db.db.Del(blockHash, key)
}

func (db *StakingDB) ranking(blockHash common.Hash, prefix []byte, ranges int) iterator.Iterator {
	return db.db.Ranking(blockHash, prefix, ranges)
}

func (db *StakingDB) GetLastKVHash(blockHash common.Hash) []byte {
	return db.db.GetLastKVHash(blockHash)
}

func (db *StakingDB) GetCandidateStore(blockHash common.Hash, addr common.Address) (*Candidate, error) {
	key := CandidateKeyByAddr(addr)

	canByte, err := db.get(blockHash, key)

	if nil != err {
		return nil, err
	}
	var can Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCandidateStoreByIrr(addr common.Address) (*Candidate, error) {
	key := CandidateKeyByAddr(addr)
	canByte, err := db.getFromCommitted(key)

	if nil != err {
		return nil, err
	}
	var can Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCandidateStoreWithSuffix(blockHash common.Hash, suffix []byte) (*Candidate, error) {
	key := CandidateKeyBySuffix(suffix)

	// todo test
	log.Debug("GetCandidateStoreWithSuffix: Query can", "key", hex.EncodeToString(key))

	canByte, err := db.get(blockHash, key)

	if nil != err {
		return nil, err
	}
	var can Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCandidateStoreByIrrWithSuffix(suffix []byte) (*Candidate, error) {
	key := CandidateKeyBySuffix(suffix)
	canByte, err := db.getFromCommitted(key)

	if nil != err {
		return nil, err
	}
	var can Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) SetCandidateStore(blockHash common.Hash, addr common.Address, can *Candidate) error {

	key := CandidateKeyByAddr(addr)

	if val, err := rlp.EncodeToBytes(can); nil != err {
		return err
	} else {

		// todo test
		log.Debug("SetCandidateStore", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key), "val", hex.EncodeToString(val))

		return db.put(blockHash, key, val)
	}
}

func (db *StakingDB) DelCandidateStore(blockHash common.Hash, addr common.Address) error {
	key := CandidateKeyByAddr(addr)

	// todo test
	log.Debug("DelCandidateStore", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key))

	return db.del(blockHash, key)
}

func (db *StakingDB) SetCanPowerStore(blockHash common.Hash, addr common.Address, can *Candidate) error {
	key := TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex, can.ProgramVersion)

	// todo test
	log.Debug("SetCanPowerStore", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key), "val", hex.EncodeToString(addr.Bytes()))

	return db.put(blockHash, key, addr.Bytes())
}

func (db *StakingDB) DelCanPowerStore(blockHash common.Hash, can *Candidate) error {
	key := TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex, can.ProgramVersion)

	// todo test
	log.Debug("DelCanPowerStore", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key))

	return db.del(blockHash, key)
}

func (db *StakingDB) AddUnStakeItemStore(blockHash common.Hash, epoch uint64, canAddr common.Address, stakeBlockNumber uint64) error {

	count_key := GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	var v uint64
	switch {
	case nil != err && err != snapshotdb.ErrNotFound:
		return err
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	v++

	// todo test
	log.Debug("AddUnStakeItemStore before put count ", "blockHash", blockHash.Hex(), "count_key", hex.EncodeToString(count_key), "val", hex.EncodeToString(common.Uint64ToBytes(v)))

	if err := db.put(blockHash, count_key, common.Uint64ToBytes(v)); nil != err {
		return err
	}
	item_key := GetUnStakeItemKey(epoch, v)

	unStakeItem := &UnStakeItem{
		NodeAddress:     canAddr,
		StakingBlockNum: stakeBlockNumber,
	}

	item, err := rlp.EncodeToBytes(unStakeItem)
	if nil != err {
		return err
	}

	// todo test
	log.Debug("AddUnStakeItemStore before put item ", "blockHash", blockHash.Hex(), "item_key", hex.EncodeToString(item_key), "val", hex.EncodeToString(item))

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
	// todo test
	log.Debug("DelUnStakeCountStore", "blockHash", blockHash.Hex(), "count_key", hex.EncodeToString(count_key))

	return db.del(blockHash, count_key)
}

func (db *StakingDB) DelUnStakeItemStore(blockHash common.Hash, epoch, index uint64) error {
	item_key := GetUnStakeItemKey(epoch, index)

	// todo test
	log.Debug("DelUnStakeItemStore", "blockHash", blockHash.Hex(), "item_key", hex.EncodeToString(item_key))

	return db.del(blockHash, item_key)
}

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

func (db *StakingDB) SetDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	stakeBlockNumber uint64, del *Delegation) error {
	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	delByte, err := rlp.EncodeToBytes(del)
	if nil != err {
		return err
	}

	// todo test
	log.Debug("SetDelegateStore", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key), "val", hex.EncodeToString(delByte))

	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) SetDelegateStoreBySuffix(blockHash common.Hash, suffix []byte, del *Delegation) error {
	key := GetDelegateKeyBySuffix(suffix)
	delByte, err := rlp.EncodeToBytes(del)
	if nil != err {
		return err
	}

	// todo test
	log.Debug("SetDelegateStoreBySuffix", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key), "val", hex.EncodeToString(delByte))

	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) DelDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	stakeBlockNumber uint64) error {
	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	// todo test
	log.Debug("DelDelegateStore", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key))

	return db.del(blockHash, key)
}

func (db *StakingDB) DelDelegateStoreBySuffix(blockHash common.Hash, suffix []byte) error {
	key := GetDelegateKeyBySuffix(suffix)

	// todo test
	log.Debug("DelDelegateStoreBySuffix", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key))

	return db.del(blockHash, key)
}

func (db *StakingDB) DelUnDelegateCountStore(blockHash common.Hash, epoch uint64) error {
	count_key := GetUnDelegateCountKey(epoch)
	// todo test
	log.Debug("DelUnDelegateCountStore", "blockHash", blockHash.Hex(), "count_key", hex.EncodeToString(count_key))

	return db.del(blockHash, count_key)
}

func (db *StakingDB) DelUnDelegateItemStore(blockHash common.Hash, epoch, index uint64) error {
	item_key := GetUnDelegateItemKey(epoch, index)

	// todo test
	log.Debug("DelUnDelegateItemStore", "blockHash", blockHash.Hex(), "item_key", hex.EncodeToString(item_key))

	return db.del(blockHash, item_key)
}

func (db *StakingDB) SetEpochValIndex(blockHash common.Hash, indexArr ValArrIndexQueue) error {
	value, err := rlp.EncodeToBytes(indexArr)
	if nil != err {
		return err
	}

	// todo test
	log.Debug("SetEpochValIndex", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(GetEpochIndexKey()), "val", hex.EncodeToString(value))

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

	// todo test
	log.Debug("SetEpochValList", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(GetEpochValArrKey(start, end)),
		"val", hex.EncodeToString(value), "start", start, "end", end)

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

	// todo test
	log.Debug("DelEpochValListByBlockHash", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(GetEpochValArrKey(start, end)),
		"start", start, "end", end)

	return db.del(blockHash, GetEpochValArrKey(start, end))
}

func (db *StakingDB) SetRoundValIndex(blockHash common.Hash, indexArr ValArrIndexQueue) error {
	value, err := rlp.EncodeToBytes(indexArr)
	if nil != err {
		return err
	}

	// todo test
	log.Debug("SetRoundValIndex", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(GetRoundIndexKey()), "val", hex.EncodeToString(value))

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

	// todo test
	log.Debug("SetRoundValList", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(GetRoundValArrKey(start, end)),
		"val", hex.EncodeToString(value), "start", start, "end", end)

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

	// todo test
	log.Debug("DelRoundValListByBlockHash", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(GetRoundValArrKey(start, end)),
		"start", start, "end", end)

	return db.del(blockHash, GetRoundValArrKey(start, end))
}

func (db *StakingDB) IteratorCandidatePowerByBlockHash(blockHash common.Hash, ranges int) iterator.Iterator {
	return db.ranking(blockHash, CanPowerKeyPrefix, ranges)
}

//func (db *StakingDB) IteratorDelegateByIrrWithAddr (addr common.Address, ranges int) iterator.Iterator {
//	prefix := append(DelegateKeyPrefix, addr.Bytes()...)
//	return db.ranking(common.ZeroHash, prefix, ranges)
//}

func (db *StakingDB) IteratorDelegateByBlockHashWithAddr(blockHash common.Hash, addr common.Address, ranges int) iterator.Iterator {
	prefix := append(DelegateKeyPrefix, addr.Bytes()...)
	return db.ranking(blockHash, prefix, ranges)
}

// add the account staking Reference Count
func (db *StakingDB) AddAccountStakeRc(blockHash common.Hash, addr common.Address) error {
	key := GetAccountStakeRcKey(addr)
	val, err := db.get(blockHash, key)
	var v uint64
	switch {
	case nil != err && err != snapshotdb.ErrNotFound:
		return err
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	// todo test
	log.Debug("AddAccountStakeRc, query rc", "blockHash", blockHash.Hex(), "addr", addr.String(), "rc", v)

	v++

	// todo test
	log.Debug("AddAccountStakeRc", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key), "val", hex.EncodeToString(common.Uint64ToBytes(v)))

	return db.put(blockHash, key, common.Uint64ToBytes(v))
}

// sub the account staking Reference Count
func (db *StakingDB) SubAccountStakeRc(blockHash common.Hash, addr common.Address) error {
	key := GetAccountStakeRcKey(addr)
	val, err := db.get(blockHash, key)
	var v uint64
	switch {
	case nil != err && err != snapshotdb.ErrNotFound:
		return err
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	// Prevent large numbers from being directly called after the uint64 overflow
	if v == 0 {
		return nil
	}

	// todo test
	log.Debug("SubAccountStakeRc, query rc", "blockHash", blockHash.Hex(), "addr", addr.String(), "rc", v)

	v--

	if v == 0 {

		// todo test
		log.Debug("SubAccountStakeRc, del", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key))

		return db.del(blockHash, key)
	} else {

		// todo test
		log.Debug("SubAccountStakeRc, put", "blockHash", blockHash.Hex(), "key", hex.EncodeToString(key), "val", hex.EncodeToString(common.Uint64ToBytes(v)))

		return db.put(blockHash, key, common.Uint64ToBytes(v))
	}
}

// check the account staking Reference Count
func (db *StakingDB) HasAccountStakeRc(blockHash common.Hash, addr common.Address) (bool, error) {
	key := GetAccountStakeRcKey(addr)
	val, err := db.get(blockHash, key)
	var v uint64
	switch {
	case nil != err && err != snapshotdb.ErrNotFound:
		return false, err
	case nil != err && err == snapshotdb.ErrNotFound:
		return false, nil
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	// todo test
	log.Debug("HasAccountStakeRc, query rc", "blockHash", blockHash.Hex(), "addr", addr.String(), "key", hex.EncodeToString(key), "rc", v)

	if v == 0 {
		return false, nil
	} else if v > 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("Account Stake Reference Count cannot be negative, account: %s", addr.String())
	}
}

func (db *StakingDB) StoreRoundValidatorAddrs(blockHash common.Hash, key []byte, arry []common.Address) error {
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
