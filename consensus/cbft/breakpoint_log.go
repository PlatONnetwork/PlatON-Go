package cbft

import (
	"context"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
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
	Timestamp uint64      `json:"timestamp"`
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

func (bp logPrepareBP) ReceiveBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
	log.Debug("ReceiveBlock", "block", block.String(), "state", state.String())
}

func (bp logPrepareBP) ReceiveVote(ctx context.Context, vote *prepareVote, state *RoundState) {
	log.Debug("ReceiveVote", "block", vote.String(), "state", state.String())

}

func (bp logPrepareBP) AcceptBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
	log.Debug("AcceptBlock", "block", block.String(), "state", state.String())
}

func (bp logPrepareBP) CacheBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
	log.Debug("CacheBlock", "block", block.String(), "state", state.String())
}

func (bp logPrepareBP) DiscardBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
	log.Debug("DiscardBlock", "block", block.String(), "state", state.String())
}

func (bp logPrepareBP) AcceptVote(ctx context.Context, vote *prepareVote, state *RoundState) {
	log.Debug("AcceptVote", "block", vote.String(), "state", state.String())
}

func (bp logPrepareBP) CacheVote(ctx context.Context, vote *prepareVote, state *RoundState) {
	log.Debug("CacheVote", "block", vote.String(), "state", state.String())
}

func (bp logPrepareBP) DiscardVote(ctx context.Context, vote *prepareVote, state *RoundState) {
	log.Debug("DiscardVote", "block", vote.String(), "state", state.String())
}

func (bp logPrepareBP) SendPrepareVote(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("SendPrepareVote", "block", ext.String(), "state", state.String())
}

func (bp logPrepareBP) InvalidBlock(ctx context.Context, block *prepareBlock, err error, state *RoundState) {
	log.Debug("InvalidBlock", "block", block.String(), "state", state.String())
}

func (bp logPrepareBP) InvalidVote(ctx context.Context, vote *prepareVote, err error, state *RoundState) {
	log.Debug("InvalidVote", "block", vote.String(), "state", state.String())
}

func (bp logPrepareBP) InvalidViewChangeVote(ctx context.Context, block *prepareBlock, err error, state *RoundState) {
	log.Debug("InvalidViewChangeVote", "block", block.String(), "state", state.String())
}

func (bp logPrepareBP) TwoThirdVotes(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("TwoThirdVotes", "block", ext.String(), "state", state.String())
}

type logViewChangeBP struct {
}

func (bp logViewChangeBP) ReceiveViewChange(ctx context.Context, view *viewChange, state *RoundState) {
	log.Debug("ReceiveViewChange", "block", view.String(), "state", state.String())
}

func (bp logViewChangeBP) ReceiveViewChangeVote(ctx context.Context, vote *viewChangeVote, state *RoundState) {
	log.Debug("ReceiveViewChangeVote", "vote", vote.String(), "state", state.String())
}

func (bp logViewChangeBP) InvalidViewChange(ctx context.Context, view *viewChange, err error, state *RoundState) {
	log.Debug("InvalidViewChange", "view", view.String(), "state", state.String())
}

func (bp logViewChangeBP) InvalidViewChangeVote(ctx context.Context, view *viewChangeVote, err error, state *RoundState) {
	log.Debug("InvalidViewChangeVote", "view", view.String(), "state", state.String())
}

func (bp logViewChangeBP) InvalidViewChangeBlock(ctx context.Context, view *viewChange, state *RoundState) {
	log.Debug("InvalidViewChangeBlock", "view", view.String(), "state", state.String())
}

func (bp logViewChangeBP) TwoThirdViewChangeVotes(ctx context.Context, state *RoundState) {
	log.Debug("TwoThirdViewChangeVotes", "state", state.String())
}

func (bp logViewChangeBP) SendViewChangeVote(ctx context.Context, vote *viewChangeVote, state *RoundState) {
	log.Debug("SendViewChangeVote", "vote", vote.String(), "state", state.String())

}

func (bp logViewChangeBP) ViewChangeTimeout(ctx context.Context, state *RoundState) {
	log.Debug("ViewChangeTimeout", "state", state.String())

}

type logSyncBlockBP struct {
}

func (bp logSyncBlockBP) SyncBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("SyncBlock", "block", ext.String(), "state", state.String())

}

func (bp logSyncBlockBP) InvalidBlock(ctx context.Context, ext *BlockExt, err error, state *RoundState) {
	log.Debug("InvalidBlock", "block", ext.String(), "state", state.String())

}

type logInternalBP struct {
}

func (bp logInternalBP) ExecuteBlock(ctx context.Context, hash common.Hash, number uint64, elapse time.Duration) {
	log.Debug("ExecuteBlock", "hash", hash, "number", number, "elapse", elapse.Seconds())
}

func (bp logInternalBP) InvalidBlock(ctx context.Context, hash common.Hash, number uint64, err error) {
	log.Debug("InvalidBlock", "hash", hash, number, number)

}

func (bp logInternalBP) ForkedResetTxPool(ctx context.Context, newHeader *types.Header, injectBlock types.Blocks, elapse time.Duration, state *RoundState) {
	log.Debug("ForkedResetTxPool",
		"newHeader", fmt.Sprintf("[hash:%s, number:%d]", newHeader.Hash().TerminalString(), newHeader.Number.Uint64()),
		"block", injectBlock.String(), "elapse", elapse.Seconds(), "state", state.String())

}

func (bp logInternalBP) ResetTxPool(ctx context.Context, ext *BlockExt, elapse time.Duration, state *RoundState) {
	log.Debug("ResetTxPool", "block", ext.String(), "elapse", elapse.Seconds(), "state", state.String())

}

func (bp logInternalBP) NewConfirmedBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("NewConfirmedBlock", "block", ext.String(), "state", state.String())

}

func (bp logInternalBP) NewLogicalBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("NewLogicalBlock", "block", ext.String(), "state", state.String())

}

func (bp logInternalBP) NewRootBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("NewRootBlock", "block", ext.String(), "state", state.String())
}

func (bp logInternalBP) NewHighestConfirmedBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("NewHighestConfirmedBlock", "block", ext.String(), "state", state.String())
}

func (bp logInternalBP) NewHighestLogicalBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("NewHighestLogicalBlock", "block", ext.String(), "state", state.String())
}

func (bp logInternalBP) NewHighestRootBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("NewHighestRootBlock", "block", ext.String(), "state", state.String())
}

func (bp logInternalBP) SwitchView(ctx context.Context, view *viewChange) {
	log.Debug("SwitchView", "view", view.String())

}

func (bp logInternalBP) Seal(ctx context.Context, ext *BlockExt, state *RoundState) {
	log.Debug("SwitchView", "block", ext.String(), "state", state.String())
}
