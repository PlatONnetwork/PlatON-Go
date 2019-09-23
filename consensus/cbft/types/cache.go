package types

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/log"
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

type CSMsgPool struct {
	prepareBlocks map[uint32]*MsgInfo
	prepareVotes  map[uint32]map[uint32]*MsgInfo
	blockMetric   map[uint32]uint32
	voteMetric    map[uint32]map[uint32]uint32
}

func NewCSMsgPool() *CSMsgPool {
	return &CSMsgPool{
		prepareBlocks: make(map[uint32]*MsgInfo),
		prepareVotes:  make(map[uint32]map[uint32]*MsgInfo),
		blockMetric:   make(map[uint32]uint32),
		voteMetric:    make(map[uint32]map[uint32]uint32),
	}
}

// Add prepare block to cache. There is no strict distinction between the blocks of the current view,
// which will be cleared from the cache after each acquisition, so that the new block can re-enter the cache.
func (cs *CSMsgPool) AddPrepareBlock(blockIndex uint32, msg *MsgInfo) {
	cs.prepareBlocks[blockIndex] = msg
}

// Get prepare block and clear it from the cache
func (cs *CSMsgPool) GetPrepareBlock(index uint32) *MsgInfo {
	if m, ok := cs.prepareBlocks[index]; ok {
		cs.addBlockMetric(index)
		delete(cs.prepareBlocks, index)
		return m
	}
	return nil
}

func (cs *CSMsgPool) addBlockMetric(index uint32) {
	if m, ok := cs.blockMetric[index]; ok {
		cs.blockMetric[index] = m + 1
	} else {
		cs.blockMetric[index] = 1
	}
}

func (cs *CSMsgPool) getBlockMetric(index uint32) uint32 {
	if m, ok := cs.blockMetric[index]; ok {
		return m
	}
	return 0
}

// Add prepare votes to cache. There is no strict distinction between the votes of the current view,
// which will be cleared from the cache after each acquisition, so that the new votes can re-enter the cache.
func (cs *CSMsgPool) AddPrepareVote(blockIndex uint32, validatorIndex uint32, msg *MsgInfo) {
	if votes, ok := cs.prepareVotes[blockIndex]; ok {
		votes[validatorIndex] = msg
	} else {
		votes := make(map[uint32]*MsgInfo)
		votes[validatorIndex] = msg
		cs.prepareVotes[blockIndex] = votes
	}
}

// Get prepare vote and clear it from the cache
func (cs *CSMsgPool) GetPrepareVote(blockIndex uint32, validatorIndex uint32) *MsgInfo {
	if p, ok := cs.prepareVotes[blockIndex]; ok {
		if m, ok := p[validatorIndex]; ok {
			cs.addVoteMetric(blockIndex, validatorIndex)
			delete(p, validatorIndex)
			return m
		}
	}
	return nil
}

func (cs *CSMsgPool) addVoteMetric(blockIndex uint32, validatorIndex uint32) {
	if votes, ok := cs.voteMetric[blockIndex]; ok {
		if m, ok := votes[validatorIndex]; ok {
			votes[validatorIndex] = m + 1
		} else {
			votes[validatorIndex] = 1
		}
	} else {
		votes := make(map[uint32]uint32)
		votes[validatorIndex] = 1
		cs.voteMetric[blockIndex] = votes
	}
}

func (cs *CSMsgPool) getVoteMetric(blockIndex uint32, validatorIndex uint32) uint32 {
	if votes, ok := cs.voteMetric[blockIndex]; ok {
		if m, ok := votes[validatorIndex]; ok {
			return m
		}
	}
	return 0
}

func (cs *CSMsgPool) Purge() {

	for k, v := range cs.blockMetric {
		log.Debug(fmt.Sprintf("pool block index:%d, count:%d", k, v))
	}

	for k, v := range cs.voteMetric {
		for vl, c := range v {
			log.Debug(fmt.Sprintf("pool vote index:%d, validator:%d, count:%d", k, vl, c))
		}
	}

	cs.prepareBlocks = make(map[uint32]*MsgInfo)
	cs.prepareVotes = make(map[uint32]map[uint32]*MsgInfo)

	cs.blockMetric = make(map[uint32]uint32)
	cs.voteMetric = make(map[uint32]map[uint32]uint32)

}
