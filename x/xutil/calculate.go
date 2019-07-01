package xutil

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
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
	size := 12 * 30 * 4 * xcom.ConsensusSize * xcom.EpochSize
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