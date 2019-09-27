package runtime

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/life/utils"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
)

func TestDefaults(t *testing.T) {
	cfg := new(Config)
	setDefaults(cfg)

	if cfg.Difficulty == nil {
		t.Error("expected difficulty to be non nil")
	}
	if cfg.Time == nil {
		t.Error("expected time to non nil")
	}
	if cfg.GasLimit == 0 {
		t.Error("didn't expect gaslimit to be zero")
	}
	if cfg.GasPrice == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.Value == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.GetHashFn == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.BlockNumber == nil {
		t.Error("expected block number to be non nil")
	}
}

func TestEVM(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	code := genCodeInput()

	Execute(code, nil, nil)
}

func TestExecute(t *testing.T) {
	code := genCodeInput()
	input := genInput()
	ret, _, err := Execute(code, input, nil)
	if err != nil {
		t.Fatal("didn't expect error", err)
	}
	// [255,255,255,255,......,158]

	/*buf := bytes.NewBuffer(ret)
	binary.Read(buf, binary.BigEndian, &result)*/
	hexRes := common.Bytes2Hex(ret)
	fmt.Println(hexRes)
	//result.UnmarshalText(ret)
	result, _ := math.ParseBig256(hexRes)
	/*if big.Int(result).Cmp(big.NewInt(-100)) != 0 {
		t.Error("Expected 10, got", big.Int(result))
	}*/
	fmt.Println("....", result)
}

func TestCall(t *testing.T) {
	state, _ := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	address := common.HexToAddress("0x0a")
	code := genCodeInput()
	state.SetCode(address, code)
	input := genInput()
	ret, _, err := Call(address, input, &Config{State: state})

	callInput := genCallInput()
	ret02, _, err := Call(address, callInput, &Config{State: state})
	if err != nil {
		//t.Fatal("didn't expect error", err)
	}
	fmt.Println("CallResponse:", string(ret02))
	num := string(ret)
	expected := "x"
	if !strings.EqualFold(num, expected) {
		//t.Error("Expected "+expected+", got", num)
	}
}

func TestCallCode(t *testing.T){
	code := genInput()
	hexcode := common.Bytes2Hex(code)
	fmt.Println("encoded(Input):", hexcode)
}

func TestGGG(t *testing.T) {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, utils.Int64ToBytes(2))
	input = append(input, []byte("set"))
	input = append(input, utils.Int64ToBytes(100))

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	if err != nil {
		fmt.Println("geninput fail.", err)
	}
	fmt.Println(common.Bytes2Hex(buffer.Bytes()))
}

func TestGetIRInput(t *testing.T) {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, utils.Int64ToBytes(9))
	input = append(input, []byte("get_participants"))
	//input = append(input, []byte("func01"))
	//input = append(input, []byte("extradata"))

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	if err != nil {
		fmt.Println("geninput fail.", err)
	}
	fmt.Println(common.Bytes2Hex(buffer.Bytes()))
}

func TestGetStartInput(t *testing.T) {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, utils.Int64ToBytes(5))
	input = append(input, []byte("TestFooAdd01"))
	//input = append(input, []byte("func01"))
	input = append(input, []byte("extradata"))

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	if err != nil {
		fmt.Println("geninput fail.", err)
	}
	fmt.Println(common.Bytes2Hex(buffer.Bytes()))
}

func genInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, utils.Int64ToBytes(1))
	input = append(input, []byte("get"))
	//input = append(input, []byte("0x0000000000000000000000000000000000000001"))
	//input = append(input, []byte("0x0000000000000000000000000000000000000002"))
	//input = append(input, utils.Int64ToBytes(-100))

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	if err != nil {
		fmt.Println("geninput fail.", err)
	}
	return buffer.Bytes()
}

func genCallInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, utils.Int64ToBytes(1))
	input = append(input, []byte("getBalance"))
	input = append(input, []byte("0x0000000000000000000000000000000000000002"))

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	if err != nil {
		fmt.Println("genCallInput fail.", err)
	}
	return buffer.Bytes()
}

func genCodeInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, utils.Int64ToBytes(2))
	code, _ := ioutil.ReadFile("../contract/inputtest.wasm")
	//fmt.Println("SrcCode:", common.Bytes2Hex(code))
	input = append(input, code)
	abi, _ := ioutil.ReadFile("../contract/inputtest.cpp.abi.json")
	//fmt.Println("SrcAbi:", common.Bytes2Hex(abi))
	input = append(input, abi)
	buffer := new(bytes.Buffer)
	rlp.Encode(buffer, input)
	return buffer.Bytes()
}

func TestCreateCode(t *testing.T){
	code := genCodeInput()
	hexcode := common.Bytes2Hex(code)
	fmt.Println("encoded :", hexcode)

	// decode
	parseRlpData(common.Hex2Bytes(hexcode))
}

func parseRlpData(rlpData []byte) (int64, []byte, []byte, error) {
	ptr := new(interface{})
	err := rlp.Decode(bytes.NewReader(rlpData), &ptr)
	if err != nil {
		return -1, nil, nil, err
	}
	rlpList := reflect.ValueOf(ptr).Elem().Interface()

	if _, ok := rlpList.([]interface{}); !ok {
		return -1, nil, nil, fmt.Errorf("invalid rlp format.")
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) <= 2 {
		return -1, nil, nil, fmt.Errorf("invalid input. ele must greater than 2")
	}
	var (
		txType int64
		code   []byte
		abi    []byte
	)
	if v, ok := iRlpList[0].([]byte); ok {
		txType = utils.BytesToInt64(v)
	}
	if v, ok := iRlpList[1].([]byte); ok {
		code = v
		//fmt.Println("dstCode: ", common.Bytes2Hex(code))
	}
	if v, ok := iRlpList[2].([]byte); ok {
		abi = v
		//fmt.Println("dstAbi:", common.Bytes2Hex(abi))
	}
	return txType, abi, code, nil
}

func TestUtils(t *testing.T) {
	fmt.Println("Code:", common.Bytes2Hex(genCodeInput()))
	fmt.Println("TxCall:", common.Bytes2Hex(genInput()))
	fmt.Println("Call:", common.Bytes2Hex(genCallInput()))
}

func TestParseCode(t *testing.T) {
	hexCode := ""
	code := common.Hex2Bytes(hexCode)
	typ, abi, _, _ := parseRlpData(code)
	fmt.Println(typ)
	fmt.Println(string(abi))
}

func TestParseInput(t *testing.T) {
	data := "f86e880000000000000002897472616e7366657231aa307861613331636139643839323830306161363733383362623838313134623631383638323231656532aa3078616133316361396438393238303061613637333833626238383131346236313836383232316565338400000014"
	input := common.Hex2Bytes(data)
	parseInputFromAbi(input,nil)
}

// parse input(payload)
func parseInputFromAbi(input []byte, abi []byte) (txType int, funcName string, params []int64, returnType string, err error) {
	if input == nil || len(input) <= 1 {
		return -1, "", nil, "", fmt.Errorf("invalid input.")
	}
	// [txType][funcName][args1][args2]
	// rlp decode
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(input), &ptr)
	if err != nil {
		return -1, "", nil, "", err
	}
	rlpList := reflect.ValueOf(ptr).Elem().Interface()

	if _, ok := rlpList.([]interface{}); !ok {
		return -1, "", nil, "", fmt.Errorf("invalid rlp format.")
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) < 2 {
		return -1, "", nil, "", fmt.Errorf("invalid input. ele must greater than 2")
	}

	params = make([]int64, 0)
	if v, ok := iRlpList[0].([]byte); ok {
		txType = int(utils.BytesToInt64(v))
	}
	if v, ok := iRlpList[1].([]byte); ok {
		funcName = string(v)
	}
	return txType, funcName, params, returnType,nil
}
