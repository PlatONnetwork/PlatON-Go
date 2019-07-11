package snapshotdb

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/robfig/cron"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math/big"
	"os"
	"sync"
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
	WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error
	Commit(hash common.Hash) error

	// Clear close db , remove all db file,set dbInstance nil
	Clear() error

	PutBaseDB(key, value []byte) error
	GetBaseDB(key []byte) ([]byte, error)

	// WriteBaseDB apply the given [][2][]byte to the baseDB.
	WriteBaseDB(kvs [][2][]byte) error

	//SetCurrent use for fast sync
	SetCurrent(highestHash common.Hash, base, height big.Int) error

	DelBaseDB(key []byte) error
	GetLastKVHash(blockHash common.Hash) []byte
	BaseNum() (*big.Int, error)
	Close() error
	Compaction() error
}

var (
	dbpath string

	dbInstance *snapshotDB

	instance sync.Mutex

	logger = log.Root().New("package", "snapshotdb")

	//ErrorSnaphotLock when db is Lock
	ErrorSnaphotLock = errors.New("can't create snapshot,snapshot is lock now")

	//ErrNotFound when db not found
	ErrNotFound = errors.New("snapshotDB: not found")

	//ErrDBNotInit when db  not init
	ErrDBNotInit = errors.New("snapshotDB: not init")
)

type snapshotDB struct {
	path string

	snapshotLockC bool
	snapshotLock  *sync.Cond

	current *current
	baseDB  *leveldb.DB

	unRecognized     *blockData
	unRecognizedLock sync.RWMutex

	recognized sync.Map

	committed  []blockData
	commitLock sync.RWMutex

	journalw          map[common.Hash]*journalWriter
	journalWriterLock sync.RWMutex

	storage storage

	corn *cron.Cron

	closed bool
}

//SetDBPath set db path
func SetDBPath(ctx *node.ServiceContext) {
	dbpath = ctx.ResolvePath(DBPath)
}

func SetDBPathWithNode(n *node.Node) {
	dbpath = n.ResolvePath(DBPath)
}

//Instance return the Instance of the db
func Instance() DB {
	instance.Lock()
	defer instance.Unlock()
	if dbInstance == nil || dbInstance.closed {
		logger.Debug("dbInstance is nil", "path", dbpath)
		if err := initDB(); err != nil {
			logger.Error(fmt.Sprint("init db fail"), err)
			panic(err)
			//return nil, errors.New("init db fail:" + err.Error())
		}
	}
	return dbInstance
}

func Open(path string) (DB, error) {
	s, err := openFile(path, false)
	if err != nil {
		logger.Error("open db file fail", "error", err, "path", dbpath)
		return nil, err
	}
	fds, err := s.List(TypeCurrent)
	if err != nil {
		logger.Error("get current file fail", "error", err)
		return nil, err
	}
	if len(fds) > 0 {
		logger.Info("begin open")
		db := new(snapshotDB)
		if err := db.recover(s); err != nil {
			logger.Error("recover db fail:", "error", err)
			return nil, err
		}
		return db, nil
	}
	return nil, nil
}

func copyDB(from, to *snapshotDB) {
	to.path = from.path
	to.snapshotLockC = from.snapshotLockC
	to.snapshotLock = from.snapshotLock
	to.current = from.current
	to.baseDB = from.baseDB
	to.unRecognized = from.unRecognized
	to.recognized = from.recognized
	to.committed = from.committed
	to.journalw = from.journalw
	to.storage = from.storage
	to.corn = from.corn
	to.closed = from.closed
}

func initDB() error {
	s, err := openFile(dbpath, false)
	if err != nil {
		logger.Error("open db file fail", "error", err, "path", dbpath)
		return err
	}
	fds, err := s.List(TypeCurrent)
	if err != nil {
		logger.Error("get current file fail", "error", err)
		return err
	}
	if dbInstance == nil {
		dbInstance = new(snapshotDB)
	}
	if len(fds) > 0 {
		logger.Info("begin recover")
		db := new(snapshotDB)
		if err := db.recover(s); err != nil {
			logger.Error("recover db fail:", "error", err)
			return err
		}
		copyDB(db, dbInstance)
	} else {
		logger.Info("begin newDB")
		db, err := newDB(s)
		if err != nil {
			logger.Error(fmt.Sprint("new db fail:", err))
			return err
		}
		copyDB(db, dbInstance)
	}
	dbInstance.corn = cron.New()
	if err := dbInstance.corn.AddFunc("@every 1s", dbInstance.schedule); err != nil {
		logger.Error(fmt.Sprint("new db fail", err))
		return err
	}
	dbInstance.corn.Start()
	return err
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
	if err := s.current.update(); err != nil {
		return err
	}
	return nil
}

