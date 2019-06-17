package plugin


type StakingDB struct {
	// todo the snapshotDB instance
	db interface{}
}


func NewStakingDB (db interface{}) *StakingDB{
	return &StakingDB {
		db: db,
	}
}