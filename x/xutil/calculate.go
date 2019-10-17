package xutil

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

func NodeId2Addr(nodeId discover.NodeID) (common.Address, error) {
	if pk, err := nodeId.Pubkey(); nil != err {
		return common.ZeroAddr, err
	} else {
		return crypto.PubkeyToAddress(*pk), nil
	}
}

// The ProgramVersion: Major.Minor.Patch eg. 1.1.0
// Calculate the LargeVersion
// eg: 1.1.x ==> 1.1.0
func CalcVersion(programVersion uint32) uint32 {
	programVersion = programVersion >> 8
	return programVersion << 8
}

func IsWorker(extra []byte) bool {
	return len(extra[32:]) >= common.ExtraSeal && bytes.Equal(extra[32:97], make([]byte, common.ExtraSeal))
}

func CheckStakeThreshold(stake *big.Int) bool {
	return stake.Cmp(xcom.StakeThreshold()) >= 0
}

func CheckMinimumThreshold(balance *big.Int) bool {
	return balance.Cmp(xcom.MinimumThreshold()) >= 0
}

// eg. 65536 => 1.0.0
func ProgramVersion2Str(programVersion uint32) string {
	if programVersion == 0 {
		return "0.0.0"
	}
	major := programVersion << 8
	major = major >> 24

	minor := programVersion << 16
	minor = minor >> 24

	patch := programVersion << 24
	patch = patch >> 24

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

// ConsensusSize returns how many blocks per consensus round.
func ConsensusSize() uint64 {
	return xcom.BlocksWillCreate() * xcom.ConsValidatorNum()
}

// EpochSize returns how many consensus rounds per epoch.
func EpochSize() uint64 {
	consensusSize := ConsensusSize()
	em := xcom.ExpectedMinutes()
	i := xcom.Interval()

	epochSize := em * 60 / (i * consensusSize)
	return epochSize
}

// EpochsPerYear returns how many epochs per year
func EpochsPerYear() uint64 {
	epochBlocks := CalcBlocksEachEpoch()
	i := xcom.Interval()
	return xcom.AdditionalCycleTime() * 60 / (i * epochBlocks)
}

// CalcBlocksEachEpoch return how many blocks per epoch
func CalcBlocksEachEpoch() uint64 {
	return ConsensusSize() * EpochSize()
}

// calculate returns how many blocks per year.
func CalcBlocksEachYear() uint64 {
	return EpochsPerYear() * CalcBlocksEachEpoch()
}

// calculate the Epoch number by blockNumber
func CalculateEpoch(blockNumber uint64) uint64 {
	size := CalcBlocksEachEpoch()

	var epoch uint64
	div := blockNumber / size
	mod := blockNumber % size

	switch {
	// first epoch
	case div == 0:
		epoch = 1
	case div > 0 && mod == 0:
		epoch = div
	case div > 0 && mod > 0:
		epoch = div + 1
	}

	return epoch
}

// calculate the Consensus number by blockNumber
func CalculateRound(blockNumber uint64) uint64 {
	size := ConsensusSize()

	var round uint64
	div := blockNumber / size
	mod := blockNumber % size
	switch {
	// first consensus round
	case div == 0:
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
func CalculateYear(blockNumber uint64) uint32 {
	size := CalcBlocksEachYear()

	div := blockNumber / size
	mod := blockNumber % size

	if mod == 0 {
		return uint32(div)
	} else {
		return uint32(div + 1)
	}
}

func CalculateLastYear(blockNumber uint64) uint32 {
	thisYear := CalculateYear(blockNumber)

	if thisYear == 0 {
		return 0
	} else {
		return thisYear - 1
	}
}

func InNodeIDList(nodeID discover.NodeID, nodeIDList []discover.NodeID) bool {
	for _, v := range nodeIDList {
		if nodeID == v {
			return true
		}
	}
	return false
}

func InHashList(hash common.Hash, hashList []common.Hash) bool {
	for _, v := range hashList {
		if hash == v {
			return true
		}
	}
	return false
}

// end-voting-block = the end block of a consensus period - electionDistance, end-voting-block must be a Consensus Election block
func CalEndVotingBlock(blockNumber uint64, endVotingRounds uint64) uint64 {
	electionDistance := xcom.ElectionDistance()
	consensusSize := ConsensusSize()
	return blockNumber + consensusSize - blockNumber%consensusSize + endVotingRounds*consensusSize - electionDistance
}

// active-block = the begin of a consensus period, so, It's possible that active-block is the begin block of a settlement epoch
func CalActiveBlock(endVotingBlock uint64) uint64 {
	return endVotingBlock + xcom.ElectionDistance() + (xcom.VersionProposalActive_ConsensusRounds()-1)*ConsensusSize() + 1
}

func IsSpecialBlock(blockNumber uint64) bool {
	if IsElection(blockNumber) || IsEndOfEpoch(blockNumber) || IsYearEnd(blockNumber) || IsBeginOfEpoch(blockNumber) || IsBeginOfConsensus(blockNumber) {
		return true
	}
	return false
}

func IsYearBegin(blockNumber uint64) bool {
	size := CalcBlocksEachYear()
	mod := blockNumber % size
	return mod == 1
}

// IsBeginOfEpoch returns true if current block is the first block of a Epoch
func IsBeginOfEpoch(blockNumber uint64) bool {
	size := CalcBlocksEachEpoch()
	mod := blockNumber % size
	return mod == 1
}

// IsBeginOfConsensus returns true if current block is the first block of a Consensus Cycle
func IsBeginOfConsensus(blockNumber uint64) bool {
	size := ConsensusSize()
	mod := blockNumber % size
	return mod == 1
}

func IsYearEnd(blockNumber uint64) bool {
	size := CalcBlocksEachYear()
	mod := blockNumber % size
	return mod == 0
}

func IsEndOfEpoch(blockNumber uint64) bool {
	size := CalcBlocksEachEpoch()
	mod := blockNumber % size
	return mod == 0
}

func IsEndOfConsensus(blockNumber uint64) bool {
	size := ConsensusSize()
	mod := blockNumber % size
	return mod == 0
}

func IsElection(blockNumber uint64) bool {
	tmp := blockNumber + xcom.ElectionDistance()
	mod := tmp % ConsensusSize()
	return mod == 0
}
