package snapshotdb

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
)

var (
	parentHash  = rlpHash("parentHash")
	currentHash = rlpHash("currentHash")
)

func TestSnapshotDB_NewBlock(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	t.Run("new recognized block", func(t *testing.T) {
		_, err := dbInstance.NewBlock(big.NewInt(30), parentHash, currentHash)
		if err != nil {
			t.Error(err)
		}
		bd, ok := dbInstance.recognized[currentHash]
		if !ok {
			t.Fatal("must find recognized")
		}
		if bd.ParentHash != parentHash {
			t.Fatal("parentHash must same:", bd.ParentHash, parentHash)
		}
		if bd.Number.Cmp(big.NewInt(30)) != 0 {
			t.Fatal("block number must same:", bd.Number, big.NewInt(30))
		}
		if bd.BlockHash != currentHash {
			t.Fatal("BlockHash must right:", bd.BlockHash, currentHash)
		}
	})
	t.Run("new unrecognized block", func(t *testing.T) {
		_, err := dbInstance.NewBlock(big.NewInt(30), parentHash, common.ZeroHash)
		if err != nil {
			t.Error(err)
		}
		bd := dbInstance.unRecognized
		if bd.ParentHash != parentHash {
			t.Fatal("parentHash must same:", bd.ParentHash, parentHash)
		}
		if bd.Number.Cmp(big.NewInt(30)) != 0 {
			t.Fatal("block number must same:", bd.Number, big.NewInt(30))
		}
	})
}

func TestSnapshotDB_Get(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr            = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}
		recognizedHash = rlpHash("recognizedHash")
		commitHash     = rlpHash("commitHash")
	)
	{
		//unRecognized
		unRecognized := blockData{
			ParentHash: recognizedHash,
			Number:     big.NewInt(50),
			data:       memdb.New(DefaultComparer, 10),
			readOnly:   false,
		}
		unRecognized.data.Put(arr[0], arr[0])
		dbInstance.unRecognized = &unRecognized

		//recognized
		Recognized := blockData{
			ParentHash: commitHash,
			Number:     big.NewInt(50),
			data:       memdb.New(DefaultComparer, 10),
			readOnly:   false,
			BlockHash:  recognizedHash,
		}
		Recognized.data.Put(arr[1], arr[1])
		dbInstance.recognized[recognizedHash] = Recognized

		//commit
		commit := blockData{
			ParentHash: parentHash,
			Number:     big.NewInt(50),
			data:       memdb.New(DefaultComparer, 10),
			readOnly:   false,
			BlockHash:  commitHash,
		}
		commit.data.Put(arr[2], arr[2])
		dbInstance.committed = append(dbInstance.committed, commit)

		//baseDB
		dbInstance.baseDB.Put(arr[3], arr[3], nil)
	}

	t.Run("with not hash", func(t *testing.T) {
		t.Run("must find", func(t *testing.T) {
			for _, key := range arr {
				val, err := dbInstance.Get(common.ZeroHash, key)
				if err != nil {
					t.Error(err)
				}
				if bytes.Compare(key, val) != 0 {
					t.Error("must find key")
				}
			}
		})
		t.Run("not find", func(t *testing.T) {
			_, err := dbInstance.Get(common.ZeroHash, []byte("e"))
			if err == nil {
				t.Error(err)
			}
		})
	})
	t.Run("with hash", func(t *testing.T) {
		t.Run("can't get unrecognized BlockData", func(t *testing.T) {
			_, err := dbInstance.Get(recognizedHash, []byte("a"))
			if err == nil {
				t.Error(err)
			}
		})
		t.Run("get from recognized BlockData", func(t *testing.T) {
			for _, key := range arr[1 : len(arr)-1] {
				val, err := dbInstance.Get(recognizedHash, key)
				if err != nil {
					t.Error(err)
				}
				if bytes.Compare(key, val) != 0 {
					t.Error("must find key")
				}
			}
		})
		t.Run("get from committed blockData", func(t *testing.T) {
			for _, key := range arr[2 : len(arr)-1] {
				val, err := dbInstance.Get(commitHash, key)
				if err != nil {
					t.Error(err)
				}
				if bytes.Compare(key, val) != 0 {
					t.Error("must find key")
				}
			}
		})
	})

}

