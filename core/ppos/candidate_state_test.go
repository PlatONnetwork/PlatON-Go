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

func TestWithdrawCandidate(t *testing.T) {

}

func TestGetChosens(t *testing.T) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	t.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

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
	canArr := candidatePool.GetChosens(state, 0)
	printObject("immediate elected candidates", canArr, t)

}


func TestGetElection(t *testing.T) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	t.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

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
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
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
	cans = append(cans, candidate2)
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
	cans = append(cans, candidate3)
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
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		t.Error("SetCandidate err:", err)
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
	t.Log("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
	t.Log("Whether election was successful err", err)

	arr, _ := candidatePool.GetWitness(state, 1)
	fmt.Println(arr)
}


func TestGetWitness (t *testing.T) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	t.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

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
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
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
	cans = append(cans, candidate2)
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
	cans = append(cans, candidate3)
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
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		t.Error("SetCandidate err:", err)
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
	t.Log("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
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
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	t.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

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
	t.Log("Set New Candidate ...")
	/** test SetCandidate */
	cans = append(cans, candidate)
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
	cans = append(cans, candidate2)
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
	cans = append(cans, candidate3)
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
	cans = append(cans, candidate4)
	if err := candidatePool.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		t.Error("SetCandidate err:", err)
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
	t.Log("test Election ...")
	_, err := candidatePool.Election(state, common.Hash{}, big.NewInt(20))
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


/** Unit Testing */
func TestCandidatePool_SetCandidate(t *testing.T) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	t.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)
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

func TestCandidatePool_GetCandidate(t *testing.T) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	t.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

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

func TestCandidatePool_WithdrawCandidate(t *testing.T) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		t.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	t.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

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

/** Unit BenchMarking */

func BenchmarkCandidatePool_SetCandidate(b *testing.B) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		b.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	b.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)
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

	b.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		b.Error("SetCandidate err:", err)
	}
}

func BenchmarkCandidatePool_GetCandidate(b *testing.B) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		b.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	b.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}

	b.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		b.Error("SetCandidate err:", err)
	}


	/** test GetCandidate */
	b.Log("test GetCandidate ...")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, b)

}

func BenchmarkCandidatePool_WithdrawCandidate(b *testing.B) {
	var candidatePool *pposm.CandidatePool
	var ticketPool *pposm.TicketPool
	var state *state.StateDB
	if st, err := newChainState(); nil != err {
		b.Error("Getting stateDB err", err)
	}else {state = st}
	/** test init candidatePool and ticketPool */
	candidatePool, ticketPool = newPool()
	b.Log("ticketPool.MaxCount", ticketPool.MaxCount, "ticketPool.ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}
	b.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		b.Error("SetCandidate err:", err)
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
	b.Log("Set New Candidate ...")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		b.Error("SetCandidate err:", err)
	}

	/** test GetCandidate */
	b.Log("test GetCandidate ...")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, b)

	/** test WithdrawCandidate */
	b.Log("test WithdrawCandidate ...")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	b.Log("error", ok1)

	/** test GetCandidate */
	b.Log("test GetCandidate ...")
	can2, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, b)
}

