package core

import (
	"encoding/json"
	"fmt"
)

func GetTxReceipt(txHash string, configPath string) error {

	config := Config{}
	parseConfigJson(configPath, &config)

	var receipt = Receipt{}
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

func GetTransactionReceipt(txHash string, ch chan string, url string) {
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
