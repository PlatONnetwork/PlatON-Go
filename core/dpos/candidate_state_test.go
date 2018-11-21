package depos_test

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
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/consensus/ethash"
	"Platon-go/p2p/discover"

	"Platon-go/core/dpos"
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
			Deposit:			new(big.Int).SetUint64(0),
			BlockNumber: 	 	new(big.Int).SetUint64(9),
			TxIndex: 		 	2,
			CandidateId:		discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"),
			Host: 			 	"10.0.0.1",
			Port: 			 	"8545",
		},
		&params.CandidateConfig{
			Deposit:			new(big.Int).SetUint64(0),
			BlockNumber: 	 	new(big.Int).SetUint64(10),
			TxIndex: 		 	2,
			CandidateId:		discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012342"),
			Host: 			 	"10.0.0.1",
			Port: 			 	"8546",
		},
		&params.CandidateConfig{
			Deposit:			new(big.Int).SetUint64(0),
			BlockNumber: 	 	new(big.Int).SetUint64(10),
			TxIndex: 		 	3,
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
		RefundBlockNumber: 	1,
		Candidates: can_Configs,
	}
	var candidatePool *depos.CandidatePool
	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		fmt.Println("reference statedb failed", err)
	}else {
		var isgenesis bool
		if blockchain.CurrentBlock().NumberU64() == blockchain.Genesis().NumberU64() {
			isgenesis = true
		}
		/** test init candidatePool */
		if pool, err := depos.NewCandidatePool(statedb, &configs, isgenesis); nil != err {
			fmt.Println("init candidatePool err", err)
		}else{
			candidatePool = pool
		}
		state = statedb
	}

	// commit trie to db
	state.Commit(false)
	//tr := state.StorageTrie(CandidateAddr)
	//candidatePool.IteratorTrie("初始化：")
	//printObject("初始化DB中的 " + WitnessPrefix, candidatePool.buildCandidatesByTrie(WitnessPrefix))
	//printObject("初始化DB中的 " + ImmediatePrefix, candidatePool.buildCandidatesByTrie(ImmediatePrefix))
	//printObject("初始化DB中的 " + DefeatPrefix, candidatePool.buildCandidateArrByTrie(DefeatPrefix))



	candidate := &types.Candidate{
		Deposit: 		new(big.Int).SetUint64(100),
		BlockNumber:    new(big.Int).SetUint64(7),
		CandidateId:   discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		TxIndex:  		6,
		Host:  			"10.0.0.1",
		Port:  			"8548",
		Owner: 			common.HexToAddress("0x12"),

	}

	fmt.Println("设置新的k-v \n")
	/** test SetCandidate */
	if err := candidatePool.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		fmt.Println("SetCandidate err:", err)
	}

	// commit trie to db
	//candidatePool.CommitTrie(false)
	//fmt.Println("提交.......")
	//printObject("更新后DB中的 " + WitnessPrefix, candidatePool.buildCandidatesByTrie(WitnessPrefix))
	//printObject("更新后DB中的 " + ImmediatePrefix, candidatePool.buildCandidatesByTrie(ImmediatePrefix))
	//printObject("更新后DB中的 " + DefeatPrefix, candidatePool.buildCandidateArrByTrie(DefeatPrefix))

	/** test GetCandidate */
	fmt.Println("test GetCandidate")
	can, _ := candidatePool.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012341"))
	fmt.Println("GetCandidate", can)


	/** test WithdrawCandidate */
	fmt.Println("test WithdrawCandidate")
	ok1 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	fmt.Println("error", ok1)
	//candidatePool.CommitTrie(false)
	//printObject("退款后DB中的 " + WitnessPrefix, candidatePool.buildCandidatesByTrie(WitnessPrefix))
	//printObject("退款后DB中的 " + ImmediatePrefix, candidatePool.buildCandidatesByTrie(ImmediatePrefix))
	//printObject("退款后DB中的 " + DefeatPrefix, candidatePool.buildCandidateArrByTrie(DefeatPrefix))

	/** test WithdrawCandidate again */
	fmt.Println("test WithdrawCandidate again")
	ok2 := candidatePool.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(10)), new(big.Int).SetUint64(uint64(11)))
	fmt.Println("error", ok2)
	//candidatePool.CommitTrie(false)
	//printObject("退款后DB中的 " + WitnessPrefix, candidatePool.buildCandidatesByTrie(WitnessPrefix))
	//printObject("退款后DB中的 " + ImmediatePrefix, candidatePool.buildCandidatesByTrie(ImmediatePrefix))
	//printObject("退款后DB中的 " + DefeatPrefix, candidatePool.buildCandidateArrByTrie(DefeatPrefix))

	/** test GetChosens */
	fmt.Println("test GetChosens")
	canArr := candidatePool.GetChosens(state)
	depos.PrintObject("入围候选人", canArr)

	/** test GetChairpersons */
	fmt.Println("test GetChairpersons")
	canArr = candidatePool.GetChairpersons(state)
	depos.PrintObject("见证人", canArr)

	/** test GetDefeat */
	fmt.Println("test GetDefeat")
	defeatArr, _ := candidatePool.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	depos.PrintObject("可以退款信息", defeatArr)

	/** test IsDefeat */
	fmt.Println("test IsDefeat")
	flag, _ := candidatePool.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	depos.PrintObject("是否落榜", flag)

	/** test Election */
	fmt.Println("test Election")
	flag = candidatePool.Election(state)
	fmt.Println("是否揭榜成功", flag)

	/** test RefundBalance */
	fmt.Println("test RefundBalance")
	err := candidatePool.RefundBalance(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(10)))
	fmt.Println("err", err)

	/** test GetOwner */
	fmt.Println("test GetOwner")
	addr := candidatePool.GetOwner(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	fmt.Println("收益地址", addr.String())

	/**  test GetWitness */
	fmt.Println("test GetWitness")
	nodeArr := candidatePool.GetWitness(state)
	fmt.Printf("nodeArr := %+v", nodeArr)
}

func TestInitCandidatePoolByTrie (t *testing.T){

}