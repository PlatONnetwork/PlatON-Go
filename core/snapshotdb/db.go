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
	"log"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func getBaseDBPath(dbpath string) string {
	return path.Join(dbpath, DBBasePath)
}

func newDB(dbpath string) (*SnapshotDB, error) {
	s, err := OpenFile(dbpath, false)
	if err != nil {
		return nil, fmt.Errorf("[SnapshotDB]open db dir fail:%v", err)
	}
	baseDB, err := leveldb.OpenFile(getBaseDBPath(dbpath), nil)
	if err != nil {
		return nil, fmt.Errorf("[SnapshotDB]open baseDB fail:%v", err)
	}
	return &SnapshotDB{
		path:         dbpath,
		storage:      s,
		unRecognized: new(blockData),
		recognized:   make(map[common.Hash]blockData),
		commited:     make([]blockData, 0),
		journalw:     make(map[common.Hash]*journal.Writer),
		baseDB:       baseDB,
		current:      newCurrent(dbpath),
	}, nil
}

func (s *SnapshotDB) findJournalFile() []string {
	matchs, _ := filepath.Glob(path.Join(s.path, "*.log"))
	return matchs
}

func (s *SnapshotDB) getBlockFromJournal(fd fileDesc) (*blockData, error) {
	reader, err := s.storage.Open(fd)
	if err != nil {
		return nil, err
	}
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
	block.BlockHash = &fd.BlockHash
	block.Number = big.NewInt(fd.Num)
	block.data = memdb.New(DefaultComparer, 0)

	switch header.From {
	case journalHeaderFromUnRecognized:
		if fd.BlockHash == s.getUnRecognizedHash() {
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
			//
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

//todo 统一current 和 log
func (s *SnapshotDB) recover(dbpath string) error {
	c, err := loadCurrent(dbpath)
	if err != nil {
		return err
	}
	s.path = dbpath
	s.current = c

	//baseDB
	baseDB, err := leveldb.OpenFile(getBaseDBPath(dbpath), nil)
	if err != nil {
		return fmt.Errorf("[SnapshotDB]open baseDB fail:%v", err)
	}
	s.baseDB = baseDB

	//storage
	storage, err := OpenFile(dbpath, false)
	if err != nil {
		return fmt.Errorf("[SnapshotDB]open db dir fail:%v", err)
	}
	s.storage = storage

	//find Journal
	matchs := s.findJournalFile()
	var fds fileDescs
	fds = make([]fileDesc, 0)
	for _, value := range matchs {
		if fd, err := fsParseName(value); err != nil {
			return err
		} else {
			fds = append(fds, fd)
		}
	}
	sort.Sort(fds)

	//初始化一些变量
	baseNum := s.current.BaseNum.Int64()
	highestNum := s.current.HighestNum.Int64()
	UnRecognizedHash := s.getUnRecognizedHash()
	s.commited = make([]blockData, 0)
	s.recognized = make(map[common.Hash]blockData)
	s.journalw = make(map[common.Hash]*journal.Writer)

	//read Journal
	for _, fd := range fds {
		//从日志读到内存
		block, err := s.getBlockFromJournal(fd)
		if err != nil {
			return err
		}
		if baseNum < fd.Num && fd.Num <= highestNum {
			s.commited = append(s.commited, *block)
		} else if fd.Num > highestNum {
			if UnRecognizedHash == fd.BlockHash {
				//UnRecognized
				s.unRecognized = block
				//2.打开writer
				w, err := s.storage.Append(fd)
				if err != nil {
					return err
				}
				s.journalw[fd.BlockHash] = journal.NewWriter(w)
			} else {
				//Recognized
				s.recognized[fd.BlockHash] = *block
				//2.根据情况writer
				if !block.readOnly {
					w, err := s.storage.Append(fd)
					if err != nil {
						return err
					}
					s.journalw[fd.BlockHash] = journal.NewWriter(w)
				}

			}

		}
	}
	return nil
}

func (s *SnapshotDB) removeJournalLessThanBaseNum() error {
	m, err := filepath.Glob(filepath.Join(s.path, "*.log"))
	if err != nil {
		return err
	}
	for _, value := range m {
		ss := strings.Split(value, "-")
		tmp2 := strings.Split(ss[0], "/")
		i, err := strconv.ParseInt(tmp2[1], 10, 64)
		if err != nil {
			return err
		}
		if i <= s.current.BaseNum.Int64() {
			if err := os.Remove(value); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

func (s *SnapshotDB) schedule() {
	if counter.get() == 60 || s.current.HighestNum.Int64()-s.current.BaseNum.Int64() >= 100 {
		if _, err := s.Compaction(); err != nil {
			log.Print("[SnapshotDB]compaction fail:", err)
		}
		counter.reset()
	} else {
		counter.increment()
	}
}

func (s *SnapshotDB) generateKVHash(k, v []byte, hash common.Hash) common.Hash {
	var buf bytes.Buffer
	buf.Write(k)
	buf.Write(v)
	buf.Write(hash.Bytes())
	return rlpHash(buf.Bytes())
}

func (s *SnapshotDB) getFromUnRecognized(key []byte) ([]byte, error) {
	return s.unRecognized.data.Get(key)
}

func (s *SnapshotDB) getFromRecognized(hash *common.Hash, key []byte) ([]byte, error) {
	if hash == nil {
		for _, value := range s.recognized {
			v, err := value.data.Get(key)
			if err == nil {
				return v, nil
			}
			if err != memdb.ErrNotFound {
				return nil, err
			}
		}
	} else {
		b, ok := s.recognized[*hash]
		if ok {
			return b.data.Get(key)
		}
	}
	return nil, memdb.ErrNotFound
}

func (s *SnapshotDB) getFromCommited(hash *common.Hash, key []byte) ([]byte, error) {
	if hash == nil {
		for _, value := range s.commited {
			v, err := value.data.Get(key)
			if err == nil {
				return v, nil
			}
			if err != memdb.ErrNotFound {
				return nil, err
			}
		}
	} else {
		for _, value := range s.commited {
			if *hash == *value.BlockHash {
				v, err := value.data.Get(key)
				return v, err
			}
		}
	}
	return nil, memdb.ErrNotFound
}

func (s *SnapshotDB) getFromBaseDB(key []byte) ([]byte, error) {
	return s.baseDB.Get(key, nil)
}

func (s *SnapshotDB) getUnRecognizedHash() common.Hash {
	return rlpHash(CURRENT)
}

func (s *SnapshotDB) closeJournalWriter(hash common.Hash) error {
	if j, ok := s.journalw[hash]; ok {
		if err := j.Close(); err != nil {
			return errors.New("[snapshotdb]close  journal writer fail:" + err.Error())
		}
		delete(s.journalw, hash)
	}
	return nil
}

func (s *SnapshotDB) rmOldRecognizedBlockData() error {
	for key, value := range s.recognized {
		if s.current.HighestNum.Cmp(value.Number) <= 0 {
			delete(s.recognized, key)
			if err := s.closeJournalWriter(key); err != nil {
				return err
			}
			if err := s.rmJournalFile(value.Number, key); err != nil {
				//todo 怎么一致性...
				return err
			}
		}
	}
	return nil
}
