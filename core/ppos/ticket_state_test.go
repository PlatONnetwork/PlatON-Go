package pposm_test

import (
	"Platon-go/common"
	"Platon-go/consensus/ethash"
	"Platon-go/core"
	"Platon-go/core/ppos"
	"Platon-go/core/state"
	"Platon-go/core/ticketcache"
	"Platon-go/core/types"
	"Platon-go/core/vm"
	"Platon-go/ethdb"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func TestTicketProcess(t *testing.T)  {
	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	ticketcache.NewTicketIdsCache(db)
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, ethash.NewFaker(), vm.Config{}, nil)

	configs := params.PposConfig{
		//MaxChair: 1,
		//MaxCount: 3,
		//RefundBlockNumber: 	1,
		Candidate: &params.CandidateConfig{
			MaxChair: 1,
			MaxCount: 3,
			RefundBlockNumber: 	1,
		},
		TicketConfig: &params.TicketConfig{
			MaxCount: 			100,
			ExpireBlockNumber:	4,
		},
	}

	candidatePool := pposm.NewCandidatePool(&configs)

	ticketPool := pposm.NewTicketPool(&configs)

	t.Log("MaxCount", ticketPool.MaxCount, "ExpireBlockNumber", ticketPool.ExpireBlockNumber)

	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		fmt.Println("reference statedb failed", err)
	} else {
		state = statedb
	}

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),
	}

	fmt.Println("设置新的k-v \n", candidate)
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		fmt.Println("SetCandidate err:", err)
	}
	// set ownerList
	ownerList := []common.Address{common.HexToAddress("0x20"), common.HexToAddress("0x21")}
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)
	voteNum := 10001
	timeMap := make(map[uint32]int64)
	var releaseTime int64 = 0
	for i := 0; i < voteNum ; i++ {
		//go func() {
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
				var tempBlockNumber uint64 = 6
				for i := 0; i < 4; i++ {
					ticketPool.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
					tempBlockNumber++
				}
			}

			_, err := ticketPool.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
			if nil != err {
				fmt.Println("vote ticket error:", err)
			}
			atomic.AddUint32(&count, 1)
			timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
		//}()
	}
	for int(count) < voteNum  {
		fmt.Println("count:", count)
	}

	candidate, err := candidatePool.GetCandidate(state, candidate.CandidateId)
	if err != nil {
		fmt.Println("GetCandidate error")
		return
	}

	ticketIds, err := ticketPool.GetCandidateTicketIds(state, candidate.CandidateId)
	if nil != err {
		t.Error("GetCandidateTicketIds error", err)
	}
	ticketList, err := ticketPool.GetTicketList(state, ticketIds)

	expireTicketIds, err := ticketPool.GetExpireTicketIds(state, blockNumber)
	if nil != err {
		t.Error("GetExpireTicketIds error", err)
	}
	//expireTickets, err := ticketPool.GetTicketList(state, expireTicketIds)

	ticketPool.GetPoolNumber(state)
	candidateAttach, err := ticketPool.GetCandidateAttach(state, candidate.CandidateId)

	t.Logf("ticketPoolSize:[%d],expireTicketListSize:[%d],candidate.TicketPool:[%d],tcount:[%d],epoch:[%d]\n",
		ticketPool.SurplusQuantity, len(expireTicketIds), len(ticketIds), state.TCount(candidate.CandidateId.String()), candidateAttach.Epoch)
	t.Logf("ticketPoolBalance[%v],ticketDetailBalance[%v]", state.GetBalance(common.TicketPoolAddr), state.GetBalance(common.TicketPoolAddr))
	fmt.Println("------all ticket-----")
	for _, ticket := range ticketList {
		fmt.Printf("ticket:%+v,ticketId:[%v]\n", ticket, ticket.TicketId.Hex())
	}
	blockHash := common.Hash{}
	blockHash.SetBytes([]byte("3b41e0aee38c1a1f959a6aaae678d86f1e6af59617d2f667bb2ef5527779c861"))
	luckyTicketId, err := ticketPool.SelectionLuckyTicket(state, candidate.CandidateId, blockHash)
	if nil != err {
		t.Error("SelectionLuckyTicket error", err)
	}
	//selectedTicketIndex := rand.Intn(len(ticketList))
	//selectedTicketId := ticketList[selectedTicketIndex].TicketId
	t.Logf("-----------开始释放一张选票【%v】-----------\n", luckyTicketId.Hex())
	tempTime := time.Now().UnixNano() / 1e6
	err = ticketPool.ReturnTicket(state, candidate.CandidateId, luckyTicketId, blockNumber)
	if nil != err {
		t.Error("ReleaseSelectedTicket error", err)
	}
	releaseTime = (time.Now().UnixNano() / 1e6) - tempTime
	ticket, err := ticketPool.GetTicket(state, luckyTicketId)
	t.Logf("幸运票:%+v", ticket)

	candidate, err = candidatePool.GetCandidate(state, candidate.CandidateId)
	ticketIds, err = ticketPool.GetCandidateTicketIds(state, candidate.CandidateId)
	if nil != err {
		t.Error("GetCandidateTicketIds error", err)
	}
	expireTicketIds, err = ticketPool.GetExpireTicketIds(state, blockNumber)
	if nil != err {
		t.Error("GetExpireTicketIds error", err)
	}
	ticketPool.GetPoolNumber(state)
	candidateAttach, err = ticketPool.GetCandidateAttach(state, candidate.CandidateId)
	t.Logf("ticketPoolSize:[%d],expireTicketListSize:[%d],candidate.TicketPool:[%d],tcount:[%d],epoch:[%d]\n",
		ticketPool.SurplusQuantity, len(expireTicketIds), len(ticketIds), state.TCount(candidate.CandidateId.String()), candidateAttach.Epoch)
	t.Logf("ticketPoolBalance[%v],ticketDetailBalance[%v]", state.GetBalance(common.TicketPoolAddr), state.GetBalance(common.TicketPoolAddr))

	if err := ticketPool.Notify(state, blockNumber); err != nil {
		t.Error("Execute HandleExpireTicket error", err)
	}

	candidate, err = candidatePool.GetCandidate(state, candidate.CandidateId)
	ticketIds, err = ticketPool.GetCandidateTicketIds(state, candidate.CandidateId)
	if nil != err {
		t.Error("GetCandidateTicketIds error", err)
	}
	expireTicketIds, err = ticketPool.GetExpireTicketIds(state, blockNumber)
	if nil != err {
		t.Error("GetExpireTicketIds error", err)
	}
	ticketPool.GetPoolNumber(state)
	candidateAttach, err = ticketPool.GetCandidateAttach(state, candidate.CandidateId)
	t.Logf("处理完过期票块高为：[%d]", blockNumber)
	t.Logf("ticketPoolSize:[%d],expireTicketListSize:[%d],candidate.TicketPool:[%d],tcount:[%d],epoch:[%d]\n",
		ticketPool.SurplusQuantity, len(expireTicketIds), len(ticketIds), state.TCount(candidate.CandidateId.String()), candidateAttach.Epoch)
	t.Logf("ticketPoolBalance[%v],ticketDetailBalance[%v]", state.GetBalance(common.TicketPoolAddr), state.GetBalance(common.TicketPoolAddr))

	var temp []string
	temp = append(temp, "string")
	fmt.Println(temp==nil, len(temp), cap(temp))

	fmt.Println("释放一张选票耗时：", releaseTime, "ms")
	fmt.Println("第10000张票时，投票所耗时：", timeMap[10000], "ms")
	fmt.Println("第5000张票时，投票所耗时：", timeMap[5000], "ms")
	fmt.Println("第1000张票时，投票所耗时：", timeMap[1000], "ms")
	fmt.Println("第500张票时，投票所耗时：", timeMap[500], "ms")
	fmt.Println("第100张票时，投票所耗时：", timeMap[100], "ms")
	fmt.Println("第50张票时，投票所耗时：", timeMap[50], "ms")
	fmt.Println("第10张票时，投票所耗时：", timeMap[10], "ms")
	fmt.Println("第1张票时，投票所耗时：", timeMap[1], "ms")
}

