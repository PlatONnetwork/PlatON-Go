package platonstats

import (
	"encoding/json"
	"math/big"
	"regexp"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	address = common.MustBech32ToAddress("lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj")
)

func makeReceipt(addr common.Address) *types.Receipt {
	receipt := types.NewReceipt(nil, false, 0)
	receipt.Logs = []*types.Log{
		{Address: addr},
	}
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	return receipt
}

func Test_blockJson(t *testing.T) {

	/*header := &types.Header{GasLimit: 9424776, Number: big.NewInt(236500), GasUsed: 0, Time: big.NewInt(1614686649157)}
	block := types.NewBlockWithHeader(header)*/

	block := types.NewSimplifiedBlock(42322, common.HexToHash("499987a73fa100f582328c92c1239262edf5c0a3479face652c89f60314aa805"))

	blockJsonMapping, err := jsonBlock(block)

	statsBlockExt := &StatsBlockExt{
		Block: blockJsonMapping,
	}

	if err != nil {
		t.Error("marshal block to json string error")
	} else {
		jsonBytes, err := json.Marshal(statsBlockExt)
		if err != nil {
			t.Error("marshal block to json string error")
		} else {
			t.Log(string(jsonBytes))
		}
	}

}
func Test_statsBlockExt(t *testing.T) {
	blockEnc := common.FromHex("f90264f901fda00000000000000000000000000000000000000000000000000000000000000000948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080832fefd8825208845506eb0780b8510376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a8023a09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b180")
	var block *types.Block
	if err := rlp.DecodeBytes(blockEnc, block); err != nil {
		t.Fatal("decode block data error: ", err)
	}

	brief := collectBrief(block)

	blockJsonMapping, err := jsonBlock(block)
	if err != nil {
		t.Fatal("marshal block to json string error", err)
	}
	statsBlockExt := &StatsBlockExt{
		BlockType:   brief.BlockType,
		Epoch:       brief.Epoch,
		NodeID:      brief.NodeID,
		NodeAddress: brief.NodeAddress,
		Receipts:    []*types.Receipt{makeReceipt(address)},
		Block:       blockJsonMapping,
	}

	jsonBytes, err := json.Marshal(statsBlockExt)
	if err != nil {
		t.Fatal("marshal platon stats block message to json string error", err)
	} else {
		t.Log("marshal platon stats block", "blockNumber", block.NumberU64(), "json", string(jsonBytes))
	}

}
func T2estUrl(t *testing.T) {
	re := regexp.MustCompile("([^:@]*)(:([^@]*))?@(.+)")
	url := "center:myPasswordd@ws://localhost:1900"
	parts := re.FindStringSubmatch(url)
	for i := 0; i < len(parts); i++ {
		t.Logf("url parts: [%d]%s", i, parts[i])
	}
}

func Test_Unmarshal_accountCheckingMessage(t *testing.T) {
	message := AccountCheckingMessage{
		BlockNumber: uint64(12132131),
		AccountList: []*AccountItem{
			&AccountItem{Addr: address, Balance: big.NewInt(1)},
			&AccountItem{Addr: address, Balance: big.NewInt(2)},
		}}
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		t.Fatal("Failed to marshal accountCheckingMessage to json format", err)
	} else {
		t.Log("accountCheckingMessage json format:" + string(jsonBytes))
	}

	/*jsonStr := "{\n    \"blockNumber\":12132131,\n    \"accountList\":[\n        {\n            \"addr\":\"lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj\",\n            \"balance\":1\n        },\n        {\n            \"addr\":\"lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj\",\n            \"balance\":2\n        }]\n}"
	jsonBytes = []byte(jsonStr)*/
	var data AccountCheckingMessage
	if len(jsonBytes) > 0 {
		if err := json.Unmarshal(jsonBytes, &data); err != nil {
			t.Fatal("Failed to unmarshal json to accountCheckingMessage", err)
		} else {
			t.Log("accountCheckingMessage.accountList", data.AccountList[1].Balance.Int64())
		}
	}
}
