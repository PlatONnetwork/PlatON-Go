package xcom

import "math/big"

// plugin rule key
const (
	DefualtRule = iota
	StakingRule
	SlashingRule
	RestrictingRule
	RewardRule
	GovernanceRule

	// ......
)

// config  TODO Configuration is all here
type EconomicModel struct {
	staking       StakingConfig
	slashing      SlashingConfig
	reward        RewardConfig
	gov           GovernanceConfig
	epochsPerYear uint32
}

var ec *EconomicModel

func SetEconomicModel(ecParams *EconomicModel) {
	if nil == ec {
		ec = ecParams
	}
}
// Getting the global EconomicModel single instance
//func GetEc() *EconomicModel {
//	return ec
//}


// Staking config
func  StakeThreshold() *big.Int {
	return ec.staking.StakeThreshold
}

func  DelegateThreshold () *big.Int {
	return ec.staking.DelegateThreshold
}

func  ConsValidatorNum () uint64 {
	return ec.staking.ConsValidatorNum
}

func  EpochValidatorNum () uint64 {
	return ec.staking.EpochValidatorNum
}

func  ShiftValidatorNum () uint64 {
	return ec.staking.ShiftValidatorNum
}

func  EpochSize () uint64 {
	return ec.staking.EpochSize
}

func  HesitateRatio () uint64 {
	return ec.staking.HesitateRatio
}

func  EffectiveRatio () uint64 {
	return ec.staking.EffectiveRatio
}

func  ElectionDistance () uint64 {
	return ec.staking.ElectionDistance
}


func  ConsensusSize () uint64 {
	return ec.staking.ConsensusSize
}

func UnStakeFreezeRatio () uint64 {
	return ec.staking.UnStakeFreezeRatio
}

func  PassiveUnDelFreezeRatio () uint64 {
	return ec.staking.PassiveUnDelegateFreezeRatio
}

func  ActiveUnDelFreezeRatio () uint64 {
	return ec.staking.ActiveUnDelegateFreezeRatio
}


// Slashing config
func BlockAmountLow () uint32 {
	return ec.slashing.BlockAmountLow
}

func  BlockAmountHigh () uint32 {
	return ec.slashing.BlockAmountHigh
}

func  BlockAmountLowSlash () uint32 {
	return ec.slashing.BlockAmountLowSlashing
}

func  BlockAmountHighSlash () uint32 {
	return ec.slashing.BlockAmountHighSlashing
}

func  DuplicateSignNum () uint32 {
	return ec.slashing.DuplicateSignNum
}

func  DuplicateSignLowSlash () uint32 {
	return ec.slashing.DuplicateSignLowSlashing
}

func  DuplicateSignHighSlash () uint32 {
	return ec.slashing.DuplicateSignHighSlashing
}


// Reward config
func  SecondYearAllowance () *big.Int {
	return ec.reward.SecondYearAllowance
}

func  ThirdYearAllowance () *big.Int {
	return ec.reward.ThirdYearAllowance
}

func  GenesisRestrictingBalance () *big.Int {
	return ec.reward.GenesisRestrictingBalance
}

func  FirstYearEndEpoch () uint64 {
	return ec.reward.FirstYearEndEpoch
}

func  SecondYearEncEpoch () uint64 {
	return ec.reward.SecondYearEncEpoch
}

// Governance config
func  SupportRateThreshold () float64 {
	return ec.gov.SupportRate_Threshold
}
func  MaxVotingDuration() uint64 {
	return ec.gov.MaxVotingDuration
}





type StakingConfig struct {
	StakeThreshold               *big.Int
	DelegateThreshold            *big.Int
	ConsValidatorNum             uint64
	EpochValidatorNum            uint64
	ShiftValidatorNum            uint64
	EpochSize                    uint64
	HesitateRatio                uint64
	EffectiveRatio               uint64
	ElectionDistance             uint64
	ConsensusSize                uint64
	UnStakeFreezeRatio           uint64
	PassiveUnDelegateFreezeRatio uint64
	ActiveUnDelegateFreezeRatio  uint64
}

type SlashingConfig struct {
	BlockAmountLow            uint32
	BlockAmountHigh           uint32
	BlockAmountLowSlashing    uint32
	BlockAmountHighSlashing   uint32
	DuplicateSignNum          uint32
	DuplicateSignLowSlashing  uint32
	DuplicateSignHighSlashing uint32
}

