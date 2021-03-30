package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/btcsuite/btcutil/bech32"

	"github.com/PlatONnetwork/PlatON-Go/common/bech32util"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"

	"gopkg.in/urfave/cli.v1"
)

var HexAccountFileFlag = cli.StringFlag{
	Name:  "hexAddressFile",
	Usage: "file bech32/hex accounts want to update to  bech32 address,file like  [hex,hex...]",
}

type addressPair struct {
	Address       string
	OriginAddress string
	Hex           string
}

var commandAddressHexToBech32 = cli.Command{
	Name:      "updateaddress",
	Usage:     "update hex/bech32 address to  bech32 address",
	ArgsUsage: "[<address> <address>...]",
	Description: `
update hex/bech32 address to  bech32 address.
`,
	Flags: []cli.Flag{
		jsonFlag,
		HexAccountFileFlag,
		utils.AddressHRPFlag,
	},
	Action: func(ctx *cli.Context) error {
		hrp := ctx.String(utils.AddressHRPFlag.Name)
		if err := common.SetAddressHRP(hrp); err != nil {
			return err
		}

		var accounts []string
		if ctx.IsSet(HexAccountFileFlag.Name) {
			accountPath := ctx.String(HexAccountFileFlag.Name)
			accountjson, err := ioutil.ReadFile(accountPath)
			if err != nil {
				utils.Fatalf("Failed to read the keyfile at '%s': %v", accountPath, err)
			}
			if err := json.Unmarshal(accountjson, &accounts); err != nil {
				utils.Fatalf("Failed to json decode '%s': %v", accountPath, err)
			}
		} else {
			for _, add := range ctx.Args() {
				if add == "" {
					utils.Fatalf("the account can't be nil")
				}
				accounts = append(accounts, add)
			}
		}
		for i, account := range accounts {
			_, _, err := bech32.Decode(account)
			var out addressPair
			var address common.Address
			if err != nil {
				address = common.HexToAddress(account)
			} else {
				_, converted, err := bech32util.DecodeAndConvert(account)
				if err != nil {
					return err
				}
				address.SetBytes(converted)

			}
			out = addressPair{
				Address:       address.String(),
				OriginAddress: account,
				Hex:           address.Hex(),
			}
			if ctx.Bool(jsonFlag.Name) {
				mustPrintJSON(out)
			} else {
				fmt.Printf("origin: %s\n", out.OriginAddress)
				fmt.Printf("bech32: %s\n", out.Address)
				fmt.Printf("hex: %s\n", out.Hex)
				if i != len(accounts)-1 {
					fmt.Println("---")
				}
			}
		}

		return nil
	},
}
