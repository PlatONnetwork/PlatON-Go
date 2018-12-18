package core

import (
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"encoding/json"
	"fmt"
)

/// contract invoke.
func ContractInvoke(contractAddr string, abiPath string, funcParams string, configPath string) {
	config := Config{}
	parseConfigJson(configPath, &config)

	if !getContractByAddress(contractAddr, config.Url) {
		fmt.Printf("the contract address is not exist ...")
		return
	}

	funcName, inputParams := GetFuncNameAndParams(funcParams)

	abiFunc := parseFuncFromAbi(abiPath, funcName)
	if abiFunc.Name == "" {
		fmt.Printf("the function not exist ,func= %s\n", funcName)
		return
	}

	if len(abiFunc.Inputs) != len(inputParams) {
		fmt.Printf("incorrect number of parameters ,request=%d,get=%d\n", len(abiFunc.Inputs), len(inputParams))
		return
	}

	paramArr := [][]byte{
		Int64ToBytes(invoke),
		[]byte(funcName),
	}

	for i, v := range inputParams {
		input := abiFunc.Inputs[i]
		p, e := StringConverter(v, input.Type)
		if e != nil {
			fmt.Printf("incorrect param type: %s,index:%d", v, i)
			return
		}
		paramArr = append(paramArr, p)
	}

	paramBytes, e := rlp.EncodeToBytes(paramArr)
	if e != nil {
		fmt.Printf("rpl encode error...\n %s", e.Error())
		panic(e)
	}

	params := TxParams{
		From:     config.From,
		To:       contractAddr,
		GasPrice: config.GasPrice,
		Gas:      config.Gas,
		Data:     hexutil.Encode(paramBytes),
	}

	var r string
	var err error
	if abiFunc.Constant == "true" {
		paramList := make(List, 2)
		paramList[0] = params
		paramList[1] = "latest"

		paramJson, _ := json.Marshal(paramList)
		fmt.Printf("\nrequest json data：%s\n", string(paramJson))
		r, err = Send(paramList, "eth_call", config.Url)
	} else {
		paramList := make(List, 1)
		paramList[0] = params

		paramJson, _ := json.Marshal(paramList)
		fmt.Printf("\nrequest json data：%s\n", string(paramJson))
		r, err = Send(paramList, "eth_sendTransaction", config.Url)
	}

	fmt.Printf("\nresponse json：%s\n", r)

	if err != nil {
		fmt.Printf("send http post to invoke contract error ")
		return
	}

	var resp = Response{}
	err = json.Unmarshal([]byte(r), &resp)
	if err != nil {
		fmt.Printf("parse eth_sendTransaction result error ! \n %s", err.Error())
		return
	}

	if resp.Error.Code != 0 {
		fmt.Printf("eth_sendTransaction error ,error:%v", resp.Error.Message)
		return
	}

	if abiFunc.Constant == "true" {
		if len(abiFunc.Outputs) != 0 && abiFunc.Outputs[0].Type != "void" {
			bytes, _ := hexutil.Decode(resp.Result)
			result := BytesConverter(bytes, abiFunc.Outputs[0].Type)
			fmt.Printf("\nresult: %v\n", result)
		}
		fmt.Printf("\nresult: []\n")
	} else {
		fmt.Printf("\ntrasaction hash: %s\n", resp.Result)
	}

}

func getContractByAddress(addr, url string) bool {

	params := []string{addr, "latest"}
	r, err := Send(params, "eth_getCode", url)
	if err != nil {
		fmt.Printf("send http post to get contract address error ")
		return false
	}

	var resp = Response{}
	err = json.Unmarshal([]byte(r), &resp)
	if err != nil {
		fmt.Printf("parse eth_getCode result error ! \n %s", err.Error())
		return false
	}

	if resp.Error.Code != 0 {
		fmt.Printf("eth_getCode error ,error:%v", resp.Error.Message)
		return false
	}
	//fmt.Printf("trasaction hash: %s\n", resp.Result)

	if resp.Result != "" && len(resp.Result) > 2 {
		return true
	} else {
		return false
	}
}
