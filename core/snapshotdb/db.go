package snapshotdb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"io"
	"math/big"
	"path"
	"sync"
)

func getBaseDBPath(dbpath string) string {
	return path.Join(dbpath, DBBasePath)
}

func newDB(stor storage) (*snapshotDB, error) {
	dbpath := stor.Path()
	baseDB, err := leveldb.OpenFile(getBaseDBPath(dbpath), nil)
	if err != nil {
		return nil, fmt.Errorf("[SnapshotDB]open baseDB fail:%v", err)
	}
	mu := sync.Mutex{}
	return &snapshotDB{
		path:          dbpath,
		storage:       stor,
		unRecognized:  new(blockData),
		recognized:    make(map[common.Hash]blockData),
		committed:     make([]blockData, 0),
		journalw:      make(map[common.Hash]*journalWriter),
		baseDB:        baseDB,
		current:       newCurrent(dbpath),
		snapshotLock:  sync.NewCond(&mu),
		snapshotLockC: false,
	}, nil
}

func (s *snapshotDB) getBlockFromJournal(fd fileDesc) (*blockData, error) {
	reader, err := s.storage.Open(fd)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	journals := journal.NewReader(reader, nil, false, false)
	j, err := journals.Next()
	if err != nil {
		return nil, err
	}
	var header journalHeader
	if err := decode(j, &header); err != nil {
		return nil, err
	}
	block := new(blockData)
	block.ParentHash = header.ParentHash
	if fd.BlockHash != s.getUnRecognizedHash() {
		block.BlockHash = fd.BlockHash
	}
	block.Number = big.NewInt(fd.Num)
	block.data = memdb.New(DefaultComparer, 0)

	switch header.From {
	case journalHeaderFromUnRecognized:
		if fd.BlockHash == s.getUnRecognizedHash() {
			block.readOnly = false
		} else {
			block.readOnly = true
		}
	case journalHeaderFromRecognized:
		if fd.Num <= s.current.HighestNum.Int64() {
			block.readOnly = true
		}
	}
	var kvhash common.Hash
	for {
		j, err := journals.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var body journalData
		if err := decode(j, &body); err != nil {
			return nil, err
		}
		switch body.FuncType {
		case funcTypePut:
			if err := block.data.Put(body.Key, body.Value); err != nil {
				return nil, err
			}
		case funcTypeDel:
			if err := block.data.Delete(body.Key); err != nil {
				return nil, err
			}
		}
		kvhash = body.Hash
	}
	block.kvHash = kvhash
	return block, nil
}

func (s *snapshotDB) recover(stor storage) error {
	dbpath := stor.Path()
	c, err := loadCurrent(dbpath)
	if err != nil {
		return fmt.Errorf("[SnapshotDB.recover]load  current fail:%v", err)
	}
	s.path = dbpath
	s.current = c

	//baseDB
	baseDB, err := leveldb.OpenFile(getBaseDBPath(dbpath), nil)
	if err != nil {
		return fmt.Errorf("[SnapshotDB.recover]open baseDB fail:%v", err)
	}
	s.baseDB = baseDB

	//storage
	s.storage = stor
	fds, err := s.storage.List(TypeJournal)
	sortFds(fds)
	baseNum := s.current.BaseNum.Int64()
	highestNum := s.current.HighestNum.Int64()
	UnRecognizedHash := s.getUnRecognizedHash()
	s.committed = make([]blockData, 0)
	s.recognized = make(map[common.Hash]blockData)
	s.journalw = make(map[common.Hash]*journalWriter)

	mu := sync.Mutex{}
	s.snapshotLock = sync.NewCond(&mu)
	s.snapshotLockC = false

	//read Journal
	for _, fd := range fds {
		block, err := s.getBlockFromJournal(fd)
		if err != nil {
			return err
		}
		if baseNum < fd.Num && fd.Num <= highestNum {
			s.committed = append(s.committed, *block)
		} else if fd.Num > highestNum {
			if UnRecognizedHash == fd.BlockHash {
				//1. UnRecognized
				s.unRecognized = block
				//2. open writer
				w, err := s.storage.Append(fd)
				if err != nil {
					return fmt.Errorf("[SnapshotDB.recover]unRecognizedHash open storage fail:%v", err)
				}
				s.journalw[fd.BlockHash] = newJournalWriter(w)
			} else {
				//1. Recognized
				s.recognized[fd.BlockHash] = *block
				//2. open writer
				if !block.readOnly {
					w, err := s.storage.Append(fd)
					if err != nil {
						return fmt.Errorf("[SnapshotDB.recover]recognized open storage fail:%v", err)
					}
					s.journalw[fd.BlockHash] = newJournalWriter(w)
				}

			}

		}
	}
	return nil
}

func (s *snapshotDB) removeJournalLessThanBaseNum() error {
	fds, err := s.storage.List(TypeJournal)
	if err != nil {
		return err
	}
	for _, fd := range fds {
		if fd.Num <= s.current.BaseNum.Int64() {
			if err := s.storage.Remove(fd); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *snapshotDB) generateKVHash(k, v []byte, hash common.Hash) common.Hash {
	var buf bytes.Buffer
	buf.Write(k)
	buf.Write(v)
	buf.Write(hash.Bytes())
	return rlpHash(buf.Bytes())
}

func (s *snapshotDB) getUnRecognizedHash() common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte("CURRENT"))
	return rlpHash(buf.Bytes())
}

