package resolver

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	"Platon-go/core/vm/life/exec"
)

var (
	cfc  				= newCfcSet()
	cgbl 				= newGlobalSet()
)

type CResolver struct{}

func (r *CResolver) ResolveFunc(module, field string) *exec.FunctionImport {
	df := &exec.FunctionImport{
		Execute: func(vm *exec.VirtualMachine) int64 {
			panic(fmt.Sprintf("unsupport func module:%s field:%s", module, field))
		},
		GasCost: func(vm *exec.VirtualMachine) (uint64, error) {
			panic(fmt.Sprintf("unsupport gas cost module:%s field:%s", module, field))
		},
	}

	if m, exist := cfc[module]; exist == true {
		if f, exist := m[field]; exist == true {
			return f
		} else {
			return df
		}
	} else {
		return df
	}
}

func (r *CResolver) ResolveGlobal(module, field string) int64 {
	if m, exist := cgbl[module]; exist == true {
		if g, exist := m[field]; exist == true {
			return g
		} else {
			return 0
			//panic("unknown field " + field)

		}
	} else {
		return 0
		//panic("unknown module " + module)
	}
}

func newCfcSet() map[string]map[string]*exec.FunctionImport {
	return map[string]map[string]*exec.FunctionImport{
		"env": {
			"malloc":  &exec.FunctionImport{Execute: envMalloc, GasCost: envMallocGasCost},
			"free":    &exec.FunctionImport{Execute: envFree, GasCost: envFreeGasCost},
			"calloc":  &exec.FunctionImport{Execute: envCalloc, GasCost: envCallocGasCost},
			"realloc": &exec.FunctionImport{Execute: envRealloc, GasCost: envReallocGasCost},

			"memcpy":  &exec.FunctionImport{Execute: envMemcpy, GasCost: envMemcpyGasCost},
			"memmove": &exec.FunctionImport{Execute: envMemmove, GasCost: envMemmoveGasCost},
			"memcmp":  &exec.FunctionImport{Execute: envMemcpy, GasCost: envMemmoveGasCost},
			"memset":  &exec.FunctionImport{Execute: envMemset, GasCost: envMemsetGasCost},

			"prints":     &exec.FunctionImport{Execute: envPrints, GasCost: envPrintsGasCost},
			"prints_l":   &exec.FunctionImport{Execute: envPrintsl, GasCost: envPrintslGasCost},
			"printi":     &exec.FunctionImport{Execute: envPrinti, GasCost: envPrintiGasCost},
			"printui":    &exec.FunctionImport{Execute: envPrintui, GasCost: envPrintuiGasCost},
			"printi128":  &exec.FunctionImport{Execute: envPrinti128, GasCost: envPrinti128GasCost},
			"printui128": &exec.FunctionImport{Execute: envPrintui128, GasCost: envPrintui128GasCost},
			"printsf":    &exec.FunctionImport{Execute: envPrintsf, GasCost: envPrintsfGasCost},
			"printdf":    &exec.FunctionImport{Execute: envPrintdf, GasCost: envPrintdfGasCost},
			"printqf":    &exec.FunctionImport{Execute: envPrintqf, GasCost: envPrintqfGasCost},
			"printn":     &exec.FunctionImport{Execute: envPrintn, GasCost: envPrintnGasCost},
			"printhex":   &exec.FunctionImport{Execute: envPrinthex, GasCost: envPrinthexGasCost},

			"abort": &exec.FunctionImport{Execute: envAbort, GasCost: envAbortGasCost},
		},
	}
}

func newGlobalSet() map[string]map[string]int64 {
	return map[string]map[string]int64{
		"env": {
			"stderr": 0,
			"stdin":  0,
			"stdout": 0,
		},
	}
}

//void * memcpy ( void * destination, const void * source, size_t num );
func envMemcpy(vm *exec.VirtualMachine) int64 {
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	return int64(dest)
}

func envMemcpyGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

//void * memmove ( void * destination, const void * source, size_t num );
func envMemmove(vm *exec.VirtualMachine) int64 {
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	return int64(dest)
}

func envMemmoveGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

//int memcmp ( const void * ptr1, const void * ptr2, size_t num );
func envMemcmp(vm *exec.VirtualMachine) int64 {
	ptr1 := int(uint32(vm.GetCurrentFrame().Locals[0]))
	ptr2 := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	pos := 0
	for pos < len {
		if vm.Memory.Memory[ptr1+pos] == vm.Memory.Memory[ptr2+pos] {
			pos += 1
		} else if vm.Memory.Memory[ptr1+pos] <= vm.Memory.Memory[ptr2+pos] {
			return -1
		} else {
			return 1
		}
	}
	return 0
}

func envMemcmpGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

//void * memset ( void * ptr, int value, size_t num );
func envMemset(vm *exec.VirtualMachine) int64 {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	value := int(uint32(vm.GetCurrentFrame().Locals[1]))
	num := int(uint32(vm.GetCurrentFrame().Locals[2]))

	pos := 0
	for pos < num {
		vm.Memory.Memory[ptr+pos] = byte(value)
	}
	return int64(ptr)
}

func envMemsetGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

//libc prints()
func envPrints(vm *exec.VirtualMachine) int64 {
	start := int(uint32(vm.GetCurrentFrame().Locals[0]))
	end := 0
	for end = start; end < len(vm.Memory.Memory); end++ {
		if vm.Memory.Memory[end] == 0 {
			break
		}
	}
	fmt.Printf("%s", string(vm.Memory.Memory[start:end]))
	return 0
}

func envPrintsGasCost(vm *exec.VirtualMachine) (uint64, error) {
	start := int(uint32(vm.GetCurrentFrame().Locals[0]))
	end := 0
	for end = start; end < len(vm.Memory.Memory); end++ {
		if vm.Memory.Memory[end] == 0 {
			break
		}
	}
	return uint64(end - start), nil
}

//libc prints_l
func envPrintsl(vm *exec.VirtualMachine) int64 {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	msg := vm.Memory.Memory[ptr : ptr+msgLen]
	fmt.Printf("%s", string(msg))
	return 0
}

func envPrintslGasCost(vm *exec.VirtualMachine) (uint64, error) {
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	return uint64(msgLen), nil
}

//libc printi()
func envPrinti(vm *exec.VirtualMachine) int64 {
	fmt.Printf("%d", int(uint32(vm.GetCurrentFrame().Locals[0])))
	return 0
}

func envPrintiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintui(vm *exec.VirtualMachine) int64 {
	fmt.Printf("%d", vm.GetCurrentFrame().Locals[0])
	return 0
}

func envPrintuiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrinti128(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	num := new(big.Int)
	num.SetBytes(vm.Memory.Memory[pos : pos+16])

	fmt.Printf("%s", num.String())
	return 0
}

func envPrinti128GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintui128(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	num := new(big.Int)
	num.SetBytes(vm.Memory.Memory[pos : pos+16])
	fmt.Printf("%s", num.String())
	return 0
}

func envPrintui128GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintsf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	bits := binary.LittleEndian.Uint32(vm.Memory.Memory[pos : pos+4])
	float := math.Float32frombits(bits)
	fmt.Printf("%f", float)
	return 0
}

func envPrintsfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintdf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	bits := binary.LittleEndian.Uint64(vm.Memory.Memory[pos : pos+8])
	float := math.Float64frombits(bits)
	fmt.Printf("%f", float)
	return 0
}

func envPrintdfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintqf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	num := new(big.Int)
	num.SetBytes(vm.Memory.Memory[pos : pos+16])
	float := new(big.Float)
	float.SetInt(num)
	fmt.Printf("%s", float.String())
	return 0
}

func envPrintqfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintn(vm *exec.VirtualMachine) int64 {
	fmt.Printf("%d", int(uint32(vm.GetCurrentFrame().Locals[0])))
	return 0
}

func envPrintnGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrinthex(vm *exec.VirtualMachine) int64 {
	data := int(uint32(vm.GetCurrentFrame().Locals[0]))
	dataLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	fmt.Printf("%x", vm.Memory.Memory[data:dataLen])
	return 0
}

func envPrinthexGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

//libc malloc()
func envMalloc(vm *exec.VirtualMachine) int64 {
	mem := vm.Memory
	size := int(uint32(vm.GetCurrentFrame().Locals[0]))
	if mem.Current+size > len(mem.Memory) {
		panic("out of memory")
	}
	pos := int64(mem.Current)
	mem.MemPoints[mem.Current] = size
	mem.Current += size

	return pos
}

func envMallocGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

//libc free()
func envFree(vm *exec.VirtualMachine) int64 {
	return 0
}

func envFreeGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

//libc calloc()
func envCalloc(vm *exec.VirtualMachine) int64 {
	mem := vm.Memory
	num := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	total := num * size
	if mem.Current+total > len(mem.Memory) {
		panic("out of memory")
	}

	for i := 0; i < total; i++ {
		mem.Memory[mem.Current+i] = 0
	}

	pos := int64(mem.Current)
	mem.MemPoints[mem.Current] = total
	mem.Current += total

	return pos
}

func envCallocGasCost(vm *exec.VirtualMachine) (uint64, error) {
	num := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	total := num * size
	return uint64(total), nil
}

func envRealloc(vm *exec.VirtualMachine) int64 {
	mem := vm.Memory
	ptr := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))

	if size == 0 {
		return 0
	}

	if _, exist := mem.MemPoints[ptr]; exist != true {
		panic("realloc error")
	}

	if mem.Current+size > len(mem.Memory) {
		panic("out of memory")
	}

	pos := int64(mem.Current)
	mem.MemPoints[mem.Current] = size
	mem.Current += size
	return pos
}

func envReallocGasCost(vm *exec.VirtualMachine) (uint64, error) {
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	return uint64(size), nil
}

func envAbort(vm *exec.VirtualMachine) int64 {
	return 0
}

func envAbortGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}
