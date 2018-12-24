package resolver

// #cgo CFLAGS: -I./softfloat/source/include -I./builtins
// #define SOFTFLOAT_FAST_INT64
// #define SOFTFLOAT_ROUND_EVEN
// #define INLINE_LEVEL 5
// #define SOFTFLOAT_FAST_DIV32TO16
// #define SOFTFLOAT_FAST_DIV64TO32
// #cgo LDFLAGS: -L./softfloat/build -lsoftfloat -L./builtins/build -lbuiltins
// #include "softfloat.h"
// #include "compiler_builtins.hpp"
// #include "int_t.h"
import "C"

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/PlatONnetwork/PlatON-Go/life/exec"
)

type uint128 struct {
	high uint64
	low  uint64
}

func (u *uint128) lsh(shift uint) {
	if shift >= 128 {
		u.low = 0
		u.high = 0
	} else {
		var halfSize uint = 128 / 2

		if shift >= halfSize {
			shift -= halfSize
			u.high = u.low
			u.low = 0
		}

		if shift != 0 {
			u.high <<= shift
		}

		var mask uint64 = ^(math.MaxUint64 >> shift)
		u.high |= (u.low & mask) >> (halfSize - shift)
		u.low <<= shift
	}
}

func (u *uint128) rsh(shift uint) {
	if shift >= 128 {
		u.high = 0
		u.low = 0
	} else {
		var halfSize uint = 128 / 2

		if shift >= halfSize {
			shift -= halfSize
			u.low = u.high
			u.high = 0
		}

		if shift != 0 {
			u.low >>= shift
		}

		var mask uint64 = ^(math.MaxUint64 << shift)
		u.low |= (u.high & mask) << (halfSize - shift)
		u.high >>= shift
	}
}

// arithmetic long double
func env__ashlti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	u := &uint128{
		low:  uint64(frame.Locals[1]),
		high: uint64(frame.Locals[2]),
	}
	shift := uint(frame.Locals[3])
	u.lsh(shift)

	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf, u.low)
	binary.LittleEndian.PutUint64(buf[8:], u.high)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__ashlti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__ashrti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	low := C.uint64_t(frame.Locals[1])
	high := C.uint64_t(frame.Locals[2])
	shift := C.uint32_t(frame.Locals[3])

	ret := C.___ashriti3(low, high, shift)
	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__ashrti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__lshlti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	u := &uint128{
		low:  uint64(frame.Locals[1]),
		high: uint64(frame.Locals[2]),
	}
	shift := uint(frame.Locals[3])
	u.lsh(shift)

	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf, u.low)
	binary.LittleEndian.PutUint64(buf[8:], u.high)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__lshlti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__lshrti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	u := &uint128{
		low:  uint64(frame.Locals[1]),
		high: uint64(frame.Locals[2]),
	}
	shift := uint(frame.Locals[3])
	u.rsh(shift)

	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf, u.low)
	binary.LittleEndian.PutUint64(buf[8:], u.high)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__lshrti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__divti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	ret := C.___divti3(
		C.uint64_t(frame.Locals[1]),
		C.uint64_t(frame.Locals[2]),
		C.uint64_t(frame.Locals[3]),
		C.uint64_t(frame.Locals[4]),
	)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__divti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__udivti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	ret := C.___udivti3(
		C.uint64_t(frame.Locals[1]),
		C.uint64_t(frame.Locals[2]),
		C.uint64_t(frame.Locals[3]),
		C.uint64_t(frame.Locals[4]),
	)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__udivti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__modti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	ret := C.___modti3(
		C.uint64_t(frame.Locals[1]),
		C.uint64_t(frame.Locals[2]),
		C.uint64_t(frame.Locals[3]),
		C.uint64_t(frame.Locals[4]),
	)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__modti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__umodti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	ret := C.___umodti3(
		C.uint64_t(frame.Locals[1]),
		C.uint64_t(frame.Locals[2]),
		C.uint64_t(frame.Locals[3]),
		C.uint64_t(frame.Locals[4]),
	)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__umodti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__multi3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	ret := C.___multi3(
		C.uint64_t(frame.Locals[1]),
		C.uint64_t(frame.Locals[2]),
		C.uint64_t(frame.Locals[3]),
		C.uint64_t(frame.Locals[4]),
	)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__multi3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__addtf3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(int32(frame.Locals[0]))
	la := uint64(frame.Locals[1])
	ha := uint64(frame.Locals[2])
	lb := uint64(frame.Locals[3])
	hb := uint64(frame.Locals[4])

	var a C.float128_t
	a.v[0] = C.uint64_t(la)
	a.v[1] = C.uint64_t(ha)
	var b C.float128_t
	b.v[0] = C.uint64_t(lb)
	b.v[1] = C.uint64_t(hb)

	sfRet := C.f128_add(a, b)
	buf := C.GoBytes(unsafe.Pointer(&sfRet), C.sizeof_float128_t)
	copy(vm.Memory.Memory[ret:ret+16], buf)
	return 0
}

