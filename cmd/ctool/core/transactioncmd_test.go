package core

import (
	"fmt"
	"Platon-go/common/hexutil"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func TestParsePKFile(t *testing.T) {
	parsePkFile(pkFilePath)
}

func TestGenerateAccount(t *testing.T) {
	generateAccount(10, pkFilePath)
}

func TestGetNonce(t *testing.T) {
	parseConfigJson(configPath)
	nonce := getNonce("0x60ceca9c1290ee56b98d4e160ef0453f7c40d219")
	fmt.Println(nonce)
}

func TestSendTransaction(t *testing.T) {

	parseConfigJson(configPath)

	from := "0x60ceca9c1290ee56b98d4e160ef0453f7c40d219"
	to := "0x3058552A64Ce86aFb57806d15Fd9612a8591b01d"
	value := "100000000000000000000000000"

	if !strings.HasPrefix(value, "0x") {
		intValue, _ := strconv.ParseInt(value, 10, 64)
		value = hexutil.EncodeBig(big.NewInt(intValue))
	}

	hash, err := SendTransaction(from, to, value)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	if hash == "" {
		t.Fatalf("error get transaction hash ")
	}
	fmt.Printf(hash)

}

func TestSendRawTransaction(t *testing.T) {

	parsePkFile(pkFilePath)
	parseConfigJson(configPath)

	hash, err := SendRawTransaction("0x9A7313f7868D9452c8d914A38340701c448072B6", "0xCF2efb592aA9FF75B5814A1211449df03B533657", "1000000", pkFilePath)
	if err != nil {
		t.Fatalf("send error,%s.", err.Error())
	}

	fmt.Printf(hash)
}
