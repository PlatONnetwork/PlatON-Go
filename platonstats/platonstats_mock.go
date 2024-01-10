package platonstats

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/PlatONnetwork/PlatON-Go/core/statsdb"

	"github.com/PlatONnetwork/PlatON-Go/rpc"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
)

type StatsServer interface {
	reportBlockMsg(block *types.Block) error
	scanGenesis(genesisBlock *types.Block) (*common.GenesisData, error)
}
type MockPlatonStatsService struct {
	server *p2p.Server // Peer-to-peer server to retrieve networking infos

	kafkaUrl        string
	kafkaBlockTopic string

	//eth      *eth.Ethereum // Full Ethereum service if monitoring a full node
	blockChain *core.BlockChain
	chainDb    ethdb.Database

	kafkaClient *ConfluentKafkaClient

	stopSampleMsg chan struct{}
	stopBlockMsg  chan struct{}
	stopOnce      sync.Once
}

func NewMockPlatonStatsService() *MockPlatonStatsService {
	statsService := new(MockPlatonStatsService)

	//statsService.chainDb = ethdb.NewMemDatabase()
	genesis := new(core.Genesis).MustCommit(statsService.chainDb)

	bft := consensus.NewFaker()
	bft.InsertChain(genesis)

	chain := makeBlockChain(bft.CurrentBlock(), 121, consensus.NewFaker(), statsService.chainDb, 0)
	statsService.blockChain = chain

	return statsService
}

// makeBlockChain creates a deterministic chain of blocks rooted at parent.
func makeBlockChain(parent *types.Block, n int, engine consensus.Engine, db ethdb.Database, seed int) *core.BlockChain {
	blockChain := core.GenerateBlockChain2(params.TestChainConfig, parent, engine, db, n, func(i int, b *core.BlockGen) {
		b.SetCoinbase(common.Address{0: byte(seed), 19: byte(i)})
	})
	return blockChain
}

func (s *MockPlatonStatsService) reportBlockMsg(block *types.Block) error {
	var genesisData *common.GenesisData
	var receipts []*types.Receipt
	var exeBlockData *common.ExeBlockData

	var err error
	if block.NumberU64() == 0 {
		if genesisData, err = s.scanGenesis(block); err != nil {
			log.Error("cannot read genesis block", err)
			return err
		}
	} else {
		receipts = s.BlockChain().GetReceiptsByHash(block.Hash())

		//tx1 := types.NewTransaction(1, common.Address{0x01, 0x02}, big.NewInt(1), 30000, big.NewInt(1), nil)
		//tx2 := types.NewTransaction(2, common.Address{0x01, 0x02}, big.NewInt(2), 30000, big.NewInt(2), nil)

		//block.SetTransactions(types.Transactions{tx1, tx2})

		exeBlockData = statsdb.Instance().ReadExeBlockData(block.Number())

	}
	brief := collectBrief(block)

	blockJsonMapping, err := jsonBlock(block)
	if err != nil {
		log.Error("marshal block to json string error")
		return err
	}

	statsBlockExt := &StatsBlockExt{
		BlockType: brief.BlockType,
		Epoch:     brief.Epoch,
		//Block:        convertBlock(block),
		Block:        blockJsonMapping,
		Receipts:     receipts,
		ExeBlockData: exeBlockData,
		GenesisData:  genesisData,
	}

	json, err := json.Marshal(statsBlockExt)
	if err != nil {
		log.Error("marshal platon stats block message to json string error")
		return err
	} else {
		log.Debug("marshal platon stats block", "json", string(json))
	}

	// send message
	var blockTopic string
	if len(s.kafkaBlockTopic) == 0 {
		blockTopic = defaultKafkaBlockTopic
	} else {
		blockTopic = s.kafkaBlockTopic
	}
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &blockTopic, Partition: 0},
		Key:            []byte(strconv.FormatUint(block.NumberU64(), 10)),
		Value:          []byte(json),
		Timestamp:      time.Now(),
	}

	err = s.kafkaClient.producer.Produce(msg, nil)
	if err != nil {
		log.Error("Failed to enqueue the block message", "blockNumber", block.NumberU64(), "err", err)
		return err
	} else {
		log.Info("Success to enqueue the block message", "blockNumber", block.NumberU64())
	}

	//statsdb.Instance().DeleteExeBlockData(block.Number())
	return nil
}
func (s *MockPlatonStatsService) sampleMsgLoop() {
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

func (s *MockPlatonStatsService) scanGenesis(genesisBlock *types.Block) (*common.GenesisData, error) {
	genesisData := &common.GenesisData{}
	/*hash := rawdb.ReadCanonicalHash(s.eth.ChainDb(), 0)
	println("genesis block hash:", hash.String())
	block := rawdb.ReadBlock(s.eth.ChainDb(), hash, 0)
	if block == nil {
		return nil, fmt.Errorf("cannot read genesis block")
	}
	*/
	root := genesisBlock.Root()
	tr, err := trie.NewSecure(root, trie.NewDatabase(s.ChainDb()))
	if err != nil {
		return nil, err
	}

	iter := tr.NodeIterator(nil)
	for iter.Next(true) {
		if iter.Leaf() {
			var obj state.DumpAccount
			err := rlp.DecodeBytes(iter.LeafBlob(), &obj)
			if err != nil {
				return nil, fmt.Errorf("parse account error:%s", err.Error())
			}
			key := iter.LeafKey()
			address := common.BytesToAddress(key)
			balance, _ := new(big.Int).SetString(obj.Balance, 10)
			genesisData.AddAllocItem(address, balance)

			log.Debug("alloc account", "address", address, "balance", balance)
		}
	}
	return genesisData, nil
}

func (s *MockPlatonStatsService) BlockChain() *core.BlockChain {
	return s.blockChain
}

func (s *MockPlatonStatsService) ChainDb() ethdb.Database {
	return s.chainDb
}

func (s *MockPlatonStatsService) Protocols() []p2p.Protocol { return nil }

// APIs implements node.Service, returning the RPC API endpoints provided by the
// stats service (nil as it doesn't provide any user callable APIs).
func (s *MockPlatonStatsService) APIs() []rpc.API { return nil }

// Start implements node.Service, starting up the monitoring and reporting daemon.
func (s *MockPlatonStatsService) Start(server *p2p.Server) error {
	s.server = server

	s.kafkaClient = NewConfluentKafkaClient(s.kafkaUrl, s.kafkaBlockTopic)

	go s.blockMsgLoop()
	go s.sampleMsgLoop()

	log.Info("PlatON stats daemon started")
	return nil
}

// Stop implements node.Service, terminating the monitoring and reporting daemon.
func (s *MockPlatonStatsService) Stop() error {
	s.stopOnce.Do(func() {
		close(s.stopSampleMsg)
		close(s.stopBlockMsg)
		if s.kafkaClient != nil {
			s.kafkaClient.Close()
		}
	})

	log.Info("PlatON stats daemon stopped")
	return nil
}

func (s *MockPlatonStatsService) blockMsgLoop() {
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
				log.Debug("Success to write stats service log", "blockNumber", nextBlock.NumberU64())
				nextBlockNumber = nextBlockNumber + 1
			} else {
				log.Error("Failed to report block message", "blockNumber", nextBlock.NumberU64())
			}
		} else {
			time.Sleep(time.Microsecond * 100)
		}
	}
}
