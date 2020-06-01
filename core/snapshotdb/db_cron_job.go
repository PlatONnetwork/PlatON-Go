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
	if counter.get() >= 60 || s.current.GetHighest(false).Num.Uint64()-s.current.GetBase(false).Num.Uint64() >= 100 {
		//only one compaction can execute
		if atomic.CompareAndSwapInt32(&s.snapshotLockC, snapshotUnLock, snapshotLock) {
			if err := s.Compaction(); err != nil {
				logger.Error("compaction fail", "err", err)
				s.dbError = err
				s.corn.Stop()
			}
			counter.reset()
			atomic.StoreInt32(&s.snapshotLockC, snapshotUnLock)
			return
		}
		logger.Info("snapshotDB is still Compaction Lock,wait for next schedule")
	}
	counter.increment()
}
