package xutil

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
)

var (
	eh = uint32(6)   // expected hours every settle epoch
	l = uint32(1)    // time of creating a new block in seconds
	u = uint32(25)   // the consensus validators count
	vn = uint32(10)  // each validator will seal blocks per view

)

func GetConsensusSize() uint32 {
	return uint32(u * vn)
}

func GetBlocksPerEpoch() uint32 {
	consensusSize := GetConsensusSize()
	expected := eh * 3600/(l * consensusSize) * consensusSize
	return expected
}

func GetExpectedEpochsPerYear() uint32 {
	blocks := GetBlocksPerEpoch()
	return 365 * 24 * 3600 / (l * blocks)
}

// calculate the Epoch number by blocknumber
func CalculateEpoch(blockNumber uint64) uint64 {
	// block counts of per epoch
	size := xcom.ConsensusSize()*xcom.EpochSize()

	var epoch uint64
	div := blockNumber / size
	mod := blockNumber % size

	switch  {
	// first epoch
	case (div == 0 && mod == 0) || (div == 0 && mod > 0):
		epoch = 1
	case div > 0 && mod == 0:
		epoch = div
	case div > 0 && mod > 0:
		epoch = div + 1
	}

	return epoch
}


// calculate the Consensus number by blocknumber
func CalculateRound (blockNumber uint64) uint64 {
	size := xcom.ConsensusSize()

	var round uint64
	div := blockNumber / size
	mod := blockNumber % size
	switch  {
	// first consensus round
	case (div == 0 && mod == 0) || (div == 0 && mod > 0):
		round = 1
	case div > 0 && mod == 0:
		round = div
	case div > 0 && mod > 0:
		round = div + 1
	}

	return round
}

// calculate the year by blockNumber.
// (V.0.1) If blockNumber eqs 0, year eqs 0 too, else rounded up the result of
// the blockNumber divided by the expected number of blocks per year
func CalculateYear (blockNumber uint64) uint64 {
	// size is expected new blocks per year
	size := GetExpectedEpochsPerYear() * GetBlocksPerEpoch()

	div := blockNumber / uint64(size)
	mod := blockNumber % uint64(size)

	if mod == 0 {
		return div
	} else {
		return div + 1
	}
}

func IsElection(blockNumber uint64) bool {
	tmp := blockNumber + xcom.ElectionDistance()
	mod := tmp % xcom.ConsensusSize()
	return mod == 0
}


func IsSwitch (blockNumber uint64) bool {
	mod := blockNumber % xcom.ConsensusSize()
	return mod == 0
}

func IsSettlementPeriod (blockNumber uint64) bool {
	// block counts of per epoch
	size := GetBlocksPerEpoch()
	mod := blockNumber % uint64(size)
	return mod == 0
}

func IsYearEnd (blockNumber uint64) bool {
	size := GetBlocksPerEpoch() * GetExpectedEpochsPerYear()
	return blockNumber > 0 && blockNumber % uint64(size) == 0
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
	return stake.Cmp(xcom.StakeThreshold()) >= 0
}

func CheckDelegateThreshold(delegate *big.Int) bool {
	return delegate.Cmp(xcom.DelegateThreshold()) >= 0
}

// The ProcessVersion: Major.Minor.Patch eg. 1.1.0
// Calculate the LargeVersion
// eg: 1.1.0 ==> 1.1
func CalcVersion (processVersion uint32) uint32 {
	return processVersion>>8
}

// eg. 65536 => 1.0.0
func ProcessVerion2Str (processVersion uint32) string {
	major := processVersion<<8
	major = major>>24

	minor := processVersion<<16
	minor = minor>>24

	patch := processVersion<<24
	patch = patch>>24

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}