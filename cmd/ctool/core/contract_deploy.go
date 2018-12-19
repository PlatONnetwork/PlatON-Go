package core

import (
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

func parseFile(file string) []byte {
	codeBytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("An error occurred on parse file :%s \n%s", file, err.Error())
		panic(err)
	}
	return codeBytes
}

func parseResponse(r string) *Response {
	var resp = Response{}
	err := json.Unmarshal([]byte(r), &resp)

	if err != nil {
		panic(fmt.Sprintf("parse result error ! error:%s \n", err.Error()))
	}

	if resp.Error.Code != 0 {
		panic(fmt.Sprintf("send transaction error ,error:%v \n", resp.Error.Message))
	}
	return &resp
}

func Deploy(abiFilePath string, codeFilePath string, configPath string) error {

	abiBytes := parseFile(abiFilePath)
	codeBytes := parseFile(codeFilePath)

	param := [3][]byte{
		Int64ToBytes(depoly),
		codeBytes,
		abiBytes,
	}
	paramBytes, _ := rlp.EncodeToBytes(param)

	config := Config{}
	parseConfigJson(configPath, &config)

	params := DeployParams{
		From:     config.From,
		GasPrice: config.GasPrice,
		Gas:      config.Gas,
		Data:     hexutil.Encode(paramBytes),
	}

	paramList := make(List, 1)
	paramList[0] = params

	//paramJson, _ := json.Marshal(paramList)
	//fmt.Printf("\n request json data：%s\n", string(paramJson))

	r, err := Send(paramList, "eth_sendTransaction", config.Url)

	//fmt.Printf("\nresponse json：%s\n", r)

	resp := parseResponse(r)
	fmt.Printf("\ntrasaction hash: %s\n", resp.Result)

	ch := make(chan string, 1)
	go GetTransactionReceipt(resp.Result, ch, config.Url)

	select {
	case address := <-ch:
		fmt.Printf("contract address: %s\n", address)
	case <-time.After(time.Second * 200):
		fmt.Printf("get contract receipt timeout...more than 200 second.\n")
	}
	return err

}
