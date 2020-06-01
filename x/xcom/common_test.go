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

package xcom

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
)

func setup(t *testing.T) *mock.Chain {
	t.Log("setup()......")
	chain := mock.NewChain()
	chain.AddBlock()
	err := chain.SnapDB.NewBlock(chain.CurrentHeader().Number, chain.CurrentHeader().ParentHash, chain.CurrentHeader().Hash())
	if err != nil {
		fmt.Println("newBlock, %", err)
	}
	StorageAvgPackTime(chain.CurrentHeader().Hash(), chain.SnapDB, uint64(2000))
	commit_sndb(chain)

	prepair_sndb(chain, chain.CurrentHeader().Hash())
	return chain
}

func clear(chain *mock.Chain, t *testing.T) {
	t.Log("tear down()......")
	if err := chain.SnapDB.Clear(); err != nil {
		t.Error("clear chain.SnapDB error", err)
	}
}

func commit_sndb(chain *mock.Chain) {
	/*
		//Flush() signs a Hash to the current block which has no hash yet. Flush() do not write the data to database.
		//in this file, all blocks in each test case has a hash already, so, do not call Flush()
				if err := chain.SnapDB.Flush(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number); err != nil {
					fmt.Println("commit_sndb error:", err)
				}
	*/
	if err := chain.SnapDB.Commit(chain.CurrentHeader().Hash()); err != nil {
		fmt.Println("commit_sndb error:", err)
	}
}

func prepair_sndb(chain *mock.Chain, txHash common.Hash) {
	if txHash == common.ZeroHash {
		chain.AddBlock()
	} else {
		chain.AddBlockWithTxHash(txHash)
	}
	if err := chain.SnapDB.NewBlock(chain.CurrentHeader().Number, chain.CurrentHeader().ParentHash, chain.CurrentHeader().Hash()); err != nil {
		fmt.Println("prepair_sndb error:", err)
	}
}

func TestCommon_StorageAvgPackTime(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	avgPackTime, err := LoadCurrentAvgPackTime()
	if err != nil {
		t.Error("load current block average pack time error", err)
	}

	assert.Equal(t, uint64(2000), avgPackTime)
	StorageAvgPackTime(chain.CurrentHeader().Hash(), snapshotdb.Instance(), uint64(3000))
	//commit_sndb(chain)

	avgPackTime, err = LoadAvgPackTime(chain.CurrentHeader().Hash(), snapshotdb.Instance())
	assert.Equal(t, uint64(3000), avgPackTime)

	avgPackTime, err = LoadCurrentAvgPackTime()
	assert.Equal(t, uint64(2000), avgPackTime)
}

func TestStateDB(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	resByte := NewResult(common.InternalError, "test")

	res := new(Result)
	if err := json.Unmarshal(resByte, res); err != nil {
		t.Error(err)
	}
	if res.Code != common.InternalError.Code {
		t.Error("code must same")
	}

	AddLog(chain.StateDB, 1, common.ZeroAddr, "aa", "bb")
	AddLogWithRes(chain.StateDB, 1, common.ZeroAddr, "aa", "bb", "cc")
}
