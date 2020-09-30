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
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"sync"
	"time"
)

type SyncCache struct {
	lock    sync.RWMutex
	items   map[interface{}]time.Time
	timeout time.Duration
}

func NewSyncCache(timeout time.Duration) *SyncCache {
	cache := &SyncCache{
		items:   make(map[interface{}]time.Time),
		timeout: timeout,
	}
	return cache
}

func (s *SyncCache) Add(v interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.items[v] = time.Now()
}

func (s *SyncCache) AddOrReplace(v interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if t, ok := s.items[v]; ok {
		if time.Since(t) < s.timeout {
			return false
		}
	}
	s.items[v] = time.Now()
	return true
}

func (s *SyncCache) Remove(v interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	delete(s.items, v)
}

func (s *SyncCache) Purge() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.items = make(map[interface{}]time.Time)
}
func (s *SyncCache) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

type viewCache struct {
	prepareBlocks map[uint32]*MsgInfo
	prepareVotes  map[uint32]map[uint32]*MsgInfo
	prepareQC     map[uint32]*MsgInfo
	blockMetric   map[uint32]uint32
	voteMetric    map[uint32]map[uint32]uint32
	qcMetric      map[uint32]uint32
}

func newViewCache() *viewCache {
	return &viewCache{
		prepareBlocks: make(map[uint32]*MsgInfo),
		prepareVotes:  make(map[uint32]map[uint32]*MsgInfo),
		prepareQC:     make(map[uint32]*MsgInfo),
		blockMetric:   make(map[uint32]uint32),
		voteMetric:    make(map[uint32]map[uint32]uint32),
		qcMetric:      make(map[uint32]uint32),
	}
}

// Add prepare block to cache.
func (v *viewCache) addPrepareBlock(blockIndex uint32, msg *MsgInfo) {
	v.prepareBlocks[blockIndex] = msg
}

// Get prepare block and clear it from the cache
func (v *viewCache) getPrepareBlock(index uint32) *MsgInfo {
	if m, ok := v.prepareBlocks[index]; ok {
		v.addBlockMetric(index)
		delete(v.prepareBlocks, index)
		return m
	}
	return nil
}

func (v *viewCache) addBlockMetric(index uint32) {
	if m, ok := v.blockMetric[index]; ok {
		v.blockMetric[index] = m + 1
	} else {
		v.blockMetric[index] = 1
	}
}

// Add prepare block to cache.
func (v *viewCache) addPrepareQC(blockIndex uint32, msg *MsgInfo) {
	v.prepareQC[blockIndex] = msg
}

// Get prepare QC and clear it from the cache
func (v *viewCache) getPrepareQC(index uint32) *MsgInfo {
	if m, ok := v.prepareQC[index]; ok {
		v.addQCMetric(index)
		delete(v.prepareQC, index)
		return m
	}
	return nil
}

func (v *viewCache) getBlockMetric(index uint32) uint32 {
	if m, ok := v.blockMetric[index]; ok {
		return m
	}
	return 0
}

func (v *viewCache) addQCMetric(index uint32) {
	if m, ok := v.qcMetric[index]; ok {
		v.qcMetric[index] = m + 1
	} else {
		v.qcMetric[index] = 1
	}
}

func (v *viewCache) getQCMetric(index uint32) uint32 {
	if m, ok := v.qcMetric[index]; ok {
		return m
	}
	return 0
}

// Add prepare votes to cache.
func (v *viewCache) addPrepareVote(blockIndex uint32, validatorIndex uint32, msg *MsgInfo) {
	if votes, ok := v.prepareVotes[blockIndex]; ok {
		votes[validatorIndex] = msg
	} else {
		votes := make(map[uint32]*MsgInfo)
		votes[validatorIndex] = msg
		v.prepareVotes[blockIndex] = votes
	}
}

// Get prepare vote and clear it from the cache
func (v *viewCache) getPrepareVote(blockIndex uint32, validatorIndex uint32) *MsgInfo {
	if p, ok := v.prepareVotes[blockIndex]; ok {
		if m, ok := p[validatorIndex]; ok {
			v.addVoteMetric(blockIndex, validatorIndex)
			delete(p, validatorIndex)
			return m
		}
	}
	return nil
}

func (v *viewCache) addVoteMetric(blockIndex uint32, validatorIndex uint32) {
	if votes, ok := v.voteMetric[blockIndex]; ok {
		if m, ok := votes[validatorIndex]; ok {
			votes[validatorIndex] = m + 1
		} else {
			votes[validatorIndex] = 1
		}
	} else {
		votes := make(map[uint32]uint32)
		votes[validatorIndex] = 1
		v.voteMetric[blockIndex] = votes
	}
}

func (v *viewCache) getVoteMetric(blockIndex uint32, validatorIndex uint32) uint32 {
	if votes, ok := v.voteMetric[blockIndex]; ok {
		if m, ok := votes[validatorIndex]; ok {
			return m
		}
	}
	return 0
}

type epochCache struct {
	views map[uint64]*viewCache
}

