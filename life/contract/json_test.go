package contract

import (
	"Platon-go/life/utils"
	"encoding/json"
	"fmt"
	"testing"
)


type ABI struct {
	Version string 			`json:"version"`
	Abi []ABI_FUNC			`json:"abi"`
}

type ABI_FUNC struct {
	Method string 			`json:"method"`
	Args []ABI_Params			`json:"args"`
}

type ABI_Params struct {
	Name string 			`json:"name"`
	TypeName string 		`json:"typeName"`
	RealTypeName string 	`json:"realTypeName"`
}

func TestJsonFormat(t *testing.T) {
	// 定义结构，进行转换
	body := `{
	"version": "0.01",
	"abi": [{
			"method": "transfer",
			"args": [{
					"name": "from",
					"typeName": "address",
					"realTypeName": "char *"
				}, {
					"name": "to",
					"typeName": "address",
					"realTypeName": "char *"
				}, {
					"name": "asset",
					"typeName": "",
					"realTypeName": "int"
				}
			]
		}
	]
}`
	var abi ABI
	err := json.Unmarshal(utils.String2bytes(body), &abi)
	if err != nil {
		fmt.Println("error", err)
	}

	fmt.Println("version:", abi.Version)
	for _, v := range abi.Abi {
		fmt.Println("method:", v.Method)
		for _, arg := range v.Args {
			fmt.Println("name:", arg.Name)
			fmt.Println("typeName:", arg.TypeName)
			fmt.Println("realTypeName:", arg.RealTypeName)
		}
	}
}

