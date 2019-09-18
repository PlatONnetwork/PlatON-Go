package vm

import "github.com/PlatONnetwork/PlatON-Go/common"

// the inner contract addr  table
var (
	RestrictingContractAddr      = common.HexToAddress("0x1000000000000000000000000000000000000001") // The PlatON Precompiled contract addr for restricting
	StakingContractAddr          = common.HexToAddress("0x1000000000000000000000000000000000000002") // The PlatON Precompiled contract addr for staking
	RewardManagerPoolAddr        = common.HexToAddress("0x1000000000000000000000000000000000000003") // The PlatON Precompiled contract addr for reward
	SlashingContractAddr         = common.HexToAddress("0x1000000000000000000000000000000000000004") // The PlatON Precompiled contract addr for slashing
	GovContractAddr              = common.HexToAddress("0x1000000000000000000000000000000000000005") // The PlatON Precompiled contract addr for governance
	ValidatorInnerContractAddr   = common.HexToAddress("0x2000000000000000000000000000000000000000") // The PlatON Precompiled contract addr for cbft inner
	CommunityDeveloperFoundation = common.HexToAddress("0x60ceca9c1290ee56b98d4e160ef0453f7c40d219") // Community development Foundation addr
	PlatONFoundationAddress      = common.HexToAddress("0x55bfd49472fd41211545b01713a9c3a97af78b05") // PlatON Foundation addr
)
