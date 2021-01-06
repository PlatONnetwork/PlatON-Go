// Copyright 2016 The go-ethereum Authors
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

package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	ethereum "github.com/PlatONnetwork/PlatON-Go"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
)

// Verify that Client implements the ethereum interfaces.
var (
	_ = ethereum.ChainReader(&Client{})
	_ = ethereum.TransactionReader(&Client{})
	_ = ethereum.ChainStateReader(&Client{})
	_ = ethereum.ChainSyncReader(&Client{})
	_ = ethereum.ContractCaller(&Client{})
	_ = ethereum.GasEstimator(&Client{})
	_ = ethereum.GasPricer(&Client{})
	_ = ethereum.LogFilterer(&Client{})
	_ = ethereum.PendingStateReader(&Client{})
	// _ = ethereum.PendingStateEventer(&Client{})
	_ = ethereum.PendingContractCaller(&Client{})
)

const (
	host = ""
)

func TestPlatonCall(t *testing.T) {
	ctx := context.Background()
	common.SetAddressPrefix("atp")
	address := common.MustBech32ToAddress("atp19rnu40l5dux6p2n2wvlft40lqt4hya0l99fspf")
	client, err := DialContext(ctx, host)
	if err != nil {
		t.Fatalf("rawurl %s err: %s", host, err)
	}

	for _, testCase := range []struct {
		name    string
		address common.Address
		output  interface{}
		err     error
	}{
		{
			"get_balance",
			address,
			"0x66248f04690c5",
			nil,
		},
		{
			"get_nonce",
			address,
			2523,
			nil,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			switch testCase.name {
			case "get_balance":
				output, err := client.BalanceAt(ctx, testCase.address, nil)
				if (testCase.err == nil) != (err == nil) {
					t.Fatalf("expected error %v but got %v", testCase.err, err)
				}
				if testCase.err != nil {
					if testCase.err.Error() != err.Error() {
						t.Fatalf("expected error %v but got %v", testCase.err, err)
					}
				} else if !reflect.DeepEqual(testCase.output, hexutil.EncodeBig(output)) {
					t.Fatalf("expected filter arg %v but got %v", testCase.output, output)
				}
			case "get_nonce":
				output, err := client.NonceAt(ctx, testCase.address, nil)
				if (testCase.err == nil) != (err == nil) {
					t.Fatalf("expected error %v but got %v", testCase.err, err)
				}
				if testCase.err != nil {
					if testCase.err.Error() != err.Error() {
						t.Fatalf("expected error %v but got %v", testCase.err, err)
					}
				} else if !reflect.DeepEqual(testCase.output, output) {
					t.Fatalf("expected filter arg %v but got %v", testCase.output, output)
				}
			}
		})
	}
}

func TestToFilterArg(t *testing.T) {
	blockHashErr := fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
	addresses := []common.Address{
		common.MustBech32ToAddress("atx16dnj9t0v8mwt98yw0ddy0u6j6uqnjdrz98vr3t"),
	}
	blockHash := common.HexToHash(
		"0xeb94bb7d78b73657a9d7a99792413f50c0a45c51fc62bdcb08a53f18e9a2b4eb",
	)

	for _, testCase := range []struct {
		name   string
		input  ethereum.FilterQuery
		output interface{}
		err    error
	}{
		{
			"without BlockHash",
			ethereum.FilterQuery{
				Addresses: addresses,
				FromBlock: big.NewInt(1),
				ToBlock:   big.NewInt(2),
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "0x1",
				"toBlock":   "0x2",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with nil fromBlock and nil toBlock",
			ethereum.FilterQuery{
				Addresses: addresses,
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "0x0",
				"toBlock":   "latest",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with blockhash",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"blockHash": blockHash,
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with blockhash and from block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				FromBlock: big.NewInt(1),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
		{
			"with blockhash and to block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				ToBlock:   big.NewInt(1),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
		{
			"with blockhash and both from / to block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				FromBlock: big.NewInt(1),
				ToBlock:   big.NewInt(2),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			output, err := toFilterArg(testCase.input)
			if (testCase.err == nil) != (err == nil) {
				t.Fatalf("expected error %v but got %v", testCase.err, err)
			}
			if testCase.err != nil {
				if testCase.err.Error() != err.Error() {
					t.Fatalf("expected error %v but got %v", testCase.err, err)
				}
			} else if !reflect.DeepEqual(testCase.output, output) {
				t.Fatalf("expected filter arg %v but got %v", testCase.output, output)
			}
		})
	}
}
