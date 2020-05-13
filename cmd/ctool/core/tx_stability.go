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
	"crypto/ecdsa"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/secp256k1"
	"gopkg.in/urfave/cli.v1"
)

var (
	//accountPool         = make(map[string]*PriAccount)
	accountPool         = make(map[common.Address]*PriAccount)
	StressTransferValue = 1000
	txCh                = make(chan *types.Transaction, 10)
	wg                  = &sync.WaitGroup{}

	DefaultPrivateKeyFilePath  = "./test/privateKeys.txt"
	DefaultAccountAddrFilePath = "./test/addr.json"

	StabilityCmd = cli.Command{
		Name:    "stability",
		Aliases: []string{"stab"},
		Usage:   "start stability test ",
		Action:  stabilityTest,
		Flags:   stabilityCmdFlags,
	}
	StabPrepareCmd = cli.Command{
		Name:    "prepare",
		Aliases: []string{"pre"},
		Usage:   "prepare some accounts are used for stability test ",
		Action:  prepareAccount,
		Flags:   stabPrepareCmdFlags,
	}
)

func prepareAccount(c *cli.Context) {
	pkFile := c.String(PKFilePathFlag.Name)
	size := c.Int(AccountSizeFlag.Name)
	value := c.String(TransferValueFlag.Name)

	parseConfigJson(c.String(ConfigPathFlag.Name))

	err := PrepareAccount(size, pkFile, value)
	if err != nil {
		panic(fmt.Errorf("send raw transaction error,%s", err.Error()))
	}
}

func stabilityTest(c *cli.Context) {
	pkFile := c.String(PKFilePathFlag.Name)
	times := c.Int(StabExecTimesFlag.Name)
	interval := c.Int(SendTxIntervalFlag.Name)

	parseConfigJson(c.String(ConfigPathFlag.Name))

	err := StabilityTest(pkFile, times, interval)
	if err != nil {
		panic(fmt.Errorf("stress test error,%s", err.Error()))
	}
}

type PriAccount struct {
	Priv    *ecdsa.PrivateKey
	Nonce   uint64
	Address common.Address
}

func generateAccount(size int, pkFile string) {
	addrs := make([]string, size)
	for i := 0; i < size; i++ {
		privateKey, _ := crypto.GenerateKey()
		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		accountPool[address] = &PriAccount{privateKey, 0, address}
		addrs[i] = address.String()
	}
	savePrivateKeyPool(pkFile)
	saveAddrs(addrs, pkFile)

}

func savePrivateKeyPool(pkFile string) {
	if pkFile == "" {
		pkFile = DefaultPrivateKeyFilePath
	}
	gob.Register(&secp256k1.BitCurve{})
	file, err := os.Create(pkFile)
	if err != nil {
		panic(fmt.Errorf("save private key err,%s,%s", pkFile, err.Error()))
	}
	os.Truncate(pkFile, 0)
	enc := gob.NewEncoder(file)
	err = enc.Encode(accountPool)
	if err != nil {
		panic(err.Error())
	}
}

func saveAddrs(addrs []string, pkFile string) {
	addrsPath := DefaultAccountAddrFilePath
	if pkFile != "" {
		addrsPath = filepath.Dir(pkFile) + "/addr.json"
	}
	os.Truncate(DefaultAccountAddrFilePath, 0)
	byts, err := json.MarshalIndent(addrs, "", "\t")
	_, err = os.Create(addrsPath)
	if err != nil {
		panic(fmt.Errorf("create addr.json error%s \n", err.Error()))
	}
	err = ioutil.WriteFile(addrsPath, byts, 0644)
	if err != nil {
		panic(fmt.Errorf("write to addr.json error%s \n", err.Error()))
	}
}

func PrepareAccount(size int, pkFile, value string) error {

	if len(accountPool) == 0 {
		generateAccount(size, pkFile)
	}

	for addr := range accountPool {
		hash, err := SendTransaction(from, addr.String(), value)
		if err != nil {
			return fmt.Errorf("prepare error,send from coinbase error,%s", err.Error())
		}
		fmt.Printf("transfer hash: %s \n", hash)
	}
	fmt.Printf("prepare %d account finish...", size)
	return nil
}

func StabilityTest(pkFile string, times, interval int) error {
	if len(accountPool) == 0 {
		parsePkFile(pkFile)
	}

	addrs := getAllAddress(pkFile)

	for i := 0; i < times; i++ {
		if interval != 0 {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
		from, to := getRandomAddr(addrs)
		if from == "" || to == "" {
			continue
		}

		acc, ok := accountPool[common.MustBech32ToAddress(from)]
		if !ok {
			return fmt.Errorf("private key not found,addr:%s", from)
		}

		wg.Add(1)
		go getTransactionGo(acc, from, to)

		wg.Add(1)
		go sendTransactionGo()
	}

	wg.Wait()

	savePrivateKeyPool(pkFile)

	return nil
}

func sendTransactionGo() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("send raw transaction failed:%s \n", err)
			//debug.PrintStack()
			wg.Done()
		}
	}()

	hash, _ := sendRawTransaction(<-txCh)
	fmt.Printf("tx hash：%s\n", hash)
	wg.Done()
}

func getTransactionGo(acc *PriAccount, from, to string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("get transaction error:", err)
			//debug.PrintStack()
			wg.Done()
		}
	}()

	newTx := getSignedTransaction(from, to, int64(StressTransferValue), acc.Priv, acc.Nonce)
	lock := &sync.Mutex{}
	lock.Lock()

	acc.Nonce++
	fmt.Printf("add =%s,Nonce:%d \n", from, acc.Nonce)

	lock.Unlock()

	txCh <- newTx
	wg.Done()
}

func parsePkFile(pkFile string) {
	if pkFile == "" {
		dir, _ := os.Getwd()
		pkFile = dir + DefaultPrivateKeyFilePath
	}
	//gob.Register(&secp256k1.BitCurve{})
	//file, err := os.Open(pkFile)
	//dec := gob.NewDecoder(file)
	//err2 := dec.Decode(&accountPool)
	//if err2 != nil {
	//	panic(err.Error())
	//}

	gob.Register(&secp256k1.BitCurve{})
	file, err := os.Open(pkFile)
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(file)
	err2 := dec.Decode(&accountPool)
	if err2 != nil {
		panic(err2)
	}

}

func getAllAddress(pkFile string) []string {
	addrsPath := ""
	if pkFile != "" {
		addrsPath = filepath.Dir(pkFile) + "/addr.json"
	} else {
		dir, _ := os.Getwd()
		addrsPath = dir + DefaultAccountAddrFilePath
	}

	bytes, err := ioutil.ReadFile(addrsPath)
	if err != nil {
		panic(fmt.Errorf("get all address array error,%s \n", err.Error()))
	}
	var addrs []string
	err = json.Unmarshal(bytes, &addrs)
	if err != nil {
		panic(fmt.Errorf("parse address to array error,%s \n", err.Error()))
	}

	return addrs
}

func getRandomAddr(addrs []string) (string, string) {
	if len(addrs) == 0 {
		return "", ""
	}
	fromIndex := rand.Intn(len(addrs))
	toIndex := rand.Intn(len(addrs))
	for toIndex == fromIndex {
		toIndex = rand.Intn(len(addrs))
	}
	return addrs[fromIndex], addrs[toIndex]
}
