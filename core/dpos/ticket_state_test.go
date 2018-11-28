package depos_test

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/consensus/ethash"
	"Platon-go/core"
	"Platon-go/core/dpos"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/core/vm"
	"Platon-go/ethdb"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestVoteTicket(t *testing.T)  {
	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, ethash.NewFaker(), vm.Config{}, nil)

	configs := params.DposConfig{
		//MaxChair: 1,
		//MaxCount: 3,
		//RefundBlockNumber: 	1,
		Candidate: &params.CandidateConfig{
			MaxChair: 1,
			MaxCount: 3,
			RefundBlockNumber: 	1,
		},
		TicketConfig: &params.TicketConfig{
			MaxCount: 			10,
			ExpireBlockNumber:	1000,
		},
	}

	candidatePool := depos.NewCandidatePool(&configs)

	ticketPool := depos.NewTicketPool(&configs, candidatePool)

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

	for i := 0; i < 50 ; i++ {
		go func() {
			err := ticketPool.VoteTicket(state, common.HexToAddress("0x20"), new(big.Int).SetUint64(10), candidate.CandidateId, new(big.Int).SetUint64(10))
			if nil != err {
				fmt.Println("vote ticket error:", err)
			}
		}()
	}
	time.Sleep(time.Second * 5)
	candidate, err := candidatePool.GetCandidate(state, candidate.CandidateId)
	if err != nil {
		fmt.Println("GetCandidate error")
		return
	}

	ticketList, err := ticketPool.GetTicketList(state, candidate.TicketPool)

	expireTicketList, err := ticketPool.GetExpireTicket(state, new(big.Int).SetUint64(10))

	//fmt.Printf("print info:\n\t%+v\n\t%+v\n\t%+v\n\t%+v,%v", candidatePool, ticketPool, candidate, ticketList, err)
	fmt.Printf("ticketPoolSize:[%d],expireTicketListSize:[%d],candidate.TicketPool:[%d],tcount:[%d],epoch:[%d]\n",
		ticketPool.SurplusQuantity, len(expireTicketList), len(candidate.TicketPool), candidate.TCount, candidate.Epoch)
	fmt.Println("------all ticket-----")
	for _, ticket := range ticketList {
		fmt.Printf("ticket:%+v,ticketId:[%v]\n", ticket, hexutil.Encode(ticket.TicketId.Bytes()))
	}
	ticketId := ticketList[rand.Intn(len(ticketList))].TicketId
	fmt.Printf("-----------开始释放一张选票【%v】-----------\n", hexutil.Encode(ticketId.Bytes()))
	ticketPool.ReleaseTicket(state, candidate.CandidateId, ticketId, new(big.Int).SetUint64(10))
	candidate, err = candidatePool.GetCandidate(state, candidate.CandidateId)
	fmt.Printf("ticketPoolSize:[%d],expireTicketListSize:[%d],candidate.TicketPool:[%d],tcount:[%d],epoch:[%d]\n",
		ticketPool.SurplusQuantity, len(expireTicketList), len(candidate.TicketPool), candidate.TCount, candidate.Epoch)
}