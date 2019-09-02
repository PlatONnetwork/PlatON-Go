package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/handler"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type Staking interface {
	GetVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.ValidatorExQueue, error)
	GetCandidateList(blockHash common.Hash, blockNumber uint64) (staking.CandidateHexQueue, error)
	GetCandidateInfo(blockHash common.Hash, addr common.Address) (*staking.Candidate, error)
	DeclarePromoteNotify(blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID, programVersion uint32) error
}

func GetVersionForStaking(state xcom.StateDB) uint32 {
	preActiveVersion := GetPreActiveVersion(state)
	if preActiveVersion > 0 {
		return preActiveVersion
	} else {
		return GetCurrentActiveVersion(state)
	}
}

func GetActiveVersion(blockNumber uint64, state xcom.StateDB) uint32 {
	avList, err := ListActiveVersion(state)
	if err != nil {
		log.Error("List active version error", "err", err)
		return 0
	}

	for _, av := range avList {
		if blockNumber >= av.ActiveBlock {
			return av.ActiveVersion
		}
	}
	return 0
}

// Get current active version record
func GetCurrentActiveVersion(state xcom.StateDB) uint32 {
	avList, err := ListActiveVersion(state)
	if err != nil {
		log.Error("Cannot find active version list")
		return 0
	}

	var version uint32
	if len(avList) == 0 {
		log.Error("cannot find current active version")
		return 0
	} else {
		version = avList[0].ActiveVersion
	}
	return version
}

func GetProgramVersion() (*ProgramVersionValue, error) {
	programVersion := uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)
	sig, err := handler.GetCryptoHandler().Sign(programVersion)
	if err != nil {
		log.Error("sign version data error", "err", err)
		return nil, err
	}
	value := &ProgramVersionValue{ProgramVersion: programVersion, ProgramVersionSign: hexutil.Encode(sig)}
	return value, nil
}

// submit a proposal
func Submit(from common.Address, proposal Proposal, blockHash common.Hash, blockNumber uint64, stk Staking, state xcom.StateDB) error {
	log.Debug("call Submit", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "proposal", proposal)

	//param check
	if err := proposal.Verify(blockNumber, blockHash, state); err != nil {
		if bizError, ok := err.(*common.BizError); ok {
			return bizError
		} else {
			log.Error("verify proposal parameters failed", "err", err)
			return common.NewBizError(err.Error())
		}
	}

	//check caller and proposer
	if err := checkVerifier(from, proposal.GetProposer(), blockHash, proposal.GetSubmitBlock(), stk); err != nil {
		return err
	}

	//handle storage
	if err := SetProposal(proposal, state); err != nil {
		log.Error("save proposal failed", "proposalID", proposal.GetProposalID())
		return err
	}
	if err := AddVotingProposalID(blockHash, proposal.GetProposalID()); err != nil {
		log.Error("add proposal ID to voting proposal ID list failed", "proposalID", proposal.GetProposalID())
		return err
	}
	return nil
}

