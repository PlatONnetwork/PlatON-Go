package core

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"gopkg.in/urfave/cli.v1"
	"math/big"
	"strconv"
	"strings"
)

var (
	SendTransactionCmd = cli.Command{
		Name:   "sendTransaction",
		Usage:  "send a transaction",
		Action: sendTransactionCmd,
		Flags:  sendTransactionCmdFlags,
	}
	SendRawTransactionCmd = cli.Command{
		Name:   "sendRawTransaction",
		Usage:  "send a raw transaction",
		Action: sendRawTransactionCmd,
		Flags:  sendRawTransactionCmdFlags,
	}
	GetTxReceiptCmd = cli.Command{
		Name:   "getTxReceipt",
		Usage:  "get transaction receipt by hash",
		Action: getTxReceiptCmd,
		Flags:  getTxReceiptCmdFlags,
	}
)

func getTxReceiptCmd(c *cli.Context) {
	hash := c.String(TransactionHashFlag.Name)
	parseConfigJson(c.String(ConfigPathFlag.Name))
	GetTxReceipt(hash)
}

func GetTxReceipt(txHash string) (Receipt, error) {
	var receipt = Receipt{}
	//根据result获取交易receipt
	res, _ := Send([]string{txHash}, "eth_getTransactionReceipt")
	e := json.Unmarshal([]byte(res), &receipt)
	if e != nil {
		panic(fmt.Sprintf("parse get receipt result error ! \n %s", e.Error()))
	}

	if receipt.Result.BlockHash == "" {
		panic("no receipt found")
	}
	out, _ := json.MarshalIndent(receipt, "", "  ")
	fmt.Println(string(out))
	return receipt, nil
}

func sendTransactionCmd(c *cli.Context) error {
	from := c.String(TxFromFlag.Name)
	to := c.String(TxToFlag.Name)
	value := c.String(TransferValueFlag.Name)
	parseConfigJson(c.String(ConfigPathFlag.Name))

	hash, err := SendTransaction(from, to, value)
	if err != nil {
		utils.Fatalf("Send transaction error: %v", err)
	}

	fmt.Printf("tx hash: %s", hash)
	return nil
}

func sendRawTransactionCmd(c *cli.Context) error {
	from := c.String(TxFromFlag.Name)
	to := c.String(TxToFlag.Name)
	value := c.String(TransferValueFlag.Name)
	pkFile := c.String(PKFilePathFlag.Name)

	parseConfigJson(c.String(ConfigPathFlag.Name))

	hash, err := SendRawTransaction(from, to, value, pkFile)
	if err != nil {
		utils.Fatalf("Send transaction error: %v", err)
	}

	fmt.Printf("tx hash: %s", hash)
	return nil
}

func SendTransaction(from, to, value string) (string, error) {
	var tx TxParams
	if from == "" {
		from = config.From
	}
	tx.From = from
	tx.To = to
	tx.Gas = config.Gas
	tx.GasPrice = config.GasPrice

	if !strings.HasPrefix(value, "0x") {
		intValue, _ := strconv.ParseInt(value, 10, 64)
		value = hexutil.EncodeBig(big.NewInt(intValue))
	}
	tx.Value = value

	params := make([]TxParams, 1)
	params[0] = tx

	res, _ := Send(params, "eth_sendTransaction")
	response := parseResponse(res)

	return response.Result, nil
}

func SendRawTransaction(from, to, value string, pkFile string) (string, error) {
	if len(accountPool) == 0 {
		parsePkFile(pkFile)
	}
	var v int64
	if strings.HasPrefix(value, "0x") {
		bigValue, _ := hexutil.DecodeBig(value)
		v = bigValue.Int64()
	} else {
		v, _ = strconv.ParseInt(value, 10, 64)
	}

	acc, ok := accountPool[from]
	if !ok {
		return "", fmt.Errorf("private key not found in private key file,addr:%s", from)
	}
	nonce := getNonce(from)
	nonce++
	newTx := getSignedTransaction(from, to, v, acc.Priv, nonce)

	hash, err := sendRawTransaction(newTx)
	if err != nil {
		panic(err)
	}
	return hash, nil
}

func sendRawTransaction(transaction *types.Transaction) (string, error) {
	bytes, _ := rlp.EncodeToBytes(transaction)
	res, err := Send([]string{hexutil.Encode(bytes)}, "eth_sendRawTransaction")
	if err != nil {
		panic(err)
	}
	response := parseResponse(res)

	return response.Result, nil
}

func getSignedTransaction(from, to string, value int64, priv *ecdsa.PrivateKey, nonce uint64) *types.Transaction {
	newTx, err := types.SignTx(types.NewTransaction(nonce, common.HexToAddress(to), big.NewInt(value), 100000, big.NewInt(90000), []byte{}), types.HomesteadSigner{}, priv)
	if err != nil {
		panic(fmt.Errorf("sign error,%s", err.Error()))
	}
	return newTx
}

func getNonce(addr string) uint64 {
	res, _ := Send([]string{addr, "latest"}, "eth_getTransactionCount")
	response := parseResponse(res)
	nonce, _ := hexutil.DecodeBig(response.Result)
	fmt.Println(addr, nonce)
	return nonce.Uint64()
}

//func getCoinbase() (error) {
//	res, _ := Send([]string{}, "eth_coinbase")
//	response := parseResponse(res)
//	coinBase = response.Result
//	return nil
//}
