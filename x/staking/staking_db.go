package staking

import (
	"math/big"

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
	return db.db.Put(blockHash, key, value)
}

func (db *StakingDB) del(blockHash common.Hash, key []byte) error {
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
		return db.put(blockHash, key, val)
	}
}

func (db *StakingDB) DelCandidateStore(blockHash common.Hash, addr common.Address) error {
	key := CandidateKeyByAddr(addr)
	return db.del(blockHash, key)
}

func (db *StakingDB) SetCanPowerStore(blockHash common.Hash, addr common.Address, can *Candidate) error {
	key := TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex, can.ProgramVersion)
	return db.put(blockHash, key, addr.Bytes())
}

func (db *StakingDB) DelCanPowerStore(blockHash common.Hash, can *Candidate) error {
	key := TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex, can.ProgramVersion)
	return db.del(blockHash, key)
}

func (db *StakingDB) AddUnStakeItemStore(blockHash common.Hash, epoch uint64, addr common.Address) error {

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

	if err := db.put(blockHash, count_key, common.Uint64ToBytes(v)); nil != err {
		return err
	}
	item_key := GetUnStakeItemKey(epoch, v)

	return db.put(blockHash, item_key, addr.Bytes())
}

func (db *StakingDB) GetUnStakeCountStore(blockHash common.Hash, epoch uint64) (uint64, error) {
	count_key := GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return 0, err
	}
	return common.BytesToUint64(val), nil
}

func (db *StakingDB) GetUnStakeItemStore(blockHash common.Hash, epoch, index uint64) (common.Address, error) {
	item_key := GetUnStakeItemKey(epoch, index)
	addrByte, err := db.get(blockHash, item_key)
	if nil != err {
		return common.ZeroAddr, err
	}
	return common.BytesToAddress(addrByte), nil
}

func (db *StakingDB) DelUnStakeCountStore(blockHash common.Hash, epoch uint64) error {
	count_key := GetUnStakeCountKey(epoch)
	return db.del(blockHash, count_key)
}

func (db *StakingDB) DelUnStakeItemStore(blockHash common.Hash, epoch, index uint64) error {
	item_key := GetUnStakeItemKey(epoch, index)
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

func (db *StakingDB) AddUnDelegateItemStore(blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	epoch, stakeBlockNumber uint64, amount *big.Int) error {

	count_key := GetUnDelegateCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	var v uint64
	switch {
	case nil != err && err != snapshotdb.ErrNotFound:
		return err
	case nil == err && len(val) != 0:
		v = common.BytesToUint64(val)
	}

	v++

	if err := db.put(blockHash, count_key, common.Uint64ToBytes(v)); nil != err {
		return err
	}
	item_key := GetUnDelegateItemKey(epoch, v)

	suffix := append(delAddr.Bytes(), append(nodeId.Bytes(), common.Uint64ToBytes(stakeBlockNumber)...)...)

	unDelegateItem := &UnDelegateItem{
		KeySuffix: suffix,
		Amount:    amount,
	}

	item, err := rlp.EncodeToBytes(unDelegateItem)
	if nil != err {
		return err
	}
	return db.put(blockHash, item_key, item)
}

func (db *StakingDB) GetUnDelegateCountStore(blockHash common.Hash, epoch uint64) (uint64, error) {

	count_key := GetUnDelegateCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return 0, err
	}

	return common.BytesToUint64(val), nil
}

func (db *StakingDB) GetUnDelegateItemStore(blockHash common.Hash, epoch, index uint64) (*UnDelegateItem, error) {

	item_key := GetUnDelegateItemKey(epoch, index)

	itemByte, err := db.get(blockHash, item_key)
	if nil != err {
		return nil, err
	}

	var unDelegateItem UnDelegateItem
	if err := rlp.DecodeBytes(itemByte, &unDelegateItem); nil != err {
		return nil, err
	}
	return &unDelegateItem, nil
}

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

func (db *StakingDB) SetEpochValList(blockHash common.Hash, start, end uint64, val_Arr ValidatorQueue) error {

	value, err := rlp.EncodeToBytes(val_Arr)
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

func (db *StakingDB) SetRoundValList(blockHash common.Hash, start, end uint64, val_Arr ValidatorQueue) error {

	value, err := rlp.EncodeToBytes(val_Arr)
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
