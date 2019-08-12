package snapshotdb

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
//parentHash  = generateHash("parentHash")
//currentHash = generateHash("currentHash")
)

func TestCommitZeroBlock(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	if err := newBlockCommited(big.NewInt(0), common.ZeroHash, generateHash("zerohash"), generatekv(1)); err != nil {
		t.Error(err)
	}
}

func TestSnapshotDB_NewBlock(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		p = generateHash("parentHash")
		c = generateHash("currentHash")
	)
	t.Run("new recognized block", func(t *testing.T) {
		err := dbInstance.NewBlock(big.NewInt(30), p, c)
		if err != nil {
			t.Error(err)
		}
		bd, ok := dbInstance.unCommit.blocks[c]
		if !ok {
			t.Fatal("must find recognized")
		}
		if bd.ParentHash != p {
			t.Fatal("parentHash must same:", bd.ParentHash, p)
		}
		if bd.Number.Cmp(big.NewInt(30)) != 0 {
			t.Fatal("block number must same:", bd.Number, big.NewInt(30))
		}
		if bd.BlockHash != c {
			t.Fatal("BlockHash must right:", bd.BlockHash, c)
		}
	})
	t.Run("new unrecognized block", func(t *testing.T) {
		err := dbInstance.NewBlock(big.NewInt(30), p, common.ZeroHash)
		if err != nil {
			t.Error(err)
		}
		bd := dbInstance.unCommit.blocks[dbInstance.getUnRecognizedHash()]
		if bd.ParentHash != p {
			t.Fatal("parentHash must same:", bd.ParentHash, p)
		}
		if bd.Number.Cmp(big.NewInt(30)) != 0 {
			t.Fatal("block number must same:", bd.Number, big.NewInt(30))
		}
	})
}

func TestSnapshotDB_GetWithNoCommit(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr = [][]kv{generatekv(10), generatekv(10)}
	)
	//recognized(unRecognized not in the chain)
	if err := newBlockRecognizedDirect(big.NewInt(2), generateHash(fmt.Sprint(1)), generateHash(fmt.Sprint(2)), arr[0]); err != nil {
		t.Error(err)
	}
	//unRecognized
	if err := newBlockUnRecognized(big.NewInt(3), generateHash(fmt.Sprint(2)), arr[1]); err != nil {
		t.Error(err)
	}

	for _, a := range arr {
		for _, kv := range a {
			val, err := dbInstance.Get(common.ZeroHash, kv.key)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(kv.value, val) != 0 {
				t.Error("must find key")
			}
		}
	}
}

func TestSnapshotDB_Get_after_del(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr = [][]kv{generatekv(10), generatekv(10), generatekv(10), generatekv(10), generatekv(10)}
	)
	{
		//baseDB
		if err := newBlockBaseDB(big.NewInt(1), generateHash(fmt.Sprint(0)), generateHash(fmt.Sprint(1)), arr[0]); err != nil {
			t.Error(err)
			return
		}
		//commit
		if err := newBlockCommited(big.NewInt(2), generateHash(fmt.Sprint(1)), generateHash(fmt.Sprint(2)), arr[1]); err != nil {
			t.Error(err)
			return
		}

		//recognized
		if err := newBlockRecognizedDirect(big.NewInt(3), generateHash(fmt.Sprint(2)), generateHash(fmt.Sprint(3)), arr[2]); err != nil {
			t.Error(err)
			return
		}

		//unRecognized
		if err := newBlockUnRecognized(big.NewInt(4), generateHash(fmt.Sprint(3)), arr[3]); err != nil {
			t.Error(err)
		}

		t.Run("delete commit", func(t *testing.T) {
			key := arr[1][0].key
			if err := dbInstance.Del(generateHash(fmt.Sprint(3)), key); err != nil {
				t.Error(err)
				return
			}
			_, err := dbInstance.Get(generateHash(fmt.Sprint(3)), key)
			if err != ErrNotFound {
				t.Error(err)
				return
			}
			if err := dbInstance.Commit(generateHash(fmt.Sprint(3))); err != nil {
				t.Error(err)
				return
			}
			_, err = dbInstance.Get(common.ZeroHash, key)
			if err != ErrNotFound {
				t.Error(err)
				return
			}
		})
	}
}

