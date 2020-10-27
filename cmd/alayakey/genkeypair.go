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

package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
)

type outputGenkeypair struct {
	PrivateKey string
	PublicKey  string
}

var commandGenkeypair = cli.Command{
	Name:      "genkeypair",
	Usage:     "generate new private key pair",
	ArgsUsage: "[ ]",
	Description: `
Generate a new private key pair.
.
`,
	Flags: []cli.Flag{
		jsonFlag,
	},
	Action: func(ctx *cli.Context) error {
		// Check if keyfile path given and make sure it doesn't already exist.
		var privateKey *ecdsa.PrivateKey
		var err error
		// generate random.
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			utils.Fatalf("Failed to generate random private key: %v", err)
		}

		// Output some information.
		out := outputGenkeypair{
			PublicKey:  hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey)[1:]),
			PrivateKey: hex.EncodeToString(crypto.FromECDSA(privateKey)),
		}
		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			fmt.Println("PrivateKey: ", out.PrivateKey)
			fmt.Println("PublicKey : ", out.PublicKey)
		}
		return nil
	},
}
