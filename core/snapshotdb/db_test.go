package snapshotdb

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestRecover(t *testing.T) {
	//os.RemoveAll(dbpath)
	//if err := initDB(); err != nil {
	//	t.Error(err)
	//	return
	//}
	//defer dbInstance.Clear()
	//var (
	//	recognizedHash                                       = generateHash("recognizedHash")
	//	parenthash                                           common.Hash
	//	baseDBArr, commitArr, recognizedArr, unrecognizedArr []kv
	//	base, high                                           int64
	//	commit, recognized, unrecognized                     *blockData
	//)
	//{
	//	commitHash := recognizedHash
	//	var str string
	//	for i := 0; i < 4; i++ {
	//		str += "a"
	//		baseDBArr = append(baseDBArr, kv{key: []byte(str), value: []byte(str)})
	//	}
	//	baseDBArr = append(baseDBArr, kv{key: []byte("d"), value: []byte("d")})
	//	if err := newBlockBaseDB(big.NewInt(1), parenthash, commitHash, baseDBArr); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	parenthash = commitHash
	//}
	//{
	//	commitHash := generateHash("recognizedHash3")
	//	str := "a"
	//	for i := 0; i < 4; i++ {
	//		str += "c"
	//		commitArr = append(commitArr, kv{key: []byte(str), value: []byte(str)})
	//	}
	//	commitArr = append(commitArr, kv{key: []byte("ddd"), value: []byte("ddd")})
	//	if err := newBlockCommited(big.NewInt(2), parenthash, commitHash, commitArr); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	commit = dbInstance.committed[0]
	//	parenthash = commitHash
	//}
	//{
	//	commitHash := generateHash("recognizedHash4")
	//	str := "a"
	//	for i := 0; i < 4; i++ {
	//		str += "e"
	//		recognizedArr = append(recognizedArr, kv{key: []byte(str), value: []byte(str)})
	//	}
	//	recognizedArr = append(recognizedArr, kv{key: []byte("ee"), value: []byte("ee")})
	//
	//	if err := newBlockRecognizedDirect(big.NewInt(3), parenthash, commitHash, recognizedArr); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//
	//	rg, ok := dbInstance.unCommit.blocks[generateHash("recognizedHash4")]
	//	if !ok {
	//		t.Error("not found recognizedHash4")
	//	}
	//	recognized = rg
	//	parenthash = commitHash
	//}
	//{
	//	str := "a"
	//	for i := 0; i < 4; i++ {
	//		str += "f"
	//		unrecognizedArr = append(unrecognizedArr, kv{key: []byte(str), value: []byte(str)})
	//	}
	//	unrecognizedArr = append(unrecognizedArr, kv{key: []byte("ff"), value: []byte("ff")})
	//
	//	if err := newBlockUnRecognized(big.NewInt(4), parenthash, unrecognizedArr); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	unrecognized = dbInstance.unCommit.blocks[dbInstance.getUnRecognizedHash()]
	//}
	//base = dbInstance.current.BaseNum.Int64()
	//high = dbInstance.current.HighestNum.Int64()
	//if err := dbInstance.Close(); err != nil {
	//	t.Error(err)
	//}
	//dbInstance = nil
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
	//	db := new(snapshotDB)
	//	if err := db.recover(s); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	dbInstance = db
	//	defer dbInstance.Clear()
	//}
	//
	//if dbInstance.path != dbpath {
	//	t.Error("path is wrong", dbInstance.path, dbpath)
	//	return
	//}
	//if dbInstance.current.BaseNum.Int64() != base {
	//	t.Error("BaseNum is wrong", dbInstance.current.BaseNum.Int64(), base)
	//	return
	//}
	//if dbInstance.current.HighestNum.Int64() != high {
	//	t.Error("HighestNum is wrong", dbInstance.current.HighestNum.Int64(), high)
	//	return
	//}
	//if dbInstance.current.path != getCurrentPath(dbpath) {
	//	t.Error("current path is wrong", dbInstance.current.path, getCurrentPath(dbpath))
	//	return
	//}
	//for _, value := range baseDBArr {
	//	v, err := dbInstance.baseDB.Get(value.key, nil)
	//	if err != nil {
	//		t.Error("should be nil", err)
	//		return
	//	}
	//	if bytes.Compare(v, value.value) != 0 {
	//		t.Error("should be equal", v, []byte(value.value))
	//		return
	//	}
	//}
	//oldarr := []*blockData{
	//	unrecognized,
	//	recognized,
	//	commit,
	//}
	//rg, ok := dbInstance.unCommit.blocks[generateHash("recognizedHash4")]
	//if !ok {
	//	t.Error("recognizedHash4 not found")
	//}
	//rg2 := rg
	//newarr := []*blockData{
	//	dbInstance.unCommit.blocks[dbInstance.getUnRecognizedHash()],
	//	rg2,
	//	dbInstance.committed[0],
	//}
	//
	//for i := 0; i < 3; i++ {
	//	if i == 0 {
	//		if newarr[i].BlockHash != common.ZeroHash {
	//			t.Error("unRecognized block hash must nil", i, newarr[i].BlockHash, oldarr[i].BlockHash)
	//			return
	//		}
	//	} else {
	//		if oldarr[i].BlockHash != newarr[i].BlockHash {
	//			t.Error("block hash must compare", i, oldarr[i].BlockHash.String(), newarr[i].BlockHash.String())
	//			return
	//		}
	//	}
	//
	//	if oldarr[i].ParentHash != newarr[i].ParentHash {
	//		t.Error("ParentHash  must compare", i)
	//		return
	//	}
	//	if oldarr[i].kvHash != newarr[i].kvHash {
	//		t.Error("kvHash  must compare", i)
	//		return
	//	}
	//	if oldarr[i].readOnly != newarr[i].readOnly {
	//		t.Error("readOnly  must compare", i, oldarr[i].readOnly, newarr[i].readOnly)
	//		return
	//	}
	//	itr := oldarr[i].data.NewIterator(nil)
	//	for itr.Next() {
	//		v, err := newarr[i].data.Get(itr.Key())
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//		if bytes.Compare(v, itr.Value()) != 0 {
	//			t.Error("kv  must compare", i)
	//			return
	//		}
	//	}
	//}
	//
	//if err := dbInstance.Put(common.ZeroHash, []byte("dddd"), []byte("dddd")); err != nil {
	//	t.Error(err)
	//}
	//if err := dbInstance.Flush(generateHash("flush"), big.NewInt(4)); err != nil {
	//	t.Error(err)
	//}
	//if err := dbInstance.Commit(generateHash("recognizedHash4")); err != nil {
	//	t.Error(err)
	//}
	//if err := dbInstance.Compaction(); err != nil {
	//	t.Error(err)
	//}
}

