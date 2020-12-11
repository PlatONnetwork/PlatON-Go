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
	ErrParamEpochInvalid                  = common.NewBizError(304001, "The initial epoch for staking cannot be zero")
	ErrCountRestrictPlansInvalid          = common.NewBizError(304002, fmt.Sprintf("The number of the restricting plan cannot be (0, %d]", RestrictTxPlanSize))
	ErrLockedAmountTooLess                = common.NewBizError(304003, "Total staking amount shall be more than 1 LAT")
	ErrBalanceNotEnough                   = common.NewBizError(304004, "Create plan,the sender balance is not enough in restrict")
	ErrAccountNotFound                    = common.NewBizError(304005, "Account is not found on restricting contract")
	ErrSlashingTooMuch                    = common.NewBizError(304006, "Slashing amount is larger than staking amount")
	ErrStakingAmountEmpty                 = common.NewBizError(304007, "Staking amount cannot be 0")
	ErrPledgeLockFundsAmountLessThanZero  = common.NewBizError(304008, "Pledge lock funds amount cannot be less than or equal to 0")
	ErrReturnLockFundsAmountLessThanZero  = common.NewBizError(304009, "Return lock funds amount cannot be less than or equal to 0")
	ErrSlashingAmountLessThanZero         = common.NewBizError(304010, "Slashing amount cannot be less than 0")
	ErrCreatePlanAmountLessThanZero       = common.NewBizError(304011, "Create plan each amount cannot be less than or equal to 0")
	ErrStakingAmountInvalid               = common.NewBizError(304012, "The staking amount is less than the return amount")
	ErrRestrictBalanceNotEnough           = common.NewBizError(304013, "The user restricting balance is not enough for pledge lock funds")
	ErrCreatePlanAmountLessThanMiniAmount = common.NewBizError(304014, "Create plan each amount should greater than mini amount")
	ErrRestrictBalanceAndFreeNotEnough    = common.NewBizError(304015, "The user restricting  and free balance is not enough for pledge lock funds")
)
