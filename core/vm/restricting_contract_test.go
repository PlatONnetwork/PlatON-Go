package vm_test

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

type ResultTest struct {
	balance *big.Int
	slash   *big.Int
	staking *big.Int
	debt    *big.Int
	entry   []byte
}

func TestRestrictingContract_createRestrictingPlan(t *testing.T) {
	contract := &vm.RestrictingContract{
		Plugin:  plugin.RestrictingInstance(),
		Contract: newContract(common.Big0),
		Evm: newEvm(blockNumber, blockHash, nil),
	}

	var plan  byteutil.RestrictingPlan
	var plans = make([]byteutil.RestrictingPlan, 5)

	for i := 0; i < 5; i++ {
		plan.Epoch = uint64(i+1)
		plan.Amount = big.NewInt(10000000)
		plans = append(plans, plan)
	}

	var params [][]byte
	// function_type
	param0, _ := rlp.EncodeToBytes(common.Uint32ToBytes(4000))
	// account release to
	param1 := addrArr[0].Bytes()
	// restricting plan
	param2, _ := rlp.EncodeToBytes(plans)

	params = append(params, param0)
	params = append(params, param1)
	params = append(params, param2)

	input, err := rlp.EncodeToBytes(params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("fail to rlp encode restricting input")
	} else {
		fmt.Println("rlp encode restricting input: ", hexutil.Encode(input))
	}

	if result, err := contract.Run(input); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(string(result))
	}
}


func TestRestrictingContract_getRestrictingInfo(t *testing.T) {
	contract := &vm.RestrictingContract{
		Plugin:  plugin.RestrictingInstance(),
		Contract: newContract(common.Big0),
		Evm: newEvm(blockNumber, blockHash, nil),
	}

	var params [][]byte
	param0, _ := rlp.EncodeToBytes(common.Uint32ToBytes(4100))
	param1 := addrArr[0].Bytes()

	params = append(params, param0)
	params = append(params, param1)

	input, err := rlp.EncodeToBytes(params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("fail to rlp encode restricting input")
	} else {
		fmt.Println("rlp encode restricting input: ", hexutil.Encode(input))
	}

	if result, err := contract.Run(input); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(string(result))

		var res ResultTest
		if err = rlp.Decode(bytes.NewBuffer(result), &res); err != nil {
			t.Log(res.balance)
			t.Log(res.staking)
			t.Log(res.slash)
			t.Log(res.debt)
			t.Log(res.entry)
		}
	}
}