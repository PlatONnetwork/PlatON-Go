// create by platon
package vm

import (
	"Platon-go/cmd/ctool/rlp"
	"Platon-go/common/byteutil"
	"bytes"
	"fmt"
	"math/big"
	"reflect"
)

// 用map封装所有的函数
var command = map[string] interface{} {
	"CandidateApply" : CandidateApply,
	"CandidateQuit" : CandidateQuit,
	"CandidateDetails" : CandidateDetails,
	"CandidateList" : CandidateList,
	"VerifiersList" : VerifiersList,
	"SayHi" : SayHi,
}

// def implemented as a native contract.
type def struct {}

/**
 * 合约gas消耗
 */
func (c *def) RequiredGas(input []byte) uint64 {
	// TODO 获取设定的预编译合约消耗
	return 3000
}

/**
 * 合约的执行
 * rlp解码、解析参数匹配对应函数并执行
 */
func (c *def) Run(input []byte) ([]byte, error) {
	// rlp解码
	source := make([][]byte, 0)
	rlp.Decode(bytes.NewReader(input), &source)
	// 获取要调用的函数
	funcValue := command[byteutil.BytesToString(source[1])]
	// 目标函数参数列表
	paramList := reflect.TypeOf(funcValue)
	// 目标函数参数个数
	paramNum := paramList.NumIn()
	// var param []interface{}
	params := make([]reflect.Value, paramNum)

	for i := 0; i < paramNum; i++ {
		// 目标参数类型的值
		targetType := paramList.In(i).Name()
		// 原始[]byte类型参数
		originByte := []reflect.Value{reflect.ValueOf(source[i+2])}
		// 转换为对应类型的参数
		if targetType != "" {
			params[i] = reflect.ValueOf(byteutil.Command[targetType]).Call(originByte)[0]
		} else {
			params[i] = originByte[0]
		}
	}
	// 传入参数调用函数
	result := reflect.ValueOf(funcValue).Call(params)

	// TODO 处理返回值 []byte
	fmt.Println(result)

	return nil, nil
}

func SayHi(a []byte, b uint64) (string) {
	fmt.Println(b)
	return "2"
}

// 候选人结构
type candidateInfo struct {
	nodeId []byte	// 节点Id
	deposit *big.Int	// 质押金
	height *big.Int	// 最后一次质押所在块高
	position uint8	// 当前块高所出排行名次
}

// 获取某个候选人详情
// 参数：nodeId 候选人ID(节点公钥)
// 返回值：candidateInfo 候选人详情节后，error{int错误码，string错误信息}
func CandidateDetails(nodeId []byte)(candidateInfo, error){
	// 根据nodeId,调用底层GetCandidate(),返回候选人信息
	fmt.Println("into ....")
	return candidateInfo{},nil
}

// 获取当前轮次候选人列表 0~200
// 参数：无
// 返回值：map->节点ID：详情
func CandidateList() map[[64]byte]candidateInfo  {
	return nil
}

// 获取当前轮次验证人列表 25个  TODO: 当前轮次确定
// 参数：无
// 返回值：map->节点ID：详情
func VerifiersList() map[[64]byte]candidateInfo {
	return nil
}

// 候选人申请 && 增加质押金
// 参数：nodeId 候选人ID (节点公钥)，owner 质押金的退回地址，deposit 质押金额
// 返回值：error{int错误码，string错误信息}
// 增加质押金操作：Value = 已锁定质押金+ 新增质押金
func CandidateApply(nodeId [64]byte, owner [20]byte, deposit big.Int) error   {
	return nil
}

// 候选人退出
// 参数：nodeId 候选人ID(节点公钥)，sig 节点签名值
// 返回值：error{int错误码，string错误信息}
func CandidateQuit(nodeId [64]byte, sig []byte) error  {
	return nil
}