func env__addtf3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__subtf3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])

	var a C.float128_t
	a.v[0] = C.uint64_t(frame.Locals[1])
	a.v[1] = C.uint64_t(frame.Locals[2])

	var b C.float128_t
	b.v[0] = C.uint64_t(frame.Locals[3])
	b.v[1] = C.uint64_t(frame.Locals[4])

	sfRet := C.f128_sub(a, b)
	buf := C.GoBytes(unsafe.Pointer(&sfRet), C.sizeof_float128_t)
	copy(vm.Memory.Memory[ret:ret+16], buf)
	return 0
}

func env__subtf3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__multf3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	var a C.float128_t
	a.v[0] = C.uint64_t(frame.Locals[1])
	a.v[1] = C.uint64_t(frame.Locals[2])
	var b C.float128_t
	b.v[0] = C.uint64_t(frame.Locals[3])
	b.v[1] = C.uint64_t(frame.Locals[4])

	sfRet := C.f128_mul(a, b)
	buf := C.GoBytes(unsafe.Pointer(&sfRet), C.sizeof_float128_t)
	copy(vm.Memory.Memory[ret:ret+16], buf)
	return 0
}

func env__multf3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__divtf3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	var a C.float128_t
	a.v[0] = C.uint64_t(frame.Locals[1])
	a.v[1] = C.uint64_t(frame.Locals[2])

	var b C.float128_t
	b.v[0] = C.uint64_t(frame.Locals[3])
	b.v[1] = C.uint64_t(frame.Locals[4])

	ret := C.f128_div(a, b)
	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof_float128_t)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__divtf3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

// conversion long double
func env__floatsitf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])
	i := C.int32_t(frame.Locals[1])

	ret := C.i32_to_f128(i)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof_float128_t)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__floatsitfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__floatunsitf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	i := C.uint32_t(frame.Locals[1])

	sfRet := C.ui32_to_f128(i)
	buf := C.GoBytes(unsafe.Pointer(&sfRet), C.sizeof_float128_t)
	copy(vm.Memory.Memory[ret:ret+16], buf)
	return 0
}

func env__floatunsitfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__floatditf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])
	a := C.int64_t(frame.Locals[1])

	ret := C.i64_to_f128(a)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof_float128_t)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__floatditfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__floatunditf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])
	a := C.uint64_t(frame.Locals[1])

	ret := C.ui64_to_f128(a)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof_float128_t)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__floatunditfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__floattidf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	l := C.uint64_t(frame.Locals[0])
	h := C.uint64_t(frame.Locals[1])

	d := C.___floattidf(l, h)
	return int64(d)
}

func env__floattidfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__floatuntidf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	l := C.uint64_t(frame.Locals[0])
	h := C.uint64_t(frame.Locals[1])

	d := C.___floatuntidf(l, h)
	return int64(d)
}

