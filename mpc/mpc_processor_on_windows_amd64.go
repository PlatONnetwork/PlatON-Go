// +build mpcon

package mpc

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	psdkdll					*syscall.DLL
	securityInitProc		*syscall.Proc
	securityCalculation		*syscall.Proc
)

type MPCParams struct {
	TaskId		string
	Pubkey 		string
	From 		common.Address
	IRAddr		common.Address
	Method 		string
	Extra 		string
}

func InitVM(icepath string, httpEndpoint string) {
	// set default path of mpcLib.
	workDir, _ := filepath.Abs("")
	libPath := filepath.Join(workDir, "mpclib")

	// init ptr
	libPathPtr, _ := syscall.UTF16PtrFromString(libPath)
	syscall.SetCurrentDirectory(libPathPtr)
	psdkdll = syscall.MustLoadDLL("mpc_vm_platonsdk.dll")
	securityInitProc = psdkdll.MustFindProc("notify_security_init")
	securityCalculation = psdkdll.MustFindProc("notify_security_calculation")

	// convert type
	cCfg := C.CString(icepath)
	cUrl := C.CString(httpEndpoint)

	cfgUintPtr := uintptr(unsafe.Pointer(cCfg))
	urlUintPtr := uintptr(unsafe.Pointer(cUrl))
	syscall.Syscall(securityInitProc.Addr(), 2, cfgUintPtr, urlUintPtr, 0)

	fmt.Println("mpc_process initVM method...")
	log.Info("Init mpc processor success", "osType", "window", "icepath", icepath, "httpEndpoint", httpEndpoint)
	defer func() {
		C.free(unsafe.Pointer(cCfg))
		C.free(unsafe.Pointer(cUrl))
	}()
}

func ExecuteMPCTx(params MPCParams) error {

	defer func() {
		if err := recover(); err != nil {
			log.Error("execute mpc tx fail.", "err", err)
		}
	}()

	cTaskId := C.CString(params.TaskId)
	cPubKey := C.CString(params.Pubkey)
	cAddr := C.CString(params.From.Hex())
	cIRAddr := C.CString(params.IRAddr.Hex())
	cMethod := C.CString(params.Method)
	cExtra := C.CString(params.Extra)

	// convert to uintptr
	taskIdUintPtr := uintptr(unsafe.Pointer(cTaskId))
	pubKeyUintPtr := uintptr(unsafe.Pointer(cPubKey))
	addrUintPtr   := uintptr(unsafe.Pointer(cAddr))
	irAddrUintPtr := uintptr(unsafe.Pointer(cIRAddr))
	methodUintPtr := uintptr(unsafe.Pointer(cMethod))
	extraUintPtr  := uintptr(unsafe.Pointer(cExtra))

	syscall.Syscall6(securityCalculation.Addr(), 6, taskIdUintPtr, pubKeyUintPtr, addrUintPtr, irAddrUintPtr, methodUintPtr, extraUintPtr)

	defer func() {
		C.free(unsafe.Pointer(cTaskId))
		C.free(unsafe.Pointer(cPubKey))
		C.free(unsafe.Pointer(cAddr))
		C.free(unsafe.Pointer(cIRAddr))
		C.free(unsafe.Pointer(cMethod))
		C.free(unsafe.Pointer(cExtra))
	}()

	log.Info("Notify mvm success, ExecuteMPCTx method invoke success.",
		"taskId", params.TaskId,
		"pubkey", params.Pubkey,
		"from", params.From.Hex(),
		"irAddr", params.IRAddr.Hex(),
		"method", params.Method)

	return nil
}