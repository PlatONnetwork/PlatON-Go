package snapshotdb

import (
	"sync/atomic"
)

var counter count32

type count32 int32

func (c *count32) increment() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *count32) reset() {
	atomic.StoreInt32((*int32)(c), 0)
}

func (c *count32) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

const (
	snapshotUnLock = 0
	snapshotLock   = 1
)

func (s *snapshotDB) schedule() {
	// Compaction condition , last Compaction execute time gt than 60s or commit block num gt than 100
	if counter.get() >= 60 || s.current.HighestNum.Int64()-s.current.BaseNum.Int64() >= 100 {
		//only one compaction can execute
		if atomic.CompareAndSwapInt32(&s.snapshotLockC, snapshotUnLock, snapshotLock) {
			if err := s.Compaction(); err != nil {
				logger.Error("compaction fail", "err", err)
			}
			counter.reset()
			atomic.StoreInt32(&s.snapshotLockC, snapshotUnLock)
			return
		}
		logger.Info("snapshotDB is still Compaction Lock,wait for next schedule")
	}
	counter.increment()
}
