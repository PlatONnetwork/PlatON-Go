package plugin

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/restriting"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	errParamPeriodInvalid = common.NewBizError("param epoch invalid")
	errBalanceNotEnough = common.NewBizError("balance not enough to restrict")
	errAccountNotFound = common.NewBizError("account is not found")
)


type restrictingInfo struct {
	balance     *big.Int `json:"balance"` // balance representation all locked amount
	debt        *big.Int `json:"debt"`    // debt representation will released amount. Positive numbers can be used instead of release, 0 means no release, negative numbers indicate not enough to release
	releaseList []uint64 `json:"list"`    // releaseList representation which epoch will release restricting
}

type releaseAmountInfo struct {
	height uint64 	 `json:"blockNumber"`  	// blockNumber representation of the block number at the released epoch
	amount *big.Int	 `json:"amount"`		// amount representation of the released amount
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
func (rp *RestrictingPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock invoke releaseRestricting
func (rp *RestrictingPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {

	// !!! get latest epoch
	// epoch := getLatestEpoch()
	// getBlockNumberByEpoch(epoch)
	// !!!

	epoch := uint64(0)
	blockNumber := uint64(0)

	if blockNumber < head.Number.Uint64() {
		return nil
	}

	log.Info("begin to release restricting", "curr", head.Number)
	return rp.releaseRestricting(epoch, state)

}

// Comfired is empty function
func (rp *RestrictingPlugin) Confirmed(block *types.Block) error {
	return nil
}

// AddRestrictingRecord stores four K-V record in StateDB:
// RestrictingInfo: the account info to be released
// ReleaseEpoch:   the number of accounts to be released on the epoch corresponding to the target block height
// ReleaseAccount: the account on the index on the target epoch
// ReleaseAmount: the amount of the account to be released on the target epoch
func (rp *RestrictingPlugin) AddRestrictingRecord(sender common.Address, account common.Address, plans []byteutil.RestrictingPlan,
	state xcom.StateDB) error {

	var (
		err         error
		totalAmount *big.Int
	)

	// pre-check
	// !!! get latest epoch
	// latest := getLatestEpoch()
	// !!!
	latest := uint64(0)
	for i := 0; i < len(plans); i++ {
		epoch  := plans[i].Epoch
		amount := plans[i].Amount

		if epoch < latest {
			log.Error("param epoch invalid", "epoch", epoch, "latest", latest)
			return errParamPeriodInvalid
		}

		totalAmount = totalAmount.Add(totalAmount, amount)
	}

	if state.GetBalance(sender).Cmp(totalAmount) == -1 {
		log.Error("sender's balance not enough", "total", totalAmount)
		return errBalanceNotEnough
	}

	// TODO
	var (
		epochList  []uint64
		index      uint32
		info       restrictingInfo
		accNumbers uint32
	)

	restrictingKey := restriting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Debug("restricting record not exist", "account", account.Bytes())

		for i := 0; i < len(plans); i++ {
			epoch := plans[i].Epoch
			amount := plans[i].Amount

			// step1: get account numbers at target epoch
			releaseEpochKey := restriting.GetReleaseEpochKey(epoch)
			bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)

			if len(bAccNumbers) == 0 {
				accNumbers = uint32(1)
			} else {
				accNumbers = byteutil.BytesToUint32(bAccNumbers) + 1
			}
			index = accNumbers

			// step2: save account numbers at target epoch
			state.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(accNumbers))

			// step3: save account at target index
			releaseAccountKey := restriting.GetReleaseAccountKey(epoch, index)
			state.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

			// step4: save restricting amount at target epoch
			releaseAmountKey := restriting.GetReleaseAmountKey(epoch, account)

			state.SetState(account, releaseAmountKey, amount.Bytes())

			epochList = append(epochList, epoch)
		}

		info.balance = totalAmount
		info.debt = big.NewInt(0)
		info.releaseList = epochList

	} else {
		log.Debug("restricting record exist", "account", account.Bytes())

		if err = rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
			log.Error("failed to rlp decode the account of release to", "err", err.Error())
			return common.NewSysError(err.Error())
		}

		for i := 0; i < len(plans); i++ {
			epoch := plans[i].Epoch
			amount := plans[i].Amount

			// step1: get restricting amount at target epoch
			releaseAmountKey := restriting.GetReleaseAmountKey(epoch, account)
			bAmount := state.GetState(account, releaseAmountKey)

			if len(bAmount) == 0 {
				log.Trace("release record not exist on curr epoch ", "account", account, "epoch", epoch)

				releaseEpochKey := restriting.GetReleaseEpochKey(epoch)
				bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)

				if len(bAccNumbers) == 0 {
					accNumbers = uint32(1)
				} else {
					accNumbers = byteutil.BytesToUint32(bAccNumbers) + 1
				}
				index = accNumbers

				// step2: save account numbers at target epoch
				state.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(accNumbers))

				// step3: save account at target index
				releaseAccountKey := restriting.GetReleaseAccountKey(epoch, index)

				state.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

				info.releaseList = append(info.releaseList, epoch)

			} else {
				log.Trace("release record exist at curr epoch", "account", account, "epoch", epoch)

				var origAmount *big.Int
				origAmount = origAmount.SetBytes(bAmount)
				amount = amount.Add(amount, origAmount)
			}

			// step4: save restricting amount at target epoch
			state.SetState(account, releaseAmountKey, amount.Bytes())
		}

		info.balance = info.balance.Add(info.balance, totalAmount)
	}

	// step5: save restricting account info at target epoch
	bAccInfo, err = rlp.EncodeToBytes(info)
	if err != nil {
		log.Error("failed to rlp encode restricting info", "account", account, "error", err)
		return common.NewSysError(err.Error())
	}
	state.SetState(account, restrictingKey, bAccInfo)

	return nil
}

