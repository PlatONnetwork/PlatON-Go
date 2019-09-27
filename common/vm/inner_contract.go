package vm

import "github.com/PlatONnetwork/PlatON-Go/common"

// the inner contract addr  table
var (
	RestrictingContractAddr    = common.HexToAddress("0x1000000000000000000000000000000000000001") // The PlatON Precompiled contract addr for restricting
	StakingContractAddr        = common.HexToAddress("0x1000000000000000000000000000000000000002") // The PlatON Precompiled contract addr for staking
	RewardManagerPoolAddr      = common.HexToAddress("0x1000000000000000000000000000000000000003") // The PlatON Precompiled contract addr for reward
	SlashingContractAddr       = common.HexToAddress("0x1000000000000000000000000000000000000004") // The PlatON Precompiled contract addr for slashing
	GovContractAddr            = common.HexToAddress("0x1000000000000000000000000000000000000005") // The PlatON Precompiled contract addr for governance
	ValidatorInnerContractAddr = common.HexToAddress("0x2000000000000000000000000000000000000000") // The PlatON Precompiled contract addr for cbft inner
)
