package vm_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"time"

	"testing"
	//"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
)

func TestCandidatePoolOverAll(t *testing.T) {

	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(7000)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// CandidateDetails(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateDetails input==>", "nodeId: ", nodeId.String())
	resByte, err := candidateContract.CandidateDetails(nodeId)
	if nil != err {
		fmt.Println("CandidateDetails fail", "err", err)
		return
	}
	if nil == resByte {
		fmt.Println("The candidate info is null")
		return
	}
	fmt.Println("The candidate info is: ", vm.ResultByte2Json(resByte))

	// CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error)
	withdraw := big.NewInt(100)
	fmt.Println("CandidateApplyWithdraw input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateApplyWithdraw(nodeId, withdraw)
	if nil != err {
		fmt.Println("CandidateApplyWithdraw fail", "err", err)
	}
	fmt.Println("CandidateApplyWithdraw success")

	// CandidateWithdraw(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateWithdraw input==>", "nodeId: ", nodeId.String())
	_, err = candidateContract.CandidateWithdraw(nodeId)
	if nil != err {
		fmt.Println("CandidateWithdraw fail", "err", err)
	}
	fmt.Println("CandidateWithdraw success")

	// CandidateWithdrawInfos(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateWithdrawInfos input==>", "nodeId: ", nodeId.String())
	resByte, err = candidateContract.CandidateWithdrawInfos(nodeId)
	if nil != err {
		fmt.Println("CandidateWithdrawInfos fail", "err", err)
		return
	}
	if nil == resByte {
		fmt.Println("The CandidateWithdrawInfos is null")
		return
	}
	fmt.Println("The CandidateWithdrawInfos is: ", vm.ResultByte2Json(resByte))

	// CandidateList() ([]byte, error)
	resByte, err = candidateContract.CandidateList()
	if nil != err {
		fmt.Println("CandidateList fail", "err", err)
	}
	if nil == resByte {
		fmt.Println("The candidate list is null")
		return
	}
	fmt.Println("The candidate list is: ", vm.ResultByte2Json(resByte))
}

func newContract() *vm.Contract {
	callerAddress := vm.AccountRef(common.HexToAddress("0x12"))
	contract := vm.NewContract(callerAddress, callerAddress, big.NewInt(1000), uint64(1))
	return contract
}

func newEvm() *vm.EVM {
	state, _ := newChainState()
	candidatePool, ticketPool := newPool()
	evm := &vm.EVM{
		StateDB:       state,
		CandidatePool: candidatePool,
		TicketPool:    ticketPool,
	}
	context := vm.Context{
		BlockNumber: big.NewInt(7),
	}
	evm.Context = context
	return evm
}

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
	} else {
		state = statedb
	}
	return state, nil
}

func newPool() (*pposm.CandidatePool, *pposm.TicketPool) {
	configs := params.PposConfig{
		Candidate: &params.CandidateConfig{
			MaxChair:          1,
			MaxCount:          3,
			RefundBlockNumber: 1,
		},
		TicketConfig: &params.TicketConfig{
			MaxCount:          100,
			ExpireBlockNumber: 2,
		},
	}
	return pposm.NewCandidatePool(&configs), pposm.NewTicketPool(&configs)
}

func TestCandidateDeposit(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}

	// CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string) ([]byte, error)
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(7000)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")
}

func TestCandidateDetails(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(7000)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// CandidateDetails(nodeId discover.NodeID) ([]byte, error)
	// nodeId = discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	fmt.Println("CandidateDetails input==>", "nodeId: ", nodeId.String())
	resByte, err := candidateContract.CandidateDetails(nodeId)
	if nil != err {
		fmt.Println("CandidateDetails fail", "err", err)
		return
	}
	if nil == resByte {
		fmt.Println("The candidate info is null")
		return
	}
	fmt.Println("The candidate info is: ", vm.ResultByte2Json(resByte))
}

