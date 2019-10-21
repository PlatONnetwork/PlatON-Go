package snapshotdb

import (
	"container/heap"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/robfig/cron"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbError "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	kvLIMIT = 2000
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

	PutBaseDB(key, value []byte) error
	GetBaseDB(key []byte) ([]byte, error)
	DelBaseDB(key []byte) error
	// WriteBaseDB apply the given [][2][]byte to the baseDB.
	WriteBaseDB(kvs [][2][]byte) error

	//SetCurrent use for fast sync
	SetCurrent(highestHash common.Hash, base, height big.Int) error
	GetCurrent() *current

	GetLastKVHash(blockHash common.Hash) []byte
	BaseNum() (*big.Int, error)
	Close() error
	Compaction() error
	SetEmpty() error
}

var (
	dbpath string

	blockchain Chain

	dbInstance *snapshotDB

	instance sync.Mutex

	logger = log.Root().New("package", "snapshotdb")

	//ErrNotFound when db not found
	ErrNotFound = errors.New("snapshotDB: not found")

	ErrBlockRepeat = errors.New("the block is exist in snapshotdb uncommit")
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

	journalSync sync.WaitGroup

	storage storage

	corn *cron.Cron

	closed bool
}

type Chain interface {
	CurrentHeader() *types.Header
	GetHeaderByHash(hash common.Hash) *types.Header
}

type unCommitBlocks struct {
	blocks map[common.Hash]*blockData
	sync.RWMutex
}

func (u *unCommitBlocks) Get(key common.Hash) *blockData {
	u.RLock()
	block, ok := u.blocks[key]
	u.RUnlock()
	if !ok {
		return nil
	}
	return block
}

func (u *unCommitBlocks) Set(key common.Hash, block *blockData) {
	u.Lock()
	u.blocks[key] = block
	u.Unlock()
}

func SetDBPathWithNode(path string) {
	dbpath = path
	logger.Info("set path", "path", dbpath)
}

func SetDBBlockChain(n Chain) {
	blockchain = n
	logger.Info("set blockchain")
}

//Instance return the Instance of the db
func Instance() DB {
	instance.Lock()
	defer instance.Unlock()
	if dbInstance == nil || dbInstance.closed {
		logger.Debug("dbInstance is nil", "path", dbpath)
		if err := initDB(); err != nil {
			logger.Error("init db fail", "err", err)
			panic(err)
		}
	}
	return dbInstance
}

func Open(path string) (DB, error) {
	s, err := openFile(path, false)
	if err != nil {
		logger.Error("open db file fail", "error", err, "path", path)
		return nil, err
	}
	baseDB, err := leveldb.OpenFile(getBaseDBPath(path), nil)
	if err != nil {
		if _, corrupted := err.(*leveldbError.ErrCorrupted); corrupted {
			baseDB, err = leveldb.RecoverFile(getBaseDBPath(path), nil)
			if err != nil {
				return nil, fmt.Errorf("[SnapshotDB.recover]RecoverFile baseDB fail:%v", err)
			}
		} else {
			return nil, err
		}
	}
	unCommitBlock := new(unCommitBlocks)
	unCommitBlock.blocks = make(map[common.Hash]*blockData)
	db := &snapshotDB{
		path:          path,
		storage:       s,
		unCommit:      unCommitBlock,
		committed:     make([]*blockData, 0),
		baseDB:        baseDB,
		snapshotLockC: snapshotUnLock,
	}

	_, getCurrentError := baseDB.Get([]byte(CurrentSet), nil)
	if getCurrentError == nil {
		logger.Info("begin recover", "path", path)
		if err := db.loadCurrent(); err != nil {
			return nil, err
		}
		logger.Info("load current", "current", db.current)
		if err := db.recover(); err != nil {
			logger.Error("recover db fail:", "error", err)
			return nil, err
		}
		return db, nil
	}
	if getCurrentError == leveldb.ErrNotFound {
		logger.Info("begin new", "path", path)
		if err := db.newCurrent(); err != nil {
			return nil, err
		}
		return db, nil
	}
	return nil, getCurrentError
}

func copyDB(from, to *snapshotDB) {
	to.path = from.path
	to.current = from.current
	to.baseDB = from.baseDB
	to.unCommit = from.unCommit
	to.committed = from.committed
	to.storage = from.storage
	to.corn = from.corn
	to.closed = from.closed
	to.snapshotLockC = from.snapshotLockC
}

func initDB() error {
	dbInterface, err := Open(dbpath)
	if err != nil {
		return err
	}
	db := dbInterface.(*snapshotDB)
	if dbInstance == nil {
		dbInstance = new(snapshotDB)
	}
	copyDB(db, dbInstance)
	//	dbInstance.writeCurrentLoop()
	if err := dbInstance.cornStart(); err != nil {
		return err
	}
	return nil
}