func TestSnapshotDB_GetFromCommitedBlock(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr        = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}
		commitHash = rlpHash("commitHash")
	)
	{
		//commit
		commit := blockData{
			ParentHash: parentHash,
			Number:     big.NewInt(50),
			data:       memdb.New(DefaultComparer, 10),
			readOnly:   false,
			BlockHash:  commitHash,
		}
		commit.data.Put(arr[2], arr[2])
		dbInstance.committed = append(dbInstance.committed, commit)

		//baseDB
		dbInstance.baseDB.Put(arr[3], arr[3], nil)
	}

	t.Run("should get", func(t *testing.T) {
		for _, key := range arr[2:3] {
			val, err := dbInstance.GetFromCommittedBlock(key)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(key, val) != 0 {
				t.Error("must find key")
			}
		}
	})
	t.Run("not found", func(t *testing.T) {
		_, err := dbInstance.GetFromCommittedBlock(arr[1])
		if err == nil {
			t.Error(err)
		}
		if err != ErrNotFound {
			t.Error("err should be ErrNotFound")
		}
	})
}

func TestSnapshotDB_Del(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr                   = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")}
		recognizedHash        = rlpHash("recognizedHash")
		recognizedByFlushHash = rlpHash("recognizedByFlush")
		commitHash            = rlpHash("commitHash")
	)
	{
		//unRecognized
		unRecognized := blockData{
			ParentHash: recognizedHash,
			Number:     big.NewInt(50),
			data:       memdb.New(DefaultComparer, 10),
			readOnly:   false,
		}
		unRecognized.data.Put(arr[0], arr[0])
		dbInstance.unRecognized = &unRecognized
		f, err := ioutil.TempFile(os.TempDir(), "test_del*.log")
		if err != nil {
			t.Error(err)
		}
		dbInstance.journalw[dbInstance.getUnRecognizedHash()] = newJournalWriter(f)
	}
	{
		//recognized
		Recognized := blockData{
			ParentHash: commitHash,
			Number:     big.NewInt(50),
			data:       memdb.New(DefaultComparer, 10),
			readOnly:   false,
			BlockHash:  recognizedHash,
		}
		Recognized.data.Put(arr[1], arr[1])
		dbInstance.recognized[recognizedHash] = Recognized
		f, err := ioutil.TempFile(os.TempDir(), "test_del*.log")
		if err != nil {
			t.Error(err)
		}
		dbInstance.journalw[recognizedHash] = newJournalWriter(f)
	}
	{
		//recognized by flush
		Recognized2 := blockData{
			ParentHash: commitHash,
			Number:     big.NewInt(50),
			data:       memdb.New(DefaultComparer, 10),
			readOnly:   true,
			BlockHash:  recognizedByFlushHash,
		}
		Recognized2.data.Put(arr[4], arr[4])
		dbInstance.recognized[recognizedByFlushHash] = Recognized2
		f, err := ioutil.TempFile(os.TempDir(), "test_del*.log")
		if err != nil {
			t.Error(err)
		}
		dbInstance.journalw[recognizedByFlushHash] = newJournalWriter(f)
	}
	{
		//commit
		commit := blockData{
			ParentHash: parentHash,
			Number:     big.NewInt(50),
			data:       memdb.New(DefaultComparer, 10),
			readOnly:   false,
			BlockHash:  commitHash,
		}
		commit.data.Put(arr[2], arr[2])
		dbInstance.committed = append(dbInstance.committed, commit)
		f, err := ioutil.TempFile(os.TempDir(), "test_del*.log")
		if err != nil {
			t.Error(err)
		}
		dbInstance.journalw[commitHash] = newJournalWriter(f)
	}
	{
		//baseDB
		dbInstance.baseDB.Put(arr[3], arr[3], nil)
	}

	t.Run("delete unrecognized", func(t *testing.T) {
		ok, err := dbInstance.Del(common.ZeroHash, arr[0])
		if !ok {
			t.Error("return must be true")
		}
		if err != nil {
			t.Error("err must be nil", err)
		}
	})
	t.Run("delete recognized", func(t *testing.T) {
		ok, err := dbInstance.Del(recognizedHash, arr[1])
		if !ok {
			t.Error("return must be true")
		}
		if err != nil {
			t.Error("err must be nil", err)
		}
	})

	t.Run("can't delete readonly", func(t *testing.T) {
		ok, err := dbInstance.Del(recognizedByFlushHash, arr[4])
		if ok {
			t.Error("return must be false")
		}
		if err == nil {
			t.Error("err must not nil", err)
		}
	})
	t.Run("can't delete commit", func(t *testing.T) {
		ok, err := dbInstance.Del(commitHash, arr[2])
		if ok {
			t.Error("return must be false")
		}
		if err == nil {
			t.Error("err must not nil", err)
		}
	})
}

