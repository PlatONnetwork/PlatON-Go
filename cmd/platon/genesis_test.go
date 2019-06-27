// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var customGenesisTests = []struct {
	genesis string
	query   string
	result  string
}{
	// Plain genesis file without anything extra
	{
		genesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00"
		}`,
		query:  "platon.getBlock(0).nonce",
		result: "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	},
	//Genesis file with only cbft config
	{
		genesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {
               "cbft": {
                       "initialNodes":  [
                          "enode://4fcc251cf6bf3ea53a748971a223f5676225ee4380b65c7889a2b491e1551d45fe9fcc19c6af54dcf0d5323b5aa8ee1d919791695082bae1f86dd282dba4150f@127.0.0.1:16701",
                          "enode://99e82e36db41e81366f644e14943bed70f03494d744a5d4f983387c2128e3fb5f2b8fa1b6555ea8eab81dd96de71cfda11ea8eb5310cefabc34357229e880a00@127.0.0.1:16702",
                          "enode://2f92a6719fb214667ebc85a12f738dae4d9dfd3b02be251512ab3bc1b240f92a58badc62774e1552f59d97fd5b52be8d182ed941db6f9954be120e680b531adf@127.0.0.1:16703",
                          "enode://bf9752bb531b9df04dcf869b6b19be56235a1771b19df81a6f874404db903ee360f41925ccfe3e12ca3020e91d10549b745507cd75a12cd327f52972fcc73ce3@127.0.0.1:16704"
                         ]
               }}
		}`,
		query:  "platon.getBlock(0).nonce",
		result: "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	},
	//Genesis file with specific chain configurations
	{
		genesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {
               "chainId": 101,
               "homesteadBlock": 1,
               "eip150Block": 2,
               "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
               "eip155Block": 3,
               "eip158Block": 3,
               "byzantiumBlock": 4,
               "interpreter":"wasm",
               "cbft": {
                       "initialNodes":  [
                          "enode://4fcc251cf6bf3ea53a748971a223f5676225ee4380b65c7889a2b491e1551d45fe9fcc19c6af54dcf0d5323b5aa8ee1d919791695082bae1f86dd282dba4150f@127.0.0.1:16701",
                          "enode://99e82e36db41e81366f644e14943bed70f03494d744a5d4f983387c2128e3fb5f2b8fa1b6555ea8eab81dd96de71cfda11ea8eb5310cefabc34357229e880a00@127.0.0.1:16702",
                          "enode://2f92a6719fb214667ebc85a12f738dae4d9dfd3b02be251512ab3bc1b240f92a58badc62774e1552f59d97fd5b52be8d182ed941db6f9954be120e680b531adf@127.0.0.1:16703",
                          "enode://bf9752bb531b9df04dcf869b6b19be56235a1771b19df81a6f874404db903ee360f41925ccfe3e12ca3020e91d10549b745507cd75a12cd327f52972fcc73ce3@127.0.0.1:16704"
                         ]
               }
			}
		}`,
		query:  "platon.getBlock(0).nonce",
		result: "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	},
}

// Tests that initializing Geth with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		// Create a temporary data directory to use and inspect later
		datadir := tmpdir(t)
		defer os.RemoveAll(datadir)

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}
		runGeth(t, "--datadir", datadir, "init", json).WaitExit()

		// Query the custom genesis block
		geth := runGeth(t,
			"--datadir", datadir, "--maxpeers", "0", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable",
			"--exec", tt.query, "console")
		geth.ExpectRegexp(tt.result)
		geth.ExpectExit()
	}
}
