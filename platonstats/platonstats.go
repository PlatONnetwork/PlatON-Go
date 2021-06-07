package platonstats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
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

	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

var (
	statsLogFile = "./platonstats.log"
	statsLogFlag = os.O_RDWR | os.O_CREATE | os.O_TRUNC

	checkingErrFile = "./checkingerr.log"
	checkingErrFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND
)

type platonStats interface {
	SubscribeSampleEvent(ch chan<- SampleEvent) event.Subscription
}

type blockdata struct {
	Number       uint64             `json:"number"    gencodec:"required"`
	Hash         common.Hash        `json:"hash"    gencodec:"required"`
	ParentHash   common.Hash        `json:"parentHash"    gencodec:"required"`
	LogsBloom    types.Bloom        `json:"logsBloom"    gencodec:"required"`
	StateRoot    common.Hash        `json:"stateRoot"    gencodec:"required"`
	ReceiptsRoot common.Hash        `json:"receiptsRoot"    gencodec:"required"`
	TxHash       common.Hash        `json:"transactionsRoot" gencodec:"required"`
	Miner        common.Address     `json:"miner"    gencodec:"required"`
	ExtraData    ExtraData          `json:"extraData"    gencodec:"required"`
	GasLimit     uint64             `json:"gasLimit"    gencodec:"required"`
	GasUsed      uint64             `json:"gasUsed"    gencodec:"required"`
	Timestamp    uint64             `json:"timestamp"    gencodec:"required"`
	Transactions types.Transactions `json:"transactions"    gencodec:"required"`
	Nonce        Nonce              `json:"nonce"    gencodec:"required"`
}

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

func convertBlock(block *types.Block) *blockdata {
	blk := new(blockdata)
	blk.Number = block.NumberU64()
	blk.Hash = block.Hash()
	blk.ParentHash = block.ParentHash()
	blk.LogsBloom = block.Bloom()
	blk.StateRoot = block.Root()
	blk.ReceiptsRoot = block.ReceiptHash()
	blk.TxHash = block.TxHash()
	blk.Miner = block.Coinbase()
	blk.ExtraData = ExtraData(block.Extra())
	blk.GasLimit = block.GasLimit()
	blk.GasUsed = block.GasUsed()
	blk.Timestamp = block.Time().Uint64()
	blk.Transactions = block.Transactions()
	blk.Nonce = block.Nonce()
	return blk
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
	ContractList []*common.Address      `json:"ContractList,omitempty"`
}

type PlatonStatsService struct {
	server                    *p2p.Server //Peer-to-peer server to retrieve networking infos
	kafkaUrl                  string
	kafkaBlockTopic           string        //统计数据消息Topic
	kafkaAccountCheckingTopic string        //对账请求消息Topic
	kafkaAccountCheckingGroup string        //对账请求消息Group
	eth                       *eth.Ethereum // Full Ethereum service if monitoring a full node
	datadir                   string
	kafkaClient               *ConfluentKafkaClient
	/*blockProducer             sarama.SyncProducer
	msgProducer               sarama.AsyncProducer
	checkingConsumer          *cluster.Consumer*/
	stopSampleMsg chan struct{}
	stopBlockMsg  chan struct{}
	stopOnce      sync.Once
}

var (
	//platonStatsServiceOnce sync.Once
	platonStatsService *PlatonStatsService
)

func New(kafkaUrl, kafkaBlockTopic, kafkaAccountCheckingTopic, kafkaAccountCheckingGroup string, ethServ *eth.Ethereum, datadir string) (*PlatonStatsService, error) {
	log.Info("new PlatON stats service", "kafkaUrl", kafkaUrl, "kafkaBlockTopic", kafkaBlockTopic, "kafkaAccountCheckingTopic", kafkaAccountCheckingTopic, "kafkaAccountCheckingGroup", kafkaAccountCheckingGroup)

	platonStatsService = &PlatonStatsService{
		kafkaUrl:                  kafkaUrl,
		kafkaBlockTopic:           kafkaBlockTopic,
		kafkaAccountCheckingTopic: kafkaAccountCheckingTopic,
		kafkaAccountCheckingGroup: kafkaAccountCheckingGroup,
		eth:                       ethServ,
		datadir:                   datadir,
	}
	if len(datadir) > 0 {
		statsLogFile = filepath.Join(datadir, statsLogFile)
		checkingErrFile = filepath.Join(datadir, checkingErrFile)
	}
	log.Debug("PlatON stats log file", "datadir", datadir)
	log.Debug("PlatON stats log file", "statsLogFile", statsLogFile)
	log.Debug("PlatON stats log file", "checkingErrFile", checkingErrFile)
	return platonStatsService, nil
}