func TestSnapshotDB_Has(t *testing.T) {
	//same as get
}

func TestSnapshotDB_Ranking(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		recognizedHash = rlpHash("recognizedHash")
		parenthash     common.Hash
		arr            []string
	)
	{
		commitHash := recognizedHash
		if _, err := dbInstance.NewBlock(big.NewInt(1), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		var str string
		for i := 0; i < 4; i++ {
			str += "a"
			arr = append(arr, str)
			if _, err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(commitHash, []byte("d"), []byte("d")); err != nil {
			t.Fatal(err)
		}
		if _, err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		parenthash = commitHash
		dbInstance.Compaction()
	}

	{
		commitHash := rlpHash("recognizedHash2")
		if _, err := dbInstance.NewBlock(big.NewInt(2), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "b"
			arr = append(arr, str)
			if _, err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(commitHash, []byte("dd"), []byte("dd")); err != nil {
			t.Fatal(err)
		}
		if _, err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		parenthash = commitHash
		dbInstance.Compaction()
	}

	{
		commitHash := rlpHash("recognizedHash3")
		if _, err := dbInstance.NewBlock(big.NewInt(3), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "c"
			arr = append(arr, str)
			if _, err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(commitHash, []byte("ddd"), []byte("ddd")); err != nil {
			t.Fatal(err)
		}
		if _, err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		parenthash = commitHash
	}
	{
		commitHash := rlpHash("recognizedHash4")
		if _, err := dbInstance.NewBlock(big.NewInt(4), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "e"
			arr = append(arr, str)
			if _, err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(commitHash, []byte("ee"), []byte("ee")); err != nil {
			t.Fatal(err)
		}
		parenthash = commitHash
	}
	{
		if _, err := dbInstance.NewBlock(big.NewInt(5), parenthash, common.ZeroHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "f"
			arr = append(arr, str)
			if _, err := dbInstance.Put(common.ZeroHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(common.ZeroHash, []byte("ff"), []byte("ff")); err != nil {
			t.Fatal(err)
		}
	}
	t.Run("with hash", func(t *testing.T) {
		t.Run("form recognized", func(t *testing.T) {
			commitHash := rlpHash("recognizedHash4")
			itr := dbInstance.Ranking(commitHash, []byte("a"), 100)
			for _, value := range arr[0:15] {
				if !itr.Next() {
					t.Error("it's must can itr")
				}
				if value != string(itr.Key()) {
					t.Errorf("itr return wrong key :%s,should return:%s ", string(itr.Key()), value)
				}
			}
			itr.Release()
		})
		t.Run("form commit", func(t *testing.T) {
			commitHash := rlpHash("recognizedHash3")
			itr := dbInstance.Ranking(commitHash, []byte("a"), 100)
			for _, value := range arr[0:11] {
				if !itr.Next() {
					t.Error("it's must can itr")
				}
				if value != string(itr.Key()) {
					t.Errorf("itr return wrong key :%s,should return:%s ", string(itr.Key()), value)
				}
			}
			itr.Release()
		})

	})
	t.Run("with out hash", func(t *testing.T) {
		itr := dbInstance.Ranking(common.ZeroHash, []byte("a"), 100)
		for _, value := range arr[0:19] {
			if !itr.Next() {
				t.Error("it's must can itr")
			}
			if value != string(itr.Key()) {
				t.Errorf("itr return wrong key :%s,should return:%s ", string(itr.Key()), value)
			}
		}
		itr.Release()
	})
}

func TestSnapshotDB_WalkBaseDB(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		recognizedHash = rlpHash("recognizedHash")
		parenthash     common.Hash
		arr            []string
	)
	{
		commitHash := recognizedHash
		if _, err := dbInstance.NewBlock(big.NewInt(1), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		var str string
		for i := 0; i < 4; i++ {
			str += "a"
			arr = append(arr, str)
			if _, err := dbInstance.Put(commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(commitHash, []byte("d"), []byte("d")); err != nil {
			t.Fatal(err)
		}
		if _, err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		parenthash = commitHash
		dbInstance.Compaction()
	}
	prefix := util.BytesPrefix([]byte("a"))
	f := func(num *big.Int, iter iterator.Iterator) error {
		if num.Int64() != 1 {
			return fmt.Errorf("basenum is wrong:%v,should be 1", num)
		}
		var i int
		for iter.Next() {
			if arr[i] != string(iter.Key()) {
				return fmt.Errorf("itr return wrong key :%s,should return:%s ,index:%d", string(iter.Key()), arr[i], i)
			}
			i++
		}
		return nil
	}
	if err := dbInstance.WalkBaseDB(prefix, f); err != nil {
		t.Error(err)
	}
}

func TestSnapshotDB_Clear(t *testing.T) {

}

func TestSnapshotDB_GetLastKVHash(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr            = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")}
		recognizedHash = rlpHash("recognizedHash")
		commitHash     = rlpHash("commitHash")
	)
	{
		dbInstance.NewBlock(big.NewInt(10), recognizedHash, common.ZeroHash)
		dbInstance.Put(common.ZeroHash, arr[0], arr[0])
		dbInstance.Put(common.ZeroHash, arr[1], arr[1])
	}
	{
		dbInstance.NewBlock(big.NewInt(10), commitHash, recognizedHash)
		dbInstance.Put(recognizedHash, arr[2], arr[2])
		dbInstance.Put(recognizedHash, arr[3], arr[3])
	}
	t.Run("get from unRecognized", func(t *testing.T) {
		var lastkvhash common.Hash
		kvhash := dbInstance.GetLastKVHash(common.ZeroHash)
		lastkvhash = dbInstance.generateKVHash(arr[0], arr[0], lastkvhash)
		lastkvhash = dbInstance.generateKVHash(arr[1], arr[1], lastkvhash)
		if bytes.Compare(kvhash, lastkvhash.Bytes()) != 0 {
			t.Error("kv hash must same", lastkvhash, kvhash)
		}
	})
	t.Run("get from recognized", func(t *testing.T) {
		var lastkvhash common.Hash
		kvhash := dbInstance.GetLastKVHash(recognizedHash)
		lastkvhash = dbInstance.generateKVHash(arr[2], arr[2], lastkvhash)
		lastkvhash = dbInstance.generateKVHash(arr[3], arr[3], lastkvhash)
		if bytes.Compare(kvhash, lastkvhash.Bytes()) != 0 {
			t.Error("kv hash must same", lastkvhash, kvhash)
		}
	})
}

func TestSnapshotDB_BaseNum(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	_, err := dbInstance.BaseNum()
	if err != nil {
		t.Error(err)
	}
}

func TestSnapshotDB_Compaction(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		recognizedHash = rlpHash("recognizedHash")
		commitHash     common.Hash
		parenthash     common.Hash
	)
	{
		_, err := dbInstance.NewBlock(big.NewInt(1), commitHash, recognizedHash)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 3000; i++ {
			_, err := dbInstance.Put(recognizedHash, []byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)))
			if err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Commit(recognizedHash); err != nil {
			t.Fatal(err)
		}
	}
	{
		currenthash := rlpHash(fmt.Sprint(2))
		if _, err := dbInstance.NewBlock(big.NewInt(int64(2)), recognizedHash, currenthash); err != nil {
			t.Fatal(err)
		}
		for i := 3000; i < 3100; i++ {
			if _, err := dbInstance.Put(currenthash, []byte(fmt.Sprint(i)), []byte(fmt.Sprint(i))); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Commit(currenthash); err != nil {
			t.Fatal(err)
		}
		parenthash = currenthash

		currenthash = rlpHash(fmt.Sprint(3))
		if _, err := dbInstance.NewBlock(big.NewInt(int64(3)), parenthash, currenthash); err != nil {
			t.Fatal(err)
		}
		for i := 3100; i < 3200; i++ {
			if _, err := dbInstance.Put(currenthash, []byte(fmt.Sprint(i)), []byte(fmt.Sprint(i))); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Commit(currenthash); err != nil {
			t.Fatal(err)
		}
		parenthash = currenthash

		currenthash = rlpHash(fmt.Sprint(4))
		if _, err := dbInstance.NewBlock(big.NewInt(int64(4)), parenthash, currenthash); err != nil {
			t.Fatal(err)
		}
		for i := 3200; i < 4998; i++ {
			if _, err := dbInstance.Put(currenthash, []byte(fmt.Sprint(i)), []byte(fmt.Sprint(i))); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Commit(currenthash); err != nil {
			t.Fatal(err)
		}
		parenthash = currenthash
	}
	{
		for i := 5; i < 16; i++ {
			currenthash := rlpHash(fmt.Sprint(i))
			if _, err := dbInstance.NewBlock(big.NewInt(int64(i)), parenthash, currenthash); err != nil {
				t.Fatal(err)
			}
			for j := 0; j < 20; j++ {
				if _, err := dbInstance.Put(currenthash, []byte(fmt.Sprint(j)), []byte(fmt.Sprint(j))); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := dbInstance.Commit(currenthash); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Run("a block kv>2000", func(t *testing.T) {
		ok, err := dbInstance.Compaction()
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("should be ok")
		}
		if dbInstance.current.BaseNum.Int64() != 1 {
			t.Error("must be 1", dbInstance.current.BaseNum)
		}
		if len(dbInstance.committed) != 14 {
			t.Error("must be 14:", len(dbInstance.committed))
		}
		for i := 0; i < 3000; i++ {
			v, err := dbInstance.baseDB.Get([]byte(fmt.Sprint(i)), nil)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(v, []byte(fmt.Sprint(i))) != 0 {
				t.Error("value not the same")
			}
		}
	})
	t.Run("kv<2000,block<10", func(t *testing.T) {
		ok, err := dbInstance.Compaction()
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("should be ok")
		}
		if dbInstance.current.BaseNum.Int64() != 4 {
			t.Error("must be 4", dbInstance.current.BaseNum)
		}
		if len(dbInstance.committed) != 11 {
			t.Error("must be 11:", len(dbInstance.committed))
		}
		for i := 3000; i < 4998; i++ {
			v, err := dbInstance.baseDB.Get([]byte(fmt.Sprint(i)), nil)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(v, []byte(fmt.Sprint(i))) != 0 {
				t.Error("value not the same")
			}
		}
	})
	t.Run("kv<2000,block=10", func(t *testing.T) {
		ok, err := dbInstance.Compaction()
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("should be ok")
		}
		if dbInstance.current.BaseNum.Int64() != 13 {
			t.Error("must be 14", dbInstance.current.BaseNum)
		}
		if len(dbInstance.committed) != 2 {
			t.Error("must be 1:", len(dbInstance.committed))
		}
	})
}

//1.put之前必须newblock，如果没有需要返回错误   -
func TestPutToUnRecognized(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		_, err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	data := [][2]string{
		[2]string{"a", "b"},
		[2]string{"b", "4421ffwef"},
		[2]string{"C", "2wgewbrw2"},
	}

	if _, err := db.Put(common.ZeroHash, []byte("a"), []byte("b")); err == nil {
		t.Error("new block must call before put to UnRecognized")
	}

	//	currentHash := rlpHash("b")
	if _, err := db.NewBlock(big.NewInt(20), parentHash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	var lastkvHash common.Hash
	var lastkvHashs []common.Hash

	for _, value := range data {
		if _, err := db.Put(common.ZeroHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(db.GetLastKVHash(common.ZeroHash), db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash).Bytes()) != 0 {
			t.Fatal("kv hash is wrong")
		}
		lastkvHash = db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash)
		lastkvHashs = append(lastkvHashs, lastkvHash)
	}
	for _, value := range data {
		v, err := db.unRecognized.data.Get([]byte(value[0]))
		if err != nil {
			t.Fatal(err)
		}
		if string(v) != value[1] {
			t.Fatal("should equal")
		}
	}

	fd := fileDesc{Type: TypeJournal, Num: db.unRecognized.Number.Int64(), BlockHash: db.getUnRecognizedHash()}
	read, err := db.storage.Open(fd)
	if err != nil {
		t.Fatal("should open storage", err)
	}
	defer read.Close()
	if err != nil {
		panic(err)
	}
	r := journal.NewReader(read, nil, true, true)
	rr, err := r.Next()
	if err != nil {
		t.Fatal("next", err)
	}
	var header journalHeader
	if err := decode(rr, &header); err != nil {
		t.Fatal(err)
	}
	if header.ParentHash.String() != parentHash.String() {
		t.Fatal("header ParentHash should same")
	}
	if header.BlockNumber.Int64() != 20 {
		t.Fatal("header BlockNumber should same")
	}
	var i int
	for _, value := range data {
		reader, err := r.Next()
		if err != nil {
			t.Fatal(err)
		}
		var body journalData
		if err := decode(reader, &body); err != nil {
			t.Fatal(err)
		}
		if body.FuncType != funcTypePut {
			t.Fatal("body FuncType should be put")
		}
		if string(body.Key) != value[0] {
			t.Fatal("body key should be same", string(body.Key), value[0])
		}
		if string(body.Value) != value[1] {
			t.Fatal("body value should be same", string(body.Value), value)
		}
		if lastkvHashs[i] != body.Hash {
			t.Fatal("kv hash is wrong")
		}
		i++
	}
}

//需要测试kv hash的正确性
func TestPutToRecognized(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		_, err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	parentHash := rlpHash("a")
	data := [][2]string{
		[2]string{"a", "b"},
		[2]string{"b", "4421ffwef"},
		[2]string{"C", "2wgewbrw2"},
	}
	currentHash := rlpHash("b")
	if _, err := db.NewBlock(big.NewInt(20), parentHash, currentHash); err != nil {
		t.Fatal(err)
	}
	var lastkvHash common.Hash
	var lastkvHashs []common.Hash
	for _, value := range data {
		if _, err := db.Put(currentHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(db.GetLastKVHash(currentHash), db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash).Bytes()) != 0 {
			t.Fatal("kv hash is wrong")
		}
		lastkvHash = db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash)
		lastkvHashs = append(lastkvHashs, lastkvHash)
	}
	recognized, ok := db.recognized[currentHash]
	if !ok {
		t.Fatal("[SnapshotDB] recognized hash should be find")
	}
	for _, value := range data {
		v, err := recognized.data.Get([]byte(value[0]))
		if err != nil {
			t.Fatal(err)
		}
		if string(v) != value[1] {
			t.Fatal("[SnapshotDB] should equal")
		}
	}

	fd := fileDesc{Type: TypeJournal, Num: recognized.Number.Int64(), BlockHash: recognized.BlockHash}
	read, err := db.storage.Open(fd)
	if err != nil {
		t.Fatal("[SnapshotDB]should open storage", err)
	}
	defer read.Close()
	if err != nil {
		panic(err)
	}
	r := journal.NewReader(read, nil, true, true)
	rr, err := r.Next()
	if err != nil {
		t.Fatal("next", err)
	}
	var header journalHeader
	if err := decode(rr, &header); err != nil {
		t.Fatal(err)
	}
	if header.ParentHash.String() != parentHash.String() {
		t.Fatal("header ParentHash should same")
	}
	if header.BlockNumber.Int64() != recognized.Number.Int64() {
		t.Fatal("header BlockNumber should same")
	}
	var i int

	for _, value := range data {
		reader, err := r.Next()
		if err != nil {
			t.Fatal(err)
		}
		var body journalData
		if err := decode(reader, &body); err != nil {
			t.Fatal(err)
		}
		if body.FuncType != funcTypePut {
			t.Fatal("body FuncType should be put")
		}
		if string(body.Key) != value[0] {
			t.Fatal("body key should be same", string(body.Key), value[0])
		}
		if string(body.Value) != value[1] {
			t.Fatal("body value should be same", string(body.Value), value)
		}
		if lastkvHashs[i] != body.Hash {
			t.Fatal("kv hash is wrong")
		}
		i++
	}
}

func TestFlush(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		_, err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	parentHash := rlpHash("a")
	blockNumber := big.NewInt(20)
	data := [][2]string{
		[2]string{"a", "b"},
		[2]string{"b", "4421ffwef"},
		[2]string{"C", "2wgewbrw2"},
	}
	if _, err := db.NewBlock(blockNumber, parentHash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if _, err := db.Put(common.ZeroHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}
	currentHash := rlpHash("b")
	if _, err := db.Flush(currentHash, blockNumber); err != nil {
		t.Fatal(err)
	}
	recognized, ok := db.recognized[currentHash]
	if !ok {
		t.Fatal("[SnapshotDB] recognized hash should be find")
	}
	if !recognized.readOnly {
		t.Fatal("[SnapshotDB] unrecognized flush to recognized , then the block must read only")
	}
	for _, value := range data {
		v, err := recognized.data.Get([]byte(value[0]))
		if err != nil {
			t.Fatal(err)
		}
		if string(v) != value[1] {
			t.Fatal("[SnapshotDB] should equal")
		}
	}

	fd := fileDesc{Type: TypeJournal, Num: recognized.Number.Int64(), BlockHash: recognized.BlockHash}
	read, err := db.storage.Open(fd)
	if err != nil {
		t.Fatal("[SnapshotDB]should open storage", err)
	}
	defer read.Close()
	if err != nil {
		panic(err)
	}
	r := journal.NewReader(read, nil, true, true)
	rr, err := r.Next()
	if err != nil {
		t.Fatal("next", err)
	}
	var header journalHeader
	if err := decode(rr, &header); err != nil {
		t.Fatal(err)
	}
	if header.ParentHash.String() != parentHash.String() {
		t.Fatal("header ParentHash should same")
	}
	if header.BlockNumber.Int64() != recognized.Number.Int64() {
		t.Fatal("header BlockNumber should same")
	}

	for _, value := range data {
		reader, err := r.Next()
		if err != nil {
			t.Fatal(err)
		}
		var body journalData
		if err := decode(reader, &body); err != nil {
			t.Fatal(err)
		}
		if body.FuncType != funcTypePut {
			t.Fatal("body FuncType should be put")
		}
		if string(body.Key) != value[0] {
			t.Fatal("body key should be same", string(body.Key), value[0])
		}
		if string(body.Value) != value[1] {
			t.Fatal("body value should be same", string(body.Value), value)
		}
	}
	if db.unRecognized != nil {
		t.Fatal("unRecognized must be nil")
	}

	if _, err := db.NewBlock(blockNumber.Add(blockNumber, big.NewInt(1)), parentHash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if _, err := db.Put(common.ZeroHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}

	//已经flush的block无法被写入
	if _, err := db.Put(currentHash, []byte("cccccccccc"), []byte("mmmmmmmmmmmm")); err == nil {
		t.Fatal("[SnapshotDB] can't update the block after flush")
	}
}

func TestCommit(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		_, err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	currentHash := rlpHash("currentHash")
	parentHash := rlpHash("parentHash")
	blockNumber := big.NewInt(1)
	data := [][2]string{
		[2]string{"a", "b"},
		[2]string{"b", "4421ffwef"},
		[2]string{"C", "2wgewbrw2"},
	}
	if _, err := db.NewBlock(blockNumber, parentHash, currentHash); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if _, err := db.Put(currentHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Commit(currentHash); err != nil {
		t.Fatal("commit fail:", err)
	}
	if db.current.HighestNum.Cmp(blockNumber) != 0 {
		t.Fatalf("current HighestNum must be :%v,but is%v", blockNumber.Int64(), db.current.HighestNum.Int64())
	}
	if db.committed[0].readOnly != true {
		t.Fatal("read only must be true")
	}
	if db.committed[0].BlockHash.String() != currentHash.String() {
		t.Fatal("BlockHash not cmp:", db.committed[0].BlockHash.String(), currentHash.String())
	}
	if db.committed[0].ParentHash.String() != parentHash.String() {
		t.Fatal("ParentHash not cmp", db.committed[0].ParentHash.String(), parentHash.String())
	}
	if db.committed[0].Number.Cmp(blockNumber) != 0 {
		t.Fatal("block number not cmp", db.committed[0].Number, blockNumber)
	}
	for _, value := range data {
		v, err := db.committed[0].data.Get([]byte(value[0]))
		if err != nil {
			t.Fatal(err)
		}
		if string(v) != value[1] {
			t.Fatal("[SnapshotDB] should equal")
		}
	}
	if _, ok := db.recognized[currentHash]; ok {
		t.Fatal("[SnapshotDB] should move to commit")
	}

}
