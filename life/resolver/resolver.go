package resolver

import (
	"github.com/PlatONnetwork/PlatON-Go/life/exec"
)

var (
	clang  int = 0x01
	golang int = 0x02
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

	pos := mem.Malloc(size)
	copy(mem.Memory[pos:pos+size], []byte(str))
	vm.ExternalParams = append(vm.ExternalParams, int64(pos))
	return int64(pos)
}
