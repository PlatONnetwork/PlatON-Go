package gov

import (
	"bytes"
	"fmt"
	"testing"
)


func TestSplit(t *testing.T) {

	a := bytes.Split([]byte("test"), []byte("e"))
	fmt.Println(a[0])
	fmt.Println(a[1])
}

