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










func (db *StakingDB) getCandidate (blockHash common.Hash, nodeId discover.NodeID) ([]byte, error) {
	key, err := xcom.CandidateKeyByNodeId(nodeId)
	if nil != err {
		return nil, err
	}
	return db.get(blockHash, key)
}


func (db *StakingDB) setCandidate2DB(blockHash common.Hash, addr common.Address, can *xcom.Candidate) error {

	key := xcom.CandidateKeyByAddr(addr)

	if val, err := rlp.EncodeToBytes(can); nil != err {
		return err
	}else {
		return db.put(blockHash, key, val)
	}
}

func (db *StakingDB) setCanPower2DB (blockHash common.Hash, addr common.Address, can *xcom.Candidate) error {
	key := xcom.TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex)
	return db.put(blockHash, key, addr.Bytes())
}

func (db *StakingDB) delCanPower2DB (blockHash common.Hash, can *xcom.Candidate) error {
	key := xcom.TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex)
	return db.del(blockHash, key)
}


func (db *StakingDB) increaseUnStakeCount (blockHash common.Hash, epoch uint64) error {

	key := xcom.GetUnStakeCountKey(epoch)

	val, err := db.get(blockHash, key)
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
	return db.put(blockHash, key, valNew)
}