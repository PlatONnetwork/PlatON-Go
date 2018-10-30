package core

import (
	"fmt"
	"testing"
)

func TestByteConvert(t *testing.T) {

	result := BytesConverter(Float64ToBytes(3434.4545), "string")
	fmt.Printf("\nresult: %v\n", result)

}

func TestStringConverter(t *testing.T) {

	result, err := StringConverter("2343234", "uint64")
	fmt.Printf("\nresult: %v\n", result)
	if err != nil {
		fmt.Printf("\nerr: %v\n", err.Error())
	}

}
