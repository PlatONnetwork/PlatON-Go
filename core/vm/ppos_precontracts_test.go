package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	//"math/big"
	"reflect"
	"testing"
)

func TestRlpEncode(t *testing.T) {
	nodeId, _ := hex.DecodeString("751f4f62fccee84fc290d0c68d673e4b0cc6975a5747d2baccb20f954d59ba3315d7bfb6d831523624d003c8c2d33451129e67c3eef3098f711ef3b3e268fd3c")

	owner, _ := hex.DecodeString("740ce31b3fac20dac379db243021a51e80ad00d7") //38
	//owner, _ := hex.DecodeString("5a5c4368e2692746b286cee36ab0710af3efa6cf") //39
	//owner, _ := hex.DecodeString("493301712671ada506ba6ca7891f436d29185821") //40

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

	nodeId := []byte("0x751f4f62fccee84fc290d0c68d673e4b0cc6975a5747d2baccb20f954d59ba3315d7bfb6d831523624d003c8c2d33451129e67c3eef3098f711ef3b3e268fd3c")
	owner := []byte("0xf216d6e4c17097a60ee2b8e5c88941cd9f07263b")

	//CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string)
	var CandidateDeposit [][]byte
	CandidateDeposit = make([][]byte, 0)
	CandidateDeposit = append(CandidateDeposit, uint64ToBytes(0xf1))
	CandidateDeposit = append(CandidateDeposit, []byte("CandidateDeposit"))
	CandidateDeposit = append(CandidateDeposit, nodeId)
	CandidateDeposit = append(CandidateDeposit, owner)

	CandidateDeposit = append(CandidateDeposit, uint64ToBytes(500))	//10000
	CandidateDeposit = append(CandidateDeposit, []byte("192.168.9.182"))
	CandidateDeposit = append(CandidateDeposit, []byte("16789"))
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
	rlpcode, _ := hex.DecodeString("f9024005853339059800829c409410000000000000000000000000000000000000019331303030303030303030303030303030303030b901c7f901c48800000000000000029043616e6469646174654465706f736974b88230786133363364313234333634366236656162663164343835316636343662353233663537303764303533636161623935303232663136383236303561636130353337656530633563313462346466613736646362636532363462376536386435396465373961343262376364613035396539643335383333366139616238643830aa3078663231366436653463313730393761363065653262386535633838393431636439663037323633628800000000000001f487302e302e302e30853330333033b8e27b226e6f64654e616d65223a22e88a82e782b9e5908de7a7b0222c226e6f64654469736372697074696f6e223a22e88a82e782b9e7ae80e4bb8b222c226e6f64654465706172746d656e74223a22e69cbae69e84e5908de7a7b0222c226f6666696369616c57656273697465223a227777772e706c61746f6e2e6e6574776f726b222c226e6f6465506f727472616974223a2255524c222c2274696d65223a313534333931333639353638352c226f776e6572223a22307866323136643665346331373039376136306565326238653563383839343163643966303732363362227d1ba0e789e2d95ed796dec19e7a40b760a9849a1ca09110e1f95e46ed0c18487cfdf3a021bfd18bdb4c32a70f2836c6b78944c88b92d09f5e5201f125444792538617e7")
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
