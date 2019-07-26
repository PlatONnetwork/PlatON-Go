package cbft

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	errors "github.com/pkg/errors"

	"reflect"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	cconsensus "github.com/PlatONnetwork/PlatON-Go/common/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/executor"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/fetcher"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/network"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/rules"
	cstate "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
)

const cbftVersion = 1

type Cbft struct {
	config           ctypes.Config
	eventMux         *event.TypeMux
	closeOnce        sync.Once
	exitCh           chan struct{}
	txPool           consensus.TxPoolReset
	blockChain       consensus.ChainReader
	blockCacheWriter consensus.BlockCacheWriter
	peerMsgCh        chan *ctypes.MsgInfo
	syncMsgCh        chan *ctypes.MsgInfo
	evPool           evidence.EvidencePool
	log              log.Logger
	network          *network.EngineManager

	start    bool
	syncing  bool
	fetching bool
	// Async call channel
	asyncCallCh chan func()

	fetcher *fetcher.Fetcher
	// Control the current view state
	state *cstate.ViewState

	// Block asyncExecutor, the block responsible for executing the current view
	asyncExecutor executor.AsyncBlockExecutor

	// Verification security rules for proposed blocks and viewchange
	safetyRules rules.SafetyRules

	// Determine when to allow voting
	voteRules rules.VoteRules

	// Validator pool
	validatorPool *validator.ValidatorPool

	// Store blocks that are not committed
	blockTree *ctypes.BlockTree

	// wal
	nodeServiceContext *node.ServiceContext
	wal                wal.Wal
	loading            int32

	// Record the number of peer requests for obtaining cbft information.
	queues map[string]int // Per peer message counts to prevent memory exhaustion.
}

func New(sysConfig *params.CbftConfig, optConfig *ctypes.OptionsConfig, eventMux *event.TypeMux, ctx *node.ServiceContext) *Cbft {
	cbft := &Cbft{
		config:             ctypes.Config{Sys: sysConfig, Option: optConfig},
		eventMux:           eventMux,
		exitCh:             make(chan struct{}),
		peerMsgCh:          make(chan *ctypes.MsgInfo, optConfig.PeerMsgQueueSize),
		syncMsgCh:          make(chan *ctypes.MsgInfo, optConfig.PeerMsgQueueSize),
		log:                log.New(),
		start:              false,
		syncing:            false,
		fetching:           false,
		asyncCallCh:        make(chan func(), optConfig.PeerMsgQueueSize),
		nodeServiceContext: ctx,
		queues:             make(map[string]int),
		state:              cstate.NewViewState(),
	}

	if evPool, err := evidence.NewEvidencePool(ctx, optConfig.EvidenceDir); err == nil {
		cbft.evPool = evPool
	} else {
		return nil
	}

	return cbft
}

// Returns the ID value of the current node
func (cbft *Cbft) NodeId() discover.NodeID {
	return discover.NodeID{}
}

