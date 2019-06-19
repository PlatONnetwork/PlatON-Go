package xcom

import "github.com/PlatONnetwork/PlatON-Go/p2p/discover"

const (

	CandidatePrefixStr = "Can"

	CanPowerPrefixStr 	= "Power"

	UnStakeCountKeyStr 	= "UnStakeCount"

	UnStakeItemKeyStr = "UnStakeItem"

	DelegatePrefixStr 	= "Del"

	UndelegateCountKeyStr = "UnDelCount"

	UndelegateItemKeyStr = "UnDelItem"


	EpochValidatorKeyStr = "EpochValidator"

	PreRoundValidatorKeyStr = "PreRoundValidator"

	CurRoundValidatorKeyStr = "CurRoundValidator"

	NextRoundValidatorKeyStr = "NextRoundValidator"

)

var (
	CandidateKeyPrefix = []byte(CandidatePrefixStr)

	CanPowerKeyPrefix = []byte(CanPowerPrefixStr)

	UnStakeCountKey = []byte(UnStakeCountKeyStr)

	UnStakeItemKey = []byte(UnStakeItemKeyStr)

	DelegateKeyPrefix = []byte(DelegatePrefixStr)

	UndelegateCountKey = []byte(UndelegateCountKeyStr)

	UndelegateItemKey = []byte(UndelegateItemKeyStr)

	EpochValidatorKey = []byte(EpochValidatorKeyStr)

	PreRoundValidatorKey = []byte(PreRoundValidatorKeyStr)

	CurRoundValidatorKey = []byte(CurRoundValidatorKeyStr)

	NextRoundValidatorKey = []byte(NextRoundValidatorKeyStr)

)



//////// TODO

func CandidateKey(nodeId discover.NodeID) []byte {
	return append(CandidateKeyPrefix, nodeId.Bytes()...)
}

