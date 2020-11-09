// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package restricting

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

const (
	RestrictTxPlanSize = 36
)

var (
	ErrParamEpochInvalid                  = common.NewBizError(304001, "param epoch can't be zero")
	ErrCountRestrictPlansInvalid          = common.NewBizError(304002, fmt.Sprintf("the number of the restricting plan can't be zero or more than %d", RestrictTxPlanSize))
	ErrLockedAmountTooLess                = common.NewBizError(304003, "total restricting amount need more than 1 LAT")
	ErrBalanceNotEnough                   = common.NewBizError(304004, "create plan,the sender balance is not enough in restrict")
	ErrAccountNotFound                    = common.NewBizError(304005, "account is not found on restricting contract")
	ErrSlashingTooMuch                    = common.NewBizError(304006, "slashing amount is larger than staking amount")
	ErrStakingAmountEmpty                 = common.NewBizError(304007, "staking amount is 0")
	ErrPledgeLockFundsAmountLessThanZero  = common.NewBizError(304008, "pledge lock funds amount should greater than 0")
	ErrReturnLockFundsAmountLessThanZero  = common.NewBizError(304009, "return lock funds amount should greater than 0")
	ErrSlashingAmountLessThanZero         = common.NewBizError(304010, "slashing amount can't less than 0")
	ErrCreatePlanAmountLessThanZero       = common.NewBizError(304011, "create plan each amount should greater than 0")
	ErrStakingAmountInvalid               = common.NewBizError(304012, "staking return amount is wrong")
	ErrRestrictBalanceNotEnough           = common.NewBizError(304013, "the user restricting balance is not enough for pledge lock funds")
	ErrCreatePlanAmountLessThanMiniAmount = common.NewBizError(304014, "create plan each amount should greater than mini amount")
)
