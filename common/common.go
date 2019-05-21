package common

import "math/big"

const (
	BaseElection = 230
	//
	BaseSwitchWitness = 250
	//
	BaseAddNextPeers = 230
	//
	BaseIrrCount = 1

	FirstRound = 1


	/*BaseElection = 50
	//
	BaseSwitchWitness = 60
	//
	BaseAddNextPeers = 50
	//
	BaseIrrCount = 20*/

)

var (
	//ppos
	InitAmount      = new(big.Int).SetBytes([]byte{3, 59, 46, 60, 159, 208, 128, 60, 232, 0, 0, 0}) //1,000,000,000 * 10^18 ADP -> 0x033b2e3c9fd0803ce8000000 -> [3 59 46 60 159 208 128 60 232 0 0 0]
	FirstYearReward = new(big.Int).SetBytes([]byte{20, 173, 244, 183, 50, 3, 52, 185, 0, 0, 0})     //25,000,000 * 10^18 ADP -> 0x14adf4b7320334b9000000-> [20 173 244 183 50 3 52 185 0 0 0]
	YearBlocks      = new(big.Int).SetInt64(31536000)                                               //24*3600*365
	Rate            = big.NewInt(1025)                                                              //Actual value Rate/Base
	Base            = big.NewInt(1000)
	FeeBase         = big.NewInt(10000)
)
