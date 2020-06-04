package core

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

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

func genesisStakingData(prevHash common.Hash, snapdb snapshotdb.BaseDB, g *Genesis, stateDB *state.StateDB) (common.Hash, error) {

	if g.Config.Cbft.ValidatorMode != common.PPOS_VALIDATOR_MODE {
		log.Info("Init staking snapshotdb data, validatorMode is not ppos")
		return prevHash, nil
	}

	var length int

	if int(xcom.MaxConsensusVals()) <= len(g.Config.Cbft.InitialNodes) {
		length = int(xcom.MaxConsensusVals())
	} else {
		length = len(g.Config.Cbft.InitialNodes)
	}

	// Check the balance of Staking Account
	needStaking := new(big.Int).Mul(xcom.GeneStakingAmount, big.NewInt(int64(length)))
	remain := stateDB.GetBalance(xcom.CDFAccount())

	if remain.Cmp(needStaking) < 0 {
		return prevHash, fmt.Errorf("Failed to store genesis staking data, the balance of '%s' is no enough. "+
			"balance: %s, need staking: %s", xcom.CDFAccount().String(), remain.String(), needStaking.String())
	}

	initQueue := g.Config.Cbft.InitialNodes

	validatorQueue := make(staking.ValidatorQueue, length)

	lastHash := prevHash

	putbasedbFn := func(key, val []byte, hash common.Hash) (common.Hash, error) {
		if err := snapdb.PutBaseDB(key, val); nil != err {
			return hash, err
		}
		newHash := common.GenerateKVHash(key, val, hash)
		return newHash, nil
	}

	for index := 0; index < length; index++ {

		node := initQueue[index]

		var keyHex bls.PublicKeyHex
		if b, err := node.BlsPubKey.MarshalText(); nil != err {
			return lastHash, err
		} else {
			if err := keyHex.UnmarshalText(b); nil != err {
				return lastHash, err
			}
		}

		base := &staking.CandidateBase{
			NodeId:          node.Node.ID,
			BlsPubKey:       keyHex,
			StakingAddress:  xcom.CDFAccount(),
			BenefitAddress:  vm.RewardManagerPoolAddr,
			StakingTxIndex:  uint32(index),           // txIndex from zero to n
			ProgramVersion:  g.Config.GenesisVersion, // genesis version
			StakingBlockNum: uint64(0),
			Description: staking.Description{
				ExternalId: "",
				NodeName:   "platon.node." + fmt.Sprint(index+1),
				Website:    "www.platon.network",
				Details:    "The PlatON Node",
			},
		}

		mutable := &staking.CandidateMutable{
			Status:             staking.Valided,
			StakingEpoch:       uint32(0),
			Shares:             new(big.Int).Set(xcom.GeneStakingAmount),
			Released:           new(big.Int).Set(xcom.GeneStakingAmount),
			ReleasedHes:        new(big.Int).SetInt64(0),
			RestrictingPlan:    new(big.Int).SetInt64(0),
			RestrictingPlanHes: new(big.Int).SetInt64(0),
		}

		nodeAddr, err := xutil.NodeId2Addr(base.NodeId)
		if err != nil {
			return lastHash, fmt.Errorf("Failed to convert nodeID to address. nodeId:%s, error:%s",
				base.NodeId.String(), err.Error())
		}

		// about CanBase ...
		baseKey := staking.CanBaseKeyByAddr(nodeAddr)
		if val, err := rlp.EncodeToBytes(base); nil != err {
			return lastHash, fmt.Errorf("Failed to Store CanBase Info: rlp encodeing failed. nodeId:%s, error:%s",
				base.NodeId.String(), err.Error())
		} else {

			lastHash, err = putbasedbFn(baseKey, val, lastHash)
			if nil != err {
				return lastHash, fmt.Errorf("Failed to Store CanBase Info: PutBaseDB failed. nodeId:%s, error:%s",
					base.NodeId.String(), err.Error())
			}

		}

		// about CanMutable ...
		mutableKey := staking.CanMutableKeyByAddr(nodeAddr)
		if val, err := rlp.EncodeToBytes(mutable); nil != err {
			return lastHash, fmt.Errorf("Failed to Store CanMutable Info: rlp encodeing failed. nodeId:%s, error:%s",
				base.NodeId.String(), err.Error())
		} else {

			lastHash, err = putbasedbFn(mutableKey, val, lastHash)
			if nil != err {
				return lastHash, fmt.Errorf("Failed to Store CanMutable Info: PutBaseDB failed. nodeId:%s, error:%s",
					base.NodeId.String(), err.Error())
			}

		}

		// about can power ...
		powerKey := staking.TallyPowerKey(base.ProgramVersion, mutable.Shares, base.NodeId, base.StakingBlockNum, base.StakingTxIndex)
		lastHash, err = putbasedbFn(powerKey, nodeAddr.Bytes(), lastHash)
		if nil != err {
			return lastHash, fmt.Errorf("Failed to Store Candidate Power: PutBaseDB failed. nodeId:%s, error:%s",
				base.NodeId.String(), err.Error())
		}

		// build validator queue for the first consensus epoch
		validator := &staking.Validator{
			NodeAddress:     nodeAddr,
			NodeId:          base.NodeId,
			BlsPubKey:       base.BlsPubKey,
			ProgramVersion:  base.ProgramVersion, // real version
			Shares:          mutable.Shares,
			StakingBlockNum: base.StakingBlockNum,
			StakingTxIndex:  base.StakingTxIndex,
			ValidatorTerm:   0,
		}
		validatorQueue[index] = validator

		stateDB.SubBalance(xcom.CDFAccount(), new(big.Int).Set(xcom.GeneStakingAmount))
		stateDB.AddBalance(vm.StakingContractAddr, new(big.Int).Set(xcom.GeneStakingAmount))
	}

	// store the account staking Reference Count
	lastHash, err := putbasedbFn(staking.GetAccountStakeRcKey(xcom.CDFAccount()), common.Uint64ToBytes(uint64(length)), lastHash)
	if nil != err {
		return lastHash, fmt.Errorf("Failed to Store Staking Account Reference Count. account: %s, error:%s",
			xcom.CDFAccount().String(), err.Error())
	}

	validatorArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		return lastHash, fmt.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
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
		return lastHash, fmt.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
	}

	lastHash, err = putbasedbFn(staking.GetEpochIndexKey(), epoch_index, lastHash)
	if nil != err {
		return lastHash, fmt.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
	}

	// Epoch validators
	lastHash, err = putbasedbFn(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), validatorArr, lastHash)
	if nil != err {
		return lastHash, fmt.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
	}

	/**
	Round
	*/
	// build previous round validators indexInfo
	pre_indexInfo := &staking.ValArrIndex{
		Start: 0,
		End:   0,
	}
	// build current round validators indexInfo
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
		return lastHash, fmt.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
	}
	lastHash, err = putbasedbFn(staking.GetRoundIndexKey(), round_index, lastHash)
	if nil != err {
		return lastHash, fmt.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
	}

	// Previous Round validator
	lastHash, err = putbasedbFn(staking.GetRoundValArrKey(pre_indexInfo.Start, pre_indexInfo.End), validatorArr, lastHash)
	if nil != err {
		return lastHash, fmt.Errorf("Failed to Store Previous Round Validators: PutBaseDB failed. error:%s", err.Error())
	}

	// Current Round validator
	lastHash, err = putbasedbFn(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), validatorArr, lastHash)
	if nil != err {
		return lastHash, fmt.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
	}

	log.Info("Call genesisStakingData, Store genesis pposHash by stake data", "pposHash", lastHash.Hex())

	stateDB.SetState(vm.StakingContractAddr, staking.GetPPOSHASHKey(), lastHash.Bytes())

	return lastHash, nil
}

