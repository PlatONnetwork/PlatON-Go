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

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type Staking interface {
	GetVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.ValidatorExQueue, error)
	ListVerifierNodeID(blockHash common.Hash, blockNumber uint64) ([]discover.NodeID, error)
	GetCanBaseList(blockHash common.Hash, blockNumber uint64) (staking.CandidateBaseQueue, error)
	GetCandidateInfo(blockHash common.Hash, addr common.NodeAddress) (*staking.Candidate, error)
	GetCanBase(blockHash common.Hash, addr common.NodeAddress) (*staking.CandidateBase, error)
	GetCanMutable(blockHash common.Hash, addr common.NodeAddress) (*staking.CandidateMutable, error)
	DeclarePromoteNotify(blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID, programVersion uint32) error
}

const (
	ModuleStaking  = "staking"
	ModuleSlashing = "slashing"
	ModuleBlock    = "block"
	ModuleTxPool   = "txPool"
	ModuleReward   = "reward"
)

const (
	KeyStakeThreshold             = "stakeThreshold"
	KeyOperatingThreshold         = "operatingThreshold"
	KeyMaxValidators              = "maxValidators"
	KeyUnStakeFreezeDuration      = "unStakeFreezeDuration"
	KeySlashFractionDuplicateSign = "slashFractionDuplicateSign"
	KeyDuplicateSignReportReward  = "duplicateSignReportReward"
	KeyMaxEvidenceAge             = "maxEvidenceAge"
	KeySlashBlocksReward          = "slashBlocksReward"
	KeyMaxBlockGasLimit           = "maxBlockGasLimit"
	KeyMaxTxDataLimit             = "maxTxDataLimit"
	KeyZeroProduceNumberThreshold = "zeroProduceNumberThreshold"
	KeyZeroProduceCumulativeTime  = "zeroProduceCumulativeTime"
	KeyRewardPerMaxChangeRange    = "rewardPerMaxChangeRange"
	KeyRewardPerChangeInterval    = "rewardPerChangeInterval"
	KeyIncreaseIssuanceRatio      = "increaseIssuanceRatio"
)

func GetVersionForStaking(blockHash common.Hash, state xcom.StateDB) uint32 {
	preActiveVersion := GetPreActiveVersion(blockHash)
	if preActiveVersion > 0 {
		return preActiveVersion
	} else {
		return GetCurrentActiveVersion(state)
	}
}

// Get current active version record
func GetCurrentActiveVersion(state xcom.StateDB) uint32 {
	avList, err := ListActiveVersion(state)
	if err != nil {
		log.Error("Cannot find active version list", "err", err)
		return 0
	}

	var version uint32
	if len(avList) == 0 {
		log.Warn("cannot find current active version, The ActiveVersion List is nil")
		return 0
	} else {
		version = avList[0].ActiveVersion
	}
	return version
}

