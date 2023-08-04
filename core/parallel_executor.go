package core

import (
	cmath "github.com/PlatONnetwork/PlatON-Go/common/math"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/panjf2000/ants/v2"

	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

const (
	// Number of contractAddress->bool associations to keep.
	contractCacheSize = 100000
)

var (
	executorOnce sync.Once
	executor     Executor
)

type Executor struct {
	chainContext ChainContext
	chainConfig  *params.ChainConfig
	vmCfg        vm.Config

	workerPool *ants.PoolWithFunc
	txpool     *TxPool
}

type TaskArgs struct {
	ctx          *ParallelContext
	idx          int
	intrinsicGas uint64
}

func NewExecutor(chainConfig *params.ChainConfig, chainContext ChainContext, vmCfg vm.Config, txpool *TxPool) {
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

		executor.vmCfg = vmCfg
		executor.txpool = txpool
	})
}

func GetExecutor() *Executor {
	return &executor
}

func (exe *Executor) ExecuteTransactions(ctx *ParallelContext) error {
	if len(ctx.txList) > 0 {
		txDag := NewTxDag(ctx.signer)
		start := time.Now()
		// load tx fromAddress from txpool by txHash
		if err := txDag.MakeDagGraph(ctx, exe); err != nil {
			return err
		}
		log.Trace("Make dag graph cost", "number", ctx.header.Number.Uint64(), "time", time.Since(start))

		start = time.Now()
		batchNo := 0
		for !ctx.IsTimeout() && txDag.HasNext() {
			parallelTxIdxs := txDag.Next()

			if len(parallelTxIdxs) <= 0 {
				break
			}

			if len(parallelTxIdxs) == 1 && txDag.IsContract(parallelTxIdxs[0]) {
				exe.executeContractTransaction(ctx, parallelTxIdxs[0])
			} else {
				for _, originIdx := range parallelTxIdxs {
					tx := ctx.GetTx(originIdx)
					if ctx.packNewBlock {
						if ctx.IsTimeout() {
							log.Warn("Parallel executor is timeout,interrupt current tx-executing")
							break
						}

						from := tx.FromAddr(ctx.signer)
						if _, popped := ctx.poppedAddresses[from]; popped {
							log.Debug("Address popped", "from", from.Bech32())
							continue
						}
					}

					intrinsicGas, err := IntrinsicGas(tx.Data(), tx.AccessList(), false)
					if err != nil {
						ctx.buildTransferFailedResult(originIdx, err, false)
						continue
					}
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
				ctx.batchMerge(parallelTxIdxs)
				batchNo++
			}
		}
		// all transactions executed
		log.Trace("Execute transactions cost", "number", ctx.header.Number, "time", time.Since(start))
		//add balance for miner
		if ctx.GetEarnings().Cmp(big.NewInt(0)) > 0 {
			ctx.state.AddMinerEarnings(ctx.header.Coinbase, ctx.GetEarnings())
		}
		start = time.Now()
		ctx.state.Finalise(true)
		log.Trace("Finalise stateDB cost", "number", ctx.header.Number, "time", time.Since(start))
	}

	// dag print info
	/*	logVerbosity := debug.GetLogVerbosity()
		if logVerbosity == log.LvlTrace {
			inf := ctx.txListInfo()
			log.Trace("TxList Info", "blockNumber", ctx.header.Number, "txList", inf)
		}*/

	return nil
}

func (exe *Executor) preCheck(msg types.Message, fromObj *state.ParallelStateObject, baseFee *big.Int, gte150 bool) error {
	// check nonce
	if fromObj.GetNonce() < msg.Nonce() {
		return ErrNonceTooHigh
	} else if fromObj.GetNonce() > msg.Nonce() {
		return ErrNonceTooLow
	}
	// check balance
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), msg.GasPrice())
	balanceCheck := mgval
	if msg.GasFeeCap() != nil {
		balanceCheck = new(big.Int).SetUint64(msg.Gas())
		balanceCheck = balanceCheck.Mul(balanceCheck, msg.GasFeeCap())
		balanceCheck.Add(balanceCheck, msg.Value())
	}
	if fromObj.GetBalance().Cmp(balanceCheck) < 0 {
		return errInsufficientBalanceForGas
	}
	// Make sure that transaction gasFeeCap is greater than the baseFee (post london)
	if gte150 {
		// Skip the checks if gas fields are zero and baseFee was explicitly disabled (eth_call)
		if !exe.vmCfg.NoBaseFee || msg.GasFeeCap().BitLen() > 0 || msg.GasTipCap().BitLen() > 0 {
			if msg.GasFeeCap().BitLen() > 256 {
				return ErrFeeCapVeryHigh
			}
			if msg.GasTipCap().BitLen() > 256 {
				return ErrTipVeryHigh
			}
			if msg.GasFeeCap().Cmp(msg.GasTipCap()) < 0 {
				return ErrTipAboveFeeCap
			}
			// This will panic if baseFee is nil, but basefee presence is verified
			// as part of header validation.
			if msg.GasFeeCap().Cmp(baseFee) < 0 {
				return ErrFeeCapTooLow
			}
		}
	}
	return nil
}