// genesisAllowancePlan writes the data of precompiled restricting contract, which used for the second year allowance
// and the third year allowance, to stateDB
func genesisAllowancePlan(statedb *state.StateDB) error {

	account := vm.RewardManagerPoolAddr
	var (
		zeroYear  = new(big.Int).Mul(big.NewInt(62215742), big.NewInt(1e18))
		oneYear   = new(big.Int).Mul(big.NewInt(55965742), big.NewInt(1e18))
		twoYear   = new(big.Int).Mul(big.NewInt(49559492), big.NewInt(1e18))
		threeYear = new(big.Int).Mul(big.NewInt(42993086), big.NewInt(1e18))
		fourYear  = new(big.Int).Mul(big.NewInt(36262520), big.NewInt(1e18))
		fiveYear  = new(big.Int).Mul(big.NewInt(29363689), big.NewInt(1e18))
		sixYear   = new(big.Int).Mul(big.NewInt(22292388), big.NewInt(1e18))
		sevenYear = new(big.Int).Mul(big.NewInt(15044304), big.NewInt(1e18))
		eightYear = new(big.Int).Mul(big.NewInt(7615018), big.NewInt(1e18))
	)

	statedb.SubBalance(xcom.CDFAccount(), zeroYear)
	statedb.AddBalance(account, zeroYear)
	needRelease := []*big.Int{oneYear, twoYear, threeYear, fourYear, fiveYear, sixYear, sevenYear, eightYear}

	restrictingPlans := make([]restricting.RestrictingPlan, 0)
	OneYearEpochs := xutil.EpochsPerYear()
	for key, value := range needRelease {
		epochs := OneYearEpochs * (uint64(key) + 1)
		restrictingPlans = append(restrictingPlans, restricting.RestrictingPlan{epochs, value})
	}

	if err := plugin.RestrictingInstance().AddRestrictingRecord(xcom.CDFAccount(), vm.RewardManagerPoolAddr, 0, restrictingPlans, statedb); err != nil {
		return err
	}
	return nil
}

