package resolver

/*
#cgo CFLAGS:-I .
#cgo LDFLAGS:-L ./libcsnark
#include "goLayer.h"
*/
import "C"
import (
	"github.com/PlatONnetwork/PlatON-Go/life/exec"
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

	return int64(0)
}

func envVerifyGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}