// vote for a proposal
func Vote(from common.Address, vote VoteInfo, blockHash common.Hash, blockNumber uint64, programVersion uint32, programVersionSign common.VersionSign, stk Staking, state xcom.StateDB) error {
	log.Debug("call Vote", "from", from, "proposalID", vote.ProposalID, "voteNodeID", vote.VoteNodeID, "voteOption", vote.VoteOption, "blockHash", blockHash, "blockNumber", blockNumber, "programVersion", programVersion, "programVersionSign", programVersionSign)
	if vote.ProposalID == common.ZeroHash {
		return ProposalIDEmpty
	}

	if vote.VoteOption != Yes && vote.VoteOption != No && vote.VoteOption != Abstention {
		return VoteOptionError
	}

	proposal, err := GetProposal(vote.ProposalID, state)
	if err != nil {
		log.Error("find proposal error", "proposalID", vote.ProposalID)
		return err
	} else if proposal == nil {
		return ProposalNotFound
	}

	//check caller and voter
	if err := checkVerifier(from, vote.VoteNodeID, blockHash, blockNumber, stk); err != nil {
		return err
	}

	if proposal.GetProposalType() == Version {
		if vp, ok := proposal.(*VersionProposal); ok {
			//The signature should be verified when node vote for a version proposal.
			if !handler.GetCryptoHandler().IsSignedByNodeID(programVersion, programVersionSign.Bytes(), vote.VoteNodeID) {
				return VersionSignError
			}

			//vote option can only be Yes for version proposal
			if vote.VoteOption != Yes {
				return VoteOptionError
			}

			if vp.GetNewVersion() != programVersion {
				log.Error("cannot vote for version proposal until node upgrade to a new version", "newVersion", vp.GetNewVersion(), "programVersion", programVersion)
				return VerifierNotUpgraded
			}
		}
	}

	//check if vote.proposalID is in voting
	votingIDs, err := ListVotingProposalID(blockHash, blockNumber, state)
	if err != nil {
		log.Error("list voting proposal error", "blockHash", blockHash, "blockNumber", blockNumber, "err", err)
		return err
	} else if len(votingIDs) == 0 {
		log.Error("there's no voting proposal ID", "blockHash", blockHash, "blockNumber", blockNumber)
		return ProposalNotAtVoting
	} else {
		var isVoting = false
		for _, votingID := range votingIDs {
			if votingID == vote.ProposalID {
				isVoting = true
			}
		}
		if !isVoting {
			return ProposalNotAtVoting
		}
	}

	//check if node has voted
	verifierList, err := ListVotedVerifier(vote.ProposalID, state)
	if err != nil {
		log.Error("list voted verifier error", "proposalID", vote.ProposalID, "blockHash", blockHash, "blockNumber", blockNumber)
		return err
	}

	if xutil.InNodeIDList(vote.VoteNodeID, verifierList) {
		return VoteDuplicated
	}

	//handle storage
	if err := SetVote(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, state); err != nil {
		log.Error("save vote error", "proposalID", vote.ProposalID)
		return err
	}

	//the proposal is version type, so add the node ID to active node list.
	if proposal.GetProposalType() == Version {
		if err := AddActiveNode(blockHash, vote.ProposalID, vote.VoteNodeID); err != nil {
			log.Error("add nodeID to active node list error", "proposalID", vote.ProposalID, "nodeID", byteutil.PrintNodeID(vote.VoteNodeID))
			return err
		}
	}

	return nil
}

