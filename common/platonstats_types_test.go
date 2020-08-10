package common

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

var (
	address     = MustBech32ToAddress("lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj")
	nodeAddress = NodeAddress(address)
	nodeId      = MustHexID("0x362003c50ed3a523cdede37a001803b8f0fed27cb402b3d6127a1a96661ec202318f68f4c76d9b0bfbabfd551a178d4335eaeaa9b7981a4df30dfc8c0bfe3384")
)

func Test_encode_Data(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	bigint, _ := new(big.Int).SetString("4238288267758744845421755", 10)
	t.Log("bigint::::::::::::", bigint)

	content := fmt.Sprintf("chainBalance=%d", bigint)
	t.Log("content::::::::::::", content)

	exeData := buildExeBlockData()
	bytes := Hex2Bytes("f855838203ec8180b842b84027eff1a24cfd76e7151c0410d6a2c9fe9660b0906862fe3c461f5bae4ce5893d4f975239187510ff4b399ec1daea4f71690e00de3de0c56a2bd896cf5cd37eca8a8902b5e3af16b1880000")
	t.Log("input::::::::::::", Bytes2Hex(bytes))

	jsonBytes, err := json.Marshal(exeData)
	if err != nil {
		t.Fatal("Failed to marshal exeData to json format", err)
	} else {
		t.Log("json format:" + string(jsonBytes))

		var data ExeBlockData
		if len(jsonBytes) > 0 {
			if err := json.Unmarshal(jsonBytes, &data); err != nil {
				t.Fatal("Failed to unmarshal json to statsData", err)
			} else {
				t.Log("ExeBlockData.RewardData.CandidateInfoList[0].NodeID", Bytes2Hex(data.RewardData.CandidateInfoList[0].NodeID[:]))
				t.Log("AdditionalIssuanceData==nil", data.AdditionalIssuanceData == nil)
				t.Log("EmbedTransferTxList[0].amount", data.EmbedTransferTxList[0].Amount)
				t.Log("EmbedTransferTxList[1].amount", data.EmbedTransferTxList[1].Amount)
			}
		}
	}
}

func buildExeBlockData() *ExeBlockData {
	blockNumber := uint64(100)

	InitExeBlockData(blockNumber)

	candidate := &CandidateInfo{nodeId, address}
	candidateInfoList := []*CandidateInfo{candidate}

	CollectRestrictingReleaseItem(blockNumber, address, big.NewInt(111), Big0)
	CollectStakingFrozenItem(blockNumber, nodeId, nodeAddress, 222, true)
	CollectStakingFrozenItem(blockNumber, nodeId, nodeAddress, 111, false)
	CollectDuplicatedSignSlashingSetting(blockNumber, 2000, 60)

	rewardData := &RewardData{BlockRewardAmount: big.NewInt(12), StakingRewardAmount: big.NewInt(12), CandidateInfoList: candidateInfoList}
	CollectRewardData(blockNumber, rewardData)

	value1, _ := new(big.Int).SetString("3000000000000000000000", 10)
	//value2, _ := new(big.Int).SetString("3000000000000000000000", 10)
	//value3, _ := new(big.Int).SetString("3000000000000000000000", 10)

	additionalIssuance := new(AdditionalIssuanceData)
	additionalIssuance.AdditionalNo = 1
	additionalIssuance.AdditionalBase = big.NewInt(1000000)
	additionalIssuance.AdditionalAmount = big.NewInt(100000)
	additionalIssuance.AdditionalRate = 10
	additionalIssuance.AddIssuanceItem(HexToAddress("0x1000000000000000000000000000000000000003"), big.NewInt(10000))
	CollectAdditionalIssuance(blockNumber, additionalIssuance)

	CollectEmbedTransferTx(blockNumber, Hash{0x01}, address, address, value1)
	CollectEmbedTransferTx(blockNumber, Hash{0x01}, address, address, value1)
	CollectEmbedTransferTx(blockNumber, Hash{0x01}, address, address, value1)
	CollectEmbedContractTx(blockNumber, Hash{0x03}, address, address, Hex2Bytes("f855838203ec8180b842b84027eff1a24cfd76e7151c0410d6a2c9fe9660b0906862fe3c461f5bae4ce5893d4f975239187510ff4b399ec1daea4f71690e00de3de0c56a2bd896cf5cd37eca8a8902b5e3af16b1880000"))

	return GetExeBlockData(blockNumber)
}
