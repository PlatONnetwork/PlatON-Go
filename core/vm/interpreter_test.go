package vm

import (
	"context"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

type account struct{}

func (account) SubBalance(amount *big.Int)                                 {}
func (account) AddBalance(amount *big.Int)                                 {}
func (account) SetAddress(common.Address)                                  {}
func (account) Value() *big.Int                                            { return nil }
func (account) SetBalance(*big.Int)                                        {}
func (account) SetNonce(uint64)                                            {}
func (account) Balance() *big.Int                                          { return nil }
func (account) Address() common.Address                                    { return common.Address{} }
func (account) ReturnGas(*big.Int)                                         {}
func (account) SetCode(common.Hash, []byte)                                {}
func (account) ForEachStorage(cb func(key common.Hash, value []byte) bool) {}

func TestEnforceRestrictions(t *testing.T) {
	var (
		env            = NewEVM(Context{}, nil, &mock.MockStateDB{}, params.TestChainConfig, Config{})
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)
	evmInterpreter.readOnly = true
	varstore := operation{
		execute:       opSstore,
		gasCost:       gasSStore,
		validateStack: makeStackFunc(2, 0),
		valid:         true,
		writes:        true,
	}
	err := evmInterpreter.enforceRestrictions(SSTORE, varstore, nil)
	if err == nil {
		t.Errorf("Test enforceRestrictions error")
	}
}

func TestRun(t *testing.T) {
	var (
		env            = NewEVM(Context{Ctx: context.TODO()}, nil, &mock.MockStateDB{}, params.TestChainConfig, Config{})
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)
	contract := NewContract(account{}, account{}, big.NewInt(0), 1)
	_, err := evmInterpreter.Run(contract, []byte{}, true)
	if err != nil {
		t.Errorf("Test Run error")
	}

	contract.Code = []byte{byte(0xfe), 0x1, byte(PUSH1), 0x1, 0x0}
	_, err = evmInterpreter.Run(contract, []byte{}, true)
	if err == nil {
		t.Errorf("Test Run error")
	}

	contract.Code = []byte{byte(PUSH1), 0x1, byte(PUSH1), 0x1, byte(SSTORE), 0x1, 0x2}
	_, err = evmInterpreter.Run(contract, []byte{}, true)
	if err == nil {
		t.Errorf("Test Run error")
	}

	contract.Code = []byte{byte(PUSH1), 0x1, byte(PUSH1), 0x1, byte(SSTORE), 0x1, 0x2}
	contract.Gas = 100000
	_, err = evmInterpreter.Run(contract, []byte{}, true)
	if err == nil {
		t.Errorf("Test Run error")
	}

	contract.Code = []byte{}
	for i := 0; i <= 1024; i++ {
		contract.Code = append(contract.Code, byte(PUSH1), 0x1)
	}
	contract.Gas = 100000
	evmInterpreter.cfg.JumpTable = constantinopleInstructionSet
	_, err = evmInterpreter.Run(contract, []byte{}, true)
	if err == nil {
		t.Errorf("Test Run error")
	}

	contract.Code = []byte{byte(PUSH1), 0x1, byte(PUSH1), 0x1, byte(REVERT)}
	contract.Gas = 100000
	_, err = evmInterpreter.Run(contract, []byte{}, true)
	if err == nil {
		t.Errorf("Test Run error")
	}
}
