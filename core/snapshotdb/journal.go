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
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

const WalKeyPrefix = "journal-"

func EncodeWalKey(blockNum *big.Int) []byte {
	return append([]byte(WalKeyPrefix), blockNum.Bytes()...)
}

func DecodeWalKey(key []byte) *big.Int {
	return new(big.Int).SetBytes(key[len([]byte(WalKeyPrefix)):])
}

type journalData struct {
	Key, Value []byte
}

type blockWal struct {
	ParentHash  common.Hash
	BlockHash   common.Hash
	BlockNumber *big.Int `rlp:"nil"`
	KvHash      common.Hash
	Data        []journalData
}

func (s *snapshotDB) loopWriteWal() {
	for {
		select {
		case block := <-s.walCh:
			if err := s.writeWal(block); err != nil {
				logger.Error("asynchronous write Journal fail", "err", err, "block", block.Number, "hash", block.BlockHash.String())
				s.dbError = err
				s.walSync.Done()
				continue
			}
			nc := newCurrent(block.Number, nil, block.BlockHash)
			if err := nc.saveCurrentToBaseDB(CurrentHighestBlock, s.baseDB, false); err != nil {
				logger.Error("asynchronous update current highest fail", "err", err, "block", block.Number, "hash", block.BlockHash.String())
				s.dbError = err
				s.walSync.Done()
				continue
			}
			s.walSync.Done()
		case <-s.walExitCh:
			logger.Info("loopWriteWal exist")
			close(s.walCh)
			return
		}
	}
}

func (s *snapshotDB) writeBlockToWalAsynchronous(block *blockData) {
	s.walSync.Add(1)
	s.walCh <- block
}

func (s *snapshotDB) writeWal(block *blockData) error {
	return s.baseDB.Put(block.BlockKey(), block.BlockVal(), nil)
}
