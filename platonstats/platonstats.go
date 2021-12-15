package platonstats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/vm"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/PlatONnetwork/PlatON-Go/internal/ethapi"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/core/statsdb"

	"github.com/syndtr/goleveldb/leveldb/errors"

	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/eth"

	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/rpc"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

var (
	statsLogFile = "./platonstats.log"
	statsLogFlag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
)

type ExtraData []byte

func (extraData ExtraData) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0x%x", extraData))
}

type Nonce []byte

func (nonce Nonce) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0x%x", nonce))
}

func jsonBlock(block *types.Block) (map[string]interface{}, error) {
	fields, err := ethapi.RPCMarshalBlock(block, true, true)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

type Brief struct {
	BlockType   common.BlockType
	Epoch       uint64
	ConsensusNo uint64
	NodeID      common.NodeID
	NodeAddress common.Address
}

type StatsBlockExt struct {
	BlockType    common.BlockType       `json:"blockType"`
	Epoch        uint64                 `json:"epoch"`
	ConsensusNo  uint64                 `json:"consensusNo"`
	NodeID       common.NodeID          `json:"nodeID,omitempty"`
	NodeAddress  common.Address         `json:"nodeAddress,omitempty"`
	Block        map[string]interface{} `json:"block,omitempty"`
	Receipts     []*types.Receipt       `json:"receipts,omitempty"`
	ExeBlockData *common.ExeBlockData   `json:"exeBlockData,omitempty"`
	GenesisData  *common.GenesisData    `json:"GenesisData,omitempty"`
	ContractList []common.Address       `json:"ContractList,omitempty"`
	StatData     *common.StatData       `json:"statData,omitempty"`
}

type PlatonStatsService struct {
	server          *p2p.Server //Peer-to-peer server to retrieve networking infos
	kafkaUrl        string
	kafkaBlockTopic string        //统计数据消息Topic
	eth             *eth.Ethereum // Full Ethereum service if monitoring a full node
	datadir         string
	kafkaClient     *ConfluentKafkaClient
	stopSampleMsg   chan struct{}
	stopBlockMsg    chan struct{}
	stopOnce        sync.Once
	quit            chan bool
	waitQuit        sync.WaitGroup
}

var (
	platonStatsService *PlatonStatsService
)

func New(kafkaUrl, kafkaBlockTopic string, ethServ *eth.Ethereum, datadir string) (*PlatonStatsService, error) {
	log.Info("new PlatON stats service", "kafkaUrl", kafkaUrl, "kafkaBlockTopic", kafkaBlockTopic)

	waitQ := sync.WaitGroup{}
	waitQ.Add(1)

	platonStatsService = &PlatonStatsService{
		kafkaUrl:        kafkaUrl,
		kafkaBlockTopic: kafkaBlockTopic,
		eth:             ethServ,
		datadir:         datadir,
		quit:            make(chan bool),
		waitQuit:        waitQ,
	}
	if len(datadir) > 0 {
		statsLogFile = filepath.Join(datadir, statsLogFile)
	}
	log.Debug("PlatON stats log file", "datadir", datadir)
	log.Debug("PlatON stats log file", "statsLogFile", statsLogFile)
	return platonStatsService, nil
}

func (s *PlatonStatsService) BlockChain() *core.BlockChain {
	return s.eth.BlockChain()
}

func (s *PlatonStatsService) ChainDb() ethdb.Database {
	return s.eth.ChainDb()
}

// Protocols implements node.Service, returning the P2P network protocols used
// by the stats service (nil as it doesn't use the devp2p overlay network).
func (s *PlatonStatsService) Protocols() []p2p.Protocol { return nil }

// APIs implements node.Service, returning the RPC API endpoints provided by the
// stats service (nil as it doesn't provide any user callable APIs).
func (s *PlatonStatsService) APIs() []rpc.API { return nil }

// Start implements node.Service, starting up the monitoring and reporting daemon.
func (s *PlatonStatsService) Start(server *p2p.Server) error {
	log.Info("PlatON stats server starting....")
	s.server = server

	s.kafkaClient = NewConfluentKafkaClient(s.kafkaUrl, s.kafkaBlockTopic)

	go s.blockMsgLoop()

	log.Info("PlatON stats daemon started")
	return nil
}

// Stop implements node.Service, terminating the monitoring and reporting daemon.
func (s *PlatonStatsService) Stop() error {
	s.stopOnce.Do(func() {
		//close(s.stopSampleMsg)
		//close(s.stopBlockMsg)
		s.quit <- true
		s.waitQuit.Wait()
		if s.kafkaClient != nil {
			s.kafkaClient.Close()
		}
	})

	log.Info("PlatON stats daemon stopped")
	return nil
}

func (s *PlatonStatsService) blockMsgLoop() {
	var nextBlockNumber uint64
	nextBlockNumber = 0

	if loggedBlockNumber, err := readBlockNumber(); err == nil {
		nextBlockNumber = loggedBlockNumber + 1
	}

	for {
		select {
		case <-s.quit:
			s.waitQuit.Done()
			break
		default:
			nextBlock := s.BlockChain().GetBlockByNumber(nextBlockNumber)
			if nextBlock != nil {
				if err := s.reportBlockMsg(nextBlock); err == nil {
					writeStatsLog(nextBlockNumber)
					nextBlockNumber = nextBlockNumber + 1
				} else {
					//
					panic(err)
				}
			} else {
				time.Sleep(time.Microsecond * 50)
			}
		}
	}
}

func (s *PlatonStatsService) reportBlockMsg(block *types.Block) error {
	var genesisData *common.GenesisData
	var receipts []*types.Receipt
	var exeBlockData *common.ExeBlockData
	var statData *common.StatData

	var err error
	if block.NumberU64() == 0 {
		if genesisData = statsdb.Instance().ReadGenesisData(); genesisData == nil {
			log.Error("cannot read genesis data", "err", err)
			return errors.New("cannot read genesis data")
		}
		exeBlockData = statsdb.Instance().ReadExeBlockData(block.Number())
		statData = statsdb.Instance().ReadStatData(block.Number())
	} else {
		receipts = s.BlockChain().GetReceiptsByHash(block.Hash())
		exeBlockData = statsdb.Instance().ReadExeBlockData(block.Number())
		statData = statsdb.Instance().ReadStatData(block.Number())
	}

	brief := collectBrief(block)
	contractList := s.filterDistinctContract(block.NumberU64(), block.Transactions())

	blockJsonMapping, err := jsonBlock(block)
	if err != nil {
		log.Error("marshal block to json string error")
		return err
	}
	statsBlockExt := &StatsBlockExt{
		BlockType:    brief.BlockType,
		Epoch:        brief.Epoch,
		ConsensusNo:  brief.ConsensusNo,
		NodeID:       brief.NodeID,
		NodeAddress:  brief.NodeAddress,
		Block:        blockJsonMapping,
		Receipts:     receipts,
		ExeBlockData: exeBlockData,
		GenesisData:  genesisData,
		ContractList: contractList,
		StatData:     statData,
	}

	jsonBytes, err := json.Marshal(statsBlockExt)
	if err != nil {
		log.Error("marshal platon stats block message to json string error", "blockNumber", block.NumberU64(), "err", err)
		return err
	} else {
		log.Info("marshal platon stats block", "blockNumber", block.NumberU64(), "json", string(jsonBytes))
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.kafkaClient.blockTopic, Partition: 0},
		Key:            []byte(strconv.FormatUint(block.NumberU64(), 10)),
		Value:          []byte(jsonBytes),
		Timestamp:      time.Now(),
	}
	err = s.kafkaClient.producer.Produce(msg, nil)
	if err != nil {
		log.Error("Failed to enqueue the block message", "blockNumber", block.NumberU64(), "err", err)
		return err
	} else {
		log.Info("Success to enqueue the block message", "blockNumber", block.NumberU64())
	}

	//不从statsdb中删除统计需要的过程数据。
	//statsdb.Instance().DeleteExeBlockData(block.Number())
	return nil
}

