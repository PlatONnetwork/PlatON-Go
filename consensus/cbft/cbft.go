// Package bft implements the BFT consensus engine.
package cbft

import (
	"bytes"
	"container/list"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"math/big"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
	lru "github.com/hashicorp/golang-lru"
)

var (
	errSign = errors.New("sign error")
	//errUnauthorizedSigner  = errors.New("unauthorized signer")
	//errIllegalBlock        = errors.New("illegal block")
	//lateBlock              = errors.New("block is late")
	//errDuplicatedBlock     = errors.New("duplicated block")
	//errBlockNumber         = errors.New("error block number")
	errUnknownBlock = errors.New("unknown block")
	errFutileBlock  = errors.New("futile block")
	//errGenesisBlock        = errors.New("cannot handle genesis block")
	//errHighestLogicalBlock = errors.New("cannot find a logical block")
	//errListConfirmedBlocks = errors.New("list confirmed blocks error")
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")

	errInitiateViewchange          = errors.New("not initiated viewChange")
	errTwoThirdViewchangeVotes     = errors.New("lower two third viewChangeVotes")
	errTwoThirdPrepareVotes        = errors.New("lower two third prepare prepareVotes")
	errNotFoundViewBlock           = errors.New("not found block")
	errInvalidViewChangeVotes      = errors.New("invalid prepare viewChangeVotes")
	errInvalidPrepareVotes         = errors.New("invalid prepare prepareVotes")
	errInvalidatorCandidateAddress = errors.New("invalid address")
	errDuplicationConsensusMsg     = errors.New("duplication message")

	//errInvalidVrfProve = errors.New("Invalid vrf prove")
	extraSeal = 65
	//windowSize         = 10

	//periodMargin is a percentum for period margin
	//periodMargin = uint64(20)

	//maxPingLatency is the time in milliseconds between Ping and Pong
	maxPingLatency = int64(5000)

	//maxAvgLatency is the time in milliseconds between two peers
	maxAvgLatency = int64(2000)

	maxResetCacheSize = 512

	// lastBlockOffsetMs is the offset in milliseconds for the last block deadline
	// calculate. (200ms)
	//lastBlockOffsetMs = 200 * time.Millisecond

	peerMsgQueueSize = 1024
	cbftVersion      = byte(0x01)

	maxBlockDist = uint64(192)

	maxQueuesLimit = 4096
)

func NewFaker() consensus.Engine {
	return new(consensus.BftMock)
}

type Cbft struct {
	config      *params.CbftConfig
	eventMux    *event.TypeMux
	handler     handler
	closeOnce   sync.Once
	exitCh      chan struct{}
	txPool      *core.TxPool
	blockChain  *core.BlockChain //the block chain
	running     int32
	peerMsgCh   chan *MsgInfo
	syncBlockCh chan *BlockExt

	executeBlockCh          chan *ExecuteBlockStatus
	baseBlockCh             chan chan *types.Block
	sealBlockCh             chan *SealBlock
	getBlockCh              chan *GetBlock
	innerUnExecutedBlockCh  chan []*BlockExt
	shouldSealCh            chan chan error
	viewChangeTimeoutCh     chan *viewChange
	viewChangeVoteTimeoutCh chan *viewChangeVote
	blockChainCache         *core.BlockChainCache
	hasBlockCh              chan *HasBlock
	statusCh                chan chan string
	getBlockByHashCh        chan *GetBlock
	fastSyncCommitHeadCh    chan chan error
	needPending             bool
	RoundState

	netLatencyMap  map[discover.NodeID]*list.List
	netLatencyLock sync.RWMutex
	//todo maybe discarded struct
	dataReceiveCh chan interface{} //a channel to receive data from miner

	//maybe
	signedSet map[uint64]struct{} //all block numbers signed by local node

	log        log.Logger
	resetCache *lru.Cache
	bp         Breakpoint
	// router
	router *router
	queues map[string]int

	queueMu sync.RWMutex

	// wal
	nodeServiceContext *node.ServiceContext
	wal                Wal
	loading            int32

	// validator
	agency     Agency
	validators atomic.Value //*Validators

	startTimeOfEpoch int64

	evPool  EvidencePool
	tracing *tracing
}

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, eventMux *event.TypeMux, ctx *node.ServiceContext) *Cbft {
	cbft := &Cbft{
		config:                  config,
		eventMux:                eventMux,
		running:                 1,
		exitCh:                  make(chan struct{}),
		signedSet:               make(map[uint64]struct{}),
		syncBlockCh:             make(chan *BlockExt, peerMsgQueueSize),
		peerMsgCh:               make(chan *MsgInfo, peerMsgQueueSize),
		executeBlockCh:          make(chan *ExecuteBlockStatus),
		baseBlockCh:             make(chan chan *types.Block),
		sealBlockCh:             make(chan *SealBlock),
		getBlockCh:              make(chan *GetBlock),
		innerUnExecutedBlockCh:  make(chan []*BlockExt, peerMsgQueueSize),
		shouldSealCh:            make(chan chan error, peerMsgQueueSize),
		viewChangeTimeoutCh:     make(chan *viewChange),
		viewChangeVoteTimeoutCh: make(chan *viewChangeVote),
		hasBlockCh:              make(chan *HasBlock, peerMsgQueueSize),
		statusCh:                make(chan chan string, peerMsgQueueSize),
		getBlockByHashCh:        make(chan *GetBlock),
		fastSyncCommitHeadCh:    make(chan chan error),
		netLatencyMap:           make(map[discover.NodeID]*list.List),
		log:                     log.New(),
		nodeServiceContext:      ctx,
	}

	evPool, err := NewEvidencePoolByCtx(ctx)
	if err != nil {
		return nil
	}
	cbft.evPool = evPool
	cbft.bp = defaultBP
	cbft.handler = NewHandler(cbft)
	cbft.router = NewRouter(cbft, cbft.handler)
	cbft.queues = make(map[string]int)
	cbft.resetCache, _ = lru.New(maxResetCacheSize)
	cbft.tracing = NewTracing()
	return cbft
}

func (cbft *Cbft) getValidators() *cbfttypes.Validators {
	if v := cbft.validators.Load(); v == nil {
		panic("Get validators fail")
	} else {
		return v.(*cbfttypes.Validators)
	}
}

func (cbft *Cbft) ReceivePeerMsg(msg *MsgInfo) {
	select {
	case cbft.peerMsgCh <- msg:
		cbft.log.Debug("Received message from peer", "peer", msg.PeerID.TerminalString(), "msgType", reflect.TypeOf(msg.Msg), "msgHash", msg.Msg.MsgHash().TerminalString(), "BHash", msg.Msg.BHash().TerminalString())
	case <-cbft.exitCh:
		cbft.log.Error("Cbft exit")
	}
}

func (cbft *Cbft) InsertChain(block *types.Block, syncState chan error) {
	var extra *BlockExtra
	var err error

	if _, extra, err = cbft.decodeExtra(block.ExtraData()); err != nil {
		if syncState != nil {
			syncState <- err
		}
		return
	}
	ext := NewBlockExt(block, block.NumberU64(), cbft.nodeLength())
	for _, vote := range extra.Prepare {
		ext.prepareVotes.Add(vote)
	}

	ext.view = extra.ViewChange
	ext.viewChangeVotes = extra.ViewChangeVotes
	ext.syncState = syncState
	cbft.log.Debug("Insert new block", "hash", block.Hash(), "number", block.NumberU64(), "view", ext.view.String())

	cbft.syncBlockCh <- ext
}

// SetPrivateKey sets local's private key by the backend.go
func (cbft *Cbft) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	cbft.config.PrivateKey = privateKey
	cbft.config.NodeID = discover.PubkeyID(&privateKey.PublicKey)
}

func (cbft *Cbft) SetBlockChainCache(blockChainCache *core.BlockChainCache) {
	cbft.blockChainCache = blockChainCache
}

func (cbft *Cbft) SetBreakpoint(t string, log string) error {
	if bp, err := getBreakpoint(t, log); err != nil {
		return err
	} else {
		cbft.bp = bp
	}
	return nil
}

// Start sets blockChain and txPool into cbft
func (cbft *Cbft) Start(blockChain *core.BlockChain, txPool *core.TxPool, agency Agency) error {
	cbft.blockChain = blockChain
	cbft.startTimeOfEpoch = int64(blockChain.Genesis().Time().Uint64())

	cbft.agency = agency

	currentBlock := blockChain.CurrentBlock()

	validators, err := cbft.agency.GetValidator(currentBlock.NumberU64())
	if err != nil {
		cbft.log.Error("Get validator fail", "error", err)
		return err
	}
	cbft.validators.Store(validators)

	genesisParentHash := bytes.Repeat([]byte{0x00}, 32)
	if bytes.Equal(currentBlock.ParentHash().Bytes(), genesisParentHash) && currentBlock.Number() == nil {
		currentBlock.Header().Number = big.NewInt(0)
	}

	cbft.log.Debug("Init highestLogicalBlock", "hash", currentBlock.Hash(), "number", currentBlock.NumberU64())

	current := NewBlockExtBySeal(currentBlock, currentBlock.NumberU64(), cbft.nodeLength())
	current.number = currentBlock.NumberU64()

	if current.number > 0 && cbft.getValidators().Len() > 1 {
		var extra *BlockExtra

		if _, extra, err = cbft.decodeExtra(current.block.ExtraData()); err != nil {
			cbft.log.Error("Block extra decode fail", "error", err)
			return err
		}
		current.view = extra.ViewChange

		for _, vote := range extra.Prepare {
			current.timestamp = vote.Timestamp
			current.prepareVotes.Add(vote)
		}

	}

	cbft.localHighestPrepareVoteNum = current.number

	cbft.blockExtMap = NewBlockExtMap(current, cbft.getThreshold())
	cbft.saveBlockExt(currentBlock.Hash(), current)

	cbft.highestConfirmed.Store(current)
	cbft.highestLogical.Store(current)

	cbft.rootIrreversible.Store(current)

	cbft.txPool = txPool
	cbft.init()

	// init wal and load wal journal
	cbft.wal = &emptyWal{}
	if cbft.config.WalEnabled {
		if cbft.wal, err = NewWal(cbft.nodeServiceContext, ""); err != nil {
			return err
		}
	}
	atomic.StoreInt32(&cbft.loading, 1)

	go cbft.receiveLoop()
	go cbft.executeBlockLoop()
	//start receive cbft message
	log.Debug("handler.Start")
	go cbft.handler.Start()
	go cbft.update()

	if err = cbft.wal.Load(cbft.AddJournal); err != nil {
		return err
	}
	atomic.StoreInt32(&cbft.loading, 0)
	return nil
}

// schedule is responsible for HighestPrepareBlock synchronization
func (cbft *Cbft) scheduleHighestPrepareBlock() {
	schedule := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-schedule.C:
			cbft.handler.SendPartBroadcast(&getHighestPrepareBlock{Lowest: cbft.getRootIrreversible().number + 1})
		}
	}
}

func (cbft *Cbft) receiveLoop() {
	for {
		select {
		case msg := <-cbft.peerMsgCh:
			count := cbft.queues[msg.PeerID.TerminalString()] + 1
			if count > maxQueuesLimit {
				cbft.log.Debug("Discarded msg, exceeded allowance", "peer", msg.PeerID.TerminalString(), "msgType", reflect.TypeOf(msg.Msg), "limit", maxQueuesLimit)
				break
			}
			cbft.queues[msg.PeerID.TerminalString()] = count
			cbft.handleMsg(msg)
			cbft.queues[msg.PeerID.TerminalString()]--
			if cbft.queues[msg.PeerID.TerminalString()] == 0 {
				delete(cbft.queues, msg.PeerID.TerminalString())
			}

		case bt := <-cbft.syncBlockCh:
			cbft.OnSyncBlock(bt)
		case bs := <-cbft.executeBlockCh:
			cbft.OnExecutedBlock(bs)
		case shouldSeal := <-cbft.shouldSealCh:
			cbft.OnShouldSeal(shouldSeal)
		case view := <-cbft.viewChangeTimeoutCh:
			cbft.OnViewChangeTimeout(view)
		case viewVote := <-cbft.viewChangeVoteTimeoutCh:
			cbft.OnViewChangeVoteTimeout(viewVote)
		case sealBlock := <-cbft.sealBlockCh:
			cbft.OnSeal(sealBlock.block, sealBlock.sealResultCh, sealBlock.stopCh)
		case block := <-cbft.getBlockCh:
			cbft.OnGetBlock(block.hash, block.number, block.ch)
		case baseBlock := <-cbft.baseBlockCh:
			cbft.OnBaseBlock(baseBlock)
		case hasBlock := <-cbft.hasBlockCh:
			cbft.OnHasBlock(hasBlock)
		case status := <-cbft.statusCh:
			cbft.OnStatus(status)
		case block := <-cbft.getBlockByHashCh:
			cbft.OnGetBlockByHash(block.hash, block.ch)
		case fastSync := <-cbft.fastSyncCommitHeadCh:
			cbft.OnFastSyncCommitHead(fastSync)
		}
	}
}

