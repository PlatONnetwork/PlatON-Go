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
	"io/ioutil"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"gopkg.in/urfave/cli.v1"
)

var (
	DeployCmd = cli.Command{
		Name:   "deploy",
		Usage:  "deploy a contract",
		Action: deploy,
		Flags:  deployCmdFlags,
	}

	InvokeCmd = cli.Command{
		Name:    "invoke",
		Aliases: []string{"i"},
		Usage:   "invoke contract function",
		Action:  invoke,
		Flags:   invokeCmdFlags,
	}
)

func deploy(c *cli.Context) error {

	abiPath := c.String("abi")
	codePath := c.String("code")

	parseConfigJson(c.String(ConfigPathFlag.Name))
	err := DeployContract(abiPath, codePath)

	if err != nil {
		panic(fmt.Errorf("deploy contract error,%s", err.Error()))
	}
	return nil
}

func DeployContract(abiFilePath string, codeFilePath string) error {
	var err error

	abiBytes := parseFileToBytes(abiFilePath)
	codeBytes := parseFileToBytes(codeFilePath)

	param := [3][]byte{
		Int64ToBytes(deployContract),
		codeBytes,
		abiBytes,
	}
	paramBytes, err := rlp.EncodeToBytes(param)
	if err != nil {
		return fmt.Errorf("rlp encode error,%s", err.Error())
	}

	deployParams := DeployParams{
		From:     config.From,
		GasPrice: config.GasPrice,
		Gas:      config.Gas,
		Data:     hexutil.Encode(paramBytes),
	}

	params := make([]interface{}, 1)
	params[0] = deployParams

	//paramJson, _ := json.Marshal(paramList)
	//fmt.Printf("\n request json data：%s\n", string(paramJson))

	r, err := Send(params, "platon_sendTransaction")

	//fmt.Printf("\nresponse json：%s\n", r)

	resp := parseResponse(r)

	fmt.Printf("\ntrasaction hash: %s\n", resp.Result)

	// Get transaction receipt according to result
	ch := make(chan string, 1)
	exit := make(chan string, 1)
	go GetTransactionReceipt(resp.Result, ch, exit)

	// Getting receipt
	select {
	case address := <-ch:
		fmt.Printf("contract address: %s\n", address)
	case <-time.After(time.Second * 200):
		exit <- "exit"
		fmt.Printf("get contract receipt timeout...more than 200 second.\n")
	}
	return err

}

func parseFileToBytes(file string) []byte {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("parse file %s error,%s", file, err.Error()))
	}
	return bytes
}

func invoke(c *cli.Context) error {
	addr := c.String("addr")
	abiPath := c.String("abi")
	funcParams := c.String("func")
	txType := c.Int("type")

	//param check
	if abiPath == "" {
		fmt.Printf("abi can't be empty!")
		return nil
	}
	if addr == "" {
		fmt.Printf("addr can't be empty!")
		return nil
	}
	if funcParams == "" {
		fmt.Printf("func can't be empty!")
		return nil
	}
	parseConfigJson(c.String(ConfigPathFlag.Name))

	err := InvokeContract(addr, abiPath, funcParams, txType)
	if err != nil {
		panic(fmt.Errorf("invokeContract contract error,%s", err.Error()))
	}
	return nil
}

/**

 */
func InvokeContract(contractAddr string, abiPath string, funcParams string, txType int) error {

	//Judging whether this contract exists or not
	if !getContractByAddress(contractAddr) {
		return fmt.Errorf("the contract address is not exist ...")
	}

	//parse the function and param
	funcName, inputParams := GetFuncNameAndParams(funcParams)

	//Judging whether this method exists or not
	abiFunc, err := parseFuncFromAbi(abiPath, funcName)
	if err != nil {
		return err
	}

	if len(abiFunc.Inputs) != len(inputParams) {
		return fmt.Errorf("incorrect number of parameters ,request=%d,get=%d\n", len(abiFunc.Inputs), len(inputParams))
	}

	if txType == 0 {
		txType = invokeContract
	}

	paramArr := [][]byte{
		Int64ToBytes(int64(txType)),
		[]byte(funcName),
	}

	for i, v := range inputParams {
		input := abiFunc.Inputs[i]
		p, e := StringConverter(v, input.Type)
		if e != nil {
			return fmt.Errorf("incorrect param type: %s,index:%d", v, i)
		}
		paramArr = append(paramArr, p)
	}

	paramBytes, e := rlp.EncodeToBytes(paramArr)
	if e != nil {
		return fmt.Errorf("rpl encode error,%s", e.Error())
	}

	txParams := TxParams{
		From:     config.From,
		To:       contractAddr,
		GasPrice: config.GasPrice,
		Gas:      config.Gas,
		Data:     hexutil.Encode(paramBytes),
	}

	var r string
	if abiFunc.Constant == "true" {
		params := make([]interface{}, 2)
		params[0] = txParams
		params[1] = "latest"

		paramJson, _ := json.Marshal(params)
		fmt.Printf("\n request json data：%s \n", string(paramJson))
		r, err = Send(params, "platon_call")
	} else {
		params := make([]interface{}, 1)
		params[0] = txParams

		paramJson, _ := json.Marshal(params)
		fmt.Printf("\n request json data：%s \n", string(paramJson))
		r, err = Send(params, "platon_sendTransaction")
	}

	fmt.Printf("\n response json：%s \n", r)

	if err != nil {
		return fmt.Errorf("send http post to invokeContract contract error,%s", e.Error())
	}
	resp := parseResponse(r)

	//parse the return type through adi
	if abiFunc.Constant == "true" {
		if len(abiFunc.Outputs) != 0 && abiFunc.Outputs[0].Type != "void" {
			bytes, _ := hexutil.Decode(resp.Result)
			result := BytesConverter(bytes, abiFunc.Outputs[0].Type)
			fmt.Printf("\nresult: %v\n", result)
			return nil
		}
		fmt.Printf("\n result: []\n")
	} else {
		fmt.Printf("\n trasaction hash: %s\n", resp.Result)
	}
	return nil
}

/**
  Judging whether a contract exists through platon_getCode
*/
func getContractByAddress(addr string) bool {

	params := []string{addr, "latest"}
	r, err := Send(params, "platon_getCode")
	if err != nil {
		fmt.Printf("send http post to get contract address error ")
		return false
	}

	var resp = Response{}
	err = json.Unmarshal([]byte(r), &resp)
	if err != nil {
		fmt.Printf("parse platon_getCode result error ! \n %s", err.Error())
		return false
	}

	if resp.Error.Code != 0 {
		fmt.Printf("platon_getCode error ,error:%v", resp.Error.Message)
		return false
	}
	//fmt.Printf("trasaction hash: %s\n", resp.Result)

	if resp.Result != "" && len(resp.Result) > 2 {
		return true
	} else {
		return false
	}
}

/*
  Loop call to get transactionReceipt... until 200s timeout
*/
func GetTransactionReceipt(txHash string, ch chan string, exit chan string) {
	var receipt = Receipt{}
	var contractAddr string
	for {
		res, e := Send([]string{txHash}, "platon_getTransactionReceipt")
		if e != nil {
			panic(fmt.Sprintf("send http post to get transaction receipt error！\n %s", e.Error()))
		}
		e = json.Unmarshal([]byte(res), &receipt)
		if e != nil {
			panic(fmt.Sprintf("parse get receipt result error ! \n %s", e.Error()))
		}
		contractAddr = receipt.Result.ContractAddress
		if contractAddr != "" {
			ch <- contractAddr
			break
		}
		select {
		case <-exit:
			break
		default:
		}
	}
}
