package snapshotdb

import (
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"
)

func TestJournal(t *testing.T) {
	initDB()
	db := dbInstance
	defer func() {
		_, err := db.Clear()
		if err != nil {
			t.Fatal(err)
		}
	}()
	blockHash := rlpHash("a")
	parentHash := rlpHash("b")
	blockNumber := big.NewInt(100)
	if err := db.writeJournalHeader(blockNumber, blockHash, parentHash, journalHeaderFromRecognized); err != nil {
		t.Error(err)
	}

	str := []string{"abcdefghijk", "kjhiughdjdi"}

	if err := db.writeJournalBody(blockHash, []byte(str[0])); err != nil {
		t.Error(err)
	}

	if err := db.writeJournalBody(blockHash, []byte(str[1])); err != nil {
		t.Error(err)
	}

	fd := fileDesc{Type: TypeJournal, Num: blockNumber.Int64(), BlockHash: blockHash}
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
	if err := db.rmJournalFile(blockNumber, blockHash); err != nil {
		t.Error(err)
	}
	_, err = db.storage.Open(fd)
	if !os.IsNotExist(os.ErrNotExist) {
		t.Error(err)
	}

}