func TestGetBatchCandidateDetail(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(7000)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit1 success")

	nodeId = discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner = common.HexToAddress("0x12")
	fee = uint64(7000)
	host = "10.0.0.2"
	port = "8548"
	extra = "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit2 success")

	// GetBatchCandidateDetail(nodeIds []discover.NodeID) ([]byte, error)
	nodeIds := []discover.NodeID{discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")}
	input, _ := json.Marshal(nodeIds)
	fmt.Println("GetBatchCandidateDetail input==>", "nodeIds: ", string(input))
	resByte, err := candidateContract.GetBatchCandidateDetail(nodeIds)
	if nil != err {
		fmt.Println("GetBatchCandidateDetail fail", "err", err)
	}
	if nil == resByte {
		fmt.Println("The candidate info is null")
		return
	}
	fmt.Println("The batch candidate info is: ", vm.ResultByte2Json(resByte))
}

func TestCandidateList(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(7000)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit1 success")

	nodeId = discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner = common.HexToAddress("0x12")
	fee = uint64(6800)
	host = "10.0.0.2"
	port = "8548"
	extra = "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit2 success")

	// CandidateList() ([]byte, error)
	resByte, err := candidateContract.CandidateList()
	if nil != err {
		fmt.Println("CandidateList fail", "err", err)
	}
	if nil == resByte {
		fmt.Println("The candidate list is null")
		return
	}
	fmt.Println("The candidate list is: ", vm.ResultByte2Json(resByte))
}

func TestSetCandidateExtra(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(1)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// SetCandidateExtra(nodeId discover.NodeID, extra string) ([]byte, error)
	extra = "this node is powerful"
	fmt.Println("SetCandidateExtra input=>", "nodeId: ", nodeId.String(), "extra: ", extra)
	_, err = candidateContract.SetCandidateExtra(nodeId, extra)
	if nil != err {
		fmt.Println("SetCandidateExtra fail", "err", err)
	}
	fmt.Println("SetCandidateExtra success")

	// CandidateDetails(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateDetails input==>", "nodeId: ", nodeId.String())
	resByte, err := candidateContract.CandidateDetails(nodeId)
	if nil != err {
		fmt.Println("CandidateDetails fail", "err", err)
		return
	}
	if nil == resByte {
		fmt.Println("The candidate info is null")
		return
	}
	fmt.Println("The candidate info is: ", vm.ResultByte2Json(resByte))
}

func TestCandidateApplyWithdraw(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(1)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error)
	withdraw := big.NewInt(100)
	fmt.Println("CandidateApplyWithdraw input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateApplyWithdraw(nodeId, withdraw)
	if nil != err {
		fmt.Println("CandidateApplyWithdraw fail", "err", err)
	}
	fmt.Println("CandidateApplyWithdraw success")
}

func TestCandidateWithdraw(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(7000)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error)
	withdraw := big.NewInt(100)
	fmt.Println("CandidateApplyWithdraw input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateApplyWithdraw(nodeId, withdraw)
	if nil != err {
		fmt.Println("CandidateApplyWithdraw fail", "err", err)
	}
	fmt.Println("CandidateApplyWithdraw success")

	// CandidateWithdraw(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateWithdraw input==>", "nodeId: ", nodeId.String())
	_, err = candidateContract.CandidateWithdraw(nodeId)
	if nil != err {
		fmt.Println("CandidateWithdraw fail", "err", err)
	}
	fmt.Println("CandidateWithdraw success")

	// CandidateDetails(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateDetails input==>", "nodeId: ", nodeId.String())
	resByte, err := candidateContract.CandidateDetails(nodeId)
	if nil != err {
		fmt.Println("CandidateDetails fail", "err", err)
		return
	}
	if nil == resByte {
		fmt.Println("The candidate info is null")
		return
	}
	fmt.Println("The candidate info is: ", vm.ResultByte2Json(resByte))
}

func TestCandidateWithdrawInfos(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(1)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error)
	withdraw := big.NewInt(100)
	fmt.Println("CandidateApplyWithdraw input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateApplyWithdraw(nodeId, withdraw)
	if nil != err {
		fmt.Println("CandidateApplyWithdraw fail", "err", err)
	}
	fmt.Println("CandidateApplyWithdraw success")

	// CandidateWithdraw(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateWithdraw input==>", "nodeId: ", nodeId.String())
	_, err = candidateContract.CandidateWithdraw(nodeId)
	if nil != err {
		fmt.Println("CandidateWithdraw fail", "err", err)
	}
	fmt.Println("CandidateWithdraw success")

	// CandidateWithdrawInfos(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateWithdrawInfos input==>", "nodeId: ", nodeId.String())
	resByte, err := candidateContract.CandidateWithdrawInfos(nodeId)
	if nil != err {
		fmt.Println("CandidateWithdrawInfos fail", "err", err)
		return
	}
	if nil == resByte {
		fmt.Println("The CandidateWithdrawInfos is null")
		return
	}
	fmt.Println("The CandidateWithdrawInfos is: ", vm.ResultByte2Json(resByte))
}

