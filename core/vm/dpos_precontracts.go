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
	"reflect"
	"Platon-go/params"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"Platon-go/p2p/discover"
	"Platon-go/core/types"
)

//error def
var (
	ErrRepeatOwner = errors.New("Node ID cannot bind multiple owners")
	ErrPermissionDenied = errors.New("Transaction from address permission denied")
	ErrWithdrawEmpyt = errors.New("No withdrawal amount")
	ErrParamsRlpDecode = errors.New("Rlp decode faile")
	ErrParamsBaselen = errors.New("Params Base length does not match")
	ErrParamsLen = errors.New("Params length does not match")
	ErrUndefFunction = errors.New("Undefined function")
	ErrCandidateEmpyt = errors.New("CandidatePool is nil")
)

var PrecompiledContractsDpos = map[common.Address]PrecompiledContract{
	common.HexToAddress("0x1000000000000000000000000000000000000111") : &candidateContract{},
}

type candidatePool interface {
	SetCandidate(state StateDB, nodeId discover.NodeID, can *types.Candidate) error
	GetCandidate(state StateDB, nodeId discover.NodeID) (*types.Candidate, error)
	WithdrawCandidate (state StateDB, nodeId discover.NodeID, price int) error
	GetChosens (state StateDB, ) []*types.Candidate
	GetChairpersons (state StateDB, ) []*types.Candidate
	GetDefeat(state StateDB, nodeId discover.NodeID) ([]*types.Candidate, error)
	IsDefeat(state StateDB, nodeId discover.NodeID) (bool, error)
	RefundBalance (state StateDB, nodeId discover.NodeID, blockNumber uint64) error
	GetOwner (state StateDB, nodeId discover.NodeID) common.Address
}

type candidateContract struct{
	contract *Contract
	evm *EVM
}

func (c *candidateContract) RequiredGas(input []byte) uint64 {
	// TODO 获取设定的预编译合约消耗
	return params.EcrecoverGas
}

func (c *candidateContract) Run(input []byte) ([]byte, error) {
	// 用map封装所有的函数
	var command = map[string] interface{}{
		"CandidateDetails" : c.CandidateDetails,
		"CandidateApplyWithdraw" : c.CandidateApplyWithdraw,
		"CandidateDeposit" : c.CandidateDeposit,
		"CandidateList" : c.CandidateList,
		"CandidateWithdraw" : c.CandidateWithdraw,
		"VerifiersList" : c.VerifiersList,
		"SayHi" : SayHi,
	}
	//rlp decode
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &source); err != nil {
		fmt.Println(err)
		return nil, ErrParamsRlpDecode
	}
	//check
	if len(source)<2 {
		return nil, ErrParamsBaselen
	}
	if c.evm.CandidatePool==nil{
		return nil, ErrCandidateEmpyt
	}

	// 获取要调用的函数
	if _, ok := command[byteutil.BytesToString(source[1])]; !ok {
		return nil, ErrUndefFunction
	}
	funcValue := command[byteutil.BytesToString(source[1])]
	// 目标函数参数列表
	paramList := reflect.TypeOf(funcValue)
	// 目标函数参数个数
	paramNum := paramList.NumIn()
	// var param []interface{}
	params := make([]reflect.Value, paramNum)

	if paramNum!=len(source) {
		return nil, ErrParamsLen
	}

	for i := 0; i < paramNum; i++ {
		// 目标参数类型的值
		targetType := paramList.In(i).String()
		fmt.Println("i: ", i, " type: ", targetType)
		// 原始[]byte类型参数
		originByte := []reflect.Value{reflect.ValueOf(source[i+2])}
		// 转换为对应类型的参数
		params[i] = reflect.ValueOf(byteutil.Command[targetType]).Call(originByte)[0]
	}
	fmt.Println("params: ", params)
	// 传入参数调用函数
	//result := reflect.ValueOf(funcValue).Call(params)
	reflect.ValueOf(funcValue).Call(params)
	// TODO
	// 返回值也是一个 Value 的 slice，同样对应反射函数类型的返回值。
	//return result[0].Bytes(), result[1].Interface().(error)
	return nil, nil
}

func SayHi(nodeId discover.NodeID, owner common.Address, fee uint64) ([]byte, error) {
	fmt.Println("into ...")
	fmt.Println("CandidateDeposit==> nodeId: ", nodeId, " owner: ", owner, "  fee: ", fee)
	return nil, nil
}

//候选人申请 && 增加质押金
func (c *candidateContract) CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64) ([]byte, error)   {



	//params parse
	deposit := *c.contract.value
	txHash := c.evm.StateDB.TxHash()
	txIdx := c.evm.StateDB.TxIdx()
	fmt.Println("CandidateDeposit==> nodeId: ", nodeId.String(), " owner: ", owner.Hex(), " deposit: ", deposit,
		"  fee: ", fee, " txhash: ", txHash.Hex(), " txIdx: ", txIdx)

	//todo
	c.evm.CandidatePool.GetCandidate(c.evm.StateDB, nodeId)


	//判断nodeid和owner是否唯一
	//先获取已有质押金，加上本次value，更新





	//调用操作db的接口如果失败，则回滚交易。申请失败的交易，钱会被扣除,需要回滚
	//返回值用json形式按照实际合约执行的返回形式格式化


	return nil, nil
}

//申请退回质押金
func (c *candidateContract) CandidateApplyWithdraw(params [][]byte) ([]byte, error)  {

	if len(params)!=1 {
		return nil, ErrParamsLen
	}
	nodeId := hex.EncodeToString(params[0])
	from := c.contract.caller.Address().Hex()
	fmt.Println("CandidateApplyWithdraw==> nodeId: ", nodeId, " from: ", from)

	//校验from和owner是否一致
	//调用接口生成退款记录


	return nil, nil
}

//质押金提现
func (c *candidateContract) CandidateWithdraw(params [][]byte) ([]byte, error)  {

	if len(params)!=1 {
		return nil, ErrParamsLen
	}
	nodeId := hex.EncodeToString(params[0])

	//调用接口退款，判断返回值

	fmt.Println("CandidateWithdraw==> nodeId: ", nodeId)
	return nil, nil
}

//获取候选人详情
func (c *candidateContract) CandidateDetails(nodeId [64]byte) ([]byte, error)  {
	//cbft.GetDpos().GetCandidate(nodeId)
	// TODO
	return nil, nil
}

//获取当前区块候选人列表 0~200
func (c *candidateContract) CandidateList() ([]byte, error) {
	// dpos.GetChosens()
	return nil, nil
}

//获取当前区块轮次验证人列表 25个
func (c *candidateContract) VerifiersList() ([]byte, error) {
	// dpos.GetChairpersons()
	return nil, nil
}

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

	encodedStr := hex.EncodeToString(finalData)
	fmt.Println("finalData: ", encodedStr)

	return finalData
}
