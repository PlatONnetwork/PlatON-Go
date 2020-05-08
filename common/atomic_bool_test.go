package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicBool(t *testing.T) {
	var ab AtomicBool
	assert.False(t, ab.IsSet(), "Expected value to be false")

	ab.Set(true)
	assert.False(t, ab.value != 1, "Set(true) did not set value to 1")
	assert.True(t, ab.IsSet(), "Expected value to be true")

	ab.Set(true)
	assert.True(t, ab.IsSet(), "Expected value to be true")

	ab.Set(false)
	assert.False(t, ab.value != 0, "Set(false) did not set value to 0")
	assert.False(t, ab.IsSet(), "Expected value to be false")

	ab.Set(false)
	assert.False(t, ab.IsSet(), "Expected value to be false")
	assert.False(t, ab.TrySet(false), "Expected TrySet(false) to fail")

	assert.True(t, ab.TrySet(true), "Exepected TrySet(true) to succed")
	assert.True(t, ab.IsSet(), "Expected value to be true")

	ab.Set(true)
	assert.True(t, ab.IsSet(), "Expected value to be true")
	assert.False(t, ab.TrySet(true), "Expected TrySet(true) to fail")
	assert.True(t, ab.TrySet(false), "Exptected TrySet(false) to succeed")
	assert.False(t, ab.IsSet(), "Expected value to be false")
	ab._noCopy.Lock()
}
