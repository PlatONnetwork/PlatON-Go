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

package staking

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (
	/**
	######   ######   ######   ######
	#	  THE CANDIDATE  STATUS     #
	######   ######   ######   ######
	*/
	Invalided     CandidateStatus = 1 << iota // 0001: The current candidate withdraws from the staking qualification (Active OR Passive)
	LowRatio                                  // 0010: The candidate was low package ratio AND no delete
	NotEnough                                 // 0100: The current candidate's von does not meet the minimum staking threshold
	DuplicateSign                             // 1000: The Duplicate package or Duplicate sign
	LowRatioDel                               // 0001,0000: The lowRatio AND must delete
	Withdrew                                  // 0010,0000: The Active withdrew
	Valided       = 0                         // 0000: The current candidate is in force
	NotExist      = 1 << 31                   // 1000,xxxx,... : The candidate is not exist
)

type CandidateStatus uint32

func (status CandidateStatus) IsValid() bool {
	return !status.IsInvalid()
}

func (status CandidateStatus) IsInvalid() bool {
	return status&Invalided == Invalided
}

func (status CandidateStatus) IsPureInvalid() bool {
	return status&Invalided == status|Invalided
}

func (status CandidateStatus) IsLowRatio() bool {
	return status&LowRatio == LowRatio
}

func (status CandidateStatus) IsPureLowRatio() bool {
	return status&LowRatio == status|LowRatio
}

func (status CandidateStatus) IsNotEnough() bool {
	return status&NotEnough == NotEnough
}

func (status CandidateStatus) IsPureNotEnough() bool {
	return status&NotEnough == status|NotEnough
}

func (status CandidateStatus) IsInvalidLowRatio() bool {
	return status&(Invalided|LowRatio) == (Invalided | LowRatio)
}

func (status CandidateStatus) IsInvalidNotEnough() bool {
	return status&(Invalided|NotEnough) == (Invalided | NotEnough)
}

func (status CandidateStatus) IsInvalidLowRatioNotEnough() bool {
	return status&(Invalided|LowRatio|NotEnough) == (Invalided | LowRatio | NotEnough)
}

func (status CandidateStatus) IsLowRatioNotEnough() bool {
	return status&(LowRatio|NotEnough) == (LowRatio | NotEnough)
}

func (status CandidateStatus) IsDuplicateSign() bool {
	return status&DuplicateSign == DuplicateSign
}

func (status CandidateStatus) IsInvalidDuplicateSign() bool {
	return status&(DuplicateSign|Invalided) == (DuplicateSign | Invalided)
}

func (status CandidateStatus) IsLowRatioDel() bool {
	return status&LowRatioDel == LowRatioDel
}

func (status CandidateStatus) IsPureLowRatioDel() bool {
	return status&LowRatioDel == status|LowRatioDel
}

func (status CandidateStatus) IsInvalidLowRatioDel() bool {
	return status&(Invalided|LowRatioDel) == (Invalided | LowRatioDel)
}

func (status CandidateStatus) IsWithdrew() bool {
	return status&Withdrew == Withdrew
}

func (status CandidateStatus) IsPureWithdrew() bool {
	return status&Withdrew == status|Withdrew
}

func (status CandidateStatus) IsInvalidWithdrew() bool {
	return status&(Invalided|Withdrew) == (Invalided | Withdrew)
}

// The Candidate info
type Candidate struct {
	*CandidateBase
	*CandidateMutable
}

func (can *Candidate) String() string {
	return fmt.Sprintf(`{"NodeId": "%s","BlsPubKey": "%s","StakingAddress": "%s","BenefitAddress": "%s","StakingTxIndex": %d,"ProgramVersion": %d,"Status": %d,"StakingEpoch": %d,"StakingBlockNum": %d,"Shares": %d,"Released": %d,"ReleasedHes": %d,"RestrictingPlan": %d,"RestrictingPlanHes": %d,"ExternalId": "%s","NodeName": "%s","Website": "%s","Details": "%s","DelegateEpoch": %d,"DelegateTotal": %d,"DelegateTotalHes": %d,"RewardPer": %d,"NextRewardPer": %d}`,
		fmt.Sprintf("%x", can.NodeId.Bytes()),
		fmt.Sprintf("%x", can.BlsPubKey.Bytes()),
		fmt.Sprintf("%x", can.StakingAddress.Bytes()),
		fmt.Sprintf("%x", can.BenefitAddress.Bytes()),
		can.StakingTxIndex,
		can.ProgramVersion,
		can.Status,
		can.StakingEpoch,
		can.StakingBlockNum,
		can.Shares,
		can.Released,
		can.ReleasedHes,
		can.RestrictingPlan,
		can.RestrictingPlanHes,
		can.ExternalId,
		can.NodeName,
		can.Website,
		can.Details,
		can.DelegateEpoch,
		can.DelegateTotal,
		can.DelegateTotalHes,
		can.RewardPer,
		can.NextRewardPer)
}

func (can *Candidate) IsNotEmpty() bool {
	return !can.IsEmpty()
}

func (can *Candidate) IsEmpty() bool {
	return nil == can
}

type CandidateBase struct {
	NodeId discover.NodeID
	// bls public key
	BlsPubKey bls.PublicKeyHex
	// The account used to initiate the staking
	StakingAddress common.Address
	// The account receive the block rewards and the staking rewards
	BenefitAddress common.Address
	// The tx index at the time of staking
	StakingTxIndex uint32
	// The version of the node program
	// (Store Large Verson: the 2.1.x large version is 2.1.0)
	ProgramVersion uint32
	// Block height at the time of staking
	StakingBlockNum uint64
	// Node desc
	Description
}

