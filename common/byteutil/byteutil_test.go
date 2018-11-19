// create by platon
package byteutil

import (
	"fmt"
	"testing"
)

func TestByteUtil(t *testing.T)  {
	/*fmt.Println(StringToBytes("abc"))
	fmt.Println(IntToBytes(50))
	fmt.Println(BytesToString([]byte{97,98,99}))
	fmt.Println(BytesToInt([]byte{0,0,0,50}))*/

	fmt.Println(BytesToInt64(StringToBytes("abc")))
	//big1 := new(big.Int).SetInt64(BytesToInt64(StringToBytes("abc")))
	//fmt.Println(reflect.TypeOf(big1))
}
