// +build test

package snapshotdb

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/robfig/cron"
)

const (
	//DBPath path of db
	DBPath = "snapshotdb_test"
	//DBBasePath path of basedb
	DBBasePath  = "base"
	currentPath = "current"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	//	logger.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	logger.Info("begin test")
	dbpath = path.Join(os.TempDir(), DBPath, fmt.Sprint(rand.Uint64()))
	testChain := new(testchain)
	header := generateHeader(big.NewInt(1000000000), common.ZeroHash)
	testChain.h = append(testChain.h, header)
	blockchain = testChain
}

type testchain struct {
	h []*types.Header
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

// New  create a new snapshotDB
func New() (DB, error) {
	p := path.Join(os.TempDir(), DBPath, fmt.Sprint(rand.Uint64()))
	logger.Info("begin newDB", "path", p)
	s, err := openFile(p, false)
	if err != nil {
		logger.Error(fmt.Sprint("open db file fail:", err))
		return nil, err
	}
	db, err := newDB(s)
	if err != nil {
		logger.Error(fmt.Sprint("new db fail:", err))
		return nil, err
	}
	db.corn = cron.New()
	if err := db.corn.AddFunc("@every 1s", dbInstance.schedule); err != nil {
		logger.Error(fmt.Sprint("new db fail", err))
		return nil, err
	}
	db.corn.Start()
	return db, nil
}
