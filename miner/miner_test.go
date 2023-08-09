package miner

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/stretchr/testify/assert"
)

func minerStart(t *testing.T) *Miner {
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

	go func() {
		select {
		case <-miner.worker.startCh:
			t.Log("Start miner done")
		case <-time.After(2 * time.Second):
			t.Fatal("Start miner timeout")
		}
	}()

	miner.Start()
	return miner
}

func TestMiner_Start(t *testing.T) {

	miner := minerStart(t)

	assert.True(t, miner.Mining())
	close(miner.worker.startCh)

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
			//startCh: make(chan struct{}),
		},
	}
	go miner.update()

	miner.Stop()
	assert.Equal(t, int32(0), miner.worker.running,
		fmt.Sprintf("After Stop, the worker flag `running` expect: %d, got: %d", int32(0), miner.worker.running))
}

func TestMiner_Mining(t *testing.T) {
	miner := minerStart(t)
	assert.True(t, miner.Mining(), "the miner is not running")
}

func TestMiner_Close(t *testing.T) {
	miner := minerStart(t)

	go func() {
		select {
		case <-miner.exitCh:

		case <-miner.worker.exitCh:

		case <-time.After(2 * time.Second):
			t.Fatal("Close miner and worker timeout")

		}
	}()

	miner.Close()

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

	go func() {
		select {
		case <-miner.worker.resubmitIntervalCh:
			t.Log("receive the resubmit signal")
		case <-time.After(interval):
			t.Fatal("resubmit timeout")
		}
	}()

	miner.SetRecommitInterval(interval)

}
