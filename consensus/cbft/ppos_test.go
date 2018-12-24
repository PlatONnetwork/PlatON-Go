package cbft

import (
	"Platon-go/core/types"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"math/big"
	"testing"
	"Platon-go/core/ppos"
	"Platon-go/ethdb"
	"Platon-go/core"
	"fmt"
	"Platon-go/consensus/ethash"
	"Platon-go/core/vm"
	"encoding/json"
	"net"
	"strconv"
	"Platon-go/core/state"
	"sync/atomic"
	"Platon-go/common"
	"math/rand"
	"time"
)

func newTesterAccountPool() ([]discover.NodeID, error) {
	var accounts []discover.NodeID
	for _, url := range params.MainnetBootnodes {
		node, err := discover.ParseNode(url)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, node.ID)
	}
	return accounts, nil
}

func TestBlockCopy(t *testing.T) {
	header := &types.Header{
		Number: big.NewInt(1),
	}
	block1 := types.NewBlock(header, nil, nil, nil)

	block2Obj := *block1

	block2Obj.Header().Number = big.NewInt(2)

	block2 := &block2Obj
	block2.Header().Number = big.NewInt(3)

	println(block1.Number().Uint64())
	println(block2.Number().Uint64())

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


func buildPpos() (*ppos, *core.BlockChain) {
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
	ppos := &ppos{
		candidatePool: pposm.NewCandidatePool(&configs),
		ticketPool: pposm.NewTicketPool(&configs),
	}

	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, ethash.NewFaker(), vm.Config{}, nil)

	ppos.setCandidatePool(blockchain, buildInitialNodes())

	return ppos, blockchain
}

func buildInitialNodes() []discover.Node {

	nodes := make([]discover.Node, 0)
	for i := 1; i <= 3; i++ {
		ip := net.ParseIP("127.0.0.1")
		// uint16
		var port uint16
		if portInt, err := strconv.Atoi("854" + fmt.Sprint(i)); nil != err {
			return nil
		} else {
			port = uint16(portInt)
		}
		nodeId := discover.MustHexID("0x0123456789012134567890112345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234" + fmt.Sprint(i))
		nodes = append(nodes, *(discover.NewNode(nodeId, ip, port, port)))
	}
	return nodes
}

func TestNewPpos (t *testing.T) {
	ppos, _ := buildPpos()
	printObject("ppos.candidatePool:", ppos.candidatePool, t)
	printObject("ppos.ticketPool:", ppos.ticketPool, t)
}

// test BlockProducerIndex
func ppos_BlockProducerIndex(logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	curr := bc.CurrentBlock()
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012342")
	num := ppos.BlockProducerIndex(curr.Number(), curr.Hash(), curr.Number(), nodeId, 1)
	logFn("BlockProducerIndex：", num)
}
func TestPpos_BlockProducerIndex(t *testing.T) {
	ppos_BlockProducerIndex(t, t.Log, t.Error)
}
func BenchmarkPpos_BlockProducerIndex(b *testing.B) {
	ppos_BlockProducerIndex(b, b.Log, b.Error)
}

/** about candidatepool */
// test SetCandidate
func ppos_SetCandidate (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test SetCandidate getting state err", err)
	}else {
		state = st
	}
	logFn("test SetCandidate ...")



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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}else {
		logFn("SetCandidate success ... ")
	}
}
func TestPpos_SetCandidate(t *testing.T) {
	ppos_SetCandidate(t, t.Log, t.Error)
}
func BenchmarkPpos_SetCandidate(b *testing.B) {
	ppos_SetCandidate(b, b.Log, b.Error)
}

// test GetCandidate
func ppos_GetCandidate (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetCandidate getting state err", err)
	}else {
		state = st
	}
	logFn("test GetCandidate ...")



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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}else {
		logFn("SetCandidate success ... ")
	}

	/** test GetCandidate ... */
	if can, err := ppos.GetCandidate(state, candidate.CandidateId); nil != err {
		errFn("GetCandidate err", err)
	}else {
		printObject("GetCandidate can:", can, logger)
	}
}
func TestPpos_GetCandidate(t *testing.T) {
	ppos_GetCandidate(t, t.Log, t.Error)
}
func BenchmarkPpos_GetCandidate(b *testing.B) {
	ppos_GetCandidate(b, b.Log, b.Error)
}

// test Election
func ppos_Election (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test Election getting state err", err)
	}else {
		state = st
	}
	logFn("test Election ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
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
				ppos.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := ppos.Election(state, big.NewInt(20))
	logFn("Whether election was successful err", err)
}
func TestPpos_Election(t *testing.T) {
	ppos_Election(t, t.Log, t.Error)
}
func BenchmarkPpos_Election(b *testing.B) {
	ppos_Election(b, b.Log, b.Error)
}

