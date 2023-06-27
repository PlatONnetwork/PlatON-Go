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
	"time"

	"github.com/PlatONnetwork/PlatON-Go/p2p/dnsdisc"

	"github.com/PlatONnetwork/PlatON-Go/eth/ethconfig"

	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"
	"github.com/PlatONnetwork/PlatON-Go/eth/protocols/eth"
	"github.com/PlatONnetwork/PlatON-Go/eth/protocols/snap"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	vrfhandler "github.com/PlatONnetwork/PlatON-Go/x/handler"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"

	"github.com/PlatONnetwork/PlatON-Go/accounts"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
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
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
	xplugin "github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

// Ethereum implements the Ethereum full node service.
type Ethereum struct {
	config *ethconfig.Config

	// Handlers
	txPool             *core.TxPool
	blockchain         *core.BlockChain
	handler            *handler
	ethDialCandidates  enode.Iterator
	snapDialCandidates enode.Iterator

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

	p2pServer *p2p.Server

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)
}

// New creates a new Ethereum object (including the
// initialisation of the common Ethereum object)
func New(stack *node.Node, config *ethconfig.Config) (*Ethereum, error) {
	// Ensure configuration values are compatible and sane
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run PlatON in light sync mode, use les.LightPlatON")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	if config.Miner.GasPrice == nil || config.Miner.GasPrice.Cmp(common.Big0) <= 0 {
		log.Warn("Sanitizing invalid miner gas price", "provided", config.Miner.GasPrice, "updated", ethconfig.Defaults.Miner.GasPrice)
		config.Miner.GasPrice = new(big.Int).Set(ethconfig.Defaults.Miner.GasPrice)
	}
	//if config.NoPruning && config.TrieDirtyCache > 0 {
	//	if config.SnapshotCache > 0 {
	//		config.TrieCleanCache += config.TrieDirtyCache * 3 / 5
	//		config.SnapshotCache += config.TrieDirtyCache * 2 / 5
	//	} else {
	//		config.TrieCleanCache += config.TrieDirtyCache
	//	}
	//	config.TrieDirtyCache = 0
	//}
	// Assemble the Ethereum object
	chainDb, err := stack.OpenDatabaseWithFreezer("chaindata", config.DatabaseCache, config.DatabaseHandles, config.DatabaseFreezer, "eth/db/chaindata/", false)
	if err != nil {
		return nil, err
	}
	snapshotdb.SetDBOptions(config.DatabaseCache, config.DatabaseHandles)

	snapshotBaseDB, err := snapshotdb.Open(stack.ResolvePath(snapshotdb.DBPath), config.DatabaseCache, config.DatabaseHandles, true)
	if err != nil {
		return nil, err
	}

	height := rawdb.ReadHeaderNumber(chainDb, rawdb.ReadHeadHeaderHash(chainDb))
	log.Debug("read header number from chain db", "height", height)
	if height != nil && *height > 0 {
		//when last fast syncing fail,we will clean chaindb,wal,snapshotdb
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
		if err == nil { // KeyFastSyncStatus == (FastSyncBegin || FastSyncFail)
			// Just commit the new block if there is no stored genesis block.
			stored := rawdb.ReadCanonicalHash(chainDb, 0)

			log.Info("last fast sync is fail,init db", "status", common.BytesToUint32(status), "prichain", config.Genesis == nil)
			chainDb.Close()
			if err := snapshotBaseDB.Close(); err != nil {
				return nil, err
			}

			if config.DatabaseFreezer != "" {
				if err := os.RemoveAll(stack.Config().ResolveFreezerPath("chaindata", config.DatabaseFreezer)); err != nil {
					return nil, err
				}
			}

			if err := os.RemoveAll(stack.ResolvePath("chaindata")); err != nil {
				return nil, err
			}

			if err := os.RemoveAll(stack.ResolvePath(wal.WalDir(stack))); err != nil {
				return nil, err
			}

			if err := os.RemoveAll(stack.ResolvePath(snapshotdb.DBPath)); err != nil {
				return nil, err
			}

			chainDb, err = stack.OpenDatabaseWithFreezer("chaindata", config.DatabaseCache, config.DatabaseHandles, config.DatabaseFreezer, "eth/db/chaindata/", false)
			if err != nil {
				return nil, err
			}

			snapshotBaseDB, err = snapshotdb.Open(stack.ResolvePath(snapshotdb.DBPath), config.DatabaseCache, config.DatabaseHandles, true)
			if err != nil {
				return nil, err
			}

			//only private net need InitGenesisAndSetEconomicConfig
			if stored != params.MainnetGenesisHash && config.Genesis == nil {
				// private net
				config.Genesis = new(core.Genesis)
				if err := config.Genesis.InitGenesisAndSetEconomicConfig(stack.GenesisPath()); err != nil {
					return nil, err
				}
			}
			log.Info("last fast sync is fail,init db finish")
		} else { // err == snapshotdb.ErrNotFound
			// Just commit the new block if there is no stored genesis block.
			stored := rawdb.ReadCanonicalHash(chainDb, 0)
			//todo 这是一个暂时的hack方法,针对我们的测试链使用,待测试链版本升级到1.5.0后此方法可以删除
			if stored != params.MainnetGenesisHash && config.Genesis == nil {
				// private net
				config.Genesis = new(core.Genesis)
				if err := config.Genesis.InitGenesisAndSetEconomicConfig(stack.GenesisPath()); err != nil {
					return nil, err
				}
			}
		}
	}

	chainConfig, _, genesisErr := core.SetupGenesisBlock(chainDb, snapshotBaseDB, config.Genesis)

	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}

	if chainConfig.Cbft.Period == 0 || chainConfig.Cbft.Amount == 0 {
		chainConfig.Cbft.Period = config.CbftConfig.Period
		chainConfig.Cbft.Amount = config.CbftConfig.Amount
	}

	log.Info("Initialised chain configuration", "config", chainConfig)
	stack.SetP2pChainID(chainConfig.ChainID, chainConfig.PIP7ChainID)

	eth := &Ethereum{
		config:            config,
		chainDb:           chainDb,
		eventMux:          stack.EventMux(),
		accountManager:    stack.AccountManager(),
		engine:            ethconfig.CreateConsensusEngine(stack, chainConfig, config.Miner.Noverify, chainDb, &config.CbftConfig, stack.EventMux()),
		closeBloomHandler: make(chan struct{}),
		networkID:         config.NetworkId,
		gasPrice:          config.Miner.GasPrice,
		bloomRequests:     make(chan chan *bloombits.Retrieval),
		bloomIndexer:      core.NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
		p2pServer:         stack.Server(),
	}

	bcVersion := rawdb.ReadDatabaseVersion(chainDb)

	var dbVer = "<nil>"
	if bcVersion != nil {
		dbVer = fmt.Sprintf("%d", *bcVersion)
	}
	log.Info("Initialising PlatON protocol", "network", config.NetworkId, "dbversion", dbVer)

	if !config.SkipBcVersionCheck {
		if bcVersion != nil && *bcVersion > core.BlockChainVersion {
			return nil, fmt.Errorf("database version is v%d, PlatON %s only supports v%d", *bcVersion, params.VersionWithMeta, core.BlockChainVersion)
		} else if bcVersion == nil || *bcVersion < core.BlockChainVersion {
			if bcVersion != nil { // only print warning on upgrade, not on init
				log.Warn("Upgrade blockchain database version", "from", dbVer, "to", core.BlockChainVersion)
			}
			rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
		}
	}

	var (
		vmConfig = vm.Config{
			ConsoleOutput: config.Debug,
			WasmType:      vm.Str2WasmType(config.VMWasmType),
		}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieDirtyLimit: config.TrieDirtyCache, TrieTimeLimit: config.TrieTimeout,
			SnapshotLimit:  config.SnapshotCache,
			BodyCacheLimit: config.BodyCacheLimit, BlockCacheLimit: config.BlockCacheLimit,
			MaxFutureBlocks: config.MaxFutureBlocks,
			TriesInMemory:   config.TriesInMemory, TrieCleanLimit: config.TrieDBCache, Preimages: config.Preimages,
			TrieCleanJournal:   stack.ResolvePath(config.TrieCleanCacheJournal),
			TrieCleanRejournal: config.TrieCleanCacheRejournal,
			DBGCInterval:       config.DBGCInterval, DBGCTimeout: config.DBGCTimeout,
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

	eth.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, chainConfig, eth.engine, vmConfig, eth.shouldPreserve, &config.TxLookupLimit)
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
		config.TxPool.Journal = stack.ResolvePath(config.TxPool.Journal)
	}
	eth.txPool = core.NewTxPool(config.TxPool, chainConfig, core.NewTxPoolBlockChain(blockChainCache))

	core.SenderCacher.SetTxPool(eth.txPool)

	currentBlock := eth.blockchain.CurrentBlock()
	currentNumber := currentBlock.NumberU64()
	currentHash := currentBlock.Hash()
	gasCeil, err := gov.GovernMaxBlockGasLimit(currentNumber, currentHash, snapshotBaseDB)
	if err := snapshotBaseDB.Close(); err != nil {
		return nil, err
	}
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
	node.GetCryptoHandler().SetPrivateKey(stack.Config().NodeKey())

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
			reactor.SetVRFhandler(vrfhandler.NewVrfHandler(eth.blockchain.Genesis().Nonce()))
			reactor.SetPluginEventMux()
			reactor.SetPrivateKey(stack.Config().NodeKey())
			handlePlugin(reactor, chainDb, chainConfig, config.DBValidatorsHistory)
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
	} else {
		log.Crit("engin not good")
	}

	// Permit the downloader to use the trie cache allowance during fast sync
	cacheLimit := cacheConfig.TrieCleanLimit + cacheConfig.TrieDirtyLimit + cacheConfig.SnapshotLimit
	checkpoint := config.Checkpoint
	if checkpoint == nil {
		checkpoint = params.TrustedCheckpoints[params.MainnetGenesisHash]
	}
	if eth.handler, err = newHandler(&handlerConfig{
		Database:   chainDb,
		Chain:      eth.blockchain,
		TxPool:     eth.txPool,
		Network:    config.NetworkId,
		Sync:       config.SyncMode,
		BloomCache: uint64(cacheLimit),
		EventMux:   eth.eventMux,
		Checkpoint: checkpoint,
		Whitelist:  config.Whitelist,
	}); err != nil {
		return nil, err
	}
	eth.APIBackend = &EthAPIBackend{stack.Config().ExtRPCEnabled(), stack.Config().AllowUnprotectedTxs, eth, nil}
	if eth.APIBackend.allowUnprotectedTxs {
		log.Info("Unprotected transactions allowed")
	}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.Miner.GasPrice
	}

	eth.APIBackend.gpo = gasprice.NewOracle(eth.APIBackend, gpoParams)
	// Setup DNS discovery iterators.
	dnsclient := dnsdisc.NewClient(dnsdisc.Config{})
	eth.ethDialCandidates, err = dnsclient.NewIterator(eth.config.EthDiscoveryURLs...)
	if err != nil {
		return nil, err
	}
	eth.snapDialCandidates, err = dnsclient.NewIterator(eth.config.SnapDiscoveryURLs...)
	if err != nil {
		return nil, err
	}
	// Start the RPC service
	eth.netRPCService = ethapi.NewPublicNetAPI(eth.p2pServer, config.NetworkId)

	// Register the backend on the node
	stack.RegisterAPIs(eth.APIs())
	stack.RegisterProtocols(eth.Protocols())
	stack.RegisterLifecycle(eth)

	// Check for unclean shutdown
	if uncleanShutdowns, discards, err := rawdb.PushUncleanShutdownMarker(chainDb); err != nil {
		log.Error("Could not update unclean-shutdown-marker list", "error", err)
	} else {
		if discards > 0 {
			log.Warn("Old unclean shutdowns found", "count", discards)
		}
		for _, tstamp := range uncleanShutdowns {
			t := time.Unix(int64(tstamp), 0)
			log.Warn("Unclean shutdown detected", "booted", t,
				"age", common.PrettyAge(t))
		}
	}
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
			Service:   downloader.NewPublicDownloaderAPI(s.handler.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "platon",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false, 5*time.Minute),
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

