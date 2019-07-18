package cbft

import (
	"crypto/ecdsa"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/executor"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/fetcher"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/rules"
	cstate "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
)

type Config struct {
	sys    *params.CbftConfig
	option *OptionsConfig
}

type Cbft struct {
	config     Config
	eventMux   *event.TypeMux
	closeOnce  sync.Once
	exitCh     chan struct{}
	txPool     consensus.TxPoolReset
	blockChain consensus.ChainReader
	peerMsgCh  chan *ctypes.MsgInfo
	syncMsgCh  chan *ctypes.MsgInfo
	evPool     evidence.EvidencePool
	log        log.Logger

	// Async call channel
	asyncCallCh chan func()

	fetcher *fetcher.Fetcher
	// Control the current view state
	state cstate.ViewState

	// Execution block function
	execute consensus.Executor
	// Block asyncExecutor, the block responsible for executing the current view
	asyncExecutor executor.AsyncBlockExecutor

	// Verification security rules for proposed blocks and viewchange
	safetyRules rules.SafetyRules

	// Determine when to allow voting
	voteRules rules.VoteRules

	// Validator pool
	validatorPool *validator.ValidatorPool

	// Store blocks that are not committed
	blockTree ctypes.BlockTree
}

func New(sysConfig *params.CbftConfig, optConfig *OptionsConfig, eventMux *event.TypeMux, ctx *node.ServiceContext) *Cbft {
	cbft := &Cbft{
		config:      Config{sysConfig, optConfig},
		eventMux:    eventMux,
		exitCh:      make(chan struct{}),
		peerMsgCh:   make(chan *ctypes.MsgInfo, optConfig.PeerMsgQueueSize),
		syncMsgCh:   make(chan *ctypes.MsgInfo, optConfig.PeerMsgQueueSize),
		log:         log.New(),
		asyncCallCh: make(chan func(), optConfig.PeerMsgQueueSize),
	}

	if evPool, err := evidence.NewEvidencePool(); err == nil {
		cbft.evPool = evPool
	} else {
		return nil
	}

	//todo init safety rules, vote rules, state, asyncExecutor
	cbft.safetyRules = rules.NewSafetyRules(&cbft.state)
	cbft.voteRules = rules.NewVoteRules(&cbft.state)

	return cbft
}

func (cbft *Cbft) Start(chain consensus.ChainReader, executorFn consensus.Executor, txPool consensus.TxPoolReset, agency consensus.Agency) error {
	cbft.blockChain = chain
	cbft.txPool = txPool
	cbft.asyncExecutor = executor.NewAsyncExecutor(executorFn)
	cbft.validatorPool = validator.NewValidatorPool(agency, chain.CurrentHeader().Number.Uint64(), cbft.config.sys.NodeID)

	//Initialize block tree
	block := chain.GetBlock(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64())

	cbft.blockTree.InsertQCBlock(block, nil)

	//Initialize view state
	cbft.state.SetHighestExecutedBlock(block)
	cbft.state.SetHighestQCBlock(block)
	cbft.state.SetHighestLockBlock(block)
	cbft.state.SetHighestCommitBlock(block)
	go cbft.receiveLoop()
	return nil
}

//Receive all consensus related messages, all processing logic in the same goroutine
func (cbft *Cbft) receiveLoop() {
	// channel Divided into read-only type, writable type
	// Read-only is the channel that gets the current CBFT status.
	// Writable type is the channel that affects the consensus state

	for {
		select {
		case msg := <-cbft.peerMsgCh:
			cbft.handleConsensusMsg(msg)
		case msg := <-cbft.syncMsgCh:
			cbft.handleSyncMsg(msg)

		case fn := <-cbft.asyncCallCh:
			fn()
		default:
		}

		// read-only channel
		select {}
	}
}

//Handling consensus messages, there are three main types of messages. prepareBlock, prepareVote, viewchagne
func (cbft *Cbft) handleConsensusMsg(info *ctypes.MsgInfo) {
	msg, peerID := info.Msg, info.PeerID
	var err error

	switch msg := msg.(type) {
	case *protocols.PrepareBlock:
		err = cbft.OnPrepareBlock(msg)
	case *protocols.PrepareVote:
		err = cbft.OnPrepareVote(msg)
	case *protocols.ViewChange:
		err = cbft.OnViewChange(msg)
	}

	if err != nil {
		cbft.log.Error("Handle msg Failed", "error", err, "type", reflect.TypeOf(msg), "peer", peerID)
	}
}

// Behind the node will be synchronized by synchronization message
func (cbft *Cbft) handleSyncMsg(info *ctypes.MsgInfo) {
	msg, peerID := info.Msg, info.PeerID

	if cbft.fetcher.MatchTask(peerID.String(), msg) {
		return
	}

	var err error
	switch msg.(type) {
	}

	if err != nil {
		cbft.log.Error("Handle msg Failed", "error", err, "type", reflect.TypeOf(msg), "peer", peerID)
	}
}

func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return cbft.validatorPool.VerifyHeader(header)
}

func (Cbft) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	panic("implement me")
}

func (Cbft) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	panic("implement me")
}

func (Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	panic("implement me")
}

func (Cbft) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	receipts []*types.Receipt) (*types.Block, error) {
	panic("implement me")
}

