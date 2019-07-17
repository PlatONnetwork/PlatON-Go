package wal

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	viewChangeNumber = uint64(100)
	viewChangeHash   = common.BytesToHash(cbft.Rand32Bytes(32))
	times            = 1000
	tempDir          string
	wal              Wal
)

func TestWal(t *testing.T) {
	tempDir, _ = ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	// Test Wal UpdateViewChange
	assert.Nil(t, testWalUpdateViewChange())

	// Test Wal Write
	beginTime := uint64(time.Now().UnixNano())
	count, err := testWalWrite()
	if err != nil {
		t.Log("write error", err)
		t.Fatalf("%s", "write error")
	}
	t.Log("write total msg info", count)
	if count != times {
		t.Fatalf("%s", "write error")
	}
	endTime := uint64(time.Now().UnixNano())
	t.Log("write elapsed time", endTime-beginTime)

	// Test Wal Write
	beginTime = uint64(time.Now().UnixNano())
	count, err = testWalLoad()
	if err != nil {
		t.Log("load error", err)
		t.Fatalf("%s", "load error")
	}
	t.Log("total msg info", count)
	if count != times {
		t.Fatalf("%s", "load error")
	}
	endTime = uint64(time.Now().UnixNano())
	t.Log("load elapsed time", endTime-beginTime)

	// Test LevelDB
	if err = testLevelDB(); err != nil {
		t.Fatalf("%s", "TestLevelDB error")
	}
}

func getWal() Wal {
	if wal == nil {
		wal, _ = NewWal(nil, tempDir)
	}
	return wal
}

func testWalUpdateViewChange() error {
	// UpdateViewChange
	return getWal().UpdateViewChange(&ViewChangeMessage{
		Hash:   viewChangeHash,
		Number: viewChangeNumber,
	})
}

func testWalWrite() (int, error) {
	var err error
	// WriteJournal
	count := 0

	for i := 0; i < times; i++ {
		ordinal := ordinalMessages()
		if ordinal == 0 {
			err = getWal().WriteSync(buildSendPrepareBlock())
		} else if ordinal == 1 {
			err = getWal().WriteSync(buildSendPrepareVote())
		} else if ordinal == 2 {
			err = getWal().WriteSync(buildSendViewChange())
		} else if ordinal == 3 {
			err = getWal().WriteSync(buildConfirmedViewChange())
		}
		if err != nil {
			return 0, err
		}
		count++
	}
	getWal().Close() // force flush
	wal = nil
	return count, nil
}

func testWalLoad() (int, error) {
	var err error
	// LoadJournal
	count := 0
	err = getWal().Load(func(msg interface{}) {
		count++
	})
	if err != nil {
		return 0, err
	}
	getWal().Close() // force flush
	wal = nil
	return count, nil
}

func testLevelDB() error {
	path := filepath.Join(tempDir, "wal_meta")
	if db, err := leveldb.OpenFile(path, nil); err != nil {
		return err
	} else {
		data, err := db.Get([]byte("view-change"), nil)
		if err != nil {
			db.Close()
			return err
		}
		var v ViewChangeMeta
		if err := rlp.DecodeBytes(data, &v); err == nil {
			db.Close()
			if v.Number != 100 || v.Hash.Hex() != viewChangeHash.Hex() {
				return errors.New("TestLevelDB error")
			}
		}
	}
	return nil
}

func TestEmptyWal(t *testing.T) {
	wal := &emptyWal{}
	assert.Nil(t, wal.Write(nil))
	assert.Nil(t, wal.Load(func(msg interface{}) {}))
	assert.Nil(t, wal.UpdateViewChange(nil))
}