// start mining
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
		atomic.StoreUint32(&s.handler.acceptTxs, 1)

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
func (s *Ethereum) Downloader() *downloader.Downloader { return s.handler.downloader }
func (s *Ethereum) Synced() bool                       { return atomic.LoadUint32(&s.handler.acceptTxs) == 1 }
func (s *Ethereum) BloomIndexer() *core.ChainIndexer   { return s.bloomIndexer }

// Protocols returns all the currently configured
// network protocols to start.
func (s *Ethereum) Protocols() []p2p.Protocol {
	protos := eth.MakeProtocols((*ethHandler)(s.handler), s.networkID, s.ethDialCandidates)
	protos = append(protos, s.Engine().Protocols()...)
	if s.config.SnapshotCache > 0 {
		protos = append(protos, snap.MakeProtocols((*snapHandler)(s.handler), s.snapDialCandidates)...)
	}
	return protos
}

// Start implements node.Lifecycle, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *Ethereum) Start() error {
	eth.StartENRUpdater(s.blockchain, s.p2pServer.LocalNode())

	// Start the bloom bits servicing goroutines
	s.startBloomHandlers(params.BloomBitsBlocks)

	// Figure out a max peers count based on the server limits
	maxPeers := s.p2pServer.MaxPeers
	// Start the networking layer and the light server if requested
	s.handler.Start(maxPeers)

	//log.Debug("node start", "srvr.Config.PrivateKey", srvr.Config.PrivateKey)
	if cbftEngine, ok := s.engine.(consensus.Bft); ok {
		if flag := cbftEngine.IsConsensusNode(); flag {
			for _, n := range s.blockchain.Config().Cbft.InitialNodes {
				// todo: Mock point.
				if !node.FakeNetEnable {
					s.p2pServer.AddConsensusPeer(n.Node)
				}
			}
		}
		s.StartMining()
	}
	s.p2pServer.StartWatching(s.eventMux)

	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *Ethereum) Stop() error {
	s.ethDialCandidates.Close()
	s.snapDialCandidates.Close()
	s.handler.Stop()

	// Then stop everything else.
	// Only the operations related to block execution are stopped here
	// and engine.Close cannot be called directly because it has a dependency on the following modules
	s.engine.Stop()
	s.bloomIndexer.Close()
	close(s.closeBloomHandler)
	s.txPool.Stop()
	s.miner.Stop()
	s.blockchain.Stop()
	s.engine.Close()
	core.GetReactorInstance().Close()
	rawdb.PopUncleanShutdownMarker(s.chainDb)
	s.chainDb.Close()
	s.eventMux.Stop()
	return nil
}