type RewardConfig struct {
	SecondYearAllowance       *big.Int
	ThirdYearAllowance        *big.Int
	GenesisRestrictingBalance *big.Int
	FirstYearEndEpoch         uint64
	SecondYearEncEpoch        uint64
}

type RestrictingConfig struct {
	// do nothings
}

type GovernanceConfig struct {
	SupportRate_Threshold float64
	MaxVotingDuration     uint64
}

// PlatON Main Network Config
var DefaultConfig = EconomicModel{
	staking:       defaultStakingConfig,
	slashing:      defaultSlashingConfig,
	reward:        defaultRewardConfig,
	gov:           defaultGovConfig,
	epochsPerYear: 1,
}

// PlatON Alpha Test Net Config
var TestnetDefaultConfig = EconomicModel{
	staking:       testnetDefaultStakingConfig,
	slashing:      testnetDefaultSlashingConfig,
	reward:        testnetDefaultRewardConfig,
	gov:           testnetDefaultGovConfig,
	epochsPerYear: 1,
}

// PlatON Beta Test Net Config
var BetaDefaultConfig = EconomicModel{
	staking:       betaDefaultStakingConfig,
	slashing:      betaDefaultSlashingConfig,
	reward:        betaDefaultRewardConfig,
	gov:           betaDefaultGovConfig,
	epochsPerYear: 1,
}

// PlatON Inner Test Net Config
var InnerTestDefaultConfig = EconomicModel{
	staking:       innerTestDefaultStakingConfig,
	slashing:      innerTestDefaultSlashingConfig,
	reward:        innerTestDefaultRewardConfig,
	gov:           innerTestDefaultGovConfig,
	epochsPerYear: 1,
}

// PlatON Inner Dev Net Config
var InnerDevDefaultConfig = EconomicModel{
	staking:       innerDevDefaultStakingConfig,
	slashing:      innerDevDefaultSlashingConfig,
	reward:        innerDevDefaultRewardConfig,
	gov:           innerDevDefaultGovConfig,
	epochsPerYear: 1,
}

//  Dev Config
var DevConfig = EconomicModel{
	staking:       defaultStakingConfig,
	slashing:      defaultSlashingConfig,
	reward:        defaultRewardConfig,
	gov:           defaultGovConfig,
	epochsPerYear: 1,
}

