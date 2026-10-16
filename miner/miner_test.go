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
			running:            0,
			startCh:            make(chan struct{}),
			exitCh:             make(chan struct{}),
			resubmitIntervalCh: make(chan time.Duration),
		},
	}

	miner.wg.Add(1)
	go miner.update()

	// Drain worker.startCh so worker.start() cannot block, and wait until the
	// start signal is observed before returning — otherwise Mining() can pass
	// while start() is still sending and a later close races with that send.
	started := make(chan struct{})
	go func() {
		select {
		case <-miner.worker.startCh:
			t.Log("Start miner done")
			close(started)
		case <-time.After(2 * time.Second):
			t.Error("Start miner timeout")
			close(started)
		}
	}()

	miner.Start()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("waiting for miner start timed out")
	}

	t.Cleanup(func() {
		miner.Close()
	})
	return miner
}

func TestMiner_Start(t *testing.T) {
	miner := minerStart(t)
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
		worker: &worker{
			running: 1,
			startCh: make(chan struct{}),
			exitCh:  make(chan struct{}),
		},
	}
	miner.wg.Add(1)
	go miner.update()
	t.Cleanup(func() {
		miner.Close()
	})

	miner.Stop()
	assert.Eventually(t, func() bool {
		return miner.worker.running == 0
	}, time.Second, 10*time.Millisecond,
		fmt.Sprintf("After Stop, the worker flag `running` expect: 0, got: %d", miner.worker.running))
}

func TestMiner_Mining(t *testing.T) {
	miner := minerStart(t)
	assert.True(t, miner.Mining(), "the miner is not running")
}

func TestMiner_Close(t *testing.T) {
	cbft := consensus.NewFaker()
	miner := &Miner{
		engine:  cbft,
		mux:     new(event.TypeMux),
		exitCh:  make(chan struct{}),
		startCh: make(chan struct{}),
		stopCh:  make(chan struct{}),
		worker: &worker{
			running:            0,
			startCh:            make(chan struct{}),
			exitCh:             make(chan struct{}),
			resubmitIntervalCh: make(chan time.Duration),
		},
	}
	miner.wg.Add(1)
	go miner.update()

	started := make(chan struct{})
	go func() {
		select {
		case <-miner.worker.startCh:
			close(started)
		case <-time.After(2 * time.Second):
			t.Error("Start miner timeout")
			close(started)
		}
	}()
	miner.Start()
	<-started

	done := make(chan struct{})
	go func() {
		select {
		case <-miner.exitCh:
			close(done)
		case <-miner.worker.exitCh:
			close(done)
		case <-time.After(2 * time.Second):
			t.Error("Close miner and worker timeout")
			close(done)
		}
	}()
	miner.Close()
	<-done
}

func TestMiner_Pending(t *testing.T) {
	miner := minerStart(t)
	b, st := miner.Pending()
	assert.Nil(t, b, "the block must be nil")
	assert.Nil(t, st, "the state must be nil")
}

func TestMiner_PendingBlock(t *testing.T) {
	miner := minerStart(t)
	b := miner.PendingBlock()
	assert.Nil(t, b, "the block must be nil")
}

func TestMiner_SetRecommitInterval(t *testing.T) {
	miner := minerStart(t)
	interval := 3 * time.Second

	got := make(chan struct{})
	go func() {
		select {
		case <-miner.worker.resubmitIntervalCh:
			t.Log("receive the resubmit signal")
			close(got)
		case <-time.After(interval):
			t.Error("resubmit timeout")
			close(got)
		}
	}()

	miner.SetRecommitInterval(interval)
	<-got
}
