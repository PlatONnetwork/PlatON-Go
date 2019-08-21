package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type Staking interface {
	GetVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.ValidatorExQueue, error)
	GetCandidateList(blockHash common.Hash, blockNumber uint64) (staking.CandidateQueue, error)
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
	sig, err := xcom.GetCryptoHandler().Sign(programVersion)
	if err != nil {
		log.Error("sign version data error")
		return nil, err
	}
	value := &ProgramVersionValue{ProgramVersion: programVersion, ProgramVersionSign: common.BytesToVersionSign(sig)}
	return value, nil
}

// submit a proposal
func Submit(from common.Address, proposal Proposal, blockHash common.Hash, blockNumber uint64, stk Staking, state xcom.StateDB) error {
	log.Debug("call Submit", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "proposal", proposal)

	//param check
	if err := proposal.Verify(blockNumber, blockHash, state); err != nil {
		log.Error("verify proposal parameters failed", "err", err)
		return common.NewBizError(err.Error())
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
	log.Debug("call Vote", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "programVersion", programVersion, "programVersionSign", programVersionSign, "voteInfo", vote)
	if vote.ProposalID == common.ZeroHash || vote.VoteOption == 0 {
		return common.NewBizError("empty parameter detected.")
	}

	proposal, err := GetProposal(vote.ProposalID, state)
	if err != nil {
		log.Error("cannot find proposal by ID", "proposalID", vote.ProposalID)
		return err
	} else if proposal == nil {
		log.Error("incorrect proposal ID.", "proposalID", vote.ProposalID)
		return common.NewBizError("incorrect proposal ID.")
	}

	//check caller and voter
	if err := checkVerifier(from, vote.VoteNodeID, blockHash, blockNumber, stk); err != nil {
		return err
	}

	//voteOption range check
	if !(vote.VoteOption >= Yes && vote.VoteOption <= Abstention) {
		return common.NewBizError("vote option is error.")
	}

	if proposal.GetProposalType() == Version {
		if vp, ok := proposal.(*VersionProposal); ok {
			//The signature should be verified when node vote for a version proposal.
			if !xcom.GetCryptoHandler().IsSignedByNodeID(programVersion, programVersionSign.Bytes(), vote.VoteNodeID) {
				return common.NewBizError("version sign error.")
			}

			//vote option can only be Yes for version proposal
			if vote.VoteOption != Yes {
				return common.NewBizError("vote option error.")
			}

			if vp.GetNewVersion() != programVersion {
				log.Error("cannot vote for version proposal until node upgrade to a new version", "newVersion", vp.GetNewVersion(), "programVersion", programVersion)
				return common.NewBizError("node have not upgraded to a new version")
			}
		}
	}

	//check if vote.proposalID is in voting
	votingIDs, err := ListVotingProposalID(blockHash, blockNumber, state)
	if err != nil {
		log.Error("list all voting proposal IDs failed", "blockHash", blockHash)
		return err
	} else if len(votingIDs) == 0 {
		log.Error("there's no voting proposal ID.", "blockHash", blockHash)
		return err
	} else {
		var isVoting = false
		for _, votingID := range votingIDs {
			if votingID == vote.ProposalID {
				isVoting = true
			}
		}
		if !isVoting {
			log.Error("proposal is not at voting stage", "proposalID", vote.ProposalID)
			return common.NewBizError("Proposal is not at voting stage.")
		}
	}

	//check if node has voted
	verifierList, err := ListVotedVerifier(vote.ProposalID, state)
	if err != nil {
		log.Error("list voted verifiers failed", "proposalID", vote.ProposalID)
		return err
	}

	if xutil.InNodeIDList(vote.VoteNodeID, verifierList) {
		log.Error("node has voted this proposal", "proposalID", vote.ProposalID, "nodeID", byteutil.PrintNodeID(vote.VoteNodeID))
		return common.NewBizError("node has voted this proposal.")
	}

	//handle storage
	if err := SetVote(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, state); err != nil {
		log.Error("save vote failed", "proposalID", vote.ProposalID)
		return err
	}

	//the proposal is version type, so add the node ID to active node list.
	if proposal.GetProposalType() == Version {
		if err := AddActiveNode(blockHash, vote.ProposalID, vote.VoteNodeID); err != nil {
			log.Error("add nodeID to active node list failed", "proposalID", vote.ProposalID, "nodeID", byteutil.PrintNodeID(vote.VoteNodeID))
			return err
		}
	}

	return nil
}

// node declares it's version
func DeclareVersion(from common.Address, declaredNodeID discover.NodeID, declaredVersion uint32, programVersionSign common.VersionSign, blockHash common.Hash, blockNumber uint64, stk Staking, state xcom.StateDB) error {

	log.Debug("call DeclareVersion", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "versionSign", programVersionSign)
	//check caller is a Verifier or Candidate
	/*if err := govPlugin.checkVerifier(from, declaredNodeID, blockHash, blockNumber); err != nil {
		return err
	}*/

	if !xcom.GetCryptoHandler().IsSignedByNodeID(declaredVersion, programVersionSign.Bytes(), declaredNodeID) {
		return common.NewBizError("version sign error.")
	}

	if err := checkCandidate(from, declaredNodeID, blockHash, blockNumber, stk); err != nil {
		return err
	}

	activeVersion := GetCurrentActiveVersion(state)
	if activeVersion <= 0 {
		return common.NewBizError("wrong current active version.")
	}

	votingVP, err := FindVotingVersionProposal(blockHash, blockNumber, state)
	if err != nil {
		log.Error("find if there's a voting version proposal failed", "blockHash", blockHash)
		return err
	}

	//there is a voting version proposal
	if votingVP != nil {
		if declaredVersion>>8 == activeVersion>>8 {
			nodeList, err := ListVotedVerifier(votingVP.ProposalID, state)
			if err != nil {
				log.Error("cannot list voted verifiers", "proposalID", votingVP.ProposalID)
				return err
			} else {
				if xutil.InNodeIDList(declaredNodeID, nodeList) && declaredVersion != votingVP.GetNewVersion() {
					log.Error("declared version should be new version",
						"declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "proposalID", votingVP.ProposalID, "newVersion", votingVP.GetNewVersion())
					return common.NewBizError("declared version should be same as proposal's version")
				} else {
					//there's a voting-version-proposal, if the declared version equals the current active version, notify staking immediately
					log.Debug("there's a voting-version-proposal, and declared version equals active version, notify staking immediately.",
						"blockNumber", blockNumber, "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "activeVersion", activeVersion)
					if err := stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion); err != nil {
						log.Error("notify staking of declared node ID failed", "err", err)
						return common.NewBizError("notify staking of declared node ID failed")
					}
				}
			}
		} else if declaredVersion>>8 == votingVP.GetNewVersion()>>8 {
			//the declared version equals the new version, will notify staking when the proposal is passed
			log.Debug("declared version equals the new version.",
				"newVersion", votingVP.GetNewVersion, "declaredVersion", declaredVersion)
			if err := AddActiveNode(blockHash, votingVP.ProposalID, declaredNodeID); err != nil {
				log.Error("add declared node ID to active node list failed", "err", err)
				return err
			}
		} else {
			log.Error("declared version neither equals active version nor new version.", "activeVersion", activeVersion, "newVersion", votingVP.GetNewVersion, "declaredVersion", declaredVersion)
			return common.NewBizError("declared version neither equals active version nor new version.")
		}
	} else {
		preActiveVersion := GetPreActiveVersion(state)
		if declaredVersion>>8 == activeVersion>>8 || (preActiveVersion != 0 && declaredVersion == preActiveVersion) {
			//there's no voting-version-proposal, if the declared version equals either the current active version or preActive version, notify staking immediately
			//stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion)
			log.Debug("there's no voting-version-proposal, the declared version equals either the current active version or preActive version, notify staking immediately.",
				"blockNumber", blockNumber, "declaredVersion", declaredVersion, "declaredNodeID", declaredNodeID, "activeVersion", activeVersion, "preActiveVersion", preActiveVersion)
			if err := stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion); err != nil {
				log.Error("notify staking of declared node ID failed", "err", err)
				return common.NewBizError("notify staking of declared node ID failed")
			}
		} else {
			log.Error("there's no version proposal at voting stage, declared version should be active or pre-active version.", "activeVersion", activeVersion, "declaredVersion", declaredVersion)
			return common.NewBizError("there's no version proposal at voting stage, declared version should be active version.")
		}
	}
	return nil
}

// check if the node a verifier, and the caller address is same as the staking address
func checkVerifier(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64, stk Staking) error {
	log.Debug("call checkVerifier", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
	verifierList, err := stk.GetVerifierList(blockHash, blockNumber, false)

	if err != nil {
		log.Error("list verifiers failed", "blockHash", blockHash, "err", err)
		return err
	}

	xcom.PrintObject("checkVerifier", verifierList)

	for _, verifier := range verifierList {
		if verifier != nil && verifier.NodeId == nodeID {
			if verifier.StakingAddress == from {
				nodeAddress, _ := xutil.NodeId2Addr(verifier.NodeId)
				candidate, err := stk.GetCandidateInfo(blockHash, nodeAddress)
				if err != nil {
					return common.NewBizError("cannot get verifier's detail info.")
				} else if staking.Is_Invalid(candidate.Status) {
					return common.NewBizError("verifier's status is invalid.")
				}
				log.Debug("tx sender is a valid verifier.", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
				return nil
			} else {
				return common.NewBizError("tx sender should be node's staking address.")
			}
		}
	}
	log.Error("tx sender is not a verifier.", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
	return common.NewBizError("tx sender is not a verifier.")
}

// query proposal list
func ListProposal(blockHash common.Hash, state xcom.StateDB) ([]Proposal, error) {
	log.Debug("call ListProposal")
	var proposalIDs []common.Hash
	var proposals []Proposal

	votingProposals, err := ListVotingProposal(blockHash)
	if err != nil {
		log.Error("list voting proposals failed.", "blockHash", blockHash)
		return nil, err
	}
	endProposals, err := ListEndProposalID(blockHash)
	if err != nil {
		log.Error("list end proposals failed.", "blockHash", blockHash)
		return nil, err
	}

	preActiveProposals, err := GetPreActiveProposalID(blockHash)
	if err != nil {
		log.Error("find pre-active proposal failed.", "blockHash", blockHash)
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
			log.Error("find proposal failed.", "proposalID", proposalID)
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
		log.Error("find voting version proposal failed", "blockHash", blockHash)
		return nil, err
	}
	return idList, nil
}

// find a cancel proposal at voting stage
func FindVotingCancelProposal(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (*CancelProposal, error) {
	log.Debug("call findVotingCancelProposal", "blockHash", blockHash, "blockNumber", blockNumber)
	idList, err := ListVotingProposal(blockHash)
	if err != nil {
		log.Error("find voting proposal failed", "blockHash", blockHash)
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
		log.Error("list candidates failed", "blockHash", blockHash)
		return err
	}

	for _, candidate := range candidateList {
		if candidate.NodeId == nodeID {
			if candidate.StakingAddress == from {
				log.Debug("tx sender is a candidate.", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
				return nil
			} else {
				return common.NewBizError("tx sender should be node's staking address.")
			}
		}
	}
	return common.NewBizError("tx sender is not candidate.")
}
