package p2p

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"strings"
	"sync"
)

type consensusDial struct {
	lock                  sync.RWMutex
	queue                 []*dialTask
	maxPeers              int
	removeConsensusPeerFn removeConsensusPeerFn
}

func NewConsensusDial(maxPeers int) *consensusDial {
	dial := &consensusDial {
		maxPeers:              maxPeers,
	}
	return dial
}

func (dial *consensusDial) InitRemoveConsensusPeerFn(removeConsensusPeerFn removeConsensusPeerFn) {
	dial.removeConsensusPeerFn = removeConsensusPeerFn
}

func (dial *consensusDial) AddTask(task *dialTask) error {
	dial.lock.Lock()
	defer dial.lock.Unlock()

	// whether the task is already in the queue
	// 1 if exists,remove task to the end of the queue;
	// 2 if not exists(not exceeding the maximum limit),add new task directly to the end of the queue
	// 3 if not exists(exceeding the maximum limit),Remove queue head task and add new task to the end of the queue
	index := dial.index(task)
	log.Info("[before add]Consensus dialed task list before AddTask operation", "task queue", dial.description())
	if index != -1 {
		log.Info("Consensus dialed task exists,Remove new task to the end of the queue", "index", index)
		dial.pollIndex(index)
		dial.offer(task)
	} else if dial.size() < dial.maxPeers {
		log.Info("Consensus dialed task not exists,Not exceeding the maximum limit,Add new task directly to the end of the queue", "tasks size", dial.size(), "maxConsensusPeers", dial.maxPeers)
		dial.offer(task)
	} else {
		log.Info("Consensus dialed task not exists,Exceeding the maximum limit，Remove queue head task and add new task to the end of the queue", "tasks size", dial.size(), "maxConsensusPeers", dial.maxPeers)
		pollTask := dial.queue[0]                 // queue head task
		dial.removeConsensusPeerFn(pollTask.dest) // disconnect head peer
		dial.queue = dial.queue[1:]               // remove queue head task
		dial.offer(task)
	}
	log.Info("[after add]Consensus dialed task list after AddTask operation", "task queue", dial.description())
	return nil
}

func (dial *consensusDial) RemoveTask(NodeID discover.NodeID) error {
	dial.lock.Lock()
	defer dial.lock.Unlock()

	log.Info("[before remove]Consensus dialed task list before RemoveTask operation", "task queue", dial.description())
	if !dial.isEmpty() {
		for i, t := range dial.queue {
			if t.dest.ID == NodeID {
				dial.queue = append(dial.queue[:i], dial.queue[i+1:]...)
				break
			}
		}
	}
	log.Info("[after remove]Consensus dialed task list after RemoveTask operation", "task queue", dial.description())
	return nil
}

func (dial *consensusDial) ListTask() []*dialTask {
	dial.lock.RLock()
	defer dial.lock.RUnlock()

	log.Info("[after list]Consensus dialed task list after ListTask operation", "task queue", dial.description())
	return dial.queue
}

// adding new task to the end of the queue
func (dial *consensusDial) offer(task *dialTask) {
	dial.queue = append(dial.queue, task)
}

// remove the first task in the queue
func (dial *consensusDial) poll() *dialTask {
	if dial.isEmpty() {
		log.Info("dialedTasks is empty!")
		return nil
	}

	pollTask := dial.queue[0]
	dial.queue = dial.queue[1:]
	return pollTask
}

// remove the specify index task in the queue
func (dial *consensusDial) pollIndex(index int) *dialTask {
	if dial.isEmpty() {
		log.Info("dialedTasks is empty!")
		return nil
	}

	pollTask := dial.queue[index]
	dial.queue = append(dial.queue[:index], dial.queue[index+1:]...)
	return pollTask
}

// index of task in the queue
func (dial *consensusDial) index(task *dialTask) int {
	for i, t := range dial.queue {
		if t.dest.ID == task.dest.ID {
			return i
		}
	}
	return -1
}

// queue size
func (dial *consensusDial) size() int {
	return len(dial.queue)
}

// clear queue
func (dial *consensusDial) clear() bool {
	if dial.isEmpty() {
		log.Info("queue is empty!")
		return false
	}
	for i := 0; i < dial.size(); i++ {
		dial.queue[i] = nil
	}
	dial.queue = nil
	return true
}

// whether the queue is empty
func (dial *consensusDial) isEmpty() bool {
	if len(dial.queue) == 0 {
		return true
	}
	return false
}

func (dial *consensusDial) description() string {
	var description []string
	for _, t := range dial.queue {
		description = append(description, fmt.Sprintf("%x", t.dest.ID[:8]))
	}
	return strings.Join(description, ",")
}
