package xcom

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
)

const (
	CandidatePrefixStr       = "Can"
	CanPowerPrefixStr        = "Power"
	UnStakeCountKeyStr       = "UnStakeCount"
	UnStakeItemKeyStr        = "UnStakeItem"
	DelegatePrefixStr        = "Del"
	UnDelegateCountKeyStr    = "UnDelCount"
	UnDelegateItemKeyStr     = "UnDelItem"
	EpochValidatorKeyStr     = "EpochValidator"
	PreRoundValidatorKeyStr  = "PreRoundValidator"
	CurRoundValidatorKeyStr  = "CurRoundValidator"
	NextRoundValidatorKeyStr = "NextRoundValidator"
	PPOSHASHStr              = "PPOS_HASH"
)

var (
	CandidateKeyPrefix    = []byte(CandidatePrefixStr)
	CanPowerKeyPrefix     = []byte(CanPowerPrefixStr)
	UnStakeCountKey       = []byte(UnStakeCountKeyStr)
	UnStakeItemKey        = []byte(UnStakeItemKeyStr)
	DelegateKeyPrefix     = []byte(DelegatePrefixStr)
	UnDelegateCountKey    = []byte(UnDelegateCountKeyStr)
	UnDelegateItemKey     = []byte(UnDelegateItemKeyStr)
	EpochValidatorKey     = []byte(EpochValidatorKeyStr)
	PreRoundValidatorKey  = []byte(PreRoundValidatorKeyStr)
	CurRoundValidatorKey  = []byte(CurRoundValidatorKeyStr)
	NextRoundValidatorKey = []byte(NextRoundValidatorKeyStr)
	PPOSHASHKey           = []byte(PPOSHASHStr)
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

// need to add ProcessVersion
func TallyPowerKey(shares *big.Int, stakeBlockNum uint64, stakeTxIndex, processVersion uint32) []byte {
	version := common.Uint32ToBytes(processVersion)
	priority := new(big.Int).Sub(math.MaxBig256, shares)
	prio := priority.String()
	num := common.Uint64ToBytes(stakeBlockNum)
	txIndex := common.Uint32ToBytes(stakeTxIndex)
	return append(version, append(CanPowerKeyPrefix, append([]byte(prio),
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

func GetEpochValidatorKey() []byte {
	return EpochValidatorKey
}

func GetPreRoundValidatorKey() []byte {
	return PreRoundValidatorKey
}

func GetCurRoundValidatorKey() []byte {
	return CurRoundValidatorKey
}

func GetNextRoundValidatorKey() []byte {
	return NextRoundValidatorKey
}


func GetPPOSHASHKey() []byte {
	return PPOSHASHKey
}