/**
PlatON Main Network Defaut Config
*/
var (
	/**
	Staking config
	**/
	// The staking minimum threshold allowed (100,0000 LAT)
	stakeThreshold, _ = new(big.Int).SetString("1000000000000000000000000", 10)
	// The delegate minimum threshold allowed (10 LAT)
	delegateThreshold, _ = new(big.Int).SetString("10", 10)
	// The consensus validators count
	consValidatorNum = uint64(25)
	// The epoch (billing cycle) validators count
	epochValidatorNum = uint64(101)
	// The number of elections and replacements for each of the consensus rounds
	shiftValidatorNum = uint64(8)
	// Each epoch (billing cycle) is a multiple of the consensus rounds
	// TODO NOTE：It should be calculated by that
	//
	//      	 /  eh * 3600  \
	// C = floor |—————————————|
	//           \	L * (u*vn) /
	//
	// C: 	the epoch (just be this)
	// eh: 	the number of hours per epoch
	// L： 	each block interval (uint: seconds)
	// u: 	the consensus validators count
	// vn:  each validator has a target number of blocks per view
	epochSize = uint64(88)
	// Each hesitation period is a multiple of the epoch
	hesitateRatio = uint64(1)
	// Each effective period is a multiple of the epoch
	effectiveRatio = uint64(1)
	// The interval of the last block of the high-distance
	// consensus round of the election block for each consensus round
	electionDistance = uint64(20)
	// Number of blocks per consensus round
	// TODO NOTE: just like that
	// this = u*vn
	// u: 	the consensus validators count
	// vn:  each validator has a target number of blocks per view
	consensusSize = uint64(250)

	// The freeze period of the withdrew staking (unit is  epochs)
	unStakeFreezeRatio = uint64(1)

	// The freeze period of the delegate was invalidated
	// due to the withdrawal of the Stake (unit is  epochs)
	passiveUnDelegateFreezeRatio = uint64(0)

	// The freeze period of the delegate was invalidated
	// due to active withdrew delegate (unit is  epochs)
	activeUnDelegateFreezeRatio = uint64(0)

	/**
	Slashing
	*/
	// The number of low exceptions per consensus round
	blockAmountLow = uint32(8)
	// Number of blocks per high consensus exception
	blockAmountHigh = uint32(5)
	// Penalty quota for each consensus round with a low
	// number of abnormal blocks, percentage
	blockAmountLowSlashing = uint32(10)
	// The penalty amount for each consensus round high
	// abnormal number of blocks, percentage
	blockAmountHighSlashing = uint32(20)
	// The conditions for the highest penalty,
	// double signing
	duplicateSignNum = uint32(2)
	// Double sign low penalty amount, percentage
	duplicateSignLowSlashing = uint32(10)
	// DuplicateSignHighSlashing
	duplicateSignHighSlashing = uint32(10)

	/**
	Restricting
	TODO
	*/

	/**
	Reward
	*/
	secondYearAllowance, _       = new(big.Int).SetString("15000000000000000000000000", 10)
	thirdYearAllowance, _        = new(big.Int).SetString("5000000000000000000000000", 10)
	genesisRestrictingBalance, _ = new(big.Int).SetString("20000000000000000000000000", 10)
	firstYearEndEpoch            = 365 * 24 * 3600 / (epochSize * consensusSize)
	secondYearEncEpoch           = 2 * 365 * 24 * 3600 / (epochSize * consensusSize)
	// initial issuance:
	// 2% used for reward
	// 0.5% used for developer foundation
	// 4.5% used for allowance
	// almost 2.5 % used for staking
	//GenesisIssue, _ = new(big.Int).SetString("1000000000‬000000000000000000", 10)

	/**
	governance
	*/
	supportRate_Threshold = float64(0.85)
	maxVotingDuration     = uint64(14*24*60*60) / consensusSize * consensusSize
)

/**
PlatON Alpha Test Net
*/
var (
	/**
	Staking config
	**/
	testnet_stakeThreshold, _            = new(big.Int).SetString("1000000000000000000000000", 10)
	testnet_delegateThreshold, _         = new(big.Int).SetString("10", 10)
	testnet_consValidatorNum             = uint64(25)
	testnet_epochValidatorNum            = uint64(101)
	testnet_shiftValidatorNum            = uint64(8)
	testnet_epochSize                    = uint64(88)
	testnet_hesitateRatio                = uint64(1)
	testnet_effectiveRatio               = uint64(1)
	testnet_electionDistance             = uint64(20)
	testnet_consensusSize                = uint64(250)
	testnet_unStakeFreezeRatio           = uint64(1)
	testnet_passiveUnDelegateFreezeRatio = uint64(0)
	testnet_activeUnDelegateFreezeRatio  = uint64(0)

	/**
	Slashing
	*/
	testnet_blockAmountLow            = uint32(8)
	testnet_blockAmountHigh           = uint32(5)
	testnet_blockAmountLowSlashing    = uint32(10)
	testnet_blockAmountHighSlashing   = uint32(20)
	testnet_duplicateSignNum          = uint32(2)
	testnet_duplicateSignLowSlashing  = uint32(10)
	testnet_duplicateSignHighSlashing = uint32(10)

	/**
	Restricting
	TODO
	*/

	/**
	Reward
	*/
	testnet_secondYearAllowance, _       = new(big.Int).SetString("15000000000000000000000000", 10)
	testnet_thirdYearAllowance, _        = new(big.Int).SetString("5000000000000000000000000", 10)
	testnet_genesisRestrictingBalance, _ = new(big.Int).SetString("20000000000000000000000000", 10)
	testnet_firstYearEndEpoch            = 365 * 24 * 3600 / (epochSize * consensusSize)
	testnet_secondYearEncEpoch           = 2 * 365 * 24 * 3600 / (epochSize * consensusSize)

	/**
	governance
	*/
	testnet_supportRate_Threshold = float64(0.85)
	testnet_maxVotingDuration     = uint64(14*24*60*60) / consensusSize * consensusSize
)

