package vm

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestRlpEncode(t *testing.T) {
	nodeId, _ := hex.DecodeString("01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012340")
	//owner, _ := hex.DecodeString("740ce31b3fac20dac379db243021a51e80ad00d7") //38
	//owner, _ := hex.DecodeString("5a5c4368e2692746b286cee36ab0710af3efa6cf") //39
	//owner, _ := hex.DecodeString("493301712671ada506ba6ca7891f436d29185821") //40
	//fmt.Println(nodeId)
	// 编码
	var source [][]byte
	source = make([][]byte, 0)
	source = append(source, common.Hex2Bytes("1011"))  // tx type
	//source = append(source, []byte("CandidateDeposit")) // func name
	source = append(source, []byte("CandidateApplyWithdraw")) // func name
	//source = append(source, []byte("CandidateWithdraw")) // func name
	//source = append(source, []byte("CandidateWithdrawInfos")) // func name
	//source = append(source, []byte("SetCandidateExtra")) // func name
	//source = append(source, []byte("CandidateDetails")) // func name
	//source = append(source, []byte("CandidateList")) // func name
	//source = append(source, []byte("VerifiersList")) // func name
	source = append(source, nodeId) // [64]byte nodeId discover.NodeID
	//source = append(source, owner) // [20]byte owner common.Address
	//source = append(source, byteutil.Uint64ToBytes(100)) // fee
	//source = append(source, []byte("0.0.0.1")) // host
	//source = append(source, []byte("6789")) // port
	//source = append(source, []byte("extra info..")) // extra
	source = append(source, new(big.Int).SetInt64(1).Bytes()) // withdraw

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, source)
	if err != nil {
		fmt.Println(err)
		t.Errorf("fail")
	}
	encodedBytes := buffer.Bytes()
	// 编码后字节数组
	fmt.Println(encodedBytes)
	// to hex as data
	fmt.Println(hexutil.Encode(encodedBytes))

	// 解码
	ptr := new(interface{})
	rlp.Decode(bytes.NewReader(encodedBytes), &ptr)

	deref := reflect.ValueOf(ptr).Elem().Interface()
	fmt.Println(deref)
	for i, v := range deref.([]interface{}) {
		// fmt.Println(i,"    ",hex.EncodeToString(v.([]byte)))
		// 类型判断，然后转换
		switch i {
		case 0:
			// fmt.Println(string(v.([]byte)))
		case 1:
			fmt.Println(string(v.([]byte)))
		case 2:
			// fmt.Println(string(v.([]byte)))
		}
	}
}

func TestAppendSlice(t *testing.T)  {
	a := []int{0, 1, 2, 3, 4}
	// 删除第i个元素
	i := 2
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)
}
