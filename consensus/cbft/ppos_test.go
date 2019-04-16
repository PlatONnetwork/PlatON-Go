package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/core/ppos_storage"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"math/big"
	"testing"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"encoding/json"
	"net"
	"strconv"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"sync/atomic"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/rand"
	"time"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/core/ticketcache"
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
		CandidateConfig: &params.CandidateConfig{
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
		candidateContext:  pposm.NewCandidatePoolContext(&configs),
		ticketContext: pposm.NewTicketPoolContext(&configs),
	}

	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)

	// Initialize ppos storage
	ppos_storage.NewPPosTemp(db)

	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, nil, vm.Config{}, nil)

	ppos.SetCandidateContextOption(blockchain, buildInitialNodes())

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
	printObject("ppos.candidatePoolText:", ppos.candidateContext, t)
	printObject("ppos.ticketPool:", ppos.ticketContext, t)
}

// test BlockProducerIndex
func ppos_BlockProducerIndex(logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	curr := bc.CurrentBlock()
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012342")
	num,_ := ppos.BlockProducerIndex(curr.Number(), curr.Hash(), curr.Number(), nodeId, 1)
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
	if can := ppos.GetCandidate(state, candidate.CandidateId, big.NewInt(1)); nil == can {
		errFn("GetCandidate err")
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

// test GetCandidateArr
func ppos_GetCandidateArr (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test GetCandidateArr getting state err", err)
	}else {
		state = st
	}
	logFn("test GetCandidateArr ...")



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

	/** test GetCandidate */
	logFn("test GetCandidateArr ...")
	canArr := ppos.GetCandidateArr(state, big.NewInt(1), []discover.NodeID{candidate.CandidateId, candidate2.CandidateId}...)
	printObject("GetCandidateArr", canArr, logger)
}
func TestPpos_GetCandidateArr(t *testing.T) {
	ppos_GetCandidateArr(t, t.Log, t.Error)
}
func BenchmarkPpos_GetCandidateArr(b *testing.B) {
	ppos_GetCandidateArr(b, b.Log, b.Error)
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
		fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	_, err := ppos.Election(state, common.Hash{}, big.NewInt(20))
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
		fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	_, err := ppos.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ := ppos.GetWitness(state, 1, big.NewInt(1))
	printObject("next Witnesses", canArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := ppos.Switch(state, big.NewInt(1))
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
		fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	_, err := ppos.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ := ppos.GetWitness(state, 1, big.NewInt(1))
	printObject("next Witnesses", canArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := ppos.Switch(state, big.NewInt(1))
	logFn("Switch was success ", flag)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ = ppos.GetWitness(state, 0, big.NewInt(1))
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
		fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	_, err := ppos.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	canArr, _ := ppos.GetWitness(state, 1, big.NewInt(1))
	printObject("next Witnesses", canArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := ppos.Switch(state, big.NewInt(1))
	logFn("Switch was success ", flag)

	/** test Election again */
	logFn("test Election again ...")
	_, err = ppos.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election again was successful err", err)

	/** test GetAllWitness */
	logFn("test GetAllWitness ...")
	preArr, canArr, nextArr, _ := ppos.GetAllWitness(state, big.NewInt(1))
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
	can := ppos.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), big.NewInt(1))
	printObject("GetCandidate", can, logger)

	/** test WithdrawCandidate */
	logFn("test WithdrawCandidate ...")
	ok1 := ppos.WithdrawCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), new(big.Int).SetUint64(uint64(99)), new(big.Int).SetUint64(uint64(10)))
	logFn("error", ok1)

	/** test GetCandidate */
	logFn("test GetCandidate ...")
	can2 := ppos.GetCandidate(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), big.NewInt(1))
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
	canArr := ppos.GetChosens(state, 0, big.NewInt(1))
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
		fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	_, err := ppos.Election(state, common.Hash{}, big.NewInt(20))
	logFn("Whether election was successful err", err)

	/** test GetWitness */
	logFn("test GetWitness ...")
	nodeIdArr, _ := ppos.GetWitness(state, 1, big.NewInt(1))
	printObject("next Witnesses", nodeIdArr, logger)

	/** test switch */
	logFn("test Switch ...")
	flag := ppos.Switch(state, big.NewInt(1))
	logFn("Switch was success ", flag)

	/** test GetChairpersons */
	logFn("test GetChairpersons ...")
	canArr := ppos.GetChairpersons(state, big.NewInt(1))
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
		fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	flag := ppos.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), big.NewInt(1))
	logFn("isdefeat", flag)

	/** test GetDefeat */
	logFn("test GetDefeat ...")
	defeatArr := ppos.GetDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), big.NewInt(1))
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
		fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	flag := ppos.IsDefeat(state, discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), big.NewInt(1))
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
		Deposit: 		new(big.Int).SetUint64(120), // 99
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
		fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	ownerAddr := ppos.GetOwner(state, candidate.CandidateId, big.NewInt(1))
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
	num := ppos.GetRefundInterval(big.NewInt(1))
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
	num := ppos.GetPoolNumber(state)
	logFn("GetPoolNumber:", num)
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
	fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
	fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
	num, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}else{
		logFn("ticket success num", num)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	/** GetTicket */
	if ticket := ppos.GetTicket(state, state.TxHash()); nil == ticket {
		errFn("GetTicket err")
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
	fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
	num, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}else{
		logFn("ticket success num", num)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	/** GetTicketList */
	if tickets := ppos.GetTicketList(state, []common.Hash{state.TxHash()}); len(tickets) == 0 {
		errFn("GetTicketList err")
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
	fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
	_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	tickIds := ppos.GetCandidateTicketIds(state, candidate.CandidateId)
	printObject("GetCandidateTicketIds:", tickIds, logger)
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
	fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
	_, err := ppos.VoteTicket(state, voteOwner, 1, deposit, candidate.CandidateId, tempBlockNumber)
	if nil != err {
		errFn("vote ticket error:", err)
	}
	atomic.AddUint32(&count, 1)
	timeMap[count] = (time.Now().UnixNano() / 1e6) - startTime
	logFn("VOTING END .............................................................")

	epoch := ppos.GetCandidateEpoch(state, candidate.CandidateId)
	logFn("GetCandidateEpoch:", epoch)
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
	num := ppos.GetTicketPrice(state)
	logFn("GetTicketPrice:", num)
}
func TestPpos_GetTicketPrice(t *testing.T) {
	ppos_GetTicketPrice(t, t.Log, t.Error)
}
func BenchmarkPpos_GetTicketPrice(b *testing.B) {
	ppos_GetTicketPrice(b, b.Log, b.Error)
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
	fmt.Println("Voting Current Candidate:", "voter is:", voteOwner.String(), " ,ticket num:", candidate.CandidateId.String(), " ,blocknumber when voted:", tempBlockNumber.String())
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
func ppos_UpdateNodeList (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	logFn("test UpdateNodeList ...")
	genesis := bc.Genesis()
	ppos.UpdateNodeList(bc, genesis.Number(), genesis.Hash())
}
func TestPpos_UpdateNodeList(t *testing.T) {
	ppos_UpdateNodeList(t, t.Log, t.Error)
}
func BenchmarkPpos_UpdateNodeList(b *testing.B) {
	ppos_UpdateNodeList(b, b.Log, b.Error)
}

// test Submit2Cache
func getBlockMaxData() (ticketcache.TicketCache, error) {
	//every nodeid has 256 ticket total has 200 nodeid
	ret := ticketcache.NewTicketCache()
	for n:=0; n<200; n++ {
		nodeid := make([]byte, 0, 64)
		nodeid = append(nodeid, crypto.Keccak256Hash([]byte("nodeid"), byteutil.IntToBytes(n)).Bytes()...)
		nodeid = append(nodeid, crypto.Keccak256Hash([]byte("nodeid"), byteutil.IntToBytes(n*10)).Bytes()...)
		NodeId, err := discover.BytesID(nodeid)
		if err!=nil {
			return ret, err
		}
		tids := make([]common.Hash, 0)
		for i:=0; i< 51200/200 ; i++ {
			tids = append(tids, crypto.Keccak256Hash([]byte("tid"), byteutil.IntToBytes(i)))
		}
		ret[NodeId] = tids
	}
	return ret, nil
}
func ppos_Submit2Cache (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ldb, err := ethdb.NewLDBDatabase("./data/platon/chaindata", 0, 0)
	if err!=nil {
		errFn("NewLDBDatabase faile")
	}
	tc := ticketcache.NewTicketIdsCache(ldb)
	for i:=0; i< 20; i++  {
		number := big.NewInt(int64(i))
		bkhash := crypto.Keccak256Hash(byteutil.IntToBytes(i))
		mapCache ,err := getBlockMaxData()
		if err!=nil {
			errFn("getMaxtickets faile", "err: ", err)
		}
		tc.Submit2Cache(number, big.NewInt(int64(20)), bkhash, mapCache)

	}
	ldb.Close()
}
func TestPpos_Submit2Cache(t *testing.T) {
	ppos_Submit2Cache(t, t.Log, t.Error)
}
func BenchmarkPpos_Submit2Cache(b *testing.B) {
	ppos_Submit2Cache(b, b.Log, b.Error)
}

// TODO Hash

// test GetFormerRound
func ppos_GetFormerRound (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	logFn("test GetFormerRound ...")
	genesis := bc.Genesis()
	round := ppos.GetFormerRound(genesis.Number(), genesis.Hash())
	printObject("GetFormerRound", round.nodes, logger)
}
func TestPpos_GetFormerRound(t *testing.T) {
	ppos_GetFormerRound(t, t.Log, t.Error)
}
func BenchmarkPpos_GetFormerRound(b *testing.B) {
	ppos_GetFormerRound(b, b.Log, b.Error)
}

// test GetCurrentRound
func ppos_GetCurrentRound (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	logFn("test GetCurrentRound ...")
	genesis := bc.Genesis()
	round := ppos.GetCurrentRound(genesis.Number(), genesis.Hash())
	printObject("GetCurrentRound", round.nodes, logger)
}
func TestPpos_GetCurrentRound(t *testing.T) {
	ppos_GetCurrentRound(t, t.Log, t.Error)
}
func BenchmarkPpos_GetCurrentRound(b *testing.B) {
	ppos_GetCurrentRound(b, b.Log, b.Error)
}


// test GetNextRound
func ppos_GetNextRound (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})) {
	ppos, bc := buildPpos()
	logFn("test GetNextRound ...")
	genesis := bc.Genesis()
	round := ppos.GetNextRound(genesis.Number(), genesis.Hash())
	printObject("GetNextRound", round.nodes, logger)
}
func TestPpos_GetNextRound(t *testing.T) {
	ppos_GetNextRound(t, t.Log, t.Error)
}
func BenchmarkPpos_GetNextRound(b *testing.B) {
	ppos_GetNextRound(b, b.Log, b.Error)
}

// test SetNodeCache
func ppos_SetNodeCache (logger interface{}, logFn func (args ... interface{}), errFn func (args ... interface{})){
	ppos, bc := buildPpos()
	var state *state.StateDB
	if st, err := bc.State(); nil != err {
		errFn("test SetNodeCache getting state err", err)
	}else {
		state = st
	}
	logFn("test SetNodeCache ...")
	genesis := bc.Genesis()
	if err := ppos.SetNodeCache(state, genesis.Number(), big.NewInt(0), big.NewInt(1), genesis.Hash(), common.HexToHash("0xa1d63b9e5f36c9b12e6aed34612bc1f6e846d1e94a53f52673f2433a30e9ac51"), common.HexToHash("0xa1d63b9e5f36c9b12e6aed34612bc1f6e846d1e94a53f52673f2433a30e9bd62")); nil != err {
		errFn("SetNodeCache err", err)
	}else {
		logFn("SetNodeCache success ... ")
	}
}
func TestPpos_SetNodeCache(t *testing.T) {
	ppos_SetNodeCache(t, t.Log, t.Error)
}
func BenchmarkPpos_SetNodeCache(b *testing.B) {
	ppos_SetNodeCache(b, b.Log, b.Error)
}


