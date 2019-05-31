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
	ReceiveBlock(ctx context.Context, block *prepareBlock, state *RoundState)
	ReceiveVote(ctx context.Context, block *prepareVote, state *RoundState)

	//accept block
	AcceptBlock(ctx context.Context, block *prepareBlock, state *RoundState)
	CacheBlock(ctx context.Context, block *prepareBlock, state *RoundState)
	DiscardBlock(ctx context.Context, block *prepareBlock, state *RoundState)

	//accept block
	AcceptVote(ctx context.Context, block *prepareVote, state *RoundState)
	CacheVote(ctx context.Context, block *prepareVote, state *RoundState)
	DiscardVote(ctx context.Context, block *prepareVote, state *RoundState)

	SendPrepareVote(ctx context.Context, ext *BlockExt, state *RoundState)
	InvalidBlock(ctx context.Context, block *prepareBlock, err error, state *RoundState)
	InvalidVote(ctx context.Context, block *prepareVote, err error, state *RoundState)
	InvalidViewChangeVote(ctx context.Context, block *prepareBlock, err error, state *RoundState)
	TwoThirdVotes(ctx context.Context, ext *BlockExt, state *RoundState)
}

type ViewChangeBP interface {
	ReceiveViewChange(ctx context.Context, view *viewChange, state *RoundState)
	ReceiveViewChangeVote(ctx context.Context, view *viewChangeVote, state *RoundState)
	InvalidViewChange(ctx context.Context, view *viewChange, err error, state *RoundState)
	InvalidViewChangeVote(ctx context.Context, view *viewChangeVote, err error, state *RoundState)
	InvalidViewChangeBlock(ctx context.Context, view *viewChange, state *RoundState)
	TwoThirdViewChangeVotes(ctx context.Context, state *RoundState)
	SendViewChangeVote(ctx context.Context, view *viewChangeVote, state *RoundState)
	ViewChangeTimeout(ctx context.Context, state *RoundState)
}

type SyncBlockBP interface {
	SyncBlock(ctx context.Context, ext *BlockExt, state *RoundState)
	InvalidBlock(ctx context.Context, ext *BlockExt, err error, state *RoundState)
}

type InternalBP interface {
	ExecuteBlock(ctx context.Context, hash common.Hash, number uint64, elapse time.Duration)
	InvalidBlock(ctx context.Context, hash common.Hash, number uint64, err error)
	ForkedResetTxPool(ctx context.Context, newHeader *types.Header, injectBlock types.Blocks, elapse time.Duration, state *RoundState)
	ResetTxPool(ctx context.Context, ext *BlockExt, elapse time.Duration, state *RoundState)
	NewConfirmedBlock(ctx context.Context, ext *BlockExt, state *RoundState)
	NewLogicalBlock(ctx context.Context, ext *BlockExt, state *RoundState)
	NewRootBlock(ctx context.Context, ext *BlockExt, state *RoundState)
	NewHighestConfirmedBlock(ctx context.Context, ext *BlockExt, state *RoundState)
	NewHighestLogicalBlock(ctx context.Context, ext *BlockExt, state *RoundState)
	NewHighestRootBlock(ctx context.Context, ext *BlockExt, state *RoundState)

	SwitchView(ctx context.Context, view *viewChange)
	Seal(ctx context.Context, ext *BlockExt, state *RoundState)
}

type defaultBreakpoint struct {
	prepareBP    PrepareBP
	viewChangeBP ViewChangeBP
	syncBlockBP  SyncBlockBP
	internalBP   InternalBP
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

func (bp defaultPrepareBP) ReceiveBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
}

func (bp defaultPrepareBP) ReceiveVote(ctx context.Context, block *prepareVote, state *RoundState) {

}

func (bp defaultPrepareBP) AcceptBlock(ctx context.Context, block *prepareBlock, state *RoundState) {

}

func (bp defaultPrepareBP) CacheBlock(ctx context.Context, block *prepareBlock, state *RoundState) {

}

func (bp defaultPrepareBP) DiscardBlock(ctx context.Context, block *prepareBlock, state *RoundState) {

}

