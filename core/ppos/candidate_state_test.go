package pposm_test

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

	"Platon-go/core/ppos"
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
		if pool, err := pposm.NewCandidatePool(*//*statedb,*//* &configs*//*, isgenesis*//*); nil != err {
			t.Log("init candidatePool err", err)
		}else{
			candidatePool = pool
		}*/
		state = statedb
	}
	return state, nil
}

func newCandidatePool() *pposm.CandidatePool {
	configs := params.PposConfig{
		Candidate: &params.CandidateConfig{
			MaxChair: 1,
			MaxCount: 3,
			RefundBlockNumber: 	1,
		},
	}
	return pposm.NewCandidatePool(&configs)
}

func printObject(title string, obj interface{}, t *testing.T){
	objs, _ := json.Marshal(obj)
	t.Log(title, string(objs), "\n")
}

func TestInitCandidatePoolByConfig (t *testing.T){

	var candidatePool *pposm.CandidatePool
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

	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate err:", err)
	}


	/** test GetCandidate */
	t.Log("test GetCandidate ...")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"))
	t.Log("GetCandidate", can)


	/** test WithdrawCandidate */
	t.Log("test WithdrawCandidate ...")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	t.Log("error", ok1)

	/** test WithdrawCandidate again */
	t.Log("test WithdrawCandidate again ...")
	ok2 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(10)), new(big.Int).SetUint64(uint64(11)))
	t.Log("error", ok2)

	/** test GetChosens */
	t.Log("test GetChosens ...")
	canArr := candidatePool.GetChosens(state)
	printObject("Elected candidates", canArr, t)

	/** test GetChairpersons */
	t.Log("test GetChairpersons ...")
	canArr = candidatePool.GetChairpersons(state)
	printObject("Witnesses", canArr, t)

	/** test GetDefeat */
	t.Log("test GetDefeat ...")
	defeatArr, _ := candidatePool.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("can be refund defeats", defeatArr, t)

	/** test IsDefeat */
	t.Log("test IsDefeat ...")
	flag, _ := candidatePool.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("isdefeat", flag, t)

	/** test Election */
	t.Log("test Election ...")
	_, err := candidatePool.Election(state)
	t.Log("whether election was successful", err)

	/** test RefundBalance */
	t.Log("test RefundBalance ...")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(11)))
	t.Log("err", err)

	/** test RefundBalance again */
	t.Log("test RefundBalance again ...")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012343"), new(big.Int).SetUint64(uint64(11)))
	t.Log("err", err)


	/** test GetOwner */
	t.Log("test GetOwner ...")
	addr := candidatePool.GetOwner(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	t.Log("Benefit address", addr.String())

	/**  test GetWitness */
	t.Log("test GetWitness ...")
	nodeArr, _ := candidatePool.GetWitness(state, 0)
	printObject("nodeArr", nodeArr, t)
}

func TestSetCandidate (t *testing.T){
	var candidatePool *pposm.CandidatePool
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

	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate err:", err)
	}

}


func TestGetCandidate (t *testing.T) {
	var candidatePool *pposm.CandidatePool
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

	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate err:", err)
	}


	/** test GetCandidate */
	t.Log("test GetCandidate ...")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, t)

}

func TestWithdrawCandidate(t *testing.T) {
	var candidatePool *pposm.CandidatePool
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
	t.Log("Set New Candidate ...")
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
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		t.Error("SetCandidate err:", err)
	}

	/** test GetCandidate */
	t.Log("test GetCandidate ...")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, t)

	/** test WithdrawCandidate */
	t.Log("test WithdrawCandidate ...")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	t.Log("error", ok1)

	/** test GetCandidate */
	t.Log("test GetCandidate ...")
	can2, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, t)
}

func TestGetChosens(t *testing.T) {
	var candidatePool *pposm.CandidatePool
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
	t.Log("Set New Candidate ...")
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
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		t.Error("SetCandidate err:", err)
	}

	/** test GetChosens */
	t.Log("test GetChosens ...")
	canArr := candidatePool.GetChosens(state)
	printObject("immediate elected candidates", canArr, t)

}


