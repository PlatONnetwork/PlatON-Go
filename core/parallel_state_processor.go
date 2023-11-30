package core

import (
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
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
		if err := bcr.BeginBlocker(header, statedb); nil != err {
			log.Error("Failed to call BeginBlocker on StateProcessor", "blockNumber", block.Number(),
				"blockHash", block.Hash(), "err", err)
			return nil, nil, 0, err
		}
	}

	// Iterate over and process the individual transactions
	if len(block.Transactions()) > 0 {
		start := time.Now()
		tempContractCache := make(map[common.Address]struct{})
		signer := types.MakeSigner(p.config, block.Number(), gov.Gte150VersionState(statedb))
		ctx := NewParallelContext(statedb, header, block.Hash(), gp, false, signer, tempContractCache)
		ctx.SetBlockGasUsedHolder(usedGas)
		ctx.SetTxList(block.Transactions())

		//wait tx from cal done
		if block.CalTxFromCH != nil {
			tasks := cap(block.CalTxFromCH)
			timeout := time.NewTimer(time.Millisecond * 800)
			txHaveCal := 0
			for tasks > 0 {
				select {
				case txs := <-block.CalTxFromCH:
					txHaveCal = txHaveCal + txs
					tasks--
				case <-timeout.C:
					log.Warn("Parallel cal tx from timeout", "num", block.Number(), "left_task", tasks, "total_task", cap(block.CalTxFromCH), "txcal", txHaveCal)
					tasks = 0
				}
			}
			timeout.Stop()
		}

		if err := GetExecutor().ExecuteTransactions(ctx); err != nil {
			return nil, nil, 0, err
		}
		receipts = sortReceipts(block.Transactions(), ctx.GetReceipts())
		allLogs = ctx.GetLogs()
		log.Trace("Process parallel execute transactions cost time", "blockNumber", block.Number(), "blockHash", block.Hash(), "time", time.Since(start))
	}

	if bcr != nil {
		// EndBlocker()
		if err := bcr.EndBlocker(header, statedb); nil != err {
			log.Error("Failed to call EndBlocker on StateProcessor", "blockNumber", block.Number(),
				"blockHash", block.Hash(), "err", err)
			return nil, nil, 0, err
		}
		log.Debug("Process end blocker cost time", "blockNumber", block.Number(), "blockHash", block.Hash())
	}

	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	//p.engine.Finalize(p.bc, header, statedb, block.Transactions(), receipts)
	statedb.IntermediateRoot(true)
	return receipts, allLogs, *usedGas, nil
}

func sortReceipts(txs types.Transactions, receipts types.Receipts) types.Receipts {
	receiptsMap := make(map[common.Hash]*types.Receipt)
	cumulativeGasUsed := uint64(0)
	sortReceipts := make([]*types.Receipt, 0, receipts.Len())

	for _, r := range receipts {
		receiptsMap[r.TxHash] = r
	}
	for _, tx := range txs {
		if r, ok := receiptsMap[tx.Hash()]; ok {
			cumulativeGasUsed += r.GasUsed
			r.CumulativeGasUsed = cumulativeGasUsed
			sortReceipts = append(sortReceipts, r)
			log.Trace("sortReceipts tx", "hash", tx.Hash(), "to", tx.To(), "data", tx.Data())
		} else {
			log.Error("GetReceipts error,the corresponding receipt was not found", "txhash", tx.Hash())
		}
	}
	return sortReceipts
}
