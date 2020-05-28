package core

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	dag3 "github.com/PlatONnetwork/PlatON-Go/core/dag"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	//"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"time"
)

type TxDag struct {
	dag       *dag3.Dag
	signer    types.Signer
	contracts map[int]struct{}
}

func NewTxDag(signer types.Signer) *TxDag {
	txDag := &TxDag{
		signer:    signer,
		contracts: make(map[int]struct{}),
	}
	return txDag
}

func (txDag *TxDag) MakeDagGraph(blockNumber uint64, state *state.StateDB, txs []*types.Transaction, start time.Time) error {
	log.Debug("MakeDagGraph begin", "number", blockNumber, "txs length", len(txs), "duration", time.Since(start))
	txDag.dag = dag3.NewDag(len(txs))
	log.Debug("NewDag object end", "number", blockNumber, "txs length", len(txs), "duration", time.Since(start))
	//save all transfer addresses between two contracts(precompiled and user defined)
	transferAddressMap := make(map[common.Address]int, 0)
	latestPrecompiledIndex := -1
	for index, tx := range txs {
		xxstart := time.Now()
		log.Debug("Get tx fromAddr begin", "number", blockNumber, "txs length", len(txs), "index", index, "duration", time.Since(start), "duration2", time.Since(xxstart))
		if tx.FromAddr(txDag.signer) == (common.Address{}) {
			log.Error("The from of the transaction cannot be resolved", "number", blockNumber, "index", index)
			continue
		}
		log.Debug("Get tx fromAddr end", "number", blockNumber, "txs length", len(txs), "index", index, "duration", time.Since(start), "duration2", time.Since(xxstart))

		sstart := time.Now()
		log.Debug("Handle tx begin", "number", blockNumber, "txs length", len(txs), "index", index, "duration", time.Since(start), "durations", time.Since(sstart))
		//if tx.To() == nil || vm.IsPrecompiledContract(*tx.To()) || state.GetCodeSize(*tx.To()) > 0 {
		//	txDag.contracts[index] = struct{}{}
		//	if index > 0 {
		//		if index-latestPrecompiledIndex > 1 {
		//			for begin := latestPrecompiledIndex + 1; begin < index; begin++ {
		//				txDag.dag.AddEdge(begin, index)
		//			}
		//		} else if index-latestPrecompiledIndex == 1 {
		//			txDag.dag.AddEdge(latestPrecompiledIndex, index)
		//		}
		//	}
		//	latestPrecompiledIndex = index
		//	//reset transferAddressMap
		//	if len(transferAddressMap) > 0 {
		//		transferAddressMap = make(map[common.Address]int, 0)
		//	}
		//} else {
		dependFound := 0

		if dependIdx, ok := transferAddressMap[tx.FromAddr(txDag.signer)]; ok {
			txDag.dag.AddEdge(dependIdx, index)
			dependFound++
		}

		if dependIdx, ok := transferAddressMap[*tx.To()]; ok {
			txDag.dag.AddEdge(dependIdx, index)
			dependFound++
		}
		if dependFound == 0 && latestPrecompiledIndex >= 0 {
			txDag.dag.AddEdge(latestPrecompiledIndex, index)
		}

		transferAddressMap[tx.FromAddr(txDag.signer)] = index
		transferAddressMap[*tx.To()] = index
		//}
		log.Debug("Handle tx end", "number", blockNumber, "txs length", len(txs), "index", index, "duration", time.Since(start), "durations", time.Since(sstart))
	}
	return nil
}

func (txDag *TxDag) HasNext() bool {
	return txDag.dag.HasNext()
}

func (txDag *TxDag) Next() []int {
	return txDag.dag.Next()
}

func (txDag *TxDag) IsContract(idx int) bool {
	if _, ok := txDag.contracts[idx]; ok {
		return true
	}
	return false
}
