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

	startTime := time.Now()
	// Iterate over and process the individual transactions
	if len(block.Transactions()) > 0 {
		ctx := NewParallelContext(statedb, header, block.Hash(), gp, startTime, false)
		ctx.SetBlockGasUsedHolder(usedGas)
		ctx.SetTxList(block.Transactions())
		if err := GetExecutor().ExecuteBlocks(ctx); err != nil {
			return nil, nil, 0, err
		}
		receipts = ctx.GetReceipts()
		allLogs = ctx.GetLogs()
	}

	if bcr != nil {
		// EndBlocker()
		if err := bcr.EndBlocker(block.Header(), statedb); nil != err {
			log.Error("Failed to call EndBlocker on StateProcessor", "blockNumber", block.Number(),
				"blockHash", block.Hash(), "err", err)
			return nil, nil, 0, err
		}
	}

	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), receipts)

	return receipts, allLogs, *usedGas, nil
}
