package vm

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"testing"
)

func TestRLP_encode (t *testing.T) {

	var params [][]byte
	params = make([][]byte, 0)

	fnType, err := rlp.EncodeToBytes(1102)
	if nil != err {
		fmt.Println("fnType err", err)
	}
	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateDeposit encode rlp data fail")
	} else {
		fmt.Println("CandidateDeposit data rlp: ", hexutil.Encode(buf.Bytes()))
	}
}


func TestRLP_2 (t *testing.T) {

	/*var GetVerifiersList [][]byte
	GetVerifiersList = make([][]byte, 0)
	GetVerifiersList = append(GetVerifiersList, byteutil.Uint64ToBytes(0xf1))
	GetVerifiersList = append(GetVerifiersList, []byte("GetVerifiersList"))
	bufGetVerifiersList := new(bytes.Buffer)
	err := rlp.Encode(bufGetVerifiersList, GetVerifiersList)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetVerifiersList encode rlp data fail")
	} else {
		fmt.Println("GetVerifiersList data rlp: ", hexutil.Encode(bufGetVerifiersList.Bytes()))
	}*/
}