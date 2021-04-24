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

package gov

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ActiveVersionError                = common.NewBizError(302001, "Current active version not found")
	VoteOptionError                   = common.NewBizError(302002, "Illegal voting option")
	ProposalTypeError                 = common.NewBizError(302003, "Illegal proposal type")
	ProposalIDEmpty                   = common.NewBizError(302004, "Proposal ID is null")
	ProposalIDExist                   = common.NewBizError(302005, "Proposal ID already existed")
	ProposalNotFound                  = common.NewBizError(302006, "Proposal not found")
	PIPIDEmpty                        = common.NewBizError(302007, "PIPID is null")
	PIPIDExist                        = common.NewBizError(302008, "PIPID already existed")
	EndVotingRoundsTooSmall           = common.NewBizError(302009, "EndVotingRounds is too small")
	EndVotingRoundsTooLarge           = common.NewBizError(302010, "EndVotingRounds is too large")
	NewVersionError                   = common.NewBizError(302011, "NewVersion should larger than current active version")
	VotingVersionProposalExist        = common.NewBizError(302012, "Another version proposal already existed at voting stage")
	PreActiveVersionProposalExist     = common.NewBizError(302013, "Another version proposal already existed at pre-active stage")
	VotingCancelProposalExist         = common.NewBizError(302014, "Another cancel proposal already existed at voting stage")
	TobeCanceledProposalNotFound      = common.NewBizError(302015, "The proposal to be canceled is not found")
	TobeCanceledProposalTypeError     = common.NewBizError(302016, "The proposal to be canceled proposal has an illegal version type")
	TobeCanceledProposalNotAtVoting   = common.NewBizError(302017, "The proposal to be canceled is not at voting stage")
	ProposerEmpty                     = common.NewBizError(302018, "The proposer is null")
	VerifierInfoNotFound              = common.NewBizError(302019, "Detailed verifier information is not found")
	VerifierStatusInvalid             = common.NewBizError(302020, "The verifier status is invalid")
	TxSenderDifferFromStaking         = common.NewBizError(302021, "Transaction account is inconsistent with the staking account")
	TxSenderIsNotVerifier             = common.NewBizError(302022, "Transaction node is not the validator")
	TxSenderIsNotCandidate            = common.NewBizError(302023, "Transaction node is not the candidate")
	VersionSignError                  = common.NewBizError(302024, "Invalid version signature")
	VerifierNotUpgraded               = common.NewBizError(302025, "Verifier does not upgraded to the latest version")
	ProposalNotAtVoting               = common.NewBizError(302026, "The proposal is not at voting stage")
	VoteDuplicated                    = common.NewBizError(302027, "Duplicated votes found")
	DeclareVersionError               = common.NewBizError(302028, "Declared version is invalid")
	NotifyStakingDeclaredVersionError = common.NewBizError(302029, "Error is found when notifying staking for the declared version")
	TallyResultNotFound               = common.NewBizError(302030, "The result of proposal is not found")
	UnsupportedGovernParam            = common.NewBizError(302031, "Unsupported governent parameter")
	VotingParamProposalExist          = common.NewBizError(302032, "Another parameter proposal already existed at voting stage")
	GovernParamValueError             = common.NewBizError(302033, "Govern parameter value error")
	ParamProposalIsSameValue          = common.NewBizError(302034, "The new value of the parameter proposal is the same as the old one")
)