// test Switch
func ppos_Switch (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test Switch getting state err", err)
	}else {
		state = st
	}
	logFn("test Switch ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
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
				ppos.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := ppos.Election(state, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ := ppos.GetWitness(state, 1)
	printObject("next Witnesses", canArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := ppos.Switch(state)
	logFn("Switch was success ", flag)
}
func TestPpos_Switch(t *testing.T) {
	ppos_Switch(t, t.Log, t.Error)
}
func BenchmarkPpos_Switch(b *testing.B) {
	ppos_Switch(b, b.Log, b.Error)
}

// test GetWitness
func ppos_GetWitness (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetWitness getting state err", err)
	}else {
		state = st
	}
	logFn("test GetWitness ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
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
				ppos.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := ppos.Election(state, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ := ppos.GetWitness(state, 1)
	printObject("next Witnesses", canArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := ppos.Switch(state)
	logFn("Switch was success ", flag)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ = ppos.GetWitness(state, 0)
	printObject(" current Witnesses", canArr, logger)
}
func TestPpos_GetWitness(t *testing.T) {
	ppos_GetWitness(t, t.Log, t.Error)
}
func BenchmarkPpos_GetWitness(b *testing.B) {
	ppos_GetWitness(b, b.Log, b.Error)
}

// test GetAllWitness
func ppos_GetAllWitness (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetAllWitness getting state err", err)
	}else {
		state = st
	}
	logFn("test GetAllWitness ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
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
				ppos.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := ppos.Election(state, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ := ppos.GetWitness(state, 1)
	printObject("next Witnesses", canArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := ppos.Switch(state)
	logFn("Switch was success ", flag)

	/** test Election again */
	logFn("test Election again ...")
	_, err = ppos.Election(state, big.NewInt(20))
	logFn("Whether election again was successful err", err)

	/** test GetAllWitness */
	logFn("test GetAllWitness ...")
	preArr, canArr, nextArr, _ := ppos.GetAllWitness(state)
	printObject(" previous Witnesses", preArr, logger)
	printObject(" current Witnesses", canArr, logger)
	printObject(" next Witnesses", nextArr, logger)
}
func TestPpos_GetAllWitness(t *testing.T) {
	ppos_GetAllWitness(t, t.Log, t.Error)
}
func BenchmarkPpos_GetAllWitness(b *testing.B) {
	ppos_GetAllWitness(b, b.Log, b.Error)
}

// test WithdrawCandidate
func ppos_WithdrawCandidate (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test WithdrawCandidate getting state err", err)
	}else {
		state = st
	}
	logFn("test WithdrawCandidate ...")


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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can, _ := ppos.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can, logger)

	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := ppos.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)

	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can2, _ := ppos.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("GetCandidate", can2, logger)

}
func TestPpos_WithdrawCandidate(t *testing.T) {
	ppos_WithdrawCandidate(t, t.Log, t.Error)
}
func BenchmarkPpos_WithdrawCandidate(b *testing.B) {
	ppos_WithdrawCandidate(b, b.Log, b.Error)
}

// test GetChosens
func ppos_GetChosens (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetChosens getting state err", err)
	}else {
		state = st
	}
	logFn("test GetChosens ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
		errFn("SetCandidate err:", err)
	}

	/** test GetChosens */
	logFn("test GetChosens ...")
	canArr := ppos.GetChosens(state, 0)
	printObject("elected candidates", canArr, logger)
}
func TestPpos_GetChosens(t *testing.T) {
	ppos_GetChosens(t, t.Log, t.Error)
}
func BenchmarkPpos_GetChosens(b *testing.B) {
	ppos_GetChosens(b, b.Log, b.Error)
}

// test GetChairpersons
func ppos_GetChairpersons (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetChairpersons getting state err", err)
	}else {
		state = st
	}
	logFn("test GetChairpersons ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
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
				ppos.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")

	/** test Election */
	logFn("test Election ...")
	_, err := ppos.Election(state, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	nodeIdArr, _ := ppos.GetWitness(state, 1)
	printObject("next Witnesses", nodeIdArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := ppos.Switch(state)
	logFn("Switch was success ", flag)

	/** test GetChairpersons */
	logFn("test GetChairpersons ...")
	canArr := ppos.GetChairpersons(state)
	printObject("GetChairpersons canArr:", canArr, logger)
}
func TestPpos_GetChairpersons(t *testing.T) {
	ppos_GetChairpersons(t, t.Log, t.Error)
}
func BenchmarkPpos_GetChairpersons(b *testing.B) {
	ppos_GetChairpersons(b, b.Log, b.Error)
}

// test GetDefeat
func ppos_GetDefeat (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetDefeat getting state err", err)
	}else {
		state = st
	}
	logFn("test GetDefeat ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
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
				ppos.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")



	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := ppos.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)


	/** test IsDefeat */
	logFn("test IsDefeat ...")
	flag, _ := ppos.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	logFn("isdefeat", flag)

	/** test GetDefeat */
	logFn("test GetDefeat ...")
	defeatArr, _ := ppos.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	printObject("can be refund defeats", defeatArr, logger)
}
func TestPpos_GetDefeat(t *testing.T) {
	ppos_GetDefeat(t, t.Log, t.Error)
}
func BenchmarkPpos_GetDefeat(b *testing.B) {
	ppos_GetDefeat(b, b.Log, b.Error)
}

// test IsDefeat
func ppos_IsDefeat (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test IsDefeat getting state err", err)
	}else {
		state = st
	}
	logFn("test IsDefeat ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
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
				ppos.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")



	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := ppos.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)


	/** test IsDefeat */
	logFn("test IsDefeat ...")
	flag, _ := ppos.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"))
	logFn("isdefeat", flag)
}
func TestPpos_IsDefeat(t *testing.T) {
	ppos_IsDefeat(t, t.Log, t.Error)
}
func BenchmarkPpos_IsDefeat(b *testing.B) {
	ppos_IsDefeat(b, b.Log, b.Error)
}

// test RefundBalance
func ppos_RefundBalance (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test RefundBalance getting state err", err)
	}else {
		state = st
	}
	logFn("test RefundBalance ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
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
	if err := ppos.SetCandidate(state, candidate2.CandidateId, candidate2); nil != err {
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
	if err := ppos.SetCandidate(state, candidate3.CandidateId, candidate3); nil != err {
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
	if err := ppos.SetCandidate(state, candidate4.CandidateId, candidate4); nil != err {
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
				ppos.Notify(state, new(big.Int).SetUint64(tempBlockNumber))
				tempBlockNumber++
			}
			fmt.Println("release ticket,end ############################################################")
		}
		fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", can.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
		_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, can.CandidateId, tempBlockNumber)
		if nil != err {
			fmt.Println("vote ticket error:", err)
		}
		atomic.AddUint32(&count, 1)
		timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime

	}
	fmt.Println("VOTING END .............................................................")



	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := ppos.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)

	/** test  RefundBalance*/
	if err := ppos.RefundBalance(state, candidate.CandidateId, big.NewInt(11)); nil != err {
		errFn("RefundBalance err", err)
	}else {
		logFn("RefundBalance success ...")
	}
}
func TestPpos_RefundBalance(t *testing.T) {
	ppos_RefundBalance(t, t.Log, t.Error)
}
func BenchmarkPpos_RefundBalance(b *testing.B) {
	ppos_RefundBalance(b, b.Log, b.Error)
}

