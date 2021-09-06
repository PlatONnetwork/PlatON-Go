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
	"os"
	"sync"
	"sync/atomic"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/x/handler"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

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

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// Ethereum implements the Ethereum full node service.
type Ethereum struct {
	config *Config

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

	bloomRequests     chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer      *core.ChainIndexer             // Bloom indexer operating during block imports
	closeBloomHandler chan struct{}

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
	if config.Miner.GasPrice == nil || config.Miner.GasPrice.Cmp(common.Big0) <= 0 {
		log.Warn("Sanitizing invalid miner gas price", "provided", config.Miner.GasPrice, "updated", DefaultConfig.Miner.GasPrice)
		config.Miner.GasPrice = new(big.Int).Set(DefaultConfig.Miner.GasPrice)
	}
	// Assemble the Ethereum object
	chainDb, err := ctx.OpenDatabaseWithFreezer("chaindata", config.DatabaseCache, config.DatabaseHandles, config.DatabaseFreezer, "eth/db/chaindata/")
	if err != nil {
		return nil, err
	}
	snapshotdb.SetDBOptions(config.DatabaseCache, config.DatabaseHandles)

	snapshotBaseDB, err := snapshotdb.Open(ctx.ResolvePath(snapshotdb.DBPath), config.DatabaseCache, config.DatabaseHandles, true)
	if err != nil {
		return nil, err
	}

	height := rawdb.ReadHeaderNumber(chainDb, rawdb.ReadHeadHeaderHash(chainDb))
	log.Debug("read header number from chain db", "height", height)
	if height != nil && *height > 0 {
		//when last  fast syncing fail,we will clean chaindb,wal,snapshotdb
		status, err := snapshotBaseDB.GetBaseDB([]byte(downloader.KeyFastSyncStatus))

		// systemError
		if err != nil && err != snapshotdb.ErrNotFound {
			if err := snapshotBaseDB.Close(); err != nil {
				return nil, err
			}
			return nil, err
		}
		//if find sync status,this means last syncing not finish,should clean all db to reinit
		//if not find sync status,no need init chain
		if err == nil {

			// Just commit the new block if there is no stored genesis block.
			stored := rawdb.ReadCanonicalHash(chainDb, 0)

			log.Info("last fast sync is fail,init  db", "status", common.BytesToUint32(status), "prichain", config.Genesis == nil)
			chainDb.Close()
			if err := snapshotBaseDB.Close(); err != nil {
				return nil, err
			}
			if err := os.RemoveAll(ctx.ResolvePath("chaindata")); err != nil {
				return nil, err
			}

			if err := os.RemoveAll(ctx.ResolvePath(wal.WalDir(ctx))); err != nil {
				return nil, err
			}

			if err := os.RemoveAll(ctx.ResolvePath(snapshotdb.DBPath)); err != nil {
				return nil, err
			}

			chainDb, err = ctx.OpenDatabaseWithFreezer("chaindata", config.DatabaseCache, config.DatabaseHandles, config.DatabaseFreezer, "eth/db/chaindata/")
			if err != nil {
				return nil, err
			}

			snapshotBaseDB, err = snapshotdb.Open(ctx.ResolvePath(snapshotdb.DBPath), config.DatabaseCache, config.DatabaseHandles, true)
			if err != nil {
				return nil, err
			}

			//only private net  need InitGenesisAndSetEconomicConfig
			if stored != params.MainnetGenesisHash && config.Genesis == nil {
				// private net
				config.Genesis = new(core.Genesis)
				if err := config.Genesis.InitGenesisAndSetEconomicConfig(ctx.GenesisPath()); err != nil {
					return nil, err
				}
			}
			log.Info("last fast sync is fail,init  db finish")
		}
	}

	chainConfig, _, genesisErr := core.SetupGenesisBlock(chainDb, snapshotBaseDB, config.Genesis)

	if err := snapshotBaseDB.Close(); err != nil {
		return nil, err
	}

	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}

	if chainConfig.Cbft.Period == 0 || chainConfig.Cbft.Amount == 0 {
		chainConfig.Cbft.Period = config.CbftConfig.Period
		chainConfig.Cbft.Amount = config.CbftConfig.Amount
	}

	log.Info("Initialised chain configuration", "config", chainConfig)

	eth := &Ethereum{
		config:            config,
		chainDb:           chainDb,
		eventMux:          ctx.EventMux,
		accountManager:    ctx.AccountManager,
		engine:            CreateConsensusEngine(ctx, chainConfig, config.Miner.Noverify, chainDb, &config.CbftConfig, ctx.EventMux),
		closeBloomHandler: make(chan struct{}),
		networkID:         config.NetworkId,
		gasPrice:          config.Miner.GasPrice,
		bloomRequests:     make(chan chan *bloombits.Retrieval),
		bloomIndexer:      NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
	}

	bcVersion := rawdb.ReadDatabaseVersion(chainDb)

	var dbVer = "<nil>"
	if bcVersion != nil {
		dbVer = fmt.Sprintf("%d", *bcVersion)
	}
	log.Info("Initialising PlatON protocol", "versions", ProtocolVersions, "network", config.NetworkId, "dbversion", dbVer)

	if !config.SkipBcVersionCheck {
		if bcVersion != nil && *bcVersion > core.BlockChainVersion {
			return nil, fmt.Errorf("database version is v%d, Geth %s only supports v%d", *bcVersion, params.VersionWithMeta, core.BlockChainVersion)
		} else if bcVersion == nil || *bcVersion < core.BlockChainVersion {
			log.Warn("Upgrade blockchain database version", "from", dbVer, "to", core.BlockChainVersion)
			rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
		}
	}

	var (
		vmConfig = vm.Config{
			ConsoleOutput: config.Debug,
			WasmType:      vm.Str2WasmType(config.VMWasmType),
		}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieDirtyLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout,
			BodyCacheLimit: config.BodyCacheLimit, BlockCacheLimit: config.BlockCacheLimit,
			MaxFutureBlocks: config.MaxFutureBlocks, BadBlockLimit: config.BadBlockLimit,
			TriesInMemory: config.TriesInMemory, TrieCleanLimit: config.TrieDBCache,
			DBGCInterval: config.DBGCInterval, DBGCTimeout: config.DBGCTimeout,
			DBGCMpt: config.DBGCMpt, DBGCBlock: config.DBGCBlock,
		}

		minningConfig = &core.MiningConfig{MiningLogAtDepth: config.MiningLogAtDepth, TxChanSize: config.TxChanSize,
			ChainHeadChanSize: config.ChainHeadChanSize, ChainSideChanSize: config.ChainSideChanSize,
			ResultQueueSize: config.ResultQueueSize, ResubmitAdjustChanSize: config.ResubmitAdjustChanSize,
			MinRecommitInterval: config.MinRecommitInterval, MaxRecommitInterval: config.MaxRecommitInterval,
			IntervalAdjustRatio: config.IntervalAdjustRatio, IntervalAdjustBias: config.IntervalAdjustBias,
			StaleThreshold: config.StaleThreshold, DefaultCommitRatio: config.DefaultCommitRatio,
		}
	)
	cacheConfig.DBDisabledGC.Set(config.DBDisabledGC)

	eth.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, chainConfig, eth.engine, vmConfig, eth.shouldPreserve)
	if err != nil {
		return nil, err
	}
	snapshotdb.SetDBBlockChain(eth.blockchain)

	blockChainCache := core.NewBlockChainCache(eth.blockchain)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		return nil, compat
		//eth.blockchain.SetHead(compat.RewindTo)
		//rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	eth.bloomIndexer.Start(eth.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	eth.txPool = core.NewTxPool(config.TxPool, chainConfig, core.NewTxPoolBlockChain(blockChainCache))

	core.SenderCacher.SetTxPool(eth.txPool)

	currentBlock := eth.blockchain.CurrentBlock()
	currentNumber := currentBlock.NumberU64()
	currentHash := currentBlock.Hash()
	gasCeil, err := gov.GovernMaxBlockGasLimit(currentNumber, currentHash)
	if nil != err {
		log.Error("Failed to query gasCeil from snapshotdb", "err", err)
		return nil, err
	}
	if config.Miner.GasFloor > uint64(gasCeil) {
		log.Error("The gasFloor must be less than gasCeil", "gasFloor", config.Miner.GasFloor, "gasCeil", gasCeil)
		return nil, fmt.Errorf("The gasFloor must be less than gasCeil, got: %d, expect range (0, %d]", config.Miner.GasFloor, gasCeil)
	}

	eth.miner = miner.New(eth, &config.Miner, eth.blockchain.Config(), minningConfig, eth.EventMux(), eth.engine,
		eth.isLocalBlock, blockChainCache, config.VmTimeoutDuration)

	reactor := core.NewBlockChainReactor(eth.EventMux(), eth.blockchain.Config().ChainID)
	node.GetCryptoHandler().SetPrivateKey(ctx.NodePriKey())

	if engine, ok := eth.engine.(consensus.Bft); ok {
		var agency consensus.Agency
		core.NewExecutor(eth.blockchain.Config(), eth.blockchain, vmConfig, eth.txPool)
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
			reactor.SetVRFhandler(handler.NewVrfHandler(eth.blockchain.Genesis().Nonce()))
			reactor.SetPluginEventMux()
			reactor.SetPrivateKey(ctx.NodePriKey())
			handlePlugin(reactor)
			agency = reactor

			//register Govern parameter verifiers
			gov.RegisterGovernParamVerifiers()
		}

		if err := recoverSnapshotDB(blockChainCache); err != nil {
			log.Error("recover SnapshotDB fail", "error", err)
			return nil, errors.New("Failed to recover SnapshotDB")
		}

		if err := engine.Start(eth.blockchain, blockChainCache, eth.txPool, agency); err != nil {
			log.Error("Init cbft consensus engine fail", "error", err)
			return nil, errors.New("Failed to init cbft consensus engine")
		}
	}

	// Permit the downloader to use the trie cache allowance during fast sync
	cacheLimit := cacheConfig.TrieCleanLimit + cacheConfig.TrieDirtyLimit
	if eth.protocolManager, err = NewProtocolManager(chainConfig, config.SyncMode, config.NetworkId, eth.eventMux, eth.txPool, eth.engine, eth.blockchain, chainDb, cacheLimit); err != nil {
		return nil, err
	}

	eth.APIBackend = &EthAPIBackend{ctx.ExtRPCEnabled(), eth, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.Miner.GasPrice
	}
	eth.APIBackend.gpo = gasprice.NewOracle(eth.APIBackend, gpoParams)

	return eth, nil
}

