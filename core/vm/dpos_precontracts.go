/**********************定义*********************************
内置合约地址：
0x1000..0011		候选池内置合约
0x1000..0012		选票池内置合约
0x1000..0010 + x	其他定义的内置合约

交易data字段定义：
data = rlp(type [8]byte, funcname string, parma1 []byte, parma2 []byte, ...)

候选池合约：
0x11 + "funcname" + parmas

选票池合约：
0x12 + "funcname" + parmas

其他合约：
(0x10+x) + "funcname" + parmas

**********************定义*********************************/
package vm

import (
	"Platon-go/common"
	"Platon-go/common/byteutil"
	"Platon-go/core/types"
	"Platon-go/crypto"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"Platon-go/rlp"
	"Platon-go/log"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"

)

//error def
var (
	ErrOwnerNotonly = errors.New("Node ID cannot bind multiple owners")
	ErrPermissionDenied = errors.New("Transaction from address permission denied")
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

var PrecompiledContractsDpos = map[common.Address]PrecompiledContract{
	common.CandidateAddr : &candidateContract{},
}

type ResultCommon struct {
	Ret bool
	ErrMsg string
}

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
	defer func() {
		if err := recover(); nil != err {
			// 捕捉反射解析参数时由input中数据源问题造成的panic
			log.Error("Run==> ", ErrCallRecode.Error())
		}
	}()
	// 用map封装所有的函数
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
		log.Error("Run==> ", err.Error())
		return nil, ErrParamsRlpDecode
	}
	//check
	if len(source)<2 {
		log.Error("Run==> ", ErrParamsBaselen.Error())
		return nil, ErrParamsBaselen
	}
	if c.evm.CandidatePool==nil{
		log.Error("Run==> ", ErrCandidateEmpyt.Error())
		return nil, ErrCandidateEmpyt
	}
	// 获取要调用的函数
	if _, ok := command[byteutil.BytesToString(source[1])]; !ok {
		log.Error("Run==> ", ErrUndefFunction.Error())
		return nil, ErrUndefFunction
	}
	funcValue := command[byteutil.BytesToString(source[1])]
	// 目标函数参数列表
	paramList := reflect.TypeOf(funcValue)
	// 目标函数参数个数
	paramNum := paramList.NumIn()
	// var param []interface{}
	params := make([]reflect.Value, paramNum)
	if paramNum!=len(source)-2 {
		log.Error("Run==> ", ErrParamsLen.Error())
		return nil, ErrParamsLen
	}
	for i := 0; i < paramNum; i++ {
		// 目标参数类型的值
		targetType := paramList.In(i).String()
		// 原始[]byte类型参数
		originByte := []reflect.Value{reflect.ValueOf(source[i+2])}
		// 转换为对应类型的参数
		params[i] = reflect.ValueOf(byteutil.Command[targetType]).Call(originByte)[0]
	}
	// 传入参数调用函数
	result := reflect.ValueOf(funcValue).Call(params)
	if _, err := result[1].Interface().(error); !err {
		log.Error("Run==> ", err)
		return result[0].Bytes(), nil
	}
	log.Info("result[0] is: ", result[0].Bytes())
	log.Info(result[1].Interface().(error).Error())
	// 返回值也是一个 Value 的 slice，同样对应反射函数类型的返回值。
	return result[0].Bytes(), result[1].Interface().(error)
}

