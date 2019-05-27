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
	errSign                = errors.New("sign error")
	errUnauthorizedSigner  = errors.New("unauthorized signer")
	errIllegalBlock        = errors.New("illegal block")
	lateBlock              = errors.New("block is late")
	errDuplicatedBlock     = errors.New("duplicated block")
	errBlockNumber         = errors.New("error block number")
	errUnknownBlock        = errors.New("unknown block")
	errFutileBlock         = errors.New("futile block")
	errGenesisBlock        = errors.New("cannot handle genesis block")
	errHighestLogicalBlock = errors.New("cannot find a logical block")
	errListConfirmedBlocks = errors.New("list confirmed blocks error")
	errMissingSignature    = errors.New("extra-data 65 byte signature suffix missing")

	errInitiateViewchange          = errors.New("not initiated viewchange")
	errTwoThirdViewchangeVotes     = errors.New("lower two third viewchange prepareVotes")
	errTwoThirdPrepareVotes        = errors.New("lower two third prepare prepareVotes")
	errNotFoundViewBlock           = errors.New("not found block")
	errInvalidViewChangeVotes      = errors.New("invalid prepare prepareVotes")
	errInvalidPrepareVotes         = errors.New("invalid prepare prepareVotes")
	errInvalidatorCandidateAddress = errors.New("invalid address")
	extraSeal                      = 65
	windowSize                     = 10

	//periodMargin is a percentum for period margin
	periodMargin = uint64(20)

	//maxPingLatency is the time in milliseconds between Ping and Pong
	maxPingLatency = int64(5000)

	//maxAvgLatency is the time in milliseconds between two peers
	maxAvgLatency = int64(2000)

	maxResetCacheSize = 512

	// lastBlockOffsetMs is the offset in milliseconds for the last block deadline
	// calculate. (200ms)
	lastBlockOffsetMs = 200 * time.Millisecond

	peerMsgQueueSize = 1024
	cbftVersion      = byte(0x01)

	maxBlockDist = uint64(192)

	msgQueuesLimit = 2048
)

