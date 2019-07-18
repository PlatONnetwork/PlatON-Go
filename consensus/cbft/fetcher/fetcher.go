package fetcher

import (
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"sync"
	"time"
)

var (
	arriveTimeout = 500 * time.Millisecond
)

type MatchFunc func(types.Message) bool
type ExecutorFunc func(types.Message)
type ExpireFunc func()

type task struct {
	match    MatchFunc
	executor ExecutorFunc
	expire   ExpireFunc
	time     time.Time
}

type Fetcher struct {
	lock    sync.Mutex
	newTask chan struct{}
	task    map[string]*task
}

func NewFetcher() *Fetcher {
	fetcher := &Fetcher{
		newTask: make(chan struct{}, 1),
		task:    make(map[string]*task),
	}
	go fetcher.loop()
	return fetcher
}

// Add a fetcher task
func (f *Fetcher) AddTask(id string, match MatchFunc, executor ExecutorFunc, expire ExpireFunc) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if len(f.task) == 0 {
		f.newTask <- struct{}{}
	}
	f.task[id] = &task{match: match, executor: executor, expire: expire, time: time.Now()}
}

func (f *Fetcher) MatchTask(id string, message types.Message) bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	if t, ok := f.task[id]; ok {
		if t.match(message) {
			go t.executor(message)
			return true
		}
	}
	return false
}

func (f *Fetcher) loop() {
	fetchTimer := time.NewTimer(0)
	for {
		select {
		case <-f.newTask:
			fetchTimer.Reset(arriveTimeout)

		case <-fetchTimer.C:
			f.lock.Lock()
			for id, task := range f.task {
				if time.Since(task.time) > arriveTimeout {
					if task.expire != nil {
						task.expire()
					}
					delete(f.task, id)
				}
			}
			if len(f.task) == 0 {
				fetchTimer.Reset(0)
			}
			f.lock.Unlock()

		}
	}
}
