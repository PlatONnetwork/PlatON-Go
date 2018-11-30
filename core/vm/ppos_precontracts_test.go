package vm

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/rlp"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"Platon-go/common/byteutil"
	//"math/big"
	"reflect"
	"testing"
)

func TestRlpEncode(t *testing.T) {
	nodeId, _ := hex.DecodeString("751f4f62fccee84fc290d0c68d673e4b0cc6975a5747d2baccb20f954d59ba3315d7bfb6d831523624d003c8c2d33451129e67c3eef3098f711ef3b3e268fd3c")

	owner, _ := hex.DecodeString("740ce31b3fac20dac379db243021a51e80ad00d7") //38
	//owner, _ := hex.DecodeString("5a5c4368e2692746b286cee36ab0710af3efa6cf") //39
	//owner, _ := hex.DecodeString("493301712671ada506ba6ca7891f436d29185821") //40
	//fmt.Println(nodeId)
	// code
	var source [][]byte
	source = make([][]byte, 0)
	source = append(source, common.Hex2Bytes("1011"))  // tx type
	source = append(source, []byte("CandidateDeposit")) // func name
	//source = append(source, []byte("CandidateApplyWithdraw")) // func name
	//source = append(source, []byte("CandidateWithdraw")) // func name
	//source = append(source, []byte("CandidateWithdrawInfos")) // func name
	//source = append(source, []byte("SetCandidateExtra")) // func name
	//source = append(source, []byte("CandidateDetails")) // func name
	//source = append(source, []byte("CandidateList")) // func name
	//source = append(source, []byte("VerifiersList")) // func name
	source = append(source, nodeId) // [64]byte nodeId discover.NodeID
	source = append(source, owner) // [20]byte owner common.Address
	source = append(source, byteutil.Uint64ToBytes(100)) // fee
	source = append(source, []byte("192.168.9.184")) // host
	source = append(source, []byte("16789")) // port
	source = append(source, []byte("extra info..")) // extra
	//source = append(source, new(big.Int).SetInt64(1).Bytes()) // withdraw

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, source)
	if err != nil {
		fmt.Println(err)
		t.Errorf("fail")
	}
	encodedBytes := buffer.Bytes()
	// result
	fmt.Println(encodedBytes)
	// to hex as data
	fmt.Println(hexutil.Encode(encodedBytes))

	// decode
	ptr := new(interface{})
	rlp.Decode(bytes.NewReader(encodedBytes), &ptr)

	deref := reflect.ValueOf(ptr).Elem().Interface()
	fmt.Println(deref)
	for i, v := range deref.([]interface{}) {
		// fmt.Println(i,"    ",hex.EncodeToString(v.([]byte)))
		// check type and switch
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
	i := 2
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)
}

