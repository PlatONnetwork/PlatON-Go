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

package console

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/internal/jsre"
	"github.com/PlatONnetwork/PlatON-Go/node"
	_ "github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	testInstance = "console-tester"
)

func init() {
	bls.Init(bls.BLS12_381)
}

// hookedPrompter implements UserPrompter to simulate use input via channels.
type hookedPrompter struct {
	scheduler chan string
}

func (p *hookedPrompter) PromptInput(prompt string) (string, error) {
	// Send the prompt to the tester
	select {
	case p.scheduler <- prompt:
	case <-time.After(time.Second):
		return "", errors.New("prompt timeout")
	}
	// Retrieve the response and feed to the console
	select {
	case input := <-p.scheduler:
		return input, nil
	case <-time.After(time.Second):
		return "", errors.New("input timeout")
	}
}

func (p *hookedPrompter) PromptPassword(prompt string) (string, error) {
	return "", errors.New("not implemented")
}
func (p *hookedPrompter) PromptConfirm(prompt string) (bool, error) {
	return false, errors.New("not implemented")
}
func (p *hookedPrompter) SetHistory(history []string)              {}
func (p *hookedPrompter) AppendHistory(command string)             {}
func (p *hookedPrompter) ClearHistory()                            {}
func (p *hookedPrompter) SetWordCompleter(completer WordCompleter) {}

// tester is a console test environment for the console tests to operate on.
type tester struct {
	workspace string
	stack     *node.Node
	ethereum  *eth.Ethereum
	console   *Console
	input     *hookedPrompter
	output    *bytes.Buffer
}

// newTester creates a test environment based on which the console can operate.
// Please ensure you call Close() on the returned tester to avoid leaks.
func newTester(t *testing.T, confOverride func(*eth.Config)) *tester {
	// Create a temporary storage for the node keys and initialize it
	workspace, err := ioutil.TempDir("", "console-tester-")
	if err != nil {
		t.Fatalf("failed to create temporary keystore: %v", err)
	}

	// Create a networkless protocol stack and start an Ethereum service within
	stack, err := node.New(&node.Config{DataDir: workspace, UseLightweightKDF: true, Name: testInstance})
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}
	snapshotdb.SetDBPathWithNode(stack.ResolvePath(snapshotdb.DBPath))
	//ethConf := &eth.Config{
	//	Genesis:   core.DeveloperGenesisBlock(15, common.Address{}),
	//}
	ethConf := &eth.DefaultConfig
	ethConf.Genesis = core.DefaultGrapeGenesisBlock()
	n, _ := discover.ParseNode("enode://73f48a69ae73b85c0a578258954936300b305cb063cbd658d680826ebc0d47cedb890f01f15df2f2e510342d16e7bf5aaf3d7be4ba05a3490de0e9663663addc@127.0.0.1:16789")

	var nodes []params.CbftNode
	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	nodes = append(nodes, params.CbftNode{Node: *n, BlsPubKey: *blsKey.GetPublicKey()})
	ethConf.Genesis.Config.Cbft = &params.CbftConfig{
		InitialNodes: nodes,
	}
	if confOverride != nil {
		confOverride(ethConf)
	}
	if err = stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {

		return eth.New(ctx, ethConf)
	}); err != nil {
		t.Fatalf("failed to register Ethereum protocol: %v", err)
	}
	// Start the node and assemble the JavaScript console around it
	if err = stack.Start(); err != nil {
		t.Fatalf("failed to start test stack: %v", err)
	}
	client, err := stack.Attach()
	if err != nil {
		t.Fatalf("failed to attach to node: %v", err)
	}
	prompter := &hookedPrompter{scheduler: make(chan string)}
	printer := new(bytes.Buffer)

	console, err := New(Config{
		DataDir:  stack.DataDir(),
		DocRoot:  "testdata",
		Client:   client,
		Prompter: prompter,
		Printer:  printer,
		Preload:  []string{"preload.js"},
	})
	if err != nil {
		t.Fatalf("failed to create JavaScript console: %v", err)
	}
	// Create the final tester and return
	var ethereum *eth.Ethereum
	stack.Service(&ethereum)

	return &tester{
		workspace: workspace,
		stack:     stack,
		ethereum:  ethereum,
		console:   console,
		input:     prompter,
		output:    printer,
	}
}

// Close cleans up any temporary data folders and held resources.
func (env *tester) Close(t *testing.T) {
	if err := env.console.Stop(false); err != nil {
		t.Errorf("failed to stop embedded console: %v", err)
	}
	if err := env.stack.Stop(); err != nil {
		t.Errorf("failed to stop embedded node: %v", err)
	}
	os.RemoveAll(env.workspace)
}

