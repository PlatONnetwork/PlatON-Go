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
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/params"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 miner:1.0 net:1.0 personal:1.0 platon:1.0 rpc:1.0 txpool:1.0 web3:1.0"
	httpAPIs = "net:1.0 platon:1.0 rpc:1.0 web3:1.0"
)

// Tests that a node embedded within a console can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	datadir := tmpdir(t)
	defer os.RemoveAll(datadir)
	platon := runPlatON(t,
		"--datadir", datadir, "--port", "0", "--ipcdisable", "--testnet", "--maxpeers", "0", "--nodiscover", "--nat", "none", "console")

	// Gather all the infos the welcome message needs to contain
	platon.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	platon.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	platon.SetTemplateFunc("gover", runtime.Version)
	platon.SetTemplateFunc("gethver", func() string { return params.VersionWithMeta })
	platon.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	platon.SetTemplateFunc("apis", func() string { return ipcAPIs })

	// Verify the actual welcome message to the required template
	platon.Expect(`
Welcome to the PlatON JavaScript console!

instance: PlatONnetwork/v{{gethver}}/{{goos}}-{{goarch}}/{{gover}}
at block: 0 ({{niltime}})
 datadir: {{.Datadir}}
 modules: {{apis}}

> {{.InputLine "exit"}}
`)
	platon.ExpectExit()
}

// Tests that a console can be attached to a running node via various means.
func TestIPCAttachWelcome(t *testing.T) {
	// Configure the instance for IPC attachement
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\platon` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "platon.ipc")
	}
	platon := runPlatON(t,
		"--port", "0", "--testnet", "--maxpeers", "0", "--nodiscover", "--nat", "none", "--ipcpath", ipc)

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, platon, "ipc:"+ipc, ipcAPIs)

	platon.Interrupt()
	platon.ExpectExit()
}

func TestHTTPAttachWelcome(t *testing.T) {
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	platon := runPlatON(t,
		"--port", "0", "--ipcdisable", "--testnet", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--rpc", "--rpcport", port)

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, platon, "http://localhost:"+port, httpAPIs)

	platon.Interrupt()
	platon.ExpectExit()
}

func TestWSAttachWelcome(t *testing.T) {
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P

	platon := runPlatON(t,
		"--port", "0", "--ipcdisable", "--testnet", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--ws", "--wsport", port /*, "--testnet"*/)

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, platon, "ws://localhost:"+port, httpAPIs)

	platon.Interrupt()
	platon.ExpectExit()
}

func testAttachWelcome(t *testing.T, platon *testplaton, endpoint, apis string) {
	// Attach to a running platon note and terminate immediately
	attach := runPlatON(t, "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("gethver", func() string { return params.VersionWithMeta })
	attach.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return platon.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

	// Verify the actual welcome message to the required template
	attach.Expect(`
Welcome to the PlatON JavaScript console!

instance: PlatONnetwork/v{{gethver}}/{{goos}}-{{goarch}}/{{gover}}
at block: 0 ({{niltime}}){{if ipc}}
 datadir: {{datadir}}{{end}}
 modules: {{apis}}

> {{.InputLine "exit" }}
`)
	attach.ExpectExit()
}

// trulyRandInt generates a crypto random integer used by the console tests to
// not clash network ports with other tests running cocurrently.
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}
