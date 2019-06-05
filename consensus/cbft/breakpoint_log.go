package cbft

import (
	"context"
	"encoding/json"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"math/big"
	"sync/atomic"
	"time"
)

const (
	flagState = byte(1)
	flagStat  = byte(2)
)

type Context struct {
	//TraceID represents globally unique ID of the trace, such view's timestamp
	TraceID uint64 `json:"trace_id"`

	// SpanID represents span ID that must be unique within its trace, such as peerID, blockNum, baseBlock
	// but does not have to be globally unique.
	SpanID string `json:"span_id"`

	// ParentID refers to the ID of the parent span.
	// Should be "" if the current span is a root span.
	ParentID string `json:"parent_id"`

	// Log type such as "state", "stat"
	Flags byte `json:"flags"`

	//message signer
	Creator string `json:"creator"`

	//local node
	Processor string `json:"processor"`
}
type Tag struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type LogRecord struct {
	Timestamp int64       `json:"timestamp"`
	Log       interface{} `json:"log"`
}

type Span struct {
	Context      Context       `json:"context"`
	StartTime    time.Time     `json:"start_time"`
	DurationTime time.Duration `json:"duration_time"`
	Tags         []Tag         `json:"tags"`
	LogRecords   []LogRecord   `json:"log_records"`
	//operation name, such as message type
	OperationName string `json:"operation_time"`
}

var logBP Breakpoint
var localAddr atomic.Value
var localID atomic.Value

func init() {
	logBP = &defaultBreakpoint{
		prepareBP:    new(logPrepareBP),
		viewChangeBP: new(logViewChangeBP),
		syncBlockBP:  new(logSyncBlockBP),
		internalBP:   new(logInternalBP),
	}
}

type logPrepareBP struct {
}

func localAddress(cbft *Cbft) string {
	addr := ""
	if v := localAddr.Load(); v == nil {
		if cbft != nil {
			pub, _ := cbft.config.NodeID.Pubkey()
			addr = crypto.PubkeyToAddress(*pub).String()
			localID.Store(cbft.config.NodeID.String())
			localAddr.Store(addr)
		}
	} else {
		addr = v.(string)
	}
	return addr
}

func localNodeID() string {
	if v := localID.Load(); v != nil {
		return v.(string)
	}
	return ""
}

