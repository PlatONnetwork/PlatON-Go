// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package eth implements the Ethereum protocol.
package eth

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"

	"github.com/PlatONnetwork/PlatON-Go/accounts"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/bloombits"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/eth/filters"
	"github.com/PlatONnetwork/PlatON-Go/eth/gasprice"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/internal/ethapi"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/miner"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
	xplugin "github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var indexMock = map[int][]int{
	1:  []int{2, 3, 4},
	2:  []int{5, 6, 7},
	3:  []int{8, 9, 10},
	4:  []int{11, 12, 13},
	5:  []int{14, 15, 16},
	6:  []int{17, 18, 19},
	7:  []int{},
	8:  []int{20, 21, 22},
	9:  []int{},
	10: []int{23, 24, 25},
}

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// Ethereum implements the Ethereum full node service.
type Ethereum struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the Ethereum

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer
	// modify
	//mpcPool *core.MPCPool
	//vcPool  *core.VCPool

	// DB interfaces
	chainDb ethdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *EthAPIBackend

	miner         *miner.Miner
	gasPrice      *big.Int
	networkID     uint64
	netRPCService *ethapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)
}

func (s *Ethereum) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new Ethereum object (including the
// initialisation of the common Ethereum object)
func New(ctx *node.ServiceContext, config *Config) (*Ethereum, error) {
	// Ensure configuration values are compatible and sane
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run eth.PlatON in light sync mode, use les.LightPlatON")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	if config.MinerGasPrice == nil || config.MinerGasPrice.Cmp(common.Big0) <= 0 || config.MinerGasPrice.Cmp(DefaultConfig.MinerGasPrice) <= 0 {
		log.Warn("Sanitizing invalid miner gas price", "provided", config.MinerGasPrice, "updated", DefaultConfig.MinerGasPrice)
		config.MinerGasPrice = new(big.Int).Set(DefaultConfig.MinerGasPrice)
	}
	// Assemble the Ethereum object
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}

	// set snapshotdb path
	//snapshotdb.SetDBPath(ctx)

	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}

	if nil != chainConfig && nil != chainConfig.Cbft {
		xcom.SetNodeBlockTimeWindow(chainConfig.Cbft.Period / 1000)
		xcom.SetPerRoundBlocks(uint64(chainConfig.Cbft.Amount))
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	eth := &Ethereum{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(ctx, chainConfig, config.MinerNotify, config.MinerNoverify, chainDb, &config.CbftConfig, ctx.EventMux),
		shutdownChan:   make(chan bool),
		networkID:      config.NetworkId,
		gasPrice:       config.MinerGasPrice,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
	}

	log.Info("Initialising Ethereum protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != config.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d).\n", bcVersion, config.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, config.BlockChainVersion)
	}
	var (
		vmConfig = vm.Config{
			EnablePreimageRecording: config.EnablePreimageRecording,
			EWASMInterpreter:        config.EWASMInterpreter,
			EVMInterpreter:          config.EVMInterpreter,
			ConsoleOutput:           config.Debug,
		}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout,
			BodyCacheLimit: config.BodyCacheLimit, BlockCacheLimit: config.BlockCacheLimit,
			MaxFutureBlocks: config.MaxFutureBlocks, BadBlockLimit: config.BadBlockLimit,
			TriesInMemory: config.TriesInMemory, DefaultTxsCacheSize: config.DefaultTxsCacheSize,
			DefaultBroadcastInterval: config.DefaultBroadcastInterval,
		}

		minningConfig = &core.MiningConfig{MiningLogAtDepth: config.MiningLogAtDepth, TxChanSize: config.TxChanSize,
			ChainHeadChanSize: config.ChainHeadChanSize, ChainSideChanSize: config.ChainSideChanSize,
			ResultQueueSize: config.ResultQueueSize, ResubmitAdjustChanSize: config.ResubmitAdjustChanSize,
			MinRecommitInterval: config.MinRecommitInterval, MaxRecommitInterval: config.MaxRecommitInterval,
			IntervalAdjustRatio: config.IntervalAdjustRatio, IntervalAdjustBias: config.IntervalAdjustBias,
			StaleThreshold: config.StaleThreshold, DefaultCommitRatio: config.DefaultCommitRatio,
		}
	)

	eth.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, eth.chainConfig, eth.engine, vmConfig, eth.shouldPreserve)
	if err != nil {
		return nil, err
	}

	blockChainCache := core.NewBlockChainCache(eth.blockchain)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		eth.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	eth.bloomIndexer.Start(eth.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	//eth.txPool = core.NewTxPool(config.TxPool, eth.chainConfig, eth.blockchain)
	eth.txPool = core.NewTxPool(config.TxPool, eth.chainConfig, blockChainCache)

	// mpcPool deal with mpc transactions
	// modify By J
	//if config.MPCPool.Journal != "" {
	//	config.MPCPool.Journal = ctx.ResolvePath(config.MPCPool.Journal)
	//} else {
	//	config.MPCPool.Journal = ctx.ResolvePath(core.DefaultMPCPoolConfig.Journal)
	//}
	//if config.MPCPool.Rejournal == 0 {
	//	config.MPCPool.Rejournal = core.DefaultMPCPoolConfig.Rejournal
	//}
	//if config.MPCPool.Lifetime == 0 {
	//	config.MPCPool.Lifetime = core.DefaultMPCPoolConfig.Lifetime
	//}
	//eth.mpcPool = core.NewMPCPool(config.MPCPool, eth.chainConfig, eth.blockchain)
	//eth.vcPool = core.NewVCPool(config.VCPool, eth.chainConfig, eth.blockchain)

	// modify by platon remove consensusCache
	//var consensusCache *cbft.Cache = cbft.NewCache(eth.blockchain)
	eth.miner = miner.New(eth, eth.chainConfig, minningConfig, eth.EventMux(), eth.engine, config.MinerRecommit,
		config.MinerGasFloor, config.MinerGasCeil, eth.isLocalBlock, blockChainCache)
	//extra data for each block will be set by worker.go
	//eth.miner.SetExtra(makeExtraData(eth.blockchain, config.MinerExtraData))

	reactor := core.NewBlockChainReactor(config.CbftConfig.NodePriKey, eth.EventMux())

	if engine, ok := eth.engine.(consensus.Bft); ok {

		var agency consensus.Agency
		// validatorMode:
		// - static (default)
		// - inner (via inner contract)eth/handler.go
		// - ppos
		log.Debug("Validator mode", "mode", chainConfig.Cbft.ValidatorMode)
		if chainConfig.Cbft.ValidatorMode == "" || chainConfig.Cbft.ValidatorMode == common.STATIC_VALIDATOR_MODE {
			agency = validator.NewStaticAgency(chainConfig.Cbft.InitialNodes)
			reactor.Start(common.STATIC_VALIDATOR_MODE)
		} else if chainConfig.Cbft.ValidatorMode == common.INNER_VALIDATOR_MODE {
			blocksPerNode := int(chainConfig.Cbft.Amount)
			offset := blocksPerNode * 2
			agency = validator.NewInnerAgency(chainConfig.Cbft.InitialNodes, eth.blockchain, blocksPerNode, offset)
			reactor.Start(common.INNER_VALIDATOR_MODE)
		} else if chainConfig.Cbft.ValidatorMode == common.PPOS_VALIDATOR_MODE {
			reactor.Start(common.PPOS_VALIDATOR_MODE)
			reactor.SetVRF_handler(xcom.NewVrfHandler(eth.blockchain.Genesis().Nonce()))
			handlePlugin(reactor)
			agency = reactor
		}

		if err := engine.Start(eth.blockchain, blockChainCache, eth.txPool, agency); err != nil {
			log.Error("Init cbft consensus engine fail", "error", err)
			return nil, errors.New("Failed to init cbft consensus engine")
		}
	}

	if eth.protocolManager, err = NewProtocolManager(eth.chainConfig, config.SyncMode, config.NetworkId, eth.eventMux, eth.txPool, eth.engine, eth.blockchain, chainDb); err != nil {
		return nil, err
	}

	eth.APIBackend = &EthAPIBackend{eth, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.MinerGasPrice
	}
	eth.APIBackend.gpo = gasprice.NewOracle(eth.APIBackend, gpoParams)

	return eth, nil
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (ethdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*ethdb.LDBDatabase); ok {
		db.Meter("eth/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an Ethereum service
func CreateConsensusEngine(ctx *node.ServiceContext, chainConfig *params.ChainConfig, notify []string, noverify bool, db ethdb.Database,
	cbftConfig *ctypes.OptionsConfig, eventMux *event.TypeMux) consensus.Engine {
	// If proof-of-authority is requested, set it up
	engine := cbft.New(chainConfig.Cbft, cbftConfig, eventMux, ctx)
	if engine == nil {
		panic("create consensus engine fail")
	}
	return engine
}

// APIs return the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Ethereum) APIs() []rpc.API {
	apis := ethapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "platon",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "platon",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *Ethereum) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

// isLocalBlock checks whether the specified block is mined
// by local miner accounts.
//
// We regard two types of accounts as local miner account: etherbase
// and accounts specified via `txpool.locals` flag.
func (s *Ethereum) isLocalBlock(block *types.Block) bool {
	author, err := s.engine.Author(block.Header())
	if err != nil {
		log.Warn("Failed to retrieve block author", "number", block.NumberU64(), "hash", block.Hash(), "err", err)
		return false
	}
	// Check whether the given address is etherbase.
	s.lock.RLock()
	etherbase := common.Address{}
	s.lock.RUnlock()
	if author == etherbase {
		return true
	}
	// Check whether the given address is specified by `txpool.local`
	// CLI flag.
	for _, account := range s.config.TxPool.Locals {
		if account == author {
			return true
		}
	}
	return false
}

// shouldPreserve checks whether we should preserve the given block
// during the chain reorg depending on whether the author of block
// is a local account.
func (s *Ethereum) shouldPreserve(block *types.Block) bool {
	// The reason we need to disable the self-reorg preserving for clique
	// is it can be probable to introduce a deadlock.
	//
	// e.g. If there are 7 available signers
	//
	// r1   A
	// r2     B
	// r3       C
	// r4         D
	// r5   A      [X] F G
	// r6    [X]
	//
	// In the round5, the inturn signer E is offline, so the worst case
	// is A, F and G sign the block of round5 and reject the block of opponents
	// and in the round6, the last available signer B is offline, the whole
	// network is stuck.
	return s.isLocalBlock(block)
}

//start mining
func (s *Ethereum) StartMining() error {
	// If the miner was not running, initialize it
	if !s.IsMining() {
		// Propagate the initial price point to the transaction pool
		s.lock.RLock()
		price := s.gasPrice
		s.lock.RUnlock()
		s.txPool.SetGasPrice(price)

		// If mining is started, we can disable the transaction rejection mechanism
		// introduced to speed sync times.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)

		go s.miner.Start()
	}
	return nil
}

// StopMining terminates the miner, both at the consensus engine level as well as
// at the block creation level.
func (s *Ethereum) StopMining() {
	s.miner.Stop()
}

func (s *Ethereum) IsMining() bool      { return s.miner.Mining() }
func (s *Ethereum) Miner() *miner.Miner { return s.miner }

func (s *Ethereum) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Ethereum) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Ethereum) TxPool() *core.TxPool               { return s.txPool }
func (s *Ethereum) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Ethereum) Engine() consensus.Engine           { return s.engine }
func (s *Ethereum) ChainDb() ethdb.Database            { return s.chainDb }
func (s *Ethereum) IsListening() bool                  { return true } // Always listening
func (s *Ethereum) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Ethereum) NetVersion() uint64                 { return s.networkID }
func (s *Ethereum) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Ethereum) Protocols() []p2p.Protocol {
	protocols := make([]p2p.Protocol, 0)
	protocols = append(protocols, s.protocolManager.SubProtocols...)
	protocols = append(protocols, s.engine.Protocols()...)

	if s.lesServer == nil {
		return protocols
	}
	protocols = append(protocols, s.lesServer.Protocols()...)
	return protocols
}

