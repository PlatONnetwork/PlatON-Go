package vm

import (
	"Platon-go/common"
	"Platon-go/common/byteutil"
	"Platon-go/common/hexutil"
	"Platon-go/core/types"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
)

// send tx from script
func TestTicketPoolEncode(t *testing.T) {
	nodeId := []byte("0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429")
	// owner := []byte("0xf216d6e4c17097a60ee2b8e5c88941cd9f07263b")

	// VoteTicket(count uint64, price *big.Int, nodeId discover.NodeID)
	var VoteTicket [][]byte
	VoteTicket = make([][]byte, 0)
	VoteTicket = append(VoteTicket, uint64ToBytes(0xf1))
	VoteTicket = append(VoteTicket, []byte("VoteTicket"))
	VoteTicket = append(VoteTicket, byteutil.Uint64ToBytes(1000))
	price, ok :=new(big.Int).SetString("14d1120d7b160000", 16)
	if !ok {
		t.Errorf("big int setstring fail")
	}
	VoteTicket = append(VoteTicket, price.Bytes())
	VoteTicket = append(VoteTicket, nodeId)
	bufVote := new(bytes.Buffer)
	err := rlp.Encode(bufVote, VoteTicket)
	if err != nil {
		fmt.Println(err)
		t.Errorf("VoteTicket encode rlp data fail")
	} else {
		fmt.Println("VoteTicket data rlp: ", hexutil.Encode(bufVote.Bytes()))
	}

	// GetCandidateTicketIds(nodeId discover.NodeID, blockNumber *big.Int)
	var GetCandidateTicketIds [][]byte
	GetCandidateTicketIds = make([][]byte, 0)
	GetCandidateTicketIds = append(GetCandidateTicketIds, uint64ToBytes(0xf1))
	GetCandidateTicketIds = append(GetCandidateTicketIds, []byte("GetCandidateTicketIds"))
	GetCandidateTicketIds = append(GetCandidateTicketIds, nodeId)
	blockNumber, ok :=new(big.Int).SetString("14d1120d7b160000", 16)
	if !ok {
		t.Errorf("big int setstring fail")
	}
	GetCandidateTicketIds = append(GetCandidateTicketIds, blockNumber.Bytes())
	bufCandidate := new(bytes.Buffer)
	err = rlp.Encode(bufCandidate, GetCandidateTicketIds)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetCandidateTicketIds encode rlp data fail")
	} else {
		fmt.Println("GetCandidateTicketIds data rlp: ", hexutil.Encode(bufCandidate.Bytes()))
	}
}

func TestTicketPoolDecode(t *testing.T)  {

	//HexString -> []byte
	rlpcode, _ := hex.DecodeString("f8aa8800000000000000f18a566f74655469636b65748800000000000003e88814d1120d7b160000b88230783166336138363732333438666636623738396534313637363261643533653639303633313338623865623464383738303130313635386632346232333639663161386530393439393232366234363764386263306334653033653164633930336466383537656562336336373733336432316236616165653238343065343239")
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(rlpcode), &source); err != nil {
		fmt.Println(err)
		t.Errorf("TestRlpDecode decode rlp data fail")
	}

	for i,v := range source {
		fmt.Println("i: ", i, " v: ", hex.EncodeToString(v))
	}
}


func TestDecodeResultStr(t *testing.T) {
	ticket := types.Ticket{
		TicketId:		common.HexToHash("0x0123456789012345678901234567890123456789012345678901234567890123"),
		Owner:			common.HexToAddress("0x0123456789012345678901234567890123456789"),
		Deposit:		big.NewInt(1),
		CandidateId:	discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		BlockNumber:	big.NewInt(100),
		State:			1,
	}
	/*ticketIds := make([]common.Hash, 0)
	ticketIds = append(ticketIds, common.BytesToHash([]byte("1")))*/
	data, _ := json.Marshal(ticket)
	sdata := DecodeResultStr(string(data))
	log.Info("GetPoolRemainder==> ", "json: ", string(data), " []byte: ", sdata)
}