func genesisPluginState(g *Genesis, statedb *state.StateDB, snapDB snapshotdb.BaseDB, genesisIssue *big.Int) error {

	if g.Config.Cbft.ValidatorMode != common.PPOS_VALIDATOR_MODE {
		log.Info("Init xxPlugin genesis statedb, validatorMode is not ppos")
		return nil
	}

	// Store genesis yearEnd reward balance item

	// Store genesis Issue for LAT
	plugin.SetYearEndCumulativeIssue(statedb, 0, genesisIssue)

	log.Info("Write genesis version into genesis block", "genesis version", fmt.Sprintf("%d/%s", g.Config.GenesisVersion, params.FormatVersion(g.Config.GenesisVersion)))

	// Store genesis governance data
	activeVersionList := []gov.ActiveVersionValue{
		{ActiveVersion: g.Config.GenesisVersion, ActiveBlock: 0},
	}
	activeVersionListBytes, _ := json.Marshal(activeVersionList)
	statedb.SetState(vm.GovContractAddr, gov.KeyActiveVersions(), activeVersionListBytes)

	err := plugin.NewRestrictingPlugin(nil).InitGenesisRestrictingPlans(statedb)
	if err != nil {
		return fmt.Errorf("Failed to init genesis restricting plans, err:%s", err.Error())
	}
	genesisReward := statedb.GetBalance(vm.RewardManagerPoolAddr)
	plugin.SetYearEndBalance(statedb, 0, genesisReward)
	log.Info("Set SetYearEndBalance", "genesisReward", genesisReward)

	return nil
}

func genesisGovernParamData(prevHash common.Hash, snapdb snapshotdb.BaseDB, genesisVersion uint32) (common.Hash, error) {
	return gov.InitGenesisGovernParam(prevHash, snapdb, genesisVersion)
}

func hashEconomicConfig(economicModel *xcom.EconomicModel, prevHash common.Hash) (common.Hash, error) {
	if economicModel != nil {
		bytes, err := rlp.EncodeToBytes(economicModel)
		if err != nil {
			return prevHash, err
		}
		prevHash = common.GenerateKVHash([]byte("economicConfig"), bytes, prevHash)
	}
	return prevHash, nil
}
