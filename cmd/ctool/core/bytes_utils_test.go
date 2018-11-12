package core

import (
	"Platon-go/common"
	"fmt"
	"testing"
)

func TestByteConvert(t *testing.T) {
	//bytes, _ := hexutil.Decode("0x0c55699c")
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
	//buf := bytes.NewBuffer([]byte{})
	//binary.Write(buf, binary.BigEndian, "true")
	//fmt.Println(buf.Bytes())
	//fmt.Println(len(buf.Bytes()))

	//fmt.Printf("%v",i)
}
