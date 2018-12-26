package vm_test

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestTicketPoolOverAll(t *testing.T) {

	ticketContract := vm.TicketContract{
		newContract(),
		newEvm(),
	}
	candidateContract := vm.CandidateContract{
		newContract(),
		newEvm(),
	}

	// CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string) ([]byte, error)
	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	owner := common.HexToAddress("0x12")
	fee := uint64(1)
	host := "10.0.0.1"
	port := "8548"
	extra := "extra data"
	_, err := candidateContract.CandidateDeposit(nodeId, owner, fee, host, port, extra)
	if nil != err {
		fmt.Println("CandidateDeposit fail", "err", err)
	}
	fmt.Println("质押成功...")

	// VoteTicket(count uint64, price *big.Int, nodeId discover.NodeID) ([]byte, error)
	count := uint64(1000)
	price := big.NewInt(1)
	_, err = ticketContract.VoteTicket(count, price, nodeId)
	if nil != err {
		fmt.Println("VoteTicket fail", "err", err)
	}

	// GetCandidateTicketIds(nodeId discover.NodeID) ([]byte, error)
	_, err = ticketContract.GetCandidateTicketIds(nodeId)
	if nil != err {
		fmt.Println("GetCandidateTicketIds fail", "err", err)
	}

	// GetTicketDetail(ticketId common.Hash) ([]byte, error)
	ticketId := common.HexToHash("e69d8e6dbc1ee87d7fb20600f3fc6744f28b637d43b5a130b2904c30d12e9b30")
	_, err = ticketContract.GetTicketDetail(ticketId)
	if nil != err {
		fmt.Println("GetTicketDetail fail", "err", err)
	}

	// GetBatchTicketDetail(ticketIds []common.Hash) ([]byte, error)
	ticketIds := []common.Hash{common.HexToHash("e69d8e6dbc1ee87d7fb20600f3fc6744f28b637d43b5a130b2904c30d12e9b30"), common.HexToHash("008674dae3f0c660158fe602589c5505b20e24be4caa8f65c0f92ff372149ccc")}
	_, err = ticketContract.GetBatchTicketDetail(ticketIds)
	if nil != err {
		fmt.Println("GetTicketDetail fail", "err", err)
	}
}

func TestTicketPoolEncode(t *testing.T) {
	nodeId := []byte("0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429")

	// VoteTicket(count uint64, price *big.Int, nodeId discover.NodeID)
	var VoteTicket [][]byte
	VoteTicket = make([][]byte, 0)
	VoteTicket = append(VoteTicket, byteutil.Uint64ToBytes(0xf1))
	VoteTicket = append(VoteTicket, []byte("VoteTicket"))
	VoteTicket = append(VoteTicket, byteutil.Uint64ToBytes(1))
	VoteTicket = append(VoteTicket, big.NewInt(1).Bytes())
	VoteTicket = append(VoteTicket, nodeId)
	bufVoteTicket := new(bytes.Buffer)
	err := rlp.Encode(bufVoteTicket, VoteTicket)
	if err != nil {
		fmt.Println(err)
		t.Errorf("VoteTicket encode rlp data fail")
	} else {
		fmt.Println("VoteTicket data rlp: ", hexutil.Encode(bufVoteTicket.Bytes()))
	}

	// GetCandidateTicketIds(nodeId discover.NodeID)
	var GetCandidateTicketIds [][]byte
	GetCandidateTicketIds = make([][]byte, 0)
	GetCandidateTicketIds = append(GetCandidateTicketIds, byteutil.Uint64ToBytes(0xf1))
	GetCandidateTicketIds = append(GetCandidateTicketIds, []byte("GetCandidateTicketIds"))
	GetCandidateTicketIds = append(GetCandidateTicketIds, nodeId)
	bufGetCandidateTicketIds := new(bytes.Buffer)
	err = rlp.Encode(bufGetCandidateTicketIds, GetCandidateTicketIds)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetCandidateTicketIds encode rlp data fail")
	} else {
		fmt.Println("GetCandidateTicketIds data rlp: ", hexutil.Encode(bufGetCandidateTicketIds.Bytes()))
	}

	// GetTicketDetail(ticketId common.Hash)
	ticketId := common.Hex2Bytes("64a25965db45003819598699d46a138b3e91ab0be338463359a21eb3e896ac68")
	var GetTicketDetail [][]byte
	GetTicketDetail = make([][]byte, 0)
	GetTicketDetail = append(GetTicketDetail, byteutil.Uint64ToBytes(0xf1))
	GetTicketDetail = append(GetTicketDetail, []byte("GetTicketDetail"))
	GetTicketDetail = append(GetTicketDetail, ticketId)
	bufGetTicketDetail := new(bytes.Buffer)
	err = rlp.Encode(bufGetTicketDetail, GetTicketDetail)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetTicketDetail encode rlp data fail")
	} else {
		fmt.Println("GetTicketDetail data rlp: ", hexutil.Encode(bufGetTicketDetail.Bytes()))
	}
}

func TestTicketPoolDecode(t *testing.T) {
	//HexString -> []byte
	rlpcode, _ := hex.DecodeString("f8a28800000000000000f18a566f74655469636b65748800000000000003e801b88230783166336138363732333438666636623738396534313637363261643533653639303633313338623865623464383738303130313635386632346232333639663161386530393439393232366234363764386263306334653033653164633930336466383537656562336336373733336432316236616165653238343065343239")
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(rlpcode), &source); err != nil {
		fmt.Println(err)
		t.Errorf("TestRlpDecode decode rlp data fail")
	}

	for i, v := range source {
		fmt.Println("i: ", i, " v: ", hex.EncodeToString(v))
	}
}
