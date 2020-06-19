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
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cespare/cp"
)

// These tests are 'smoke tests' for the account related
// subcommands and flags.
//
// For most tests, the test files from package accounts
// are copied into a temporary keystore directory.

func tmpDatadirWithKeystore(t *testing.T) string {
	datadir := tmpdir(t)
	keystore := filepath.Join(datadir, "keystore")
	source := filepath.Join("..", "..", "accounts", "keystore", "testdata", "keystore")
	if err := cp.CopyAll(keystore, source); err != nil {
		t.Fatal(err)
	}
	return datadir
}

func TestAccountListEmpty(t *testing.T) {
	platon := runPlatON(t, "account", "list")
	platon.ExpectExit()
}

func TestAccountList(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	log.Print(datadir)
	platon := runPlatON(t, "account", "list", "--datadir", datadir)
	defer platon.ExpectExit()
	if runtime.GOOS == "windows" {
		platon.Expect(`
Account #0: {mainnet:lat10m66vy6lrlt2qfvnamwgd8rdg8vnfthcd74p32,testnet:lax10m66vy6lrlt2qfvnamwgd8rdg8vnfthczm8wl9} keystore://{{.Datadir}}\keystore\UTC--2016-03-22T12-57-55.920751759Z--7ef5a6135f1fd6a02593eedc869c6d41d934aef8
Account #1: {mainnet:lat173ngt84dryedws7kyt9hflq93zpwsey2m0wqp6,testnet:lax173ngt84dryedws7kyt9hflq93zpwsey252u004} keystore://{{.Datadir}}\keystore\aaa
Account #2: {mainnet:lat19zw5shvhw9c5en536vun6ajwzvgeq7kvh7rqmg,testnet:lax19zw5shvhw9c5en536vun6ajwzvgeq7kvcm3048} keystore://{{.Datadir}}\keystore\zzz
`)
	} else {
		platon.Expect(`
Account #0: {mainnet:lat10m66vy6lrlt2qfvnamwgd8rdg8vnfthcd74p32,testnet:lax10m66vy6lrlt2qfvnamwgd8rdg8vnfthczm8wl9} keystore://{{.Datadir}}/keystore/UTC--2016-03-22T12-57-55.920751759Z--7ef5a6135f1fd6a02593eedc869c6d41d934aef8
Account #1: {mainnet:lat173ngt84dryedws7kyt9hflq93zpwsey2m0wqp6,testnet:lax173ngt84dryedws7kyt9hflq93zpwsey252u004} keystore://{{.Datadir}}/keystore/aaa
Account #2: {mainnet:lat19zw5shvhw9c5en536vun6ajwzvgeq7kvh7rqmg,testnet:lax19zw5shvhw9c5en536vun6ajwzvgeq7kvcm3048} keystore://{{.Datadir}}/keystore/zzz
`)
	}
}

func TestAccountNew(t *testing.T) {
	platon := runPlatON(t, "account", "new", "--lightkdf")
	defer platon.ExpectExit()
	platon.Expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Repeat passphrase: {{.InputLine "foobar"}}
`)

	platon.ExpectRegexp(`main net Address: lat1[0-9a-z]{38}\nother net Address: lax1[0-9a-z]{38}\n`)
}

func TestAccountNewBadRepeat(t *testing.T) {
	platon := runPlatON(t, "account", "new", "--lightkdf")
	defer platon.ExpectExit()
	platon.Expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "something"}}
Repeat passphrase: {{.InputLine "something else"}}
Fatal: Passphrases do not match
`)
}

func TestAccountUpdate(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	platon := runPlatON(t, "account", "update",
		"--datadir", datadir, "--lightkdf",
		"lat173ngt84dryedws7kyt9hflq93zpwsey2m0wqp6")
	defer platon.ExpectExit()
	platon.Expect(`
Unlocking account lat173ngt84dryedws7kyt9hflq93zpwsey2m0wqp6 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Please give a new password. Do not forget this password.
Passphrase: {{.InputLine "foobar2"}}
Repeat passphrase: {{.InputLine "foobar2"}}
`)
}

