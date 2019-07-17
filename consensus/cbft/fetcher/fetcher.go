package fetcher

import "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"

type MatchFunc func(types.Message) bool
type ExecutorFunc func(types.Message)
type task struct {
	match    MatchFunc
	executor ExecutorFunc
}

type Fetcher struct {
	task map[string]*task
}

func (f *Fetcher) AddTask(id string, match MatchFunc, executor ExecutorFunc) {
	f.task[id] = &task{match: match, executor: executor}
}

func (f *Fetcher) MatchTask(id string, message types.Message) bool {
	if t, ok := f.task[id]; ok {
		if t.match(message) {
			go t.executor(message)
			return true
		}
	}
	return false
}
