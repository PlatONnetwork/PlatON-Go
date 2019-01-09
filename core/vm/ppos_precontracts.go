package vm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"reflect"
)

// error def
var (
	ErrOwnerNotonly       = errors.New("Node ID cannot bind multiple owners")
	ErrPermissionDenied   = errors.New("Transaction from address permission denied")
	ErrFeeIllegal         = errors.New("The fee is illegal")
	ErrDepositEmpyt       = errors.New("Deposit balance not zero")
	ErrWithdrawEmpyt      = errors.New("No withdrawal amount")
	ErrParamsRlpDecode    = errors.New("Rlp decode faile")
	ErrParamsBaselen      = errors.New("Params Base length does not match")
	ErrParamsLen          = errors.New("Params length does not match")
	ErrUndefFunction      = errors.New("Undefined function")
	ErrTxType             = errors.New("Transaction type does not match the function")
	ErrCandidatePoolEmpyt = errors.New("CandidatePool is nil")
	ErrCandidateNotExist  = errors.New("The candidate is not exist")
)

const (
	CandidateDepositEvent       = "CandidateDepositEvent"
	CandidateApplyWithdrawEvent = "CandidateApplyWithdrawEvent"
	CandidateWithdrawEvent      = "CandidateWithdrawEvent"
	SetCandidateExtraEvent      = "SetCandidateExtraEvent"
)

var PrecompiledContractsPpos = map[common.Address]PrecompiledContract{
	common.CandidateAddr: &CandidateContract{},
}

type ResultCommon struct {
	Ret    bool
	ErrMsg string
}

