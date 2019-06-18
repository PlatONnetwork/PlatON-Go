package vm

import (
	"errors"
	xcommon "github.com/PlatONnetwork/PlatON-Go/x/common"
	xcore "github.com/PlatONnetwork/PlatON-Go/x/core"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)



type stakingContract struct {
	plugin 		*plugin.StakingPlugin
	Contract 	*Contract
	Evm      	*EVM
}






var (
	DecodeStakingTxDataErr = errors.New("decode staking tx data is err")
)


func (stkc *stakingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (stkc *stakingContract) Run(input []byte) ([]byte, error) {

	reactor := xcore.GetReactorInstance()
	plugin := reactor.GetPlugin(xcommon.SlashingRule)
	plugin.GetVal(stkc.Evm.StateDB)
	return nil, nil
}

//func execute (input []byte) (ret []byte, err error) {
//
//	var args [][]byte
//	if err := rlp.Decode(bytes.NewReader(input), &args); nil != err {
//		err = DecodeStakingTxDataErr
//		return
//	}
//
//
//
//	return nil, nil
//}


