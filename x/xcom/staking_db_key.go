package xcom

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
	"strconv"
)

const (

	CandidatePrefixStr = "Can"

	CanPowerPrefixStr 	= "Power"

	UnStakeCountKeyStr 	= "UnStakeCount"

	UnStakeItemKeyStr = "UnStakeItem"

	DelegatePrefixStr 	= "Del"

	UnDelegateCountKeyStr = "UnDelCount"

	UnDelegateItemKeyStr = "UnDelItem"


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

	UnDelegateCountKey = []byte(UnDelegateCountKeyStr)

	UnDelegateItemKey = []byte(UnDelegateItemKeyStr)

	EpochValidatorKey = []byte(EpochValidatorKeyStr)

	PreRoundValidatorKey = []byte(PreRoundValidatorKeyStr)

	CurRoundValidatorKey = []byte(CurRoundValidatorKeyStr)

	NextRoundValidatorKey = []byte(NextRoundValidatorKeyStr)





)



//////// TODO

func CandidateKeyByNodeId(nodeId discover.NodeID) ([]byte, error) {

	if pk, err := nodeId.Pubkey(); nil != err {
		return nil, err
	}else {
		addr := crypto.PubkeyToAddress(*pk)
		return append(CandidateKeyPrefix, addr.Bytes()...), nil
	}
}

func CandidateKeyByPubKey(p ecdsa.PublicKey) []byte {
	addr :=  crypto.PubkeyToAddress(p)
	return append(CandidateKeyPrefix, addr.Bytes()...)
}

func CandidateKeyByAddr (addr common.Address) []byte {
	return append(CandidateKeyPrefix, addr.Bytes()...)
}

func CandidateKeyBySuffix (addr []byte) []byte {
	return append(CandidateKeyPrefix, addr...)
}


func TallyPowerKey(shares *big.Int, stakeBlockNum, stakeTxIndex  int) []byte {

	priority := new(big.Int).Sub(math.MaxBig256, shares)
	prio := priority.String()
	num := strconv.Itoa(stakeBlockNum)
	index := strconv.Itoa(stakeTxIndex)
	return append(CanPowerKeyPrefix, append([]byte(prio), append([]byte(num), []byte(index)...)...)...)
}



func GetUnStakeCountKey (epoch int) []byte {
	epochStr := strconv.Itoa(epoch)
	return  append(UnStakeCountKey, []byte(epochStr)...)
}


func GetUnStakeItemKey (epoch, index int) []byte {
	epochStr := strconv.Itoa(epoch)
	indexStr := strconv.Itoa(index)
	return append(UnStakeItemKey, append([]byte(epochStr), []byte(indexStr)...)...)
}


func GetDelegateKey(delAddr common.Address, nodeId discover.NodeID, stakeBlockNumber int) []byte {
	num := strconv.Itoa(stakeBlockNumber)
	return append(DelegateKeyPrefix, append(delAddr.Bytes(), append(nodeId.Bytes(), []byte(num)...)...)...)
}

func GetDelegateKeyBySuffix(suffix []byte) []byte {
	return append(DelegateKeyPrefix, suffix...)
}

func GetUnDelegateCountKey (epoch int) []byte {
	epochStr := strconv.Itoa(epoch)
	return  append(UnDelegateCountKey, []byte(epochStr)...)
}

func GetUnDelegateItemKey (epoch, index int) []byte {
	epochStr := strconv.Itoa(epoch)
	indexStr := strconv.Itoa(index)
	return append(UnDelegateItemKey, append([]byte(epochStr), []byte(indexStr)...)...)
}

func GetEpochValidatorKey () []byte {
	return EpochValidatorKey
}

func GetPreRoundValidatorKey () []byte {
	return PreRoundValidatorKey
}

func GetCurRoundValidatorKey () []byte {
	return CurRoundValidatorKey
}

func GetNextRoundValidatorKey () []byte {
	return NextRoundValidatorKey
}

