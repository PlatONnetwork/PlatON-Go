package vm

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
)

//void invoke() {}
var emptyInvokeWasm = hexutil.MustDecode("0x0061736d01000000018480808000016000000382808080000100048480808000017000000583808080000100010681808080000007938080800002066d656d6f7279020006696e766f6b6500000a8880808000018280808000000b")

func TestWagonEngine_Run(t *testing.T) {

	//contract := &Contract{
	//	Code: emptyInvokeWasm,
	//}
	//engine := &wagonEngine{
	//	contract: contract,
	//}
	//p := &WasmInoke{VM: 1, Args: &WasmParams{
	//	FuncName: []byte{},
	//	Args: [][]byte{
	//		[]byte{1, 2},
	//	},
	//}}
	//input, _ := rlp.EncodeToBytes(p)
	//_, err := engine.Run(input, false)
	//assert.Nil(t, err)
}
