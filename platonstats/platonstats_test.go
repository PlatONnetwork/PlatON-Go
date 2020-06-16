package platonstats

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/statsdb"

	"gotest.tools/assert"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/Shopify/sarama"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	address     = common.MustBech32ToAddress("lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj")
	nodeAddress = common.NodeAddress(address)
	nodeId      = discover.MustHexID("0x362003c50ed3a523cdede37a001803b8f0fed27cb402b3d6127a1a96661ec202318f68f4c76d9b0bfbabfd551a178d4335eaeaa9b7981a4df30dfc8c0bfe3384")
)

func TestUrl(t *testing.T) {
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

	candidate := &common.CandidateInfo{nodeId, address}
	candidateInfoList := []*common.CandidateInfo{candidate}

	common.CollectRestrictingReleaseItem(blockNumber, address, 111)
	common.CollectUnstakingRefundItem(blockNumber, nodeId, nodeAddress, 222)
	common.CollectDuplicatedSignSlashingSetting(blockNumber, 2000, 60)

	rewardData := &common.RewardData{BlockRewardAmount: 12, StakingRewardAmount: 12, CandidateInfoList: candidateInfoList}
	common.CollectRewardData(blockNumber, rewardData)

	return common.GetExeBlockData(blockNumber)
}
func Test_rlp_Data(t *testing.T) {
	NewMockPlatonStatsService()
	blockData := buildExeBlockData()
	json, _ := json.Marshal(blockData)
	t.Log("blockData", string(json))

	bytes := common.MustRlpEncode(blockData)
	t.Log("encode data", "hex", common.Bytes2Hex(bytes))

	var data common.ExeBlockData
	if len(bytes) > 0 {
		if err := rlp.DecodeBytes(bytes, &data); err != nil {
			t.Fatal("Failed to rlp decode bytes to ExeBlockData", err)
		} else {
			t.Log("ExeBlockData.RewardData.CandidateInfoList[0].NodeID", common.Bytes2Hex(data.RewardData.CandidateInfoList[0].NodeID[:]))
			t.Log("AdditionalIssuanceData==nil", data.AdditionalIssuanceData == nil)

		}
	}
}

func Test_Kafka_producer(t *testing.T) {
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

func Test_StatsService(t *testing.T) {
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

func Test_Log(t *testing.T) {
	blockNumber := uint64(121)
	writeBlockNumber(blockNumber)

	if blockNo, err := readBlockNumber(); err != nil {
		t.Fatal("Failed to read stats service log", "error", err)
	} else {
		t.Log("read the number from log", "number", blockNo)
		assert.Equal(t, blockNumber, blockNo)
	}

}
