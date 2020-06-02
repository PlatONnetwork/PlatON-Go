package core

import (
	"github.com/PlatONnetwork/PlatON-Go/log"
	"math/big"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
)

type Result struct {
	fromStateObject   *state.ParallelStateObject
	toStateObject     *state.ParallelStateObject
	receipt           *types.Receipt
	minerEarnings     *big.Int
	err               error
	needRefundGasPool bool
}

type ParallelContext struct {
	state           *state.StateDB
	header          *types.Header
	blockHash       common.Hash
	gp              *GasPool
	txList          []*types.Transaction
	packedTxList    []*types.Transaction
	resultList      []*Result
	receipts        types.Receipts
	poppedAddresses map[common.Address]struct{}

	blockGasUsedHolder *uint64
	earnings           *big.Int
	blockDeadline      time.Time
	packNewBlock       bool
	wg                 sync.WaitGroup
	signer             types.Signer
}

func NewParallelContext(state *state.StateDB, header *types.Header, blockHash common.Hash, gp *GasPool, packNewBlock bool, signer types.Signer) *ParallelContext {
	ctx := &ParallelContext{
		state:           state,
		header:          header,
		blockHash:       blockHash,
		gp:              gp,
		poppedAddresses: make(map[common.Address]struct{}),
		earnings:        big.NewInt(0),
		packNewBlock:    packNewBlock,
		signer:          signer,
	}
	return ctx
}

func (ctx *ParallelContext) GetState() *state.StateDB {
	return ctx.state
}

func (ctx *ParallelContext) GetHeader() *types.Header {
	return ctx.header
}

func (ctx *ParallelContext) GetBlockHash() common.Hash {
	return ctx.blockHash
}

func (ctx *ParallelContext) GetGasPool() *GasPool {
	return ctx.gp
}

func (ctx *ParallelContext) SetTxList(txs []*types.Transaction) {
	ctx.txList = txs
	ctx.resultList = make([]*Result, len(txs))
}

func (ctx *ParallelContext) GetTxList() []*types.Transaction {
	return ctx.txList
}
func (ctx *ParallelContext) GetTx(idx int) *types.Transaction {
	if len(ctx.txList) > idx {
		return ctx.txList[idx]
	} else {
		return nil
	}
}

func (ctx *ParallelContext) GetPackedTxList() []*types.Transaction {
	return ctx.packedTxList
}
func (ctx *ParallelContext) AddPackedTx(tx *types.Transaction) {
	ctx.packedTxList = append(ctx.packedTxList, tx)
}

func (ctx *ParallelContext) SetResult(idx int, result *Result) {
	if idx > len(ctx.resultList)-1 {
		return
	}
	ctx.resultList[idx] = result
}

func (ctx *ParallelContext) GetResults() []*Result {
	return ctx.resultList
}

func (ctx *ParallelContext) GetReceipts() types.Receipts {
	return ctx.receipts
}

func (ctx *ParallelContext) AddReceipt(receipt *types.Receipt) {
	ctx.receipts = append(ctx.receipts, receipt)
}

func (ctx *ParallelContext) SetPoppedAddress(poppedAddress common.Address) {
	ctx.poppedAddresses[poppedAddress] = struct{}{}
}

func (ctx *ParallelContext) GetLogs() []*types.Log {
	return ctx.state.Logs()
}

func (ctx *ParallelContext) CumulateBlockGasUsed(txGasUsed uint64) {
	*ctx.blockGasUsedHolder += txGasUsed
}

func (ctx *ParallelContext) GetBlockGasUsed() uint64 {
	return *ctx.blockGasUsedHolder
}

func (ctx *ParallelContext) GetBlockGasUsedHolder() *uint64 {
	return ctx.blockGasUsedHolder
}

func (ctx *ParallelContext) SetBlockGasUsedHolder(blockGasUsedHolder *uint64) {
	ctx.blockGasUsedHolder = blockGasUsedHolder
}

func (ctx *ParallelContext) GetEarnings() *big.Int {
	return ctx.earnings
}

func (ctx *ParallelContext) AddEarnings(earning *big.Int) {
	ctx.earnings = new(big.Int).Add(ctx.earnings, earning)
}

func (ctx *ParallelContext) SetBlockDeadline(blockDeadline time.Time) {
	ctx.blockDeadline = blockDeadline
}

func (ctx *ParallelContext) IsTimeout() bool {
	if ctx.packNewBlock {
		if ctx.blockDeadline.Equal(time.Now()) || ctx.blockDeadline.Before(time.Now()) {
			return true
		}
	}
	return false
}

