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
	"io/ioutil"
	"path"

	"github.com/PlatONnetwork/PlatON-Go/metrics"
)

var (
	dbSizeGauge = metrics.NewRegisteredGauge("snapshotdb/basedb/size", nil)
	dbForkGauge = metrics.NewRegisteredGauge("snapshotdb/fork", nil)
)

func walkDir(dir string) int64 {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Error("read dir fail", "err", err)
		return 0
	}
	var dirSize int64
	for _, f := range entries {
		if f.IsDir() {
			dirSize = walkDir(path.Join(dir, f.Name())) + dirSize
		} else {
			dirSize = f.Size() + dirSize
		}
	}
	return dirSize
}

func (s *snapshotDB) metrics() {
	// metric size
	size := walkDir(s.path)
	dbSizeGauge.Update(size)
	// metric fork num
	forkNumList := make(map[int64]int)
	var forkMax int
	s.unCommit.RLock()
	for _, value := range s.unCommit.blocks {
		if forkSum, ok := forkNumList[value.Number.Int64()]; ok {
			forkIncr := forkSum + 1
			forkNumList[value.Number.Int64()] = forkIncr
			if forkIncr > forkMax {
				forkMax = forkIncr
			}
		} else {
			forkNumList[value.Number.Int64()] = 1
		}
	}
	s.unCommit.RUnlock()
	dbForkGauge.Update(int64(forkMax))
}
