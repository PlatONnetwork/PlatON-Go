// Copyright 2018 The go-ethereum Authors
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

package graphql

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/eth/ethconfig"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core"

	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/node"
)

func TestBuildSchema(t *testing.T) {
	ddir, err := ioutil.TempDir("", "graphql-buildschema")
	if err != nil {
		t.Fatalf("failed to create temporary datadir: %v", err)
	}
	// Copy config
	conf := node.DefaultConfig
	conf.DataDir = ddir
	stack, err := node.New(&conf)
	if err != nil {
		t.Fatalf("could not create new node: %v", err)
	}
	// Make sure the schema can be parsed and matched up to the object model.
	if err := newHandler(stack, nil, []string{}, []string{}); err != nil {
		t.Errorf("Could not construct GraphQL handler: %v", err)
	}
}

// Tests that a graphQL request is successfully handled when graphql is enabled on the specified endpoint
func TestGraphQLBlockSerialization(t *testing.T) {
	stack := createNode(t, true, false)
	defer stack.Close()
	// start node
	if err := stack.Start(); err != nil {
		t.Fatalf("could not start node: %v", err)
	}

	for i, tt := range []struct {
		body string
		want string
		code int
	}{
		{ // Should return latest block
			body: `{"query": "{block{number}}","variables": null}`,
			want: `{"data":{"block":{"number":10}}}`,
			code: 200,
		},
		{ // Should return info about latest block
			body: `{"query": "{block{number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":{"number":10,"gasUsed":0,"gasLimit":11500000}}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:0){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":{"number":0,"gasUsed":0,"gasLimit":11500000}}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:-1){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":null}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:-500){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":null}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:\"0\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":{"number":0,"gasUsed":0,"gasLimit":11500000}}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:\"-33\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":null}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:\"1337\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":null}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:\"0xbad\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"errors":[{"message":"strconv.ParseInt: parsing \"0xbad\": invalid syntax"}],"data":{}}`,
			code: 400,
		},
		{ // hex strings are currently not supported. If that's added to the spec, this test will need to change
			body: `{"query": "{block(number:\"0x0\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"errors":[{"message":"strconv.ParseInt: parsing \"0x0\": invalid syntax"}],"data":{}}`,
			code: 400,
		},
		{
			body: `{"query": "{block(number:\"a\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"errors":[{"message":"strconv.ParseInt: parsing \"a\": invalid syntax"}],"data":{}}`,
			code: 400,
		},
		{
			body: `{"query": "{bleh{number}}","variables": null}"`,
			want: `{"errors":[{"message":"Cannot query field \"bleh\" on type \"Query\".","locations":[{"line":1,"column":2}]}]}`,
			code: 400,
		},
		// should return `estimateGas` as decimal
		{
			body: `{"query": "{block{ estimateGas(data:{}) }}"}`,
			want: `{"data":{"block":{"estimateGas":53000}}}`,
			code: 200,
		},
		// should return `status` as decimal
		{
			body: `{"query": "{block {number call (data : {from : \"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b\", to: \"0x6295ee1b4f6dd65047762f924ecd367c17eabf8f\", data :\"0x12a7b914\"}){data status}}}"}`,
			want: `{"data":{"block":{"number":10,"call":{"data":"0x","status":1}}}}`,
			code: 200,
		},
	} {
		resp, err := http.Post(fmt.Sprintf("%s/graphql", stack.HTTPEndpoint()), "application/json", strings.NewReader(tt.body))
		if err != nil {
			t.Fatalf("could not post: %v", err)
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read from response body: %v", err)
		}
		if have := string(bodyBytes); have != tt.want {
			t.Errorf("testcase %d %s,\nhave:\n%v\nwant:\n%v", i, tt.body, have, tt.want)
		}
		if tt.code != resp.StatusCode {
			t.Errorf("testcase %d %s,\nwrong statuscode, have: %v, want: %v", i, tt.body, resp.StatusCode, tt.code)
		}
	}
}