func (cbft *Cbft) Start(chain consensus.ChainReader, blockCacheWriter consensus.BlockCacheWriter, txPool consensus.TxPoolReset, agency consensus.Agency) error {
	cbft.blockChain = chain
	cbft.txPool = txPool
	cbft.blockCacheWriter = blockCacheWriter
	cbft.asyncExecutor = executor.NewAsyncExecutor(blockCacheWriter.Execute)
	cbft.validatorPool = validator.NewValidatorPool(agency, chain.CurrentHeader().Number.Uint64(), cbft.config.Option.NodeID)

	cbft.state = cstate.NewViewState()
	//Initialize block tree
	block := chain.GetBlock(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64())

	isGenesis := func() bool {
		return block.NumberU64() == 0
	}

	var qc *ctypes.QuorumCert
	if !isGenesis() {
		var err error
		_, qc, err = ctypes.DecodeExtra(block.ExtraData())

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("start cbft failed"))
		}
	}

	cbft.blockTree = ctypes.NewBlockTree(block, qc)
	atomic.StoreInt32(&cbft.loading, 1)
	if isGenesis() {
		cbft.changeView(cbft.config.Sys.Epoch, 1, block, qc, nil)
	} else {
		cbft.changeView(qc.Epoch, qc.ViewNumber, block, qc, nil)
	}

	//Initialize view state
	cbft.state.SetHighestExecutedBlock(block)
	cbft.state.SetHighestQCBlock(block)
	cbft.state.SetHighestLockBlock(block)
	cbft.state.SetHighestCommitBlock(block)

	//Initialize rules
	cbft.safetyRules = rules.NewSafetyRules(cbft.state, cbft.blockTree)
	cbft.voteRules = rules.NewVoteRules(cbft.state)

	// load consensus state
	if err := cbft.LoadWal(); err != nil {
		return err
	}
	atomic.StoreInt32(&cbft.loading, 0)

	go cbft.receiveLoop()

	// init handler and router to process message.
	// cbft -> handler -> router.
	cbft.network = network.NewEngineManger(cbft) // init engineManager as handler.

	// Start the handler to process the message.
	go cbft.network.Start()

	cbft.start = true
	cbft.log.Info("Cbft engine start")
	return nil
}

// Entrance: The messages related to the consensus are entered from here.
// The message sent from the peer node is sent to the CBFT message queue and
// there is a loop that will distribute the incoming message.
func (cbft *Cbft) ReceiveMessage(msg *ctypes.MsgInfo) {
	select {
	case cbft.peerMsgCh <- msg:
		cbft.log.Debug("Received message from peer", "peer", msg.PeerID, "msgType", reflect.TypeOf(msg.Msg), "msgHash", msg.Msg.MsgHash().TerminalString(), "BHash", msg.Msg.BHash().TerminalString())
	case <-cbft.exitCh:
		cbft.log.Error("Cbft exit")
	}
}

// ReceiveSyncMsg is used to receive messages that are synchronized from other nodes.
//
// Possible message types are:
//  PrepareBlockVotesMsg/GetLatestStatusMsg/LatestStatusMsg/
func (cbft *Cbft) ReceiveSyncMsg(msg *ctypes.MsgInfo) {
	select {
	case cbft.syncMsgCh <- msg:
		cbft.log.Debug("Receive synchronization related messages from peer", "peer", msg.PeerID, "msgType", reflect.TypeOf(msg.Msg), "msgHash", msg.Msg.MsgHash().TerminalString(), "BHash", msg.Msg.BHash().TerminalString())
	case <-cbft.exitCh:
		cbft.log.Error("Cbft exit")
	}
}

// LoadWal tries to recover consensus state and view msg from the wal.
func (cbft *Cbft) LoadWal() error {
	// init wal and load wal state
	var err error
	if cbft.wal, err = wal.NewWal(cbft.nodeServiceContext, ""); err != nil {
		return err
	}
	//cbft.wal = &emptyWal{}
	// load consensus chainState
	if err = cbft.wal.LoadChainState(cbft.recoveryChainState); err != nil {
		cbft.log.Error(err.Error())
		return err
	}
	// load consensus message
	if err = cbft.wal.Load(cbft.recoveryMsg); err != nil {
		cbft.log.Error(err.Error())
		return err
	}
	return nil
}