func recoverSnapshotDB(blockChainCache *core.BlockChainCache) error {
	sdb := snapshotdb.Instance()
	ch := sdb.GetCurrent().GetHighest(false).Num.Uint64()
	blockChanHegiht := blockChainCache.CurrentHeader().Number.Uint64()
	if ch < blockChanHegiht {
		for i := ch + 1; i <= blockChanHegiht; i++ {
			block, parentBlock := blockChainCache.GetBlockByNumber(i), blockChainCache.GetBlockByNumber(i-1)
			log.Debug("snapshotdb recover block from blockchain", "num", block.Number(), "hash", block.Hash())
			if err := blockChainCache.Execute(block, parentBlock); err != nil {
				log.Error("snapshotdb recover block from blockchain  execute fail", "error", err)
				return err
			}
			if err := sdb.Commit(block.Hash()); err != nil {
				log.Error("snapshotdb recover block from blockchain  Commit fail", "error", err)
				return err
			}
		}
	}
	return nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an Ethereum service
func CreateConsensusEngine(ctx *node.ServiceContext, chainConfig *params.ChainConfig, noverify bool, db ethdb.Database,
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
			Service:   NewPrivateDebugAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   xplugin.NewPublicPPOSAPI(),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
		{
			Namespace: "txgen",
			Version:   "1.0",
			Service:   NewTxGenAPI(s),
			Public:    true,
		},
	}...)
}

