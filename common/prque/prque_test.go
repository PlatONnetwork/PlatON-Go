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

package prque

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrque(t *testing.T) {
	queue := New(nil)
	assert.True(t, queue.Empty())
	queue.Push("item1", 1)
	queue.Push("item5", 5)
	queue.Push("item3", 3)
	queue.Push("item2", 2)
	queue.Push("item4", 4)
	assert.False(t, queue.Empty())
	assert.Equal(t, queue.Size(), 5)
	value, priority := queue.Pop()
	assert.Equal(t, value, "item5")
	assert.Equal(t, priority, int64(5))
	assert.Equal(t, queue.Size(), 4)

	value = queue.PopItem()
	assert.Equal(t, value, "item4")
	assert.Equal(t, queue.Size(), 3)

	queue.Remove(0) // remove item3
	value, priority = queue.Pop()
	assert.Equal(t, value, "item2")
	assert.Equal(t, priority, int64(2))
	assert.Equal(t, queue.Size(), 1)

	queue.Reset()
	assert.True(t, queue.Empty())
	assert.Equal(t, queue.Size(), 0)
}