func (bp defaultPrepareBP) AcceptVote(ctx context.Context, block *prepareVote, state *RoundState) {

}

func (bp defaultPrepareBP) CacheVote(ctx context.Context, block *prepareVote, state *RoundState) {

}

func (bp defaultPrepareBP) DiscardVote(ctx context.Context, block *prepareVote, state *RoundState) {

}

func (bp defaultPrepareBP) SendPrepareVote(ctx context.Context, ext *BlockExt, state *RoundState) {

}

func (bp defaultPrepareBP) InvalidBlock(ctx context.Context, block *prepareBlock, err error, state *RoundState) {

}

func (bp defaultPrepareBP) InvalidVote(ctx context.Context, block *prepareVote, err error, state *RoundState) {

}

func (bp defaultPrepareBP) InvalidViewChangeVote(ctx context.Context, block *prepareBlock, err error, state *RoundState) {

}

func (bp defaultPrepareBP) TwoThirdVotes(ctx context.Context, ext *BlockExt, state *RoundState) {

}

type defaultViewChangeBP struct {
}

func (bp defaultViewChangeBP) ReceiveViewChange(ctx context.Context, view *viewChange, state *RoundState) {

}

func (bp defaultViewChangeBP) ReceiveViewChangeVote(ctx context.Context, view *viewChangeVote, state *RoundState) {

}

func (bp defaultViewChangeBP) InvalidViewChange(ctx context.Context, view *viewChange, err error, state *RoundState) {

}

func (bp defaultViewChangeBP) InvalidViewChangeVote(ctx context.Context, view *viewChangeVote, err error, state *RoundState) {

}

func (bp defaultViewChangeBP) InvalidViewChangeBlock(ctx context.Context, view *viewChange, state *RoundState) {

}

func (bp defaultViewChangeBP) TwoThirdViewChangeVotes(ctx context.Context, state *RoundState) {

}

func (bp defaultViewChangeBP) SendViewChangeVote(ctx context.Context, view *viewChangeVote, state *RoundState) {

}

func (bp defaultViewChangeBP) ViewChangeTimeout(ctx context.Context, state *RoundState) {

}

type defaultSyncBlockBP struct {
}

func (bp defaultSyncBlockBP) SyncBlock(ctx context.Context, ext *BlockExt, state *RoundState) {

}

func (bp defaultSyncBlockBP) InvalidBlock(ctx context.Context, ext *BlockExt, err error, state *RoundState) {

}

type defaultInternalBP struct {
}

func (bp defaultInternalBP) InvalidBlock(ctx context.Context, hash common.Hash, number uint64, err error) {

}
func (bp defaultInternalBP) ExecuteBlock(ctx context.Context, hash common.Hash, number uint64, elapse time.Duration) {

}

func (bp defaultInternalBP) ForkedResetTxPool(ctx context.Context, newHeader *types.Header, injectBlock types.Blocks, elapse time.Duration, state *RoundState) {

}

func (bp defaultInternalBP) ResetTxPool(ctx context.Context, ext *BlockExt, elapse time.Duration, state *RoundState) {

}

func (bp defaultInternalBP) NewConfirmedBlock(ctx context.Context, ext *BlockExt, state *RoundState) {

}

func (bp defaultInternalBP) NewLogicalBlock(ctx context.Context, ext *BlockExt, state *RoundState) {

}

func (bp defaultInternalBP) NewRootBlock(ctx context.Context, ext *BlockExt, state *RoundState) {

}

func (bp defaultInternalBP) NewHighestConfirmedBlock(ctx context.Context, ext *BlockExt, state *RoundState) {

}

func (bp defaultInternalBP) NewHighestLogicalBlock(ctx context.Context, ext *BlockExt, state *RoundState) {

}

func (bp defaultInternalBP) NewHighestRootBlock(ctx context.Context, ext *BlockExt, state *RoundState) {

}

func (bp defaultInternalBP) SwitchView(ctx context.Context, view *viewChange) {

}

func (bp defaultInternalBP) Seal(ctx context.Context, ext *BlockExt, state *RoundState) {

}
