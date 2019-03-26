package vm_test

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"testing"
)

func TestTicketPoolOverAll(t *testing.T) {
	contract := newContract()
	evm := newEvm()

	ticketContract := vm.TicketContract{
		contract,
		evm,
	}
	candidateContract := vm.CandidateContract{
		contract,
		evm,
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
	fmt.Println("setcandidate successfully...")

	// VoteTicket(count uint32, price *big.Int, nodeId discover.NodeID) ([]byte, error)
	count := uint32(1000)
	price := big.NewInt(1)
	fmt.Println("VoteTicket input==>", "count: ", count, "price: ", price, "nodeId: ", nodeId.String())
	resByte, err := ticketContract.VoteTicket(count, price, nodeId)
	if nil != err {
		fmt.Println("VoteTicket fail", "err", err)
	}
	fmt.Println("The list of generated ticketId is: ", vm.ResultByte2Json(resByte))
}

func TestVoteTicket(t *testing.T) {
	contract := newContract()
	evm := newEvm()

	ticketContract := vm.TicketContract{
		contract,
		evm,
	}
	candidateContract := vm.CandidateContract{
		contract,
		evm,
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

	// VoteTicket(count uint64, price *big.Int, nodeId discover.NodeID) ([]byte, error)
	count := uint32(1000)
	price := big.NewInt(1)
	fmt.Println("VoteTicket input==>", "count: ", count, "price: ", price, "nodeId: ", nodeId.String())
	resByte, err := ticketContract.VoteTicket(count, price, nodeId)
	if nil != err {
		fmt.Println("VoteTicket fail", "err", err)
	}
	fmt.Println("The list of generated ticketId is: ", vm.ResultByte2Json(resByte))
}

func TestGetCandidateTicketCount(t *testing.T) {
	contract := newContract()
	evm := newEvm()

	ticketContract := vm.TicketContract{
		contract,
		evm,
	}
	candidateContract := vm.CandidateContract{
		contract,
		evm,
	}
	nodeId1 := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint32(7000)
	host := "192.168.9.184"
	port := "16789"
	extra := "{\"nodeName\": \"Platon-Beijing\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Gravitational area\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId1: ", nodeId1.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err := candidateContract.CandidateDeposit(nodeId1, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit1 success")

	nodeId2 := discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner = common.HexToAddress("0x12")
	fee = uint32(8000)
	host = "192.168.9.185"
	port = "16789"
	extra = "{\"nodeName\": \"Platon-Shenzhen\", \"nodePortrait\": \"\",\"nodeDiscription\": \"PlatON-Cosmic wave\",\"nodeDepartment\": \"JUZIX\",\"officialWebsite\": \"https://www.platon.network/sz\",\"time\":1546503651190}"
	fmt.Println("CandidateDeposit input==>", "nodeId2: ", nodeId2.String(), "owner: ", owner.Hex(), "fee: ", fee, "host: ", host, "port: ", port, "extra: ", extra)
	_, err = candidateContract.CandidateDeposit(nodeId2, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("CandidateDeposit2 success")

	// CandidateList() ([]byte, error)
	resByte, err := candidateContract.GetCandidateList()
	if nil != err {
		fmt.Println("CandidateList fail", "err", err)
	}
	if nil == resByte {
		fmt.Println("The candidate list is null")
		return
	}
	fmt.Println("The candidate list is: ", vm.ResultByte2Json(resByte))

	// Vote to Candidate1
	count := uint32(100)
	price := big.NewInt(1)
	fmt.Println("VoteTicket input==>", "count: ", count, "price: ", price, "nodeId1: ", nodeId1.String())
	resByte, err = ticketContract.VoteTicket(count, price, nodeId1)
	if nil != err {
		fmt.Println("VoteTicket fail", "err", err)
	}
	fmt.Println("The list of generated ticketId is: ", vm.ResultByte2Json(resByte))

	// Vote to Candidate2
	count = uint32(101)
	price = big.NewInt(1)
	fmt.Println("VoteTicket input==>", "count: ", count, "price: ", price, "nodeId2: ", nodeId2.String())
	resByte, err = ticketContract.VoteTicket(count, price, nodeId2)
	if nil != err {
		fmt.Println("VoteTicket fail", "err", err)
	}
	fmt.Println("The list of generated ticketId is: ", vm.ResultByte2Json(resByte))

	// GetBatchCandidateTicketCount(nodeIds []discover.NodeID) ([]byte, error)
	fmt.Println("GetBatchCandidateTicketCount input==>", "nodeIds: ", nodeId1.String(), nodeId2.String())
	var nodeIds []discover.NodeID
	nodeId1 = discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	nodeId2 = discover.MustHexID("0x11234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	nodeIds = append(append(nodeIds, nodeId1), nodeId2)
	resByte, err = ticketContract.GetCandidateTicketCount(nodeIds)
	if nil != err {
		fmt.Println("GetBatchCandidateTicketCount fail", "err", err)
	}
	if nil == resByte {
		fmt.Println("The candidates's ticket list is null")
		return
	}
	fmt.Println("The number of candidate's ticket is: ", vm.ResultByte2Json(resByte))
}

func TestGetCandidateEpoch(t *testing.T) {
	contract := newContract()
	evm := newEvm()

	ticketContract := vm.TicketContract{
		contract,
		evm,
	}
	candidateContract := vm.CandidateContract{
		contract,
		evm,
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
	count := uint32(1000)
	price := big.NewInt(1)
	fmt.Println("VoteTicket input==>", "count: ", count, "price: ", price, "nodeId: ", nodeId.String())
	resByte, err := ticketContract.VoteTicket(count, price, nodeId)
	if nil != err {
		fmt.Println("VoteTicket fail", "err", err)
	}

	// GetCandidateEpoch(nodeId discover.NodeID) ([]byte, error)
	fmt.Println("GetCandidateEpoch input==>", "nodeId: ", nodeId.String())
	resByte, err = ticketContract.GetCandidateEpoch(nodeId)
	if nil != err {
		fmt.Println("GetCandidateEpoch fail", "err", err)
	}
	fmt.Println("The candidate's epoch is: ", vm.ResultByte2Json(resByte))

}

func TestGetTicketPrice(t *testing.T) {
	ticketContract := vm.TicketContract{
		newContract(),
		newEvm(),
	}

	// GetTicketPrice() ([]byte, error)
	resByte, err := ticketContract.GetTicketPrice()
	if nil != err {
		fmt.Println("GetTicketPrice fail", "err", err)
	}
	fmt.Println("The ticket price is: ", vm.ResultByte2Json(resByte))
}

func TestTicketPoolEncode(t *testing.T) {
	nodeId := []byte("1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429")
	// VoteTicket(count uint32, price *big.Int, nodeId discover.NodeID)
	price, _ := new(big.Int).SetString("100000000000000000000", 10)
	var VoteTicket [][]byte
	VoteTicket = make([][]byte, 0)
	VoteTicket = append(VoteTicket, byteutil.Uint64ToBytes(1000))
	VoteTicket = append(VoteTicket, []byte("VoteTicket"))
	VoteTicket = append(VoteTicket, byteutil.Uint32ToBytes(100))
	VoteTicket = append(VoteTicket, price.Bytes())
	VoteTicket = append(VoteTicket, nodeId)
	bufVoteTicket := new(bytes.Buffer)
	err := rlp.Encode(bufVoteTicket, VoteTicket)
	if err != nil {
		fmt.Println(err)
		t.Errorf("VoteTicket encode rlp data fail")
	} else {
		fmt.Println("VoteTicket data rlp: ", hexutil.Encode(bufVoteTicket.Bytes()))
	}
	// GetCandidateTicketCount(nodeIds []discover.NodeID)
	nodeId1 := "0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"
	nodeId2 := "0x2f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"
	nodeId3 := "0x3f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"
	nodeIds := nodeId1 + ":" + nodeId2 + ":" + nodeId3
	var GetCandidateTicketCount [][]byte
	GetCandidateTicketCount = make([][]byte, 0)
	GetCandidateTicketCount = append(GetCandidateTicketCount, byteutil.Uint64ToBytes(0xf1))
	GetCandidateTicketCount = append(GetCandidateTicketCount, []byte("GetCandidateTicketCount"))
	GetCandidateTicketCount = append(GetCandidateTicketCount, []byte(nodeIds))
	bufGetCandidateTicketCount := new(bytes.Buffer)
	err = rlp.Encode(bufGetCandidateTicketCount, GetCandidateTicketCount)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetCandidateTicketCount encode rlp data fail")
	} else {
		fmt.Println("GetCandidateTicketCount data rlp: ", hexutil.Encode(bufGetCandidateTicketCount.Bytes()))
	}

	// GetTicketCountByTxHash(ticketIds []common.Hash)
	txHash1 := "0xb26d49ba1bc4246f2456ee03c947276b8a8d2c95f1207b8e7bfaa333246e98c6"
	//txHash2 := "0xbf1145c872f38d50f53cbc15040c9131102a43d6895d239b97938cb1b109e03e"
	txHashs := txHash1 /*+ ":" + txHash2*/
	var GetTicketCountByTxHash [][]byte
	GetTicketCountByTxHash = make([][]byte, 0)
	GetTicketCountByTxHash = append(GetTicketCountByTxHash, byteutil.Uint64ToBytes(0xf1))
	GetTicketCountByTxHash = append(GetTicketCountByTxHash, []byte("GetTicketCountByTxHash"))
	GetTicketCountByTxHash = append(GetTicketCountByTxHash, []byte(txHashs))
	bufGetTicketCountByTxHash := new(bytes.Buffer)
	err = rlp.Encode(bufGetTicketCountByTxHash, GetTicketCountByTxHash)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetTicketCountByTxHash encode rlp data fail")
	} else {
		fmt.Println("GetTicketCountByTxHash data rlp: ", hexutil.Encode(bufGetTicketCountByTxHash.Bytes()))
	}

	// GetCandidateEpoch(nodeId discover.NodeID) ([]byte, error)G
	var GetCandidateEpoch [][]byte
	GetCandidateEpoch = make([][]byte, 0)
	GetCandidateEpoch = append(GetCandidateEpoch, byteutil.Uint64ToBytes(0xf1))
	GetCandidateEpoch = append(GetCandidateEpoch, []byte("GetCandidateEpoch"))
	GetCandidateEpoch = append(GetCandidateEpoch, nodeId)
	bufGetCandidateEpoch := new(bytes.Buffer)
	err = rlp.Encode(bufGetCandidateEpoch, GetCandidateEpoch)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetCandidateEpoch encode rlp data fail")
	} else {
		fmt.Println("GetCandidateEpoch data rlp: ", hexutil.Encode(bufGetCandidateEpoch.Bytes()))
	}

	// GetPoolRemainder() ([]byte, error)
	var GetPoolRemainder [][]byte
	GetPoolRemainder = make([][]byte, 0)
	GetPoolRemainder = append(GetPoolRemainder, byteutil.Uint64ToBytes(0xf1))
	GetPoolRemainder = append(GetPoolRemainder, []byte("GetPoolRemainder"))
	bufGetPoolRemainder := new(bytes.Buffer)
	err = rlp.Encode(bufGetPoolRemainder, GetPoolRemainder)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetPoolRemainder encode rlp data fail")
	} else {
		fmt.Println("GetPoolRemainder data rlp: ", hexutil.Encode(bufGetPoolRemainder.Bytes()))
	}

	// GetTicketPrice() ([]byte, error)
	var GetTicketPrice [][]byte
	GetTicketPrice = make([][]byte, 0)
	GetTicketPrice = append(GetTicketPrice, byteutil.Uint64ToBytes(0xf1))
	GetTicketPrice = append(GetTicketPrice, []byte("GetTicketPrice"))
	bufGetTicketPrice := new(bytes.Buffer)
	err = rlp.Encode(bufGetTicketPrice, GetTicketPrice)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetTicketPrice encode rlp data fail")
	} else {
		fmt.Println("GetTicketPrice data rlp: ", hexutil.Encode(bufGetTicketPrice.Bytes()))
	}

}
