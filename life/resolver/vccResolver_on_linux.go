// +build vcon

package resolver

/*
#cgo CFLAGS:-I .
#cgo LDFLAGS:-L ./libcsnark -lcsnark -lsnark -lff -lm -lgmp -lgmpxx -lcrypto -lprocps -lstdc++
#include "goLayer.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/PlatONnetwork/PlatON-Go/life/exec"
)

// define: void vc_InitGadgetEnv();
func envInitGadgetEnv(vm *exec.VirtualMachine) int64 {
	fmt.Println("begin init gadget env")
	C.gadget_initEnv()
	fmt.Println("end init gadget env")
	return 0
}

func envInitGadgetEnvGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: void vc_UninitGadgetEnv();
func envUninitGadgetEnv(vm *exec.VirtualMachine) int64 {
	fmt.Println("begin uninit gadget env")
	C.gadget_uninitEnv()
	fmt.Println("end uninit gadget env")
	return 0
}

func envUninitGadgetEnvGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: int64_t vc_CreatePBVar(void *varAddr);
func envCreatePBVarEnv(vm *exec.VirtualMachine) int64 {
	// get parameters
	varAddr := int64(vm.GetCurrentFrame().Locals[0])
	cvarAddr := C.longlong(varAddr)

	// call c func
	C.gadget_createPBVar(cvarAddr)

	return 0
}

func envCreatePBVarGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: uint8_t vc_CreateGadget(void *input0, void *input1,
//                 void *input2, void *res, int32_t Type);
func envCreateGadgetEnv(vm *exec.VirtualMachine) int64 {
	// get parameters
	input0Addr := int64(vm.GetCurrentFrame().Locals[0])
	input1Addr := int64(vm.GetCurrentFrame().Locals[1])
	input2Addr := int64(vm.GetCurrentFrame().Locals[2])
	resAddr := int64(vm.GetCurrentFrame().Locals[3])
	gType := int32(vm.GetCurrentFrame().Locals[4])

	cinput0 := C.longlong(input0Addr)
	cinput1 := C.longlong(input1Addr)
	cinput2 := C.longlong(input2Addr)
	cres := C.longlong(resAddr)

	// call c func
	retVal := uint8(C.gadget_createGadget(cinput0, cinput1, cinput2, cres, C.int(gType)))

	return int64(retVal)
}

func envCreateGadgetGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: void vc_SetVar(void *var, uint64_t Val);
func envSetVarEnv(vm *exec.VirtualMachine) int64 {
	// get parameters
	varAddr := int64(vm.GetCurrentFrame().Locals[0])
	varVal := int64(vm.GetCurrentFrame().Locals[1])
	varUnsign := int8(vm.GetCurrentFrame().Locals[2])
	cvarAddr := C.longlong(varAddr)

	// call c func
	C.gadget_setVar(cvarAddr, C.longlong(varVal), C.uchar(varUnsign))

	return 0
}

func envSetVarGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: void vc_SetRetIndex(int64_t RetAddr);
func envSetRetIndexEnv(vm *exec.VirtualMachine) int64 {
	// get parameters
	retAddr := int64(vm.GetCurrentFrame().Locals[0])
	cretAddr := C.longlong(retAddr)

	// call c func
	C.gadget_setRetIndex(cretAddr)

	return 0
}

func envSetRetIndexGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: void vc_GenerateWitness();
func envGenWitnessEnv(vm *exec.VirtualMachine) int64 {
	C.gadget_generateWitness()
	return 0
}

func envGenWitnessGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: uint8_t vc_GenerateProofAndResult(const char *pPKEY, int32_t pkSize, char *pProof,
//									int32_t prSize, char *pResult, int32_t resSize);
func envGenProofAndResultEnv(vm *exec.VirtualMachine) int64 {
	// get parameter
	pkOffset := int(int32(vm.GetCurrentFrame().Locals[0]))
	pkSize := int(int32(vm.GetCurrentFrame().Locals[1]))
	prOffset := int(int32(vm.GetCurrentFrame().Locals[2]))
	prSize := int(int32(vm.GetCurrentFrame().Locals[3]))
	resOffset := int(int32(vm.GetCurrentFrame().Locals[4]))
	resSize := int(int32(vm.GetCurrentFrame().Locals[5]))
	pkData := vm.Memory.Memory[pkOffset : pkOffset+pkSize]
	prData := vm.Memory.Memory[prOffset : prOffset+prSize]
	resData := vm.Memory.Memory[resOffset : resOffset+resSize]
	go_pk := string(pkData[:])
	c_pk := C.CString(go_pk)
	go_pr := string(prData[:])
	c_pr := C.CString(go_pr)
	go_res := string(resData[:])
	c_res := C.CString(go_res)

	// call c function
	retVal := C.GenerateProofAndResult(c_pk, c_pr, C.int(prSize), c_res, C.int(resSize))

	// copy
	proof := C.GoString(c_pr)
	result := C.GoString(c_res)
	copy(vm.Memory.Memory[prOffset:], proof)
	copy(vm.Memory.Memory[resOffset:], result)

	// release memory
	defer C.free(unsafe.Pointer(c_pk))
	defer C.free(unsafe.Pointer(c_pr))
	defer C.free(unsafe.Pointer(c_res))

	return int64(retVal)
}

func envGenProofAndResultGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: uint8_t vc_Verify(const char *pVKEY, int32_t pkSize, const char *pPoorf, int32_t prSize,
//					const char *pInput, int32_t inSize, const char *pOutput, int32_t outSize);
func envVerifyEnv(vm *exec.VirtualMachine) int64 {
	// get parameters
	vkOffset := int(int32(vm.GetCurrentFrame().Locals[0]))
	vkSize := int(int32(vm.GetCurrentFrame().Locals[1]))
	prOffset := int(int32(vm.GetCurrentFrame().Locals[2]))
	prSize := int(int32(vm.GetCurrentFrame().Locals[3]))
	inOffset := int(int32(vm.GetCurrentFrame().Locals[4]))
	inSize := int(int32(vm.GetCurrentFrame().Locals[5]))
	outOffset := int(int32(vm.GetCurrentFrame().Locals[6]))
	outSize := int(int32(vm.GetCurrentFrame().Locals[7]))
	vkData := vm.Memory.Memory[vkOffset : vkOffset+vkSize]
	prData := vm.Memory.Memory[prOffset : prOffset+prSize]
	inData := vm.Memory.Memory[inOffset : inOffset+inSize]
	outData := vm.Memory.Memory[outOffset : outOffset+outSize]
	go_vk := string(vkData[:])
	go_pr := string(prData[:])
	go_in := string(inData[:])
	go_out := string(outData[:])
	c_vk := C.CString(go_vk)
	c_pr := C.CString(go_pr)
	c_in := C.CString(go_in)
	c_out := C.CString(go_out)

	// call c func
	retVal := uint8(C.Verify(c_vk, c_pr, c_in, c_out))

	// release memory
	defer C.free(unsafe.Pointer(c_vk))
	defer C.free(unsafe.Pointer(c_pr))
	defer C.free(unsafe.Pointer(c_in))
	defer C.free(unsafe.Pointer(c_out))
	return int64(retVal)
}

func envVerifyGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}
