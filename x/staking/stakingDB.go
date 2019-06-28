package staking

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"math/big"
)





type StakingDB struct {
	// todo the snapshotDB instance
	db snapshotdb.DB
}


func NewStakingDB (db snapshotdb.DB) *StakingDB{
	return &StakingDB {
		db: db,
	}
}

func (db *StakingDB) get (blockHash common.Hash, key []byte) ([]byte, error) {
	return db.db.Get(blockHash, key)
}

func (db *StakingDB) getFromCommitted (key []byte) ([]byte, error) {
	return db.db.GetFromCommittedBlock(key)
}

func (db *StakingDB) put (blockHash common.Hash, key, value []byte) error {
	return db.db.Put(blockHash, key, value)
}

func (db *StakingDB) del (blockHash common.Hash, key []byte) error {
	return db.db.Del(blockHash, key)
}

func (db *StakingDB) ranking (blockHash common.Hash, prefix []byte, ranges int) iterator.Iterator {
	return db.db.Ranking(blockHash, prefix, ranges)
}








func (db *StakingDB) GetCandidateStore (blockHash common.Hash, addr common.Address) (*xcom.Candidate, error) {
	key := xcom.CandidateKeyByAddr(addr)
	canByte, err := db.get(blockHash, key)

	if nil != err {
		return nil, err
	}
	var can xcom.Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCandidateStoreByIrr (addr common.Address) (*xcom.Candidate, error) {
	key := xcom.CandidateKeyByAddr(addr)
	canByte, err := db.getFromCommitted(key)

	if nil != err {
		return nil, err
	}
	var can xcom.Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}



func (db *StakingDB) GetCandidateStoreWithSuffix (blockHash common.Hash, suffix []byte) (*xcom.Candidate, error) {
	key := xcom.CandidateKeyBySuffix(suffix)
	canByte, err := db.get(blockHash, key)

	if nil != err {
		return nil, err
	}
	var can xcom.Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) GetCandidateStoreByIrrWithSuffix (suffix []byte) (*xcom.Candidate, error) {
	key := xcom.CandidateKeyBySuffix(suffix)
	canByte, err := db.getFromCommitted(key)

	if nil != err {
		return nil, err
	}
	var can xcom.Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (db *StakingDB) SetCandidateStore (blockHash common.Hash, addr common.Address, can *xcom.Candidate) error {

	key := xcom.CandidateKeyByAddr(addr)

	if val, err := rlp.EncodeToBytes(can); nil != err {
		return err
	}else {
		return db.put(blockHash, key, val)
	}
}

func (db *StakingDB) DelCandidateStore (blockHash common.Hash, addr common.Address) error {
	key := xcom.CandidateKeyByAddr(addr)
	return db.del(blockHash, key)
}

func (db *StakingDB) SetCanPowerStore (blockHash common.Hash, addr common.Address, can *xcom.Candidate) error {
	key := xcom.TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex)
	return db.put(blockHash, key, addr.Bytes())
}

func (db *StakingDB) DelCanPowerStore (blockHash common.Hash, can *xcom.Candidate) error {
	key := xcom.TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex)
	return db.del(blockHash, key)
}



func (db *StakingDB) AddUnStakeItemStore (blockHash common.Hash, epoch uint64, addr common.Address) error {


	count_key := xcom.GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return err
	}

	v := common.BytesToUint64(val)

	v++

	if err := db.put(blockHash, count_key, common.Uint64ToBytes(v)); nil != err {
		return err
	}
	item_key := xcom.GetUnStakeItemKey(epoch, v)

	return db.put(blockHash, item_key, addr.Bytes())
}

func (db *StakingDB) GetUnStakeCountStore (blockHash common.Hash, epoch uint64) (uint64, error) {
	count_key := xcom.GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return 0, err
	}
	return common.BytesToUint64(val), nil
}