//Receive all consensus related messages, all processing logic in the same goroutine
func (cbft *Cbft) receiveLoop() {
	// channel Divided into read-only type, writable type
	// Read-only is the channel that gets the current CBFT status.
	// Writable type is the channel that affects the consensus state.
	for {
		select {
		case msg := <-cbft.peerMsgCh:

			// Prevent Dos attacks and limit the number of messages sent by each node.
			count := cbft.queues[msg.PeerID] + 1
			if int64(count) > cbft.config.Option.MaxQueuesLimit {
				log.Error("Discarded message, exceeded allowance for the layer of cbft", "peer", msg.PeerID, "msgHash", msg.Msg.MsgHash().TerminalString())
				// Need further confirmation.
				// todo: Is the program exiting or dropping the message here?
				break
			}
			cbft.queues[msg.PeerID] = count
			cbft.handleConsensusMsg(msg)
			// After the message is processed, the counter is decremented by one.
			// If it is reduced to 0, the mapping relationship of the corresponding
			// node will be deleted.
			cbft.queues[msg.PeerID]--
			if cbft.queues[msg.PeerID] == 0 {
				delete(cbft.queues, msg.PeerID)
			}

		case msg := <-cbft.syncMsgCh:
			cbft.handleSyncMsg(msg)
		case msg := <-cbft.asyncExecutor.ExecuteStatus():
			cbft.onAsyncExecuteStatus(msg)
		case fn := <-cbft.asyncCallCh:
			fn()

		case <-cbft.state.ViewTimeout():
			cbft.OnViewTimeout()
		default:
		}

		// read-only channel
		select {
		default:
		}
	}
}

//Handling consensus messages, there are three main types of messages. prepareBlock, prepareVote, viewChange
func (cbft *Cbft) handleConsensusMsg(info *ctypes.MsgInfo) {
	if cbft.running() {
		return
	}
	msg, id := info.Msg, info.PeerID
	var err error

	// Forward the message before processing the message.
	go cbft.network.Forwarding(id, msg)

	switch msg := msg.(type) {
	case *protocols.PrepareBlock:
		err = cbft.OnPrepareBlock(id, msg)
	case *protocols.PrepareVote:
		err = cbft.OnPrepareVote(id, msg)
	case *protocols.ViewChange:
		err = cbft.OnViewChange(id, msg)
	}

	if err != nil {
		cbft.log.Error("Handle msg Failed", "error", err, "type", reflect.TypeOf(msg), "peer", id)
	}
}

// Behind the node will be synchronized by synchronization message
func (cbft *Cbft) handleSyncMsg(info *ctypes.MsgInfo) {
	msg, id := info.Msg, info.PeerID

	switch msg := msg.(type) {
	case *protocols.GetPrepareBlock:
		cbft.OnGetPrepareBlock(id, msg)
	case *protocols.GetBlockQuorumCert:
		cbft.OnGetBlockQuorumCert(id, msg)
	case *protocols.BlockQuorumCert:
		cbft.OnBlockQuorumCert(id, msg)
	case *protocols.GetQCBlockList:
		cbft.OnGetQCBlockList(id, msg)
	default:
		cbft.fetcher.MatchTask(id, msg)
	}
}

func (cbft *Cbft) running() bool {
	return !cbft.syncing && !cbft.fetching
}

func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	if header.Number == nil {
		cbft.log.Error("Verify header fail, unknown block")
		return errors.New("unknown block")
	}

	cbft.log.Trace("Verify header", "number", header.Number, "hash", header.Hash, "seal", seal)
	if len(header.Extra) < consensus.ExtraSeal {
		cbft.log.Error("Verify header fail, missing signature", "number", header.Number, "hash", header.Hash)
	}

	if err := cbft.validatorPool.VerifyHeader(header); err != nil {
		cbft.log.Error("Verify header fail", "number", header.Number, "hash", header.Hash(), "err", err)
	}
	return nil
}

func (Cbft) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	panic("implement me")
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (cbft *Cbft) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	cbft.log.Trace("Verify seal", "hash", header.Hash(), "number", header.Number)
	if header.Number.Uint64() == 0 {
		return errors.New("unknown block")
	}
	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header of running the transactions on top.
func (cbft *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	cbft.log.Debug("Prepare", "hash", header.Hash(), "number", header.Number.Uint64())

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
	cbft.log.Debug("Finalize block", "hash", header.Hash(), "number", header.Number, "txs", len(txs), "receipts", len(receipts))
	header.Root = state.IntermediateRoot(true)
	return types.NewBlock(header, txs, receipts), nil
}