// PledgeLockFunds transfer the money from the restricting contract account to the staking contract account,
// the first output returns true when business is success, else return false
func (rp *RestrictingPlugin) PledgeLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey := restriting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Error("record not found in PledgeLockFunds", "account", account, "funds", amount.Uint64())
		return errAccountNotFound
	}

	var (
		err  error
		info restrictingInfo
	)

	if err := rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
		log.Error("rlp decode failed", "error", err.Error())
		return common.NewSysError(err.Error())
	}

	if info.balance.Cmp(amount) == -1 {
		log.Error("restricting balance is not enough", "error", errBalanceNotEnough)
		return errBalanceNotEnough
	}

	// sub balance
	info.balance = info.balance.Sub(info.balance, amount)

	// save restricting account info
	if bAccInfo, err = rlp.EncodeToBytes(info); err != nil {
		log.Error("failed to encode restricting info", "account", account, "error", err)
		return common.NewSysError(err.Error())
	}
	state.SetState(account, restrictingKey, bAccInfo)

	state.SubBalance(vm.RestrictingContractAddr, amount)
	state.AddBalance(vm.StakingContractAddr, amount)

	return nil
}

// ReturnLockFunds transfer the money from the staking contract account  to the restricting contract account,
func (rp *RestrictingPlugin) ReturnLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey := restriting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Error("record not found in ReturnLockFunds", "account", account, "funds", amount.Uint64())
		return errAccountNotFound
	}

	var (
		err  error
		info restrictingInfo
	)

	if err = rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
		log.Error("rlp decode failed", "error", err.Error())
		return common.NewSysError(err.Error())
	}

	if info.debt.Sign() == -1 {
		log.Trace("debt is negative", "account", account, "debt", info.debt, "amount", amount)

		if amount.CmpAbs(info.debt) < 1 {
			info.debt = info.debt.Add(info.debt, amount)
		} else {
			info.balance.Add(info.balance, info.debt.Add(info.debt, amount))
			info.debt = big.NewInt(0)
		}

	} else {
		// directly add balance while debt ge 0
		log.Trace("debt is not negative", "account", account, "debt", info.debt)
		info.balance = info.balance.Add(info.balance, amount)
	}

	// save restricting account info
	if bAccInfo, err = rlp.EncodeToBytes(info); err != nil {
		log.Error("failed to encode restricting info", "account", account, "error", err)
		return common.NewSysError(err.Error())
	}
	state.SetState(account, restrictingKey, bAccInfo)

	state.SubBalance(vm.StakingContractAddr, amount)
	state.AddBalance(vm.RestrictingContractAddr, amount)

	return nil
}

// SlashingNotify modify debt of restricting account
func (rp *RestrictingPlugin) SlashingNotify(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey := restriting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Error("record not found in SlashingNotify", "account", account, "funds", amount.Uint64())
		return errAccountNotFound
	}

	var (
		err  error
		info restrictingInfo
	)

	if err = rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
		log.Error("failed to rlp decode restricting info", "error", err.Error())
		return common.NewSysError(err.Error())
	}

	info.debt = info.debt.Sub(info.debt, amount)

	// save restricting account info
	if bAccInfo, err = rlp.EncodeToBytes(info); err != nil {
		log.Error("failed to encode restricting info", "account", account, "error", err)
		return common.NewSysError(err.Error())
	}
	state.SetState(account, restrictingKey, bAccInfo)

	return nil
}

