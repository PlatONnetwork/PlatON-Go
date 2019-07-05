package plugin_test

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)


type restrictingInfo struct {
	balance     *big.Int `json:"balance"` // balance representation all locked amount
	debt        *big.Int `json:"debt"`    // debt representation will released amount. Positive numbers can be used instead of release, 0 means no release, negative numbers indicate not enough to release
	releaseList []uint64 `json:"list"`    // releaseList representation which epoch will release restricting
}

type releaseAmountInfo struct {
	height uint64 	 `json:"blockNumber"`  	// blockNumber representation of the block number at the released epoch
	amount *big.Int	 `json:"amount"`		// amount representation of the released amount
}


func buildStateDB (t *testing.T) xcom.StateDB{
	db := ethdb.NewMemDatabase()
	stateDb, err := state.New(common.Hash{}, state.NewDatabase(db))

	if err != nil {
		t.Errorf("new state db failed: %s", err.Error())
	}

	return stateDb
}

func buildDbRestrictingPlan(t *testing.T) {
	stateDb := buildStateDB(t)
	account := common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80aadd24")

	const Epochs = 5
	var list = make([]uint64, Epochs)

	for epoch := 0; epoch < Epochs; epoch++ {
		// build release account record
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(epoch), 1)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

		// build release amount record
		releaseAmount := big.NewInt(10000000)
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(epoch), account)
		stateDb.SetState(account, releaseAmountKey, releaseAmount.Bytes())

		// build release epoch record
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(epoch))
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint64ToBytes(1))

		list = append(list, uint64(epoch))
	}

	// build restricting user info
	var user restrictingInfo
	user.balance = big.NewInt(50000000)
	user.debt = big.NewInt(0)
	user.releaseList = list

	bUser, err := rlp.EncodeToBytes(user)
	if err != nil {
		t.Errorf("failed to rlp encode restricting info: %s", err.Error())
	}

	restrictingKey := restricting.GetRestrictingKey(account)
	stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bUser)

	stateDb.AddBalance(vm.RestrictingContractAddr, big.NewInt(50000000))
}


func TestRestrictingPlugin_EndBlock(t *testing.T) {
	stateDB := buildStateDB(t)
	head := types.Header{ Number: big.NewInt(1),}

	// case1: blockNumber not arrived settle block height

	// case2: blockNumber arrived settle block height, restricting plan not exist

	// case3: blockNumber arrived settle block height, restricting plan exist


	if err := plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDB); err != nil {
		t.Errorf("endblock returns err: %s", err.Error())
	}
}

func TestRestrictingPlugin_AddRestrictingRecord(t *testing.T) {
	sender := common.HexToAddress("0x11")

	// case1: account is new user to restricting

	// case2: restricting account exist, but restricting epoch not intersect

	// case3: restricting account exist, and restricting epoch intersect


	if err := plugin.RestrictingInstance().AddRestrictingRecord(); err != nil {

	}

}

func TestRestrictingPlugin_PledgeLockFunds(t *testing.T) {



	// case1: restricting account not exist

	// case2: restricting account exist, but balance not enough

	// case3: restricting account exist, and balance is enough

}


func TestRestrictingPlugin_ReturnLockFunds(t *testing.T) {

	// case1: restricting account not exist


	// case2: restricting account exist


}

func TestRestrictingPlugin_SlashingNotify(t *testing.T) {

	// case1: restricting account not exist

	// case2: restricting account exist

}

func TestRestrictingPlugin_GetRestrictingInfo(t *testing.T) {
	// case1: restricting account not exist
	account := common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80aadd24")
	stateDb := buildStateDB(t)

	// case2: restricting account exist



	if result, err := plugin.RestrictingInstance().GetRestrictingInfo(account, stateDb); err != nil {
		t.Errorf("GetRestrictingInfo returns err: %s", err.Error())
	} else {
		t.Log(string(result))
	}

}