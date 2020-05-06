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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	parseConfig(t)
	assert.NotEqual(t, config.Gas, "", "the Gas is empty")
	assert.NotEqual(t, config.GasPrice, "", "the GasPrice is empty")
	assert.NotEqual(t, config.From, "", "the From is empty")
}

func TestParseFuncFromAbi(t *testing.T) {
	funcDesc, err := parseFuncFromAbi(abiFilePath, "atransfer")
	assert.Nil(t, err, fmt.Sprintf("%v", err))
	assert.NotEqual(t, funcDesc, "", "the funcDesc is nil in abi")
}

func TestParseAbiFromJson(t *testing.T) {
	a, e := parseAbiFromJson(abiFilePath)
	assert.Nil(t, e, fmt.Sprintf("parse abi json error! \n，%v", e))
	marshal, e := json.Marshal(a)
	assert.Nil(t, e, fmt.Sprintf("parse data to json error! \n，%v", e))
	assert.NotEqual(t, marshal, "", "the data is nil")
}

func TestHttpPostTransfer(t *testing.T) {

	//platon, datadir := prepare(t)
	//
	//param := JsonParam{
	//	Jsonrpc: "2.0",
	//	Method:  "platon_sendTransaction",
	//	Params: []TxParams{
	//		{
	//			From:     from,
	//			To:       to,
	//			Value:    "0xf4240",        // 1000000
	//			Gas:      "0x5208",         // 21000
	//			GasPrice: "0x2d79883d2000", // 50000000000000
	//		},
	//	},
	//	Id: 1,
	//}
	//
	//r, e := HttpPost(param)
	//assert.Nil(t, e, fmt.Sprintf("test http post error: %v", e))
	//assert.NotEqual(t, r, "", "the result is nil")
	//t.Log("the result ", r)
	//clean(platon, datadir)
}

func TestHttpPostDeploy(t *testing.T) {
	//platon, datadir := prepare(t)
	//
	//deployParams := DeployParams{
	//	From:     from,
	//	Gas:      "0x400000",
	//	GasPrice: "0x9184e72a000",
	//}
	//
	//params := make([]interface{}, 1)
	//params[0] = deployParams
	//param := JsonParam{
	//	Jsonrpc: "2.0",
	//	Method:  "platon_sendTransaction",
	//	Params: []TxParams{
	//		{
	//			From:     from,
	//			To:       to,
	//			Value:    "0xf4240",        // 1000000
	//			Gas:      "0x5208",         // 21000
	//			GasPrice: "0x2d79883d2000", // 50000000000000
	//		},
	//	},
	//	Id: 1,
	//}
	//
	//r, e := HttpPost(param)
	//assert.Nil(t, e, fmt.Sprintf("test http post error: %v", e))
	//assert.NotEqual(t, r, "", "the result is nil")
	//t.Log("the result ", r)
	//
	//var resp = Response{}
	//err := json.Unmarshal([]byte(r), &resp)
	//if err != nil {
	//	t.Fatalf("parse result error ! \n %s", err.Error())
	//}
	//
	//if resp.Error.Code != 0 {
	//	t.Fatalf("send transaction error ,error:%v", resp.Error.Message)
	//}
	//fmt.Printf("trasaction hash: %s\n", resp.Result)
	//
	//// Get transaction receipt according to result
	//ch := make(chan string, 1)
	//exit := make(chan string, 1)
	//go GetTransactionReceipt(resp.Result, ch, exit)
	//
	//// Then, we use the timeout channel
	//select {
	//case address := <-ch:
	//	fmt.Printf("contract address:%s\n", address)
	//case <-time.After(time.Second * 10):
	//	exit <- "exit"
	//	fmt.Printf("get contract receipt timeout...more than 100 second.\n")
	//}
	//
	//clean(platon, datadir)
}

func TestHttpCallContact(t *testing.T) {
	//platon, datadir := prepare(t)
	//
	//param1 := uint(33)
	//b := new(bytes.Buffer)
	//rlp.Encode(b, param1)
	//
	//params := TxParams{
	//	From:     from,
	//	To:       "0xace6bdba54c8c359e70f541bfc1cabaf0244b916",
	//	Value:    "0x2710",
	//	Gas:      "0x76c00",
	//	GasPrice: "0x9184e72a000",
	//	Data:     "0x60fe47b10000000000000000000000000000000000000000000000000000000000000011",
	//}
	//
	//param := JsonParam{
	//	Jsonrpc: "2.0",
	//	Method:  "platon_sendTransaction",
	//	Params:  []TxParams{params},
	//	Id:      1,
	//}
	//paramJson, _ := json.Marshal(param)
	//fmt.Println(string(paramJson))
	//r, e := HttpPost(param)
	//assert.Nil(t, e, fmt.Sprintf("test http post error: %v", e))
	//assert.NotEqual(t, r, "", "the result is nil")
	//t.Log("the result ", r)
	//clean(platon, datadir)

}

func TestGetFuncParam(t *testing.T) {
	f := "set()"
	s, strings := GetFuncNameAndParams(f)
	assert.Equal(t, s, "set", fmt.Sprintf("the result is not `set`, but it is %s", s))
	assert.Equal(t, len(strings), 0, fmt.Sprintf("the params len is not 0, but it is %d", len(strings)))
}
