package staking

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

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
	Invalided     = 1 << iota // 0001: The current candidate withdraws from the staking qualification (Active OR Passive)
	LowRatio                  // 0010: The candidate was low package ratio AND no delete
	NotEnough                 // 0100: The current candidate's von does not meet the minimum staking threshold
	DuplicateSign             // 1000: The Duplicate package or Duplicate sign
	LowRatioDel               // 0001,0000: The lowRatio AND must delete
	Withdrew                  // 0010,0000: The Active withdrew
	Valided       = 0         // 0000: The current candidate is in force
	NotExist      = 1 << 31   // 1000,xxxx,... : The candidate is not exist
)

const SWeightItem = 4

func Is_Valid(status uint32) bool {
	return !Is_Invalid(status)
}

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

func Is_Invalid_DuplicateSign(status uint32) bool {
	return status&(DuplicateSign|Invalided) == (DuplicateSign | Invalided)
}

func Is_LowRatioDel(status uint32) bool {
	return status&LowRatioDel == LowRatioDel
}

func Is_PureLowRatioDel(status uint32) bool {
	return status&LowRatioDel == status|LowRatioDel
}

func Is_Invalid_LowRatioDel(status uint32) bool {
	return status&(Invalided|LowRatioDel) == (Invalided | LowRatioDel)
}

func Is_Withdrew(status uint32) bool {
	return status&Withdrew == Withdrew
}

func Is_PureWithdrew(status uint32) bool {
	return status&Withdrew == status|Withdrew
}

func Is_Invalid_Withdrew(status uint32) bool {
	return status&(Invalided|Withdrew) == (Invalided | Withdrew)
}