// GetCommittedBlock    get value from committed blockdata > baseDB
func (s *snapshotDB) GetFromCommittedBlock(key []byte) ([]byte, error) {
	s.commitLock.RLock()
	defer s.commitLock.RUnlock()
	var (
		v   []byte
		err error
	)
	for i := len(s.committed) - 1; i >= 0; i-- {
		v, err = s.committed[i].data.Get(key)
		if err == nil {
			break
		} else if err != memdb.ErrNotFound {
			logger.Error(fmt.Sprintf(" find from committed hash:%s fail,%v", string(key), err))
			return nil, err
		}
	}
	if err != nil {
		v, err = s.baseDB.Get(key, nil)
		if err != nil {
			if err == leveldb.ErrNotFound {
				return nil, ErrNotFound
			} else {
				return nil, err
			}
		}
	}
	if v == nil || len(v) == 0 {
		return nil, ErrNotFound
	} else {
		return v, nil
	}
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
	if blockHash == common.ZeroHash {
		return s.unRecognized.kvHash.Bytes()
	}
	block, ok := s.recognized.Load(blockHash)
	if !ok {
		return nil
	}
	b := block.(blockData)
	return b.kvHash.Bytes()
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
	s.commitLock.Lock()
	s.snapshotLock.L.Lock()
	s.snapshotLockC = true
	defer func() {
		s.snapshotLockC = false
		s.snapshotLock.Broadcast()
		s.snapshotLock.L.Unlock()
		s.commitLock.Unlock()
	}()
	var (
		kvsize    int
		commitNum int
	)
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
	if commitNum == 0 {
		commitNum++
	}
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
	}
	if err := s.baseDB.Write(batch, nil); err != nil {
		logger.Error(fmt.Sprint("write to baseDB fail:", err))
		return errors.New("[SnapshotDB]write to baseDB fail:" + err.Error())
	}
	s.current.BaseNum.Add(s.current.BaseNum, big.NewInt(int64(commitNum)))
	if err := s.current.update(); err != nil {
		logger.Error(fmt.Sprint("update to current fail:", err))
		return errors.New("[SnapshotDB]update to current fail:" + err.Error())
	}
	s.committed = s.committed[commitNum:len(s.committed)]
	if err := s.removeJournalLessThanBaseNum(); err != nil {
		logger.Error(fmt.Sprint("remove journal less than baseNum fail:", err))
		return errors.New("[SnapshotDB]remove journal less than baseNum fail:" + err.Error())
	}
	return nil
}

