package plugin

import "github.com/PlatONnetwork/PlatON-Go/common"

type StakingDB struct {
	// todo the snapshotDB instance
	db interface{}
}


func NewStakingDB (db interface{}) *StakingDB{
	return &StakingDB {
		db: db,
	}
}

func (db *StakingDB) Get (blockHash common.Hash, key []byte) ([]byte, error) {


	return nil, nil
}