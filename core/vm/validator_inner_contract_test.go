// Copyright 2026 The PlatON Network Authors
// This file is part of the PlatON-Go library.

package vm

import (
	"encoding/binary"
	"errors"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
)

func encodeValidatorInnerInput(t *testing.T, parts ...[]byte) []byte {
	t.Helper()
	data, err := rlp.EncodeToBytes(parts)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func newValidatorInnerForTest(t *testing.T, caller common.Address, activeVersion uint32) *validatorInnerContract {
	t.Helper()
	chain := newMockChain()
	t.Cleanup(func() { chain.SnapDB.Clear() })
	evm := newEvm(blockNumber, blockHash, chain)
	if activeVersion != 0 {
		gov.NewGovDB(chain.SnapDB).AddActiveVersion(activeVersion, 0, chain.StateDB)
	}
	return &validatorInnerContract{
		Contract: newContract(common.Big0, caller),
		Evm:      evm,
	}
}

func setInnerAdmin(evm *EVM, admin common.Address) {
	if evm.chainConfig == nil {
		evm.chainConfig = &params.ChainConfig{}
	}
	if evm.chainConfig.Cbft == nil {
		evm.chainConfig.Cbft = &params.CbftConfig{}
	}
	evm.chainConfig.Cbft.InnerAdmin = admin
}

func TestValidatorInnerCurrentValidatorsDispatchSpelling(t *testing.T) {
	vic := newValidatorInnerForTest(t, sender, 0)

	assertNotUndefined := func(t *testing.T, contract *validatorInnerContract, name string, input []byte) {
		t.Helper()
		_, err := contract.Run(input)
		if err != nil && err.Error() == "Function undefined" {
			t.Fatalf("%s rejected as undefined: %v", name, err)
		}
	}

	// Correct spelling is Hawking-only. Before 1.6.0 it stays undefined.
	_, err := vic.Run(encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeCurrent),
		[]byte("CurrentValidators"),
	))
	if err == nil || err.Error() != "Function undefined" {
		t.Fatalf("pre-Hawking CurrentValidators = %v, want Function undefined", err)
	}
	_, err = vic.Run(encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeSwitch),
		[]byte("CurrentValidators"),
		[]byte{0, 0, 0, 0, 0, 0, 0, 1},
	))
	if err == nil || err.Error() != "Function undefined" {
		t.Fatalf("pre-Hawking switch via new alias = %v, want Function undefined", err)
	}

	hawking := newValidatorInnerForTest(t, sender, params.FORKVERSION_1_6_0)
	assertNotUndefined(t, hawking, "CurrentValidators", encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeCurrent),
		[]byte("CurrentValidators"),
	))

	// Historical misspelling remains accepted as an alias.
	assertNotUndefined(t, vic, "CurrentValidatros", encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeCurrent),
		[]byte("CurrentValidatros"),
	))
}

func TestValidatorInnerShortCalldataNoPanic(t *testing.T) {
	vic := newValidatorInnerForTest(t, sender, 0)

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
	vic := newValidatorInnerForTest(t, sender, 0)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 100)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	// Empty validator state may return a decode error; must not panic.
	// Pre-Hawking: any caller is accepted (historical replay).
	_, _ = vic.Run(encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeSwitch),
		[]byte("SwitchValidators"),
		buf,
	))
}

func TestValidatorInnerHawkingACL(t *testing.T) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 100)
	switchInput := encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeSwitch),
		[]byte("SwitchValidators"),
		buf,
	)
	// Minimal valid UpdateValidators JSON (empty nodes rejected after ACL).
	updateInput := encodeValidatorInnerInput(t,
		common.Int64ToBytes(txTypeUpdate),
		[]byte("UpdateValidators"),
		[]byte(`{"validateNodes":[]}`),
	)

	t.Run("unauthorized rejected after Hawking", func(t *testing.T) {
		vic := newValidatorInnerForTest(t, sender, params.FORKVERSION_1_6_0)
		setInnerAdmin(vic.Evm, common.HexToAddress("0x0101010101010101010101010101010101010101"))
		_, err := vic.Run(switchInput)
		if !errors.Is(err, errUnauthorizedValidatorInner) {
			t.Fatalf("SwitchValidators: want unauthorized, got %v", err)
		}
		_, err = vic.Run(updateInput)
		if !errors.Is(err, errUnauthorizedValidatorInner) {
			t.Fatalf("UpdateValidators: want unauthorized, got %v", err)
		}
	})

	t.Run("admin accepted after Hawking", func(t *testing.T) {
		admin := common.HexToAddress("0x0101010101010101010101010101010101010101")
		vic := newValidatorInnerForTest(t, admin, params.FORKVERSION_1_6_0)
		setInnerAdmin(vic.Evm, admin)
		_, err := vic.Run(switchInput)
		if errors.Is(err, errUnauthorizedValidatorInner) {
			t.Fatalf("admin SwitchValidators must not be unauthorized: %v", err)
		}
		_, err = vic.Run(updateInput)
		if errors.Is(err, errUnauthorizedValidatorInner) {
			t.Fatalf("admin UpdateValidators must not be unauthorized: %v", err)
		}
		// Empty nodes still rejected after ACL passes.
		if err == nil || err.Error() != "Empty validator nodes" {
			t.Fatalf("admin UpdateValidators empty nodes: want Empty validator nodes, got %v", err)
		}
	})

	t.Run("reads unrestricted after Hawking", func(t *testing.T) {
		vic := newValidatorInnerForTest(t, sender, params.FORKVERSION_1_6_0)
		_, err := vic.Run(encodeValidatorInnerInput(t,
			common.Int64ToBytes(txTypeCurrent),
			[]byte("CurrentValidators"),
		))
		if errors.Is(err, errUnauthorizedValidatorInner) {
			t.Fatalf("CurrentValidators must not require admin: %v", err)
		}
	})

	t.Run("pre-Hawking any caller still accepted", func(t *testing.T) {
		vic := newValidatorInnerForTest(t, sender, 0)
		_, err := vic.Run(switchInput)
		if errors.Is(err, errUnauthorizedValidatorInner) {
			t.Fatalf("pre-Hawking must not enforce ACL: %v", err)
		}
	})
}
