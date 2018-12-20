package vm

import (
	"Platon-go/common"
	"Platon-go/core/types"
	"Platon-go/p2p/discover"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
)

func TestDecodeResultStr(t *testing.T) {
	ticket := types.Ticket{
		TicketId:    common.HexToHash("0x0123456789012345678901234567890123456789012345678901234567890123"),
		Owner:       common.HexToAddress("0x0123456789012345678901234567890123456789"),
		Deposit:     big.NewInt(1),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		BlockNumber: big.NewInt(100),
		State:       1,
	}
	ticketIds := make([]common.Hash, 0)
	ticketIds = append(ticketIds, common.BytesToHash([]byte("1")))
	data, _ := json.Marshal(ticket)
	sdata := DecodeResultStr(string(data))
	fmt.Println("GetTicketDetail==> ", "json: ", string(data), " []byte: ", sdata)
}