func TestGetElection(t *testing.T) {
	var candidatePool *pposm.CandidatePool
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
	t.Log("Set New Candidate ...")
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
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		t.Error("SetCandidate err:", err)
	}

	candidate3 := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(99),
		BlockNumber:    new(big.Int).SetUint64(6),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012342"),
		TxIndex:  		5,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x15"),

	}
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		t.Error("SetCandidate err:", err)
	}

	candidate4 := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(99),
		BlockNumber:    new(big.Int).SetUint64(6),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012343"),
		TxIndex:  		4,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x15"),

	}
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		t.Error("SetCandidate err:", err)
	}


	/** test Election */
	t.Log("test Election ...")
	_, err := candidatePool.Election(state)
	t.Log("Whether election was successful err", err)

}


func TestGetWitness (t *testing.T) {
	var candidatePool *pposm.CandidatePool
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
	t.Log("Set New Candidate ...")
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
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		t.Error("SetCandidate err:", err)
	}

	candidate3 := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(99),
		BlockNumber:    new(big.Int).SetUint64(6),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012342"),
		TxIndex:  		5,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x15"),

	}
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		t.Error("SetCandidate err:", err)
	}

	candidate4 := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(99),
		BlockNumber:    new(big.Int).SetUint64(6),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012343"),
		TxIndex:  		4,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x15"),

	}
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		t.Error("SetCandidate err:", err)
	}


	/** test Election */
	t.Log("test Election ...")
	_, err := candidatePool.Election(state)
	t.Log("Whether election was successful err", err)

	/** test switch */
	t.Log("test Switch ...")
	flag := candidatePool.Switch(state)
	t.Log("Switch was success ", flag)

	/** test GetChairpersons */
	t.Log("test GetChairpersons ...")
	canArr := candidatePool.GetChairpersons(state)
	printObject("Witnesses", canArr, t)
}


func TestGetDefeat(t *testing.T) {
	var candidatePool *pposm.CandidatePool
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
	t.Log("Set New Candidate ...")
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
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		t.Error("SetCandidate err:", err)
	}

	candidate3 := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(99),
		BlockNumber:    new(big.Int).SetUint64(6),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012342"),
		TxIndex:  		5,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x15"),

	}
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		t.Error("SetCandidate err:", err)
	}

	candidate4 := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(99),
		BlockNumber:    new(big.Int).SetUint64(6),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012343"),
		TxIndex:  		4,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x15"),

	}
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		t.Error("SetCandidate err:", err)
	}


	/** test Election */
	t.Log("test Election ...")
	_, err := candidatePool.Election(state)
	t.Log("Whether election was successful err", err)

	/**  */
	printObject("candidatePool:", *candidatePool, t)
	/** test MaxChair */
	t.Log("test MaxChair:", candidatePool.MaxChair())
	/**test Interval*/
	t.Log("test Interval:", candidatePool.GetRefundInterval())

	/** test switch */
	t.Log("test Switch ...")
	flag := candidatePool.Switch(state)
	t.Log("Switch was success ", flag)

	/** test GetChairpersons */
	t.Log("test GetChairpersons ...")
	canArr := candidatePool.GetChairpersons(state)
	printObject("Witnesses", canArr, t)


	/** test WithdrawCandidate */
	t.Log("test WithdrawCandidate ...")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	t.Log("error", ok1)

	/** test GetCandidate */
	t.Log("test GetCandidate ...")
	can2, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, t)


	/** test GetDefeat */
	t.Log("test GetDefeat ...")
	defeatArr, _ := candidatePool.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("can be refund defeats", defeatArr, t)

	/** test IsDefeat */
	t.Log("test IsDefeat ...")
	flag, _ = candidatePool.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	t.Log("isdefeat", flag)



	/** test RefundBalance */
	t.Log("test RefundBalance ...")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(11)))
	t.Log("RefundBalance err", err)

	/** test RefundBalance again */
	t.Log("test RefundBalance again ...")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(11)))
	t.Log("RefundBalance again err", err)


	/** test GetOwner */
	t.Log("test GetOwner ...")
	addr := candidatePool.GetOwner(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	t.Log("Benefit address", addr.String())

	/**  test GetWitness */
	t.Log("test GetWitness ...")
	nodeArr, _ := candidatePool.GetWitness(state, 0)
	printObject("nodeArr", nodeArr, t)
}