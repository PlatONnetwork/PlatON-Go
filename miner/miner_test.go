package miner

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
)

func minerStart(t *testing.T) *Miner {
	cbft := cbft.NewFaker()

	miner := &Miner{
		engine:   cbft,
		exitCh:   make(chan struct{}),
		canStart: 1,
		worker: &worker{
			running:            0,
			startCh:            make(chan struct{}),
			exitCh:             make(chan struct{}),
			resubmitIntervalCh: make(chan time.Duration),
		},
	}

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

	assert.Equal(t, int32(1), miner.shouldStart,
		fmt.Sprintf("After Start, the miner flag `shouldStart` expect: %d, got: %d", int32(1), miner.shouldStart))
	close(miner.worker.startCh)

}

func TestMiner_Stop(t *testing.T) {
	cbft := cbft.NewFaker()

	miner := &Miner{
		engine:      cbft,
		exitCh:      make(chan struct{}),
		shouldStart: 1,
		worker: &worker{
			running: 1,
			//startCh: make(chan struct{}),
		},
	}

	miner.Stop()
	assert.Equal(t, int32(0), miner.shouldStart,
		fmt.Sprintf("After Stop, the miner flag `shouldStart` expect: %d, got: %d", int32(0), miner.shouldStart))
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
