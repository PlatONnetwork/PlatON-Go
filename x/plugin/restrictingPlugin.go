package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

var (
	errParamPeriodInvalid = errors.New("param period invalid")
	errBalanceNotEnough   = errors.New("balance not enough to restrict")
	errAccountNotFound    = errors.New("account is not found")
)

type restrictingInfo struct {
	balance     *big.Int `json:"balance"` // balance representation all locked amount
	debt        *big.Int `json:"debt"`    // debt representation will released amount. Positive numbers can be used instead of release, 0 means no release, negative numbers indicate not enough to release
	releaseList []uint64 `json:"list"`    // releaseList representation
}

type releaseAmountInfo struct {
	height uint64   `json:"blockNumber"` // blockNumber representation of the block number at the released lock-repo period
	amount *big.Int `json:"amount"`      // amount representation of the released amount
}

type restrictingPlan struct {
	period uint64   `json:"period"` // period representation of the released period at the target blockNumber
	amount *big.Int `json:"amount"` // amount representation of the released amount
}

type Result struct {
	balance *big.Int
	slash   *big.Int
	pledge  *big.Int
	debt    *big.Int
	entry   []byte
}

type RestrictingPlugin struct {
}

var RestrictingPtr *RestrictingPlugin = nil

func GetRestrictingInstance() *RestrictingPlugin {
	if RestrictingPtr == nil {
		RestrictingPtr = &RestrictingPlugin{}
	}
	return RestrictingPtr
}

// BeginBlock does something like check input params before execute transactions,
// in RestrictingPlugin it does nothing.
func (rp *RestrictingPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) (bool, error) {
	return false, nil
}

// EndBlock invoke releaseRestricting
func (rp *RestrictingPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) (bool, error) {

	if xutil.IsSettlementPeriod(head.Number.Uint64()) {

		log.Info("begin to release restricting", "curr", head.Number)
		return rp.releaseRestricting(head, state)
	}

	return true, nil
}

// Comfired is empty function
func (rp *RestrictingPlugin) Confirmed(block *types.Block) error {
	return nil
}

// AddRestrictingRecord stores four K-V record in StateDB:
// RestrictingInfo: the information released to the account
// ReleaseNumbers: the number of accounts to be released at the target block height
// ReleaseAccount: the account on the index at the target block height
// ReleaseAmount: the released account amount at the target block height
func (rp *RestrictingPlugin) AddRestrictingRecord(sender common.Address, account common.Address, plan string,
	state xcom.StateDB) (bool, error) {

	// decode plan from input param
	var plans []restrictingPlan

	if err := rlp.Decode(bytes.NewBuffer([]byte(plan)), plans); err != nil {
		log.Error("failed to rlp encode input plan", "plan", plan)
		return false, err
	}

	// pre-check
	var totalLock = big.NewInt(0)
	for i := 0; i < len(plans); i++ {
		period := plans[i].period
		amount := plans[i].amount

		if vm.RestrictingContractAddr == account && period%120 != 0 {
			log.Error("param period invalid", "period", plans[i].period)
			return false, errParamPeriodInvalid
		}

		totalLock = totalLock.Add(totalLock, amount)
	}

	if state.GetBalance(sender).Cmp(totalLock) == -1 {
		log.Error("sender's balance not enough", "lock", totalLock)
		return false, errBalanceNotEnough
	}

	var (
		recordsNum uint32
		index      uint32
		heightList []uint64
		info       restrictingInfo
	)

	restrictingKey := xcom.GetRestrictingKey(account)
	bRecord := state.GetState(account, restrictingKey)

	if len(bRecord) == 0 { // restricting not exist

		log.Debug("restricting record not exist", "account", account.Bytes())

		for i := 0; i < len(plans); i++ {
			height := plans[i].period * xcom.EpochSize * xcom.ConsensusSize
			amount := plans[i].amount

			releaseNumberKey := xcom.GetReleaseNumberKey(height)
			numbers := state.GetState(vm.RestrictingContractAddr, releaseNumberKey)
			if len(numbers) == 0 {
				recordsNum = uint32(1)
				index = uint32(1)
			} else {
				recordsNum += 1
				index = byteutil.BytesToUint32(numbers) + 1
			}

			state.SetState(vm.RestrictingContractAddr, releaseNumberKey, common.Uint32ToBytes(recordsNum))

			releaseAccountKey := xcom.GetReleaseAccountKey(height, index)
			state.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

			releaseAmountKey := xcom.GetReleaseAmountKey(account, height)
			state.SetState(account, releaseAmountKey, amount.Bytes())

			heightList = append(heightList, height)
		}

		info.balance = totalLock
		info.debt = big.NewInt(0)
		info.releaseList = heightList

	} else { // restricting exist

		log.Debug("restricting record exist", "account", account.Bytes())

		if err := rlp.Decode(bytes.NewReader(bRecord), &info); err != nil {
			log.Error("rlp decode failed", "err", err.Error())
			return false, err
		}

		for i := 0; i < len(plans); i++ {
			// release info
			height := plans[i].period * xcom.EpochSize * xcom.ConsensusSize
			amount := plans[i].amount

			releaseAmountKey := xcom.GetReleaseAmountKey(account, height)
			bAmount := state.GetState(account, releaseAmountKey)

			if len(bAmount) == 0 {

				releaseNumberKey := xcom.GetReleaseNumberKey(height)
				numbers := state.GetState(vm.RestrictingContractAddr, releaseNumberKey)

				index = byteutil.BytesToUint32(numbers) + 1
				state.SetState(vm.RestrictingContractAddr, releaseNumberKey, common.Uint32ToBytes(index))

				releaseAccountKey := xcom.GetReleaseAccountKey(height, index)
				state.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

				info.releaseList = append(info.releaseList, height)

			} else {
				var tmpAmount *big.Int
				if err := rlp.Decode(bytes.NewBuffer(bAmount), tmpAmount); err != nil {
					log.Error("failed to rlp decode release amount", "error", err)
					return false, err
				}
				amount = amount.Add(amount, tmpAmount)
			}

			state.SetState(account, releaseAmountKey, amount.Bytes())
		}

		info.balance = info.balance.Add(info.balance, totalLock)
	}

	bInfo, err := rlp.EncodeToBytes(info)
	if err != nil {
		log.Error("failed to rlp encode restricting info", "account", account, "error", err)
		return true, err
	}
	state.SetState(account, restrictingKey, bInfo)

	state.SubBalance(sender, totalLock)
	state.AddBalance(vm.RestrictingContractAddr, totalLock)

	return true, nil
}