func (db *StakingDB) GetUnStakeItemStore (blockHash common.Hash, epoch, index uint64) (common.Address, error) {
	item_key := xcom.GetUnStakeItemKey(epoch, index)
	addrByte, err := db.get(blockHash, item_key)
	if nil != err {
		return common.ZeroAddr, err
	}
	return common.BytesToAddress(addrByte), nil
}


func (db *StakingDB) DelUnStakeCountStore (blockHash common.Hash, epoch uint64) error {
	count_key := xcom.GetUnStakeCountKey(epoch)
	return db.del(blockHash, count_key)
}

func (db *StakingDB) DelUnStakeItemStore (blockHash common.Hash, epoch, index uint64) error {
	item_key := xcom.GetUnStakeItemKey(epoch, index)
	return db.del(blockHash, item_key)
}



func (db *StakingDB) GetDelegateStore (blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber uint64) (*xcom.Delegation, error) {
	key := xcom.GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	delByte, err := db.get(blockHash, key)
	if nil != err {
		return nil, err
	}

	var del xcom.Delegation
	if err := rlp.DecodeBytes(delByte, &del); nil != err {
		return nil, err
	}
	return &del, nil
}


func (db *StakingDB) GetDelegateStoreBySuffix (blockHash common.Hash, keySuffix[]byte) (*xcom.Delegation, error) {
	key := xcom.GetDelegateKeyBySuffix(keySuffix)
	delByte, err := db.get(blockHash, key)
	if nil != err {
		return nil, err
	}

	var del xcom.Delegation
	if err := rlp.DecodeBytes(delByte, &del); nil != err {
		return nil, err
	}
	return &del, nil
}

func (db *StakingDB) SetDelegateStore (blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	stakeBlockNumber uint64, del  *xcom.Delegation) error {
	key := xcom.GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	delByte, err := rlp.EncodeToBytes(del)
	if nil != err {
		return err
	}
	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) SetDelegateStoreBySuffix (blockHash common.Hash, suffix []byte, del *xcom.Delegation) error {
	key := xcom.GetDelegateKeyBySuffix(suffix)
	delByte, err := rlp.EncodeToBytes(del)
	if nil != err {
		return err
	}
	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) DelDelegateStore (blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	stakeBlockNumber uint64) error {
	key := xcom.GetDelegateKey(delAddr, nodeId, stakeBlockNumber)
	return db.del(blockHash, key)
}

func (db *StakingDB) DelDelegateStoreBySuffix (blockHash common.Hash, suffix []byte) error {
	key := xcom.GetDelegateKeyBySuffix(suffix)
	return db.del(blockHash, key)
}

func (db *StakingDB) AddUnDelegateItemStore (blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	epoch, stakeBlockNumber uint64, amount *big.Int) error {


	count_key := xcom.GetUnDelegateCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return err
	}

	v := common.BytesToUint64(val)
	v++


	if err := db.put(blockHash, count_key, common.Uint64ToBytes(v)); nil != err {
		return err
	}
	item_key := xcom.GetUnDelegateItemKey(epoch, v)

	suffix :=  append(delAddr.Bytes(), append(nodeId.Bytes(), common.Uint64ToBytes(stakeBlockNumber)...)...)

	unDelegateItem := &xcom.UnDelegateItem{
		KeySuffix: 	suffix,
		Amount: 	amount,
	}

	item, err := rlp.EncodeToBytes(unDelegateItem)
	if nil != err {
		return err
	}
	return db.put(blockHash, item_key, item)
}

func (db *StakingDB) GetUnDelegateCountStore (blockHash common.Hash, epoch uint64) (uint64, error) {

	count_key := xcom.GetUnDelegateCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return 0, err
	}

	return common.BytesToUint64(val), nil
}

