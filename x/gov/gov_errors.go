package gov

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ActiveVersionError = common.NewBizError(302001, "current active version not found")

	VoteOptionError = common.NewBizError(302002, "vote option error")

	ProposalTypeError = common.NewBizError(302003, "proposal type error")

	ProposalIDEmpty  = common.NewBizError(302004, "proposal ID is empty")
	ProposalIDExist  = common.NewBizError(302005, "proposal ID already exists")
	ProposalNotFound = common.NewBizError(302006, "proposal not found")

	PIPIDEmpty = common.NewBizError(302007, "PIPID is empty")
	PIPIDExist = common.NewBizError(302008, "PIPID already exists")

	EndVotingRoundsTooSmall = common.NewBizError(302009, "endVotingRounds too small")
	EndVotingRoundsTooLarge = common.NewBizError(302010, "endVotingRounds too large")

	NewVersionError               = common.NewBizError(302011, "newVersion should larger than current active version")
	VotingVersionProposalExist    = common.NewBizError(302012, "another version proposal at voting stage")
	PreActiveVersionProposalExist = common.NewBizError(302013, "another version proposal at pre-active stage")

	VotingCancelProposalExist       = common.NewBizError(302014, "another cancel proposal at voting stage")
	TobeCanceledProposalNotFound    = common.NewBizError(302015, "to be canceled proposal not found")
	TobeCanceledProposalTypeError   = common.NewBizError(302016, "to be canceled proposal not version type")
	TobeCanceledProposalNotAtVoting = common.NewBizError(302017, "to be canceled proposal not at voting stage")

	ProposerEmpty = common.NewBizError(302018, "proposer is empty")

	VerifierInfoNotFound  = common.NewBizError(302019, "verifier detail info not found")
	VerifierStatusInvalid = common.NewBizError(302020, "verifier status is invalid")

	TxSenderDifferFromStaking = common.NewBizError(302021, "Tx caller differ from staking")
	TxSenderIsNotVerifier     = common.NewBizError(302022, "Tx caller is not verifier")
	TxSenderIsNotCandidate    = common.NewBizError(302023, "Tx caller is not candidate")

	VersionSignError    = common.NewBizError(302024, "version sign error")
	VerifierNotUpgraded = common.NewBizError(302025, "verifier not upgraded")

	ProposalNotAtVoting = common.NewBizError(302026, "proposal not at voting stage")
	VoteDuplicated      = common.NewBizError(302027, "vote duplicated")

	DeclareVersionError               = common.NewBizError(302028, "declared version error")
	NotifyStakingDeclaredVersionError = common.NewBizError(302029, "notify staking declared version error")

	TallyResultNotFound = common.NewBizError(302030, "tally result not found")
)
