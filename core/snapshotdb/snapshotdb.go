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
	Clear() error
	PutBaseDB(key, value []byte) error
	DelBaseDB(key []byte) error
	GetLastKVHash(blockHash common.Hash) []byte
	BaseNum() (*big.Int, error)
	Close() error
	Compaction() error
}

var (
	dbpath string

	dbInstance *snapshotDB

	logger = log.Root().New("package", "snapshotdb")

	//ErrorSnaphotLock when db is Lock
	ErrorSnaphotLock = errors.New("can't create snapshot,snapshot is lock now")

	//ErrNotFound when db not found
	ErrNotFound = errors.New("snapshotDB: not found")
)

type snapshotDB struct {
	path string

	snapshotLockC bool
	snapshotLock  *sync.Cond

	current *current
	baseDB  *leveldb.DB

	unRecognized     *blockData
	unRecognizedLock sync.RWMutex

	recognized map[common.Hash]blockData

	committed  []blockData
	commitLock sync.RWMutex

	journalw map[common.Hash]*journalWriter
	storage  storage

	corn *cron.Cron

	closed bool
}

//SetDBPath set db path
func SetDBPath(ctx *node.ServiceContext) {
	dbpath = ctx.ResolvePath(DBPath)
}

//Instance return the Instance of the db
func Instance() DB {
	if dbInstance == nil || dbInstance.closed {
		if err := initDB(); err != nil {
			logger.Error(fmt.Sprint("init db fail"), err)
			panic(err)
			//return nil, errors.New("init db fail:" + err.Error())
		}
	}
	return dbInstance
}

func initDB() error {
	s, err := openFile(dbpath, false)
	if err != nil {
		logger.Error(fmt.Sprint("open db file fail:", err))
		return err
	}
	fds, err := s.List(TypeCurrent)
	if err != nil {
		logger.Error(fmt.Sprint("get current file fail:", err))
		return err
	}
	if len(fds) > 0 {
		db := new(snapshotDB)
		if err := db.recover(s); err != nil {
			logger.Error(fmt.Sprint("recover  db fail:", err))
			return err
		}
		dbInstance = db
	} else {
		db, err := newDB(s)
		if err != nil {
			logger.Error(fmt.Sprint("new db fail:", err))
			return err
		}
		dbInstance = db
	}
	dbInstance.corn = cron.New()
	if err := dbInstance.corn.AddFunc("@every 1s", dbInstance.schedule); err != nil {
		logger.Error(fmt.Sprint("new db fail", err))
		return err
	}
	dbInstance.corn.Start()
	return err
}

// GetCommittedBlock    get value from committed blockdata > baseDB
func (s *snapshotDB) GetFromCommittedBlock(key []byte) ([]byte, error) {
	s.commitLock.RLock()
	defer s.commitLock.RUnlock()
	for _, value := range s.committed {
		if v, err := value.data.Get(key); err == nil {
			return v, nil
		} else if err != memdb.ErrNotFound {
			logger.Error(fmt.Sprintf(" find from committed hash:%s fail,%v", string(key), err))
			return nil, err
		}
	}
	v, err := s.baseDB.Get(key, nil)
	if err == nil {
		return v, nil
	} else if err == leveldb.ErrNotFound {
		return nil, ErrNotFound
	}
	logger.Error(fmt.Sprintf("hash:%s find from base fail,%v", string(key), err))
	return nil, err
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
	block, ok := s.recognized[blockHash]
	if !ok {
		return nil
	}
	return block.kvHash.Bytes()
}

// Del del key,val from  snapshotDB
// if hash is nil, unRecognizedBlockData > recognizedBlockData
// if hash is not nil,it will del in recognized BlockData
func (s *snapshotDB) Del(hash common.Hash, key []byte) error {
	if err := s.put(hash, key, nil, funcTypeDel); err != nil {
		return err
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
			commitNum = 9
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
			batch.Put(itr.Key(), itr.Value())
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
			return errors.New("[SnapshotDB]can't  new unRecognized block,it's readonly now")
		}
	}
	block := new(blockData)
	block.Number = blockNumber
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
			return fmt.Errorf("[SnapshotDB] write Journal body fail:%v", err)
		}
		s.recognized[hash] = *block
	}
	return nil
}

// Put sets the value for the given key. It overwrites any previous value
// for that key; a DB is not a multi-map.
func (s *snapshotDB) Put(hash common.Hash, key, value []byte) error {
	if err := s.put(hash, key, value, funcTypePut); err != nil {
		return err
	}
	return nil
}

