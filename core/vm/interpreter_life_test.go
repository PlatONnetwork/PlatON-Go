package vm

import (
	"fmt"
	"testing"
)

func TestInt64ToBytes(t *testing.T) {
	var v int64 = 1000
	bytes := Int64ToBytes(v)
	fmt.Printf(string(bytes))
}