func TestGraphQLBlockSerializationEIP2718(t *testing.T) {
	stack := createNode(t, true, true)
	defer stack.Close()
	// start node
	if err := stack.Start(); err != nil {
		t.Fatalf("could not start node: %v", err)
	}

	for i, tt := range []struct {
		body string
		want string
		code int
	}{
		{
			body: `{"query": "{block {number transactions { from { address } to { address } value hash type accessList { address storageKeys } index}}}"}`,
			want: `{"data":{"block":{"number":1,"transactions":[{"from":{"address":"0x71562b71999873DB5b286dF957af199Ec94617F7"},"to":{"address":"0x0000000000000000000000000000000000000DAd"},"value":"0x64","hash":"0xa20f53352272dcf4acb84bd1364de8240a53bb7c7725d8516b626107c0ff77af","type":0,"accessList":[],"index":0},{"from":{"address":"0x71562b71999873DB5b286dF957af199Ec94617F7"},"to":{"address":"0x0000000000000000000000000000000000000DAd"},"value":"0x32","hash":"0xb3ea15151ca9997a2e6db0f1d94193d60f4b4c5353f437460b92579e3900b0e8","type":1,"accessList":[{"address":"0x0000000000000000000000000000000000000DAd","storageKeys":["0x0000000000000000000000000000000000000000000000000000000000000000"]}],"index":1}]}}}`,
			code: 200,
		},
	} {
		resp, err := http.Post(fmt.Sprintf("%s/graphql", stack.HTTPEndpoint()), "application/json", strings.NewReader(tt.body))
		if err != nil {
			t.Fatalf("could not post: %v", err)
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read from response body: %v", err)
		}
		if have := string(bodyBytes); have != tt.want {
			t.Errorf("testcase %d %s,\nhave:\n%v\nwant:\n%v", i, tt.body, have, tt.want)
		}
		if tt.code != resp.StatusCode {
			t.Errorf("testcase %d %s,\nwrong statuscode, have: %v, want: %v", i, tt.body, resp.StatusCode, tt.code)
		}
	}
}