func env__floatuntidfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__floatsidf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()

	ret := C.i32_to_f64(C.int32_t(frame.Locals[0]))
	return int64(ret.v)
}

func env__floatsidfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__extendsftf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	var f C.float32_t
	f.v = C.uint32_t(frame.Locals[1])

	ret := C.f32_to_f128(f)
	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof_float128_t)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__extendsftf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__extenddftf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(int32(frame.Locals[0]))
	double := uint64(frame.Locals[1])

	var sf64 C.float64_t
	sf64.v = C.uint64_t(double)
	sf128 := C.f64_to_f128(sf64)

	buf := C.GoBytes(unsafe.Pointer(&sf128), C.sizeof_float128_t)
	copy(vm.Memory.Memory[ret:ret+16], buf)
	return 0
}

func env__extenddftf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__fixtfti(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	var f C.float128_t
	f.v[0] = C.uint64_t(frame.Locals[1])
	f.v[1] = C.uint64_t(frame.Locals[2])

	ret := C.___fixtfti(f)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__fixtftiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__fixtfdi(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()

	var f C.float128_t
	f.v[0] = C.uint64_t(frame.Locals[0])
	f.v[1] = C.uint64_t(frame.Locals[1])

	return int64(C.f128_to_i64(f, 1, false))
}

func env__fixtfdiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__fixtfsi(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	var f C.float128_t
	f.v[0] = C.uint64_t(frame.Locals[0])
	f.v[1] = C.uint64_t(frame.Locals[1])
	ret := C.f128_to_i32(f, 1, false)
	return int64(ret)
}

func env__fixtfsiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__fixunstfti(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	var f C.float128_t
	f.v[0] = C.uint64_t(frame.Locals[1])
	f.v[1] = C.uint64_t(frame.Locals[2])

	ret := C.___fixunstfti(f)

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__fixunstftiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__fixunstfdi(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()

	var f C.float128_t
	f.v[0] = C.uint64_t(frame.Locals[0])
	f.v[1] = C.uint64_t(frame.Locals[1])

	return int64(uint64(C.f128_to_ui64(f, 1, false)))
}

func env__fixunstfdiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__fixunstfsi(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	var f C.float128_t
	f.v[0] = C.uint64_t(frame.Locals[0])
	f.v[1] = C.uint64_t(frame.Locals[1])
	/*
	rounding:
	二进制值00: 近似到最近的偶数(默认)
	二进制值01: 向下近似趋向于-∞
	二进制值10: 向上近似趋向于+∞
	二进制值11: 近似趋向于0（剪裁)

	这里使用rounding使用1, 目的是保证__subtf3(a,b)中(a-b) >= 0.
	根据调试结果, 当rounding为0, double/long double转换为字符串时,可能
	会出现死循环.
	*/
	ret := uint32(C.f128_to_ui32(f, 1, false))
	return int64(ret)
}

func env__fixunstfsiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__fixsfti(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	ret := C.___fixsfti(C.uint32_t(frame.Locals[1]))

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__fixsftiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__fixdfti(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	ret := C.___fixdfti(C.uint64_t(frame.Locals[1]))

	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__fixdftiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__trunctfdf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()

	var f C.float128_t
	f.v[0] = C.uint64_t(frame.Locals[0])
	f.v[1] = C.uint64_t(frame.Locals[1])

	ret := C.f128_to_f64(f)

	return int64(ret.v)
}

func env__trunctfdf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__trunctfsf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()

	var f C.float128_t
	f.v[0] = C.uint64_t(frame.Locals[0])
	f.v[1] = C.uint64_t(frame.Locals[1])

	ret := C.f128_to_f32(f)
	return int64(ret.v)
}

func env__trunctfsf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__eqtf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := C.uint64_t(frame.Locals[0])
	ha := C.uint64_t(frame.Locals[1])
	lb := C.uint64_t(frame.Locals[2])
	hb := C.uint64_t(frame.Locals[3])
	return int64(cmptf2(la, ha, lb, hb, 1))
}

func env__eqtf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__netf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := C.uint64_t(frame.Locals[0])
	ha := C.uint64_t(frame.Locals[1])
	lb := C.uint64_t(frame.Locals[2])
	hb := C.uint64_t(frame.Locals[3])
	return int64(cmptf2(la, ha, lb, hb, 1))
}

