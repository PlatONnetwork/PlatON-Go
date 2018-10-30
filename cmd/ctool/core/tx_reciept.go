package core

import (
	"encoding/json"
	"fmt"
)

/**
命令入口
*/
func GetTxReceipt(txHash string, configPath string) error {

	config := Config{}
	parseConfigJson(configPath, &config)

	var receipt = Receipt{}
	//根据result获取交易receipt
	res, _ := Send([]string{txHash}, "eth_getTransactionReceipt", config.Url)
	e := json.Unmarshal([]byte(res), &receipt)
	if e != nil {
		panic(fmt.Sprintf("parse get receipt result error ! \n %s", e.Error()))
	}

	if receipt.Result.BlockHash == "" {
		fmt.Print("no receipt found")
	} else {
		fmt.Printf("\nget transaction reciept result: \n\n %s", res)
	}
	return nil

}

/*
  单次调用获取transactionReceipt...
*/
//func getReceipt(url string,param JsonParam) (Receipt)  {
//	//根据result获取交易receipt
//	r, e := HttpPost(param,url)
//	if e != nil {
//		panic(fmt.Sprintf("send http post error .\n %s" + e.Error()))
//	}
//
//	var receipt = Receipt{}
//	e = json.Unmarshal([]byte(r), &receipt)
//	if e != nil {
//		panic(fmt.Sprintf("parse receipt result error ! \n %s", e.Error()))
//	}
//	return receipt
//}

/*
  循环调用获取transactionReceipt....直至200s超时
*/
func GetTransactionReceipt(txHash string, ch chan string, url string) {
	//根据result获取交易receipt
	var receipt = Receipt{}
	var contractAddr string
	for {
		res, _ := Send([]string{txHash}, "eth_getTransactionReceipt", url)
		e := json.Unmarshal([]byte(res), &receipt)
		if e != nil {
			panic(fmt.Sprintf("parse get receipt result error ! \n %s", e.Error()))
		}
		contractAddr = receipt.Result.ContractAddress
		if contractAddr != "" {
			ch <- contractAddr
			break
		}
	}
}