func (cbft *Cbft) handleMsg(info *MsgInfo) {
	msg, peerID := info.Msg, info.PeerID
	var err error

	// record the message for received
	cbft.tracing.RecordReceive(cbft.config.NodeID.TerminalString(),
		info.PeerID.TerminalString(),
		info.Msg.MsgHash().TerminalString(),
		fmt.Sprintf("%T", info.Msg))

	if !cbft.isRunning() {
		switch msg.(type) {
		case *prepareBlock,
		*prepareBlockHash,
		*prepareVote,
		*viewChange,
		*viewChangeVote:
			cbft.log.Debug("Cbft is not running, discard consensus message")
			return
		}
	}
	isWriteWal := true
	switch msg := msg.(type) {
	case *prepareBlock:
		err = cbft.OnNewPrepareBlock(peerID, msg, true)
	case *prepareVote:
		err = cbft.OnPrepareVote(peerID, msg, true)
	case *viewChange:
		err = cbft.OnViewChange(peerID, msg)
	case *viewChangeVote:
		err = cbft.OnViewChangeVote(peerID, msg)
	case *confirmedPrepareBlock:
		err = cbft.OnConfirmedPrepareBlock(peerID, msg)
	case *getPrepareBlock:
		err = cbft.OnGetPrepareBlock(peerID, msg)
	case *getPrepareVote:
		err = cbft.OnGetPrepareVote(peerID, msg)
	case *prepareVotes:
		err = cbft.OnPrepareVotes(peerID, msg)
	case *getHighestPrepareBlock:
		err = cbft.OnGetHighestPrepareBlock(peerID, msg)
	case *highestPrepareBlock:
		err = cbft.OnHighestPrepareBlock(peerID, msg)
	case *prepareBlockHash:
		err = cbft.OnPrepareBlockHash(peerID, msg)
	case *getLatestStatus:
		err = cbft.OnGetLatestStatus(peerID, msg)
		isWriteWal = false
	case *latestStatus:
		err = cbft.OnLatestStatus(peerID, msg)
		isWriteWal = false
	}

	if err != nil {
		cbft.log.Error("Handle msg Failed", "error", err, "type", reflect.TypeOf(msg), "peer", peerID)
	} else if !cbft.isLoading() && isWriteWal {
		// write journal msg if cbft is not loading
		cbft.wal.Write(info)
	}
}

func (cbft *Cbft) isRunning() bool {
	return atomic.LoadInt32(&cbft.running) == 1
}

func (cbft *Cbft) isLoading() bool {
	return atomic.LoadInt32(&cbft.loading) == 1
}

func (cbft *Cbft) OnShouldSeal(shouldSeal chan error) {
	//clear all invalid data
	if len(cbft.shouldSealCh) > 0 {
		for {
			select {
			case shouldSeal = <-cbft.shouldSealCh:
			default:
				goto END
			}
		}
	}
END:
	if cbft.hadSendViewChange() {
		validator, err := cbft.getValidators().NodeIndexAddress(cbft.config.NodeID)

		if err != nil {
			log.Debug("Get node index and address failed", "error", err)
			shouldSeal <- err
			return
		}

		// May be currently only me produce block, so current view
		// should be invalid, making a new one.
		if !cbft.validViewChange() && (time.Now().Unix()-int64(cbft.viewChange.Timestamp)) > cbft.config.Duration {
			// need send viewchange
			cbft.OnSendViewChange()

			oldCount := viewChangeGauage.Value()
			viewChangeGauage.Update(oldCount + 1)
			viewChangeCounter.Inc(1)

			shouldSeal <- errTwoThirdViewchangeVotes
			return
		}

		if cbft.isRunning() && cbft.agreeViewChange() &&
			cbft.viewChange.ProposalAddr == validator.Address &&
			uint32(validator.Index) == cbft.viewChange.ProposalIndex {
			// do something check
			shouldSeal <- nil
		} else {
			shouldSeal <- errInitiateViewchange
		}
	} else {
		// need send viewchange
		cbft.OnSendViewChange()

		oldCount := viewChangeGauage.Value()
		viewChangeGauage.Update(oldCount + 1)
		viewChangeCounter.Inc(1)

		shouldSeal <- errTwoThirdViewchangeVotes
	}
}

//Sync block from p2p
//BlockNum is must higher current highest confirmed block
//If verify block success, cbft will change sync mode
//It don't process new PrepareBlock, PrepareVote, Viewchange, ViewchangeVote
//It stop sync mode when new view's parent block is in memory. Now it had catch up other peer
func (cbft *Cbft) OnSyncBlock(ext *BlockExt) {
	cbft.bp.SyncBlockBP().SyncBlock(context.TODO(), ext, cbft)
	//todo verify block
	if ext.block.NumberU64() < cbft.getHighestConfirmed().number {
		cbft.log.Debug("Sync block too lower", "hash", ext.block.Hash(), "number", ext.number, "highestNumber", cbft.getHighestConfirmed().number, "highestHash", cbft.getHighestConfirmed().block.Hash(), "root", cbft.getRootIrreversible().number)
		ext.SetSyncState(nil)
		cbft.bp.SyncBlockBP().InvalidBlock(context.TODO(), ext, fmt.Errorf("sync block too lower"), cbft)
		return
	}

	if cbft.blockExtMap.findBlock(ext.block.Hash(), ext.block.NumberU64()) != nil {
		cbft.log.Debug("Sync block had exist", "hash", ext.block.Hash(), "number", ext.number, "highest", cbft.getHighestConfirmed().number, "root", cbft.getRootIrreversible().number)
		ext.SetSyncState(nil)
		cbft.bp.SyncBlockBP().InvalidBlock(context.TODO(), ext, fmt.Errorf("sync block had exist"), cbft)
		return
	}

	cbft.log.Debug("Sync block success", "hash", ext.block.Hash(), "number", ext.number)

	if (cbft.viewChange != nil && !cbft.viewChange.Equal(ext.view)) || !cbft.agreeViewChange() {
		cbft.viewChange = ext.view
		if len(ext.viewChangeVotes) >= cbft.getThreshold() {
			if err := cbft.checkViewChangeVotes(ext.viewChangeVotes); err != nil {
				log.Error("Receive prepare invalid block", "err", err)
				cbft.bp.SyncBlockBP().InvalidBlock(context.TODO(), ext, err, cbft)
				ext.SetSyncState(err)
				return
			}
			for _, v := range ext.viewChangeVotes {
				cbft.viewChangeVotes[v.ValidatorAddr] = v
			}

			cbft.clearPending()
			cbft.ClearChildren(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum, cbft.viewChange.Timestamp)
			cbft.producerBlocks = NewProducerBlocks(cbft.getValidators().NodeID(int(ext.view.ProposalIndex)), ext.block.NumberU64())
			if cbft.producerBlocks != nil {
				cbft.producerBlocks.AddBlock(ext.block)
				cbft.log.Debug("Add producer block", "hash", ext.block.Hash(), "number", ext.block.Number(), "producer", cbft.producerBlocks.String())
			}
		}
	}
	ext.timestamp = cbft.viewChange.Timestamp
	cbft.OnNewBlock(ext)
}

//Sync confirmed prepare prepareVotes, not sync when local node has enough prepare prepareVotes
func (cbft *Cbft) OnConfirmedPrepareBlock(peerID discover.NodeID, pb *confirmedPrepareBlock) error {
	cbft.log.Debug("Received confirmed prepareBlock ", "peer", peerID, "confirmedPrepareBlock", pb.String())
	ext := cbft.blockExtMap.findBlock(pb.Hash, pb.Number)
	if ext == nil || ext.block == nil {
		cbft.handler.Send(peerID, &getPrepareBlock{Hash: pb.Hash, Number: pb.Number})
	}
	if ext != nil && ext.prepareVotes.Len() < cbft.getThreshold() {
		sub := pb.VoteBits.Sub(ext.prepareVotes.voteBits)
		if !sub.IsEmpty() {
			cbft.handler.Send(peerID, &getPrepareVote{Hash: pb.Hash, Number: pb.Number, VoteBits: sub})
		}
	}

	cbft.syncMissingBlock(peerID, pb.Number)

	if cbft.needBroadcast(peerID, pb) {
		go cbft.handler.SendBroadcast(pb)
	}

	return nil
}

func (cbft *Cbft) syncMissingBlock(peerID discover.NodeID, highest uint64) {
	hs := cbft.blockExtMap.GetHasVoteWithoutBlock(highest)
	for _, h := range hs {
		cbft.handler.Send(peerID, &getPrepareBlock{Hash: h.hash, Number: h.number})
	}

	hs = cbft.blockExtMap.GetWithoutTwoThirdVotes(highest)
	for _, h := range hs {
		cbft.handler.Send(peerID, &getPrepareVote{Hash: h.hash, Number: h.number, VoteBits: h.bits})
	}
}
func (cbft *Cbft) OnGetPrepareBlock(peerID discover.NodeID, g *getPrepareBlock) error {
	ext := cbft.blockExtMap.findBlock(g.Hash, g.Number)
	if ext != nil {
		pb, err := ext.PrepareBlock()
		if err == nil {
			cbft.handler.Send(peerID, pb)
			cbft.log.Debug("Send Block", "peer", peerID, "hash", g.Hash, "number", g.Number, "msgHash", pb.MsgHash().TerminalString())
		}
	}
	return nil
}

func (cbft *Cbft) OnGetPrepareVote(peerID discover.NodeID, pv *getPrepareVote) error {
	ext := cbft.blockExtMap.findBlock(pv.Hash, pv.Number)
	votes := make([]*prepareVote, 0)

	if ext != nil {
		for i := uint32(0); i < pv.VoteBits.Size(); i++ {
			if !pv.VoteBits.GetIndex(i) {
				if v := ext.prepareVotes.Get(i); v != nil {
					votes = append(votes, v)
				}
			}
		}
	} else {
		block := cbft.blockChain.GetBlock(pv.Hash, pv.Number)
		if block != nil {
			_, ed, err := cbft.decodeExtra(block.ExtraData())
			if err == nil {
				votes = append(votes, ed.Prepare...)
			} else {
				cbft.log.Error("Get Block error", "hash", block.Hash(), "number", block.NumberU64(), "err", err)
			}
		}
	}
	if len(votes) != 0 {
		cbft.handler.Send(peerID, &prepareVotes{Hash: pv.Hash, Number: pv.Number, Votes: votes})
		cbft.log.Debug("Send PrepareVotes", "peer", peerID, "hash", pv.Hash, "number", pv.Number)
	}
	return nil
}

func (cbft *Cbft) OnPrepareVotes(peerID discover.NodeID, view *prepareVotes) error {
	for _, vote := range view.Votes {
		if err := cbft.OnPrepareVote(peerID, vote, false); err != nil {
			cbft.log.Error("Handle PrepareVotes failed", "peer", peerID, "err", err)
			return err
		}
	}
	return nil
}

func (cbft *Cbft) OnGetHighestPrepareBlock(peerID discover.NodeID, msg *getHighestPrepareBlock) error {
	highest := cbft.getHighestLogical().number
	commit := cbft.getRootIrreversible().number
	commitedBlock := make([]*types.Block, 0)
	unconfirmedBlock := make([]*prepareBlock, 0)
	votes := make([]*prepareVotes, 0)
	cbft.log.Debug("Receive GetHighestPrepareBlock", "peer", peerID.TerminalString(), "msg", msg.String())
	if commit > msg.Lowest && commit-msg.Lowest > maxBlockDist {
		log.Debug("Discard GetHighestPrepareBlock msg, too far away", "peer", peerID.TerminalString(), "lowest", msg.Lowest, "root", commit)
		return errors.New("peer's block too far away")
	}
	for i := msg.Lowest; i <= commit; i++ {
		if b := cbft.blockChain.GetBlockByNumber(i); b != nil {
			commitedBlock = append(commitedBlock, b)
		}
	}

	exts := cbft.blockExtMap.findBlockExtByNumber(commit+1, highest)
	for _, ext := range exts {
		if prepare, err := ext.PrepareBlock(); err == nil {
			unconfirmedBlock = append(unconfirmedBlock, prepare)
			votes = append(votes, &prepareVotes{Hash: ext.block.Hash(), Number: ext.number, Votes: ext.Votes()})
		} else {
			cbft.log.Error("Retrieve block fail", "err", err)
		}
	}
	cbft.log.Debug("Send highestPrepareBlock")
	cbft.handler.Send(peerID, &highestPrepareBlock{
		CommitedBlock:    commitedBlock,
		UnconfirmedBlock: unconfirmedBlock,
		Votes:            votes,
	})
	return nil
}

