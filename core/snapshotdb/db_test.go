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
	"bytes"
	"fmt"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestRecover(t *testing.T) {
	ch := newTestchain(dbpath)
	defer ch.clear()
	var (
		baseDBArr              []kv
		base, high             int64
		chainCommitHighestHash common.Hash
	)
	{
		baseDBArr = generatekv(100)
		if err := ch.insert(true, baseDBArr, newBlockBaseDB); err != nil {
			t.Error(err)
			return
		}
	}
	{
		if err := ch.insert(true, generatekv(100), newBlockCommited); err != nil {
			t.Error(err)
			return
		}
	}
	{
		if err := ch.insert(true, generatekv(100), newBlockCommited); err != nil {
			t.Error(err)
			return
		}
		chainCommitHighestHash = ch.CurrentHeader().Hash()
	}
	{
		if err := ch.insert(true, generatekv(100), newBlockRecognizedByFlush); err != nil {
			t.Error(err)
			return
		}

		_, ok := ch.db.unCommit.blocks[ch.CurrentHeader().Hash()]
		if !ok {
			t.Error("not found recognizedHash4")
		}
	}
	{
		if err := ch.insert(false, generatekv(100), newBlockUnRecognized); err != nil {
			t.Error(err)
			return
		}
	}
	base = ch.db.current.base.Num.Int64()
	high = ch.db.current.highest.Num.Int64()
	commitLength := len(ch.db.committed)

	oldarr := []blockData{
		*ch.db.committed[0],
		*ch.db.committed[1],
	}

	if err := ch.db.Close(); err != nil {
		t.Error(err)
	}

	ch.reOpenSnapshotDB()

	if ch.db.path != ch.path {
		t.Error("path is wrong", ch.db.path, ch.path)
		return
	}
	if ch.db.current.base.Num.Int64() != base {
		t.Error("BaseNum is wrong", ch.db.current.base.Num.Int64(), base)
		return
	}
	if ch.db.current.highest.Num.Int64() != high {
		t.Error("HighestNum is wrong", ch.db.current.highest.Num.Int64(), high)
		return
	}
	if ch.db.current.highest.Hash != chainCommitHighestHash {
		t.Error("Highest Hash is wrong", ch.db.current.highest.Hash.String(), ch.h[1].Hash().String())
		return
	}
	if len(ch.db.unCommit.blocks) != 0 {
		t.Error("recover uncommit should empty", len(ch.db.unCommit.blocks))
		return
	}
	for _, value := range baseDBArr {
		v, err := ch.db.baseDB.Get(value.key, nil)
		if err != nil {
			t.Error("should be nil", err)
			return
		}
		if bytes.Compare(v, value.value) != 0 {
			t.Error("should be equal", v, []byte(value.value))
			return
		}
	}

	if len(ch.db.committed) != commitLength {
		t.Error("should recover commit ", len(ch.db.committed))
		return
	}
	newarr := []blockData{
		*ch.db.committed[0],
		*ch.db.committed[1],
	}

	for i := 0; i < 2; i++ {
		if oldarr[i].BlockHash != newarr[i].BlockHash {
			t.Error("block hash must compare", i, oldarr[i].BlockHash.String(), newarr[i].BlockHash.String())
			return
		}

		if oldarr[i].ParentHash != newarr[i].ParentHash {
			t.Error("ParentHash  must compare", i)
			return
		}
		if oldarr[i].kvHash != newarr[i].kvHash {
			t.Error("kvHash  must compare", i)
			return
		}
		if !newarr[i].readOnly {
			t.Error("readOnly  must true", i, newarr[i].readOnly)
			return
		}
		itr := oldarr[i].data.NewIterator(nil)
		for itr.Next() {
			v, err := newarr[i].data.Get(itr.Key())
			if err != nil {
				t.Error(err)
				return
			}
			if bytes.Compare(v, itr.Value()) != 0 {
				t.Error("kv  must compare", i)
				return
			}
		}
	}

}

/*
func TestRMOldRecognizedBlockData(t *testing.T) {
	ch := new(testchain)
	blockchain = ch
	initDB()
	defer dbInstance.Clear()
	ch.addBlock()
	if err := newBlockBaseDB(ch.CurrentHeader(), generatekv(100)); err != nil {
		t.Error(err)
	}
	ch.addBlock()
	if err := newBlockCommited(ch.CurrentHeader(), generatekv(100)); err != nil {
		t.Error(err)
	}

	ch.addBlock()
	if err := newBlockRecognizedDirect(ch.currentForkHeader(), generatekv(100)); err != nil {
		t.Error(err)
	}
	if err := newBlockRecognizedDirect(ch.currentForkHeader(), generatekv(100)); err != nil {
		t.Error(err)
	}
	if err := newBlockCommited(ch.CurrentHeader(), generatekv(100)); err != nil {
		t.Error(err)
	}

	if err := dbInstance.rmOldRecognizedBlockData(); err != nil {
		t.Error(err)
	}
	if len(dbInstance.unCommit.blocks) != 0 {
		t.Error("not rm old data")
	}
}
*/

