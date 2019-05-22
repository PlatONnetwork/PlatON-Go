package cbft

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeOrderViewChange_Add(t *testing.T) {
	var p TimeOrderViewChange
	p.Add(nil)
	assert.Len(t, p, 1)
	assert.Nil(t, p)
}
