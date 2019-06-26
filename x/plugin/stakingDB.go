package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	"strconv"
)





type StakingDB struct {
	// todo the snapshotDB instance
	db xcom.SnapshotDB
}


func NewStakingDB (db xcom.SnapshotDB) *StakingDB{
	return &StakingDB {
		db: db,
	}
}

func (db *StakingDB) get (blockHash common.Hash, key []byte) ([]byte, error) {
	return db.db.Get(blockHash, key)
}

func (db *StakingDB) put (blockHash common.Hash, key, value []byte) error {
	return db.db.Put(blockHash, key, value)
}

func (db *StakingDB) del (blockHash common.Hash, key []byte) error {
	return db.db.Del(blockHash, key)
}










func (db *StakingDB) getCandidateStore (blockHash common.Hash, addr common.Address) (*xcom.Candidate, error) {
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


func (db *StakingDB) setCandidateStore (blockHash common.Hash, addr common.Address, can *xcom.Candidate) error {

	key := xcom.CandidateKeyByAddr(addr)

	if val, err := rlp.EncodeToBytes(can); nil != err {
		return err
	}else {
		return db.put(blockHash, key, val)
	}
}

func (db *StakingDB) delCandidateStore (blockHash common.Hash, addr common.Address) error {
	key := xcom.CandidateKeyByAddr(addr)
	return db.del(blockHash, key)
}

func (db *StakingDB) setCanPowerStore (blockHash common.Hash, addr common.Address, can *xcom.Candidate) error {
	key := xcom.TallyPowerKey(can.Shares, int(can.StakingBlockNum), int(can.StakingTxIndex))
	return db.put(blockHash, key, addr.Bytes())
}

func (db *StakingDB) delCanPowerStore (blockHash common.Hash, can *xcom.Candidate) error {
	key := xcom.TallyPowerKey(can.Shares, int(can.StakingBlockNum), int(can.StakingTxIndex))
	return db.del(blockHash, key)
}



func (db *StakingDB) addUnStakeItemStore (blockHash common.Hash, epoch int, addr common.Address) error {


	count_key := xcom.GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return err
	}

	var v int
	if len(val) != 0 {
		v, err = strconv.Atoi(string(val))
		if nil != err {
			return err
		}
	}

	v++

	valNew := []byte(strconv.Itoa(v))
	if err := db.put(blockHash, count_key, valNew); nil != err {
		return err
	}
	item_key := xcom.GetUnStakeItemKey(epoch, v)

	return db.put(blockHash, item_key, addr.Bytes())
}

func (db *StakingDB) getUnStakeCountStore (blockHash common.Hash, epoch int) (int, error) {
	count_key := xcom.GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return 0, err
	}

	var v int
	if len(val) != 0 {
		v, err = strconv.Atoi(string(val))
		if nil != err {
			return 0, err
		}
	}

	return v, nil
}

func (db *StakingDB) getUnStakeItemStore (blockHash common.Hash, epoch, index int) (common.Address, error) {
	item_key := xcom.GetUnStakeItemKey(epoch, index)
	addrByte, err := db.get(blockHash, item_key)
	if nil != err {
		return common.ZeroAddr, err
	}
	return common.BytesToAddress(addrByte), nil
}


func (db *StakingDB) delUnStakeCountStore (blockHash common.Hash, epoch int) error {
	count_key := xcom.GetUnStakeCountKey(epoch)
	return db.del(blockHash, count_key)
}

func (db *StakingDB) delUnStakeItemStore (blockHash common.Hash, epoch, index int) error {
	item_key := xcom.GetUnStakeItemKey(epoch, index)
	return db.del(blockHash, item_key)
}



func (db *StakingDB) getDelegateStore (blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber int) (*xcom.Delegation, error) {
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


func (db *StakingDB) getDelegateStoreBySuffix (blockHash common.Hash, keySuffix[]byte) (*xcom.Delegation, error) {
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

func (db *StakingDB) setDelegateStore (blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber int, del  *xcom.Delegation) error {
	key := xcom.GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	delByte, err := rlp.EncodeToBytes(del)
	if nil != err {
		return err
	}
	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) setDelegateStoreBySuffix (blockHash common.Hash, suffix []byte, del *xcom.Delegation) error {
	key := xcom.GetDelegateKeyBySuffix(suffix)
	delByte, err := rlp.EncodeToBytes(del)
	if nil != err {
		return err
	}
	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) delDelegateStore (blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber int) error {
	key := xcom.GetDelegateKey(delAddr, nodeId, stakeBlockNumber)
	return db.del(blockHash, key)
}

func (db *StakingDB) delDelegateStoreBySuffix (blockHash common.Hash, suffix []byte) error {
	key := xcom.GetDelegateKeyBySuffix(suffix)
	return db.del(blockHash, key)
}

func (db *StakingDB) addUnDelegateItemStore (blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	epoch, stakeBlockNumber int, amount *big.Int) error {


	count_key := xcom.GetUnDelegateCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return err
	}

	var v int
	if len(val) != 0 {
		v, err = strconv.Atoi(string(val))
		if nil != err {
			return err
		}
	}

	v++

	valNew := []byte(strconv.Itoa(v))
	if err := db.put(blockHash, count_key, valNew); nil != err {
		return err
	}
	item_key := xcom.GetUnDelegateItemKey(epoch, v)

	num := strconv.Itoa(stakeBlockNumber)

	suffix :=  append(delAddr.Bytes(), append(nodeId.Bytes(), []byte(num)...)...)

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

func (db *StakingDB) getUnDelegateCountStore (blockHash common.Hash, epoch int) (int, error) {

	count_key := xcom.GetUnDelegateCountKey(epoch)

	val, err := db.get(blockHash, count_key)
	if nil != err {
		return 0, err
	}

	var v int
	if len(val) != 0 {
		v, err = strconv.Atoi(string(val))
		if nil != err {
			return 0, err
		}
	}
	return v, nil
}

func (db *StakingDB) getUnDelegateItemStore (blockHash common.Hash, epoch, index int) (*xcom.UnDelegateItem, error) {

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