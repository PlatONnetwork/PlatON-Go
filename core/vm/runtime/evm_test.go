package runtime

import (
	"context"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
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

func NewEVMWithCtx(cfg *Config) *vm.EVM {
	vmenv := NewEnv(cfg)
	vmenv.Ctx = context.TODO()
	return vmenv
}

func TestEVMCallError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	cfg := new(Config)
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	}
	var (
		vmenv  = NewEVMWithCtx(cfg)
		sender = vm.AccountRef(cfg.Origin)
	)
	vmenv.Call(
		sender,
		common.BytesToAddress([]byte("contract")),
		nil,
		1,
		cfg.Value,
	)
}

func baseConfig(cfg *Config) {
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	}
	var (
		address = common.BytesToAddress([]byte("contract"))
	)
	cfg.State.CreateAccount(address)
	code := []byte{
		byte(vm.DIFFICULTY),
		byte(vm.TIMESTAMP),
		byte(vm.GASLIMIT),
		byte(vm.PUSH1),
		byte(vm.ORIGIN),
		byte(vm.BLOCKHASH),
		byte(vm.COINBASE),
	}
	// set the receiver's (the executing contract) code for execution.
	cfg.State.SetCode(address, code)
}

func TestEVMCall(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	cfg := new(Config)
	baseConfig(cfg)
	vmenv := NewEVMWithCtx(cfg)
	sender := vm.AccountRef(cfg.Origin)
	// Call the code with the given configuration.
	vmenv.Call(
		sender,
		common.BytesToAddress([]byte("contract")),
		nil,
		1,
		cfg.Value,
	)
}

func TestCallCode(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	cfg := new(Config)
	baseConfig(cfg)
	vmenv := NewEVMWithCtx(cfg)
	sender := vm.AccountRef(cfg.Origin)
	// Call the code with the given configuration.
	vmenv.CallCode(
		sender,
		common.BytesToAddress([]byte("contract")),
		nil,
		cfg.GasLimit,
		cfg.Value,
	)
}

func TestDelegateCall(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	cfg := new(Config)
	baseConfig(cfg)
	vmenv := NewEVMWithCtx(cfg)
	sender := vm.NewContract(account{}, account{}, big.NewInt(0), cfg.GasLimit)
	vmenv.DelegateCall(
		sender,
		common.BytesToAddress([]byte("contract")),
		nil,
		1,
	)
}

func TestStaticCall(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	cfg := new(Config)
	baseConfig(cfg)
	vmenv := NewEVMWithCtx(cfg)
	sender := vm.NewContract(account{}, account{}, big.NewInt(0), cfg.GasLimit)
	vmenv.StaticCall(
		sender,
		common.BytesToAddress([]byte("contract")),
		nil,
		1,
	)
}

func TestCreate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	cfg := new(Config)
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	}
	var (
		address = common.BytesToAddress([]byte("contract"))
		vmenv   = NewEVMWithCtx(cfg)
		sender  = vm.NewContract(account{}, account{}, big.NewInt(0), cfg.GasLimit)
	)
	cfg.State.CreateAccount(address)
	code := []byte{
		byte(vm.DIFFICULTY),
		byte(vm.TIMESTAMP),
		byte(vm.GASLIMIT),
		byte(vm.PUSH1),
		byte(vm.ORIGIN),
		byte(vm.BLOCKHASH),
		byte(vm.COINBASE),
	}
	// set the receiver's (the executing contract) code for execution.
	vmenv.Create(
		sender,
		code,
		cfg.GasLimit,
		cfg.Value,
	)
	vmenv.Create2(
		sender,
		code,
		cfg.GasLimit,
		cfg.Value,
		cfg.Value,
	)
	vmenv.Create(
		sender,
		code,
		1,
		cfg.Value,
	)
}

func TestPre(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	cfg := new(Config)
	baseConfig(cfg)
	vmenv := NewEVMWithCtx(cfg)
	sender := vm.AccountRef(cfg.Origin)
	// Call the code with the given configuration.
	vmenv.CallCode(
		sender,
		common.HexToAddress("0x1000000000000000000000000000000000000001"),
		nil,
		cfg.GasLimit,
		cfg.Value,
	)
	vmenv.CallCode(
		sender,
		common.HexToAddress("0x1000000000000000000000000000000000000002"),
		nil,
		cfg.GasLimit,
		cfg.Value,
	)
	vmenv.CallCode(
		sender,
		common.HexToAddress("0x1000000000000000000000000000000000000003"),
		nil,
		cfg.GasLimit,
		cfg.Value,
	)
	vmenv.CallCode(
		sender,
		common.HexToAddress("0x1000000000000000000000000000000000000004"),
		nil,
		cfg.GasLimit,
		cfg.Value,
	)
	vmenv.CallCode(
		sender,
		common.HexToAddress("0x1000000000000000000000000000000000000005"),
		nil,
		cfg.GasLimit,
		cfg.Value,
	)
	vmenv.CallCode(
		sender,
		common.HexToAddress("0x2000000000000000000000000000000000000000"),
		nil,
		cfg.GasLimit,
		cfg.Value,
	)
}

func TestOthers(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	cfg := new(Config)
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	}
	var (
		vmenv = NewEVMWithCtx(cfg)
	)

	vmenv.Interpreter()
	vmenv.GetStateDB()
	vmenv.GetEvm()
	vmenv.GetVMConfig()
	vmenv.Cancel()
}
