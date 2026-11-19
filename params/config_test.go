// Copyright 2017 The go-ethereum Authors
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

package params

import (
	"bytes"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestLegacyHawkingHeight(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  *big.Int
	}{
		{`{}`, nil},
		{`{"diracBlock":0}`, big.NewInt(0)},
		{`{"diracBlock":123}`, big.NewInt(123)},
		{`{"diracBlock":1,"hawkingBlock":20}`, big.NewInt(20)},
		{`{"diracBlock":"obsolete","hawkingBlock":20}`, big.NewInt(20)},
		{`{"hawkingBlock":null,"diracBlock":1}`, nil},
	} {
		var cfg ChainConfig
		if err := json.Unmarshal([]byte(tc.input), &cfg); err != nil || !configNumEqual(cfg.HawkingBlock, tc.want) {
			t.Fatalf("%s: height %v, error %v", tc.input, cfg.HawkingBlock, err)
		}
		out, err := json.Marshal(&cfg)
		if err != nil || bytes.Contains(out, []byte("diracBlock")) || bytes.Contains(out, []byte("hawkingBlock")) != (tc.want != nil) {
			t.Fatalf("unexpected config: %s, %v", out, err)
		}
	}
}

func TestInnerAdminCompatibilityAtActivation(t *testing.T) {
	stored := &ChainConfig{PauliBlock: big.NewInt(0), HawkingBlock: big.NewInt(10)}
	next := &ChainConfig{PauliBlock: big.NewInt(0), HawkingBlock: big.NewInt(10), Cbft: &CbftConfig{InnerAdmin: common.Address{1}}}
	if stored.CheckCompatible(next, 9) != nil || stored.CheckCompatible(next, 10) == nil {
		t.Fatal("admin migration must be allowed only before activation")
	}
}

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new *ChainConfig
		head        uint64
		wantErr     *ConfigCompatError
	}
	tests := []test{
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 0, wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 100, wantErr: nil},
		{
			stored:  &ChainConfig{EIP155Block: big.NewInt(10)},
			new:     &ChainConfig{EIP155Block: big.NewInt(20)},
			head:    9,
			wantErr: nil,
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &ChainConfig{EIP155Block: nil},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "EIP155 fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    nil,
				RewindTo:     0,
			},
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &ChainConfig{EIP155Block: big.NewInt(1)},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "EIP155 fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    big.NewInt(1),
				RewindTo:     0,
			},
		},
	}

	for _, test := range tests {
		err := test.stored.CheckCompatible(test.new, test.head)
		if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		}
	}
}