func GetPlatonStatsService() *PlatonStatsService {
	return platonStatsService
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
	//urls := []string{s.kafkaUrl}

	//s.kafkaClient = NewKafkaClient(s.kafkaUrl, s.kafkaBlockTopic, s.kafkaAccountCheckingTopic, s.kafkaAccountCheckingGroup)
	s.kafkaClient = NewConfluentKafkaClient(s.kafkaUrl, s.kafkaBlockTopic, s.kafkaAccountCheckingTopic, s.kafkaAccountCheckingGroup)

	/*if msgProducer, err := sarama.NewAsyncProducer(urls, msgProducerConfig()); err != nil {
		log.Error("Failed to init msg Kafka async producer....", "err", err)
		return err
	} else {
		log.Info("Success to init msg Kafka async producer....")
		s.msgProducer = msgProducer
	}

	if blockProducer, err := sarama.NewSyncProducer(urls, blockProducerConfig()); err != nil {
		log.Error("Failed to init msg Kafka sync producer....", "err", err)
		return err
	} else {
		log.Info("Success to init msg Kafka sync producer....")
		s.blockProducer = blockProducer
	}
	*/
	go s.blockMsgLoop()
	//go s.sampleMsgLoop()

	go s.accountCheckingLoop()
	log.Info("PlatON stats daemon started")
	return nil
}

// Stop implements node.Service, terminating the monitoring and reporting daemon.
func (s *PlatonStatsService) Stop() error {
	s.stopOnce.Do(func() {
		close(s.stopSampleMsg)
		//close(s.stopBlockMsg)
		if s.kafkaClient != nil {
			s.kafkaClient.Close()
		}
	})

	log.Info("PlatON stats daemon stopped")
	return nil
}

