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
<<<<<<< HEAD
	"Platon-go/common/byteutil"
=======
>>>>>>> 6147a1bfddb57d02d70d7edb5db56bd27fc4bff6
	//"Platon-go/consensus/cbft"
	"Platon-go/params"
	"Platon-go/rlp"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
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
)

var PrecompiledContractsDpos = map[common.Address]PrecompiledContract{
	common.HexToAddress("0x1000000000000000000000000000000000000111") : &candidateContract{},
}

type candidateContract struct{
	contract *Contract
	evm *EVM
}

<<<<<<< HEAD
=======
// 用map封装所有的函数
/*var command = map[string] interface{} {
	"CandidateDetails" : candidateContract.CandidateDetails,
	"CandidateApplyWithdraw" : candidateContract.CandidateApplyWithdraw,
	"CandidateDeposit" : candidateContract.CandidateDeposit,
	"CandidateList" : candidateContract.CandidateList,
	"CandidateWithdraw" : candidateContract.CandidateWithdraw,
	"VerifiersList" : candidateContract.VerifiersList,
	// TODO test delete
	"SayHi" : SayHi,
}*/

>>>>>>> 6147a1bfddb57d02d70d7edb5db56bd27fc4bff6
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
		// TODO test delete
		"SayHi" : SayHi,
	}

	//rlp decode
	var params [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &params); err != nil {
		fmt.Println(err)
		return nil, ErrParamsRlpDecode
	}
	//function call
	if len(params)<2 {
		return nil, ErrParamsBaselen
	}

	switch string(params[1]) {
	case "CandidateDeposit":
		return c.CandidateDeposit(params[2:])
	case "CandidateApplyWithdraw":
		return c.CandidateApplyWithdraw(params[2:])
	case "CandidateWithdraw":
		return c.CandidateWithdraw(params[2:])
	default:
		fmt.Println("Undefined function")
		return nil, ErrUndefFunction
	}

		// 获取要调用的函数
	/*if _, ok := command[byteutil.BytesToString(source[1])]; !ok {
		return nil, ErrUndefFunction
	}
	funcValue := command[byteutil.BytesToString(source[1])]
	// 目标函数参数列表
	paramList := reflect.TypeOf(funcValue)
	// 目标函数参数个数
	paramNum := paramList.NumIn()
	// var param []interface{}
	params := make([]reflect.Value, paramNum)

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
	// TODO
	// 返回值也是一个 Value 的 slice，同样对应反射函数类型的返回值。
<<<<<<< HEAD
	return result[0].Bytes(), result[1].Interface().(error)
}

func SayHi(a []byte, b [64]byte) (string) {
	fmt.Println(b)
	return "2"
}

//获取候选人详情
func (c *candidateContract) CandidateDetails(nodeId [64]byte) ([]byte, error)  {
	//cbft.GetDpos().GetCandidate(nodeId)
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
=======
	return result[0].Bytes(), result[1].Interface().(error)*/
>>>>>>> 6147a1bfddb57d02d70d7edb5db56bd27fc4bff6
}

//候选人申请 && 增加质押金
func (c *candidateContract) CandidateDeposit(params [][]byte) ([]byte, error)   {

	//params parse
	if len(params)!=3 {
		return nil, ErrParamsLen
	}
	nodeId := hex.EncodeToString(params[0])
	owner := hex.EncodeToString(params[1])
	fee := binary.BigEndian.Uint64(params[2])
	deposit := *c.contract.value
	fmt.Println("CandidateDeposit==> nodeId: ", nodeId, " owner: ", owner, " deposit: ", deposit, "  fee: ", fee)

	//todo
	//dpos := cbft.GetDpos()
	//dpos.Switch()
	//cbft.GetDpos()




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

func SayHi(a []byte, b [64]byte) (string) {
	fmt.Println(b)
	return "2"
}