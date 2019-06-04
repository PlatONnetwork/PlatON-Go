package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"testing"
)

func BenchmarkLogInternalBP_ExecuteBlock(b *testing.B) {
	l := logBP
	for i := 0; i < 10000; i++ {
		l.InternalBP().ExecuteBlock(nil, common.BytesToHash(Rand32Bytes(32)), 40, 100, 40)
	}
}
