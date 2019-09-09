package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/pborman/uuid"
	"gopkg.in/urfave/cli.v1"
)

type outputFound struct {
	Address      string
	AddressEIP55 string
}

var commandFound = cli.Command{
	Name:      "found",
	Usage:     "generate a set of accounts for PlatON network",
	ArgsUsage: "<summaryFile>",
	Description: `
Generate a set of for PlatON network.

All these accounts info will be saved into summary file.
`,
	Flags: []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		// Check if textFileName path given and make sure it doesn't already exist.
		summaryFile := ctx.Args().First()
		if _, err := os.Stat(summaryFile); err == nil {
			utils.Fatalf("File already exists at %s.", summaryFile)
		} else if !os.IsNotExist(err) {
			utils.Fatalf("Error checking if file exists: %v", err)
		}

		summaryFilePath := filepath.Dir(summaryFile)
		if err := os.MkdirAll(summaryFilePath, 0700); err != nil {
			utils.Fatalf("Could not create directory %s", filepath.Dir(summaryFilePath))
		}

		keystoreFilePath := filepath.Join(summaryFilePath, string(filepath.Separator), "keystore")
		if err := os.MkdirAll(keystoreFilePath, 0700); err != nil {
			utils.Fatalf("Could not create directory %s", filepath.Dir(keystoreFilePath))
		}

		summaryFileObj, err := os.OpenFile(summaryFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			utils.Fatalf("Could not create directory %s", filepath.Dir(summaryFile))
		}
		defer closeSummaryFile(summaryFileObj)

		//generate key for Main Account
		mainAccountKeystoreFile := filepath.Join(keystoreFilePath, string(filepath.Separator), "main_account_keystore.json")
		if key, keystore, password, err := genKey(); err != nil {
			utils.Fatalf("Failed to generate main account: %v", err)
		} else if err := ioutil.WriteFile(mainAccountKeystoreFile, keystore, 0600); err != nil {
			utils.Fatalf("Failed to write main account keystore file %s: %v", mainAccountKeystoreFile, err)
		} else if err := gatherAccountInfo(summaryFileObj, key, mainAccountKeystoreFile, password, "main_account"); err != nil {
			utils.Fatalf("Failed to write main account info into file %s: %v", summaryFileObj, err)
		}

		//generate key for Foundation Account
		foundationAccountKeystoreFile := filepath.Join(keystoreFilePath, string(filepath.Separator), "foundation_account_keystore.json")
		if key, keystore, password, err := genKey(); err != nil {
			utils.Fatalf("Failed to generate foundation account: %v", err)
		} else if err := ioutil.WriteFile(foundationAccountKeystoreFile, keystore, 0600); err != nil {
			utils.Fatalf("Failed to write foundation account keystore file %s: %v", foundationAccountKeystoreFile, err)
		} else if err := gatherAccountInfo(summaryFileObj, key, foundationAccountKeystoreFile, password, "foundation_account"); err != nil {
			utils.Fatalf("Failed to write foundation account info into file %s: %v", summaryFileObj, err)
		}

		//generate staking account for each node
		for i := 1; i <= 70; i++ {
			nodeKeystoreFile := filepath.Join(keystoreFilePath, string(filepath.Separator), "staking_node_"+fmt.Sprintf("%d", i)+"_keystore.json")
			if key, keystore, password, err := genKey(); err != nil {
				utils.Fatalf("Failed to generate node staking account: %v", err)
			} else if err := ioutil.WriteFile(nodeKeystoreFile, keystore, 0600); err != nil {
				utils.Fatalf("Failed to write node staking account keystore file %s: %v", nodeKeystoreFile, err)
			} else if err := gatherAccountInfo(summaryFileObj, key, nodeKeystoreFile, password, "staking_node_"+fmt.Sprintf("%d", i)+"_account"); err != nil {
				utils.Fatalf("Failed to write node staking account info into file %s: %v", summaryFileObj, err)
			}
		}

		//generate reward account for each node
		for i := 1; i <= 70; i++ {
			nodeKeystoreFile := filepath.Join(keystoreFilePath, string(filepath.Separator), "reward_node_"+fmt.Sprintf("%d", i)+"_keystore.json")
			if key, keystore, password, err := genKey(); err != nil {
				utils.Fatalf("Failed to generate node reward account: %v", err)
			} else if err := ioutil.WriteFile(nodeKeystoreFile, keystore, 0600); err != nil {
				utils.Fatalf("Failed to write node reward account keystore file %s: %v", nodeKeystoreFile, err)
			} else if err := gatherAccountInfo(summaryFileObj, key, nodeKeystoreFile, password, "reward_node_"+fmt.Sprintf("%d", i)+"_account"); err != nil {
				utils.Fatalf("Failed to write node reward account info into file %s: %v", summaryFileObj, err)
			}
		}
		return nil
	},
}

func genKey() (k *keystore.Key, ks []byte, p string, e error) {
	var privateKey *ecdsa.PrivateKey
	var err error
	// generate random.
	privateKey, err = crypto.GenerateKey()
	if err != nil {
		utils.Fatalf("Failed to generate random private key: %v", err)
	}

	// Create the keyfile object with a random UUID.
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}
	passphrase := randomPwd()

	keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	return key, keyjson, passphrase, err
}

func randomPwd() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" + "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func gatherAccountInfo(f *os.File, key *keystore.Key, keystoreFile string, password string, accountName string) error {

	PublicKey := hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey)[1:])
	PrivateKey := hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))

	data := []byte(accountName + ":\n" + key.Address.Hex() + "	" + PublicKey + "	" + PrivateKey + "	" + keystoreFile + "	" + password + "\n")

	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return err
}

func closeSummaryFile(f *os.File) {
	if err := f.Close(); err != nil {
		utils.Fatalf("Failed to close Account info file: %v", err)
	}
}