/**
PlatON Beta Test Net
*/
var (
	/**
	Staking config
	**/
	beta_stakeThreshold, _            = new(big.Int).SetString("1000000000000000000000000", 10)
	beta_delegateThreshold, _         = new(big.Int).SetString("10", 10)
	beta_consValidatorNum             = uint64(25)
	beta_epochValidatorNum            = uint64(101)
	beta_shiftValidatorNum            = uint64(8)
	beta_epochSize                    = uint64(88)
	beta_hesitateRatio                = uint64(1)
	beta_effectiveRatio               = uint64(1)
	beta_electionDistance             = uint64(20)
	beta_consensusSize                = uint64(250)
	beta_unStakeFreezeRatio           = uint64(1)
	beta_passiveUnDelegateFreezeRatio = uint64(0)
	beta_activeUnDelegateFreezeRatio  = uint64(0)

	/**
	Slashing
	*/
	beta_blockAmountLow            = uint32(8)
	beta_blockAmountHigh           = uint32(5)
	beta_blockAmountLowSlashing    = uint32(10)
	beta_blockAmountHighSlashing   = uint32(20)
	beta_duplicateSignNum          = uint32(2)
	beta_duplicateSignLowSlashing  = uint32(10)
	beta_duplicateSignHighSlashing = uint32(10)

	/**
	Restricting
	TODO
	*/

	/**
	Reward
	*/
	beta_secondYearAllowance, _       = new(big.Int).SetString("15000000000000000000000000", 10)
	beta_thirdYearAllowance, _        = new(big.Int).SetString("5000000000000000000000000", 10)
	beta_genesisRestrictingBalance, _ = new(big.Int).SetString("20000000000000000000000000", 10)
	beta_firstYearEndEpoch            = 365 * 24 * 3600 / (epochSize * consensusSize)
	beta_secondYearEncEpoch           = 2 * 365 * 24 * 3600 / (epochSize * consensusSize)

	/**
	governance
	*/
	beta_supportRate_Threshold = float64(0.85)
	beta_maxVotingDuration     = uint64(14*24*60*60) / consensusSize * consensusSize
)

/**
PlatON Inner Test Net
*/
var (
	/**
	Staking config
	**/
	innerTest_stakeThreshold, _            = new(big.Int).SetString("1000000000000000000000000", 10)
	innerTest_delegateThreshold, _         = new(big.Int).SetString("10", 10)
	innerTest_consValidatorNum             = uint64(25)
	innerTest_epochValidatorNum            = uint64(101)
	innerTest_shiftValidatorNum            = uint64(8)
	innerTest_epochSize                    = uint64(88)
	innerTest_hesitateRatio                = uint64(1)
	innerTest_effectiveRatio               = uint64(1)
	innerTest_electionDistance             = uint64(20)
	innerTest_consensusSize                = uint64(250)
	innerTest_unStakeFreezeRatio           = uint64(1)
	innerTest_passiveUnDelegateFreezeRatio = uint64(0)
	innerTest_activeUnDelegateFreezeRatio  = uint64(0)

	/**
	Slashing
	*/
	innerTest_blockAmountLow            = uint32(8)
	innerTest_blockAmountHigh           = uint32(5)
	innerTest_blockAmountLowSlashing    = uint32(10)
	innerTest_blockAmountHighSlashing   = uint32(20)
	innerTest_duplicateSignNum          = uint32(2)
	innerTest_duplicateSignLowSlashing  = uint32(10)
	innerTest_duplicateSignHighSlashing = uint32(10)

	/**
	Restricting
	TODO
	*/

	/**
	Reward
	*/
	innerTest_secondYearAllowance, _       = new(big.Int).SetString("15000000000000000000000000", 10)
	innerTest_thirdYearAllowance, _        = new(big.Int).SetString("5000000000000000000000000", 10)
	innerTest_genesisRestrictingBalance, _ = new(big.Int).SetString("20000000000000000000000000", 10)
	innerTest_firstYearEndEpoch            = 365 * 24 * 3600 / (epochSize * consensusSize)
	innerTest_secondYearEncEpoch           = 2 * 365 * 24 * 3600 / (epochSize * consensusSize)

	/**
	governance
	*/
	innerTest_supportRate_Threshold = float64(0.85)
	innerTest_maxVotingDuration     = uint64(14*24*60*60) / consensusSize * consensusSize
)

