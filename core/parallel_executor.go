package core

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"

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

	log.Debug(fmt.Sprintf("PackBlockTxs begin blockNumber=%d, gasPool=%d", ctx.header.Number.Uint64(), ctx.gp.Gas()))

	if len(ctx.txList) > 0 {
		var bftEngine = exe.chainConfig.Cbft != nil
		txDag := NewTxDag(exe.signer)

		for idx, tx := range ctx.txList {
			toAddr := common.ZeroAddr.Hex()
			if tx.To() != nil {
				toAddr = tx.To().Hex()
			}
			log.Debug(fmt.Sprintf("PackBlockTxs(all tx info), blockNumber=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s, txGas=%d", ctx.header.Number.Uint64(), idx, tx.GetFromAddr().Hex(), toAddr, tx.Hash().Hex(), tx.Gas()))
		}

		if err := txDag.MakeDagGraph(ctx.GetState(), ctx.txList); err != nil {
			return err
		}

		for idx, tx := range ctx.txList {
			toAddr := common.ZeroAddr.Hex()
			if tx.To() != nil {
				toAddr = tx.To().Hex()
			}
			log.Debug(fmt.Sprintf("PackBlockTxs(all tx info), blockNumber=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s, txGas=%d", ctx.header.Number.Uint64(), idx, tx.GetFromAddr().Hex(), toAddr, tx.Hash().Hex(), tx.Gas()))
		}

		batchNo := 0

		for !ctx.IsTimeout() && txDag.HasNext() {
			parallelTxIdxs := txDag.Next()
			for _, idx := range parallelTxIdxs {
				tx := ctx.GetTx(idx)
				toAddr := common.ZeroAddr.Hex()
				if tx.To() != nil {
					toAddr = tx.To().Hex()
				}
				log.Debug(fmt.Sprintf("PackBlockTxs(batch tx info), blockNumber=%d, batchNo=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s", ctx.header.Number.Uint64(), batchNo, idx, tx.GetFromAddr().Hex(), toAddr, tx.Hash().Hex()))
			}
			if len(parallelTxIdxs) > 0 {
				if len(parallelTxIdxs) == 1 && txDag.IsContract(parallelTxIdxs[0]) {
					exe.executeTransaction(parallelTxIdxs[0])
				} else {
					for _, originIdx := range parallelTxIdxs {
						if bftEngine && ctx.IsTimeout() {
							log.Debug("ctx.IsTimeout() is TRUE")
							break
						}

						tx := exe.ctx.GetTx(originIdx)
						log.Debug(fmt.Sprintf("PackBlockTxs(to execute tx info), blockNumber=%d, batchNo=%d, idx=%d, txFrom=%s, txTo=%s, txHash=%s", ctx.header.Number.Uint64(), batchNo, originIdx, tx.GetFromAddr().Hex(), tx.To().Hex(), tx.Hash().Hex()))

						from := tx.GetFromAddr()
						if _, popped := ctx.poppedAddresses[from]; popped {
							log.Debug("popped!", "from", from.Hex())
							continue
						}

						if err := ctx.gp.SubGas(tx.Gas()); err != nil {
							exe.buildTransferFailedResult(originIdx, err)
							continue
						}

						exe.wg.Add(1)
						_ = exe.workerPool.Invoke(originIdx)
					}
					// waiting for current batch done
					exe.wg.Wait()
					exe.batchMerge(batchNo, parallelTxIdxs, true)

				}
			}
			batchNo++
		}
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
	}
	return nil
}

