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

package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := time.NewTimer(0)
	timer.Stop()

	viewTimer := &viewTimer{timer: timer, timeInterval: viewTimeInterval{baseMs: uint64(1 * time.Second), exponentBase: exponentBase, maxExponent: maxExponent}}
	viewTimer.setupTimer(1)
	assert.False(t, viewTimer.isDeadline())
	select {
	case <-viewTimer.timerChan():
		assert.True(t, viewTimer.isDeadline())
	}
}

func TestCalViewInterval(t *testing.T) {
	type views struct {
		in  uint64
		out uint64
	}
	testcases := [][]views{
		{{2, 2}, {1, 1}, {2, 2}},
		{{3, 3}, {1, 2}, {1, 1}},
		{{1, 1}, {2, 2}, {2, 3}, {3, 4}, {2, 3}, {2, 2}, {2, 3}},
		{{3, 3}, {2, 2}, {2, 3}, {1, 2}, {1, 1}, {1, 1}},
		{{1, 1}, {1, 1}, {1, 1}, {1, 1}},
		{{2, 2}, {2, 3}, {2, 2}, {2, 3}},
	}
	for row, test := range testcases {
		timer := newViewTimer(10)
		timer.calViewInterval(1)
		for cul, c := range test {
			//fmt.Printf("row:%d, cul:%d, pre:%d in:%d, out:%d\n", row, cul, timer.preViewInterval, c.in, c.out)
			assert.Equal(t, c.out, timer.calViewInterval(c.in), "row:%d, cul:%d, pre:%d in:%d, out:%d", row, cul, timer.preViewInterval, c.in, c.out)
		}
	}

}

func TestInterval(t *testing.T) {
	in := viewTimeInterval{
		uint64(10), 2, 2,
	}

	assert.Equal(t, uint64(10), uint64(in.getViewTimeInterval(1)))
	assert.Equal(t, uint64(20), uint64(in.getViewTimeInterval(2)))

	in = viewTimeInterval{
		uint64(10 * time.Second), 1.5, 2,
	}

	assert.Equal(t, uint64(10*time.Second), uint64(in.getViewTimeInterval(1)))
	assert.Equal(t, uint64(15*time.Second), uint64(in.getViewTimeInterval(2)))
	assert.Equal(t, uint64(22*time.Second+500*time.Millisecond), uint64(in.getViewTimeInterval(3)))

}