func (db *StakingDB) GetUnDelegateItemStore (blockHash common.Hash, epoch, index uint64) (*xcom.UnDelegateItem, error) {

	item_key := xcom.GetUnDelegateItemKey(epoch, index)

	itemByte, err := db.get(blockHash, item_key)
	if nil != err {
		return nil, err
	}

	var unDelegateItem xcom.UnDelegateItem
	if err := rlp.DecodeBytes(itemByte, &unDelegateItem); nil != err {
		return nil, err
	}
	return &unDelegateItem, nil
}


func (db *StakingDB) GetVerifierListByIrr () (*xcom.Validator_array, error) {

	arrByte, err := db.getFromCommitted(xcom.GetEpochValidatorKey())
	if nil != err {
		return nil, err
	}

	var arr *xcom.Validator_array
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) SetVerfierList (blockHash common.Hash, val_Arr *xcom.Validator_array) error {

	value, err := rlp.EncodeToBytes(val_Arr)
	if nil != err {
		return err
	}
	return db.put(blockHash, xcom.GetEpochValidatorKey(), value)
}

func (db *StakingDB) GetVerifierListByBlockHash (blockHash common.Hash) (*xcom.Validator_array, error) {

	arrByte, err := db.get(blockHash, xcom.GetEpochValidatorKey())
	if nil != err {
		return nil, err
	}

	var arr *xcom.Validator_array
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}


func (db *StakingDB) GetPreviousValidatorListByIrr () (*xcom.Validator_array, error) {
	arrByte, err := db.getFromCommitted(xcom.GetPreRoundValidatorKey())
	if nil != err {
		return nil, err
	}

	var arr *xcom.Validator_array
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) GetPreviousValidatorListByBlockHash (blockHash common.Hash) (*xcom.Validator_array, error) {
	arrByte, err := db.get(blockHash, xcom.GetPreRoundValidatorKey())
	if nil != err {
		return nil, err
	}

	var arr *xcom.Validator_array
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) GetCurrentValidatorListByIrr () (*xcom.Validator_array, error) {
	arrByte, err := db.getFromCommitted(xcom.GetCurRoundValidatorKey())
	if nil != err {
		return nil, err
	}

	var arr *xcom.Validator_array
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) GetCurrentValidatorListByBlockHash (blockHash common.Hash) (*xcom.Validator_array, error) {
	arrByte, err := db.get(blockHash, xcom.GetCurRoundValidatorKey())
	if nil != err {
		return nil, err
	}

	var arr *xcom.Validator_array
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) GetNextValidatorListByIrr () (*xcom.Validator_array, error) {
	arrByte, err := db.getFromCommitted(xcom.GetNextRoundValidatorKey())
	if nil != err {
		return nil, err
	}

	var arr *xcom.Validator_array
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) GetNextValidatorListByBlockHash (blockHash common.Hash) (*xcom.Validator_array, error) {
	arrByte, err := db.get(blockHash, xcom.GetNextRoundValidatorKey())
	if nil != err {
		return nil, err
	}

	var arr *xcom.Validator_array
	if err := rlp.DecodeBytes(arrByte, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (db *StakingDB) IteratorCandidatePowerByBlockHash (blockHash common.Hash, ranges int) iterator.Iterator {
	return db.ranking(blockHash, xcom.CanPowerKeyPrefix, ranges)
}

func (db *StakingDB) IteratorCandidatePowerByIrr (ranges int) iterator.Iterator {
	return db.ranking(common.ZeroHash, xcom.CanPowerKeyPrefix, ranges)
}

func (db *StakingDB) IteratorDelegateByBlockHashWithAddr (blockHash common.Hash, addr common.Address, ranges int) iterator.Iterator {
	prefix := append(xcom.DelegateKeyPrefix, addr.Bytes()...)
	return db.ranking(blockHash, prefix, ranges)
}

func (db *StakingDB) IteratorDelegateByIrrWithAddr (addr common.Address, ranges int) iterator.Iterator {
	prefix := append(xcom.DelegateKeyPrefix, addr.Bytes()...)
	return db.ranking(common.ZeroHash, prefix, ranges)
}