package vm

import (
	"Platon-go/params"
)

type ticketContract struct {
	contract *Contract
	evm *EVM
}

func (t *ticketContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (t *ticketContract) Run(input []byte) ([]byte, error) {
	var command = map[string] interface{}{
		// 接口列表
	}
	return execute(input, command)
}