func (ctx *ParallelContext) AddGasPool(amount uint64) {
	ctx.gp.AddGas(amount)
}

func (ctx *ParallelContext) buildTransferFailedResult(idx int, err error, needRefundGasPool bool) {
	result := &Result{
		err:               err,
		needRefundGasPool: needRefundGasPool,
	}
	ctx.SetResult(idx, result)
	tx := ctx.GetTx(idx)
	log.Debug("Execute trasnfer failed", "blockNumber", ctx.header.Number.Uint64(), "txIdx", idx, "txHash", ctx.GetTx(idx).Hash().TerminalString(),
		"gasPool", ctx.gp.Gas(), "txGasLimit", tx.Gas(), "txFrom", tx.FromAddr(ctx.signer).String(), "txTo", tx.To().String(),
		"txValue", tx.Value().Uint64(), "needRefundGasPool", needRefundGasPool, "error", err.Error())
}

func (ctx *ParallelContext) buildTransferSuccessResult(idx int, fromStateObject, toStateObject *state.ParallelStateObject, txGasUsed uint64, minerEarnings *big.Int) {
	tx := ctx.GetTx(idx)
	var root []byte
	receipt := types.NewReceipt(root, false, txGasUsed)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = txGasUsed
	// Set the receipt logs and create a bloom for filtering
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	//update root here instead of in state.Merge()
	fromStateObject.UpdateRoot()

	result := &Result{
		fromStateObject: fromStateObject,
		toStateObject:   toStateObject,
		receipt:         receipt,
		minerEarnings:   minerEarnings,
		err:             nil,
	}
	ctx.SetResult(idx, result)
	log.Debug("Execute trasnfer success", "blockNumber", ctx.header.Number.Uint64(), "txIdx", idx, "txHash", tx.Hash().TerminalString(),
		"gasPool", ctx.gp.Gas(), "txGasLimit", tx.Gas(), "txUsedGas", txGasUsed, "txFrom", tx.FromAddr(ctx.signer).String(), "txTo", tx.To().String(),
		"txValue", tx.Value().Uint64(), "minerEarnings", minerEarnings.Uint64())
}

func (ctx *ParallelContext) batchMerge(batchNo int, originIdxList []int, deleteEmptyObjects bool) {
	resultList := ctx.GetResults()
	for _, idx := range originIdxList {
		if resultList[idx] != nil {
			if resultList[idx].err == nil && resultList[idx].receipt != nil {
				originState := ctx.GetState()
				originState.Merge(idx, resultList[idx].fromStateObject, resultList[idx].toStateObject, true)

				// Set the receipt logs and create a bloom for filtering
				// reset log's logIndex and txIndex
				receipt := resultList[idx].receipt
				//tx := ctx.GetTx(idx)

				//total with all txs(not only all parallel txs)
				ctx.CumulateBlockGasUsed(receipt.GasUsed)
				//log.Debug("tx packed success", "txHash", exe.ctx.GetTx(idx).Hash().Hex(), "txUsedGas", receipt.GasUsed)

				//reset receipt.CumulativeGasUsed
				receipt.CumulativeGasUsed = ctx.GetBlockGasUsed()

				//receipt.Logs = originState.GetLogs(exe.ctx.GetTx(idx).Hash())
				//receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
				ctx.AddReceipt(resultList[idx].receipt)

				ctx.AddPackedTx(ctx.GetTx(idx))

				ctx.GetState().IncreaseTxIdx()

				// Cumulate the miner's earnings
				ctx.AddEarnings(resultList[idx].minerEarnings)

				// if transfer ok, needn't refund to gasPool
				/*if tx.Gas() >= receipt.GasUsed {
					ctx.AddGasPool(tx.Gas() - receipt.GasUsed)
				} else {
					log.Error("gas < gasUsed", "txIdx", idx, "gas", tx.Gas(), "gasUsed", receipt.GasUsed)
					panic("gas < gasUsed")
				}*/
			} else {
				if resultList[idx].needRefundGasPool {
					tx := ctx.GetTx(idx)
					ctx.AddGasPool(tx.GetIntrinsicGas())
				}
				switch resultList[idx].err {
				case ErrGasLimitReached, ErrNonceTooHigh, vm.ErrAbort:
					// pop error
					ctx.SetPoppedAddress(ctx.GetTx(idx).FromAddr(ctx.signer))
				default:
					//shift
				}
			}
		}
	}
}
