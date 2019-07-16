package snapshotdb

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
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
		b, ok := dbInstance.recognized.Load(c)
		if !ok {
			t.Fatal("must find recognized")
		}
		bd := b.(blockData)
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
		bd := dbInstance.unRecognized
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
	commitDBkv1 := generatekvWithPrefix(30, "ac")
	if err := newBlockCommited(big.NewInt(2), commitHash, commitHash1, commitDBkv1); err != nil {
		t.Error(err)
		return
	}

	f := func(hash common.Hash, arr []kvs, num int) error {
		itr := dbInstance.Ranking(hash, []byte("a"), 100)
		var i int
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
				i++
			}
		}
		if num != i {
			return fmt.Errorf("num not the same want %d,have %d", num, i)
		}
		itr.Release()
		return nil
	}
	t.Run("commitHash", func(t *testing.T) {
		if err := f(commitHash, []kvs{commitDBkv}, 30); err != nil {
			t.Error(err)
			return
		}
	})
	t.Run("commitHash2", func(t *testing.T) {
		if err := f(commitHash1, []kvs{commitDBkv, commitDBkv1}, 60); err != nil {
			t.Error(err)
			return
		}
	})
}

func TestSnapshotDB_Ranking(t *testing.T) {
	initDB()
	defer dbInstance.Clear()

	parenthash := generateHash("parenthash")
	baseDBBlockhash := generateHash("baseDBBlockhash")
	baseDBkv := generatekvWithPrefix(20, "aa")
	if err := newBlockBaseDB(big.NewInt(1), parenthash, baseDBBlockhash, baseDBkv); err != nil {
		t.Error(err)
		return
	}

	commitHash := generateHash("commitHash")
	commitDBkv := generatekvWithPrefix(30, "ab")
	if err := newBlockCommited(big.NewInt(2), baseDBBlockhash, commitHash, commitDBkv); err != nil {
		t.Error(err)
		return
	}

	recognizedHash := generateHash("recognizedHash")
	recognizedDBkv := generatekvWithPrefix(40, "ac")
	if err := newBlockRecognizedDirect(big.NewInt(3), commitHash, recognizedHash, recognizedDBkv); err != nil {
		t.Error(err)
		return
	}

	recognizedHash2 := generateHash("recognizedHash2")
	recognizedDBkv2 := generatekvWithPrefix(40, "ae")
	if err := newBlockRecognizedDirect(big.NewInt(3), commitHash, recognizedHash2, recognizedDBkv2); err != nil {
		t.Error(err)
		return
	}

	recognizedHash3 := generateHash("recognizedHash3")
	if err := newBlockRecognizedDirect(big.NewInt(3), commitHash, recognizedHash3, nil); err != nil {
		t.Error(err)
		return
	}
	if err := dbInstance.Del(recognizedHash3, commitDBkv[0].key); err != nil {
		t.Error(err)
		return
	}

	unrecognizedDBkv := generatekvWithPrefix(50, "ad")
	if err := newBlockUnRecognized(big.NewInt(3), recognizedHash, unrecognizedDBkv); err != nil {
		t.Error(err)
		return
	}

	f := func(hash common.Hash, arr []kvs) error {
		itr := dbInstance.Ranking(hash, []byte("a"), 100)
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

	t.Run("with hash", func(t *testing.T) {
		t.Run("from recognized", func(t *testing.T) {
			if err := f(recognizedHash, []kvs{baseDBkv, commitDBkv, recognizedDBkv}); err != nil {
				t.Error(err)
				return
			}
		})
		t.Run("from other recognized", func(t *testing.T) {
			if err := f(recognizedHash2, []kvs{baseDBkv, commitDBkv, recognizedDBkv2}); err != nil {
				t.Error(err)
				return
			}
		})
		t.Run("form commit", func(t *testing.T) {
			if err := f(commitHash, []kvs{baseDBkv, commitDBkv}); err != nil {
				t.Error(err)
				return
			}
		})
		t.Run("delete should not be see", func(t *testing.T) {
			if err := f(recognizedHash3, []kvs{baseDBkv, commitDBkv[1 : len(commitDBkv)-1]}); err != nil {
				t.Error(err)
				return
			}
		})

	})
	t.Run("with out hash", func(t *testing.T) {
		if err := f(common.ZeroHash, []kvs{baseDBkv, commitDBkv, recognizedDBkv, unrecognizedDBkv}); err != nil {
			t.Error(err)
			return
		}
	})
}

func TestSnapshotDB_WalkBaseDB(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		recognizedHash = generateHash("recognizedHash")
		parenthash     common.Hash
		arr            []string
	)
	{
		commitHash := recognizedHash
		if err := dbInstance.NewBlock(big.NewInt(1), parenthash, commitHash); err != nil {
			t.Fatal(err)
		}
		var str string
		for i := 0; i < 4; i++ {
			str += "a"
			arr = append(arr, str)
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
	t.Run("walk base db when compaction,should lock", func(t *testing.T) {
		dbInstance.snapshotLockC = true
		go func() {
			time.Sleep(time.Millisecond * 500)
			dbInstance.snapshotLock.Send(struct{}{})
			dbInstance.snapshotLockC = false
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

//1.put之前必须newblock，如果没有需要返回错误   -
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
	rg, ok := db.recognized.Load(currentHash)
	if !ok {
		t.Fatal("[SnapshotDB] recognized hash should be find")
	}
	recognized := rg.(blockData)
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
	rg, ok := db.recognized.Load(currentHash)
	if !ok {
		t.Fatal("[SnapshotDB] recognized hash should be find")
	}
	recognized := rg.(blockData)
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

	if err := db.NewBlock(blockNumber.Add(blockNumber, big.NewInt(1)), parentHash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if err := db.Put(common.ZeroHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}

	//已经flush的block无法被写入
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
	if _, ok := db.recognized.Load(currentHash); ok {
		t.Fatal("[SnapshotDB] should move to commit")
	}

}

func TestNewBlockOpenTooMany(t *testing.T) {
	//os.RemoveAll(dbpath)
	//initDB()
	//db := dbInstance
	//defer func() {
	//	err := db.Close()
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//}()
	//for i := 0; i < 110; i++ {
	//	if err := db.NewBlock(big.NewInt(int64(i+1)), generateHash(fmt.Sprint(i)), generateHash(fmt.Sprint(i+1))); err != nil {
	//		t.Fatal(err)
	//		break
	//	}
	//	if err := db.Put(generateHash(fmt.Sprint(i+1)), []byte(fmt.Sprint(i)), []byte(fmt.Sprint(i))); err != nil {
	//		t.Fatal(err)
	//		break
	//	}
	//}
	//for i := 0; i < 110; i++ {
	//	if err := db.Commit(generateHash(fmt.Sprint(i + 1))); err != nil {
	//		t.Fatal(err)
	//		break
	//	}
	//}
	//for len(dbInstance.committed) > 90 {
	//
	//}
	//b, err := exec.Command("bash", "-c", fmt.Sprintf("lsof -p %d | wc -l", os.Getpid())).Output()
	//if err != nil {
	//	t.Error(err)
	//}
	//log.Print("open file num:", string(b))
}
