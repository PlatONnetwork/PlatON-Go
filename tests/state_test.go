// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/vm"
)

func TestState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
	// Long tests:
	st.skipShortMode(`^stQuadraticComplexityTest/`)
	// Very time consuming
	st.skipLoad(`^stTimeConsuming/`)
	// Broken tests:
	st.skipLoad(`^stTransactionTest/OverflowGasRequire\.json`) // gasLimit > 256 bits
	st.skipLoad(`^stTransactionTest/zeroSigTransa[^/]*\.json`) // EIP-86 is not supported yet
	// Expected failures:
	st.fails(`^stRevertTest/RevertPrecompiledTouch\.json/EIP158`, "bug in test")
	st.fails(`^stRevertTest/RevertPrecompiledTouch\.json/Byzantium`, "bug in test")

	// For Istanbul, older tests were moved into LegacyTests
	for _, dir := range []string{
		stateTestDir,
		legacyStateTestDir,
	} {
		st.walk(t, dir, func(t *testing.T, name string, test *StateTest) {
			for _, subtest := range test.Subtests() {
				subtest := subtest
				key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
				name := name + "/" + key
				t.Run(key, func(t *testing.T) {
					withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
						_, err := test.Run(subtest, vmconfig)
						return st.checkFailure(t, name, err)
					})
				})
			}
		})
	}
}

// Transactions with gasLimit above this value will not get a VM trace on failure.
const traceErrorLimit = 400000

func withTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	err := test(vm.Config{})
	if err == nil {
		return
	}
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for EVM trace")
		return
	}
	tracer := vm.NewStructLogger(nil)
	err2 := test(vm.Config{Debug: true, Tracer: tracer})
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	buf := new(bytes.Buffer)
	vm.WriteTrace(buf, tracer.StructLogs())
	if buf.Len() == 0 {
		t.Log("no EVM operation logs generated")
	} else {
		t.Log("EVM operation log:\n" + buf.String())
	}
	t.Logf("EVM output: 0x%x", tracer.Output())
	t.Logf("EVM error: %v", tracer.Error())
}
