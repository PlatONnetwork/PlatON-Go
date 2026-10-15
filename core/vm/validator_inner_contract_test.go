// Copyright 2026 The PlatON Network Authors
// This file is part of the PlatON-Go library.

package vm

import (
	"encoding/binary"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func encodeValidatorInnerInput(t *testing.T, parts ...[]byte) []byte {
	t.Helper()
	data, err := rlp.EncodeToBytes(parts)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func newValidatorInnerForTest(t *testing.T) *validatorInnerContract {
	t.Helper()
	chain := newMockChain()
	t.Cleanup(func() { chain.SnapDB.Clear() })
	return &validatorInnerContract{
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, chain),
	}
}

func TestValidatorInnerCurrentValidatorsDispatchSpelling(t *testing.T) {
	vic := newValidatorInnerForTest(t)

	assertNotUndefined := func(t *testing.T, name string, input []byte) {
		t.Helper()
		_, err := vic.Run(input)
		if err != nil && err.Error() == "Function undefined" {
			t.Fatalf("%s rejected as undefined: %v", name, err)
		}
	}

	// Correct spelling must be accepted (issue #2359).
	assertNotUndefined(t, "CurrentValidators", encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeCurrent),
		[]byte("CurrentValidators"),
	))

	// Historical misspelling remains accepted as an alias.
	assertNotUndefined(t, "CurrentValidatros", encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeCurrent),
		[]byte("CurrentValidatros"),
	))
}

func TestValidatorInnerShortCalldataNoPanic(t *testing.T) {
	vic := newValidatorInnerForTest(t)

	cases := []struct {
		name  string
		parts [][]byte
	}{
		{
			name: "UpdateValidators missing arg",
			parts: [][]byte{
				common.Int64ToBytes(txTypeUpdate),
				[]byte("UpdateValidators"),
			},
		},
		{
			name: "SwitchValidators missing arg",
			parts: [][]byte{
				common.Int64ToBytes(txTypeSwitch),
				[]byte("SwitchValidators"),
			},
		},
		{
			name: "SwitchValidators short number",
			parts: [][]byte{
				common.Int64ToBytes(txTypeSwitch),
				[]byte("SwitchValidators"),
				{0x01, 0x02}, // < 8 bytes
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("unexpected panic: %v", r)
				}
			}()
			_, err := vic.Run(encodeValidatorInnerInput(t, tc.parts...))
			if err == nil {
				t.Fatal("expected error for malformed calldata")
			}
		})
	}
}

func TestValidatorInnerSwitchValidatorsWithValidArg(t *testing.T) {
	vic := newValidatorInnerForTest(t)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 100)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	// Empty validator state may return a decode error; must not panic.
	_, _ = vic.Run(encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeSwitch),
		[]byte("SwitchValidators"),
		buf,
	))
}
