package depos

import (
	"testing"
	"path/filepath"
	"fmt"
	"Platon-go/params"
	"math/big"
	"Platon-go/common"
	"Platon-go/ethdb"
	"Platon-go/core"
	"Platon-go/core/vm"
	"Platon-go/consensus/ethash"
	"Platon-go/p2p/discover"
)

func TestInitCandidatePoolByConfig (t *testing.T){
	fmt.Println(filepath.Join("E:/platon-data", "dpos.json"))
	//fmt.Println(common.BytesToHash([]byte("1")))
	//fmt.Println(string(common.HexToHash("1").Bytes()), common.HexToHash("1").String(), string(common.BytesToHash([]byte("1")).Bytes()))

	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, ethash.NewFaker(), vm.Config{}, nil)
	//if nil != err {
	//	fmt.Println("init blockchain err", err)
	//	return
	//}

	can_Configs := []*params.CandidateConfig {
		&params.CandidateConfig{
			Deposit:			0,
			BlockNumber: 	 	new(big.Int).SetUint64(10),
			TxIndex: 		 	2,
			CandidateId:		discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"),
			Host: 			 	"10.0.0.1",
			Port: 			 	"8545",
		},
		&params.CandidateConfig{
			Deposit:			0,
			BlockNumber: 	 	new(big.Int).SetUint64(10),
			TxIndex: 		 	2,
			CandidateId:		discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012342"),
			Host: 			 	"10.0.0.1",
			Port: 			 	"8546",
		},
		&params.CandidateConfig{
			Deposit:			0,
			BlockNumber: 	 	new(big.Int).SetUint64(10),
			TxIndex: 		 	2,
			CandidateId:		discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012343"),
			Host: 			 	"10.0.0.1",
			Port: 			 	"8547",
		},
	}
	//a := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341")
	//b := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012343")
	//c := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012344")
	//fmt.Println("hex byte == \n", a.Bytes(), "\n", b.Bytes(), "\n", c.Bytes())


	configs := params.DposConfig{
		MaxChair: 1,
		MaxCount: 3,
		Candidates: can_Configs,
	}
	var candidatePool *CandidatePool
	//if statedb, err := blockchain.State(); nil != err {
	//	fmt.Println("reference statedb failed", err)
	//}else {
	//
	//	state = statedb
	//}
	/** test init candidatePool */
	if pool, err := NewCandidatePool(blockchain, &configs); nil != err {
		fmt.Println("init candidatePool err", err)
	}else{
		candidatePool = pool
	}
	// commit trie to db
	candidatePool.CommitTrie(false)
	//tr := state.StorageTrie(CandidateAddr)
	candidatePool.IteratorTrie("初始化：")
	printObject("初始化DB中的 " + WitnessPrefix, candidatePool.buildCandidatesByTrie(WitnessPrefix))
	printObject("初始化DB中的 " + ImmediatePrefix, candidatePool.buildCandidatesByTrie(ImmediatePrefix))
	printObject("初始化DB中的 " + DefeatPrefix, candidatePool.buildCandidateArrByTrie(DefeatPrefix))



	candidate := &Candidate{
		Deposit: 		100,
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}

	fmt.Println("设置新的k-v \n")
	/** test SetCandidate */
	candidatePool.SetCandidate(candidate.CandidateId, candidate)

	// commit trie to db
	candidatePool.CommitTrie(false)
	fmt.Println("提交.......")
	printObject("更新后DB中的 " + WitnessPrefix, candidatePool.buildCandidatesByTrie(WitnessPrefix))
	printObject("更新后DB中的 " + ImmediatePrefix, candidatePool.buildCandidatesByTrie(ImmediatePrefix))
	printObject("更新后DB中的 " + DefeatPrefix, candidatePool.buildCandidateArrByTrie(DefeatPrefix))

	/** test GetCandidate */
	can, _ := candidatePool.GetCandidate(discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"))
	fmt.Println("GetCandidate", can)


	/** test WithdrawCandidate */
	ok1 := candidatePool.WithdrawCandidate(discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), 99)
	fmt.Println("error", ok1)
	candidatePool.CommitTrie(false)
	printObject("退款后DB中的 " + WitnessPrefix, candidatePool.buildCandidatesByTrie(WitnessPrefix))
	printObject("退款后DB中的 " + ImmediatePrefix, candidatePool.buildCandidatesByTrie(ImmediatePrefix))
	printObject("退款后DB中的 " + DefeatPrefix, candidatePool.buildCandidateArrByTrie(DefeatPrefix))

	/** test WithdrawCandidate again */
	ok2 := candidatePool.WithdrawCandidate(discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), 12)
	fmt.Println("error", ok2)
	candidatePool.CommitTrie(false)
	printObject("退款后DB中的 " + WitnessPrefix, candidatePool.buildCandidatesByTrie(WitnessPrefix))
	printObject("退款后DB中的 " + ImmediatePrefix, candidatePool.buildCandidatesByTrie(ImmediatePrefix))
	printObject("退款后DB中的 " + DefeatPrefix, candidatePool.buildCandidateArrByTrie(DefeatPrefix))

	/** test GetChosens */
	canArr := candidatePool.GetChosens()
	printObject("入围候选人", canArr)

	/** test GetChairpersons */
	canArr = candidatePool.GetChairpersons()
	printObject("见证人", canArr)

	/** test GetDefeat */


	/** test IsDefeat */

	/** test Election */


	/** test RefundBalance */


	/** test GetOwner */
}

func TestInitCandidatePoolByTrie (t *testing.T){

}