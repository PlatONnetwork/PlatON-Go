package vm

import "github.com/PlatONnetwork/PlatON-Go/common"

// the inner contract addr  table
var (
	LockRepoContractAddr = common.HexToAddress("0x1000000000000000000000000000000000000001")
	StakingContractAddr = common.HexToAddress("0x1000000000000000000000000000000000000002")
	AwardMgrContractAddr = common.HexToAddress("0x1000000000000000000000000000000000000003")
	SlashingContractAddr = common.HexToAddress("0x1000000000000000000000000000000000000004")
	ValidatorInnerContractAddr = common.HexToAddress("0x2000000000000000000000000000000000000000")
)
