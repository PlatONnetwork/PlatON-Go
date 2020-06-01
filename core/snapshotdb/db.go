// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package snapshotdb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"path"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/memdb"
)

func getBaseDBPath(dbpath string) string {
	return path.Join(dbpath, DBBasePath)
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
	block.kvHash = header.KvHash
	if fd.BlockHash != s.getUnRecognizedHash() {
		block.BlockHash = fd.BlockHash
	}
	block.Number = new(big.Int).SetUint64(fd.Num)
	block.data = memdb.New(DefaultComparer, 0)
	block.readOnly = true

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
		//kvhash = body.Hash
	}
	//block.kvHash = kvhash
	return block, nil
}

func (s *snapshotDB) recover() error {
	//storage
	fds, err := s.storage.List(TypeJournal)
	if err != nil {
		return err
	}
	sortFds(fds)
	baseNum := s.current.GetBase(false).Num.Uint64()
	highestNum := s.current.GetHighest(false).Num.Uint64()
	//read Journal
	if len(fds) == 0 {
		if baseNum != highestNum {
			return errors.New("current baseNum and highestNum is not same,but not wal find")
		}
		return nil
	}
	if fds[0].Num > baseNum+1 {
		return fmt.Errorf("wal is not enough,want recover  from %v,have %v", baseNum+1, fds[0].Num)
	}
	fileToRecover, fileToRemove := make([]fileDesc, 0), make([]fileDesc, 0)
	if blockchain != nil {
		for _, fd := range fds {
			if baseNum < fd.Num && fd.Num <= highestNum {
				if header := blockchain.GetHeaderByHash(fd.BlockHash); header == nil {
					fileToRemove = append(fileToRemove, fd)
				} else {
					fileToRecover = append(fileToRecover, fd)
				}
			} else {
				fileToRemove = append(fileToRemove, fd)
			}
		}
	} else {
		for _, fd := range fds {
			if baseNum < fd.Num && fd.Num <= highestNum {
				fileToRecover = append(fileToRecover, fd)
			} else {
				fileToRemove = append(fileToRemove, fd)
			}
		}
	}
	for _, fd := range fileToRemove {
		logger.Info("recovering, removeing journal no need", "num", fd.Num)
		if err := s.storage.Remove(fd); err != nil {
			return err
		}
	}

	var (
		journalBroken      bool = false
		lastNotBrokenBlock *blockData
	)

	for _, fd := range fileToRecover {
		if journalBroken {
			logger.Info("recovering, some block is broken,remove left", "num", fd.Num)
			if err := s.storage.Remove(fd); err != nil {
				return err
			}
			continue
		}
		block, err := s.getBlockFromJournal(fd)
		if err != nil {
			journalBroken = true
			logger.Info("recovering, block is broken,remove it", "num", fd.Num, "err", err)
			if err := s.storage.Remove(fd); err != nil {
				return err
			}
			continue
		}
		s.committed = append(s.committed, block)
		lastNotBrokenBlock = block
		logger.Debug("recover block ", "num", block.Number, "hash", block.BlockHash.String())
	}

	if journalBroken {
		base := big.NewInt(int64(baseNum))
		if err := s.SetCurrent(lastNotBrokenBlock.BlockHash, *base, *lastNotBrokenBlock.Number); err != nil {
			return err
		}
	} else {
		if len(fileToRecover) > 0 {
			base := big.NewInt(int64(baseNum))
			block := fileToRecover[len(fileToRecover)-1]
			highest := big.NewInt(int64(block.Num))
			if err := s.SetCurrent(block.BlockHash, *base, *highest); err != nil {
				return err
			}
		} else {
			//no recover block,so set current highest and base the same
			base := big.NewInt(int64(baseNum))
			if err := s.SetCurrent(common.ZeroHash, *base, *base); err != nil {
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

	block.kvHash = s.generateKVHash(key, value, block.kvHash)
	if err := block.data.Put(key, value); err != nil {
		return err
	}
	return nil
}
