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
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

const (
	CurrentHighestBlock = "snapshotdbCurrentHighestBlock"
	CurrentBaseNum      = "snapshotdbCurrentBaseNum"
	CurrentAll          = "snapshotdbCurrentAll"
	CurrentSet          = "snapshotdbCurrentSet"
)

func (s *snapshotDB) loadCurrent() error {
	ct := new(current)
	if err := ct.loadFromBaseDB(s.baseDB); err != nil {
		return err
	}

	if err := ct.Valid(); err != nil {
		return err
	}
	if blockchain != nil {
		currentHead := blockchain.CurrentHeader()
		if err := ct.resetHighestByChainCurrentHeader(currentHead, s.baseDB); err != nil {
			return err
		}
	}
	s.current = ct
	return nil
}

func newCurrent(highestNum, baseNum *big.Int, highestHash common.Hash) *current {
	c := new(current)
	c.highest = new(CurrentHighest)
	if highestNum != nil {
		c.highest.Num = new(big.Int).Set(highestNum)
	}
	c.highest.Hash = highestHash
	c.base = new(CurrentBase)
	if baseNum != nil {
		c.base.Num = new(big.Int).Set(baseNum)
	}
	return c
}

type current struct {
	base    *CurrentBase
	highest *CurrentHighest
}

func (c *current) GetHighest(copy bool) *CurrentHighest {
	if copy {
		h := new(CurrentHighest)
		h.Num = new(big.Int).Set(c.highest.Num)
		h.Hash = c.highest.Hash
		return h
	}
	return c.highest
}

func (c *current) GetHighestFromDB(baseDB *leveldb.DB) (*CurrentHighest, error) {
	hight, err := baseDB.Get([]byte(CurrentHighestBlock), nil)
	if err != nil {
		return nil, fmt.Errorf("get current highest block fail:%v", err)
	}
	ch := new(CurrentHighest)
	if err := rlp.DecodeBytes(hight, ch); err != nil {
		return nil, fmt.Errorf("decode current highest block fail:%v", err)
	}
	return ch, nil
}

func (c *current) GetBase(copy bool) *CurrentBase {
	if copy {
		h := new(CurrentBase)
		h.Num = new(big.Int).Set(c.base.Num)
		return h
	}
	return c.base
}

// notice:this valid  assume current have already loadFromBaseDB
func (c *current) Valid() error {
	if c.base.Num.Cmp(c.highest.Num) > 0 {
		return fmt.Errorf("base num %v can't be greater than highest Num %v", c.base.Num, c.highest.Num)
	}
	return nil
}

//the current highest  must not  greater than block chain current
func (c *current) resetHighestByChainCurrentHeader(currentHead *types.Header, baseDB *leveldb.DB) error {
	if c.base.Num.Cmp(currentHead.Number) > 0 {
		return fmt.Errorf("base num %v can't be greater than currentHead Number %v", c.base.Num, currentHead.Number)
	}
	if c.highest.Num.Cmp(currentHead.Number) > 0 {
		c.highest.Num = new(big.Int).Set(currentHead.Number)
		c.highest.Hash = currentHead.Hash()
		if err := c.saveCurrentToBaseDB(CurrentHighestBlock, baseDB, false); err != nil {
			return err
		}
	}
	return nil
}

func (c *current) loadFromBaseDB(baseDB *leveldb.DB) error {
	base, err := baseDB.Get([]byte(CurrentBaseNum), nil)
	if err != nil {
		return fmt.Errorf("get current base num fail:%v", err)
	}
	c.base = new(CurrentBase)
	if err := rlp.DecodeBytes(base, c.base); err != nil {
		return fmt.Errorf("decode current base num fail:%v", err)
	}
	hight, err := baseDB.Get([]byte(CurrentHighestBlock), nil)
	if err != nil {
		return fmt.Errorf("get current highest block fail:%v", err)
	}
	c.highest = new(CurrentHighest)
	if err := rlp.DecodeBytes(hight, c.highest); err != nil {
		return fmt.Errorf("decode current highest block fail:%v", err)
	}
	return nil
}

func (c *current) increaseBase(commitNum uint64, baseDB *leveldb.DB) error {
	c.base.Num.Add(c.base.Num, new(big.Int).SetUint64(commitNum))
	if err := c.saveCurrentToBaseDB(CurrentBaseNum, baseDB, false); err != nil {
		return err
	}
	return nil
}

func (c *current) increaseHighest(hash common.Hash) {
	c.highest.Num.Add(c.highest.Num, common.Big1)
	c.highest.Hash = hash
	logger.Debug("increase current highest", "hash", hash, "num", c.highest.Num)
}

func (c *current) saveCurrentToBaseDB(currentType string, db *leveldb.DB, init bool) error {
	batch := new(leveldb.Batch)
	switch currentType {
	case CurrentHighestBlock:
		height := c.EncodeHighest()
		batch.Put([]byte(CurrentHighestBlock), height)
	case CurrentBaseNum:
		base := c.EncodeBase()
		batch.Put([]byte(CurrentBaseNum), base)
	default:
		height, base := c.EncodeHighest(), c.EncodeBase()
		batch.Put([]byte(CurrentHighestBlock), height)
		batch.Put([]byte(CurrentBaseNum), base)
	}
	if init {
		batch.Put([]byte(CurrentSet), []byte(CurrentSet))
	}
	if c.highest != nil && c.base != nil {
		logger.Debug("save current to baseDB", "height", c.highest.Num, "hash", c.highest.Hash, "base", c.base.Num, "type", currentType)
	} else {
		logger.Debug("save current to baseDB", "height", c.highest, "base", c.base, "type", currentType)
	}
	if err := db.Write(batch, nil); err != nil {
		return fmt.Errorf("write %v  to base db fail:%v", currentType, err)
	}
	return nil
}

func (c *current) EncodeHighest() []byte {
	highest, err := rlp.EncodeToBytes(c.highest)
	if err != nil {
		panic(err)
	}
	return highest
}

func (c *current) EncodeBase() []byte {
	base, err := rlp.EncodeToBytes(c.base)
	if err != nil {
		panic(err)
	}
	return base
}

type CurrentHighest struct {
	Num  *big.Int `rlp:"nil"`
	Hash common.Hash
}

type CurrentBase struct {
	Num *big.Int `rlp:"nil"`
}
