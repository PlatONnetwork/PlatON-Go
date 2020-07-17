// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package snapshotdb

import (
	"container/heap"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/metrics"

	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/robfig/cron"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbError "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

const (
	kvLIMIT                   = 2000
	JournalRemain             = 200
	UnBlockNeedClean          = 200
	MaxBlockCompaction        = 10
	MaxBlockCompactionSync    = 100
	MaxBlockTriggerCompaction = 200
)

//DB the main snapshotdb interface
//  example
//  new a recognized blockData(sync from other peer)
//  dbInstance.NewBlock(blockNumber, parentHash, hash)
//  dbInstance.Put(hash, kv.key, kv.value)
//  dbInstance.Commit(hash)
//
//  new a unrecognized blockData(a block produce by self)
//  dbInstance.NewBlock(blockNumber, parentHash, common.ZeroHash)
//  dbInstance.Put(hash, kv.key, kv.value)
//  dbInstance.Flush(hash common.Hash, blockNumber *big.Int)
//  dbInstance.Commit(hash)
//  get a  blockData with hash
//  dbInstance.Get(hash, key)
//  get a  blockData without hash
//  dbInstance.Get(common.zerohash, key)
type DB interface {
	Put(hash common.Hash, key, value []byte) error
	NewBlock(blockNumber *big.Int, parentHash common.Hash, hash common.Hash) error
	Get(hash common.Hash, key []byte) ([]byte, error)
	GetFromCommittedBlock(key []byte) ([]byte, error)
	Del(hash common.Hash, key []byte) error
	Has(hash common.Hash, key []byte) (bool, error)
	Flush(hash common.Hash, blocknumber *big.Int) error
	Ranking(hash common.Hash, key []byte, ranges int) iterator.Iterator
	//notice , iter.key or iter.value is slice，if you want to save it to a slice,you can use copy
	// container:=make([]byte,0)
	// for iter.next{
	// 	tosave:= make([]byte,len(iter.value))
	//  copy(tosave,iter.value)
	// 	container = append(container,tosave)
	// }
	//
	WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error
	Commit(hash common.Hash) error

	// Clear close db , remove all db file
	Clear() error

	BaseDB

	GetLastKVHash(blockHash common.Hash) []byte
	BaseNum() (*big.Int, error)
	Close() error
	Compaction() error
	SetEmpty() error

	//ues to Revert failed tx
	RevertToSnapshot(hash common.Hash, revid int)
	Snapshot(hash common.Hash) int
}

type BaseDB interface {
	PutBaseDB(key, value []byte) error
	GetBaseDB(key []byte) ([]byte, error)
	DelBaseDB(key []byte) error
	// WriteBaseDB apply the given [][2][]byte to the baseDB.
	WriteBaseDB(kvs [][2][]byte) error
	//SetCurrent use for fast sync
	SetCurrent(highestHash common.Hash, base, height big.Int) error
	GetCurrent() *current
}

var (
	dbpath string

	blockchain Chain

	dbInstance *snapshotDB

	instance sync.Mutex

	baseDBcache   int
	baseDBhandles int

	logger = log.Root().New("package", "snapshotdb")

	//ErrNotFound when db not found
	ErrNotFound = errors.New("snapshotDB: not found")

	ErrBlockTooLow = errors.New("the block is less than commit highest block")
)

type snapshotDB struct {
	path string

	snapshotLockC int32

	current *current

	baseDB *leveldb.DB

	unCommit *unCommitBlocks

	committed  []*blockData
	commitLock sync.RWMutex

	walCh     chan *blockData
	walExitCh chan struct{}
	walSync   sync.WaitGroup

	corn *cron.Cron

	closed bool

	dbError error
}

type Chain interface {
	CurrentHeader() *types.Header
	GetHeaderByHash(hash common.Hash) *types.Header
	GetHeaderByNumber(number uint64) *types.Header
}

func SetDBPathWithNode(path string) {
	dbpath = path
	logger.Info("set path", "path", dbpath)
}

func SetDBBlockChain(c Chain) {
	blockchain = c
	logger.Info("set blockchain")
}

func GetDBBlockChain() Chain {
	return blockchain
}

func SetDBOptions(cache int, handles int) {
	baseDBcache = cache
	baseDBhandles = handles
}

//Instance return the Instance of the db
func Instance() DB {
	instance.Lock()
	defer instance.Unlock()
	if dbInstance == nil || dbInstance.closed {
		logger.Debug("dbInstance is nil", "path", dbpath)
		if dbInstance == nil {
			dbInstance = new(snapshotDB)
		}
		if err := initDB(dbpath, dbInstance); err != nil {
			logger.Error("init db fail", "err", err)
			panic(err)
		}
	}
	return dbInstance
}

func openBaseDB(snapshotDBPath string, cache int, handles int) (*leveldb.DB, error) {
	leveldbPath := getBaseDBPath(snapshotDBPath)
	baseDB, err := leveldb.OpenFile(leveldbPath, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if err != nil {
		if _, corrupted := err.(*leveldbError.ErrCorrupted); corrupted {
			baseDB, err = leveldb.RecoverFile(leveldbPath, nil)
			if err != nil {
				return nil, fmt.Errorf("[SnapshotDB.recover]RecoverFile baseDB fail:%v", err)
			}
		} else {
			return nil, err
		}
	}
	return baseDB, nil
}

func open(path string, cache int, handles int, baseOnly bool) (*snapshotDB, error) {
	logger.Info("open snapshot db Allocated cache and file handles", "cache", cache, "handles", handles, "baseDB", baseOnly)

	baseDB, err := openBaseDB(path, cache, handles)
	if err != nil {
		return nil, err
	}

	unCommitBlock := new(unCommitBlocks)
	unCommitBlock.blocks = make(map[common.Hash]*blockData)
	db := &snapshotDB{
		path:          path,
		unCommit:      unCommitBlock,
		committed:     make([]*blockData, 0),
		baseDB:        baseDB,
		snapshotLockC: snapshotUnLock,
		walCh:         make(chan *blockData, 2),
		walExitCh:     make(chan struct{}),
	}
	if baseOnly {
		return db, nil
	}

	_, getCurrentError := baseDB.Get([]byte(CurrentSet), nil)
	if getCurrentError == nil {
		logger.Info("begin recover", "path", path)
		if err := db.loadCurrent(); err != nil {
			return nil, err
		}
		logger.Info("load current", "base", db.current.base, "high", db.current.highest)
		if err := db.recover(); err != nil {
			logger.Error("recover db fail:", "error", err)
			return nil, err
		}
	} else if getCurrentError == leveldb.ErrNotFound {
		logger.Info("begin init db current", "path", path)
		if err := db.SetCurrent(common.ZeroHash, *common.Big0, *common.Big0); err != nil {
			return nil, err
		}
	} else {
		return nil, getCurrentError
	}
	return db, nil
}

func Open(path string, cache int, handles int, baseOnly bool) (DB, error) {
	db, err := open(path, cache, handles, baseOnly)
	if err != nil {
		return nil, err
	}
	if !baseOnly {
		if err := db.Start(); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func copyDB(from, to *snapshotDB) {
	to.path = from.path
	to.current = from.current
	to.baseDB = from.baseDB
	to.unCommit = from.unCommit
	to.committed = from.committed
	to.corn = from.corn
	to.closed = from.closed
	to.snapshotLockC = from.snapshotLockC
	to.walExitCh = from.walExitCh
	to.walCh = from.walCh
}

func initDB(path string, sdb *snapshotDB) error {
	dbInterface, err := open(path, baseDBcache, baseDBhandles, false)
	if err != nil {
		return err
	}
	copyDB(dbInterface, sdb)
	if err := sdb.Start(); err != nil {
		return err
	}
	return nil
}

func (s *snapshotDB) Start() error {
	if err := s.cornStart(); err != nil {
		return err
	}
	go s.loopWriteWal()
	return nil
}

func (s *snapshotDB) cornStart() error {
	s.corn = cron.New()
	if err := s.corn.AddFunc("@every 1s", s.schedule); err != nil {
		logger.Error("set db corn compaction fail", "err", err)
		return err
	}
	if metrics.Enabled {
		if err := s.corn.AddFunc("@every 3s", s.metrics); err != nil {
			logger.Error("set db corn metrics fail", "err", err)
			return err
		}
	}
	s.corn.Start()
	return nil
}

func (s *snapshotDB) WriteBaseDB(kvs [][2][]byte) error {
	batch := new(leveldb.Batch)
	for _, value := range kvs {
		batch.Put(value[0], value[1])
	}
	if err := s.baseDB.Write(batch, nil); err != nil {
		return err
	}
	return nil
}

func (s *snapshotDB) SetCurrent(highestHash common.Hash, base, height big.Int) error {
	current := newCurrent(&height, &base, highestHash)
	if err := current.saveCurrentToBaseDB(CurrentAll, s.baseDB, true); err != nil {
		return err
	}
	s.current = current
	logger.Debug("SetCurrent", "base", s.current.base, "height", s.current.highest)
	return nil
}

func (s *snapshotDB) GetCurrent() *current {
	return s.current
}

// GetCommittedBlock    get value from committed blockdata > baseDB
func (s *snapshotDB) GetFromCommittedBlock(key []byte) ([]byte, error) {
	v, err := s.getFromCommit(key)
	if err == nil {
		if v == nil || len(v) == 0 {
			return nil, ErrNotFound
		}
		return v, nil
	}
	return s.GetBaseDB(key)
}

func (s *snapshotDB) SetEmpty() error {
	logger.Debug("set snapshotDB empty", "path", s.path)
	instance.Lock()
	defer instance.Unlock()
	if err := s.Clear(); err != nil {
		return err
	}
	return initDB(s.path, s)
}

func (s *snapshotDB) PutBaseDB(key, value []byte) error {
	err := s.baseDB.Put(key, value, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *snapshotDB) DelBaseDB(key []byte) error {
	err := s.baseDB.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetLastKVHash return the last kv hash
// if hash is nil ,get unRecognized block lastkv hash,
// else, get recognized block lastkv  hash
func (s *snapshotDB) GetLastKVHash(blockHash common.Hash) []byte {
	block := s.unCommit.Get(blockHash)
	if block == nil {
		return nil
	}
	return block.kvHash.Bytes()
}

// Del del key,val from  snapshotDB
// if hash is nil, unRecognizedBlockData > recognizedBlockData
// if hash is not nil,it will del in recognized BlockData
func (s *snapshotDB) Del(hash common.Hash, key []byte) error {
	if err := s.put(hash, key, nil); err != nil {
		return fmt.Errorf("[SnapshotDB]del fail:%v", err)
	}
	return nil
}

// Compaction ,write commit to baseDB,and then removeJournal lessThan BaseNum
// it will write to baseDB
// case kv>2000 and block == 1
// case kv<2000,block... <9
// case kv<2000,block...=9
func (s *snapshotDB) Compaction() error {
	if len(s.committed) == 0 {
		return nil
	}
	commitNum := s.findToWrite()
	if commitNum == 0 {
		return nil
	}
	if err := s.writeToBasedb(commitNum); err != nil {
		return err
	}
	s.commitLock.Lock()
	s.committed = s.committed[commitNum:]
	if err := s.current.increaseBase(uint64(commitNum), s.baseDB); err != nil {
		s.commitLock.Unlock()
		logger.Error("save base to current fail", "err", err)
		return err
	}
	s.commitLock.Unlock()
	//delete block no use  in unCommit
	s.rmExpireForkBlock()
	return nil
}

func (s *snapshotDB) rmExpireForkBlock() {
	s.unCommit.Lock()
	// if unCommit blocks length > 200,rm block below BaseNum
	if len(s.unCommit.blocks) > UnBlockNeedClean {
		currentBase := s.current.GetBase(false).Num
		for key, value := range s.unCommit.blocks {
			if currentBase.Cmp(value.Number) >= 0 {
				delete(s.unCommit.blocks, key)
				logger.Debug("compaction delete no need blocks", "num", value.Number, "hash", value.BlockHash.String())
			}
		}
	}
	s.unCommit.Unlock()
}

func (s *snapshotDB) findToWrite() int {
	s.commitLock.RLock()
	defer s.commitLock.RUnlock()
	var (
		kvsize    int
		commitNum int
	)
	if len(s.committed) > MaxBlockTriggerCompaction {
		commitNum = MaxBlockCompactionSync
		return commitNum
	} else {
		for i := 0; i < len(s.committed); i++ {
			if i < MaxBlockCompaction {
				if kvsize > kvLIMIT {
					commitNum = i - 1
					break
				}
				kvsize += s.committed[i].data.Len()
			} else {
				commitNum = i
				break
			}
		}
	}
	if commitNum == 0 {
		commitNum++
	}

	//the commit block can't grate than highest from baseDB and blockchain.Current
	highest, err := s.current.GetHighestFromDB(s.baseDB)
	if err != nil {
		s.dbError = err
		return 0
	}
	minimumHeight := new(big.Int).Set(highest.Num)
	if blockchain != nil {
		header := blockchain.CurrentHeader()
		if minimumHeight.Cmp(header.Number) > 0 {
			minimumHeight = new(big.Int).Set(header.Number)
		}
	}

	var length = commitNum
	for i := 0; i < length; i++ {
		if s.committed[i].Number.Cmp(minimumHeight) > 0 {
			commitNum--
		}
	}
	return commitNum
}

func (s *snapshotDB) writeToBasedb(commitNum int) error {
	batch := new(leveldb.Batch)
	for i := 0; i < commitNum; i++ {
		itr := s.committed[i].data.NewIterator(nil)
		for itr.Next() {
			if itr.Value() == nil || len(itr.Value()) == 0 {
				batch.Delete(itr.Key())
			} else {
				batch.Put(itr.Key(), itr.Value())
			}
		}
		batch.Delete(s.committed[i].BlockKey())
		itr.Release()
	}
	logger.Debug("write to basedb", "from", s.committed[0].Number, "to", s.committed[commitNum-1].Number, "len", len(s.committed), "commitNum", commitNum)
	if err := s.baseDB.Write(batch, nil); err != nil {
		logger.Error("write to baseDB fail", "err", err)
		return errors.New("[SnapshotDB]write to baseDB fail:" + err.Error())
	}
	return nil
}

//NewBlock call when you need a new unRecognized or recognized  block data
//it will set JournalHeader for the block
//if hash nil ,new unRecognized data
//if hash not nul,new Recognized data
func (s *snapshotDB) NewBlock(blockNumber *big.Int, parentHash common.Hash, hash common.Hash) error {
	if blockNumber == nil {
		return errors.New("[SnapshotDB]the blockNumber must not be nil ")
	}
	if s.current.GetHighest(false).Num.Cmp(blockNumber) >= 0 {
		logger.Error("the block is less than commit highest", "commit", s.current.GetHighest(false).Num, "new", blockNumber)
		return ErrBlockTooLow
	}
	block := new(blockData)
	block.Number = new(big.Int).Set(blockNumber)
	block.ParentHash = parentHash
	block.BlockHash = hash
	block.data = memdb.New(DefaultComparer, 100)
	block.journal = make([]journalEntry, 0)
	block.validRevisions = make([]revision, 0)

	s.unCommit.Set(hash, block)
	logger.Info("NewBlock", "num", block.Number, "hash", hash, "parent", parentHash)
	return nil
}

func (s *snapshotDB) RevertToSnapshot(hash common.Hash, revid int) {
	s.unCommit.Lock()
	defer s.unCommit.Unlock()
	block, ok := s.unCommit.blocks[hash]
	if ok {
		block.RevertToSnapshot(revid)
	}
}
func (s *snapshotDB) Snapshot(hash common.Hash) int {
	s.unCommit.Lock()
	defer s.unCommit.Unlock()
	block, ok := s.unCommit.blocks[hash]
	if !ok {
		return 0
	}
	return block.Snapshot()
}

// Put sets the value for the given key. It overwrites any previous value
// for that key; a DB is not a multi-map.
func (s *snapshotDB) Put(hash common.Hash, key, value []byte) error {
	if err := s.put(hash, key, value); err != nil {
		return fmt.Errorf("[SnapshotDB]put fail:%v", err)
	}
	return nil
}

func (s *snapshotDB) getFromUnCommit(hash common.Hash, key []byte) ([]byte, error) {
	s.unCommit.RLock()
	defer s.unCommit.RUnlock()
	for {
		if block, ok := s.unCommit.blocks[hash]; ok {
			if hash == block.ParentHash {
				return nil, errors.New("getFromRecognized loop error")
			}
			v, err := block.data.Get(key)
			if err == nil {
				return v, nil
			}
			if err == memdb.ErrNotFound {
				hash = block.ParentHash
				continue
			}
			return nil, err
		} else {
			break
		}
	}
	return nil, ErrNotFound
}

func (s *snapshotDB) getFromCommit(key []byte) ([]byte, error) {
	s.commitLock.RLock()
	defer s.commitLock.RUnlock()
	if len(s.committed) > 0 {
		for i := len(s.committed) - 1; i >= 0; i-- {
			v, err := s.committed[i].data.Get(key)
			if err == nil {
				return v, nil
			}
			if err == memdb.ErrNotFound {
				continue
			}
			return nil, err
		}
	}
	return nil, ErrNotFound
}

// Get get key,val from  snapshotDB
// if hash is nil, unRecognizedBlockData > RecognizedBlockData > CommittedBlockData > baseDB
// if hash is not nil,it will find from the chain, RecognizedBlockData > CommittedBlockData > baseDB
func (s *snapshotDB) Get(hash common.Hash, key []byte) ([]byte, error) {
	v, err := s.getFromUnCommit(hash, key)
	if err != nil && err != ErrNotFound {
		return nil, err
	}
	if err == nil {
		if v == nil || len(v) == 0 {
			return nil, ErrNotFound
		}
		return v, nil
	}
	valueFromCommit, errFromCommit := s.getFromCommit(key)
	if errFromCommit != nil && errFromCommit != ErrNotFound {
		return nil, errFromCommit
	}
	if errFromCommit == nil {
		if valueFromCommit == nil || len(valueFromCommit) == 0 {
			return nil, ErrNotFound
		}
		return valueFromCommit, nil
	}
	return s.GetBaseDB(key)
}

func (s *snapshotDB) GetBaseDB(key []byte) ([]byte, error) {
	if v, err := s.baseDB.Get(key, nil); err == nil {
		return v, nil
	} else if err != leveldb.ErrNotFound {
		return nil, err
	} else {
		return nil, ErrNotFound
	}
}

//Has check the key is exist in chain
//same logic with get
func (s *snapshotDB) Has(hash common.Hash, key []byte) (bool, error) {
	_, err := s.Get(hash, key)
	if err == nil {
		return true, nil
	} else if err == ErrNotFound {
		return true, ErrNotFound
	} else {
		return false, err
	}
}

// Flush move unRecognized to Recognized data
func (s *snapshotDB) Flush(hash common.Hash, blockNumber *big.Int) error {
	if s.dbError != nil {
		return s.dbError
	}
	s.unCommit.RLock()
	block, ok := s.unCommit.blocks[s.getUnRecognizedHash()]
	if !ok {
		s.unCommit.RUnlock()
		return errors.New("[snapshotdb]the unRecognized is nil, can't flush")
	}
	if _, ok := s.unCommit.blocks[hash]; ok {
		s.unCommit.RUnlock()
		return errors.New("the hash is exist in recognized data")
	}
	s.unCommit.RUnlock()

	if block.Number == nil {
		return errors.New("[snapshotdb]the unRecognized Number is nil, can't flush")
	}
	if blockNumber.Uint64() != block.Number.Uint64() {
		return fmt.Errorf("[snapshotdb]blocknumber not compare the unRecognized blocknumber=%v,unRecognizedNumber=%v", blockNumber.Uint64(), block.Number.Uint64())
	}
	s.unCommit.Lock()
	block.BlockHash = hash
	block.readOnly = true
	block.cleanJournal()
	s.unCommit.blocks[hash] = block
	delete(s.unCommit.blocks, common.ZeroHash)
	s.unCommit.Unlock()

	return nil
}

func (s *snapshotDB) theBlockIsCommit(block *blockData) bool {
	if block.Number.Cmp(s.current.GetHighest(false).Num) != 0 {
		return false
	}
	if block.BlockHash != s.current.GetHighest(false).Hash {
		return false
	}
	return true
}

// Commit move blockdata from recognized to commit
func (s *snapshotDB) Commit(hash common.Hash) error {
	if s.dbError != nil {
		return s.dbError
	}
	s.unCommit.RLock()
	block, ok := s.unCommit.blocks[hash]
	s.unCommit.RUnlock()
	if !ok {
		return errors.New("[snapshotdb]commit fail, not found block from recognized :" + hash.String())
	}
	if s.theBlockIsCommit(block) {
		s.unCommit.Lock()
		delete(s.unCommit.blocks, hash)
		s.unCommit.Unlock()
		logger.Info("commit block", "num", block.Number, "hash", hash.String())
		return nil
	}

	isFirstBlock := s.current.GetHighest(false).Num.Uint64() == 0 && block.Number.Uint64() == 0
	if !isFirstBlock {
		hight := s.current.GetHighest(true)
		if hight.Num.Cmp(block.Number) >= 0 {
			return fmt.Errorf("[snapshotdb]commit fail,the commit block num  %v is less or eq than HighestNum %v", block.Number, hight.Num)
		}
		if (block.Number.Uint64() - hight.Num.Uint64()) != 1 {
			return fmt.Errorf("[snapshotdb]commit fail,the commit block num %v - HighestNum %v should be eq 1", block.Number, hight.Num)
		}
		if hight.Hash != common.ZeroHash {
			if block.ParentHash != s.current.GetHighest(false).Hash {
				return fmt.Errorf("[snapshotdb]commit fail,the commit block ParentHash %v not eq HighestHash of commit hash %v ", block.ParentHash.String(), hight.Hash.String())
			}
		}
	}

	block.readOnly = true
	s.writeBlockToWalAsynchronous(block)

	s.commitLock.Lock()
	s.current.increaseHighest(hash)
	block.cleanJournal()
	s.committed = append(s.committed, block)
	s.commitLock.Unlock()

	s.unCommit.Lock()
	delete(s.unCommit.blocks, hash)
	s.unCommit.Unlock()
	logger.Info("commit block", "num", block.Number, "hash", hash.String())
	return nil
}

func (s *snapshotDB) BaseNum() (*big.Int, error) {
	if s.current == nil {
		return nil, errors.New("current is nil")
	}
	return s.current.GetBase(true).Num, nil
}

// WalkBaseDB returns a latest snapshot of the underlying DB. A snapshot
// is a frozen snapshot of a DB state at a particular point in time. The
// content of snapshot are guaranteed to be consistent.
// slice
func (s *snapshotDB) WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
	logger.Debug("begin walkbase db")
	snapshot, err := s.baseDB.GetSnapshot()
	if err != nil {
		return errors.New("[snapshotdb] get snapshot fail:" + err.Error())
	}
	defer snapshot.Release()
	t := snapshot.NewIterator(slice, nil)
	defer func() {
		logger.Debug("WalkBaseDB release ")
		t.Release()
	}()
	return f(s.current.GetBase(true).Num, t)
}

// Clear close db , remove all db file
func (s *snapshotDB) Clear() error {
	if s == nil {
		return errors.New("snapshotDB is nil")
	}
	if err := s.Close(); err != nil {
		return err
	}
	logger.Info("begin clear file", "path", s.path)
	if err := os.RemoveAll(s.path); err != nil {
		return err
	}
	return nil
}

// Ranking return iterates  of the DB.
// the key range that satisfy the given prefix
// the hash means from  unRecognized or recognized
// return iterates ,iterates over a DB's key/value pairs in key order.
// The iterator must be released after use, by calling Release method.t
// Also read Iterator documentation of the leveldb/iterator package.
func (s *snapshotDB) Ranking(hash common.Hash, key []byte, rangeNumber int) iterator.Iterator {
	prefix := util.BytesPrefix(key)
	var itrs []iterator.Iterator
	var parentHash common.Hash
	parentHash = hash
	s.unCommit.RLock()
	for {
		if block, ok := s.unCommit.blocks[parentHash]; ok {
			itrs = append(itrs, block.data.NewIterator(prefix))
			parentHash = block.ParentHash
		} else {
			break
		}
	}
	s.unCommit.RUnlock()
	s.commitLock.RLock()
	for i := len(s.committed) - 1; i >= 0; i-- {
		block := s.committed[i]
		if block.BlockHash == hash || block.BlockHash == parentHash {
			itrs = append(itrs, block.data.NewIterator(prefix))
			parentHash = block.ParentHash
		}
	}
	s.commitLock.RUnlock()
	//put  unCommit and commit itr to heap
	rankingHeap := newRankingHeap(rangeNumber)
	for i := 0; i < len(itrs); i++ {
		rankingHeap.itr2Heap(itrs[i], false, false)
	}
	//put baseDB itr to heap
	itr := s.baseDB.NewIterator(prefix, nil)
	rankingHeap.itr2Heap(itr, true, true)
	//generate memdb Iterator
	mdb := memdb.New(DefaultComparer, rangeNumber)
	for rankingHeap.heap.Len() > 0 {
		kv := heap.Pop(&rankingHeap.heap).(kv)
		if err := mdb.Put(kv.key, kv.value); err != nil {
			return iterator.NewEmptyIterator(errors.New("put to mdb fail" + err.Error()))
		}
	}
	rankingHeap = nil
	return mdb.NewIterator(nil)
}

func (s *snapshotDB) Close() error {
	logger.Info("begin close snapshotdb", "path", s.path)
	//	runtime.SetFinalizer(s, nil)
	if s == nil {
		return nil
	}
	if s.closed {
		return nil
	}
	if s.corn != nil {
		s.corn.Stop()
	}
	s.walSync.Wait()
	close(s.walExitCh)

	if s.baseDB != nil {
		if err := s.baseDB.Close(); err != nil {
			return fmt.Errorf("[snapshotdb]close base db fail:%v", err)
		}
	}

	s.current = nil
	s.unCommit = nil
	s.committed = nil
	s.closed = true
	logger.Info("snapshotdb closed")
	return nil
}

func IsDbNotFoundErr(err error) bool {
	return nil != err && err == ErrNotFound
}

func NonDbNotFoundErr(err error) bool {
	return nil != err && err != ErrNotFound
}
