package cbfttypes

import (
	"Platon-go/common"
	"math/big"
	"Platon-go/core/state"
	"Platon-go/core/types"
)

// modify by platon
// Block's Signature info
type BlockSignature struct {
	Hash        common.Hash
	Number      *big.Int
	Signature   *common.BlockConfirmSign
}

// modify by platon
type CbftResult struct {
	Block       		*types.Block
	Receipts    		types.Receipts
	State       		*state.StateDB
	BlockConfirmSigns 	[]*common.BlockConfirmSign
}