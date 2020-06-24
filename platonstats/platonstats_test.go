package platonstats

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/PlatONnetwork/PlatON-Go/core/statsdb"

	"gotest.tools/assert"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/Shopify/sarama"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	address     = common.MustBech32ToAddress("lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj")
	nodeAddress = common.NodeAddress(address)
	nodeId      = discover.MustHexID("0x362003c50ed3a523cdede37a001803b8f0fed27cb402b3d6127a1a96661ec202318f68f4c76d9b0bfbabfd551a178d4335eaeaa9b7981a4df30dfc8c0bfe3384")
)

func T2estUrl(t *testing.T) {
	re := regexp.MustCompile("([^:@]*)(:([^@]*))?@(.+)")
	url := "center:myPasswordd@ws://localhost:1900"
	parts := re.FindStringSubmatch(url)
	for i := 0; i < len(parts); i++ {
		t.Logf("url parts: [%d]%s", i, parts[i])
	}
}

func buildExeBlockData() *common.ExeBlockData {
	blockNumber := uint64(100)

	common.InitExeBlockData(blockNumber)

	candidate := &common.CandidateInfo{common.NodeID(nodeId), address}
	candidateInfoList := []*common.CandidateInfo{candidate}

	common.CollectRestrictingReleaseItem(blockNumber, address, big.NewInt(111))
	common.CollectUnstakingRefundItem(blockNumber, common.NodeID(nodeId), nodeAddress, 222)
	common.CollectDuplicatedSignSlashingSetting(blockNumber, 2000, 60)

	rewardData := &common.RewardData{BlockRewardAmount: big.NewInt(12), StakingRewardAmount: big.NewInt(12), CandidateInfoList: candidateInfoList}
	common.CollectRewardData(blockNumber, rewardData)

	common.CollectEmbedTransferTx(blockNumber, common.Hash{0x01}, address, address, big.NewInt(1))
	common.CollectEmbedTransferTx(blockNumber, common.Hash{0x02}, address, address, big.NewInt(2))
	common.CollectEmbedContractTx(blockNumber, common.Hash{0x03}, address, address, []byte{0x01, 0x02, 0x03, 0x04, 0x05})

	return common.GetExeBlockData(blockNumber)
}
func T2est_encode_Data(t *testing.T) {
	NewMockPlatonStatsService()
	exeData := buildExeBlockData()

	jsonBytes, err := json.Marshal(exeData)
	if err != nil {
		t.Fatal("Failed to marshal exeData to json format", err)
	} else {
		t.Log("json format:" + string(jsonBytes))

		var data common.ExeBlockData
		if len(jsonBytes) > 0 {
			if err := json.Unmarshal(jsonBytes, &data); err != nil {
				t.Fatal("Failed to unmarshal json to statsData", err)
			} else {
				t.Log("ExeBlockData.RewardData.CandidateInfoList[0].NodeID", common.Bytes2Hex(data.RewardData.CandidateInfoList[0].NodeID[:]))
				t.Log("AdditionalIssuanceData==nil", data.AdditionalIssuanceData == nil)
			}
		}
	}

}

func Test_encode_Data(t *testing.T) {
	//NewMockPlatonStatsService()
	exeData := buildExeBlockData()

	jsonBytes, err := json.Marshal(exeData)
	if err != nil {
		t.Fatal("Failed to marshal exeData to json format", err)
	} else {
		t.Log("json format:" + string(jsonBytes))

		var data common.ExeBlockData
		if len(jsonBytes) > 0 {
			if err := json.Unmarshal(jsonBytes, &data); err != nil {
				t.Fatal("Failed to unmarshal json to statsData", err)
			} else {
				t.Log("ExeBlockData.RewardData.CandidateInfoList[0].NodeID", common.Bytes2Hex(data.RewardData.CandidateInfoList[0].NodeID[:]))
				t.Log("AdditionalIssuanceData==nil", data.AdditionalIssuanceData == nil)
			}
		}
	}

}

func T2est_Kafka_producer(t *testing.T) {
	s := NewMockPlatonStatsService()

	statsLogFile = "d:\\swap\\statsdb\\platonstats.log"
	statsdb.SetDBPath("d:\\swap\\statsdb")

	var blockProducer sarama.SyncProducer
	var err error
	if blockProducer, err = sarama.NewSyncProducer([]string{"192.168.112.32:9092"}, blockProducerConfig()); err != nil {
		t.Fatal("Failed to create kafka sync producer", "error", err)
	}
	s.blockProducer = blockProducer

	defer func() {
		if blockProducer != nil {
			blockProducer.Close()
		}
	}()

	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	nextBlock := s.BlockChain().GetBlockByNumber(10)

	statsdb.Instance().WriteExeBlockData(nextBlock.Number(), buildExeBlockData())

	if err = s.reportBlockMsg(nextBlock); err != nil {
		t.Fatal("Failed to report BlockMsg", "error", err)
	} else {
		t.Log("ok.......")
	}
}