func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	cbft.log.Info("Seal block", "number", block.Number(), "parentHash", block.ParentHash())
	header := block.Header()
	if block.NumberU64() == 0 {
		return errors.New("unknown block")
	}

	sign, err := cbft.signFn(header.SealHash().Bytes())
	if err != nil {
		cbft.log.Error("Seal block sign fail", "number", block.Number(), "parentHash", block.ParentHash(), "err", err)
		return err
	}

	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sign[:])

	sealBlock := block.WithSeal(header)

	cbft.asyncCallCh <- func() {
		cbft.OnSeal(sealBlock, results, stop)
	}
	return nil
}

func (cbft *Cbft) OnSeal(block *types.Block, results chan<- *types.Block, stop <-chan struct{}) {
	// TODO: check is turn to seal block

	if cbft.state.HighestQCBlock().Hash() != block.ParentHash() ||
		cbft.state.HighestExecutedBlock().Hash() != block.ParentHash() {
		cbft.log.Warn("Futile block cause highest executed block changed", "nubmer", block.Number(), "parentHash", block.ParentHash(),
			"qcNumber", cbft.state.HighestQCBlock().Number(), "qcHash", cbft.state.HighestQCBlock().Hash(),
			"exectedNumber", cbft.state.HighestExecutedBlock().Number(), "exectedHash", cbft.state.HighestExecutedBlock().Hash())
		return
	}

	// TODO: seal process
	prepareBlock := &protocols.PrepareBlock{
		Epoch:      cbft.state.Epoch(),
		ViewNumber: cbft.state.ViewNumber(),
		Block:      block,
		BlockIndex: cbft.state.NumViewBlocks(),
	}

	if cbft.state.NumViewBlocks() == 0 {
		parentBlock, parentQC := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
		if parentBlock == nil {
			cbft.log.Error("Can not find parent block", "number", block.Number(), "parentHash", block.ParentHash())
			return
		}
		prepareBlock.PrepareQC = parentQC
	}

	cbft.log.Info("Seal New Block", "prepareBlock", prepareBlock.String())

	// TODO: signature block - fake verify.
	if err := cbft.signMsgByBls(prepareBlock); err != nil {
		cbft.log.Error("Sign PrepareBlock failed", "err", err, "hash", block.Hash(), "number", block.NumberU64())
		return
	}

	cbft.state.SetExecuting(prepareBlock.BlockIndex, true)

	if err := cbft.OnPrepareBlock("", prepareBlock); err != nil {
		cbft.log.Error("Check Seal Block failed", "err", err, "hash", block.Hash(), "number", block.NumberU64())
		return
	}

	if err := cbft.signBlock(block.Hash(), block.NumberU64(), prepareBlock.BlockIndex); err != nil {
		cbft.log.Error("Sign PrepareBlock failed", "err", err, "hash", block.Hash(), "number", block.NumberU64())
		return
	}

	cbft.findQCBlock()

	cbft.state.SetHighestExecutedBlock(block)

	// write sendPrepareBlock info to wal
	cbft.sendPrepareBlock(prepareBlock)

	cbft.network.Broadcast(prepareBlock)

	go func() {
		select {
		case <-stop:
			return
		case results <- block:
		default:
			cbft.log.Warn("Sealing result channel is not ready by miner", "sealHash", block.Header().SealHash())
		}
	}()
}

// SealHash returns the hash of a block prior to it being sealed.
func (cbft *Cbft) SealHash(header *types.Header) common.Hash {
	cbft.log.Debug("Seal hash", "hash", header.Hash(), "number", header.Number)
	return header.SealHash()
}

func (Cbft) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{}
}

func (cbft *Cbft) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