// test GetOwner
func ppos_GetOwner (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetOwner getting state err", err)
	}else {
		state = st
	}
	logFn("test GetOwner ...")

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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}

	/** test GetOwner */
	ownerAddr := ppos.GetOwner(state, candidate.CandidateId)
	logFn("Getting Onwer's Address:", ownerAddr.String())
}
func TestPpos_GetOwner(t *testing.T) {
	ppos_GetOwner(t, t.Log, t.Error)
}
func BenchmarkPpos_GetOwner(b *testing.B) {
	ppos_GetOwner(b, b.Log, b.Error)
}

// test GetRefundInterval
func ppos_GetRefundInterval (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, _ := buildPpos()

	logFn("test GetRefundInterval ...")

	/** test  GetRefundInterval*/
	num := ppos.GetRefundInterval()
	logFn("RefundInterval:", num)
}
func TestPpos_GetRefundInterval(t *testing.T) {
	ppos_GetRefundInterval(t, t.Log, t.Error)
}
func BenchmarkPpos_GetRefundInterval(b *testing.B) {
	ppos_GetRefundInterval(b, b.Log, b.Error)
}

/** about tickpool */
// test GetPoolNumber
func ppos_GetPoolNumber (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetPoolNumber getting state err", err)
	}else {
		state = st
	}
	logFn("test GetPoolNumber ...")

	/** test  GetPoolNumber */
	logFn("test GetPoolNumber ...")
	if num, err := ppos.GetPoolNumber(state); nil != err {
		errFn("GetPoolNumber err", err)
	}else {
		logFn("GetPoolNumber:", num)
	}
}
func TestPpos_GetPoolNumber(t *testing.T) {
	ppos_GetPoolNumber(t, t.Log, t.Error)
}
func BenchmarkPpos_GetPoolNumber(b *testing.B) {
	ppos_GetPoolNumber(b, b.Log, b.Error)
}

