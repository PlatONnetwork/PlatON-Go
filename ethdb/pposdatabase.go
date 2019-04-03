package ethdb

import (
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)


func NewPPosDatabase (file string) (*LDBDatabase, error)  {


	logger := log.New("database", file)

	logger.Info("Allocated cache and file handles")

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(file, &opt.Options{

		//  100 block  45.1 M  45.2 M; 200 block 45.3 M  45.4 M  45.1 M
		OpenFilesCacheCapacity: 1000,
		BlockCacheCapacity:     500,
		BlockSize:			1,
		WriteBuffer:		64,

		CompactionSourceLimitFactor: 1000,
		Filter:                 filter.NewBloomFilter(10),

		//Compression:  opt.SnappyCompression,

		/*// 100 block 65.9 M 88.1 M； 200 block 73.5 M 88.2 M
		DisableBlockCache:      true,
		BlockRestartInterval:   5,
		BlockSize:              80,
		Compression:            opt.NoCompression,
		OpenFilesCacheCapacity: -1,
		Strict:                 opt.StrictAll,
		WriteBuffer:            1000,
		CompactionTableSize:    2000,
		Filter:                 filter.NewBloomFilter(10),*/

		/*// myself 100 block 42.9  M 45.7 45.8 M ; 200 block 45.8M  45.9 M
		CompactionSourceLimitFactor: 	1 * opt.MiB,
		DisableLargeBatchTransaction: 	true,
		CompactionTableSize:          	1 * opt.MiB,
		WriteBuffer:                  	10 * opt.MiB,
		Strict: 						opt.StrictCompaction,
		Filter:                 		filter.NewBloomFilter(10),*/


		/*// ethdb config 100 block 356 M 615 M ； 200 block 586
		OpenFilesCacheCapacity: 1024,
		BlockCacheCapacity:     768 / 2 * opt.MiB,
		WriteBuffer:            1024 / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),*/
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


