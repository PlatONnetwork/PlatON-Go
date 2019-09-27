package reward

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	HistoryIncreasePrefix    = []byte("RewardHistory")
	LastYearEndBalancePrefix = []byte("RewardBalance")
)

// GetHistoryIncreaseKey used for search the balance of reward pool at last year
func GetHistoryIncreaseKey(year uint32) []byte {
	return append(HistoryIncreasePrefix, common.Uint32ToBytes(year)...)
}

//
func HistoryBalancePrefix(year uint32) []byte {
	return append(LastYearEndBalancePrefix, common.Uint32ToBytes(year)...)
}
