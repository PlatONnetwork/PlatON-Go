package snapshotdb

import (
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/journal"
)

func TestJournal(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	blockHash := generateHash("a")
	parentHash := generateHash("b")
	blockNumber := big.NewInt(100)

	err := dbInstance.NewBlock(blockNumber, parentHash, blockHash)
	if err != nil {
		t.Error(err)
	}

	str := []string{"abcdefghijk", "kjhiughdjdi", "acadadwdfqwrwq"}

	if err := db.unCommit.blocks[blockHash].writeJournalBody([]byte(str[0])); err != nil {
		t.Error(err)
	}

	if err := db.unCommit.blocks[blockHash].writeJournalBody([]byte(str[1])); err != nil {
		t.Error(err)
	}
	if err := db.unCommit.blocks[blockHash].writeJournalBody([]byte(str[2])); err != nil {
		t.Error(err)
	}
	fd := fileDesc{Type: TypeJournal, Num: blockNumber.Uint64(), BlockHash: blockHash}
	file, err := db.storage.Open(fd)
	if err != nil {
		t.Error(err)
	}
	journals := journal.NewReader(file, nil, true, true)
	r, err := journals.Next()
	if err != nil {
		t.Error(err)
	}
	var jHead journalHeader
	if err := rlp.Decode(r, &jHead); err != nil {
		t.Error(err)
	}
	if jHead.From != journalHeaderFromRecognized {
		t.Error("from is wrong", jHead.From)
	}
	if jHead.ParentHash != parentHash {
		t.Error("parent hash is wrong", jHead.ParentHash)
	}
	if jHead.BlockNumber.Int64() != blockNumber.Int64() {
		t.Error("block number is wrong", jHead.BlockNumber)
	}
	i := 0
	for {
		jo, err := journals.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal("next", err)
		}
		x, err := ioutil.ReadAll(jo)
		b := string(x)
		if strings.Compare(str[i], b) != 0 {
			t.Errorf("should eq:%v,%v", str[i], b)
		}
		i++
	}
	if err := file.Close(); err != nil {
		t.Error(err)
	}
	if err := db.rmJournalFile(blockNumber, blockHash); err != nil {
		t.Error(err)
	}
	_, err = db.storage.Open(fd)
	if !os.IsNotExist(os.ErrNotExist) {
		t.Error(err)
	}

}

func TestCloseJournalWriter(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "test_close*.log")
	if err != nil {
		t.Error(err)
	}
	jw := newJournalWriter(f)
	writer, err := jw.journal.Next()
	if err != nil {
		t.Error(err)
	}
	if _, err := writer.Write([]byte("a")); err != nil {
		t.Error("should write", err)
	}
	if err := jw.Close(); err != nil {
		t.Error("should can close", err)
	}
	if _, err := jw.journal.Next(); err == nil {
		t.Fatal(err)
	}
	if err := jw.writer.Close(); err == nil {
		t.Error("should have be closed")
	}
}

//
//func TestRMthan(t *testing.T) {
//	Instance()
//	if err := dbInstance.writeJournalHeader(big.NewInt(1), generateHash("aaaa"), common.ZeroHash, journalHeaderFromUnRecognized); err != nil {
//		t.Error(err)
//	}
//	if err := dbInstance.writeJournalBody(generateHash("aaaa"), []byte("abcdefg")); err != nil {
//		t.Error(err)
//	}
//	fds, err := dbInstance.storage.List(TypeJournal)
//	if err != nil {
//		t.Error(err)
//	}
//
//	for _, fd := range fds {
//		reader, err := dbInstance.storage.Open(fd)
//		if err != nil {
//			t.Error(err)
//		}
//		journals := journal.NewReader(reader, nil, false, false)
//		j, err := journals.Next()
//		if err != nil {
//			t.Error(err)
//		}
//		var header journalHeader
//		if err := decode(j, &header); err != nil {
//			t.Error(err)
//		}
//		logger.Debug("header", "v", header)
//		for {
//			j, err := journals.Next()
//			if err == io.EOF {
//				break
//			}
//			if err != nil {
//				t.Error(err)
//			}
//			var body string
//			if err := decode(j, &body); err != nil {
//				t.Error(err)
//			}
//			logger.Debug("body", "v", body, "byte", j)
//		}
//	}
//
//	if err := dbInstance.writeJournalHeader(big.NewInt(1), generateHash("aaaa"), common.ZeroHash, journalHeaderFromRecognized); err != nil {
//		t.Error(err)
//	}
//	if err := dbInstance.writeJournalBody(generateHash("aaaa"), []byte("bbbbbbbbbbbbbbbbbbbb")); err != nil {
//		t.Error(err)
//	}
//
//	for _, fd := range fds {
//		reader, err := dbInstance.storage.Open(fd)
//		if err != nil {
//			t.Error(err)
//		}
//		journals := journal.NewReader(reader, nil, false, false)
//		j, err := journals.Next()
//		if err != nil {
//			t.Error(err)
//		}
//		var header journalHeader
//		if err := decode(j, &header); err != nil {
//			t.Error(err)
//		}
//		logger.Debug("header", "v", header)
//		for {
//			j, err := journals.Next()
//			if err == io.EOF {
//				break
//			}
//			if err != nil {
//				t.Error(err)
//			}
//			var body string
//			if err := decode(j, &body); err != nil {
//				t.Error(err)
//			}
//			logger.Debug("body", "v", body, "byte", j)
//		}
//	}
//}