func (can *CandidateBase) String() string {
	return fmt.Sprintf(`{"NodeId": "%s","BlsPubKey": "%s","StakingAddress": "%s","BenefitAddress": "%s","StakingTxIndex": %d,"ProgramVersion": %d,"StakingBlockNum": %d,"ExternalId": "%s","NodeName": "%s","Website": "%s","Details": "%s"}`,
		fmt.Sprintf("%x", can.NodeId.Bytes()),
		fmt.Sprintf("%x", can.BlsPubKey.Bytes()),
		fmt.Sprintf("%x", can.StakingAddress.Bytes()),
		fmt.Sprintf("%x", can.BenefitAddress.Bytes()),
		can.StakingTxIndex,
		can.ProgramVersion,
		can.StakingBlockNum,
		can.ExternalId,
		can.NodeName,
		can.Website,
		can.Details)
}

func (can *CandidateBase) IsNotEmpty() bool {
	return !can.IsEmpty()
}

func (can *CandidateBase) IsEmpty() bool {
	return nil == can
}

type CandidateMutableOld struct {
	// The candidate status
	// Reference `THE CANDIDATE  STATUS`
	Status CandidateStatus
	// The epoch number at staking or edit
	StakingEpoch uint32
	// All vons of staking and delegated
	Shares *big.Int
	// The staking von  is circulating for effective epoch (in effect)
	Released *big.Int
	// The staking von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *big.Int
	// The staking von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *big.Int
	// The staking von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *big.Int
	// The epoch number at delegate or edit
	DelegateEpoch uint32
	// Effective total delegate
	DelegateTotal *big.Int
	// hesitant total delegate
	DelegateTotalHes *big.Int
	// Delegate reward amount percent for current settlement cycle
	RewardPer uint16
	// Delegate reward amount percent for next settlement cycle
	NextRewardPer uint16
	// Number of settlement cycles when changing the commission reward percentage
	RewardPerChangeEpoch uint32

	// current epoch  total Delegate reward
	CurrentEpochDelegateReward *big.Int

	// total  Delegate reward have give
	DelegateRewardTotal *big.Int
}

type CandidateMutable struct {
	// The candidate status
	// Reference `THE CANDIDATE  STATUS`
	Status CandidateStatus
	// The epoch number at staking or edit
	StakingEpoch uint32
	// All vons of staking and delegated
	Shares *big.Int
	// The staking von  is circulating for effective epoch (in effect)
	Released *big.Int
	// The staking von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *big.Int
	// The staking von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *big.Int
	// The staking von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *big.Int
	// The epoch number at delegate or edit
	DelegateEpoch uint32
	// Effective total delegate
	DelegateTotal *big.Int
	// hesitant total delegate
	DelegateTotalHes *big.Int
	// Delegate reward amount percent for current settlement cycle
	RewardPer uint16
	// Delegate reward amount percent for next settlement cycle
	NextRewardPer uint16
	// Number of settlement cycles when changing the commission reward percentage
	RewardPerChangeEpoch uint32

	// current epoch  total Delegate reward
	CurrentEpochDelegateReward *big.Int

	// total  Delegate reward have give
	DelegateRewardTotal *big.Int

	// The settlement period to which the effective delegation is cancelled
	WithdrewDelegateEpoch uint32

	// The number of effective cancellations (the total number of cancellations in a settlement period)
	WithdrewDelegateAmount *big.Int
}

func (can *CandidateMutable) SubDelegateTotalHes(amount *big.Int) {
	can.DelegateTotalHes = new(big.Int).Sub(can.DelegateTotalHes, amount)
}

func (can *CandidateMutable) SubDelegateTotal(amount *big.Int) {
	can.DelegateTotal = new(big.Int).Sub(can.DelegateTotal, amount)
}

func (can *CandidateMutable) PrepareNextEpoch() bool {
	var changed bool
	if can.CurrentEpochDelegateReward.Cmp(common.Big0) > 0 {
		can.DelegateRewardTotal.Add(can.DelegateRewardTotal, can.CurrentEpochDelegateReward)
		//clean  CurrentEpochDelegateReward
		can.CleanCurrentEpochDelegateReward()
		changed = true
	}
	if can.WithdrewDelegateAmount != nil && can.WithdrewDelegateAmount.Cmp(common.Big0) > 0 {
		can.CleanWithdrewDelegateAmount()
		changed = true
	}
	if can.NextRewardPer != can.RewardPer {
		can.RewardPer = can.NextRewardPer
		changed = true
	}
	return changed
}

func (can *CandidateMutable) HasWithdrewDelegation() bool {
	if can.WithdrewDelegateAmount != nil && can.WithdrewDelegateAmount.Cmp(common.Big0) > 0 {
		return true
	}
	return false
}

func (can *CandidateMutable) ShouldGiveDelegateReward(epoch uint32) bool {
	if can.EffectiveDelegateTotal(epoch).Cmp(common.Big0) > 0 && can.RewardPer > 0 {
		return true
	}
	return false
}

func (can *CandidateMutable) EffectiveDelegateTotal(epoch uint32) *big.Int {
	if epoch == can.WithdrewDelegateEpoch {
		return new(big.Int).Add(can.DelegateTotal, can.WithdrewDelegateAmount)
	}
	return can.DelegateTotal
}

func (can *CandidateMutable) AppendWithdrewDelegateAmount(epoch uint32, amount *big.Int) {
	if epoch == can.WithdrewDelegateEpoch {
		can.WithdrewDelegateAmount = new(big.Int).Add(can.WithdrewDelegateAmount, amount)
	} else {
		can.WithdrewDelegateEpoch = epoch
		can.WithdrewDelegateAmount = new(big.Int).Set(amount)
	}
}

