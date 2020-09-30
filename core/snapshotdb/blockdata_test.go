package snapshotdb

import (
	"bytes"
	"testing"

	"github.com/syndtr/goleveldb/leveldb/memdb"
)

func TestBlockData_RevertToSnapshot(t *testing.T) {
	block := new(blockData)
	block.data = memdb.New(DefaultComparer, 10)
	keyA, keyB := []byte("a"), []byte("b")

	a := block.Snapshot()
	block.Write(keyA, keyA)
	b := block.Snapshot()
	block.Write(keyB, keyB)

	for _, valOld := range [][]byte{keyA, keyB} {
		if valNew, err := block.data.Get(valOld); err != nil {
			t.Error(err)
			return
		} else {
			if bytes.Compare(valOld, valNew) != 0 {
				t.Error("key not same")
			}
		}
	}

	block.RevertToSnapshot(b)

	if _, err := block.data.Get(keyB); err == nil {
		t.Error(err)
		return
	}

	if valNew, err := block.data.Get(keyA); err != nil {
		t.Error(err)
		return
	} else {
		if bytes.Compare(keyA, valNew) != 0 {
			t.Error("key not same")
		}
	}

	block.RevertToSnapshot(a)

	if _, err := block.data.Get(keyB); err == nil {
		t.Error(err)
		return
	}
	if _, err := block.data.Get(keyA); err == nil {
		t.Error(err)
		return
	}

	c := block.Snapshot()
	block.Write(keyA, keyA)
	block.Write(keyB, keyB)
	block.RevertToSnapshot(c)
	if v, err := block.data.Get(keyB); err == nil {
		t.Errorf("should not nil,%v", string(v))
		return
	}
	if v, err := block.data.Get(keyA); err == nil {
		t.Errorf("should not nil,%v", string(v))
		return
	}

	d := block.Snapshot()
	block.Write(keyA, nil)

	if valNew, err := block.data.Get(keyA); err != nil {
		t.Error(err)
		return
	} else {
		if bytes.Compare(nil, valNew) != 0 {
			t.Error("key not same")
		}
	}

	block.RevertToSnapshot(d)
	if _, err := block.data.Get(keyA); err == nil {
		t.Error(err)
		return
	}

	e := block.Snapshot()
	block.Write(keyA, nil)
	m := block.Snapshot()

	block.Write(keyA, keyA)

	block.RevertToSnapshot(m)
	if valNew, err := block.data.Get(keyA); err != nil {
		t.Error(err)
		return
	} else {
		if bytes.Compare(nil, valNew) != 0 {
			t.Error("key not same")
		}
	}
	block.RevertToSnapshot(e)
	if _, err := block.data.Get(keyA); err == nil {
		t.Error(err)
		return
	}

}
