package core

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	dag3 "github.com/PlatONnetwork/PlatON-Go/core/dag"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"math/big"
	"sync"
)

var (
	contractCacheOnce sync.Once
	contractCacher    *ContractCacher
)

// ContractCacher is a global contract address cacher
type ContractCacher struct {
	contractCache map[common.Address]struct{}
	lock          sync.Mutex
}

func ContractCacherInstance() *ContractCacher {
	contractCacheOnce.Do(func() {
		log.Info("Init Contract Cacher ...")
		contractCacher = &ContractCacher{
			contractCache: make(map[common.Address]struct{}),
		}
	})
	return contractCacher
}

func (cc *ContractCacher) IsContractTx(chainConfig *params.ChainConfig, blockNumber *big.Int, msg Message, state vm.StateDB) bool {
	if state == nil {
		log.Info("state is nil")
		return false
	}
	address := msg.To()
	cc.lock.Lock()
	defer cc.lock.Unlock()
	if address == nil { // create contract
		contractAddress := crypto.CreateAddress(msg.From(), msg.Nonce())
		cc.contractCache[contractAddress] = struct{}{}
		return true
	}
	if _, ok := cc.contractCache[*address]; ok {
		return true
	}

	isContract := vm.IsPrecompiledContract(*address, chainConfig.Rules(blockNumber), gov.Gte150VersionState(state)) || state.GetCodeSize(*address) > 0
	if isContract {
		cc.contractCache[*address] = struct{}{}
	}
	return isContract
}

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

func (txDag *TxDag) MakeDagGraph(ctx *ParallelContext, exe *Executor) error {
	blockNumber, state, txs := ctx.header.Number.Uint64(), ctx.GetState(), ctx.txList

	txDag.dag = dag3.NewDag(len(txs))
	//save all transfer addresses between two contracts(precompiled and user defined)
	transferAddressMap := make(map[common.Address]int, 0)
	latestPrecompiledIndex := -1
	for index, tx := range txs {
		if tx.FromAddr(txDag.signer) == (common.Address{}) {
			log.Error("The from of the transaction cannot be resolved", "number", blockNumber, "index", index)
			continue
		}
		msg, err := tx.AsMessage(txDag.signer, ctx.header.BaseFee)
		if err != nil {
			log.Error("could not apply tx %d [%v]: %w", "number", blockNumber, "index", index)
			continue
		}
		if ContractCacherInstance().IsContractTx(exe.chainConfig, ctx.header.Number, msg, state) {
			txDag.contracts[index] = struct{}{}
			if index > 0 {
				if index-latestPrecompiledIndex > 1 {
					for begin := latestPrecompiledIndex + 1; begin < index; begin++ {
						txDag.dag.AddEdge(begin, index)
					}
				} else if index-latestPrecompiledIndex == 1 {
					txDag.dag.AddEdge(latestPrecompiledIndex, index)
				}
			}
			latestPrecompiledIndex = index
			//reset transferAddressMap
			if len(transferAddressMap) > 0 {
				transferAddressMap = make(map[common.Address]int, 0)
			}
		} else {
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
		}
	}
	/*
		// dag print info
		logVerbosity := debug.GetLogVerbosity()
		if logVerbosity == log.LvlTrace {
			buff, err := txDag.dag.Print()
			if err != nil {
				log.Error("print DAG Graph error!", "blockNumber", blockNumber, "err", err)
				return nil
			}
			log.Trace("DAG Graph", "blockNumber", blockNumber, "info", buff.String())
		}
	*/
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
