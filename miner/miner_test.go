package miner

import (
	"fmt"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/event"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/consensus"
)

func minerStart(t *testing.T) *Miner {
	t.Helper()

	cbft := consensus.NewFaker()

	miner := &Miner{
		engine:  cbft,
		mux:     new(event.TypeMux),
		exitCh:  make(chan struct{}),
		startCh: make(chan struct{}),
		stopCh:  make(chan struct{}),
		worker: &worker{
			startCh:            make(chan struct{}, 1),
			exitCh:             make(chan struct{}),
			resubmitIntervalCh: make(chan time.Duration),
		},
	}

	miner.wg.Add(1)
	go miner.update()

	miner.Start()
	return miner
}

func TestMiner_Start(t *testing.T) {
	miner := minerStart(t)
	defer miner.Close()

	deadline := time.Now().Add(2 * time.Second)
	for !miner.Mining() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	assert.True(t, miner.Mining())
}

func TestMiner_Stop(t *testing.T) {
	cbft := consensus.NewFaker()

	miner := &Miner{
		mux:     new(event.TypeMux),
		engine:  cbft,
		exitCh:  make(chan struct{}),
		startCh: make(chan struct{}),
		stopCh:  make(chan struct{}),
		worker:  &worker{
			//startCh: make(chan struct{}),
		},
	}
	miner.worker.running.Store(true)
	go miner.update()

	miner.Stop()
	assert.False(t, miner.worker.running.Load(),
		fmt.Sprintf("After Stop, the worker flag `running` expect false, got %v", miner.worker.running.Load()))
}

func TestMiner_Mining(t *testing.T) {
	miner := minerStart(t)
	defer miner.Close()
	assert.True(t, miner.Mining(), "the miner is not running")
}

func TestMiner_Close(t *testing.T) {
	miner := minerStart(t)

	go func() {
		select {
		case <-miner.exitCh:

		case <-miner.worker.exitCh:

		case <-time.After(2 * time.Second):
			t.Error("Close miner and worker timeout")
		}
	}()
	miner.Close()
}

func TestMiner_Pending(t *testing.T) {
	miner := minerStart(t)
	defer miner.Close()
	b, st := miner.Pending()
	assert.Nil(t, b, "the block must be nil")
	assert.Nil(t, st, "the state must be nil")
}

func TestMiner_PendingBlock(t *testing.T) {
	miner := minerStart(t)
	defer miner.Close()
	b := miner.PendingBlock()
	assert.Nil(t, b, "the block must be nil")
}

func TestMiner_SetRecommitInterval(t *testing.T) {
	miner := minerStart(t)
	defer miner.Close()
	interval := 3 * time.Second

	go func() {
		select {
		case <-miner.worker.resubmitIntervalCh:
			t.Log("receive the resubmit signal")
		case <-time.After(interval):
			t.Error("resubmit timeout")
		}
	}()

	miner.SetRecommitInterval(interval)
}
