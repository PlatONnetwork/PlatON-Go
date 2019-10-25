package restricting

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ErrParamEpochInvalid                 = common.NewBizError(304001, "param epoch can't be zero")
	ErrCountRestrictPlansInvalid         = common.NewBizError(304002, "the number of the restricting plan can't be zero or more than 36")
	ErrLockedAmountTooLess               = common.NewBizError(304003, "total restricting amount need more than 1 LAT")
	ErrBalanceNotEnough                  = common.NewBizError(304004, "create plan,the sender balance is not enough in restrict")
	ErrAccountNotFound                   = common.NewBizError(304005, "account is not found on restricting contract")
	ErrSlashingTooMuch                   = common.NewBizError(304006, "slashing amount is larger than staking amount")
	ErrStakingAmountEmpty                = common.NewBizError(304007, "staking amount is 0")
	ErrPledgeLockFundsAmountLessThanZero = common.NewBizError(304008, "pledge lock funds amount can't less than 0")
	ErrReturnLockFundsAmountLessThanZero = common.NewBizError(304009, "return lock funds amount can't less than 0")
	ErrSlashingAmountLessThanZero        = common.NewBizError(304010, "slashing amount can't less than 0")
	ErrCreatePlanAmountLessThanZero      = common.NewBizError(304011, "create plan each amount can't less than 0")
	ErrStakingAmountInvalid              = common.NewBizError(304012, "staking return amount is wrong")
	ErrRestrictBalanceNotEnough          = common.NewBizError(304013, "the user restricting balance is not enough for pledge lock funds")
)
