package restricting

import "math/big"

// for genesis and plugin test
type RestrictingInfo struct {
	Balance     *big.Int // Balance representation all locked amount
	Debt        *big.Int // Debt representation will released amount.
	DebtSymbol  bool     // Debt is owed to release in the past while symbol is true, else Debt can be used instead of release
	ReleaseList []uint64 // ReleaseList representation which epoch will release restricting
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
	Balance *big.Int `json:"balance"`
	Debt    *big.Int `json:"debt"`
	Symbol  bool     `json:"symbol"`
	Entry   []ReleaseAmountInfo
}
