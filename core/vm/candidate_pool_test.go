package vm_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
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
	"testing"
	"time"
)

func TestCandidatePoolOverAll(t *testing.T) {

	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// GetCandidateDetails(nodeIds []discover.NodeID) ([]byte, error)
	fmt.Println("GetCandidateDetails input==>", "nodeIds: ", nodeId.String())
	var nodeIds []discover.NodeID
	nodeIds = append(nodeIds, nodeId)
	resByte, err := candidateContract.GetCandidateDetails(nodeIds)
	if nil != err {
		fmt.Println("GetCandidateDetails fail", "err", err)
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

	// GetCandidateWithdrawInfos(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("GetCandidateWithdrawInfos input==>", "nodeId: ", nodeId.String())
	resByte, err = candidateContract.GetCandidateWithdrawInfos(nodeId)
	if nil != err {
		fmt.Println("GetCandidateWithdrawInfos fail", "err", err)
		return
	}
	if nil == resByte {
		fmt.Println("The GetCandidateWithdrawInfos is null")
		return
	}
	fmt.Println("The GetCandidateWithdrawInfos is: ", vm.ResultByte2Json(resByte))

	// GetCandidateList() ([]byte, error)
	resByte, err = candidateContract.GetCandidateList()
	if nil != err {
		fmt.Println("GetCandidateList fail", "err", err)
	}
	if nil == resByte {
		fmt.Println("The candidate list is null")
		return
	}
	fmt.Println("The candidate list is: ", vm.ResultByte2Json(resByte))
}

func newChainState() (*state.StateDB, error) {
	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, nil, vm.Config{}, nil)

	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		return nil, errors.New("reference statedb failed" + err.Error())
	} else {
		state = statedb
	}
	return state, nil
}