func (can *CandidateMutable) CleanWithdrewDelegateAmount() {
	can.WithdrewDelegateEpoch = 0
	can.WithdrewDelegateAmount = new(big.Int).SetUint64(0)
}

func (can *CandidateMutable) String() string {
	return fmt.Sprintf(`{"Status": %d,"StakingEpoch": %d,"Shares": %d,"Released": %d,"ReleasedHes": %d,"RestrictingPlan": %d,"RestrictingPlanHes": %d,"DelegateEpoch": %d,"DelegateTotal": %d,"DelegateTotalHes": %d,"RewardPer": %d,"NextRewardPer": %d,"RewardPerChangeEpoch": %d,"WithdrewDelegateEpoch": %d,"WithdrewDelegateAmount": %d}`,
		can.Status,
		can.StakingEpoch,
		can.Shares,
		can.Released,
		can.ReleasedHes,
		can.RestrictingPlan,
		can.RestrictingPlanHes,
		can.DelegateEpoch,
		can.DelegateTotal,
		can.DelegateTotalHes,
		can.RewardPer,
		can.NextRewardPer,
		can.RewardPerChangeEpoch,
		can.WithdrewDelegateEpoch,
		can.WithdrewDelegateAmount)
}

func (can *CandidateMutable) SetValided() {
	can.Status = Valided
}

func (can *CandidateMutable) SetStatus(status CandidateStatus) {
	can.Status = status
}

func (can *CandidateMutable) AppendStatus(status CandidateStatus) {
	can.Status |= status
}

func (can *CandidateMutable) CleanLowRatioStatus() {
	can.Status &^= LowRatio
}

func (can *CandidateMutable) CleanShares() {
	can.Shares = new(big.Int).SetInt64(0)
}

func (can *CandidateMutable) CleanCurrentEpochDelegateReward() {
	can.CurrentEpochDelegateReward = new(big.Int).SetInt64(0)
}

func (can *CandidateMutable) AddShares(amount *big.Int) {
	can.Shares = new(big.Int).Add(can.Shares, amount)
}

func (can *CandidateMutable) SubShares(amount *big.Int) {
	can.Shares = new(big.Int).Sub(can.Shares, amount)
}

func (can *CandidateMutable) IsNotEmpty() bool {
	return !can.IsEmpty()
}

func (can *CandidateMutable) IsEmpty() bool {
	return nil == can
}

func (can *CandidateMutable) IsValid() bool {
	return can.Status.IsValid()
}

func (can *CandidateMutable) IsInvalid() bool {
	return can.Status.IsInvalid()
}

func (can *CandidateMutable) IsPureInvalid() bool {
	return can.Status.IsPureInvalid()
}

func (can *CandidateMutable) IsLowRatio() bool {
	return can.Status.IsLowRatio()
}

func (can *CandidateMutable) IsPureLowRatio() bool {
	return can.Status.IsPureLowRatio()
}

func (can *CandidateMutable) IsNotEnough() bool {
	return can.Status.IsNotEnough()
}

func (can *CandidateMutable) IsPureNotEnough() bool {
	return can.Status.IsPureNotEnough()
}

func (can *CandidateMutable) IsInvalidLowRatio() bool {
	return can.Status.IsInvalidLowRatio()
}

func (can *CandidateMutable) IsInvalidNotEnough() bool {
	return can.Status.IsInvalidNotEnough()
}

func (can *CandidateMutable) IsInvalidLowRatioNotEnough() bool {
	return can.Status.IsInvalidLowRatioNotEnough()
}

func (can *CandidateMutable) IsLowRatioNotEnough() bool {
	return can.Status.IsLowRatioNotEnough()
}

func (can *CandidateMutable) IsDuplicateSign() bool {
	return can.Status.IsDuplicateSign()
}

func (can *CandidateMutable) IsInvalidDuplicateSign() bool {
	return can.Status.IsInvalidDuplicateSign()
}

func (can *CandidateMutable) IsLowRatioDel() bool {
	return can.Status.IsLowRatioDel()
}

func (can *CandidateMutable) IsPureLowRatioDel() bool {
	return can.Status.IsPureLowRatioDel()
}

func (can *CandidateMutable) IsInvalidLowRatioDel() bool {
	return can.Status.IsInvalidLowRatioDel()
}

func (can *CandidateMutable) IsWithdrew() bool {
	return can.Status.IsWithdrew()
}

func (can *CandidateMutable) IsPureWithdrew() bool {
	return can.Status.IsPureWithdrew()
}

func (can *CandidateMutable) IsInvalidWithdrew() bool {
	return can.Status.IsInvalidWithdrew()
}

// Display amount field using 0x hex
type CandidateHex struct {
	NodeId                 discover.NodeID
	BlsPubKey              bls.PublicKeyHex
	StakingAddress         common.Address
	BenefitAddress         common.Address
	RewardPer              uint16
	NextRewardPer          uint16
	RewardPerChangeEpoch   uint32
	StakingTxIndex         uint32
	ProgramVersion         uint32
	Status                 CandidateStatus
	StakingEpoch           uint32
	StakingBlockNum        uint64
	Shares                 *hexutil.Big
	Released               *hexutil.Big
	ReleasedHes            *hexutil.Big
	RestrictingPlan        *hexutil.Big
	RestrictingPlanHes     *hexutil.Big
	DelegateEpoch          uint32
	DelegateTotal          *hexutil.Big
	DelegateTotalHes       *hexutil.Big
	DelegateRewardTotal    *hexutil.Big
	WithdrewDelegateEpoch  uint32
	WithdrewDelegateAmount *hexutil.Big
	Description
}

