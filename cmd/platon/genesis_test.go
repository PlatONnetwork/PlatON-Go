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
    "alloc":{
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp7pn3ep":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp7pn3ep":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v":{
            "balance":"200000000000000000000000000"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyva9ztf":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq93t3hkm":{
            "balance":"0"
        },
        "lat1vr8v48qjjrh9dwvdfctqauz98a7yp5se3mm5yk":{
            "balance":"8050000000000000000000000000"
        },
        "lat12klaf9rjl4qjz929kqt382wr49a00zc9ng995w":{
            "balance":"2000000000000000000000000000"
        }
    },
    "economicModel":{
        "common":{
            "maxEpochMinutes":4,
            "maxConsensusVals":4,
            "additionalCycleTime":28
        },
        "staking":{
            "stakeThreshold": 1000000000000000000000000,
            "operatingThreshold": 10000000000000000000,
            "maxValidators": 30,
            "unStakeFreezeDuration": 2
        },
        "slashing":{
           "slashFractionDuplicateSign": 100,
           "duplicateSignReportReward": 50,
           "maxEvidenceAge":1,
           "slashBlocksReward":20
        },
         "gov": {
            "versionProposalVoteDurationSeconds": 160,
            "versionProposalSupportRate": 6670,
            "textProposalVoteDurationSeconds": 160,
            "textProposalVoteRate": 5000,
            "textProposalSupportRate": 6670,          
            "cancelProposalVoteRate": 5000,
            "cancelProposalSupportRate": 6670,
            "paramProposalVoteDurationSeconds": 160,
            "paramProposalVoteRate": 5000,
            "paramProposalSupportRate": 6670      
        },
        "reward":{
            "newBlockRate": 50,
            "platonFoundationYear": 10 
        },
        "innerAcc":{
            "platonFundAccount": "lat1fyeszufxwxk62p46djncj86rd553skppy4qgz4",
            "platonFundBalance": 0,
            "cdfAccount": "lat1c8enpvs5v6974shxgxxav5dsn36e5jl4r0hwhh",
            "cdfBalance": 331811981000000000000000000
        }
    },
    "coinbase":"lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a",
    "extraData":"",
    "gasLimit":"0x2fefd8",
    "nonce":"0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
    "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp":"0x00",
    "config":{
        "cbft":{
            "initialNodes":[
                {
                    "node":"enode://4fcc251cf6bf3ea53a748971a223f5676225ee4380b65c7889a2b491e1551d45fe9fcc19c6af54dcf0d5323b5aa8ee1d919791695082bae1f86dd282dba4150f@0.0.0.0:16789",
                    "blsPubKey":"d341a0c485c9ec00cecf7ea16323c547900f6a1bacb9daacb00c2b8bacee631f75d5d31b75814b7f1ae3a4e18b71c617bc2f230daa0c893746ed87b08b2df93ca4ddde2816b3ac410b9980bcc048521562a3b2d00e900fd777d3cf88ce678719"
                }
            ],
            "amount":10,
			"period":10000,
            "validatorMode":"ppos"
        }
    }
}`,
		query:  "platon.getBlock(0).nonce",
		result: "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	},
	//Genesis file with only cbft config
	{
		genesis: `{
    "alloc":{
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp7pn3ep":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp7pn3ep":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v":{
            "balance":"200000000000000000000000000"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyva9ztf":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq93t3hkm":{
            "balance":"0"
        },
        "lat1vr8v48qjjrh9dwvdfctqauz98a7yp5se3mm5yk":{
            "balance":"8050000000000000000000000000"
        },
        "lat12klaf9rjl4qjz929kqt382wr49a00zc9ng995w":{
            "balance":"2000000000000000000000000000"
        }
    },
    "economicModel":{
        "common":{
            "maxEpochMinutes":4,
            "maxConsensusVals":4,
            "additionalCycleTime":28
        },
        "staking":{
            "stakeThreshold": 1000000000000000000000000,
            "operatingThreshold": 10000000000000000000,
            "maxValidators": 30,
            "unStakeFreezeDuration": 2
        },
        "slashing":{
           "slashFractionDuplicateSign": 100,
           "duplicateSignReportReward": 50,
           "maxEvidenceAge":1,
           "slashBlocksReward":20
        },
         "gov": {
            "versionProposalVoteDurationSeconds": 160,
            "versionProposalSupportRate": 6670,
            "textProposalVoteDurationSeconds": 160,
            "textProposalVoteRate": 5000,
            "textProposalSupportRate": 6670,          
            "cancelProposalVoteRate": 5000,
            "cancelProposalSupportRate": 6670,
            "paramProposalVoteDurationSeconds": 160,
            "paramProposalVoteRate": 5000,
            "paramProposalSupportRate": 6670      
        },
        "reward":{
            "newBlockRate": 50,
            "platonFoundationYear": 10 
        },
        "innerAcc":{
            "platonFundAccount": "lat1fyeszufxwxk62p46djncj86rd553skppy4qgz4",
            "platonFundBalance": 0,
            "cdfAccount": "lat1c8enpvs5v6974shxgxxav5dsn36e5jl4r0hwhh",
            "cdfBalance": 331811981000000000000000000
        }
    },
    "coinbase":"lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a",
    "extraData":"",
    "gasLimit":"0x2fefd8",
    "nonce":"0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
    "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp":"0x00",
    "config":{
        "cbft":{
            "initialNodes":[
                {
                    "node":"enode://4fcc251cf6bf3ea53a748971a223f5676225ee4380b65c7889a2b491e1551d45fe9fcc19c6af54dcf0d5323b5aa8ee1d919791695082bae1f86dd282dba4150f@0.0.0.0:16789",
                    "blsPubKey":"d341a0c485c9ec00cecf7ea16323c547900f6a1bacb9daacb00c2b8bacee631f75d5d31b75814b7f1ae3a4e18b71c617bc2f230daa0c893746ed87b08b2df93ca4ddde2816b3ac410b9980bcc048521562a3b2d00e900fd777d3cf88ce678719"
                }
            ],
            "amount":10,
			"period":10000,
            "validatorMode":"ppos"
        }
    }
}`,
		query:  "platon.getBlock(0).nonce",
		result: "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	},
	//Genesis file with specific chain configurations
	{
		genesis: `{
    "alloc":{
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp7pn3ep":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzsjx8h7":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v":{
            "balance":"200000000000000000000000000"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyva9ztf":{
            "balance":"0"
        },
        "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq93t3hkm":{
            "balance":"0"
        },
        "lat1vr8v48qjjrh9dwvdfctqauz98a7yp5se3mm5yk":{
            "balance":"8050000000000000000000000000"
        },
        "lat12klaf9rjl4qjz929kqt382wr49a00zc9ng995w":{
            "balance":"2000000000000000000000000000"
        }
    },
    "economicModel":{
        "common":{
            "maxEpochMinutes":4,
            "maxConsensusVals":4,
            "additionalCycleTime":28
        },
        "staking":{
            "stakeThreshold": 1000000000000000000000000,
            "operatingThreshold": 10000000000000000000,
            "maxValidators": 30,
            "unStakeFreezeDuration": 2
        },
        "slashing":{
           "slashFractionDuplicateSign": 100,
           "duplicateSignReportReward": 50,
           "maxEvidenceAge":1,
           "slashBlocksReward":20
        },
         "gov": {
            "versionProposalVoteDurationSeconds": 160,
            "versionProposalSupportRate": 6670,
            "textProposalVoteDurationSeconds": 160,
            "textProposalVoteRate": 5000,
            "textProposalSupportRate": 6670,          
            "cancelProposalVoteRate": 5000,
            "cancelProposalSupportRate": 6670,
            "paramProposalVoteDurationSeconds": 160,
            "paramProposalVoteRate": 5000,
            "paramProposalSupportRate": 6670      
        },
        "reward":{
            "newBlockRate": 50,
            "platonFoundationYear": 10 
        },
        "innerAcc":{
            "platonFundAccount": "lat1fyeszufxwxk62p46djncj86rd553skppy4qgz4",
            "platonFundBalance": 0,
            "cdfAccount": "lat1c8enpvs5v6974shxgxxav5dsn36e5jl4r0hwhh",
            "cdfBalance": 331811981000000000000000000
        }
    },
    "coinbase":"lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a",
    "extraData":"",
    "gasLimit":"0x2fefd8",
    "nonce":"0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
    "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp":"0x00",
    "config":{
        "chainId":101,
        "eip155Block":0,
        "interpreter":"wasm",
        "cbft":{
            "initialNodes":[
                {
                    "node":"enode://4fcc251cf6bf3ea53a748971a223f5676225ee4380b65c7889a2b491e1551d45fe9fcc19c6af54dcf0d5323b5aa8ee1d919791695082bae1f86dd282dba4150f@0.0.0.0:16789",
                    "blsPubKey":"d341a0c485c9ec00cecf7ea16323c547900f6a1bacb9daacb00c2b8bacee631f75d5d31b75814b7f1ae3a4e18b71c617bc2f230daa0c893746ed87b08b2df93ca4ddde2816b3ac410b9980bcc048521562a3b2d00e900fd777d3cf88ce678719"
                }
            ],
            "amount":10,
			"period":10000,
            "validatorMode":"ppos"
        }
    }
}`,
		query:  "platon.getBlock(0).nonce",
		result: "0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	},
}

// Tests that initializing PlatON with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		// Create a temporary data directory to use and inspect later
		datadir := t.TempDir()

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := os.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}
		runPlatON(t, "--datadir", datadir, "init", json).WaitExit()

		// Query the custom genesis block
		platon := runPlatON(t,
			"--datadir", datadir, "--maxpeers", "60", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable", "--testnet",
			"--exec", tt.query, "console")
		t.Log("testi", i)
		platon.ExpectRegexp(tt.result)
		platon.ExpectExit()
	}
}
