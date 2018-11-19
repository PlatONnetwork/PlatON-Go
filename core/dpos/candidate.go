package depos

import (
	"math/big"
	"Platon-go/common"
	"Platon-go/p2p/discover"
)

// 候选人
type Candidate struct {

	// 抵押金额(保证金)数目
	Deposit			uint64
	// 发生抵押时的当前块高
	BlockNumber 	*big.Int
	// 发生抵押时的tx index
	TxIndex 		uint32
	// 候选人Id
	CandidateId 	discover.NodeID
	//CandidateId 	string			`json:"candidateid"`
	//
	Host 			string
	Port 			string

	// 质押收益账户
	Owner 			common.Address
	// 发起质押的账户
	From 			common.Address


	// 被投的票Id集
	//ticketPool		[]common.Hash
	// 被投票数目
	//TCount    		uint64				`json:"tcount"`
	// 票龄
	//Epoch			*big.Int			`json:"epoch"`
	// 佣金
	//Brokerage		uint64				`json:"brokerage"`
}

func newCandidate (){

}
