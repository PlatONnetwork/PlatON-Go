package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	viewChangeNumber = uint64(100)
	viewChangeHash   = common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065747")
	times            = 1000000
	tempDir          string
	wal              Wal
)

func TestMain(m *testing.M) {
	fmt.Println("begin test wal")
	tempDir, _ = ioutil.TempDir("", "wal")
	m.Run()
	os.RemoveAll(tempDir)
	fmt.Println("end test wal")
}

func getWal() Wal {
	if wal == nil {
		wal, _ = NewWal(nil, tempDir)
	}
	return wal
}

func TestWalUpdateViewChange(t *testing.T) {
	// UpdateViewChange
	getWal().UpdateViewChange(&ViewChangeMessage{
		Hash:   viewChangeHash,
		Number: viewChangeNumber,
	})
}

func TestWalWrite(t *testing.T) {
	var err error
	// WriteJournal
	beginTime := uint64(time.Now().UnixNano())
	count := 0

	for i := 0; i < times; i++ {
		ordinal := ordinalMessages()
		if ordinal == 0 {
			err = getWal().Write(&MsgInfo{
				Msg: buildPrepareBlock(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 1 {
			err = getWal().Write(&MsgInfo{
				Msg: buildPrepareVote(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 2 {
			err = getWal().Write(&MsgInfo{
				Msg: buildViewChange(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 3 {
			err = getWal().Write(&MsgInfo{
				Msg: buildviewChangeVote(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 4 {
			err = getWal().Write(&MsgInfo{
				Msg: buildConfirmedPrepareBlock(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 5 {
			err = getWal().Write(&MsgInfo{
				Msg: buildGetPrepareVote(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 6 {
			err = getWal().Write(&MsgInfo{
				Msg: buildPrepareVotes(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 7 {
			err = getWal().Write(&MsgInfo{
				Msg: buildGetPrepareBlock(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 8 {
			err = getWal().Write(&MsgInfo{
				Msg: buildGetHighestPrepareBlock(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 9 {
			err = getWal().Write(&MsgInfo{
				Msg: buildHighestPrepareBlock(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 10 {
			err = getWal().Write(&MsgInfo{
				Msg: buildCbftStatusData(),
				PeerID: buildPeerId(),
			})
		} else if ordinal == 11 {
			err = getWal().Write(&MsgInfo{
				Msg: buildPrepareBlockHash(),
				PeerID: buildPeerId(),
			})
		}
		if err != nil {
			t.Log("write error", err)
			t.Errorf("%s", "write error")
		}
		count ++
	}
	getWal().Close() // force flush
	wal = nil
	t.Log("write total msg info", count)
	if count != times {
		t.Errorf("%s", "write error")
	}
	endTime := uint64(time.Now().UnixNano())
	t.Log("write elapsed time", endTime-beginTime)
}

func TestWalLoad(t *testing.T) {
	var err error
	// LoadJournal
	beginTime := uint64(time.Now().UnixNano())
	count := 0
	err = getWal().Load(func(info *MsgInfo) {
		count ++
	})
	if err != nil {
		t.Log("load error", err)
		t.Errorf("%s", "load error")
	}
	getWal().Close() // force flush
	wal = nil
	t.Log("total msg info", count)
	if count != times {
		t.Errorf("%s", "load error")
	}
	endTime := uint64(time.Now().UnixNano())
	t.Log("load elapsed time", endTime-beginTime)

}

func TestLevelDB(t *testing.T) {
	path := filepath.Join(tempDir, "wal_meta")
	if db, err := leveldb.OpenFile(path, nil); err != nil {
		t.Errorf("%s", "TestLevelDB error")
	} else {
		data, err := db.Get([]byte("view-change"), nil)
		if err == nil {
			var v ViewChangeMeta
			if err := rlp.DecodeBytes(data, &v); err == nil {
				t.Log(v.Number)
				t.Log(v.Hash.Hex())
				t.Log(v.FileID)
				t.Log(v.Seq)
				db.Close()
				if v.Number != 100 || v.Hash.Hex() != viewChangeHash.Hex() {
					t.Errorf("%s", "TestLevelDB error")
				}
			}
		} else {
			db.Close()
			t.Errorf("%s", "TestLevelDB error")
		}
	}
}

func TestEmptyWal(t *testing.T) {
	wal := &emptyWal{}
	assert.Nil(t, wal.Write(nil))
	assert.Nil(t, wal.Load(func(info *MsgInfo) {

	}))
	assert.Nil(t, wal.UpdateViewChange(nil))
	assert.Nil(t, wal.UpdateViewChange(nil))
}