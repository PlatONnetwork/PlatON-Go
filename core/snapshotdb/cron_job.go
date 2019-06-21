package snapshotdb

import "sync/atomic"

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
