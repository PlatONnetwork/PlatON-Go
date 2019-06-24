package snapshotdb

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/robfig/cron"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"math/big"
	"os"
	"sync"
)

const (
	kvLIMIT = 2000
)

//DB the main snapshotdb interface
type DB interface {
	Put(hash *common.Hash, key, value []byte) (bool, error)
	NewBlock(blockNumber *big.Int, parentHash common.Hash, hash *common.Hash) (bool, error)
	Get(hash *common.Hash, key []byte) ([]byte, error)
	GetFromCommitedBlock(key []byte) ([]byte, error)
	Del(hash *common.Hash, key []byte) (bool, error)
	Has(hash *common.Hash, key []byte) (bool, error)
	Flush(hash common.Hash, blocknumber *big.Int) (bool, error)
	Ranking(hash *common.Hash, key []byte, ranges int) iterator.Iterator
	WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error
	Commit(hash common.Hash) (bool, error)
	Clear() (bool, error)
	PutBaseDB(key, value []byte) (bool, error)
	GetLastKVHash(blockHash *common.Hash) []byte
	BaseNum() (*big.Int, error)
	Close() (bool, error)
	Compaction() (bool, error)
}

var dbInstance *snapshotDB

var (
	//ErrorSnaphotLock when db is Lock
	ErrorSnaphotLock = errors.New("can't create snapshot,snapshot is lock now")
	ErrNotFound      = errors.New("snapshotDB: not found")
)

type snapshotDB struct {
	path         string
	mu           sync.RWMutex
	snapshotLock bool
	current      *current
	baseDB       *leveldb.DB
	unRecognized *blockData
	recognized   map[common.Hash]blockData

	commited   []blockData
	commitLock sync.RWMutex

	journalw map[common.Hash]*journal.Writer
	storage  storage

	corn *cron.Cron

	closed bool
}

//Instance return the Instance of the db
func Instance() DB {
	if dbInstance == nil {
		initDB()
	}
	if dbInstance.closed {
		initDB()
	}
	return dbInstance
}

//GetCommitedBlock    get value from commited blockdata > baseDB
func (s *snapshotDB) GetFromCommitedBlock(key []byte) ([]byte, error) {
	for _, value := range s.commited {
		if v, err := value.data.Get(key); err == nil {
			return v, nil
		} else if err != memdb.ErrNotFound {
			return nil, err
		}
	}
	v, err := s.baseDB.Get(key, nil)
	if err == nil {
		return v, nil
	} else if err == leveldb.ErrNotFound {
		return nil, ErrNotFound
	}
	return nil, err
}

