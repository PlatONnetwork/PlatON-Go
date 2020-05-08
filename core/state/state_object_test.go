package state

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"testing"
)

func TestStateObject(t *testing.T) {
	x := Account{
		Root: common.HexToHash("0x1000000000000000000000000000000000000001"),
	}
	x2 := newObject(nil, common.HexToAddress("0x1000000000000000000000000000000000000001"), x)
	x3 := x2.deepCopy(nil)
	x2.data.Root = common.HexToHash("0x1000000000000000000000000000000000000012")
	t.Log(x2.data.Root.String())
	t.Log(x3.data.Root.String())
}
