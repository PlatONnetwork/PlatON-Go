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
	"Platon-go/consensus/cbft"
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
	state StateDB
}

func (c *candidateContract) RequiredGas(input []byte) uint64 {
	// TODO 获取设定的预编译合约消耗
	return params.EcrecoverGas
}

func (c *candidateContract) Run(input []byte) ([]byte, error) {
	//rlp decode
	var params [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &params); err!=nil {
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
		case "CandidateDetails":
			return c.CandidateDetails(params[2:])
		case "CandidateList":
			return c.CandidateList(params[2:])
		case "VerifiersList":
			return c.VerifiersList(params[2:])
		default:
			fmt.Println("Undefined function")
			return nil, ErrUndefFunction
	}
}

var dpos *cbft.Dpos

// 初始化获取dpos实例
func init() {
	dpos = cbft.GetDpos()
}

//获取候选人详情
func (c *candidateContract) CandidateDetails(params [][]byte)([]byte, error)  {
	// TODO nodeId discover.NodeID 参数校验
	// dpos.GetCandidate()
	return nil, nil
}

//获取当前区块候选人列表 0~200
func (c *candidateContract) CandidateList(params [][]byte) ([]byte, error) {
	// dpos.GetChosens()
	return nil, nil
}

//获取当前区块轮次验证人列表 25个
func (c *candidateContract) VerifiersList(params [][]byte) ([]byte, error) {
	// dpos.GetChairpersons()
	return nil, nil
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
	return nil, nil
}

//质押金提现
func (c *candidateContract) CandidateWithdraw(params [][]byte) ([]byte, error)  {

	if len(params)!=1 {
		return nil, ErrParamsLen
	}
	nodeId := hex.EncodeToString(params[0])

	fmt.Println("CandidateWithdraw==> nodeId: ", nodeId)
	return nil, nil
}