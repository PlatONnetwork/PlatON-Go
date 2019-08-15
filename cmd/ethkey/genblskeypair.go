package main

import (
	"encoding/hex"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"gopkg.in/urfave/cli.v1"
)

type outputGenblskeypair struct {
	PrivateKey string
	PublicKey  string
}

var commandGenblskeypair = cli.Command{
	Name:      "genblskeypair",
	Usage:     "generate new bls private key pair",
	ArgsUsage: "[  ]",
	Description: `
Generate a new bls private key pair.
`,
	Flags: []cli.Flag{
		jsonFlag,
	},
	Action: func(ctx *cli.Context) error {
		bls.Init(bls.CurveFp254BNb)
		var privateKey bls.SecretKey
		privateKey.SetByCSPRNG()
		pubKey := privateKey.GetPublicKey()
		out := outputGenblskeypair{
			PrivateKey: hex.EncodeToString(privateKey.GetLittleEndian()),
			PublicKey:  hex.EncodeToString(pubKey.Serialize()),
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
