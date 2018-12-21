package cbft

import (
	"Platon-go/core/types"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"math/big"
	"testing"
	"Platon-go/core/ppos"
	"Platon-go/core/state"
	"Platon-go/ethdb"
	"Platon-go/core"
	"fmt"
	"errors"
	"Platon-go/consensus/ethash"
	"Platon-go/core/vm"
	"encoding/json"
)

func newTesterAccountPool() ([]discover.NodeID, error) {
	var accounts []discover.NodeID
	for _, url := range params.MainnetBootnodes {
		node, err := discover.ParseNode(url)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, node.ID)
	}
	return accounts, nil
}

func TestBlockCopy(t *testing.T) {
	header := &types.Header{
		Number: big.NewInt(1),
	}
	block1 := types.NewBlock(header, nil, nil, nil)

	block2Obj := *block1

	block2Obj.Header().Number = big.NewInt(2)

	block2 := &block2Obj
	block2.Header().Number = big.NewInt(3)

	println(block1.Number().Uint64())
	println(block2.Number().Uint64())

}

func newChainState() (*state.StateDB, error) {
	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, ethash.NewFaker(), vm.Config{}, nil)

	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		return nil, errors.New("reference statedb failed" + err.Error())
	}else {
		state = statedb
	}
	return state, nil
}


func printObject(title string, obj, logger interface{}){
	objs, _ := json.Marshal(obj)
	switch logger.(type) {
	case *testing.T:
		t := logger.(*testing.T)
		t.Log(title, string(objs), "\n")
	case *testing.B:
		b := logger.(*testing.B)
		b.Log(title, string(objs), "\n")
	}
}

func newPool() (*pposm.CandidatePool, *pposm.TicketPool) {
	configs := params.PposConfig{
		Candidate: &params.CandidateConfig{
			MaxChair: 1,
			MaxCount: 3,
			RefundBlockNumber: 	1,
		},
		TicketConfig: &params.TicketConfig {
			MaxCount: 100,
			ExpireBlockNumber: 2,
		},
	}
	return pposm.NewCandidatePool(&configs), pposm.NewTicketPool(&configs)
}

func buildPpos() *ppos {
	configs := params.PposConfig{
		Candidate: &params.CandidateConfig{
			MaxChair: 1,
			MaxCount: 3,
			RefundBlockNumber: 	1,
		},
		TicketConfig: &params.TicketConfig {
			MaxCount: 100,
			ExpireBlockNumber: 2,
		},
	}
	ppos := &ppos{
		candidatePool: pposm.NewCandidatePool(&configs),
		ticketPool: pposm.NewTicketPool(&configs),
	}

	return ppos
}

func TestNewPpos (t *testing.T) {
	configs := params.PposConfig{
		Candidate: &params.CandidateConfig{
			MaxChair: 1,
			MaxCount: 3,
			RefundBlockNumber: 	1,
		},
		TicketConfig: &params.TicketConfig {
			MaxCount: 100,
			ExpireBlockNumber: 2,
		},
	}
	ppos := &ppos{
		candidatePool: pposm.NewCandidatePool(&configs),
		ticketPool: pposm.NewTicketPool(&configs),
	}
	printObject("ppos.candidatePool:", ppos.candidatePool, t)
	printObject("ppos.ticketPool:", ppos.ticketPool, t)
}

// test BlockProducerIndex
func ppos_BlockProducerIndex(logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){

}