func TestRMOldRecognizedBlockData(t *testing.T) {
	//initDB()
	//defer dbInstance.Clear()
	//var (
	//	recognizedHash       = generateHash("recognizedHash")
	//	parenthash           common.Hash
	//	baseDBArr, commitArr []string
	//)
	//{
	//	commitHash := recognizedHash
	//	if err := dbInstance.NewBlock(big.NewInt(1), parenthash, commitHash); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	var str string
	//	for i := 0; i < 4; i++ {
	//		str += "a"
	//		baseDBArr = append(baseDBArr, str)
	//		if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
	//			t.Error(err)
	//			return
	//		}
	//	}
	//	if err := dbInstance.Put(commitHash, []byte("d"), []byte("d")); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	if err := dbInstance.Commit(commitHash); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	parenthash = commitHash
	//	dbInstance.Compaction()
	//}
	//{
	//	commitHash := generateHash("recognizedHash3")
	//	if err := dbInstance.NewBlock(big.NewInt(2), parenthash, commitHash); err != nil {
	//		t.Error(err)
	//	}
	//	str := "a"
	//	for i := 0; i < 4; i++ {
	//		str += "c"
	//		commitArr = append(commitArr, str)
	//		if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
	//			t.Error(err)
	//		}
	//	}
	//	if err := dbInstance.Put(commitHash, []byte("ddd"), []byte("ddd")); err != nil {
	//		t.Error(err)
	//	}
	//	if err := dbInstance.Commit(commitHash); err != nil {
	//		t.Error(err)
	//	}
	//	parenthash = commitHash
	//}
	//{
	//	commitHash := generateHash("recognizedHash4")
	//	if err := dbInstance.NewBlock(big.NewInt(1), parenthash, commitHash); err != nil {
	//		t.Error(err)
	//	}
	//	str := "a"
	//	for i := 0; i < 4; i++ {
	//		str += "c"
	//		commitArr = append(commitArr, str)
	//		if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
	//			t.Error(err)
	//		}
	//	}
	//	if err := dbInstance.Put(commitHash, []byte("ddd"), []byte("ddd")); err != nil {
	//		t.Error(err)
	//	}
	//}
	//{
	//	commitHash := generateHash("recognizedHash5")
	//	if err := dbInstance.NewBlock(big.NewInt(2), parenthash, commitHash); err != nil {
	//		t.Error(err)
	//	}
	//	str := "a"
	//	for i := 0; i < 4; i++ {
	//		str += "c"
	//		commitArr = append(commitArr, str)
	//		if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
	//			t.Error(err)
	//		}
	//	}
	//	if err := dbInstance.Put(commitHash, []byte("ddd"), []byte("ddd")); err != nil {
	//		t.Error(err)
	//	}
	//}
	//if err := dbInstance.rmOldRecognizedBlockData(); err != nil {
	//	t.Error(err)
	//}
	//if len(dbInstance.unCommit.blocks) != 0 {
	//	t.Error("not rm old data")
	//}
}

