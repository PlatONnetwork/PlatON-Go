package utils

import (
	"encoding/json"
	"fmt"
)

type WasmAbi struct {
	Version string 			`json:"version"`
	Abi 	[]Func			`json:"abi"`
}

type Func struct {
	Method string 			`json:"method"`
	Args []Args				`json:"args"`
}

type Args struct {
	Name string 			`json:"name"`
	TypeName string 		`json:"typeName"`
	RealTypeName string 	`json:"realTypeName"`
}

func (abi *WasmAbi) FromJson(body []byte) error {
	if body == nil {
		return fmt.Errorf("invalid param. %v", body)
	}
	err := json.Unmarshal(body, abi)
	return err
}