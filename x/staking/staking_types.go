package staking

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (
	/**
	######   ######   ######   ######
	#	  THE CANDIDATE  STATUS     #
	######   ######   ######   ######
	*/
	Invalided     = 1 << iota // 0001: The current candidate withdraws from the staking qualification (Active OR Passive)
	LowRatio                  // 0010: The candidate was low package ratio
	NotEnough                 // 0100: The current candidate's von does not meet the minimum staking threshold
	DuplicateSign             // 1000: The Duplicate package or Duplicate sign
	Valided       = 0         // 0000: The current candidate is in force
	NotExist      = 1 << 31   // 1000,xxxx,... : The candidate is not exist
)

const SWeightItem = 4

func Is_Invalid(status uint32) bool {
	return status&Invalided == Invalided
}

func Is_PureInvalid(status uint32) bool {
	return status&Invalided == status|Invalided
}

func Is_LowRatio(status uint32) bool {
	return status&LowRatio == LowRatio
}

func Is_PureLowRatio(status uint32) bool {
	return status&LowRatio == status|LowRatio
}

func Is_NotEnough(status uint32) bool {
	return status&NotEnough == NotEnough
}

func Is_PureNotEnough(status uint32) bool {
	return status&NotEnough == status|NotEnough
}

func Is_Invalid_LowRatio(status uint32) bool {
	return status&(Invalided|LowRatio) == (Invalided | LowRatio)
}

func Is_Invalid_NotEnough(status uint32) bool {
	return status&(Invalided|NotEnough) == (Invalided | NotEnough)
}

func Is_Invalid_LowRatio_NotEnough(status uint32) bool {
	return status&(Invalided|LowRatio|NotEnough) == (Invalided | LowRatio | NotEnough)
}

func Is_LowRatio_NotEnough(status uint32) bool {
	return status&(LowRatio|NotEnough) == (LowRatio | NotEnough)
}

func Is_DuplicateSign(status uint32) bool {
	return status&DuplicateSign == DuplicateSign
}

func Is_DuplicateSign_Invalid(status uint32) bool {
	return status&(DuplicateSign|Invalided) == (DuplicateSign | Invalided)
}