func TestTime(t *testing.T) {
	fmt.Printf("时间戳（毫秒）：%v;\n", time.Now().UnixNano()/1e6)
}

func TestCandidatePoolEncode(t *testing.T) {
	//"enode://1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429@192.168.9.181:16789",
	//"enode://751f4f62fccee84fc290d0c68d673e4b0cc6975a5747d2baccb20f954d59ba3315d7bfb6d831523624d003c8c2d33451129e67c3eef3098f711ef3b3e268fd3c@192.168.9.182:16789",
	//"enode://b6c8c9f99bfebfa4fb174df720b9385dbd398de699ec36750af3f38f8e310d4f0b90447acbef64bdf924c4b59280f3d42bb256e6123b53e9a7e99e4c432549d6@192.168.9.183:16789",
	//"enode://97e424be5e58bfd4533303f8f515211599fd4ffe208646f7bfdf27885e50b6dd85d957587180988e76ae77b4b6563820a27b16885419e5ba6f575f19f6cb36b0@192.168.9.184:16789"
	nodeId := []byte("0x97e424be5e58bfd4533303f8f515211599fd4ffe208646f7bfdf27885e50b6dd85d957587180988e76ae77b4b6563820a27b16885419e5ba6f575f19f6cb36b0")
	owner := []byte("0x740ce31b3fac20dac379db243021a51e80ad00d7")
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-引力区\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
	// CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string)
	var CandidateDeposit [][]byte
	CandidateDeposit = make([][]byte, 0)
	CandidateDeposit = append(CandidateDeposit, uint64ToBytes(1001))
	CandidateDeposit = append(CandidateDeposit, []byte("CandidateDeposit"))
	CandidateDeposit = append(CandidateDeposit, nodeId)
	CandidateDeposit = append(CandidateDeposit, owner)
	CandidateDeposit = append(CandidateDeposit, uint64ToBytes(8000))
	//CandidateDeposit = append(CandidateDeposit, bigIntStrToBytes("130000000000000000000"))
	CandidateDeposit = append(CandidateDeposit, []byte("192.168.9.184"))
	CandidateDeposit = append(CandidateDeposit, []byte("16789"))
	CandidateDeposit = append(CandidateDeposit, []byte(extra))
	bufDeposit := new(bytes.Buffer)
	err := rlp.Encode(bufDeposit, CandidateDeposit)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateDeposit encode rlp data fail")
	} else {
		fmt.Println("CandidateDeposit data rlp: ", hexutil.Encode(bufDeposit.Bytes()))
	}

	// CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int)
	var CandidateApplyWithdraw [][]byte
	CandidateApplyWithdraw = make([][]byte, 0)
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, uint64ToBytes(1002))
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, []byte("CandidateApplyWithdraw"))
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, nodeId)
	withdraw, ok := new(big.Int).SetString("14d1120d7b160000", 16)
	if !ok {
		t.Errorf("big int setstring fail")
	}
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, withdraw.Bytes())
	bufApply := new(bytes.Buffer)
	err = rlp.Encode(bufApply, CandidateApplyWithdraw)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateApplyWithdraw encode rlp data fail")
	} else {
		fmt.Println("CandidateApplyWithdraw data rlp: ", hexutil.Encode(bufApply.Bytes()))
	}

	// CandidateWithdraw(nodeId discover.NodeID)
	var CandidateWithdraw [][]byte
	CandidateWithdraw = make([][]byte, 0)
	CandidateWithdraw = append(CandidateWithdraw, uint64ToBytes(1003))
	CandidateWithdraw = append(CandidateWithdraw, []byte("CandidateWithdraw1"))
	CandidateWithdraw = append(CandidateWithdraw, nodeId)
	bufWith := new(bytes.Buffer)
	err = rlp.Encode(bufWith, CandidateWithdraw)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateWithdraw encode rlp data fail")
	} else {
		fmt.Println("CandidateWithdraw data rlp: ", hexutil.Encode(bufWith.Bytes()))
	}

	// CandidateWithdrawInfos(nodeId discover.NodeID)
	var CandidateWithdrawInfos [][]byte
	CandidateWithdrawInfos = make([][]byte, 0)
	CandidateWithdrawInfos = append(CandidateWithdrawInfos, uint64ToBytes(0xf1))
	CandidateWithdrawInfos = append(CandidateWithdrawInfos, []byte("CandidateWithdrawInfos"))
	CandidateWithdrawInfos = append(CandidateWithdrawInfos, nodeId)
	bufWithInfos := new(bytes.Buffer)
	err = rlp.Encode(bufWithInfos, CandidateWithdrawInfos)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateWithdrawInfos encode rlp data fail")
	} else {
		fmt.Println("CandidateWithdrawInfos data rlp: ", hexutil.Encode(bufWithInfos.Bytes()))
	}

	// CandidateDetails(nodeId discover.NodeID)
	var CandidateDetails [][]byte
	CandidateDetails = make([][]byte, 0)
	CandidateDetails = append(CandidateDetails, uint64ToBytes(0xf1))
	CandidateDetails = append(CandidateDetails, []byte("CandidateDetails"))
	CandidateDetails = append(CandidateDetails, nodeId)
	bufDetails := new(bytes.Buffer)
	err = rlp.Encode(bufDetails, CandidateDetails)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateDetails encode rlp data fail")
	} else {
		fmt.Println("CandidateDetails data rlp: ", hexutil.Encode(bufDetails.Bytes()))
	}

	// CandidateList()
	var CandidateList [][]byte
	CandidateList = make([][]byte, 0)
	CandidateList = append(CandidateList, uint64ToBytes(0xf1))
	CandidateList = append(CandidateList, []byte("CandidateList"))
	bufCList := new(bytes.Buffer)
	err = rlp.Encode(bufCList, CandidateList)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateList encode rlp data fail")
	} else {
		fmt.Println("CandidateList data rlp: ", hexutil.Encode(bufCList.Bytes()))
	}
}

