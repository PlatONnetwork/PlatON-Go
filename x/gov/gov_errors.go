package gov

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ActiveVersionError = common.NewSysError("current active version not found")

	VoteOptionError = common.NewBizError("vote option error")

	ProposalTypeError = common.NewBizError("proposal type error")

	ProposalIDEmpty  = common.NewBizError("proposal ID is empty")
	ProposalIDExist  = common.NewBizError("proposal ID already exists")
	ProposalNotFound = common.NewBizError("proposal not found")

	PIPIDEmpty = common.NewBizError("PIPID is empty")
	PIPIDExist = common.NewBizError("PIPID already exists")

	EndVotingRoundsTooSmall = common.NewBizError("endVotingRounds too small")
	EndVotingRoundsTooLarge = common.NewBizError("endVotingRounds too large")

	NewVersionError               = common.NewBizError("newVersion should larger than current active version")
	VotingVersionProposalExist    = common.NewBizError("another version proposal at voting stage")
	PreActiveVersionProposalExist = common.NewBizError("another version proposal at pre-active stage")

	VotingCancelProposalExist       = common.NewBizError("another cancel proposal at voting stage")
	TobeCanceledProposalNotFound    = common.NewBizError("to be canceled proposal not found")
	TobeCanceledProposalTypeError   = common.NewBizError("to be canceled proposal not version type")
	TobeCanceledProposalNotAtVoting = common.NewBizError("to be canceled proposal not at voting stage")

	ProposerEmpty = common.NewBizError("proposer is empty")

	VerifierInfoNotFound  = common.NewBizError("verifier detail info not found")
	VerifierStatusInvalid = common.NewBizError("verifier status is invalid")

	TxSenderDifferFromStaking = common.NewBizError("tx sender differ from staking")
	TxSenderIsNotVerifier     = common.NewBizError("tx sender is not verifier")
	TxSenderIsNotCandidate    = common.NewBizError("tx sender is not candidate")

	VersionSignError    = common.NewBizError("version sign error")
	VerifierNotUpgraded = common.NewBizError("verifier not upgraded")

	ProposalNotAtVoting = common.NewBizError("proposal not at voting stage")
	VoteDuplicated      = common.NewBizError("vote duplicated")

	DeclareVersionError               = common.NewBizError("declared version error")
	NotifyStakingDeclaredVersionError = common.NewBizError("notify staking declared version error")
)
