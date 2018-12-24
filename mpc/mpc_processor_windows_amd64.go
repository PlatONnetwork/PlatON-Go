package mpc

// test part - no libso

/*
#include <stdio.h>
#include <stdlib.h>

void notify_security_init(const char* icecfg, const char* url) {
	//printf("init : icecfg : %s, url : %s", icecfg, url);
}

void notify_security_calculation(const char* taskid, const char* pubkey, const char* address, const char* ir_address, const char* method, const char* extra) {
	printf("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\n");
	printf("Received Params: taskId:%s, pubkey:%s, addr:%s, irAddr:%s, method:%s, extra:%s", taskid, pubkey, address, ir_address, method, extra);
}
*/
import "C"

// Main part, call library form

/*
#cgo LDFLAGS: -Wl,-rpath="./libs"
#cgo LDFLAGS: -L./libs
#cgo LDFLAGS: -ljuzixmpcvm_platonsdk_demo
#include <stdio.h>
#include <stdlib.h>

extern int notify_security_calculation(const char* taskid, const char* pubkey, const char* address, const char* ir_address, const char* method, const char* extra);
*/
//import "C"

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"fmt"
	"unsafe"
)

//void notify_security_calculation(const char* taskid, const char* pubkey, const char* address, const char* ir_address, const char* method, const char* extra)

type MPCParams struct {
	TaskId		string
	Pubkey 		string
	From 		common.Address
	IRAddr		common.Address
	Method 		string
	Extra 		string
}

func InitVM(icepath string, httpEndpoint string) {
	cCfg := C.CString(icepath)
	cUrl := C.CString(httpEndpoint)
	defer func() {
		C.free(unsafe.Pointer(cCfg))
		C.free(unsafe.Pointer(cUrl))
	}()
	C.notify_security_init(cCfg, cUrl)
	fmt.Println("mpc_process initVM method...")
	log.Info("Init mpc processor success", "osType", "window", "icepath", icepath, "httpEndpoint", httpEndpoint)
}

// for test
func ExecuteMPCTxForRedis(params MPCParams) (err error) {

	myRedis, err := NewRedis("192.168.9.14:6379")
	if err != nil {
		log.Error("Create connection of redis not success.", "err", err.Error())
		return err
	}

	var (
		KEY_TASK_ID = "taskId"
		KEY_PUB_KEY = "pubKey"
		KEY_ADDRESS = "address"
		KEY_IR_ADDRESS = "irAddress"
		KEY_METHOD = "method"
		KEY_EXTRA = "extra"
	)

	jsonMap := make(map[string]string)
	jsonMap[KEY_TASK_ID] = params.TaskId
	jsonMap[KEY_PUB_KEY] = params.Pubkey
	jsonMap[KEY_ADDRESS] = params.From.Hex()
	jsonMap[KEY_IR_ADDRESS] = params.IRAddr.Hex()
	jsonMap[KEY_METHOD] = params.Method
	jsonMap[KEY_EXTRA] = params.Extra

	err = myRedis.RPush(MPC_TASK_KEY_ALICE, jsonMap)
	if err != nil {
		fmt.Println("add mpc task to queue fail : -> to Alice")
		return err
	}
	myRedis.RPush(MPC_TASK_KEY_BOB, jsonMap)
	if err != nil {
		fmt.Println("add mpc task to queue fail : -> to Bob")
		return err
	}

	log.Trace("Notify mvm success, ExecuteMPCTx method invoke success.",
		"taskId", params.TaskId,
		"pubkey", params.Pubkey,
		"from", params.From.Hex(),
		"irAddr", params.IRAddr.Hex(),
		"method", params.Method)
	return nil
}

func ExecuteMPCTx(params MPCParams) error {

	ExecuteMPCTxForRedis(params)

	cTaskId := C.CString(params.TaskId)
	cPubKey := C.CString(params.Pubkey)
	cAddr := C.CString(params.From.Hex())
	cIRAddr := C.CString(params.IRAddr.Hex())
	cMethod := C.CString(params.Method)
	cExtra := C.CString(params.Extra)

	// call interface
	C.notify_security_calculation(cTaskId, cPubKey, cAddr, cIRAddr, cMethod, cExtra)

	defer func() {
		// free memory
		C.free(unsafe.Pointer(cTaskId))
		C.free(unsafe.Pointer(cPubKey))
		C.free(unsafe.Pointer(cAddr))
		C.free(unsafe.Pointer(cIRAddr))
		C.free(unsafe.Pointer(cMethod))
		C.free(unsafe.Pointer(cExtra))
	}()

	fmt.Printf("02->Received param, the taskId: %v, the pubkey: %v, the from: %v, the irAddr: %v, the method: %v, the extra: %v \n",
		params.TaskId, params.Pubkey, params.From.Hex(), params.IRAddr.Hex(), params.Method, params.Extra)

	log.Trace("Notify mvm success, ExecuteMPCTx method invoke success.",
		"taskId", params.TaskId,
		"pubkey", params.Pubkey,
		"from", params.From.Hex(),
		"irAddr", params.IRAddr.Hex(),
		"method", params.Method)

	return nil

}