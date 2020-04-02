package core

import (
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
	ants "github.com/panjf2000/ants/v2"
)

var (
	executorOnce sync.Once
	executor     Executor
)

const TIMEOUT = int32(1)

type Executor struct {
	chainContext ChainContext
	chainConfig  *params.ChainConfig
	vmCfg        vm.Config
	signer       types.Signer

	workerPool *ants.PoolWithFunc
	wg         sync.WaitGroup
	ctx        Context
}

func NewExecutor(chainConfig *params.ChainConfig, chainContext ChainContext, vmCfg vm.Config) {
	executorOnce.Do(func() {
		log.Info("Init parallel executor ...")
		executor = Executor{}
		executor.workerPool, _ = ants.NewPoolWithFunc(runtime.NumCPU(), func(i interface{}) {
			executor.executeParallel(i)
			executor.wg.Done()
		})
		executor.chainConfig = chainConfig
		executor.chainContext = chainContext
		executor.signer = types.NewEIP155Signer(chainConfig.ChainID)
		executor.vmCfg = vmCfg
	})
}

func GetExecutor() *Executor {
	return &executor
}

func SetExecutor() *Executor {
	return &executor
}

func (exe *Executor) PackBlockTxs(ctx *PackBlockContext) (err error) {
	exe.ctx = ctx
	gasPoolEnough := true
	mergeCost := int64(0)
	finaliseCost := int64(0)
	if len(ctx.txList) > 0 {
		var bftEngine = exe.chainConfig.Cbft != nil
		txDag := NewTxDag(exe.signer)

		if err := txDag.MakeDagGraph(ctx.GetState(), ctx.txList); err != nil {
			return err
		}
		batchNo := 0

		for gasPoolEnough && !ctx.IsTimeout() && txDag.HasNext() {
			parallelTxIdxs := txDag.Next()
			//call executeTransaction if batch length == 1
			if len(parallelTxIdxs) == 1 {
				exe.executeTransaction(parallelTxIdxs[0])
			} else if len(parallelTxIdxs) > 1 {
				for _, originIdx := range parallelTxIdxs {
					from := ctx.GetTx(originIdx).GetFromAddr()
					if _, popped := ctx.poppedAddresses[from]; popped {
						break
					}

					if bftEngine && ctx.IsTimeout() {
						break
					}

					if ctx.gp.Gas() < params.TxGas {
						gasPoolEnough = false
						break
					}
					exe.wg.Add(1)
					_ = exe.workerPool.Invoke(originIdx)
				}
				// waiting for current batch done
				exe.wg.Wait()
				mergeStart := time.Now()
				exe.batchMerge(batchNo, parallelTxIdxs, true)
				mergeCost += time.Since(mergeStart).Milliseconds()

			} else {
				break
			}

			batchNo++
		}
		//add balance for miner
		if ctx.GetEarnings().Cmp(big.NewInt(0)) > 0 {
			//log.Debug("add miner balance", "minerAddr", ctx.header.Coinbase.Hex(), "amount", ctx.GetEarnings().Uint64())
			ctx.state.AddMinerEarnings(ctx.header.Coinbase, ctx.GetEarnings())
			//exe.ctx.GetHeader().GasUsed = ctx.GetBlockGasUsed()
		}
		finaliseStart := time.Now()
		ctx.state.Finalise(true)
		finaliseCost += time.Since(finaliseStart).Milliseconds()
	}
	return nil
}

func (exe *Executor) VerifyBlockTxs(ctx *VerifyBlockContext) error {
	exe.ctx = ctx
	if len(ctx.txList) > 0 {
		txDag := NewTxDag(exe.signer)
		if err := txDag.MakeDagGraph(ctx.GetState(), ctx.txList); err != nil {
			return err
		}

		batchNo := 0
		for txDag.HasNext() {
			parallelTxIdxs := txDag.Next()

			if len(parallelTxIdxs) == 1 {
				exe.executeTransaction(parallelTxIdxs[0])
			} else if len(parallelTxIdxs) > 1 {
				for _, originIdx := range parallelTxIdxs {
					exe.wg.Add(1)
					//submit task
					_ = exe.workerPool.Invoke(originIdx)
				}
				// waiting for current batch done
				exe.wg.Wait()

				exe.batchMerge(batchNo, parallelTxIdxs, true)
				batchNo++
			} else {
				break
			}
		}

		if ctx.GetEarnings().Cmp(big.NewInt(0)) > 0 {
			ctx.state.AddMinerEarnings(ctx.header.Coinbase, ctx.GetEarnings())
			exe.ctx.GetHeader().GasUsed = ctx.GetBlockGasUsed()
		}
		exe.ctx.GetState().Finalise(true)
	}
	return nil
}

