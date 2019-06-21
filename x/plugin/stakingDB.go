package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
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










func (db *StakingDB) getCandidateStore (blockHash common.Hash, nodeId discover.NodeID) ([]byte, error) {
	key, err := xcom.CandidateKeyByNodeId(nodeId)
	if nil != err {
		return nil, err
	}
	return db.get(blockHash, key)
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
	key := xcom.TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex)
	return db.put(blockHash, key, addr.Bytes())
}

func (db *StakingDB) delCanPowerStore (blockHash common.Hash, can *xcom.Candidate) error {
	key := xcom.TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex)
	return db.del(blockHash, key)
}



func (db *StakingDB) addUnStakeItemStore (blockHash common.Hash, epoch uint64, addr common.Address) error {


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
	item_key := xcom.GetUnStakeItemKey(epoch, uint64(v))
	return db.put(blockHash, item_key, addr.Bytes())
}

func (db *StakingDB) getUnStakeCountStore (blockHash common.Hash, epoch uint64) (int, error) {


	return 0, nil
}

func (db *StakingDB) getUnStakeItemStore (blockHash common.Hash, epoch uint64) (common.Address, error) {


	return common.Address{}, nil
}