//todo: 服务如何退出？整个Node如何停止？
func (s *PlatonStatsService) blockMsgLoop() {
	var nextBlockNumber uint64
	nextBlockNumber = 0

	if loggedBlockNumber, err := readBlockNumber(); err == nil {
		nextBlockNumber = loggedBlockNumber + 1
	}

	for {
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

func (s *PlatonStatsService) reportBlockMsg(block *types.Block) error {
	var genesisData *common.GenesisData
	var receipts []*types.Receipt
	var exeBlockData *common.ExeBlockData

	var err error
	if block.NumberU64() == 0 {
		if genesisData = statsdb.Instance().ReadGenesisData(); genesisData == nil {
			log.Error("cannot read genesis data", "err", err)
			return errors.New("cannot read genesis data")
		}
	} else {
		receipts = s.BlockChain().GetReceiptsByHash(block.Hash())
		exeBlockData = statsdb.Instance().ReadExeBlockData(block.Number())
	}

	brief := collectBrief(block)
	contractList := s.filterContract(block.NumberU64(), block.Transactions())

	blockJsonMapping, err := jsonBlock(block)
	if err != nil {
		log.Error("marshal block to json string error")
		return err
	}
	statsBlockExt := &StatsBlockExt{
		BlockType:   brief.BlockType,
		Epoch:       brief.Epoch,
		ConsensusNo: brief.ConsensusNo,
		NodeID:      brief.NodeID,
		NodeAddress: brief.NodeAddress,
		//Block:        convertBlock(block),
		Block:        blockJsonMapping,
		Receipts:     receipts,
		ExeBlockData: exeBlockData,
		GenesisData:  genesisData,
		ContractList: contractList,
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
		/*
			} else if yes, err := xcom.IsYearEnd(common.ZeroHash, bn); err != nil {
				panic(err)
			} else if yes {
				brief.BlockType = common.EndOfYear
				} else if xutil.IsElection(bn) {
				brief.BlockType = common.ConsensusElectionBlock
					} else if xutil.IsBeginOfConsensus(bn) {
					brief.BlockType = common.ConsensusBeginBlock
		*/
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

func (s *PlatonStatsService) filterContract(blockNumber uint64, txs types.Transactions) []*common.Address {
	contractTxList := make([]*common.Address, 0)
	for _, tx := range txs {
		if s.isContract(*tx.To(), blockNumber) {
			contractTxList = append(contractTxList, tx.To())
		}
	}
	return contractTxList
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

func (s *PlatonStatsService) accountCheckingLoop() {
	for {
		msg, err := s.kafkaClient.consumer.ReadMessage(-1)
		if err == nil {
			key := string(msg.Key)
			value := string(msg.Value)
			log.Debug("received account-checking message by group consumer", "key", key, "value", value)

			if len(key) > 0 {
				checkingNumber, err := strconv.ParseUint(key, 10, 64)
				if err != nil {
					log.Error("Failed to parse block number", "key", key, "err", err)
					panic(err)
				}

				for {
					currentNumber := s.eth.BlockChain().CurrentBlock().NumberU64()
					log.Debug("current block number of block chain", "blockNumber", currentNumber)
					if currentNumber >= checkingNumber {
						break
					} else {
						time.Sleep(1 * time.Second)
					}
				}
			}

			err := s.accountChecking(key, msg.Value)
			if err != nil {
				log.Crit("Failed to check account balance", "err", err)
				//panic(err)
			} else {
				log.Debug("Success to check account balance", "key", key)
			}

		} else {
			// The client will automatically try to recover from all errors.
			log.Error("Consumer error", "msg", msg, "err", err)
		}
	}
}

var (
	ErrKey             = errors.New("account checking: cannot convert key to block number")
	ErrValue           = errors.New("account checking: cannot unmarshal value to message struct")
	ErrKeyValue        = errors.New("account checking: key is not matched to value")
	ErrChain           = errors.New("account checking: failed to get account chain balance")
	ErrAccountNotFound = errors.New("account checking: failed to find account in current block chain")
	ErrAccountChecking = errors.New("account checking: Account chain and tracking balances are not equal")
)

func (s *PlatonStatsService) accountChecking(key string, value []byte) error {
	keyNumber, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		log.Error("Failed to convert key to block number", "key", key, "err", err)
		return ErrKey
	}

	var message AccountCheckingMessage
	if len(value) > 0 {
		err := json.Unmarshal(value, &message)
		if err != nil {
			log.Error("Failed to unmarshal value to accountCheckingMessage", "value", string(value), "err", err)
			return ErrValue
		}
	}

	accountCheckingError := false
	if message.BlockNumber == uint64(keyNumber) {
		for _, item := range message.AccountList {
			bech32 := item.Addr.Bech32()
			chainBalance, err := getBalance(s.eth.APIBackend, item.Addr, rpc.BlockNumber(keyNumber))
			if err != nil {
				log.Error("Failed to get account chain balance", "blockNumber", keyNumber, "address", bech32, "err", err)
				return ErrChain
			}
			//the current stateDB's block number is higher than checking request.
			//so, the account must exists in current stateDB.
			if chainBalance == nil {
				log.Error("Failed to find account in current block chain", "blockNumber", keyNumber, "address", bech32)
				return ErrAccountNotFound
			}

			log.Debug("account checking", "blockNumber", keyNumber, "address", bech32, "chainBalance", chainBalance, "trackingBalance", item.Balance)
			if item.Balance.Cmp(chainBalance) != 0 {
				writeCheckingErr(bech32, message.BlockNumber, chainBalance, item.Balance)
				accountCheckingError = true
			}
		}
	} else {
		log.Error("Block number of Kafka message is invalid", "key", keyNumber, "blockNumber", message.BlockNumber)
		return ErrKeyValue
	}

	if accountCheckingError {
		return ErrAccountChecking
	} else {
		return nil
	}
}

func getBalance(backend *eth.EthAPIBackend, address common.Address, blockNr rpc.BlockNumber) (*big.Int, error) {
	state, _, err := backend.StateAndHeaderByNumber(nil, blockNr)
	if state == nil || err != nil {
		return nil, err
	}
	state.ClearParentReference()
	return state.GetBalance(address), state.Error()
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
	if code, err := s.getCode(to, blockNumber); err == nil && len(code) > 0 {
		return vm.CanUseEVMInterp(code) || vm.CanUseWASMInterp(code)
	}
	return false
}

func writeCheckingErr(bech32 string, blockNumber uint64, chainBalance, trackingBalance *big.Int) {
	log.Error("Account chain and tracking balances are not equal", "blockNumber", blockNumber, "address", bech32, "chainBalance", chainBalance, "trackingBalance", trackingBalance)
	content := fmt.Sprintf("blockNumber=%d    account=%s    chainBalance=%d    trackingBalance=%d\n", blockNumber, bech32, chainBalance, trackingBalance)
	err := common.WriteFile(checkingErrFile, []byte(content), checkingErrFlag, os.ModePerm)
	if err != nil {
		log.Error("Failed to log account-checking-error", "content", content)
	}
}

type AccountCheckingMessage struct {
	BlockNumber uint64
	AccountList []*AccountItem
}

type AccountItem struct {
	Addr    common.Address
	Balance *big.Int
}

/*func convertTxs(transactions types.Transactions) []*Tx {
	txs := make([]*Tx, transactions.Len())
	for idx, t := range transactions {
		tx := new(Tx)
		tx.Hash = t.Hash()
		tx.Nonce = tx.Nonce
		tx.From = t.GetFromAddr()
		tx.To = t.To()
		tx.Value = t.Value().Uint64()
		tx.gas = t.Gas()
		tx.GasPrice = t.GasPrice().Uint64()
		tx.Input = t.Data()
		txs[idx] = tx
	}
	return txs
}*/
