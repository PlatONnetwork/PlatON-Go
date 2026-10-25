// Copyright 2021 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
)

// pdPrecompileInput builds the restricting precompile calldata:
// [TxCreateRestrictingPlan, account, plans].
func pdPrecompileInput(t *testing.T, account common.Address, plans []restricting.RestrictingPlan) []byte {
	t.Helper()
	p0, err := rlp.EncodeToBytes(common.Uint16ToBytes(uint16(TxCreateRestrictingPlan)))
	if err != nil {
		t.Fatal(err)
	}
	p1, err := rlp.EncodeToBytes(account.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	p2, err := rlp.EncodeToBytes(plans)
	if err != nil {
		t.Fatal(err)
	}
	in, err := rlp.EncodeToBytes([][]byte{p0, p1, p2})
	if err != nil {
		t.Fatal(err)
	}
	if len(in) >= 256 {
		t.Fatalf("precompile input must fit PUSH1, got %d bytes", len(in))
	}
	return in
}

// pdDelegateCallCode returns runtime bytecode that CODECOPYs the embedded calldata
// and DELEGATECALLs it to target. Header is 40 bytes.
func pdDelegateCallCode(t *testing.T, target common.Address, calldata []byte) []byte {
	t.Helper()
	l := len(calldata)
	code := []byte{
		0x60, byte(l), // PUSH1 l
		0x60, 0x28, // PUSH1 40 (code offset of the calldata)
		0x60, 0x00, // PUSH1 0
		0x39, // CODECOPY
		0x60, 0x00, // PUSH1 0 (retSize)
		0x60, 0x00, // PUSH1 0 (retOffset)
		0x60, byte(l), // PUSH1 l (inSize)
		0x60, 0x00, // PUSH1 0 (inOffset)
		0x73, // PUSH20
	}
	code = append(code, target.Bytes()...)
	code = append(code,
		0x5a, // GAS
		0xf4, // DELEGATECALL
		0x50, // POP
		0x00, // STOP
	)
	if len(code) != 40 {
		t.Fatalf("delegatecall header must be 40 bytes, got %d", len(code))
	}
	return append(code, calldata...)
}

// pdEVM builds an EVM over the mock chain with the restricting minimum amount
// cleared so the plan amount is the only variable.
func pdEVM(t *testing.T, origin common.Address, chain *mock.Chain) *EVM {
	t.Helper()
	sdb := snapshotdb.Instance()
	gov.InitGenesisGovernParam(common.ZeroHash, sdb, 2048)
	gov.NewGovDB(sdb).AddActiveVersion(initProgramVersion, 0, chain.StateDB)
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{Value: new(big.Int).SetInt64(0).String()})
	if err := sdb.PutBaseDB(key, value); err != nil {
		t.Fatal(err)
	}
	blockCtx := BlockContext{
		CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
		GetHash:     func(uint64) common.Hash { return blockHash },
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
		Difficulty:  new(big.Int),
	}
	return NewEVM(blockCtx, TxContext{Origin: origin, GasPrice: big.NewInt(0)}, sdb, chain.StateDB, &params.ChainConfig{}, Config{})
}

// TestSystemPrecompileRejectsDelegateCall verifies that a contract cannot make a
// system precompile act on another account's behalf by delegatecalling it. The
// delegatecall must fail and must not move any funds.
func TestSystemPrecompileRejectsDelegateCall(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	victim := common.HexToAddress("0x1111111111111111111111111111111111111111")
	beneficiary := common.HexToAddress("0x2222222222222222222222222222222222222222")
	attackerContract := common.HexToAddress("0x3333333333333333333333333333333333333333")

	drain := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	chain.StateDB.AddBalance(victim, new(big.Int).Set(drain))
	// A non-zero tx hash is required for the restricting plugin to take effect.
	chain.StateDB.SetTxContext(common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000dEaD"), 0)

	evm := pdEVM(t, victim, chain)

	input := pdPrecompileInput(t, beneficiary, []restricting.RestrictingPlan{{Epoch: 1, Amount: new(big.Int).Set(drain)}})
	chain.StateDB.SetCode(attackerContract, pdDelegateCallCode(t, cvm.RestrictingContractAddr, input))

	// The victim calls the attacker contract, which delegatecalls the precompile.
	if _, _, err := evm.Call(AccountRef(victim), attackerContract, nil, 8000000, big.NewInt(0)); err != nil {
		t.Fatalf("outer call failed: %v", err)
	}

	if got := chain.StateDB.GetBalance(victim); got.Cmp(drain) != 0 {
		t.Fatalf("delegatecall into system precompile moved funds: victim balance %s, want %s", got, drain)
	}
	if _, bizErr := plugin.RestrictingInstance().GetRestrictingInfo(beneficiary, chain.StateDB); bizErr == nil {
		t.Fatalf("delegatecall into system precompile created a restricting plan")
	}
}

// TestSystemPrecompileAcceptsDirectCall verifies the fix does not break the
// legitimate path: an account calling the precompile directly still works.
func TestSystemPrecompileAcceptsDirectCall(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	from := common.HexToAddress("0x4444444444444444444444444444444444444444")
	beneficiary := common.HexToAddress("0x5555555555555555555555555555555555555555")

	amount := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	chain.StateDB.AddBalance(from, new(big.Int).Set(amount))
	chain.StateDB.SetTxContext(common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000bEEF"), 0)

	evm := pdEVM(t, from, chain)
	input := pdPrecompileInput(t, beneficiary, []restricting.RestrictingPlan{{Epoch: 1, Amount: new(big.Int).Set(amount)}})

	if _, _, err := evm.Call(AccountRef(from), cvm.RestrictingContractAddr, input, 8000000, big.NewInt(0)); err != nil {
		t.Fatalf("direct call to the restricting precompile failed: %v", err)
	}
	if got := chain.StateDB.GetBalance(from); got.Sign() != 0 {
		t.Fatalf("direct call did not debit the caller: balance %s", got)
	}
	if _, bizErr := plugin.RestrictingInstance().GetRestrictingInfo(beneficiary, chain.StateDB); bizErr != nil {
		t.Fatalf("direct call did not create the restricting plan: %v", bizErr)
	}
}
