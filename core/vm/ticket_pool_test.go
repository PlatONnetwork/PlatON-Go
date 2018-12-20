package vm

import (
	"Platon-go/common/byteutil"
	"Platon-go/common/hexutil"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

// send tx from script
func TestTicketPoolEncode(t *testing.T) {
	nodeId := []byte("0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429")

	// VoteTicket(count uint64, price *big.Int, nodeId discover.NodeID)
	var VoteTicket [][]byte
	VoteTicket = make([][]byte, 0)
	VoteTicket = append(VoteTicket, uint64ToBytes(0xf1))
	VoteTicket = append(VoteTicket, []byte("VoteTicket"))
	VoteTicket = append(VoteTicket, byteutil.Uint64ToBytes(1000))
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

	// GetTicketDetail(ticketId common.Hash)
	/*ticketId := []byte("")
	var GetTicketDetail [][]byte
	GetTicketDetail = make([][]byte, 0)
	GetTicketDetail = append(GetTicketDetail, uint64ToBytes(0xf1))
	GetTicketDetail = append(GetTicketDetail, []byte("GetTicketDetail"))
	GetTicketDetail = append(GetTicketDetail, ticketId)
	bufGetTicketDetail := new(bytes.Buffer)
	err = rlp.Encode(bufGetTicketDetail, GetTicketDetail)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetTicketDetail encode rlp data fail")
	} else {
		fmt.Println("GetTicketDetail data rlp: ", hexutil.Encode(bufGetTicketDetail.Bytes()))
	}*/

	// GetCandidateTicketIds(nodeId discover.NodeID)
	var GetCandidateTicketIds [][]byte
	GetCandidateTicketIds = make([][]byte, 0)
	GetCandidateTicketIds = append(GetCandidateTicketIds, uint64ToBytes(0xf1))
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