func (cbft *Cbft) NextBaseBlock() *types.Block {
	result := make(chan *types.Block, 1)
	cbft.asyncCallCh <- func() {
		block := cbft.state.HighestExecutedBlock()
		cbft.log.Debug("Base block", "hash", block.Hash(), "number", block.Number())
		result <- block
	}
	return <-result
}

func (Cbft) InsertChain(block *types.Block, errCh chan error) {
	panic("implement me")
}

// HashBlock check if the specified block exists in block tree.
func (cbft *Cbft) HasBlock(hash common.Hash, number uint64) bool {
	has := false
	cbft.checkStart(func() {
		if cbft.state.HighestExecutedBlock().NumberU64() >= number {
			has = true
		}
	})

	return has
}

func (Cbft) Status() string {
	panic("implement me")
}

// GetBlockByHash get the specified block by hash.
func (cbft *Cbft) GetBlockByHash(hash common.Hash) *types.Block {
	result := make(chan *types.Block, 1)
	cbft.asyncCallCh <- func() {
		block := cbft.blockTree.FindBlockByHash(hash)
		result <- block
	}
	return <-result
}

// CurrentBlock get the current lock block.
func (cbft *Cbft) CurrentBlock() *types.Block {
	var block *types.Block
	cbft.checkStart(func() {
		block = cbft.state.HighestLockBlock()
	})
	return block
}

func (cbft *Cbft) checkStart(exe func()) {
	if cbft.start {
		exe()
	}
}

func (cbft *Cbft) FastSyncCommitHead() <-chan error {
	result := make(chan error, 1)

	cbft.asyncCallCh <- func() {
		currentBlock := cbft.blockChain.GetBlock(cbft.blockChain.CurrentHeader().Hash(), cbft.blockChain.CurrentHeader().Number.Uint64())

		// TODO: update view
		cbft.state.SetHighestExecutedBlock(currentBlock)
		cbft.state.SetHighestQCBlock(currentBlock)
		cbft.state.SetHighestLockBlock(currentBlock)
		cbft.state.SetHighestCommitBlock(currentBlock)

		result <- nil
	}
	return result
}

func (cbft *Cbft) Close() error {
	cbft.log.Info("Close cbft consensus")
	cbft.start = false
	cbft.closeOnce.Do(func() {
		// Short circuit if the exit channel is not allocated.
		if cbft.exitCh == nil {
			return
		}
		close(cbft.exitCh)
	})
	if cbft.asyncExecutor != nil {
		cbft.asyncExecutor.Stop()
	}
	return nil
}

func (cbft *Cbft) ConsensusNodes() ([]discover.NodeID, error) {
	return cbft.validatorPool.ValidatorList(cbft.state.HighestQCBlock().NumberU64()), nil
}

// ShouldSeal check if we can seal block.
func (cbft *Cbft) ShouldSeal(curTime time.Time) (bool, error) {
	if cbft.isLoading() {
		return false, nil
	}
	currentExecutedBlockNumber := cbft.state.HighestExecutedBlock().NumberU64()
	if !cbft.validatorPool.IsValidator(currentExecutedBlockNumber, cbft.config.Option.NodeID) {
		return false, errors.New("current node not a validator")
	}

	result := make(chan error, 2)
	cbft.asyncCallCh <- func() {
		cbft.OnShouldSeal(result)
	}
	select {
	case err := <-result:
		return err == nil, err
	case <-time.After(2 * time.Millisecond):
		result <- errors.New("timeout")
		return false, errors.New("CBFT engine busy")
	}
}