// test VoteTicket
func ppos_VoteTicket (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test VoteTicket getting state err", err)
	}else {
		state = st
	}
	logFn("test VoteTicket ...")


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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)

	timeMap := make(map[uint32]int64)
	logFn("VOTING START .............................................................")

	startTime := time.Now().UnixNano() / 1e6
	voteOwner := common.HexToAddress("0x20")
	deposit := new(big.Int).SetUint64(10)
	state.SubBalance(voteOwner, deposit)
	state.AddBalance(common.TicketPoolAddr, deposit)
	tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
	fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", candidate.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
	_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")
}
func TestPpos_VoteTicket(t *testing.T) {
	ppos_VoteTicket(t, t.Log, t.Error)
}
func BenchmarkPpos_VoteTicket(b *testing.B) {
	ppos_VoteTicket(b, b.Log, b.Error)
}

// test GetTicket
func ppos_GetTicket (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetTicket getting state err", err)
	}else {
		state = st
	}
	logFn("test GetTicket ...")


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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)

	timeMap := make(map[uint32]int64)
	logFn("VOTING START .............................................................")

	startTime := time.Now().UnixNano() / 1e6
	voteOwner := common.HexToAddress("0x20")
	deposit := new(big.Int).SetUint64(10)
	state.SubBalance(voteOwner, deposit)
	state.AddBalance(common.TicketPoolAddr, deposit)
	tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
	fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", candidate.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
	tickList, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	/** GetTicket */
	if ticket, err := ppos.GetTicket(state, tickList[0]); nil != err {
		errFn("GetTicket err", err)
	}else {
		printObject("GetTicket ticketInfo:", ticket, logger)
	}
}
func TestPpos_GetTicket(t *testing.T) {
	ppos_GetTicket(t, t.Log, t.Error)
}
func BenchmarkPpos_GetTicket(b *testing.B) {
	ppos_GetTicket(b, b.Log, b.Error)
}

// test GetTicketList
func ppos_GetTicketList (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetTicketList getting state err", err)
	}else {
		state = st
	}
	logFn("test GetTicketList ...")


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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)

	timeMap := make(map[uint32]int64)
	logFn("VOTING START .............................................................")

	startTime := time.Now().UnixNano() / 1e6
	voteOwner := common.HexToAddress("0x20")
	deposit := new(big.Int).SetUint64(10)
	state.SubBalance(voteOwner, deposit)
	state.AddBalance(common.TicketPoolAddr, deposit)
	tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
	fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", candidate.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
	tickIdList, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	/** GetTicketList */
	if tickets, err := ppos.GetTicketList(state, tickIdList); nil != err {
		errFn("GetTicketList err", err)
	}else {
		printObject("GetTicketList ticketArr:", tickets, logger)
	}
}
func TestPpos_GetTicketList(t *testing.T) {
	ppos_GetTicketList(t, t.Log, t.Error)
}
func BenchmarkPpos_GetTicketList(b *testing.B) {
	ppos_GetTicketList(b, b.Log, b.Error)
}

// test GetCandidateTicketIds
func ppos_GetCandidateTicketIds (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetCandidateTicketIds getting state err", err)
	}else {
		state = st
	}
	logFn("test GetCandidateTicketIds ...")


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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)

	timeMap := make(map[uint32]int64)
	logFn("VOTING START .............................................................")

	startTime := time.Now().UnixNano() / 1e6
	voteOwner := common.HexToAddress("0x20")
	deposit := new(big.Int).SetUint64(10)
	state.SubBalance(voteOwner, deposit)
	state.AddBalance(common.TicketPoolAddr, deposit)
	tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
	fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", candidate.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
	_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	tickIds, err := ppos.GetCandidateTicketIds(state, candidate.CandidateId)
	if nil != err {
		errFn("GetCandidateTicketIds err", err)
	}else {
		printObject("GetCandidateTicketIds:", tickIds, logger)
	}
}
func TestPpos_GetCandidateTicketIds(t *testing.T) {
	ppos_GetCandidateTicketIds(t, t.Log, t.Error)
}
func BenchmarkPpos_GetCandidateTicketIds(b *testing.B) {
	ppos_GetCandidateTicketIds(b, b.Log, b.Error)
}

