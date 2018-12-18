package core

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"fmt"
	"testing"
)

func TestByteConvert(t *testing.T) {
	hash := common.BytesToHash(Int32ToBytes(121))
	result := BytesConverter(hash.Bytes(), "int32")
	fmt.Printf("\nresult: %v\n", result)
}

func TestStringConverter(t *testing.T) {
	result, err := StringConverter("false", "bool")
	fmt.Printf("\nresult: %v\n", result)
	if err != nil {
		fmt.Printf("\nerr: %v\n", err.Error())
	}
}