func (cbft *Cbft) OnHighestPrepareBlock(peerID discover.NodeID, msg *highestPrepareBlock) error {
	cbft.log.Debug("Receive HighestPrepareBlock", "peer", peerID.TerminalString(), "msg", msg.String())

	if len(msg.CommitedBlock) > int(maxBlockDist) {
		cbft.log.Debug("Discard HighestPrepareBlock msg, exceeded allowance", "peer", peerID.TerminalString(), "CommitedBlock", len(msg.CommitedBlock), "limited", maxBlockDist)
		atomic.StoreInt32(&cbft.running, 0)
		return errors.New("exceeded allowance")
	}

	if len(msg.CommitedBlock)+len(cbft.syncBlockCh) < cap(cbft.syncBlockCh) {
		var largestNum int64 = 0
		for _, block := range msg.CommitedBlock {
			cbft.log.Debug("Sync Highest Block", "number", block.NumberU64())
			cbft.InsertChain(block, nil)
			if block.Number().Int64() > largestNum {
				largestNum = block.Number().Int64()
			}
		}
		if largestNum != 0 {
			p, err := cbft.handler.PeerSet().Get(peerID.TerminalString())
			if err == nil {
				p.SetConfirmedHighestBn(new(big.Int).SetInt64(largestNum))
			}
		}
	}

	for _, prepare := range msg.UnconfirmedBlock {
		cbft.log.Debug("Sync Highest Block", "number", prepare.Block.NumberU64())
		cbft.OnNewPrepareBlock(peerID, prepare, false)
	}
	for _, votes := range msg.Votes {
		cbft.log.Debug("Sync Highest Block", "number", votes.Number)
		cbft.OnPrepareVotes(peerID, votes)
	}
	atomic.StoreInt32(&cbft.running, 1)
	return nil
}

func (cbft *Cbft) OnViewChangeTimeout(view *viewChange) {
	//todo viewchange timeout

	//cbft.mux.Lock()
	//defer cbft.mux.Lock()
	cbft.log.Debug(fmt.Sprintf("Check view change timeout send:%v agree:%v msgHash:%v", cbft.hadSendViewChange(), cbft.agreeViewChange(), view.MsgHash().TerminalString()))
	if cbft.viewChange != nil && view.Equal(cbft.viewChange) {
		if cbft.hadSendViewChange() && !cbft.agreeViewChange() {
			cbft.handleCache()
			cbft.log.Info("View change timeout", "current view", cbft.viewChange.String(), "msgHash", view.MsgHash().TerminalString())
			cbft.resetViewChange()
			cbft.bp.ViewChangeBP().ViewChangeTimeout(context.TODO(), view, cbft)
		}
	}

	viewChangeTimeoutMeter.Mark(1)
}

//Current view change timeout
//Need reset view
func (cbft *Cbft) OnViewChangeVoteTimeout(view *viewChangeVote) {
	//todo viewchange vote timeout
	//cbft.mux.Lock()
	//defer cbft.mux.Lock()
	if cbft.viewChange != nil && view.EqualViewChange(cbft.viewChange) {
		if !cbft.agreeViewChange() {
			cbft.log.Warn("Waiting master response timeout", "view", cbft.viewChange.String())
			cbft.handleCache()
			cbft.resetViewChange()
			cbft.needPending = true
			cbft.handler.SendPartBroadcast(&getHighestPrepareBlock{cbft.getHighestConfirmed().number})
		}
	}

	viewChangeVoteTimeoutMeter.Mark(1)
}

func (cbft *Cbft) OnPrepareBlockHash(peerID discover.NodeID, msg *prepareBlockHash) error {
	cbft.log.Debug("Received message of prepareBlockHash", "FromPeerId", peerID.TerminalString(),
		"BlockHash", msg.Hash.TerminalString(), "Number", msg.Number)
	// Prerequisite: Nodes with PrepareBlock data can forward Hash
	if cbft.blockExtMap.findBlock(msg.Hash, msg.Number) == nil {
		cbft.handler.Send(peerID, &getPrepareBlock{Hash: msg.Hash, Number: msg.Number})
	}

	// then: to forward msg
	if ok := cbft.needBroadcast(peerID, msg); ok {
		go cbft.handler.SendBroadcast(msg)
	}

	return nil
}

func (cbft *Cbft) OnGetLatestStatus(peerID discover.NodeID, msg *getLatestStatus) error {
	cbft.log.Debug("Received message of getHighestConfirmedStatus", "FromPeerId", peerID.TerminalString(), "Number", msg.Highest, "Type", msg.Type, "msgHash", msg.MsgHash().TerminalString())
	curConfirmedNum, curLogicNum := cbft.getHighestConfirmed().number, cbft.getHighestLogical().number
	if msg.Type == HIGHEST_CONFIRMED_BLOCK {
		if curConfirmedNum < msg.Highest {
			p, err := cbft.handler.GetPeer(peerID.TerminalString())
			if err != nil {
				cbft.log.Error("Failed to get peerID", "peerID", peerID.TerminalString())
				return err
			}
			p.SetConfirmedHighestBn(new(big.Int).SetUint64(msg.Highest))
			cbft.log.Debug("The current confirmed block height is lower than the specified block height", "current", curConfirmedNum, "specified", msg.Highest)
			cbft.handler.Send(peerID, &getHighestPrepareBlock{Lowest: cbft.getRootIrreversible().number + 1})
		} else {
			cbft.log.Debug("Current confirmed highest larger and make reply highestConfirmedStatus msg", "highest", msg.Highest, "currentNum", curConfirmedNum)
			cbft.handler.Send(peerID, &latestStatus{Highest: curConfirmedNum, Type: msg.Type})
		}
	}
	if msg.Type == HIGHEST_LOGIC_BLOCK {
		if curLogicNum < msg.Highest {
			p, err := cbft.handler.GetPeer(peerID.TerminalString())
			if err != nil {
				cbft.log.Error("Failed to get peerID", "peerID", peerID.TerminalString())
				return err
			}
			p.SetLogicHighestBn(new(big.Int).SetUint64(msg.Highest))
			cbft.log.Debug("The current logic block height is lower than the specified block height", "current", curLogicNum, "specified", msg.Highest)
			cbft.syncMissingBlock(peerID, msg.Highest)
		} else {
			cbft.log.Debug("Current logic highest larger and make reply highestConfirmedStatus msg", "highest", msg.Highest, "currentNum", curLogicNum)
			cbft.handler.Send(peerID, &latestStatus{Highest: curConfirmedNum, Type: msg.Type})
		}
	}
	return nil
}

func (cbft *Cbft) OnLatestStatus(peerID discover.NodeID, msg *latestStatus) error {
	cbft.log.Debug("Received message of highestConfirmedStatus", "FromPeerId", peerID.TerminalString(), "Number", msg.Highest, "Type", msg.Type, "msgHash", msg.MsgHash().TerminalString())
	curConfirmedNum, curLogicNum := cbft.getHighestConfirmed().number, cbft.getHighestLogical().number
	switch msg.Type {
	case HIGHEST_CONFIRMED_BLOCK:
		if curConfirmedNum < msg.Highest {
			p, err := cbft.handler.GetPeer(peerID.TerminalString())
			if err != nil {
				cbft.log.Error("Failed to get peerID for confirmed highest number", "peerID", peerID.TerminalString())
				return err
			}
			p.SetConfirmedHighestBn(new(big.Int).SetUint64(msg.Highest))
			cbft.log.Debug("The current confirmed block height is lower than the specified block height and getPrepareBlock")
			cbft.handler.Send(peerID, &getHighestPrepareBlock{Lowest: cbft.getRootIrreversible().number + 1})
		}
	case HIGHEST_LOGIC_BLOCK:
		if curLogicNum < msg.Highest {
			p, err := cbft.handler.GetPeer(peerID.TerminalString())
			if err != nil {
				cbft.log.Error("Failed to get peerID for logic highest number", "peerID", peerID.TerminalString())
				return err
			}
			p.SetLogicHighestBn(new(big.Int).SetUint64(msg.Highest))
			cbft.log.Debug("The current logic block height is lower than the specified block height and getPrepareBlock")
			cbft.syncMissingBlock(peerID, msg.Highest)
		}
	}
	return nil
}

func (cbft *Cbft) NextBaseBlock() *types.Block {
	ch := make(chan *types.Block, 1)
	cbft.baseBlockCh <- ch
	return <-ch
}

func (cbft *Cbft) OnBaseBlock(ch chan *types.Block) {
	if cbft.master && cbft.agreeViewChange() && (cbft.producerBlocks == nil || cbft.producerBlocks.Len() == 0) {
		block := cbft.getHighestConfirmed().block
		cbft.log.Debug("Base block", "hash", block.Hash(), "number", block.Number())
		ch <- block
	} else {
		block := cbft.getHighestLogical().block
		cbft.log.Debug("Base block", "hash", block.Hash(), "number", block.Number())
		ch <- block
	}
}

//to sign the block, and store the sign to header.Extra[32:], send the sign to chanel to broadcast to other consensus nodes
func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) error {
	cbft.log.Info("Seal block", "number", block.NumberU64(), "parentHash", block.ParentHash())
	header := block.Header()

	number := block.NumberU64()

	if number == 0 {
		return errUnknownBlock
	}

	// sign the seal hash
	sign, err := cbft.signFn(header.SealHash().Bytes())
	if err != nil {
		log.Error("Seal block sign failed", "err", err)
		return err
	}

	//store the sign in  header.Extra[32:]
	copy(header.Extra[len(header.Extra)-extraSeal:], sign[:])

	sealedBlock := block.WithSeal(header)

	cbft.sealBlockCh <- NewSealBlock(sealedBlock, sealResultCh, stopCh)

	return nil
}

func (cbft *Cbft) OnSeal(sealedBlock *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) {
	if (cbft.getHighestLogical() != nil && !cbft.getHighestLogical().IsParent(sealedBlock.ParentHash())) &&
		(cbft.getHighestConfirmed() != nil && !cbft.getHighestConfirmed().IsParent(sealedBlock.ParentHash())) {
		cbft.log.Warn("Futile block cause highest logical block changed",
			"number", sealedBlock.Number(),
			"parentHash", sealedBlock.ParentHash(),
			"state", cbft.blockState())
		return
	}
	logicNum := cbft.getHighestLogical().number
	if logicNum == sealedBlock.NumberU64() {
		cbft.log.Warn("logicNum must not equal sealedBlock", "logicNum", logicNum, "sealedNum", sealedBlock.NumberU64())
		return
	}

	current := cbft.sealBlockProcess(sealedBlock)

	cbft.bp.InternalBP().Seal(context.TODO(), current, cbft)
	cbft.bp.InternalBP().NewHighestLogicalBlock(context.TODO(), current, cbft)

	cbft.broadcastBlock(current)
	//todo change sign and block state
	go func() {
		select {
		case <-stopCh:
			return
		case sealResultCh <- sealedBlock:
			//reset pool when seal block
			//start := time.Now()
			//cbft.reset(sealedBlock)
			//cbft.bp.InternalBP().ResetTxPool(context.TODO(), current, time.Now().Sub(start), cbft)

		default:
			cbft.log.Warn("Sealing result is not ready by miner", "sealHash", sealedBlock.Header().SealHash())
		}
	}()
}

