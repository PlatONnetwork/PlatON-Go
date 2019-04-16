package pposm_test

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/ticketcache"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"math/big"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos_storage"
)

func TestTicketProcess(t *testing.T) {
	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	ppos_storage.NewPPosTemp(db)
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, nil, vm.Config{}, nil)

	configs := params.PposConfig{
		CandidateConfig: &params.CandidateConfig{
			Threshold:         "100",
			DepositLimit:      10,
			MaxChair:          1,
			MaxCount:          3,
			RefundBlockNumber: 1,
		},
		TicketConfig: &params.TicketConfig{
			TicketPrice:       "1",
			MaxCount:          10000,
			ExpireBlockNumber: 4,
		},
	}

	candidatePoolContext := pposm.NewCandidatePoolContext(&configs)

	ticketPoolContext := pposm.NewTicketPoolContext(&configs)

	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		fmt.Println("reference statedb failed", err)
	} else {
		state = statedb
	}

	candidate := &types.Candidate{
		Deposit:     new(big.Int).SetUint64(100),
		BlockNumber: new(big.Int).SetUint64(7),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:     6,
		Host:        "10.0.0.1",
		Port:        "8548",
		Owner:       common.HexToAddress("0x12"),
	}

	fmt.Println("seting new k-v \n", candidate)
	// test SetCandidate
	if err := candidatePoolContext.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		fmt.Println("SetCandidate err:", err)
	}

	// set ownerList
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 100
	timeMap := make(map[uint32]int64)
	var releaseTime int64 = 0
	for i := 0; i < voteNum; i++ {
		//go func() {
		startTime := time.Now().UnixNano() / 1e6
		voteOwner := ownerList[rand.Intn(2)]
		deposit := new(big.Int).SetUint64(1)
		var voteNum uint32 = 2
		state.SubBalance(voteOwner, new(big.Int).Mul(deposit, new(big.Int).SetUint64(uint64(voteNum))))
		state.AddBalance(common.TicketPoolAddr, new(big.Int).Mul(deposit, new(big.Int).SetUint64(uint64(voteNum))))
		tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
		if i < 2 {
			tempBlockNumber.SetUint64(6)
			t.Logf("vote blockNumber[%v]", tempBlockNumber.Uint64())
		}

		txHash := common.Hash{}
		txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
		state.Prepare(txHash, common.Hash{}, 1)

		_, err := ticketPoolContext.VoteTicket(state, voteOwner, voteNum, deposit, candidate.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		if i == 1 {
			var tempBlockNumber uint64 = 6
			for i := 0; i < 4; i++ {
				ticketPoolContext.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
		//}()
	}

	candidate = candidatePoolContext.GetCandidate(state, candidate.CandidateId, blockNumber)


	ticketIds := ticketPoolContext.GetCandidateTicketIds(state, candidate.CandidateId)
	ticketList := ticketPoolContext.GetTicketList(state, ticketIds)
	t.Logf("ticketListSize=%d\n", len(ticketList))
	expireTicketIds := ticketPoolContext.GetExpireTicketIds(state, blockNumber)
	//expireTickets, err := ticketPool.GetTicketList(state, expireTicketIds)

	surplusQuantity := ticketPoolContext.GetPoolNumber(state)
	epoch := ticketPoolContext.GetCandidateEpoch(state, candidate.CandidateId)

	t.Logf("ticketPoolSize:[%d],expireTicketListSize:[%d],candidate.TicketPool:[%d],tcount:[%d],epoch:[%d]\n",
		surplusQuantity, len(expireTicketIds), len(ticketIds), state.GetPPOSCache().GetCandidateTicketCount(candidate.CandidateId), epoch)
	t.Logf("ticketPoolBalance[%v]", state.GetBalance(common.TicketPoolAddr))
	fmt.Println("------all ticket-----")
	//for _, ticket := range ticketList {
	//	fmt.Printf("ticket:%+v,ticketId:[%v]\n", ticket, ticket.TicketId.Hex())
	//}

	candidate = candidatePoolContext.GetCandidate(state, candidate.CandidateId, blockNumber)
	ticketIds = ticketPoolContext.GetCandidateTicketIds(state, candidate.CandidateId)


	if err := ticketPoolContext.Notify(state, blockNumber); err != nil {
		t.Error("Execute HandleExpireTicket error", err)
	}
	expireTicketIds = ticketPoolContext.GetExpireTicketIds(state, blockNumber)
	surplusQuantity = ticketPoolContext.GetPoolNumber(state)
	epoch = ticketPoolContext.GetCandidateEpoch(state, candidate.CandidateId)
	t.Logf("ticketPoolSize:[%d],expireTicketListSize:[%d],candidate.TicketPool:[%d],tcount:[%d],epoch:[%d]\n",
		surplusQuantity, len(expireTicketIds), len(ticketIds), ticketPoolContext.GetCandidateTicketCount(state, candidate.CandidateId), epoch)
	t.Logf("ticketPoolBalance[%v]", state.GetBalance(common.TicketPoolAddr))

	blockHash := common.Hash{}
	blockHash.SetBytes([]byte("3b41e0aee38c1a1f959a6aaae678d86f1e6af59617d2f667bb2ef5527779c861"))
	luckyTicketId, err := ticketPoolContext.SelectionLuckyTicket(state, candidate.CandidateId, blockHash)
	if nil != err {
		t.Error("SelectionLuckyTicket error", err)
	}
	//selectedTicketIndex := rand.Intn(len(ticketList))
	//selectedTicketId := ticketList[selectedTicketIndex].TicketId
	t.Logf("-----------Start releasing a ticket 【%v】-----------\n", luckyTicketId.Hex())
	tempTime := time.Now().UnixNano() / 1e6
	err = ticketPoolContext.ReturnTicket(state, candidate.CandidateId, luckyTicketId, blockNumber)
	if nil != err {
		t.Error("ReleaseSelectedTicket error", err)
	}
	releaseTime = (time.Now().UnixNano() / 1e6) - tempTime
	ticket := ticketPoolContext.GetTicket(state, luckyTicketId)
	t.Logf("lucky ticket :%+v", ticket)

	expireTicketIds = ticketPoolContext.GetExpireTicketIds(state, blockNumber)
	surplusQuantity = ticketPoolContext.GetPoolNumber(state)
	epoch = ticketPoolContext.GetCandidateEpoch(state, candidate.CandidateId)
	t.Logf("After processing the expired ticket block height：[%d]", blockNumber)
	t.Logf("ticketPoolSize:[%d],expireTicketListSize:[%d],candidate.TicketPool:[%d],tcount:[%d],epoch:[%d]\n",
		surplusQuantity, len(expireTicketIds), len(ticketIds), ticketPoolContext.GetCandidateTicketCount(state, candidate.CandidateId), epoch)
	t.Log("expireTicketList info", "blockNumber", )
	t.Logf("ticketPoolBalance[%v]", state.GetBalance(common.TicketPoolAddr))

	fmt.Println("It takes time to release a vote：", releaseTime, "ms")
	fmt.Println("When the 10,000th ticket is voted, it takes time：", timeMap[10000], "ms")
	fmt.Println("When the 5000th ticket is used, the voting takes time：", timeMap[5000], "ms")
	fmt.Println("When the 1000th ticket is taken, the time taken for voting：", timeMap[1000], "ms")
	fmt.Println("When the 500th ticket is taken, the time taken for voting：", timeMap[500], "ms")
	fmt.Println("When the 100th ticket is taken, the time taken for voting：", timeMap[100], "ms")
	fmt.Println("When voting for the 50th ticket, it takes time to vote.：", timeMap[50], "ms")
	fmt.Println("When the 10th ticket is taken, the voting takes time：", timeMap[10], "ms")
	fmt.Println("When the first ticket is taken, the time taken for voting：", timeMap[1], "ms")
}

func initParam() (*state.StateDB, *pposm.CandidatePoolContext, *pposm.TicketPoolContext) {
	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, nil, vm.Config{}, nil)
	ticketcache.NewTicketIdsCache(db)
	configs := params.PposConfig{
		CandidateConfig: &params.CandidateConfig{
			Threshold:         "100",
			DepositLimit:      10,
			MaxChair:          1,
			MaxCount:          3,
			RefundBlockNumber: 1,
		},
		TicketConfig: &params.TicketConfig{
			TicketPrice:       "1",
			MaxCount:          10000,
			ExpireBlockNumber: 100,
		},
	}

	candidatePoolContextContext := pposm.NewCandidatePoolContext(&configs)

	ticketPoolContext := pposm.NewTicketPoolContext(&configs)

	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		fmt.Println("reference statedb failed", "err", err)
	} else {
		state = statedb
	}
	txHash := common.Hash{}
	txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Second()))))
	state.Prepare(txHash, common.Hash{}, 1)
	return state, candidatePoolContextContext, ticketPoolContext
}

