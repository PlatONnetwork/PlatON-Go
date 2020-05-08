// Copyright 2018-2019 The PlatON Network Authors
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
	"sync"

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
}

type unCommitBlocks struct {
	blocks map[common.Hash]*blockData
	sync.RWMutex
}

func (u *unCommitBlocks) Get(key common.Hash) *blockData {
	u.RLock()
	block, ok := u.blocks[key]
	u.RUnlock()
	if !ok {
		return nil
	}
	return block
}

func (u *unCommitBlocks) Set(key common.Hash, block *blockData) {
	u.Lock()
	u.blocks[key] = block
	u.Unlock()
}