// The Candidate info
type Candidate struct {
	NodeId discover.NodeID
	// bls public key
	BlsPubKey bls.PublicKey
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

func (can *Candidate) String() string {
	return fmt.Sprintf(`
	{
		"NodeId": "%s", 
		"BlsPubKey": "%s", 
		"StakingAddress": "%s", 
		"BenefitAddress": "%s", 
		"StakingTxIndex": %d, 
		"ProgramVersion": %d, 
		"Status": %d, 
		"StakingEpoch": %d, 
		"StakingBlockNum": %d,
		"Shares": %d,
		"Released": %d,
		"ReleasedHes": %d,
		"RestrictingPlan": %d,
		"RestrictingPlanHes": %d,
		"ExternalId": "%s",
		"NodeName": "%s",
		"Website": "%s",
		"Details": "%s"
	}`,
		can.NodeId.String(),
		fmt.Sprintf("%x", can.BlsPubKey.Serialize()),
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
		can.Details)
}

// Display amount field using 0x hex
type CandidateHex struct {
	NodeId             discover.NodeID
	BlsPubKey          bls.PublicKey
	StakingAddress     common.Address
	BenefitAddress     common.Address
	StakingTxIndex     uint32
	ProgramVersion     uint32
	Status             uint32
	StakingEpoch       uint32
	StakingBlockNum    uint64
	Shares             *hexutil.Big
	Released           *hexutil.Big
	ReleasedHes        *hexutil.Big
	RestrictingPlan    *hexutil.Big
	RestrictingPlanHes *hexutil.Big
	Description
}

func (can *CandidateHex) String() string {
	return fmt.Sprintf(`
	{
		"NodeId": "%s", 
		"BlsPubKey": "%s", 
		"StakingAddress": "%s", 
		"BenefitAddress": "%s", 
		"StakingTxIndex": %d, 
		"ProgramVersion": %d, 
		"Status": %d, 
		"StakingEpoch": %d, 
		"StakingBlockNum": %d,
		"Shares": "%s",
		"Released": "%s",
		"ReleasedHes": "%s",
		"RestrictingPlan": "%s",
		"RestrictingPlanHes": "%s",
		"ExternalId": "%s",
		"NodeName": "%s",
		"Website": "%s",
		"Details": "%s"
	}`,
		can.NodeId.String(),
		fmt.Sprintf("%x", can.BlsPubKey.Serialize()),
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
		can.Details)
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
type CandidateHexQueue []*CandidateHex

// the Validator info
// They are Simplified Candidate
// They are consensus nodes and Epoch nodes snapshot
type Validator struct {
	NodeAddress common.Address
	NodeId      discover.NodeID
	// bls public key
	BlsPubKey bls.PublicKey
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

func (val *Validator) String() string {
	return fmt.Sprintf(`
	{
		"NodeId": "%s", 
		"NodeAddress": "%s",
		"BlsPubKey": "%s", 
		"StakingWeight": %s, 
		"ValidatorTerm": %d
	}`,
		val.NodeId.String(),
		fmt.Sprintf("%x", val.NodeAddress.Bytes()),
		fmt.Sprintf("%x", val.BlsPubKey.Serialize()),
		fmt.Sprintf(`[%s,%s,%s,%s]`, val.StakingWeight[0], val.StakingWeight[1], val.StakingWeight[2], val.StakingWeight[3]),
		val.ValidatorTerm)
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

	_, leftOk := removes[left.NodeId]
	_, rightOk := removes[right.NodeId]

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
		leftTxIndex, _ := l.GetStakingTxIndex()
		rightTxIndex, _ := r.GetStakingTxIndex()
		switch {
		case leftTxIndex > rightTxIndex:
			return 1
		case leftTxIndex < rightTxIndex:
			return -1
		default:
			return 0
		}
	}

	// Compare BlockNumber
	compareBlockNumberFunc := func(l, r *Validator) int {
		leftNum, _ := l.GetStakingBlockNumber()
		rightNum, _ := r.GetStakingBlockNumber()
		switch {
		case leftNum > rightNum:
			return 1
		case leftNum < rightNum:
			return -1
		default:
			return compareTxIndexFunc(l, r)
		}
	}

	// Compare Shares
	compareSharesFunc := func(l, r *Validator) int {
		leftShares, _ := l.GetShares()
		rightShares, _ := r.GetShares()

		switch {
		case leftShares.Cmp(rightShares) < 0:
			return 1
		case leftShares.Cmp(rightShares) > 0:
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
		lversion, _ := l.GetProgramVersion()
		rversion, _ := r.GetProgramVersion()
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
		case Is_DuplicateSign(lCan.Status) && !Is_DuplicateSign(rCan.Status):
			return 1
		case !Is_DuplicateSign(lCan.Status) && Is_DuplicateSign(rCan.Status):
			return -1
		case Is_DuplicateSign(lCan.Status) && Is_DuplicateSign(rCan.Status):
			// compare Shares
			return compareSharesFunc(left, right)
		default:
			// compare low ratio delete
			// compare low ratio
			switch {
			case Is_LowRatioDel(lCan.Status) && !Is_LowRatioDel(rCan.Status):
				return 1
			case !Is_LowRatioDel(lCan.Status) && Is_LowRatioDel(rCan.Status):
				return -1
			case Is_LowRatioDel(lCan.Status) && Is_LowRatioDel(rCan.Status):
				// compare Shares
				return compareSharesFunc(left, right)
			default:
				switch {
				case Is_LowRatio(lCan.Status) && !Is_LowRatio(rCan.Status):
					return 1
				case !Is_LowRatio(lCan.Status) && Is_LowRatio(rCan.Status):
					return -1
				case Is_LowRatio(lCan.Status) && Is_LowRatio(rCan.Status):
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

	//// 3. Shares
	//compareSharesFunc := func(l, r *Validator) int {
	//	leftShares, _ := l.GetShares()
	//	rightShares, _ := r.GetShares()
	//
	//	switch {
	//	case leftShares.Cmp(rightShares) < 0:
	//		return -1
	//	case leftShares.Cmp(rightShares) > 0:
	//		return 1
	//	default:
	//		return compareBlockNumberFunc(l, r)
	//	}
	//}

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
type ValidatorArray struct {
	// the round start blockNumber or epoch start blockNumber
	Start uint64
	// the round end blockNumber or epoch blockNumber
	End uint64
	// the round validators or epoch validators
	Arr ValidatorQueue
}

type ValidatorEx struct {
	//NodeAddress common.Address
	NodeId discover.NodeID
	// bls public key
	BlsPubKey bls.PublicKey
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
	//Shares *big.Int
	Shares *hexutil.Big
	// Node desc
	Description
	// this is the term of validator in consensus round
	// [0, N]
	ValidatorTerm uint32
}

func (vex *ValidatorEx) String() string {
	return fmt.Sprintf(`
	{
		"NodeId": "%s", 
		"NodeAddress": "%s",
		"BlsPubKey": "%s", 
		"StakingAddress": "%s", 
		"BenefitAddress": "%s", 
		"StakingTxIndex": %d, 
		"ProgramVersion": %d,
		"StakingBlockNum": %d,
		"Shares": "%s",
		"ExternalId": "%s",
		"NodeName": "%s",
		"Website": "%s",
		"Details": "%s",
		"ValidatorTerm": %d
	}`,
		vex.NodeId.String(),
		fmt.Sprintf("%x", vex.StakingAddress.Bytes()),
		fmt.Sprintf("%x", vex.BlsPubKey.Serialize()),
		fmt.Sprintf("%x", vex.StakingAddress.Bytes()),
		fmt.Sprintf("%x", vex.BenefitAddress.Bytes()),
		vex.StakingTxIndex,
		vex.ProgramVersion,
		vex.StakingBlockNum,
		vex.Shares,
		vex.ExternalId,
		vex.NodeName,
		vex.Website,
		vex.Details,
		vex.ValidatorTerm)
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
}

func (del *Delegation) String() string {
	return fmt.Sprintf(`
	{
		"DelegateEpoch": "%d", 
		"Released": "%d", 
		"ReleasedHes": %d, 
		"RestrictingPlan": %d,
		"RestrictingPlanHes": %d
	}`,
		del.DelegateEpoch,
		del.Released,
		del.ReleasedHes,
		del.RestrictingPlan,
		del.RestrictingPlanHes)
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
}

func (delHex *DelegationHex) String() string {
	return fmt.Sprintf(`
	{
		"DelegateEpoch": "%d", 
		"Released": "%s", 
		"ReleasedHes": %s, 
		"RestrictingPlan": %s,
		"RestrictingPlanHes": %s
	}`,
		delHex.DelegateEpoch,
		delHex.Released,
		delHex.ReleasedHes,
		delHex.RestrictingPlan,
		delHex.RestrictingPlanHes)
}

type DelegationEx struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
	DelegationHex
}

func (dex *DelegationEx) String() string {
	return fmt.Sprintf(`
	{
		"Addr": "%s", 
		"NodeId": "%s",
		"StakingBlockNum": "%d", 
		"DelegateEpoch": "%d", 
		"Released": "%s", 
		"ReleasedHes": %s, 
		"RestrictingPlan": %s,
		"RestrictingPlanHes": %s
	}`,
		dex.Addr.String(),
		fmt.Sprintf("%x", dex.NodeId.Bytes()),
		dex.StakingBlockNum,
		dex.DelegateEpoch,
		dex.Released,
		dex.ReleasedHes,
		dex.RestrictingPlan,
		dex.RestrictingPlanHes)
}

type DelegateRelated struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
}

func (dr *DelegateRelated) String() string {
	return fmt.Sprintf(`
	{
		"Addr": "%s", 
		"NodeId": "%s",
		"StakingBlockNum": "%d"
	}`,
		dr.Addr.String(),
		fmt.Sprintf("%x", dr.NodeId.Bytes()),
		dr.StakingBlockNum)
}

type DelRelatedQueue = []*DelegateRelated

type UnStakeItem struct {
	// this is the nodeAddress
	NodeAddress     common.Address
	StakingBlockNum uint64
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

type ValArrIndexQueue []*ValArrIndex

func (queue ValArrIndexQueue) ConstantAppend(index *ValArrIndex, size int) (*ValArrIndex, ValArrIndexQueue) {

	queue = append(queue, index)
	if size < len(queue) {
		return queue[0], queue[1:]
	}
	return nil, queue
}

// An item that exists for slash
type SlashNodeItem struct {
	// the nodeId will be slashed
	NodeId discover.NodeID
	// the amount of von with slashed
	Amount *big.Int
	// slash type
	SlashType int
	// the benefit adrr who will receive the slash amount of von
	BenefitAddr common.Address
}

type SlashQueue []*SlashNodeItem