type candidatePool interface {
	SetCandidate(state StateDB, nodeId discover.NodeID, can *types.Candidate) error
	GetCandidate(state StateDB, nodeId discover.NodeID) (*types.Candidate, error)
	WithdrawCandidate(state StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error
	GetChosens(state StateDB) []*types.Candidate
	GetChairpersons(state StateDB) []*types.Candidate
	GetDefeat(state StateDB, nodeId discover.NodeID) ([]*types.Candidate, error)
	IsDefeat(state StateDB, nodeId discover.NodeID) (bool, error)
	RefundBalance(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) error
	GetOwner(state StateDB, nodeId discover.NodeID) common.Address
	SetCandidateExtra(state StateDB, nodeId discover.NodeID, extra string) error
	GetRefundInterval() uint64
	MaxCount() uint64
}

type CandidateContract struct {
	Contract *Contract
	Evm      *EVM
}

func (c *CandidateContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (c *CandidateContract) Run(input []byte) ([]byte, error) {
	c.logInfo("Input to Run==> ", "input: ", hex.EncodeToString(input))
	defer func() {
		if err := recover(); nil != err {
			// catch call panic
			c.logError("Failed to Run==> ", "err", fmt.Sprint(err))
		}
	}()
	var command = map[string]interface{}{
		"CandidateDetails":       c.CandidateDetails,
		"CandidateApplyWithdraw": c.CandidateApplyWithdraw,
		"CandidateDeposit":       c.CandidateDeposit,
		"CandidateList":          c.CandidateList,
		"CandidateWithdraw":      c.CandidateWithdraw,
		"SetCandidateExtra":      c.SetCandidateExtra,
		"CandidateWithdrawInfos": c.CandidateWithdrawInfos,
		"VerifiersList":          c.VerifiersList,
	}
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &source); nil != err {
		c.logError("Failed to Run==> ", "ErrParamsRlpDecode: ", ErrParamsRlpDecode.Error())
		return nil, ErrParamsRlpDecode
	}
	// check
	if len(source) < 2 {
		c.logError("Failed to Run==> ", "ErrParamsBaselen: ", ErrParamsBaselen.Error())
		return nil, ErrParamsBaselen
	}
	if c.Evm.CandidatePool == nil {
		c.logError("Failed to Run==> ", "ErrCandidateEmpyt: ", ErrCandidatePoolEmpyt.Error())
		return nil, ErrCandidatePoolEmpyt
	}
	// get func and param list
	if _, ok := command[byteutil.BytesToString(source[1])]; !ok {
		c.logError("Failed to Run==> ", "ErrUndefFunction: ", ErrUndefFunction.Error())
		return nil, ErrUndefFunction
	}
	funcValue := command[byteutil.BytesToString(source[1])]
	// validate transaction type
	var txTypeMap = map[string]uint64{
		"CandidateDeposit":       1001,
		"CandidateApplyWithdraw": 1002,
		"CandidateWithdraw":      1003,
		"SetCandidateExtra":      1004,
	}
	if txType, ok := txTypeMap[byteutil.BytesToString(source[1])]; ok {
		if txType != binary.BigEndian.Uint64(source[0]) {
			log.Error("Failed to Run==> ", "ErrTxType: ", ErrTxType.Error())
			return nil, ErrTxType
		}
	}
	paramList := reflect.TypeOf(funcValue)
	paramNum := paramList.NumIn()
	// var param []interface{}
	params := make([]reflect.Value, paramNum)
	if paramNum != len(source)-2 {
		c.logError("Failed to Run==> ", "ErrParamsLen: ", ErrParamsLen.Error())
		return nil, ErrParamsLen
	}
	for i := 0; i < paramNum; i++ {
		targetType := paramList.In(i).String()
		originByte := []reflect.Value{reflect.ValueOf(source[i+2])}
		params[i] = reflect.ValueOf(byteutil.Command[targetType]).Call(originByte)[0]
	}
	// call func
	result := reflect.ValueOf(funcValue).Call(params)
	c.logInfo("Result of Run==> ", "result[0]: ", result[0].Bytes())
	if _, errOk := result[1].Interface().(error); !errOk {
		return result[0].Bytes(), nil
	}
	c.logError("Result of Run==> ", "result[1]: ", result[1].Interface().(error).Error())
	return result[0].Bytes(), result[1].Interface().(error)
}

// Candidate Application && Increase Quality Deposit.
func (c *CandidateContract) CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string) ([]byte, error) {
	deposit := c.Contract.value
	txHash := c.Evm.StateDB.TxHash()
	txIdx := c.Evm.StateDB.TxIdx()
	height := c.Evm.Context.BlockNumber
	from := c.Contract.caller.Address()
	c.logInfo("Input to CandidateDeposit==> ", "nodeId: ", nodeId.String(), " owner: ", owner.Hex(), " deposit: ", deposit,
		"  fee: ", fee, " txhash: ", txHash.Hex(), " txIdx: ", txIdx, " height: ", height, " from: ", from.Hex(),
		" host: ", host, " port: ", port, " extra: ", extra)
	if fee > 10000 {
		c.logError("Failed to CandidateDeposit==> ", "ErrFeeIllegal: ", ErrFeeIllegal.Error())
		return nil, ErrFeeIllegal
	}
	if deposit.Cmp(big.NewInt(0)) < 1 {
		c.logError("Failed to CandidateDeposit==> ", "ErrDepositEmpyt: ", ErrDepositEmpyt.Error())
		return nil, ErrDepositEmpyt
	}
	can, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if nil != err {
		c.logError("Failed to CandidateDeposit==> ", "GetCandidate return err: ", err.Error())
		return nil, err
	}
	var alldeposit *big.Int
	if nil != can {
		if ok := bytes.Equal(can.Owner.Bytes(), owner.Bytes()); !ok {
			c.logError("Failed to CandidateDeposit==> ", "ErrOwnerNotonly: ", ErrOwnerNotonly.Error())
			return nil, ErrOwnerNotonly
		}
		alldeposit = new(big.Int).Add(can.Deposit, deposit)
		c.logInfo("CandidateDeposit==> ", "alldeposit: ", alldeposit, " can.Deposit: ", can.Deposit, " deposit: ", deposit)
	} else {
		alldeposit = deposit
	}
	canDeposit := types.Candidate{
		alldeposit,
		height,
		txIdx,
		nodeId,
		host,
		port,
		owner,
		from,
		extra,
		fee,
	}
	c.logInfo("CandidateDeposit==> ", "canDeposit: ", canDeposit)
	if err = c.Evm.CandidatePool.SetCandidate(c.Evm.StateDB, nodeId, &canDeposit); nil != err {
		c.logError("Failed to CandidateDeposit==> ", "SetCandidate return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateDepositEvent, string(data))
	c.logInfo("Result of CandidateDeposit==> ", "json: ", string(data))
	return nil, nil
}

// Apply for a refund of the deposit.
func (c *CandidateContract) CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	from := c.Contract.caller.Address()
	height := c.Evm.Context.BlockNumber
	c.logInfo("Input to CandidateApplyWithdraw==> ", "nodeId: ", nodeId.String(), " from: ", from.Hex(), " txHash: ", txHash.Hex(), " withdraw: ", withdraw, " height: ", height)
	can, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if nil == can {
		c.logError("Failed to CandidateApplyWithdraw==> ", "ErrCandidateNotExist: ", ErrCandidateNotExist.Error())
		return nil, ErrCandidateNotExist
	}
	if nil != err {
		c.logError("Failed to CandidateApplyWithdraw==> ", "GetCandidate return err: ", err.Error())
		return nil, err
	}
	if can.Deposit.Cmp(big.NewInt(0)) < 1 {
		c.logError("Failed to CandidateApplyWithdraw==> ", "ErrWithdrawEmpyt: ", err.Error())
		return nil, ErrWithdrawEmpyt
	}
	if ok := bytes.Equal(can.Owner.Bytes(), from.Bytes()); !ok {
		c.logError("Failed to CandidateApplyWithdraw==> ", "ErrPermissionDenied: ", err.Error())
		return nil, ErrPermissionDenied
	}
	if withdraw.Cmp(can.Deposit) > 0 {
		withdraw = can.Deposit
	}
	if err := c.Evm.CandidatePool.WithdrawCandidate(c.Evm.StateDB, nodeId, withdraw, height); nil != err {
		c.logError("Failed to CandidateApplyWithdraw==> ", "WithdrawCandidate return err:", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateApplyWithdrawEvent, string(data))
	c.logInfo("Result of CandidateApplyWithdraw==> ", "json: ", string(data))
	return nil, nil
}

// Deposit withdrawal.
func (c *CandidateContract) CandidateWithdraw(nodeId discover.NodeID) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	height := c.Evm.Context.BlockNumber
	c.logInfo("Input to CandidateWithdraw==> ", "nodeId: ", nodeId.String(), " height: ", height, " txHash: ", txHash.Hex())
	if err := c.Evm.CandidatePool.RefundBalance(c.Evm.StateDB, nodeId, height); nil != err {
		c.logError("Failed to CandidateWithdraw==> ", "RefundBalance return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateWithdrawEvent, string(data))
	c.logInfo("Result of CandidateWithdraw==> ", "json: ", string(data))
	return nil, nil
}

// Get the refund history you have applied for.
func (c *CandidateContract) CandidateWithdrawInfos(nodeId discover.NodeID) ([]byte, error) {
	c.logInfo("Input to CandidateWithdrawInfos==> ", "nodeId: ", nodeId.String())
	infos, err := c.Evm.CandidatePool.GetDefeat(c.Evm.StateDB, nodeId)
	if nil != err {
		c.logError("Failed to CandidateWithdrawInfos==> ", "GetDefeat return err: ", err.Error())
		return nil, err
	}
	type WithdrawInfo struct {
		Balance        *big.Int
		LockNumber     *big.Int
		LockBlockCycle uint64
	}
	type WithdrawInfos struct {
		Ret    bool
		ErrMsg string
		Infos  []WithdrawInfo
	}
	r := WithdrawInfos{true, "success", make([]WithdrawInfo, len(infos))}
	for i, v := range infos {
		r.Infos[i] = WithdrawInfo{v.Deposit, v.BlockNumber, c.Evm.CandidatePool.GetRefundInterval()}
	}
	data, _ := json.Marshal(r)
	sdata := DecodeResultStr(string(data))
	c.logInfo("Result of CandidateWithdrawInfos==> ", "json: ", string(data))
	return sdata, nil
}

// Set up additional information.
func (c *CandidateContract) SetCandidateExtra(nodeId discover.NodeID, extra string) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	from := c.Contract.caller.Address()
	c.logInfo("Input to SetCandidateExtra==> ", "nodeId: ", nodeId.String(), " extra: ", extra, " from: ", from.Hex(), " txHash: ", txHash.Hex())
	owner := c.Evm.CandidatePool.GetOwner(c.Evm.StateDB, nodeId)
	if ok := bytes.Equal(owner.Bytes(), from.Bytes()); !ok {
		c.logError("Failed to SetCandidateExtra==> ", "ErrPermissionDenied: ", ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if err := c.Evm.CandidatePool.SetCandidateExtra(c.Evm.StateDB, nodeId, extra); err != nil {
		c.logError("Failed to SetCandidateExtra==> ", "SetCandidateExtra return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(SetCandidateExtraEvent, string(data))
	c.logInfo("Result of SetCandidateExtra==> ", "json: ", string(data))
	return nil, nil
}

// Get candidate details.
func (c *CandidateContract) CandidateDetails(nodeId discover.NodeID) ([]byte, error) {
	c.logInfo("Input to CandidateDetails==> ", "nodeId: ", nodeId.String())
	candidate, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if nil != err {
		c.logError("Failed to CandidateDetails==> ", "GetCandidate return err: ", err.Error())
		return nil, err
	}
	if nil == candidate {
		c.logError("Failed to CandidateDetails==> ", "The query does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(candidate)
	sdata := DecodeResultStr(string(data))
	c.logInfo("Result of CandidateDetails==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// Get the current block candidate list.
func (c *CandidateContract) CandidateList() ([]byte, error) {
	arr := c.Evm.CandidatePool.GetChosens(c.Evm.StateDB)
	if 0 == len(arr) {
		c.logError("Failed to CandidateList==> ", "The query does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(arr)
	sdata := DecodeResultStr(string(data))
	c.logInfo("Result of CandidateList==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// Get the current block round certifier list.
func (c *CandidateContract) VerifiersList() ([]byte, error) {
	arr := c.Evm.CandidatePool.GetChairpersons(c.Evm.StateDB)
	if 0 == len(arr) {
		c.logError("Failed to VerifiersList==> ", "The query does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(arr)
	sdata := DecodeResultStr(string(data))
	c.logInfo("Result of VerifiersList==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// transaction add event.
func (c *CandidateContract) addLog(event, data string) {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); err != nil {
		c.logError("addlog Err==> ", "rlp encode fail: ", err.Error())
	}
	c.Evm.StateDB.AddLog(&types.Log{
		Address:     common.CandidateAddr,
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data:        buf.Bytes(),
		BlockNumber: c.Evm.Context.BlockNumber.Uint64(),
	})
}

//debug log
func (c *CandidateContract) logInfo(msg string, ctx ...interface{}) {
	log.Info(msg, ctx...)
	//args := []interface{}{msg}
	//args = append(args, ctx...)
	//fmt.Println(args...)
	/*if c.evm.vmConfig.ConsoleOutput {
		//console output
		args := []interface{}{msg}
		args = append(args, ctx...)
		fmt.Println(args...)
	}else {
		//log output
		log.Info(msg, ctx...)
	}*/
}
func (c *CandidateContract) logError(msg string, ctx ...interface{}) {
	log.Error(msg, ctx...)
	//args := []interface{}{msg}
	//args = append(args, ctx...)
	//fmt.Println(args...)
	/*if c.evm.vmConfig.ConsoleOutput {
		//console output
		args := []interface{}{msg}
		args = append(args, ctx...)
		fmt.Println(args...)
	}else {
		//log output
		log.Error(msg, ctx...)
	}*/
}
func (c *CandidateContract) logPrint(level log.Lvl, msg string, ctx ...interface{}) {
	if c.Evm.vmConfig.ConsoleOutput {
		//console output
		args := make([]interface{}, len(ctx)+1)
		args[0] = msg
		for i, v := range ctx {
			args[i+1] = v
		}
		fmt.Println(args...)
	} else {
		//log output
		switch level {
		case log.LvlCrit:
			log.Crit(msg, ctx...)
		case log.LvlError:
			log.Error(msg, ctx...)
		case log.LvlWarn:
			log.Warn(msg, ctx...)
		case log.LvlInfo:
			log.Info(msg, ctx...)
		case log.LvlDebug:
			log.Debug(msg, ctx...)
		case log.LvlTrace:
			log.Trace(msg, ctx...)
		}
	}
}

//return string format
func DecodeResultStr(result string) []byte {
	// 0x0000000000000000000000000000000000000020
	// 00000000000000000000000000000000000000000d
	// 00000000000000000000000000000000000000000

	resultBytes := []byte(result)
	strHash := common.BytesToHash(common.Int32ToBytes(32))
	sizeHash := common.BytesToHash(common.Int64ToBytes(int64((len(resultBytes)))))
	var dataRealSize = len(resultBytes)
	if (dataRealSize % 32) != 0 {
		dataRealSize = dataRealSize + (32 - (dataRealSize % 32))
	}
	dataByt := make([]byte, dataRealSize)
	copy(dataByt[0:], resultBytes)

	finalData := make([]byte, 0)
	finalData = append(finalData, strHash.Bytes()...)
	finalData = append(finalData, sizeHash.Bytes()...)
	finalData = append(finalData, dataByt...)
	//encodedStr := hex.EncodeToString(finalData)
	//fmt.Println("finalData: ", encodedStr)
	return finalData
}