func randomString2(s string) []byte {
	b := new(bytes.Buffer)
	if s != "" {
		b.Write([]byte(s))
	}
	for i := 0; i < 4; i++ {
		b.WriteByte(' ' + byte(rand.Int()))
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

type testchain struct {
	h *types.Header
}

func (c *testchain) CurrentHeader() *types.Header {
	return c.h
}

func (c *testchain) GetHeaderByHash(hash common.Hash) *types.Header {
	return c.h
}

func newBlockBaseDB(blockNumber *big.Int, parentHash common.Hash, hash common.Hash, kvs []kv) error {
	tchain := new(testchain)
	tchain.h = &types.Header{ParentHash: hash, Number: blockNumber}
	blockchain = tchain
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

func newBlockCommited(blockNumber *big.Int, parentHash common.Hash, hash common.Hash, kvs []kv) error {
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

func newBlockRecognizedDirect(blockNumber *big.Int, parentHash common.Hash, hash common.Hash, kvs []kv) error {
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

func newBlockRecognizedByFlush(blockNumber *big.Int, parentHash common.Hash, hash common.Hash, kvs []kv) error {
	if err := newBlockUnRecognized(blockNumber, parentHash, kvs); err != nil {
		return err
	}
	return dbInstance.Flush(hash, blockNumber)
}

func newBlockUnRecognized(blockNumber *big.Int, parentHash common.Hash, kvs []kv) error {
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

func TestCheckHashChain(t *testing.T) {
	initDB()
	defer dbInstance.Clear()

	//1
	if err := newBlockBaseDB(big.NewInt(1), generateHash(fmt.Sprint(0)), generateHash(fmt.Sprint(1)), generatekv(1)); err != nil {
		t.Error(err)
	}

	//2-11
	for i := 2; i < 11; i++ {
		if err := newBlockCommited(big.NewInt(int64(i)), generateHash(fmt.Sprint(i-1)), generateHash(fmt.Sprint(i)), generatekv(1)); err != nil {
			t.Error(err)
		}
	}
	//12-20
	for i := 11; i < 21; i++ {
		if err := newBlockRecognizedDirect(big.NewInt(int64(i)), generateHash(fmt.Sprint(i-1)), generateHash(fmt.Sprint(i)), generatekv(1)); err != nil {
			t.Error(err)
		}
	}

	if err := newBlockUnRecognized(big.NewInt(int64(21)), generateHash(fmt.Sprint(20)), generatekv(1)); err != nil {
		t.Error(err)
	}
	t.Run("find from recognized", func(t *testing.T) {
		for i := 11; i < 21; i++ {
			location, ok := dbInstance.checkHashChain(generateHash(fmt.Sprint(i)))
			if !ok {
				t.Error("should be ok")
			}
			if location != hashLocationUnCommitted {
				t.Error("should be locate Recognized", location)
			}
		}

	})

	t.Run("find from commit", func(t *testing.T) {
		for i := 2; i < 11; i++ {
			location, ok := dbInstance.checkHashChain(generateHash(fmt.Sprint(i)))
			if !ok {
				t.Error("should be ok")
			}
			if location != hashLocationCommitted {
				t.Error("should be locate Recognized", location)
			}
		}
	})

	t.Run("find from unrecognized", func(t *testing.T) {
		location, ok := dbInstance.checkHashChain(common.ZeroHash)
		if !ok {
			t.Error("should be ok")
		}
		if location != hashLocationUnCommitted {
			t.Error("should be locate Recognized", location)
		}
	})

	t.Run("not found", func(t *testing.T) {
		location, ok := dbInstance.checkHashChain(generateHash(fmt.Sprint(1)))
		if !ok {
			t.Error("should be ok")
		}
		if location != hashLocationNotFound {
			t.Error("should be locate Recognized", location)
		}
	})
}
