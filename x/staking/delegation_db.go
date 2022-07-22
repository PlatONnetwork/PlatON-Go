package staking

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

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

func (db *StakingDB) DelDelegateStore(blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	stakeBlockNumber uint64) error {
	key := GetDelegateKey(delAddr, nodeId, stakeBlockNumber)

	return db.del(blockHash, key)
}

func (db *StakingDB) IteratorDelegateByBlockHashWithAddr(blockHash common.Hash, addr common.Address, ranges int) iterator.Iterator {
	prefix := append(DelegateKeyPrefix, addr.Bytes()...)
	return db.ranking(blockHash, prefix, ranges)
}