func TestSnapshotDB_Get(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr = [][]kv{generatekv(10), generatekv(10), generatekv(10), generatekv(10), generatekv(10)}
	)
	{
		//baseDB
		if err := newBlockBaseDB(big.NewInt(1), generateHash(fmt.Sprint(0)), generateHash(fmt.Sprint(1)), arr[0]); err != nil {
			t.Error(err)
		}
		//commit
		if err := newBlockCommited(big.NewInt(2), generateHash(fmt.Sprint(1)), generateHash(fmt.Sprint(2)), arr[1]); err != nil {
			t.Error(err)
		}

		//recognized
		if err := newBlockRecognizedDirect(big.NewInt(3), generateHash(fmt.Sprint(2)), generateHash(fmt.Sprint(3)), arr[2]); err != nil {
			t.Error(err)
		}

		//recognized(unRecognized not in the chain)
		if err := newBlockRecognizedDirect(big.NewInt(3), generateHash(fmt.Sprint(2)), generateHash(fmt.Sprint(5)), arr[4]); err != nil {
			t.Error(err)
		}

		//unRecognized
		if err := newBlockUnRecognized(big.NewInt(4), generateHash(fmt.Sprint(3)), arr[3]); err != nil {
			t.Error(err)
		}

	}

	t.Run("UnRecognized", func(t *testing.T) {
		t.Run("must find all in the chain", func(t *testing.T) {
			for _, a := range arr[0:3] {
				for _, kv := range a {
					val, err := dbInstance.Get(common.ZeroHash, kv.key)
					if err != nil {
						t.Error(err)
					}
					if bytes.Compare(kv.value, val) != 0 {
						t.Error("must find key")
					}
				}
			}
		})
		t.Run("must not find key not exist", func(t *testing.T) {
			_, err := dbInstance.Get(common.ZeroHash, []byte("e"))
			if err == nil {
				t.Error(err)
			}
		})
		t.Run("not in the chain,must not find ", func(t *testing.T) {
			for _, kv := range arr[4] {
				_, err := dbInstance.Get(common.ZeroHash, kv.key)

				if err == nil {
					t.Error(err)
				}
				if err != ErrNotFound {
					t.Error("must not find")
				}
			}
		})
	})

	t.Run("Recognized", func(t *testing.T) {
		t.Run("must find all in the chain", func(t *testing.T) {
			for _, a := range arr[0:2] {
				for _, kv := range a {
					val, err := dbInstance.Get(generateHash(fmt.Sprint(3)), kv.key)
					if err != nil {
						t.Error(err)
					}
					if bytes.Compare(kv.value, val) != 0 {
						t.Error("must find key")
					}
				}
			}
		})
		t.Run("must not find", func(t *testing.T) {
			for _, a := range arr[3:4] {
				for _, kv := range a {
					_, err := dbInstance.Get(generateHash(fmt.Sprint(3)), kv.key)
					if err == nil {
						t.Error(err)
					}
					if err != ErrNotFound {
						t.Error("must not find")
					}
				}
			}
		})
	})

	t.Run("committed", func(t *testing.T) {
		t.Run("must find all in the chain", func(t *testing.T) {
			for _, a := range arr[0:1] {
				for _, kv := range a {
					val, err := dbInstance.Get(generateHash(fmt.Sprint(2)), kv.key)
					if err != nil {
						t.Error(err)
					}
					if bytes.Compare(kv.value, val) != 0 {
						t.Error("must find key")
					}
				}
			}
		})
	})

	t.Run("baseDB", func(t *testing.T) {
		t.Run("must find all in the chain", func(t *testing.T) {
			for _, kv := range arr[0] {
				val, err := dbInstance.Get(generateHash(fmt.Sprint(1)), kv.key)
				if err != nil {
					t.Error(err)
				}
				if bytes.Compare(kv.value, val) != 0 {
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
		baseDBHash  = generateHash("hash1")
		commit1hash = generateHash("hash2")
		commit2hash = generateHash("hash3")
		baseDBkv    = kv{key: []byte("a"), value: []byte("a")}
		commit1KV   = kv{key: []byte("b"), value: []byte("b")}
		commit2KV   = kv{key: []byte("b"), value: []byte("c")}
	)
	if err := newBlockBaseDB(big.NewInt(1), common.ZeroHash, baseDBHash, []kv{baseDBkv}); err != nil {
		t.Error(err)
		return
	}
	if err := newBlockCommited(big.NewInt(2), baseDBHash, commit1hash, []kv{commit1KV}); err != nil {
		t.Error(err)
		return
	}

	t.Run("should get", func(t *testing.T) {
		for _, key := range []kv{commit1KV, baseDBkv} {
			val, err := dbInstance.GetFromCommittedBlock(key.key)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(key.value, val) != 0 {
				t.Error("val not compare", key.value, val)
			}
		}
	})
	t.Run("should change", func(t *testing.T) {
		if err := newBlockCommited(big.NewInt(3), commit1hash, commit2hash, []kv{commit2KV}); err != nil {
			t.Error(err)
			return
		}
		val, err := dbInstance.GetFromCommittedBlock(commit2KV.key)
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(commit2KV.value, val) != 0 {
			t.Error("must find key")
		}
	})
	t.Run("not find from any path", func(t *testing.T) {
		_, err := dbInstance.GetFromCommittedBlock([]byte("ccccc"))
		if err != ErrNotFound {
			t.Error(err)
		}
	})
}

func TestSnapshotDB_Del(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr                   = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")}
		recognizedHash        = generateHash("recognizedHash")
		recognizedByFlushHash = generateHash("recognizedByFlush")
		commitHash            = generateHash("commitHash")
		baseHash              = generateHash("baseHash")
	)
	if err := newBlockBaseDB(big.NewInt(1), common.ZeroHash, baseHash, []kv{kv{key: arr[3], value: arr[3]}}); err != nil {
		t.Error(err)
		return
	}
	if err := newBlockCommited(big.NewInt(2), baseHash, commitHash, []kv{kv{key: arr[2], value: arr[2]}}); err != nil {
		t.Error(err)
		return
	}
	if err := newBlockRecognizedDirect(big.NewInt(3), commitHash, recognizedHash, []kv{kv{key: arr[1], value: arr[1]}}); err != nil {
		t.Error(err)
		return
	}
	if err := newBlockRecognizedByFlush(big.NewInt(3), commitHash, recognizedByFlushHash, []kv{kv{key: arr[4], value: arr[4]}}); err != nil {
		t.Error(err)
		return
	}
	if err := newBlockUnRecognized(big.NewInt(4), recognizedByFlushHash, []kv{kv{key: arr[0], value: arr[0]}}); err != nil {
		t.Error(err)
		return
	}
	t.Run("delete unrecognized", func(t *testing.T) {
		err := dbInstance.Del(common.ZeroHash, arr[0])
		if err != nil {
			t.Error("err must be nil", err)
		}
	})
	t.Run("delete recognized", func(t *testing.T) {
		err := dbInstance.Del(recognizedHash, arr[1])
		if err != nil {
			t.Error("err must be nil", err)
		}
	})

	t.Run("can't delete readonly", func(t *testing.T) {
		err := dbInstance.Del(recognizedByFlushHash, arr[4])
		if err == nil {
			t.Error("err must not nil", err)
		}
	})
	t.Run("can't delete commit", func(t *testing.T) {
		err := dbInstance.Del(commitHash, arr[2])
		if err == nil {
			t.Error("err must not nil", err)
		}
	})
}

func TestSnapshotDB_Has(t *testing.T) {
	//same as get
}

func TestSnapshotDB_Ranking2(t *testing.T) {
	initDB()
	defer dbInstance.Clear()

	commitHash := generateHash("commitHash1")
	commitDBkv := generatekvWithPrefix(30, "ab")
	if err := newBlockCommited(big.NewInt(1), common.ZeroHash, commitHash, commitDBkv); err != nil {
		t.Error(err)
		return
	}

	commitHash1 := generateHash("commitHash2")
	commitDBkv1 := generatekvWithPrefix(40, "ac")
	if err := newBlockCommited(big.NewInt(2), commitHash, commitHash1, commitDBkv1); err != nil {
		t.Error(err)
		return
	}
	t.Run("ranking kv should compare", func(t *testing.T) {
		f := func(hash common.Hash, prefix string, arr []kvs, num int) error {
			itr := dbInstance.Ranking(hash, []byte(prefix), num)
			for _, kvs := range arr {
				for _, kv := range kvs {
					if !itr.Next() {
						return errors.New("it's must can itr")
					}
					if bytes.Compare(kv.value, itr.Value()) != 0 {
						return fmt.Errorf("itr return wrong value :%v,should return:%v ", itr.Key(), kv.value)
					}
					if bytes.Compare(kv.key, itr.Key()) != 0 {
						return fmt.Errorf("itr return wrong key :%v,should return:%v ", itr.Key(), kv.key)
					}
				}
			}
			itr.Release()
			return nil
		}
		if err := f(commitHash, "ab", []kvs{commitDBkv}, len(commitDBkv)); err != nil {
			t.Error(err)
			return
		}
		if err := f(commitHash1, "ac", []kvs{commitDBkv1}, len(commitDBkv1)); err != nil {
			t.Error(err)
			return
		}
	})
	t.Run("ranking num should compare", func(t *testing.T) {
		type testData struct {
			input  int
			expect int
			des    string
		}
		datas := []testData{
			{10, 10, "num less than total"},
			{40, 30, "num large than total"},
			{30, 30, "num eq total"},
			{0, 30, "if input num is 0,should return total"},
		}
		f := func(prefix string, num int) int {
			itr := dbInstance.Ranking(commitHash, []byte(prefix), num)
			var i int
			for itr.Next() {
				i++
			}
			itr.Release()
			return i
		}
		for _, data := range datas {
			i := f("ab", data.input)
			if i != data.expect {
				t.Errorf("%s,ranking num must compare,want %d,have %d", data.des, data.expect, i)
			}
		}
	})
}

func TestSnapshotDB_Ranking4(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	generatekvs := generatekvWithPrefix(1000, "aaa")
	if err := newBlockBaseDB(big.NewInt(1), common.ZeroHash, generateHash("baseDBBlockhash"), generatekvs); err != nil {
		t.Error(err)
	}
	f := func(hash common.Hash, prefix string, num int) kvs {
		itr := dbInstance.Ranking(hash, []byte(prefix), num)
		err := itr.Error()
		if err != nil {
			t.Fatal(err)
		}
		o := make(kvs, 0)
		for itr.Next() {
			o = append(o, kv{itr.Key(), itr.Value()})
		}
		itr.Release()
		return o
	}
	v := f(common.ZeroHash, "aaa", 1000)
	if err := v.compareWithkvs(generatekvs); err != nil {
		t.Error(err)
	}
	v2 := f(common.ZeroHash, "aaa", 0)
	if err := v2.compareWithkvs(generatekvs); err != nil {
		t.Error(err)
	}
}

func TestSnapshotDB_Ranking5(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	generatekvs := generatekvWithPrefix(4, "aaa")
	if err := newBlockBaseDB(big.NewInt(1), common.ZeroHash, generateHash("baseDBBlockhash"), generatekvs); err != nil {
		t.Error(err)
	}
	if err := dbInstance.NewBlock(big.NewInt(2), generateHash("baseDBBlockhash"), generateHash("baseDBBlockhash2")); err != nil {
		t.Error(err)
	}
	if err := dbInstance.Del(generateHash("baseDBBlockhash2"), generatekvs[0].key); err != nil {
		t.Error(err)

	}
	if err := dbInstance.Commit(generateHash("baseDBBlockhash2")); err != nil {
		t.Error(err)

	}
	var parentHash common.Hash
	parentHash = generateHash("baseDBBlockhash2")
	for i := 3; i < 11; i++ {
		generatekvs := generatekvWithPrefix(1, "bbb")
		hash := generateHash(fmt.Sprintf("baseDBBlockhash%v", i))
		newBlockCommited(big.NewInt(int64(i)), parentHash, hash, generatekvs)
		parentHash = hash
	}

	itr := dbInstance.Ranking(generateHash("baseDBBlockhash10"), []byte("aaa"), 20)
	err := itr.Error()
	if err != nil {
		t.Fatal(err)
	}
	o := make(kvs, 0)
	for itr.Next() {
		o = append(o, kv{itr.Key(), itr.Value()})
	}
	itr.Release()
	if len(o) != 3 {
		t.Errorf("must equql 3,have %v", len(o))
	}
	for i := 0; i < 3; i++ {
		if !bytes.Equal(o[i].key, generatekvs[i+1].key) {
			t.Errorf("not compare want %v,have %v", generatekvs[i+1].key, o[i].key)
		}
		if !bytes.Equal(o[i].value, generatekvs[i+1].value) {
			t.Errorf("not compare want %v,have %v", generatekvs[i+1].value, o[i].value)
		}
	}

}

func TestSnapshotDB_Ranking3(t *testing.T) {
	initDB()
	defer dbInstance.Clear()

	generatekvs := generatekvWithPrefix(1000, "aaa")

	parenthash := generateHash("parenthash")
	baseDBBlockhash := generateHash("baseDBBlockhash")
	if err := newBlockBaseDB(big.NewInt(1), parenthash, baseDBBlockhash, generatekvs[0:150]); err != nil {
		t.Error(err)
		return
	}

	commitHash := generateHash("commitHash")
	if err := newBlockCommited(big.NewInt(2), baseDBBlockhash, commitHash, generatekvs[150:200]); err != nil {
		t.Error(err)
		return
	}

	recognizedHash := generateHash("recognizedHash")
	if err := newBlockRecognizedDirect(big.NewInt(3), commitHash, recognizedHash, generatekvs[200:400]); err != nil {
		t.Error(err)
		return
	}

	recognizedHash2 := generateHash("recognizedHash2")
	if err := newBlockRecognizedDirect(big.NewInt(4), recognizedHash, recognizedHash2, generatekvs[400:700]); err != nil {
		t.Error(err)
		return
	}

	recognizedHash3 := generateHash("recognizedHash3")
	if err := newBlockRecognizedDirect(big.NewInt(5), recognizedHash2, recognizedHash3, nil); err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 50; i++ {
		if err := dbInstance.Del(recognizedHash3, generatekvs[i].key); err != nil {
			t.Error(err)
			return
		}
	}

	if err := newBlockUnRecognized(big.NewInt(6), recognizedHash3, generatekvs[700:]); err != nil {
		t.Error(err)
		return
	}

	f := func(hash common.Hash, prefix string, num int) kvs {
		itr := dbInstance.Ranking(hash, []byte(prefix), num)
		o := make(kvs, 0)
		for itr.Next() {
			o = append(o, kv{itr.Key(), itr.Value()})
		}
		itr.Release()
		return o
	}
	v := f(common.ZeroHash, "aaa", 1000)
	if err := v.compareWithkvs(generatekvs[50:]); err != nil {
		t.Error(err)
	}
}

func TestSnapshotDB_WalkBaseDB(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		baseHash1 = generateHash("baseHash1")
		baseHash2 = generateHash("baseHash2")
		prefix    = util.BytesPrefix([]byte("a"))
	)
	kvsWithA := generatekvWithPrefix(100, "a")
	if err := newBlockBaseDB(big.NewInt(1), common.ZeroHash, baseHash1, kvsWithA); err != nil {
		t.Error(err)
	}
	if err := newBlockBaseDB(big.NewInt(2), baseHash1, baseHash2, generatekvWithPrefix(100, "b")); err != nil {
		t.Error(err)
	}
	t.Run("kv should compare", func(t *testing.T) {
		var kvGetFromWalk kvs
		f := func(num *big.Int, iter iterator.Iterator) error {
			if num.Int64() != 2 {
				return fmt.Errorf("basenum is wrong:%v,should be 2", num)
			}
			for iter.Next() {
				k, v := make([]byte, len(iter.Key())), make([]byte, len(iter.Value()))
				copy(k, iter.Key())
				copy(v, iter.Value())
				kvGetFromWalk = append(kvGetFromWalk, kv{
					key:   k,
					value: v,
				})
			}
			return nil
		}
		if err := dbInstance.WalkBaseDB(prefix, f); err != nil {
			t.Error(err)
		}
		sort.Sort(kvGetFromWalk)
		if err := kvsWithA.compareWithkvs(kvGetFromWalk); err != nil {
			t.Error(err)
		}
	})

	t.Run("walk base db when compaction,should lock", func(t *testing.T) {
		dbInstance.snapshotLockC = snapshotLock
		go func() {
			time.Sleep(time.Millisecond * 500)
			dbInstance.snapshotLock.Send(struct{}{})
			dbInstance.snapshotLockC = snapshotUnLock
		}()
		f2 := func(num *big.Int, iter iterator.Iterator) error {
			return nil
		}
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				if err := dbInstance.WalkBaseDB(prefix, f2); err != nil {
					t.Error(err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func TestSnapshotDB_Clear(t *testing.T) {

}

func TestSnapshotDB_GetLastKVHash(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		arr            = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")}
		recognizedHash = generateHash("recognizedHash")
		commitHash     = generateHash("commitHash")
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

func TestSnapshotDB_Compaction_del(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	baseDBHash := generateHash("base")
	baseDBkv := generatekv(10)
	if err := newBlockBaseDB(big.NewInt(1), common.ZeroHash, baseDBHash, baseDBkv); err != nil {
		t.Error(err)
		return
	}
	delkey := baseDBkv[0].key
	delVal := baseDBkv[0].value
	v, err := dbInstance.GetBaseDB(delkey)
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(v, delVal) != 0 {
		t.Error("must same")
		return
	}

	baseDBHash2 := generateHash("base2")
	if err := dbInstance.NewBlock(big.NewInt(2), baseDBHash, baseDBHash2); err != nil {
		t.Error(err)
		return
	}
	if err := dbInstance.Del(baseDBHash2, delkey); err != nil {
		t.Error(err)
		return
	}
	if err := dbInstance.Commit(baseDBHash2); err != nil {
		t.Error(err)
		return
	}
	if err := dbInstance.Compaction(); err != nil {
		t.Error(err)
		return
	}

	_, err = dbInstance.GetBaseDB(delkey)
	if err != ErrNotFound {
		t.Error(err)
		return
	}
}

func TestSnapshotDB_Compaction(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		commitHash1 = generateHash(fmt.Sprint(1))
		commitHash2 = generateHash(fmt.Sprint(2))
		commitHash3 = generateHash(fmt.Sprint(3))
		commitHash4 = generateHash(fmt.Sprint(4))
		kvs1        = generatekv(3000)
		kvs2        = generatekv(100)
		kvs3        = generatekv(100)
		kvs4        = generatekv(1798)
	)
	t.Run("0 commit block with Compaction", func(t *testing.T) {
		err := dbInstance.Compaction()
		if err != nil {
			t.Error(err)
		}
	})
	{
		if err := newBlockCommited(big.NewInt(1), common.ZeroHash, commitHash1, kvs1); err != nil {
			t.Error(err)
			return
		}
	}
	{
		if err := newBlockCommited(big.NewInt(2), commitHash1, commitHash2, kvs2); err != nil {
			t.Error(err)
			return
		}
		if err := newBlockCommited(big.NewInt(3), commitHash2, commitHash3, kvs3); err != nil {
			t.Error(err)
			return
		}

		if err := newBlockCommited(big.NewInt(4), commitHash3, commitHash4, kvs4); err != nil {
			t.Error(err)
			return
		}
	}
	{

		for i := 5; i < 16; i++ {
			if err := newBlockCommited(big.NewInt(int64(i)), generateHash(fmt.Sprint(i-1)), generateHash(fmt.Sprint(i)), generatekv(20)); err != nil {
				t.Error(err)
				return
			}
		}
	}
	t.Run("a block kv>2000,commit 1", func(t *testing.T) {
		err := dbInstance.Compaction()
		if err != nil {
			t.Error(err)
		}
		if dbInstance.current.BaseNum.Int64() != 1 {
			t.Error("must be 1", dbInstance.current.BaseNum)
		}
		if len(dbInstance.committed) != 14 {
			t.Error("must be 14:", len(dbInstance.committed))
		}
		for _, kv := range kvs1 {
			v, err := dbInstance.baseDB.Get(kv.key, nil)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(v, kv.value) != 0 {
				t.Error("value not the same")
			}
		}
	})
	t.Run("kv<2000,block<10,commit 2,3,4", func(t *testing.T) {
		err := dbInstance.Compaction()
		if err != nil {
			t.Error(err)
		}
		if dbInstance.current.BaseNum.Int64() != 4 {
			t.Error("must be 4", dbInstance.current.BaseNum)
		}
		if len(dbInstance.committed) != 11 {
			t.Error("must be 11:", len(dbInstance.committed))
		}
		for _, kvs := range [][]kv{kvs2, kvs3, kvs4} {
			for _, kv := range kvs {
				v, err := dbInstance.baseDB.Get(kv.key, nil)
				if err != nil {
					t.Error(err)
				}
				if bytes.Compare(v, kv.value) != 0 {
					t.Error("value not the same")
				}
			}
		}
	})
	t.Run("kv<2000,block=10,commit 5-15", func(t *testing.T) {
		err := dbInstance.Compaction()
		if err != nil {
			t.Error(err)
		}
		if dbInstance.current.BaseNum.Int64() != 14 {
			t.Error("must be 14", dbInstance.current.BaseNum)
		}
		if len(dbInstance.committed) != 1 {
			t.Error("must be 1:", len(dbInstance.committed))
		}
	})
}

//  put  must before newblock
func TestPutToUnRecognized(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	data := [][2]string{
		[2]string{"a", "b"},
		[2]string{"b", "4421ffwef"},
		[2]string{"C", "2wgewbrw2"},
	}
	parentHash := generateHash("parentHash")
	if err := db.Put(common.ZeroHash, []byte("a"), []byte("b")); err == nil {
		t.Error("new block must call before put to UnRecognized")
	}

	//	currentHash := rlpHash("b")
	if err := db.NewBlock(big.NewInt(20), parentHash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	var lastkvHash common.Hash
	var lastkvHashs []common.Hash

	for _, value := range data {
		if err := db.Put(common.ZeroHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(db.GetLastKVHash(common.ZeroHash), db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash).Bytes()) != 0 {
			t.Fatal("kv hash is wrong")
		}
		lastkvHash = db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash)
		lastkvHashs = append(lastkvHashs, lastkvHash)
	}
	block := db.unCommit.blocks[db.getUnRecognizedHash()]
	for _, value := range data {
		v, err := block.data.Get([]byte(value[0]))
		if err != nil {
			t.Fatal(err)
		}
		if string(v) != value[1] {
			t.Fatal("should equal")
		}
	}

	fd := fileDesc{Type: TypeJournal, Num: block.Number.Uint64(), BlockHash: db.getUnRecognizedHash()}
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

func TestPutToRecognized(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	parentHash := generateHash("a")
	data := [][2]string{
		[2]string{"a", "b"},
		[2]string{"b", "4421ffwef"},
		[2]string{"C", "2wgewbrw2"},
	}
	currentHash := generateHash("b")
	if err := db.NewBlock(big.NewInt(20), parentHash, currentHash); err != nil {
		t.Fatal(err)
	}
	var lastkvHash common.Hash
	var lastkvHashs []common.Hash
	for _, value := range data {
		if err := db.Put(currentHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(db.GetLastKVHash(currentHash), db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash).Bytes()) != 0 {
			t.Fatal("kv hash is wrong")
		}
		lastkvHash = db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash)
		lastkvHashs = append(lastkvHashs, lastkvHash)
	}
	recognized, ok := db.unCommit.blocks[currentHash]
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

	fd := fileDesc{Type: TypeJournal, Num: recognized.Number.Uint64(), BlockHash: recognized.BlockHash}
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
		err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	parentHash := generateHash("a")
	blockNumber := big.NewInt(20)
	data := [][2]string{
		[2]string{"a", "b"},
		[2]string{"b", "4421ffwef"},
		[2]string{"C", "2wgewbrw2"},
	}
	if err := db.NewBlock(blockNumber, parentHash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if err := db.Put(common.ZeroHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}
	currentHash := generateHash("b")
	if err := db.Flush(currentHash, blockNumber); err != nil {
		t.Fatal(err)
	}
	recognized, ok := db.unCommit.blocks[currentHash]
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

	fd := fileDesc{Type: TypeJournal, Num: recognized.Number.Uint64(), BlockHash: recognized.BlockHash}
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
		if string(body.Key) != value[0] {
			t.Fatal("body key should be same", string(body.Key), value[0])
		}
		if string(body.Value) != value[1] {
			t.Fatal("body value should be same", string(body.Value), value)
		}
	}

	if _, ok := db.unCommit.blocks[db.getUnRecognizedHash()]; ok {
		t.Fatal("unRecognized must be nil")
	}

	if err := db.NewBlock(blockNumber.Add(blockNumber, big.NewInt(1)), parentHash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if err := db.Put(common.ZeroHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}

	// can't write to flushed block
	if err := db.Put(currentHash, []byte("cccccccccc"), []byte("mmmmmmmmmmmm")); err == nil {
		t.Fatal("[SnapshotDB] can't update the block after flush")
	}
}

func TestCommit(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	currentHash := generateHash("currentHash")
	parentHash := generateHash("parentHash")
	blockNumber := big.NewInt(1)
	data := [][2]string{
		[2]string{"a", "b"},
		[2]string{"b", "4421ffwef"},
		[2]string{"C", "2wgewbrw2"},
	}
	if err := db.NewBlock(blockNumber, parentHash, currentHash); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if err := db.Put(currentHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Commit(currentHash); err != nil {
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
	if _, ok := db.unCommit.blocks[currentHash]; ok {
		t.Fatal("[SnapshotDB] should move to commit")
	}

}