func TestUnlockFlag(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	platon := runPlatON(t,
		"--datadir", datadir, "--ipcdisable", "--testnet", "--nat", "none", "--nodiscover", "--maxpeers", "0", "--port", "0",
		"--unlock", "lax10m66vy6lrlt2qfvnamwgd8rdg8vnfthczm8wl9",
		"js", "testdata/empty.js")
	platon.Expect(`
Unlocking account lax10m66vy6lrlt2qfvnamwgd8rdg8vnfthczm8wl9 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
`)
	platon.ExpectExit()

	wantMessages := []string{
		"Unlocked account",
		"=lax10m66vy6lrlt2qfvnamwgd8rdg8vnfthczm8wl9",
	}
	for _, m := range wantMessages {
		if !strings.Contains(platon.StderrText(), m) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagWrongPassword(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	platon := runPlatON(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--maxpeers", "0", "--port", "0", "--ipcdisable", "--testnet",
		"--unlock", "lax173ngt84dryedws7kyt9hflq93zpwsey252u004")
	defer platon.ExpectExit()
	platon.Expect(`
Unlocking account lax173ngt84dryedws7kyt9hflq93zpwsey252u004 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "wrong1"}}
Unlocking account lax173ngt84dryedws7kyt9hflq93zpwsey252u004 | Attempt 2/3
Passphrase: {{.InputLine "wrong2"}}
Unlocking account lax173ngt84dryedws7kyt9hflq93zpwsey252u004 | Attempt 3/3
Passphrase: {{.InputLine "wrong3"}}
Fatal: Failed to unlock account lax173ngt84dryedws7kyt9hflq93zpwsey252u004 (could not decrypt key with given passphrase)
`)
}

// https://github.com/ethereum/go-ethereum/issues/1785
func TestUnlockFlagMultiIndex(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	platon := runPlatON(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--maxpeers", "0", "--port", "0", "--ipcdisable", "--testnet",
		"--unlock", "0,2",
		"js", "testdata/empty.js")
	platon.Expect(`
Unlocking account 0 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Unlocking account 2 | Attempt 1/3
Passphrase: {{.InputLine "foobar"}}
`)
	platon.ExpectExit()

	wantMessages := []string{
		"Unlocked account",
		"=lax10m66vy6lrlt2qfvnamwgd8rdg8vnfthczm8wl9",
		"=lax19zw5shvhw9c5en536vun6ajwzvgeq7kvcm3048",
	}
	for _, m := range wantMessages {
		if !strings.Contains(platon.StderrText(), m) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagPasswordFile(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	platon := runPlatON(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--maxpeers", "0", "--port", "0",
		"--password", "testdata/passwords.txt", "--unlock", "0,2", "--ipcdisable", "--testnet",
		"js", "testdata/empty.js")
	platon.ExpectExit()

	wantMessages := []string{
		"Unlocked account",
		"=lax10m66vy6lrlt2qfvnamwgd8rdg8vnfthczm8wl9",
		"=lax19zw5shvhw9c5en536vun6ajwzvgeq7kvcm3048",
	}
	for _, m := range wantMessages {
		if !strings.Contains(platon.StderrText(), m) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagPasswordFileWrongPassword(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	platon := runPlatON(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--maxpeers", "0", "--port", "0", "--ipcdisable", "--testnet",
		"--password", "testdata/wrong-passwords.txt", "--unlock", "0,2")
	defer platon.ExpectExit()
	platon.Expect(`
Fatal: Failed to unlock account 0 (could not decrypt key with given passphrase)
`)
}

func TestUnlockFlagAmbiguous(t *testing.T) {
	store := filepath.Join("..", "..", "accounts", "keystore", "testdata", "dupes")
	platon := runPlatON(t,
		"--keystore", store, "--nat", "none", "--nodiscover", "--maxpeers", "0", "--port", "0", "--ipcdisable", "--testnet",
		"--unlock", "lax173ngt84dryedws7kyt9hflq93zpwsey252u004",
		"js", "testdata/empty.js")
	defer platon.ExpectExit()

	// Helper for the expect template, returns absolute keystore path.
	platon.SetTemplateFunc("keypath", func(file string) string {
		abs, _ := filepath.Abs(filepath.Join(store, file))
		return abs
	})
	platon.Expect(`
Unlocking account lax173ngt84dryedws7kyt9hflq93zpwsey252u004 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Multiple key files exist for address lax173ngt84dryedws7kyt9hflq93zpwsey252u004:
   keystore://{{keypath "1"}}
   keystore://{{keypath "2"}}
Testing your passphrase against all of them...
Your passphrase unlocked keystore://{{keypath "1"}}
In order to avoid this warning, you need to remove the following duplicate key files:
   keystore://{{keypath "2"}}
`)
	platon.ExpectExit()

	wantMessages := []string{
		"Unlocked account",
		"=lax173ngt84dryedws7kyt9hflq93zpwsey252u004",
	}
	for _, m := range wantMessages {
		if !strings.Contains(platon.StderrText(), m) {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagAmbiguousWrongPassword(t *testing.T) {
	store := filepath.Join("..", "..", "accounts", "keystore", "testdata", "dupes")
	platon := runPlatON(t,
		"--keystore", store, "--nat", "none", "--nodiscover", "--maxpeers", "0", "--port", "0", "--ipcdisable", "--testnet",
		"--unlock", "lax173ngt84dryedws7kyt9hflq93zpwsey252u004")
	defer platon.ExpectExit()

	// Helper for the expect template, returns absolute keystore path.
	platon.SetTemplateFunc("keypath", func(file string) string {
		abs, _ := filepath.Abs(filepath.Join(store, file))
		return abs
	})
	platon.Expect(`
Unlocking account lax173ngt84dryedws7kyt9hflq93zpwsey252u004 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "wrong"}}
Multiple key files exist for address lax173ngt84dryedws7kyt9hflq93zpwsey252u004:
   keystore://{{keypath "1"}}
   keystore://{{keypath "2"}}
Testing your passphrase against all of them...
Fatal: None of the listed files could be unlocked.
`)
	platon.ExpectExit()
}