func collectBrief(block *types.Block) *Brief {
	bn := block.NumberU64()

	brief := new(Brief)
	brief.BlockType = common.GeneralBlock
	brief.Epoch = xutil.CalculateEpoch(bn)
	brief.ConsensusNo = xutil.CalculateRound(bn)
	if bn == 0 {
		brief.BlockType = common.GenesisBlock
		return brief
	} else if xutil.IsBeginOfConsensus(bn) && !xutil.IsBeginOfEpoch(bn) {
		brief.BlockType = common.ConsensusBeginBlock
	} else if xutil.IsElection(bn) {
		brief.BlockType = common.ConsensusElectionBlock
	} else if xutil.IsBeginOfEpoch(bn) {
		brief.BlockType = common.EpochBeginBlock
	} else if xutil.IsEndOfEpoch(bn) {
		brief.BlockType = common.EpochEndBlock
	}
	if nodeID, nodeAddress, err := discover.ExtractNode(block.Header().SealHash(), block.Header().Extra[32:97]); err != nil {
		log.Error("cannot extract node info from block seal hash and signature")
		panic(err)
	} else {
		brief.NodeID = common.NodeID(nodeID)
		brief.NodeAddress = common.Address(nodeAddress)
	}

	return brief
}

func (s *PlatonStatsService) filterDistinctContract(blockNumber uint64, txs types.Transactions) []common.Address {
	contractAddrMap := make(map[common.Address]interface{})
	for _, tx := range txs {
		if tx.To() != nil {
			if _, exist := contractAddrMap[*tx.To()]; !exist && !vm.IsPrecompiledContract(*tx.To()) && !vm.IsPlatONPrecompiledContract(*tx.To()) && len(tx.Data()) > 0 && s.isContract(*tx.To(), blockNumber) {
				contractAddrMap[*tx.To()] = nil
			}
		}
	}
	var contractAddrList []common.Address
	for addr, _ := range contractAddrMap {
		contractAddrList = append(contractAddrList, addr)
	}
	return contractAddrList
}

