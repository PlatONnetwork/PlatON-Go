package fetcher

import (
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	keepTimeout = 100 * time.Millisecond
)

// LimitFetcher tag the fetching request that is happening
// Limit the frequency of the same requests
type LimitFetcher struct {
	lock     sync.Mutex
	fetching map[common.Hash]time.Time
	quit     chan struct{}
}

// NewLimitFetcher returns a new pointer to the LimitFetcher.
func NewLimitFetcher() *LimitFetcher {
	fetcher := &LimitFetcher{
		fetching: make(map[common.Hash]time.Time),
		quit:     make(chan struct{}),
	}

	go fetcher.loop()
	return fetcher
}

// AddTask adds a fetcher task.
func (f *LimitFetcher) AddTask(id common.Hash) bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	if _, ok := f.fetching[id]; ok {
		return false
	}
	f.fetching[id] = time.Now()
	return true
}

func (f *LimitFetcher) loop() {
	keepTimer := time.NewTicker(keepTimeout)
	for {
		select {
		case <-keepTimer.C:
			f.lock.Lock()
			for id, t := range f.fetching {
				if time.Since(t) > keepTimeout {
					delete(f.fetching, id)
				}
			}
			f.lock.Unlock()
		case <-f.quit:
			f.lock.Lock()
			f.fetching = make(map[common.Hash]time.Time)
			keepTimer.Stop()
			f.lock.Unlock()
			return
		}
	}
}

// Stop turns off for LimitFetcher.
func (f *LimitFetcher) Stop() {
	close(f.quit)
}