func (cbft *Cbft) sealBlockProcess(sealedBlock *types.Block) *BlockExt {
	log.Debug("sealBlockProcess", "number", sealedBlock.NumberU64(), "hash", sealedBlock.Hash())
	current := NewBlockExt(sealedBlock, sealedBlock.NumberU64(), cbft.nodeLength())
	//this block is produced by local node, so need not execute in cbft.
	current.view = cbft.viewChange
	current.timestamp = cbft.viewChange.Timestamp
	current.inTree = true
	current.executing = true
	current.isExecuted = true
	current.isSigned = true

	//save the block to cbft.blockExtMap
	cbft.saveBlockExt(sealedBlock.Hash(), current)

	//log this signed block's number
	cbft.signedSet[sealedBlock.NumberU64()] = struct{}{}

	cbft.SetLocalHighestPrepareNum(current.number)
	cbft.reset(sealedBlock)
	if cbft.getValidators().Len() == 1 {
		cbft.log.Info("Seal complete", "hash", sealedBlock.Hash(), "number", sealedBlock.NumberU64())
		cbft.log.Debug("Single node mode, confirm now")
		//only one consensus node, so, each block is highestConfirmed. (lock is needless)
		current.isConfirmed = true
		cbft.highestLogical.Store(current)
		cbft.highestConfirmed.Store(current)
		cbft.flushReadyBlock()
		return current
	}

	//reset cbft.highestLogicalBlockExt cause this block is produced by myself
	cbft.highestLogical.Store(current)
	cbft.AddPrepareBlock(sealedBlock)
	cbft.log.Info("Seal complete", "nodeID", cbft.config.NodeID, "hash", sealedBlock.Hash(), "number", sealedBlock.NumberU64(), "timestamp", sealedBlock.Time(), "producerBlocks", cbft.producerBlocks.Len())
	return current
}

// ShouldSeal checks if it's local's turn to package new block at current time.
func (cbft *Cbft) ShouldSeal(curTime int64) (bool, error) {
	inturn := !cbft.isLoading() && cbft.inTurn(curTime)
	if inturn {
		cbft.netLatencyLock.RLock()
		peersCount := len(cbft.netLatencyMap)
		cbft.netLatencyLock.RUnlock()
		if peersCount < cbft.getThreshold() {
			//inturn = false
		}
	}
	//cbft.log.Debug("Should Seal", "time", curTime, "inturn", inturn, "peers", len(cbft.netLatencyMap))
	if inturn {
		if cbft.lastViewChange != nil {
			lastViewChangeTime := time.Unix(int64(cbft.lastViewChange.Timestamp), 0)
			viewChangeTimer.UpdateSince(lastViewChangeTime)
		}
		// if first block of mine, send viewchange message, return false
		// if viewchange success , return true
		// if viewchange failed , wait timeout until re-send message
		//cbft.mux.Lock()
		//defer cbft.mux.Unlock()
		shouldSeal := make(chan error, 1)
		cbft.shouldSealCh <- shouldSeal
		select {
		case err := <-shouldSeal:
			return err == nil, err
		case <-time.After(2 * time.Millisecond):
			return false, fmt.Errorf("waiting for ShouldSeal timeout")
		}
	}

	return inturn, nil
}

func (cbft *Cbft) OnSendViewChange() {
	view, err := cbft.newViewChange()
	if err != nil {
		cbft.log.Error("New view change failed", "err", err)
		return
	}
	cbft.log.Info("Send new view", "nodeID", cbft.config.NodeID, "view", view.String(), "msgHash", view.MsgHash().TerminalString())
	cbft.bp.ViewChangeBP().SendViewChange(context.TODO(), view, cbft)

	// write new viewChange info to wal journal
	cbft.wal.WriteSync(&MsgInfo{
		Msg:    &sendViewChange{ViewChange: view, Master: true},
		PeerID: cbft.config.NodeID,
	})
	cbft.handler.SendAllConsensusPeer(view)

	// gauage
	blockHighNumConfirmedGauage.Update(int64(cbft.getHighestConfirmed().number))
	blockHighNumLogicGauage.Update(int64(cbft.getHighestLogical().number))

	time.AfterFunc(time.Duration(cbft.config.Period)*time.Second*2, func() {
		cbft.viewChangeTimeoutCh <- view
	})
}

// Receive view from other nodes
// Need verify timestamp , signature, promise highest confirmed block
func (cbft *Cbft) OnViewChange(peerID discover.NodeID, view *viewChange) error {
	cbft.log.Debug("Receive view change", "peer", peerID, "nodeID", cbft.getValidators().NodeID(int(view.ProposalIndex)), "view", view.String(), "msgHash", view.MsgHash().TerminalString())

	if cbft.viewChange != nil && cbft.viewChange.Equal(view) {
		cbft.log.Debug("Duplication view change message, discard this")
		return errDuplicationConsensusMsg
	}

	if view != nil {
		// priority forwarding
		cbft.handler.SendAllConsensusPeer(view)
	}

	bpCtx := context.WithValue(context.Background(), "peer", peerID)
	cbft.bp.ViewChangeBP().ReceiveViewChange(bpCtx, view, cbft)
	if err := cbft.VerifyAndViewChange(view); err != nil {
		if view.BaseBlockNum > cbft.getHighestConfirmed().number {
			if view.BaseBlockNum-cbft.getHighestConfirmed().number > maxBlockDist {
				atomic.StoreInt32(&cbft.running, 0)
			} else {
				cbft.log.Warn(fmt.Sprintf("Local is too slower, need to sync block to %s", peerID.TerminalString()))

				cbft.handler.Send(peerID, &getHighestPrepareBlock{Lowest: cbft.getRootIrreversible().number + 1})
			}
		}

		cbft.bp.ViewChangeBP().InvalidViewChange(bpCtx, view, err, cbft)
		cbft.log.Error("Verify view failed", "err", err, "peer", peerID, "view", view.String(), "local", cbft.viewChange.String())
		return err
	}

	validator, err := cbft.getValidators().NodeIndexAddress(cbft.config.NodeID)
	if err != nil {
		cbft.bp.ViewChangeBP().InvalidViewChange(bpCtx, view, errInvalidatorCandidateAddress, cbft)
		return errInvalidatorCandidateAddress
	}

	resp := &viewChangeVote{
		ValidatorIndex: uint32(validator.Index),
		ValidatorAddr:  validator.Address,
		Timestamp:      view.Timestamp,
		BlockHash:      view.BaseBlockHash,
		BlockNum:       view.BaseBlockNum,
		ProposalIndex:  view.ProposalIndex,
		ProposalAddr:   view.ProposalAddr,
	}

	sign, err := cbft.signMsg(resp)
	if err != nil {
		cbft.log.Error("Signature view vote failed", "err", err)
		return err
	}

	resp.Signature.SetBytes(sign)
	cbft.viewChangeResp = resp
	cbft.log.Info("Response viewChangeVote", "viewChangeResp", resp, "msgHash", resp.MsgHash())
	time.AfterFunc(time.Duration(cbft.config.Period)*time.Second, func() {
		cbft.viewChangeVoteTimeoutCh <- resp
	})
	cbft.agreeViewChangeProcess(view, resp)
	//cbft.viewChangeResp = resp
	//cbft.setViewChange(view)
	//// add to viewChangeVote When the viewChange is approved by self
	//cbft.viewChangeVotes[resp.ValidatorAddr] = resp
	cbft.bp.InternalBP().SwitchView(bpCtx, view, cbft)
	cbft.bp.ViewChangeBP().SendViewChangeVote(bpCtx, resp, cbft)
	cbft.handler.SendAllConsensusPeer(resp)
	return nil
}

func (cbft *Cbft) agreeViewChangeProcess(view *viewChange, viewChangeResp *viewChangeVote) {
	cbft.setViewChange(view)
	cbft.viewChangeResp = viewChangeResp
	// add vote to viewChangeVotes When the viewChange is agree by self
	cbft.viewChangeVotes[viewChangeResp.ValidatorAddr] = viewChangeResp
}

// flushReadyBlock finds ready blocks and flush them to chain
func (cbft *Cbft) flushReadyBlock() bool {
	cbft.log.Debug("Flush to chain", "state", cbft.blockState())
	if cbft.viewChange == nil {
		return false
	}
	//todo verify state
	//todo direct flush block if node is no-consensus node
	if ext := cbft.blockExtMap.findChild(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum); ext == nil || !ext.isConfirmed {
		cbft.log.Debug("No block need flush db", "ext", ext, "viewChange", cbft.viewChange)
		return false
	}

	flush := cbft.blockExtMap.GetSubChainWithTwoThirdVotes(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum)

	if len(flush) == 0 {
		cbft.log.Debug("Enable flushed block is empty")
		return false
	}

	cbft.log.Debug("Flush block", "total", len(flush))
	blockMinedMeter.Mark(int64(len(flush)))
	cbft.storeBlocks(flush)
	chainBlock := cbft.blockChain.CurrentBlock()
	highestBlockHash, highestBlockNum := chainBlock.Hash(), chainBlock.NumberU64()

	cbft.blockExtMap.ClearParents(highestBlockHash, highestBlockNum)
	newRoot := cbft.blockExtMap.BaseBlock(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum)
	cbft.log.Debug("Set new root", "hash", newRoot.block.Hash(), "number", newRoot.block.NumberU64())
	cbft.rootIrreversible.Store(newRoot)
	cbft.bp.InternalBP().NewHighestRootBlock(context.TODO(), newRoot, cbft)

	blockConfirmedTimer.UpdateSince(common.MillisToTime(newRoot.rcvTime))

	cbft.evPool.Clear(cbft.viewChange.Timestamp, cbft.viewChange.BaseBlockNum)
	return true

}