// node declares it's version
func DeclareVersion(from common.Address, declaredNodeID discover.NodeID, declaredVersion uint32, programVersionSign common.VersionSign, blockHash common.Hash, blockNumber uint64, stk Staking, state xcom.StateDB) error {
	log.Debug("call DeclareVersion", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "versionSign", programVersionSign)

	if !handler.GetCryptoHandler().IsSignedByNodeID(declaredVersion, programVersionSign.Bytes(), declaredNodeID) {
		return VersionSignError
	}

	if err := checkCandidate(from, declaredNodeID, blockHash, blockNumber, stk); err != nil {
		return err
	}

	activeVersion := GetCurrentActiveVersion(state)
	if activeVersion <= 0 {
		return ActiveVersionError
	}

	votingVP, err := FindVotingVersionProposal(blockHash, state)
	if err != nil {
		log.Error("find voting version proposal error", "blockHash", blockHash)
		return err
	}

	//there is a voting version proposal
	if votingVP != nil {
		log.Debug("there is a version proposal at voting stage", "proposal", votingVP)

		votedList, err := ListVotedVerifier(votingVP.ProposalID, state)
		if err != nil {
			log.Error("list voted verifier error", "proposalID", votingVP.ProposalID)
			return err
		}
		if xutil.InNodeIDList(declaredNodeID, votedList) {
			if declaredVersion != votingVP.GetNewVersion() {
				return DeclareVersionError
			}
		} else if declaredVersion>>8 == activeVersion>>8 {
			//there's a voting-version-proposal, if the declared version equals the current active version, notify staking immediately
			log.Debug("call stk.DeclarePromoteNotify(not voted, declaredVersion==activeVersion)", "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "activeVersion", activeVersion, "blockHash", blockHash, "blockNumber", blockNumber)
			if err := stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion); err != nil {
				log.Error("call stk.DeclarePromoteNotify failed", "err", err)
				return NotifyStakingDeclaredVersionError
			}
		} else if declaredVersion>>8 == votingVP.GetNewVersion()>>8 {
			//the declared version equals the new version, will notify staking when the proposal is passed
			log.Debug("add node to activeNodeList(not voted, declaredVersion==newVersion.", "newVersion", votingVP.GetNewVersion, "declaredVersion", declaredVersion)
			if err := AddActiveNode(blockHash, votingVP.ProposalID, declaredNodeID); err != nil {
				log.Error("add declared node ID to active node list failed", "err", err)
				return err
			}
		} else {
			log.Error("declared version should be either active version or new version", "activeVersion", activeVersion, "newVersion", votingVP.GetNewVersion, "declaredVersion", declaredVersion)
			return DeclareVersionError
		}
	} else {
		log.Debug("there is no version proposal at voting stage")
		preActiveVersion := GetPreActiveVersion(state)
		if preActiveVersion == 0 && declaredVersion>>8 == activeVersion>>8 {
			log.Debug("there is no pre-active version proposal")
			log.Debug("call stk.DeclarePromoteNotify", "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "activeVersion", activeVersion, "blockHash", blockHash, "blockNumber", blockNumber)
			if err := stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion); err != nil {
				log.Error("call stk.DeclarePromoteNotify failed", "err", err)
				return NotifyStakingDeclaredVersionError
			}
		} else if preActiveVersion != 0 && declaredVersion == preActiveVersion {
			log.Debug("there is a version proposal at voting stage")
			log.Debug("call stk.DeclarePromoteNotify", "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "activeVersion", activeVersion, "blockHash", blockHash, "blockNumber", blockNumber)
			if err := stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion); err != nil {
				log.Error("call stk.DeclarePromoteNotify failed", "err", err)
				return NotifyStakingDeclaredVersionError
			}
		} else {
			log.Error("declared version should be either active version or pre-active version", "activeVersion", activeVersion, "declaredVersion", declaredVersion)
			return DeclareVersionError
		}
	}
	return nil
}

