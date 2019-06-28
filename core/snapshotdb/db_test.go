package snapshotdb

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestRecover(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		recognizedHash                                       = generateHash("recognizedHash")
		parenthash                                           common.Hash
		baseDBArr, commitArr, recognizedArr, unrecognizedArr []string
		base, high                                           int64
		commit, recognized, unrecognized                     blockData
	)
	{
		commitHash := recognizedHash
		if err := dbInstance.NewBlock(big.NewInt(1), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		var str string
		for i := 0; i < 4; i++ {
			str += "a"
			baseDBArr = append(baseDBArr, str)
			if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if err := dbInstance.Put(commitHash, []byte("d"), []byte("d")); err != nil {
			t.Fatal(err)
		}
		if err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		parenthash = commitHash
		dbInstance.Compaction()
	}
	{
		commitHash := generateHash("recognizedHash3")
		if err := dbInstance.NewBlock(big.NewInt(2), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "c"
			commitArr = append(commitArr, str)
			if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if err := dbInstance.Put(commitHash, []byte("ddd"), []byte("ddd")); err != nil {
			t.Fatal(err)
		}
		if err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		commit = dbInstance.committed[0]
		parenthash = commitHash
	}
	{
		commitHash := generateHash("recognizedHash4")
		if err := dbInstance.NewBlock(big.NewInt(3), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "e"
			recognizedArr = append(recognizedArr, str)
			if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if err := dbInstance.Put(commitHash, []byte("ee"), []byte("ee")); err != nil {
			t.Fatal(err)
		}
		recognized = dbInstance.recognized[generateHash("recognizedHash4")]
		parenthash = commitHash
	}
	{
		if err := dbInstance.NewBlock(big.NewInt(4), parenthash, common.ZeroHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "f"
			unrecognizedArr = append(unrecognizedArr, str)
			if err := dbInstance.Put(common.ZeroHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if err := dbInstance.Put(common.ZeroHash, []byte("ff"), []byte("ff")); err != nil {
			t.Fatal(err)
		}
		unrecognized = *dbInstance.unRecognized
	}
	base = dbInstance.current.BaseNum.Int64()
	high = dbInstance.current.HighestNum.Int64()
	if err := dbInstance.Close(); err != nil {
		t.Error(err)
	}
	dbInstance = nil
	s, err := openFile(dbpath, false)
	if err != nil {
		t.Fatal(err)
	}
	fds, err := s.List(TypeCurrent)
	if err != nil {
		t.Fatal(err)
	}
	if len(fds) > 0 {
		db := new(snapshotDB)
		if err := db.recover(s); err != nil {
			t.Fatal(err)
		}
		dbInstance = db
		defer dbInstance.Clear()
	}

	if dbInstance.path != dbpath {
		t.Error("path is wrong", dbInstance.path, dbpath)
	}
	if dbInstance.current.BaseNum.Int64() != base {
		t.Error("BaseNum is wrong", dbInstance.current.BaseNum.Int64(), base)
	}
	if dbInstance.current.HighestNum.Int64() != high {
		t.Error("HighestNum is wrong", dbInstance.current.HighestNum.Int64(), high)
	}
	if dbInstance.current.path != getCurrentPath(dbpath) {
		t.Error("current path is wrong", dbInstance.current.path, getCurrentPath(dbpath))
	}
	for _, value := range baseDBArr {
		v, err := dbInstance.baseDB.Get([]byte(value), nil)
		if err != nil {
			t.Error("should be nil", err)
		}
		if bytes.Compare(v, []byte(value)) != 0 {
			t.Error("should be equal", v, []byte(value))
		}
	}
	oldarr := []blockData{
		unrecognized,
		recognized,
		commit,
	}
	newarr := []blockData{
		*dbInstance.unRecognized,
		dbInstance.recognized[generateHash("recognizedHash4")],
		dbInstance.committed[0],
	}

	for i := 0; i < 3; i++ {
		if i == 0 {
			if newarr[i].BlockHash != common.ZeroHash {
				t.Error("unRecognized block hash must nil", i, newarr[i].BlockHash, oldarr[i].BlockHash)
			}
		} else {
			if oldarr[i].BlockHash != newarr[i].BlockHash {
				t.Error("block hash must compare", i, oldarr[i].BlockHash.String(), newarr[i].BlockHash.String())
			}
		}

		if oldarr[i].ParentHash != newarr[i].ParentHash {
			t.Error("ParentHash  must compare", i)
		}
		if oldarr[i].kvHash != newarr[i].kvHash {
			t.Error("kvHash  must compare", i)
		}
		if oldarr[i].readOnly != newarr[i].readOnly {
			t.Error("readOnly  must compare", i, oldarr[i].readOnly, newarr[i].readOnly)
		}
		itr := oldarr[i].data.NewIterator(nil)
		for itr.Next() {
			v, err := newarr[i].data.Get(itr.Key())
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(v, itr.Value()) != 0 {
				t.Error("kv  must compare", i)
			}
		}
	}

	if err := dbInstance.Put(common.ZeroHash, []byte("dddd"), []byte("dddd")); err != nil {
		t.Error(err)
	}
	if err := dbInstance.Flush(generateHash("flush"), big.NewInt(4)); err != nil {
		t.Error(err)
	}
	if err := dbInstance.Commit(generateHash("recognizedHash4")); err != nil {
		t.Error(err)
	}
	if err := dbInstance.Compaction(); err != nil {
		t.Error(err)
	}
}

func TestRMOldRecognizedBlockData(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		recognizedHash       = generateHash("recognizedHash")
		parenthash           common.Hash
		baseDBArr, commitArr []string
	)
	{
		commitHash := recognizedHash
		if err := dbInstance.NewBlock(big.NewInt(1), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		var str string
		for i := 0; i < 4; i++ {
			str += "a"
			baseDBArr = append(baseDBArr, str)
			if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if err := dbInstance.Put(commitHash, []byte("d"), []byte("d")); err != nil {
			t.Fatal(err)
		}
		if err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		parenthash = commitHash
		dbInstance.Compaction()
	}
	{
		commitHash := generateHash("recognizedHash3")
		if err := dbInstance.NewBlock(big.NewInt(2), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "c"
			commitArr = append(commitArr, str)
			if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if err := dbInstance.Put(commitHash, []byte("ddd"), []byte("ddd")); err != nil {
			t.Fatal(err)
		}
		if err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		parenthash = commitHash
	}
	{
		commitHash := generateHash("recognizedHash4")
		if err := dbInstance.NewBlock(big.NewInt(1), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "c"
			commitArr = append(commitArr, str)
			if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if err := dbInstance.Put(commitHash, []byte("ddd"), []byte("ddd")); err != nil {
			t.Fatal(err)
		}
	}
	{
		commitHash := generateHash("recognizedHash5")
		if err := dbInstance.NewBlock(big.NewInt(2), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "c"
			commitArr = append(commitArr, str)
			if err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if err := dbInstance.Put(commitHash, []byte("ddd"), []byte("ddd")); err != nil {
			t.Fatal(err)
		}
	}
	if err := dbInstance.rmOldRecognizedBlockData(); err != nil {
		t.Error(err)
	}
	if len(dbInstance.recognized) != 0 {
		t.Error("not rm old data")
	}
}

type kv struct {
	key   []byte
	value []byte
}

func randomString2() []byte {
	b := new(bytes.Buffer)
	for i := 0; i < 4; i++ {
		b.WriteByte(' ' + byte(rand.Int()))
	}
	return b.Bytes()
}

func generatekv(n int) []kv {
	rand.Seed(time.Now().UnixNano())
	kvs := make([]kv, n)
	for i := 0; i < n; i++ {
		kvs[i] = kv{
			key:   randomString2(),
			value: randomString2(),
		}
	}
	return kvs
}

func newBlockBaseDB(blockNumber *big.Int, parentHash common.Hash, hash common.Hash, kvs []kv) error {
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
		t.Fatal(err)
	}

	//2-11
	for i := 2; i < 11; i++ {
		if err := newBlockCommited(big.NewInt(int64(i)), generateHash(fmt.Sprint(i-1)), generateHash(fmt.Sprint(i)), generatekv(1)); err != nil {
			t.Fatal(err)
		}
	}
	//12-20
	for i := 11; i < 21; i++ {
		if err := newBlockRecognizedDirect(big.NewInt(int64(i)), generateHash(fmt.Sprint(i-1)), generateHash(fmt.Sprint(i)), generatekv(1)); err != nil {
			t.Fatal(err)
		}
	}

	if err := newBlockUnRecognized(big.NewInt(int64(21)), generateHash(fmt.Sprint(20)), generatekv(1)); err != nil {
		t.Fatal(err)
	}
	t.Run("find from recognized", func(t *testing.T) {
		for i := 11; i < 21; i++ {
			location, ok := dbInstance.checkHashChain(generateHash(fmt.Sprint(i)))
			if !ok {
				t.Error("should be ok")
			}
			if location != hashLocationRecognized {
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
		if location != hashLocationUnRecognized {
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