// Receive prepare block from the other consensus node.
// Need check something ,such as validator index, address, view is equal local view , and last verify signature
func (cbft *Cbft) OnNewPrepareBlock(nodeId discover.NodeID, request *prepareBlock, propagation bool) error {
	cbft.log.Info("Received a PrepareBlockMsg", "FromPeerId", nodeId.TerminalString(), "prepare", request.String(), "msgHash", request.MsgHash())
	bpCtx := context.WithValue(context.TODO(), "peer", nodeId)
	cbft.bp.PrepareBP().ReceiveBlock(bpCtx, request, cbft)

	//discard block when view.Timestamp != request.Timestamp && request.BlockNum > view.BlockNum
	if cbft.viewChange != nil && len(request.ViewChangeVotes) < cbft.getThreshold() && request.Timestamp > cbft.viewChange.Timestamp {
		log.Debug("Invalid prepare block", "hash", request.Block.Hash(), "number", request.Block.NumberU64(), "timestamp", request.Timestamp, "view", cbft.viewChange)
		return errFutileBlock
	}

	if err := cbft.VerifyHeader(cbft.blockChain, request.Block.Header(), false); err != nil {
		cbft.bp.PrepareBP().InvalidBlock(bpCtx, request, err, cbft)
		log.Error("Failed to verify header in PrepareBlockMsg, discard this msg", "peer", nodeId, "err", err)
		return err
	}

	ext := cbft.blockExtMap.findBlock(request.Block.Hash(), request.Block.NumberU64())
	if (ext != nil && ext.block != nil) || cbft.blockChain.HasBlock(request.Block.Hash(), request.Block.NumberU64()) {
		log.Warn("Block already in blockchain, discard this msg", "prepare block", request.String())
		return nil
	}

	err := cbft.verifyValidatorSign(request.Block.NumberU64(), request.ProposalIndex, request.ProposalAddr, request, request.Signature[:])
	if err != nil {
		cbft.bp.PrepareBP().InvalidBlock(bpCtx, request, err, cbft)
		cbft.log.Error("Verify prepareBlock signature fail", "number", request.Block.NumberU64(), "hash", request.Block.Hash())
		return err
	}

	if !cbft.IsConsensusNode() && !cbft.agency.IsCandidateNode(cbft.config.NodeID) {
		log.Warn("Local node is not consensus node,discard this msg")
		return errInvalidatorCandidateAddress
	} else if !cbft.CheckConsensusNode(request.ProposalAddr) {
		cbft.bp.PrepareBP().InvalidBlock(bpCtx, request,
			fmt.Errorf("remote node is not consensus node address:%s", request.ProposalAddr.String()), cbft)
		log.Warn("Remote node is not consensus node,discard this msg", "address", request.ProposalAddr)
		return errInvalidatorCandidateAddress
	}

	cbft.log.Debug("Receive prepare block", "number", request.Block.NumberU64(), "hash", request.Block.Hash(), "peer", nodeId, "ViewChangeVotes", len(request.ViewChangeVotes))
	ext = NewBlockExtByPrepareBlock(request, cbft.nodeLength())

	if len(request.ViewChangeVotes) != 0 && request.View != nil {
		if len(request.ViewChangeVotes) < cbft.getThreshold() {
			cbft.bp.PrepareBP().InvalidBlock(bpCtx, request, errTwoThirdPrepareVotes, cbft)
			cbft.log.Error(fmt.Sprintf("Receive not enough prepareVotes %d threshold %d", len(request.ViewChangeVotes), cbft.getThreshold()))
			return errTwoThirdViewchangeVotes
		}

		if cbft.getHighestLogical().number < request.View.BaseBlockNum ||
			cbft.blockExtMap.findBlock(request.View.BaseBlockHash, request.View.BaseBlockNum) == nil {
			cbft.bp.PrepareBP().InvalidBlock(bpCtx, request, errNotFoundViewBlock, cbft)
			cbft.handler.Send(nodeId, &getHighestPrepareBlock{Lowest: cbft.getRootIrreversible().number + 1})
			cbft.log.Error(fmt.Sprintf("View Block is not found, hash:%s, number:%d, logical:%d", request.View.BaseBlockHash.TerminalString(), request.View.BaseBlockNum, cbft.getHighestLogical().number))
			return errNotFoundViewBlock
		}

		oldViewChange := cbft.viewChange
		viewChange := request.View
		if cbft.viewChange == nil || cbft.viewChange.Timestamp <= viewChange.Timestamp {
			cbft.log.Debug("New PrepareBlock is not match current view, need change")
			cbft.viewChange = viewChange
		}
		if err := cbft.checkViewChangeVotes(request.ViewChangeVotes); err != nil {
			cbft.bp.PrepareBP().InvalidViewChangeVote(bpCtx, request, err, cbft)
			cbft.viewChange = oldViewChange
			cbft.log.Error("Receive prepare invalid block", "err", err)
			viewChangeVoteVerifyFailMeter.Mark(1)
			return err
		}

		cbft.log.Debug("Receive prepare block, check view change prepareVotes success, view change prepareVotes need added")

		for _, v := range request.ViewChangeVotes {
			cbft.viewChangeVotes[v.ValidatorAddr] = v
		}
		//todo check fork, clear all block larger than the request block
		//change producer
		cbft.producerBlocks = NewProducerBlocks(nodeId, request.Block.NumberU64())

		//receive 2f+1 view vote , clear last view state
		if cbft.agreeViewChange() {
			viewChangeConfirmedTimer.UpdateSince(time.Unix(int64(cbft.viewChange.Timestamp), 0))
			cbft.bp.ViewChangeBP().TwoThirdViewChangeVotes(bpCtx, cbft.viewChange, cbft.viewChangeVotes, cbft)
			var newHeader *types.Header
			viewBlock := cbft.blockExtMap.findBlock(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum)

			if viewBlock == nil {
				cbft.bp.ViewChangeBP().InvalidViewChangeBlock(bpCtx, cbft.viewChange, cbft)
				log.Error("ViewChange block find error", "BaseBlockHash", cbft.viewChange.BaseBlockHash,
					"BaseBlockNum", cbft.viewChange.BaseBlockNum, "blockMap", cbft.blockExtMap.BlockString())
				cbft.handler.Send(nodeId, &getPrepareBlock{Hash: cbft.viewChange.BaseBlockHash, Number: cbft.viewChange.BaseBlockNum})
				//panic("Find nil block")
			} else {
				newHeader = viewBlock.block.Header()
				injectBlock := cbft.blockExtMap.findBlockByNumber(cbft.viewChange.BaseBlockNum+1, cbft.getHighestLogical().number)
				start := time.Now()
				cbft.txPool.ForkedReset(newHeader, injectBlock)
				cbft.bp.InternalBP().ForkedResetTxPool(bpCtx, newHeader, injectBlock, time.Now().Sub(start), cbft)

			}

			cbft.clearPending()
			cbft.ClearChildren(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum, cbft.viewChange.Timestamp)
		}
		ext.view = cbft.viewChange
		ext.viewChangeVotes = request.ViewChangeVotes
	}

	switch cbft.AcceptPrepareBlock(request) {
	case Accept:
		cbft.bp.PrepareBP().AcceptBlock(bpCtx, request, cbft)
		if cbft.producerBlocks != nil {
			if cbft.producerBlocks.Limited(cbft) {
				cbft.log.Error("The producer produce block has over limit", "hash", ext.block.Hash(), "number", ext.block.Number(), "proposalIndex", request.ProposalIndex, "proposalAddr", request.ProposalAddr)
				return errors.New("over limit")
			}

			cbft.producerBlocks.AddBlock(ext.block)
			cbft.log.Debug("Add producer block", "hash", ext.block.Hash(), "number", ext.block.Number(), "producer", cbft.producerBlocks.String())
		}

		// if accept the block then forward the message
		if propagation && cbft.needBroadcast(nodeId, request) {
			go cbft.handler.SendBroadcast(&prepareBlockHash{Hash: request.Block.Hash(), Number: request.Block.NumberU64()})
		}
		consensusJoinCounter.Inc(1)

		return cbft.OnNewBlock(ext)
	case Cache:
		cbft.bp.PrepareBP().CacheBlock(bpCtx, request, cbft)
		cbft.log.Info("Cache block", "hash", ext.block.Hash(), "number", ext.block.NumberU64())
		// if cache the block then forward the message
		if propagation && cbft.needBroadcast(nodeId, request) {
			go cbft.handler.SendBroadcast(&prepareBlockHash{Hash: request.Block.Hash(), Number: request.Block.NumberU64()})
		}
	case Discard:
		cbft.bp.PrepareBP().DiscardBlock(bpCtx, request, cbft)
		//todo changing view discard block
		cbft.log.Info("Discard block", "hash", ext.block.Hash(), "number", ext.block.NumberU64())
	}
	return nil
}

// OnNewBlock is called by protocol handler when it received a new block by P2P.
func (cbft *Cbft) OnNewBlock(ext *BlockExt) error {
	rcvBlock := ext.block
	cbft.log.Debug("Receive new block", "hash", rcvBlock.Hash(), "number",
		rcvBlock.NumberU64(), "ParentHash", rcvBlock.ParentHash())

	cbft.blockReceiver(ext)
	return nil
}

//blockReceiver handles the new block
func (cbft *Cbft) blockReceiver(ext *BlockExt) {

	cbft.blockExtMap.Add(ext.block.Hash(), ext.block.NumberU64(), ext)
	blocks := cbft.blockExtMap.GetSubChainUnExecuted()
	log.Debug("Receive block", "unexecuted", len(blocks), "block map", cbft.blockExtMap.Len())
	for _, ext := range blocks {
		ext.executing = true
	}
	cbft.innerUnExecutedBlockCh <- blocks
}

func (cbft *Cbft) executeBlockLoop() {
	for {
		select {
		case blocks := <-cbft.innerUnExecutedBlockCh:

			//execute block from small to large
			cbft.executeBlock(blocks)
		}
	}
}

// signReceiver handles the received block signature
func (cbft *Cbft) prepareVoteReceiver(peerID discover.NodeID, vote *prepareVote) {
	cbft.log.Debug("Receive new vote",
		"vote", vote.String(),
		"state", cbft.blockState())
	ext := cbft.blockExtMap.findBlock(vote.Hash, vote.Number)
	if ext == nil {
		cbft.handler.Send(peerID, &getPrepareBlock{Hash: vote.Hash, Number: vote.Number})
		cbft.log.Warn("Have not received the corresponding block", "hash", vote.Hash, "number", vote.Number)
		//the block is nil
		ext = NewBlockExtByPeer(nil, vote.Number, cbft.nodeLength())

		ext.timestamp = vote.Timestamp
	}

	hadSend := (ext.inTree && ext.isExecuted && ext.isConfirmed)
	ext.prepareVotes.Add(vote)

	cbft.log.Info("Add prepare vote success", "number", ext.number, "hash", vote.Hash, "votes", ext.prepareVotes.Len(), "voteBits", ext.prepareVotes.voteBits.String())

	cbft.saveBlockExt(vote.Hash, ext)

	//receive enough signature broadcast
	if ext.inTree && ext.isExecuted && ext.isConfirmed {
		cbft.bp.PrepareBP().TwoThirdVotes(context.TODO(), vote, cbft)
		if h := cbft.blockExtMap.FindHighestConfirmedWithHeader(); h != nil {
			cbft.highestConfirmed.Store(h)
			blockConfirmedMeter.Mark(1)
			blockConfirmedTimer.UpdateSince(time.Unix(int64(ext.timestamp), 0))
			cbft.flushReadyBlock()
			cbft.updateValidator()
		}
		if !hadSend {
			cbft.log.Debug("Send Confirmed Block", "hash", ext.block.Hash(), "number", ext.block.NumberU64())
			cbft.bp.InternalBP().NewHighestConfirmedBlock(context.TODO(), ext, cbft)
			cbft.handler.SendAllConsensusPeer(&confirmedPrepareBlock{Hash: ext.block.Hash(), Number: ext.block.NumberU64(), VoteBits: ext.prepareVotes.voteBits})
		}
	}

}

//Receive executed block status, remove block if status is error
//If status is nil, broadcast itself PrepareVote about this block
//Reset highest logical (because receive block is order)
//Reset highest confirmed if this block had 2f+1 prepareVotes
//Flush Block (2f+1 ViewChangeVote and 2f+1 BaseBlock's PrepareVote
func (cbft *Cbft) OnExecutedBlock(bs *ExecuteBlockStatus) {
	if bs.err != nil {
		bs.block.inTree = false
		cbft.blockExtMap.RemoveBlock(bs.block)
		cbft.log.Error("Execute block failed", "err", bs.err, "block", bs.block.String())
	} else {
		bs.block.inTree = true
		bs.block.isExecuted = true
		//If blockExtMap is removed when viewchange reseted, so stop reset txpool and send PrepareVote
		if cbft.blockExtMap.findBlock(bs.block.block.Hash(), bs.block.number) != nil {
			start := time.Now()
			cbft.reset(bs.block.block)
			cbft.bp.InternalBP().ResetTxPool(context.TODO(), bs.block, time.Now().Sub(start), cbft)
			cbft.highestLogical.Store(bs.block)
			cbft.bp.InternalBP().NewHighestLogicalBlock(context.TODO(), bs.block, cbft)
			cbft.sendPrepareVote(bs.block)
			//cbft.bp.PrepareBP().SendPrepareVote(context.TODO(), bs.block, cbft)

			highest := cbft.blockExtMap.FindHighestConfirmed(cbft.getHighestConfirmed().block.Hash(), cbft.getHighestConfirmed().block.NumberU64())
			if bs.block.isConfirmed {
				if highest != nil && highest.number > cbft.getHighestConfirmed().number {
					cbft.highestConfirmed.Store(highest)
					cbft.bp.InternalBP().NewHighestConfirmedBlock(context.TODO(), highest, cbft)
				}
				cbft.log.Debug("Send Confirmed Block", "hash", bs.block.block.Hash(), "number", bs.block.block.NumberU64())
				cbft.handler.SendAllConsensusPeer(&confirmedPrepareBlock{Hash: bs.block.block.Hash(), Number: bs.block.block.NumberU64(), VoteBits: bs.block.prepareVotes.voteBits})
				blockConfirmedMeter.Mark(1)
			}

			if cbft.viewChange != nil && len(cbft.viewChangeVotes) >= cbft.getThreshold() && cbft.blockExtMap.head.number != cbft.viewChange.BaseBlockNum {
				cbft.flushReadyBlock()
			}

			if bs.block.isConfirmed {
				cbft.updateValidator()
			}
			cbft.log.Debug("Execute block success", "block", bs.block.String())
		}
	}
}

//Send PrepareVote if execute block success
func (cbft *Cbft) sendPrepareVote(ext *BlockExt) {
	cbft.log.Debug("Need send prepare vote", "hash", ext.block.Hash(), "number", ext.block.NumberU64())

	validator, err := cbft.getValidators().NodeIndexAddress(cbft.config.NodeID)
	if ext.number <= cbft.localHighestPrepareVoteNum {
		cbft.log.Warn("May happen double prepare vote")
		return
	}
	if err == nil {
		pv := &prepareVote{
			Timestamp:      ext.view.Timestamp,
			Hash:           ext.block.Hash(),
			Number:         ext.block.NumberU64(),
			ValidatorIndex: uint32(validator.Index),
			ValidatorAddr:  validator.Address,
		}

		sign, err := cbft.signMsg(pv)
		if err == nil {
			pv.Signature.SetBytes(sign)
			if cbft.viewChange != nil && !cbft.agreeViewChange() && cbft.viewChange.BaseBlockNum < ext.block.NumberU64() {
				cbft.log.Debug("Cache prepareVote, view is changing", "prepareVote", pv.String(), "view", cbft.viewChange.String(), "len", len(cbft.viewChangeVotes))
				cbft.pendingVotes.Add(pv.Hash, pv)
			} else {
				ext.prepareVotes.Add(pv)
				cbft.blockExtMap.Add(pv.Hash, pv.Number, ext)
				cbft.log.Debug("Broadcast prepare vote", "vote", pv.String())
				cbft.handler.SendAllConsensusPeer(pv)
				cbft.bp.PrepareBP().SendPrepareVote(context.TODO(), pv, cbft)
				cbft.SetLocalHighestPrepareNum(pv.Number)
			}
		} else {
			log.Error("Signature failed", "hash", ext.block.Hash(), "number", ext.block.NumberU64(), "err", err)
		}
	} else {
		log.Error("Local node is not consensus node", "err", err)
	}
}

