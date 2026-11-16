package plugin

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
)

func TestSetForkBlockAfterLaterFork(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := gov.NewGovDB(chain.SnapDB)
	gp := &GovPlugin{gov: gov.NewGov(chain.SnapDB), gdb: gdb}
	pauliBlock := func() *big.Int {
		var got *big.Int
		if err := gp.setForkBlock(common.Hash{}, chain.StateDB, params.FORKVERSION_1_5_0, func(b *big.Int) { got = b }); err != nil {
			t.Fatal(err)
		}
		return got
	}

	if err := gdb.AddActiveVersion(params.FORKVERSION_1_5_0, 500, chain.StateDB); err != nil {
		t.Fatal(err)
	}
	if got := pauliBlock(); got == nil || got.Uint64() != 500 {
		t.Fatalf("1.5.0 current: PauliBlock = %v, want 500", got)
	}

	if err := gdb.AddActiveVersion(params.FORKVERSION_1_6_0, 900, chain.StateDB); err != nil {
		t.Fatal(err)
	}
	if got := pauliBlock(); got == nil || got.Uint64() != 500 {
		t.Fatalf("1.6.0 current: PauliBlock = %v, want 500", got)
	}
}