// PledgeLockFunds does nothing, output[0] return true when business is success, else return false
func (rp *RestrictingPlugin) PledgeLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) (bool, error) {

	restrictingKey := xcom.GetRestrictingKey(account)
	record := state.GetState(account, restrictingKey)

	if len(record) == 0 {
		log.Error("record not found in PledgeLockFunds", "account", account, "funds", amount.Uint64())
		return false, errAccountNotFound
	}

	var (
		info restrictingInfo
	)

	if err := rlp.Decode(bytes.NewReader(record), &info); err != nil {
		log.Error("rlp decode failed", "error", err.Error())
		return false, err
	}

	if info.balance.Cmp(amount) == -1 {
		log.Error("restricting balance is not enough", "error", errBalanceNotEnough)
		return false, errBalanceNotEnough
	}

	// sub balance
	info.balance = info.balance.Sub(info.balance, amount)

	// store user info
	bInfo, err := rlp.EncodeToBytes(info)
	if err != nil {
		log.Error("failed to encode restricting info", "account", account, "error", err)
		return true, err
	}
	state.SetState(account, restrictingKey, bInfo)

	state.SubBalance(vm.RestrictingContractAddr, amount)
	state.AddBalance(vm.StakingContractAddr, amount)

	return true, nil
}

// ReturnLockFunds does nothing
func (rp *RestrictingPlugin) ReturnLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) (bool, error) {

	restrictingKey := xcom.GetRestrictingKey(account)
	record := state.GetState(account, restrictingKey)

	if len(record) == 0 {
		log.Error("record not found in ReturnLockFunds", "account", account, "funds", amount.Uint64())
		return false, errAccountNotFound
	}

	var (
		info restrictingInfo
	)

	if err := rlp.Decode(bytes.NewReader(record), &info); err != nil {
		log.Error("rlp decode failed", "error", err.Error())
		return false, err
	}

	if info.debt.Sign() == 1 {
		// add balance while debt is positive
		log.Trace("debt is positive", "account", account, "debt", info.debt)
		info.balance = info.balance.Add(info.balance, amount)

	} else {
		log.Trace("debt is not positive", "account", account, "debt", info.debt, "amount", amount)

		if amount.CmpAbs(info.debt) < 1 {
			info.debt = info.debt.Add(info.debt, amount)
		} else {
			info.balance.Add(info.balance, info.debt.Add(info.debt, amount))
			info.debt = big.NewInt(0)
		}
	}

	bInfo, err := rlp.EncodeToBytes(info)
	if err != nil {
		log.Error("failed to encode restricting info", "account", account, "error", err)
		return true, err
	}
	state.SetState(account, restrictingKey, bInfo)

	state.SubBalance(vm.StakingContractAddr, amount)
	state.AddBalance(vm.RestrictingContractAddr, amount)

	return true, nil
}

// SlashingNotify does nothing
func (rp *RestrictingPlugin) SlashingNotify(account common.Address, amount *big.Int, state xcom.StateDB) (bool, error) {

	restrictingKey := xcom.GetRestrictingKey(account)
	record := state.GetState(account, restrictingKey)

	if len(record) == 0 {
		log.Error("record not found in SlashingNotify", "account", account, "funds", amount.Uint64())
		return false, errAccountNotFound
	}

	var (
		info restrictingInfo
	)

	if err := rlp.Decode(bytes.NewReader(record), &info); err != nil {
		log.Error("rlp decode failed", "error", err.Error())
		return false, err
	}

	info.debt = info.debt.Sub(info.debt, amount)

	bInfo, err := rlp.EncodeToBytes(info)
	if err != nil {
		log.Error("failed to encode restricting info", "account", account, "error", err)
		return true, err
	}
	state.SetState(account, restrictingKey, bInfo)

	return true, nil
}

