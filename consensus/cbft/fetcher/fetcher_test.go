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
