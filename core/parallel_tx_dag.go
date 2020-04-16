package core

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	dag3 "github.com/PlatONnetwork/PlatON-Go/core/dag"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
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

func (txDag *TxDag) MakeDagGraph(state *state.StateDB, txs []*types.Transaction) error {
	txDag.dag = dag3.NewDag(len(txs))
	//save all transfer addresses between two contracts(precompiled and user defined)
	transferAddressMap := make(map[common.Address]int, 0)
	latestPrecompiledIndex := -1
	for curIdx, cur := range txs {
		if cur.GetFromAddr() == nil {
			if from, err := types.Sender(txDag.signer, cur); err != nil {
				return err
			} else {
				cur.SetFromAddr(&from)
			}
		}

		if cur.To() == nil || vm.IsPrecompiledContract(*cur.To()) || state.GetCodeSize(*cur.To()) > 0 {
			log.Debug("found contract tx", "idx", curIdx, "txHash", cur.Hash(), "txGas", cur.Gas(), "fromAddr", *cur.GetFromAddr(), "toAddr", *cur.To())
			txDag.contracts[curIdx] = struct{}{}
			if curIdx > 0 {
				if curIdx-latestPrecompiledIndex > 1 {
					for begin := latestPrecompiledIndex + 1; begin < curIdx; begin++ {
						txDag.dag.AddEdge(begin, curIdx)
					}
				} else if curIdx-latestPrecompiledIndex == 1 {
					txDag.dag.AddEdge(latestPrecompiledIndex, curIdx)
				}
			}
			latestPrecompiledIndex = curIdx
			//reset transferAddressMap
			if len(transferAddressMap) > 0 {
				transferAddressMap = make(map[common.Address]int, 0)
			}
		} else {
			log.Debug("found transfer tx", "idx", curIdx, "txHash", cur.Hash(), "txGas", cur.Gas(), "fromAddr", *cur.GetFromAddr(), "toAddr", *cur.To())
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

func (txDag *TxDag) IsContract(idx int) bool {
	if _, ok := txDag.contracts[idx]; ok {
		return true
	}
	return false
}
