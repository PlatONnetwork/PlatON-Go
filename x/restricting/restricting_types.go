package restricting

import "math/big"

type RestrictingInfo struct {
	Balance     *big.Int `json:"balance"` // balance representation all locked amount
	Debt        *big.Int `json:"debt"`    // debt representation will released amount. Positive numbers can be used instead of release, 0 means no release, negative numbers indicate not enough to release
	ReleaseList []uint64 `json:"list"`    // releaseList representation which epoch will release restricting
}


type RestrictingPlan struct {
	Epoch   uint64  `json:"epoch"`			// epoch representation of the released epoch at the target blockNumber
	Amount	*big.Int `json:"amount"`		// amount representation of the released amount
}
