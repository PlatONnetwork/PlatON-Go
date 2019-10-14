package snapshotdb

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestRecover(t *testing.T) {
	ch := new(testchain)
	blockchain = ch
	os.RemoveAll(dbpath)
	if err := initDB(); err != nil {
		t.Error(err)
		return
	}
	defer dbInstance.Clear()
	var (
		baseDBArr, commitArr, recognizedArr, unrecognizedArr []kv
		base, high                                           int64
		commit, recognized                                   *blockData
	)
	{
		baseDBArr = generatekv(100)
		ch.addBlock()
		if err := newBlockBaseDB(ch.CurrentHeader(), baseDBArr); err != nil {
			t.Error(err)
			return
		}
	}
	{
		commitArr = generatekv(100)
		ch.addBlock()
		if err := newBlockCommited(ch.CurrentHeader(), commitArr); err != nil {
			t.Error(err)
			return
		}
		commit = dbInstance.committed[0]
	}
	{
		recognizedArr = generatekv(100)
		ch.addBlock()
		if err := newBlockRecognizedByFlush(ch.CurrentHeader().Hash(), ch.GetHeaderByHash(ch.CurrentHeader().ParentHash), recognizedArr); err != nil {
			t.Error(err)
			return
		}

		rg, ok := dbInstance.unCommit.blocks[ch.CurrentHeader().Hash()]
		if !ok {
			t.Error("not found recognizedHash4")
		}
		recognized = rg
	}
	{
		unrecognizedArr = generatekv(100)
		if err := newBlockUnRecognized(ch.CurrentHeader(), unrecognizedArr); err != nil {
			t.Error(err)
			return
		}
	}
	base = dbInstance.current.BaseNum.Int64()
	high = ch.CurrentHeader().Number.Int64()
	if err := dbInstance.Close(); err != nil {
		t.Error(err)
	}
	dbInstance = nil
	//s, err := openFile(dbpath, false)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//fds, err := s.List(TypeCurrent)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//if len(fds) > 0 {
	//
	//}
	initDB()
	defer dbInstance.Clear()
	if dbInstance.path != dbpath {
		t.Error("path is wrong", dbInstance.path, dbpath)
		return
	}
	if dbInstance.current.BaseNum.Int64() != base {
		t.Error("BaseNum is wrong", dbInstance.current.BaseNum.Int64(), base)
		return
	}
	if dbInstance.current.HighestNum.Int64() != high {
		t.Error("HighestNum is wrong", dbInstance.current.HighestNum.Int64(), high)
		return
	}
	if dbInstance.current.HighestHash != ch.CurrentHeader().Hash() {
		t.Error("Highest Hash is wrong", dbInstance.current.HighestHash.String(), ch.CurrentHeader().Hash().String())
		return
	}
	if len(dbInstance.unCommit.blocks) != 0 {
		t.Error("recover uncommit should empty", len(dbInstance.unCommit.blocks))
		return
	}
	for _, value := range baseDBArr {
		v, err := dbInstance.baseDB.Get(value.key, nil)
		if err != nil {
			t.Error("should be nil", err)
			return
		}
		if bytes.Compare(v, value.value) != 0 {
			t.Error("should be equal", v, []byte(value.value))
			return
		}
	}
	oldarr := []*blockData{
		recognized,
		commit,
	}
	if len(dbInstance.committed) != 2 {
		t.Error("should recover commit ", len(dbInstance.committed))
		return
	}
	newarr := []*blockData{
		dbInstance.committed[1],
		dbInstance.committed[0],
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
func randomString2(s string) []byte {
	b := new(bytes.Buffer)
	if s != "" {
		b.Write([]byte(s))
	}
	for i := 0; i < 8; i++ {
		b.WriteByte(' ' + byte(rand.Uint64()))
	}
	return b.Bytes()
}

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

func generatekv(n int) kvs {
	rand.Seed(time.Now().UnixNano())
	kvs := make(kvs, n)
	for i := 0; i < n; i++ {
		kvs[i] = kv{
			key:   randomString2(""),
			value: randomString2(""),
		}
	}
	sort.Sort(kvs)
	return kvs
}

func generatekvWithPrefix(n int, p string) kvs {
	rand.Seed(time.Now().UnixNano())
	kvs := make(kvs, n)
	for i := 0; i < n; i++ {
		kvs[i] = kv{
			key:   randomString2(p),
			value: randomString2(p),
		}
	}
	sort.Sort(kvs)
	return kvs
}

func newBlockBaseDB(head *types.Header, kvs []kv) error {
	hash := head.Hash()
	blockNumber := head.Number
	parentHash := head.ParentHash
	if err := dbInstance.NewBlock(blockNumber, parentHash, hash); err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := dbInstance.Put(hash, kv.key, kv.value); err != nil {
			return err
		}
	}
	if err := dbInstance.Commit(hash); err != nil {
		return err
	}
	return dbInstance.Compaction()
}

func newBlockCommited(head *types.Header, kvs []kv) error {
	hash := head.Hash()
	blockNumber := head.Number
	parentHash := head.ParentHash
	if err := dbInstance.NewBlock(blockNumber, parentHash, hash); err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := dbInstance.Put(hash, kv.key, kv.value); err != nil {
			return err
		}
	}
	if err := dbInstance.Commit(hash); err != nil {
		return err
	}
	return nil
}

func newBlockRecognizedDirect(head *types.Header, kvs []kv) error {
	hash := head.Hash()
	blockNumber := head.Number
	parentHash := head.ParentHash
	if err := dbInstance.NewBlock(blockNumber, parentHash, hash); err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := dbInstance.Put(hash, kv.key, kv.value); err != nil {
			return err
		}
	}
	return nil
}

func newBlockRecognizedByFlush(hash common.Hash, parentHead *types.Header, kvs []kv) error {
	if err := newBlockUnRecognized(parentHead, kvs); err != nil {
		return err
	}
	return dbInstance.Flush(hash, new(big.Int).Add(parentHead.Number, common.Big1))
}

func newBlockUnRecognized(parentHead *types.Header, kvs []kv) error {
	blockNumber := new(big.Int).Add(parentHead.Number, common.Big1)
	parentHash := parentHead.Hash()
	if err := dbInstance.NewBlock(blockNumber, parentHash, common.ZeroHash); err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := dbInstance.Put(common.ZeroHash, kv.key, kv.value); err != nil {
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
