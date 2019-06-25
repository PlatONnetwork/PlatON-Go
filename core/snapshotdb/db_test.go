package snapshotdb

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/big"
	"os"
	"path"
	"testing"
)

func TestRecover(t *testing.T) {
	initDB()
	defer dbInstance.Clear()
	var (
		recognizedHash                                       = rlpHash("recognizedHash")
		parenthash                                           common.Hash
		baseDBArr, commitArr, recognizedArr, unrecognizedArr []string
		base, high                                           int64
		commit, recognized, unrecognized                     blockData
	)
	{
		commitHash := recognizedHash
		if _, err := dbInstance.NewBlock(big.NewInt(1), parenthash, &commitHash); err != nil {
			t.Fatal(err)
		}
		var str string
		for i := 0; i < 4; i++ {
			str += "a"
			baseDBArr = append(baseDBArr, str)
			if _, err := dbInstance.Put(&commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(&commitHash, []byte("d"), []byte("d")); err != nil {
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
		if _, err := dbInstance.NewBlock(big.NewInt(2), parenthash, &commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "c"
			commitArr = append(commitArr, str)
			if _, err := dbInstance.Put(&commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(&commitHash, []byte("ddd"), []byte("ddd")); err != nil {
			t.Fatal(err)
		}
		if _, err := dbInstance.Commit(commitHash); err != nil {
			t.Fatal(err)
		}
		commit = dbInstance.committed[0]
		parenthash = commitHash
	}
	{
		commitHash := rlpHash("recognizedHash4")
		if _, err := dbInstance.NewBlock(big.NewInt(3), parenthash, &commitHash); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "e"
			recognizedArr = append(recognizedArr, str)
			if _, err := dbInstance.Put(&commitHash, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(&commitHash, []byte("ee"), []byte("ee")); err != nil {
			t.Fatal(err)
		}
		recognized = dbInstance.recognized[rlpHash("recognizedHash4")]
		parenthash = commitHash
	}
	{
		if _, err := dbInstance.NewBlock(big.NewInt(4), parenthash, nil); err != nil {
			t.Fatal(err)
		}
		str := "a"
		for i := 0; i < 4; i++ {
			str += "f"
			unrecognizedArr = append(unrecognizedArr, str)
			if _, err := dbInstance.Put(nil, []byte(str), []byte(str)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := dbInstance.Put(nil, []byte("ff"), []byte("ff")); err != nil {
			t.Fatal(err)
		}
		unrecognized = *dbInstance.unRecognized
	}
	base = dbInstance.current.BaseNum.Int64()
	high = dbInstance.current.HighestNum.Int64()
	if _, err := dbInstance.Close(); err != nil {
		t.Error(err)
	}
	dbInstance = nil
	dbPath := path.Join(os.TempDir(), DBTestPath)
	s, err := openFile(dbPath, false)
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

	if dbInstance.path != dbPath {
		t.Error("path is wrong", dbInstance.path, dbPath)
	}
	if dbInstance.current.BaseNum.Int64() != base {
		t.Error("BaseNum is wrong", dbInstance.current.BaseNum.Int64(), base)
	}
	if dbInstance.current.HighestNum.Int64() != high {
		t.Error("HighestNum is wrong", dbInstance.current.HighestNum.Int64(), high)
	}
	if dbInstance.current.path != getCurrentPath(dbPath) {
		t.Error("current path is wrong", dbInstance.current.path, getCurrentPath(dbPath))
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
		dbInstance.recognized[rlpHash("recognizedHash4")],
		dbInstance.committed[0],
	}

	for i := 0; i < 3; i++ {
		if i == 0 {
			if newarr[i].BlockHash != nil {
				t.Error("unRecognized block hash must nil", i, newarr[i].BlockHash, oldarr[i].BlockHash)
			}
		} else {
			if *oldarr[i].BlockHash != *newarr[i].BlockHash {
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

	if _, err := dbInstance.Put(nil, []byte("dddd"), []byte("dddd")); err != nil {
		t.Error(err)
	}
	if _, err := dbInstance.Flush(rlpHash("flush"), big.NewInt(4)); err != nil {
		t.Error(err)
	}
	if _, err := dbInstance.Commit(rlpHash("recognizedHash4")); err != nil {
		t.Error(err)
	}
	if _, err := dbInstance.Compaction(); err != nil {
		t.Error(err)
	}
}

func TestRMOldRecognizedBlockData(t *testing.T) {

}

func TestCron(t *testing.T) {

}

func TestCheckHashChain(t *testing.T) {

}