// Start implements node.Service, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *Ethereum) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers(params.BloomBitsBlocks)

	// Start the RPC service
	s.netRPCService = ethapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= srvr.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)

	if cbftEngine, ok := s.engine.(consensus.Bft); ok {
		core.GetReactorInstance().SetPrivateKey(srvr.Config.PrivateKey)
		if flag := cbftEngine.IsConsensusNode(); flag {
			// self: s.chainConfig.Cbft.NodeID
			// list: s.chainConfig.Cbft.InitialNodes
			// dep: test
			/*ok, idxs := needAdd(s.chainConfig.Cbft.NodeID, s.chainConfig.Cbft.InitialNodes)
			for idx, n := range s.chainConfig.Cbft.InitialNodes {
				if idxs == nil {
					break
				}
				for _, i := range idxs {
					if ok && i == (idx+1) {
						srvr.AddConsensusPeer(discover.NewNode(n.ID, n.IP, n.UDP, n.TCP))
						break
					}
				}
			}*/
			for _, n := range s.chainConfig.Cbft.InitialNodes {
				srvr.AddConsensusPeer(discover.NewNode(n.Node.ID, n.Node.IP, n.Node.UDP, n.Node.TCP))
			}
		}
		s.StartMining()
	}
	srvr.StartWatching(s.eventMux)

	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	return nil
}

