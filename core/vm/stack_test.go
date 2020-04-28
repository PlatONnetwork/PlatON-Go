package vm

import (
	"math/big"
	"testing"
)

func TestSwap(t *testing.T) {
	stack := newstack()
	element0 := new(big.Int).SetUint64(100)
	element1 := new(big.Int).SetUint64(200)
	stack.push(element0)
	stack.push(element1)
	stack.swap(2)
	actual := stack.pop()
	if actual.Cmp(element0) != 0 {
		t.Errorf("Test swap, expected  %v, got %v", element0, actual)
	}
	actual = stack.pop()
	if actual.Cmp(element1) != 0 {
		t.Errorf("Test swap, expected  %v, got %v", element1, actual)
	}
}

func TestPrint(t *testing.T) {
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.Print()
}