func (k kvs) compareWithkvs(s kvs) error {
	if len(k) != len(s) {
		return fmt.Errorf("kv length not compare,want %d have %d", len(k), len(s))
	}
	for i := 0; i < len(k); i++ {
		if bytes.Compare(k[i].key, s[i].key) != 0 {
			return fmt.Errorf("key not compare,want %v have %v", k[i].key, s[i].key)
		}
		if bytes.Compare(k[i].value, s[i].value) != 0 {
			return fmt.Errorf("value not compare,want %v have %v", k[i].value, s[i].value)
		}
	}
	return nil
}

func newBlockBaseDB(db *snapshotDB, kvs kvs, head *types.Header) error {
	hash := head.Hash()
	blockNumber := head.Number
	parentHash := head.ParentHash
	if err := db.NewBlock(blockNumber, parentHash, hash); err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := db.Put(hash, kv.key, kv.value); err != nil {
			return err
		}
	}
	if err := db.Commit(hash); err != nil {
		return err
	}
	db.walSync.Wait()
	return db.Compaction()
}

func newBlockCommited(db *snapshotDB, kvs kvs, head *types.Header) error {
	hash := head.Hash()
	blockNumber := head.Number
	parentHash := head.ParentHash
	if err := db.NewBlock(blockNumber, parentHash, hash); err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := db.Put(hash, kv.key, kv.value); err != nil {
			return err
		}
	}
	if err := db.Commit(hash); err != nil {
		return err
	}
	return nil
}

func newBlockRecognizedDirect(db *snapshotDB, kvs kvs, head *types.Header) error {
	hash := head.Hash()
	blockNumber := head.Number
	parentHash := head.ParentHash
	if err := db.NewBlock(blockNumber, parentHash, hash); err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := db.Put(hash, kv.key, kv.value); err != nil {
			return err
		}
	}
	return nil
}

func newBlockRecognizedByFlush(db *snapshotDB, kvs kvs, head *types.Header) error {
	if err := newBlockUnRecognized(db, kvs, head); err != nil {
		return err
	}
	return db.Flush(head.Hash(), head.Number)
}

func newBlockUnRecognized(db *snapshotDB, kvs kvs, head *types.Header) error {
	if err := db.NewBlock(head.Number, head.ParentHash, common.ZeroHash); err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := db.Put(common.ZeroHash, kv.key, kv.value); err != nil {
			return err
		}
	}
	return nil
}

//
//func TestCheckHashChain(t *testing.T) {
//	ch := new(testchain)
//	blockchain = ch
//	initDB()
//	defer dbInstance.Clear()
//
//	ch.addBlock()
//	if err := newBlockBaseDB(ch.CurrentHeader(), generatekv(1)); err != nil {
//		t.Error(err)
//	}
//
//	//2-11
//	for i := 2; i < 11; i++ {
//		ch.addBlock()
//		if err := newBlockCommited(ch.CurrentHeader(), generatekv(1)); err != nil {
//			t.Error(err)
//		}
//	}
//	//12-20
//	for i := 11; i < 21; i++ {
//		ch.addBlock()
//		if err := newBlockRecognizedDirect(ch.CurrentHeader(), generatekv(1)); err != nil {
//			t.Error(err)
//		}
//	}
//
//	if err := newBlockUnRecognized(ch.CurrentHeader(), generatekv(1)); err != nil {
//		t.Error(err)
//	}
//	t.Run("find from recognized", func(t *testing.T) {
//		for i := 11; i < 21; i++ {
//			location, ok := dbInstance.checkHashChain(ch.h[i-1].Hash())
//			if !ok {
//				t.Error("should be ok")
//			}
//			if location != hashLocationUnCommitted {
//				t.Error("should be locate Recognized", location)
//			}
//		}
//
//	})
//
//	t.Run("find from commit", func(t *testing.T) {
//		for i := 2; i < 11; i++ {
//			location, ok := dbInstance.checkHashChain(ch.h[i-1].Hash())
//			if !ok {
//				t.Error("should be ok")
//			}
//			if location != hashLocationCommitted {
//				t.Error("should be locate Recognized", location)
//			}
//		}
//	})
//
//	t.Run("find from unrecognized", func(t *testing.T) {
//		location, ok := dbInstance.checkHashChain(common.ZeroHash)
//		if !ok {
//			t.Error("should be ok")
//		}
//		if location != hashLocationUnCommitted {
//			t.Error("should be locate Recognized", location)
//		}
//	})
//
//	t.Run("not found", func(t *testing.T) {
//		location, ok := dbInstance.checkHashChain(generateHash(fmt.Sprint(1)))
//		if !ok {
//			t.Error("should be ok")
//		}
//		if location != hashLocationNotFound {
//			t.Error("should be locate Recognized", location)
//		}
//	})
//}
