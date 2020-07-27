package platonstats

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"os"
	"regexp"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/log"

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

	common.CollectRestrictingReleaseItem(blockNumber, address, big.NewInt(111), common.Big0)
	common.CollectUnstakingRefundItem(blockNumber, common.NodeID(nodeId), nodeAddress, 222)
	common.CollectDuplicatedSignSlashingSetting(blockNumber, 2000, 60)

	rewardData := &common.RewardData{BlockRewardAmount: big.NewInt(12), StakingRewardAmount: big.NewInt(12), CandidateInfoList: candidateInfoList}
	common.CollectRewardData(blockNumber, rewardData)

	value1, _ := new(big.Int).SetString("3000000000000000000000", 10)
	//value2, _ := new(big.Int).SetString("3000000000000000000000", 10)
	//value3, _ := new(big.Int).SetString("3000000000000000000000", 10)

	additionalIssuance := new(common.AdditionalIssuanceData)
	additionalIssuance.AdditionalNo = 1
	additionalIssuance.AdditionalBase = big.NewInt(1000000)
	additionalIssuance.AdditionalAmount = big.NewInt(100000)
	additionalIssuance.AdditionalRate = 10
	additionalIssuance.AddIssuanceItem(vm.RewardManagerPoolAddr, big.NewInt(10000))
	common.CollectAdditionalIssuance(blockNumber, additionalIssuance)

	common.CollectEmbedTransferTx(blockNumber, common.Hash{0x01}, address, address, value1)
	common.CollectEmbedTransferTx(blockNumber, common.Hash{0x01}, address, address, value1)
	common.CollectEmbedTransferTx(blockNumber, common.Hash{0x01}, address, address, value1)
	common.CollectEmbedContractTx(blockNumber, common.Hash{0x03}, address, address, []byte{0x01, 0x02, 0x03, 0x04, 0x05})

	return common.GetExeBlockData(blockNumber)
}

func Test_Unmarshal_accountCheckingMessage(t *testing.T) {
	message := AccountCheckingMessage{
		BlockNumber: uint64(12132131),
		AccountList: []*AccountItem{
			&AccountItem{Addr: address, Balance: big.NewInt(1)},
			&AccountItem{Addr: address, Balance: big.NewInt(2)},
		}}
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		t.Fatal("Failed to marshal accountCheckingMessage to json format", err)
	} else {
		t.Log("accountCheckingMessage json format:" + string(jsonBytes))
	}

	/*jsonStr := "{\n    \"blockNumber\":12132131,\n    \"accountList\":[\n        {\n            \"addr\":\"lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj\",\n            \"balance\":1\n        },\n        {\n            \"addr\":\"lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj\",\n            \"balance\":2\n        }]\n}"
	jsonBytes = []byte(jsonStr)*/
	var data AccountCheckingMessage
	if len(jsonBytes) > 0 {
		if err := json.Unmarshal(jsonBytes, &data); err != nil {
			t.Fatal("Failed to unmarshal json to accountCheckingMessage", err)
		} else {
			t.Log("accountCheckingMessage.accountList", data.AccountList[1].Balance.Int64())
		}
	}
}
func Test_encode_Data(t *testing.T) {

	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

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
				t.Log("EmbedTransferTxList[0].amount", data.EmbedTransferTxList[0].Amount)
				t.Log("EmbedTransferTxList[1].amount", data.EmbedTransferTxList[1].Amount)
			}
		}
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