func TestCandidatePoolDecode(t *testing.T) {

	//HexString -> []byte
	rlpcode, _ := hex.DecodeString("f901a18800000000000003e99043616e6469646174654465706f736974b88230783166336138363732333438666636623738396534313637363261643533653639303633313338623865623464383738303130313635386632346232333639663161386530393439393232366234363764386263306334653033653164633930336466383537656562336336373733336432316236616165653238343065343239aa30783734306365333162336661633230646163333739646232343330323161353165383061643030643789070c1cc73b00c800008d3139322e3136382e392e313833853136373839b8b87b226e6f64654e616d65223a2022506c61746f6e2d4265696a696e67222c20226e6f6465506f727472616974223a2022222c226e6f64654469736372697074696f6e223a2022506c61744f4e2de5bc95e58a9be58cba222c226e6f64654465706172746d656e74223a20224a555a4958222c226f6666696369616c57656273697465223a202268747470733a2f2f7777772e706c61746f6e2e6e6574776f726b2f222c2274696d65223a313534363530333635313139307d")
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(rlpcode), &source); err != nil {
		fmt.Println(err)
		t.Errorf("TestRlpDecode decode rlp data fail")
	}

	for i, v := range source {
		switch i {
		case 0:
			fmt.Println("i: ", i, " v: ", bytesToUint64(v))
		case 4:
			fmt.Println("i: ", i, " v: ", bigIntByteToStr(v))
		default:
			fmt.Println("i: ", i, " v: ",  string(v))
		}


	}
}

func TestAppendSlice(t *testing.T) {
	a := []int{0, 1, 2, 3, 4}
	i := 2
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)
}

func uint64ToBytes(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return buf[:]
}

func bytesToUint64 (by []byte) uint64 {
	return binary.BigEndian.Uint64(by)
}

func bigIntStrToBytes(str string) []byte {
	a, _ := new(big.Int).SetString(str, 10)
	return a.Bytes()
}

func bigIntByteToStr (by []byte) string {
	a := new(big.Int).SetBytes(by)
	return a.String()
}