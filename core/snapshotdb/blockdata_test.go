// Copyright 2021 The PlatON Network Authors
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
	"bytes"
	"testing"

	"github.com/syndtr/goleveldb/leveldb/memdb"
)

func TestBlockData_RevertToSnapshot(t *testing.T) {
	block := new(BlockData)
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
			if !bytes.Equal(valOld, valNew) {
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
		if !bytes.Equal(keyA, valNew) {
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
		if !bytes.Equal(nil, valNew) {
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
		if !bytes.Equal(nil, valNew) {
			t.Error("key not same")
		}
	}
	block.RevertToSnapshot(e)
	if _, err := block.data.Get(keyA); err == nil {
		t.Error(err)
		return
	}
}