// Get get key,val from  snapshotDB
// if hash is nil, unRecognizedBlockData > RecognizedBlockData > CommittedBlockData > baseDB
// if hash is not nil,it will find from the chain, RecognizedBlockData > CommittedBlockData > baseDB
func (s *snapshotDB) Get(hash common.Hash, key []byte) ([]byte, error) {
	var parentHash common.Hash
	if hash == common.ZeroHash {
		//from unRecognizedBlockData
		if s.unRecognized == nil {
			return nil, errors.New("[SnapshotDB]unRecognized is not find now")
		}
		if v, err := s.unRecognized.data.Get(key); err == nil {
			return v, nil
		} else if err != memdb.ErrNotFound {
			return nil, err
		}
		parentHash = s.unRecognized.ParentHash
	} else {
		parentHash = hash
	}
	//from RecognizedBlockData
	for {
		if block, ok := s.recognized[parentHash]; ok {
			if v, err := block.data.Get(key); err == nil {
				return v, nil
			} else if err != memdb.ErrNotFound {
				return nil, err
			}
			parentHash = block.ParentHash
		} else {
			break
		}
	}

	//from committed
	if len(s.committed) > 0 {
		block := s.committed[len(s.committed)-1]
		if block.BlockHash != parentHash {
			return nil, ErrNotFound
		}
		for i := len(s.committed) - 1; i >= 0; i-- {
			if v, err := s.committed[i].data.Get(key); err == nil {
				return v, nil
			} else if err != memdb.ErrNotFound {
				return nil, err
			}
			continue
		}
	}

	//from baseDB
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
	if blocknumber.Int64() != s.unRecognized.Number.Int64() {
		return errors.New("[snapshotdb]blocknumber not compare the unRecognized blocknumber")
	}
	if _, ok := s.recognized[hash]; ok {
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
	s.recognized[hash] = *s.unRecognized

	s.unRecognized = nil
	return nil
}

// Commit move blockdata from recognized to commit
func (s *snapshotDB) Commit(hash common.Hash) error {
	s.commitLock.Lock()
	defer s.commitLock.Unlock()
	block, ok := s.recognized[hash]
	if !ok {
		return errors.New("[snapshotdb]not found form commit block:" + hash.String())

	}
	if s.current.HighestNum.Cmp(block.Number) >= 0 {
		return fmt.Errorf("[snapshotdb]the commit block num  %v is less than HighestNum %v", block.Number, s.current.HighestNum)
	}
	if (block.Number.Int64() - s.current.HighestNum.Int64()) != 1 {
		return fmt.Errorf("[snapshotdb]the commit block num %v - HighestNum %v should be eq 1", block.Number, s.current.HighestNum)
	}
	block.readOnly = true
	s.committed = append(s.committed, block)
	s.current.HighestNum = block.Number
	if err := s.current.update(); err != nil {
		return errors.New("[snapshotdb]update current fail:" + err.Error())
	}

	delete(s.recognized, hash)

	if err := s.rmOldRecognizedBlockData(); err != nil {
		return err
	}
	return nil
}

func (s *snapshotDB) BaseNum() (*big.Int, error) {
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
	if err := s.Close(); err != nil {
		return err
	}
	logger.Info(fmt.Sprint("begin clear file:", s.path))
	if err := os.RemoveAll(s.path); err != nil {
		return err
	}
	return nil
}

func itrToMdb(itr iterator.Iterator, mdb *memdb.DB) error {
	for itr.Next() {
		if err := mdb.Put(itr.Key(), itr.Value()); err != nil {
			return err
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
	var itrs []iterator.Iterator
	m := memdb.New(comparer.DefaultComparer, rangeNumber)
	prefix := util.BytesPrefix(key)
	var parentHash common.Hash
	if hash != common.ZeroHash {
		parentHash = hash
		location, ok := s.checkHashChain(parentHash)
		if !ok {
			return iterator.NewEmptyIterator(errors.New("this hash not in chain:" + parentHash.String()))
		}
		switch location {
		case hashLocationRecognized:
			for {
				if block, ok := s.recognized[parentHash]; ok {
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
					parentHash = block.BlockHash
				} else if block.ParentHash == parentHash {
					itrs = append(itrs, block.data.NewIterator(prefix))
					parentHash = block.BlockHash
				}
			}
		}
	} else {
		if s.unRecognized != nil {
			itrs = append(itrs, s.unRecognized.data.NewIterator(prefix))
			parentHash = s.unRecognized.ParentHash
		}
		for {
			if block, ok := s.recognized[parentHash]; ok {
				itrs = append(itrs, block.data.NewIterator(prefix))
				parentHash = block.ParentHash
			} else {
				break
			}
		}
		for _, block := range s.committed {
			itrs = append(itrs, block.data.NewIterator(prefix))
		}
	}
	itrs = append(itrs, s.baseDB.NewIterator(prefix, nil))
	for _, value := range itrs {
		if err := itrToMdb(value, m); err != nil {
			return iterator.NewEmptyIterator(err)
		}
	}
	return m.NewIterator(nil)
}

func (s *snapshotDB) Close() error {
	logger.Info("db is closing")
	//	runtime.SetFinalizer(s, nil)
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
