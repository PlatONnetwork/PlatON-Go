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

package types

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSyncCache(t *testing.T) {
	cache := NewSyncCache(2 * time.Millisecond)
	cache.Add(1)
	assert.False(t, cache.AddOrReplace(1))
	time.Sleep(3 * time.Millisecond)
	assert.True(t, cache.AddOrReplace(1))
	cache.Remove(1)
	assert.True(t, cache.AddOrReplace(1))
	cache.AddOrReplace(2)
	cache.AddOrReplace(3)
	assert.Equal(t, 3, cache.Len())
	cache.Purge()
	assert.Equal(t, 0, cache.Len())

}

type mockCSMsg struct {
	epoch uint64
	view  uint64
}

func (m mockCSMsg) EpochNum() uint64 {
	return m.epoch
}

func (m mockCSMsg) ViewNum() uint64 {
	return m.view
}

func (m mockCSMsg) BlockNum() uint64 {
	panic("implement me")
}

func (m mockCSMsg) NodeIndex() uint32 {
	panic("implement me")
}

func (m mockCSMsg) CannibalizeBytes() ([]byte, error) {
	panic("implement me")
}

func (m mockCSMsg) Sign() []byte {
	panic("implement me")
}

func (m mockCSMsg) SetSign([]byte) {
	panic("implement me")
}
func (m mockCSMsg) String() string {
	panic("implement me")
}
func (m mockCSMsg) MsgHash() common.Hash {
	panic("implement me")
}
func (m mockCSMsg) BHash() common.Hash {
	panic("implement me")
}

func TestCSMsgPool(t *testing.T) {
	pool := NewCSMsgPool()

	pool.Purge(1, 1)

	defaultMsg := &MsgInfo{Msg: &mockCSMsg{epoch: 1, view: 1}}
	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareBlock(i, defaultMsg)
	}

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareVote(i, i+1, defaultMsg)
	}

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareQC(1, 1, i, defaultMsg)
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareBlock(1, 1, uint32(i)))
	}

	for i := uint32(10); i < 11; i++ {
		assert.Nil(t, pool.GetPrepareBlock(1, 1, uint32(i)))
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareVote(1, 1, i, i+1))
	}

	for i := uint32(10); i < 11; i++ {
		assert.Nil(t, pool.GetPrepareVote(1, 1, i, i+1))
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareQC(1, 1, uint32(i)))
	}

	for i := uint32(10); i < 11; i++ {
		assert.Nil(t, pool.GetPrepareQC(1, 1, uint32(i)))
	}

	//re-add
	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareBlock(i, defaultMsg)
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareBlock(1, 1, uint32(i)))
	}

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareQC(1, 1, i, defaultMsg)
	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareQC(1, 1, uint32(i)))
	}

	for i := uint32(0); i < 10; i++ {
		pool.AddPrepareVote(i, i+1, defaultMsg)
		pool.AddPrepareVote(i, i+2, defaultMsg)

	}

	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, pool.GetPrepareVote(1, 1, i, i+1))
		assert.NotNil(t, pool.GetPrepareVote(1, 1, i, i+2))
	}

}

func TestCSMsgPoolInvalidEpoch(t *testing.T) {
	pool := NewCSMsgPool()
	pool.Purge(1, 1)
	assert.Equal(t, false, pool.invalidEpochView(1, 2))
	assert.Equal(t, false, pool.invalidEpochView(1, 1))
	assert.Equal(t, false, pool.invalidEpochView(2, 0))
	assert.Equal(t, true, pool.invalidEpochView(1, 3))
	assert.Equal(t, true, pool.invalidEpochView(1, 0))
}
func TestCSMsgPoolPurge(t *testing.T) {
	pool := NewCSMsgPool()
	pool.Purge(1, 1)
	msgInfo := func(epoch, view uint64) *MsgInfo {
		return &MsgInfo{Msg: &mockCSMsg{epoch: epoch, view: view}}
	}
	pool.AddPrepareBlock(1, msgInfo(1, 1))
	pool.AddPrepareBlock(1, msgInfo(1, 2))

	assert.NotNil(t, pool.GetPrepareBlock(1, 1, 1))

	pool.Purge(1, 2)
	assert.Nil(t, pool.GetPrepareBlock(1, 1, 1))
	assert.NotNil(t, pool.GetPrepareBlock(1, 2, 1))

}

func TestCSMsgPoolAdd(t *testing.T) {
	pool := NewCSMsgPool()
	pool.Purge(1, 1)
	msgInfo := func(epoch, view uint64, inner bool) *MsgInfo {
		return &MsgInfo{Msg: &mockCSMsg{epoch: epoch, view: view}, Inner: inner}
	}

	pool.AddPrepareBlock(1, msgInfo(1, 1, true))
	assert.Nil(t, pool.GetPrepareBlock(1, 1, 1))
}
