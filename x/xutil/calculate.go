// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package xutil

import (
	"bytes"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

func NodeId2Addr(nodeId discover.NodeID) (common.NodeAddress, error) {
	if pk, err := nodeId.Pubkey(); nil != err {
		return common.ZeroNodeAddr, err
	} else {
		return common.NodeAddress(crypto.PubkeyToAddress(*pk)), nil
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
	return xcom.ConsensusSize()
}

// EpochSize returns how many consensus rounds per epoch.
func EpochSize() uint64 {
	return xcom.EpochSize()
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

func EstimateConsensusRoundsForGov(seconds uint64) uint64 {
	//v0.7.5, hard code 1 second for block interval for estimating.
	blockInterval := uint64(1)
	return seconds / (blockInterval * ConsensusSize())
}

func EstimateEndVotingBlockForParaProposal(blockNumber uint64, seconds uint64) uint64 {
	consensusSize := ConsensusSize()
	epochMaxDuration := xcom.MaxEpochMinutes() //minutes
	//estimate how many consensus rounds in a epoch.
	consensusRoundsEachEpoch := epochMaxDuration * 60 / (xcom.Interval() * consensusSize)
	blocksEachEpoch := consensusRoundsEachEpoch * consensusSize

	//v0.7.5, hard code 1 second for block interval for estimating.
	blockInterval := uint64(1)
	durationEachEpoch := blocksEachEpoch * blockInterval

	epochRounds := seconds / durationEachEpoch
	return blockNumber + blocksEachEpoch - blockNumber%blocksEachEpoch + epochRounds*blocksEachEpoch
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

// active-block = the begin of a consensus period, so, it is possible that active-block also is the begin of a epoch.
func CalActiveBlock(endVotingBlock uint64) uint64 {
	//return endVotingBlock + xcom.ElectionDistance() + (xcom.VersionProposalActive_ConsensusRounds()-1)*ConsensusSize() + 1
	return endVotingBlock + xcom.ElectionDistance() + 1
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

func IsEndOfEpoch(blockNumber uint64) bool {
	size := CalcBlocksEachEpoch()
	mod := blockNumber % size
	return mod == 0
}

func IsElection(blockNumber uint64) bool {
	tmp := blockNumber + xcom.ElectionDistance()
	mod := tmp % ConsensusSize()
	return mod == 0
}