func (s *snapshotDB) closeJournalWriter(hash common.Hash) error {
	if j, ok := s.journalw[hash]; ok {
		if err := j.Close(); err != nil {
			return errors.New("[snapshotdb]close  journal writer fail:" + err.Error())
		}
		delete(s.journalw, hash)
	}
	return nil
}

func (s *snapshotDB) rmOldRecognizedBlockData() error {
	for key, value := range s.recognized {
		if s.current.HighestNum.Cmp(value.Number) >= 0 {
			delete(s.recognized, key)
			if err := s.closeJournalWriter(key); err != nil {
				return err
			}
			if err := s.rmJournalFile(value.Number, key); err != nil {
				return err
			}
		}
	}
	return nil
}

const (
	hashLocationUnRecognized = 1
	hashLocationRecognized   = 2
	hashLocationCommitted    = 3
	hashLocationNotFound     = 4
)

func (s *snapshotDB) checkHashChain(hash common.Hash) (int, bool) {
	lastblockNumber := big.NewInt(0)
	lastParenthash := common.ZeroHash
	if hash == common.ZeroHash {
		if s.unRecognized == nil {
			return 0, false
		}
		lastParenthash = s.unRecognized.ParentHash
		lastblockNumber = s.unRecognized.Number
		if s.current.HighestNum.Int64() >= s.unRecognized.Number.Int64() {
			return 0, false
		}
		for {
			if data, ok := s.recognized[lastParenthash]; ok {
				if lastParenthash == data.ParentHash {
					logger.Error("loop error")
					return 0, false
				}
				lastblockNumber = data.Number
				lastParenthash = data.ParentHash
			} else {
				break
			}
		}
		if lastblockNumber.Int64() > 0 {
			if s.current.LastHash == common.ZeroHash {
				return hashLocationUnRecognized, true
			}
			if s.current.HighestNum.Int64() != lastblockNumber.Int64()-1 {
				return 0, false
			}

			if s.current.LastHash != lastParenthash {
				return 0, false
			}
			return hashLocationUnRecognized, true
		}
		return 0, false
	}
	{
		// find from recognized
		lastParenthash = hash
		for {
			if data, ok := s.recognized[lastParenthash]; ok {
				if lastParenthash == data.ParentHash {
					logger.Error("loop error")
					return 0, false
				}
				lastParenthash = data.ParentHash
				lastblockNumber = data.Number
			} else {
				break
			}
		}

		//check find from recognized is right
		if lastblockNumber.Int64() > 0 {
			if s.current.HighestNum.Int64() != lastblockNumber.Int64()-1 {
				return 0, false
			}
			if s.current.LastHash == common.ZeroHash {
				return hashLocationRecognized, true
			}
			if s.current.LastHash != lastParenthash {
				return 0, false
			}
			return hashLocationRecognized, true
		}
	}

	// find from committed
	for _, value := range s.committed {
		if value.BlockHash == hash {
			return hashLocationCommitted, true
		}
	}

	return hashLocationNotFound, true
}

func (s *snapshotDB) put(hash common.Hash, key, value []byte, funcType uint64) error {
	var (
		blockHash  common.Hash
		kvhash     common.Hash
		recognized blockData
	)
	if hash == common.ZeroHash {
		s.unRecognizedLock.Lock()
		defer s.unRecognizedLock.Unlock()
		if s.unRecognized == nil {
			return errors.New("[SnapshotDB]can't put to unRecognized,it was nil")
		}
		if s.unRecognized.readOnly {
			return errors.New("[SnapshotDB]can't put read only block")
		}
		blockHash = s.getUnRecognizedHash()
		kvhash = s.unRecognized.kvHash
	} else {
		bb, ok := s.recognized[hash]
		if !ok {
			return errors.New("[SnapshotDB]get recognized block data by hash fail")
		}
		blockHash = hash
		if bb.readOnly {
			return errors.New("[SnapshotDB]can't put read only block")
		}
		recognized = bb
		kvhash = recognized.kvHash
	}

	jData := journalData{
		Key:      key,
		Value:    value,
		Hash:     s.generateKVHash(key, value, kvhash),
		FuncType: funcType,
	}
	body, err := encode(jData)
	if err != nil {
		return errors.New("[SnapshotDB]encode fail:" + err.Error())
	}
	if err := s.writeJournalBody(blockHash, body); err != nil {
		return errors.New("[SnapshotDB]write journalBody fail:" + err.Error())
	}
	if hash != common.ZeroHash {
		switch funcType {
		case funcTypePut:
			if err := recognized.data.Put(key, value); err != nil {
				return err
			}
		case funcTypeDel:
			if err := recognized.data.Delete(key); err != nil {
				return err
			}
		}
		recognized.kvHash = jData.Hash
		s.recognized[hash] = recognized
	} else {
		switch funcType {
		case funcTypePut:
			if err := s.unRecognized.data.Put(key, value); err != nil {
				return err
			}
		case funcTypeDel:
			if err := s.unRecognized.data.Delete(key); err != nil {
				return err
			}
		}
		s.unRecognized.kvHash = jData.Hash
	}
	return nil
}