func initParam() (*state.StateDB, *pposm.CandidatePool, *pposm.TicketPool) {
	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, ethash.NewFaker(), vm.Config{}, nil)
	ticketcache.NewTicketIdsCache(db)
	configs := params.PposConfig{
		//MaxChair: 1,
		//MaxCount: 3,
		//RefundBlockNumber: 	1,
		Candidate: &params.CandidateConfig{
			MaxChair: 1,
			MaxCount: 3,
			RefundBlockNumber: 	1,
		},
		TicketConfig: &params.TicketConfig{
			MaxCount: 			10000,
			ExpireBlockNumber:	4,
		},
	}

	candidatePool := pposm.NewCandidatePool(&configs)

	ticketPool := pposm.NewTicketPool(&configs)

	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		log.Error("reference statedb failed", "err", err)
	} else {
		state = statedb
	}
	return state, candidatePool, ticketPool
}

func TestTicketPool_VoteTicket(t *testing.T) {
	state, candidatePool, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),
	}

	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		log.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(10))
	if nil != err {
		log.Error("vote ticket fail", "err", err)
	}
}

func TestTicketPool_GetExpireTicketIds(t *testing.T) {
	state, _, ticketPool := initParam()

	_, err := ticketPool.GetExpireTicketIds(state, new(big.Int).SetUint64(10))
	if nil != err {
		log.Error("getExpireTicketIds fail", "err", err)
	}
}

