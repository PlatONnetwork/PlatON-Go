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
	ErrWrongBlsPubKey            = common.NewBizError(301000, "Wrong bls public key")
	ErrWrongBlsPubKeyProof       = common.NewBizError(301001, "Wrong bls public key proof")
	ErrDescriptionLen            = common.NewBizError(301002, "The Description length is wrong")
	ErrWrongProgramVersionSign   = common.NewBizError(301003, "The program version sign is wrong")
	ErrProgramVersionTooLow      = common.NewBizError(301004, "The program version of the relates node's is too low")
	ErrDeclVsFialedCreateCan     = common.NewBizError(301005, "DeclareVersion is failed on create staking")
	ErrNoSameStakingAddr         = common.NewBizError(301006, "The address must be the same as initiated staking")
	ErrInvalidRewardPer          = common.NewBizError(301007, "Invalid param RewardPer")
	ErrRewardPerInterval         = common.NewBizError(301008, "Modify the commission reward ratio too frequently")
	ErrRewardPerChangeRange      = common.NewBizError(301009, "The modification range exceeds the limit")
	ErrStakeVonTooLow            = common.NewBizError(301100, "Staking deposit too low")
	ErrCanAlreadyExist           = common.NewBizError(301101, "This candidate is already exist")
	ErrCanNoExist                = common.NewBizError(301102, "This candidate is not exist")
	ErrCanStatusInvalid          = common.NewBizError(301103, "This candidate status was invalided")
	ErrIncreaseStakeVonTooLow    = common.NewBizError(301104, "IncreaseStake von is too low")
	ErrDelegateVonTooLow         = common.NewBizError(301105, "Delegate deposit too low")
	ErrAccountNoAllowToDelegate  = common.NewBizError(301106, "The account is not allowed to be used for delegating")
	ErrCanNoAllowDelegate        = common.NewBizError(301107, "The candidate does not accept the delegation")
	ErrWithdrewDelegateVonTooLow = common.NewBizError(301108, "Withdrew delegation von is too low")
	ErrDelegateNoExist           = common.NewBizError(301109, "This delegation is not exist")
	ErrWrongVonOptType           = common.NewBizError(301110, "The von operation type is wrong")
	ErrAccountVonNoEnough        = common.NewBizError(301111, "The von of account is not enough")
	ErrBlockNumberDisordered     = common.NewBizError(301112, "The blockNumber is disordered")
	ErrDelegateVonNoEnough       = common.NewBizError(301113, "The von of delegation is not enough")
	ErrWrongWithdrewDelVonCalc   = common.NewBizError(301114, "Withdrew delegation von calculation is wrong")
	ErrValidatorNoExist          = common.NewBizError(301115, "The validator is not exist")
	ErrWrongFuncParams           = common.NewBizError(301116, "The fn params is wrong")
	ErrWrongSlashType            = common.NewBizError(301117, "The slashing type is wrong")
	ErrSlashVonOverflow          = common.NewBizError(301118, "Slashing amount is overflow")
	ErrWrongSlashVonCalc         = common.NewBizError(301119, "Slashing candidate von calculate is wrong")
	ErrGetVerifierList           = common.NewBizError(301200, "Getting verifierList is failed")
	ErrGetValidatorList          = common.NewBizError(301201, "Getting validatorList is failed")
	ErrGetCandidateList          = common.NewBizError(301202, "Getting candidateList is failed")
	ErrGetDelegateRelated        = common.NewBizError(301203, "Getting related of delegate is failed")
	ErrQueryCandidateInfo        = common.NewBizError(301204, "Query candidate info failed")
	ErrQueryDelegateInfo         = common.NewBizError(301205, "Query delegate info failed")
)
