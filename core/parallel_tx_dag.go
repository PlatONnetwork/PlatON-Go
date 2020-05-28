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
	for curIdx, cur := range txs {
		if cur.GetFromAddr() == nil {
			xxstart := time.Now()
			log.Debug("Tx fromAddr is nil begin", "number", blockNumber, "txs length", len(txs), "index", curIdx, "duration", time.Since(start), "duration2", time.Since(xxstart))
			if from, err := types.Sender(txDag.signer, cur); err != nil {
				return err
			} else {
				cur.SetFromAddr(&from)
				log.Debug("Tx fromAddr is nil end", "number", blockNumber, "txs length", len(txs), "index", curIdx, "duration", time.Since(start), "duration2", time.Since(xxstart))
			}
		}

		sstart := time.Now()
		log.Debug("Handle tx begin", "number", blockNumber, "txs length", len(txs), "index", curIdx, "duration", time.Since(start), "durations", time.Since(sstart))
		//if cur.To() == nil || vm.IsPrecompiledContract(*cur.To()) || state.GetCodeSize(*cur.To()) > 0 {
		//	txDag.contracts[curIdx] = struct{}{}
		//	if curIdx > 0 {
		//		if curIdx-latestPrecompiledIndex > 1 {
		//			for begin := latestPrecompiledIndex + 1; begin < curIdx; begin++ {
		//				txDag.dag.AddEdge(begin, curIdx)
		//			}
		//		} else if curIdx-latestPrecompiledIndex == 1 {
		//			txDag.dag.AddEdge(latestPrecompiledIndex, curIdx)
		//		}
		//	}
		//	latestPrecompiledIndex = curIdx
		//	//reset transferAddressMap
		//	if len(transferAddressMap) > 0 {
		//		transferAddressMap = make(map[common.Address]int, 0)
		//	}
		//} else {
		dependFound := 0

		if dependIdx, ok := transferAddressMap[*cur.GetFromAddr()]; ok {
			txDag.dag.AddEdge(dependIdx, curIdx)
			dependFound++
		}

		//if cur.GetFromAddr().Hex() != cur.To().Hex() {
		if dependIdx, ok := transferAddressMap[*cur.To()]; ok {
			txDag.dag.AddEdge(dependIdx, curIdx)
			dependFound++
		}
		//}
		if dependFound == 0 && latestPrecompiledIndex >= 0 {
			txDag.dag.AddEdge(latestPrecompiledIndex, curIdx)
		}

		transferAddressMap[*cur.GetFromAddr()] = curIdx
		transferAddressMap[*cur.To()] = curIdx
		//}
		log.Debug("Handle tx end", "number", blockNumber, "txs length", len(txs), "index", curIdx, "duration", time.Since(start), "durations", time.Since(sstart))
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
