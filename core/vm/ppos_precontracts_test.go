package vm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
	"testing"
)

func TestExecute(t *testing.T) {
	var command = map[string]interface{}{
		"recalled": recalled,
	}
	input, _ := hex.DecodeString("f8c28800000000000000f188726563616c6c6564b88230783166336138363732333438666636623738396534313637363261643533653639303633313338623865623464383738303130313635386632346232333639663161386530393439393232366234363764386263306334653033653164633930336466383537656562336336373733336432316236616165653238343065343239aa30786632313664366534633137303937613630656532623865356338383934316364396630373236336201")
	_, err := execute(input, command)
	if nil != err {
		fmt.Println("execute fail", "err", err)
	}
}

func TestDecodeResultStr(t *testing.T) {
	ticket := types.Ticket{
		Owner:       common.HexToAddress("0x0123456789012345678901234567890123456789"),
		Deposit:     big.NewInt(1),
		CandidateId: discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"),
		BlockNumber: big.NewInt(100),
		Remaining:   1000,
	}
	data, _ := json.Marshal(ticket)
	sdata := DecodeResultStr(string(data))
	json := ResultByte2Json(sdata)
	fmt.Println("origin: ", string(data), "[]byte: ", sdata, "json: ", json)
}

func recalled(nodeId discover.NodeID, owner common.Address, deposit *big.Int) ([]byte, error) {
	fmt.Println("nodeId:", nodeId, "owner:", owner, "deposit:", deposit)
	return nil, nil
}

func ResultByte2Json(origin []byte) string {
	resultByte := origin[64:]
	return string(resultByte)
}