//Candidate Application && Increase Quality Deposit
func (c *candidateContract) CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string) ([]byte, error)   {
	//debug
	deposit := c.contract.value
	txHash := c.evm.StateDB.TxHash()
	txIdx := c.evm.StateDB.TxIdx()
	height := c.evm.Context.BlockNumber
	from := c.contract.caller.Address()
	log.Info("CandidateDeposit==> nodeId: ", nodeId.String(), " owner: ", owner.Hex(), " deposit: ", deposit,
		"  fee: ", fee, " txhash: ", txHash.Hex(), " txIdx: ", txIdx, " height: ", height, " from: ", from.Hex(),
		" host: ", host, " port: ", port, " extra: ", extra)
	//todo
	can, err := c.evm.CandidatePool.GetCandidate(c.evm.StateDB, nodeId)
	if err!=nil {
		log.Error("GetCandidate err!=nill: ", err.Error())
		return nil, err
	}
	var alldeposit *big.Int
	if can!=nil {
		if ok := bytes.Equal(can.Owner.Bytes(), owner.Bytes()); !ok {
			log.Error(ErrOwnerNotonly.Error())
			return nil, ErrOwnerNotonly
		}
		alldeposit = new(big.Int).Add(can.Deposit, deposit)
		log.Info("alldeposit: ", alldeposit,  " can.Deposit: ", can.Deposit, " deposit: ", deposit)
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
	}
	log.Info("canDeposit: ", canDeposit)
	if err = c.evm.CandidatePool.SetCandidate(c.evm.StateDB, nodeId, &canDeposit); err!=nil {
		//rollback transaction
		//......
		log.Error("SetCandidate return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateDepositEvent, string(data))
	log.Info("json: ", string(data))
	return nil, nil
}

//Apply for a refund of the deposit
func (c *candidateContract) CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error)  {
	//debug
	from := c.contract.caller.Address()
	height := c.evm.Context.BlockNumber
	log.Info("CandidateApplyWithdraw==> nodeId: ", nodeId.String(), " from: ", from.Hex(), " withdraw: ", withdraw, " height: ", height)
	//todo
	owner :=  c.evm.CandidatePool.GetOwner(c.evm.StateDB, nodeId)
	if ok := bytes.Equal(owner.Bytes(), from.Bytes()); !ok {
		log.Error(ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if err := c.evm.CandidatePool.WithdrawCandidate(c.evm.StateDB, nodeId, withdraw, height); err!=nil {
		log.Error(err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateApplyWithdrawEvent, string(data))
	log.Info("json: ", string(data))
	return nil, nil
}

//Deposit withdrawal
func (c *candidateContract) CandidateWithdraw(nodeId discover.NodeID) ([]byte, error)  {
	//debug
	height := c.evm.Context.BlockNumber
	log.Info("CandidateWithdraw==> nodeId: ", nodeId.String(), " height: ", 	height)
	//todo
	if err :=c.evm.CandidatePool.RefundBalance(c.evm.StateDB, nodeId, height); err!=nil{
		log.Error(err.Error())
		return nil, err
	}
	//return
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(CandidateWithdrawEvent, string(data))
	log.Info("json: ", string(data))
	return nil, nil
}

//Get the refund history you have applied for
func (c *candidateContract) CandidateWithdrawInfos(nodeId discover.NodeID)([]byte, error){
	//debug
	log.Info("CandidateWithdrawInfos==> nodeId: ", nodeId.String())
	//todo
	infos, err := c.evm.CandidatePool.GetDefeat(c.evm.StateDB, nodeId)
	if err!=nil{
		log.Error(err.Error())
		return nil, err
	}
	//return
	type WithdrawInfo struct {
		Balance *big.Int 
		LockNumber *big.Int
		LockBlockCycle uint64
	} 
	type WithdrawInfos struct {
		Ret ResultCommon
		Infos []WithdrawInfo
	}
	r := WithdrawInfos{}
	r.Ret = ResultCommon{true, "success"}
	r.Infos = make([]WithdrawInfo, len(infos))
	for i, v := range infos {
		r.Infos[i] = WithdrawInfo{v.Deposit, v.BlockNumber, c.evm.CandidatePool.GetRefundInterval()}
	}
	data, _ := json.Marshal(r)
	sdata := DecodeResultStr(string(data))
	log.Info("json: ", string(data))
	return sdata, nil
}

//Set up additional information
func (c *candidateContract) SetCandidateExtra(nodeId discover.NodeID, extra string)([]byte, error){
	//debug
	from := c.contract.caller.Address()
	log.Info("SetCandidate==> nodeId: ", nodeId.String(), " extra: ", extra, " from: ", from.Hex())
	//todo
	owner :=  c.evm.CandidatePool.GetOwner(c.evm.StateDB, nodeId)
	if ok := bytes.Equal(owner.Bytes(), from.Bytes()); !ok {
		log.Error(ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if err := c.evm.CandidatePool.SetCandidateExtra(c.evm.StateDB, nodeId, extra); err!=nil{
		log.Error(err.Error())
		return nil, err
	}
	r := ResultCommon{true, "success"}
	data, _ := json.Marshal(r)
	c.addLog(SetCandidateExtraEvent, string(data))
	log.Info("json: ", string(data))
	return nil, nil
}

//Get candidate details
func (c *candidateContract) CandidateDetails(nodeId discover.NodeID) ([]byte, error)  {
	fmt.Println("into func CandidateDetails... ")
	fmt.Println(nodeId.String())
	candidate, err := c.evm.CandidatePool.GetCandidate(c.evm.StateDB, nodeId)
	if err != nil{
		fmt.Println("get CandidateDetails() occured error: ", err.Error())
		return nil, err
	}
	if nil == candidate {
		fmt.Println("The candidate for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(candidate)
	sdata := DecodeResultStr(string(data))
	fmt.Println("json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

//Get the current block candidate list 0~200
func (c *candidateContract) CandidateList() ([]byte, error) {
	fmt.Println("into func CandidateList... ")
	arr := c.evm.CandidatePool.GetChosens(c.evm.StateDB)
	if nil == arr {
		fmt.Println("The candidateList for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(arr)
	sdata := DecodeResultStr(string(data))
	fmt.Println("json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

//Get the current block round certifier list 25个
func (c *candidateContract) VerifiersList() ([]byte, error) {
	fmt.Println("into func VerifiersList... ")
	arr := c.evm.CandidatePool.GetChairpersons(c.evm.StateDB)
	if nil == arr {
		fmt.Println("The verifiersList for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(arr)
	sdata := DecodeResultStr(string(data))
	fmt.Println("json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

//transaction add event
func (c *candidateContract) addLog(event, data string) {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); err!=nil {
		log.Error("addlog rlp encode fail: ", err.Error())
	}
	c.evm.StateDB.AddLog(&types.Log{
		Address:common.CandidateAddr,
		Topics: []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data: buf.Bytes(),
		BlockNumber: c.evm.Context.BlockNumber.Uint64(),
	})
}

//return string format
func DecodeResultStr (result string) []byte {
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
