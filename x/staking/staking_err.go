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

package staking

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ErrWrongBlsPubKey            = common.NewBizError(301000, "Invalid BLS public key length")
	ErrWrongBlsPubKeyProof       = common.NewBizError(301001, "The BLS proof is incorrect")
	ErrDescriptionLen            = common.NewBizError(301002, "The Description length is incorrect")
	ErrWrongProgramVersionSign   = common.NewBizError(301003, "The program version signature is invalid")
	ErrProgramVersionTooLow      = common.NewBizError(301004, "The program version is too low")
	ErrDeclVsFialedCreateCan     = common.NewBizError(301005, "The Version Declaration is failed when creating staking")
	ErrNoSameStakingAddr         = common.NewBizError(301006, "The address must be the same as the one initiated staking")
	ErrInvalidRewardPer          = common.NewBizError(301007, "Invalid param RewardPer")
	ErrRewardPerInterval         = common.NewBizError(301008, "Modify the commission reward ratio too frequently")
	ErrRewardPerChangeRange      = common.NewBizError(301009, "The modification range exceeds the limit")
	ErrStakeVonTooLow            = common.NewBizError(301100, "Staking deposit is insufficient")
	ErrCanAlreadyExist           = common.NewBizError(301101, "The candidate already existed")
	ErrCanNoExist                = common.NewBizError(301102, "The candidate does not exist")
	ErrCanStatusInvalid          = common.NewBizError(301103, "This candidate status is expired")
	ErrIncreaseStakeVonTooLow    = common.NewBizError(301104, "Increased stake is insufficient")
	ErrDelegateVonTooLow         = common.NewBizError(301105, "Delegate deposit is insufficient")
	ErrAccountNoAllowToDelegate  = common.NewBizError(301106, "The account is not allowed to delegate")
	ErrCanNoAllowDelegate        = common.NewBizError(301107, "The candidate is not allowed to delegate")
	ErrWithdrewDelegateVonTooLow = common.NewBizError(301108, "Withdrawal of delegation is insufficient")
	ErrDelegateNoExist           = common.NewBizError(301109, "The delegate does not exist")
	ErrWrongVonOptType           = common.NewBizError(301110, "The von operation type is incorrect")
	ErrAccountVonNoEnough        = common.NewBizError(301111, "The account balance is insufficient")
	ErrBlockNumberDisordered     = common.NewBizError(301112, "The blockNumber is inconsistent with the expected number")
	ErrDelegateVonNoEnough       = common.NewBizError(301113, "The balance of delegate is insufficient")
	ErrWrongWithdrewDelVonCalc   = common.NewBizError(301114, "The amount of delegate withdrawal is incorrect")
	ErrValidatorNoExist          = common.NewBizError(301115, "The validator does not exist")
	ErrWrongFuncParams           = common.NewBizError(301116, "The fn params is invalid")
	ErrWrongSlashType            = common.NewBizError(301117, "The slash type is illegal")
	ErrSlashVonOverflow          = common.NewBizError(301118, "The amount of slash is overflowed")
	ErrWrongSlashVonCalc         = common.NewBizError(301119, "The amount of slash for decreasing staking is incorrect")
	ErrGetVerifierList           = common.NewBizError(301200, "Retreiving verifier list failed")
	ErrGetValidatorList          = common.NewBizError(301201, "Retreiving validator list failed")
	ErrGetCandidateList          = common.NewBizError(301202, "Retreiving candidate list failed")
	ErrGetDelegateRelated        = common.NewBizError(301203, "Retreiving delegation related mapping failed")
	ErrQueryCandidateInfo        = common.NewBizError(301204, "Query candidate info failed")
	ErrQueryDelegateInfo         = common.NewBizError(301205, "Query delegate info failed")
)
