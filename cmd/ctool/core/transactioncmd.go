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
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"gopkg.in/urfave/cli.v1"
)

var (
	SendTransactionCmd = cli.Command{
		Name:   "sendTransaction",
		Usage:  "send a transaction",
		Action: sendTransactionCmd,
		Flags:  sendTransactionCmdFlags,
	}
	SendRawTransactionCmd = cli.Command{
		Name:   "sendRawTransaction",
		Usage:  "send a raw transaction",
		Action: sendRawTransactionCmd,
		Flags:  sendRawTransactionCmdFlags,
	}
	GetTxReceiptCmd = cli.Command{
		Name:   "getTxReceipt",
		Usage:  "get transaction receipt by hash",
		Action: getTxReceiptCmd,
		Flags:  getTxReceiptCmdFlags,
	}
)

func getTxReceiptCmd(c *cli.Context) {
	hash := c.String(TransactionHashFlag.Name)
	parseConfigJson(c.String(ConfigPathFlag.Name))
	GetTxReceipt(hash)
}

func GetTxReceipt(txHash string) (Receipt, error) {
	var receipt = Receipt{}
	res, _ := Send([]string{txHash}, "platon_getTransactionReceipt")
	e := json.Unmarshal([]byte(res), &receipt)
	if e != nil {
		panic(fmt.Sprintf("parse get receipt result error ! \n %s", e.Error()))
	}

	if receipt.Result.BlockHash == "" {
		panic("no receipt found")
	}
	out, _ := json.MarshalIndent(receipt, "", "  ")
	fmt.Println(string(out))
	return receipt, nil
}

func sendTransactionCmd(c *cli.Context) error {
	from := c.String(TxFromFlag.Name)
	to := c.String(TxToFlag.Name)
	value := c.String(TransferValueFlag.Name)
	parseConfigJson(c.String(ConfigPathFlag.Name))

	hash, err := SendTransaction(from, to, value)
	if err != nil {
		utils.Fatalf("Send transaction error: %v", err)
	}

	fmt.Printf("tx hash: %s", hash)
	return nil
}

func sendRawTransactionCmd(c *cli.Context) error {
	from := c.String(TxFromFlag.Name)
	to := c.String(TxToFlag.Name)
	value := c.String(TransferValueFlag.Name)
	pkFile := c.String(PKFilePathFlag.Name)

	parseConfigJson(c.String(ConfigPathFlag.Name))

	hash, err := SendRawTransaction(from, to, value, pkFile)
	if err != nil {
		utils.Fatalf("Send transaction error: %v", err)
	}

	fmt.Printf("tx hash: %s", hash)
	return nil
}

func SendTransaction(from, to, value string) (string, error) {
	var tx TxParams
	if from == "" {
		from = config.From
	}
	tx.From = from
	tx.To = to
	tx.Gas = config.Gas
	tx.GasPrice = config.GasPrice

	//todo
	if !strings.HasPrefix(value, "0x") {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("transfer value to int error.%s", err))
		}
		value = hexutil.EncodeBig(big.NewInt(intValue))
	}
	tx.Value = value

	params := make([]TxParams, 1)
	params[0] = tx

	res, _ := Send(params, "platon_sendTransaction")
	response := parseResponse(res)

	return response.Result, nil
}

func SendRawTransaction(from, to, value string, pkFilePath string) (string, error) {
	if len(accountPool) == 0 {
		parsePkFile(pkFilePath)
	}
	var v int64
	var err error
	if strings.HasPrefix(value, "0x") {
		bigValue, _ := hexutil.DecodeBig(value)
		v = bigValue.Int64()
	} else {
		v, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("transfer value to int error.%s", err))
		}
	}

	////
	//
	//for k, v := range accountPool {
	//	fmt.Println("acc", k.Hex())
	//	fmt.Println("value", fmt.Sprintf("%+v", v))
	//}

	acc, ok := accountPool[common.MustBech32ToAddress(from)]
	if !ok {
		return "", fmt.Errorf("private key not found in private key file,addr:%s", from)
	}
	nonce := getNonce(from)
	nonce++

	//// getBalance
	//
	//unlock := JsonParam{
	//	Jsonrpc: "2.0",
	//	Method:  "personal_unlockAccount",
	//	// {"method": "platon_getBalance", "params": [account, pwd, expire]}
	//	// {"jsonrpc":"2.0", "method":"eth_getBalance","params":["0xde1e758511a7c67e7db93d1c23c1060a21db4615","latest"],"id":67}
	//	Params: []interface{}{from, "latest"},
	//	Id:     1,
	//}
	//
	//// unlock
	//s, e := HttpPost(unlock)
	//if nil != e {
	//	fmt.Println("the gat balance err:", e)
	//}
	//fmt.Println("the balance:", s)

	newTx := getSignedTransaction(from, to, v, acc.Priv, nonce)

	hash, err := sendRawTransaction(newTx)
	if err != nil {
		panic(err)
	}
	return hash, nil
}

func sendRawTransaction(transaction *types.Transaction) (string, error) {
	bytes, _ := rlp.EncodeToBytes(transaction)
	res, err := Send([]string{hexutil.Encode(bytes)}, "platon_sendRawTransaction")
	if err != nil {
		panic(err)
	}
	response := parseResponse(res)

	return response.Result, nil
}

func getSignedTransaction(from, to string, value int64, priv *ecdsa.PrivateKey, nonce uint64) *types.Transaction {
	gas, _ := strconv.Atoi(config.Gas)
	gasPrice, _ := new(big.Int).SetString(config.GasPrice, 10)
	newTx, err := types.SignTx(types.NewTransaction(nonce, common.MustBech32ToAddress(to), big.NewInt(value), uint64(gas), gasPrice, []byte{}), types.NewEIP155Signer(new(big.Int).SetInt64(100)), priv)
	if err != nil {
		panic(fmt.Errorf("sign error,%s", err.Error()))
	}
	return newTx
}

func getNonce(addr string) uint64 {
	res, _ := Send([]string{addr, "latest"}, "platon_getTransactionCount")
	response := parseResponse(res)
	nonce, _ := hexutil.DecodeBig(response.Result)
	fmt.Println(addr, nonce)
	return nonce.Uint64()
}
