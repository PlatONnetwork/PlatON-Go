package cbft

import (
	"context"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"time"
)

var elkBP Breakpoint

func init() {
	elkBP = &defaultBreakpoint{
		prepareBP:    new(elkPrepareBP),
		viewChangeBP: new(elkViewChangeBP),
		syncBlockBP:  new(elkSyncBlockBP),
		internalBP:   new(elkInternalBP),
	}
}

type elkPrepareBP struct {
}

func (bp elkPrepareBP) ReceiveBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
	//log.Debug("ReceiveBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) ReceiveVote(ctx context.Context, vote *prepareVote, state *RoundState) {
	//log.Debug("ReceiveVote", "block", vote.String(), "state", state.String())

}

func (bp elkPrepareBP) AcceptBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
	//log.Debug("AcceptBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) CacheBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
	//log.Debug("CacheBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) DiscardBlock(ctx context.Context, block *prepareBlock, state *RoundState) {
	//log.Debug("DiscardBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) AcceptVote(ctx context.Context, vote *prepareVote, state *RoundState) {
	//log.Debug("AcceptVote", "block", vote.String(), "state", state.String())
}

func (bp elkPrepareBP) CacheVote(ctx context.Context, vote *prepareVote, state *RoundState) {
	//log.Debug("CacheVote", "block", vote.String(), "state", state.String())
}

func (bp elkPrepareBP) DiscardVote(ctx context.Context, vote *prepareVote, state *RoundState) {
	//log.Debug("DiscardVote", "block", vote.String(), "state", state.String())
}

func (bp elkPrepareBP) SendPrepareVote(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("SendPrepareVote", "block", ext.String(), "state", state.String())
}

func (bp elkPrepareBP) InvalidBlock(ctx context.Context, block *prepareBlock, err error, state *RoundState) {
	//log.Debug("InvalidBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) InvalidVote(ctx context.Context, vote *prepareVote, err error, state *RoundState) {
	//log.Debug("InvalidVote", "block", vote.String(), "state", state.String())
}

func (bp elkPrepareBP) InvalidViewChangeVote(ctx context.Context, block *prepareBlock, err error, state *RoundState) {
	//log.Debug("InvalidViewChangeVote", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) TwoThirdVotes(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("TwoThirdVotes", "block", ext.String(), "state", state.String())
}

type elkViewChangeBP struct {
}

func (bp elkViewChangeBP) ReceiveViewChange(ctx context.Context, view *viewChange, state *RoundState) {
	//log.Debug("ReceiveViewChange", "block", view.String(), "state", state.String())
}

func (bp elkViewChangeBP) ReceiveViewChangeVote(ctx context.Context, vote *viewChangeVote, state *RoundState) {
	//log.Debug("ReceiveViewChangeVote", "vote", vote.String(), "state", state.String())
}

func (bp elkViewChangeBP) InvalidViewChange(ctx context.Context, view *viewChange, err error, state *RoundState) {
	//log.Debug("InvalidViewChange", "view", view.String(), "state", state.String())
}

func (bp elkViewChangeBP) InvalidViewChangeVote(ctx context.Context, view *viewChangeVote, err error, state *RoundState) {
	//log.Debug("InvalidViewChangeVote", "view", view.String(), "state", state.String())
}

func (bp elkViewChangeBP) InvalidViewChangeBlock(ctx context.Context, view *viewChange, state *RoundState) {
	//log.Debug("InvalidViewChangeBlock", "view", view.String(), "state", state.String())
}

func (bp elkViewChangeBP) TwoThirdViewChangeVotes(ctx context.Context, state *RoundState) {
	//log.Debug("TwoThirdViewChangeVotes", "state", state.String())
}

func (bp elkViewChangeBP) SendViewChangeVote(ctx context.Context, vote *viewChangeVote, state *RoundState) {
	//log.Debug("SendViewChangeVote", "vote", vote.String(), "state", state.String())

}

func (bp elkViewChangeBP) ViewChangeTimeout(ctx context.Context, state *RoundState) {
	//log.Debug("ViewChangeTimeout", "state", state.String())

}

type elkSyncBlockBP struct {
}

func (bp elkSyncBlockBP) SyncBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("SyncBlock", "block", ext.String(), "state", state.String())

}

func (bp elkSyncBlockBP) InvalidBlock(ctx context.Context, ext *BlockExt, err error, state *RoundState) {
	//log.Debug("InvalidBlock", "block", ext.String(), "state", state.String())

}

type elkInternalBP struct {
}

func (bp elkInternalBP) ExecuteBlock(ctx context.Context, hash common.Hash, number uint64, elapse time.Duration) {
	//log.Debug("ExecuteBlock", "hash", hash, "number", number, "elapse", elapse.Seconds())
}

func (bp elkInternalBP) InvalidBlock(ctx context.Context, hash common.Hash, number uint64, err error) {
	//log.Debug("InvalidBlock", "hash", hash, number, number)

}

func (bp elkInternalBP) ForkedResetTxPool(ctx context.Context, newHeader *types.Header, injectBlock types.Blocks, elapse time.Duration, state *RoundState) {
	//log.Debug("ForkedResetTxPool", "newHeader", fmt.Sprintf("[hash:%s, number:%d]", newHeader.Hash().TerminalString(), newHeader.Number.Uint64()), "block", injectBlock.String(), "elapse", elapse.Seconds(), "state", state.String())

}

func (bp elkInternalBP) ResetTxPool(ctx context.Context, ext *BlockExt, elapse time.Duration, state *RoundState) {
	//log.Debug("ResetTxPool", "block", ext.String(), "elapse", elapse.Seconds(), "state", state.String())

}

func (bp elkInternalBP) NewConfirmedBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("NewConfirmedBlock", "block", ext.String(), "state", state.String())

}

func (bp elkInternalBP) NewLogicalBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("NewLogicalBlock", "block", ext.String(), "state", state.String())

}

func (bp elkInternalBP) NewRootBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("NewRootBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) NewHighestConfirmedBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("NewHighestConfirmedBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) NewHighestLogicalBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("NewHighestLogicalBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) NewHighestRootBlock(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("NewHighestRootBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) SwitchView(ctx context.Context, view *viewChange) {
	//log.Debug("SwitchView", "view", view.String())

}

func (bp elkInternalBP) Seal(ctx context.Context, ext *BlockExt, state *RoundState) {
	//log.Debug("SwitchView", "block", ext.String(), "state", state.String())
}
