package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := time.NewTimer(0)
	timer.Stop()

	viewTimer := &viewTimer{timer: timer, timeInterval: viewTimeInterval{baseMs: uint64(1 * time.Second), exponentBase: exponentBase, maxExponent: maxExponent}}
	viewTimer.setupTimer(1)
	assert.False(t, viewTimer.isDeadline())
	select {
	case <-viewTimer.timerChan():
		assert.True(t, viewTimer.isDeadline())
	}
}

func TestInterval(t *testing.T) {
	in := viewTimeInterval{
		uint64(10), 2, 2,
	}

	assert.Equal(t, uint64(10), uint64(in.getViewTimeInterval(1)))
	assert.Equal(t, uint64(20), uint64(in.getViewTimeInterval(2)))

	in = viewTimeInterval{
		uint64(10 * time.Second), 1.5, 2,
	}

	assert.Equal(t, uint64(10*time.Second), uint64(in.getViewTimeInterval(1)))
	assert.Equal(t, uint64(15*time.Second), uint64(in.getViewTimeInterval(2)))
	assert.Equal(t, uint64(22*time.Second+500*time.Millisecond), uint64(in.getViewTimeInterval(3)))

}
