package core

import (
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

func genesisStakingData(g *Genesis, genesisHash common.Hash, programVersion uint32) error {

	isDone := false
	switch {
	case nil == g.Config:
		isDone = true
	case nil == g.Config.Cbft:
		isDone = true
	case len(g.Config.Cbft.InitialNodes) == 0:
		isDone = true
	}

	if isDone {
		log.Warn("Genesis StakingData, the genesis config or cbft or initialNodes is nil, Not building staking data")
		return nil
	}

	if g.Config.Cbft.ValidatorMode != common.PPOS_VALIDATOR_MODE {
		log.Info("Init staking snapshotdb data, validatorMode is not ppos")
		return nil
	}

	snapdb := snapshotdb.Instance()

	version := xutil.CalcVersion(programVersion)

	var length int

	if int(xcom.ConsValidatorNum()) <= len(g.Config.Cbft.InitialNodes) {
		length = int(xcom.ConsValidatorNum())
	} else {
		length = len(g.Config.Cbft.InitialNodes)
	}
	initQueue := g.Config.Cbft.InitialNodes

	validatorQueue := make(staking.ValidatorQueue, length)

	for index := 0; index < length; index++ {

		node := initQueue[index]

		can := &staking.Candidate{
			NodeId:             node.Node.ID,
			StakingAddress:     vm.PlatONFoundationAddress,
			BenefitAddress:     vm.RewardManagerPoolAddr,
			StakingTxIndex:     uint32(index + 1),
			ProgramVersion:     version,
			Status:             staking.Valided,
			StakingEpoch:       uint32(0),
			StakingBlockNum:    uint64(0),
			Shares:             xcom.StakeThreshold(),
			Released:           xcom.StakeThreshold(),
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,
			Description: staking.Description{
				ExternalId: "",
				NodeName:   "platon.node." + fmt.Sprint(index+1),
				Website:    "www.platon.network",
				Details:    "The PlatON Node",
			},
		}

		nodeAddr, err := xutil.NodeId2Addr(can.NodeId)
		if err != nil {
			return fmt.Errorf("Failed to convert nodeID to address. nodeId:%s, error:%s", can.NodeId.String(), err.Error())
		}

		key := staking.CandidateKeyByAddr(nodeAddr)

		if val, err := rlp.EncodeToBytes(can); nil != err {
			return fmt.Errorf("Failed to Store Candidate Info: rlp encodeing failed. nodeId:%s, error:%s", can.NodeId.String(), err.Error())
		} else {
			if err := snapdb.PutBaseDB(key, val); nil != err {
				return fmt.Errorf("Failed to Store Candidate Info: PutBaseDB failed. nodeId:%s, error:%s", can.NodeId.String(), err.Error())
			}
		}

		powerKey := staking.TallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex, can.ProgramVersion)
		if err := snapdb.PutBaseDB(powerKey, nodeAddr.Bytes()); nil != err {
			return fmt.Errorf("Failed to Store Candidate Power: PutBaseDB failed. nodeId:%s, error:%s", can.NodeId.String(), err.Error())
		}

		// build validator queue for the first consensus epoch
		validator := &staking.Validator{
			NodeAddress:   nodeAddr,
			NodeId:        node.Node.ID,
			StakingWeight: [staking.SWeightItem]string{fmt.Sprint(version), xcom.StakeThreshold().String(), "0", fmt.Sprint(index + 1)},
			ValidatorTerm: 0,
		}
		validatorQueue[index] = validator

	}

	// store the account staking Reference Count
	err := snapdb.PutBaseDB(staking.GetAccountStakeRcKey(vm.PlatONFoundationAddress), common.Uint64ToBytes(uint64(length)))
	if nil != err {
		return fmt.Errorf("Failed to Store Staking Account Reference Count. account: %s, error:%s", vm.PlatONFoundationAddress.Hex(), err.Error())
	}

	validatorArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		return fmt.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
	}

	/**
	Epoch
	*/
	// build epoch validators indexInfo
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}
	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		return fmt.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
	}
	if err := snapdb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		return fmt.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
	}

	// Epoch validators
	if err := snapdb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), validatorArr); nil != err {
		return fmt.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
	}

	/**
	Round
	*/
	// build previous validators indexInfo
	pre_indexInfo := &staking.ValArrIndex{
		Start: 0,
		End:   0,
	}
	// build current validators indexInfo
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, pre_indexInfo)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		return fmt.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
	}
	if err := snapdb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		return fmt.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
	}

	// Previous Round validator
	if err := snapdb.PutBaseDB(staking.GetRoundValArrKey(pre_indexInfo.Start, pre_indexInfo.End), validatorArr); nil != err {
		return fmt.Errorf("Failed to Store Previous Round Validators: PutBaseDB failed. error:%s", err.Error())
	}
	// Current Round validator
	if err := snapdb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), validatorArr); nil != err {
		return fmt.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
	}

	if err := snapdb.SetCurrent(genesisHash, *common.Big0, *common.Big0); nil != err {
		return fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error())
	}
	return nil
}