func TestTicketPool_GetTicketList(t *testing.T) {
	state, candidatePool, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),
	}

	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		log.Error("SetCandidate fail", "err", err)
	}

	ticketIdList, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(10))
	if nil != err {
		log.Error("vote ticket fail", "err", err)
	}

	_, err = ticketPool.GetTicketList(state, ticketIdList)
	if nil != err {
		log.Error("getTicketList fail", "err", err)
	}
}

func TestTicketPool_GetTicket(t *testing.T) {
	state, candidatePool, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),
	}

	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		log.Error("SetCandidate fail", "err", err)
	}

	ticketIdList, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(10))
	if nil != err {
		log.Error("vote ticket fail", "err", err)
	}

	_, err = ticketPool.GetTicket(state, ticketIdList[0])
	if nil != err {
		log.Error("getTicket fail", "err", err)
	}
}

func TestTicketPool_DropReturnTicket(t *testing.T) {
	state, candidatePool, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),
	}

	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		log.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(6))
	if nil != err {
		log.Error("vote ticket fail", "err", err)
	}

	err = ticketPool.DropReturnTicket(state, candidate.CandidateId)
	if nil != err {
		log.Error("dropReturnTicket fail", "err", err)
	}
}

func TestTicketPool_ReturnTicket(t *testing.T) {
	state, candidatePool, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),
	}

	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		log.Error("SetCandidate fail", "err", err)
	}

	ticketIdList, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(10))
	if nil != err {
		log.Error("vote ticket fail", "err", err)
	}

	err = ticketPool.ReturnTicket(state, candidate.CandidateId, ticketIdList[0], new(big.Int).SetUint64(10))
	if nil != err {
		log.Error("returnTicket fail", "err", err)
	}
}

func TestTicketPool_Notify(t *testing.T) {
	state, candidatePool, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),
	}

	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		log.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(6))
	if nil != err {
		log.Error("vote ticket fail", "err", err)
	}

	err = ticketPool.Notify(state, new(big.Int).SetUint64(10))
	if nil != err {
		log.Error("notify fail", "err", err)
	}
}

func TestTicketPool_SelectionLuckyTicket(t *testing.T) {
	state, candidatePool, ticketPool := initParam()

	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),
	}

	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		log.Error("SetCandidate fail", "err", err)
	}

	_, err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), 10, new(big.Int).SetUint64(100), candidate.CandidateId, new(big.Int).SetUint64(6))
	if nil != err {
		log.Error("vote ticket fail", "err", err)
	}

	blockHash := common.Hash{}
	blockHash.SetBytes([]byte("3b41e0aee38c1a1f959a6aaae678d86f1e6af59617d2f667bb2ef5527779c861"))
	_, err = ticketPool.SelectionLuckyTicket(state, candidate.CandidateId, blockHash)
	if nil != err {
		log.Error("selectionLuckyTicket fail", "err", err)
	}
}

func TestTicketPool_GetPoolNumber(t *testing.T) {
	state, _, ticketPool := initParam()

	_, err := ticketPool.GetPoolNumber(state)
	if nil != err {
		log.Error("getPoolNumber fail", "err", err)
	}
}

func TestTicketPool_GetCandidateTicketIds(t *testing.T) {
	state, _, ticketPool := initParam()

	_, err := ticketPool.GetCandidateTicketIds(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	if nil != err {
		log.Error("getCandidateTicketIds fail", "err", err)
	}
}

func TestTicketPool_GetCandidateAttach(t *testing.T) {
	state, _, ticketPool := initParam()

	_, err := ticketPool.GetCandidateAttach(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	if nil != err {
		log.Error("getCandidateAttach fail", "err", err)
	}
}

func TestTicketPool_GetCandidateEpoch(t *testing.T) {
	state, _, ticketPool := initParam()

	_, err := ticketPool.GetCandidateEpoch(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	if nil != err {
		log.Error("getCandidateEpoch fail", "err", err)
	}
}

func TestTicketPool_GetTicketPrice(t *testing.T) {
	state, _, ticketPool := initParam()

	_, err := ticketPool.GetTicketPrice(state)
	if nil != err {
		log.Error("getTicketPrice fail", "err", err)
	}
}