// The Candidate info
type Candidate struct {
	NodeId discover.NodeID
	// The account used to initiate the staking
	StakingAddress common.Address
	// The account receive the block rewards and the staking rewards
	BenefitAddress common.Address
	// The tx index at the time of staking
	StakingTxIndex uint32
	// The version of the node program
	// (Store Large Verson: the 2.1.x large version is 2.1.0)
	ProgramVersion uint32
	// The candidate status
	// Reference `THE CANDIDATE  STATUS`
	Status uint32
	// The epoch number at staking or edit
	StakingEpoch uint32
	// Block height at the time of staking
	StakingBlockNum uint64
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

	// Node desc
	Description
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

type CandidateQueue []*Candidate

// the Validator info
// They are Simplified Candidate
// They are consensus nodes and Epoch nodes snapshot
type Validator struct {
	NodeAddress common.Address
	NodeId      discover.NodeID
	// The weight snapshot
	// NOTE:
	// converted from the weight snapshot of Candidate, they array order is:
	//
	// programVersion, candidate.shares, stakingBlocknum, stakingTxindex
	//
	// They origin type is: uint32, *big.int, uint64, uint32
	StakingWeight [SWeightItem]string
	// Validator's term in the consensus round
	ValidatorTerm uint32
}

func (val *Validator) GetProgramVersion() (uint32, error) {
	version := val.StakingWeight[0]
	v, err := strconv.Atoi(version)
	if nil != err {
		return 0, err
	}
	return uint32(v), nil
}
func (val *Validator) GetShares() (*big.Int, error) {
	shares, ok := new(big.Int).SetString(val.StakingWeight[1], 10)
	if !ok {
		return nil, errors.New("parse bigInt failed from validator's shares")
	}
	return shares, nil
}

func (val *Validator) GetStakingBlockNumber() (uint64, error) {
	stakingBlockNumber := val.StakingWeight[2]
	num, err := strconv.ParseUint(stakingBlockNumber, 10, 64)
	if nil != err {
		return 0, err
	}
	return uint64(num), nil
}

func (val *Validator) GetStakingTxIndex() (uint32, error) {
	txIndex := val.StakingWeight[3]
	index, err := strconv.Atoi(txIndex)
	if nil != err {
		return 0, err
	}
	return uint32(index), nil
}

type ValidatorQueue []*Validator

//type SlashMark map[discover.NodeID]struct{}
type SlashCandidate map[common.Address]*Candidate

func (arr ValidatorQueue) ValidatorSort(slashs SlashCandidate,
	compare func(slashs SlashCandidate, c, can *Validator) int) {
	if len(arr) <= 1 {
		return
	}

	if nil == compare {
		arr.quickSort(slashs, 0, len(arr)-1, CompareDefault)
	} else {
		arr.quickSort(slashs, 0, len(arr)-1, compare)
	}
}
func (arr ValidatorQueue) quickSort(slashs SlashCandidate, left, right int,
	compare func(slashs SlashCandidate, c, can *Validator) int) {
	if left < right {
		pivot := arr.partition(slashs, left, right, compare)
		arr.quickSort(slashs, left, pivot-1, compare)
		arr.quickSort(slashs, pivot+1, right, compare)
	}
}
func (arr ValidatorQueue) partition(slashs SlashCandidate, left, right int,
	compare func(slashs SlashCandidate, c, can *Validator) int) int {
	for left < right {
		for left < right && compare(slashs, arr[left], arr[right]) >= 0 {
			right--
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			left++
		}
		for left < right && compare(slashs, arr[left], arr[right]) >= 0 {
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
func CompareDefault(slashs SlashCandidate, left, right *Validator) int {

	compareTxIndexFunc := func(l, r *Validator) int {
		leftTxIndex, _ := l.GetStakingTxIndex()
		rightTxIndex, _ := r.GetStakingTxIndex()
		switch {
		case leftTxIndex > rightTxIndex:
			return -1
		case leftTxIndex < rightTxIndex:
			return 1
		default:
			return 0
		}
	}

	compareBlockNumberFunc := func(l, r *Validator) int {
		leftNum, _ := l.GetStakingBlockNumber()
		rightNum, _ := r.GetStakingBlockNumber()
		switch {
		case leftNum > rightNum:
			return -1
		case leftNum < rightNum:
			return 1
		default:
			return compareTxIndexFunc(l, r)
		}
	}

	compareSharesFunc := func(l, r *Validator) int {
		leftShares, _ := l.GetShares()
		rightShares, _ := r.GetShares()

		switch {
		case leftShares.Cmp(rightShares) < 0:
			return -1
		case leftShares.Cmp(rightShares) > 0:
			return 1
		default:
			return compareBlockNumberFunc(l, r)
		}
	}

	_, leftOk := slashs[left.NodeAddress]
	_, rightOk := slashs[right.NodeAddress]

	if leftOk && !rightOk {
		return -1
	} else if !leftOk && rightOk {
		return 1
	} else {
		leftVersion, _ := left.GetProgramVersion()
		rightVersion, _ := right.GetProgramVersion()

		switch {
		case leftVersion < rightVersion:
			return -1
		case leftVersion > rightVersion:
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
// DuplicateSign > ProgramVersion > LowPackageRatio > validaotorTerm  > Shares > BlockNumber > TxIndex
//
// DuplicateSign: From yes to no (When both are double-signed, priority is given to removing high weights [Shares. BlockNumber. TxIndex].)
// ProgramVersion: From small to big
// validaotorTerm: From big to small
// LowPackageRatio: From small to big (When both are zero package, priority is given to removing high weights [Shares. BlockNumber. TxIndex].)
// Shares： From small to bigLowPackageRatio
// BlockNumber: From big to small
// TxIndex: From big to small
//
// Compare Left And Right
// 1: Left > Right
// 0: Left == Right
// -1:Left < Right
func CompareForDel(slashs SlashCandidate, left, right *Validator) int {

	// some funcs

	// 7. TxIndex
	compareTxIndexFunc := func(l, r *Validator) int {
		leftTxIndex, _ := l.GetStakingTxIndex()
		rightTxIndex, _ := r.GetStakingTxIndex()
		switch {
		case leftTxIndex > rightTxIndex:
			return -1
		case leftTxIndex < rightTxIndex:
			return 1
		default:
			return 0
		}
	}

	// 6. BlockNumber
	compareBlockNumberFunc := func(l, r *Validator) int {
		leftNum, _ := l.GetStakingBlockNumber()
		rightNum, _ := r.GetStakingBlockNumber()
		switch {
		case leftNum > rightNum:
			return -1
		case leftNum < rightNum:
			return 1
		default:
			return compareTxIndexFunc(l, r)
		}
	}

	// 5. Shares
	compareSharesFunc := func(l, r *Validator) int {
		leftShares, _ := l.GetShares()
		rightShares, _ := r.GetShares()

		switch {
		case leftShares.Cmp(rightShares) < 0:
			return -1
		case leftShares.Cmp(rightShares) > 0:
			return 1
		default:
			return compareBlockNumberFunc(l, r)
		}
	}

	// 4. Term
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

	lCan, lOK := slashs[left.NodeAddress]
	rCan, rOK := slashs[right.NodeAddress]

	// 1. Duplicate Sign
	if lOK && Is_DuplicateSign(lCan.Status) { // left is duplicateSign
		if !rOK || (rOK && !Is_DuplicateSign(rCan.Status)) { // right is not duplicateSign
			return 1
		} else {

			// When both duplicateSign
			/*lversion, _ := left.GetProgramVersion()
			rversion, _ := right.GetProgramVersion()
			switch {
			case lversion > rversion:
				return -1
			case lversion < rversion:
				return 1
			default:
				return compareSharesFunc(left, right)
			}*/
			return compareSharesFunc(left, right)
		}
	} else { // left is not duplicateSign

		if rOK && Is_DuplicateSign(rCan.Status) { // right is duplicateSign
			return -1
		} else { // When both no duplicateSign

			// 2. ProgramVersion
			lversion, _ := left.GetProgramVersion()
			rversion, _ := right.GetProgramVersion()
			switch {
			case lversion > rversion:
				return -1
			case lversion < rversion:
				return 1
			default:

				// 3. LowPackageRatio
				if lOK && Is_LowRatio(lCan.Status) { // left is LowRatio
					if !rOK { // right is not LowRatio
						return 1
					} else { // When both LowRatio

						switch {
						// left.Status(xxxxx1) && right.Status(xxxxx0)
						case Is_Invalid(lCan.Status) && !Is_Invalid(rCan.Status):
							return 1
						// left.Status(xxxxx0) && right.Status(xxxxx1)
						case !Is_Invalid(lCan.Status) && Is_Invalid(rCan.Status):
							return -1
						// When both valid OR both Invalid
						default:
							return compareTermFunc(left, right)
						}

					}

				} else { // left is not LowRatio

					if rOK && Is_LowRatio(rCan.Status) { // right is LowRatio
						return -1
					} else { // When both no LowRatio
						return compareTermFunc(left, right)
					}
				}

			}
		}
	}
}

// NOTE: Sort when doing storage
//
// The priorities just like that:
// ProgramVersion > validaotorTerm > Shares > BlockNumber > TxIndex
//
// Compare Left And Right
// 1: Left > Right
// 0: Left == Right
// -1:Left < Right
func CompareForStore(_ SlashCandidate, left, right *Validator) int {
	// some funcs

	// 5. TxIndex
	compareTxIndexFunc := func(l, r *Validator) int {
		leftTxIndex, _ := l.GetStakingTxIndex()
		rightTxIndex, _ := r.GetStakingTxIndex()
		switch {
		case leftTxIndex > rightTxIndex:
			return -1
		case leftTxIndex < rightTxIndex:
			return 1
		default:
			return 0
		}
	}

	// 4. BlockNumber
	compareBlockNumberFunc := func(l, r *Validator) int {
		leftNum, _ := l.GetStakingBlockNumber()
		rightNum, _ := r.GetStakingBlockNumber()
		switch {
		case leftNum > rightNum:
			return -1
		case leftNum < rightNum:
			return 1
		default:
			return compareTxIndexFunc(l, r)
		}
	}

	// 3. Shares
	compareSharesFunc := func(l, r *Validator) int {
		leftShares, _ := l.GetShares()
		rightShares, _ := r.GetShares()

		switch {
		case leftShares.Cmp(rightShares) < 0:
			return -1
		case leftShares.Cmp(rightShares) > 0:
			return 1
		default:
			return compareBlockNumberFunc(l, r)
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
			return compareSharesFunc(l, r)
		}
	}

	// 1. ProgramVersion
	lVersion, _ := left.GetProgramVersion()
	rVersion, _ := right.GetProgramVersion()
	if lVersion < rVersion {
		return -1
	} else if lVersion > rVersion {
		return 1
	} else {
		return compareTermFunc(left, right)
	}
}

// some consensus round validators or current epoch validators
type Validator_array struct {
	// the round start blockNumber or epoch start blockNumber
	Start uint64
	// the round end blockNumber or epoch blockNumber
	End uint64
	// the round validators or epoch validators
	Arr ValidatorQueue
}

type ValidatorEx struct {
	NodeId discover.NodeID
	// The account used to initiate the staking
	StakingAddress common.Address
	// The account receive the block rewards and the staking rewards
	BenefitAddress common.Address
	// The tx index at the time of staking
	StakingTxIndex uint32
	// The version of the node process
	ProgramVersion uint32
	// Block height at the time of staking
	StakingBlockNum uint64
	// All vons of staking and delegated
	Shares *big.Int
	// Node desc
	Description
	// this is the term of validator in consensus round
	// [0, N]
	ValidatorTerm uint32
}

type ValidatorExQueue = []*ValidatorEx

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
	// Total amount in all cancellation plans
	Reduction *big.Int
}

type DelegationEx struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
	Delegation
}

type DelegateRelated struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
}

type DelRelatedQueue = []*DelegateRelated

/*type UnStakeItem struct {
	// this is the nodeAddress
	KeySuffix  	[]byte
	Amount 		*big.Int
}*/

type UnDelegateItem struct {
	// this is the `delegateAddress` + `nodeAddress` + `stakeBlockNumber`
	KeySuffix []byte
	Amount    *big.Int
}

type ValArrIndex struct {
	Start uint64
	End   uint64
}

type ValArrIndexQueue []*ValArrIndex

func (queue ValArrIndexQueue) ConstantAppend(index *ValArrIndex, size int) (*ValArrIndex, ValArrIndexQueue) {
	queue = append(queue, index)
	if size < len(queue) {
		return queue[0], queue[1:]
	}
	return nil, queue
}
