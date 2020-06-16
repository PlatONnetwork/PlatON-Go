package platonstats

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

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
	"github.com/Shopify/sarama"
)

type StatsServer interface {
	reportBlockMsg(block *types.Block) error
	scanGenesis(genesisBlock *types.Block) (*common.GenesisData, error)
}
type MockPlatonStatsService struct {
	server *p2p.Server // Peer-to-peer server to retrieve networking infos

	kafkaUrl string
	//eth      *eth.Ethereum // Full Ethereum service if monitoring a full node
	blockChain    *core.BlockChain
	chainDb       ethdb.Database
	blockProducer sarama.SyncProducer
	msgProducer   sarama.AsyncProducer

	stopSampleMsg chan struct{}
	stopBlockMsg  chan struct{}
	stopOnce      sync.Once
}

func NewMockPlatonStatsService() *MockPlatonStatsService {
	statsService := new(MockPlatonStatsService)

	statsService.chainDb = ethdb.NewMemDatabase()
	genesis := new(core.Genesis).MustCommit(statsService.chainDb)

	bft := consensus.NewFaker()
	bft.InsertChain(genesis)

	chain := makeBlockChain(bft.CurrentBlock(), 121, consensus.NewFaker(), statsService.chainDb, 0)
	statsService.blockChain = chain

	common.PlatONStatsServiceRunning = true
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
		exeBlockData = statsdb.Instance().ReadExeBlockData(block.Number())
	}

	statsBlockExt := &StatsBlockExt{
		Brief:        collectBrief(block),
		Header:       block.Header(),
		Txs:          convertTxs(block.Transactions()),
		Receipts:     receipts,
		ExeBlockData: exeBlockData,
		GenesisData:  genesisData,
	}

	blockExtBytes, err := rlp.EncodeToBytes(statsBlockExt)
	if err != nil {
		log.Error("encode platon stats block ext message error")
		return err
	}
	data := common.Bytes2Hex(blockExtBytes)
	log.Info("Prepare to send PlatON block ext message to Kafka", "blockNumber", block.NumberU64(), "hex", data)

	// send message
	msg := &sarama.ProducerMessage{
		Topic:     kafkaBlockTopic,
		Partition: 0,
		Key:       sarama.StringEncoder(strconv.FormatUint(block.NumberU64(), 10)),
		Value:     sarama.StringEncoder(common.Bytes2Hex(blockExtBytes)),
		Timestamp: time.Now(),
	}

	partition, offset, err := s.blockProducer.SendMessage(msg)

	if err != nil {
		log.Error("send block message error.", "blockNumber=", block.NumberU64(), "error", err)
	} else {
		log.Info("send block message success.", "blockNumber=", block.NumberU64(), "partition", partition, "offset", offset)
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
	tr, err := trie.NewSecure(root, trie.NewDatabase(s.ChainDb()), 0)
	if err != nil {
		return nil, err
	}

	iter := tr.NodeIterator(nil)
	for iter.Next(true) {
		if iter.Leaf() {
			var obj state.Account
			err := rlp.DecodeBytes(iter.LeafBlob(), &obj)
			if err != nil {
				return nil, fmt.Errorf("parse account error:%s", err.Error())
			}
			key := iter.LeafKey()
			address := common.BytesToAddress(key)
			balance := obj.Balance.Uint64()
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
	urls := strings.Split(s.kafkaUrl, ",")

	if msgProducer, err := sarama.NewAsyncProducer(urls, msgProducerConfig()); err != nil {
		return err
	} else {
		s.msgProducer = msgProducer
	}

	if blockProducer, err := sarama.NewSyncProducer(urls, blockProducerConfig()); err != nil {
		return err
	} else {
		s.blockProducer = blockProducer
	}

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

		if s.msgProducer != nil {
			s.msgProducer.AsyncClose()
		}
		if s.blockProducer != nil {
			s.blockProducer.Close()
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
				if err := writeBlockNumber(nextBlockNumber); err == nil {
					log.Debug("Success to write stats service log", "blockNumber", nextBlock.NumberU64())
					nextBlockNumber = nextBlockNumber + 1
				} else {
					log.Error("Failed to write stats service log", "blockNumber", nextBlock.NumberU64())
				}
			} else {
				log.Error("Failed to report block message", "blockNumber", nextBlock.NumberU64())
			}
		} else {
			time.Sleep(time.Microsecond * 100)
		}
	}
}