func TestTicketPool_VoteTicket(t *testing.T) {
	state, candidatePoolContext, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit:     new(big.Int).SetUint64(100),
		BlockNumber: new(big.Int).SetUint64(7),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:     6,
		Host:        "10.0.0.1",
		Port:        "8548",
		Owner:       common.HexToAddress("0x12"),
	}

	if err := candidatePoolContext.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate fail", "err", err)
	}

	successNumber, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(10))
	if nil != err {
		t.Error("vote ticket fail", "err", err)
	} else {
		t.Log("vote ticket success", "successNumber", successNumber)
	}
}

func TestTicketPool_GetExpireTicketIds(t *testing.T) {
	state, _, ticketPool := initParam()

	list := ticketPool.GetExpireTicketIds(state, new(big.Int).SetUint64(10))
	t.Logf("Execute TestTicketPool_GetExpireTicketIds success, listSize=%d", len(list))
}

func TestTicketPool_GetTicketList(t *testing.T) {
	state, candidatePoolContext, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit:     new(big.Int).SetUint64(100),
		BlockNumber: new(big.Int).SetUint64(7),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:     6,
		Host:        "10.0.0.1",
		Port:        "8548",
		Owner:       common.HexToAddress("0x12"),
	}

	if err := candidatePoolContext.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(10))
	if nil != err {
		t.Error("vote ticket fail", "err", err)
	}

	list := ticketPool.GetCandidateTicketIds(state, candidate.CandidateId)
	t.Logf("Execute TestTicketPool_GetTicketList success, list=%d", len(list))
}

