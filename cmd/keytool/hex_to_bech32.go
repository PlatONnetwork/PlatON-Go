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
	Usage: "file bech32/hex accounts want to update to mainnet/testnet bech32 address,file like  [hex,hex...]",
}

type addressPair struct {
	Address       common.AddressOutput
	OriginAddress string
}

var commandAddressHexToBech32 = cli.Command{
	Name:      "updateaddress",
	Usage:     "update hex/bech32 address to mainnet/testnet bech32 address",
	ArgsUsage: "[<address> <address>...]",
	Description: `
update hex/bech32 address to mainnet/testnet bech32 address.
`,
	Flags: []cli.Flag{
		jsonFlag,
		HexAccountFileFlag,
	},
	Action: func(ctx *cli.Context) error {
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
		var outAddress []addressPair
		for _, account := range accounts {
			_, _, err := bech32.Decode(account)
			if err != nil {
				address := common.HexToAddress(account)
				out := addressPair{
					Address:       common.NewAddressOutput(address),
					OriginAddress: account,
				}
				outAddress = append(outAddress, out)
			} else {
				_, converted, err := bech32util.DecodeAndConvert(account)
				if err != nil {
					return err
				}
				var a common.Address
				a.SetBytes(converted)
				out := addressPair{
					Address:       common.NewAddressOutput(a),
					OriginAddress: account,
				}
				outAddress = append(outAddress, out)
			}
		}

		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(outAddress)
		} else {
			for i, address := range outAddress {
				fmt.Println("originAddress: ", address.OriginAddress)
				address.Address.Print()
				if i != len(outAddress)-1 {
					fmt.Println("---")
				}
			}
		}
		return nil
	},
}