func (s *snapshotDB) PutBaseDB(key, value []byte) (bool, error) {
	err := s.baseDB.Put(key, value, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

//if hash is nil ,get unRecognized block lastkv hash,
//else, get recognized block lastkv  hash
func (s *snapshotDB) GetLastKVHash(blockHash *common.Hash) []byte {
	if blockHash == nil {
		return s.unRecognized.kvHash.Bytes()
	}
	block, ok := s.recognized[*blockHash]
	if !ok {
		return nil
	}
	return block.kvHash.Bytes()
}

func (s *snapshotDB) Del(hash *common.Hash, key []byte) (bool, error) {
	if err := s.put(hash, key, nil, funcTypeDel); err != nil {
		return false, err
	}
	return true, nil
}

//Compaction ,write commit to baseDB,and then removeJournal lessThan BaseNum
// it will write to baseDB
// case kv>2000 and block == 1
// case kv<2000,block... <9
// case kv<2000,block...=9
func (s *snapshotDB) Compaction() (bool, error) {
	s.commitLock.Lock()
	s.snapshotLock = true
	defer func() {
		s.snapshotLock = false
		s.commitLock.Unlock()
	}()
	var (
		kvsize    int
		commitNum int
	)
	for i := 0; i < len(s.commited); i++ {
		if i < 10 {
			if kvsize > kvLIMIT {
				commitNum = i - 1
				break
			}
			kvsize += s.commited[i].data.Len()
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
		itr := s.commited[i].data.NewIterator(nil)
		for itr.Next() {
			batch.Put(itr.Key(), itr.Value())
		}
	}
	if err := s.baseDB.Write(batch, nil); err != nil {
		return false, errors.New("[SnapshotDB]write to baseDB fail:" + err.Error())
	}
	s.current.BaseNum.Add(s.current.BaseNum, big.NewInt(int64(commitNum)))
	if err := s.current.update(); err != nil {
		return false, err
	}
	s.commited = s.commited[commitNum:len(s.commited)]
	if err := s.removeJournalLessThanBaseNum(); err != nil {
		return false, err
	}
	return true, nil
}

//NewBlock call when you need a new unRecognized or recognized  block data
//it will set JournalHeader for the block
//if hash nil ,new unRecognized data
//if hash not nul,new Recognized data
func (s *snapshotDB) NewBlock(blockNumber *big.Int, parentHash common.Hash, hash *common.Hash) (bool, error) {
	if hash == nil {
		if s.unRecognized != nil && s.unRecognized.readOnly {
			return false, errors.New("[SnapshotDB]can't  new unRecognized block,it's readonly now")
		}
	}
	block := new(blockData)
	block.Number = blockNumber
	block.ParentHash = parentHash
	block.BlockHash = hash
	block.data = memdb.New(DefaultComparer, 100)
	if hash == nil {
		if err := s.writeJournalHeader(blockNumber, s.getUnRecognizedHash(), parentHash, journalHeaderFromUnRecognized); err != nil {
			return false, fmt.Errorf("[SnapshotDB] write Journal Header fail:%v", err)
		}
		s.unRecognized = block
	} else {
		if err := s.writeJournalHeader(blockNumber, *hash, parentHash, journalHeaderFromRecognized); err != nil {
			return false, fmt.Errorf("[SnapshotDB] write Journal body fail:%v", err)
		}
		s.recognized[*hash] = *block
	}
	return true, nil
}

func (s *snapshotDB) Put(hash *common.Hash, key, value []byte) (bool, error) {
	if err := s.put(hash, key, value, funcTypePut); err != nil {
		return false, err
	}
	return true, nil
}

//Get get key,val from  snapshotDB
// if hash is nil, unRecognizedBlockData > RecognizedBlockData > CommitedBlockData > baseDB
// if hash is not nil,it will find from the chain, RecognizedBlockData > CommitedBlockData > baseDB
func (s *snapshotDB) Get(hash *common.Hash, key []byte) ([]byte, error) {
	var parentHash common.Hash
	if hash == nil {
		//from unRecognizedBlockData
		if s.unRecognized == nil {
			return nil, errors.New("unRecognized is not find now")
		}
		if v, err := s.unRecognized.data.Get(key); err == nil {
			return v, nil
		} else if err != memdb.ErrNotFound {
			return nil, err
		}
		parentHash = s.unRecognized.ParentHash
	} else {
		parentHash = *hash
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

	//from commited
	if len(s.commited) > 0 {
		block := s.commited[len(s.commited)-1]
		if *block.BlockHash != parentHash {
			return nil, ErrNotFound
		}
		for i := len(s.commited) - 1; i >= 0; i-- {
			if v, err := s.commited[i].data.Get(key); err == nil {
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
func (s *snapshotDB) Has(hash *common.Hash, key []byte) (bool, error) {
	_, err := s.Get(hash, key)
	if err == nil {
		return true, nil
	} else if err == ErrNotFound {
		return true, ErrNotFound
	} else {
		return false, err
	}
}

//move unRecognized to Recognized data
func (s *snapshotDB) Flush(hash common.Hash, blocknumber *big.Int) (bool, error) {
	if blocknumber.Int64() != s.unRecognized.Number.Int64() {
		return false, errors.New("[snapshotdb]blocknumber not compare the unRecognized blocknumber")
	}
	if _, ok := s.recognized[hash]; ok {
		return false, errors.New("the hash is exist in recognized data")
	}
	currentHash := s.getUnRecognizedHash()
	oldFd := fileDesc{Type: TypeJournal, Num: blocknumber.Int64(), BlockHash: currentHash}
	newFd := fileDesc{Type: TypeJournal, Num: blocknumber.Int64(), BlockHash: hash}
	s.unRecognized.mu.Lock()
	if err := s.storage.Rename(oldFd, newFd); err != nil {
		s.unRecognized.mu.Unlock()
		return false, errors.New("[snapshotdb]rename fiel fail:" + oldFd.String() + "," + newFd.String() + "," + err.Error())
	}
	s.unRecognized.BlockHash = &hash
	s.unRecognized.readOnly = true
	s.recognized[hash] = *s.unRecognized
	if err := s.closeJournalWriter(currentHash); err != nil {
		s.unRecognized.mu.Unlock()
		return false, err
	}
	s.unRecognized.mu.Unlock()
	s.unRecognized = nil
	return true, nil
}

func (s *snapshotDB) Commit(hash common.Hash) (bool, error) {
	s.commitLock.Lock()
	defer s.commitLock.Unlock()
	block, ok := s.recognized[hash]
	if !ok {
		return false, errors.New("[snapshotdb]not found form commit block:" + hash.String())

	}
	if s.current.HighestNum.Cmp(block.Number) >= 0 {
		return false, fmt.Errorf("[snapshotdb]the commit block num  %v is less than HighestNum %v", block.Number, s.current.HighestNum)
	}
	if (block.Number.Int64() - s.current.HighestNum.Int64()) != 1 {
		return false, fmt.Errorf("[snapshotdb]the commit block num %v - HighestNum %v should be eq 1", block.Number, s.current.HighestNum)
	}
	block.readOnly = true
	s.commited = append(s.commited, block)
	s.current.HighestNum = block.Number
	if err := s.current.update(); err != nil {
		return false, errors.New("[snapshotdb]update current fail:" + err.Error())
	}

	delete(s.recognized, hash)

	if err := s.rmOldRecognizedBlockData(); err != nil {
		return false, err
	}
	return true, nil
}

func (s *snapshotDB) BaseNum() (*big.Int, error) {
	return s.current.BaseNum, nil
}

func (s *snapshotDB) WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
	if s.snapshotLock {
		return errors.New("[snapshotdb] snapshot is lock now,can't get")
	}
	snapshot, err := s.baseDB.GetSnapshot()
	if err != nil {
		return errors.New("[snapshotdb] get snapshot fail:" + err.Error())
	}
	t := snapshot.NewIterator(slice, nil)
	return f(s.current.BaseNum, t)
}

//todo 需要确认clear的执行流程
func (s *snapshotDB) Clear() (bool, error) {
	if _, err := s.Close(); err != nil {
		return false, err
	}
	log.Print("begin remove", s.path)
	if err := os.RemoveAll(s.path); err != nil {
		return false, err
	}
	return true, nil
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

//1.hash为空的时候，从unRecognized开始查（假设unRecognized的parentHash必为真），如果unRecognized也为空，从commited开始查
//2.hash不为空时，从Recognized开始逐级往上查
func (s *snapshotDB) Ranking(hash *common.Hash, key []byte, rangeNumber int) iterator.Iterator {
	var itrs []iterator.Iterator
	m := memdb.New(comparer.DefaultComparer, rangeNumber)
	prefix := util.BytesPrefix(key)
	var parentHash common.Hash
	if hash != nil {
		parentHash = *hash
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
			for _, block := range s.commited {
				itrs = append(itrs, block.data.NewIterator(prefix))
			}
		case hashLocationCommited:
			for i := len(s.commited) - 1; i >= 0; i-- {
				block := s.commited[i]
				if block.BlockHash == hash {
					itrs = append(itrs, block.data.NewIterator(prefix))
					parentHash = *block.BlockHash
				} else if block.ParentHash == parentHash {
					itrs = append(itrs, block.data.NewIterator(prefix))
					parentHash = *block.BlockHash
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
		for _, block := range s.commited {
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

func (s *snapshotDB) Close() (bool, error) {
	s.corn.Stop()
	if err := s.baseDB.Close(); err != nil {
		return false, fmt.Errorf("[snapshotdb]close base db fail:%v", err)
	}
	if err := s.storage.Close(); err != nil {
		return false, fmt.Errorf("[snapshotdb]close storage fail:%v", err)
	}
	if s.current != nil {
		s.current.f.Close()
	}

	for key := range s.journalw {
		if err := s.journalw[key].Close(); err != nil {
			return false, fmt.Errorf("[snapshotdb]close journalw fail:%v", err)
		}
	}
	s.closed = true
	return true, nil
}
