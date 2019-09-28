package staking

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (
	CandidatePrefixStr       = "Can"
	CanPowerPrefixStr        = "Power"
	UnStakeCountKeyStr       = "UnStakeCount"
	UnStakeItemKeyStr        = "UnStakeItem"
	DelegatePrefixStr        = "Del"
	UnDelegateCountKeyStr    = "UnDelCount"
	UnDelegateItemKeyStr     = "UnDelItem"
	EpochIndexKeyStr         = "EpochIndex"
	EpochValArrPrefixStr     = "EpochValArr"
	RoundIndexKeyStr         = "RoundIndex"
	RoundValArrPrefixStr     = "RoundValArr"
	AccountStakeRcPrefixStr  = "AccStakeRc"
	PPOSHASHStr              = "PPOS_HASH"
	RoundValAddrArrPrefixStr = "RoundValAddrArr"
)

var (
	CandidateKeyPrefix    = []byte(CandidatePrefixStr)
	CanPowerKeyPrefix     = []byte(CanPowerPrefixStr)
	UnStakeCountKey       = []byte(UnStakeCountKeyStr)
	UnStakeItemKey        = []byte(UnStakeItemKeyStr)
	DelegateKeyPrefix     = []byte(DelegatePrefixStr)
	UnDelegateCountKey    = []byte(UnDelegateCountKeyStr)
	UnDelegateItemKey     = []byte(UnDelegateItemKeyStr)
	EpochIndexKey         = []byte(EpochIndexKeyStr)
	EpochValArrPrefix     = []byte(EpochValArrPrefixStr)
	RoundIndexKey         = []byte(RoundIndexKeyStr)
	RoundValArrPrefix     = []byte(RoundValArrPrefixStr)
	AccountStakeRcPrefix  = []byte(AccountStakeRcPrefixStr)
	PPOSHASHKey           = []byte(PPOSHASHStr)
	b104Len               = len(math.MaxBig104.Bytes())
	RoundValAddrArrPrefix = []byte(RoundValAddrArrPrefixStr)
)

func CandidateKeyByNodeId(nodeId discover.NodeID) ([]byte, error) {

	if pk, err := nodeId.Pubkey(); nil != err {
		return nil, err
	} else {
		addr := crypto.PubkeyToAddress(*pk)
		return append(CandidateKeyPrefix, addr.Bytes()...), nil
	}
}

func CandidateKeyByPubKey(p ecdsa.PublicKey) []byte {
	addr := crypto.PubkeyToAddress(p)
	return append(CandidateKeyPrefix, addr.Bytes()...)
}

func CandidateKeyByAddr(addr common.Address) []byte {
	return append(CandidateKeyPrefix, addr.Bytes()...)
}

func CandidateKeyBySuffix(addr []byte) []byte {
	return append(CandidateKeyPrefix, addr...)
}

// need to add ProgramVersion
func TallyPowerKey(shares *big.Int, stakeBlockNum uint64, stakeTxIndex, programVersion uint32) []byte {

	subVersion := math.MaxInt32 - programVersion

	sortVersion := common.Uint32ToBytes(subVersion)
	priority := new(big.Int).Sub(math.MaxBig104, shares)

	zeros := make([]byte, b104Len)
	prio := append(zeros, priority.Bytes()...)

	num := common.Uint64ToBytes(stakeBlockNum)
	txIndex := common.Uint32ToBytes(stakeTxIndex)
	return append(CanPowerKeyPrefix, append(sortVersion, append(prio,
		append(num, txIndex...)...)...)...)
}

func GetUnStakeCountKey(epoch uint64) []byte {
	return append(UnStakeCountKey, common.Uint64ToBytes(epoch)...)
}

func GetUnStakeItemKey(epoch, index uint64) []byte {
	return append(UnStakeItemKey, append(common.Uint64ToBytes(epoch), common.Uint64ToBytes(index)...)...)
}

func GetDelegateKey(delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber uint64) []byte {
	return append(DelegateKeyPrefix, append(delAddr.Bytes(), append(nodeId.Bytes(),
		common.Uint64ToBytes(stakeBlockNumber)...)...)...)
}

func GetDelegateKeyBySuffix(suffix []byte) []byte {
	return append(DelegateKeyPrefix, suffix...)
}

func GetUnDelegateCountKey(epoch uint64) []byte {
	return append(UnDelegateCountKey, common.Uint64ToBytes(epoch)...)
}

func GetUnDelegateItemKey(epoch, index uint64) []byte {
	return append(UnDelegateItemKey, append(common.Uint64ToBytes(epoch), common.Uint64ToBytes(index)...)...)
}

//func GetEpochValidatorKey() []byte {
//	return EpochValidatorKey
//}
//
//func GetPreRoundValidatorKey() []byte {
//	return PreRoundValidatorKey
//}
//
//func GetCurRoundValidatorKey() []byte {
//	return CurRoundValidatorKey
//}
//
//func GetNextRoundValidatorKey() []byte {
//	return NextRoundValidatorKey
//}

func GetEpochIndexKey() []byte {
	return EpochIndexKey
}

func GetEpochValArrKey(start, end uint64) []byte {
	startByte := common.Uint64ToBytes(start)
	endByte := common.Uint64ToBytes(end)
	return append(EpochValArrPrefix, append(startByte, endByte...)...)
}

func GetRoundIndexKey() []byte {
	return RoundIndexKey
}

func GetRoundValArrKey(start, end uint64) []byte {
	startByte := common.Uint64ToBytes(start)
	endByte := common.Uint64ToBytes(end)
	return append(RoundValArrPrefix, append(startByte, endByte...)...)
}

func GetAccountStakeRcKey(addr common.Address) []byte {
	return append(AccountStakeRcPrefix, addr.Bytes()...)
}

func GetPPOSHASHKey() []byte {
	return PPOSHASHKey
}

func GetRoundValAddrArrKey(round uint64) []byte {
	return append(RoundValAddrArrPrefix, common.Uint64ToBytes(round)...)
}