/**
PlatON Inner Dev Net
*/
var (
	/**
	Staking config
	**/
	innerDev_stakeThreshold, _            = new(big.Int).SetString("1000000000000000000000000", 10)
	innerDev_delegateThreshold, _         = new(big.Int).SetString("10", 10)
	innerDev_consValidatorNum             = uint64(25)
	innerDev_epochValidatorNum            = uint64(101)
	innerDev_shiftValidatorNum            = uint64(8)
	innerDev_epochSize                    = uint64(88)
	innerDev_hesitateRatio                = uint64(1)
	innerDev_effectiveRatio               = uint64(1)
	innerDev_electionDistance             = uint64(20)
	innerDev_consensusSize                = uint64(250)
	innerDev_unStakeFreezeRatio           = uint64(1)
	innerDev_passiveUnDelegateFreezeRatio = uint64(0)
	innerDev_activeUnDelegateFreezeRatio  = uint64(0)

	/**
	Slashing
	*/
	innerDev_blockAmountLow            = uint32(8)
	innerDev_blockAmountHigh           = uint32(5)
	innerDev_blockAmountLowSlashing    = uint32(10)
	innerDev_blockAmountHighSlashing   = uint32(20)
	innerDev_duplicateSignNum          = uint32(2)
	innerDev_duplicateSignLowSlashing  = uint32(10)
	innerDev_duplicateSignHighSlashing = uint32(10)

	/**
	Restricting
	TODO
	*/

	/**
	Reward
	*/
	innerDev_secondYearAllowance, _       = new(big.Int).SetString("15000000000000000000000000", 10)
	innerDev_thirdYearAllowance, _        = new(big.Int).SetString("5000000000000000000000000", 10)
	innerDev_genesisRestrictingBalance, _ = new(big.Int).SetString("20000000000000000000000000", 10)
	innerDev_firstYearEndEpoch            = 365 * 24 * 3600 / (epochSize * consensusSize)
	innerDev_secondYearEncEpoch           = 2 * 365 * 24 * 3600 / (epochSize * consensusSize)

	/**
	governance
	*/
	innerDev_supportRate_Threshold = float64(0.85)
	innerDev_maxVotingDuration     = uint64(14*24*60*60) / consensusSize * consensusSize
)

/**
PlatON Main Net
*/
var defaultStakingConfig = StakingConfig{
	StakeThreshold:               stakeThreshold,
	DelegateThreshold:            delegateThreshold,
	ConsValidatorNum:             consValidatorNum,
	EpochValidatorNum:            epochValidatorNum,
	ShiftValidatorNum:            shiftValidatorNum,
	EpochSize:                    epochSize,
	HesitateRatio:                hesitateRatio,
	EffectiveRatio:               effectiveRatio,
	ElectionDistance:             electionDistance,
	ConsensusSize:                consensusSize,
	UnStakeFreezeRatio:           unStakeFreezeRatio,
	PassiveUnDelegateFreezeRatio: passiveUnDelegateFreezeRatio,
	ActiveUnDelegateFreezeRatio:  activeUnDelegateFreezeRatio,
}
var defaultSlashingConfig = SlashingConfig{
	BlockAmountLow:            blockAmountLow,
	BlockAmountHigh:           blockAmountHigh,
	BlockAmountLowSlashing:    blockAmountLowSlashing,
	BlockAmountHighSlashing:   blockAmountHighSlashing,
	DuplicateSignNum:          duplicateSignNum,
	DuplicateSignLowSlashing:  duplicateSignLowSlashing,
	DuplicateSignHighSlashing: duplicateSignHighSlashing,
}
var defaultRewardConfig = RewardConfig{
	SecondYearAllowance:       secondYearAllowance,
	ThirdYearAllowance:        thirdYearAllowance,
	GenesisRestrictingBalance: genesisRestrictingBalance,
	FirstYearEndEpoch:         firstYearEndEpoch,
	SecondYearEncEpoch:        secondYearEncEpoch,
}
var defaultGovConfig = GovernanceConfig{
	SupportRate_Threshold: supportRate_Threshold,
	MaxVotingDuration:     maxVotingDuration,
}

