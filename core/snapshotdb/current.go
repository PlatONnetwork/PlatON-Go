package snapshotdb

import (
	"fmt"
	"math/big"

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

func (s *snapshotDB) saveCurrentToBaseDB(t string, c *current) error {
	batch := new(leveldb.Batch)
	height, base := c.ToByte()
	switch t {
	case CurrentHighestBlock:
		batch.Put([]byte(CurrentHighestBlock), height)
	case CurrentBaseNum:
		batch.Put([]byte(CurrentBaseNum), base)
	default:
		batch.Put([]byte(CurrentHighestBlock), height)
		batch.Put([]byte(CurrentBaseNum), base)
	}
	if err := s.baseDB.Write(batch, nil); err != nil {
		return fmt.Errorf("write %v  to base db fail:%v", t, err)
	}
	return nil
}

func (s *snapshotDB) newCurrent() error {
	c := new(current)
	c.HighestNum = big.NewInt(0)
	c.BaseNum = big.NewInt(0)
	c.HighestHash = common.ZeroHash
	s.current = c
	if err := s.saveCurrentToBaseDB(CurrentAll, c); err != nil {
		return err
	}
	if err := s.baseDB.Put([]byte(CurrentSet), []byte(CurrentSet), nil); err != nil {
		return err
	}
	return nil
}

func (s *snapshotDB) loadCurrent() error {
	hight, err := s.baseDB.Get([]byte(CurrentHighestBlock), nil)
	if err != nil {
		return fmt.Errorf("get current highest block fail:%v", err)
	}
	var ch CurrentHighest
	if err := rlp.DecodeBytes(hight, &ch); err != nil {
		return fmt.Errorf("decode current highest block fail:%v", err)
	}
	base, err := s.baseDB.Get([]byte(CurrentBaseNum), nil)
	if err != nil {
		return fmt.Errorf("get current base num fail:%v", err)
	}
	var cb CurrentBase
	if err := rlp.DecodeBytes(base, &cb); err != nil {
		return fmt.Errorf("decode current base num fail:%v", err)
	}
	c := new(current)
	c.BaseNum = cb.Num
	c.HighestHash = ch.Hash
	c.HighestNum = ch.Num
	if blockchain != nil {
		currentHead := blockchain.CurrentHeader()
		if c.HighestNum.Cmp(currentHead.Number) > 0 {
			c.HighestNum = currentHead.Number
			c.HighestHash = currentHead.Hash()
			if err := s.saveCurrentToBaseDB(CurrentHighestBlock, c); err != nil {
				return err
			}
		}
	}
	s.current = c
	return nil
}

type current struct {
	//	f           *os.File    `rlp:"-"`
	//	path        string      `rlp:"-"`
	HighestNum  *big.Int    `rlp:"nil"`
	HighestHash common.Hash `rlp:"nil"`
	BaseNum     *big.Int    `rlp:"nil"`
	//	sync.RWMutex `rlp:"-"`
}

type CurrentHighest struct {
	Num  *big.Int    `rlp:"nil"`
	Hash common.Hash `rlp:"nil"`
}

type CurrentBase struct {
	Num *big.Int `rlp:"nil"`
}

func (c *current) ToByte() ([]byte, []byte) {
	highest, err := rlp.EncodeToBytes(CurrentHighest{c.HighestNum, c.HighestHash})
	if err != nil {
		panic(err)
	}
	base, err := rlp.EncodeToBytes(CurrentBase{c.BaseNum})
	if err != nil {
		panic(err)
	}
	return highest, base
}

//
//func (c *current) update() error {
//	b := new(bytes.Buffer)
//	if err := rlp.Encode(b, c); err != nil {
//		return err
//	}
//	if err := c.f.Truncate(0); err != nil {
//		return err
//	}
//	c.f.Seek(io.SeekStart, io.SeekEnd)
//	_, err := c.f.Write(b.Bytes())
//	if err != nil {
//		return err
//	}
//	if err := c.f.Sync(); err != nil {
//		return err
//	}
//	return nil
//}
//
//func getCurrentPath(dir string) string {
//	return path.Join(dir, currentPath)
//}
//
//func loadCurrent(dir string) (*current, error) {
//	cpath := getCurrentPath(dir)
//	f, err := os.OpenFile(cpath, os.O_RDWR, 0666)
//	if err != nil {
//		return nil, err
//	}
//	currentBytes, err := ioutil.ReadAll(f)
//	if err != nil {
//		return nil, err
//	}
//	c := new(current)
//	if err := rlp.DecodeBytes(currentBytes, c); err != nil {
//		return nil, err
//	}
//	c.path = cpath
//	c.f = f
//	return c, nil
//}
