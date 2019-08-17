package executor

import (
	"errors"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type Executor func(block *types.Block, parent *types.Block) error

type BlockExecutor interface {
	//Execution block, you need to pass in the parent block to find the parent block state
	Execute(block *types.Block, parent *types.Block) error
}

//Block execution results, including block Hash, block Number, error message
type BlockExecuteStatus struct {
	Hash   common.Hash
	Number uint64
	Err    error
}

type AsyncBlockExecutor interface {
	BlockExecutor
	Stop()
	//Asynchronous acquisition block execution results
	ExecuteStatus() <-chan *BlockExecuteStatus
}

type executeTask struct {
	parent *types.Block
	block  *types.Block
}

// asyncExecutor async block executor implement.
type AsyncExecutor struct {
	AsyncBlockExecutor

	// executeFn is a function use to execute block.
	executeFn Executor

	executeTasks   chan *executeTask        // A channel for notify execute task.
	executeResults chan *BlockExecuteStatus // A channel for notify execute result.

	// A channel for notify stop signal
	closed chan struct{}
}

// NewAsyncExecutor new a async block executor.
func NewAsyncExecutor(executeFn Executor) *AsyncExecutor {
	exe := &AsyncExecutor{
		executeFn:      executeFn,
		executeTasks:   make(chan *executeTask, 64),
		executeResults: make(chan *BlockExecuteStatus, 64),
		closed:         make(chan struct{}),
	}

	go exe.loop()

	return exe
}

// Stop stop async exector.
func (exe *AsyncExecutor) Stop() {
	close(exe.closed)
}

// Execute async execute block.
func (exe *AsyncExecutor) Execute(block *types.Block, parent *types.Block) error {
	return exe.newTask(block, parent)
}

// executeStatus return a channel for notify block execute result.
func (exe *AsyncExecutor) ExecuteStatus() <-chan *BlockExecuteStatus {
	return exe.executeResults
}

// newTask new a block execute task and push in execute channel.
// If execute channel if full, will return a error.
func (exe *AsyncExecutor) newTask(block *types.Block, parent *types.Block) error {
	select {
	case exe.executeTasks <- &executeTask{parent: parent, block: block}:
		return nil
	default:
		// FIXME: blocking if channel is full?
		return errors.New("execute task queue is full")
	}
}

// loop process task from execute channel until executor stopped.
func (exe *AsyncExecutor) loop() {
	for {
		select {
		case <-exe.closed:
			return
		case task := <-exe.executeTasks:
			err := exe.executeFn(task.block, task.parent)
			exe.executeResults <- &BlockExecuteStatus{
				Hash:   task.block.Hash(),
				Number: task.block.Number().Uint64(),
				Err:    err,
			}
		}
	}
}