// submit a proposal
func Submit(from common.Address, proposal Proposal, blockHash common.Hash, blockNumber uint64, stk Staking, state xcom.StateDB, chainID *big.Int) error {
	log.Debug("call Submit", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "proposal", proposal)

	//param check
	if err := proposal.Verify(blockNumber, blockHash, state, chainID); err != nil {
		if bizError, ok := err.(*common.BizError); ok {
			return bizError
		} else {
			log.Error("verify proposal parameters failed", "err", err)
			return common.InvalidParameter.Wrap(err.Error())
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

	verifierList, err := stk.ListVerifierNodeID(blockHash, blockNumber)
	if err != nil {
		return err
	}
	log.Debug("verifiers count of current settlement", "verifierCount", len(verifierList))

	if err := AccuVerifiers(blockHash, proposal.GetProposalID(), verifierList); err != nil {
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
			if !node.GetCryptoHandler().IsSignedByNodeID(programVersion, programVersionSign.Bytes(), vote.VoteNodeID) {
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
	votingIDs, err := ListVotingProposalID(blockHash)
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
	votedMap, err := GetVotedVerifierMap(vote.ProposalID, blockHash)
	if err != nil {
		log.Error("get voted verifier map error", "proposalID", vote.ProposalID, "blockHash", blockHash, "blockNumber", blockNumber)
		return err
	}

	if _, exist := votedMap[vote.VoteNodeID]; exist {
		return VoteDuplicated
	}

	//handle storage
	if err := AddVoteValue(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, blockHash); err != nil {
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

	if !node.GetCryptoHandler().IsSignedByNodeID(declaredVersion, programVersionSign.Bytes(), declaredNodeID) {
		return VersionSignError
	}

	if err := checkCandidate(from, declaredNodeID, blockHash, blockNumber, stk); err != nil {
		return err
	}

	activeVersion := GetCurrentActiveVersion(state)
	if activeVersion <= 0 {
		return ActiveVersionError
	}

	proposal, err := FindVotingProposal(blockHash, state, Version)
	if err != nil {
		log.Error("find voting version proposal error", "blockHash", blockHash)
		return err
	}

	//there is a voting version proposal
	if proposal != nil {
		votingVP := proposal.(*VersionProposal)

		log.Debug("there is a version proposal at voting stage", "proposal", votingVP)

		votedMap, err := GetVotedVerifierMap(votingVP.ProposalID, blockHash)
		if err != nil {
			log.Error("get voted verifier map error", "proposalID", votingVP.ProposalID)
			return err
		}
		//if xutil.InNodeIDList(declaredNodeID, votedList) {
		if _, exist := votedMap[declaredNodeID]; exist {
			if declaredVersion>>8 != votingVP.GetNewVersion()>>8 {
				log.Error("node voted new version, then declared version, the major is different between the declared version and new version")
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
		preActiveVersion := GetPreActiveVersion(blockHash)
		if preActiveVersion <= 0 {
			log.Debug("there is no version proposal at pre-active stage")
			if declaredVersion>>8 == activeVersion>>8 {
				log.Debug("call stk.DeclarePromoteNotify", "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "activeVersion", activeVersion, "blockHash", blockHash, "blockNumber", blockNumber)
				if err := stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion); err != nil {
					log.Error("call stk.DeclarePromoteNotify failed", "err", err)
					return NotifyStakingDeclaredVersionError
				}
			} else {
				log.Error("declared version should be active version", "activeVersion", activeVersion, "declaredVersion", declaredVersion)
				return DeclareVersionError
			}
		} else {
			log.Debug("there is a version proposal at pre-active stage", "preActiveVersion", preActiveVersion)
			if declaredVersion>>8 == preActiveVersion>>8 {
				log.Debug("call stk.DeclarePromoteNotify", "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "activeVersion", activeVersion, "blockHash", blockHash, "blockNumber", blockNumber)
				if err := stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion); err != nil {
					log.Error("call stk.DeclarePromoteNotify failed", "err", err)
					return NotifyStakingDeclaredVersionError
				}
			} else {
				log.Error("declared version should be pre-active version", "activeVersion", activeVersion, "declaredVersion", declaredVersion)
				return DeclareVersionError
			}
		}
	}
	return nil
}

// check if the node a verifier, and the caller address is same as the staking address
func checkVerifier(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64, stk Staking) error {
	log.Debug("call checkVerifier", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)

	_, err := xutil.NodeId2Addr(nodeID)
	if nil != err {
		log.Error("parse nodeID error", "err", err)
		return err
	}

	verifierList, err := stk.GetVerifierList(blockHash, blockNumber, false)
	if err != nil {
		log.Error("list verifiers error", "blockHash", blockHash, "err", err)
		return err
	}

	//xcom.PrintObject("checkVerifier", verifierList)

	for _, verifier := range verifierList {
		if verifier != nil && verifier.NodeId == nodeID {
			if verifier.StakingAddress == from {
				nodeAddress, err := xutil.NodeId2Addr(verifier.NodeId)
				if err != nil {
					return err
				}
				candidate, err := stk.GetCanMutable(blockHash, nodeAddress)
				if candidate == nil || err != nil {
					return VerifierInfoNotFound
				} else if candidate.IsInvalid() {
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
func ListVotingProposalID(blockHash common.Hash) ([]common.Hash, error) {
	log.Debug("call ListVotingProposalID", "blockHash", blockHash)
	idList, err := ListVotingProposal(blockHash)
	if err != nil {
		log.Error("find voting version proposal error", "blockHash", blockHash)
		return nil, err
	}
	return idList, nil
}

// find a proposal at voting stage
func FindVotingProposal(blockHash common.Hash, state xcom.StateDB, proposalTypes ...ProposalType) (Proposal, error) {

	if len(proposalTypes) == 0 {
		return nil, common.InvalidParameter
	}
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

		for _, typ := range proposalTypes {
			if p.GetProposalType() == typ {
				return p, nil
			}
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
			if voteValueList, err := ListVoteValue(proposalID, blockHash); err != nil {
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

// NotifyPunishedVerifiers receives punished verifies notification from Staking
func NotifyPunishedVerifiers(blockHash common.Hash, punishedVerifierMap map[discover.NodeID]struct{}, state xcom.StateDB) error {
	if punishedVerifierMap == nil || len(punishedVerifierMap) == 0 {
		return nil
	}
	if votingProposalIDList, err := ListVotingProposalID(blockHash); err != nil {
		return err
	} else if len(votingProposalIDList) > 0 {
		for _, proposalID := range votingProposalIDList {
			if voteValueList, err := ListVoteValue(proposalID, blockHash); err != nil {
				return err
			} else if len(voteValueList) > 0 {
				idx := 0 // output index
				removed := make([]VoteValue, 0)
				for _, voteValue := range voteValueList {
					//if !xutil.InNodeIDList(voteValue.VoteNodeID, punishedVerifiers) {
					if _, isPunished := punishedVerifierMap[voteValue.VoteNodeID]; !isPunished {
						voteValueList[idx] = voteValue
						idx++
					} else {
						removed = append(removed, voteValue)
					}
				}
				if len(removed) > 0 && idx < len(voteValueList) {
					voteValueList = voteValueList[:idx]
					log.Debug(fmt.Sprintf("remove voted value, proposalID:%s, removedVoteValue:%+v", proposalID.Hex(), removed))
					if err := UpdateVoteValue(proposalID, voteValueList, blockHash); err != nil {
						return err
					}
				}
			}

			/*if verifierList, err := ListAccuVerifier(blockHash, proposalID); err != nil {
				return err
			} else if len(verifierList) > 0 {
				idx := 0 // output index
				for _, verifier := range verifierList {
					if !xutil.InNodeIDList(verifier, punishedVerifiers) {
						verifierList[idx] = verifier
						idx++
					}
				}
				verifierList = verifierList[:idx]
				//UpdateAccuVerifiers(blockHash, voteList)
			}*/
		}
	}
	return nil
}

func ClearProcessingProposals(blockHash common.Hash, state xcom.StateDB) error {
	if votingIDList, err := ListVotingProposalID(blockHash); err != nil {
		return err
	} else {
		for _, votingID := range votingIDList {
			if err := clearProcessingProposal(votingID, true, blockHash, state); err != nil {
				return err
			}
		}
	}

	if preactiveID, err := GetPreActiveProposalID(blockHash); err != nil {
		log.Error(" find pre-active proposal ID failed", "blockHash", blockHash)
		return err
	} else if preactiveID != common.ZeroHash {
		if err := clearProcessingProposal(preactiveID, false, blockHash, state); err != nil {
			return err
		}
	}
	return nil
}

func clearProcessingProposal(proposalID common.Hash, isVoting bool, blockHash common.Hash, state xcom.StateDB) error {
	if isVoting {
		if err := MoveVotingProposalIDToEnd(proposalID, blockHash); err != nil {
			log.Error("move proposalID from voting proposalID list to end list failed", "proposalID", proposalID, "blockHash", blockHash)
			return err
		}
	} else {
		if err := MovePreActiveProposalIDToEnd(blockHash, proposalID); err != nil {
			log.Error("move pre-active proposal ID to end list failed", "proposalID", proposalID, "blockHash", blockHash)
			return err
		}

		if err := delPreActiveVersion(blockHash); err != nil {
			log.Error("delete pre-active version failed", "blockHash", blockHash)
			return err
		}
	}

	if err := ClearVoteValue(proposalID, blockHash); err != nil {
		log.Error("clear vote values failed", "proposalID", proposalID, "blockHash", blockHash)
		return err
	}
	if err := ClearAccuVerifiers(blockHash, proposalID); err != nil {
		log.Error("clear voted verifiers failed", "proposalID", proposalID, "blockHash", blockHash.Hex(), "error", err)
		return err
	}
	tallyResult := &TallyResult{
		ProposalID:    proposalID,
		Yeas:          0x0,
		Nays:          0x0,
		Abstentions:   0x0,
		AccuVerifiers: 0x0,
		Status:        Failed,
	}
	if err := SetTallyResult(*tallyResult, state); err != nil {
		log.Error("save tally result failed", "proposalID", proposalID, "blockHash", blockHash)
		return err
	}
	return nil
}

func SetGovernParam(module, name, desc, initValue string, activeBlockNumber uint64, currentBlockHash common.Hash) error {
	paramValue := &ParamValue{"", initValue, activeBlockNumber}
	return addGovernParam(module, name, desc, paramValue, currentBlockHash)
}

func UpdateGovernParamValue(module, name string, newValue string, activeBlock uint64, blockHash common.Hash) error {
	return updateGovernParamValue(module, name, newValue, activeBlock, blockHash)
}

func ListGovernParam(module string, blockHash common.Hash) ([]*GovernParam, error) {
	return listGovernParam(module, blockHash)
}

func FindGovernParam(module, name string, blockHash common.Hash) (*GovernParam, error) {
	itemList, err := listGovernParamItem(module, blockHash)
	if err != nil {
		return nil, err
	}
	for _, item := range itemList {
		if item.Name == name {
			if value, err := findGovernParamValue(module, name, blockHash); err != nil {
				return nil, err
			} else if value != nil {
				param := &GovernParam{item, value, nil}
				return param, nil
			}
		}
	}
	return nil, nil
}

// check if the node a candidate, and the caller address is same as the staking address
func checkCandidate(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64, stk Staking) error {

	_, err := xutil.NodeId2Addr(nodeID)
	if nil != err {
		log.Error("parse nodeID error", "err", err)
		return err
	}

	candidateList, err := stk.GetCanBaseList(blockHash, blockNumber)
	if err != nil {
		log.Error("list candidates error", "blockHash", blockHash)
		return err
	}

	for _, candidate := range candidateList {
		if candidate.NodeId == nodeID {
			if candidate.StakingAddress == from {
				return nil
			} else {
				return TxSenderDifferFromStaking
			}
		}
	}
	return TxSenderIsNotCandidate
}

type ParamVerifier func(blockNumber uint64, blockHash common.Hash, value string) error

func GetGovernParamValue(module, name string, blockNumber uint64, blockHash common.Hash) (string, error) {
	paramValue, err := findGovernParamValue(module, name, blockHash)
	if err != nil {
		log.Error("get govern parameter value failed", "module", module, "name", name, "blockNumber", blockNumber, "blockHash", blockHash, "err", err)
		return "", err
	}
	if paramValue == nil {
		log.Error("govern parameter value is nil", "module", module, "name", name, "blockNumber", blockNumber, "blockHash", blockHash, "err", err)
		return "", UnsupportedGovernParam
	} else {
		if blockNumber >= paramValue.ActiveBlock {
			return paramValue.Value, nil
		} else {
			return paramValue.StaleValue, nil
		}
	}
}

func GovernStakeThreshold(blockNumber uint64, blockHash common.Hash) (*big.Int, error) {
	thresholdStr, err := GetGovernParamValue(ModuleStaking, KeyStakeThreshold, blockNumber, blockHash)
	if nil != err {
		return new(big.Int).SetInt64(0), err
	}

	threshold, ok := new(big.Int).SetString(thresholdStr, 10)
	if !ok {
		return new(big.Int).SetInt64(0), fmt.Errorf("Failed to parse the govern stakethreshold")
	}

	return threshold, nil
}

func GovernOperatingThreshold(blockNumber uint64, blockHash common.Hash) (*big.Int, error) {
	thresholdStr, err := GetGovernParamValue(ModuleStaking, KeyOperatingThreshold, blockNumber, blockHash)
	if nil != err {
		return new(big.Int).SetInt64(0), err
	}

	threshold, ok := new(big.Int).SetString(thresholdStr, 10)
	if !ok {
		return new(big.Int).SetInt64(0), fmt.Errorf("Failed to parse the govern operatingthreshold")
	}

	return threshold, nil
}

func GovernMaxValidators(blockNumber uint64, blockHash common.Hash) (uint64, error) {
	maxvalidatorsStr, err := GetGovernParamValue(ModuleStaking, KeyMaxValidators, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	maxvalidators, err := strconv.Atoi(maxvalidatorsStr)
	if nil != err {
		return 0, err
	}

	return uint64(maxvalidators), nil
}

func GovernUnStakeFreezeDuration(blockNumber uint64, blockHash common.Hash) (uint64, error) {
	durationStr, err := GetGovernParamValue(ModuleStaking, KeyUnStakeFreezeDuration, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	duration, err := strconv.Atoi(durationStr)
	if nil != err {
		return 0, err
	}

	return uint64(duration), nil
}

func GovernSlashFractionDuplicateSign(blockNumber uint64, blockHash common.Hash) (uint32, error) {
	fractionStr, err := GetGovernParamValue(ModuleSlashing, KeySlashFractionDuplicateSign, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	fraction, err := strconv.Atoi(fractionStr)
	if nil != err {
		return 0, err
	}

	return uint32(fraction), nil
}

func GovernDuplicateSignReportReward(blockNumber uint64, blockHash common.Hash) (uint32, error) {
	rewardStr, err := GetGovernParamValue(ModuleSlashing, KeyDuplicateSignReportReward, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	reward, err := strconv.Atoi(rewardStr)
	if nil != err {
		return 0, err
	}

	return uint32(reward), nil
}

func GovernMaxEvidenceAge(blockNumber uint64, blockHash common.Hash) (uint32, error) {
	ageStr, err := GetGovernParamValue(ModuleSlashing, KeyMaxEvidenceAge, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	age, err := strconv.Atoi(ageStr)
	if nil != err {
		return 0, err
	}

	return uint32(age), nil
}

func GovernSlashBlocksReward(blockNumber uint64, blockHash common.Hash) (uint32, error) {
	rewardStr, err := GetGovernParamValue(ModuleSlashing, KeySlashBlocksReward, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	reward, err := strconv.Atoi(rewardStr)
	if nil != err {
		return 0, err
	}

	return uint32(reward), nil
}

func GovernMaxBlockGasLimit(blockNumber uint64, blockHash common.Hash) (int, error) {
	gasLimitStr, err := GetGovernParamValue(ModuleBlock, KeyMaxBlockGasLimit, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	gasLimit, err := strconv.Atoi(gasLimitStr)
	if nil != err {
		return 0, err
	}

	return gasLimit, nil
}

//func GovernMaxTxDataLimit(blockNumber uint64, blockHash common.Hash) (int, error) {
//	sizeStr, err := GetGovernParamValue(ModuleTxPool, KeyMaxTxDataLimit, blockNumber, blockHash)
//	if nil != err {
//		return 0, err
//	}
//
//	size, err := strconv.Atoi(sizeStr)
//	if nil != err {
//		return 0, err
//	}
//
//	return size, nil
//}

func GovernZeroProduceNumberThreshold(blockNumber uint64, blockHash common.Hash) (uint16, error) {
	valueStr, err := GetGovernParamValue(ModuleSlashing, KeyZeroProduceNumberThreshold, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	value, err := strconv.Atoi(valueStr)
	if nil != err {
		return 0, err
	}

	return uint16(value), nil
}

func GovernZeroProduceCumulativeTime(blockNumber uint64, blockHash common.Hash) (uint16, error) {
	valueStr, err := GetGovernParamValue(ModuleSlashing, KeyZeroProduceCumulativeTime, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	value, err := strconv.Atoi(valueStr)
	if nil != err {
		return 0, err
	}

	return uint16(value), nil
}

func GovernRewardPerMaxChangeRange(blockNumber uint64, blockHash common.Hash) (uint16, error) {
	valueStr, err := GetGovernParamValue(ModuleStaking, KeyRewardPerMaxChangeRange, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	value, err := strconv.Atoi(valueStr)
	if nil != err {
		return 0, err
	}

	return uint16(value), nil
}

func GovernRewardPerChangeInterval(blockNumber uint64, blockHash common.Hash) (uint16, error) {
	valueStr, err := GetGovernParamValue(ModuleStaking, KeyRewardPerChangeInterval, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	value, err := strconv.Atoi(valueStr)
	if nil != err {
		return 0, err
	}

	return uint16(value), nil
}

func GovernIncreaseIssuanceRatio(blockNumber uint64, blockHash common.Hash) (uint16, error) {
	valueStr, err := GetGovernParamValue(ModuleReward, KeyIncreaseIssuanceRatio, blockNumber, blockHash)
	if nil != err {
		return 0, err
	}

	value, err := strconv.Atoi(valueStr)
	if nil != err {
		return 0, err
	}

	return uint16(value), nil
}
