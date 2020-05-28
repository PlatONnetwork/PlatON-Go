package core

import (
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/panjf2000/ants/v2"
)

var (
	executorOnce sync.Once
	executor     Executor
)

type Executor struct {
	chainContext ChainContext
	chainConfig  *params.ChainConfig
	vmCfg        vm.Config
	signer       types.Signer

	workerPool *ants.PoolWithFunc
}

type TaskArgs struct {
	ctx          *ParallelContext
	idx          int
	intrinsicGas uint64
}

var intrinsicGasCache atomic.Value

func EstimateTransferIntrinsicGas(txData []byte) (uint64, error) {
	intrinsicGasInf := intrinsicGasCache.Load()
	if intrinsicGasInf != nil {
		return intrinsicGasInf.(uint64), nil
	} else {
		if gas, err := IntrinsicGas(txData, false, nil); err != nil {
			return uint64(0), err
		} else {
			intrinsicGasCache.Store(gas)
			return gas, nil
		}
	}
}

func NewExecutor(chainConfig *params.ChainConfig, chainContext ChainContext, vmCfg vm.Config) {
	executorOnce.Do(func() {
		log.Info("Init parallel executor ...")
		executor = Executor{}
		executor.workerPool, _ = ants.NewPoolWithFunc(runtime.NumCPU(), func(i interface{}) {
			args := i.(TaskArgs)
			ctx := args.ctx
			idx := args.idx
			intrinsicGas := args.intrinsicGas
			executor.executeParallelTx(ctx, idx, intrinsicGas)
			ctx.wg.Done()
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

func (exe *Executor) Signer() types.Signer {
	return exe.signer
}

func (exe *Executor) ExecuteTransactions(ctx *ParallelContext) error {
	log.Debug("Execute transactions begin", "number", ctx.header.Number.Uint64(), "packNewBlock", ctx.packNewBlock, "gasPool", ctx.gp.Gas())
	log.Debug("Execute transactions goroutine info(start)", "number", ctx.header.Number, "cap", exe.workerPool.Cap(), "free", exe.workerPool.Free(), "running", exe.workerPool.Running())
	if len(ctx.txList) > 0 {
		txDag := NewTxDag(exe.signer)

		start := time.Now()

		if err := txDag.MakeDagGraph(ctx.header.Number.Uint64(), ctx.GetState(), ctx.txList, start); err != nil {
			return err
		}
		log.Debug("Make dag graph cost", "number", ctx.header.Number.Uint64(), "time", time.Since(start))
		start = time.Now()

		batchNo := 0
		for !ctx.IsTimeout() && txDag.HasNext() {
			parallelTxIdxs := txDag.Next()

			if len(parallelTxIdxs) > 0 {
				if len(parallelTxIdxs) == 1 && txDag.IsContract(parallelTxIdxs[0]) {
					exe.executeContractTransaction(ctx, parallelTxIdxs[0])
					//log.Trace(fmt.Sprintf("ExecuteBlocks(tx type:contract) done, blockNumber=%d, batchNo=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s", ctx.header.Number.Uint64(), batchNo, originIdx, tx.FromAddr().Hex(), toAddr, tx.Hash().Hex()))
				} else {
					for _, originIdx := range parallelTxIdxs {
						tx := ctx.GetTx(originIdx)
						if ctx.packNewBlock {
							if ctx.IsTimeout() {
								log.Debug("Ctx is timeout")
								break
							}

							from := tx.FromAddr(exe.signer)
							if _, popped := ctx.poppedAddresses[from]; popped {
								log.Debug("Address popped", "from", from.Hex())
								continue
							}
						}

						intrinsicGas, err := EstimateTransferIntrinsicGas(tx.Data())
						if err != nil {
							ctx.buildTransferFailedResult(originIdx, err, false)
							continue
						}
						tx.SetIntrinsicGas(intrinsicGas)

						//if err := ctx.gp.SubGas(tx.Gas()); err != nil {
						if err := ctx.gp.SubGas(intrinsicGas); err != nil {
							ctx.buildTransferFailedResult(originIdx, err, false)
							continue
						}

						ctx.wg.Add(1)
						args := TaskArgs{ctx, originIdx, intrinsicGas}
						_ = exe.workerPool.Invoke(args)
					}
					// waiting for current batch done
					ctx.wg.Wait()
					ctx.batchMerge(batchNo, parallelTxIdxs, true)
				}
			} else {
				//log.Error("DAG has unconsumed vertexes.")
			}
			batchNo++
		}
		log.Debug("Execute transactions cost", "number", ctx.header.Number, "time", time.Since(start))

		//add balance for miner
		if ctx.GetEarnings().Cmp(big.NewInt(0)) > 0 {
			ctx.state.AddMinerEarnings(ctx.header.Coinbase, ctx.GetEarnings())
		}
		start = time.Now()
		ctx.state.Finalise(true)
		log.Debug("Finalise stateDB cost", "number", ctx.header.Number, "time", time.Since(start))
	}
	log.Debug("Execute transactions goroutine info(end)", "number", ctx.header.Number, "cap", exe.workerPool.Cap(), "free", exe.workerPool.Free(), "running", exe.workerPool.Running())

	return nil
}

func (exe *Executor) executeParallelTx(ctx *ParallelContext, idx int, intrinsicGas uint64) {
	if ctx.IsTimeout() {
		return
	}
	tx := ctx.GetTx(idx)
	//log.Debug("execute tx in parallel", "txHash", tx.Hash(), "txIdx", idx, "gasPool", ctx.gp.Gas(), "txGasLimit", tx.Gas())

	msg, err := tx.AsMessage(exe.signer)
	if err != nil {
		//gas pool is subbed
		ctx.buildTransferFailedResult(idx, err, true)
		return
	}
	start := time.Now()
	fromObj := ctx.GetState().GetOrNewParallelStateObject(msg.From())
	if start.Add(30 * time.Millisecond).Before(time.Now()) {
		log.Debug("Get state object overtime", "address", msg.From().String(), "duration", time.Since(start))
	}

	mgval := new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas()), tx.GasPrice())
	if fromObj.GetBalance().Cmp(mgval) < 0 {
		ctx.buildTransferFailedResult(idx, errInsufficientBalanceForGas, true)
		return
	}

	if fromObj.GetNonce() < msg.Nonce() {
		ctx.buildTransferFailedResult(idx, ErrNonceTooHigh, true)
		return
	} else if fromObj.GetNonce() > msg.Nonce() {
		ctx.buildTransferFailedResult(idx, ErrNonceTooLow, true)
		return
	}

	if msg.Gas() < intrinsicGas {
		ctx.buildTransferFailedResult(idx, vm.ErrOutOfGas, true)
		return
	}

	minerEarnings := new(big.Int).Mul(new(big.Int).SetUint64(intrinsicGas), msg.GasPrice())
	subTotal := new(big.Int).Add(msg.Value(), minerEarnings)
	if fromObj.GetBalance().Cmp(subTotal) < 0 {
		ctx.buildTransferFailedResult(idx, errInsufficientBalanceForGas, true)
		return
	}

	fromObj.SubBalance(subTotal)
	fromObj.SetNonce(fromObj.GetNonce() + 1)

	toObj := ctx.GetState().GetOrNewParallelStateObject(*msg.To())
	toObj.AddBalance(msg.Value())

	ctx.buildTransferSuccessResult(idx, fromObj, toObj, intrinsicGas, minerEarnings)
	return
}

func (exe *Executor) executeContractTransaction(ctx *ParallelContext, idx int) {
	if ctx.IsTimeout() {
		return
	}
	snap := ctx.GetState().Snapshot()
	tx := ctx.GetTx(idx)

	//log.Debug("execute contract", "txHash", tx.Hash(), "txIdx", idx, "gasPool", ctx.gp.Gas(), "txGasLimit", tx.Gas())
	ctx.GetState().Prepare(tx.Hash(), ctx.GetBlockHash(), int(ctx.GetState().TxIdx()))
	receipt, _, err := ApplyTransaction(exe.chainConfig, exe.chainContext, ctx.GetGasPool(), ctx.GetState(), ctx.GetHeader(), tx, ctx.GetBlockGasUsedHolder(), exe.vmCfg)
	if err != nil {
		log.Debug("Execute contract transaction failed", "blockNumber", ctx.GetHeader().Number.Uint64(), "txHash", tx.Hash(), "gasPool", ctx.GetGasPool().Gas(), "txGasLimit", tx.Gas(), "err", err.Error())
		ctx.GetState().RevertToSnapshot(snap)
		return
	}
	ctx.AddPackedTx(tx)
	ctx.GetState().IncreaseTxIdx()
	ctx.AddReceipt(receipt)
	log.Debug("Execute contract transaction success", "blockNumber", ctx.GetHeader().Number.Uint64(), "txHash", tx.Hash().Hex(), "gasPool", ctx.gp.Gas(), "txGasLimit", tx.Gas(), "gasUsed", receipt.GasUsed)
}
