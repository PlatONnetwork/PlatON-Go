package snapshotdb

import (
	"fmt"
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

func (s *snapshotDB) schedule() {
	if counter.get() == 60 || s.current.HighestNum.Int64()-s.current.BaseNum.Int64() >= 100 {
		if _, err := s.Compaction(); err != nil {
			logger.Error(fmt.Sprint("[SnapshotDB]compaction fail:", err))
		}
		counter.reset()
	} else {
		counter.increment()
	}
}
