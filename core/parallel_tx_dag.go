package core

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	dag3 "github.com/PlatONnetwork/PlatON-Go/core/dag"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
)

type TxDag struct {
	dag    *dag3.Dag
	signer types.Signer
}

func NewTxDag(signer types.Signer) *TxDag {
	txDag := &TxDag{
		signer: signer,
	}
	return txDag
}

func (txDag *TxDag) MakeDagGraph(txs []*types.Transaction) error {
	txDag.dag = dag3.NewDag(len(txs))
	tempMap := make(map[common.Address]int, 0)
	latestPrecompiledIndex := -1
	for curIdx, cur := range txs {
		if vm.IsPrecompiled(*cur.To()) {
			if curIdx > 0 {
				for begin := latestPrecompiledIndex + 1; begin < curIdx; begin++ {
					txDag.dag.AddEdge(begin, curIdx)
				}
				//txDag.dag.AddEdge(curIdx-1, curIdx)
			}
			latestPrecompiledIndex = curIdx
			//reset tempMap
			if len(tempMap) > 0 {
				tempMap = make(map[common.Address]int, 0)
			}
		} else {
			dependFound := 0
			if cur.GetFromAddr() == nil {
				if from, err := types.Sender(txDag.signer, cur); err != nil {
					return err
				} else {
					cur.SetFromAddr(&from)
				}
			}
			if dependIdx, ok := tempMap[*cur.GetFromAddr()]; ok {
				txDag.dag.AddEdge(dependIdx, curIdx)
				dependFound++
			}
			if dependIdx, ok := tempMap[*cur.To()]; ok {
				txDag.dag.AddEdge(dependIdx, curIdx)
				dependFound++
			}

			if dependFound == 0 && latestPrecompiledIndex >= 0 {
				txDag.dag.AddEdge(latestPrecompiledIndex, curIdx)
			}

			tempMap[*cur.GetFromAddr()] = curIdx
			tempMap[*cur.To()] = curIdx
		}
	}
	return nil
}

func (txDag *TxDag) HasNext() bool {
	return txDag.dag.HasNext()
}

func (txDag *TxDag) Next() []int {
	return txDag.dag.Next()
}
