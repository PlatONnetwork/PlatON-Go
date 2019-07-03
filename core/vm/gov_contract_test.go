package vm

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"testing"
)


func rlpEncode(value interface{}) []byte{
	bytes, err := rlp.EncodeToBytes(value)
	if err == nil {
		return bytes
	}else{
		return nil
	}
}

var nodeID = discover.MustHexID("bf5317e3a60a55e9c7d09cb20d4381f579c4318eb1031426612959ab5fa7a9d3f3e362b58887e83df8048115501f0b0390b4cdab4548b2728b6633ab692f9ca1")


func buildSubmitTextInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, rlpEncode(2000))			// func type code
	input = append(input, rlpEncode("submitText"))	// func name
	input = append(input, rlpEncode(nodeID))				// param 1 ...
	input = append(input, rlpEncode("githubID"))
	input = append(input, rlpEncode("textTopic"))
	input = append(input, rlpEncode("textDesc"))
	input = append(input, rlpEncode("textUrl"))
	input = append(input, rlpEncode(uint64(1000)))
	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, input)
	if err != nil {
		panic(err)
	} else {
		return hexutil.Encode(buf.Bytes())
	}
}

var govContractTests = []precompiledTest{
	{
		input:   	buildSubmitTextInput(),
		expected: 	"",
		name:     	"submitText1",
	},
}

// Tests the sample inputs from the elliptic curve pairing check EIP 197.
func TestPrecompiledGovContract(t *testing.T) {
	for _, test := range govContractTests {
		testPrecompiled("08", test, t)
	}
}