func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	cbft.log.Info("Seal block", "number", block.Number(), "parentHash", block.ParentHash())
	if block.NumberU64() == 0 {
		return errors.New("unknow block")
	}

	// TODO signature block

	cbft.asyncCallCh <- func() {
		cbft.OnSeal(block, results, stop)
	}
	return nil
}

func (cbft *Cbft) OnSeal(block *types.Block, results chan<- *types.Block, stop <-chan struct{}) {
	// TODO: check is turn to seal block

	if cbft.state.HighestExecutedBlock().Hash() != block.ParentHash() {
		cbft.log.Warn("Futile block cause highest executed block changed", "nubmer", block.Number(), "parentHash", block.ParentHash())
		return
	}

	me, _ := cbft.validatorPool.GetValidatorByNodeID(cbft.state.HighestQCBlock().NumberU64(), cbft.config.sys.NodeID)

	// TODO: seal process
	prepareBlock := &protocols.PrepareBlock{
		Epoch:         cbft.state.Epoch(),
		ViewNumber:    cbft.state.ViewNumber(),
		Block:         block,
		BlockIndex:    cbft.state.NumViewBlocks(),
		ProposalIndex: uint32(me.Index),
		ProposalAddr:  me.Address,
	}

	if cbft.state.NumViewBlocks() == 0 {
	}

	// TODO: add viewchange qc

	// TODO: signature block

	cbft.state.AddPrepareBlock(prepareBlock)

	// TODO: broadcast block

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

func (Cbft) SealHash(header *types.Header) common.Hash {
	panic("implement me")
}

func (Cbft) APIs(chain consensus.ChainReader) []rpc.API {
	panic("implement me")
}

func (Cbft) Protocols() []p2p.Protocol {
	panic("implement me")
}

func (Cbft) NextBaseBlock() *types.Block {
	panic("implement me")
}

func (Cbft) InsertChain(block *types.Block, errCh chan error) {
	panic("implement me")
}

func (Cbft) HasBlock(hash common.Hash, number uint64) bool {
	panic("implement me")
}

func (Cbft) Status() string {
	panic("implement me")
}

func (Cbft) GetBlockByHash(hash common.Hash) *types.Block {
	panic("implement me")
}

func (Cbft) CurrentBlock() *types.Block {
	panic("implement me")
}

func (Cbft) FastSyncCommitHead() <-chan error {
	panic("implement me")
}

func (Cbft) Close() error {
	panic("implement me")
}

func (Cbft) ConsensusNodes() ([]discover.NodeID, error) {
	panic("implement me")
}

// ShouldSeal check if we can seal block.
func (cbft *Cbft) ShouldSeal(curTime time.Time) (bool, error) {
	result := make(chan error, 1)
	// FIXME: should use independent channel?
	cbft.asyncCallCh <- func() {
		cbft.OnShouldSeal(result)
	}
	select {
	case err := <-result:
		return err == nil, err
	case <-time.After(2 * time.Millisecond):
		return false, errors.New("CBFT engine busy")
	}
}

func (cbft *Cbft) OnShouldSeal(result chan error) {
	currentExecutedBlockNumber := cbft.state.HighestExecutedBlock().NumberU64()
	if !cbft.validatorPool.IsValidator(currentExecutedBlockNumber, cbft.config.sys.NodeID) {
		result <- errors.New("current node not a validator")
		return
	}

	numValidators := cbft.validatorPool.Len(currentExecutedBlockNumber)
	currentProposer := cbft.state.ViewNumber() % uint64(numValidators)
	validator, _ := cbft.validatorPool.GetValidatorByNodeID(currentExecutedBlockNumber, cbft.config.sys.NodeID)
	if currentProposer != uint64(validator.Index) {
		result <- errors.New("current node not the proposer")
		return
	}

	if cbft.state.NumViewBlocks() >= cbft.config.sys.Amount {
		result <- errors.New("produce block over limit")
		return
	}
	result <- nil
}

func (cbft *Cbft) CalcBlockDeadline(timePoint time.Time) time.Time {
	// FIXME: condition race
	produceInterval := time.Duration(cbft.config.sys.Period/uint64(cbft.config.sys.Amount)) * time.Millisecond
	if cbft.state.Deadline().Sub(timePoint) > produceInterval {
		return timePoint.Add(produceInterval)
	}
	return cbft.state.Deadline()
}

func (cbft *Cbft) CalcNextBlockTime(blockTime time.Time) time.Time {
	// FIXME: condition race
	produceInterval := time.Duration(cbft.config.sys.Period/uint64(cbft.config.sys.Amount)) * time.Millisecond
	if time.Now().Sub(blockTime) < produceInterval {
		// TODO: add network latency
		return time.Now().Add(time.Now().Sub(blockTime))
	}
	return time.Now()
}

func (Cbft) IsConsensusNode() bool {
	panic("implement me")
}

func (Cbft) GetBlock(hash common.Hash, number uint64) *types.Block {
	panic("implement me")
}

func (Cbft) GetBlockWithoutLock(hash common.Hash, number uint64) *types.Block {
	panic("implement me")
}

func (Cbft) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	panic("implement me")
}

func (Cbft) IsSignedBySelf(sealHash common.Hash, signature []byte) bool {
	panic("implement me")
}

func (Cbft) Evidences() string {
	panic("implement me")
}

func (Cbft) TracingSwitch(flag int8) {
	panic("implement me")
}

func (cbft *Cbft) OnPong(nodeID discover.NodeID, netLatency int64) error {
	panic("need to be improved")
	return nil
}