func readBlockNumber() (uint64, error) {
	if bytes, err := ioutil.ReadFile(statsLogFile); err != nil || len(bytes) == 0 {
		return 0, errors.New("Failed to read PlatON stats service log")
	} else {
		if blockNumber, err := strconv.ParseUint(strings.Trim(string(bytes), "\n"), 10, 64); err != nil {
			log.Warn("Failed to read PlatON stats service log", "error", err)
			return 0, errors.New("Failed to read PlatON stats service log")
		} else {
			log.Info("Success to read PlatON stats service log", "blockNumber", blockNumber)
			return blockNumber, nil
		}
	}
}

func writeStatsLog(blockNumber uint64) {
	if err := common.WriteFile(statsLogFile, []byte(strconv.FormatUint(blockNumber, 10)), statsLogFlag, os.ModePerm); err != nil {
		log.Error("Failed to log stats block number", "file", statsLogFile, "blockNumber", blockNumber, "err", err)
	} else {
		log.Debug("Success to log stats block number", "file", statsLogFile, "blockNumber", blockNumber)
	}
}

func (s *PlatonStatsService) sampleMsgLoop() {
	var sampleEventProducer SampleEventProducer
	sampleEventCh := make(chan SampleEvent, sampleEventChanSize)
	sampleEventSub := sampleEventProducer.SubscribeSampleEvent(sampleEventCh)
	defer sampleEventSub.Unsubscribe()

	for {
		select {
		case sampleEvent := <-sampleEventCh:
			log.Debug("received a sample event", sampleEvent)
		case <-sampleEventSub.Err():
			return
		case <-s.stopSampleMsg:
			return
		}
	}
}

func (s *PlatonStatsService) getCode(to common.Address, blockNumber uint64) ([]byte, error) {
	state, _, err := s.eth.APIBackend.StateAndHeaderByNumber(nil, rpc.BlockNumber(blockNumber))
	if state == nil || err != nil {
		return nil, err
	}
	state.ClearParentReference()
	return state.GetCode(to), state.Error()
}

func (s *PlatonStatsService) isContract(to common.Address, blockNumber uint64) bool {
	if code, err := s.getCode(to, blockNumber); err == nil && len(code) >= 4 {
		return vm.CanUseEVMInterp(code) || vm.CanUseWASMInterp(code)
	}
	return false
}