// releaseRestricting will release restricting plans on target epoch
func (rp *RestrictingPlugin) releaseRestricting(epoch uint64, state xcom.StateDB) error {

	releaseEpochKey := restriting.GetReleaseEpochKey(epoch)
	bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)

	if len(bAccNumbers) == 0 {
		log.Debug("there is no release record on curr epoch", "epoch", epoch)
		return nil
	}

	numbers := byteutil.BytesToUint32(bAccNumbers)

	// TODO
	var (
		info    restrictingInfo
		release *big.Int
	)

	for index := numbers; index > 0; index-- {

		releaseAccountKey := restriting.GetReleaseAccountKey(epoch, index)
		bAccount := state.GetState(vm.RestrictingContractAddr, releaseAccountKey)
		account := byteutil.BytesToAddress(bAccount)

		restrictingKey := restriting.GetRestrictingKey(account)
		bAccInfo := state.GetState(account, restrictingKey)

		if err := rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
			log.Error("failed to rlp decode restricting info", "error", err.Error())
			return common.NewSysError(err.Error())
		}

		releaseAmountKey := restriting.GetReleaseAmountKey(epoch, account)
		bRelease := state.GetState(account, releaseAmountKey)
		release = release.SetBytes(bRelease)

		if info.balance.Uint64() == 0 {
			log.Trace("balance of release account equals zero", "account", account, "epoch", epoch)
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
					tmpBalance := info.balance
					release = release.Sub(release, info.balance)
					info.balance = big.NewInt(0)
					info.debt = info.debt.Sub(info.debt, release)

					state.SubBalance(vm.RestrictingContractAddr, tmpBalance)
					state.AddBalance(account, tmpBalance)
				}

			} else {
				log.Trace("debt of release account le zero", "account", account, "debt", info.debt)

				if release.Cmp(info.balance) <= 0 {
					info.balance = info.balance.Sub(info.balance, release)

					state.SubBalance(vm.RestrictingContractAddr, release)
					state.AddBalance(account, release)

				} else {
					tmpBalance := info.balance
					release = release.Sub(release, info.balance)
					info.balance = big.NewInt(0)
					info.debt = info.debt.Sub(info.debt, release)

					state.SubBalance(vm.RestrictingContractAddr, tmpBalance)
					state.AddBalance(account, tmpBalance)

				}
			}
		}


		// delete ReleaseAmount
		state.SetState(account, releaseAmountKey, []byte{})

		// delete ReleaseAccount
		state.SetState(vm.RestrictingContractAddr, releaseAccountKey, []byte{})

		// delete epoch in releaseList
		for i, target := range info.releaseList {
			if target == epoch {
				info.releaseList = append(info.releaseList[:i], info.releaseList[i+1:]...)
				break
			}
		}
	}

	// delete ReleaseEpoch
	state.SetState(vm.RestrictingContractAddr, releaseEpochKey, []byte{})

	return nil
}


func (rp *RestrictingPlugin) GetRestrictingInfo(account common.Address, state xcom.StateDB) ([]byte, error) {

	restrictingKey := restriting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Error("record not found in GetRestrictingInfo", "account", account)
		return []byte{}, errAccountNotFound
	}

	var (
		amount           *big.Int
		bAmount          []byte
		info             restrictingInfo
		plan             releaseAmountInfo
		plans            []releaseAmountInfo
		releaseAmountKey []byte
		result           Result
	)

	if err := rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
		log.Error("rlp decode failed", "error", err.Error())
		return []byte{}, common.NewSysError(err.Error())
	}

	for i := 0; i < len(info.releaseList); i++ {
		epoch := info.releaseList[i]

		releaseAmountKey = restriting.GetReleaseAmountKey(epoch, account)
		bAmount = state.GetState(account, releaseAmountKey)
		amount = amount.SetBytes(bAmount)

		plan.height = epoch
		// !!!
		// plan.height = getBlockNumberByEpoch(epoch)
		// !!!
		plan.amount = amount

		plans = append(plans, plan)
	}

	bPlans, err := json.Marshal(plans)
	if err != nil {
		log.Error("falied to Marshal restricting result")
		return []byte{}, err
	}

	// !!!
	// getslash
	// getpledge
	// !!!

	result.balance = info.balance
	result.debt = info.debt
	result.slash = big.NewInt(0)
	result.pledge = big.NewInt(0)
	result.entry = bPlans

	log.Trace("get restricting result", "account", account, "result", result)

	return rlp.EncodeToBytes(result)
}
