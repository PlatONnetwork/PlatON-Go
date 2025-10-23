// Copyright 2021 The PlatON Network Authors
// This file is part of PlatON-Go.
//
// PlatON-Go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// PlatON-Go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PlatON-Go. If not, see <http://www.gnu.org/licenses/>.

package ppos

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	platon "github.com/PlatONnetwork/PlatON-Go"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/ethclient"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func CallPPosContract(client *ethclient.Client, funcType uint16, params ...interface{}) ([]byte, error) {
	send, to := EncodePPOS(funcType, params...)
	var msg platon.CallMsg
	msg.Data = send
	msg.To = &to
	return client.CallContract(context.Background(), msg, nil)
}

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	To   *common.Address // the destination contract (nil for contract creation)
	Data hexutil.Bytes   // input data, usually an ABI-encoded contract method invocation
}

func BuildPPosContract(funcType uint16, params ...interface{}) ([]byte, error) {
	send, to := EncodePPOS(funcType, params...)
	var msg CallMsg
	msg.Data = send
	msg.To = &to
	return json.Marshal(msg)
}

func EncodePPOS(funcType uint16, params ...interface{}) ([]byte, common.Address) {
	par := buildParams(funcType, params...)
	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, par)
	if err != nil {
		panic(fmt.Errorf("encode rlp data fail: %v", err))
	}
	return buf.Bytes(), funcTypeToContractAddress(funcType)
}

func buildParams(funcType uint16, params ...interface{}) [][]byte {
	var res [][]byte
	res = make([][]byte, 0)
	fnType, _ := rlp.EncodeToBytes(funcType)
	res = append(res, fnType)
	for _, param := range params {
		val, err := rlp.EncodeToBytes(param)
		if err != nil {
			panic(err)
		}
		res = append(res, val)
	}
	return res
}

func funcTypeToContractAddress(funcType uint16) common.Address {
	toadd := common.ZeroAddr
	switch {
	case 0 < funcType && funcType < 2000:
		toadd = vm.StakingContractAddr
	case funcType >= 2000 && funcType < 3000:
		toadd = vm.GovContractAddr
	case funcType >= 3000 && funcType < 4000:
		toadd = vm.SlashingContractAddr
	case funcType >= 4000 && funcType < 5000:
		toadd = vm.RestrictingContractAddr
	case funcType >= 5000 && funcType < 6000:
		toadd = vm.DelegateRewardPoolAddr
	}
	return toadd
}

func netCheck(context *cli.Context) error {
	hrp := context.String(addressHRPFlag.Name)
	if err := common.SetAddressHRP(hrp); err != nil {
		return err
	}
	return nil
}

func query(c *cli.Context, funcType uint16, params ...interface{}) error {
	url := c.String(rpcUrlFlag.Name)
	if url == "" {
		return errors.New("rpc url not set")
	}
	if c.Bool(jsonFlag.Name) {
		res, err := BuildPPosContract(funcType, params...)
		if err != nil {
			return err
		}
		fmt.Println(string(res))
		return nil
	} else {
		client, err := ethclient.Dial(url)
		if err != nil {
			return err
		}
		res, err := CallPPosContract(client, funcType, params...)
		if err != nil {
			return err
		}
		fmt.Println(string(res))
		return nil
	}
}