func (s *snapshotDB) cornStart() error {
	s.corn = cron.New()
	if err := s.corn.AddFunc("@every 1s", s.schedule); err != nil {
		logger.Error("set db corn compaction fail", "err", err)
		return err
	}
	if err := s.corn.AddFunc("@every 3s", s.metrics); err != nil {
		logger.Error("set db corn metrics fail", "err", err)
		return err
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
	s.current.HighestNum = &height
	s.current.BaseNum = &base
	s.current.HighestHash = highestHash
	logger.Debug("SetCurrent", "current", s.current)
	if err := s.saveCurrentToBaseDB("", s.current); err != nil {
		return err
	}
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
	path := s.path
	if err := s.Clear(); err != nil {
		return err
	}
	dbInterface, err := Open(path)
	if err != nil {
		return err
	}
	db := dbInterface.(*snapshotDB)
	copyDB(db, s)
	if err := s.cornStart(); err != nil {
		return err
	}
	return nil
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
	s.current.BaseNum.Add(s.current.BaseNum, big.NewInt(int64(commitNum)))
	s.committed = s.committed[commitNum:]
	s.commitLock.Unlock()
	if err := s.saveCurrentToBaseDB(CurrentBaseNum, &current{
		HighestNum:  nil,
		HighestHash: common.Hash{},
		BaseNum:     new(big.Int).Set(s.current.BaseNum),
	}); err != nil {
		logger.Error("save base to current fail", "err", err)
	}
	if err := s.rmExpireForkBlockJournal(); err != nil {
		return err
	}
	//delete block no use  in unCommit
	s.rmExpireForkBlock()
	return nil
}

func (s *snapshotDB) rmExpireForkBlockJournal() error {
	fds, err := s.storage.List(TypeJournal)
	if err != nil {
		return err
	}
	for _, fd := range fds {
		if s.current.BaseNum.Uint64() >= fd.Num {
			if err := s.storage.Remove(fd); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *snapshotDB) rmExpireForkBlock() {
	s.unCommit.Lock()
	// if unCommit blocks length > 200,rm block below BaseNum
	if len(s.unCommit.blocks) > 200 {
		for key, value := range s.unCommit.blocks {
			if s.current.BaseNum.Cmp(value.Number) >= 0 {
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
	if len(s.committed) > 200 {
		commitNum = 100
	} else {
		for i := 0; i < len(s.committed); i++ {
			if i < 10 {
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

	header := blockchain.CurrentHeader()
	var length = commitNum
	for i := 0; i < length; i++ {
		if s.committed[i].Number.Cmp(header.Number) > 0 {
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
		itr.Release()
	}
	//logger.Debug("write to basedb", "from", s.committed[0].Number, "to", s.committed[commitNum-1].Number, "len", len(s.committed), "commitNum", commitNum)
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
	findBlock := s.unCommit.Get(hash)
	//a block can't new twice
	if findBlock != nil {
		//  if block num is different,hash is same as common.ZeroHash,the exsist block may have commit ,so just cover it
		newBlockWithDiffNumber := findBlock.BlockHash == common.ZeroHash && findBlock.Number.Cmp(blockNumber) != 0
		if !newBlockWithDiffNumber {
			logger.Error("the block is exist in snapshotdb uncommit,can't NewBlock", "hash", hash)
			return ErrBlockRepeat
		}
	}
	if s.current.HighestNum.Cmp(blockNumber) >= 0 {
		logger.Error("the block is less than commit highest", "commit", s.current.HighestNum, "new", blockNumber)
		return ErrBlockTooLow
	}

	block := new(blockData)
	block.Number = new(big.Int).Set(blockNumber)
	block.ParentHash = parentHash
	block.BlockHash = hash
	block.data = memdb.New(DefaultComparer, 100)
	s.unCommit.Set(hash, block)
	logger.Info("NewBlock", "num", block.Number, "hash", hash.String())
	return nil
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
	v2, err2 := s.getFromCommit(key)
	if err2 != nil && err2 != ErrNotFound {
		return nil, err2
	}
	if err2 == nil {
		if v2 == nil || len(v2) == 0 {
			return nil, ErrNotFound
		}
		return v2, nil
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
func (s *snapshotDB) Flush(hash common.Hash, blocknumber *big.Int) error {
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
	if blocknumber.Uint64() != block.Number.Uint64() {
		return fmt.Errorf("[snapshotdb]blocknumber not compare the unRecognized blocknumber=%v,unRecognizedNumber=%v", blocknumber.Uint64(), block.Number.Uint64())
	}
	block.BlockHash = hash
	block.readOnly = true
	s.unCommit.Lock()
	s.writeBlockToJournalAsynchronous(block)
	s.unCommit.blocks[hash] = block
	delete(s.unCommit.blocks, common.ZeroHash)
	s.unCommit.Unlock()

	return nil
}

func (s *snapshotDB) theBlockIsCommit(block *blockData) bool {
	if block.Number.Cmp(s.current.HighestNum) != 0 {
		return false
	}
	if block.BlockHash != s.current.HighestHash {
		return false
	}
	return true
}

// Commit move blockdata from recognized to commit
func (s *snapshotDB) Commit(hash common.Hash) error {
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
		logger.Info("[snapshotDB]commit block", "num", block.Number, "hash", hash.String())
		return nil
	}

	isFirstBlock := s.current.HighestNum.Int64() == 0 && block.Number.Int64() == 0
	if !isFirstBlock {
		if s.current.HighestNum.Cmp(block.Number) >= 0 {
			return fmt.Errorf("[snapshotdb]commit fail,the commit block num  %v is less or eq than HighestNum %v", block.Number, s.current.HighestNum)
		}
		if (block.Number.Int64() - s.current.HighestNum.Int64()) != 1 {
			return fmt.Errorf("[snapshotdb]commit fail,the commit block num %v - HighestNum %v should be eq 1", block.Number, s.current.HighestNum)
		}
		if s.current.HighestHash != common.ZeroHash {
			if block.ParentHash != s.current.HighestHash {
				return fmt.Errorf("[snapshotdb]commit fail,the commit block ParentHash %v not eq HighestHash of commit hash %v ", block.ParentHash.String(), s.current.HighestHash.String())
			}
		}
	}
	if !block.readOnly {
		s.writeBlockToJournalAsynchronous(block)
	} else {
		if err := s.saveCurrentToBaseDB(CurrentHighestBlock, &current{
			HighestNum:  block.Number,
			HighestHash: block.BlockHash,
			BaseNum:     s.current.BaseNum,
		}); err != nil {
			return err
		}
	}
	block.readOnly = true
	s.commitLock.Lock()
	s.committed = append(s.committed, block)
	s.current.HighestHash = hash
	s.current.HighestNum = new(big.Int).Set(block.Number)
	s.commitLock.Unlock()

	s.unCommit.Lock()
	delete(s.unCommit.blocks, hash)
	s.unCommit.Unlock()
	logger.Info("[snapshotDB]commit block", "num", block.Number, "hash", hash.String())
	return nil
}

func (s *snapshotDB) BaseNum() (*big.Int, error) {
	if s.current == nil {
		return nil, errors.New("current is nil")
	}
	return s.current.BaseNum, nil
}

// WalkBaseDB returns a latest snapshot of the underlying DB. A snapshot
// is a frozen snapshot of a DB state at a particular point in time. The
// content of snapshot are guaranteed to be consistent.
// slice
func (s *snapshotDB) WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
	logger.Info("begin walkbase db")
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
	return f(s.current.BaseNum, t)
}

// Clear close db , remove all db file
func (s *snapshotDB) Clear() error {
	if s == nil {
		return errors.New("snapshotDB is nil")
	}
	if err := s.Close(); err != nil {
		return err
	}
	logger.Info(fmt.Sprint("begin clear file:", s.path))
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
	//	t := time.Now()
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
	//	logger.Info("Ranking uncommit", "rangeNumber", rangeNumber, "hash", hash, "duration", time.Since(t))

	//	t = time.Now()
	s.commitLock.RLock()
	for i := len(s.committed) - 1; i >= 0; i-- {
		block := s.committed[i]
		if block.BlockHash == hash || block.BlockHash == parentHash {
			itrs = append(itrs, block.data.NewIterator(prefix))
			parentHash = block.ParentHash
		}
	}
	s.commitLock.RUnlock()
	//	logger.Info("Ranking commit", "rangeNumber", rangeNumber, "hash", hash, "duration", time.Since(t))

	//	t = time.Now()
	//put  unCommit and commit itr to heap
	rankingHeap := newRankingHeap(rangeNumber)
	for i := 0; i < len(itrs); i++ {
		rankingHeap.itr2Heap(itrs[i], false, false)
	}
	//	logger.Info("Ranking heap", "rangeNumber", rangeNumber, "hash", hash, "duration", time.Since(t))

	//	t = time.Now()
	//put baseDB itr to heap
	itr := s.baseDB.NewIterator(prefix, nil)
	rankingHeap.itr2Heap(itr, true, true)
	//	logger.Info("Ranking base", "rangeNumber", rangeNumber, "hash", hash, "duration", time.Since(t))

	//generate memdb Iterator
	//	t = time.Now()
	mdb := memdb.New(DefaultComparer, rangeNumber)
	for rankingHeap.heap.Len() > 0 {
		kv := heap.Pop(&rankingHeap.heap).(kv)
		if err := mdb.Put(kv.key, kv.value); err != nil {
			return iterator.NewEmptyIterator(errors.New("put to mdb fail" + err.Error()))
		}
	}
	rankingHeap = nil
	//	logger.Info("Ranking pop", "rangeNumber", rangeNumber, "hash", hash, "duration", time.Since(t))
	return mdb.NewIterator(nil)
}

func (s *snapshotDB) Close() error {
	logger.Info("begin close snapshotdb", "path", s.path)
	//	runtime.SetFinalizer(s, nil)
	if s == nil {
		return nil
	}
	if s.corn != nil {
		s.corn.Stop()
	}
	s.journalSync.Wait()

	if s.baseDB != nil {
		if err := s.baseDB.Close(); err != nil {
			return fmt.Errorf("[snapshotdb]close base db fail:%v", err)
		}
	}

	if err := s.storage.Close(); err != nil {
		return fmt.Errorf("[snapshotdb]close storage fail:%v", err)
	}
	s.current = nil
	s.unCommit = nil
	s.committed = nil
	s.closed = true
	logger.Info("snapshotdb closed")
	return nil
}
