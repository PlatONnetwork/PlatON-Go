package cbft

import (
	"context"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"time"
)

//Breakpoint is help resolve cbft consensus state
//Self implement breakpoint, watch all state about cbft
type Breakpoint interface {
	PrepareBP() PrepareBP
	ViewChangeBP() ViewChangeBP
	InternalBP() InternalBP
	SyncBlockBP() SyncBlockBP
}

type PrepareBP interface {
	ReceiveBlock(ctx context.Context, block *prepareBlock, cbft *Cbft)
	ReceiveVote(ctx context.Context, block *prepareVote, cbft *Cbft)

	//accept block
	AcceptBlock(ctx context.Context, block *prepareBlock, cbft *Cbft)
	CacheBlock(ctx context.Context, block *prepareBlock, cbft *Cbft)
	DiscardBlock(ctx context.Context, block *prepareBlock, cbft *Cbft)

	//accept block
	AcceptVote(ctx context.Context, block *prepareVote, cbft *Cbft)
	CacheVote(ctx context.Context, block *prepareVote, cbft *Cbft)
	DiscardVote(ctx context.Context, block *prepareVote, cbft *Cbft)

	SendPrepareVote(ctx context.Context, ext *prepareVote, cbft *Cbft)
	InvalidBlock(ctx context.Context, block *prepareBlock, err error, cbft *Cbft)
	InvalidVote(ctx context.Context, block *prepareVote, err error, cbft *Cbft)
	InvalidViewChangeVote(ctx context.Context, block *prepareBlock, err error, cbft *Cbft)
	TwoThirdVotes(ctx context.Context, ext *prepareVote, cbft *Cbft)
}

type ViewChangeBP interface {
	SendViewChange(ctx context.Context, view *viewChange, cbft *Cbft)
	ReceiveViewChange(ctx context.Context, view *viewChange, cbft *Cbft)
	ReceiveViewChangeVote(ctx context.Context, view *viewChangeVote, cbft *Cbft)
	InvalidViewChange(ctx context.Context, view *viewChange, err error, cbft *Cbft)
	InvalidViewChangeVote(ctx context.Context, view *viewChangeVote, err error, cbft *Cbft)
	InvalidViewChangeBlock(ctx context.Context, view *viewChange, cbft *Cbft)
	TwoThirdViewChangeVotes(ctx context.Context, view *viewChange, votes ViewChangeVotes,  cbft *Cbft)
	SendViewChangeVote(ctx context.Context, view *viewChangeVote, cbft *Cbft)
	ViewChangeTimeout(ctx context.Context, view *viewChange, cbft *Cbft)
}

type SyncBlockBP interface {
	SyncBlock(ctx context.Context, ext *BlockExt, cbft *Cbft)
	InvalidBlock(ctx context.Context, ext *BlockExt, err error, cbft *Cbft)
}

type InternalBP interface {
	ExecuteBlock(ctx context.Context, hash common.Hash, number uint64, timestamp uint64, elapse time.Duration)
	InvalidBlock(ctx context.Context, hash common.Hash, number uint64, timestamp uint64, err error)
	ForkedResetTxPool(ctx context.Context, newHeader *types.Header, injectBlock types.Blocks, elapse time.Duration, cbft *Cbft)
	ResetTxPool(ctx context.Context, ext *BlockExt, elapse time.Duration, cbft *Cbft)
	NewConfirmedBlock(ctx context.Context, ext *BlockExt, cbft *Cbft)
	NewLogicalBlock(ctx context.Context, ext *BlockExt, cbft *Cbft)
	NewRootBlock(ctx context.Context, ext *BlockExt, cbft *Cbft)
	NewHighestConfirmedBlock(ctx context.Context, ext *BlockExt, cbft *Cbft)
	NewHighestLogicalBlock(ctx context.Context, ext *BlockExt, cbft *Cbft)
	NewHighestRootBlock(ctx context.Context, ext *BlockExt, cbft *Cbft)

	SwitchView(ctx context.Context, view *viewChange, cbft *Cbft)
	Seal(ctx context.Context, ext *BlockExt, cbft *Cbft)
}

type defaultBreakpoint struct {
	prepareBP    PrepareBP
	viewChangeBP ViewChangeBP
	syncBlockBP  SyncBlockBP
	internalBP   InternalBP
}

func getBreakpoint(t string) Breakpoint {
	switch t {
	case "tracing":
		return logBP
	case "default":
		return defaultBP
	case "elk":
		return elkBP
	}
	return defaultBP
}

var defaultBP Breakpoint

func init() {
	defaultBP = &defaultBreakpoint{
		prepareBP:    new(defaultPrepareBP),
		viewChangeBP: new(defaultViewChangeBP),
		syncBlockBP:  new(defaultSyncBlockBP),
		internalBP:   new(defaultInternalBP),
	}
}

