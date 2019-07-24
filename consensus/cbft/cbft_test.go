package cbft

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThreshold(t *testing.T) {
	f := &Cbft{}
	assert.Equal(t, 1, f.threshold(1))
	assert.Equal(t, 2, f.threshold(2))
	assert.Equal(t, 3, f.threshold(3))
	assert.Equal(t, 3, f.threshold(4))
	assert.Equal(t, 4, f.threshold(5))
	assert.Equal(t, 5, f.threshold(6))
	assert.Equal(t, 5, f.threshold(7))

}