//NewBlock call when you need a new unRecognized or recognized  block data
//it will set JournalHeader for the block
//if hash nil ,new unRecognized data
//if hash not nul,new Recognized data
func (s *snapshotDB) NewBlock(blockNumber *big.Int, parentHash common.Hash, hash common.Hash) error {
	if hash == common.ZeroHash {
		if s.unRecognized != nil && s.unRecognized.readOnly {
			return errors.New("[SnapshotDB]can't  new unRecognized block,it's have value,must flush it before NewBlock ")
		}
	}
	if blockNumber == nil {
		return errors.New("[SnapshotDB]the blockNumber must not be nil ")
	}
	block := new(blockData)
	block.Number = big.NewInt(blockNumber.Int64())
	block.ParentHash = parentHash
	block.BlockHash = hash
	block.data = memdb.New(DefaultComparer, 100)
	if hash == common.ZeroHash {
		if err := s.writeJournalHeader(blockNumber, s.getUnRecognizedHash(), parentHash, journalHeaderFromUnRecognized); err != nil {
			return fmt.Errorf("[SnapshotDB] write Journal Header fail:%v", err)
		}
		s.unRecognized = block
	} else {
		if err := s.writeJournalHeader(blockNumber, hash, parentHash, journalHeaderFromRecognized); err != nil {
			return fmt.Errorf("[SnapshotDB] write Journal Header fail:%v", err)
		}
		s.recognized.Store(hash, *block)
	}
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

// Get get key,val from  snapshotDB
// if hash is nil, unRecognizedBlockData > RecognizedBlockData > CommittedBlockData > baseDB
// if hash is not nil,it will find from the chain, RecognizedBlockData > CommittedBlockData > baseDB
func (s *snapshotDB) Get(hash common.Hash, key []byte) ([]byte, error) {
	//found from unRecognized
	location, ok := s.checkHashChain(hash)
	if !ok {
		return nil, errors.New("[SnapshotDB]the hash not in chain " + hash.String())
	}
	//found from Recognized
	var (
		v   []byte
		err error
	)

	switch location {
	case hashLocationUnRecognized:
		if s.unRecognized == nil {
			return nil, errors.New("[SnapshotDB]unRecognized is not find now")
		}
		if v, err = s.unRecognized.data.Get(key); err == nil {

		} else if err != memdb.ErrNotFound {
			return nil, err
		} else {
			v, err = s.getFromRecognized(s.unRecognized.ParentHash, key)
		}
	case hashLocationRecognized:
		v, err = s.getFromRecognized(hash, key)
	case hashLocationCommitted:
		v, err = s.getFromCommitted(hash, key)
	default:
		v, err = s.GetBaseDB(key)
	}
	if err != nil {
		return nil, err
	}
	if v == nil || len(v) == 0 {
		return v, ErrNotFound
	}
	return v, nil
}

func (s *snapshotDB) getFromRecognized(hash common.Hash, key []byte) ([]byte, error) {
	for {
		if b, ok := s.recognized.Load(hash); ok {
			block := b.(blockData)
			if hash == block.ParentHash {
				return nil, errors.New("getFromRecognized loop error")
			}
			if v, err := block.data.Get(key); err == nil {
				return v, nil
			} else if err != memdb.ErrNotFound {
				return nil, err
			}
			hash = block.ParentHash
		} else {
			break
		}
	}
	s.commitLock.RLock()
	defer s.commitLock.RUnlock()
	if len(s.committed) > 0 {
		block := s.committed[len(s.committed)-1]
		if block.BlockHash != hash {
			return nil, ErrNotFound
		}
		v, err := s.getFromCommitted(hash, key)
		if err == nil {
			return v, nil
		} else if err != ErrNotFound {
			return nil, err
		}
	}
	return s.GetBaseDB(key)
}

func (s *snapshotDB) getFromCommitted(hash common.Hash, key []byte) ([]byte, error) {
	s.commitLock.RLock()
	defer s.commitLock.RUnlock()
	for i := len(s.committed) - 1; i >= 0; i-- {
		if s.committed[i].BlockHash == hash {
			if v, err := s.committed[i].data.Get(key); err == nil {
				return v, nil
			} else if err != memdb.ErrNotFound {
				return nil, err
			}
			hash = s.committed[i].ParentHash
		}
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
	if s.unRecognized == nil {
		return errors.New("[snapshotdb]the unRecognized is nil, can't flush")
	}
	if s.unRecognized.Number == nil {
		return errors.New("[snapshotdb]the unRecognized Number is nil, can't flush")
	}
	if blocknumber.Int64() != s.unRecognized.Number.Int64() {
		return errors.New("[snapshotdb]blocknumber not compare the unRecognized blocknumber")
	}
	if _, ok := s.recognized.Load(hash); ok {
		return errors.New("the hash is exist in recognized data")
	}
	currentHash := s.getUnRecognizedHash()
	oldFd := fileDesc{Type: TypeJournal, Num: blocknumber.Int64(), BlockHash: currentHash}
	newFd := fileDesc{Type: TypeJournal, Num: blocknumber.Int64(), BlockHash: hash}
	s.unRecognizedLock.Lock()
	defer s.unRecognizedLock.Unlock()
	if err := s.closeJournalWriter(currentHash); err != nil {
		return err
	}
	if err := s.storage.Rename(oldFd, newFd); err != nil {
		return errors.New("[snapshotdb]rename file fail:" + oldFd.String() + "," + newFd.String() + "," + err.Error())
	}
	s.unRecognized.BlockHash = hash
	s.unRecognized.readOnly = true
	s.recognized.Store(hash, *s.unRecognized)

	s.unRecognized = nil
	return nil
}

// Commit move blockdata from recognized to commit
func (s *snapshotDB) Commit(hash common.Hash) error {
	s.commitLock.Lock()
	defer s.commitLock.Unlock()

	b, ok := s.recognized.Load(hash)
	if !ok {
		return errors.New("[snapshotdb]commit fail, not found block from recognized :" + hash.String())
	}
	block := b.(blockData)
	if s.current.HighestNum.Int64() == 0 && block.Number.Int64() == 0 {

	} else {
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
	block.readOnly = true
	s.committed = append(s.committed, block)
	s.current.HighestNum = block.Number
	s.current.HighestHash = hash
	if err := s.current.update(); err != nil {
		return errors.New("[snapshotdb]commit fail,update current fail:" + err.Error())
	}

	if err := s.closeJournalWriter(hash); err != nil {
		return err
	}
	s.recognized.Delete(hash)
	if err := s.rmOldRecognizedBlockData(); err != nil {
		return err
	}
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
	if s.snapshotLockC {
		s.snapshotLock.L.Lock()
		defer s.snapshotLock.L.Unlock()
		for s.snapshotLockC {
			s.snapshotLock.Wait()
		}
	}
	snapshot, err := s.baseDB.GetSnapshot()
	if err != nil {
		return errors.New("[snapshotdb] get snapshot fail:" + err.Error())
	}
	defer snapshot.Release()
	t := snapshot.NewIterator(slice, nil)
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
	//	s = nil
	return nil
}

func (s *snapshotDB) Reset() {

}

func itrToMdb(itr iterator.Iterator, mdb *memdb.DB) error {
	for itr.Next() {
		if itr.Value() == nil || len(itr.Value()) == 0 {
			if err := mdb.Delete(itr.Key()); err != nil && err != memdb.ErrNotFound {
				return err
			}
		} else {
			if err := mdb.Put(itr.Key(), itr.Value()); err != nil {
				return err
			}
		}
	}
	itr.Release()
	return nil
}

// Ranking return iterates  of the DB.
// the key range that satisfy the given prefix
// the hash means from  unRecognized or recognized
// return iterates ,iterates over a DB's key/value pairs in key order.
// The iterator must be released after use, by calling Release method.
// Also read Iterator documentation of the leveldb/iterator package.
func (s *snapshotDB) Ranking(hash common.Hash, key []byte, rangeNumber int) iterator.Iterator {
	location, ok := s.checkHashChain(hash)
	if !ok {
		return iterator.NewEmptyIterator(errors.New("this hash not in chain:" + hash.String()))
	}
	prefix := util.BytesPrefix(key)
	m := memdb.New(comparer.DefaultComparer, rangeNumber)
	var itrs []iterator.Iterator
	var parentHash common.Hash
	switch location {
	case hashLocationUnRecognized:
		if s.unRecognized.data != nil {
			itrs = append(itrs, s.unRecognized.data.NewIterator(prefix))
			parentHash = s.unRecognized.ParentHash
		}
		for {
			if b, ok := s.recognized.Load(parentHash); ok {
				block := b.(blockData)
				itrs = append(itrs, block.data.NewIterator(prefix))
				parentHash = block.ParentHash
			} else {
				break
			}
		}
		for _, block := range s.committed {
			itrs = append(itrs, block.data.NewIterator(prefix))
		}
	case hashLocationRecognized:
		parentHash = hash
		for {
			if b, ok := s.recognized.Load(parentHash); ok {
				block := b.(blockData)
				itrs = append(itrs, block.data.NewIterator(prefix))
				parentHash = block.ParentHash
			} else {
				break
			}
		}
		for _, block := range s.committed {
			itrs = append(itrs, block.data.NewIterator(prefix))
		}
	case hashLocationCommitted:
		for i := len(s.committed) - 1; i >= 0; i-- {
			block := s.committed[i]
			if block.BlockHash == hash {
				itrs = append(itrs, block.data.NewIterator(prefix))
				parentHash = block.ParentHash
			} else if block.BlockHash == parentHash {
				itrs = append(itrs, block.data.NewIterator(prefix))
				parentHash = block.BlockHash
			}
		}
	}
	itrs = append(itrs, s.baseDB.NewIterator(prefix, nil))
	for i := len(itrs) - 1; i >= 0; i-- {
		if err := itrToMdb(itrs[i], m); err != nil {
			return iterator.NewEmptyIterator(err)
		}
	}
	return m.NewIterator(nil)
}

func (s *snapshotDB) Close() error {
	//	runtime.SetFinalizer(s, nil)
	if s == nil {
		return nil
	}
	if s.corn != nil {
		s.corn.Stop()
	}
	if s.baseDB != nil {
		if err := s.baseDB.Close(); err != nil {
			return fmt.Errorf("[snapshotdb]close base db fail:%v", err)
		}
	}

	if err := s.storage.Close(); err != nil {
		return fmt.Errorf("[snapshotdb]close storage fail:%v", err)
	}
	if s.current != nil {
		s.current.f.Close()
	}

	for key := range s.journalw {
		if err := s.journalw[key].Close(); err != nil {
			return fmt.Errorf("[snapshotdb]close journalw fail:%v", err)
		}
	}
	s.closed = true
	return nil
}
