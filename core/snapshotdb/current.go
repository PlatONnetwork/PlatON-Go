package snapshotdb

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"sync"
)

func newCurrent(dir string) *current {
	c := new(current)
	c.HighestNum = big.NewInt(0)
	c.BaseNum = big.NewInt(0)
	c.path = getCurrentPath(dir)
	f, err := os.OpenFile(c.path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	c.f = f
	c.update()
	return c
}

type current struct {
	f            *os.File `rlp:"-"`
	path         string   `rlp:"-"`
	HighestNum   *big.Int `rlp:"nil"`
	BaseNum      *big.Int `rlp:"nil"`
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

func findCurrent(dir string) bool {
	p := getCurrentPath(dir)
	matches, err := filepath.Glob(p)
	if err != nil {
		log.Error("find current fail:", err)
		return false
	}
	if len(matches) == 0 {
		return false
	}
	return true
}

func loadCurrent(dir string) (*current, error) {
	f, err := os.OpenFile(getCurrentPath(dir), os.O_RDWR, 0666)
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
	c.path = dir
	c.f = f
	return c, nil
}
