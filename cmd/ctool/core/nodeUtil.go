// Copyright 2018-2020 The PlatON Network Authors
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

package core

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/pkg/reexec"

	"github.com/PlatONnetwork/PlatON-Go/internal/cmdtest"
	"github.com/stretchr/testify/assert"
)

type testPlatON struct {
	*cmdtest.TestCmd

	// template variables for expect
	Datadir string
}

var (
	//todo
	from = `0xc1f330b214668beac2e6418dd651b09c759a4bf5`

	to = `0x1000000000000000000000000000000000000007`

	genesis = `{
    "alloc":{
        "1000000000000000000000000000000000000003":{
            "balance":"200000000000000000000000000"
        },
        "c1f330b214668beac2e6418dd651b09c759a4bf5":{
            "balance":"8050000000000000000000000000"
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
            "platonFundAccount": "0x493301712671ada506ba6ca7891f436d29185821",
            "platonFundBalance": 0,
            "cdfAccount": "0xc1f330b214668beac2e6418dd651b09c759a4bf5",
            "cdfBalance": 331811981000000000000000000
        }
    },
    "coinbase":"0x0000000000000000000000000000000000000000",
    "extraData":"",
    "gasLimit":"0x2fefd8",
    "nonce":"0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
    "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp":"0x00",
    "config":{
		"chainId": 100,
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
}`
)

var (
	dir, _      = os.Getwd()
	abiFilePath = "../test/contracta.cpp.abi.json"
	//configPath  = "../config.json"
	configPath = filepath.Join(dir, "../config.json")
	pkFilePath = "../test/privateKeys.txt"
)

func TestMain(m *testing.M) {
	// check if we have been reexec'd
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

func parseConfig(t *testing.T) {
	err := parseConfigJson(configPath)
	assert.Nil(t, err, fmt.Sprintf("%v", err))
}

func prepare(t *testing.T) (*testPlatON, string) {
	parseConfig(t)
	datadir := tmpdir(t)
	json := filepath.Join(datadir, "genesis.json")
	err := ioutil.WriteFile(json, []byte(genesis), 0600)

	assert.Nil(t, err, fmt.Sprintf("failed to write genesis file: %v", err))

	runPlatON(t, "--datadir", datadir, "init", json).WaitExit()

	//time.Sleep(2 * time.Second)

	port := strings.Split(config.Url, ":")[2] // http://localhost:6789
	platon := runPlatON(t,
		"--datadir", datadir, "--port", "0", "--nodiscover", "--nat", "none",
		"--rpc", "--rpcaddr", "0.0.0.0", "--rpcport", port, "--rpcapi", "txpool,platon,net,web3,miner,admin,personal,version")

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open

	unlock := JsonParam{
		Jsonrpc: "2.0",
		Method:  "personal_unlockAccount",
		// {"method": "personal_unlockAccount", "params": [account, pwd, expire]}
		Params: []interface{}{from, "123456", 2222222},
		Id:     1,
	}

	// unlock
	_, e := HttpPost(unlock)

	assert.Nil(t, e, fmt.Sprintf("test http post error: %v", e))
	return platon, datadir
}

func clean(platon *testPlatON, datadir string) {

	platon.Interrupt()
	platon.ExpectExit()
	os.RemoveAll(datadir)
}

func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}

func tmpdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "platon-test")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func runPlatON(t *testing.T, args ...string) *testPlatON {
	tt := &testPlatON{}
	tt.TestCmd = cmdtest.NewTestCmd(t, tt)
	for i, arg := range args {
		switch {
		case arg == "-datadir" || arg == "--datadir":
			if i < len(args)-1 {
				tt.Datadir = args[i+1]
			}
		}
	}
	if tt.Datadir == "" {
		tt.Datadir = tmpdir(t)
		tt.Cleanup = func() { os.RemoveAll(tt.Datadir) }
		args = append([]string{"-datadir", tt.Datadir}, args...)
		// Remove the temporary datadir if something fails below.
		defer func() {
			if t.Failed() {
				tt.Cleanup()
			}
		}()
	}
	t.Log("run platon args: ", strings.Join(args, " "))

	tt.Run("platon-test", args...)

	return tt
}