func (can *CandidateHex) String() string {
	return fmt.Sprintf(`{"NodeId": "%s","BlsPubKey": "%s","StakingAddress": "%s","BenefitAddress": "%s","RewardPer": "%d","NextRewardPer": "%d","RewardPerChangeEpoch": "%d","StakingTxIndex": %d,"ProgramVersion": %d,"Status": %d,"StakingEpoch": %d,"StakingBlockNum": %d,"Shares": "%s","Released": "%s","ReleasedHes": "%s","RestrictingPlan": "%s","RestrictingPlanHes": "%s","DelegateEpoch": "%d","DelegateTotal": "%s","DelegateTotalHes": "%s","ExternalId": "%s","NodeName": "%s","Website": "%s","Details": "%s","DelegateRewardTotal": "%s","WithdrewDelegateEpoch": "%d","WithdrewDelegateAmount": "%s"}`,
		fmt.Sprintf("%x", can.NodeId.Bytes()),
		fmt.Sprintf("%x", can.BlsPubKey.Bytes()),
		fmt.Sprintf("%x", can.StakingAddress.Bytes()),
		fmt.Sprintf("%x", can.BenefitAddress.Bytes()),
		can.RewardPer,
		can.NextRewardPer,
		can.RewardPerChangeEpoch,
		can.StakingTxIndex,
		can.ProgramVersion,
		can.Status,
		can.StakingEpoch,
		can.StakingBlockNum,
		can.Shares,
		can.Released,
		can.ReleasedHes,
		can.RestrictingPlan,
		can.RestrictingPlanHes,
		can.DelegateEpoch,
		can.DelegateTotal,
		can.DelegateTotalHes,
		can.ExternalId,
		can.NodeName,
		can.Website,
		can.Details,
		can.DelegateRewardTotal,
		can.WithdrewDelegateEpoch,
		can.WithdrewDelegateAmount)
}

func (can *CandidateHex) IsNotEmpty() bool {
	return !can.IsEmpty()
}

func (can *CandidateHex) IsEmpty() bool {
	return nil == can
}

//// EncodeRLP implements rlp.Encoder
//func (c *Candidate) EncodeRLP(w io.Writer) error {
//	return rlp.Encode(w, &c)
//}
//
//
//// DecodeRLP implements rlp.Decoder
//func (c *Candidate) DecodeRLP(s *rlp.Stream) error {
//	if err := s.Decode(&c); err != nil {
//		return err
//	}
//	return nil
//}

const (
	MaxExternalIdLen = 70
	MaxNodeNameLen   = 30
	MaxWebsiteLen    = 140
	MaxDetailsLen    = 280
)

type Description struct {
	// External Id for the third party to pull the node description (with length limit)
	ExternalId string
	// The Candidate Node's Name  (with a length limit)
	NodeName string
	// The third-party home page of the node (with a length limit)
	Website string
	// Description of the node (with a length limit)
	Details string
}

func (desc *Description) CheckLength() error {

	if len(desc.ExternalId) > MaxExternalIdLen {
		return fmt.Errorf("ExternalId overflow, got len is: %d, max len is: %d", len(desc.ExternalId), MaxExternalIdLen)
	}
	if len(desc.NodeName) > MaxNodeNameLen {
		return fmt.Errorf("NodeName overflow, got len is: %d, max len is: %d", len(desc.NodeName), MaxNodeNameLen)
	}
	if len(desc.Website) > MaxWebsiteLen {
		return fmt.Errorf("Website overflow, got len is: %d, max len is: %d", len(desc.Website), MaxWebsiteLen)
	}
	if len(desc.Details) > MaxDetailsLen {
		return fmt.Errorf("Details overflow, got len is: %d, max len is: %d", len(desc.Details), MaxDetailsLen)
	}
	return nil
}

type CandidateQueue []*Candidate