func (exe *Executor) batchMerge(batchNo int, originIdxList []int, deleteEmptyObjects bool) {
	resultList := exe.ctx.GetResults()
	for _, idx := range originIdxList {
		if resultList[idx] != nil {
			if resultList[idx].err == nil {
				if resultList[idx].receipt != nil && resultList[idx].err == nil {
					originState := exe.ctx.GetState()
					originState.Merge(idx, resultList[idx].fromStateObject, resultList[idx].toStateObject, true)

					// Set the receipt logs and create a bloom for filtering
					// reset log's logIndex and txIndex
					receipt := resultList[idx].receipt

					//total with all txs(not only all parallel txs)
					exe.ctx.CumulateBlockGasUsed(receipt.GasUsed)

					//reset receipt.CumulativeGasUsed
					receipt.CumulativeGasUsed = exe.ctx.GetBlockGasUsed()

					//receipt.Logs = originState.GetLogs(exe.ctx.GetTx(idx).Hash())
					//receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
					exe.ctx.AddReceipt(resultList[idx].receipt)

					exe.ctx.AddPackedTx(exe.ctx.GetTx(idx))

					exe.ctx.GetState().IncreaseTxIdx()

					// Cumulate the miner's earnings
					exe.ctx.AddEarnings(resultList[idx].minerEarnings)
				} else {
					//log.Debug("to merge result, stateCpy/receipt is nil", "stateCpy is Nil", resultList[idx].stateCpy != nil, "receipt is Nil", resultList[idx].receipt != nil)
				}
			} else {
				switch resultList[idx].err {
				case ErrGasLimitReached, ErrNonceTooHigh, vm.ErrAbort:
					// pop error
					exe.ctx.SetPoppedAddress(exe.ctx.GetTx(idx).GetFromAddr())
				default:
					//shift
				}
			}
		}
	}
}

func (exe *Executor) executeParallel(arg interface{}) {
	defer exe.setPackBlockTimeout()
	if exe.ctx.IsTimeout() {
		return
	}

	idx := arg.(int)
	tx := exe.ctx.GetTx(idx)
	msg, err := tx.AsMessage(exe.signer)
	if err != nil {
		exe.buildTransferFailedResult(idx, err)
		return
	}
	fromObj := exe.ctx.GetState().GetOrNewParallelStateObject(msg.From())

	if fromObj.GetNonce() < msg.Nonce() {
		exe.buildTransferFailedResult(idx, ErrNonceTooHigh)
		return
	} else if fromObj.GetNonce() > msg.Nonce() {
		exe.buildTransferFailedResult(idx, ErrNonceTooLow)
		return
	}

	intrinsicGas, err := IntrinsicGas(msg.Data(), false)
	if err != nil {
		exe.buildTransferFailedResult(idx, err)
		return
	}

	minerEarnings := new(big.Int).Mul(new(big.Int).SetUint64(intrinsicGas), msg.GasPrice())
	fromObj.SubBalance(minerEarnings)
	fromObj.SetNonce(fromObj.GetNonce() + 1)
	if fromObj.GetBalance().Cmp(msg.Value()) < 0 {
		exe.buildTransferFailedResult(idx, errInsufficientBalanceForGas)
		return
	}
	fromObj.SubBalance(msg.Value())

	toObj := exe.ctx.GetState().GetOrNewParallelStateObject(*msg.To())
	toObj.AddBalance(msg.Value())

	exe.buildTransferSuccessResult(idx, fromObj, toObj, intrinsicGas, minerEarnings)
	return
}
func (exe *Executor) buildTransferFailedResult(idx int, err error) {
	result := &Result{
		err: err,
	}
	exe.ctx.SetResult(idx, result)
}
func (exe *Executor) buildTransferSuccessResult(idx int, fromStateObject, toStateObject *state.ParallelStateObject, txGasUsed uint64, minerEarnings *big.Int) {
	tx := exe.ctx.GetTx(idx)
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
	exe.ctx.SetResult(idx, result)
}

func (exe *Executor) executeTransaction(idx int) {
	defer exe.setPackBlockTimeout()
	if exe.ctx.IsTimeout() {
		return
	}
	snap := exe.ctx.GetState().Snapshot()
	tx := exe.ctx.GetTx(idx)
	exe.ctx.GetState().Prepare(tx.Hash(), exe.ctx.GetBlockHash(), int(exe.ctx.GetState().TxIdx()))
	receipt, _, err := ApplyTransaction(exe.chainConfig, exe.chainContext, exe.ctx.GetGasPool(), exe.ctx.GetState(), exe.ctx.GetHeader(), tx, exe.ctx.GetBlockGasUsedHolder(), vm.Config{})
	if err != nil {
		log.Error("Failed to commitTransaction on worker", "blockNumber", exe.ctx.GetHeader().Number.Uint64(), "err", err)
		exe.ctx.GetState().RevertToSnapshot(snap)
		return
	}
	exe.ctx.AddPackedTx(tx)
	exe.ctx.GetState().IncreaseTxIdx()
	exe.ctx.AddReceipt(receipt)
}

func (exe *Executor) setPackBlockTimeout() {
	if exe.ctx.IsTimeout() {
		return
	} else {
		ctx, ok := exe.ctx.(*PackBlockContext)
		if ok {
			if ctx.blockDeadline.Equal(time.Now()) || ctx.blockDeadline.Before(time.Now()) {
				ctx.SetTimeout(true)
			}
		}
	}
}
