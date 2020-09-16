package platonstats

import (
	"encoding/json"
	"math/big"
	"regexp"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	address = common.MustBech32ToAddress("lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj")
)

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