func newPool() (*pposm.CandidatePoolContext, *pposm.TicketPoolContext) {
	configs := &params.PposConfig{
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
	return pposm.NewCandidatePoolContext(configs), pposm.NewTicketPoolContext(configs)
}

func newEvm() *vm.EVM {
	state, _ := newChainState()
	candidatePoolContext, ticketPoolContext := newPool()
	evm := &vm.EVM{
		StateDB:              state,
		CandidatePoolContext: candidatePoolContext,
		TicketPoolContext:    ticketPoolContext,
	}
	context := vm.Context{
		BlockNumber: big.NewInt(7),
	}
	evm.Context = context
	return evm
}

func newContract() *vm.Contract {
	callerAddress := vm.AccountRef(common.HexToAddress("0x12"))
	contract := vm.NewContract(callerAddress, callerAddress, big.NewInt(1000), uint64(1))
	return contract
}

func TestCandidateDeposit(t *testing.T) {
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}

	// CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint32, host, port, extra string) ([]byte, error)
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
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
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// GetCandidateDetails(nodeIds []discover.NodeID) ([]byte, error)
	// nodeId = discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	var nodeIds []discover.NodeID
	nodeIds = append(nodeIds, nodeId)
	fmt.Println("GetCandidateDetails input==>", "nodeIds: ", nodeId.String())
	resByte, err := candidateContract.GetCandidateDetails(nodeIds)
	if nil != err {
		fmt.Println("GetCandidateDetails fail", "err", err)
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
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit1 success")

	nodeId = discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner = common.HexToAddress("0x12")
	fee = uint32(8000)
	host = "192.168.9.185"
	port = "16789"
	extra = "{\"nodeName\": \"Platon-Shenzhen\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Cosmic wave\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/sz\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit2 success")

	// GetCandidateDetails(nodeIds []discover.NodeID) ([]byte, error)
	nodeIds := []discover.NodeID{discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"), discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")}
	input, _ := json.Marshal(nodeIds)
	fmt.Println("GetBatchCandidateDetail input==>", "nodeIds: ", string(input))
	resByte, err := candidateContract.GetCandidateDetails(nodeIds)
	if nil != err {
		fmt.Println("GetCandidateDetails fail", "err", err)
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
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit1 success")

	nodeId = discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner = common.HexToAddress("0x12")
	fee = uint32(8000)
	host = "192.168.9.185"
	port = "16789"
	extra = "{\"nodeName\": \"Platon-Shenzhen\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Cosmic wave\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/sz\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit2 success")

	// GetCandidateList() ([]byte, error)
	resByte, err := candidateContract.GetCandidateList()
	if nil != err {
		fmt.Println("GetCandidateList fail", "err", err)
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
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId: ", nodeId.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit success")

	// SetCandidateExtra(nodeId discover.NodeID, extra string) ([]byte, error)
	extra = "this node is powerful!!"
	fmt.Println("SetCandidateExtra input=>", "nodeId: ", nodeId.String(), "extra: ", extra)
	_, err = candidateContract.SetCandidateExtra(nodeId, extra)
	if nil != err {
		fmt.Println("SetCandidateExtra fail", "err", err)
	}
	fmt.Println("SetCandidateExtra success")

	// GetCandidateDetails(nodeIds []discover.NodeID) ([]byte, error)
	fmt.Println("GetCandidateDetails input==>", "nodeIds: ", nodeId.String())
	var nodeIds []discover.NodeID
	nodeIds = append(nodeIds, nodeId)
	resByte, err := candidateContract.GetCandidateDetails(nodeIds)
	if nil != err {
		fmt.Println("GetCandidateDetails fail", "err", err)
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
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
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
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
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

	// GetCandidateDetails(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("CandidateDetails input==>", "nodeId: ", nodeId.String())
	var nodeIds []discover.NodeID
	nodeIds = append(nodeIds, nodeId)
	resByte, err := candidateContract.GetCandidateDetails(nodeIds)
	if nil != err {
		fmt.Println("GetCandidateDetails fail", "err", err)
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
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
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

	// GetCandidateWithdrawInfos(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("GetCandidateWithdrawInfos input==>", "nodeId: ", nodeId.String())
	resByte, err := candidateContract.GetCandidateWithdrawInfos(nodeId)
	if nil != err {
		fmt.Println("GetCandidateWithdrawInfos fail", "err", err)
		return
	}
	if nil == resByte {
		fmt.Println("The GetCandidateWithdrawInfos is null")
		return
	}
	fmt.Println("The GetCandidateWithdrawInfos is: ", vm.ResultByte2Json(resByte))
}

func TestTime(t *testing.T) {
	fmt.Printf("Timestamp (ms)：%v;\n", time.Now().UnixNano()/1e6)
}

func TestCandidatePoolEncode(t *testing.T) {
	//"enode://1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429@192.168.9.181:16789",
	//"enode://751f4f62fccee84fc290d0c68d673e4b0cc6975a5747d2baccb20f954d59ba3315d7bfb6d831523624d003c8c2d33451129e67c3eef3098f711ef3b3e268fd3c@192.168.9.182:16789",
	//"enode://b6c8c9f99bfebfa4fb174df720b9385dbd398de699ec36750af3f38f8e310d4f0b90447acbef64bdf924c4b59280f3d42bb256e6123b53e9a7e99e4c432549d6@192.168.9.183:16789",
	//"enode://97e424be5e58bfd4533303f8f515211599fd4ffe208646f7bfdf27885e50b6dd85d957587180988e76ae77b4b6563820a27b16885419e5ba6f575f19f6cb36b0@192.168.9.184:16789"
	nodeId := []byte("1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429")
	owner := []byte("0x740ce31b3fac20dac379db243021a51e80ad00d7")
	extra := "{\"nodeName\": \"Platon-Shanghai\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-eastern area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/sz\",\"time\":1546503651100}"
	// CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint32, host, port, extra string)
	var CandidateDeposit [][]byte
	CandidateDeposit = make([][]byte, 0)
	CandidateDeposit = append(CandidateDeposit, byteutil.Uint64ToBytes(1001))
	CandidateDeposit = append(CandidateDeposit, []byte("CandidateDeposit"))
	CandidateDeposit = append(CandidateDeposit, nodeId)
	CandidateDeposit = append(CandidateDeposit, owner)
	CandidateDeposit = append(CandidateDeposit, byteutil.Uint32ToBytes(7900))
	//CandidateDeposit = append(CandidateDeposit, bigIntStrToBytes("130000000000000000000"))
	CandidateDeposit = append(CandidateDeposit, []byte("0.0.0.0"))
	CandidateDeposit = append(CandidateDeposit, []byte("30303"))
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
	CandidateApplyWithdraw = append(CandidateApplyWithdraw, byteutil.Uint64ToBytes(1002))
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
	CandidateWithdraw = append(CandidateWithdraw, byteutil.Uint64ToBytes(1003))
	CandidateWithdraw = append(CandidateWithdraw, []byte("CandidateWithdraw"))
	CandidateWithdraw = append(CandidateWithdraw, nodeId)
	bufWith := new(bytes.Buffer)
	err = rlp.Encode(bufWith, CandidateWithdraw)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateWithdraw encode rlp data fail")
	} else {
		fmt.Println("CandidateWithdraw data rlp: ", hexutil.Encode(bufWith.Bytes()))
	}

	// GetCandidateWithdrawInfos(nodeId discover.NodeID)
	var GetCandidateWithdrawInfos [][]byte
	GetCandidateWithdrawInfos = make([][]byte, 0)
	GetCandidateWithdrawInfos = append(GetCandidateWithdrawInfos, byteutil.Uint64ToBytes(0xf1))
	GetCandidateWithdrawInfos = append(GetCandidateWithdrawInfos, []byte("GetCandidateWithdrawInfos"))
	GetCandidateWithdrawInfos = append(GetCandidateWithdrawInfos, nodeId)
	bufGetCandidateWithdrawInfos := new(bytes.Buffer)
	err = rlp.Encode(bufGetCandidateWithdrawInfos, GetCandidateWithdrawInfos)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetCandidateWithdrawInfos encode rlp data fail")
	} else {
		fmt.Println("GetCandidateWithdrawInfos data rlp: ", hexutil.Encode(bufGetCandidateWithdrawInfos.Bytes()))
	}

	// GetCandidateDetails(nodeIds []discover.NodeID)
	var GetCandidateDetails [][]byte
	nodeId1 := "0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"
	nodeId2 := "0x2f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"
	nodeId3 := "0x3f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"
	nodeIds := nodeId1 + ":" + nodeId2 + ":" + nodeId3
	GetCandidateDetails = make([][]byte, 0)
	GetCandidateDetails = append(GetCandidateDetails, byteutil.Uint64ToBytes(0xf1))
	GetCandidateDetails = append(GetCandidateDetails, []byte("GetCandidateDetails"))
	GetCandidateDetails = append(GetCandidateDetails, []byte(nodeIds))
	bufGetCandidateDetails := new(bytes.Buffer)
	err = rlp.Encode(bufGetCandidateDetails, GetCandidateDetails)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetCandidateDetails encode rlp data fail")
	} else {
		fmt.Println("GetCandidateDetails data rlp: ", hexutil.Encode(bufGetCandidateDetails.Bytes()))
	}

	// GetCandidateList()
	var GetCandidateList [][]byte
	GetCandidateList = make([][]byte, 0)
	GetCandidateList = append(GetCandidateList, byteutil.Uint64ToBytes(0xf1))
	GetCandidateList = append(GetCandidateList, []byte("GetCandidateList"))
	bufGetCandidateList := new(bytes.Buffer)
	err = rlp.Encode(bufGetCandidateList, GetCandidateList)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetCandidateList encode rlp data fail")
	} else {
		fmt.Println("GetCandidateList data rlp: ", hexutil.Encode(bufGetCandidateList.Bytes()))
	}

	// GetVerifiersList()
	var GetVerifiersList [][]byte
	GetVerifiersList = make([][]byte, 0)
	GetVerifiersList = append(GetVerifiersList, byteutil.Uint64ToBytes(0xf1))
	GetVerifiersList = append(GetVerifiersList, []byte("GetVerifiersList"))
	bufGetVerifiersList := new(bytes.Buffer)
	err = rlp.Encode(bufGetVerifiersList, GetVerifiersList)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetVerifiersList encode rlp data fail")
	} else {
		fmt.Println("GetVerifiersList data rlp: ", hexutil.Encode(bufGetVerifiersList.Bytes()))
	}

}

func TestCandidatePoolDecode(t *testing.T) {

	//HexString -> []byte
	rlpcode, _ := hex.DecodeString("f8a28800000000000003e88a566f74655469636b6574820008880de0b6b3a7640000b8803131346534386632316434643833656339616333396136323036326138303461303536363734326438306231393164653562613233613464633235663762656461306537386464313639333532613761643362313135383464303661303161303963653034376164383864653962646362363338383565383164653030613464")
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(rlpcode), &source); err != nil {
		fmt.Println(err)
		t.Errorf("TestRlpDecode decode rlp data fail")
	}

	for i, v := range source {
		switch i {
		case 0:
			fmt.Println("i: ", i, " v: ", byteutil.BytesTouint64(v))
		case 2:
			fmt.Println("i: ", i, " v: ", byteutil.BytesTouint32(v))
		default:
			fmt.Println("i: ", i, " v: ", string(v))
		}

	}
}
