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
	Status  bool
	Data    string
	ErrCode uint16
	ErrMsg  string
}

func SuccessResult(data string, errMsg string) []byte {
	return BuildResult(true, data, common.Success)
}

func FailResult(data string, errMsg string) []byte {
	return BuildResult(false, data, common.InternalError.Wrap(errMsg))
}

func BuildResult(status bool, data string, err *common.BizError) []byte {
	res := Result{status, data, err.Code, err.Msg}
	bs, _ := json.Marshal(res)
	return bs
}

func NewResult(data string, err *common.BizError) []byte {
	if err == nil {
		err = common.Success
	}
	res := &Result{err.Code == 0, data, err.Code, err.Msg}
	bs, _ := json.Marshal(res)
	return bs
}

var (
	NewDefaultSuccessResult, _ = json.Marshal(&Result{true, "", common.Success.Code, common.Success.Msg})
)

func NewSuccessResult(data string) []byte {
	res := &Result{true, data, common.Success.Code, common.Success.Msg}
	bs, _ := json.Marshal(res)
	return bs
}

func NewFailResult(err error) []byte {
	code, message := common.DecodeError(err)
	res := &Result{false, "", code, message}
	bs, _ := json.Marshal(res)
	return bs
}

func NewFailResultString(errorMessage string) []byte {
	err := common.InternalError.Wrap(errorMessage)
	res := &Result{false, "", err.Code, err.Msg}
	bs, _ := json.Marshal(res)
	return bs
}

/*// EncodeRLP implements rlp.Encoder
func (r *Result) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, Result{
		Status:    	r.Status,
		Data:       r.Data,
		ErrMsg: 	r.ErrMsg,
	})
}


// DecodeRLP implements rlp.Decoder
func (r *Result) DecodeRLP(s *rlp.Stream) error {
	var rs Result
	if err := s.Decode(&rs); err != nil {
		return err
	}

	ty := reflect.ValueOf(r.Data).Elem()

	if dByte, err := rlp.EncodeToBytes(r.Data); nil != err {
		return err
	}else {
		if err := rlp.DecodeBytes(dByte, &ty); nil != err {
			return err
		}
	}
	r.Status, r.Data, r.ErrMsg = rs.Status, ty, rs.ErrMsg
	return nil
}*/

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
