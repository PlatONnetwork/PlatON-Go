package snapshotdb

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"path"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbError "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/memdb"
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
	unCommitBlock := new(unCommitBlocks)
	unCommitBlock.blocks = make(map[common.Hash]*blockData)
	return &snapshotDB{
		path:          dbpath,
		storage:       stor,
		unCommit:      unCommitBlock,
		committed:     make([]*blockData, 0),
		journalw:      make(map[common.Hash]*journalWriter),
		baseDB:        baseDB,
		current:       newCurrent(dbpath),
		snapshotLockC: snapshotUnLock,
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
	block.Number = new(big.Int).SetUint64(fd.Num)
	block.data = memdb.New(DefaultComparer, 0)

	switch header.From {
	case journalHeaderFromUnRecognized:
		if fd.BlockHash == s.getUnRecognizedHash() {
			block.readOnly = false
		} else {
			block.readOnly = true
		}
	case journalHeaderFromRecognized:
		if fd.Num <= s.current.HighestNum.Uint64() {
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
		if err := block.data.Put(body.Key, body.Value); err != nil {
			return nil, err
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
	if _, corrupted := err.(*leveldbError.ErrCorrupted); corrupted {
		baseDB, err = leveldb.RecoverFile(getBaseDBPath(dbpath), nil)
		if err != nil {
			return fmt.Errorf("[SnapshotDB.recover]RecoverFile baseDB fail:%v", err)
		}
	}
	if err != nil {
		return fmt.Errorf("[SnapshotDB.recover]open baseDB fail:%v", err)
	}
	s.baseDB = baseDB

	//storage
	s.storage = stor
	fds, err := s.storage.List(TypeJournal)
	sortFds(fds)
	baseNum := s.current.BaseNum.Uint64()
	highestNum := s.current.HighestNum.Uint64()
	UnRecognizedHash := s.getUnRecognizedHash()
	s.committed = make([]*blockData, 0)
	s.journalw = make(map[common.Hash]*journalWriter)
	unCommitBlock := new(unCommitBlocks)
	unCommitBlock.blocks = make(map[common.Hash]*blockData)
	s.unCommit = unCommitBlock
	s.snapshotLockC = snapshotUnLock

	//read Journal
	for _, fd := range fds {
		block, err := s.getBlockFromJournal(fd)
		if err != nil {
			return err
		}
		if (baseNum < fd.Num && fd.Num <= highestNum) || (baseNum == 0 && highestNum == 0 && fd.Num == 0) {
			s.committed = append(s.committed, block)
		} else if fd.Num > highestNum {
			if UnRecognizedHash == fd.BlockHash {
				//1. UnRecognized
				s.unCommit.blocks[common.ZeroHash] = block
				//2. open writer
				w, err := s.storage.Append(fd)
				if err != nil {
					return fmt.Errorf("[SnapshotDB.recover]unRecognizedHash open storage fail:%v", err)
				}
				s.journalw[fd.BlockHash] = newJournalWriter(w)
			} else {
				//1. Recognized
				s.unCommit.blocks[fd.BlockHash] = block
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
		if fd.Num <= s.current.BaseNum.Uint64() {
			if _, ok := s.unCommit.blocks[fd.BlockHash]; ok {
				delete(s.unCommit.blocks, fd.BlockHash)
			}
			if err := s.closeJournalWriter(fd.BlockHash); err != nil {
				return err
			}
			if err := s.storage.Remove(fd); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *snapshotDB) rmOldRecognizedBlockData() error {
	var err error
	for key, value := range s.unCommit.blocks {
		if s.current.HighestNum.Int64() >= value.Number.Int64() {
			delete(s.unCommit.blocks, key)
			err = s.closeJournalWriter(key)
			if err != nil {
				return err
			}
			err := s.rmJournalFile(value.Number, key)
			if err != nil {
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
	return common.ZeroHash
}

func (s *snapshotDB) closeJournalWriter(hash common.Hash) error {
	s.journalWriterLock.Lock()
	defer s.journalWriterLock.Unlock()
	if j, ok := s.journalw[hash]; ok {
		if err := j.Close(); err != nil {
			return errors.New("[snapshotdb]close  journal writer fail:" + err.Error())
		}
		delete(s.journalw, hash)
	}
	return nil
}

const (
	hashLocationUnCommitted = 2
	hashLocationCommitted   = 3
	hashLocationNotFound    = 4
)

func (s *snapshotDB) checkHashChain(hash common.Hash) (int, bool) {
	var (
		lastBlockNumber = big.NewInt(0)
		lastParentHash  = hash
	)
	// find from unCommit
	for {
		if block, ok := s.unCommit.blocks[lastParentHash]; ok {
			if lastParentHash == block.ParentHash {
				logger.Error("loop error")
				return 0, false
			}
			lastParentHash = block.ParentHash
			lastBlockNumber = block.Number
		} else {
			break
		}
	}

	//check  last block find from unCommit is right
	if lastBlockNumber.Int64() > 0 {
		if s.current.HighestNum.Int64() != lastBlockNumber.Int64()-1 {
			logger.Error("[snapshotDB] find lastblock  fail ,num not compare", "current", s.current.HighestNum, "last", lastBlockNumber.Int64()-1)
			return 0, false
		}
		if s.current.HighestHash == common.ZeroHash {
			return hashLocationUnCommitted, true
		}
		if s.current.HighestHash != lastParentHash {
			logger.Error("[snapshotDB] find lastblock  fail ,hash not compare", "current", s.current.HighestHash.String(), "last", lastParentHash.String())
			return 0, false
		}
		return hashLocationUnCommitted, true
	}
	// if not find from unCommit, find from committed
	for _, value := range s.committed {
		if value.BlockHash == hash {
			return hashLocationCommitted, true
		}
	}
	return hashLocationNotFound, true
}

func (s *snapshotDB) put(hash common.Hash, key, value []byte) error {
	s.unCommit.Lock()
	defer s.unCommit.Unlock()

	block, ok := s.unCommit.blocks[hash]
	if !ok {
		return fmt.Errorf("not find the block by hash:%v", hash.String())
	}
	if block.readOnly {
		return errors.New("can't put read only block")
	}
	// TODO test
	log.Debug("old pposHash", "key", hex.EncodeToString(key), "val", hex.EncodeToString(value), "pposHash", block.kvHash.Hex())

	jData := journalData{
		Key:   key,
		Value: value,
		Hash:  s.generateKVHash(key, value, block.kvHash),
	}
	body, err := encode(jData)
	if err != nil {
		return errors.New("encode fail:" + err.Error())
	}
	if err := s.writeJournalBody(hash, body); err != nil {
		return errors.New("write journalBody fail:" + err.Error())
	}
	if err := block.data.Put(key, value); err != nil {
		return err
	}
	block.kvHash = jData.Hash
	// TODO test
	log.Debug("new pposHash", "key", hex.EncodeToString(key), "val", hex.EncodeToString(value), "pposHash", block.kvHash.Hex())
	return nil
}