type Cbft struct {
	config      *params.CbftConfig
	eventMux    *event.TypeMux
	handler     *handler
	closeOnce   sync.Once
	exitCh      chan struct{}
	txPool      *core.TxPool
	blockChain  *core.BlockChain //the block chain
	running     int32
	peerMsgCh   chan *MsgInfo
	syncBlockCh chan *BlockExt

	highestLogical   atomic.Value //highest block in logical path, local packages new block will base on it
	highestConfirmed atomic.Value //highest confirmed block in logical path
	rootIrreversible atomic.Value //the latest block has stored in chain

	executeBlockCh          chan *ExecuteBlockStatus
	baseBlockCh             chan chan *types.Block
	sealBlockCh             chan *SealBlock
	getBlockCh              chan *GetBlock
	sendViewChangeCh        chan struct{}
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
	Syncing

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

	// wal
	nodeServiceContext *node.ServiceContext
	wal                Wal
	loading            int32

	// validator
	agency     Agency
	validators *Validators

	startTimeOfEpoch int64

	evPool *EvidencePool
}

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, eventMux *event.TypeMux, ctx *node.ServiceContext) *Cbft {
	//todo need dynamic change consensus nodes
	initialNodesID := make([]discover.NodeID, 0, len(config.InitialNodes))
	for _, n := range config.InitialNodes {
		initialNodesID = append(initialNodesID, n.ID)
	}

	//dpos := newDpos(initialNodesID)

	cbft := &Cbft{
		config:   config,
		eventMux: eventMux,
		//dpos:                    dpos,
		//rotating:                newRotating(dpos, config.Duration),
		running:                 1,
		exitCh:                  make(chan struct{}),
		signedSet:               make(map[uint64]struct{}),
		syncBlockCh:             make(chan *BlockExt, peerMsgQueueSize),
		peerMsgCh:               make(chan *MsgInfo, peerMsgQueueSize),
		executeBlockCh:          make(chan *ExecuteBlockStatus),
		baseBlockCh:             make(chan chan *types.Block),
		sealBlockCh:             make(chan *SealBlock),
		getBlockCh:              make(chan *GetBlock),
		sendViewChangeCh:        make(chan struct{}),
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

	evPool, err := NewEvidencePool(ctx.ResolvePath(evidenceDir))
	if err != nil {
		return nil
	}
	cbft.evPool = evPool
	cbft.bp = defaultBP
	cbft.handler = NewHandler(cbft)
	cbft.router = NewRouter(cbft.handler)
	cbft.resetCache, _ = lru.New(maxResetCacheSize)

	return cbft
}

func (cbft *Cbft) getRootIrreversible() *BlockExt {
	if v := cbft.rootIrreversible.Load(); v == nil {
		panic("Get root block failed")
	} else {
		return v.(*BlockExt)
	}
}

func (cbft *Cbft) getHighestConfirmed() *BlockExt {
	if v := cbft.highestConfirmed.Load(); v == nil {
		panic("Get highest confirmed block failed")
	} else {
		return v.(*BlockExt)
	}
}
func (cbft *Cbft) getHighestLogical() *BlockExt {
	if v := cbft.highestLogical.Load(); v == nil {
		panic("Get highest logical block failed")
	} else {
		return v.(*BlockExt)
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

// Start sets blockChain and txPool into cbft
func (cbft *Cbft) Start(blockChain *core.BlockChain, txPool *core.TxPool, agency Agency) error {
	cbft.blockChain = blockChain
	cbft.startTimeOfEpoch = int64(blockChain.Genesis().Time().Uint64())

	cbft.agency = agency

	currentBlock := blockChain.CurrentBlock()

	var err error
	cbft.validators, err = cbft.agency.GetValidator(currentBlock.NumberU64())
	if err != nil {
		cbft.log.Error("Get validator fail", "error", err)
		return err
	}

	genesisParentHash := bytes.Repeat([]byte{0x00}, 32)
	if bytes.Equal(currentBlock.ParentHash().Bytes(), genesisParentHash) && currentBlock.Number() == nil {
		currentBlock.Header().Number = big.NewInt(0)
	}

	cbft.log.Debug("Init highestLogicalBlock", "hash", currentBlock.Hash(), "number", currentBlock.NumberU64())

	current := NewBlockExtBySeal(currentBlock, currentBlock.NumberU64(), cbft.nodeLength())
	current.number = currentBlock.NumberU64()

	if current.number > 0 && cbft.validators.Len() > 1 {
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
	//if cbft.wal, err = NewWal(cbft.nodeServiceContext); err != nil {
	//	return err
	//}
	cbft.wal = &emptyWal{}
	atomic.StoreInt32(&cbft.loading, 1)

	go cbft.receiveLoop()
	go cbft.executeBlockLoop()
	//start receive cbft message
	go cbft.handler.Start()
	go cbft.update()

	if err = cbft.wal.Load(cbft.AddJournal); err != nil {
		return err
	}
	atomic.StoreInt32(&cbft.loading, 0)
	return nil
}

func (cbft *Cbft) receiveLoop() {
	for {
		select {
		case msg := <-cbft.peerMsgCh:
			cbft.handleMsg(msg)
		case bt := <-cbft.syncBlockCh:
			cbft.OnSyncBlock(bt)
		case bs := <-cbft.executeBlockCh:
			cbft.OnExecutedBlock(bs)
		case shouldSeal := <-cbft.shouldSealCh:
			cbft.OnShouldSeal(shouldSeal)
		case view := <-cbft.viewChangeTimeoutCh:
			cbft.OnViewChangeTimeout(view)
		case <-cbft.sendViewChangeCh:
			cbft.OnSendViewChange()
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

	if !cbft.isRunning() {
		switch msg.(type) {
		case *prepareBlock,
			*prepareVote,
			*viewChange,
			*viewChangeVote:
			cbft.log.Debug("Cbft is not running, discard consensus message")
			return
		}
	}

	// write journal msg if cbft is not loading
	if !cbft.isLoading() {
		cbft.wal.Write(info)
	}

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
	}
	if err != nil {
		cbft.log.Error("Handle msg Failed", "error", err, "type", reflect.TypeOf(msg), "peer", peerID)
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
		validator, err := cbft.validators.NodeIndexAddress(cbft.config.NodeID)

		if err != nil {
			log.Debug("Get node index and address failed", "error", err)
			shouldSeal <- err
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
	cbft.bp.SyncBlockBP().SyncBlock(context.TODO(), ext, &cbft.RoundState)
	//todo verify block
	if ext.block.NumberU64() < cbft.getHighestConfirmed().number {
		cbft.log.Debug("Sync block too lower", "hash", ext.block.Hash(), "number", ext.number, "highest", cbft.getHighestConfirmed().number, "root", cbft.getRootIrreversible().number)
		ext.SetSyncState(nil)
		cbft.bp.SyncBlockBP().InvalidBlock(context.TODO(), ext, fmt.Errorf("sync block too lower"), &cbft.RoundState)
		return
	}

	if cbft.blockExtMap.findBlock(ext.block.Hash(), ext.block.NumberU64()) != nil {
		cbft.log.Debug("Sync block had exist", "hash", ext.block.Hash(), "number", ext.number, "highest", cbft.getHighestConfirmed().number, "root", cbft.getRootIrreversible().number)
		ext.SetSyncState(nil)
		cbft.bp.SyncBlockBP().InvalidBlock(context.TODO(), ext, fmt.Errorf("sync block had exist"), &cbft.RoundState)
		return
	}

	cbft.log.Debug("Sync block success", "hash", ext.block.Hash(), "number", ext.number)

	cbft.viewChange = ext.view
	if len(ext.viewChangeVotes) >= cbft.getThreshold() {
		if err := cbft.checkViewChangeVotes(ext.viewChangeVotes); err != nil {
			log.Error("Receive prepare invalid block", "err", err)
			cbft.bp.SyncBlockBP().InvalidBlock(context.TODO(), ext, err, &cbft.RoundState)
			ext.SetSyncState(err)
			return
		}
		for _, v := range ext.viewChangeVotes {
			cbft.viewChangeVotes[v.ValidatorAddr] = v
		}

		cbft.clearPending()
		cbft.ClearChildren(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum, cbft.viewChange.Timestamp)
		cbft.producerBlocks = NewProducerBlocks(cbft.validators.NodeID(int(ext.view.ProposalIndex)), ext.block.NumberU64())
		if cbft.producerBlocks != nil {
			cbft.producerBlocks.AddBlock(ext.block)
			cbft.log.Debug("Add producer block", "hash", ext.block.Hash(), "number", ext.block.Number(), "producer", cbft.producerBlocks.String())
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
			cbft.log.Debug("Send Block", "peer", peerID, "hash", g.Hash, "number", g.Number)
		}
	}
	return nil
}

func (cbft *Cbft) OnGetPrepareVote(peerID discover.NodeID, pv *getPrepareVote) error {
	ext := cbft.blockExtMap.findBlock(pv.Hash, pv.Number)
	votes := make([]*prepareVote, 0)

	if ext != nil {
		for i := uint32(0); i < pv.VoteBits.Size(); i++ {
			if pv.VoteBits.GetIndex(i) {
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
		for _, block := range msg.CommitedBlock {
			cbft.log.Debug("Sync Highest Block", "number", block.NumberU64())
			cbft.InsertChain(block, nil)
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
	cbft.log.Debug(fmt.Sprintf("Check view change timeout send:%v agree:%v", cbft.hadSendViewChange(), cbft.agreeViewChange()))
	if cbft.viewChange != nil && view.Equal(cbft.viewChange) {
		if cbft.hadSendViewChange() && !cbft.agreeViewChange() {
			cbft.handleCache()
			cbft.log.Info("View change timeout", "current view", cbft.viewChange.String())
			cbft.resetViewChange()
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
	if cbft.viewChange != nil && !view.EqualViewChange(cbft.viewChange) {
		if !cbft.agreeViewChange() {
			cbft.log.Warn("Waiting master response timeout", "view", cbft.viewChange.String())
			cbft.handleCache()
			cbft.resetViewChange()
			cbft.needPending = true
		}
	}

	viewChangeVoteTimeoutMeter.Mark(1)
}

func (cbft *Cbft) OnPrepareBlockHash(peerID discover.NodeID, msg *prepareBlockHash) error {
	cbft.log.Debug("Received message of prepareBlockHash", "FromPeerId", peerID.String(),
		"BlockHash", msg.Hash.Hex(), "Number", msg.Number)
	// Prerequisite: Nodes with PrepareBlock data can forward Hash
	cbft.handler.Send(peerID, &getPrepareBlock{Hash: msg.Hash, Number: msg.Number})

	// then: to forward msg
	if ok := cbft.needBroadcast(peerID, msg); ok {
		go cbft.handler.SendBroadcast(msg)
	}

	return nil
}

func (cbft *Cbft) NextBaseBlock() *types.Block {
	ch := make(chan *types.Block, 1)
	cbft.baseBlockCh <- ch
	return <-ch
}

func (cbft *Cbft) OnBaseBlock(ch chan *types.Block) {
	if cbft.master && cbft.agreeViewChange() && (cbft.producerBlocks == nil || len(cbft.producerBlocks.blocks) == 0) {
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

	cbft.log.Debug("Seal complete", "hash", sealedBlock.Hash(), "number", sealedBlock.NumberU64())

	cbft.bp.InternalBP().Seal(context.TODO(), current, &cbft.RoundState)
	cbft.bp.InternalBP().NewHighestLogicalBlock(context.TODO(), current, &cbft.RoundState)
	cbft.SetLocalHighestPrepareNum(current.number)
	if cbft.validators.Len() == 1 {
		cbft.log.Debug("Single node mode, confirm now")
		//only one consensus node, so, each block is highestConfirmed. (lock is needless)
		current.isConfirmed = true
		cbft.highestLogical.Store(current)
		cbft.highestConfirmed.Store(current)
		cbft.flushReadyBlock()
		return
	}

	//reset cbft.highestLogicalBlockExt cause this block is produced by myself
	cbft.highestLogical.Store(current)
	cbft.AddPrepareBlock(sealedBlock)

	cbft.broadcastBlock(current)
	//todo change sign and block state
	go func() {
		select {
		case <-stopCh:
			return
		case sealResultCh <- sealedBlock:
			//reset pool when seal block
			//start := time.Now()
			cbft.reset(sealedBlock)
			//cbft.bp.InternalBP().ResetTxPool(context.TODO(), current, time.Now().Sub(start), &cbft.RoundState)

		default:
			cbft.log.Warn("Sealing result is not ready by miner", "sealHash", sealedBlock.Header().SealHash())
		}
	}()
}

// ShouldSeal checks if it's local's turn to package new block at current time.
func (cbft *Cbft) ShouldSeal(curTime int64) (bool, error) {

	inturn := cbft.inTurn(curTime)
	if inturn {
		cbft.netLatencyLock.RLock()
		peersCount := len(cbft.netLatencyMap)
		cbft.netLatencyLock.RUnlock()
		if peersCount < cbft.getThreshold() {
			inturn = false
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

func (cbft *Cbft) sendViewChange() {
	cbft.sendViewChangeCh <- struct{}{}
}

func (cbft *Cbft) OnSendViewChange() {
	view, err := cbft.newViewChange()
	if err != nil {
		cbft.log.Error("New view change failed", "err", err)
		return
	}
	cbft.log.Debug("Send new view", "view", view.String())
	cbft.handler.SendAllConsensusPeer(view)

	// gauage
	blockHighNumConfirmedGauage.Update(int64(cbft.getHighestConfirmed().number))
	blockHighNumLogicGauage.Update(int64(cbft.getHighestLogical().number))

	time.AfterFunc(time.Duration(cbft.config.Period)*time.Second, func() {
		cbft.viewChangeTimeoutCh <- view
	})
}

// Receive view from other nodes
// Need verify timestamp , signature, promise highest confirmed block
func (cbft *Cbft) OnViewChange(peerID discover.NodeID, view *viewChange) error {
	cbft.log.Debug("Receive view change", "peer", peerID, "view", view.String())

	if cbft.viewChange != nil && cbft.viewChange.Equal(view) {
		cbft.log.Debug("Duplication view change message, discard this")
		return nil
	}

	bpCtx := context.WithValue(context.Background(), "peer", peerID)
	cbft.bp.ViewChangeBP().ReceiveViewChange(bpCtx, view, &cbft.RoundState)
	if err := cbft.VerifyAndViewChange(view); err != nil {
		if view.BaseBlockNum > cbft.getHighestConfirmed().number {
			if view.BaseBlockNum-cbft.getHighestConfirmed().number > maxBlockDist {
				atomic.StoreInt32(&cbft.running, 0)
			} else {
				cbft.log.Warn(fmt.Sprintf("Local is too slower, need to sync block to %s", peerID.TerminalString()))

				cbft.handler.Send(peerID, &getHighestPrepareBlock{Lowest: cbft.getRootIrreversible().number + 1})
			}
		}

		cbft.bp.ViewChangeBP().InvalidViewChange(bpCtx, view, err, &cbft.RoundState)
		cbft.log.Error("Verify view failed", "err", err, "peer", peerID, "view", view.String())
		return err
	}

	validator, err := cbft.validators.NodeIndexAddress(cbft.config.NodeID)
	if err != nil {
		cbft.bp.ViewChangeBP().InvalidViewChange(bpCtx, view, errInvalidatorCandidateAddress, &cbft.RoundState)
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

	time.AfterFunc(time.Duration(cbft.config.Period)*time.Second*2, func() {
		cbft.viewChangeVoteTimeoutCh <- resp
	})
	cbft.setViewChange(view)
	cbft.bp.InternalBP().SwitchView(bpCtx, view)
	cbft.bp.ViewChangeBP().SendViewChangeVote(bpCtx, resp, &cbft.RoundState)
	cbft.handler.SendAllConsensusPeer(view)
	cbft.handler.SendAllConsensusPeer(resp)

	//cbft.handler.Send(peerID, cbft.viewChangeResp)
	return nil

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
		cbft.log.Error("Flush block error")
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
	cbft.bp.InternalBP().NewHighestRootBlock(context.TODO(), newRoot, &cbft.RoundState)

	blockConfirmedTimer.UpdateSince(time.Unix(int64(newRoot.timestamp), 0))

	cbft.evPool.Clear(cbft.viewChange.Timestamp, cbft.viewChange.BaseBlockNum)
	return true

}

// Receive prepare block from the other consensus node.
// Need check something ,such as validator index, address, view is equal local view , and last verify signature
func (cbft *Cbft) OnNewPrepareBlock(nodeId discover.NodeID, request *prepareBlock, propagation bool) error {
	bpCtx := context.WithValue(context.TODO(), "peer", nodeId)
	cbft.bp.PrepareBP().ReceiveBlock(bpCtx, request, &cbft.RoundState)

	//discard block when view.Timestamp != request.Timestamp && request.BlockNum > view.BlockNum
	if cbft.viewChange != nil && len(request.ViewChangeVotes) < cbft.getThreshold() && request.Timestamp != cbft.viewChange.Timestamp && request.Block.NumberU64() > cbft.viewChange.BaseBlockNum {
		log.Debug("Invalid prepare block", "number", request.Block.NumberU64(), "hash", request.Block.Hash(), "view", cbft.viewChange)
		return errFutileBlock
	}

	if err := cbft.VerifyHeader(cbft.blockChain, request.Block.Header(), false); err != nil {
		cbft.bp.PrepareBP().InvalidBlock(bpCtx, request, err, &cbft.RoundState)
		log.Error("Failed to verify header in PrepareBlockMsg, discard this msg", "peer", nodeId, "err", err)
		return nil
	}

	ext := cbft.blockExtMap.findBlock(request.Block.Hash(), request.Block.NumberU64())
	if cbft.blockChain.HasBlock(request.Block.Hash(), request.Block.NumberU64()) || (ext != nil && ext.block != nil) {
		log.Warn("Block already in blockchain, discard this msg", "prepare block", request.String())
		return nil
	}

	if !cbft.IsConsensusNode() {
		log.Warn("Local node is not consensus node,discard this msg")
		return errInvalidatorCandidateAddress
	} else if !cbft.CheckConsensusNode(request.ProposalAddr) {
		cbft.bp.PrepareBP().InvalidBlock(bpCtx, request,
			fmt.Errorf("remote node is not consensus node addr:%s", request.ProposalAddr.String()), &cbft.RoundState)
		log.Warn("Remote node is not consensus node,discard this msg", "addr", request.ProposalAddr)
		return errInvalidatorCandidateAddress
	}

	cbft.log.Debug("Receive prepare block", "number", request.Block.NumberU64(), "prepareVotes", len(request.ViewChangeVotes))
	ext = NewBlockExtByPrepareBlock(request, cbft.nodeLength())

	if len(request.ViewChangeVotes) != 0 && request.View != nil {
		if len(request.ViewChangeVotes) < cbft.getThreshold() {
			cbft.bp.PrepareBP().InvalidBlock(bpCtx, request, errTwoThirdPrepareVotes, &cbft.RoundState)
			cbft.log.Error(fmt.Sprintf("Receive not enough prepareVotes %d threshold %d", len(request.ViewChangeVotes), cbft.getThreshold()))
			return errTwoThirdViewchangeVotes
		}

		if cbft.getHighestLogical().number < request.View.BaseBlockNum ||
			cbft.blockExtMap.findBlock(request.View.BaseBlockHash, request.View.BaseBlockNum) == nil {
			cbft.bp.PrepareBP().InvalidBlock(bpCtx, request, errNotFoundViewBlock, &cbft.RoundState)
			cbft.handler.Send(nodeId, &getHighestPrepareBlock{Lowest: cbft.getRootIrreversible().number + 1})
			cbft.log.Error(fmt.Sprintf("View Block is not found, hash:%s, number:%d", request.View.BaseBlockHash.TerminalString(), request.View.BaseBlockNum))
			return errNotFoundViewBlock
		}

		oldViewChange := cbft.viewChange
		viewChange := request.View
		if cbft.viewChange == nil || cbft.viewChange.Timestamp <= viewChange.Timestamp {
			cbft.log.Debug("New PrepareBlock is not match current view, need change")
			cbft.viewChange = viewChange
		}
		if err := cbft.checkViewChangeVotes(request.ViewChangeVotes); err != nil {
			cbft.bp.PrepareBP().InvalidViewChangeVote(bpCtx, request, err, &cbft.RoundState)
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
			cbft.bp.ViewChangeBP().TwoThirdViewChangeVotes(bpCtx, &cbft.RoundState)
			var newHeader *types.Header
			viewBlock := cbft.blockExtMap.findBlock(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum)

			if viewBlock == nil {
				cbft.bp.ViewChangeBP().InvalidViewChangeBlock(bpCtx, cbft.viewChange, &cbft.RoundState)
				log.Error("ViewChange block find error", "BaseBlockHash", cbft.viewChange.BaseBlockHash,
					"BaseBlockNum", cbft.viewChange.BaseBlockNum, "blockMap", cbft.blockExtMap.BlockString())
				cbft.handler.Send(nodeId, &getPrepareBlock{Hash: cbft.viewChange.BaseBlockHash, Number: cbft.viewChange.BaseBlockNum})
				//panic("Find nil block")
			} else {
				newHeader = viewBlock.block.Header()
				injectBlock := cbft.blockExtMap.findBlockByNumber(cbft.viewChange.BaseBlockNum+1, cbft.getHighestLogical().number)
				start := time.Now()
				cbft.txPool.ForkedReset(newHeader, injectBlock)
				cbft.bp.InternalBP().ForkedResetTxPool(bpCtx, newHeader, injectBlock, time.Now().Sub(start), &cbft.RoundState)

			}

			cbft.clearPending()
			cbft.ClearChildren(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum, cbft.viewChange.Timestamp)
		}
		ext.view = cbft.viewChange
		ext.viewChangeVotes = request.ViewChangeVotes
	}

	switch cbft.AcceptPrepareBlock(request) {
	case Accept:
		cbft.bp.PrepareBP().AcceptBlock(bpCtx, request, &cbft.RoundState)
		if cbft.producerBlocks != nil {
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
		cbft.bp.PrepareBP().CacheBlock(bpCtx, request, &cbft.RoundState)
		cbft.log.Info("Cache block", "hash", ext.block.Hash(), "number", ext.block.NumberU64())
	case Discard:
		cbft.bp.PrepareBP().DiscardBlock(bpCtx, request, &cbft.RoundState)
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

	ext.prepareVotes.Add(vote)

	cbft.saveBlockExt(vote.Hash, ext)

	//receive enough signature broadcast
	if ext.inTree && ext.isExecuted && ext.isConfirmed {
		cbft.bp.PrepareBP().TwoThirdVotes(context.TODO(), ext, &cbft.RoundState)
		if h := cbft.blockExtMap.FindHighestConfirmedWithHeader(); h != nil {
			cbft.bp.InternalBP().NewHighestConfirmedBlock(context.TODO(), ext, &cbft.RoundState)
			cbft.highestConfirmed.Store(h)
			blockConfirmedMeter.Mark(1)
			blockConfirmedTimer.UpdateSince(time.Unix(int64(ext.timestamp), 0))
			cbft.flushReadyBlock()
			cbft.updateValidator()
		}
		cbft.log.Debug("Send Confirmed Block", "hash", ext.block.Hash(), "number", ext.block.NumberU64())
		cbft.handler.SendAllConsensusPeer(&confirmedPrepareBlock{Hash: ext.block.Hash(), Number: ext.block.NumberU64(), VoteBits: ext.prepareVotes.voteBits})
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
			cbft.bp.InternalBP().ResetTxPool(context.TODO(), bs.block, time.Now().Sub(start), &cbft.RoundState)
			cbft.highestLogical.Store(bs.block)
			cbft.bp.InternalBP().NewHighestLogicalBlock(context.TODO(), bs.block, &cbft.RoundState)
			cbft.sendPrepareVote(bs.block)
			cbft.bp.PrepareBP().SendPrepareVote(context.TODO(), bs.block, &cbft.RoundState)

			if bs.block.isConfirmed {
				cbft.highestConfirmed.Store(bs.block)
				cbft.bp.InternalBP().NewHighestConfirmedBlock(context.TODO(), bs.block, &cbft.RoundState)
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

	validator, err := cbft.validators.NodeIndexAddress(cbft.config.NodeID)
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
			cbft.SetLocalHighestPrepareNum(pv.Number)
			pv.Signature.SetBytes(sign)
			if cbft.viewChange != nil && !cbft.agreeViewChange() && cbft.viewChange.BaseBlockNum < ext.block.NumberU64() {
				cbft.pendingVotes.Add(pv.Hash, pv)
			} else {
				ext.prepareVotes.Add(pv)
				cbft.blockExtMap.Add(pv.Hash, pv.Number, ext)
				cbft.log.Debug("Broadcast prepare vote", "vote", pv.String())
				cbft.handler.SendAllConsensusPeer(pv)
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
			cbft.bp.InternalBP().InvalidBlock(context.TODO(), ext.block.Hash(), ext.block.NumberU64(), err)
		}
		cbft.bp.InternalBP().ExecuteBlock(context.TODO(), ext.block.Hash(), ext.block.NumberU64(), time.Now().Sub(start))
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
	_, err := cbft.validators.AddressIndex(address)
	return err == nil
}

// IsConsensusNode check if local is a consensus node.
func (cbft *Cbft) IsConsensusNode() bool {
	_, err := cbft.validators.NodeIndex(cbft.config.NodeID)
	return err == nil
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	cbft.log.Trace("Verify header", "hash", header.Hash(), "number", header.Number.Uint64(), "seal", seal)

	if header.Number == nil {
		cbft.log.Warn(fmt.Sprintf("Verify header failed, unknow block"))
		return errUnknownBlock
	}

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
		return errors.New("execute block error")
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

func (cbft *Cbft) CalcBlockDeadline() (time.Time, error) {
	node, err := cbft.validators.NodeIndex(cbft.config.NodeID)
	if err != nil {
		return time.Time{}, err
	}
	nodeIdx := node.Index
	startEpoch := cbft.startTimeOfEpoch * 1000
	timePoint := time.Now().UnixNano() / int64(time.Millisecond)

	if nodeIdx >= 0 {
		if cbft.validators.Len() == 1 {
			return time.Now().Add(time.Duration(cbft.config.Period) * time.Second), err
		}
		durationPerNode := cbft.config.Duration * 1000
		durationPerTurn := durationPerNode * int64(cbft.validators.Len())

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

		remaing := time.Duration(nextSlotValue-value) * time.Millisecond
		if lastBlock {
			if remaing > lastBlockOffsetMs {
				remaing = remaing - lastBlockOffsetMs
			} else {
				remaing = 50 * time.Millisecond // 50ms
			}
		}
		return time.Now().Add(remaing), err
	}
	return time.Now().Add(50 * time.Millisecond), err // 50ms

}

func (cbft *Cbft) CalcNextBlockTime() (time.Time, error) {
	vn, err := cbft.validators.NodeIndex(cbft.config.NodeID)
	if err != nil {
		return time.Time{}, err
	}
	nodeIdx := vn.Index
	startEpoch := cbft.startTimeOfEpoch * 1000
	timePoint := time.Now().UnixNano() / int64(time.Millisecond)

	if nodeIdx >= 0 {
		if cbft.validators.Len() == 1 {
			return time.Now().Add(time.Duration(cbft.config.Period) * time.Second), nil
		}
		durationPerNode := cbft.config.Duration * 1000
		durationPerTurn := durationPerNode * int64(cbft.validators.Len())

		min := int64(nodeIdx) * (durationPerNode)
		value := (timePoint - startEpoch) % durationPerTurn
		max := int64(nodeIdx+1) * durationPerNode

		log.Trace("Calc next block time", "min", min, "value", value, "max", max)

		var offset int64
		if value >= min && value <= max {
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
			last := int64(cbft.validators.Len()) * durationPerNode
			offset = last - value + min
		}
		return time.Now().Add(time.Duration(offset) * time.Millisecond), nil
	}
	return time.Now().Add(time.Duration(cbft.config.Period) * time.Second), nil // 1s

}

// ConsensusNodes returns all consensus nodes.
func (cbft *Cbft) ConsensusNodes() ([]discover.NodeID, error) {
	cbft.log.Trace(fmt.Sprintf("dposNodeCount:%d", cbft.validators.Len()))
	return cbft.validators.NodeList(), nil
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
	return []rpc.API{{
		Namespace: "cbft",
		Version:   "1.0",
		Service:   &API{chain: chain, cbft: cbft},
		Public:    false,
	}}
}

func (cbft *Cbft) Protocols() []p2p.Protocol {
	return cbft.handler.Protocols()
}

// OnBlockSignature is called by by protocol handler when it received a new block signature by P2P.
func (cbft *Cbft) OnPrepareVote(peerID discover.NodeID, vote *prepareVote, propagation bool) error {
	cbft.log.Debug("Receive prepare vote", "peer", peerID, "vote", vote.String())
	bpCtx := context.WithValue(context.Background(), "peer", peerID)
	cbft.bp.PrepareBP().ReceiveVote(bpCtx, vote, &cbft.RoundState)
	err := cbft.verifyValidatorSign(cbft.viewChange.BaseBlockNum, vote.ValidatorIndex, vote.ValidatorAddr, vote, vote.Signature[:])
	if err != nil {
		cbft.bp.PrepareBP().InvalidVote(bpCtx, vote, err, &cbft.RoundState)
		cbft.log.Error("Verify vote error", "err", err)
		return err
	}

	switch cbft.AcceptPrepareVote(vote) {
	case Accept:
		cbft.log.Debug("Accept block vote", "vote", vote.String())
		cbft.bp.PrepareBP().AcceptVote(bpCtx, vote, &cbft.RoundState)
		if err := cbft.evPool.AddPrepareVote(vote); err != nil {
			if _, ok := err.(*DuplicatePrepareVoteEvidence); ok {
				cbft.log.Warn("Receive DuplicatePrepareVoteEvidence msg", "err", err.Error())
				return err
			}
		}

		cbft.prepareVoteReceiver(peerID, vote)
	case Cache:
		cbft.log.Debug("View changing, add vote into process queue", "vote", vote.String())
		cbft.bp.PrepareBP().CacheVote(bpCtx, vote, &cbft.RoundState)
		//changing view
		cbft.AddProcessingVote(peerID, vote)
	case Discard:
		cbft.bp.PrepareBP().DiscardVote(bpCtx, vote, &cbft.RoundState)
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
	nodeIdx, err := cbft.validators.AddressIndex(addr)
	if err != nil {
		cbft.log.Error("Get address index failed", "err", err)
		return false
	}
	return cbft.calTurnIndex(rcvTime, nodeIdx.Index)
}

func (cbft *Cbft) calTurn(timePoint int64, nodeID discover.NodeID) bool {
	vn, err := cbft.validators.NodeIndex(nodeID)
	if err != nil {
		return false
	}
	return cbft.calTurnIndex(timePoint, vn.Index)
}

func (cbft *Cbft) calTurnIndex(timePoint int64, nodeIdx int) bool {

	startEpoch := cbft.startTimeOfEpoch * 1000

	if nodeIdx >= 0 {
		if cbft.validators.Len() == 1 {
			return true
		}
		durationPerNode := cbft.config.Duration * 1000
		durationPerTurn := durationPerNode * int64(cbft.validators.Len())

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
func ecrecover(header *types.Header) (discover.NodeID, []byte, error) {
	var nodeID discover.NodeID
	if len(header.Extra) < extraSeal {
		return nodeID, []byte{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]
	sealHash := header.SealHash()

	pubkey, err := crypto.Ecrecover(sealHash.Bytes(), signature)
	if err != nil {
		return nodeID, []byte{}, err
	}

	nodeID, err = discover.BytesID(pubkey[1:])
	if err != nil {
		return nodeID, []byte{}, err
	}
	return nodeID, signature, nil
}

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
	trunc := cbft.validators.Len() * 2 / 3
	return trunc
}

func (cbft *Cbft) nodeLength() int {
	return cbft.validators.Len()
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

	if current.number > 0 && cbft.validators.Len() > 1 {
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
	oldVds := cbft.validators
	cbft.validators = newVds
	cbft.log.Info("Update validators success", "highestConfirmed", hc.number, "hash", hc.block.Hash(), "validators", cbft.validators)

	if _, ok := cbft.validators.Nodes[cbft.config.NodeID]; ok {
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
			for nodeID, _ := range cbft.validators.Nodes {
				if _, ok := oldVds.Nodes[nodeID]; !ok {
					cbft.eventMux.Post(cbfttypes.AddValidatorEvent{NodeID: nodeID})
					cbft.log.Trace("Post AddValidatorEvent", "nodeID", nodeID.String())
				}
			}

			for nodeID, _ := range oldVds.Nodes {
				if _, ok := cbft.validators.Nodes[nodeID]; !ok {
					cbft.eventMux.Post(cbfttypes.RemoveValidatorEvent{NodeID: nodeID})
					cbft.log.Trace("Post RemoveValidatorEvent", "nodeID", nodeID.String())
				}
			}
		} else {
			for nodeID, _ := range oldVds.Nodes {
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
			for nodeID, _ := range cbft.validators.Nodes {
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
	//isCsusNode := cbft.IsConsensusNode() fmt.Sprintf("%x", p.ID().Bytes()[:8]),
	peers := cbft.handler.peers.Peers()
	for _, peer := range peers {
		if peer.knownMessageHash.Contains(msg.MsgHash()) {
			cbft.log.Debug("needn't to broadcast", "type", reflect.TypeOf(msg), "hash", msg.MsgHash(), "BHash", msg.BHash().TerminalString())
			messageRepeatMeter.Mark(1)
			return false
		}
	}
	cbft.log.Debug("need to broadcast", "type", reflect.TypeOf(msg), "hash", msg.MsgHash(), "BHash", msg.BHash().TerminalString())
	messageGossipMeter.Mark(1)
	return true
}

func (cbft *Cbft) AddJournal(msg *MsgInfo) {
	cbft.log.Debug("Method:LoadPeerMsg received message from peer", "peer", msg.PeerID.TerminalString(), "msgType", reflect.TypeOf(msg.Msg), "msgHash", msg.Msg.MsgHash().TerminalString(), "BHash", msg.Msg.BHash().TerminalString())
	cbft.handleMsg(msg)
}