// Tests that the node lists the correct welcome message, notably that it contains
// the instance name, coinbase account, block number, data directory and supported
// console modules.
func TestWelcome(t *testing.T) {
	tester := newTester(t, nil)
	defer tester.Close(t)

	tester.console.Welcome()

	output := tester.output.String()
	if want := "Welcome"; !strings.Contains(output, want) {
		t.Fatalf("console output missing welcome message: have\n%s\nwant also %s", output, want)
	}
	if want := fmt.Sprintf("instance: %s", testInstance); !strings.Contains(output, want) {
		t.Fatalf("console output missing instance: have\n%s\nwant also %s", output, want)
	}
	if want := "at block: 0"; !strings.Contains(output, want) {
		t.Fatalf("console output missing sync status: have\n%s\nwant also %s", output, want)
	}
	if want := fmt.Sprintf("datadir: %s", tester.workspace); !strings.Contains(output, want) {
		t.Fatalf("console output missing coinbase: have\n%s\nwant also %s", output, want)
	}
}

func TestApi(t *testing.T) {
	tester := newTester(t, nil)
	defer tester.Close(t)
	fmt.Fprintf(tester.console.printer, "Welcome to the PlatON JavaScript console!\n\n")
	_, err := tester.console.jsre.Run(`
		console.log("aaaaaaa");
		console.log("instance: " + web3.version.node);
		console.log("at block: " + platon.blockNumber + " (" + new Date(1000 * platon.getBlock(platon.blockNumber).timestamp) + ")");
		console.log(" datadir: " + admin.datadir);
		console.log(" protocolVersion: " + platon.protocolVersion);
		console.log(" sync: " + platon.syncing);
		console.log("",platon.protocolVersion)
		console.log("syncing",platon.syncing)
		console.log("gasPrice",platon.gasPrice)
		console.log("accounts",platon.accounts)
		console.log("blockNumber",platon.blockNumber)
		console.log("getBalance",platon.getBalance("lax1sczumw7md5ny4f6zuaczph9utr7decvzstap70"))
		console.log("getStorageAt",platon.getStorageAt("lax1sczumw7md5ny4f6zuaczph9utr7decvzstap70"))
		console.log("getTransactionCount",platon.getTransactionCount("lax1sczumw7md5ny4f6zuaczph9utr7decvzstap70"))
		console.log("getBlockTransactionCountByHash or ByNumber",platon.getBlockTransactionCount("1234"))
		//console.log("getCode",platon.getCode("0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"))
		//adr = personal.newAccount("123456")
		//personal.unlockAccount(adr,"123456",30)
		//console.log("sign",platon.sign(adr, "0xdeadbeaf"))
		//console.log("sendTransaction",platon.sendTransaction({from:adr,to:adr,value:0,gas:0,gasPrice:0}))
		//console.log("sendRawTransaction",platon.sendRawTransaction({from:platon.accounts[0],to:platon.accounts[1],value:10,gas:88888,gasPrice:3333}))
		//console.log("call",platon.call({from:platon.accounts[0],to:platon.accounts[1],value:10,gas:88888,gasPrice:3333}))
		//console.log("estimateGas",platon.estimateGas({from:platon.accounts[0],to:platon.accounts[1],value:10,gas:88888,gasPrice:3333}))
		//console.log("getBlockByHash or number",platon.getBlock("0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331"))
		//console.log("getTransactionByHash",platon.getTransaction("0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331"))
		//console.log("getTransactionByBlockHashAndIndex",platon.getTransactionFromBlock(["0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b", "1"]))
		//console.log("getTransactionReceipt",platon.getTransactionReceipt())
		//console.log("newFilter",platon.newFilter())
		//console.log("newBlockFilter",platon.newBlockFilter())
		//console.log("newPendingTransactionFilter",platon.newPendingTransactionFilter())
		//console.log("uninstallFilter",platon.uninstallFilter())
		//console.log("getFilterChanges",platon.getFilterChanges())
		//console.log("getFilterLogs",platon.getFilterLogs())
		//console.log("getLogs",platon.getLogs({"topics":["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"]}))
		//console.log("signTransaction",platon.signTransaction())
		//console.log("test personal",personal.openWallet("adad"))
	`)
	if err != nil {
		t.Error(err)
	}
	t.Log(tester.output.String())
}

// Tests that JavaScript statement evaluation works as intended.
func TestEvaluate(t *testing.T) {
	tester := newTester(t, nil)
	defer tester.Close(t)

	tester.console.Evaluate("2 + 2")
	if output := tester.output.String(); !strings.Contains(output, "4") {
		t.Fatalf("statement evaluation failed: have %s, want %s", output, "4")
	}
}

