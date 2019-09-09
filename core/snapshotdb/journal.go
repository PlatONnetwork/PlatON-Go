package snapshotdb

import (
	"io"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb/journal"
)

// fileType represent a file type.
type fileType int

// File types.
const (
	TypeCurrent fileType = 1 << iota
	TypeJournal
	TypeAll = TypeCurrent | TypeJournal
)

type journalData struct {
	Key, Value []byte
	Hash       common.Hash
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

func newJournalWriter(w io.WriteCloser) *journalWriter {
	j := new(journalWriter)
	j.writer = w
	j.journal = journal.NewWriter(w)
	return j
}

type journalWriter struct {
	writer  io.WriteCloser
	journal *journal.Writer
}

func (j *journalWriter) Close() error {
	if err := j.journal.Close(); err != nil {
		return err
	}
	if err := j.writer.Close(); err != nil {
		return err
	}
	return nil
}

func (s *snapshotDB) writeJournalHeader(blockNumber *big.Int, hash, parentHash common.Hash, comeFrom string) (*journalWriter, error) {
	fd := fileDesc{Type: TypeJournal, Num: blockNumber.Uint64(), BlockHash: hash}
	file, err := s.storage.Create(fd)
	if err != nil {
		return nil, err
	}
	writers := newJournalWriter(file)
	jHeader := journalHeader{
		ParentHash:  parentHash,
		BlockNumber: blockNumber,
		From:        comeFrom,
	}
	h, err := encode(jHeader)
	if err != nil {
		return nil, err
	}

	writer, err := writers.journal.Next()
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(h); err != nil {
		return nil, err
	}
	if err := writers.journal.Flush(); err != nil {
		return nil, err
	}
	return writers, nil
}

func (s *snapshotDB) rmJournalFile(blockNumber *big.Int, hash common.Hash) error {
	fd := fileDesc{Type: TypeJournal, Num: blockNumber.Uint64(), BlockHash: hash}
	return s.storage.Remove(fd)
}
