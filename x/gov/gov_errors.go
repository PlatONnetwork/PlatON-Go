package gov

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ActiveVersionError = common.NewBizError(200, "current active version not found")

	VoteOptionError = common.NewBizError(201, "vote option error")

	ProposalTypeError = common.NewBizError(202, "proposal type error")

	ProposalIDEmpty  = common.NewBizError(203, "proposal ID is empty")
	ProposalIDExist  = common.NewBizError(204, "proposal ID already exists")
	ProposalNotFound = common.NewBizError(205, "proposal not found")

	PIPIDEmpty = common.NewBizError(206, "PIPID is empty")
	PIPIDExist = common.NewBizError(207, "PIPID already exists")

	EndVotingRoundsTooSmall = common.NewBizError(208, "endVotingRounds too small")
	EndVotingRoundsTooLarge = common.NewBizError(209, "endVotingRounds too large")

	NewVersionError               = common.NewBizError(210, "newVersion should larger than current active version")
	VotingVersionProposalExist    = common.NewBizError(211, "another version proposal at voting stage")
	PreActiveVersionProposalExist = common.NewBizError(212, "another version proposal at pre-active stage")

	VotingCancelProposalExist       = common.NewBizError(213, "another cancel proposal at voting stage")
	TobeCanceledProposalNotFound    = common.NewBizError(214, "to be canceled proposal not found")
	TobeCanceledProposalTypeError   = common.NewBizError(215, "to be canceled proposal not version type")
	TobeCanceledProposalNotAtVoting = common.NewBizError(206, "to be canceled proposal not at voting stage")

	ProposerEmpty = common.NewBizError(217, "proposer is empty")

	VerifierInfoNotFound  = common.NewBizError(218, "verifier detail info not found")
	VerifierStatusInvalid = common.NewBizError(219, "verifier status is invalid")

	TxSenderDifferFromStaking = common.NewBizError(220, "tx sender differ from staking")
	TxSenderIsNotVerifier     = common.NewBizError(221, "tx sender is not verifier")
	TxSenderIsNotCandidate    = common.NewBizError(222, "tx sender is not candidate")

	VersionSignError    = common.NewBizError(223, "version sign error")
	VerifierNotUpgraded = common.NewBizError(224, "verifier not upgraded")

	ProposalNotAtVoting = common.NewBizError(225, "proposal not at voting stage")
	VoteDuplicated      = common.NewBizError(226, "vote duplicated")

	DeclareVersionError               = common.NewBizError(227, "declared version error")
	NotifyStakingDeclaredVersionError = common.NewBizError(228, "notify staking declared version error")

	TallyResultNotFound = common.NewBizError(229, "tally result not found")
)
