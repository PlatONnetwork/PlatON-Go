package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/life/utils"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
)

type ContractRefSelf struct {
}

func (c ContractRefSelf) Address() common.Address {
	return common.BigToAddress(big.NewInt(66666))
}

type ContractRefCaller struct {
}

func (c ContractRefCaller) Address() common.Address {
	return common.BigToAddress(big.NewInt(77777))
}

func genInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, utils.Int64ToBytes(1))
	input = append(input, []byte("transfer"))
	input = append(input, []byte("0x0000000000000000000000000000000000000001"))
	input = append(input, []byte("0x0000000000000000000000000000000000000002"))
	input = append(input, utils.Int64ToBytes(100))

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	if err != nil {
		fmt.Println("geninput fail.", err)
	}
	return buffer.Bytes()
}

type StateDBTest struct {
}

func bytes2int64(byt []byte) int64 {
	bytesBuf := bytes.NewBuffer(byt)
	var tmp int64
	binary.Read(bytesBuf, binary.BigEndian, &tmp)
	return tmp
}
