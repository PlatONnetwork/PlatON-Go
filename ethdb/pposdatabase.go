package ethdb

import (
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	//"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

//type PPosDatabase struct {
//	db *leveldb.DB // LevelDB instance
//}



func NewPPosDatabase (file string) (*LDBDatabase, error)  {


	logger := log.New("database", file)

	logger.Info("Allocated cache and file handles")

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(file, &opt.Options{
		/*//OpenFilesCacheCapacity: -1,	// Use -1 for zero, this has same effect as specifying NoCacher to OpenFilesCacher.
		//BlockCacheCapacity:     -1, // Use -1 for zero, this has same effect as specifying NoCacher to BlockCacher.
		DisableBlockCache:		true,
		DisableBufferPool: 		true,
		//NoSync:		true,
		////CompactionL0Trigger: 	0,
		//DisableBufferPool:		true,
		//DisableLargeBatchTransaction: true,
		//NoWriteMerge: 		true,
		//CompactionTotalSizeMultiplier: 0,
		//WriteBuffer:	0,
		BlockSize: 0 * opt.MiB,
		BlockCacheCapacity: 0 * opt.MiB,
		//CompactionTotalSizeMultiplier: 0,
		Filter:                 filter.NewBloomFilter(10),*/



		/*DisableBlockCache:      true,
		BlockRestartInterval:   5,
		BlockSize:              80,
		Compression:            opt.NoCompression,
		OpenFilesCacheCapacity: -1,
		Strict:                 opt.StrictAll,
		WriteBuffer:            1000,
		CompactionTableSize:    2000,*/


		DisableLargeBatchTransaction: true,
		Compression:                  opt.NoCompression,
		CompactionTableSize:          1 * opt.MiB,
		WriteBuffer:                  1 * opt.MiB,
	})

	//db, err := leveldb.OpenFile(file,nil)

	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}
	// (Re)check for errors and abort if opening of the db failed
	if err != nil {
		return nil, err
	}
	return &LDBDatabase{
		fn:  file,
		db:  db,
		log: logger,
	}, nil
}


