package snapshotdb

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
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
	CURRENT = "CURRENT"
	KVLIMIT = 2000
)

type DB interface {
	Put(hash *common.Hash, key, value []byte) (bool, error)
	NewBlock(blockNumber *big.Int, parentHash common.Hash, hash *common.Hash) (bool, error)
	Get(hash *common.Hash, key []byte) ([]byte, error)
	GetCommitedBlock(key []byte) ([]byte, error)
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

var dbInstance *SnapshotDB

var (
	ErrorSnaphotLock = errors.New("can't create snapshot,snapshot is lock now")
)

//需要初始化一些纸
type SnapshotDB struct {
	path         string
	mu           sync.RWMutex
	snapshotLock bool
	current      *current
	baseDB       *leveldb.DB
	unRecognized *blockData
	recognized   map[common.Hash]blockData
	commited     []blockData
	journalw     map[common.Hash]*journal.Writer
	storage      Storage
}

func Instance() DB {
	return dbInstance
}

func (s *SnapshotDB) GetCommitedBlock(key []byte) ([]byte, error) {
	v, err := s.getFromCommited(nil, key)
	if err == nil {
		return v, nil
	}
	if err != memdb.ErrNotFound {
		return nil, err
	}
	v, err = s.getFromBaseDB(key)
	if err == nil {
		return v, nil
	}
	return nil, err
}

func (s *SnapshotDB) PutBaseDB(key, value []byte) (bool, error) {
	err := s.baseDB.Put(key, value, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SnapshotDB) GetLastKVHash(blockHash *common.Hash) []byte {
	if blockHash == nil {
		return s.unRecognized.kvHash.Bytes()
	}
	block := s.recognized[*blockHash]
	return block.kvHash.Bytes()
}

func (s *SnapshotDB) Del(hash *common.Hash, key []byte) (bool, error) {
	if err := s.put(hash, key, nil, funcTypeDel); err != nil {
		return false, err
	}
	return true, nil
}

func (s *SnapshotDB) Compaction() (bool, error) {
	s.mu.Lock()
	s.snapshotLock = true
	defer func() {
		s.snapshotLock = false
		s.mu.Unlock()
	}()
	var size int
	var writeBlockNum int
	for i := 0; i < len(s.commited); i++ {
		if i < 10 {
			if size > KVLIMIT {
				writeBlockNum = i - 1
				break
			}
			size += s.commited[i].data.Size()
		} else {
			writeBlockNum = 9
			break
		}
	}
	batch := new(leveldb.Batch)
	for i := 0; i <= writeBlockNum; i++ {
		itr := s.commited[i].data.NewIterator(nil)
		for itr.Next() {
			batch.Put(itr.Key(), itr.Value())
		}
	}
	if err := s.baseDB.Write(batch, nil); err != nil {
		return false, errors.New("[SnapshotDB]write to baseDB fail:" + err.Error())
	}
	s.current.BaseNum.Add(s.current.BaseNum, big.NewInt(int64(writeBlockNum)))
	if err := s.current.update(); err != nil {
		return false, err
	}
	s.commited = s.commited[writeBlockNum:len(s.commited)]
	if err := s.removeJournalLessThanBaseNum(); err != nil {
		return false, err
	}
	return true, nil
}

//if hash nil ,new unRecognized data
//if hash not nul,new Recognized data
//todo parentHash是否在chain上，需要校验
func (s *SnapshotDB) NewBlock(blockNumber *big.Int, parentHash common.Hash, hash *common.Hash) (bool, error) {
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

func (s *SnapshotDB) put(hash *common.Hash, key, value []byte, funcType uint64) error {
	var (
		block     *blockData
		blockHash common.Hash
	)
	if hash == nil {
		block = s.unRecognized
		blockHash = s.getUnRecognizedHash()
	} else {
		bb, ok := s.recognized[*hash]
		if !ok {
			return errors.New("[SnapshotDB]get recognized block data by hash fail")
		}
		block = &bb
		blockHash = *hash
	}

	if block.readOnly {
		return errors.New("[SnapshotDB]can't put read only block")
	}

	jData := journalData{
		Key:      key,
		Value:    value,
		Hash:     s.generateKVHash(key, value, block.kvHash),
		FuncType: funcType,
	}
	body, err := encode(jData)
	if err != nil {
		return errors.New("encode fail:" + err.Error())
	}
	if err := s.writeJournalBody(blockHash, body); err != nil {
		return err
	}
	switch funcType {
	case funcTypePut:
		if err := block.data.Put(key, value); err != nil {
			return err
		}
	case funcTypeDel:
		if err := block.data.Delete(key); err != nil {
			return err
		}
	}
	block.kvHash = jData.Hash
	return nil
}

func (s *SnapshotDB) Put(hash *common.Hash, key, value []byte) (bool, error) {
	if err := s.put(hash, key, value, funcTypePut); err != nil {
		return false, err
	}
	return true, nil
}

func (s *SnapshotDB) Get(hash *common.Hash, key []byte) ([]byte, error) {
	if hash == nil {
		v, err := s.getFromUnRecognized(key)
		if err == nil {
			return v, nil
		}
		if err != memdb.ErrNotFound {
			return nil, err
		}
	}
	v, err := s.getFromRecognized(hash, key)
	if err == nil {
		return v, nil
	}
	if err != memdb.ErrNotFound {
		return nil, err
	}
	v, err = s.getFromCommited(hash, key)
	if err == nil {
		return v, nil
	}
	if err != memdb.ErrNotFound {
		return nil, err
	}
	v, err = s.getFromBaseDB(key)
	if err == nil {
		return v, nil
	}
	return nil, err
}

func (s *SnapshotDB) Has(hash *common.Hash, key []byte) (bool, error) {
	if hash == nil {
		_, err := s.getFromUnRecognized(key)
		if err == nil {
			return true, nil
		}
		if err != memdb.ErrNotFound {
			return false, err
		}
	}
	_, err := s.getFromRecognized(hash, key)
	if err == nil {
		return true, nil
	}
	if err != memdb.ErrNotFound {
		return false, err
	}
	_, err = s.getFromCommited(hash, key)
	if err == nil {
		return true, nil
	}
	if err != memdb.ErrNotFound {
		return false, err
	}
	return s.baseDB.Has(key, nil)
}

//move unRecognized to Recognized data
func (s *SnapshotDB) Flush(hash common.Hash, blocknumber *big.Int) (bool, error) {
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

func (s *SnapshotDB) Commit(hash common.Hash) (bool, error) {
	block, ok := s.recognized[hash]
	if !ok {
		return false, errors.New("[snapshotdb]not found form commit block:" + hash.String())

	}
	if s.current.HighestNum.Cmp(block.Number) >= 0 {
		return false, fmt.Errorf("[snapshotdb]should blockNum %v >= HighestNum %v", block.Number, s.current.HighestNum)
	}
	if (block.Number.Int64() - s.current.HighestNum.Int64()) != 1 {
		return false, fmt.Errorf("[snapshotdb]blockNum %v - HighestNum %v should be eq 1", block.Number, s.current.HighestNum)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
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

func (s *SnapshotDB) BaseNum() (*big.Int, error) {
	return s.current.BaseNum, nil
}

func (s *SnapshotDB) WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
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
func (s *SnapshotDB) Clear() (bool, error) {
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

const (
	HashLocationRecognized = 1
	HashLocationCommited   = 2
)

func (s *SnapshotDB) checkHashChain(hash common.Hash) (int, bool) {
	lastblockNumber := big.NewInt(0)
	//在recognized中找
	for {
		if data, ok := s.recognized[hash]; ok {
			hash = data.ParentHash
			lastblockNumber = data.Number
		} else {
			break
		}
	}
	//在recognized中找到
	if lastblockNumber.Int64() > 0 {
		if len(s.commited) > 0 {
			commitBlock := s.commited[len(s.commited)-1]
			if lastblockNumber.Int64()-1 != commitBlock.Number.Int64() {
				return 0, false
			}
			if commitBlock.BlockHash.String() != hash.String() {
				return 0, false
			}
			return HashLocationRecognized, true
		} else {
			if s.current.HighestNum.Int64() == lastblockNumber.Int64()-1 {
				return HashLocationRecognized, true
			}
		}
	}
	//在recognized中没有找到,在commit中找
	for _, value := range s.commited {
		if *value.BlockHash == hash {
			return HashLocationCommited, true
		}
	}
	return 0, false
}

//1.hash为空的时候，从unRecognized开始查（假设unRecognized的parentHash必为真），如果unRecognized也为空，从commited开始查
//2.hash不为空时，从Recognized开始逐级网上查
func (s *SnapshotDB) Ranking(hash *common.Hash, key []byte, rangeNumber int) iterator.Iterator {
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
		case HashLocationRecognized:
			for {
				if block, ok := s.recognized[parentHash]; ok {
					itrs = append(itrs, block.data.NewIterator(prefix))
					parentHash = s.unRecognized.ParentHash
				} else {
					break
				}
			}
			for _, block := range s.commited {
				itrs = append(itrs, block.data.NewIterator(prefix))
			}
		case HashLocationCommited:
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
		} else {
			for _, block := range s.commited {
				itrs = append(itrs, block.data.NewIterator(prefix))
			}
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

func (s *SnapshotDB) Close() (bool, error) {
	if err := s.baseDB.Close(); err != nil {
		return false, fmt.Errorf("[snapshotdb]close base db fail:%v", err)
	}
	if err := s.storage.Close(); err != nil {
		return false, fmt.Errorf("[snapshotdb]close storage fail:%v", err)
	}
	if s.current != nil {
		s.current.f.Close()
	}

	for key, _ := range s.journalw {
		if err := s.journalw[key].Close(); err != nil {
			return false, fmt.Errorf("[snapshotdb]close journalw fail:%v", err)
		}
	}
	return true, nil
}