func env__netf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__getf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := C.uint64_t(frame.Locals[0])
	ha := C.uint64_t(frame.Locals[1])
	lb := C.uint64_t(frame.Locals[2])
	hb := C.uint64_t(frame.Locals[3])
	return int64(cmptf2(la, ha, lb, hb, -1))
}

func env__getf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__gttf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := C.uint64_t(frame.Locals[0])
	ha := C.uint64_t(frame.Locals[1])
	lb := C.uint64_t(frame.Locals[2])
	hb := C.uint64_t(frame.Locals[3])
	return int64(cmptf2(la, ha, lb, hb, 0))
}

func env__gttf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__lttf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := C.uint64_t(frame.Locals[0])
	ha := C.uint64_t(frame.Locals[1])
	lb := C.uint64_t(frame.Locals[2])
	hb := C.uint64_t(frame.Locals[3])
	return int64(cmptf2(la, ha, lb, hb, 0))
}

func env__lttf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__letf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := C.uint64_t(frame.Locals[0])
	ha := C.uint64_t(frame.Locals[1])
	lb := C.uint64_t(frame.Locals[2])
	hb := C.uint64_t(frame.Locals[3])
	return int64(cmptf2(la, ha, lb, hb, 1))
}

func env__letf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__cmptf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := C.uint64_t(frame.Locals[0])
	ha := C.uint64_t(frame.Locals[1])
	lb := C.uint64_t(frame.Locals[2])
	hb := C.uint64_t(frame.Locals[3])
	return int64(cmptf2(la, ha, lb, hb, 1))
}

func env__cmptf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__unordtf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	var a C.float128_t
	a.v[0] = C.uint64_t(frame.Locals[0])
	a.v[1] = C.uint64_t(frame.Locals[1])

	var b C.float128_t
	b.v[0] = C.uint64_t(frame.Locals[2])
	b.v[1] = C.uint64_t(frame.Locals[3])

	if isNan(a) || isNan(b) {
		return 1
	}
	return 0
}

func env__unordtf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func env__negtf2(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])
	la := C.uint64_t(frame.Locals[1])
	ha := C.uint64_t(frame.Locals[2])

	var f C.float128_t
	f.v[0] = la
	f.v[1] = ha ^ (C.uint64_t(1) << 63)

	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf, uint64(f.v[0]))
	binary.LittleEndian.PutUint64(buf[8:], uint64(f.v[1]))
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

func env__negtf2GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func isNan(f C.float128_t) bool {
	return ((^(f.v[1]) & C.uint64_t(0x7FFF000000000000)) == 0) && (f.v[0] != 0 || ((f.v[1])&C.uint64_t(0x0000FFFFFFFFFFFF)) != 0)
}

func unordtf2(la, ha, lb, hb C.uint64_t) int {
	var a C.float128_t
	a.v[0] = la
	a.v[1] = ha
	var b C.float128_t
	b.v[0] = lb
	b.v[1] = hb
	if isNan(a) || isNan(b) {
		return 1
	}
	return 0
}

func cmptf2(la, ha, lb, hb C.uint64_t, returnValueIfNan int) int {
	var a C.float128_t
	a.v[0] = la
	a.v[1] = ha
	var b C.float128_t
	b.v[0] = lb
	b.v[1] = hb

	if unordtf2(la, ha, lb, hb) == 1 {
		return returnValueIfNan
	}
	if C.f128_lt(a, b) {
		return -1
	}
	if C.f128_eq(a, b) {
		return 0
	}
	return 1
}
