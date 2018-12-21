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
	"sync/atomic"
	"time"
	"math/rand"
	//"Platon-go/core/ticketcache"
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
		state = statedb
	}
	return state, nil
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

func TestCandidatePoolAllCircle (t *testing.T){

	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()

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

	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 10
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			t.Logf("vote blockNumber[%v]", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", candidate.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

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
	canArr := candidatePool.GetChosens(state, 0)
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
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
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
	nodeArr, _ := candidatePool.GetWitness(state, 1)
	printObject("nodeArr", nodeArr, t)
}


/** Unit Test */

// test SetCandidate
func candidate_SetCandidate (logFn func (args ...interface{}), errFn func (args ...interface{})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)
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

	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}
}
func TestCandidatePool_SetCandidate(t *testing.T) {
	candidate_SetCandidate(t.Log, t.Error)
}
func BenchmarkCandidatePool_SetCandidate(b *testing.B) {
	candidate_SetCandidate(b.Log, b.Error)
}


// test GetCandidate
func candidate_GetCandidate(logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}

	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, logger)
}
func TestCandidatePool_GetCandidate(t *testing.T) {
	candidate_GetCandidate(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_GetCandidate(b *testing.B) {
	candidate_GetCandidate(b, b.Log, b.Error)
}

// test SetCandidateExtra
func candidate_SetCandidateExtra (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)
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

	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}

	/** test SetCndidateExtra */
	if err := candidatePool.SetCandidateExtra(state, candidate.CandidateId, "LALALALA"); nil != err {
		errFn("SetCndidateExtra err:", err)
	}
	/** test GetCandidate  */
	if can, err := candidatePool.GetCandidate(state, candidate.CandidateId); nil != err {
		errFn("GetCandidate err:", err)
	}else {
		logFn("candidate'extra:", can.Extra)
	}
}
func TestCandidatePool_SetCandidateExtra(t *testing.T) {
	candidate_SetCandidateExtra(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_SetCandidateExtra(b *testing.B) {
	candidate_SetCandidateExtra(b, b.Log, b.Error)
}

// test WithdrawCndidate
func candidate_WithdrawCandidate (logger interface {}, logFn func (args ... interface{}), errFn func (args ... interface {})) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
	}

	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, logger)

	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)

	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can2, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, logger)
}
func TestCandidatePool_WithdrawCandidate(t *testing.T) {
	candidate_WithdrawCandidate(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_WithdrawCandidate(b *testing.B) {
	candidate_WithdrawCandidate(b, b.Log, b.Error)
}


// test GetChosens
func candidate_GetChosens (logger interface {}, logFn func (args ... interface{}), errFn func (args ... interface {})) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
	}

	/** test GetChosens */
	logFn("test GetChosens ...")
	canArr := candidatePool.GetChosens(state, 0)
	printObject("immediate elected candidates", canArr, logger)
}
func TestCandidatePool_GetChosens(t *testing.T) {
	candidate_GetChosens(t, t.Log, t.Error)
}
func  BenchmarkCandidatePool_GetChosens(b *testing.B) {
	candidate_GetChosens(b, b.Log, b.Error)
}