func newEpochCache() *epochCache {
	return &epochCache{
		views: make(map[uint64]*viewCache),
	}
}

func (e *epochCache) matchViewCache(view uint64) *viewCache {
	for k, v := range e.views {
		if k == view {
			return v
		}
	}
	newView := newViewCache()
	e.views[view] = newView
	return newView
}

func (e *epochCache) findViewCache(view uint64) *viewCache {
	for k, v := range e.views {
		if k == view {
			return v
		}
	}
	return nil
}

func (e *epochCache) purge(view uint64) {
	for k, _ := range e.views {
		if k < view {
			delete(e.views, k)
		}
	}
}

type CSMsgPool struct {
	epochs   map[uint64]*epochCache
	minEpoch uint64
	minView  uint64
}

func NewCSMsgPool() *CSMsgPool {
	return &CSMsgPool{
		epochs:   make(map[uint64]*epochCache),
		minEpoch: math.MaxInt64,
		minView:  math.MaxInt64,
	}
}

func (cs *CSMsgPool) invalidEpochView(epoch, view uint64) bool {
	if cs.minEpoch == epoch && cs.minView == view ||
		cs.minEpoch == epoch && cs.minView+1 == view ||
		cs.minEpoch+1 == epoch && view == 0 {
		return false
	}
	return true
}

// Add prepare block to cache.
func (cs *CSMsgPool) AddPrepareBlock(blockIndex uint32, msg *MsgInfo) {
	if csMsg, ok := msg.Msg.(ConsensusMsg); ok {
		if cs.invalidEpochView(csMsg.EpochNum(), csMsg.ViewNum()) || msg.Inner {
			return
		}
		cs.matchEpochCache(csMsg.EpochNum()).
			matchViewCache(csMsg.ViewNum()).
			addPrepareBlock(blockIndex, msg)
	}
}

// Get prepare block and clear it from the cache
func (cs *CSMsgPool) GetPrepareBlock(epoch, view uint64, index uint32) *MsgInfo {
	if cs.invalidEpochView(epoch, view) {
		return nil
	}

	viewCache := cs.findViewCache(epoch, view)
	if viewCache != nil {
		return viewCache.getPrepareBlock(index)
	}
	return nil
}

// Add prepare block to cache.
func (cs *CSMsgPool) AddPrepareQC(epoch, view uint64, blockIndex uint32, msg *MsgInfo) {
	if cs.invalidEpochView(epoch, view) || msg.Inner {
		return
	}

	cs.matchEpochCache(epoch).
		matchViewCache(epoch).
		addPrepareQC(blockIndex, msg)
}

// Get prepare QC and clear it from the cache
func (cs *CSMsgPool) GetPrepareQC(epoch, view uint64, index uint32) *MsgInfo {
	if cs.invalidEpochView(epoch, view) {
		return nil
	}

	viewCache := cs.findViewCache(epoch, view)
	if viewCache != nil {
		return viewCache.getPrepareQC(index)
	}
	return nil
}

// Add prepare votes to cache.
func (cs *CSMsgPool) AddPrepareVote(blockIndex uint32, validatorIndex uint32, msg *MsgInfo) {
	if csMsg, ok := msg.Msg.(ConsensusMsg); ok {
		if cs.invalidEpochView(csMsg.EpochNum(), csMsg.ViewNum()) || msg.Inner {
			return
		}
		cs.matchEpochCache(csMsg.EpochNum()).
			matchViewCache(csMsg.ViewNum()).
			addPrepareVote(blockIndex, validatorIndex, msg)
	}
}

// Get prepare vote and clear it from the cache
func (cs *CSMsgPool) GetPrepareVote(epoch, view uint64, blockIndex uint32, validatorIndex uint32) *MsgInfo {
	if cs.invalidEpochView(epoch, view) {
		return nil
	}

	viewCache := cs.findViewCache(epoch, view)
	if viewCache != nil {
		return viewCache.getPrepareVote(blockIndex, validatorIndex)
	}
	return nil
}

func (cs *CSMsgPool) Purge(epoch, view uint64) {

	for k, v := range cs.epochs {
		if k < epoch {
			delete(cs.epochs, k)
		} else {
			v.purge(view)
		}
	}
	cs.minEpoch = epoch
	cs.minView = view
}

func (cs *CSMsgPool) matchEpochCache(epoch uint64) *epochCache {
	for k, v := range cs.epochs {
		if k == epoch {
			return v
		}
	}
	newEpoch := newEpochCache()
	cs.epochs[epoch] = newEpoch
	return newEpoch
}

func (cs *CSMsgPool) findEpochCache(epoch uint64) *epochCache {
	for k, v := range cs.epochs {
		if k == epoch {
			return v
		}
	}
	return nil
}

func (cs *CSMsgPool) findViewCache(epoch, view uint64) *viewCache {
	epochCache := cs.findEpochCache(epoch)
	if epochCache != nil {
		return epochCache.findViewCache(view)
	}
	return nil
}