/**
Alpha Test Net
*/
var testnetDefaultStakingConfig = StakingConfig{
	StakeThreshold:               testnet_delegateThreshold,
	DelegateThreshold:            testnet_delegateThreshold,
	ConsValidatorNum:             testnet_consValidatorNum,
	EpochValidatorNum:            testnet_epochValidatorNum,
	ShiftValidatorNum:            testnet_shiftValidatorNum,
	EpochSize:                    testnet_epochSize,
	HesitateRatio:                testnet_hesitateRatio,
	EffectiveRatio:               testnet_effectiveRatio,
	ElectionDistance:             testnet_electionDistance,
	ConsensusSize:                testnet_consensusSize,
	UnStakeFreezeRatio:           testnet_unStakeFreezeRatio,
	PassiveUnDelegateFreezeRatio: testnet_passiveUnDelegateFreezeRatio,
	ActiveUnDelegateFreezeRatio:  testnet_activeUnDelegateFreezeRatio,
}
var testnetDefaultSlashingConfig = SlashingConfig{
	BlockAmountLow:            testnet_blockAmountLow,
	BlockAmountHigh:           testnet_blockAmountHigh,
	BlockAmountLowSlashing:    testnet_blockAmountLowSlashing,
	BlockAmountHighSlashing:   testnet_blockAmountHighSlashing,
	DuplicateSignNum:          testnet_duplicateSignNum,
	DuplicateSignLowSlashing:  testnet_duplicateSignLowSlashing,
	DuplicateSignHighSlashing: testnet_duplicateSignHighSlashing,
}
var testnetDefaultRewardConfig = RewardConfig{
	SecondYearAllowance:       testnet_secondYearAllowance,
	ThirdYearAllowance:        testnet_thirdYearAllowance,
	GenesisRestrictingBalance: testnet_genesisRestrictingBalance,
	FirstYearEndEpoch:         testnet_firstYearEndEpoch,
	SecondYearEncEpoch:        testnet_secondYearEncEpoch,
}
var testnetDefaultGovConfig = GovernanceConfig{
	SupportRate_Threshold: testnet_supportRate_Threshold,
	MaxVotingDuration:     testnet_maxVotingDuration,
}

/**
Beta Test Net
*/
var betaDefaultStakingConfig = StakingConfig{
	StakeThreshold:               beta_stakeThreshold,
	DelegateThreshold:            beta_delegateThreshold,
	ConsValidatorNum:             beta_consValidatorNum,
	EpochValidatorNum:            beta_epochValidatorNum,
	ShiftValidatorNum:            beta_shiftValidatorNum,
	EpochSize:                    beta_epochSize,
	HesitateRatio:                beta_hesitateRatio,
	EffectiveRatio:               beta_effectiveRatio,
	ElectionDistance:             beta_electionDistance,
	ConsensusSize:                beta_consensusSize,
	UnStakeFreezeRatio:           beta_unStakeFreezeRatio,
	PassiveUnDelegateFreezeRatio: beta_passiveUnDelegateFreezeRatio,
	ActiveUnDelegateFreezeRatio:  beta_activeUnDelegateFreezeRatio,
}
var betaDefaultSlashingConfig = SlashingConfig{
	BlockAmountLow:            beta_blockAmountLow,
	BlockAmountHigh:           beta_blockAmountHigh,
	BlockAmountLowSlashing:    beta_blockAmountLowSlashing,
	BlockAmountHighSlashing:   beta_blockAmountHighSlashing,
	DuplicateSignNum:          beta_duplicateSignNum,
	DuplicateSignLowSlashing:  beta_duplicateSignLowSlashing,
	DuplicateSignHighSlashing: beta_duplicateSignHighSlashing,
}
var betaDefaultRewardConfig = RewardConfig{
	SecondYearAllowance:       beta_secondYearAllowance,
	ThirdYearAllowance:        beta_thirdYearAllowance,
	GenesisRestrictingBalance: beta_genesisRestrictingBalance,
	FirstYearEndEpoch:         beta_firstYearEndEpoch,
	SecondYearEncEpoch:        beta_secondYearEncEpoch,
}
var betaDefaultGovConfig = GovernanceConfig{
	SupportRate_Threshold: beta_supportRate_Threshold,
	MaxVotingDuration:     beta_maxVotingDuration,
}

