package resolver

import (
	"Platon-go/core/vm/life/exec"
	"fmt"
)

var (
	clang		int = 0x01
	golang		int = 0x02
)

// new import resolver
func NewResolver(lang int) exec.ImportResolver {
	switch lang {
	case clang:
		return &CResolver{}
	case golang:
	default:
	}
	return nil
}

func MallocString(vm *exec.VirtualMachine, str string) int64 {
	mem := vm.Memory
	size := len([]byte(str)) + 1

	if mem.Current+size > len(mem.Memory) {
		panic(fmt.Sprintf("out of memory  current:%d len:%d memory len:%d", mem.Current, size, len(mem.Memory)))
	}

	pos := int64(mem.Current)
	mem.MemPoints[mem.Current] = size
	copy(mem.Memory[mem.Current:], []byte(str))
	mem.Current += size
	return pos
}
