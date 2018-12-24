package p2p

import (
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"fmt"
	"strings"
	"sync"
)

type dialedTasks struct {
	lock                  sync.RWMutex
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
	tasks.lock.Lock()
	defer tasks.lock.Unlock()

	// 判断任务是否已经在队列中
	// 1 如果存在,将任务重新移到队尾;
	// 2 不存在（未超过最大限制数）,将新任务直接添加到队尾
	// 3 不存在（超过最大限制数）,移除队头任务并将新任务添加到队尾
	index := tasks.index(task)
	log.Info("AddTask操作前consensus任务列表", "列表", tasks.description())
	if index != -1 {
		log.Info("consensus任务存在,将新任务重新移到队尾", "index", index)
		tasks.pollIndex(index)
		tasks.offer(task)

	} else if tasks.size() < tasks.maxPeers {
		log.Info("consensus任务不存在,未超过最大限制数，将新任务直接添加到队尾", "tasks size", tasks.size(), "maxConsensusPeers", tasks.maxPeers)
		tasks.offer(task)
	} else {
		log.Info("consensus任务不存在,超过最大限制数，移除队头任务并将新任务添加到队尾", "tasks size", tasks.size(), "maxConsensusPeers", tasks.maxPeers)
		pollTask := tasks.queue[0]                 // 队头任务
		tasks.removeConsensusPeerFn(pollTask.dest) // 断开连接
		tasks.queue = tasks.queue[1:]              // 移除队头任务
		tasks.offer(task)
	}
	log.Info("AddTask操作后consensus任务列表", "列表", tasks.description())
	return nil
}

func (tasks *dialedTasks) RemoveTask(NodeID discover.NodeID) error {
	tasks.lock.Lock()
	defer tasks.lock.Unlock()

	log.Info("RemoveTask操作前consensus任务列表", "列表", tasks.description())
	if !tasks.isEmpty() {
		for i, t := range tasks.queue {
			if t.dest.ID == NodeID {
				tasks.queue = append(tasks.queue[:i], tasks.queue[i+1:]...)
				break
			}
		}
	}
	log.Info("RemoveTask操作后consensus任务列表", "列表", tasks.description())
	return nil
}

func (tasks *dialedTasks) ListTask() []*dialTask {
	tasks.lock.RLock()
	defer tasks.lock.RUnlock()

	log.Info("ListTask操作后consensus任务列表", "列表", tasks.description())
	return tasks.queue
}

// 向队列尾部添加新任务
func (tasks *dialedTasks) offer(task *dialTask) {
	tasks.queue = append(tasks.queue, task)
}

// 移除队列中最前面的任务
func (tasks *dialedTasks) poll() *dialTask {
	if tasks.isEmpty() {
		log.Info("dialedTasks is empty!")
		return nil
	}

	pollTask := tasks.queue[0]
	tasks.queue = tasks.queue[1:]
	return pollTask
}

// 移除队列中指定下标位置的任务
func (tasks *dialedTasks) pollIndex(index int) *dialTask {
	if tasks.isEmpty() {
		log.Info("dialedTasks is empty!")
		return nil
	}

	pollTask := tasks.queue[index]
	tasks.queue = append(tasks.queue[:index], tasks.queue[index+1:]...)
	return pollTask
}

// 判断任务在队列中的位置
func (tasks *dialedTasks) index(task *dialTask) int {
	for i, t := range tasks.queue {
		if t.dest.ID == task.dest.ID {
			return i
		}
	}
	return -1
}

// 返回队列长度
func (tasks *dialedTasks) size() int {
	return len(tasks.queue)
}

// 清空队列
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

// 判断队列是否为空
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
