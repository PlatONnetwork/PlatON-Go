package snapshotdb

import (
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
)

func TestCurrentUpdate(t *testing.T) {
	dir := os.TempDir()
	c := newCurrent(dir)
	c.HighestNum = big.NewInt(10)
	c.BaseNum = big.NewInt(5)
	if err := c.update(); err != nil {
		t.Error(err)
	}
	c.f.Seek(0, 0)
	cb1, err := ioutil.ReadAll(c.f)
	if err != nil {
		t.Error(err)
	}
	var cu1 current
	if err := rlp.DecodeBytes(cb1, &cu1); err != nil {
		t.Error(err)
	}
	if cu1.HighestNum.Int64() != 10 {
		t.Fatal("HighestNum not compare 10", cu1.HighestNum.Int64())
	}
	if cu1.BaseNum.Int64() != 5 {
		t.Fatal("BaseNum not compare")
	}

	c.HighestNum = big.NewInt(20000000)
	c.BaseNum = big.NewInt(10000)
	if err := c.update(); err != nil {
		t.Error(err)
	}
	c.f.Seek(0, 0)
	cb2, err := ioutil.ReadAll(c.f)
	if err != nil {
		t.Error(err)
	}
	var cu2 current
	if err := rlp.DecodeBytes(cb2, &cu2); err != nil {
		t.Error(err)
	}
	if cu2.HighestNum.Int64() != 20000000 {
		t.Fatal("HighestNum not compare")
	}
	if cu2.BaseNum.Int64() != 10000 {
		t.Fatal("BaseNum not compare")
	}
	c.f.Close()

}