//func (s *Ethereum) ResetWithGenesisBlock(gb *types.Block) {
//	s.blockchain.ResetWithGenesisBlock(gb)
//}

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

	//log.Debug("node start", "srvr.Config.PrivateKey", srvr.Config.PrivateKey)
	if cbftEngine, ok := s.engine.(consensus.Bft); ok {
		if flag := cbftEngine.IsConsensusNode(); flag {
			for _, n := range s.blockchain.Config().Cbft.InitialNodes {
				// todo: Mock point.
				if !node.FakeNetEnable {
					srvr.AddConsensusPeer(discover.NewNode(n.Node.ID, n.Node.IP, n.Node.UDP, n.Node.TCP))
				}
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

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *Ethereum) Stop() error {
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}

	// Then stop everything else.
	s.bloomIndexer.Close()
	close(s.closeBloomHandler)
	s.txPool.Stop()
	s.miner.Stop()
	s.blockchain.Stop()
	s.engine.Close()
	core.GetReactorInstance().Close()
	s.chainDb.Close()
	s.eventMux.Stop()
	return nil
}

// RegisterPlugin one by one
func handlePlugin(reactor *core.BlockChainReactor) {
	xplugin.RewardMgrInstance().SetCurrentNodeID(reactor.NodeId)

	reactor.RegisterPlugin(xcom.SlashingRule, xplugin.SlashInstance())
	xplugin.SlashInstance().SetDecodeEvidenceFun(evidence.NewEvidence)
	reactor.RegisterPlugin(xcom.StakingRule, xplugin.StakingInstance())
	reactor.RegisterPlugin(xcom.RestrictingRule, xplugin.RestrictingInstance())
	reactor.RegisterPlugin(xcom.RewardRule, xplugin.RewardMgrInstance())

	xplugin.GovPluginInstance().SetChainID(reactor.GetChainID())
	reactor.RegisterPlugin(xcom.GovernanceRule, xplugin.GovPluginInstance())

	// set rule order
	reactor.SetBeginRule([]int{xcom.StakingRule, xcom.SlashingRule, xcom.CollectDeclareVersionRule, xcom.GovernanceRule})
	reactor.SetEndRule([]int{xcom.CollectDeclareVersionRule, xcom.RestrictingRule, xcom.RewardRule, xcom.GovernanceRule, xcom.StakingRule})

}