// Tests that a graphQL request is not handled successfully when graphql is not enabled on the specified endpoint
func TestGraphQLHTTPOnSamePort_GQLRequest_Unsuccessful(t *testing.T) {
	stack := createNode(t, false, false)
	defer stack.Close()
	if err := stack.Start(); err != nil {
		t.Fatalf("could not start node: %v", err)
	}
	body := strings.NewReader(`{"query": "{block{number}}","variables": null}`)
	resp, err := http.Post(fmt.Sprintf("%s/graphql", stack.HTTPEndpoint()), "application/json", body)
	if err != nil {
		t.Fatalf("could not post: %v", err)
	}
	// make sure the request is not handled successfully
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func createNode(t *testing.T, gqlEnabled bool, txEnabled bool) *node.Node {
	stack, err := node.New(&node.Config{
		HTTPHost: "127.0.0.1",
		HTTPPort: 0,
		WSHost:   "127.0.0.1",
		WSPort:   0,
	})
	if err != nil {
		t.Fatalf("could not create node: %v", err)
	}
	if !gqlEnabled {
		return stack
	}
	if !txEnabled {
		createGQLService(t, stack)
	} else {
		createGQLServiceWithTransactions(t, stack)
	}
	return stack
}

func createGQLService(t *testing.T, stack *node.Node) {
	// create backend
	ethConf := &ethconfig.Config{
		Genesis: &core.Genesis{
			Config:        params.TestChainConfig,
			GasLimit:      11500000,
			EconomicModel: xcom.GetEc(xcom.DefaultUnitTestNet),
		},
		NetworkId:               1337,
		TrieCleanCache:          5,
		TrieCleanCacheJournal:   "triecache",
		TrieCleanCacheRejournal: 60 * time.Minute,
		TrieDirtyCache:          5,
		TrieTimeout:             60 * time.Minute,
		SnapshotCache:           5,
		BlockCacheLimit:         256,
		MaxFutureBlocks:         256,
	}
	gov.InitGenesisGovernParam(common.ZeroHash, snapshotdb.Instance(), 1)
	ethBackend, err := eth.New(stack, ethConf)
	if err != nil {
		t.Fatalf("could not create eth backend: %v", err)
	}
	// Create some blocks and import them
	/*chain, _ := core.GenerateChain(params.TestChainConfig, ethBackend.BlockChain().Genesis(),
		consensus.NewFaker(), ethBackend.ChainDb(), 10, func(i int, gen *core.BlockGen) {})
	_, err = ethBackend.BlockChain().InsertChain(chain)
	if err != nil {
		t.Fatalf("could not create import blocks: %v", err)
	}*/

	core.GenerateBlockChain3(params.TestChainConfig, ethBackend.BlockChain().Genesis(),
		consensus.NewFakerWithDataBase(ethBackend.ChainDb()), ethBackend.BlockChain(), 10, func(i int, gen *core.BlockGen) {})
	// create gql service
	err = New(stack, ethBackend.APIBackend, []string{}, []string{})
	if err != nil {
		t.Fatalf("could not create graphql service: %v", err)
	}
}

func createGQLServiceWithTransactions(t *testing.T, stack *node.Node) {
	// create backend
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	address := crypto.PubkeyToAddress(key.PublicKey)
	funds := big.NewInt(1000000000)
	dad := common.HexToAddress("0x0000000000000000000000000000000000000dad")

	ethConf := &ethconfig.Config{
		Genesis: &core.Genesis{
			Config:   params.TestChainConfig,
			GasLimit: 11500000,
			Alloc: core.GenesisAlloc{
				address: {Balance: funds},
				// The address 0xdad sloads 0x00 and 0x01
				dad: {
					Code: []byte{
						byte(vm.PC),
						byte(vm.PC),
						byte(vm.SLOAD),
						byte(vm.SLOAD),
					},
					Nonce:   0,
					Balance: big.NewInt(0),
				},
			},
		},
		NetworkId:               1337,
		TrieCleanCache:          5,
		TrieCleanCacheJournal:   "triecache",
		TrieCleanCacheRejournal: 60 * time.Minute,
		TrieDirtyCache:          5,
		TrieTimeout:             60 * time.Minute,
		SnapshotCache:           5,
		BlockCacheLimit:         256,
		MaxFutureBlocks:         256,
	}
	gov.InitGenesisGovernParam(common.ZeroHash, snapshotdb.Instance(), 1)

	ethBackend, err := eth.New(stack, ethConf)
	if err != nil {
		t.Fatalf("could not create eth backend: %v", err)
	}

	signer := types.NewPIP7Signer(ethConf.Genesis.Config.ChainID, ethConf.Genesis.Config.PIP7ChainID)

	legacyTx, _ := types.SignNewTx(key, signer, &types.LegacyTx{
		Nonce:    uint64(0),
		To:       &dad,
		Value:    big.NewInt(100),
		Gas:      50000,
		GasPrice: big.NewInt(1),
	})

	signer1 := types.NewLondonSigner(ethConf.Genesis.Config.PIP7ChainID)
	envelopTx, _ := types.SignNewTx(key, signer1, &types.AccessListTx{
		ChainID:  ethConf.Genesis.Config.PIP7ChainID,
		Nonce:    uint64(1),
		To:       &dad,
		Gas:      30000,
		GasPrice: big.NewInt(1),
		Value:    big.NewInt(50),
		AccessList: types.AccessList{{
			Address:     dad,
			StorageKeys: []common.Hash{{0}},
		}},
	})

	// Create some blocks and import them
	core.GenerateBlockChain3(params.TestChainConfig, ethBackend.BlockChain().Genesis(),
		consensus.NewFakerWithDataBase(ethBackend.ChainDb()), ethBackend.BlockChain(), 1, func(i int, b *core.BlockGen) {
			b.SetCoinbase(common.Address{1})
			b.SetActiveVersion(params.FORKVERSION_1_5_0)
			b.AddTx(legacyTx)
			b.AddTx(envelopTx)
		})
	// create gql service
	err = New(stack, ethBackend.APIBackend, []string{}, []string{})
	if err != nil {
		t.Fatalf("could not create graphql service: %v", err)
	}
}