func (cbft *Cbft) OnShouldSeal(result chan error) {
	// todo: need add remark.
	select {
	case <-result:
		cbft.log.Trace("Should seal timeout")
		return
	default:
	}

	currentExecutedBlockNumber := cbft.state.HighestExecutedBlock().NumberU64()
	if !cbft.validatorPool.IsValidator(currentExecutedBlockNumber, cbft.config.Option.NodeID) {
		result <- errors.New("current node not a validator")
		return
	}

	numValidators := cbft.validatorPool.Len(currentExecutedBlockNumber)
	currentProposer := cbft.state.ViewNumber() % uint64(numValidators)
	validator, err := cbft.validatorPool.GetValidatorByNodeID(currentExecutedBlockNumber, cbft.config.Option.NodeID)
	if err != nil {
		cbft.log.Error("Should seal fail", "err", err)
		result <- err
		return
	}
	if currentProposer != uint64(validator.Index) {
		result <- errors.New("current node not the proposer")
		return
	}

	if cbft.state.NumViewBlocks() >= cbft.config.Sys.Amount {
		result <- errors.New("produce block over limit")
		return
	}
	result <- nil
}

func (cbft *Cbft) CalcBlockDeadline(timePoint time.Time) time.Time {
	produceInterval := time.Duration(cbft.config.Sys.Period/uint64(cbft.config.Sys.Amount)) * time.Millisecond
	cbft.log.Debug("Calc block deadline", "timePoint", timePoint, "stateDeadline", cbft.state.Deadline(), "produceInterval", produceInterval)
	if cbft.state.Deadline().Sub(timePoint) > produceInterval {
		return timePoint.Add(produceInterval)
	}
	return cbft.state.Deadline()
}

func (cbft *Cbft) CalcNextBlockTime(blockTime time.Time) time.Time {
	produceInterval := time.Duration(cbft.config.Sys.Period/uint64(cbft.config.Sys.Amount)) * time.Millisecond
	cbft.log.Debug("Calc next block time",
		"blockTime", blockTime, "now", time.Now(), "produceInterval", produceInterval,
		"period", cbft.config.Sys.Period, "amount", cbft.config.Sys.Amount,
		"interval", time.Since(blockTime))
	if time.Since(blockTime) < produceInterval {
		// TODO: add network latency
		return time.Now().Add(produceInterval - time.Since(blockTime))
	}
	return time.Now()
}

func (cbft *Cbft) IsConsensusNode() bool {
	return cbft.validatorPool.IsValidator(cbft.state.HighestQCBlock().NumberU64(), cbft.config.Option.NodeID)
}

func (cbft *Cbft) GetBlock(hash common.Hash, number uint64) *types.Block {
	result := make(chan *types.Block, 1)
	cbft.asyncCallCh <- func() {
		block, _ := cbft.blockTree.FindBlockAndQC(hash, number)
		result <- block
	}
	return <-result
}

func (cbft *Cbft) GetBlockWithoutLock(hash common.Hash, number uint64) *types.Block {
	block, _ := cbft.blockTree.FindBlockAndQC(hash, number)
	return block
}

func (Cbft) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	//panic("implement me")
}

func (cbft *Cbft) IsSignedBySelf(sealHash common.Hash, header *types.Header) bool {
	return cbft.verifySelfSigned(sealHash.Bytes(), header.Signature())
}

func (Cbft) TracingSwitch(flag int8) {
	panic("implement me")
}

func (cbft *Cbft) OnPong(nodeID discover.NodeID, netLatency int64) error {
	//panic("need to be improved")
	return nil
}

func (cbft *Cbft) Config() *ctypes.Config {
	panic("need to be improved")
	return nil
}

// Return the highest submitted block number of the current node.
func (cbft *Cbft) HighestCommitBlockBn() uint64 {
	return cbft.state.HighestQCBlock().NumberU64()
}

// Return the highest locked block number of the current node.
func (cbft *Cbft) HighestLockBlockBn() uint64 {
	return cbft.state.HighestLockBlock().NumberU64()
}

// Return the highest QC block number of the current node.
func (cbft *Cbft) HighestQCBlockBn() uint64 {
	return cbft.state.HighestQCBlock().NumberU64()
}

func (cbft *Cbft) threshold(num int) int {
	return num - (num-1)/3
}

