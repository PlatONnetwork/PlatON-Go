package snapshotdb

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"sync"
)

func newCurrent(dir string) *current {
	c := new(current)
	c.HighestNum = big.NewInt(0)
	c.BaseNum = big.NewInt(0)
	c.path = getCurrentPath(dir)
	c.HighestHash = common.ZeroHash
	f, err := os.OpenFile(c.path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	c.f = f
	c.update()
	return c
}

type current struct {
	f            *os.File    `rlp:"-"`
	path         string      `rlp:"-"`
	HighestNum   *big.Int    `rlp:"nil"`
	HighestHash  common.Hash `rlp:"nil"`
	BaseNum      *big.Int    `rlp:"nil"`
	sync.RWMutex `rlp:"-"`
}

func (c *current) update() error {
	c.Lock()
	defer c.Unlock()
	if err := c.f.Truncate(0); err != nil {
		return err
	}
	b := new(bytes.Buffer)
	if err := rlp.Encode(b, c); err != nil {
		return err
	}
	c.f.Seek(io.SeekStart, io.SeekEnd)
	_, err := c.f.Write(b.Bytes())
	if err != nil {
		return err
	}
	if err := c.f.Sync(); err != nil {
		return err
	}
	return nil
}

func getCurrentPath(dir string) string {
	return path.Join(dir, currentPath)
}

func loadCurrent(dir string) (*current, error) {
	cpath := getCurrentPath(dir)
	f, err := os.OpenFile(cpath, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	currentBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	c := new(current)
	if err := rlp.DecodeBytes(currentBytes, c); err != nil {
		return nil, err
	}
	c.path = cpath
	c.f = f
	return c, nil
}
