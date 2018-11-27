package depos_test

import (
	"testing"
	"fmt"
	"Platon-go/params"
	"math/big"
	"Platon-go/common"
	"Platon-go/ethdb"
	"Platon-go/core"
	"Platon-go/core/vm"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/consensus/ethash"
	"Platon-go/p2p/discover"

	"Platon-go/core/dpos"
	"errors"
	"encoding/json"
)

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
		/*var isgenesis bool
		if blockchain.CurrentBlock().NumberU64() == blockchain.Genesis().NumberU64() {
			isgenesis = true
		}
		*//** test init candidatePool *//*
		if pool, err := depos.NewCandidatePool(*//*statedb,*//* &configs*//*, isgenesis*//*); nil != err {
			fmt.Println("init candidatePool err", err)
		}else{
			candidatePool = pool
		}*/
		state = statedb
	}
	return state, nil
}

func newCandidatePool() *depos.CandidatePool {
	configs := params.DposConfig{
		Candidate: &params.CandidateConfig{
			MaxChair: 1,
			MaxCount: 3,
			RefundBlockNumber: 	1,
		},
	}
	return depos.NewCandidatePool(&configs)
}

func printObject(title string, obj interface{}, t *testing.T){
	objs, _ := json.Marshal(obj)
	t.Log(title, string(objs), "\n")
}

func TestInitCandidatePoolByConfig (t *testing.T){

	var candidatePool *depos.CandidatePool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool */
	candidatePool = newCandidatePool()

	//state.Commit(false)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}

	fmt.Println("Set New Candidate \n")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		fmt.Println("SetCandidate err:", err)
	}


	/** test GetCandidate */
	fmt.Println("test GetCandidate")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"))
	fmt.Println("GetCandidate", can)


	/** test WithdrawCandidate */
	fmt.Println("test WithdrawCandidate")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	fmt.Println("error", ok1)

	/** test WithdrawCandidate again */
	fmt.Println("test WithdrawCandidate again")
	ok2 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(10)), new(big.Int).SetUint64(uint64(11)))
	fmt.Println("error", ok2)

	/** test GetChosens */
	fmt.Println("test GetChosens")
	canArr := candidatePool.GetChosens(state)
	depos.PrintObject("Elected candidates", canArr)

	/** test GetChairpersons */
	fmt.Println("test GetChairpersons")
	canArr = candidatePool.GetChairpersons(state)
	depos.PrintObject("Witnesses", canArr)

	/** test GetDefeat */
	fmt.Println("test GetDefeat")
	defeatArr, _ := candidatePool.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	depos.PrintObject("can be refund defeats", defeatArr)

	/** test IsDefeat */
	fmt.Println("test IsDefeat")
	flag, _ := candidatePool.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	depos.PrintObject("isdefeat", flag)

	/** test Election */
	fmt.Println("test Election")
	_, err := candidatePool.Election(state)
	fmt.Println("whether election was successful", err)

	/** test RefundBalance */
	fmt.Println("test RefundBalance")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(11)))
	fmt.Println("err", err)

	/** test RefundBalance again */
	fmt.Println("test RefundBalance again")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012343"), new(big.Int).SetUint64(uint64(11)))
	fmt.Println("err", err)


	/** test GetOwner */
	fmt.Println("test GetOwner")
	addr := candidatePool.GetOwner(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	fmt.Println("收益地址", addr.String())

	/**  test GetWitness */
	fmt.Println("test GetWitness")
	nodeArr, _ := candidatePool.GetWitness(state, 0)
	fmt.Printf("nodeArr := %+v", nodeArr)
}

func TestSetCandidate (t *testing.T){
	var candidatePool *depos.CandidatePool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool */
	candidatePool = newCandidatePool()

	//state.Commit(false)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}

	t.Log("Set New Candidate \n")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate err:", err)
	}

}


func TestGetCandidate (t *testing.T) {
	var candidatePool *depos.CandidatePool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool */
	candidatePool = newCandidatePool()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}

	t.Log("Set New Candidate \n")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate err:", err)
	}


	/** test GetCandidate */
	t.Log("test GetCandidate")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, t)

}

func TestWithdrawCandidate(t *testing.T) {
	var candidatePool *depos.CandidatePool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool */
	candidatePool = newCandidatePool()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	t.Log("Set New Candidate \n")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate err:", err)
	}

	candidate2 := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(99),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"),
		TxIndex:  		5,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x15"),

	}
	t.Log("Set New Candidate \n")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		t.Error("SetCandidate err:", err)
	}

	/** test GetCandidate */
	t.Log("test GetCandidate")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, t)

	/** test WithdrawCandidate */
	t.Log("test WithdrawCandidate")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	t.Log("error", ok1)

	/** test GetCandidate */
	t.Log("test GetCandidate")
	can2, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, t)
}

func TestGetChosens(t *testing.T) {
	var candidatePool *depos.CandidatePool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool */
	candidatePool = newCandidatePool()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	t.Log("Set New Candidate \n")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate err:", err)
	}

	candidate2 := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(99),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"),
		TxIndex:  		5,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x15"),

	}
	t.Log("Set New Candidate \n")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		t.Error("SetCandidate err:", err)
	}

	/** test GetChosens */
	t.Log("test GetChosens")
	canArr := candidatePool.GetChosens(state)
	printObject("immediate elected candidates", canArr, t)



}