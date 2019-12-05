package xcom

import (
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