// test GetCandidateEpoch
func ppos_GetCandidateEpoch  (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetCandidateEpoch getting state err", err)
	}else {
		state = st
	}
	logFn("test GetCandidateEpoch ...")


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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)

	timeMap := make(map[uint32]int64)
	logFn("VOTING START .............................................................")

	startTime := time.Now().UnixNano() / 1e6
	voteOwner := common.HexToAddress("0x20")
	deposit := new(big.Int).SetUint64(10)
	state.SubBalance(voteOwner, deposit)
	state.AddBalance(common.TicketPoolAddr, deposit)
	tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
	fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", candidate.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
	_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	epoch, err := ppos.GetCandidateEpoch(state, candidate.CandidateId)
	if nil != err {
		errFn("GetCandidateEpoch err", err)
	}else {
		logFn("GetCandidateEpoch:", epoch)
	}
}
func TestPpos_GetCandidateEpoch(t *testing.T) {
	ppos_GetCandidateEpoch(t, t.Log, t.Error)
}
func BenchmarkPpos_GetCandidateEpoch(b *testing.B) {
	ppos_GetCandidateEpoch(b, b.Log, b.Error)
}

// test GetTicketPrice
func ppos_GetTicketPrice (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetTicketPrice getting state err", err)
	}else {
		state = st
	}
	logFn("test GetTicketPrice ...")

	/** test  GetTicketPrice */
	logFn("test GetTicketPrice ...")
	if num, err := ppos.GetTicketPrice(state); nil != err {
		errFn("GetTicketPrice err", err)
	}else {
		logFn("GetTicketPrice:", num)
	}
}
func TestPpos_GetTicketPrice(t *testing.T) {
	ppos_GetTicketPrice(t, t.Log, t.Error)
}
func BenchmarkPpos_GetTicketPrice(b *testing.B) {
	ppos_GetTicketPrice(b, b.Log, b.Error)
}

// test GetCandidateAttach
func ppos_GetCandidateAttach (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetCandidateAttach getting state err", err)
	}else {
		state = st
	}
	logFn("test GetCandidateAttach ...")


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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)

	timeMap := make(map[uint32]int64)
	logFn("VOTING START .............................................................")

	startTime := time.Now().UnixNano() / 1e6
	voteOwner := common.HexToAddress("0x20")
	deposit := new(big.Int).SetUint64(10)
	state.SubBalance(voteOwner, deposit)
	state.AddBalance(common.TicketPoolAddr, deposit)
	tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
	fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", candidate.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
	_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	attach, err := ppos.GetCandidateAttach(state, candidate.CandidateId)
	if nil != err {
		errFn("GetCandidateAttach err", err)
	}else {
		printObject("GetCandidateAttach:", attach, logger)
	}
}
func TestPpos_GetCandidateAttach(t *testing.T) {
	ppos_GetCandidateAttach(t, t.Log, t.Error)
}
func BenchmarkPpos_GetCandidateAttach(b *testing.B) {
	ppos_GetCandidateAttach(b, b.Log, b.Error)
}

// test Notify
func ppos_Notify (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test Notify getting state err", err)
	}else {
		state = st
	}
	logFn("test Notify ...")


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
	if err := ppos.SetCandidate(state, candidate.CandidateId, candidate); nil != err {
		errFn("SetCandidate err:", err)
	}


	/** vote ticket */
	var count uint32 = 0
	var blockNumber = new(big.Int).SetUint64(10)

	timeMap := make(map[uint32]int64)
	logFn("VOTING START .............................................................")

	startTime := time.Now().UnixNano() / 1e6
	voteOwner := common.HexToAddress("0x20")
	deposit := new(big.Int).SetUint64(10)
	state.SubBalance(voteOwner, deposit)
	state.AddBalance(common.TicketPoolAddr, deposit)
	tempBlockNumber := new(big.Int).SetUint64(blockNumber.Uint64())
	fmt.Println("给当前候选人投票为:", "投票人为:", voteOwner.String(), " ,投了1张票给:", candidate.CandidateId.String(), " ,投票时的块高为:", tempBlockNumber.String())
	_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	if err := ppos.Notify(state, new(big.Int).SetUint64(10)); nil != err {
		errFn("Notify err", err)
	}else {
		logFn("Notify success ... ")
	}
}
func TestPpos_Notify(t *testing.T) {
	ppos_Notify(t, t.Log, t.Error)
}
func BenchmarkPpos_Notify(b *testing.B) {
	ppos_Notify(b, b.Log, b.Error)
}

/** about other */

// test UpdateNodeList

// test Submit2Cache

// test GetFormerRound

// test GetCurrentRound

// test GetNextRound

// test SetNodeCache



