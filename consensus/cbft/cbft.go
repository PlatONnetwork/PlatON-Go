package cbft

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
	"sync"
	"time"
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
	peerMsgCh  chan *MsgInfo
	syncMsgCh  chan *MsgInfo
	evPool     EvidencePool
	log        log.Logger

	//Control the current view state
	state viewState

	//Block executor, the block responsible for executing the current view
	executor blockExecutor

	//Verification security rules for proposed blocks and viewchange
	safetyRules safetyRules

	//Determine when to allow voting
	voteRules voteRules

	//Store blocks that are not committed
	blockTree blockTree
}

func New(sysConfig *params.CbftConfig, optConfig *OptionsConfig, eventMux *event.TypeMux, ctx *node.ServiceContext) *Cbft {
	cbft := &Cbft{
		config:    Config{sysConfig, optConfig},
		eventMux:  eventMux,
		exitCh:    make(chan struct{}),
		peerMsgCh: make(chan *MsgInfo, optConfig.PeerMsgQueueSize),
		syncMsgCh: make(chan *MsgInfo, optConfig.PeerMsgQueueSize),
		log:       log.New(),
	}

	if evPool, err := NewEvidencePool(); err == nil {
		cbft.evPool = evPool
	} else {
		return nil
	}

	//todo init safety rules, vote rules, state, executor

	return cbft
}

func (cbft *Cbft) Start(chain consensus.ChainReader, executor consensus.Executor, pool consensus.TxPoolReset, agency consensus.Agency) error {
	panic("implement me")
}

func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return nil
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

func (Cbft) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	panic("implement me")
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

func (Cbft) ShouldSeal(curTime int64) (bool, error) {
	panic("implement me")
}

func (Cbft) CalcBlockDeadline(timePoint int64) (time.Time, error) {
	panic("implement me")
}

func (Cbft) CalcNextBlockTime(timePoint int64) (time.Time, error) {
	panic("implement me")
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
