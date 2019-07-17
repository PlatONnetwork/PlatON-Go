package executor

import (
	"errors"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type BlockExecutor interface {
	//Execution block, you need to pass in the parent block to find the parent block state
	execute(block *types.Block, parent *types.Block) error
}

//Block execution results, including block hash, block number, error message
type BlockExecuteStatus struct {
	hash   common.Hash
	number uint64
	err    error
}

type AsyncBlockExecutor interface {
	BlockExecutor
	//Asynchronous acquisition block execution results
	executeStatus() chan<- BlockExecuteStatus
}

type executeTask struct {
	parent *types.Block
	block  *types.Block
}

// asyncExecutor async block executor implement.
type asyncExecutor struct {
	AsyncBlockExecutor

	// executeFn is a function use to execute block.
	executeFn consensus.Executor

	executeTasks   chan *executeTask       // A channel for notify execute task.
	executeResults chan BlockExecuteStatus // A channel for notify execute result.

	// A channel for notify stop signal
	closed chan struct{}
}

// NewAsyncExecutor new a async block executor.
func NewAsyncExecutor(executeFn consensus.Executor) *asyncExecutor {
	exe := &asyncExecutor{
		executeFn:      executeFn,
		executeTasks:   make(chan *executeTask, 64),
		executeResults: make(chan BlockExecuteStatus, 64),
		closed:         make(chan struct{}),
	}

	go exe.loop()

	return exe
}

// stop stop async exector.
func (exe *asyncExecutor) stop() {
	close(exe.closed)
}

// execute async execute block.
func (exe *asyncExecutor) execute(block *types.Block, parent *types.Block) error {
	return exe.newTask(block, parent)
}

// executeStatus return a channel for notify block execute result.
func (exe *asyncExecutor) executeStatus() chan<- BlockExecuteStatus {
	return exe.executeResults
}

// newTask new a block execute task and push in execute channel.
// If execute channel if full, will return a error.
func (exe *asyncExecutor) newTask(block *types.Block, parent *types.Block) error {
	select {
	case exe.executeTasks <- &executeTask{parent: parent, block: block}:
		return nil
	default:
		// FIXME: blocking if channel is full?
		return errors.New("execute task queue is full")
	}
}

// loop process task from execute channel until executor stopped.
func (exe *asyncExecutor) loop() {
	for {
		select {
		case <-exe.closed:
			return
		case task := <-exe.executeTasks:
			err := exe.executeFn(task.block, task.parent)
			exe.executeResults <- BlockExecuteStatus{
				hash:   task.block.Hash(),
				number: task.block.Number().Uint64(),
				err:    err,
			}
		}
	}
}