/**
PlatON Inner Test Net
*/
var innerTestDefaultStakingConfig = StakingConfig{
	StakeThreshold:               innerTest_stakeThreshold,
	DelegateThreshold:            innerTest_delegateThreshold,
	ConsValidatorNum:             innerTest_consValidatorNum,
	EpochValidatorNum:            innerTest_epochValidatorNum,
	ShiftValidatorNum:            innerTest_shiftValidatorNum,
	EpochSize:                    innerTest_epochSize,
	HesitateRatio:                innerTest_hesitateRatio,
	EffectiveRatio:               innerTest_effectiveRatio,
	ElectionDistance:             innerTest_electionDistance,
	ConsensusSize:                innerTest_consensusSize,
	UnStakeFreezeRatio:           innerTest_unStakeFreezeRatio,
	PassiveUnDelegateFreezeRatio: innerTest_passiveUnDelegateFreezeRatio,
	ActiveUnDelegateFreezeRatio:  innerTest_activeUnDelegateFreezeRatio,
}
var innerTestDefaultSlashingConfig = SlashingConfig{
	BlockAmountLow:            innerTest_blockAmountLow,
	BlockAmountHigh:           innerTest_blockAmountHigh,
	BlockAmountLowSlashing:    innerTest_blockAmountLowSlashing,
	BlockAmountHighSlashing:   innerTest_blockAmountHighSlashing,
	DuplicateSignNum:          innerTest_duplicateSignNum,
	DuplicateSignLowSlashing:  innerTest_duplicateSignLowSlashing,
	DuplicateSignHighSlashing: innerTest_duplicateSignHighSlashing,
}
var innerTestDefaultRewardConfig = RewardConfig{
	SecondYearAllowance:       innerTest_secondYearAllowance,
	ThirdYearAllowance:        innerTest_thirdYearAllowance,
	GenesisRestrictingBalance: innerTest_genesisRestrictingBalance,
	FirstYearEndEpoch:         innerTest_firstYearEndEpoch,
	SecondYearEncEpoch:        innerTest_secondYearEncEpoch,
}
var innerTestDefaultGovConfig = GovernanceConfig{
	SupportRate_Threshold: innerTest_supportRate_Threshold,
	MaxVotingDuration:     innerTest_maxVotingDuration,
}

/**
PlatON Inner Dev Net
*/
var innerDevDefaultStakingConfig = StakingConfig{
	StakeThreshold:               innerDev_stakeThreshold,
	DelegateThreshold:            innerDev_delegateThreshold,
	ConsValidatorNum:             innerDev_consValidatorNum,
	EpochValidatorNum:            innerDev_epochValidatorNum,
	ShiftValidatorNum:            innerDev_shiftValidatorNum,
	EpochSize:                    innerDev_epochSize,
	HesitateRatio:                innerDev_hesitateRatio,
	EffectiveRatio:               innerDev_effectiveRatio,
	ElectionDistance:             innerDev_electionDistance,
	ConsensusSize:                innerDev_consensusSize,
	UnStakeFreezeRatio:           innerDev_unStakeFreezeRatio,
	PassiveUnDelegateFreezeRatio: innerDev_passiveUnDelegateFreezeRatio,
	ActiveUnDelegateFreezeRatio:  innerDev_activeUnDelegateFreezeRatio,
}
var innerDevDefaultSlashingConfig = SlashingConfig{
	BlockAmountLow:            innerDev_blockAmountLow,
	BlockAmountHigh:           innerDev_blockAmountHigh,
	BlockAmountLowSlashing:    innerDev_blockAmountLowSlashing,
	BlockAmountHighSlashing:   innerDev_blockAmountHighSlashing,
	DuplicateSignNum:          innerDev_duplicateSignNum,
	DuplicateSignLowSlashing:  innerDev_duplicateSignLowSlashing,
	DuplicateSignHighSlashing: innerDev_duplicateSignHighSlashing,
}
var innerDevDefaultRewardConfig = RewardConfig{
	SecondYearAllowance:       innerDev_secondYearAllowance,
	ThirdYearAllowance:        innerDev_thirdYearAllowance,
	GenesisRestrictingBalance: innerDev_genesisRestrictingBalance,
	FirstYearEndEpoch:         innerDev_firstYearEndEpoch,
	SecondYearEncEpoch:        innerDev_secondYearEncEpoch,
}
var innerDevDefaultGovConfig = GovernanceConfig{
	SupportRate_Threshold: innerDev_supportRate_Threshold,
	MaxVotingDuration:     innerDev_maxVotingDuration,
}
