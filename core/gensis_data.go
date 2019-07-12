package core

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
)

func genesisStakingData(g *Genesis, db *staking.StakingDB, genesisHash common.Hash, version uint32) error {


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
			Shares:             xcom.StakeThreshold,
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

		if err = db.SetCandidateStore(genesisHash, nodeAddr, can); err != nil {
			return fmt.Errorf("Failed to Store Candidate Info. ID:%v, error:%s", can.NodeId, err)
		}

		if err = db.SetCanPowerStore(genesisHash, nodeAddr, can); err != nil {
			return fmt.Errorf("Failed to Store Candidate Power. ID:%v, error:%s", can.NodeId, err)
		}

		if err != nil {
			return fmt.Errorf("failed to exchange nodeID to address. ID:%v, error:%s", nodeID, err)
		}


		// build validator queue for the first consensus epoch
		validator := &staking.Validator{
			NodeAddress:   nodeAddr,
			NodeId:        node.ID,
			StakingWeight: [staking.SWeightItem]string{fmt.Sprint(version), xcom.StakeThreshold.String(), "0", fmt.Sprint(index+1)},
			ValidatorTerm: 0,
		}
		validatorQueue = append(validatorQueue, validator)

	}


	// build epoch validators
	verifierList := &staking.Validator_array{
		Start: 1,
		End:   xcom.EpochSize * xcom.ConsensusSize,
		Arr:   validatorQueue,
	}

	// build current validators
	validatorLIst := &staking.Validator_array{
		Start: 1,
		End:   xcom.ConsensusSize,
		Arr:   validatorQueue,
	}

	if err := db.SetVerfierList(genesisHash, verifierList); err != nil {
		return fmt.Errorf("Failed to Store Epoch Validators. error:%s", err)
	}

	if err := db.SetCurrentValidatorList(genesisHash, validatorLIst); err != nil {
		return fmt.Errorf("Failed to Store Current Round Validators. error:%s", err)
	}

	return nil
}



// buildAllowancePlan writes the data of precompiled restricting contract, which used for the second year allowance
// and the third year allowance, to stateDB
func buildAllowancePlan(stateDb *state.StateDB) error {

	account := vm.RewardManagerPoolAddr
	firstYearEndEpoch := 365 * 24 * 3600 / (xcom.EpochSize * xcom.ConsensusSize)
	secondYearEncEpoch := 2 * 365 * 24 * 3600 / (xcom.EpochSize * xcom.ConsensusSize)
	stableEpochs := []uint64{firstYearEndEpoch, secondYearEncEpoch}

	secondYearAllowance, _ := new(big.Int).SetString("15000000000000000000000000", 10)
	thirdYearAllowance, _ := new(big.Int).SetString("5000000000000000000000000", 10)

	epochList := make([]uint64, len(stableEpochs))
	for i, epoch := range stableEpochs {
		// store restricting account record
		releaseAccountKey := restricting.GetReleaseAccountKey(epoch, 1)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

		// store release amount record
		releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
		switch {
		case i == 0:
			stateDb.SetState(account, releaseAmountKey, secondYearAllowance.Bytes())
		case i == 1:
			stateDb.SetState(account, releaseAmountKey, thirdYearAllowance.Bytes())
		}

		// store release epoch record
		releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint64ToBytes(1))

		epochList = append(epochList, uint64(epoch))
	}

	// build restricting account info
	var restrictInfo restricting.RestrictingInfo
	restrictInfo.Balance, _ = new(big.Int).SetString("20000000000000000000000000", 10)
	restrictInfo.Debt = big.NewInt(0)
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
