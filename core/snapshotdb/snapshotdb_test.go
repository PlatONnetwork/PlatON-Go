package snapshotdb

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"math/big"
	"testing"
)

const (
	TestDBPath = "./snapshotdb_test"
)

//1.put之前必须newblock，如果没有需要返回错误   -
func TestPutToUnRecognized(t *testing.T) {
	db, err := newDB(TestDBPath)
	if err != nil {
		t.Fatal("new test db fail:", err)
	}
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

	//	currentHash := rlpHash("b")
	if _, err := db.NewBlock(big.NewInt(20), parentHash, nil); err != nil {
		t.Fatal(err)
	}
	var lastkvHash common.Hash
	var lastkvHashs []common.Hash

	for _, value := range data {
		if _, err := db.Put(nil, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(db.GetLastKVHash(nil), db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash).Bytes()) != 0 {
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
	db, err := newDB(TestDBPath)
	if err != nil {
		t.Fatal("new test db fail:", err)
	}
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
	if _, err := db.NewBlock(big.NewInt(20), parentHash, &currentHash); err != nil {
		t.Fatal(err)
	}
	var lastkvHash common.Hash
	var lastkvHashs []common.Hash
	for _, value := range data {
		if _, err := db.Put(&currentHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(db.GetLastKVHash(&currentHash), db.generateKVHash([]byte(value[0]), []byte(value[1]), lastkvHash).Bytes()) != 0 {
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

	fd := fileDesc{Type: TypeJournal, Num: recognized.Number.Int64(), BlockHash: *recognized.BlockHash}
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
	db, err := newDB(TestDBPath)
	if err != nil {
		t.Fatal("new test db fail:", err)
	}
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
	if _, err := db.NewBlock(blockNumber, parentHash, nil); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if _, err := db.Put(nil, []byte(value[0]), []byte(value[1])); err != nil {
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

	fd := fileDesc{Type: TypeJournal, Num: recognized.Number.Int64(), BlockHash: *recognized.BlockHash}
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

	if _, err := db.NewBlock(blockNumber.Add(blockNumber, big.NewInt(1)), parentHash, nil); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if _, err := db.Put(nil, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}

	//已经flush的block无法被写入
	if _, err := db.Put(&currentHash, []byte("cccccccccc"), []byte("mmmmmmmmmmmm")); err == nil {
		t.Fatal("[SnapshotDB] can't update the block after flush")
	}
}

func TestCommit(t *testing.T) {
	db, err := newDB(TestDBPath)
	if err != nil {
		t.Fatal("new test db fail:", err)
	}
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
	if _, err := db.NewBlock(blockNumber, parentHash, &currentHash); err != nil {
		t.Fatal(err)
	}
	for _, value := range data {
		if _, err := db.Put(&currentHash, []byte(value[0]), []byte(value[1])); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Commit(currentHash); err != nil {
		t.Fatal("commit fail:", err)
	}
	if db.current.HighestNum.Cmp(blockNumber) != 0 {
		t.Fatalf("current HighestNum must be :%v,but is%v", blockNumber.Int64(), db.current.HighestNum.Int64())
	}
	if db.commited[0].readOnly != true {
		t.Fatal("read only must be true")
	}
	if db.commited[0].BlockHash.String() != currentHash.String() {
		t.Fatal("BlockHash not cmp:", db.commited[0].BlockHash.String(), currentHash.String())
	}
	if db.commited[0].ParentHash.String() != parentHash.String() {
		t.Fatal("ParentHash not cmp", db.commited[0].ParentHash.String(), parentHash.String())
	}
	if db.commited[0].Number.Cmp(blockNumber) != 0 {
		t.Fatal("block number not cmp", db.commited[0].Number, blockNumber)
	}
	for _, value := range data {
		v, err := db.commited[0].data.Get([]byte(value[0]))
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

func TestRMOldRecognizedBlockData(t *testing.T) {

}

func TestCron(t *testing.T) {

}