// test GetChairpersons
func candidate_GetChairpersons (logger interface {}, logFn func (args ... interface {}), errFn func (args ... interface{})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test switch */
	logFn("test Switch ...")
	flag := candidatePool.Switch(state)
	logFn("Switch was success ", flag)

	/** test GetChairpersons */
	logFn("test GetChairpersons ...")
	canArr := candidatePool.GetChairpersons(state)
	printObject("Witnesses", canArr, logger)
}
func TestGet_Chairpersons (t *testing.T) {
	candidate_GetChairpersons(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_GetChairpersons(b *testing.B) {
	candidate_GetChairpersons(b, b.Log, b.Error)
}


// test GetWitness
func candidate_GetWitness (logger interface {}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ := candidatePool.GetWitness(state, 1)
	printObject("next Witnesses", canArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := candidatePool.Switch(state)
	logFn("Switch was success ", flag)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ = candidatePool.GetWitness(state, 0)
	printObject(" current Witnesses", canArr, logger)
}
func TestCandidatePool_GetWitness(t *testing.T) {
	candidate_GetWitness(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_GetWitness(b *testing.B) {
	candidate_GetWitness(b, b.Log, b.Error)
}

// test GetAllWitness
func candidate_GetAllWitness (logger interface {}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ := candidatePool.GetWitness(state, 1)
	printObject("next Witnesses", canArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := candidatePool.Switch(state)
	logFn("Switch was success ", flag)


	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ = candidatePool.GetWitness(state, 0)
	printObject(" current Witnesses", canArr, logger)

	/** test Election */
	logFn("test Election again ...")
	_, err = candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election again was successful err", err)

	/** test GetAllWitness */
	logFn("test GetAllWitness ...")
	preArr, curArr, nextArr, _ := candidatePool.GetAllWitness(state)
	printObject("previous Witness", preArr, logger)
	printObject(" current Witnesses", curArr, logger)
	printObject(" next Witnesses", nextArr, logger)
}
func TestCandidatePool_GetAllWitness(t *testing.T) {
	candidate_GetAllWitness(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_GetAllWitness(b *testing.B) {
	candidate_GetAllWitness(b, b.Log, b.Error)
}

// test GetDefeat
func candidate_GetDefeat (logger interface {}, logFn func (args ... interface {}), errFn func (args ... interface {})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")


	/** test Election */
	logFn("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/**  */
	printObject("candidatePool:", *candidatePool, logger)
	/** test MaxChair */
	logFn("test MaxChair:", candidatePool.MaxChair())
	/**test Interval*/
	logFn("test Interval:", candidatePool.GetRefundInterval())

	/** test switch */
	logFn("test Switch ...")
	flag := candidatePool.Switch(state)
	logFn("Switch was success ", flag)

	/** test GetChairpersons */
	logFn("test GetChairpersons ...")
	canArr := candidatePool.GetChairpersons(state)
	printObject("Witnesses", canArr, logger)


	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)

	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can2, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, logger)


	/** test GetDefeat */
	logFn("test GetDefeat ...")
	defeatArr, _ := candidatePool.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("can be refund defeats", defeatArr, logger)

	/** test IsDefeat */
	logFn("test IsDefeat ...")
	flag, _ = candidatePool.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	logFn("isdefeat", flag)



	/** test RefundBalance */
	logFn("test RefundBalance ...")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(11)))
	logFn("RefundBalance err", err)

	/** test RefundBalance again */
	logFn("test RefundBalance again ...")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(11)))
	logFn("RefundBalance again err", err)


	/** test GetOwner */
	logFn("test GetOwner ...")
	addr := candidatePool.GetOwner(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	logFn("Benefit address", addr.String())

	/**  test GetWitness */
	logFn("test GetWitness ...")
	nodeArr, _ := candidatePool.GetWitness(state, 0)
	printObject("nodeArr", nodeArr, logger)
}
func TestGetDefeat(t *testing.T) {
	candidate_GetDefeat(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_GetDefeat(b *testing.B) {
	candidate_GetDefeat(b, b.Log, b.Error)
}

// test GetOwner
func candidate_GetOwner (logger interface {}, logFn func (args ... interface {}), errFn func (args ... interface {})) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}

	/** test GetOwner */
	ownerAddr := candidatePool.GetOwner(state, candidate.CandidateId)
	logFn("Getting Onwer's Address:", ownerAddr.String())
}
func TestCandidatePool_GetOwner(t *testing.T) {
	candidate_GetOwner(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_GetOwner(b *testing.B) {
	candidate_GetOwner(b, b.Log, b.Error)
}

// test GetRefundInterval
func candidate_GetRefundInterval(logger interface{}, logFn func (args ... interface {}), errFn func (args ... interface {})) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	/** test  GetRefundInterval*/
	num := candidatePool.GetRefundInterval()
	logFn("RefundInterval:", num)
	fmt.Println(state.Error())
}
func TestCandidatePool_GetRefundInterval(t *testing.T) {
	candidate_GetRefundInterval(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_GetRefundInterval(b *testing.B) {
	candidate_GetRefundInterval(b, b.Log, b.Error)
}

// test MaxChair
func candidate_MaxChair (logger interface{}, logFn func (args ... interface {}), errFn func (args ... interface {})) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	/** test  MaxChair*/
	num := candidatePool.MaxChair()
	logFn("MaxChair:", num)
	fmt.Println(state.Error())
}
func TestCandidatePool_MaxChair(t *testing.T) {
	candidate_MaxChair(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_MaxChair(b *testing.B) {
	candidate_MaxChair(b, b.Log,b.Error)
}

// test IsChosens
func candidate_IsChosens(logger interface{}, logFn func (args ... interface {}), errFn func (args ... interface {})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
	}

	/** test GetChosens */
	logFn("test IsChosens ...")
	if flag, err := candidatePool.IsChosens(state, candidate2.CandidateId); nil != err {
		errFn("IsChosens err:", err)
	}else {
		logFn("IsChosens success", flag)
	}
}
func TestCandidatePool_IsChosens(t *testing.T) {
	candidate_IsChosens(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_IsChosens(b *testing.B) {
	candidate_IsChosens(b, b.Log, b.Error)
}

// test IsDefeat
func candidate_IsDefeat (logger interface{}, logFn func (args ... interface {}), errFn func (args ... interface {})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")


	/** test Election */
	logFn("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/**  */
	printObject("candidatePool:", *candidatePool, logger)
	/** test MaxChair */
	logFn("test MaxChair:", candidatePool.MaxChair())
	/**test Interval*/
	logFn("test Interval:", candidatePool.GetRefundInterval())

	/** test switch */
	logFn("test Switch ...")
	flag := candidatePool.Switch(state)
	logFn("Switch was success ", flag)

	/** test GetChairpersons */
	logFn("test GetChairpersons ...")
	canArr := candidatePool.GetChairpersons(state)
	printObject("Witnesses", canArr, logger)


	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)

	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can2, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, logger)

	/** test GetDefeat */
	logFn("test GetDefeat ...")
	defeatArr, _ := candidatePool.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("can be refund defeats", defeatArr, logger)

	/** test IsDefeat */
	logFn("test IsDefeat ...")
	flag, _ = candidatePool.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	logFn("isdefeat", flag)
}
func TestCandidatePool_IsDefeat(t *testing.T) {
	candidate_IsDefeat(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_IsDefeat(b *testing.B) {
	candidate_IsDefeat(b, b.Log, b.Error)
}

// test RefundBalance
func candidate_RefundBalance(logger interface{}, logFn func (args ... interface {}), errFn func (args ... interface {})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")


	/** test Election */
	logFn("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/**  */
	printObject("candidatePool:", *candidatePool, logger)
	/** test MaxChair */
	logFn("test MaxChair:", candidatePool.MaxChair())
	/**test Interval*/
	logFn("test Interval:", candidatePool.GetRefundInterval())

	/** test switch */
	logFn("test Switch ...")
	flag := candidatePool.Switch(state)
	logFn("Switch was success ", flag)

	/** test GetChairpersons */
	logFn("test GetChairpersons ...")
	canArr := candidatePool.GetChairpersons(state)
	printObject("Witnesses", canArr, logger)


	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)

	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can2, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, logger)


	/** test GetDefeat */
	logFn("test GetDefeat ...")
	defeatArr, _ := candidatePool.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("can be refund defeats", defeatArr, logger)

	/** test IsDefeat */
	logFn("test IsDefeat ...")
	flag, _ = candidatePool.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	logFn("isdefeat", flag)



	/** test RefundBalance */
	logFn("test RefundBalance ...")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(11)))
	logFn("RefundBalance err", err)

	/** test RefundBalance again */
	logFn("test RefundBalance again ...")
	err = candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(11)))
	logFn("RefundBalance again err", err)
}
func TestCandidatePool_RefundBalance(t *testing.T) {
	candidate_RefundBalance(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_RefundBalance(b *testing.B) {
	candidate_RefundBalance(b, b.Log, b.Error)
}

// test UpdateElectedQueue
func candidate_UpdateElectedQueue(logger interface{}, logFn func (args ... interface {}), errFn func (args ... interface {})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")
	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"))
	logFn("GetCandidate", can)

	/** test UpdateElectedQueue */
	logFn("test UpdateElectedQueue")
	if err := candidatePool.UpdateElectedQueue(state, []discover.NodeID{candidate2.CandidateId, candidate3.CandidateId}...); nil != err {
		errFn("UpdateElectedQueue err", err)
	}else {
		logFn("test UpdateElectedQueue success")
	}
}
func TestCandidatePool_UpdateElectedQueue(t *testing.T) {
	candidate_UpdateElectedQueue(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_UpdateElectedQueue(b *testing.B) {
	candidate_UpdateElectedQueue(b, b.Log, b.Error)
}

// test Election
func candidate_Election (logger interface{}, logFn func (args ...interface{}), errFn func (args ...interface{})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}

	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	arr, _ := candidatePool.GetWitness(state, 1)
	fmt.Println(arr)
}
func TestCandidatePool_Election(t *testing.T) {
	candidate_Election(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_Election(b *testing.B) {
	candidate_Election(b, b.Log, b.Error)
}

// test Switch
func candidate_Switch (logger interface{}, logFn func (args ... interface {}), errFn func (args ... interface{})){
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		errFn("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	logFn("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	// cache
	cans := make([]*types.Candidate, 0)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate2)
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate3)
	if err := candidatePool.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
		errFn("SetCandidate err:", err)
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
	logFn("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 13
	timeMap := make(map[uint32]int64)
	fmt.Println("VOTING START .............................................................")
	for i := 0; i < voteNum ; i++ {
		can := cans[rand.Intn(4)]

		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(10)
		state.SubBalance(voteOwner, deposit)
		state.AddBalance(common.TicketPoolAddr, deposit)
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			logFn("vote blockNumber:", tempBlockNumber.Uint64())
		}

		if i == 2 {
			fmt.Println("release ticket,start ############################################################")
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")


	/** test Election */
	logFn("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/**  */
	printObject("candidatePool:", *candidatePool, logger)
	/** test MaxChair */
	logFn("test MaxChair:", candidatePool.MaxChair())
	/**test Interval*/
	logFn("test Interval:", candidatePool.GetRefundInterval())

	/** test switch */
	logFn("test Switch ...")
	flag := candidatePool.Switch(state)
	logFn("Switch was success ", flag)
}
func TestCandidatePool_Switch(t *testing.T) {
	candidate_Switch(t, t.Log, t.Error)
}
func BenchmarkCandidatePool_Switch(b *testing.B) {
	candidate_Switch(b, b.Log, b.Error)
}