// executeBlockAndDescendant executes the block's transactions and its descendant
func (cbft *Cbft) executeBlock(blocks []*BlockExt) {
	for _, ext := range blocks {
		//Execute blocks is async, clear all children block when new view change was confirmed.
		if ext == nil || ext.parent == nil {
			cbft.log.Warn("Block was cleared, block is invalid block. stop execute all children block")
			return
		}

		start := time.Now()
		err := cbft.execute(ext, ext.parent)
		if err != nil {
			cbft.bp.InternalBP().InvalidBlock(context.TODO(), ext.block.Hash(), ext.timestamp, ext.block.NumberU64(), err)
		}
		cbft.bp.InternalBP().ExecuteBlock(context.TODO(), ext.block.Hash(), ext.block.NumberU64(), ext.timestamp, time.Now().Sub(start))
		blockExecuteTimer.UpdateSince(start)
		//send syncState after execute block
		ext.SetSyncState(err)

		cbft.executeBlockCh <- &ExecuteBlockStatus{
			block: ext,
			err:   err,
		}
	}
}

// saveBlockExt saves block in memory
func (cbft *Cbft) saveBlockExt(hash common.Hash, ext *BlockExt) {
	cbft.blockExtMap.Add(hash, ext.number, ext)
	cbft.log.Debug("Save block in memory", "hash", hash, "number", ext.number, "had block", ext.block != nil, "total", cbft.blockExtMap.Total())
}

// CheckConsensusNode check if the nodeID is a consensus node.
func (cbft *Cbft) CheckConsensusNode(address common.Address) bool {
	_, err := cbft.getValidators().AddressIndex(address)
	return err == nil
}

// IsConsensusNode check if local is a consensus node.
func (cbft *Cbft) IsConsensusNode() bool {
	_, err := cbft.getValidators().NodeIndex(cbft.config.NodeID)
	return err == nil
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	if header.Number == nil {
		cbft.log.Warn(fmt.Sprintf("Verify header failed, unknow block"))
		return errUnknownBlock
	}

	cbft.log.Trace("Verify header", "hash", header.Hash(), "number", header.Number.Uint64(), "seal", seal)
	if len(header.Extra) < extraSeal {
		cbft.log.Warn(fmt.Sprintf("Verify header failed, miss sign  number:%d", header.Number.Uint64()))
		return errMissingSignature
	}
	return nil
}

// execute executes the block's transactions based on its parent
// if success then save the receipts and state to consensusCache
func (cbft *Cbft) execute(ext *BlockExt, parent *BlockExt) error {
	cbft.log.Debug("execute block based on parent", "block", ext.String(), "parent", parent.String())
	state, err := cbft.blockChainCache.MakeStateDB(parent.block)
	if err != nil {
		cbft.log.Error("execute block error, cannot make state based on parent", "err", err, "block", ext.String(), "parent", parent.String(), "err", err)
		blockVerifyFailMeter.Mark(1)
		return errors.New("execute block error")
	}

	//to execute
	receipts, err := cbft.blockChain.ProcessDirectly(ext.block, state, parent.block)
	if err == nil {
		//save the receipts and state to consensusCache
		stateIsNil := state == nil
		cbft.log.Debug("execute block success", "block", ext.String(), "parent", parent.String(), "lenReceipts", len(receipts), "stateIsNil", stateIsNil, "root", ext.block.Root())
		sealHash := ext.block.Header().SealHash()
		cbft.blockChainCache.WriteReceipts(sealHash, receipts, ext.block.NumberU64())
		cbft.blockChainCache.WriteStateDB(sealHash, state, ext.block.NumberU64())

	} else {
		cbft.log.Error("execute block error", "err", err, "block", ext.String(), "parent", parent.String())
		blockVerifyFailMeter.Mark(1)
		return fmt.Errorf("execute block error, err:%s", err.Error())
	}
	return nil
}

func (cbft *Cbft) blockState() string {
	return fmt.Sprintf(`hightestLogical hash:%s number:%d,highestConfirmed hash:%s number:%d,root hash:%s, number:%d`,
		cbft.getHighestLogical().block.Hash().TerminalString(),
		cbft.getHighestLogical().number,
		cbft.getHighestConfirmed().block.Hash().TerminalString(),
		cbft.getHighestConfirmed().number,
		cbft.getRootIrreversible().block.Hash().TerminalString(),
		cbft.getRootIrreversible().number,
	)
}

//func extraBlocks(exts []*BlockExt) []*types.Block {
//	blocks := make([]*types.Block, len(exts))
//	for idx, ext := range exts {
//		blocks[idx] = ext.block
//	}
//	return blocks
//}

// checkFork checks if the logical path is changed cause the newConfirmed, if changed, this is a new fork.
func (cbft *Cbft) checkFork(newConfirmed *BlockExt) {
	//newHighestConfirmed := cbft.findLastClosestConfirmedIncludingSelf(newConfirmed)
	//if newHighestConfirmed != nil && newHighestConfirmed.block.Hash() != cbft.getHighestConfirmed().block.Hash() {
	//	//fork
	//	newHighestLogical := cbft.findHighestLogical(newHighestConfirmed)
	//
	//	//newPath := cbft.backTrackBlocks(newHighestLogical, cbft.getRootIrreversible(), true)
	//	//origPath := cbft.backTrackBlocks(cbft.getHighestLogical(), cbft.getRootIrreversible(), true)
	//
	//	//cbft.log.Debug("the block chain in memory forked",
	//	//	"newHighestConfirmedHash", newHighestConfirmed.block.Hash().TerminalString(),
	//	//	"newHighestConfirmedNumber", newHighestConfirmed.number,
	//	//	"newHighestLogicalHash", newHighestLogical.block.Hash().TerminalString(),
	//	//	"newHighestLogicalNumber", newHighestLogical.number,
	//	//	"len(newPath)", len(newPath), "len(origPath)", len(origPath),
	//	//)
	//	//oldTress, newTress := cbft.forked(origPath, newPath)
	//
	//	//cbft.txPool.ForkedReset(extraBlocks(oldTress), extraBlocks(newTress))
	//
	//	//fork from to lower block
	//	cbft.highestConfirmed.Store(newHighestConfirmed)
	//	cbft.highestLogical.Store(newHighestLogical)
	//	cbft.log.Warn("chain is forked")
	//}
}

func (cbft *Cbft) HasTwoThirdsMajorityViewChangeVotes() bool {
	//cbft.mux.Lock()
	//defer cbft.mux.Unlock()
	cbft.log.Debug(fmt.Sprintf("receive prepareVotes:%d threshold:%d", len(cbft.viewChangeVotes), cbft.getThreshold()))
	return cbft.agreeViewChange()
}

func (cbft *Cbft) CalcBlockDeadline(timePoint int64) (time.Time, error) {
	node, err := cbft.getValidators().NodeIndex(cbft.config.NodeID)
	if err != nil {
		return time.Time{}, err
	}
	nodeIdx := node.Index
	startEpoch := cbft.startTimeOfEpoch * 1000

	tm := common.MillisToTime(timePoint)

	if nodeIdx >= 0 {
		if cbft.getValidators().Len() == 1 {
			return tm.Add(time.Duration(cbft.config.Period) * time.Second), err
		}
		durationPerNode := cbft.config.Duration * 1000
		durationPerTurn := durationPerNode * int64(cbft.getValidators().Len())

		min := int64(nodeIdx) * (durationPerNode)
		value := (timePoint - startEpoch) % durationPerTurn
		max := int64(nodeIdx+1) * durationPerNode

		cnt := int64(cbft.config.Duration) / int64(cbft.config.Period)
		slots := make([]int64, cnt)
		var i int64
		for i = 0; i < cnt; i++ {
			slots[i] = min + (i*1000)*int64(cbft.config.Period)
		}

		curIdx := (value % durationPerNode) / (1000 * int64(cbft.config.Period))
		lastBlock := int(curIdx+1) == len(slots)
		nextSlotValue := max
		if !lastBlock {
			nextSlotValue = slots[curIdx+1]
		}

		remaining := time.Duration(nextSlotValue-value) * time.Millisecond
		interval := time.Duration(cbft.config.BlockInterval) * time.Millisecond
		cbft.log.Trace("Calc block deadline", "remaining", remaining, "interval", interval, "curIdx", curIdx)
		if remaining > interval {
			remaining = remaining - interval
		} else {
			remaining = 50 * time.Millisecond // 50ms
		}
		return tm.Add(remaining), err
	}
	return tm.Add(50 * time.Millisecond), err // 50ms

}

func (cbft *Cbft) CalcNextBlockTime(timePoint int64) (time.Time, error) {
	vn, err := cbft.getValidators().NodeIndex(cbft.config.NodeID)
	if err != nil {
		return time.Time{}, err
	}
	nodeIdx := vn.Index
	startEpoch := cbft.startTimeOfEpoch * 1000
	tm := common.MillisToTime(timePoint)

	if nodeIdx >= 0 {
		if cbft.getValidators().Len() == 1 {
			return tm.Add(time.Duration(cbft.config.Period) * time.Second), nil
		}
		durationPerNode := cbft.config.Duration * 1000
		durationPerTurn := durationPerNode * int64(cbft.getValidators().Len())

		min := int64(nodeIdx) * (durationPerNode)
		value := (timePoint - startEpoch) % durationPerTurn
		max := int64(nodeIdx+1) * durationPerNode

		log.Trace("Calc next block time", "min", min, "value", value, "max", max)
		var offset int64
		if value >= min && value < max {
			cnt := int64(cbft.config.Duration) / int64(cbft.config.Period)
			slots := make([]int64, cnt)
			var i int64
			for i = 0; i < cnt; i++ {
				slots[i] = min + (i*1000)*int64(cbft.config.Period)
			}
			curIdx := (value % durationPerNode) / (1000 * int64(cbft.config.Period))
			cbft.log.Trace("Calc next block time", "min", min, "value", value, "max", max, "curIdx", curIdx, "slots", len(slots))
			lastBlock := int(curIdx+1) == len(slots)
			nextSlotValue := max
			if !lastBlock {
				nextSlotValue = slots[curIdx+1]
			}
			remaining := nextSlotValue - value
			offset = remaining + durationPerTurn - durationPerNode
			if !lastBlock {
				offset = remaining
			}
		} else if value < min {
			offset = min - value
		} else {
			// value > max
			last := int64(cbft.getValidators().Len()) * durationPerNode
			offset = last - value + min
		}
		return tm.Add(time.Duration(offset) * time.Millisecond), nil
	}
	return tm.Add(time.Duration(cbft.config.Period) * time.Second), nil // 1s

}

// ConsensusNodes returns all consensus nodes.
func (cbft *Cbft) ConsensusNodes() ([]discover.NodeID, error) {
	cbft.log.Trace(fmt.Sprintf("dposNodeCount:%d", cbft.getValidators().Len()))
	return cbft.getValidators().NodeList(), nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (cbft *Cbft) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	cbft.log.Trace("verify headers", "total", len(headers))

	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for _, header := range headers {
			err := cbft.VerifyHeader(chain, header, false)

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (cbft *Cbft) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	cbft.log.Trace("verify seal", "hash", header.Hash(), "number", header.Number.String())

	return cbft.verifySeal(chain, header, nil)
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (cbft *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	cbft.log.Debug("prepare", "hash", header.Hash(), "number", header.Number.Uint64())

	if cbft.getHighestLogical().block == nil {
		cbft.log.Error("highest logical block is empty")
		return errors.New("highest logical block is empty")
	}

	//header.Extra[0:31] to store block's version info etc. and right pad with 0x00;
	//header.Extra[32:] to store block's sign of producer, the length of sign is 65.
	if len(header.Extra) < 32 {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, 32-len(header.Extra))...)
	}
	header.Extra = header.Extra[:32]

	//init header.Extra[32: 32+65]
	header.Extra = append(header.Extra, make([]byte, consensus.ExtraSeal)...)
	return nil
}

