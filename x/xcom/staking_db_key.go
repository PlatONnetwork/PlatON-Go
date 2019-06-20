package xcom

import (
	"crypto/ecdsa"
	"fmt"
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


func TallyPowerKey(shares *big.Int, stakeBlockNum uint64, stakeTxIndex uint32) []byte {

	priority := new(big.Int).Sub(math.MaxBig256, shares)
	prio := priority.String()
	num := fmt.Sprint(stakeBlockNum)
	index := fmt.Sprint(stakeTxIndex)
	return append(CanPowerKeyPrefix, append([]byte(prio), append([]byte(num), []byte(index)...)...)...)
}



func GetUnStakeCountKey (epoch uint64) []byte {
	epochStr := strconv.Itoa(int(epoch))
	return  append(UnStakeCountKey, []byte(epochStr)...)
}


