package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	transfer = iota
	depoly
	invoke
	vote
	permission
)

type Element interface{}
type List []Element

type JsonParam struct {
	Jsonrpc string				`json:"jsonrpc"`
	Method  string				`json:"method"`
	Params  interface{}			`json:"params"`
	Id      int					`json:"id"`
}
type BasicParam struct {
	From     string				`json:"from"`
	Gas      string				`json:"gas"`
	GasPrice string				`json:"gas_price"`
}

type TxParams struct {
	From     string				`json:"from"`
	To       string				`json:"to"`
	Gas      string				`json:"gas"`
	GasPrice string				`json:"gasPrice"`
	Value    string				`json:"value"`
	Data     string				`json:"data"`
}

type DeployParams struct {
	From     string				`json:"from"`
	Gas      string				`json:"gas"`
	GasPrice string				`json:"gasPrice"`
	Data     string				`json:"data"`
}

type Config struct {
	From     string				`json:"from"`
	Gas      string				`json:"gas"`
	GasPrice string				`json:"gasPrice"`
	Url      string				`json:"url"`
}

type AbiS struct {
	//Bytecode string
	ContractName string			`json:"contractName"`
	AbiJson      string			`json:"abiJson"`
	Abi          []FuncDesc		`json:"abi"`
}

type FuncDesc struct {
	Method string				`json:"method"`
	Args   []struct {
		Name         string		`json:"name"`
		TypeName     string		`json:"typeName"`
		RealTypeName string		`json:"realTypeName"`
	}							`json:"args"`
	Return   string				`json:"return"`
	FuncType string				`json:"funcType"`
}

type Response struct {
	Jsonrpc string				`json:"jsonrpc"`
	Result  string				`json:"result"`
	Id      int					`json:"id"`
	Error   struct {
		Code    int32			`json:"code"`
		Message string			`json:"message"`
	}							`json:"error"`
}

type Receipt struct {
	Jsonrpc string						`json:"jsonrpc"`
	Id      int							`json:"id"`
	Result  struct {
		BlockHash         string		`json:"blockHash"`
		BlockNumber       string		`json:"blockNumber"`
		ContractAddress   string 		`json:"contractAddress"`
		CumulativeGasUsed string 		`json:"cumulativeGas_used"`
		From              string 		`json:"from"`
		GasUsed           string 		`json:"gasUsed"`
		Root              string 		`json:"root"`
		To                string 		`json:"to"`
		TransactionHash   string 		`json:"transactionHash"`
		TransactionIndex  string 		`json:"transactionIndex"`
	} 									`json:"result"`
}

func parseConfigJson(configPath string, param interface{}) {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bytes, &param); err != nil {
		panic(err)
	}
}

func parseAbiFromJson(fileName string) (AbiS, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(fmt.Sprintf("ReadFile: %s", err.Error()))
	}
	a := AbiS{}
	if err := json.Unmarshal(bytes, &a); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		panic(fmt.Sprintf("Unmarshal: %s", err.Error()))
	}
	abijson, _ := json.Marshal(a.Abi)
	a.AbiJson = string(abijson)
	return a, nil
}

func parseFuncFromAbi(fileName string, funcName string) FuncDesc {
	abis, _ := parseAbiFromJson(fileName)

	for _, value := range abis.Abi {
		if value.Method == funcName {
			return value
		}
	}
	return FuncDesc{}
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

/**
  通过解析abi查找所调用方法
*/
func GetFuncNameAndParams(f string) (string, []string) {
	funcName := string(f[0:strings.Index(f, "(")])

	paramString := string(f[strings.Index(f, "(")+1 : strings.LastIndex(f, ")")])
	if paramString == "" {
		return funcName, []string{}
	}

	params := strings.Split(paramString, ",")
	for index, param := range params {
		if strings.HasPrefix(param, "\"") {
			params[index] = param[strings.Index(param, "\"")+1 : strings.LastIndex(param, "\"")]
		}
	}
	return funcName, params

}