func (exe *Executor) executeParallelTx(ctx *ParallelContext, idx int, intrinsicGas uint64) {
	if ctx.IsTimeout() {
		return
	}
	tx := ctx.GetTx(idx)

	msg, err := tx.AsMessage(ctx.signer, ctx.header.BaseFee)
	if err != nil {
		//gas pool is subbed
		ctx.buildTransferFailedResult(idx, err, true)
		return
	}

	if msg.Gas() < intrinsicGas {
		ctx.buildTransferFailedResult(idx, vm.ErrOutOfGas, true)
		return
	}

	fromObj := ctx.GetState().GetOrNewParallelStateObject(msg.From())
	// preCheck
	pauli := gov.Gte150VersionState(ctx.state)
	if err := exe.preCheck(msg, fromObj, ctx.header.BaseFee, pauli); err != nil {
		ctx.buildTransferFailedResult(idx, err, true)
		return
	}

	// miner tip
	effectiveTip := msg.GasPrice()
	if pauli {
		effectiveTip = cmath.BigMin(msg.GasTipCap(), new(big.Int).Sub(msg.GasFeeCap(), ctx.header.BaseFee))
	}
	minerEarnings := new(big.Int).Mul(new(big.Int).SetUint64(intrinsicGas), effectiveTip)
	// sender fee
	fee := new(big.Int).Mul(new(big.Int).SetUint64(intrinsicGas), msg.GasPrice())
	cost := new(big.Int).Add(msg.Value(), fee)
	if fromObj.GetBalance().Cmp(cost) < 0 {
		ctx.buildTransferFailedResult(idx, errInsufficientBalanceForGas, true)
		return
	}

	fromObj.SubBalance(cost)
	fromObj.SetNonce(fromObj.GetNonce() + 1)

	var toObj *state.ParallelStateObject
	if msg.From() == *msg.To() {
		toObj = fromObj
	} else {
		toObj = ctx.GetState().GetOrNewParallelStateObject(*msg.To())
	}
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
	ctx.GetState().Prepare(tx.Hash(), int(ctx.GetState().TxIdx()))
	receipt, err := ApplyTransaction(exe.chainConfig, exe.chainContext, ctx.GetGasPool(), ctx.GetState(), ctx.GetHeader(), tx, ctx.GetBlockGasUsedHolder(), exe.vmCfg)
	if err != nil {
		log.Warn("Execute contract transaction failed", "blockNumber", ctx.GetHeader().Number.Uint64(), "txHash", tx.Hash(), "gasPool", ctx.GetGasPool().Gas(), "txGasLimit", tx.Gas(), "err", err.Error())
		ctx.GetState().RevertToSnapshot(snap)
		return
	}
	ctx.AddPackedTx(tx)
	ctx.GetState().IncreaseTxIdx()
	ctx.AddReceipt(receipt)
	log.Debug("Execute contract transaction success", "blockNumber", ctx.GetHeader().Number.Uint64(), "txHash", tx.Hash().Hex(), "gasPool", ctx.gp.Gas(), "txGasLimit", tx.Gas(), "gasUsed", receipt.GasUsed)
}

func (exe *Executor) isContract(tx *types.Transaction, state *state.StateDB, ctx *ParallelContext) bool {
	address := tx.To()
	if address == nil { // create contract
		contractAddress := crypto.CreateAddress(tx.FromAddr(ctx.signer), tx.Nonce())
		ctx.tempContractCache[contractAddress] = struct{}{}
		return true
	}
	if _, ok := ctx.tempContractCache[*address]; ok {
		return true
	}
	isContract := vm.IsPrecompiledContract(*address, exe.chainConfig.Rules(ctx.header.Number), gov.Gte150VersionState(state)) || state.GetCodeSize(*address) > 0
	return isContract
}
