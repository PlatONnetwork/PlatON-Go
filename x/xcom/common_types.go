package xcom

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

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

	IntermediateRoot(deleteEmptyObjects bool) common.Hash
}

// inner contract event data
type Result struct {
	Code   uint32
	Data   string
	ErrMsg string
}

func SuccessResult(data string) []byte {
	return BuildResult(data, common.NoErr)
}

func FailResult(data string, err *common.BizError) []byte {
	return BuildResult(data, err)
}

func BuildResult(data string, err *common.BizError) []byte {
	res := Result{err.Code, data, err.Msg}
	bs, _ := json.Marshal(res)
	return bs
}

func NewResult(data string, err *common.BizError) []byte {
	if err == nil {
		err = common.NoErr
	}
	res := &Result{err.Code, data, err.Msg}
	bs, _ := json.Marshal(res)
	return bs
}

var (
	OkResultByte, _ = json.Marshal(&Result{common.NoErr.Code, "", common.NoErr.Msg})
)

func NewSuccessResult(data string) []byte {
	res := &Result{common.NoErr.Code, data, common.NoErr.Msg}
	bs, _ := json.Marshal(res)
	return bs
}

func NewFailResultByBiz(err *common.BizError) []byte {
	res := &Result{err.Code, "", err.Msg}
	bs, _ := json.Marshal(res)
	return bs
}

func NewFailResult(err error) []byte {
	code, message := common.DecodeError(err)
	res := &Result{code, "", message}
	bs, _ := json.Marshal(res)
	return bs
}

// addLog let the result add to event.
func AddLog(state StateDB, blockNumber uint64, contractAddr common.Address, event, data string) {
	logdata := make([][]byte, 0)
	logdata = append(logdata, []byte(data))

	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); nil != err {
		log.Error("Cannot RlpEncode the log data, data", "data", data)
		panic("Cannot RlpEncode the log data")
	}

	//encoded := common.MustRlpEncode(logdata)

	state.AddLog(&types.Log{
		Address:     contractAddr,
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data:        buf.Bytes(),
		BlockNumber: blockNumber,
	})
}

func PrintObject(s string, obj interface{}) {
	objs, _ := json.Marshal(obj)
	log.Debug(s + " == " + string(objs))
}

func PrintObjForErr(s string, obj interface{}) {
	objs, _ := json.Marshal(obj)
	log.Error(s + " == " + string(objs))
}
