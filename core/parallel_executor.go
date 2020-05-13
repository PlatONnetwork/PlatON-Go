package core

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

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
}
type TaskArgs struct {
	ctx          *ParallelContext
	idx          int
	intrinsicGas uint64
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
			executor.executeParallel(ctx, idx, intrinsicGas)
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

func (exe *Executor) ExecuteBlocks(ctx *ParallelContext) error {
	log.Debug(fmt.Sprintf("ExecuteBlocks begin blockNumber=%d, packNewBlock=%t, gasPool=%d", ctx.header.Number.Uint64(), ctx.packNewBlock, ctx.gp.Gas()))

	if len(ctx.txList) > 0 {
		var bftEngine = exe.chainConfig.Cbft != nil
		txDag := NewTxDag(exe.signer)

		start := time.Now()

		if err := txDag.MakeDagGraph(ctx.header.Number.Uint64(), ctx.GetState(), ctx.txList); err != nil {
			return err
		}
		log.Debug("make dag graph cost", "blockNumber", ctx.header.Number, "time", time.Since(start))
		start = time.Now()

		batchNo := 0
		for !ctx.IsTimeout() && txDag.HasNext() {
			parallelTxIdxs := txDag.Next()
			//log.Debug(fmt.Sprintf("DAG info, blockNumber=%d, batchNo=%d, txIdxs=%+v", ctx.header.Number.Uint64(), batchNo, parallelTxIdxs))

			/*for _, idx := range parallelTxIdxs {
				tx := ctx.GetTx(idx)
				toAddr := common.ZeroAddr.Hex()
				if tx.To() != nil {
					toAddr = tx.To().Hex()
				}
				log.Debug(fmt.Sprintf("ExecuteBlocks(batch info), blockNumber=%d, batchNo=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s", ctx.header.Number.Uint64(), batchNo, idx, tx.GetFromAddr().Hex(), toAddr, tx.Hash().Hex()))
			}*/

			if len(parallelTxIdxs) > 0 {
				if len(parallelTxIdxs) == 1 && txDag.IsContract(parallelTxIdxs[0]) {
					//originIdx := parallelTxIdxs[0]
					//tx := ctx.GetTx(originIdx)
					/*toAddr := common.ZeroAddr.Hex()
					if tx.To() != nil {
						toAddr = tx.To().Hex()
					}*/
					//log.Debug(fmt.Sprintf("ExecuteBlocks(tx type:contract) begin, blockNumber=%d, batchNo=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s", ctx.header.Number.Uint64(), batchNo, originIdx, tx.GetFromAddr().Hex(), toAddr, tx.Hash().Hex()))
					//fmt.Printf("execute contract,  txHash=%s, txIdx=%d, gasPool=%d, txGasLimit=%d\n", tx.Hash().Hex(), originIdx, ctx.gp.Gas(), tx.Gas())
					exe.executeTransaction(ctx, parallelTxIdxs[0])
					//log.Debug(fmt.Sprintf("ExecuteBlocks(tx type:contract) done, blockNumber=%d, batchNo=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s", ctx.header.Number.Uint64(), batchNo, originIdx, tx.GetFromAddr().Hex(), toAddr, tx.Hash().Hex()))
				} else {
					for _, originIdx := range parallelTxIdxs {
						tx := ctx.GetTx(originIdx)
						/*toAddr := common.ZeroAddr.Hex()
						if tx.To() != nil {
							toAddr = tx.To().Hex()
						}
						log.Debug(fmt.Sprintf("ExecuteBlocks(tx type:transfer), blockNumber=%d, batchNo=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s", ctx.header.Number.Uint64(), batchNo, originIdx, tx.GetFromAddr().Hex(), toAddr, tx.Hash().Hex()))*/

						if ctx.packNewBlock {
							if bftEngine && ctx.IsTimeout() {
								log.Debug("ctx.IsTimeout() is TRUE")
								break
							}

							from := tx.GetFromAddr()
							if _, popped := ctx.poppedAddresses[from]; popped {
								log.Debug("address popped!", "from", from.Hex())
								continue
							}
						}

						//fmt.Printf("execute transfer,  txHash=%s, txIdx=%d, gasPool=%d, txGasLimit=%d\n", tx.Hash().Hex(), originIdx, ctx.gp.Gas(), tx.Gas())
						intrinsicGas, err := IntrinsicGas(tx.Data(), false, ctx.GetState())
						if err != nil {
							ctx.buildTransferFailedResult(originIdx, err, false)
							continue
						}
						//if err := ctx.gp.SubGas(tx.Gas()); err != nil {
						tx.SetIntrinsicGas(intrinsicGas)
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
		log.Debug("execute block cost", "blockNumber", ctx.header.Number, "time", time.Since(start))
		start = time.Now()

		//add balance for miner
		if ctx.GetEarnings().Cmp(big.NewInt(0)) > 0 {
			//log.Debug("add miner balance", "minerAddr", ctx.header.Coinbase.Hex(), "amount", ctx.GetEarnings().Uint64())
			ctx.state.AddMinerEarnings(ctx.header.Coinbase, ctx.GetEarnings())
			//exe.ctx.GetHeader().GasUsed = ctx.GetBlockGasUsed()
		}
		/*for idx, tx := range ctx.GetPackedTxList() {
			log.Debug("packed tx", "blockNumber", ctx.header.Number.Uint64(), "idx", idx, "txHash", tx.Hash())
		}*/
		ctx.state.Finalise(true)
		log.Debug("finalise block cost", "blockNumber", ctx.header.Number, "time", time.Since(start))
	}
	return nil
}

func (exe *Executor) executeParallel(ctx *ParallelContext, idx int, intrinsicGas uint64) {
	/*start := time.Now()
	defer func() {
		log.Debug("executeParallel cost", "blockNumber", ctx.header.Number, "txIdx", idx, "time", time.Since(start))
	}()*/
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
	fromObj := ctx.GetState().GetOrNewParallelStateObject(msg.From())

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

	/*intrinsicGas, err := IntrinsicGas(msg.Data(), false, ctx.GetState())
	if err != nil {
		ctx.buildTransferFailedResult(idx, err)
		return
	}*/

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

func (exe *Executor) executeTransaction(ctx *ParallelContext, idx int) {
	if ctx.IsTimeout() {
		return
	}
	snap := ctx.GetState().Snapshot()
	tx := ctx.GetTx(idx)

	//log.Debug("execute contract", "txHash", tx.Hash(), "txIdx", idx, "gasPool", ctx.gp.Gas(), "txGasLimit", tx.Gas())
	ctx.GetState().Prepare(tx.Hash(), ctx.GetBlockHash(), int(ctx.GetState().TxIdx()))
	receipt, _, err := ApplyTransaction(exe.chainConfig, exe.chainContext, ctx.GetGasPool(), ctx.GetState(), ctx.GetHeader(), tx, ctx.GetBlockGasUsedHolder(), exe.vmCfg)
	if err != nil {
		log.Debug("execute contract failed", "blockNumber", ctx.GetHeader().Number.Uint64(), "txHash", tx.Hash(), "gasPool", ctx.GetGasPool().Gas(), "txGasLimit", tx.Gas(), "err", err.Error())
		ctx.GetState().RevertToSnapshot(snap)
		return
	}
	ctx.AddPackedTx(tx)
	ctx.GetState().IncreaseTxIdx()
	ctx.AddReceipt(receipt)
	//log.Debug("execute contract ok", "blockNumber", ctx.GetHeader().Number.Uint64(), "txHash", tx.Hash(), "gasPool", ctx.GetGasPool().Gas(), "txGasLimit", tx.Gas(), "gasUsed", receipt.GasUsed)
	fmt.Printf("execute contract ok, blockNo=%d txHash=%s, gasPool=%d, txGasLimit=%d, gasUsed=%d\n", ctx.GetHeader().Number.Uint64(), tx.Hash().Hex(), ctx.gp.Gas(), tx.Gas(), receipt.GasUsed)
}