func (exe *Executor) VerifyBlockTxs(ctx *VerifyBlockContext) error {
	exe.ctx = ctx

	log.Debug(fmt.Sprintf("VerifyBlockTxs begin blockNumber=%d, gasPool=%d", ctx.header.Number.Uint64(), ctx.gp.Gas()))

	if len(ctx.txList) > 0 {
		txDag := NewTxDag(exe.signer)
		if err := txDag.MakeDagGraph(ctx.GetState(), ctx.txList); err != nil {
			return err
		}

		batchNo := 0
		for txDag.HasNext() {
			parallelTxIdxs := txDag.Next()
			//log.Debug(fmt.Sprintf("VerifyBlockTxs blockNumber=%d, batch=%d, parallTxIds=%+v", ctx.header.Number.Uint64(), batchNo, parallelTxIdxs))
			if len(parallelTxIdxs) > 0 {
				if len(parallelTxIdxs) == 1 && txDag.IsContract(parallelTxIdxs[0]) {
					exe.executeTransaction(parallelTxIdxs[0])
				} else {
					for _, originIdx := range parallelTxIdxs {
						tx := exe.ctx.GetTx(originIdx)

						if err := ctx.gp.SubGas(tx.Gas()); err != nil {
							exe.buildTransferFailedResult(originIdx, err)
							continue
						}

						exe.wg.Add(1)
						//submit task
						_ = exe.workerPool.Invoke(originIdx)
					}
					// waiting for current batch done
					exe.wg.Wait()

					exe.batchMerge(batchNo, parallelTxIdxs, true)
				}
			}
			batchNo++
		}

		if ctx.GetEarnings().Cmp(big.NewInt(0)) > 0 {
			ctx.state.AddMinerEarnings(ctx.header.Coinbase, ctx.GetEarnings())
			exe.ctx.GetHeader().GasUsed = ctx.GetBlockGasUsed()
		}

		/*for idx, tx := range ctx.GetPackedTxList() {
			log.Debug("verified tx", "blockNumber", ctx.header.Number.Uint64(), "idx", idx, "txHash", tx.Hash())
		}*/

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
					tx := exe.ctx.GetTx(idx)

					//total with all txs(not only all parallel txs)
					exe.ctx.CumulateBlockGasUsed(receipt.GasUsed)
					//log.Debug("tx packed success", "txHash", exe.ctx.GetTx(idx).Hash().Hex(), "txUsedGas", receipt.GasUsed)

					//reset receipt.CumulativeGasUsed
					receipt.CumulativeGasUsed = exe.ctx.GetBlockGasUsed()

					//receipt.Logs = originState.GetLogs(exe.ctx.GetTx(idx).Hash())
					//receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
					exe.ctx.AddReceipt(resultList[idx].receipt)

					exe.ctx.AddPackedTx(exe.ctx.GetTx(idx))

					exe.ctx.GetState().IncreaseTxIdx()

					// Cumulate the miner's earnings
					exe.ctx.AddEarnings(resultList[idx].minerEarnings)

					// refund to gasPool
					if tx.Gas() >= receipt.GasUsed {
						exe.ctx.AddGasPool(tx.Gas() - receipt.GasUsed)
					} else {
						log.Error("gas < gasUsed", "txIdx", idx, "gas", tx.Gas(), "gasUsed", receipt.GasUsed)
						panic("gas < gasUsed")
					}

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
		//exe.ctx.GetState().Finalise(true)
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

	mgval := new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas()), tx.GasPrice())
	if fromObj.GetBalance().Cmp(mgval) < 0 {
		exe.buildTransferFailedResult(idx, errInsufficientBalanceForGas)
		return
	}

	if fromObj.GetNonce() < msg.Nonce() {
		exe.buildTransferFailedResult(idx, ErrNonceTooHigh)
		return
	} else if fromObj.GetNonce() > msg.Nonce() {
		exe.buildTransferFailedResult(idx, ErrNonceTooLow)
		return
	}

	intrinsicGas, err := IntrinsicGas(msg.Data(), false, exe.ctx.GetState())
	if err != nil {
		exe.buildTransferFailedResult(idx, err)
		return
	}

	if msg.Gas() < intrinsicGas {
		exe.buildTransferFailedResult(idx, vm.ErrOutOfGas)
		return
	}

	minerEarnings := new(big.Int).Mul(new(big.Int).SetUint64(intrinsicGas), msg.GasPrice())
	subTotal := new(big.Int).Add(msg.Value(), minerEarnings)
	if fromObj.GetBalance().Cmp(subTotal) < 0 {
		exe.buildTransferFailedResult(idx, errInsufficientBalanceForGas)
		return
	}

	fromObj.SubBalance(subTotal)
	fromObj.SetNonce(fromObj.GetNonce() + 1)

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

	//log.Info("buildTransferFailedResult", "blockNumber", exe.ctx.GetHeader().Number.Uint64(), "gasPool", exe.ctx.GetGasPool().Gas(), "txIdx", idx, "txHash", exe.ctx.GetTx(idx).Hash(), "txTo", *exe.ctx.GetTx(idx).To(), "txGas", exe.ctx.GetTx(idx).Gas(), "err", err)

	//fmt.Println(fmt.Sprintf("---------- Fail. tx no=%d", idx))
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
	//log.Info("buildTransferSuccessResult", "blockNumber", exe.ctx.GetHeader().Number.Uint64(), "gasPool", exe.ctx.GetGasPool().Gas(), "txIdx", idx, "txHash", tx.Hash(), "txTo", *tx.To(), "txGas", exe.ctx.GetTx(idx).Gas(), "txUsedGas", txGasUsed)
	//fmt.Println(fmt.Sprintf("============ Success. tx no=%d", idx))
}

func (exe *Executor) executeTransaction(idx int) {
	defer exe.setPackBlockTimeout()
	if exe.ctx.IsTimeout() {
		return
	}
	snap := exe.ctx.GetState().Snapshot()
	tx := exe.ctx.GetTx(idx)
	exe.ctx.GetState().Prepare(tx.Hash(), exe.ctx.GetBlockHash(), int(exe.ctx.GetState().TxIdx()))
	receipt, _, err := ApplyTransaction(exe.chainConfig, exe.chainContext, exe.ctx.GetGasPool(), exe.ctx.GetState(), exe.ctx.GetHeader(), tx, exe.ctx.GetBlockGasUsedHolder(), exe.vmCfg)
	if err != nil {
		log.Error("execute tx failed", "blockNumber", exe.ctx.GetHeader().Number.Uint64(), "gasPool", exe.ctx.GetGasPool().Gas(), "txHash", tx.Hash(), "txGas", tx.Gas(), "err", err)
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