// Finalize implements consensus.Engine, no block
// rewards given, and returns the final block.
func (cbft *Cbft) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	cbft.log.Debug("finalize block", "hash", header.Hash(), "number", header.Number.Uint64(), "txs", len(txs), "receipts", len(receipts))
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	return types.NewBlock(header, txs, receipts), nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (cbft *Cbft) SealHash(header *types.Header) common.Hash {
	cbft.log.Debug("seal", "hash", header.Hash(), "number", header.Number.Uint64())
	return header.SealHash()
}

// Close implements consensus.Engine. It's a noop for cbft as there is are no background threads.
func (cbft *Cbft) Close() error {
	cbft.log.Info("close cbft consensus")
	cbft.closeOnce.Do(func() {
		// Short circuit if the exit channel is not allocated.
		if cbft.exitCh == nil {
			return
		}
		cbft.exitCh <- struct{}{}
		close(cbft.exitCh)
	})
	if cbft.wal != nil {
		cbft.wal.Close()
	}
	cbft.bp.Close()
	return nil
}

//receive eth downloader event
func (cbft *Cbft) update() {
	events := cbft.eventMux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.FailedEvent{})
	defer events.Unsubscribe()

	for {
		select {
		case ev := <-events.Chan():
			if ev == nil {
				return
			}
			switch ev.Data.(type) {
			case downloader.StartEvent:
				atomic.StoreInt32(&cbft.running, 0)
				cbft.log.Debug("Start sync block from download, start cbft consensus",
					"logical", cbft.getHighestLogical().number,
					"confirm", cbft.getHighestConfirmed().number,
					"root", cbft.getRootIrreversible().number)
			case downloader.DoneEvent, downloader.FailedEvent:
				cbft.log.Debug("Sync block from download is finish, start cbft consensus",
					"logical", cbft.getHighestLogical().number,
					"confirm", cbft.getHighestConfirmed().number,
					"root", cbft.getRootIrreversible().number)
				cbft.handler.SendAllConsensusPeer(&getHighestPrepareBlock{Lowest: cbft.getHighestConfirmed().number})
				atomic.StoreInt32(&cbft.running, 1)
				// stop immediately and ignore all further pending events
				return
			}
		case <-cbft.exitCh:
			return
		}
	}
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (cbft *Cbft) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{}
}

func (cbft *Cbft) Protocols() []p2p.Protocol {
	return cbft.handler.Protocols()
}

// OnBlockSignature is called by by protocol handler when it received a new block signature by P2P.
func (cbft *Cbft) OnPrepareVote(peerID discover.NodeID, vote *prepareVote, propagation bool) error {
	cbft.log.Debug("Receive prepare vote", "peer", peerID, "vote", vote.String())
	bpCtx := context.WithValue(context.Background(), "peer", peerID)
	cbft.bp.PrepareBP().ReceiveVote(bpCtx, vote, cbft)
	err := cbft.verifyValidatorSign(vote.Number, vote.ValidatorIndex, vote.ValidatorAddr, vote, vote.Signature[:])
	if err != nil {
		cbft.bp.PrepareBP().InvalidVote(bpCtx, vote, err, cbft)
		cbft.log.Error("Verify vote error", "err", err)
		return err
	}

	switch cbft.AcceptPrepareVote(vote) {
	case Accept:
		cbft.log.Debug("Accept block vote", "vote", vote.String())
		cbft.bp.PrepareBP().AcceptVote(bpCtx, vote, cbft)
		if err := cbft.evPool.AddPrepareVote(vote); err != nil {
			if _, ok := err.(*DuplicatePrepareVoteEvidence); ok {
				cbft.log.Warn("Receive DuplicatePrepareVoteEvidence msg", "err", err.Error())
				return err
			}
		}

		cbft.prepareVoteReceiver(peerID, vote)
	case Cache:
		cbft.log.Debug("View changing, add vote into process queue", "vote", vote.String())
		cbft.bp.PrepareBP().CacheVote(bpCtx, vote, cbft)
		//changing view
		cbft.AddProcessingVote(peerID, vote)
	case Discard:
		cbft.bp.PrepareBP().DiscardVote(bpCtx, vote, cbft)
		cbft.log.Debug("Discard PrepareVote", "vote", vote.String())
	}
	cbft.log.Trace("Processing vote end", "hash", vote.Hash, "number", vote.Number)

	// rule:
	if propagation && cbft.needBroadcast(peerID, vote) {
		cbft.log.Debug("Broadcast the message of prepareVote", "FromPeerId", peerID.String())
		go cbft.handler.SendBroadcast(vote)
	}

	return nil
}

// OnPong is called by protocol handler when it received a new Pong message by P2P.
func (cbft *Cbft) OnPong(nodeID discover.NodeID, netLatency int64) error {
	cbft.log.Trace("pong", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]), "netLatency", netLatency)

	cbft.netLatencyLock.Lock()
	defer cbft.netLatencyLock.Unlock()

	if netLatency >= maxPingLatency {
		return nil
	}

	latencyList, exist := cbft.netLatencyMap[nodeID]
	if !exist {
		cbft.netLatencyMap[nodeID] = list.New()
		cbft.netLatencyMap[nodeID].PushBack(netLatency)
	} else {
		if latencyList.Len() > 5 {
			e := latencyList.Front()
			cbft.netLatencyMap[nodeID].Remove(e)
		}
		cbft.netLatencyMap[nodeID].PushBack(netLatency)
	}
	return nil
}

// avgLatency statistics the net latency between local and other peers.
func (cbft *Cbft) avgLatency(nodeID discover.NodeID) int64 {
	if latencyList, exist := cbft.netLatencyMap[nodeID]; exist {
		sum := int64(0)
		counts := int64(0)
		for e := latencyList.Front(); e != nil; e = e.Next() {
			if latency, ok := e.Value.(int64); ok {
				counts++
				sum += latency
			}
		}
		if counts > 0 {
			return sum / counts
		}
	}
	return cbft.config.MaxLatency
}

// HighestLogicalBlock returns the cbft.highestLogical.block.
func (cbft *Cbft) HighestLogicalBlock() *types.Block {
	if v := cbft.highestLogical.Load(); v != nil {
		return v.(*BlockExt).block
	}
	return nil
}

// HighestConfirmedBlock returns the cbft.highestConfirmed.block.
func (cbft *Cbft) HighestConfirmedBlock() *types.Block {
	if cbft.getHighestConfirmed() == nil {
		return nil
	} else {
		return cbft.getHighestConfirmed().block
	}
}

// storeBlocks sends the blocks to cbft.cbftResultOutCh, the receiver will write them into chain
func (cbft *Cbft) storeBlocks(blocksToStore []*BlockExt) {
	for _, ext := range blocksToStore {
		extra, err := cbft.encodeExtra(ext.BlockExtra())
		if err != nil {
			cbft.log.Error("Encode ExtraData failed", "err", err)
			continue
		}
		cbftResult := cbfttypes.CbftResult{
			Block:     ext.block,
			ExtraData: extra,
			SyncState: ext.syncState,
		}
		cbft.log.Debug("Send consensus result to worker", "block", ext.String())
		cbft.bp.InternalBP().StoreBlock(context.TODO(), ext, cbft)
		cbft.eventMux.Post(cbftResult)
	}
}

func (cbft *Cbft) encodeExtra(bx *BlockExtra) ([]byte, error) {
	extra := []byte{cbftVersion}
	bxBytes, err := rlp.EncodeToBytes(bx)
	if err != nil {
		return nil, err
	}
	extra = append(extra, bxBytes...)
	return extra, nil
}

func (cbft *Cbft) decodeExtra(extra []byte) (byte, *BlockExtra, error) {
	if len(extra) == 0 {
		return 0, nil, fmt.Errorf("empty extra")
	}
	version := extra[0]
	var bx BlockExtra
	err := rlp.DecodeBytes(extra[1:], &bx)
	if err != nil {
		return 0, nil, err
	}
	return version, &bx, nil
}

// inTurn return if it is local's turn to package new block.
func (cbft *Cbft) inTurn(curTime int64) bool {
	//curTime := toMilliseconds(time.Now())
	/*
		inturn := cbft.calTurn(curTime-25, cbft.config.NodeID)
		if inturn {
			inturn = cbft.calTurn(curTime+300, cbft.config.NodeID)
		}
	*/
	return cbft.calTurn(curTime, cbft.config.NodeID)

}

// inTurnVerify verifies the time is in the time-window of the nodeID to package new block.
func (cbft *Cbft) inTurnVerify(rcvTime int64, nodeID discover.NodeID) bool {
	latency := cbft.avgLatency(nodeID)
	if latency >= maxAvgLatency {
		cbft.log.Warn("check if peer's turn to commit block", "result", false, "peerID", nodeID, "high latency ", latency)
		return false
	}
	inTurnVerify := cbft.calTurn(rcvTime-latency, nodeID)
	cbft.log.Debug("check if peer's turn to commit block", "result", inTurnVerify, "peerID", nodeID, "latency", latency)
	return inTurnVerify
}

//isLegal verifies the time is legal to package new block for the nodeID.
func (cbft *Cbft) isLegal(rcvTime int64, addr common.Address) bool {
	nodeIdx, err := cbft.getValidators().AddressIndex(addr)
	if err != nil {
		cbft.log.Error("Get address index failed", "err", err)
		return false
	}
	return cbft.calTurnIndex(rcvTime, nodeIdx.Index)
}

func (cbft *Cbft) calTurn(timePoint int64, nodeID discover.NodeID) bool {
	vn, err := cbft.getValidators().NodeIndex(nodeID)
	if err != nil {
		return false
	}
	return cbft.calTurnIndex(timePoint, vn.Index)
}

func (cbft *Cbft) calTurnIndex(timePoint int64, nodeIdx int) bool {

	startEpoch := cbft.startTimeOfEpoch * 1000

	if nodeIdx >= 0 {
		if cbft.getValidators().Len() == 1 {
			return true
		}
		durationPerNode := cbft.config.Duration * 1000
		durationPerTurn := durationPerNode * int64(cbft.getValidators().Len())

		min := int64(nodeIdx) * (durationPerNode)

		value := (timePoint - startEpoch) % durationPerTurn

		max := int64(nodeIdx+1) * durationPerNode

		if value >= min && value < max {
			//cbft.log.Debug("calTurn return true", "idx", nodeIdx, "min", min, "value", value, "max", max, "timePoint", common.MillisToString(timePoint), "startEpoch", common.MillisToString(startEpoch))
			return true
		} else {
			//cbft.log.Debug("calTurn return false", "idx", nodeIdx, "min", min, "value", value, "max", max, "timePoint", common.MillisToString(timePoint), "startEpoch", common.MillisToString(startEpoch))
		}
	}
	return false
}

// producer's signature = header.Extra[32:]
// public key can be recovered from signature, the length of public key is 65,
// the length of NodeID is 64, nodeID = publicKey[1:]
//func ecrecover(header *types.Header) (discover.NodeID, []byte, error) {
//	var nodeID discover.NodeID
//	if len(header.Extra) < extraSeal {
//		return nodeID, []byte{}, errMissingSignature
//	}
//	signature := header.Extra[len(header.Extra)-extraSeal:]
//	sealHash := header.SealHash()
//
//	pubkey, err := crypto.Ecrecover(sealHash.Bytes(), signature)
//	if err != nil {
//		return nodeID, []byte{}, err
//	}
//
//	nodeID, err = discover.BytesID(pubkey[1:])
//	if err != nil {
//		return nodeID, []byte{}, err
//	}
//	return nodeID, signature, nil
//}

// verify sign, check the sign is from the right node.
func verifySign(expectedNodeID discover.NodeID, sealHash common.Hash, signature []byte) error {
	pubkey, err := crypto.SigToPub(sealHash.Bytes(), signature)

	if err != nil {
		return err
	}

	nodeID := discover.PubkeyID(pubkey)
	if bytes.Equal(nodeID.Bytes(), expectedNodeID.Bytes()) {
		return nil
	}
	signatureVerifyFailMeter.Mark(1)
	return fmt.Errorf("verify sign failed")
}

func (cbft *Cbft) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		cbft.log.Warn("Verify seal unknown block")
		return errUnknownBlock
	}
	return nil
}

func (cbft *Cbft) signMsg(msg ConsensusMsg) (sign []byte, err error) {
	buf, err := msg.CannibalizeBytes()
	if err != nil {
		return nil, err
	}
	return crypto.Sign(buf, cbft.config.PrivateKey)
}

