package resolver

/*
#cgo CFLAGS:-I .
#cgo LDFLAGS:-L ./libcsnark
#include "goLayer.h"
*/
import "C"
import (
	"github.com/PlatONnetwork/PlatON-Go/life/exec"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

// define: void vc_InitGadgetEnv();
func envInitGadgetEnv(vm *exec.VirtualMachine) int64 {
	return 0
}

func envInitGadgetEnvGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: void vc_UninitGadgetEnv();
func envUninitGadgetEnv(vm *exec.VirtualMachine) int64 {
	return 0
}

func envUninitGadgetEnvGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: int64_t vc_CreatePBVar(void *varAddr);
func envCreatePBVarEnv(vm *exec.VirtualMachine) int64 {

	return 0
}

func envCreatePBVarGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: uint8_t vc_CreateGadget(void *input0, void *input1,
//                 void *input2, void *res, int32_t Type);
func envCreateGadgetEnv(vm *exec.VirtualMachine) int64 {

	return int64(0)
}

func envCreateGadgetGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: void vc_SetVar(void *var, uint64_t Val);
func envSetVarEnv(vm *exec.VirtualMachine) int64 {

	return 0
}

func envSetVarGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: void vc_SetRetIndex(int64_t RetAddr);
func envSetRetIndexEnv(vm *exec.VirtualMachine) int64 {
	return 0
}

func envSetRetIndexGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: void vc_GenerateWitness();
func envGenWitnessEnv(vm *exec.VirtualMachine) int64 {
	return 0
}

func envGenWitnessGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: uint8_t vc_GenerateProofAndResult(const char *pPKEY, int32_t pkSize, char *pProof,
//									int32_t prSize, char *pResult, int32_t resSize);
func envGenProofAndResultEnv(vm *exec.VirtualMachine) int64 {

	return int64(0)
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
	log.Debug("Vc Verify:", go_vk, go_pr, go_in, go_out)
	//c_vk := C.CString(go_vk)
	//c_pr := C.CString(go_pr)
	//c_in := C.CString(go_in)
	//c_out := C.CString(go_out)

	// call c func
	//retVal := uint8(C.Verify(c_vk, c_pr, c_in, c_out))
	retVal := 1

	// release memory
	//defer C.free(unsafe.Pointer(c_vk))
	//defer C.free(unsafe.Pointer(c_pr))
	//defer C.free(unsafe.Pointer(c_in))
	//defer C.free(unsafe.Pointer(c_out))
	return int64(retVal)
}

func envVerifyGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}