func (queue CandidateQueue) String() string {
	arr := make([]string, len(queue))
	for i, c := range queue {
		arr[i] = c.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}

type CandidateHexQueue []*CandidateHex

func (queue CandidateHexQueue) IsNotEmpty() bool {
	return !queue.IsEmpty()
}

func (queue CandidateHexQueue) IsEmpty() bool {
	return len(queue) == 0
}

func (queue CandidateHexQueue) String() string {
	arr := make([]string, len(queue))
	for i, c := range queue {
		arr[i] = c.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}

type CandidateBaseQueue []*CandidateBase

func (queue CandidateBaseQueue) IsNotEmpty() bool {
	return !queue.IsEmpty()
}

func (queue CandidateBaseQueue) IsEmpty() bool {
	return len(queue) == 0
}

// the Validator info
// They are Simplified Candidate
// They are consensus nodes and Epoch nodes snapshot
/*type Validator struct {
	NodeAddress common.Address
	NodeId      discover.NodeID
	// bls public key
	BlsPubKey bls.PublicKeyHex
	// The weight snapshot
	// NOTE:
	// converted from the weight snapshot of Candidate, they array order is:
	//
	// programVersion, candidate.shares, stakingBlocknum, stakingTxindex
	//
	// They origin type is: uint32, *big.Int, uint64, uint32
	StakingWeight [SWeightItem]string
	// Validator's term in the consensus round
	ValidatorTerm uint32
}*/
type Validator struct {
	ProgramVersion  uint32
	StakingTxIndex  uint32
	ValidatorTerm   uint32 // Validator's term in the consensus round
	StakingBlockNum uint64
	NodeAddress     common.NodeAddress
	NodeId          discover.NodeID
	BlsPubKey       bls.PublicKeyHex
	Shares          *big.Int
}

func (val *Validator) String() string {
	return fmt.Sprintf(`{"NodeId": "%s","NodeAddress": "%s","BlsPubKey": "%s","ProgramVersion": %d,"Shares": %d,"StakingBlockNum": %d,"StakingTxIndex": %d,"ValidatorTerm": %d}`,
		val.NodeId.String(),
		fmt.Sprintf("%x", val.NodeAddress.Bytes()),
		fmt.Sprintf("%x", val.BlsPubKey.Bytes()),
		val.ProgramVersion,
		val.Shares,
		val.StakingBlockNum,
		val.StakingTxIndex,
		val.ValidatorTerm)
}

type ValidatorQueue []*Validator

func (queue ValidatorQueue) IsNotEmpty() bool {
	return !queue.IsEmpty()
}

func (queue ValidatorQueue) IsEmpty() bool {
	return len(queue) == 0
}

func (queue ValidatorQueue) String() string {
	arr := make([]string, len(queue))
	for i, v := range queue {
		arr[i] = v.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}

type CandidateMap map[discover.NodeID]*Candidate

type NeedRemoveCans map[discover.NodeID]*Candidate

func (arr ValidatorQueue) ValidatorSort(removes NeedRemoveCans,
	compare func(slashs NeedRemoveCans, c, can *Validator) int) {
	if len(arr) <= 1 {
		return
	}

	if nil == compare {
		arr.quickSort(removes, 0, len(arr)-1, CompareDefault)
	} else {
		arr.quickSort(removes, 0, len(arr)-1, compare)
	}
}
func (arr ValidatorQueue) quickSort(removes NeedRemoveCans, left, right int,
	compare func(slashs NeedRemoveCans, c, can *Validator) int) {
	if left < right {
		pivot := arr.partition(removes, left, right, compare)
		arr.quickSort(removes, left, pivot-1, compare)
		arr.quickSort(removes, pivot+1, right, compare)
	}
}
func (arr ValidatorQueue) partition(removes NeedRemoveCans, left, right int,
	compare func(slashs NeedRemoveCans, c, can *Validator) int) int {
	for left < right {
		for left < right && compare(removes, arr[left], arr[right]) >= 0 {
			right--
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			left++
		}
		for left < right && compare(removes, arr[left], arr[right]) >= 0 {
			left++
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			right--
		}
	}
	return left
}

// NOTE: Sort By Default
//
// When sorting is done by default,
// it is slashed and is sorted to the end.
//
// The priorities just like that:
// Slashing > ProgramVersion > Shares > BlockNumber > TxIndex
//
// Slashing: From no to yes
// ProgramVersion: From big to small
// Shares： From big to small
// BlockNumber: From small to big
// TxIndex: From small to big
//
// Compare Left And Right
// 1: Left > Right
// 0: Left == Right
// -1:Left < Right
func CompareDefault(removes NeedRemoveCans, left, right *Validator) int {

	compareTxIndexFunc := func(l, r *Validator) int {
		switch {
		case l.StakingTxIndex > r.StakingTxIndex:
			return -1
		case l.StakingTxIndex < r.StakingTxIndex:
			return 1
		default:
			return 0
		}
	}

	compareBlockNumberFunc := func(l, r *Validator) int {

		switch {
		case l.StakingBlockNum > r.StakingBlockNum:
			return -1
		case l.StakingBlockNum < r.StakingBlockNum:
			return 1
		default:
			return compareTxIndexFunc(l, r)
		}
	}

	compareSharesFunc := func(l, r *Validator) int {

		switch {
		case l.Shares.Cmp(r.Shares) < 0:
			return -1
		case l.Shares.Cmp(r.Shares) > 0:
			return 1
		default:
			return compareBlockNumberFunc(l, r)
		}
	}

	_, leftOk := removes[left.NodeId]
	_, rightOk := removes[right.NodeId]

	if leftOk && !rightOk {
		return -1
	} else if !leftOk && rightOk {
		return 1
	} else {

		lversion := xutil.CalcVersion(left.ProgramVersion)
		rversion := xutil.CalcVersion(right.ProgramVersion)

		switch {
		case lversion < rversion:
			return -1
		case lversion > rversion:
			return 1
		default:
			return compareSharesFunc(left, right)
		}
	}

}

// NOTE: These are sorted by priority that will be removed
//
// When sorting is done by delete slashing,
// it is slashed and is sorted to the front.
//
// The priorities just like that:
// Invalid > ProgramVersion > validaotorTerm  > Shares > BlockNumber > TxIndex
//
// What is the invalid ?  That are DuplicateSign and lowRatio&invalid and lowVersion and withdrew&NotInEpochValidators
//
//
//
// Invalid Status: From invalid to valid
// ProgramVersion: From small to big
// validaotorTerm: From big to small
// Shares： From small to big
// BlockNumber: From big to small
// TxIndex: From big to small
//
//
// Compare Left And Right
// 1: Left > Right
// 0: Left == Right
// -1:Left < Right
func CompareForDel(removes NeedRemoveCans, left, right *Validator) int {

	// some funcs

	// Compare TxIndex
	compareTxIndexFunc := func(l, r *Validator) int {

		switch {
		case l.StakingTxIndex > r.StakingTxIndex:
			return 1
		case l.StakingTxIndex < r.StakingTxIndex:
			return -1
		default:
			return 0
		}
	}

	// Compare BlockNumber
	compareBlockNumberFunc := func(l, r *Validator) int {
		switch {
		case l.StakingBlockNum > r.StakingBlockNum:
			return 1
		case l.StakingBlockNum < r.StakingBlockNum:
			return -1
		default:
			return compareTxIndexFunc(l, r)
		}
	}

	// Compare Shares
	compareSharesFunc := func(l, r *Validator) int {

		switch {
		case l.Shares.Cmp(r.Shares) < 0:
			return 1
		case l.Shares.Cmp(r.Shares) > 0:
			return -1
		default:
			return compareBlockNumberFunc(l, r)
		}
	}

	// Compare Term
	compareTermFunc := func(l, r *Validator) int {
		switch {
		case l.ValidatorTerm < r.ValidatorTerm:
			return -1
		case l.ValidatorTerm > r.ValidatorTerm:
			return 1
		default:
			return compareSharesFunc(l, r)
		}
	}

	compareVersionFunc := func(l, r *Validator) int {
		lversion := xutil.CalcVersion(l.ProgramVersion)
		rversion := xutil.CalcVersion(r.ProgramVersion)
		switch {
		case lversion > rversion:
			return -1
		case lversion < rversion:
			return 1
		default:
			return compareTermFunc(left, right)
		}
	}

	lCan, lOK := removes[left.NodeId]
	rCan, rOK := removes[right.NodeId]

	/**
	Start Compare
	*/

	switch {
	case !lOK && rOK: // left need not removed AND right need removed
		return -1
	case !lOK && !rOK: // both need not removed

		// 2. ProgramVersion
		return compareVersionFunc(left, right)

	case lOK && !rOK: // left need removed AND right need not removed
		return 1
	default: // both need removed

		// compare slash
		switch {
		case lCan.IsDuplicateSign() && !rCan.IsDuplicateSign():
			return 1
		case !lCan.IsDuplicateSign() && rCan.IsDuplicateSign():
			return -1
		case lCan.IsDuplicateSign() && rCan.IsDuplicateSign():
			// compare Shares
			return compareSharesFunc(left, right)
		default:
			// compare low ratio delete
			// compare low ratio
			switch {
			case lCan.IsLowRatioDel() && !rCan.IsLowRatioDel():
				return 1
			case !lCan.IsLowRatioDel() && rCan.IsLowRatioDel():
				return -1
			case lCan.IsLowRatioDel() && rCan.IsLowRatioDel():
				// compare Shares
				return compareSharesFunc(left, right)
			default:
				switch {
				case lCan.IsLowRatio() && !rCan.IsLowRatio():
					return 1
				case !lCan.IsLowRatio() && rCan.IsLowRatio():
					return -1
				case lCan.IsLowRatio() && rCan.IsLowRatio():
					// compare Shares
					return compareSharesFunc(left, right)
				default:
					// compare Version
					return compareVersionFunc(left, right)
				}
			}

		}

	}
}

// NOTE: Sort when doing storage
//
// The priorities just like that: (No  shares)
// ProgramVersion > validaotorTerm > BlockNumber > TxIndex
//
// Compare Left And Right
// 1: Left > Right
// 0: Left == Right
// -1:Left < Right
func CompareForStore(_ NeedRemoveCans, left, right *Validator) int {
	// some funcs

	// 5. TxIndex
	compareTxIndexFunc := func(l, r *Validator) int {

		switch {
		case l.StakingTxIndex > r.StakingTxIndex:
			return -1
		case l.StakingTxIndex < r.StakingTxIndex:
			return 1
		default:
			return 0
		}
	}

	// 4. BlockNumber
	compareBlockNumberFunc := func(l, r *Validator) int {

		switch {
		case l.StakingBlockNum > r.StakingBlockNum:
			return -1
		case l.StakingBlockNum < r.StakingBlockNum:
			return 1
		default:
			return compareTxIndexFunc(l, r)
		}
	}

	// 2. Term
	compareTermFunc := func(l, r *Validator) int {
		switch {
		case l.ValidatorTerm < r.ValidatorTerm:
			return -1
		case l.ValidatorTerm > r.ValidatorTerm:
			return 1
		default:
			//return compareSharesFunc(l, r)
			return compareBlockNumberFunc(l, r)
		}
	}

	// 1. ProgramVersion
	lVersion := xutil.CalcVersion(left.ProgramVersion)
	rVersion := xutil.CalcVersion(right.ProgramVersion)
	if lVersion < rVersion {
		return -1
	} else if lVersion > rVersion {
		return 1
	} else {
		return compareTermFunc(left, right)
	}
}

// some consensus round validators or current epoch validators
type ValidatorArray struct {
	// the round start blockNumber or epoch start blockNumber
	Start uint64
	// the round end blockNumber or epoch blockNumber
	End uint64
	// the round validators or epoch validators
	Arr ValidatorQueue
}

func (v ValidatorArray) String() string {
	return fmt.Sprintf(`{"Start": %d, "End": %d, "Arr": %s}`, v.Start, v.End, v.Arr.String())
}

type ValidatorEx struct {
	//NodeAddress common.Address
	NodeId discover.NodeID
	// bls public key
	BlsPubKey bls.PublicKeyHex
	// The account used to initiate the staking
	StakingAddress common.Address
	// The account receive the block rewards and the staking rewards
	BenefitAddress common.Address
	// Delegate reward amount percent for current settlement cycle
	RewardPer uint16
	// Delegate reward amount percent for next settlement cycle
	NextRewardPer uint16
	// Number of settlement cycles when changing the commission reward percentage
	RewardPerChangeEpoch uint32
	// The tx index at the time of staking
	StakingTxIndex uint32
	// The version of the node process
	ProgramVersion uint32
	// Block height at the time of staking
	StakingBlockNum uint64
	// All vons of staking and delegated
	//Shares *big.Int
	Shares *hexutil.Big
	// Node desc
	Description
	// this is the term of validator in consensus round
	// [0, N]
	ValidatorTerm uint32
	// Effective total delegate
	DelegateTotal *hexutil.Big

	DelegateRewardTotal    *hexutil.Big
	WithdrewDelegateEpoch  uint32
	WithdrewDelegateAmount *hexutil.Big
}

func (vex *ValidatorEx) String() string {
	return fmt.Sprintf(`{"NodeId": "%s","NodeAddress": "%s","BlsPubKey": "%s","StakingAddress": "%s","BenefitAddress": "%s","RewardPer": "%d","NextRewardPer": "%d","RewardPerChangeEpoch": "%d","StakingTxIndex": %d,"ProgramVersion": %d,"StakingBlockNum": %d,"Shares": "%s","ExternalId": "%s","NodeName": "%s","Website": "%s","Details": "%s","ValidatorTerm": %d,"DelegateTotal": "%s","WithdrewDelegateEpoch": "%d","WithdrewDelegateAmount": "%s"}`,
		vex.NodeId.String(),
		fmt.Sprintf("%x", vex.StakingAddress.Bytes()),
		fmt.Sprintf("%x", vex.BlsPubKey.Bytes()),
		fmt.Sprintf("%x", vex.StakingAddress.Bytes()),
		fmt.Sprintf("%x", vex.BenefitAddress.Bytes()),
		vex.RewardPer,
		vex.NextRewardPer,
		vex.RewardPerChangeEpoch,
		vex.StakingTxIndex,
		vex.ProgramVersion,
		vex.StakingBlockNum,
		vex.Shares,
		vex.ExternalId,
		vex.NodeName,
		vex.Website,
		vex.Details,
		vex.ValidatorTerm,
		vex.DelegateTotal,
		vex.WithdrewDelegateEpoch,
		vex.WithdrewDelegateAmount)
}

type ValidatorExQueue []*ValidatorEx

func (queue ValidatorExQueue) IsNotEmpty() bool {
	return !queue.IsEmpty()
}

func (queue ValidatorExQueue) IsEmpty() bool {
	return len(queue) == 0
}

func (queue ValidatorExQueue) String() string {
	arr := make([]string, len(queue))
	for i, v := range queue {
		arr[i] = v.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}

func NewDelegation() *Delegation {
	del := new(Delegation)
	// Prevent null pointer initialization
	del.Released = new(big.Int).SetInt64(0)
	del.RestrictingPlan = new(big.Int).SetInt64(0)
	del.ReleasedHes = new(big.Int).SetInt64(0)
	del.RestrictingPlanHes = new(big.Int).SetInt64(0)
	del.CumulativeIncome = new(big.Int).SetInt64(0)
	del.WithdrewAmount = new(big.Int).SetInt64(0)
	del.WithdrewEpoch = 0
	del.UnLockEpoch = 0
	return del
}

// the Delegate information
type Delegation struct {
	// The epoch number at delegate or edit
	DelegateEpoch uint32
	// The delegate von  is circulating for effective epoch (in effect)
	Released *big.Int
	// The delegate von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *big.Int
	// The delegate von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *big.Int
	// The delegate von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *big.Int
	// Cumulative delegate income (Waiting for withdrawal)
	CumulativeIncome *big.Int
	// Settlement period when the delegate is cancelled, Used to calculate delegation reward revenue
	WithdrewEpoch uint32
	// Number of effective delegate cancelled
	WithdrewAmount *big.Int
	// Unlock point for revoking effective commission (Billing cycle)
	UnLockEpoch uint32
}

type DelegationOld struct {
	// The epoch number at delegate or edit
	DelegateEpoch uint32
	// The delegate von  is circulating for effective epoch (in effect)
	Released *big.Int
	// The delegate von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *big.Int
	// The delegate von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *big.Int
	// The delegate von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *big.Int
	// Cumulative delegate income (Waiting for withdrawal)
	CumulativeIncome *big.Int
}

func (del *Delegation) CleanCumulativeIncome(epoch uint32) {
	del.CumulativeIncome = new(big.Int)
	del.DelegateEpoch = epoch
}

func (del *Delegation) CleanWithdrewInfo() {
	del.WithdrewEpoch = 0
	del.WithdrewAmount = new(big.Int).SetUint64(0)
	del.UnLockEpoch = 0
}

func (del *Delegation) Withdrew(epoch uint64, refundAmount *big.Int, duration uint64) {
	del.WithdrewEpoch = uint32(epoch)
	del.WithdrewAmount = new(big.Int).Set(refundAmount)
	del.UnLockEpoch = uint32(epoch) + uint32(duration)
}

func (del *Delegation) TotalHes() *big.Int {
	return new(big.Int).Add(del.ReleasedHes, del.RestrictingPlanHes)
}

func (del *Delegation) String() string {
	return fmt.Sprintf(`{"DelegateEpoch": "%d","Released": "%d","ReleasedHes": %d,"RestrictingPlan": %d,"RestrictingPlanHes": %d,"CumulativeIncome": %d,"WithdrewEpoch": %d,"WithdrewAmount": %d,"UnLockEpoch": %d}`,
		del.DelegateEpoch,
		del.Released,
		del.ReleasedHes,
		del.RestrictingPlan,
		del.RestrictingPlanHes,
		del.CumulativeIncome,
		del.WithdrewEpoch,
		del.WithdrewAmount,
		del.UnLockEpoch)
}

func (del *Delegation) IsNotEmpty() bool {
	return !del.IsEmpty()
}

func (del *Delegation) IsEmpty() bool {
	return nil == del
}

// If there is an delegate that is being withdrawn, the settlement period that it belongs to at the time of cancellation has a profit,
//and the value of "Withdrew Amount" needs to be subtracted from the subsequent settlement period,
//and the cancelled delegate will have no profit.
func (del *Delegation) TotalEffective(epoch uint64) *big.Int {
	totalReleased := new(big.Int).Add(del.Released, del.RestrictingPlan)
	if del.WithdrewEpoch > 0 && epoch > uint64(del.WithdrewEpoch) {
		totalReleased = totalReleased.Sub(totalReleased, del.WithdrewAmount)
	}
	return totalReleased
}

type DelegationHex struct {
	// The epoch number at delegate or edit
	DelegateEpoch uint32
	// The delegate von  is circulating for effective epoch (in effect)
	Released *hexutil.Big
	// The delegate von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *hexutil.Big
	// The delegate von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *hexutil.Big
	// The delegate von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *hexutil.Big
	// Cumulative delegate income (Waiting for withdrawal)
	CumulativeIncome *hexutil.Big
	WithdrewEpoch    uint32
	WithdrewAmount   *hexutil.Big
	UnLockEpoch      uint32
}

func (delHex *DelegationHex) String() string {
	return fmt.Sprintf(`{"DelegateEpoch": %d,"Released": "%s","ReleasedHes": "%s","RestrictingPlan": "%s","RestrictingPlanHes": "%s","CumulativeIncome": "%s","WithdrewEpoch": %d,"WithdrewAmount": "%s","UnLockEpoch": %d}`,
		delHex.DelegateEpoch,
		delHex.Released,
		delHex.ReleasedHes,
		delHex.RestrictingPlan,
		delHex.RestrictingPlanHes,
		delHex.CumulativeIncome,
		delHex.WithdrewEpoch,
		delHex.WithdrewAmount,
		delHex.UnLockEpoch)
}

func (del *DelegationHex) IsNotEmpty() bool {
	return !del.IsEmpty()
}

func (del *DelegationHex) IsEmpty() bool {
	return nil == del
}

type DelegationEx struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
	DelegationHex
}

func (dex *DelegationEx) String() string {
	return fmt.Sprintf(`{"Addr": "%s","NodeId": "%s","StakingBlockNum": %d,"DelegateEpoch": %d,"Released": "%s","ReleasedHes": "%s","RestrictingPlan": "%s","RestrictingPlanHes": "%s","CumulativeIncome": "%s","WithdrewEpoch": %d,"WithdrewAmount": "%s","UnLockEpoch": %d}`,
		dex.Addr.String(),
		fmt.Sprintf("%x", dex.NodeId.Bytes()),
		dex.StakingBlockNum,
		dex.DelegateEpoch,
		dex.Released,
		dex.ReleasedHes,
		dex.RestrictingPlan,
		dex.RestrictingPlanHes,
		dex.CumulativeIncome,
		dex.WithdrewEpoch,
		dex.WithdrewAmount,
		dex.UnLockEpoch)
}

func (dex *DelegationEx) IsNotEmpty() bool {
	return !dex.IsEmpty()
}

func (dex *DelegationEx) IsEmpty() bool {
	return nil == dex
}

type DelegateRelated struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
}

func (dr *DelegateRelated) String() string {
	return fmt.Sprintf(`{"Addr": "%s","NodeId": "%s","StakingBlockNum": "%d"}`,
		dr.Addr.String(),
		fmt.Sprintf("%x", dr.NodeId.Bytes()),
		dr.StakingBlockNum)
}

func (dr *DelegateRelated) IsNotEmpty() bool {
	return !dr.IsEmpty()
}

func (dr *DelegateRelated) IsEmpty() bool {
	return nil == dr
}

type DelRelatedQueue []*DelegateRelated

func (queue DelRelatedQueue) IsNotEmpty() bool {
	return !queue.IsEmpty()
}

func (queue DelRelatedQueue) IsEmpty() bool {
	return len(queue) == 0
}

func (queue DelRelatedQueue) String() string {
	arr := make([]string, len(queue))
	for i, r := range queue {
		arr[i] = r.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}

type UnStakeItem struct {
	// this is the nodeAddress
	NodeAddress     common.NodeAddress
	StakingBlockNum uint64
	// Return to normal staking state
	Recovery bool
}

//type UnDelegateItem struct {
//	// this is the `delegateAddress` + `nodeAddress` + `stakeBlockNumber`
//	KeySuffix []byte
//	Amount    *big.Int
//}

type ValArrIndex struct {
	Start uint64
	End   uint64
}

func (vindex *ValArrIndex) String() string {
	return fmt.Sprintf(`{"Start": %d, "End": %d}`, vindex.Start, vindex.End)
}

type ValArrIndexQueue []*ValArrIndex

func (queue ValArrIndexQueue) ConstantAppend(index *ValArrIndex, size int) (*ValArrIndex, ValArrIndexQueue) {

	queue = append(queue, index)
	if size < len(queue) {
		return queue[0], queue[1:]
	}
	return nil, queue
}

func (queue ValArrIndexQueue) String() string {
	arr := make([]string, len(queue))
	for i, vi := range queue {
		arr[i] = vi.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}

// An item that exists for slash
type SlashNodeItem struct {
	// the nodeId will be slashed
	NodeId discover.NodeID
	// the amount of von with slashed
	Amount *big.Int
	// slash type
	SlashType CandidateStatus
	// the benefit adrr who will receive the slash amount of von
	BenefitAddr common.Address
}

func (s *SlashNodeItem) String() string {
	return fmt.Sprintf(`{"nodeId": %s, "amount": %d, "slashType": %d, "benefitAddr": %s}`, s.NodeId.String(), s.Amount, s.SlashType, s.BenefitAddr)
}

type SlashQueue []*SlashNodeItem

func (queue SlashQueue) String() string {
	arr := make([]string, len(queue))
	for i, s := range queue {
		arr[i] = s.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}
