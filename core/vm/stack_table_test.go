package vm

import "testing"

func TestMakeStackFunc(t *testing.T) {
	validateStack := makeStackFunc(0, 1025)
	stack := newstack()
	err := validateStack(stack)
	if err == nil {
		t.Errorf("Test makeStackFunc error")
	}
}