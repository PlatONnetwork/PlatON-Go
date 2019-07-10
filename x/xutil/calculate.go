package xutil

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
)

// calculate the Epoch number by blocknumber
func CalculateEpoch(blockNumber uint64) uint64 {

	// block counts of per epoch
	size := xcom.ConsensusSize*xcom.EpochSize


	var epoch uint64
	div := blockNumber / size
	mod := blockNumber % size

	switch  {
	// first epoch
	case (div == 0 && mod == 0) || (div == 0 && mod > 0 && mod < size):
		epoch = 1
	case div > 0 && mod == 0:
		epoch = div
	case div > 0 && mod > 0 && mod < size:
		epoch = div + 1
	}

	return epoch
}


// calculate the Consensus number by blocknumber
func CalculateRound (blockNumber uint64) uint64 {

	size := xcom.ConsensusSize

	var round uint64
	div := blockNumber / size
	mod := blockNumber % size
	switch  {
	// first consensus round
	case (div == 0 && mod == 0) || (div == 0 && mod > 0 && mod < size):
		round = 1
	case div > 0 && mod == 0:
		round = div
	case div > 0 && mod > 0 && mod < size:
		round = div + 1
	}

	return round
}


func CalculateYear (blockNumber uint64) uint64 {
	// !!!
	// epochs := EpochsPerYear()
	// !!!
	epochs := uint64(1440)
	size := epochs * xcom.EpochSize * xcom.ConsensusSize

	var year uint64

	div := blockNumber / size
	mod := blockNumber % size

	switch {
	case mod == 0:
		year = div
	case mod > 0:
		year = div +1
	}

	return year
}

func IsElection(blockNumber uint64) bool {
	tmp := blockNumber + xcom.ElectionDistance
	mod := tmp % xcom.ConsensusSize
	return mod == 0
}


func IsSwitch (blockNumber uint64) bool {
	mod := blockNumber % xcom.ConsensusSize
	return mod == 0
}

func IsSettlementPeriod (blockNumber uint64) bool {
	// block counts of per epoch
	size := xcom.ConsensusSize*xcom.EpochSize
	mod := blockNumber % size
	return mod == 0
}


func IsYearEnd (blockNumber uint64) bool {
	// Calculate epoch
	eh := uint64(6)
	L := uint64(1)
	u := uint64(25)
	vn := uint64(10)
	epoch := eh*3600/(L*u*vn)*xcom.ConsensusSize

	size := uint64(365)*24*3600/(L*epoch)*epoch
	return blockNumber > 0 && blockNumber % size == 0
}

func NodeId2Addr (nodeId discover.NodeID) (common.Address, error) {

	if pk, err := nodeId.Pubkey(); nil != err {
		return common.ZeroAddr, err
	} else {
		return crypto.PubkeyToAddress(*pk), nil
	}
}


func IsWorker(extra []byte) bool {
	return len(extra[32:]) >= common.ExtraSeal && bytes.Equal(extra[32:97], make([]byte, common.ExtraSeal))
}



func CheckStakeThreshold(stake *big.Int) bool {
	return stake.Cmp(xcom.StakeThreshold) >= 0
}

func CheckDelegateThreshold(delegate *big.Int) bool {
	return delegate.Cmp(xcom.DelegateThreshold) >= 0
}

// The ProcessVersion: Major.Minor.Patch eg. 1.1.0
// Calculate the LargeVersion
// eg: 1.1.0 ==> 1.1
func CalcLargeVersion (processVersion uint32) uint32 {
	return processVersion>>8
}