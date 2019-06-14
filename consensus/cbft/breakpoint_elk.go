package cbft

import (
	"context"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"time"
)

const (
	NONE = "NONE"
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

func (bp elkPrepareBP) CommitBlock(ctx context.Context, block *types.Block, txs int, gasUsed uint64, elapse time.Duration, cbft *Cbft) {
}

func (bp elkPrepareBP) SendBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
}

func (bp elkPrepareBP) ReceiveBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Info("Reporting-ReceiveBlock", "from", peerId,
		"mark", "receiveBlock",
		"msgHash", block.MsgHash().TerminalString(),
		"hash", block.Block.Hash().TerminalString(),
		"number", block.Block.Number())
}

func (bp elkPrepareBP) ReceiveVote(ctx context.Context, vote *prepareVote, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-ReceiveVote", "from", peerId,
		"mark", "receiveVote",
		"msgHash", vote.MsgHash(),
		"hash", vote.Hash.TerminalString(),
		"number", vote.Number)
}

func (bp elkPrepareBP) AcceptBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-AcceptBlock", "from", peerId,
		"mark", "acceptBlock",
		"msgHash", block.MsgHash(),
		"hash", block.Block.Hash().TerminalString(),
		"number", block.Block.Number())
	//log.Debug("AcceptBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) CacheBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
	//log.Debug("CacheBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) DiscardBlock(ctx context.Context, block *prepareBlock, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-DiscardBlock", "from", peerId,
		"mark", "discardBlock",
		"msgHash", block.MsgHash(),
		"hash", block.Block.Hash().TerminalString(),
		"number", block.Block.Number())
	//log.Debug("DiscardBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) AcceptVote(ctx context.Context, vote *prepareVote, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-AcceptVote", "from", peerId,
		"mark", "acceptVote",
		"msgHash", vote.MsgHash(),
		"hash", vote.Hash.TerminalString(),
		"number", vote.Number)
	//log.Debug("AcceptVote", "block", vote.String(), "state", state.String())
}

func (bp elkPrepareBP) CacheVote(ctx context.Context, vote *prepareVote, cbft *Cbft) {
	//log.Debug("CacheVote", "block", vote.String(), "state", state.String())
}

func (bp elkPrepareBP) DiscardVote(ctx context.Context, vote *prepareVote, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-DiscardVote", "from", peerId,
		"mark", "discardVote",
		"msgHash", vote.MsgHash(),
		"hash", vote.Hash.TerminalString(),
		"number", vote.Number)
	//log.Debug("DiscardVote", "block", vote.String(), "state", state.String())
}

func (bp elkPrepareBP) SendPrepareVote(ctx context.Context, ext *prepareVote, cbft *Cbft) {
	//log.Debug("SendPrepareVote", "block", ext.String(), "state", state.String())
}

func (bp elkPrepareBP) InvalidBlock(ctx context.Context, block *prepareBlock, err error, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-InvalidBlock", "from", peerId,
		"mark", "invalidBlock",
		"msgHash", block.MsgHash(),
		"hash", block.Block.Hash().TerminalString(),
		"number", block.Block.Number())
	//log.Debug("InvalidBlock", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) InvalidVote(ctx context.Context, vote *prepareVote, err error, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-InvalidVote", "from", peerId,
		"mark", "invalidVote",
		"msgHash", vote.MsgHash(),
		"hash", vote.Hash.TerminalString(),
		"number", vote.Number)
	//log.Debug("InvalidVote", "block", vote.String(), "state", state.String())
}

func (bp elkPrepareBP) InvalidViewChangeVote(ctx context.Context, block *prepareBlock, err error, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-InvalidViewChangeVote", "from", peerId,
		"mark", "invalidViewChangeVote",
		"msgHash", block.MsgHash(),
		"hash", block.Block.Hash().TerminalString(),
		"number", block.Block.Number())
	//log.Debug("InvalidViewChangeVote", "block", block.String(), "state", state.String())
}

func (bp elkPrepareBP) TwoThirdVotes(ctx context.Context, ext *prepareVote, cbft *Cbft) {
	//log.Debug("TwoThirdVotes", "block", ext.String(), "state", state.String())
}