// releaseRestricting does nothing
func (rp *RestrictingPlugin) releaseRestricting(head *types.Header, state xcom.StateDB) (bool, error) {

	var blockNumber = head.Number.Uint64()

	releaseNumberKey := xcom.GetReleaseNumberKey(blockNumber)
	bNumbers := state.GetState(vm.RestrictingContractAddr, releaseNumberKey)

	if len(bNumbers) == 0 {
		return true, nil
	}

	numbers := byteutil.BytesToUint32(bNumbers)

	var (
		info    restrictingInfo
		release *big.Int
	)

	for index := numbers; index > 0; index++ {

		releaseAccountKey := xcom.GetReleaseAccountKey(blockNumber, index)
		bAccount := state.GetState(vm.RestrictingContractAddr, releaseAccountKey)
		account := byteutil.BytesToAddress(bAccount)

		releaseAmountKey := xcom.GetReleaseAmountKey(account, blockNumber)
		bRelease := state.GetState(account, releaseAmountKey)

		if err := rlp.Decode(bytes.NewBuffer(bRelease), release); err != nil {
			log.Error("rlp decode failed", "origin", bRelease)
			return false, err
		}

		restrictingKey := xcom.GetRestrictingKey(account)
		record := state.GetState(account, restrictingKey)

		if err := rlp.Decode(bytes.NewReader(record), &info); err != nil {
			log.Error("rlp decode failed", "error", err.Error())
			return false, err
		}

		if info.balance.Uint64() == 0 {
			log.Trace("balance of release account equals zero", "account", account, "blockNumber", blockNumber)
			info.debt = info.debt.Sub(info.debt, release)

		} else {

			if info.debt.Sign() == 1 {
				log.Trace("debt of release account lt zero", "account", account, "debt", info.debt)

				if release.Cmp(info.debt) <= 0 {
					info.debt = info.debt.Sub(info.debt, release)

				} else if release.Cmp(info.debt.Add(info.debt, info.balance)) <= 0 {
					release = release.Sub(release, info.debt)
					info.balance = info.balance.Sub(info.balance, release)
					info.debt = big.NewInt(0)

					state.SubBalance(vm.RestrictingContractAddr, release)
					state.AddBalance(account, release)

				} else {
					release = release.Sub(release, info.balance)

					state.SubBalance(vm.RestrictingContractAddr, info.balance)
					state.AddBalance(account, info.balance)

					info.balance = big.NewInt(0)
					info.debt = info.debt.Sub(info.debt, release)
				}

			} else {
				log.Trace("debt of release account le zero", "account", account, "debt", info.debt)

				if release.Cmp(info.balance) <= 0 {
					info.balance = info.balance.Sub(info.balance, release)

					state.SubBalance(vm.RestrictingContractAddr, release)
					state.AddBalance(account, release)

				} else {
					release = release.Sub(release, info.balance)

					state.SubBalance(vm.RestrictingContractAddr, info.balance)
					state.AddBalance(account, info.balance)

					info.balance = big.NewInt(0)
					info.debt = info.debt.Sub(info.debt, release)
				}
			}
		}
		// '''
		// 删除记录
		// '''
	}

	return true, nil
}

func (rp *RestrictingPlugin) GetRestrictingInfo(account common.Address, state xcom.StateDB) ([]byte, error) {

	restrictingKey := xcom.GetRestrictingKey(account)
	record := state.GetState(account, restrictingKey)

	if len(record) == 0 {
		log.Error("record not found in GetRestrictingInfo", "account", account)
		return []byte{}, errAccountNotFound
	}

	var (
		info             restrictingInfo
		releaseAmountKey []byte
		bAmount          []byte
		amount           *big.Int
		plans            []releaseAmountInfo
		plan             releaseAmountInfo
		result           Result
	)

	if err := rlp.Decode(bytes.NewReader(record), &info); err != nil {
		log.Error("rlp decode failed", "error", err.Error())
		return []byte{}, err
	}

	for i := 0; i < len(info.releaseList); i++ {
		blockNumber := info.releaseList[i]

		releaseAmountKey = xcom.GetReleaseAmountKey(account, blockNumber)
		bAmount = state.GetState(account, releaseAmountKey)

		if err := rlp.Encode(bytes.NewBuffer(bAmount), &amount); err != nil {
			log.Error("failed to rlp decode amount in GetRestrictingInfo")
			return []byte{}, err
		}

		plan.height = blockNumber
		plan.amount = amount

		plans = append(plans, plan)
	}

	bPlans, err := json.Marshal(plans)
	if err != nil {
		log.Error("falied to Marshal restricting result")
		return []byte{}, err
	}

	result.balance = info.balance
	result.debt = info.debt
	result.slash = big.NewInt(0)
	result.pledge = big.NewInt(0)
	result.entry = bPlans

	log.Trace("get restricting result", "result", result)

	return rlp.EncodeToBytes(result)
}
