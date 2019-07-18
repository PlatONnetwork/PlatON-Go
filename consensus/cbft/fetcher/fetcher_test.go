package fetcher

import (
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"sync"
	"testing"
)

func TestFetcher_AddTask(t *testing.T) {
	fetcher := NewFetcher()
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

}