type elkViewChangeBP struct {
}

func (bp elkViewChangeBP) SendViewChange(ctx context.Context, view *viewChange, cbft *Cbft) {

}

func (bp elkViewChangeBP) ReceiveViewChange(ctx context.Context, view *viewChange, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-ReceiveViewChange", "from", peerId,
		"mark", "receiveViewChange",
		"msgHash", view.MsgHash(),
		"hash", view.BaseBlockHash.TerminalString(),
		"number", view.BaseBlockNum)
	//log.Debug("ReceiveViewChange", "block", view.String(), "state", state.String())
}

func (bp elkViewChangeBP) ReceiveViewChangeVote(ctx context.Context, vote *viewChangeVote, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-ReceiveViewChangeVote", "from", peerId,
		"mark", "receiveViewChangeVote",
		"msgHash", vote.MsgHash(),
		"hash", vote.BlockHash.TerminalString(),
		"number", vote.BlockNum)
	//log.Debug("ReceiveViewChangeVote", "vote", vote.String(), "state", state.String())
}

func (bp elkViewChangeBP) InvalidViewChange(ctx context.Context, view *viewChange, err error, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-InvalidViewChange", "from", peerId,
		"mark", "invalidViewChange",
		"msgHash", view.MsgHash(),
		"hash", view.BaseBlockHash.TerminalString(),
		"number", view.BaseBlockNum)
	//log.Debug("InvalidViewChange", "view", view.String(), "state", state.String())
}

func (bp elkViewChangeBP) InvalidViewChangeVote(ctx context.Context, view *viewChangeVote, err error, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-InvalidViewChangeVote", "from", peerId,
		"mark", "invalidViewChangeVote",
		"msgHash", view.MsgHash(),
		"hash", view.BlockHash.TerminalString(),
		"number", view.BlockNum)
	//log.Debug("InvalidViewChangeVote", "view", view.String(), "state", state.String())
}

func (bp elkViewChangeBP) InvalidViewChangeBlock(ctx context.Context, view *viewChange, cbft *Cbft) {
	peerId := ctx.Value("peer")
	log.Debug("Reporting-InvalidViewChangeBlock", "from", peerId,
		"mark", "invalidViewChangeBlock",
		"msgHash", view.MsgHash(),
		"hash", view.BaseBlockHash.TerminalString(),
		"number", view.BaseBlockNum)
	//log.Debug("InvalidViewChangeBlock", "view", view.String(), "state", state.String())
}

func (bp elkViewChangeBP) TwoThirdViewChangeVotes(ctx context.Context, view *viewChange, votes ViewChangeVotes, cbft *Cbft) {
	//log.Debug("TwoThirdViewChangeVotes", "state", state.String())
}

func (bp elkViewChangeBP) SendViewChangeVote(ctx context.Context, vote *viewChangeVote, cbft *Cbft) {
	//log.Debug("SendViewChangeVote", "vote", vote.String(), "state", state.String())

}

func (bp elkViewChangeBP) ViewChangeTimeout(ctx context.Context, view *viewChange, cbft *Cbft) {
	//log.Debug("ViewChangeTimeout", "state", state.String())
}

type elkSyncBlockBP struct {
}

func (bp elkSyncBlockBP) SyncBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	//log.Debug("SyncBlock", "block", ext.String(), "state", state.String())

}

func (bp elkSyncBlockBP) InvalidBlock(ctx context.Context, ext *BlockExt, err error, cbft *Cbft) {
	log.Debug("Reporting-InvalidViewChangeBlock", "from", NONE,
		"mark", "invalidBlock",
		"msgHash", NONE,
		"hash", ext.block.Hash().TerminalString(),
		"number", ext.block.Number())
	//log.Debug("InvalidBlock", "block", ext.String(), "state", state.String())
}

type elkInternalBP struct {
}

func (bp elkInternalBP) ExecuteBlock(ctx context.Context, hash common.Hash, number uint64, timestamp uint64, elapse time.Duration) {
	/*log.Debug("Reporting-ExecuteBlock", "from", NONE,
	"mark", "executeBlock",
	"msgHash", NONE,
	"hash", hash.TerminalString(),
	"number", number)*/
	//log.Debug("ExecuteBlock", "hash", hash, "number", number, "elapse", elapse.Seconds())
}

