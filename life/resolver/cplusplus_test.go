package resolver

import (
	"testing"
	"fmt"
)

func TestCfcSet(t *testing.T) {
	cfgSet := newCfcSet()
	for k, v := range cfgSet {
		fmt.Println("key:", k)
		for k1, v1 := range v {
			fmt.Printf("key1: %v, v1 type: %T \n", k1, v1)
		}
	}
}
