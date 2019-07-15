package core

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
)

func genesisStakingData(g *Genesis, genesisHash common.Hash, version uint32) error {

	snapdb := snapshotdb.Instance()

	validatorQueue := make(staking.ValidatorQueue, len(g.Config.Cbft.InitialNodes))
	for index, node := range g.Config.Cbft.InitialNodes {

		can := &staking.Candidate{
			NodeId:             node.ID,
			StakingAddress:     vm.RewardManagerPoolAddr,
			BenifitAddress:     vm.RewardManagerPoolAddr,
			StakingTxIndex:     uint32(index+1),
			ProcessVersion:     version,
			Status:             staking.Valided,
			StakingEpoch:       uint32(0),
			StakingBlockNum:    uint64(0),
			Shares:             xcom.StakeThreshold(),
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,
			Description: staking.Description{
				ExternalId: "",
				NodeName:   "platon.node."+ fmt.Sprint(index+1),
				Website:    "www.platon.network",
				Details:    "The PlatON Node",
			},
		}

		nodeAddr, err := xutil.NodeId2Addr(can.NodeId)
		if err != nil {
			return fmt.Errorf("Failed to convert nodeID to address. ID:%v, error:%s", can.NodeId, err)
		}



		key := staking.CandidateKeyByAddr(nodeAddr)

		if val, err := rlp.EncodeToBytes(can); nil != err {
			return fmt.Errorf("Failed to Store Candidate Info: rlp encodeing failed. ID:%v, error:%s", can.NodeId, err)
		}else {
			if err := snapdb.PutBaseDB(key, val); nil != err {
				return fmt.Errorf("Failed to Store Candidate Info: PutBaseDB failed. ID:%v, error:%s", can.NodeId, err)
			}
		}


		powerKey := staking.TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex, can.ProcessVersion)
		if err := snapdb.PutBaseDB(powerKey, nodeAddr.Bytes()); nil != err {
			return fmt.Errorf("Failed to Store Candidate Power: PutBaseDB failed. ID:%v, error:%s", can.NodeId, err)
		}


		// build validator queue for the first consensus epoch
		validator := &staking.Validator{
			NodeAddress:   nodeAddr,
			NodeId:        node.ID,
			StakingWeight: [staking.SWeightItem]string{fmt.Sprint(version), xcom.StakeThreshold().String(), "0", fmt.Sprint(index+1)},
			ValidatorTerm: 0,
		}
		validatorQueue[index] = validator

	}


	// build epoch validators
	verifierList := &staking.Validator_array{
		Start: 1,
		End:   xcom.EpochSize() * xcom.ConsensusSize(),
		Arr:   validatorQueue,
	}


	// build current validators
	validatorLIst := &staking.Validator_array{
		Start: 1,
		End:   xcom.ConsensusSize(),
		Arr:   validatorQueue,
	}

	// build pre validators
	pre_validatorLIst := &staking.Validator_array{
		Start: 0,
		End:   0,
		Arr:   validatorQueue,
	}


	// current epoch
	verifiers, err := rlp.EncodeToBytes(verifierList)
	if nil != err {
		return fmt.Errorf("Failed to Store Epoch Validators: rlp encodeing failed. error:%s", err)
	}
	if err := snapdb.PutBaseDB(staking.GetEpochValidatorKey(), verifiers); nil != err {
		return fmt.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err)
	}


	// pre round
	pre_vals, err := rlp.EncodeToBytes(pre_validatorLIst)
	if nil != err {
		return fmt.Errorf("Failed to Store Previous Round Validators: rlp encodeing failed. error:%s", err)
	}
	if err := snapdb.PutBaseDB(staking.GetPreRoundValidatorKey(), pre_vals); nil != err {
		return fmt.Errorf("Failed to Store Previous Round Validators: PutBaseDB failed. error:%s", err)
	}

	// current round
	vals, err := rlp.EncodeToBytes(validatorLIst)
	if nil != err {
		return fmt.Errorf("Failed to Store Current Round Validators: rlp encodeing failed. error:%s", err)
	}
	if err := snapdb.PutBaseDB(staking.GetCurRoundValidatorKey(), vals); nil != err {
		return fmt.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err)
	}

	if err := snapdb.SetCurrent(genesisHash, *common.Big0, *common.Big0); nil != err {
		return fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err)
	}
	return nil
}


// buildAllowancePlan writes the data of precompiled restricting contract, which used for the second year allowance
// and the third year allowance, to stateDB
func buildAllowancePlan(stateDb *state.StateDB) error {

	account := vm.RewardManagerPoolAddr

	stableEpochs := []uint64{xcom.FirstYearEndEpoch(), xcom.SecondYearEncEpoch()}

	epochList := make([]uint64, len(stableEpochs))
	for i, epoch := range stableEpochs {
		// store restricting account record
		releaseAccountKey := restricting.GetReleaseAccountKey(epoch, 1)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

		// store release amount record
		releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
		switch {
		case i == 0:
			stateDb.SetState(account, releaseAmountKey, xcom.SecondYearAllowance().Bytes())
		case i == 1:
			stateDb.SetState(account, releaseAmountKey, xcom.ThirdYearAllowance().Bytes())
		}

		// store release epoch record
		releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint64ToBytes(1))

		epochList = append(epochList, uint64(epoch))
	}

	// build restricting account info
	var restrictInfo restricting.RestrictingInfo
	restrictInfo.Balance = xcom.GenesisRestrictingBalance()
	restrictInfo.Debt = big.NewInt(0)
	restrictInfo.DebtSymbol = false
	restrictInfo.ReleaseList = epochList

	bRestrictInfo, err := rlp.EncodeToBytes(restrictInfo)
	if err != nil {
		return fmt.Errorf("failed to rlp encode restricting info. info:%v, error:%s", restrictInfo, err.Error())
	}

	// store restricting account info
	restrictingKey := restricting.GetRestrictingKey(account)
	stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bRestrictInfo)

	return nil
}
