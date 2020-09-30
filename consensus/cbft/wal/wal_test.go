// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package wal

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
)

var (
	epoch      = uint64(1)
	viewNumber = uint64(1)
	times      = 5000
)

func TestUpdateChainState(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	wal, _ := NewWal(nil, tempDir)
	// Test Wal UpdateChainState
	chainStateW, err := testWalUpdateChainState(wal)
	assert.Nil(t, err)

	// Test Wal LoadChainState
	chainStateR, err := testWalLoadChainState(wal)
	assert.Nil(t, err)

	if !chainStateR.Commit.EqualState(chainStateW.Commit) || !chainStateR.Lock.EqualState(chainStateW.Lock) || !chainStateR.QC[0].EqualState(chainStateW.QC[0]) {
		t.Fatalf("%s", "load chain state error")
	}
}

func TestWriteMsg(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	wal, _ := NewWal(nil, tempDir)
	wal.SetMockJournalLimitSize(1 * 1024)

	// Test Wal UpdateViewChange
	assert.Nil(t, testWalUpdateViewChange(wal))

	// Test Wal Write
	beginTime := uint64(time.Now().UnixNano())
	count, err := testWalWrite(wal)
	assert.Nil(t, err)
	t.Log("write total msg info", count)
	if count != times {
		t.Fatalf("%s", "write error")
	}
	endTime := uint64(time.Now().UnixNano())
	t.Log("write elapsed time", endTime-beginTime)

	// Test Wal load msg
	wal, _ = NewWal(nil, tempDir)
	beginTime = uint64(time.Now().UnixNano())
	count, err = testWalLoad(wal)
	assert.Nil(t, err)
	t.Log("total msg info", count)
	if count != times {
		t.Fatalf("%s", "load error")
	}
	endTime = uint64(time.Now().UnixNano())
	t.Log("load elapsed time", endTime-beginTime)
}

func testWalUpdateChainState(wal Wal) (*protocols.ChainState, error) {
	// UpdateChainState
	chainState := &protocols.ChainState{
		Commit: &protocols.State{
			Block:      newBlock(),
			QuorumCert: buildQuorumCert(),
		},
		Lock: &protocols.State{
			Block:      newBlock(),
			QuorumCert: buildQuorumCert(),
		},
		QC: []*protocols.State{{
			Block:      newBlock(),
			QuorumCert: buildQuorumCert(),
		}},
	}
	err := wal.UpdateChainState(chainState)
	return chainState, err
}

func testWalLoadChainState(wal Wal) (*protocols.ChainState, error) {
	var err error
	// Load chainState
	var chainState *protocols.ChainState
	err = wal.LoadChainState(func(cs *protocols.ChainState) error {
		chainState = cs
		return nil
	})
	if err != nil {
		return nil, err
	}
	return chainState, nil
}

func testWalUpdateViewChange(wal Wal) error {
	// UpdateViewChange
	return wal.UpdateViewChange(&ViewChangeMessage{
		Epoch:      epoch,
		ViewNumber: viewNumber,
	})
}

func TestUpdateViewChangeQC(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	wal, _ := NewWal(nil, tempDir)
	epochBegin, epochEnd := uint64(1), uint64(20)
	viewBegin, viewEnd := uint64(0), uint64(25)
	var i, j uint64

	// write
	for i = epochBegin; i < epochEnd; i++ {
		for j = viewBegin; j < viewEnd; j++ {
			wal.UpdateViewChangeQC(i, j, buildViewChangeQC())
			_, err := wal.GetViewChangeQC(i, j)
			assert.Nil(t, err)
		}
	}
	time.Sleep(time.Second)
	// Only keep last epoch
	for i = epochBegin; i < epochEnd; i++ {
		for j = viewBegin; j < viewEnd; j++ {
			if i == epochEnd-1 { // Only keep last epoch
				_, err := wal.GetViewChangeQC(i, j)
				assert.Nil(t, err)
			} else {
				_, err := wal.GetViewChangeQC(i, j)
				assert.NotNil(t, err)
			}
		}
	}
}

func testWalWrite(wal Wal) (int, error) {
	var err error
	// WriteJournal
	count := 0

	for i := 0; i < times; i++ {
		ordinal := ordinalMessages()
		if ordinal == 0 {
			err = wal.WriteSync(buildConfirmedViewChange())
		} else if ordinal == 1 {
			err = wal.WriteSync(buildSendViewChange())
		} else if ordinal == 2 {
			err = wal.WriteSync(buildSendPrepareBlock())
		} else if ordinal == 3 {
			//err = getWal().WriteSync(buildSendPrepareVote())
			err = wal.Write(buildSendPrepareVote())
		}
		if err != nil {
			return 0, err
		}
		count++
	}
	wal.Close() // force flush
	wal = nil
	return count, nil
}

func testWalLoad(wal Wal) (int, error) {
	var err error
	// LoadJournal
	count := 0
	err = wal.Load(func(msg interface{}) error {
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
	wal.Close() // force flush
	wal = nil
	return count, nil
}

func TestEmptyWal(t *testing.T) {
	wal, _ := NewWal(nil, "")
	assert.Nil(t, wal.UpdateChainState(nil))
	assert.Nil(t, wal.LoadChainState(func(chainState *protocols.ChainState) error { return nil }))
	assert.Nil(t, wal.UpdateViewChangeQC(1, 0, nil))
	_, err := wal.GetViewChangeQC(1, 0)
	assert.Nil(t, err)
	assert.Nil(t, wal.Write(nil))
	assert.Nil(t, wal.WriteSync(nil))
	assert.Nil(t, wal.Load(func(msg interface{}) error { return nil }))
	assert.Nil(t, wal.UpdateViewChange(nil))
	wal.SetMockJournalLimitSize(0)
	wal.Close()
}

func TestWalDecoder(t *testing.T) {
	timestamp := uint64(time.Now().UnixNano())
	// SendPrepareBlockMsg
	prepare := &MessageSendPrepareBlock{
		Timestamp: timestamp,
		Data:      buildSendPrepareBlock(),
	}
	data, _ := rlp.EncodeToBytes(prepare)
	_, err := WALDecode(data, protocols.SendPrepareBlockMsg)
	assert.Nil(t, err)
	_, err = WALDecode(data, protocols.SendPrepareVoteMsg)
	assert.NotNil(t, err)
	// MessageSendPrepareVote
	vote := &MessageSendPrepareVote{
		Timestamp: timestamp,
		Data:      buildSendPrepareVote(),
	}
	data, _ = rlp.EncodeToBytes(vote)
	_, err = WALDecode(data, protocols.SendPrepareVoteMsg)
	assert.Nil(t, err)
	_, err = WALDecode(data, protocols.SendViewChangeMsg)
	assert.NotNil(t, err)
	// MessageSendPrepareVote
	view := &MessageSendViewChange{
		Timestamp: timestamp,
		Data:      buildSendViewChange(),
	}
	data, _ = rlp.EncodeToBytes(view)
	_, err = WALDecode(data, protocols.SendViewChangeMsg)
	assert.Nil(t, err)
	_, err = WALDecode(data, protocols.ConfirmedViewChangeMsg)
	assert.NotNil(t, err)
	// MessageSendPrepareVote
	confirm := &MessageConfirmedViewChange{
		Timestamp: timestamp,
		Data:      buildConfirmedViewChange(),
	}
	data, _ = rlp.EncodeToBytes(confirm)
	_, err = WALDecode(data, protocols.ConfirmedViewChangeMsg)
	assert.Nil(t, err)
	_, err = WALDecode(data, protocols.SendPrepareBlockMsg)
	assert.NotNil(t, err)
}