func (bp logPrepareBP) ReceiveBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   block.Timestamp,
			SpanID:    block.Block.Number().String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   block.ProposalAddr.String(),
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "prepare_block",
		Tags: []Tag{
			{
				Key:   "peer_id",
				Value: ctx.Value("peer"),
			},
			{
				Key:   "action",
				Value: "receive_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       block,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logPrepareBP) ReceiveVote(ctx context.Context, vote *prepareVote, cbft *Cbft) {
	log.Debug("ReceiveVote", "block", vote.String(), "cbft", cbft.String())

}

func (bp logPrepareBP) AcceptBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   block.Timestamp,
			SpanID:    block.Block.Number().String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   block.ProposalAddr.String(),
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "prepare_block",
		Tags: []Tag{
			{
				Key:   "peer_id",
				Value: ctx.Value("peer"),
			},
			{
				Key:   "action",
				Value: "accept_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       block,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logPrepareBP) CacheBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   block.Timestamp,
			SpanID:    block.Block.Number().String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   block.ProposalAddr.String(),
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "prepare_block",
		Tags: []Tag{
			{
				Key:   "peer_id",
				Value: ctx.Value("peer"),
			},
			{
				Key:   "action",
				Value: "cache_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       block,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logPrepareBP) DiscardBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   block.Timestamp,
			SpanID:    block.Block.Number().String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   block.ProposalAddr.String(),
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "discard_prepare_block",
		Tags: []Tag{
			{
				Key:   "peer_id",
				Value: ctx.Value("peer"),
			},
			{
				Key:   "action",
				Value: "discard_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       block,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logPrepareBP) AcceptVote(ctx context.Context, vote *prepareVote, cbft *Cbft) {
	log.Debug("AcceptVote", "block", vote.String(), "cbft", cbft.String())
}

func (bp logPrepareBP) CacheVote(ctx context.Context, vote *prepareVote, cbft *Cbft) {
	log.Debug("CacheVote", "block", vote.String(), "cbft", cbft.String())
}

func (bp logPrepareBP) DiscardVote(ctx context.Context, vote *prepareVote, cbft *Cbft) {
	log.Debug("DiscardVote", "block", vote.String(), "cbft", cbft.String())
}

func (bp logPrepareBP) SendPrepareVote(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	log.Debug("SendPrepareVote", "block", ext.String(), "cbft", cbft.String())
}

func (bp logPrepareBP) InvalidBlock(ctx context.Context, block *prepareBlock, err error, cbft *Cbft) {
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   block.Timestamp,
			SpanID:    block.Block.Number().String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   block.ProposalAddr.String(),
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "prepare_block",
		Tags: []Tag{
			{
				Key:   "peer_id",
				Value: ctx.Value("peer"),
			},
			{
				Key:   "action",
				Value: "invalid_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       block,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logPrepareBP) InvalidVote(ctx context.Context, vote *prepareVote, err error, cbft *Cbft) {
	log.Debug("InvalidVote", "block", vote.String(), "cbft", cbft.String())
}

func (bp logPrepareBP) InvalidViewChangeVote(ctx context.Context, block *prepareBlock, err error, cbft *Cbft) {
	log.Debug("InvalidViewChangeVote", "block", block.String(), "cbft", cbft.String())
}

func (bp logPrepareBP) TwoThirdVotes(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	log.Debug("TwoThirdVotes", "block", ext.String(), "cbft", cbft.String())
}

type logViewChangeBP struct {
}

func (bp logViewChangeBP) ReceiveViewChange(ctx context.Context, view *viewChange, cbft *Cbft) {
	log.Debug("ReceiveViewChange", "block", view.String(), "cbft", cbft.String())
}

func (bp logViewChangeBP) ReceiveViewChangeVote(ctx context.Context, vote *viewChangeVote, cbft *Cbft) {
	log.Debug("ReceiveViewChangeVote", "vote", vote.String(), "cbft", cbft.String())
}

func (bp logViewChangeBP) InvalidViewChange(ctx context.Context, view *viewChange, err error, cbft *Cbft) {
	log.Debug("InvalidViewChange", "view", view.String(), "cbft", cbft.String())
}

func (bp logViewChangeBP) InvalidViewChangeVote(ctx context.Context, view *viewChangeVote, err error, cbft *Cbft) {
	log.Debug("InvalidViewChangeVote", "view", view.String(), "cbft", cbft.String())
}

func (bp logViewChangeBP) InvalidViewChangeBlock(ctx context.Context, view *viewChange, cbft *Cbft) {
	log.Debug("InvalidViewChangeBlock", "view", view.String(), "cbft", cbft.String())
}

func (bp logViewChangeBP) TwoThirdViewChangeVotes(ctx context.Context, cbft *Cbft) {
	log.Debug("TwoThirdViewChangeVotes", "cbft", cbft.String())
}

func (bp logViewChangeBP) SendViewChangeVote(ctx context.Context, vote *viewChangeVote, cbft *Cbft) {
	log.Debug("SendViewChangeVote", "vote", vote.String(), "cbft", cbft.String())

}

func (bp logViewChangeBP) ViewChangeTimeout(ctx context.Context, cbft *Cbft) {
	log.Debug("ViewChangeTimeout", "cbft", cbft.String())

}

type logSyncBlockBP struct {
}

func (bp logSyncBlockBP) SyncBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	processor := localAddress(cbft)

	creator := ""
	if ext.view != nil {
		creator = ext.view.ProposalAddr.String()
	}
	span := &Span{
		Context: Context{
			TraceID:   ext.timestamp,
			SpanID:    ext.block.Number().String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   creator,
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "sync_block",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "sync_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       ext,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}

}

func (bp logSyncBlockBP) InvalidBlock(ctx context.Context, ext *BlockExt, err error, cbft *Cbft) {
	processor := localAddress(cbft)

	creator := ""
	if ext.view != nil {
		creator = ext.view.ProposalAddr.String()
	}
	span := &Span{
		Context: Context{
			TraceID:   ext.timestamp,
			SpanID:    ext.block.Number().String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   creator,
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "sync_block",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "sync_invalid_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       ext,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}

}

type logInternalBP struct {
}

func (bp logInternalBP) ExecuteBlock(ctx context.Context, hash common.Hash, number uint64, timestamp uint64, elapse time.Duration) {
	type HashNumber struct {
		Hash   common.Hash `json:"block_hash"`
		Number uint64      `json:"block_number"`
	}

	span := &Span{
		Context: Context{
			TraceID:   timestamp,
			SpanID:    big.NewInt(int64(number)).String(),
			ParentID:  localNodeID(),
			Creator:   "",
			Processor: localAddress(nil),
		},
		StartTime:     time.Now(),
		DurationTime:  elapse,
		OperationName: "execute_block",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "execute_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log: &HashNumber{
					Hash:   hash,
					Number: number,
				},
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logInternalBP) InvalidBlock(ctx context.Context, hash common.Hash, number uint64, timestamp uint64, err error) {
	type HashNumber struct {
		Hash   common.Hash `json:"block_hash"`
		Number uint64      `json:"block_number"`
		Error  string      `json:"error"`
	}

	span := &Span{
		Context: Context{
			TraceID:   timestamp,
			SpanID:    big.NewInt(int64(number)).String(),
			ParentID:  localNodeID(),
			Creator:   "",
			Processor: localAddress(nil),
		},
		StartTime:     time.Now(),
		OperationName: "execute_block",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "execute_invalid_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log: &HashNumber{
					Hash:   hash,
					Number: number,
					Error:  err.Error(),
				},
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}

}

func (bp logInternalBP) ForkedResetTxPool(ctx context.Context, newHeader *types.Header, injectBlock types.Blocks, elapse time.Duration, cbft *Cbft) {
	if cbft.viewChange == nil {
		return
	}
	type HashNumber struct {
		Hash   common.Hash `json:"block_hash"`
		Number uint64      `json:"block_number"`
	}
	processor := localAddress(cbft)
	var hash []common.Hash
	for _, b := range injectBlock {
		hash = append(hash, b.Hash())
	}
	span := &Span{
		Context: Context{
			TraceID:   cbft.viewChange.Timestamp,
			SpanID:    injectBlock.String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   processor,
			Processor: processor,
		},
		StartTime:     time.Now(),
		DurationTime:  elapse,
		OperationName: "tx_pool",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "forked_reset_tx_pool",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       newHeader,
			},
			{
				Timestamp: time.Now().UnixNano(),
				Log:       hash,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}

}

func (bp logInternalBP) ResetTxPool(ctx context.Context, ext *BlockExt, elapse time.Duration, cbft *Cbft) {
	type HashNumber struct {
		Hash   common.Hash `json:"block_hash"`
		Number uint64      `json:"block_number"`
	}
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   ext.timestamp,
			SpanID:    big.NewInt(int64(ext.number)).String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   processor,
			Processor: processor,
		},
		StartTime:     time.Now(),
		DurationTime:  elapse,
		OperationName: "tx_pool",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "reset_tx_pool",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log: &HashNumber{
					Hash:   ext.block.Hash(),
					Number: ext.number,
				},
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}

}

func (bp logInternalBP) NewConfirmedBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	log.Debug("NewConfirmedBlock", "block", ext.String(), "cbft", cbft.String())

}

func (bp logInternalBP) NewLogicalBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	log.Debug("NewLogicalBlock", "block", ext.String(), "cbft", cbft.String())

}

func (bp logInternalBP) NewRootBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	log.Debug("NewRootBlock", "block", ext.String(), "cbft", cbft.String())
}

func (bp logInternalBP) NewHighestConfirmedBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	type HashNumber struct {
		Hash   common.Hash `json:"block_hash"`
		Number uint64      `json:"block_number"`
	}
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   ext.timestamp,
			SpanID:    big.NewInt(int64(ext.number)).String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   processor,
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "chain_state",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "new_highest_confirmed_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log: &HashNumber{
					Hash:   ext.block.Hash(),
					Number: ext.number,
				},
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logInternalBP) NewHighestLogicalBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	type HashNumber struct {
		Hash   common.Hash `json:"block_hash"`
		Number uint64      `json:"block_number"`
	}
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   ext.timestamp,
			SpanID:    big.NewInt(int64(ext.number)).String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   processor,
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "chain_state",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "new_highest_logical_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log: &HashNumber{
					Hash:   ext.block.Hash(),
					Number: ext.number,
				},
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logInternalBP) NewHighestRootBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	type HashNumber struct {
		Hash   common.Hash `json:"block_hash"`
		Number uint64      `json:"block_number"`
	}
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   ext.timestamp,
			SpanID:    big.NewInt(int64(ext.number)).String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   processor,
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "chain_state",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "new_highest_root_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log: &HashNumber{
					Hash:   ext.block.Hash(),
					Number: ext.number,
				},
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}

func (bp logInternalBP) SwitchView(ctx context.Context, view *viewChange, cbft *Cbft) {
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   view.Timestamp,
			SpanID:    big.NewInt(int64(view.BaseBlockNum)).String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   processor,
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "view_state",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "switch_view",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       view,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}

}

func (bp logInternalBP) Seal(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	processor := localAddress(cbft)

	span := &Span{
		Context: Context{
			TraceID:   ext.timestamp,
			SpanID:    ext.block.Number().String(),
			ParentID:  cbft.config.NodeID.String(),
			Creator:   processor,
			Processor: processor,
		},
		StartTime:     time.Now(),
		OperationName: "seal_block",
		Tags: []Tag{
			{
				Key:   "action",
				Value: "new_highest_root_block",
			},
		},
		LogRecords: []LogRecord{
			{
				Timestamp: time.Now().UnixNano(),
				Log:       ext,
			},
		},
	}
	msg, err := json.Marshal(span)
	if err == nil {
		log.Info(string(msg))
	}
}