func TestTicketPool_GetTicket(t *testing.T) {
	state, candidatePoolContext, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit:     new(big.Int).SetUint64(100),
		BlockNumber: new(big.Int).SetUint64(7),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:     6,
		Host:        "10.0.0.1",
		Port:        "8548",
		Owner:       common.HexToAddress("0x12"),
	}

	if err := candidatePoolContext.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(10))
	if nil != err {
		t.Error("vote ticket fail", "err", err)
	}

	ticket := ticketPool.GetTicket(state, state.TxHash())
	if ticket == nil {
		t.Error("getTicket fail", "err", err)
	}
}

func TestTicketPool_DropReturnTicket(t *testing.T) {
	state, candidatePoolContext, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit:     new(big.Int).SetUint64(100),
		BlockNumber: new(big.Int).SetUint64(7),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:     6,
		Host:        "10.0.0.1",
		Port:        "8548",
		Owner:       common.HexToAddress("0x12"),
	}

	if err := candidatePoolContext.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(6))
	if nil != err {
		t.Error("vote ticket fail", "err", err)
	}

	err = ticketPool.DropReturnTicket(state, candidate.BlockNumber, candidate.CandidateId)
	if nil != err {
		t.Error("dropReturnTicket fail", "err", err)
	}
}

func TestTicketPool_ReturnTicket(t *testing.T) {
	state, candidatePoolContext, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit:     new(big.Int).SetUint64(100),
		BlockNumber: new(big.Int).SetUint64(7),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:     6,
		Host:        "10.0.0.1",
		Port:        "8548",
		Owner:       common.HexToAddress("0x12"),
	}

	if err := candidatePoolContext.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(10))
	if nil != err {
		t.Error("vote ticket fail", "err", err)
	}

	err = ticketPool.ReturnTicket(state, candidate.CandidateId, state.TxHash(), new(big.Int).SetUint64(10))
	if nil != err {
		t.Error("returnTicket fail", "err", err)
	}
}

func TestTicketPool_Notify(t *testing.T) {
	state, candidatePoolContext, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit:     new(big.Int).SetUint64(100),
		BlockNumber: new(big.Int).SetUint64(7),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:     6,
		Host:        "10.0.0.1",
		Port:        "8548",
		Owner:       common.HexToAddress("0x12"),
	}

	if err := candidatePoolContext.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(6))
	if nil != err {
		t.Error("vote ticket fail", "err", err)
	}

	err = ticketPool.Notify(state, new(big.Int).SetUint64(10))
	if nil != err {
		t.Error("notify fail", "err", err)
	}
}

func TestTicketPool_SelectionLuckyTicket(t *testing.T) {
	state, candidatePoolContext, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit:     new(big.Int).SetUint64(100),
		BlockNumber: new(big.Int).SetUint64(7),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:     6,
		Host:        "10.0.0.1",
		Port:        "8548",
		Owner:       common.HexToAddress("0x12"),
	}

	if err := candidatePoolContext.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		t.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(6))
	if nil != err {
		t.Error("vote ticket fail", "err", err)
	}

	blockHash := common.Hash{}
	blockHash.SetBytes([]byte("3b41e0aee38c1a1f959a6aaae678d86f1e6af59617d2f667bb2ef5527779c861"))
	_, err = ticketPool.SelectionLuckyTicket(state, candidate.CandidateId, blockHash)
	if nil != err {
		t.Error("selectionLuckyTicket fail", "err", err)
	}
}

func TestTicketPool_GetPoolNumber(t *testing.T) {
	state, _, ticketPool := initParam()

	sum := ticketPool.GetPoolNumber(state)
	t.Log("getPoolNumber success", "sum", sum)
}

func TestTicketPool_GetCandidateTicketIds(t *testing.T) {
	state, _, ticketPool := initParam()

	list := ticketPool.GetCandidateTicketIds(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	t.Logf("Execute TestTicketPool_GetCandidateTicketIds success, listSize=%d", len(list))
}

func TestTicketPool_GetCandidateEpoch(t *testing.T) {
	state, _, ticketPool := initParam()

	epoch := ticketPool.GetCandidateEpoch(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	t.Logf("Execute TestTicketPool_GetCandidateTicketIds success, epoch=%d", epoch)
}

func TestTicketPool_GetTicketPrice(t *testing.T) {
	state, _, ticketPool := initParam()

	price := ticketPool.GetTicketPrice(state)
	if price.Cmp(big.NewInt(0)) <= 0 {
		t.Error("getTicketPrice fail")
	}
}