func (cbft *Cbft) signFn(headerHash []byte) (sign []byte, err error) {
	return crypto.Sign(headerHash, cbft.config.PrivateKey)
}

func (cbft *Cbft) getThreshold() int {
	trunc := cbft.getValidators().Len() * 2 / 3
	return trunc
}

func (cbft *Cbft) nodeLength() int {
	return cbft.getValidators().Len()
}

func (cbft *Cbft) reset(block *types.Block) {
	if _, ok := cbft.resetCache.Get(block.Hash()); !ok {
		cbft.log.Debug("Reset txpool", "hash", block.Hash(), "number", block.NumberU64(), "parentHash", block.ParentHash())
		cbft.resetCache.Add(block.Hash(), struct{}{})
		cbft.txPool.Reset(block)
	}
}
func (cbft *Cbft) OnGetBlock(hash common.Hash, number uint64, ch chan *types.Block) {
	if ext := cbft.blockExtMap.findBlock(hash, number); ext != nil {
		ch <- ext.block
	} else {
		ch <- nil
	}
}

func (cbft *Cbft) GetBlockWithoutLock(hash common.Hash, number uint64) *types.Block {
	ext := cbft.blockExtMap.findBlock(hash, number)
	if ext != nil {
		return ext.block
	}
	return nil
}
func (cbft *Cbft) GetBlock(hash common.Hash, number uint64) *types.Block {
	ch := make(chan *types.Block, 1)
	cbft.getBlockCh <- &GetBlock{hash: hash, number: number, ch: ch}
	return <-ch
}

func (cbft *Cbft) IsSignedBySelf(sealHash common.Hash, signature []byte) bool {
	return verifySign(cbft.config.NodeID, sealHash, signature) == nil
}
func (cbft *Cbft) OnHasBlock(block *HasBlock) {
	if cbft.getHighestConfirmed().number > block.number {
		block.hasCh <- true
	} else {
		block.hasCh <- false
	}
}
func (cbft *Cbft) HasBlock(hash common.Hash, number uint64) bool {
	if cbft.getHighestLogical().number >= number {
		return true
	}
	hasBlock := &HasBlock{
		hash:   hash,
		number: number,
		hasCh:  make(chan bool, 1),
	}
	cbft.hasBlockCh <- hasBlock
	has := <-hasBlock.hasCh
	if !has {
		cbft.log.Debug("Without block", "hash", hash, "number", number, "highestConfirm", cbft.getHighestConfirmed().number, "highestLogical", cbft.getHighestLogical().number, "root", cbft.getRootIrreversible().number)
	}
	return has
}

func (cbft *Cbft) Status() string {
	status := make(chan string, 1)
	cbft.statusCh <- status
	return <-status
}

func (cbft *Cbft) OnStatus(status chan string) {
	status <- cbft.RoundState.String()
}

func (cbft *Cbft) Evidences() string {
	evs := cbft.evPool.Evidences()
	if len(evs) == 0 {
		return "{}"
	}
	evds := ClassifyEvidence(evs)
	js, err := json.MarshalIndent(evds, "", "  ")
	if err != nil {
		return ""
	}
	return string(js)
}

func (cbft *Cbft) OnGetBlockByHash(hash common.Hash, ch chan *types.Block) {
	ch <- cbft.blockExtMap.findBlockByHash(hash)
}

func (cbft *Cbft) GetBlockByHash(hash common.Hash) *types.Block {
	ch := make(chan *types.Block)
	cbft.getBlockByHashCh <- &GetBlock{hash: hash, ch: ch}
	return <-ch
}

func (cbft *Cbft) CurrentBlock() *types.Block {
	return cbft.getHighestConfirmed().block
}

func (cbft *Cbft) FastSyncCommitHead() <-chan error {
	errCh := make(chan error, 1)
	cbft.fastSyncCommitHeadCh <- errCh
	return errCh
}

func (cbft *Cbft) OnFastSyncCommitHead(errCh chan error) {
	currentBlock := cbft.blockChain.CurrentBlock()
	cbft.log.Debug("Fast sync commit highestLogicalBlock", "hash", currentBlock.Hash(), "number", currentBlock.NumberU64())
	current := NewBlockExtBySeal(currentBlock, currentBlock.NumberU64(), cbft.getThreshold())
	current.number = currentBlock.NumberU64()

	if current.number > 0 && cbft.getValidators().Len() > 1 {
		var extra *BlockExtra
		var err error

		if _, extra, err = cbft.decodeExtra(current.block.ExtraData()); err != nil {
			errCh <- err
			return
		}
		current.view = extra.ViewChange

		for _, vote := range extra.Prepare {
			current.timestamp = vote.Timestamp
			current.prepareVotes.Add(vote)
		}
	}

	cbft.blockExtMap = NewBlockExtMap(current, cbft.getThreshold())
	cbft.saveBlockExt(currentBlock.Hash(), current)

	cbft.highestConfirmed.Store(current)
	cbft.highestLogical.Store(current)
	cbft.rootIrreversible.Store(current)

	errCh <- nil
}

func (cbft *Cbft) updateValidator() {
	hc := cbft.getHighestConfirmed()
	if hc.number != cbft.agency.GetLastNumber(hc.number) {
		return
	}

	// Check if we are a consensus node before updated.
	isValidatorBefore := cbft.IsConsensusNode()

	newVds, err := cbft.agency.GetValidator(hc.number + 1)
	if err != nil {
		cbft.log.Error("Get validators fail", "number", hc.number, "hash", hc.block.Hash())
		return
	}
	if newVds.Len() <= 0 {
		cbft.log.Error("Empty validators")
		return
	}
	oldVds := cbft.getValidators()
	cbft.validators.Store(newVds)
	cbft.log.Info("Update validators success", "highestConfirmed", hc.number, "hash", hc.block.Hash(), "validators", cbft.getValidators())

	cbft.afterUpdateValidator()

	if _, e := cbft.getValidators().NodeIndex(cbft.config.NodeID); e == nil && !newVds.Equal(oldVds) {
		cbft.eventMux.Post(cbfttypes.UpdateValidatorEvent{})
		log.Trace("Post UpdateValidatorEvent", "nodeID", cbft.config.NodeID)
	}

	// Check if we are become a consensus node after update.
	isValidatorAfter := cbft.IsConsensusNode()

	cbft.log.Trace("After update validators", "isValidatorBefore", isValidatorBefore, "isValidator", isValidatorAfter)
	if isValidatorBefore {
		// If we are still a consensus node, that adding
		// new validators as consensus peer, and removing
		// validators. Added as consensus peersis because
		// we need to keep connect with other validators
		// in the consensus stages. Also we are not needed
		// to keep connect with old validators.
		if isValidatorAfter {
			newNodeList := cbft.getValidators().NodeList()
			for _, nodeID := range newNodeList {
				if node, _ := oldVds.NodeIndex(nodeID); node == nil {
					cbft.eventMux.Post(cbfttypes.AddValidatorEvent{NodeID: nodeID})
					cbft.log.Trace("Post AddValidatorEvent", "nodeID", nodeID.String())
				}
			}

			oldNodeList := oldVds.NodeList()
			for _, nodeID := range oldNodeList {
				if node, _ := cbft.getValidators().NodeIndex(nodeID); node == nil {
					cbft.eventMux.Post(cbfttypes.RemoveValidatorEvent{NodeID: nodeID})
					cbft.log.Trace("Post RemoveValidatorEvent", "nodeID", nodeID.String())
				}
			}
		} else {
			oldNodeList := oldVds.NodeList()
			for _, nodeID := range oldNodeList {
				cbft.eventMux.Post(cbfttypes.RemoveValidatorEvent{NodeID: nodeID})
				cbft.log.Trace("Post RemoveValidatorEvent", "nodeID", nodeID.String())
			}
		}
	} else {
		// We are become a consensus node, that adding all
		// validators as consensus peer except us. Added as
		// consensus peers is because we need to keep connecting
		// with other validators in the consensus stages.
		if isValidatorAfter {
			newNodeList := cbft.getValidators().NodeList()
			for _, nodeID := range newNodeList {
				if cbft.config.NodeID == nodeID {
					// Ignore myself
					continue
				}
				cbft.eventMux.Post(cbfttypes.AddValidatorEvent{NodeID: nodeID})
				cbft.log.Trace("Post AddValidatorEvent", "nodeID", nodeID.String())
			}
		}

		// We are still not a consensus node, just update validator list.
	}
}

func (cbft *Cbft) needBroadcast(nodeId discover.NodeID, msg Message) bool {
	peers := cbft.handler.PeerSet().Peers()
	if len(peers) == 0 {
		return false
	}
	for _, peer := range peers {
		// exclude currently send peer.
		if peer.id == nodeId.TerminalString() {
			continue
		}
		if peer.knownMessageHash.Contains(msg.MsgHash()) {
			cbft.log.Debug("Needn't to broadcast", "type", reflect.TypeOf(msg), "hash", msg.MsgHash(), "BHash", msg.BHash().TerminalString())
			messageRepeatMeter.Mark(1)
			return false
		}
	}
	cbft.log.Debug("Need to broadcast", "type", reflect.TypeOf(msg), "hash", msg.MsgHash(), "BHash", msg.BHash().TerminalString())
	messageGossipMeter.Mark(1)
	return true
}

func (cbft *Cbft) AddJournal(info *MsgInfo) {
	msg, peerID := info.Msg, info.PeerID
	cbft.log.Debug("Load journal message from wal", "peer", peerID.TerminalString(), "msgType", reflect.TypeOf(msg))

	switch msg := msg.(type) {
	case *sendPrepareBlock:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareBlock", msg.PrepareBlock.String())
		blockExt := NewBlockExtByPrepareBlock(msg.PrepareBlock, cbft.nodeLength())
		blockExt.view = cbft.viewChange
		blockExt.viewChangeVotes = cbft.viewChangeVotes.Flatten()

		cbft.blockExtMap.Add(blockExt.block.Hash(), blockExt.block.NumberU64(), blockExt)
		blocks := cbft.blockExtMap.GetSubChainUnExecuted()
		for _, ext := range blocks {
			if ext == nil || ext.parent == nil {
				panic("add block to blockExtMap error when loading sendPrepareBlock message from wal")
			}
			err := cbft.execute(ext, ext.parent)
			if err != nil {
				panic("execute block error when loading sendPrepareBlock message from wal")
			}
		}
		cbft.sealBlockProcess(msg.PrepareBlock.Block)
	case *sendViewChange:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendViewChange", msg.ViewChange.String(), "master", msg.Master)
		cbft.newViewChangeProcess(msg.ViewChange)
	case *confirmedViewChange:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "confirmedViewChange", msg.ViewChange.String())
		if msg.Master {
			cbft.newViewChangeProcess(msg.ViewChange)
		} else {
			cbft.agreeViewChangeProcess(msg.ViewChange, msg.ViewChangeResp)
		}
		viewChangeVotes := make(map[common.Address]*viewChangeVote)
		for _, v := range msg.ViewChangeVotes {
			viewChangeVotes[v.ValidatorAddr] = v
		}
		cbft.confirmedViewChangeProcess(viewChangeVotes)
	default:
		cbft.ReceivePeerMsg(info)
	}
}

func (cbft *Cbft) isForwarded(nodeId discover.NodeID, msg Message) bool {
	peers := cbft.handler.PeerSet().Peers()
	// message of prepareBlock cannot be filtered
	for _, peer := range peers {
		if peer.id == fmt.Sprintf("%x", nodeId.Bytes()[:8]) {
			continue
		}
		if peer.knownMessageHash.Contains(msg.MsgHash()) {
			return true
		}
	}
	return false
}

/*func (cbft *Cbft) isRepeated(nodeId discover.NodeID, msg Message) bool {
	peers := cbft.handler.PeerSet().Peers()
	for _, peer := range peers {
		if peer.id == fmt.Sprintf("%x", nodeId.Bytes()[:8]) && peer.knownMessageHash.Contains(msg.MsgHash()) {
			return true
		}
	}
	return false
}*/

func (cbft *Cbft) CommitBlockBP(block *types.Block, txs int, gasUsed uint64, elapse time.Duration) {
	cbft.bp.PrepareBP().CommitBlock(context.TODO(), block, txs, gasUsed, elapse)
}

func (cbft *Cbft) TracingSwitch(flag int8) {
	if flag == 1 {
		cbft.tracing.On()
	} else {
		cbft.tracing.Off()
	}
}
