package snapshotdb

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	jo "github.com/syndtr/goleveldb/leveldb/journal"
	"math/big"
)

// fileType represent a file type.
type fileType int

// File types.
const (
	TypeCurrent fileType = 1 << iota
	TypeJournal
	TypeAll = TypeCurrent | TypeJournal
)

const (
	funcTypePut = iota
	funcTypeDel
)

type journalData struct {
	Key, Value []byte
	Hash       common.Hash
	FuncType   uint64
}

const (
	journalHeaderFromUnRecognized = "unRecognized"
	journalHeaderFromRecognized   = "recognized"
)

type journalHeader struct {
	ParentHash  common.Hash
	BlockNumber *big.Int `rlp:"nil"`
	From        string
}

func (s *snapshotDB) writeJournalHeader(blockNumber *big.Int, hash, parentHash common.Hash, comeFrom string) error {
	fd := fileDesc{Type: TypeJournal, Num: blockNumber.Int64(), BlockHash: hash}
	file, err := s.storage.Create(fd)
	if err != nil {
		return err
	}
	writers := jo.NewWriter(file)
	jHeader := journalHeader{
		ParentHash:  parentHash,
		BlockNumber: blockNumber,
		From:        comeFrom,
	}
	h, err := encode(jHeader)
	if err != nil {
		return err
	}
	writer, err := writers.Next()
	if err != nil {
		return err
	}
	if _, err := writer.Write(h); err != nil {
		return err
	}
	writers.Flush()
	s.journalw[hash] = writers
	return nil
}

func (s *snapshotDB) writeJournalBody(hash common.Hash, value []byte) error {
	var jw *jo.Writer
	var ok bool
	jw, ok = s.journalw[hash]
	if !ok {
		return errors.New("not found journal writer")
	}
	toWrite, err := jw.Next()
	if err != nil {
		return err
	}
	if _, err := toWrite.Write(value); err != nil {
		return err
	}
	jw.Flush()
	return nil
}

func (s *snapshotDB) rmJournalFile(blockNumber *big.Int, hash common.Hash) error {
	fd := fileDesc{Type: TypeJournal, Num: blockNumber.Int64(), BlockHash: hash}
	return s.storage.Remove(fd)
}