func (cbft *Cbft) commitBlock(block *types.Block, qc *ctypes.QuorumCert) {
	extra, err := ctypes.EncodeExtra(byte(cbftVersion), qc)
	if err != nil {
		cbft.log.Error("Encode extra error", "nubmer", block.Number(), "hash", block.Hash(), "cbftVersion", cbftVersion)
		return
	}

	cbft.log.Debug("Send consensus result to worker", "number", block.Number(), "hash", block.Hash())
	cbft.eventMux.Post(cbfttypes.CbftResult{
		Block:     block,
		ExtraData: extra,
		SyncState: nil,
	})
}

func (cbft *Cbft) Evidences() string {
	evs := cbft.evPool.Evidences()
	if len(evs) == 0 {
		return "{}"
	}
	evds := evidence.ClassifyEvidence(evs)
	js, err := json.MarshalIndent(evds, "", "  ")
	if err != nil {
		return ""
	}
	return string(js)
}

func (cbft *Cbft) UnmarshalEvidence(data []byte) (cconsensus.Evidences, error) {
	return cbft.evPool.UnmarshalEvidence(data)
}

// verifySelfSigned
func (cbft *Cbft) verifySelfSigned(m []byte, sig []byte) bool {
	recPubKey, err := crypto.Ecrecover(m, sig)
	if err != nil {
		return false
	}

	pubKey := cbft.config.Option.NodePriKey.PublicKey
	pbytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return false
	}
	return true
}

// signFn use private key to sign byte slice.
func (cbft *Cbft) signFn(m []byte) ([]byte, error) {
	return crypto.Sign(m, cbft.config.Option.NodePriKey)
}

// signFn use bls private key to sign byte slice.
func (cbft *Cbft) signFnByBls(m []byte) ([]byte, error) {
	sign := cbft.config.Option.BlsPriKey.Sign(string(m))
	return sign.Serialize(), nil
}

// signMsg use bls private key to sign msg.
func (cbft *Cbft) signMsgByBls(msg ctypes.ConsensusMsg) error {
	buf, err := msg.CannibalizeBytes()
	if err != nil {
		return err
	}
	sign, err := cbft.signFnByBls(buf)
	if err != nil {
		return err
	}
	msg.SetSign(sign)
	return nil
}

func (cbft *Cbft) isLoading() bool {
	return atomic.LoadInt32(&cbft.loading) == 1
}

func (cbft *Cbft) verifyConsensusMsg(msg ctypes.ConsensusMsg) (*cbfttypes.ValidateNode, error) {
	digest, err := msg.CannibalizeBytes()
	if err != nil {
		return nil, errors.Wrap(err, "get msg's cannibalize bytes failed")
	}

	if !cbft.validatorPool.Verify(msg.BlockNum(), msg.NodeIndex(), digest, msg.Sign()) {
		return nil, fmt.Errorf("signature verification failed")
	}

	vnode, err := cbft.validatorPool.GetValidatorByIndex(cbft.state.HighestQCBlock().NumberU64(), msg.NodeIndex())

	if err != nil {
		return nil, errors.Wrap(err, "get validator failed")
	}

	var prepareQC *ctypes.QuorumCert

	switch cm := msg.(type) {
	case *protocols.PrepareBlock:
		prepareQC = cm.PrepareQC
		if cm.ViewChangeQC != nil {
			if err := cbft.verifyViewChangeQC(cm.ViewChangeQC); err != nil {
				return nil, err
			}
		}
	case *protocols.PrepareVote:
		prepareQC = cm.ParentQC
	case *protocols.ViewChange:
		prepareQC = cm.PrepareQC
	}

	if err := cbft.verifyPrepareQC(prepareQC); err != nil {
		return nil, err
	}

	return vnode, nil
}

