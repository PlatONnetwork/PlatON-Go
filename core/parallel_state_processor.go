package core

import (
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

type ParallelStateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

func NewParallelStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *ParallelStateProcessor {
	return &ParallelStateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

func (p *ParallelStateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		allLogs  []*types.Log
		gp       = new(GasPool).AddGas(block.GasLimit())
	)

	if bcr != nil {
		// BeginBlocker()
		if err := bcr.BeginBlocker(block.Header(), statedb); nil != err {
			log.Error("Failed to call BeginBlocker on StateProcessor", "blockNumber", block.Number(),
				"blockHash", block.Hash(), "err", err)
			return nil, nil, 0, err
		}
	}
	if block.CalTxFromCH != nil {
		tasks := cap(block.CalTxFromCH)
		timeout := time.NewTimer(time.Second)
		for tasks > 0 {
			select {
			case <-block.CalTxFromCH:
				tasks--
			case <-timeout.C:
				log.Error("Parallel cal tx from time out")
				tasks = 0
			}
		}
		timeout.Stop()
	}
	// Iterate over and process the individual transactions
	if len(block.Transactions()) > 0 {
		start := time.Now()
		ctx := NewParallelContext(statedb, header, block.Hash(), gp, false, GetExecutor().Signer())
		ctx.SetBlockGasUsedHolder(usedGas)
		ctx.SetTxList(block.Transactions())
		if err := GetExecutor().ExecuteTransactions(ctx); err != nil {
			return nil, nil, 0, err
		}
		receipts = ctx.GetReceipts()
		allLogs = ctx.GetLogs()
		log.Trace("Process parallel execute transactions cost time", "blockNumber", block.Number(), "blockHash", block.Hash().Hex(), "time", time.Since(start))
	}

	if bcr != nil {
		// EndBlocker()
		if err := bcr.EndBlocker(block.Header(), statedb); nil != err {
			log.Error("Failed to call EndBlocker on StateProcessor", "blockNumber", block.Number(),
				"blockHash", block.Hash(), "err", err)
			return nil, nil, 0, err
		}
		log.Debug("Process end blocker cost time", "blockNumber", block.Number(), "blockHash", block.Hash().Hex())
	}

	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	//p.engine.Finalize(p.bc, header, statedb, block.Transactions(), receipts)
	statedb.IntermediateRoot(true)
	return receipts, allLogs, *usedGas, nil
}