func T2est_StatsService(t *testing.T) {
	s := NewMockPlatonStatsService()

	s.kafkaUrl = "192.168.112.32:9092"
	statsLogFile = "d:\\swap\\platonstats.log"
	statsdb.SetDBPath("d:\\swap\\statsdb")

	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	err := s.Start(nil)
	if err != nil {
		t.Fatal("Failed to start stats service", "error", err)
	}

	//合建chan
	c := make(chan int)
	//阻塞直到有信号传入
	fmt.Println("启动")
	q := <-c
	fmt.Println("退出信号", q)
}

func T2est_Log(t *testing.T) {
	blockNumber := uint64(121)
	writeBlockNumber(blockNumber)

	if blockNo, err := readBlockNumber(); err != nil {
		t.Fatal("Failed to read stats service log", "error", err)
	} else {
		t.Log("read the number from log", "number", blockNo)
		assert.Equal(t, blockNumber, blockNo)
	}
}

// 交易发起方keystore文件地址
var fromKeyStoreFile = "D:\\swap\\keystore\\UTC--2020-06-15T06-46-38.833974074Z--493301712671ada506ba6ca7891f436d29185821"

// keystore文件对应的密码
var password = "88888888"

// http服务地址, 例:http://localhost:8545
var httpUrl = "http://192.168.112.31:6901"

type PriAccount struct {
	Priv    *ecdsa.PrivateKey
	Nonce   uint64
	Address common.Address
}

func generateAccount(size int) []*PriAccount {
	addrs := make([]*PriAccount, size)
	for i := 0; i < size; i++ {
		privateKey, _ := crypto.GenerateKey()
		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		addrs[i] = &PriAccount{privateKey, 0, address}
	}
	return addrs
}

func T2est_SendTx(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	accountsize := 1000
	signer := types.NewEIP155Signer(big.NewInt(1021))

	// 数量
	initAmount, _ := new(big.Int).SetString("100000000000000000000", 10)
	// 数量
	amount := big.NewInt(1)

	// gasLimit
	var gasLimit uint64 = 30000

	// gasPrice
	var gasPrice = new(big.Int).SetInt64(10000)

	// 创建客户端
	client, err := ethclient.Dial(httpUrl)
	require.NoError(t, err)

	// 交易发送方
	// 获取私钥方式一，通过keystore文件
	fromKeystore, err := ioutil.ReadFile(fromKeyStoreFile)
	require.NoError(t, err)
	fromKey, err := keystore.DecryptKey(fromKeystore, password)
	fromPrivkey := fromKey.PrivateKey
	fromPubkey := fromPrivkey.PublicKey
	fromAddr := crypto.PubkeyToAddress(fromPubkey)

	// nonce获取
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)

	accounts := generateAccount(accountsize)

	//init
	for _, account := range accounts {
		// 交易创建
		tx := types.NewTransaction(nonce, account.Address, initAmount, gasLimit, gasPrice, []byte{})

		// 交易签名
		signedTx, err := types.SignTx(tx, signer, fromPrivkey)

		//signedTx ,err := types.SignTx(tx,types.HomesteadSigner{},fromPrivkey)
		require.NoError(t, err)
		// 交易发送
		serr := client.SendTransaction(context.Background(), signedTx)

		require.NoError(t, serr)

		nonce = nonce + 1
	}

	time.Sleep(30 * time.Second)

	for {
		from := accounts[rand.Int31n(int32(accountsize))]
		to := accounts[rand.Int31n(int32(accountsize))]
		for to.Address == from.Address {
			to = accounts[rand.Int31n(int32(accountsize))]
		}

		// 交易创建
		tx := types.NewTransaction(from.Nonce, accounts[rand.Int31n(int32(accountsize))].Address, amount, gasLimit, gasPrice, []byte{})

		// 交易签名
		signedTx, err := types.SignTx(tx, signer, from.Priv)

		//signedTx ,err := types.SignTx(tx,types.HomesteadSigner{},fromPrivkey)
		require.NoError(t, err)
		// 交易发送
		serr := client.SendTransaction(context.Background(), signedTx)

		require.NoError(t, serr)

		from.Nonce++
	}

	// 等待挖矿完成
	//bind.WaitMined(context.Background(), client, signedTx)
}
