package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	gerr "github.com/go-errors/errors"
	"math/big"
	"reflect"
)

type BasePlugin interface {
	BeginBlock(header *types.Header, state StateDB) (bool, error)
	EndBlock(header *types.Header, state StateDB) (bool, error)
	Confirmed(block *types.Block) error
	//Verify_tx_data(source [][]byte) (fn interface{}, FnParams []reflect.Value, err error)
}


// StateDB is an Plugin database for full state querying.
type StateDB interface {
	CreateAccount(common.Address)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	// todo: new func for abi of contract.
	GetAbiHash(common.Address) common.Hash
	GetAbi(common.Address) []byte
	SetAbi(common.Address, []byte)

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	// todo: hash -> bytes
	GetCommittedState(common.Address, []byte) []byte
	//GetState(common.Address, common.Hash) common.Hash
	//SetState(common.Address, common.Hash, common.Hash)
	GetState(common.Address, []byte) []byte
	SetState(common.Address, []byte, []byte)


	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)

	//ppos add
	TxHash() common.Hash
	TxIdx() uint32
}




var (
	DecodeTxDataErr = errors.New("decode tx data is err")
	FuncNotExistErr = errors.New("the func is not exist")
	FnParamsLenErr 	= errors.New("the params len and func params len is not equal")
)

func  Verify_tx_data(input []byte, command map[uint16]interface{} ) (fn interface{}, FnParams []reflect.Value, err error)  {

	defer func() {
		if er := recover(); nil != er {
			fn, FnParams, err = nil, nil,fmt.Errorf("parse tx data is panic: %s, txHash: %s", gerr.Wrap(er, 2).ErrorStack())
			log.Error("Failed to Verify Staking tx data", "error", err)
		}
	}()

	var args [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &args); nil != err {
		return nil, nil, DecodeTxDataErr
	}

	if fn, ok := command[byteutil.BytesToUint16(args[0])]; !ok {
			return nil, nil, FuncNotExistErr
	}else {

		// the func params type list
		paramList := reflect.TypeOf(fn)
		// the func params len
		paramNum := paramList.NumIn()

		if paramNum != len(args)-1 {
			return nil, nil, FnParamsLenErr
		}

		params := make([]reflect.Value, paramNum)

		for i := 0; i < paramNum; i++ {
			targetType := paramList.In(i).String()
			inputByte := []reflect.Value{reflect.ValueOf(args[i+1])}
			params[i] = reflect.ValueOf(byteutil.Command[targetType]).Call(inputByte)[0]
		}

		return fn, params, nil
	}
}