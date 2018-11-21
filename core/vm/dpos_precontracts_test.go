package vm

import (
	"Platon-go/common"
	"Platon-go/common/byteutil"
	"Platon-go/common/hexutil"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func TestRlpEncode(t *testing.T) {
	nodeId, _ := hex.DecodeString("e152be5f5f0167250592a12a197ab19b215c5295d5eb0bb1133673dc8607530db1bfa5415b2ec5e94113f2fce0c4a60e697d5d703a29609b197b836b020446c7")
	owner, _ := hex.DecodeString("4FED1fC4144c223aE3C1553be203cDFcbD38C581")

	// 编码
	var source [][]byte
	source = make([][]byte, 0)
	source = append(source, common.Hex2Bytes("1011"))  // tx type
	source = append(source, []byte("SayHi")) // func name
	source = append(source, nodeId) // [64]byte nodeId discover.NodeID
	source = append(source, owner) // [20]byte owner common.Address
	source = append(source, byteutil.Uint64ToBytes(100))
	//source = append(source, []byte("abc"))

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
