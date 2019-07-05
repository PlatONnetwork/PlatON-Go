package restriting

import "math/big"

type RestrictingPlan struct {
	Epoch   uint64  `json:"epoch"`			// epoch representation of the released epoch at the target blockNumber
	Amount	*big.Int `json:"amount"`		// amount representation of the released amount
}
