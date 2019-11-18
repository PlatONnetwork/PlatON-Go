package core

import (
	"testing"
)

func TestParsePKFile(t *testing.T) {
	parsePkFile(pkFilePath)
}

func TestGenerateAccount(t *testing.T) {
	generateAccount(10, pkFilePath)
}

func TestGetNonce(t *testing.T) {
	//parseConfigJson(configPath)
	//nonce := getNonce("0x60ceca9c1290ee56b98d4e160ef0453f7c40d219")
	//fmt.Println(nonce)
}

func TestSendTransaction(t *testing.T) {

	//platon, datadir := prepare(t)
	//
	//to := "0x3058552A64Ce86aFb57806d15Fd9612a8591b01d"
	//value := "100000000000000000000000000"
	//
	//if !strings.HasPrefix(value, "0x") {
	//	intValue, _ := strconv.ParseInt(value, 10, 64)
	//	value = hexutil.EncodeBig(big.NewInt(intValue))
	//}
	//
	//hash, err := SendTransaction(from, to, value)
	//
	//assert.Nil(t, err, fmt.Sprintf("error %v", err))
	//
	//assert.NotEqual(t, hash, "", fmt.Sprintf("the transaction hash is empty"))
	//
	//clean(platon, datadir)

}

//func TestSendRawTransaction(t *testing.T) {
//
//	platon, datadir := prepare(t)
//
//	hash, err := SendRawTransaction(from, "0xD71DaAA3ce55F52a4D820460d09A801C5D487a16", "1000000", pkFilePath)
//
//	assert.Nil(t, err, fmt.Sprintf("error %v", err))
//
//	assert.NotEqual(t, hash, "", fmt.Sprintf("the transaction hash is empty"))
//
//	clean(platon, datadir)
//}

//func TestGetTxReceipt(t *testing.T) {
//	platon, datadir := prepare(t)
//
//	to := "0x3058552A64Ce86aFb57806d15Fd9612a8591b01d"
//	value := "100000000000000000000000000"
//
//	if !strings.HasPrefix(value, "0x") {
//		intValue, _ := strconv.ParseInt(value, 10, 64)
//		value = hexutil.EncodeBig(big.NewInt(intValue))
//	}
//
//	hash, err := SendTransaction(from, to, value)
//
//	assert.Nil(t, err, fmt.Sprintf("error %v", err))
//
//	assert.NotEqual(t, hash, "", fmt.Sprintf("the transaction hash is empty"))
//
//	time.Sleep(10 * time.Second)
//
//	//
//	r, err := GetTxReceipt(hash)
//	assert.Nil(t, err, fmt.Sprintf("error %v", err))
//
//	fmt.Println(r)
//
//	assert.NotEqual(t, r, "", fmt.Sprintf("the transaction hash is empty"))
//
//	clean(platon, datadir)
//}
