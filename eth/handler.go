// Copyright 2015 The PlatON-Go Authors
// This file is part of the PlatON-Go library.
//
// The go-PlatON library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/eth/fetcher"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"
)

const (
	// softResponseLimit is the target maximum size of replies to data retrievals.
	softResponseLimit = 2 * 1024 * 1024

	estHeaderRlpSize = 500 // Approximate size of an RLP encoded block header

	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096

	numBroadcastTxPeers     = 5 // Maximum number of peers for broadcast transactions
	numBroadcastTxHashPeers = 5 // Maximum number of peers for broadcast transactions hash
	numBroadcastBlockPeers  = 5 // Maximum number of peers for broadcast new block

	defaultTxsCacheSize      = 20
	defaultBroadcastInterval = 100 * time.Millisecond
)

// errIncompatibleConfig is returned if the requested protocols and configs are
// not compatible (low protocol version restrictions and high requirements).
var errIncompatibleConfig = errors.New("incompatible configuration")

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

type ProtocolManager struct {
	networkID uint64

	fastSync        uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	acceptTxs       uint32 // Flag whether we're considered synchronised (enables transaction processing)
	acceptRemoteTxs uint32 // Flag whether we're accept remote txs

	txpool      txPool
	blockchain  *core.BlockChain
	chainconfig *params.ChainConfig
	chaindb     ethdb.Database
	maxPeers    int

	downloader   *downloader.Downloader
	blockFetcher *fetcher.BlockFetcher
	txFetcher    *fetcher.TxFetcher
	peers        *peerSet

	SubProtocols []p2p.Protocol

	eventMux      *event.TypeMux
	txsCh         chan core.NewTxsEvent
	txsCache      []*types.Transaction
	txsSub        event.Subscription
	minedBlockSub *event.TypeMuxSubscription

	prepareMinedBlockSub *event.TypeMuxSubscription
	blockSignatureSub    *event.TypeMuxSubscription

	// channels for fetcher, syncer, txsyncLoop
	txsyncCh chan *txsync
	quitSync chan struct{}

	chainSync *chainSyncer
	wg        sync.WaitGroup
	peerWG    sync.WaitGroup

	engine consensus.Engine
}

