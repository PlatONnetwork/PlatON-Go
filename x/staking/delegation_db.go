package staking

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// about delegate ...

func (db *StakingDB) GetDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId enode.IDv0, stakeBlockNumber uint64) (*Delegation, error) {
	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	delByte, err := db.get(blockHash, key)
	if nil != err {
		return nil, err
	}

	del := new(DelegationForStorage)
	if err := rlp.DecodeBytes(delByte, del); nil != err {
		return nil, err
	}

	return (*Delegation)(del), nil
}

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

		tmp := new(DelegationForStorage)
		if err := rlp.DecodeBytes(itr.Value(), tmp); err != nil {
			return nil, err
		}
		info.Delegation = (*Delegation)(tmp)
		infos = append(infos, info)
	}
	return infos, nil
}

func (db *StakingDB) SetDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId enode.IDv0,
	stakeBlockNumber uint64, del *Delegation, isEinstein bool) error {
	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)
	if isEinstein {
		delByte, err := encodeStoredDelegateRLP(del)
		if nil != err {
			return err
		}

		return db.put(blockHash, key, delByte)
	}
	delByte, err := encodeV1StoredDelegateRLP(del)
	if nil != err {
		return err
	}

	return db.put(blockHash, key, delByte)
}

func (db *StakingDB) DelDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId enode.IDv0,
	stakeBlockNumber uint64) error {
	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	return db.del(blockHash, key)
}

func (db *StakingDB) IteratorDelegateByBlockHashWithAddr(blockHash common.Hash, addr common.Address, ranges int) iterator.Iterator {
	prefix := append(DelegateKeyPrefix, addr.Bytes()...)
	return db.ranking(blockHash, prefix, ranges)
}

func (db *StakingDB) GetDelegationLock(blockHash common.Hash, delAddr common.Address, currentEpoch uint32) (*DelegationLock, error) {
	key := GetDelegationLockKey(delAddr)
	delByte, err := db.get(blockHash, key)
	if nil != err {
		return nil, err
	}

	dell := new(DelegationLock)
	if err := rlp.DecodeBytes(delByte, dell); nil != err {
		return nil, err
	}
	dell.update(currentEpoch)
	return dell, nil

}

func (db *StakingDB) PutDelegationLock(blockHash common.Hash, delAddr common.Address, infos *DelegationLock) error {
	key := GetDelegationLockKey(delAddr)
	if infos.shouldDel() {
		return db.Del(blockHash, key)
	} else {
		delByte, err := rlp.EncodeToBytes(infos)
		if err != nil {
			return err
		}
		return db.put(blockHash, key, delByte)
	}
}
