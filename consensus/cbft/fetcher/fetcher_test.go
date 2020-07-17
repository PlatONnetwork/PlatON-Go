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

package fetcher

import (
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestFetcher_AddTask(t *testing.T) {
	fetcher := NewFetcher()
	fetcher.Start()
	var w sync.WaitGroup
	w.Add(1)
	fetcher.AddTask("add", func(message types.Message) bool {
		return true
	}, func(message types.Message) {
		t.Log("add")
	}, func() {
		t.Log("timeout add")
		w.Done()
	})

	w.Wait()
	w.Add(1)
	fetcher.AddTask("add", func(message types.Message) bool {
		return true
	}, func(message types.Message) {
		t.Log("add")
	}, func() {
		t.Log("timeout add")
		w.Done()
	})

	w.Wait()

	assert.Equal(t, fetcher.Len(), 0)
	fetcher.Stop()
}

func TestFetcher_MatchTask(t *testing.T) {
	fetcher := NewFetcher()
	fetcher.Start()
	var w sync.WaitGroup
	w.Add(1)
	fetcher.AddTask("add", func(message types.Message) bool {
		_, ok := message.(*protocols.PrepareBlock)
		return ok
	}, func(message types.Message) {
		w.Done()
	}, func() {
		t.Error("timeout add ")
		w.Done()

	})
	time.Sleep(10 * time.Millisecond)
	fetcher.MatchTask("add", &protocols.PrepareBlock{})
	w.Wait()
	assert.Equal(t, fetcher.Len(), 0)
	fetcher.Stop()
}