func (bp elkInternalBP) InvalidBlock(ctx context.Context, hash common.Hash, number uint64, timestamp uint64, err error) {
	/*log.Debug("Reporting-InvalidBlock", "from", NONE,
	"mark", "invalidBlock",
	"msgHash", NONE,
	"hash", hash.TerminalString(),
	"number", number)*/
	//log.Debug("InvalidBlock", "hash", hash, number, number)
}

func (bp elkInternalBP) ForkedResetTxPool(ctx context.Context, newHeader *types.Header, injectBlock types.Blocks, elapse time.Duration, cbft *Cbft) {
	//log.Debug("ForkedResetTxPool", "newHeader", fmt.Sprintf("[hash:%s, number:%d]", newHeader.Hash().TerminalString(), newHeader.Number.Uint64()), "block", injectBlock.String(), "elapse", elapse.Seconds(), "state", state.String())

}

func (bp elkInternalBP) ResetTxPool(ctx context.Context, ext *BlockExt, elapse time.Duration, cbft *Cbft) {
	//log.Debug("ResetTxPool", "block", ext.String(), "elapse", elapse.Seconds(), "state", state.String())

}

func (bp elkInternalBP) NewConfirmedBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	/*log.Debug("Reporting-NewConfirmedBlock", "from", NONE,
	"mark", "newConfirmedBlock",
	"msgHash", NONE,
	"hash", ext.block.Hash().TerminalString(),
	"number", ext.block.Number())*/
	//log.Debug("NewConfirmedBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) NewLogicalBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	/*log.Debug("Reporting-NewLogicalBlock", "from", NONE,
	"mark", "newLogicalBlock",
	"msgHash", NONE,
	"hash", ext.block.Hash().TerminalString(),
	"number", ext.block.Number())*/
	//log.Debug("NewLogicalBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) NewRootBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	/*log.Debug("Reporting-NewRootBlock", "from", NONE,
	"mark", "newRootBlock",
	"msgHash", NONE,
	"hash", ext.block.Hash().TerminalString(),
	"number", ext.block.Number())*/
	//log.Debug("NewRootBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) NewHighestConfirmedBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	/*log.Debug("Reporting-NewHighestConfirmedBlock", "from", NONE,
	"mark", "newHighestConfirmedBlock",
	"msgHash", NONE,
	"hash", ext.block.Hash().TerminalString(),
	"number", ext.block.Number())*/
	//log.Debug("NewHighestConfirmedBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) NewHighestLogicalBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	/*log.Debug("Reporting-NewHighestLogicalBlock", "from", NONE,
	"mark", "newHighestLogicalBlock",
	"msgHash", NONE,
	"hash", ext.block.Hash().TerminalString(),
	"number", ext.block.Number())*/
	//log.Debug("NewHighestLogicalBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) NewHighestRootBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	/*log.Debug("Reporting-NewHighestRootBlock", "from", NONE,
	"mark", "newHighestRootBlock",
	"msgHash", NONE,
	"hash", ext.block.Hash().TerminalString(),
	"number", ext.block.Number())*/
	//log.Debug("NewHighestRootBlock", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) SwitchView(ctx context.Context, view *viewChange, cbft *Cbft) {
	/*log.Debug("Reporting-SwitchView", "from", NONE,
	"mark", "switchView",
	"msgHash", view.MsgHash().TerminalString(),
	"hash", view.BaseBlockHash.TerminalString(),
	"number", view.BaseBlockNum)*/
	//log.Debug("SwitchView", "view", view.String())
}

func (bp elkInternalBP) Seal(ctx context.Context, ext *BlockExt, cbft *Cbft) {
	/*log.Debug("Reporting-Seal", "from", NONE,
	"mark", "seal",
	"msgHash", NONE,
	"hash", ext.block.Hash().TerminalString(),
	"number", ext.block.Number())*/
	//log.Debug("SwitchView", "block", ext.String(), "state", state.String())
}

func (bp elkInternalBP) StoreBlock(ctx context.Context, ext *BlockExt, cbft *Cbft) {

}
