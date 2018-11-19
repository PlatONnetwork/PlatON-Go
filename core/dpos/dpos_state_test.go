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
	"Platon-go/core/state"
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

	//cfg := &params.ChainConfig{}
	//cfg.Cbft.Dposes = configs
	var candidatePool *CandidatePool
	currentBlock := blockchain.CurrentBlock()
	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		fmt.Println("reference statedb failed", err)
	}else {
		var isgenesis bool
		//fmt.Printf("genesis: %+v, current: %+v \n", blockchain.Genesis(), currentBlock)
		if blockchain.Genesis().NumberU64() == currentBlock.NumberU64() {
			isgenesis = true
			//fmt.Println("YES")
		}
		/** test init candidatePool */
		if pool, err := NewCandidatePool(statedb, &configs, isgenesis); nil != err {
			fmt.Println("init candidatePool err", err)
		}else{
			candidatePool = pool
		}
		state = statedb

	}
	// commit trie to db
	state.Commit(false)
	tr := state.StorageTrie(CandidateAddr)
	iteratorTrie("初始化：", tr)
	printObject("初始化DB中的 " + WitnessPrefix, buildCandidatesByTrie(tr, WitnessPrefix))
	printObject("初始化DB中的 " + ImmediatePrefix, buildCandidatesByTrie(tr, ImmediatePrefix))
	printObject("初始化DB中的 " + DefeatPrefix, buildCandidateArrByTrie(tr, DefeatPrefix))



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
	state.Commit(false)
	tr = state.StorageTrie(CandidateAddr)
	fmt.Println("提交.......")
	printObject("更新后DB中的 " + WitnessPrefix, buildCandidatesByTrie(tr, WitnessPrefix))
	printObject("更新后DB中的 " + ImmediatePrefix, buildCandidatesByTrie(tr, ImmediatePrefix))
	printObject("更新后DB中的 " + DefeatPrefix, buildCandidateArrByTrie(tr, DefeatPrefix))

	/** test GetCandidate */
	can := candidatePool.GetCandidate(discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"))
	fmt.Println(can)


	/** test WithdrawCandidate */
	ok1 := candidatePool.WithdrawCandidate(discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), 99)
	fmt.Println(ok1)
	state.Commit(false)
	tr = state.StorageTrie(CandidateAddr)
	printObject("退款后DB中的 " + WitnessPrefix, buildCandidatesByTrie(tr, WitnessPrefix))
	printObject("退款后DB中的 " + ImmediatePrefix, buildCandidatesByTrie(tr, ImmediatePrefix))
	printObject("退款后DB中的 " + DefeatPrefix, buildCandidateArrByTrie(tr, DefeatPrefix))

	/** test WithdrawCandidate again */
	ok2 := candidatePool.WithdrawCandidate(discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), 1)
	fmt.Println(ok2)
	state.Commit(false)
	tr = state.StorageTrie(CandidateAddr)
	printObject("退款后DB中的 " + WitnessPrefix, buildCandidatesByTrie(tr, WitnessPrefix))
	printObject("退款后DB中的 " + ImmediatePrefix, buildCandidatesByTrie(tr, ImmediatePrefix))
	printObject("退款后DB中的 " + DefeatPrefix, buildCandidateArrByTrie(tr, DefeatPrefix))

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

}

func TestInitCandidatePoolByTrie (t *testing.T){

}