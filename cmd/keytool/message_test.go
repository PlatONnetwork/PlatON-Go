// Copyright 2018 The go-ethereum Authors
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

func TestMessageSignVerify(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "keytool-test")
	if err != nil {
		t.Fatal("Can't create temporary directory:", err)
	}
	defer os.RemoveAll(tmpdir)

	keyfile := filepath.Join(tmpdir, "the-keyfile")
	message := "test message"

	// Create the key.
	generate := runKeytool(t, "generate", keyfile)
	generate.Expect(`
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Repeat passphrase: {{.InputLine "foobar"}}
`)
	_, matches := generate.ExpectRegexp(`main net Address: (lat1[0-9a-z]{38})\nother net Address: lax1[0-9a-z]{38}\n`)

	address := matches[1]
	generate.ExpectExit()

	// Sign a message.
	sign := runKeytool(t, "signmessage", keyfile, message)
	sign.Expect(`
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
`)
	_, matches = sign.ExpectRegexp(`Signature: ([0-9a-f]+)\n`)
	signature := matches[1]
	sign.ExpectExit()
	// Verify the message.
	verify := runKeytool(t, "verifymessage", address, signature, message)
	_, matches = verify.ExpectRegexp(`
Signature verification successful!
Recovered public key: [0-9a-f]+
Recovered main net address: (lat1[0-9a-z]{38})
Recovered test net address: (lax1[0-9a-z]{38})
`)
	recovered := matches[1]
	verify.ExpectExit()

	if recovered != address {
		t.Error("recovered address doesn't match generated key")
	}
}
