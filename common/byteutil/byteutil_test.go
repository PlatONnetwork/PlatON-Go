// create by platon
package byteutil

import (
	"fmt"
	"math/big"
	"testing"
)

func TestByteUtil(t *testing.T)  {
	/*fmt.Println(IntToBytes(50))
	fmt.Println(BytesToString([]byte{97,98,99}))
	fmt.Println(BytesToInt([]byte{0,0,0,50}))*/

	big1 := new(big.Int).SetInt64(1)
	fmt.Println(big1.Bytes())

	//fmt.Println(BytesToInt64(StringToBytes("abc")))
	//fmt.Println(reflect.TypeOf(big1))
}