// mock
func needAdd(self discover.NodeID, nodes []discover.Node) (bool, []int) {
	selfIndex := -1
	for idx, n := range nodes {
		if n.ID.TerminalString() == self.TerminalString() {
			selfIndex = idx
			break
		}
	}
	if selfIndex == -1 {
		return false, nil
	}
	selfIndex++
	return true, indexMock[selfIndex]
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *Ethereum) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.engine.Close()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)
	core.GetReactorInstance().Close()
	return nil
}

// RegisterPlugin one by one
func handlePlugin(reactor *core.BlockChainReactor) {
	reactor.RegisterPlugin(xcom.SlashingRule, xplugin.SlashInstance())
	xplugin.SlashInstance().SetDecodeEvidenceFun(evidence.NewEvidences)
	reactor.RegisterPlugin(xcom.StakingRule, xplugin.StakingInstance())
	reactor.RegisterPlugin(xcom.RestrictingRule, xplugin.RestrictingInstance())
	reactor.RegisterPlugin(xcom.RewardRule, xplugin.RewardMgrInstance())
	reactor.RegisterPlugin(xcom.GovernanceRule, xplugin.GovPluginInstance())

	reactor.SetPluginEventMux()

	// set rule order
	reactor.SetBeginRule([]int{xcom.SlashingRule})
	reactor.SetEndRule([]int{xcom.RestrictingRule, xcom.RewardRule, xcom.GovernanceRule, xcom.StakingRule})

}