// check if the node a verifier, and the caller address is same as the staking address
func checkVerifier(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64, stk Staking) error {
	log.Debug("call checkVerifier", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
	verifierList, err := stk.GetVerifierList(blockHash, blockNumber, false)
	if err != nil {
		log.Error("list verifiers error", "blockHash", blockHash, "err", err)
		return err
	}

	xcom.PrintObject("checkVerifier", verifierList)

	for _, verifier := range verifierList {
		if verifier != nil && verifier.NodeId == nodeID {
			if verifier.StakingAddress == from {
				nodeAddress, _ := xutil.NodeId2Addr(verifier.NodeId)
				candidate, err := stk.GetCandidateInfo(blockHash, nodeAddress)
				if err != nil {
					return VerifierInfoNotFound
				} else if staking.Is_Invalid(candidate.Status) {
					return VerifierStatusInvalid
				}
				log.Debug("tx sender is a valid verifier.", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
				return nil
			} else {
				return TxSenderDifferFromStaking
			}
		}
	}
	log.Error("tx sender is not a verifier", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
	return TxSenderIsNotVerifier
}

// query proposal list
func ListProposal(blockHash common.Hash, state xcom.StateDB) ([]Proposal, error) {
	log.Debug("call ListProposal")
	var proposalIDs []common.Hash
	var proposals []Proposal

	votingProposals, err := ListVotingProposal(blockHash)
	if err != nil {
		log.Error("list voting proposal error", "blockHash", blockHash)
		return nil, err
	}
	endProposals, err := ListEndProposalID(blockHash)
	if err != nil {
		log.Error("list end proposals error", "blockHash", blockHash)
		return nil, err
	}

	preActiveProposals, err := GetPreActiveProposalID(blockHash)
	if err != nil {
		log.Error("find pre-active proposal error", "blockHash", blockHash)
		return nil, err
	}

	proposalIDs = append(proposalIDs, votingProposals...)
	proposalIDs = append(proposalIDs, endProposals...)
	if preActiveProposals != common.ZeroHash {
		proposalIDs = append(proposalIDs, preActiveProposals)
	}

	for _, proposalID := range proposalIDs {
		proposal, err := GetExistProposal(proposalID, state)
		if err != nil {
			log.Error("find proposal error", "proposalID", proposalID)
			return nil, err
		}
		proposals = append(proposals, proposal)
	}
	return proposals, nil
}

// list all proposal IDs at voting stage
func ListVotingProposalID(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) ([]common.Hash, error) {
	log.Debug("call ListVotingProposalID", "blockHash", blockHash, "blockNumber", blockNumber)
	idList, err := ListVotingProposal(blockHash)
	if err != nil {
		log.Error("find voting version proposal error", "blockHash", blockHash)
		return nil, err
	}
	return idList, nil
}

// find a cancel proposal at voting stage
func FindVotingCancelProposal(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (*CancelProposal, error) {
	log.Debug("call findVotingCancelProposal", "blockHash", blockHash, "blockNumber", blockNumber)
	idList, err := ListVotingProposal(blockHash)
	if err != nil {
		log.Error("find voting proposal error", "blockHash", blockHash)
		return nil, err
	}
	for _, proposalID := range idList {
		p, err := GetExistProposal(proposalID, state)
		if err != nil {
			return nil, err
		}
		if p.GetProposalType() == Cancel {
			vp := p.(*CancelProposal)
			return vp, nil
		}
	}
	return nil, nil
}

// GetMaxEndVotingBlock returns the max endVotingBlock of proposals those are at voting stage, and the nodeID has voted for those proposals.
// or returns 0 if there's no proposal at voting stage, or nodeID didn't voted for any proposal.
// if any error happened, return 0 and the error
func GetMaxEndVotingBlock(nodeID discover.NodeID, blockHash common.Hash, state xcom.StateDB) (uint64, error) {
	if proposalIDList, err := ListVotingProposal(blockHash); err != nil {
		return 0, err
	} else {
		var maxEndVotingBlock = uint64(0)
		for _, proposalID := range proposalIDList {
			if voteValueList, err := ListVoteValue(proposalID, state); err != nil {
				return 0, err
			} else {
				for _, voteValue := range voteValueList {
					if voteValue.VoteNodeID == nodeID {
						if proposal, err := GetExistProposal(proposalID, state); err != nil {
							return 0, err
						} else if proposal.GetEndVotingBlock() > maxEndVotingBlock {
							maxEndVotingBlock = proposal.GetEndVotingBlock()
						}
					}
				}
			}
		}
		return maxEndVotingBlock, nil
	}
}

// check if the node a candidate, and the caller address is same as the staking address
func checkCandidate(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64, stk Staking) error {
	log.Debug("call checkCandidate", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
	candidateList, err := stk.GetCandidateList(blockHash, blockNumber)
	if err != nil {
		log.Error("list candidates error", "blockHash", blockHash)
		return err
	}

	for _, candidate := range candidateList {
		if candidate.NodeId == nodeID {
			if candidate.StakingAddress == from {
				//log.Debug("tx sender is a candidate", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
				return nil
			} else {
				return TxSenderDifferFromStaking
			}
		}
	}
	return TxSenderIsNotCandidate
}
