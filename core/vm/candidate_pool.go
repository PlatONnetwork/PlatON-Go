package vm

import (
	"Platon-go/common"
	"Platon-go/common/byteutil"
	"Platon-go/core/types"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"reflect"
)

// error def
var (
	ErrOwnerNotonly = errors.New("Node ID cannot bind multiple owners")
	ErrPermissionDenied = errors.New("Transaction from address permission denied")
	ErrDepositEmpyt = errors.New("Deposit balance not zero")
	ErrWithdrawEmpyt = errors.New("No withdrawal amount")
	ErrParamsRlpDecode = errors.New("Rlp decode faile")
	ErrParamsBaselen = errors.New("Params Base length does not match")
	ErrParamsLen = errors.New("Params length does not match")
	ErrUndefFunction = errors.New("Undefined function")
	ErrCandidateEmpyt = errors.New("CandidatePool is nil")
	ErrCallRecode = errors.New("Call recode error, panic...")
)

const (
	CandidateDepositEvent = "CandidateDepositEvent"
	CandidateApplyWithdrawEvent = "CandidateApplyWithdrawEvent"
	CandidateWithdrawEvent = "CandidateWithdrawEvent"
	SetCandidateExtraEvent = "SetCandidateExtraEvent"
)

type candidatePool interface {
	SetCandidate(state StateDB, nodeId discover.NodeID, can *types.Candidate) error
	GetCandidate(state StateDB, nodeId discover.NodeID) (*types.Candidate, error)
	WithdrawCandidate (state StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error
	GetChosens (state StateDB, ) []*types.Candidate
	GetChairpersons (state StateDB, ) []*types.Candidate
	GetDefeat(state StateDB, nodeId discover.NodeID) ([]*types.Candidate, error)
	IsDefeat(state StateDB, nodeId discover.NodeID) (bool, error)
	RefundBalance (state StateDB, nodeId discover.NodeID, blockNumber *big.Int) error
	GetOwner (state StateDB, nodeId discover.NodeID) common.Address
	SetCandidateExtra(state StateDB, nodeId discover.NodeID, extra string) error
	GetRefundInterval() uint64
}

type candidateContract struct{
	contract *Contract
	evm *EVM
}

func (c *candidateContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (c *candidateContract) Run(input []byte) ([]byte, error) {

	// debug
	logError("Run==> ", "input: ", hex.EncodeToString(input))

	defer func() {
		if err := recover(); nil != err {
			// catch call panic
			logError("Run==> ", "ErrCallRecode: ", ErrCallRecode.Error())
		}
	}()
	var command = map[string] interface{}{
		"CandidateDetails" : c.CandidateDetails,
		"CandidateApplyWithdraw" : c.CandidateApplyWithdraw,
		"CandidateDeposit" : c.CandidateDeposit,
		"CandidateList" : c.CandidateList,
		"CandidateWithdraw" : c.CandidateWithdraw,
		"SetCandidateExtra" : c.SetCandidateExtra,
		"CandidateWithdrawInfos": c.CandidateWithdrawInfos,
		"VerifiersList" : c.VerifiersList,
	}
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &source); err != nil {
		logError("Run==> ", err.Error())
		return nil, ErrParamsRlpDecode
	}
	// check
	if len(source)<2 {
		logError("Run==> ", "ErrParamsBaselen: ", ErrParamsBaselen.Error())
		return nil, ErrParamsBaselen
	}
	if c.evm.CandidatePool==nil{
		logError("Run==> ", "ErrCandidateEmpyt: ", ErrCandidateEmpyt.Error())
		return nil, ErrCandidateEmpyt
	}
	// get func and param list
	if _, ok := command[byteutil.BytesToString(source[1])]; !ok {
		logError("Run==> ", "ErrUndefFunction: ", ErrUndefFunction.Error())
		return nil, ErrUndefFunction
	}
	funcValue := command[byteutil.BytesToString(source[1])]
	paramList := reflect.TypeOf(funcValue)
	paramNum := paramList.NumIn()
	// var param []interface{}
	params := make([]reflect.Value, paramNum)
	if paramNum!=len(source)-2 {
		logError("Run==> ", "ErrParamsLen: ",ErrParamsLen.Error())
		return nil, ErrParamsLen
	}
	for i := 0; i < paramNum; i++ {
		targetType := paramList.In(i).String()
		originByte := []reflect.Value{reflect.ValueOf(source[i+2])}
		params[i] = reflect.ValueOf(byteutil.Command[targetType]).Call(originByte)[0]
	}
	// call func
	result := reflect.ValueOf(funcValue).Call(params)
	logInfo("Run==> ", "result[0]: ", result[0].Bytes())
	if _, err := result[1].Interface().(error); !err {
		return result[0].Bytes(), nil
	}
	logInfo(result[1].Interface().(error).Error())
	return result[0].Bytes(), result[1].Interface().(error)
}


// Candidate Application && Increase Quality Deposit
func (c *candidateContract) CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string) ([]byte, error)   {
	// debug
	deposit := c.contract.value
	txHash := c.evm.StateDB.TxHash()
	txIdx := c.evm.StateDB.TxIdx()
	height := c.evm.Context.BlockNumber
	from := c.contract.caller.Address()
	logInfo("CandidateDeposit==> ", "nodeId: ", nodeId.String(), " owner: ", owner.Hex(), " deposit: ", deposit,
		"  fee: ", fee, " txhash: ", txHash.Hex(), " txIdx: ", txIdx, " height: ", height, " from: ", from.Hex(),
		" host: ", host, " port: ", port, " extra: ", extra)
	//todo
	if deposit.Cmp(big.NewInt(0))<1 {
		return nil, ErrDepositEmpyt
	}
	can, err := c.evm.CandidatePool.GetCandidate(c.evm.StateDB, nodeId)
	if err!=nil {
		logError("CandidateDeposit==> ","err!=nill: ", err.Error())
		return nil, err
	}
	var alldeposit *big.Int
	if can!=nil {
		if ok := bytes.Equal(can.Owner.Bytes(), owner.Bytes()); !ok {
			logError(ErrOwnerNotonly.Error())
			return nil, ErrOwnerNotonly
		}
		alldeposit = new(big.Int).Add(can.Deposit, deposit)
		logInfo("CandidateDeposit==> ","alldeposit: ", alldeposit,  " can.Deposit: ", can.Deposit, " deposit: ", deposit)
	}else {
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
		0,
		new(big.Int).SetUint64(0),
		common.Hash{},
	}
	logInfo("CandidateDeposit==> ","canDeposit: ", canDeposit)
	if err = c.evm.CandidatePool.SetCandidate(c.evm.StateDB, nodeId, &canDeposit); err!=nil {
		// rollback transaction
		// ......
		logError("CandidateDeposit==> ","SetCandidate return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateDepositEvent, string(data))
	logInfo("CandidateDeposit==> ","json: ", string(data))
	return nil, nil
}

// Apply for a refund of the deposit
func (c *candidateContract) CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error)  {
	// debug
	txHash := c.evm.StateDB.TxHash()
	from := c.contract.caller.Address()
	height := c.evm.Context.BlockNumber
	logInfo("CandidateApplyWithdraw==> ","nodeId: ", nodeId.String(), " from: ", from.Hex(), " txHash: ", txHash.Hex(), " withdraw: ", withdraw, " height: ", height)
	// todo
	can, err := c.evm.CandidatePool.GetCandidate(c.evm.StateDB, nodeId)
	if err!=nil {
		logError("CandidateApplyWithdraw==> ","err!=nill: ", err.Error())
		return nil, err
	}
	if can.Deposit.Cmp(big.NewInt(0))<1 {
		return nil, ErrWithdrawEmpyt
	}
	if ok := bytes.Equal( can.Owner.Bytes(), from.Bytes()); !ok {
		logError(ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if withdraw.Cmp(can.Deposit)>0 {
		withdraw = can.Deposit
	}
	if err := c.evm.CandidatePool.WithdrawCandidate(c.evm.StateDB, nodeId, withdraw, height); err!=nil {
		logError(err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateApplyWithdrawEvent, string(data))
	logInfo("CandidateApplyWithdraw==> ","json: ", string(data))
	return nil, nil
}

// Deposit withdrawal
func (c *candidateContract) CandidateWithdraw(nodeId discover.NodeID) ([]byte, error)  {
	// debug
	txHash := c.evm.StateDB.TxHash()
	height := c.evm.Context.BlockNumber
	logInfo("CandidateWithdraw==> ","nodeId: ", nodeId.String(), " height: ", 	height, " txHash: ", txHash.Hex())
	// todo
	if err :=c.evm.CandidatePool.RefundBalance(c.evm.StateDB, nodeId, height); err!=nil{
		logError(err.Error())
		return nil, err
	}
	// return
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateWithdrawEvent, string(data))
	logInfo("CandidateWithdraw==> ", "json: ", string(data))
	return nil, nil
}

// Get the refund history you have applied for
func (c *candidateContract) CandidateWithdrawInfos(nodeId discover.NodeID)([]byte, error){
	// debug
	logInfo("CandidateWithdrawInfos==> ","nodeId: ", nodeId.String())
	// todo
	infos, err := c.evm.CandidatePool.GetDefeat(c.evm.StateDB, nodeId)
	if err!=nil{
		logError(err.Error())
		return nil, err
	}
	// return
	type WithdrawInfo struct {
		Balance *big.Int
		LockNumber *big.Int
		LockBlockCycle uint64
	}
	type WithdrawInfos struct {
		Ret bool
		ErrMsg string
		Infos []WithdrawInfo
	}
	r := WithdrawInfos{true, "success", make([]WithdrawInfo, len(infos))}
	for i, v := range infos {
		r.Infos[i] = WithdrawInfo{v.Deposit, v.BlockNumber, c.evm.CandidatePool.GetRefundInterval()}
	}
	data, _ := json.Marshal(r)
	sdata := DecodeResultStr(string(data))
	logInfo("CandidateWithdrawInfos==> ","json: ", string(data))
	return sdata, nil
}

// Set up additional information
func (c *candidateContract) SetCandidateExtra(nodeId discover.NodeID, extra string)([]byte, error){
	// debug
	txHash := c.evm.StateDB.TxHash()
	from := c.contract.caller.Address()
	logInfo("SetCandidate==> ","nodeId: ", nodeId.String(), " extra: ", extra, " from: ", from.Hex(), " txHash: ", txHash.Hex())
	// todo
	owner :=  c.evm.CandidatePool.GetOwner(c.evm.StateDB, nodeId)
	if ok := bytes.Equal(owner.Bytes(), from.Bytes()); !ok {
		logError(ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if err := c.evm.CandidatePool.SetCandidateExtra(c.evm.StateDB, nodeId, extra); err!=nil{
		logError(err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(SetCandidateExtraEvent, string(data))
	logInfo("SetCandidate==> ","json: ", string(data))
	return nil, nil
}

// Get candidate details
func (c *candidateContract) CandidateDetails(nodeId discover.NodeID) ([]byte, error)  {
	logInfo("CandidateDetails==> ","nodeId: ", nodeId.String())
	candidate, err := c.evm.CandidatePool.GetCandidate(c.evm.StateDB, nodeId)
	if err != nil{
		logError("CandidateDetails==> ","get CandidateDetails() occured error: ", err.Error())
		return nil, err
	}
	if nil == candidate {
		logError("CandidateDetails==> The candidate for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(candidate)
	sdata := DecodeResultStr(string(data))
	logInfo("CandidateDetails==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// Get the current block candidate list
func (c *candidateContract) CandidateList() ([]byte, error) {
	logInfo("CandidateList==> into func CandidateList... ")
	arr := c.evm.CandidatePool.GetChosens(c.evm.StateDB)
	if nil == arr {
		logError("CandidateList==> The candidateList for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(arr)
	sdata := DecodeResultStr(string(data))
	logInfo("CandidateList==>","json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// Get the current block round certifier list
func (c *candidateContract) VerifiersList() ([]byte, error) {
	logInfo("VerifiersList==> into func VerifiersList... ")
	arr := c.evm.CandidatePool.GetChairpersons(c.evm.StateDB)
	if nil == arr {
		logError("VerifiersList==> The verifiersList for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(arr)
	sdata := DecodeResultStr(string(data))
	logInfo("VerifiersList==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}