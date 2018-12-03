package cbfttypes

import (
	"Platon-go/common"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"math/big"
)

// Block's Signature info
type BlockSignature struct {
	SignHash  common.Hash //签名hash，header[0:32]
	Hash      common.Hash //块hash，header[:]
	Number    *big.Int
	Signature *common.BlockConfirmSign
	ParentHash common.Hash
}

type CbftResult struct {
	Block *types.Block
	Receipts          types.Receipts
	State             *state.StateDB
	BlockConfirmSigns []*common.BlockConfirmSign
}