// genesisAllowancePlan writes the data of precompiled restricting contract, which used for the second year allowance
// and the third year allowance, to stateDB
func genesisAllowancePlan(stateDb *state.StateDB, issue *big.Int) error {

	account := vm.RewardManagerPoolAddr

	OneYearEpochs := xutil.EpochsPerYear()
	stableEpochs := []uint64{OneYearEpochs, 2 * OneYearEpochs}

	epochList := make([]uint64, len(stableEpochs))
	for i, epoch := range stableEpochs {
		// store restricting account record
		releaseAccountKey := restricting.GetReleaseAccountKey(epoch, 1)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

		// store release amount record
		releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
		// At the first year end, 1.5% of genesis issuance need be released
		// At the twice year end, 0.5% of genesis issuance need be released
		switch {
		case i == 0:
			allowance := new(big.Int).Mul(issue, big.NewInt(15))
			allowance = allowance.Div(allowance, big.NewInt(1000))
			stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, allowance.Bytes())
		case i == 1:
			allowance := new(big.Int).Mul(issue, big.NewInt(5))
			allowance = allowance.Div(allowance, big.NewInt(1000))
			stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, allowance.Bytes())
		}

		// store release epoch record
		releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		epochList[i] = uint64(epoch)
	}

	// build restricting account info
	var restrictInfo restricting.RestrictingInfo
	restrictInfo.Balance = stateDb.GetBalance(vm.RestrictingContractAddr)
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

func genesisPluginState(g *Genesis, statedb *state.StateDB, genesisReward, genesisIssue *big.Int, programVersion uint32) error {

	isDone := false
	switch {
	case nil == g.Config:
		isDone = true
	case nil == g.Config.Cbft:
		isDone = true
	}

	if isDone {
		log.Warn("Genesis xxPlugin statedb, the genesis config or cbft is nil, Not Store plugin genesis state")
		return nil
	}

	if g.Config.Cbft.ValidatorMode != common.PPOS_VALIDATOR_MODE {
		log.Info("Init xxPlugin genesis statedb, validatorMode is not ppos")
		return nil
	}

	// Store genesis yearEnd reward balance item
	plugin.SetYearEndBalance(statedb, 0, genesisReward)

	// Store genesis Issue for LAT
	plugin.SetYearEndCumulativeIssue(statedb, 0, genesisIssue)

	log.Info("Store version for gov into genesis statedb", "real version", fmt.Sprintf("%d.%d.%d",
		params.VersionMajor, params.VersionMinor, params.VersionPatch), "uint32 version", programVersion)

	// Store genesis governance data
	statedb.SetState(vm.GovContractAddr, gov.KeyActiveVersion(), common.Uint32ToBytes(programVersion))
	// Store restricting plans for increase issue for second and third year
	if err := genesisAllowancePlan(statedb, genesisIssue); nil != err {
		return err
	}
	// Store genesis last Epoch
	log.Info("Set latest epoch", "blockNumber", g.Number, "epoch", 0)
	plugin.SetLatestEpoch(statedb, uint64(0))
	return nil
}