func (bp defaultBreakpoint) PrepareBP() PrepareBP {
	return bp.prepareBP
}

func (bp defaultBreakpoint) ViewChangeBP() ViewChangeBP {
	return bp.viewChangeBP
}

func (bp defaultBreakpoint) InternalBP() InternalBP {
	return bp.internalBP
}

func (bp defaultBreakpoint) SyncBlockBP() SyncBlockBP {
	return bp.syncBlockBP
}

type defaultPrepareBP struct {
}

func (bp defaultPrepareBP) ReceiveBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
}

func (bp defaultPrepareBP) ReceiveVote(ctx context.Context, block *prepareVote, cbft *Cbft) {

}

func (bp defaultPrepareBP) AcceptBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {

}

func (bp defaultPrepareBP) CacheBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {

}

func (bp defaultPrepareBP) DiscardBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {

}

func (bp defaultPrepareBP) AcceptVote(ctx context.Context, block *prepareVote, cbft *Cbft) {

}

func (bp defaultPrepareBP) CacheVote(ctx context.Context, block *prepareVote, cbft *Cbft) {

}

func (bp defaultPrepareBP) DiscardVote(ctx context.Context, block *prepareVote, cbft *Cbft) {

}

func (bp defaultPrepareBP) SendPrepareVote(ctx context.Context, ext *prepareVote, cbft *Cbft) {

}

func (bp defaultPrepareBP) InvalidBlock(ctx context.Context, block *prepareBlock, err error, cbft *Cbft) {

}

func (bp defaultPrepareBP) InvalidVote(ctx context.Context, block *prepareVote, err error, cbft *Cbft) {

}

func (bp defaultPrepareBP) InvalidViewChangeVote(ctx context.Context, block *prepareBlock, err error, cbft *Cbft) {

}

func (bp defaultPrepareBP) TwoThirdVotes(ctx context.Context, ext *prepareVote, cbft *Cbft) {

}

type defaultViewChangeBP struct {
}

func (bp defaultViewChangeBP) SendViewChange(ctx context.Context, view *viewChange, cbft *Cbft) {

}

func (bp defaultViewChangeBP) ReceiveViewChange(ctx context.Context, view *viewChange, cbft *Cbft) {

}

func (bp defaultViewChangeBP) ReceiveViewChangeVote(ctx context.Context, view *viewChangeVote, cbft *Cbft) {

}

func (bp defaultViewChangeBP) InvalidViewChange(ctx context.Context, view *viewChange, err error, cbft *Cbft) {

}

func (bp defaultViewChangeBP) InvalidViewChangeVote(ctx context.Context, view *viewChangeVote, err error, cbft *Cbft) {

}

func (bp defaultViewChangeBP) InvalidViewChangeBlock(ctx context.Context, view *viewChange, cbft *Cbft) {

}

func (bp defaultViewChangeBP) TwoThirdViewChangeVotes(ctx context.Context, view *viewChange, votes ViewChangeVotes, cbft *Cbft) {

}

func (bp defaultViewChangeBP) SendViewChangeVote(ctx context.Context, view *viewChangeVote, cbft *Cbft) {

}

func (bp defaultViewChangeBP) ViewChangeTimeout(ctx context.Context, view *viewChange, cbft *Cbft) {

}

type defaultSyncBlockBP struct {
}

func (bp defaultSyncBlockBP) SyncBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}

func (bp defaultSyncBlockBP) InvalidBlock(ctx context.Context, ext *BlockExt, err error, cbft *Cbft) {

}

type defaultInternalBP struct {
}

func (bp defaultInternalBP) InvalidBlock(ctx context.Context, hash common.Hash, number uint64, timestamp uint64, err error) {

}
func (bp defaultInternalBP) ExecuteBlock(ctx context.Context, hash common.Hash, number uint64, timestamp uint64, elapse time.Duration) {

}

func (bp defaultInternalBP) ForkedResetTxPool(ctx context.Context, newHeader *types.Header, injectBlock types.Blocks, elapse time.Duration, cbft *Cbft) {

}

func (bp defaultInternalBP) ResetTxPool(ctx context.Context, ext *BlockExt, elapse time.Duration, cbft *Cbft) {

}

func (bp defaultInternalBP) NewConfirmedBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}

func (bp defaultInternalBP) NewLogicalBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}

func (bp defaultInternalBP) NewRootBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}

func (bp defaultInternalBP) NewHighestConfirmedBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}

func (bp defaultInternalBP) NewHighestLogicalBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}

func (bp defaultInternalBP) NewHighestRootBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}

func (bp defaultInternalBP) SwitchView(ctx context.Context, view *viewChange, cbft *Cbft) {

}

func (bp defaultInternalBP) Seal(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}