// Tests that the console can be used in interactive mode.
func TestInteractive(t *testing.T) {
	// Create a tester and run an interactive console in the background
	tester := newTester(t, nil)
	defer tester.Close(t)

	go tester.console.Interactive()

	// Wait for a prompt and send a statement back
	select {
	case <-tester.input.scheduler:
	case <-time.After(time.Second):
		t.Fatalf("initial prompt timeout")
	}
	select {
	case tester.input.scheduler <- "2+2":
	case <-time.After(time.Second):
		t.Fatalf("input feedback timeout")
	}
	// Wait for the second prompt and ensure first statement was evaluated
	select {
	case <-tester.input.scheduler:
	case <-time.After(time.Second):
		t.Fatalf("secondary prompt timeout")
	}
	if output := tester.output.String(); !strings.Contains(output, "4") {
		t.Fatalf("statement evaluation failed: have %s, want %s", output, "4")
	}
}

// Tests that preloaded JavaScript files have been executed before user is given
// input.
func TestPreload(t *testing.T) {
	tester := newTester(t, nil)
	defer tester.Close(t)

	tester.console.Evaluate("preloaded")
	if output := tester.output.String(); !strings.Contains(output, "some-preloaded-string") {
		t.Fatalf("preloaded variable missing: have %s, want %s", output, "some-preloaded-string")
	}
}

// Tests that JavaScript scripts can be executes from the configured asset path.
func TestExecute(t *testing.T) {
	tester := newTester(t, nil)
	defer tester.Close(t)

	tester.console.Execute("exec.js")

	tester.console.Evaluate("execed")
	if output := tester.output.String(); !strings.Contains(output, "some-executed-string") {
		t.Fatalf("execed variable missing: have %s, want %s", output, "some-executed-string")
	}
}

// Tests that the JavaScript objects returned by statement executions are properly
// pretty printed instead of just displaying "[object]".
func TestPrettyPrint(t *testing.T) {
	tester := newTester(t, nil)
	defer tester.Close(t)

	tester.console.Evaluate("obj = {int: 1, string: 'two', list: [3, 3, 3], obj: {null: null, func: function(){}}}")

	// Define some specially formatted fields
	var (
		one   = jsre.NumberColor("1")
		two   = jsre.StringColor("\"two\"")
		three = jsre.NumberColor("3")
		null  = jsre.SpecialColor("null")
		fun   = jsre.FunctionColor("function()")
	)
	// Assemble the actual output we're after and verify
	want := `{
  int: ` + one + `,
  list: [` + three + `, ` + three + `, ` + three + `],
  obj: {
    null: ` + null + `,
    func: ` + fun + `
  },
  string: ` + two + `
}
`
	if output := tester.output.String(); output != want {
		t.Fatalf("pretty print mismatch: have %s, want %s", output, want)
	}
}

// Tests that the JavaScript exceptions are properly formatted and colored.
func TestPrettyError(t *testing.T) {
	tester := newTester(t, nil)
	defer tester.Close(t)
	tester.console.Evaluate("throw 'hello'")

	want := jsre.ErrorColor("hello") + "\n"
	if output := tester.output.String(); output != want {
		t.Fatalf("pretty error mismatch: have %s, want %s", output, want)
	}
}

// Tests that tests if the number of indents for JS input is calculated correct.
func TestIndenting(t *testing.T) {
	testCases := []struct {
		input               string
		expectedIndentCount int
	}{
		{`var a = 1;`, 0},
		{`"some string"`, 0},
		{`"some string with (parenthesis`, 0},
		{`"some string with newline
		("`, 0},
		{`function v(a,b) {}`, 0},
		{`function f(a,b) { var str = "asd("; };`, 0},
		{`function f(a) {`, 1},
		{`function f(a, function(b) {`, 2},
		{`function f(a, function(b) {
		     var str = "a)}";
		  });`, 0},
		{`function f(a,b) {
		   var str = "a{b(" + a, ", " + b;
		   }`, 0},
		{`var str = "\"{"`, 0},
		{`var str = "'("`, 0},
		{`var str = "\\{"`, 0},
		{`var str = "\\\\{"`, 0},
		{`var str = 'a"{`, 0},
		{`var obj = {`, 1},
		{`var obj = { {a:1`, 2},
		{`var obj = { {a:1}`, 1},
		{`var obj = { {a:1}, b:2}`, 0},
		{`var obj = {}`, 0},
		{`var obj = {
			a: 1, b: 2
		}`, 0},
		{`var test = }`, -1},
		{`var str = "a\""; var obj = {`, 1},
	}

	for i, tt := range testCases {
		counted := countIndents(tt.input)
		if counted != tt.expectedIndentCount {
			t.Errorf("test %d: invalid indenting: have %d, want %d", i, counted, tt.expectedIndentCount)
		}
	}
}
