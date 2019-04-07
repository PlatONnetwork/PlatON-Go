package types

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
)

type CanConditions map[discover.NodeID]*big.Int

type KindCanQueue []CandidateQueue

type CandidateQueue []*Candidate

func (queue CandidateQueue) DeepCopy() CandidateQueue {
	copyCandidateQueue := make(CandidateQueue, len(queue))
	if len(queue) == 0 {
		return copyCandidateQueue
	}
	for i, can := range queue {
		deposit, _ := new(big.Int).SetString(can.Deposit.String(), 10)
		canCopy := &Candidate{
			Deposit:     deposit,
			BlockNumber: big.NewInt(can.BlockNumber.Int64()),
			TxIndex:     can.TxIndex,
			CandidateId: can.CandidateId,
			Host:        can.Host,
			Port:        can.Port,
			Owner:       can.Owner,
			Extra:       can.Extra,
			Fee:         can.Fee,
			TxHash:      can.TxHash,
			TOwner:      can.TOwner,
		}
		copyCandidateQueue[i] = canCopy
	}
	return copyCandidateQueue
}

func CompareCan(curr, target *Candidate, currMoney, targetMoney *big.Int) int {
	return compare2(curr, target, currMoney, targetMoney)
}

func compare(cand CanConditions, c, can *Candidate) int {
	// put the larger deposit in front
	if cand[c.CandidateId].Cmp(cand[can.CandidateId]) > 0 /* c.Deposit.Cmp(can.Deposit) > 0*/ {
		return 1
	} else if cand[c.CandidateId].Cmp(cand[can.CandidateId]) == 0 /* c.Deposit.Cmp(can.Deposit) == 0 */ {
		// put the smaller blocknumber in front
		if c.BlockNumber.Cmp(can.BlockNumber) > 0 {
			return -1
		} else if c.BlockNumber.Cmp(can.BlockNumber) == 0 {
			// put the smaller tx'index in front
			if c.TxIndex > can.TxIndex {
				return -1
			} else if c.TxIndex == can.TxIndex {
				return 0
			} else {
				return 1
			}
		} else {
			return 1
		}
	} else {
		return -1
	}
}

func compare2(curr, target *Candidate, currMoney, targetMoney *big.Int) int {
	// put the larger deposit in front
	if currMoney.Cmp(targetMoney) > 0 {
		return 1
	} else if currMoney.Cmp(targetMoney) == 0 {
		// put the smaller blocknumber in front
		if curr.BlockNumber.Cmp(target.BlockNumber) > 0 {
			return -1
		} else if curr.BlockNumber.Cmp(target.BlockNumber) == 0 {
			// put the smaller tx'index in front
			if curr.TxIndex > target.TxIndex {
				return -1
			} else if curr.TxIndex == target.TxIndex {
				return 0
			} else {
				return 1
			}
		} else {
			return 1
		}
	} else {
		return -1
	}
}

// sorted candidates
func (arr CandidateQueue) CandidateSort(cand CanConditions) {
	if len(arr) <= 1 {
		return
	}
	arr.quickSort(cand, 0, len(arr)-1)
}
func (arr CandidateQueue) quickSort(cand CanConditions, left, right int) {
	if left < right {
		pivot := arr.partition(cand, left, right)
		arr.quickSort(cand, left, pivot-1)
		arr.quickSort(cand, pivot+1, right)
	}
}
func (arr CandidateQueue) partition(cand CanConditions, left, right int) int {
	for left < right {
		for left < right && compare(cand, arr[left], arr[right]) >= 0 {
			right--
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			left++
		}
		for left < right && compare(cand, arr[left], arr[right]) >= 0 {
			left++
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			right--
		}
	}
	return left
}

// candiate info
type Candidate struct {
	// Mortgage amount (margin)
	Deposit *big.Int
	// Current block height number at the time of the mortgage
	BlockNumber *big.Int
	// Current tx'index at the time of the mortgage
	TxIndex uint32
	// candidate's server nodeId
	CandidateId discover.NodeID
	Host        string
	Port        string
	// Mortgage beneficiary's account address
	Owner common.Address
	Extra string
	// brokerage   example: (fee/10000) * 100% == x%
	Fee uint32
	// Selected TicketId
	TxHash common.Hash
	TOwner common.Address
}

type RefundQueue []*CandidateRefund

func (queue RefundQueue) DeepCopy() RefundQueue {
	copyRefundQueue := make(RefundQueue, len(queue))
	if len(queue) == 0 {
		return copyRefundQueue
	}
	for i, refund := range queue {
		deposit, _ := new(big.Int).SetString(refund.Deposit.String(), 10)
		refundCopy := &CandidateRefund{
			Deposit:     deposit,
			BlockNumber: big.NewInt(refund.BlockNumber.Int64()),
			Owner:       refund.Owner,
		}
		copyRefundQueue[i] = refundCopy
	}
	return copyRefundQueue
}

// Refund Info
type CandidateRefund struct {
	// Mortgage amount (margin)
	Deposit *big.Int
	// Current block height number at the time of the mortgage
	BlockNumber *big.Int
	// Mortgage beneficiary's account address
	Owner common.Address
}
