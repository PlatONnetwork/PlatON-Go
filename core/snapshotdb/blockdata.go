package snapshotdb

import (
	"errors"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb/memdb"
)

type blockData struct {
	BlockHash  common.Hash
	ParentHash common.Hash
	Number     *big.Int
	data       *memdb.DB
	readOnly   bool
	kvHash     common.Hash
	jWriter    *journalWriter
}

func (s *blockData) writeJournalBody(value []byte) error {
	toWrite, err := s.jWriter.journal.Next()
	if err != nil {
		return errors.New("next err:" + err.Error())
	}
	if _, err := toWrite.Write(value); err != nil {
		return errors.New("write err:" + err.Error())
	}
	if err := s.jWriter.journal.Flush(); err != nil {
		return errors.New("flush err:" + err.Error())
	}
	return nil
}
