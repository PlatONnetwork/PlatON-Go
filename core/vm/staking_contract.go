package vm



type stakingContract struct {
	Contract *Contract
	Evm      *EVM
}



func (stkc *stakingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (stkc *stakingContract) Run(input []byte) ([]byte, error) {
	return nil, nil
}