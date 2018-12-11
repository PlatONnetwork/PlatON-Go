package runtime

import (
	"Platon-go/common/math"
	"Platon-go/life/utils"
	"Platon-go/rlp"
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"Platon-go/common"
	"Platon-go/core/state"
	"Platon-go/ethdb"
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
		t.Fatal("didn't expect error", err)
	}
	fmt.Println("CallResponse:", string(ret02))
	num := string(ret)
	expected := "我是你大爷"
	if !strings.EqualFold(num, expected) {
		t.Error("Expected "+expected+", got", num)
	}
}

func BenchmarkCall(b *testing.B) {
	//var definition = `[{"constant":true,"inputs":[],"name":"seller","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"abort","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"value","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[],"name":"refund","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"buyer","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"confirmReceived","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"state","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":false,"inputs":[],"name":"confirmPurchase","outputs":[],"type":"function"},{"inputs":[],"type":"constructor"},{"anonymous":false,"inputs":[],"name":"Aborted","type":"event"},{"anonymous":false,"inputs":[],"name":"PurchaseConfirmed","type":"event"},{"anonymous":false,"inputs":[],"name":"ItemReceived","type":"event"},{"anonymous":false,"inputs":[],"name":"Refunded","type":"event"}]`
	//
	//var code = common.Hex2Bytes("6060604052361561006c5760e060020a600035046308551a53811461007457806335a063b4146100865780633fa4f245146100a6578063590e1ae3146100af5780637150d8ae146100cf57806373fac6f0146100e1578063c19d93fb146100fe578063d696069714610112575b610131610002565b610133600154600160a060020a031681565b610131600154600160a060020a0390811633919091161461015057610002565b61014660005481565b610131600154600160a060020a039081163391909116146102d557610002565b610133600254600160a060020a031681565b610131600254600160a060020a0333811691161461023757610002565b61014660025460ff60a060020a9091041681565b61013160025460009060ff60a060020a9091041681146101cc57610002565b005b600160a060020a03166060908152602090f35b6060908152602090f35b60025460009060a060020a900460ff16811461016b57610002565b600154600160a060020a03908116908290301631606082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517f72c874aeff0b183a56e2b79c71b46e1aed4dee5e09862134b8821ba2fddbf8bf9250a150565b80546002023414806101dd57610002565b6002805460a060020a60ff021973ffffffffffffffffffffffffffffffffffffffff1990911633171660a060020a1790557fd5d55c8a68912e9a110618df8d5e2e83b8d83211c57a8ddd1203df92885dc881826060a15050565b60025460019060a060020a900460ff16811461025257610002565b60025460008054600160a060020a0390921691606082818181858883f150508354604051600160a060020a0391821694503090911631915082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517fe89152acd703c9d8c7d28829d443260b411454d45394e7995815140c8cbcbcf79250a150565b60025460019060a060020a900460ff1681146102f057610002565b6002805460008054600160a060020a0390921692909102606082818181858883f150508354604051600160a060020a0391821694503090911631915082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517f8616bbbbad963e4e65b1366f1d75dfb63f9e9704bbbf91fb01bec70849906cf79250a15056")
	//
	//abi, err := abi.JSON(strings.NewReader(definition))
	//if err != nil {
	//	b.Fatal(err)
	//}
	//
	//cpurchase, err := abi.Pack("confirmPurchase")
	//if err != nil {
	//	b.Fatal(err)
	//}
	//creceived, err := abi.Pack("confirmReceived")
	//if err != nil {
	//	b.Fatal(err)
	//}
	//refund, err := abi.Pack("refund")
	//if err != nil {
	//	b.Fatal(err)
	//}
	//
	//b.ResetTimer()
	//for i := 0; i < b.N; i++ {
	//	for j := 0; j < 400; j++ {
	//		Execute(code, cpurchase, nil)
	//		Execute(code, creceived, nil)
	//		Execute(code, refund, nil)
	//	}
	//}
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
	fmt.Println("encoded(组合后):", hexcode)

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

func TestT(t *testing.T) {
	fmt.Println(common.Bytes2Hex([]byte("来自ethcall的返回")))
	fmt.Println(string(common.Hex2Bytes("313030303030303030")))
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
