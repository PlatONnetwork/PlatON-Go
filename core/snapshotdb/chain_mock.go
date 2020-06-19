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
	"math/rand"
	"os"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

func newTestchain(path string) *testchain {
	os.RemoveAll(path)
	ch := new(testchain)
	ch.path = path
	db, err := open(path, 0, 0, false)
	if err != nil {
		panic(err)
	}
	ch.db = db
	SetDBBlockChain(ch)
	go ch.db.loopWriteWal()

	return ch
}

type testchain struct {
	h    []*types.Header
	db   *snapshotDB
	path string
}

func (c *testchain) reOpenSnapshotDB() {
	db, err := open(c.path, 0, 0, false)
	if err != nil {
		panic(err)
	}
	c.db = db
	go c.db.loopWriteWal()

}

func (c *testchain) insert(addBlock bool, kvs kvs, f func(db *snapshotDB, kvs kvs, header *types.Header) error) error {
	if addBlock {
		c.addBlock()
	}
	head := c.CurrentHeader()
	if err := f(c.db, kvs, head); err != nil {
		return err
	}
	return nil
}

func (c *testchain) clear() {
	c.db.Clear()
}

func (c *testchain) addBlock() {
	if len(c.h) == 0 {
		c.h = make([]*types.Header, 0)
		c.h = append(c.h, generateHeader(big.NewInt(1), common.ZeroHash))
		return
	}

	header := generateHeader(new(big.Int).Add(c.h[len(c.h)-1].Number, common.Big1), c.h[len(c.h)-1].Hash())
	c.h = append(c.h, header)
}

func (c *testchain) CurrentHeader() *types.Header {
	if len(c.h) != 0 {
		return c.h[len(c.h)-1]
	}
	return nil
}

func (c *testchain) currentForkHeader() *types.Header {
	if len(c.h) != 0 {
		newhead := new(types.Header)
		newhead.Number = c.h[len(c.h)-1].Number
		newhead.ParentHash = c.h[len(c.h)-1].ParentHash
		newhead.GasUsed = rand.Uint64()
		return newhead
	}
	return nil
}

func (c *testchain) GetHeaderByHash(hash common.Hash) *types.Header {
	for i := len(c.h) - 1; i >= 0; i-- {
		if c.h[i].Hash() == hash {
			return c.h[i]
		}
	}
	return nil
}

func (c *testchain) GetHeaderByNumber(number uint64) *types.Header {
	for i := len(c.h) - 1; i >= 0; i-- {
		if c.h[i].Number.Uint64() == number {
			return c.h[i]
		}
	}
	return nil
}
