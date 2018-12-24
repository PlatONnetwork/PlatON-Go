package cbfttypes

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"math/big"
)

// Block's Signature info
type BlockSignature struct {
	SignHash  common.Hash // Signature hash，header[0:32]
	Hash      common.Hash // Block hash，header[:]
	Number    *big.Int
	Signature *common.BlockConfirmSign
}

type CbftResult struct {
	Block *types.Block
	//Receipts          types.Receipts
	//State             *state.StateDB
	BlockConfirmSigns []*common.BlockConfirmSign
}
