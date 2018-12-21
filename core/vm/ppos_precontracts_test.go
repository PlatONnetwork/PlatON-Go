package vm

import (
	"Platon-go/common"
	"Platon-go/core/types"
	"Platon-go/p2p/discover"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
)

func TestExecute(t *testing.T) {
	var command = map[string]interface{}{
		"recalled": recalled,
	}
	input, _ := hex.DecodeString("f8c28800000000000000f188726563616c6c6564b88230783166336138363732333438666636623738396534313637363261643533653639303633313338623865623464383738303130313635386632346232333639663161386530393439393232366234363764386263306334653033653164633930336466383537656562336336373733336432316236616165653238343065343239aa30786632313664366534633137303937613630656532623865356338383934316364396630373236336201")
	result, error := execute(input, command)
	fmt.Println(result, error)
}

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

func recalled(nodeId discover.NodeID, owner common.Address, deposit *big.Int) ([]byte, error) {
	fmt.Println("nodeId:", nodeId, "owner:", owner, "deposit:", deposit)
	return nil, nil
}