// NewProtocolManager returns a new PlatON sub protocol manager. The PlatON sub protocol manages peers capable
// with the PlatON network.
func NewProtocolManager(config *params.ChainConfig, mode downloader.SyncMode, networkID uint64, mux *event.TypeMux, txpool txPool, engine consensus.Engine, blockchain *core.BlockChain, chaindb ethdb.Database, cacheLimit int) (*ProtocolManager, error) {
	// Create the protocol manager with the base fields
	manager := &ProtocolManager{
		networkID:   networkID,
		eventMux:    mux,
		txpool:      txpool,
		blockchain:  blockchain,
		chaindb:     chaindb,
		chainconfig: config,
		peers:       newPeerSet(),
		txsyncCh:    make(chan *txsync),
		quitSync:    make(chan struct{}),
		engine:      engine,
	}
	// If fast sync was requested and our database is empty, grant it
	if mode == downloader.FastSync && blockchain.CurrentBlock().NumberU64() == 0 {
		manager.fastSync = uint32(1)
	}
	// Initiate a sub-protocol for every implemented version we can handle
	manager.SubProtocols = make([]p2p.Protocol, 0, len(ProtocolVersions))
	for i, version := range ProtocolVersions {
		// Skip protocol version if incompatible with the mode of operation
		if atomic.LoadUint32(&manager.fastSync) == 1 && version < eth63 {
			continue
		}
		// Compatible; initialise the sub-protocol
		version := version // Closure for the run
		manager.SubProtocols = append(manager.SubProtocols, p2p.Protocol{
			Name:    protocolName,
			Version: version,
			Length:  protocolLengths[i],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				return manager.runPeer(manager.newPeer(int(version), p, rw, manager.txpool.Get))
			},
			NodeInfo: func() interface{} {
				return manager.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := manager.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}
	if len(manager.SubProtocols) == 0 {
		return nil, errIncompatibleConfig
	}

	decodeExtra := func(extra []byte) (common.Hash, uint64, error) {
		return manager.engine.DecodeExtra(extra)
	}

	// Construct the downloader (long sync) and its backing state bloom if fast
	// sync is requested. The downloader is responsible for deallocating the state
	// bloom when it's done.
	var stateBloom *trie.SyncBloom
	if atomic.LoadUint32(&manager.fastSync) == 1 {
		stateBloom = trie.NewSyncBloom(uint64(cacheLimit), chaindb)
	}
	manager.downloader = downloader.New(chaindb, snapshotdb.Instance(), stateBloom, manager.eventMux, blockchain, nil, manager.removePeer, decodeExtra)

	// Construct the fetcher (short sync)
	validator := func(header *types.Header) error {
		return engine.VerifyHeader(blockchain, header, true)
	}
	heighter := func() uint64 {
		return manager.blockchain.Engine().CurrentBlock().NumberU64() + 1
	}
	inserter := func(blocks types.Blocks) (int, error) {
		// If fast sync is running, deny importing weird blocks
		if atomic.LoadUint32(&manager.fastSync) == 1 {
			log.Warn("Discarded bad propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		atomic.StoreUint32(&manager.acceptTxs, 1) // Mark initial sync done on any fetcher import
		return manager.blockchain.InsertChain(blocks)
	}
	getBlockByHash := func(hash common.Hash) *types.Block {
		return manager.blockchain.GetBlockByHash(hash)
	}
	manager.blockFetcher = fetcher.NewBlockFetcher(getBlockByHash, validator, manager.BroadcastBlock, heighter, inserter, manager.removePeer, decodeExtra)

	fetchTx := func(peer string, hashes []common.Hash) error {
		p := manager.peers.Peer(peer)
		if p == nil {
			return errors.New("unknown peer")
		}
		return p.RequestTxs(hashes)
	}
	manager.txFetcher = fetcher.NewTxFetcher(manager.txpool.Has, manager.txpool.AddRemotes, fetchTx)

	manager.chainSync = newChainSyncer(manager)

	return manager, nil
}

func (pm *ProtocolManager) removePeer(id string) {
	// Short circuit if the peer was already removed
	peer := pm.peers.Peer(id)
	if peer == nil {
		return
	}
	log.Debug("Removing PlatON peer", "peer", id)

	// Unregister the peer from the downloader and PlatON peer set
	pm.downloader.UnregisterPeer(id)
	pm.txFetcher.Drop(id)

	if err := pm.peers.Unregister(id); err != nil {
		log.Error("Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

func (pm *ProtocolManager) Start(maxPeers int) {
	pm.maxPeers = maxPeers

	// broadcast transactions
	pm.txsCh = make(chan core.NewTxsEvent, txChanSize)
	pm.txsSub = pm.txpool.SubscribeNewTxsEvent(pm.txsCh)
	pm.wg.Add(1)
	go pm.txBroadcastLoop()

	// broadcast mined blocks
	pm.minedBlockSub = pm.eventMux.Subscribe(core.NewMinedBlockEvent{})
	pm.wg.Add(1)
	go pm.minedBroadcastLoop()

	// start sync handlers
	pm.wg.Add(2)
	go pm.chainSync.loop()
	go pm.txsyncLoop()
}

func (pm *ProtocolManager) Stop() {
	log.Info("Stopping PlatON protocol")

	pm.txsSub.Unsubscribe()        // quits txBroadcastLoop
	pm.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop

	// Quit chainSync and txsync.
	// After this send has completed, no new peers will be accepted.
	close(pm.quitSync)
	pm.wg.Wait()

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to pm.peers yet
	// will exit when they try to register.
	pm.peers.Close()
	pm.peerWG.Wait()

	log.Info("PlatON protocol stopped")
}

func (pm *ProtocolManager) runPeer(p *peer) error {
	if !pm.chainSync.handlePeerEvent(p) {
		return p2p.DiscQuitting
	}
	pm.peerWG.Add(1)
	defer pm.peerWG.Done()
	return pm.handle(p)
}

func (pm *ProtocolManager) newPeer(pv int, p *p2p.Peer, rw p2p.MsgReadWriter, getPooledTx func(hash common.Hash) *types.Transaction) *peer {
	return newPeer(pv, p, rw, getPooledTx)
}

// handle is the callback invoked to manage the life cycle of an eth peer. When
// this function terminates, the peer is disconnected.
func (pm *ProtocolManager) handle(p *peer) error {
	// Ignore maxPeers if this is a trusted peer
	if pm.peers.Len() >= pm.maxPeers && !p.Peer.Info().Network.Trusted && !p.Peer.Info().Network.Static {
		return p2p.DiscTooManyPeers
	}
	p.Log().Debug("PlatON peer connected", "name", p.Name())

	// Execute the PlatON handshake
	var (
		genesis = pm.blockchain.Genesis()
		head    = pm.blockchain.CurrentHeader()
		hash    = head.CacheHash()
	)
	if err := p.Handshake(pm.networkID, head.Number, hash, genesis.Hash(), pm); err != nil {
		p.Log().Debug("PlatON handshake failed", "err", err)
		return err
	}
	// Register the peer locally
	if err := pm.peers.Register(p); err != nil {
		p.Log().Error("PlatON peer registration failed", "err", err)
		return err
	}
	defer pm.removePeer(p.id)

	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	if err := pm.downloader.RegisterPeer(p.id, p.version, p); err != nil {
		return err
	}
	pm.chainSync.handlePeerEvent(p)

	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	pm.syncTransactions(p)

	// main loop. handle incoming messages.
	for {
		if err := pm.handleMsg(p); err != nil {
			p.Log().Error("PlatON message handling failed", "err", err)
			return err
		}
	}
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (pm *ProtocolManager) handleMsg(p *peer) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := p.rw.ReadMsg()
	if err != nil {
		p.Log().Error("read peer message error", "err", err)
		return err
	}
	if msg.Size > protocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, protocolMaxMsgSize)
	}
	defer msg.Discard()

	// Handle the message depending on its contents
	switch {
	case msg.Code == StatusMsg:
		// Status messages should never arrive after the handshake
		return errResp(ErrExtraStatusMsg, "uncontrolled status message")

	// Block header query, collect the requested headers and reply
	case msg.Code == GetBlockHeadersMsg:
		// Decode the complex header query
		var query getBlockHeadersData
		if err := msg.Decode(&query); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		hashMode := query.Origin.Hash != (common.Hash{})
		p.Log().Debug("[GetBlockHeadersMsg]Received a broadcast message", "origin.Number", query.Origin.Number,
			"origin.Hash", query.Origin.Hash, "skip", query.Skip, "amount", query.Amount,
			"reverse", query.Reverse, "number", pm.blockchain.CurrentBlock().Number(),
			"hash", pm.blockchain.CurrentBlock().Hash())
		first := true
		maxNonCanonical := uint64(100)

		// Gather headers until the fetch or network limits is reached
		var (
			bytes   common.StorageSize
			headers []*types.Header
			unknown bool
		)
		for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit && len(headers) < downloader.MaxHeaderFetch {
			// Retrieve the next header satisfying the query
			var origin *types.Header
			if hashMode {
				if first {
					first = false
					origin = pm.blockchain.GetHeaderByHash(query.Origin.Hash)
					if origin != nil {
						query.Origin.Number = origin.Number.Uint64()
					}
				} else {
					origin = pm.blockchain.GetHeader(query.Origin.Hash, query.Origin.Number)
				}
			} else {
				origin = pm.blockchain.GetHeaderByNumber(query.Origin.Number)
			}
			if origin == nil {
				break
			}
			headers = append(headers, origin)
			bytes += estHeaderRlpSize

			// Advance to the next header of the query
			switch {
			case hashMode && query.Reverse:
				// Hash based traversal towards the genesis block
				ancestor := query.Skip + 1
				if ancestor == 0 {
					unknown = true
				} else {
					query.Origin.Hash, query.Origin.Number = pm.blockchain.GetAncestor(query.Origin.Hash, query.Origin.Number, ancestor, &maxNonCanonical)
					unknown = (query.Origin.Hash == common.Hash{})
				}
			case hashMode && !query.Reverse:
				// Hash based traversal towards the leaf block
				var (
					current = origin.Number.Uint64()
					next    = current + query.Skip + 1
				)
				if next <= current {
					infos, _ := json.MarshalIndent(p.Peer.Info(), "", "  ")
					p.Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
					unknown = true
				} else {
					if header := pm.blockchain.GetHeaderByNumber(next); header != nil {
						nextHash := header.Hash()
						expOldHash, _ := pm.blockchain.GetAncestor(nextHash, next, query.Skip+1, &maxNonCanonical)
						if expOldHash == query.Origin.Hash {
							query.Origin.Hash, query.Origin.Number = nextHash, next
						} else {
							unknown = true
						}
					} else {
						unknown = true
					}
				}
			case query.Reverse:
				// Number based traversal towards the genesis block
				if query.Origin.Number >= query.Skip+1 {
					query.Origin.Number -= query.Skip + 1
				} else {
					unknown = true
				}

			case !query.Reverse:
				// Number based traversal towards the leaf block
				query.Origin.Number += query.Skip + 1
			}
		}
		p.Log().Debug("Send headers", "headers", len(headers))
		return p.SendBlockHeaders(headers)
	case p.version >= eth63 && msg.Code == GetOriginAndPivotMsg:
		p.Log().Info("[GetOriginAndPivotMsg]Received a broadcast message")
		var query uint64
		if err := msg.Decode(&query); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		oHead := pm.blockchain.GetHeaderByNumber(query)
		pivot, err := snapshotdb.Instance().BaseNum()
		if err != nil {
			p.Log().Error("GetOriginAndPivot get snapshotdb baseNum fail", "err", err)
			return errors.New("GetOriginAndPivot get snapshotdb baseNum fail")
		}
		if pivot == nil {
			p.Log().Error("[GetOriginAndPivot] pivot should not be nil")
			return errors.New("[GetOriginAndPivot] pivot should not be nil")
		}
		pHead := pm.blockchain.GetHeaderByNumber(pivot.Uint64())

		data := make([]*types.Header, 0)
		data = append(data, oHead, pHead)
		if err := p.SendOriginAndPivot(data); err != nil {
			p.Log().Error("[GetOriginAndPivotMsg]send data meassage fail", "error", err)
			return err
		}
	case p.version >= eth63 && msg.Code == OriginAndPivotMsg:
		p.Log().Debug("[OriginAndPivotMsg]Received a broadcast message")
		var data []*types.Header
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverOriginAndPivot(p.id, data); err != nil {
			p.Log().Error("Failed to deliver ppos storage data", "err", err)
			return err
		}
	case p.version >= eth63 && msg.Code == GetPPOSStorageMsg:
		p.Log().Info("[GetPPOSStorageMsg]Received a broadcast message")
		var query []interface{}
		if err := msg.Decode(&query); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		f := func(num *big.Int, iter iterator.Iterator) error {
			var psInfo PPOSInfo
			if num == nil {
				return errors.New("num should not be nil")
			}
			psInfo.Pivot = pm.blockchain.GetHeaderByNumber(num.Uint64())
			psInfo.Latest = pm.blockchain.CurrentHeader()
			if err := p.SendPPOSInfo(psInfo); err != nil {
				p.Log().Error("[GetPPOSStorageMsg]send last ppos meassage fail", "error", err)
				return err
			}
			var (
				byteSize int
				ps       PPOSStorage
				count    int
			)
			ps.KVs = make([]downloader.PPOSStorageKV, 0)
			for iter.Next() {
				if bytes.Equal(iter.Key(), []byte(snapshotdb.CurrentHighestBlock)) || bytes.Equal(iter.Key(), []byte(snapshotdb.CurrentBaseNum)) || bytes.HasPrefix(iter.Key(), []byte(snapshotdb.WalKeyPrefix)) {
					continue
				}
				byteSize = byteSize + len(iter.Key()) + len(iter.Value())
				if count >= downloader.PPOSStorageKVSizeFetch || byteSize > softResponseLimit {
					if err := p.SendPPOSStorage(ps); err != nil {
						p.Log().Error("[GetPPOSStorageMsg]send ppos message fail", "error", err, "kvnum", ps.KVNum)
						return err
					}
					count = 0
					ps.KVs = make([]downloader.PPOSStorageKV, 0)
					byteSize = 0
				}
				k, v := make([]byte, len(iter.Key())), make([]byte, len(iter.Value()))
				copy(k, iter.Key())
				copy(v, iter.Value())
				ps.KVs = append(ps.KVs, [2][]byte{
					k, v,
				})
				ps.KVNum++
				count++
			}
			ps.Last = true
			if err := p.SendPPOSStorage(ps); err != nil {
				p.Log().Error("[GetPPOSStorageMsg]send last ppos message fail", "error", err)
				return err
			}
			return nil
		}
		go func() {
			if err := snapshotdb.Instance().WalkBaseDB(nil, f); err != nil {
				p.Log().Error("[GetPPOSStorageMsg]send  ppos storage fail", "error", err)
			}
		}()

	case p.version >= eth63 && msg.Code == PPOSStorageMsg:
		p.Log().Debug("Received a broadcast message[PposStorageMsg]")
		var data PPOSStorage
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverPposStorage(p.id, data.KVs, data.Last, data.KVNum); err != nil {
			p.Log().Error("Failed to deliver ppos storage data", "err", err)
		}
	case p.version >= eth63 && msg.Code == PPOSInfoMsg:
		p.Log().Debug("Received a broadcast message[PPOSInfoMsg]")
		var data PPOSInfo
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverPposInfo(p.id, data.Latest, data.Pivot); err != nil {
			p.Log().Error("Failed to deliver ppos storage data", "err", err)
		}
	case msg.Code == BlockHeadersMsg:
		p.Log().Debug("Receive BlockHeadersMsg")
		// A batch of headers arrived to one of our previous requests
		var headers []*types.Header
		if err := msg.Decode(&headers); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		p.Log().Debug("Receive BlockHeadersMsg, before filter", "headers", len(headers))
		// Filter out any explicitly requested headers, deliver the rest to the downloader
		filter := len(headers) == 1
		if filter {
			// Irrelevant of the fork checks, send the header to the fetcher just in case
			headers = pm.blockFetcher.FilterHeaders(p.id, headers, time.Now())
		}
		p.Log().Debug("Receive BlockHeadersMsg, after filter", "headers", len(headers))
		if len(headers) > 0 || !filter {
			err := pm.downloader.DeliverHeaders(p.id, headers)
			if err != nil {
				log.Debug("Failed to deliver headers", "err", err)
			}
		}

	case msg.Code == GetBlockBodiesMsg:
		p.Log().Debug("Receive GetBlockBodiesMsg", "number", pm.blockchain.CurrentBlock().Number(), "hash", pm.blockchain.CurrentBlock().Hash())
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather blocks until the fetch or network limits is reached
		var (
			hash   common.Hash
			bytes  int
			bodies []rlp.RawValue
		)
		for bytes < softResponseLimit && len(bodies) < downloader.MaxBlockFetch {
			// Retrieve the hash of the next block
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested block body, stopping if enough was found
			log.Debug(fmt.Sprintf("Send block body peer:%s,hash:%v", p.id, hash.Hex()))
			if data := pm.blockchain.GetBodyRLP(hash); len(data) != 0 {
				bodies = append(bodies, data)
				bytes += len(data)
			} else {
				log.Debug(fmt.Sprintf("Block body empty peer:%s hash:%s", p.id, hash.TerminalString()))
			}
		}

		log.Debug(fmt.Sprintf("Send block body peer:%s,bytes:%d,bodies:%d", p.id, bytes, len(bodies)))
		return p.SendBlockBodiesRLP(bodies)

	case msg.Code == BlockBodiesMsg:
		log.Debug("Receive BlockBodiesMsg", "peer", p.id)
		// A batch of block bodies arrived to one of our previous requests
		var request blockBodiesData
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver them all to the downloader for queuing
		transactions := make([][]*types.Transaction, len(request))
		extraData := make([][]byte, len(request))

		for i, body := range request {
			transactions[i] = body.Transactions
			extraData[i] = body.ExtraData
		}
		// Filter out any explicitly requested bodies, deliver the rest to the downloader
		filter := len(transactions) > 0 || len(extraData) > 0
		log.Debug("Receive BlockBodiesMsg", "peer", p.id, "txslen", len(transactions), "extradata", len(extraData))
		if filter {
			transactions, extraData = pm.blockFetcher.FilterBodies(p.id, transactions, extraData, time.Now())
		}
		log.Debug("Receive BlockBodiesMsg", "peer", p.id, "txslen", len(transactions), "extradata", len(extraData))

		if len(transactions) > 0 || len(extraData) > 0 || !filter {
			err := pm.downloader.DeliverBodies(p.id, transactions, extraData)
			if err != nil {
				log.Debug("Failed to deliver bodies", "peer", p.id, "err", err)
			}
		}

	case p.version >= eth63 && msg.Code == GetNodeDataMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather state data until the fetch or network limits is reached
		var (
			hash  common.Hash
			bytes int
			data  [][]byte
		)
		for bytes < softResponseLimit && len(data) < downloader.MaxStateFetch {
			// Retrieve the hash of the next state entry
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested state entry, stopping if enough was found
			// todo now the code and trienode is mixed in the protocol level,
			// separate these two types.
			if !pm.downloader.SyncBloomContains(hash[:]) {
				// Only lookup the trie node if there's chance that we actually have it
				continue
			}
			// Retrieve the requested state entry, stopping if enough was found
			// todo now the code and trienode is mixed in the protocol level,
			// separate these two types.
			entry, err := pm.blockchain.TrieNode(hash)
			if len(entry) == 0 || err != nil {
				// Read the contract code with prefix only to save unnecessary lookups.
				entry, err = pm.blockchain.ContractCodeWithPrefix(hash)
			}
			if err == nil && len(entry) > 0 {
				data = append(data, entry)
				bytes += len(entry)
			}
		}
		return p.SendNodeData(data)

	case p.version >= eth63 && msg.Code == NodeDataMsg:
		// A batch of node state data arrived to one of our previous requests
		var data [][]byte
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverNodeData(p.id, data); err != nil {
			log.Debug("Failed to deliver node state data", "err", err)
		}

	case p.version >= eth63 && msg.Code == GetReceiptsMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather state data until the fetch or network limits is reached
		var (
			hash     common.Hash
			bytes    int
			receipts []rlp.RawValue
		)
		for bytes < softResponseLimit && len(receipts) < downloader.MaxReceiptFetch {
			// Retrieve the hash of the next block
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested block's receipts, skipping if unknown to us
			results := pm.blockchain.GetReceiptsByHash(hash)
			if results == nil {
				if header := pm.blockchain.GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
					continue
				}
			}
			// If known, encode and queue for response packet
			if encoded, err := rlp.EncodeToBytes(results); err != nil {
				log.Error("Failed to encode receipt", "err", err)
			} else {
				receipts = append(receipts, encoded)
				bytes += len(encoded)
			}
		}
		return p.SendReceiptsRLP(receipts)

	case p.version >= eth63 && msg.Code == ReceiptsMsg:
		// A batch of receipts arrived to one of our previous requests
		var receipts [][]*types.Receipt
		if err := msg.Decode(&receipts); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverReceipts(p.id, receipts); err != nil {
			log.Debug("Failed to deliver receipts", "err", err)
		}

	case msg.Code == NewBlockHashesMsg:
		var announces newBlockHashesData
		if err := msg.Decode(&announces); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}

		// Mark the hashes as present at the remote node
		for _, block := range announces {
			p.MarkBlock(block.Hash)
			log.Debug("Received a message[NewBlockHashesMsg]------------", "receiveAt", msg.ReceivedAt.Unix(), "peerId", p.id, "hash", block.Hash, "number", block.Number)
		}
		// Schedule all the unknown hashes for retrieval
		unknown := make(newBlockHashesData, 0, len(announces))
		for _, block := range announces {
			if !pm.blockchain.Engine().HasBlock(block.Hash, block.Number) {
				unknown = append(unknown, block)
			}
		}
		for _, block := range unknown {
			log.Debug("Unknown block", "hash", block.Hash, "number", block.Number)
			pm.blockFetcher.Notify(p.id, block.Hash, block.Number, time.Now(), p.RequestOneHeader, p.RequestBodies)
		}

	case msg.Code == NewBlockMsg:
		// Retrieve and decode the propagated block
		var request newBlockData
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		request.Block.ReceivedAt = msg.ReceivedAt
		request.Block.ReceivedFrom = p

		log.Debug("Received a message[NewBlockMsg]------------", "receiveAt", request.Block.ReceivedAt.Unix(), "peerId", p.id, "hash", request.Block.Hash(), "number", request.Block.NumberU64())

		// Mark the peer as owning the block and schedule it for import
		p.MarkBlock(request.Block.Hash())
		if pm.blockchain.Engine().HasBlock(request.Block.Hash(), request.Block.NumberU64()) {
			return nil
		}
		pm.blockFetcher.Enqueue(p.id, request.Block)

		// Assuming the block is importable by the peer, but possibly not yet done so,
		// calculate the head hash and block number that the peer truly must have.
		var (
			trueHead = request.Block.ParentHash()
			trueBn   = new(big.Int).Sub(request.Block.Number(), big.NewInt(1))
		)
		// Update the peers block number if better than the previous

		if _, bn := p.Head(); trueBn.Cmp(bn) > 0 {
			p.SetHead(trueHead, trueBn)
			pm.chainSync.handlePeerEvent(p)
		}

	case msg.Code == TransactionMsg:
		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		// if txmaker is started,the chain should not accept RemoteTxs,to reduce produce tx cost
		if atomic.LoadUint32(&pm.acceptTxs) == 0 || atomic.LoadUint32(&pm.acceptRemoteTxs) == 1 {
			break
		}
		// Transactions can be processed, parse all of them and deliver to the pool
		var txs []*types.Transaction
		if err := msg.Decode(&txs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		for i, tx := range txs {
			// Validate and mark the remote transaction
			if tx == nil {
				return errResp(ErrDecode, "transaction %d is nil", i)
			}
			p.MarkTransaction(tx.Hash())
		}

		if p.version < eth65 {
			go pm.txpool.AddRemotes(txs)
		} else {
			// PooledTransactions and Transactions are all handled by txFetcher
			return pm.txFetcher.Enqueue(p.id, txs, false)
		}

	case p.version >= eth65 && msg.Code == NewPooledTransactionHashesMsg:
		ann := new(NewPooledTransactionHashesPacket)
		if err := msg.Decode(ann); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Schedule all the unknown hashes for retrieval
		for _, hash := range *ann {
			p.MarkTransaction(hash)
		}
		return pm.txFetcher.Notify(p.id, *ann)

	case p.version >= eth65 && msg.Code == GetPooledTransactionsMsg:
		// Decode the pooled transactions retrieval message
		var query GetPooledTransactionsPacket
		if err := msg.Decode(&query); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		log.Trace("Handler Receive GetPooledTransactions", "peer", p.id, "hashes", len(query))
		hashes, txs := pm.answerGetPooledTransactions(query, p)
		if len(txs) > 0 {
			log.Trace("Handler Send PooledTransactions", "peer", p.id, "txs", len(txs))
			return p.SendPooledTransactionsRLP(hashes, txs)
		}

	case p.version >= eth65 && msg.Code == PooledTransactionsMsg:
		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		// if txmaker is started,the chain should not accept RemoteTxs,to reduce produce tx cost
		if atomic.LoadUint32(&pm.acceptTxs) == 0 || atomic.LoadUint32(&pm.acceptRemoteTxs) == 1 {
			break
		}
		// Transactions can be processed, parse all of them and deliver to the pool
		var txs PooledTransactionsPacket
		if err := msg.Decode(&txs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		for i, tx := range txs {
			// Validate and mark the remote transaction
			if tx == nil {
				return errResp(ErrDecode, "transaction %d is nil", i)
			}
			p.MarkTransaction(tx.Hash())
		}
		log.Trace("Handler Receive PooledTransactions", "peer", p.id, "txs", len(txs))
		return pm.txFetcher.Enqueue(p.id, txs, true)

	default:
		return errResp(ErrInvalidMsgCode, "%v", msg.Code)
	}
	return nil
}

// BroadcastBlock will either propagate a block to a subset of its peers, or
// will only announce its availability (depending what's requested).
func (pm *ProtocolManager) BroadcastBlock(block *types.Block, propagate bool) {
	hash := block.Hash()
	peers := pm.peers.PeersWithoutBlock(hash)
	//var peers []*peer
	//if _, ok := pm.engine.(consensus.Bft); ok {
	//	peers = pm.peers.PeersWithoutConsensus(pm.engine)
	//} else {
	//	peers = pm.peers.PeersWithoutBlock(hash)
	//}

	// If propagation is requested, send to a subset of the peer
	if propagate {
		// Calculate the TD of the block (it's not imported yet, so block.Td is not valid)
		if parent := pm.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent != nil {
		} else {
			log.Warn("Propagating dangling block", "number", block.Number(), "hash", hash)
			return
		}

		var transfer []*peer
		if len(peers) <= numBroadcastBlockPeers {
			// Send the block to all peers
			transfer = peers
		} else {
			// Send the block to a subset of our peers
			rd := rand.New(rand.NewSource(time.Now().UnixNano()))
			indexes := rd.Perm(len(peers))
			maxPeers := int(math.Sqrt(float64(len(peers))))
			transfer = make([]*peer, 0, maxPeers)
			for i := 0; i < maxPeers; i++ {
				transfer = append(transfer, peers[indexes[i]])
			}
		}
		for _, peer := range transfer {
			peer.AsyncSendNewBlock(block)
		}
		log.Trace("Propagated block", "hash", hash, "recipients", len(transfer), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
		return
	}
	// Otherwise if the block is indeed in out own chain, announce it
	if pm.blockchain.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			peer.AsyncSendNewBlockHash(block)
		}
		log.Trace("Announced block", "hash", hash, "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
	}
}

// BroadcastTxs will propagate a batch of transactions to all peers which are not known to
// already have the given transaction.
func (pm *ProtocolManager) BroadcastTxs(txs types.Transactions) {
	var (
		annoCount   int // Count of announcements made
		annoPeers   int
		directCount int // Count of the txs sent directly to peers
		directPeers int // Count of the peers that were sent transactions directly

		txset = make(map[*peer][]common.Hash) // Set peer->transaction to transfer directly
		annos = make(map[*peer][]common.Hash) // Set peer->hash to announce

	)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Broadcast transactions to a batch of peers not knowing about it
	for _, tx := range txs {
		peers := pm.peers.PeersWithoutTx(tx.Hash())
		if len(peers) <= numBroadcastTxPeers {
			for _, peer := range peers {
				txset[peer] = append(txset[peer], tx.Hash())
			}
		} else {
			indexes := rd.Perm(len(peers))
			numAnnos := int(math.Sqrt(float64(len(peers) - numBroadcastTxPeers)))
			countAnnos := 0
			if numAnnos > numBroadcastTxHashPeers {
				numAnnos = numBroadcastTxHashPeers
			}
			for i, c := 0, 0; i < len(peers) && countAnnos < numAnnos; i, c = i+1, c+1 {
				peer := peers[indexes[i]]
				if c < numBroadcastTxPeers {
					txset[peer] = append(txset[peer], tx.Hash())
				} else {
					// For the remaining peers, send announcement only
					annos[peer] = append(annos[peer], tx.Hash())
					countAnnos++
				}
			}
		}
	}

	for peer, txs := range txset {
		directPeers++
		directCount += len(txs)
		peer.AsyncSendTransactions(txs)
	}
	for peer, hashes := range annos {
		annoPeers++
		annoCount += len(hashes)
		if peer.version >= eth65 {
			peer.AsyncSendPooledTransactionHashes(hashes)
		}
	}
	log.Trace("Transaction broadcast", "txs", len(txs),
		"transaction packs", directPeers, "broadcast transaction", directCount,
		"announce packs", annoPeers, "announced hashes", annoCount)
}

func (pm *ProtocolManager) answerGetPooledTransactions(query GetPooledTransactionsPacket, peer *peer) ([]common.Hash, []rlp.RawValue) {
	// Gather transactions until the fetch or network limits is reached
	var (
		bytes  int
		hashes []common.Hash
		txs    []rlp.RawValue
	)
	for _, hash := range query {
		if bytes >= softResponseLimit {
			break
		}
		// Retrieve the requested transaction, skipping if unknown to us
		tx := pm.txpool.Get(hash)
		if tx == nil {
			continue
		}
		// If known, encode and queue for response packet
		if encoded, err := rlp.EncodeToBytes(tx); err != nil {
			log.Error("Failed to encode transaction", "err", err)
		} else {
			hashes = append(hashes, hash)
			txs = append(txs, encoded)
			bytes += len(encoded)
		}
	}
	return hashes, txs
}

// minedBroadcastLoop sends mined blocks to connected peers.
func (pm *ProtocolManager) minedBroadcastLoop() {
	defer pm.wg.Done()

	for obj := range pm.minedBlockSub.Chan() {
		if ev, ok := obj.Data.(core.NewMinedBlockEvent); ok {
			pm.BroadcastBlock(ev.Block, true)  // First propagate block to peers
			pm.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}

}

func (pm *ProtocolManager) txBroadcastLoop() {
	defer pm.wg.Done()
	timer := time.NewTimer(defaultBroadcastInterval)

	for {
		select {
		case event := <-pm.txsCh:
			pm.txsCache = append(pm.txsCache, event.Txs...)
			if len(pm.txsCache) >= defaultTxsCacheSize {
				//log.Trace("broadcast txs", "count", len(pm.txsCache))
				pm.BroadcastTxs(pm.txsCache)
				pm.txsCache = make([]*types.Transaction, 0)
				timer.Reset(defaultBroadcastInterval)
			}
		case <-timer.C:
			if len(pm.txsCache) > 0 {
				//log.Trace("broadcast txs", "count", len(pm.txsCache))
				pm.BroadcastTxs(pm.txsCache)
				pm.txsCache = make([]*types.Transaction, 0)
			}
			timer.Reset(defaultBroadcastInterval)

			// Err() channel will be closed when unsubscribing.
		case <-pm.txsSub.Err():
			return
		}
	}
}

// NodeInfo represents a short summary of the PlatON sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network uint64              `json:"network"` // PlatON network ID (1=Frontier, 2=Morden, Ropsten=3, Rinkeby=4)
	Genesis common.Hash         `json:"genesis"` // SHA3 hash of the host's genesis block
	Config  *params.ChainConfig `json:"config"`  // Chain configuration for the fork rules
	Head    common.Hash         `json:"head"`    // SHA3 hash of the host's best owned block
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (pm *ProtocolManager) NodeInfo() *NodeInfo {
	currentBlock := pm.blockchain.CurrentBlock()
	return &NodeInfo{
		Network: pm.networkID,
		Genesis: pm.blockchain.Genesis().Hash(),
		Config:  pm.blockchain.Config(),
		Head:    currentBlock.Hash(),
	}
}
