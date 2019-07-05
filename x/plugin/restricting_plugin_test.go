package plugin_test

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

func TestRestrictingPlugin_EndBlock(t *testing.T) {
	db := ethdb.NewMemDatabase()
	stateDb, err := state.New(common.Hash{}, state.NewDatabase(db))

	if err != nil {
		t.Errorf(err.Error())
	}

	head := types.Header{ Number: big.NewInt(1),}

	if err := plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb); err != nil {
		t.Error(err)
	}
}

func TestRestrictingPlugin_AddRestrictingRecord(t *testing.T) {

	sender := common.HexToAddress("0x11")

	if err := plugin.RestrictingInstance().AddRestrictingRecord(); err != nil {

	}

}

func TestRestrictingPlugin_PledgeLockFunds(t *testing.T) {

}


func TestRestrictingPlugin_ReturnLockFunds(t *testing.T) {

}

func TestRestrictingPlugin_SlashingNotify(t *testing.T) {

}

func TestRestrictingPlugin_GetRestrictingInfo(t *testing.T) {

}