func (cbft *Cbft) generatePrepareQC(votes map[uint32]*protocols.PrepareVote) *ctypes.QuorumCert {
	if len(votes) == 0 {
		return nil
	}

	var vote *protocols.PrepareVote

	for _, v := range votes {
		vote = v
	}

	// Validator set prepareQC is the same as highestQC
	total := cbft.validatorPool.Len(vote.BlockNum())

	vSet := utils.NewBitArray(uint32(total))
	vSet.SetIndex(vote.NodeIndex(), true)

	var aggSig bls.Sign
	if err := aggSig.Deserialize(vote.Sign()); err != nil {
		return nil
	}

	qc := &ctypes.QuorumCert{
		Epoch:       vote.Epoch,
		ViewNumber:  vote.ViewNumber,
		BlockHash:   vote.BlockHash,
		BlockNumber: vote.BlockNumber,
		BlockIndex:  vote.BlockIndex,
	}
	for _, p := range votes {
		if p.NodeIndex() != vote.NodeIndex() {
			var sig bls.Sign
			err := sig.Deserialize(vote.Sign())
			if err != nil {
				return nil
			}

			aggSig.Add(&sig)
			vSet.SetIndex(p.NodeIndex(), true)
		}
	}
	qc.Signature.SetBytes(aggSig.Serialize())
	qc.ValidatorSet.Update(vSet)
	return qc
}

func (cbft *Cbft) generateViewChangeQC(viewChanges map[uint32]*protocols.ViewChange) *ctypes.ViewChangeQC {
	type ViewChangeQC struct {
		cert   *ctypes.ViewChangeQuorumCert
		aggSig *bls.Sign
		ba     *utils.BitArray
	}

	total := uint32(cbft.validatorPool.Len(cbft.state.HighestQCBlock().NumberU64()))

	qcs := make(map[common.Hash]*ViewChangeQC)

	for _, v := range viewChanges {
		var aggSig bls.Sign
		if err := aggSig.Deserialize(v.Sign()); err != nil {
			return nil
		}

		if vc, ok := qcs[v.BlockHash]; !ok {
			qc := &ViewChangeQC{
				cert: &ctypes.ViewChangeQuorumCert{
					Epoch:       v.Epoch,
					ViewNumber:  v.ViewNumber,
					BlockHash:   v.BlockHash,
					BlockNumber: v.BlockNumber,
				},
				aggSig: &aggSig,
				ba:     utils.NewBitArray(total),
			}
			qcs[v.BlockHash] = qc
		} else {
			vc.aggSig.Add(&aggSig)
			vc.ba.SetIndex(v.NodeIndex(), true)
		}
	}

	qc := &ctypes.ViewChangeQC{QCs: make([]*ctypes.ViewChangeQuorumCert, 0)}
	for _, q := range qcs {
		q.cert.Signature.SetBytes(q.aggSig.Serialize())
		q.cert.ValidatorSet.Update(q.ba)
		qc.QCs = append(qc.QCs, q.cert)
	}
	return qc
}

func (cbft *Cbft) verifyPrepareQC(qc *ctypes.QuorumCert) error {
	var cb []byte
	var err error
	if cb, err = qc.CannibalizeBytes(); err != nil {
		return err
	}
	if cbft.validatorPool.VerifyAggSigByBA(qc.BlockNumber, qc.ValidatorSet, cb, qc.Signature.Bytes()) {
		return fmt.Errorf("verify prepare qc failed")
	}
	return err
}

func (cbft *Cbft) verifyViewChangeQC(viewChangeQC *ctypes.ViewChangeQC) error {
	var err error
	for _, vc := range viewChangeQC.QCs {
		var cb []byte
		if cb, err = vc.CannibalizeBytes(); err != nil {
			break
		}
		if cbft.validatorPool.VerifyAggSigByBA(vc.BlockNumber, vc.ValidatorSet, cb, vc.Signature.Bytes()) {
			err = fmt.Errorf("verify viewchange qc failed")
			break
		}
	}

	return err
}
