package gov

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ActiveVersionError = common.NewSysError("gov:parameter error:current active version not found")

	VoteOptionError = common.NewBizError("gov:parameter error:vote option error")

	ProposalTypeError = common.NewBizError("gov:parameter error:proposal type error")

	ProposalIDEmpty  = common.NewBizError("gov:parameter error:proposal ID is empty")
	ProposalIDExist  = common.NewBizError("gov:parameter error:proposal ID already exists")
	ProposalNotFound = common.NewBizError("gov:parameter error:proposal not found")

	PIPIDEmpty = common.NewBizError("gov:parameter error:PIPID is empty")
	PIPIDExist = common.NewBizError("gov:parameter error:PIPID already exists")

	EndVotingRoundsZero     = common.NewBizError("gov:voting consensus rounds should > 0")
	EndVotingRoundsTooSmall = common.NewBizError("gov:parameter error:EndVotingRounds too small")
	EndVotingRoundsTooLarge = common.NewBizError("gov:parameter error:EndVotingRounds too large")

	NewVersionError               = common.NewBizError("gov:parameter error:NewVersion should larger than current active version")
	VotingVersionProposalExist    = common.NewBizError("gov:parameter error:another version proposal at voting stage")
	PreActiveVersionProposalExist = common.NewBizError("gov:parameter error:another pre-active version proposal")

	VotingCancelProposalExist       = common.NewBizError("gov:parameter error:another cancel proposal at voting stage")
	TobeCanceledProposalNotFound    = common.NewBizError("gov:parameter error:to be canceled proposal not found")
	TobeCanceledProposalTypeError   = common.NewBizError("gov:parameter error:to be canceled proposal not version type")
	TobeCanceledProposalNotAtVoting = common.NewBizError("gov:parameter error:to be canceled proposal not at voting stage")

	ProposerEmpty = common.NewBizError("gov:parameter error:proposer is empty")

	VerifierInfoNotFound  = common.NewBizError("gov:parameter error:verifier detail info not found")
	VerifierStatusInvalid = common.NewBizError("gov:parameter error:verifier status is invalid")

	TxSenderDifferFromStaking = common.NewBizError("gov:parameter error:tx sender differ from staking")
	TxSenderIsNotVerifier     = common.NewBizError("gov:parameter error:tx sender is not verifier")
	TxSenderIsNotCandidate    = common.NewBizError("gov:parameter error:tx sender is not candidate")

	VersionSignError    = common.NewBizError("gov:parameter error:version sign error")
	VerifierNotUpgraded = common.NewBizError("gov:parameter error:verifier not upgraded")

	ProposalNotAtVoting = common.NewBizError("gov:parameter error:proposal not at voting stage")
	VoteDuplicated      = common.NewBizError("gov:parameter error:vote duplicated")

	DeclareVersionError               = common.NewBizError("gov:parameter error:declared version error")
	NotifyStakingDeclaredVersionError = common.NewBizError("gov:parameter error:notify staking declared version error")
)