func TestRlpData(t *testing.T)  {

	nodeId, _ := hex.DecodeString("a363d1243646b6eabf1d4851f646b523f5707d053caab95022f1682605aca0537ee0c5c14b4dfa76dcbce264b7e68d59de79a42b7cda059e9d358336a9ab8d80")
	owner, _ := hex.DecodeString("f216d6e4c17097a60ee2b8e5c88941cd9f07263b")

	//CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string)
	var CandidateDeposit [][]byte
	CandidateDeposit = make([][]byte, 0)
	CandidateDeposit = append(CandidateDeposit, uint64ToBytes(0xf1))
	CandidateDeposit = append(CandidateDeposit, []byte("CandidateDeposit"))
	CandidateDeposit = append(CandidateDeposit, nodeId)
	CandidateDeposit = append(CandidateDeposit, owner)
	CandidateDeposit = append(CandidateDeposit, uint64ToBytes(500))	//10000
	CandidateDeposit = append(CandidateDeposit, []byte("0.0.0.0"))
	CandidateDeposit = append(CandidateDeposit, []byte("30303"))
	CandidateDeposit = append(CandidateDeposit, []byte("extra data"))
	bufDeposit := new(bytes.Buffer)
	err := rlp.Encode(bufDeposit, CandidateDeposit)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateDeposit encode rlp data fail")
	} else {
		fmt.Println("CandidateDeposit data rlp: ", hexutil.Encode(bufDeposit.Bytes()))
	}

	//CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int)
	var CandidateApplyWithdraw [][]byte
	CandidateApplyWithdraw = make([][]byte, 0)
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, uint64ToBytes(0xf1))
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, []byte("CandidateApplyWithdraw"))
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, nodeId)
	withdraw, ok :=new(big.Int).SetString("14d1120d7b160000", 16)
	if !ok {
		t.Errorf("big int setstring fail")
	}
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, withdraw.Bytes())
	bufApply := new(bytes.Buffer)
	err = rlp.Encode(bufApply, CandidateApplyWithdraw)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateApplyWithdraw encode rlp data fail")
	} else {
		fmt.Println("CandidateApplyWithdraw data rlp: ", hexutil.Encode(bufApply.Bytes()))
	}

	//CandidateWithdraw(nodeId discover.NodeID)
	var CandidateWithdraw [][]byte
	CandidateWithdraw = make([][]byte, 0)
	CandidateWithdraw = append(CandidateWithdraw, uint64ToBytes(0xf1))
	CandidateWithdraw = append(CandidateWithdraw, []byte("CandidateWithdraw"))
	CandidateWithdraw = append(CandidateWithdraw, nodeId)
	bufWith := new(bytes.Buffer)
	err = rlp.Encode(bufWith, CandidateWithdraw)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateWithdraw encode rlp data fail")
	} else {
		fmt.Println("CandidateWithdraw data rlp: ", hexutil.Encode(bufWith.Bytes()))
	}

	//CandidateWithdrawInfos(nodeId discover.NodeID)
	var CandidateWithdrawInfos [][]byte
	CandidateWithdrawInfos = make([][]byte, 0)
	CandidateWithdrawInfos = append(CandidateWithdrawInfos, uint64ToBytes(0xf1))
	CandidateWithdrawInfos = append(CandidateWithdrawInfos, []byte("CandidateWithdrawInfos"))
	CandidateWithdrawInfos = append(CandidateWithdrawInfos, nodeId)
	bufWithInfos := new(bytes.Buffer)
	err = rlp.Encode(bufWithInfos, CandidateWithdrawInfos)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateWithdrawInfos encode rlp data fail")
	} else {
		fmt.Println("CandidateWithdrawInfos data rlp: ", hexutil.Encode(bufWithInfos.Bytes()))
	}

	//CandidateDetails(nodeId discover.NodeID)
	var CandidateDetails [][]byte
	CandidateDetails = make([][]byte, 0)
	CandidateDetails = append(CandidateDetails, uint64ToBytes(0xf1))
	CandidateDetails = append(CandidateDetails, []byte("CandidateDetails"))
	CandidateDetails = append(CandidateDetails, nodeId)
	bufDetails := new(bytes.Buffer)
	err = rlp.Encode(bufDetails, CandidateDetails)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateDetails encode rlp data fail")
	} else {
		fmt.Println("CandidateDetails data rlp: ", hexutil.Encode(bufDetails.Bytes()))
	}
}

func TestRlpDecode(t *testing.T)  {

	//HexString -> []byte
	rlpcode, _ := hex.DecodeString("f85c8800000000000000049043616e64696461746544657461696c73b840e152be5f5f0167250592a12a197ab19b215c5295d5eb0bb1133673dc8607530db1bfa5415b2ec5e94113f2fce0c4a60e697d5d703a29609b197b836b020446c7")
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(rlpcode), &source); err != nil {
		fmt.Println(err)
		t.Errorf("TestRlpDecode decode rlp data fail")
	}

	for i,v := range source {
		fmt.Println("i: ", i, " v: ", hex.EncodeToString(v))
	}
}

func uint64ToBytes(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return buf[:]
}
