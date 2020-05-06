// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

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
