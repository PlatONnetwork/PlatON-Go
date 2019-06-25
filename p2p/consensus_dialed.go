package p2p

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"strings"
)

type dialedTasks struct {
	queue                 []*dialTask
	maxPeers              int
	removeConsensusPeerFn removeConsensusPeerFn
}

func NewDialedTasks(maxPeers int, removeConsensusPeerFn removeConsensusPeerFn) *dialedTasks {
	tasks := &dialedTasks{
		maxPeers:              maxPeers,
		removeConsensusPeerFn: removeConsensusPeerFn,
	}
	return tasks
}

func (tasks *dialedTasks) InitRemoveConsensusPeerFn(removeConsensusPeerFn removeConsensusPeerFn) {
	tasks.removeConsensusPeerFn = removeConsensusPeerFn
}

func (tasks *dialedTasks) AddTask(task *dialTask) error {

	// whether the task is already in the queue
	// 1 if exists,remove task to the end of the queue;
	// 2 if not exists(not exceeding the maximum limit),add new task directly to the end of the queue
	// 3 if not exists(exceeding the maximum limit),Remove queue head task and add new task to the end of the queue
	index := tasks.index(task)
	log.Info("[before add]Consensus dialed task list before AddTask operation", "task queue", tasks.description())
	if index != -1 {
		log.Info("Consensus dialed task exists,Remove new task to the end of the queue", "index", index)
		tasks.pollIndex(index)
		tasks.offer(task)
	} else if tasks.size() < tasks.maxPeers {
		log.Info("Consensus dialed task not exists,Not exceeding the maximum limit,Add new task directly to the end of the queue", "tasks size", tasks.size(), "maxConsensusPeers", tasks.maxPeers)
		tasks.offer(task)
	} else {
		log.Info("Consensus dialed task not exists,Exceeding the maximum limit，Remove queue head task and add new task to the end of the queue", "tasks size", tasks.size(), "maxConsensusPeers", tasks.maxPeers)
		pollTask := tasks.queue[0]                 // queue head task
		tasks.removeConsensusPeerFn(pollTask.dest) // disconnect head peer
		tasks.queue = tasks.queue[1:]              // remove queue head task
		tasks.offer(task)
	}
	log.Info("[after add]Consensus dialed task list after AddTask operation", "task queue", tasks.description())
	return nil
}

func (tasks *dialedTasks) RemoveTask(NodeID discover.NodeID) error {

	log.Info("[before remove]Consensus dialed task list before RemoveTask operation", "task queue", tasks.description())
	if !tasks.isEmpty() {
		for i, t := range tasks.queue {
			if t.dest.ID == NodeID {
				tasks.queue = append(tasks.queue[:i], tasks.queue[i+1:]...)
				break
			}
		}
	}
	log.Info("[after remove]Consensus dialed task list after RemoveTask operation", "task queue", tasks.description())
	return nil
}

func (tasks *dialedTasks) ListTask() []*dialTask {

	log.Info("[after list]Consensus dialed task list after ListTask operation", "task queue", tasks.description())
	return tasks.queue
}

// adding new task to the end of the queue
func (tasks *dialedTasks) offer(task *dialTask) {
	tasks.queue = append(tasks.queue, task)
}

// remove the first task in the queue
func (tasks *dialedTasks) poll() *dialTask {
	if tasks.isEmpty() {
		log.Info("dialedTasks is empty!")
		return nil
	}

	pollTask := tasks.queue[0]
	tasks.queue = tasks.queue[1:]
	return pollTask
}

// remove the specify index task in the queue
func (tasks *dialedTasks) pollIndex(index int) *dialTask {
	if tasks.isEmpty() {
		log.Info("dialedTasks is empty!")
		return nil
	}

	pollTask := tasks.queue[index]
	tasks.queue = append(tasks.queue[:index], tasks.queue[index+1:]...)
	return pollTask
}

// index of task in the queue
func (tasks *dialedTasks) index(task *dialTask) int {
	for i, t := range tasks.queue {
		if t.dest.ID == task.dest.ID {
			return i
		}
	}
	return -1
}

// queue size
func (tasks *dialedTasks) size() int {
	return len(tasks.queue)
}

// clear queue
func (tasks *dialedTasks) clear() bool {
	if tasks.isEmpty() {
		log.Info("queue is empty!")
		return false
	}
	for i := 0; i < tasks.size(); i++ {
		tasks.queue[i] = nil
	}
	tasks.queue = nil
	return true
}

// whether the queue is empty
func (tasks *dialedTasks) isEmpty() bool {
	if len(tasks.queue) == 0 {
		return true
	}
	return false
}

func (tasks *dialedTasks) description() string {
	var description []string
	for _, t := range tasks.queue {
		description = append(description, fmt.Sprintf("%x", t.dest.ID[:8]))
	}
	return strings.Join(description, ",")
}
