// Copyright 2026 The PlatON Network Authors
// This file is part of the PlatON-Go library.

package vm

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
)

// encodeRestrictingCreate builds RLP calldata for createRestrictingPlan (fn 4000).
func encodeRestrictingCreate(t *testing.T, account common.Address, plans []restricting.RestrictingPlan) []byte {
	t.Helper()
	p0, err := rlp.EncodeToBytes(common.Uint16ToBytes(4000))
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

// delegatecallRestrictingCode returns runtime bytecode that CODECOPYs the
// embedded calldata and DELEGATECALLs (0xF4) to target. Header is 40 bytes.
func delegatecallRestrictingCode(t *testing.T, target common.Address, calldata []byte) []byte {
	t.Helper()
	L := len(calldata)
	code := []byte{
		0x60, byte(L), // PUSH1 L
		0x60, 0x28, // PUSH1 40 (data offset)
		0x60, 0x00, // PUSH1 0
		0x39,       // CODECOPY
		0x60, 0x00, // PUSH1 0 (retSize)
		0x60, 0x00, // PUSH1 0 (retOffset)
		0x60, byte(L), // PUSH1 L (inSize)
		0x60, 0x00, // PUSH1 0 (inOffset)
		0x73, // PUSH20
	}
	code = append(code, target.Bytes()...)
	code = append(code,
		0x5A, // GAS
		0xF4, // DELEGATECALL
		0x50, // POP
		0x00, // STOP
	)
	if len(code) != 40 {
		t.Fatalf("delegate header must be 40 bytes, got %d", len(code))
	}
	return append(code, calldata...)
}

// staticcallRestrictingCode is the same pattern with STATICCALL (0xFA).
func staticcallRestrictingCode(t *testing.T, target common.Address, calldata []byte) []byte {
	t.Helper()
	L := len(calldata)
	code := []byte{
		0x60, byte(L),
		0x60, 0x28,
		0x60, 0x00,
		0x39,
		0x60, 0x00,
		0x60, 0x00,
		0x60, byte(L),
		0x60, 0x00,
		0x73,
	}
	code = append(code, target.Bytes()...)
	code = append(code,
		0x5A, // GAS
		0xFA, // STATICCALL
		0x50, // POP
		0x00, // STOP
	)
	if len(code) != 40 {
		t.Fatalf("staticcall header must be 40 bytes, got %d", len(code))
	}
	return append(code, calldata...)
}

// callRestrictingCode CALLs (0xF1) target with value 0. Header is 42 bytes.
func callRestrictingCode(t *testing.T, target common.Address, calldata []byte) []byte {
	t.Helper()
	L := len(calldata)
	code := []byte{
		0x60, byte(L),
		0x60, 42,
		0x60, 0x00,
		0x39,
		0x60, 0x00,
		0x60, 0x00,
		0x60, byte(L),
		0x60, 0x00,
		0x60, 0x00, // value
		0x73,
	}
	code = append(code, target.Bytes()...)
	code = append(code,
		0x5A, // GAS
		0xF1, // CALL
		0x50, // POP
		0x00, // STOP
	)
	if len(code) != 42 {
		t.Fatalf("call header must be 42 bytes, got %d", len(code))
	}
	return append(code, calldata...)
}

func restrictingGuardEVM(t *testing.T, origin common.Address, activeVersion uint32, sdb snapshotdb.DB, state StateDB) *EVM {
	t.Helper()
	gov.InitGenesisGovernParam(common.ZeroHash, sdb, 2048)
	gov.NewGovDB(snapshotdb.Instance()).AddActiveVersion(activeVersion, 0, state)
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
	return NewEVM(blockCtx, TxContext{Origin: origin, GasPrice: big.NewInt(0)}, sdb, state, &params.ChainConfig{}, Config{})
}

func TestHawkingDelegatecallRestrictingRejected(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	victim := common.HexToAddress("0x1111111111111111111111111111111111111111")
	attackerEOA := common.HexToAddress("0x2222222222222222222222222222222222222222")
	attackerContract := common.HexToAddress("0x3333333333333333333333333333333333333333")

	drain := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	chain.StateDB.AddBalance(victim, new(big.Int).Set(drain))
	chain.StateDB.SetTxContext(common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000dEaD"), 0)

	evm := restrictingGuardEVM(t, victim, params.FORKVERSION_1_6_0, chain.SnapDB, chain.StateDB)

	plans := []restricting.RestrictingPlan{{Epoch: 1, Amount: new(big.Int).Set(drain)}}
	calldata := encodeRestrictingCreate(t, attackerEOA, plans)
	chain.StateDB.SetCode(attackerContract, delegatecallRestrictingCode(t, cvm.RestrictingContractAddr, calldata))

	before := new(big.Int).Set(chain.StateDB.GetBalance(victim))
	if _, _, err := evm.Call(AccountRef(victim), attackerContract, nil, 8000000, big.NewInt(0)); err != nil {
		t.Fatalf("victim->attacker CALL failed: %v", err)
	}
	after := chain.StateDB.GetBalance(victim)
	if after.Cmp(before) != 0 {
		t.Fatalf("hawking must reject delegatecall debit: before=%s after=%s", before, after)
	}
	if _, bizErr := plugin.RestrictingInstance().GetRestrictingInfo(attackerEOA, chain.StateDB); bizErr == nil {
		t.Fatal("hawking must not record a restricting plan via delegatecall")
	}
}

func TestHawkingDirectCallRestrictingStillWorks(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	victim := common.HexToAddress("0x4444444444444444444444444444444444444444")
	drain := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	chain.StateDB.AddBalance(victim, new(big.Int).Set(drain))
	chain.StateDB.SetTxContext(common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000c0e7"), 0)

	evm := restrictingGuardEVM(t, victim, params.FORKVERSION_1_6_0, chain.SnapDB, chain.StateDB)

	plans := []restricting.RestrictingPlan{{Epoch: 1, Amount: new(big.Int).Set(drain)}}
	calldata := encodeRestrictingCreate(t, victim, plans)

	if _, _, err := evm.Call(AccountRef(victim), cvm.RestrictingContractAddr, calldata, 8000000, big.NewInt(0)); err != nil {
		t.Fatalf("direct CALL to restricting failed: %v", err)
	}
	if after := chain.StateDB.GetBalance(victim); after.Sign() != 0 {
		t.Fatalf("direct call should debit victim to 0, got %s", after)
	}
	info, bizErr := plugin.RestrictingInstance().GetRestrictingInfo(victim, chain.StateDB)
	if bizErr != nil {
		t.Fatalf("expected restricting plan for victim: %v", bizErr)
	}
	if info.Balance.ToInt().Cmp(drain) != 0 {
		t.Fatalf("victim plan = %s, want %s", info.Balance.ToInt(), drain)
	}
}

func TestHawkingStaticcallRestrictingRejected(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	victim := common.HexToAddress("0x5555555555555555555555555555555555555555")
	attackerEOA := common.HexToAddress("0x6666666666666666666666666666666666666666")
	callerContract := common.HexToAddress("0x7777777777777777777777777777777777777777")

	drain := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	// Fund the calling contract: STATICCALL sets CallerAddress to the contract itself.
	chain.StateDB.AddBalance(callerContract, new(big.Int).Set(drain))
	chain.StateDB.SetTxContext(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000005ba1"), 0)

	evm := restrictingGuardEVM(t, victim, params.FORKVERSION_1_6_0, chain.SnapDB, chain.StateDB)

	plans := []restricting.RestrictingPlan{{Epoch: 1, Amount: new(big.Int).Set(drain)}}
	calldata := encodeRestrictingCreate(t, attackerEOA, plans)
	chain.StateDB.SetCode(callerContract, staticcallRestrictingCode(t, cvm.RestrictingContractAddr, calldata))

	before := new(big.Int).Set(chain.StateDB.GetBalance(callerContract))
	if _, _, err := evm.Call(AccountRef(victim), callerContract, nil, 8000000, big.NewInt(0)); err != nil {
		t.Fatalf("outer CALL failed: %v", err)
	}
	after := chain.StateDB.GetBalance(callerContract)
	if after.Cmp(before) != 0 {
		t.Fatalf("hawking must reject staticcall state write: before=%s after=%s", before, after)
	}
	if _, bizErr := plugin.RestrictingInstance().GetRestrictingInfo(attackerEOA, chain.StateDB); bizErr == nil {
		t.Fatal("hawking must not record a restricting plan via staticcall")
	}
}

func TestRunInheritsWasmReadOnly(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()
	origin := common.HexToAddress("0x1212121212121212121212121212121212121212")
	evm := restrictingGuardEVM(t, origin, params.FORKVERSION_1_6_0, chain.SnapDB, chain.StateDB)
	evm.interpreter = &WASMInterpreter{readOnly: true}
	contract := NewContract(AccountRef(origin), AccountRef(cvm.RestrictingContractAddr), big.NewInt(0), 100000)
	contract.SetCallCode(&cvm.RestrictingContractAddr, common.Hash{}, nil)
	_, err := run(evm, contract, nil, false)
	if err != ErrWriteProtection {
		t.Fatalf("wasm readOnly call = %v, want ErrWriteProtection", err)
	}
}

// encodePlatONInput builds PPOS calldata: RLP list of RLP-encoded fn code and params.
func encodePlatONInput(t *testing.T, fn uint16, args ...[]byte) []byte {
	t.Helper()
	p0, err := rlp.EncodeToBytes(common.Uint16ToBytes(fn))
	if err != nil {
		t.Fatal(err)
	}
	parts := [][]byte{p0}
	for _, a := range args {
		enc, err := rlp.EncodeToBytes(a)
		if err != nil {
			t.Fatal(err)
		}
		parts = append(parts, enc)
	}
	in, err := rlp.EncodeToBytes(parts)
	if err != nil {
		t.Fatal(err)
	}
	return in
}

func TestIsPlatONQuery(t *testing.T) {
	inner := func(txType int64, name string) []byte {
		return encodeValidatorInnerInput(t, common.Int64ToBytes(txType), []byte(name))
	}
	cases := []struct {
		name string
		p    PrecompiledContract
		in   []byte
		want bool
	}{
		{"staking query", &StakingContract{}, encodePlatONInput(t, QueryCandidateInfo), true},
		{"staking reward query", &StakingContract{}, encodePlatONInput(t, GetAvgPackTime), true},
		{"staking delegate", &StakingContract{}, encodePlatONInput(t, TxDelegate), false},
		{"staking with gov query code", &StakingContract{}, encodePlatONInput(t, GetProposal), false},
		{"gov query", &GovContract{}, encodePlatONInput(t, ListProposal), true},
		{"gov vote", &GovContract{}, encodePlatONInput(t, Vote), false},
		{"slashing check", &SlashingContract{}, encodePlatONInput(t, CheckDuplicateSign), true},
		{"slashing report", &SlashingContract{}, encodePlatONInput(t, TxReportDuplicateSign), false},
		{"restricting query", &RestrictingContract{}, encodePlatONInput(t, QueryRestrictingInfo), true},
		{"restricting create", &RestrictingContract{}, encodePlatONInput(t, TxCreateRestrictingPlan), false},
		{"reward query", &DelegateRewardContract{}, encodePlatONInput(t, QueryDelegateReward), true},
		{"reward withdraw", &DelegateRewardContract{}, encodePlatONInput(t, TxWithdrawDelegateReward), false},
		{"inner current", &validatorInnerContract{}, inner(txTypeCurrent, "CurrentValidatros"), true},
		{"inner next", &validatorInnerContract{}, inner(txTypeNext, "NextValidators"), true},
		{"inner switch under query name", &validatorInnerContract{}, inner(txTypeSwitch, "CurrentValidatros"), false},
		{"inner update", &validatorInnerContract{}, inner(txTypeUpdate, "UpdateValidators"), false},
		{"empty", &StakingContract{}, nil, false},
		{"not a list", &StakingContract{}, []byte{0x01, 0x02}, false},
	}
	for _, c := range cases {
		if got := isPlatONQuery(c.p, c.in); got != c.want {
			t.Errorf("%s: isPlatONQuery = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestHawkingStaticcallRestrictingQueryAllowed(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	caller := common.HexToAddress("0x1313131313131313131313131313131313131313")
	evm := restrictingGuardEVM(t, caller, params.FORKVERSION_1_6_0, chain.SnapDB, chain.StateDB)

	query := encodePlatONInput(t, QueryRestrictingInfo, caller.Bytes())
	ret, _, err := evm.StaticCall(AccountRef(caller), cvm.RestrictingContractAddr, query, 1000000)
	if err != nil {
		t.Fatalf("static query = %v, want success", err)
	}
	if len(ret) == 0 {
		t.Fatal("static query returned no result")
	}

	plans := []restricting.RestrictingPlan{{Epoch: 1, Amount: big.NewInt(1)}}
	write := encodeRestrictingCreate(t, caller, plans)
	if _, _, err := evm.StaticCall(AccountRef(caller), cvm.RestrictingContractAddr, write, 1000000); err != ErrWriteProtection {
		t.Fatalf("static write = %v, want ErrWriteProtection", err)
	}
}

func TestHawkingNestedReadOnlyQueryAllowed(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	origin := common.HexToAddress("0x1515151515151515151515151515151515151515")
	evm := restrictingGuardEVM(t, origin, params.FORKVERSION_1_6_0, chain.SnapDB, chain.StateDB)
	evm.interpreter = &WASMInterpreter{readOnly: true}
	contract := NewContract(AccountRef(origin), AccountRef(cvm.RestrictingContractAddr), big.NewInt(0), 1000000)
	contract.SetCallCode(&cvm.RestrictingContractAddr, common.Hash{}, nil)

	if _, err := run(evm, contract, encodePlatONInput(t, QueryRestrictingInfo, origin.Bytes()), false); err != nil {
		t.Fatalf("nested read-only query = %v, want success", err)
	}
	plans := []restricting.RestrictingPlan{{Epoch: 1, Amount: big.NewInt(1)}}
	if _, err := run(evm, contract, encodeRestrictingCreate(t, origin, plans), false); err != ErrWriteProtection {
		t.Fatalf("nested read-only write = %v, want ErrWriteProtection", err)
	}
}

// TestRunReadOnlyInheritanceIsHawkingOnly: before Hawking a CALL out of a
// static frame into new code stays writable, so historical blocks replay.
func TestRunReadOnlyInheritanceIsHawkingOnly(t *testing.T) {
	sstore := []byte{0x60, 0x01, 0x60, 0x00, 0x55, 0x00} // PUSH1 1 PUSH1 0 SSTORE STOP
	cases := []struct {
		name    string
		version uint32
		want    error
	}{
		{"pre-hawking", params.FORKVERSION_1_5_0, nil},
		{"hawking", params.FORKVERSION_1_6_0, ErrWriteProtection},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			chain := newMockChain()
			defer chain.SnapDB.Clear()

			origin := common.HexToAddress("0x1616161616161616161616161616161616161616")
			target := common.HexToAddress("0x1717171717171717171717171717171717171717")
			evm := restrictingGuardEVM(t, origin, c.version, chain.SnapDB, chain.StateDB)
			chain.StateDB.SetCode(target, sstore)
			chain.StateDB.AddAddressToAccessList(target)
			evm.interpreter = &WASMInterpreter{readOnly: true}

			contract := NewContract(AccountRef(origin), AccountRef(target), big.NewInt(0), 100000)
			contract.SetCallCode(&target, chain.StateDB.GetCodeHash(target), sstore)
			if _, err := run(evm, contract, nil, false); err != c.want {
				t.Fatalf("run = %v, want %v", err, c.want)
			}
		})
	}
}

func TestHawkingNestedStaticCallRestrictingRejected(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	victim := common.HexToAddress("0xabcabcabcabcabcabcabcabcabcabcabcabcabca")
	attackerEOA := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	bridge := common.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")

	drain := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	chain.StateDB.AddBalance(bridge, new(big.Int).Set(drain))
	chain.StateDB.SetTxContext(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000e51"), 0)

	evm := restrictingGuardEVM(t, victim, params.FORKVERSION_1_6_0, chain.SnapDB, chain.StateDB)
	plans := []restricting.RestrictingPlan{{Epoch: 1, Amount: new(big.Int).Set(drain)}}
	calldata := encodeRestrictingCreate(t, attackerEOA, plans)
	chain.StateDB.SetCode(bridge, callRestrictingCode(t, cvm.RestrictingContractAddr, calldata))

	before := new(big.Int).Set(chain.StateDB.GetBalance(bridge))
	if _, _, err := evm.StaticCall(AccountRef(victim), bridge, nil, 8000000); err != nil {
		t.Fatalf("outer STATICCALL failed: %v", err)
	}
	after := chain.StateDB.GetBalance(bridge)
	if after.Cmp(before) != 0 {
		t.Fatalf("nested static call mutated bridge: before=%s after=%s", before, after)
	}
	if _, bizErr := plugin.RestrictingInstance().GetRestrictingInfo(attackerEOA, chain.StateDB); bizErr == nil {
		t.Fatal("nested static call must not record a restricting plan")
	}
}

// TestPreHawkingDelegatecallRestrictingStillAllowed keeps historical semantics
// before 1.6.0 so chain replay of old blocks is unchanged.
func TestPreHawkingDelegatecallRestrictingStillAllowed(t *testing.T) {
	chain := newMockChain()
	defer chain.SnapDB.Clear()

	victim := common.HexToAddress("0x8888888888888888888888888888888888888888")
	attackerEOA := common.HexToAddress("0x9999999999999999999999999999999999999999")
	attackerContract := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	drain := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	chain.StateDB.AddBalance(victim, new(big.Int).Set(drain))
	chain.StateDB.SetTxContext(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000001d1"), 0)

	// Active version below 1.6.0 — guard must not apply.
	evm := restrictingGuardEVM(t, victim, params.FORKVERSION_1_5_0, chain.SnapDB, chain.StateDB)

	plans := []restricting.RestrictingPlan{{Epoch: 1, Amount: new(big.Int).Set(drain)}}
	calldata := encodeRestrictingCreate(t, attackerEOA, plans)
	chain.StateDB.SetCode(attackerContract, delegatecallRestrictingCode(t, cvm.RestrictingContractAddr, calldata))

	if _, _, err := evm.Call(AccountRef(victim), attackerContract, nil, 8000000, big.NewInt(0)); err != nil {
		t.Fatalf("pre-hawking victim->attacker CALL failed: %v", err)
	}
	if after := chain.StateDB.GetBalance(victim); after.Sign() != 0 {
		t.Fatalf("pre-hawking delegatecall should still debit victim, got %s", after)
	}
	info, bizErr := plugin.RestrictingInstance().GetRestrictingInfo(attackerEOA, chain.StateDB)
	if bizErr != nil {
		t.Fatalf("pre-hawking should still record attacker plan: %v", bizErr)
	}
	if info.Balance.ToInt().Cmp(drain) != 0 {
		t.Fatalf("attacker plan = %s, want %s", info.Balance.ToInt(), drain)
	}
}