// RegisterPlugin one by one
func handlePlugin(reactor *core.BlockChainReactor, chainDB ethdb.Database, chainConfig *params.ChainConfig, isValidatorsHistory bool) {
	xplugin.RewardMgrInstance().SetCurrentNodeID(reactor.NodeId)

	reactor.RegisterPlugin(xcom.SlashingRule, xplugin.SlashInstance())
	xplugin.SlashInstance().SetDecodeEvidenceFun(evidence.NewEvidence)
	reactor.RegisterPlugin(xcom.StakingRule, xplugin.StakingInstance())
	reactor.RegisterPlugin(xcom.RestrictingRule, xplugin.RestrictingInstance())
	reactor.RegisterPlugin(xcom.RewardRule, xplugin.RewardMgrInstance())

	xplugin.GovPluginInstance().SetChainID(reactor.GetChainID())
	xplugin.GovPluginInstance().SetChainDB(chainDB)
	reactor.RegisterPlugin(xcom.GovernanceRule, xplugin.GovPluginInstance())

	xplugin.StakingInstance().SetChainDB(chainDB, chainDB)
	xplugin.StakingInstance().SetChainConfig(chainConfig)
	if isValidatorsHistory {
		xplugin.StakingInstance().EnableValidatorsHistory()
	}

	// set rule order
	reactor.SetBeginRule([]int{xcom.StakingRule, xcom.SlashingRule, xcom.CollectDeclareVersionRule, xcom.GovernanceRule})
	reactor.SetEndRule([]int{xcom.CollectDeclareVersionRule, xcom.RestrictingRule, xcom.RewardRule, xcom.GovernanceRule, xcom.StakingRule})

}
