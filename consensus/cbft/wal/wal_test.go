package wal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	epoch      = uint64(1)
	viewNumber = uint64(1)
	times      = 1000
	tempDir    string
	wal        Wal
)

func TestWal(t *testing.T) {
	tempDir, _ = ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	// Test Wal UpdateChainState
	chainStateW, err := testWalUpdateChainState()
	if err != nil {
		t.Fatalf("%s", "update chain state error")
	}

	// Test Wal LoadChainState
	chainStateR, err := testWalLoadChainState()
	if err != nil {
		t.Fatalf("%s", "load chain state error")
	}
	if chainStateR.Commit.Block.Hash() != chainStateW.Commit.Block.Hash() || chainStateR.Commit.QuorumCert.ViewNumber != chainStateW.Commit.QuorumCert.ViewNumber ||
		chainStateR.Lock.Block.Hash() != chainStateW.Lock.Block.Hash() || chainStateR.Lock.QuorumCert.ViewNumber != chainStateW.Lock.QuorumCert.ViewNumber ||
		chainStateR.QC[0].Block.Hash() != chainStateW.QC[0].Block.Hash() || chainStateR.QC[0].QuorumCert.ViewNumber != chainStateW.QC[0].QuorumCert.ViewNumber {
		t.Fatalf("%s", "load chain state error")
	}

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

	// Test Wal load msg
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

func testWalUpdateChainState() (*protocols.ChainState, error) {
	// UpdateChainState
	chainState := &protocols.ChainState{
		Commit: &protocols.State{
			Block:      block,
			QuorumCert: buildQuorumCert(),
		},
		Lock: &protocols.State{
			Block:      block,
			QuorumCert: buildQuorumCert(),
		},
		QC: []*protocols.State{{
			Block:      block,
			QuorumCert: buildQuorumCert(),
		}},
	}
	err := getWal().UpdateChainState(chainState)
	return chainState, err
}

func testWalLoadChainState() (*protocols.ChainState, error) {
	var err error
	// Load chainState
	var chainState *protocols.ChainState
	err = getWal().LoadChainState(func(cs *protocols.ChainState) error {
		chainState = cs
		return nil
	})
	if err != nil {
		return nil, err
	}
	return chainState, nil
}

func testWalUpdateViewChange() error {
	// UpdateViewChange
	return getWal().UpdateViewChange(&ViewChangeMessage{
		Epoch:      epoch,
		ViewNumber: viewNumber,
	})
}

func testWalWrite() (int, error) {
	var err error
	// WriteJournal
	count := 0

	for i := 0; i < times; i++ {
		ordinal := ordinalMessages()
		if ordinal == 0 {
			err = getWal().WriteSync(buildConfirmedViewChange())
		} else if ordinal == 1 {
			err = getWal().WriteSync(buildSendViewChange())
		} else if ordinal == 2 {
			err = getWal().WriteSync(buildSendPrepareBlock())
		} else if ordinal == 3 {
			//err = getWal().WriteSync(buildSendPrepareVote())
			err = getWal().Write(buildSendPrepareVote())
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
	err = getWal().Load(func(msg interface{}) error {
		switch msg.(type) {
		case *protocols.ConfirmedViewChange:
			count++
		case *protocols.SendViewChange:
			count++
		case *protocols.SendPrepareBlock:
			count++
		case *protocols.SendPrepareVote:
			count++
		}
		return nil
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
		data, err := db.Get(viewChangeKey, nil)
		if err != nil {
			db.Close()
			return err
		}
		var vc ViewChangeMessage
		if err := rlp.DecodeBytes(data, &vc); err == nil {
			db.Close()
			if vc.Epoch != epoch || vc.ViewNumber != viewNumber {
				return errors.New("TestLevelDB error")
			}
		}
	}
	return nil
}

func TestEmptyWal(t *testing.T) {
	wal := &emptyWal{}
	assert.Nil(t, wal.Write(nil))
	assert.Nil(t, wal.Load(func(msg interface{}) error { return nil }))
	assert.Nil(t, wal.UpdateViewChange(nil))
}

func TestRlpDecode(t *testing.T) {
	msg := &Message{
		Timestamp: uint64(time.Now().UnixNano()),
		Data:      buildConfirmedViewChange(),
	}
	data, _ := rlp.EncodeToBytes(msg)

	var j MessageConfirmedViewChange
	if err := rlp.DecodeBytes(data, &j); err == nil {
		fmt.Println(fmt.Sprintf("decode msg,epoch:%d,viewNumber:%d", j.Data.Epoch, j.Data.ViewNumber))
	} else {
		t.Fatalf("%s", "rlp decode error")
	}
}
