package restricting

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

// for genesis and plugin test
type RestrictingInfo struct {
	NeedRelease     *big.Int
	StakingAmount   *big.Int
	CachePlanAmount *big.Int
	//	SlashingMount *big.Int

	//Balance     *big.Int // Balance representation all locked amount
	//Debt        *big.Int // Debt representation will released amount.
	//DebtSymbol  bool     // Debt is owed to release in the past while symbol is true, else Debt can be used instead of release
	ReleaseList []uint64 // ReleaseList representation which epoch will release restricting
}

func (r *RestrictingInfo) RemoveEpoch(epoch uint64) {
	for i, target := range r.ReleaseList {
		if target == epoch {
			r.ReleaseList = append(r.ReleaseList[:i], r.ReleaseList[i+1:]...)
			break
		}
	}
}

// for contract, plugin test, byte util
type RestrictingPlan struct {
	Epoch  uint64   `json:"epoch"`  // epoch representation of the released epoch at the target blockNumber
	Amount *big.Int `json:"amount"` // amount representation of the released amount
}

// for plugin test
type ReleaseAmountInfo struct {
	Height uint64   `json:"blockNumber"` // blockNumber representation of the block number at the released epoch
	Amount *big.Int `json:"amount"`      // amount representation of the released amount
}

// for plugin test
type Result struct {
	Balance *big.Int            `json:"balance"`
	Debt    *big.Int            `json:"debt"`
	Entry   []ReleaseAmountInfo `json:"plans"`
	Pledge  *big.Int            `json:"Pledge"`
}

type EpochInfo struct {
	Account common.Address
	Amount  *big.Int
}
