package plugin

import (
	"errors"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

var (
	govPluginOnce sync.Once
)

type GovPlugin struct {
	govDB *gov.GovDB
}

var govp *GovPlugin

func GovPluginInstance() *GovPlugin {
	govPluginOnce.Do(func() {
		log.Info("Init Governance plugin ...")
		govp = &GovPlugin{govDB: gov.GovDBInstance()}
	})
	return govp
}

func (govPlugin *GovPlugin) Confirmed(block *types.Block) error {
	return nil
}

//implement BasePlugin
func (govPlugin *GovPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	return nil
}

//implement BasePlugin
func (govPlugin *GovPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	var blockNumber = header.Number.Uint64()
	log.Debug("call EndBlock()", "blockNumber", blockNumber, "blockHash", blockHash)

	//if current block is a settlement block, to accumulate current verifiers for each voting proposal.
	if xutil.IsSettlementPeriod(blockNumber) {
		log.Debug("settlement block", "blockNumber", blockNumber, "blockHash", blockHash)
		verifierList, err := stk.ListVerifierNodeID(blockHash, blockNumber)
		if err != nil {
			return err
		}

		votingProposalIDs, err := govPlugin.govDB.ListVotingProposal(blockHash, state)
		if err != nil {
			return err
		}

		for _, votingProposalID := range votingProposalIDs {
			if err := govPlugin.govDB.AccuVerifiers(blockHash, votingProposalID, verifierList); err != nil {
				return err
			}
		}

		//According to the proposal's rules, the settlement block must not be the end-voting block, so, just return.
		return nil
	}

	votingProposalIDs, err := govPlugin.govDB.ListVotingProposal(blockHash, state)
	if err != nil {
		return err
	}

	var getNewPreActiveProposal = false
	//iterate each voting proposal, to check if current block is proposal's end-voting block.
	for _, votingProposalID := range votingProposalIDs {
		log.Debug("iterate each voting proposal", "proposalID", votingProposalID)
		votingProposal, err := govPlugin.govDB.GetExistProposal(votingProposalID, state)
		if nil != err {
			return err
		}
		if votingProposal.GetEndVotingBlock() == blockNumber {
			log.Debug("proposal's end-voting block", "proposalID", votingProposal.GetProposalID(), "blockNumber", blockNumber)
			//According to the proposal's rules, the end-voting block must not be the end-voting block, so, to accumulate current verifiers for current voting proposal.
			verifierList, err := stk.ListVerifierNodeID(blockHash, blockNumber)
			if err != nil {
				return err
			}

			if err := govPlugin.govDB.AccuVerifiers(blockHash, votingProposalID, verifierList); err != nil {
				return err
			}

			//tally the results
			if votingProposal.GetProposalType() == gov.Text {
				_, err := govPlugin.tallyBasic(votingProposal.GetProposalID(), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
			} else if votingProposal.GetProposalType() == gov.Version {
				getNewPreActiveProposal, err = govPlugin.tallyForVersionProposal(votingProposal.(gov.VersionProposal), blockHash, blockNumber, state)
				if err != nil {
					return err
				}

			} else if votingProposal.GetProposalType() == gov.Param {
				pass, err := govPlugin.tallyBasic(votingProposal.GetProposalID(), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
				if pass {
					if err := govPlugin.updateParam(votingProposal.(gov.ParamProposal), blockHash, blockNumber, state); err != nil {
						return err
					}
				}
			} else {
				log.Error("invalid proposal type", "type", votingProposal.GetProposalType())
				err = errors.New("invalid proposal type")
				return err
			}
		}
	}

	if getNewPreActiveProposal == true {
		return nil
	}

	preActiveProposalID, err := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)
	if err != nil {
		log.Error("check if there's a preactive proposal failed.", "blockHash", blockHash)
		return err
	}
	if preActiveProposalID == common.ZeroHash {
		return nil
	}

	//handle a PreActiveProposal
	preActiveProposal, err := govPlugin.govDB.GetProposal(preActiveProposalID, state)
	if err != nil {
		return err
	}
	versionProposal, isVersionProposal := preActiveProposal.(gov.VersionProposal)

	if isVersionProposal {
		log.Debug("found pre-active version proposal", "proposalID", preActiveProposalID, "blockHash", blockHash, "blockNumber", blockNumber, "activeBlockNumber", versionProposal.GetActiveBlock())
		if blockNumber >= versionProposal.GetActiveBlock() && (blockNumber-versionProposal.GetActiveBlock())%xutil.ConsensusSize() == 0 {
			currentValidatorList, err := stk.ListCurrentValidatorID(blockHash, blockNumber)
			if err != nil {
				log.Error("list current round validators failed.", "blockHash", blockHash, "blockNumber", blockNumber)
				return err
			}
			var updatedNodes int = 0
			var totalValidators int = len(currentValidatorList)

			//all active validators
			activeList, err := govPlugin.govDB.GetActiveNodeList(blockHash, preActiveProposalID)
			if err != nil {
				log.Error("list all active nodes failed.", "blockHash", blockHash, "blockNumber", blockNumber, "preActiveProposalID", preActiveProposalID)
				return err
			}

			//check if all validators are active
			for _, val := range currentValidatorList {
				if inNodeList(val, activeList) {
					updatedNodes++
				}
			}

			log.Debug("check active criteria", "pre-active nodes", updatedNodes, "total validators", totalValidators)
			if updatedNodes == totalValidators {
				log.Debug("the pre-active version proposal has passed")
				tallyResult, err := govPlugin.govDB.GetTallyResult(preActiveProposalID, state)
				if err != nil {
					log.Error("find tally result by proposal ID failed.", "preActiveProposalID", preActiveProposalID)
					return err
				}
				//change tally status to "active"
				tallyResult.Status = gov.Active
				if err := govPlugin.govDB.SetTallyResult(*tallyResult, state); err != nil {
					log.Error("update tally result failed.", "preActiveProposalID", preActiveProposalID)
					return err
				}
				if err = govPlugin.govDB.MovePreActiveProposalIDToEnd(blockHash, preActiveProposalID, state); err != nil {
					log.Error("move proposal ID from preActiveProposal list to endProposal list failed.", "blockHash", blockHash, "preActiveProposalID", preActiveProposalID)
					return err
				}

				if err = govPlugin.govDB.ClearActiveNodes(blockHash, preActiveProposalID); err != nil {
					log.Error("clear active nodes failed.", "blockHash", blockHash, "preActiveProposalID", preActiveProposalID)
					return err
				}

				if err = govPlugin.govDB.SetActiveVersion(versionProposal.NewVersion, state); err != nil {
					log.Error("save active version to stateDB failed.", "blockHash", blockHash, "preActiveProposalID", preActiveProposalID)
					return err
				}
				log.Debug("PlatON is ready to upgrade to new version.")
			}
		}
	}
	return nil
}

// nil is allowed
func (govPlugin *GovPlugin) GetPreActiveVersion(state xcom.StateDB) uint32 {
	if nil == govPlugin {
		log.Error("The gov instance is nil on GetPreActiveVersion")
	}
	if nil == govPlugin.govDB {
		log.Error("The govDB instance is nil on GetPreActiveVersion")
	}
	return govPlugin.govDB.GetPreActiveVersion(state)
}

// should not be a nil value
func (govPlugin *GovPlugin) GetActiveVersion(state xcom.StateDB) uint32 {
	if nil == govPlugin {
		log.Error("The gov instance is nil on GetActiveVersion")
	}
	if nil == govPlugin.govDB {
		log.Error("The govDB instance is nil on GetActiveVersion")
	}
	return govPlugin.govDB.GetActiveVersion(state)
}

// submit a proposal
func (govPlugin *GovPlugin) Submit(from common.Address, proposal gov.Proposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {
	log.Debug("call Submit", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "proposal", proposal)

	//param check
	if err := proposal.Verify(proposal.GetSubmitBlock(), state); err != nil {
		log.Error("verify proposal parameters failed", "err", err)
		return common.NewBizError(err.Error())
	}

	//check caller and proposer
	if err := govPlugin.checkVerifier(from, proposal.GetProposer(), blockHash, proposal.GetSubmitBlock()); err != nil {
		return err
	}

	//handle version proposal
	_, isVP := proposal.(gov.VersionProposal)
	if isVP {
		//another versionProposal in voting, exit.
		vp, err := govPlugin.findVotingVersionProposal(blockHash, blockNumber, state)
		if err != nil {
			log.Error("to find if there's a voting version proposal failed", "blockHash", blockHash)
			return err
		} else if vp != nil {
			log.Error("existing a voting version proposal.", "proposalID", vp.GetProposalID())
			return common.NewBizError("existing a version proposal at voting stage.")
		}
		//another VersionProposal in Pre-active process，exit
		proposalID, err := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)
		if err != nil {
			log.Error("to check if there's a pre-active version proposal failed.", "blockHash", blockHash)
			return err
		}
		if proposalID != common.ZeroHash {
			return common.NewBizError("existing a pre-active version proposal")
		}
	}

	//handle storage
	if err := govPlugin.govDB.SetProposal(proposal, state); err != nil {
		log.Error("save proposal failed", "proposalID", proposal.GetProposalID())
		return err
	}
	if err := govPlugin.govDB.AddVotingProposalID(blockHash, proposal.GetProposalID(), state); err != nil {
		log.Error("add proposal ID to voting proposal ID list failed", "proposalID", proposal.GetProposalID())
		return err
	}
	return nil
}

// vote for a proposal
func (govPlugin *GovPlugin) Vote(from common.Address, vote gov.Vote, blockHash common.Hash, blockNumber uint64, programVersion uint32, state xcom.StateDB) error {
	log.Debug("call Vote", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "programVersion", programVersion, "voteInfo", vote)
	if len(vote.ProposalID) == 0 || len(vote.VoteNodeID) == 0 || vote.VoteOption == 0 {
		return common.NewBizError("empty parameter detected.")
	}

	proposal, err := govPlugin.govDB.GetProposal(vote.ProposalID, state)
	if err != nil {
		log.Error("cannot find proposal by ID", "proposalID", vote.ProposalID)
		return err
	} else if proposal == nil {
		log.Error("incorrect proposal ID.", "proposalID", vote.ProposalID)
		return common.NewBizError("incorrect proposal ID.")
	}

	//check caller and voter
	if err := govPlugin.checkVerifier(from, proposal.GetProposer(), blockHash, proposal.GetSubmitBlock()); err != nil {
		return err
	}

	//voteOption range check
	if !(vote.VoteOption >= gov.Yes && vote.VoteOption <= gov.Abstention) {
		return common.NewBizError("vote option is error.")
	}

	if proposal.GetProposalType() == gov.Version {
		if vp, ok := proposal.(gov.VersionProposal); ok {
			if vp.GetNewVersion() != programVersion {
				log.Error("cannot vote for version proposal until node upgrade to a new version", "newVersion", vp.GetNewVersion(), "programVersion", programVersion)
				return common.NewBizError("node have not upgraded to a new version")
			}
		}
	}

	//check if vote.proposalID is in voting
	votingIDs, err := govPlugin.listVotingProposalID(blockHash, blockNumber, state)
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
	verifierList, err := govPlugin.govDB.ListVotedVerifier(vote.ProposalID, state)
	if err != nil {
		log.Error("list voted verifiers failed", "proposalID", vote.ProposalID)
		return err
	}

	if inNodeList(vote.VoteNodeID, verifierList) {
		log.Error("node has voted this proposal", "proposalID", vote.ProposalID, "nodeID", byteutil.PrintNodeID(vote.VoteNodeID))
		return common.NewBizError("node has voted this proposal.")
	}

	//handle storage
	if err := govPlugin.govDB.SetVote(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, state); err != nil {
		log.Error("save vote failed", "proposalID", vote.ProposalID)
		return err
	}

	//the proposal is version type, so add the node ID to active node list.
	if proposal.GetProposalType() == gov.Version {
		if err := govPlugin.govDB.AddActiveNode(blockHash, vote.ProposalID, vote.VoteNodeID); err != nil {
			log.Error("add nodeID to active node list failed", "proposalID", vote.ProposalID, "nodeID", byteutil.PrintNodeID(vote.VoteNodeID))
			return err
		}
	}

	return nil
}

// node declares it's version
func (govPlugin *GovPlugin) DeclareVersion(from common.Address, declaredNodeID discover.NodeID, declaredVersion uint32, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {
	log.Debug("call DeclareVersion", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion)
	//check caller is a Verifier or Candidate
	/*if err := govPlugin.checkVerifier(from, declaredNodeID, blockHash, blockNumber); err != nil {
		return err
	}*/
	if err := govPlugin.checkCandidate(from, declaredNodeID, blockHash, blockNumber); err != nil {
		return err
	}

	activeVersion := uint32(govPlugin.govDB.GetActiveVersion(state))
	if activeVersion <= 0 {
		return common.NewBizError("wrong active version.")
	}

	votingVP, err := govPlugin.findVotingVersionProposal(blockHash, blockNumber, state)
	if err != nil {
		log.Error("find if there's a voting version proposal failed", "blockHash", blockHash)
		return err
	}

	//there is a voting version proposal
	if votingVP != nil {
		nodeList, err := govPlugin.govDB.ListVotedVerifier(votingVP.ProposalID, state)
		if err != nil {
			log.Error("cannot list voted verifiers", "proposalID", votingVP.ProposalID)
			return err
		} else {
			if inNodeList(declaredNodeID, nodeList) && declaredVersion != votingVP.GetNewVersion() {
				log.Error("declared version should be same as proposal's version",
					"declaredNodeID", declaredNodeID, "declaredVersion", declaredVersion, "proposalID", votingVP.ProposalID, "newVersion", votingVP.GetNewVersion())
				return common.NewBizError("declared version should be same as proposal's version")
			}
		}
		if declaredVersion>>8 == activeVersion>>8 {
			//the declared version equals the current active version, notify staking immediately
			log.Debug("declared version equals active version.", "activeVersion", activeVersion, "declaredVersion", declaredVersion)
			stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion)
		} else if declaredVersion>>8 == votingVP.GetNewVersion()>>8 {
			//the declared version equals the new version, will notify staking when the proposal is passed
			log.Debug("declared version equals the new version.", "newVersion", votingVP.GetNewVersion, "declaredVersion", declaredVersion)
			govPlugin.govDB.AddActiveNode(blockHash, votingVP.ProposalID, declaredNodeID)
		} else {
			log.Error("declared version neither equals active version nor new version.", "activeVersion", activeVersion, "newVersion", votingVP.GetNewVersion, "declaredVersion", declaredVersion)
			return common.NewBizError("declared version neither equals active version nor new version.")
		}
	} else {
		preActiveVersion := govPlugin.govDB.GetPreActiveVersion(state)
		if declaredVersion>>8 == activeVersion>>8 || (preActiveVersion != 0 && declaredVersion == preActiveVersion) {
			//the declared version is the current active version, notify staking immediately
			stk.DeclarePromoteNotify(blockHash, blockNumber, declaredNodeID, declaredVersion)
		} else {
			log.Error("there's no version proposal at voting stage, declared version should be active or pre-active version.", "activeVersion", activeVersion, "declaredVersion", declaredVersion)
			return common.NewBizError("there's no version proposal at voting stage, declared version should be active version.")
		}
	}

	return nil
}

// client query a specified proposal
func (govPlugin *GovPlugin) GetProposal(proposalID common.Hash, state xcom.StateDB) (gov.Proposal, error) {
	log.Debug("call GetProposal", "proposalID", proposalID)

	proposal, err := govPlugin.govDB.GetProposal(proposalID, state)
	if err != nil {
		log.Error("get proposal by ID failed", "proposalID", proposalID, "msg", err.Error())
		return nil, err
	}
	if proposal == nil {
		return nil, common.NewBizError("incorrect proposal ID.")
	}
	return proposal, nil
}

// query a specified proposal's tally result
func (govPlugin *GovPlugin) GetTallyResult(proposalID common.Hash, state xcom.StateDB) (*gov.TallyResult, error) {
	tallyResult, err := govPlugin.govDB.GetTallyResult(proposalID, state)
	if err != nil {
		log.Error("get tallyResult by proposal ID failed.", "proposalID", proposalID, "msg", err.Error())
		return nil, err
	}
	if nil == tallyResult {
		return nil, common.NewBizError("get tallyResult by proposal ID failed.")
	}

	return tallyResult, nil
}

// query proposal list
func (govPlugin *GovPlugin) ListProposal(blockHash common.Hash, state xcom.StateDB) ([]gov.Proposal, error) {
	log.Debug("call ListProposal")
	var proposalIDs []common.Hash
	var proposals []gov.Proposal

	votingProposals, err := govPlugin.govDB.ListVotingProposal(blockHash, state)
	if err != nil {
		log.Error("list voting proposals failed.", "blockHash", blockHash)
		return nil, err
	}
	endProposals, err := govPlugin.govDB.ListEndProposalID(blockHash, state)
	if err != nil {
		log.Error("list end proposals failed.", "blockHash", blockHash)
		return nil, err
	}

	preActiveProposals, err := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)
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
		proposal, err := govPlugin.govDB.GetExistProposal(proposalID, state)
		if err != nil {
			log.Error("find proposal failed.", "proposalID", proposalID)
			return nil, err
		}
		proposals = append(proposals, proposal)
	}
	return proposals, nil
}

// tally a version proposal
func (govPlugin *GovPlugin) tallyForVersionProposal(proposal gov.VersionProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (bool, error) {
	proposalID := proposal.ProposalID
	log.Debug("call tallyForVersionProposal", "blockHash", blockHash, "blockNumber", blockNumber, "proposalID", proposal.ProposalID)

	verifiersCnt, err := govPlugin.govDB.AccuVerifiersLength(blockHash, proposalID)
	if err != nil {
		log.Error("count accumulated verifiers failed", "proposalID", proposalID, "blockHash", blockHash)
		return false, err
	}

	voteList, err := govPlugin.govDB.ListVoteValue(proposalID, state)
	if err != nil {
		log.Error("list voted values failed", "proposalID", proposalID)
		return false, err
	}

	voteCnt := uint16(len(voteList))
	yeas := voteCnt //`voteOption` can be ignored in version proposal, set voteCount to passCount as default.

	status := gov.Voting
	supportRate := float64(yeas) * 100 / float64(verifiersCnt)
	log.Debug("version proposal's supportRate", "supportRate", supportRate)

	if supportRate > xcom.SupportRateThreshold() {
		status = gov.PreActive

		activeList, err := govPlugin.govDB.GetActiveNodeList(blockHash, proposalID)
		if err != nil {
			log.Error("list active nodes failed", "blockHash", blockHash, "proposalID", proposalID)
			return false, err
		}
		if err := govPlugin.govDB.MoveVotingProposalIDToPreActive(blockHash, proposalID); err != nil {
			log.Error("move proposalID from voting proposalID list to pre-active list failed", "blockHash", blockHash, "proposalID", proposalID)
			return false, err
		}

		if err := govPlugin.govDB.SetPreActiveVersion(proposal.NewVersion, state); err != nil {
			log.Error("save pre-active version to state failed", "blockHash", blockHash, "proposalID", proposalID, "newVersion", proposal.NewVersion)
			return false, err
		}

		if err := stk.ProposalPassedNotify(blockHash, blockNumber, activeList, proposal.NewVersion); err != nil {
			log.Error("notify stating of the upgraded node list failed", "blockHash", blockHash, "proposalID", proposalID, "newVersion", proposal.NewVersion)
			return false, err
		}

	} else {
		status = gov.Failed
		if err := govPlugin.govDB.MoveVotingProposalIDToEnd(blockHash, proposalID, state); err != nil {
			log.Error("move proposalID from voting proposalID list to end list failed", "blockHash", blockHash, "proposalID", proposalID)
			return false, err
		}
	}

	tallyResult := &gov.TallyResult{
		ProposalID:    proposalID,
		Yeas:          yeas,
		Nays:          0x0,
		Abstentions:   0x0,
		AccuVerifiers: verifiersCnt,
		Status:        status,
	}

	if err := govPlugin.govDB.SetTallyResult(*tallyResult, state); err != nil {
		log.Error("save tally result failed", "tallyResult", tallyResult)
		return false, err
	}
	return status == gov.PreActive, nil
}

func (govPlugin *GovPlugin) tallyBasic(proposalID common.Hash, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (pass bool, err error) {
	log.Debug("call tallyBasic", "blockHash", blockHash, "blockNumber", blockNumber, "proposalID", proposalID)

	verifiersCnt, err := govPlugin.govDB.AccuVerifiersLength(blockHash, proposalID)
	if err != nil {
		log.Error("count accumulated verifiers failed", "proposalID", proposalID, "blockHash", blockHash)
		return false, err
	}

	status := gov.Voting
	yeas := uint16(0)
	nays := uint16(0)
	abstentions := uint16(0)

	voteList, err := govPlugin.govDB.ListVoteValue(proposalID, state)
	if err != nil {
		log.Error("list voted value failed.", "blockHash", blockHash)
		return false, err
	}
	for _, v := range voteList {
		if v.VoteOption == gov.Yes {
			yeas++
		}
		if v.VoteOption == gov.No {
			nays++
		}
		if v.VoteOption == gov.Abstention {
			abstentions++
		}
	}
	supportRate := float64(yeas) / float64(verifiersCnt)

	if supportRate >= xcom.SupportRateThreshold() {
		status = gov.Pass
	} else {
		status = gov.Failed
	}

	tallyResult := &gov.TallyResult{
		ProposalID:    proposalID,
		Yeas:          yeas,
		Nays:          nays,
		Abstentions:   abstentions,
		AccuVerifiers: verifiersCnt,
		Status:        status,
	}

	//govPlugin.govDB.MoveVotingProposalIDToEnd(blockHash, proposalID, state)
	if err := govPlugin.govDB.MoveVotingProposalIDToEnd(blockHash, proposalID, state); err != nil {
		log.Error("move proposalID from voting proposalID list to end list failed", "blockHash", blockHash, "proposalID", proposalID)
		return false, err
	}

	if err := govPlugin.govDB.SetTallyResult(*tallyResult, state); err != nil {
		log.Error("save tally result failed", "tallyResult", tallyResult)
		return false, err
	}
	return status == gov.Pass, nil
}

// check if the node a verifier, and the caller address is same as the staking address
func (govPlugin *GovPlugin) checkVerifier(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64) error {
	log.Debug("call checkVerifier", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
	verifierList, err := stk.GetVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if err != nil {
		log.Error("list verifiers failed", "blockHash", blockHash, "err", err)
		return err
	}

	for _, verifier := range verifierList {
		if verifier != nil && verifier.NodeId == nodeID {
			if verifier.StakingAddress == from {
				log.Debug("tx sender is a verifier.", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
				return nil
			} else {
				return common.NewBizError("tx sender should be node's staking address.")
			}
		}
	}
	log.Error("tx sender is not a verifier.", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
	return common.NewBizError("tx sender is not a verifier.")
}

// check if the node a candidate, and the caller address is same as the staking address
func (govPlugin *GovPlugin) checkCandidate(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64) error {
	log.Debug("call checkCandidate", "from", from, "blockHash", blockHash, "blockNumber", blockNumber, "nodeID", nodeID)
	candidateList, err := stk.GetCandidateList(blockHash)
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

// list all proposal IDs at voting stage
func (govPlugin *GovPlugin) listVotingProposalID(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) ([]common.Hash, error) {
	log.Debug("call checkCandidate", "blockHash", blockHash, "blockNumber", blockNumber)
	idList, err := govPlugin.govDB.ListVotingProposal(blockHash, state)
	if err != nil {
		log.Error("find voting version proposal failed", "blockHash", blockHash)
		return nil, err
	}
	return idList, nil
}

// find a version proposal at voting stage
func (govPlugin *GovPlugin) findVotingVersionProposal(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (*gov.VersionProposal, error) {
	log.Debug("call findVotingVersionProposal", "blockHash", blockHash, "blockNumber", blockNumber)
	idList, err := govPlugin.govDB.ListVotingProposal(blockHash, state)
	if err != nil {
		log.Error("find voting version proposal failed", "blockHash", blockHash)
		return nil, err
	}
	for _, proposalID := range idList {
		p, err := govPlugin.govDB.GetExistProposal(proposalID, state)
		if err != nil {
			return nil, err
		}
		if p.GetProposalType() == gov.Version {
			vp := p.(gov.VersionProposal)
			return &vp, nil
		}
	}
	return nil, nil
}

func (govPlugin *GovPlugin) SetParam(paraMap map[string]string, state xcom.StateDB) error {
	log.Debug("call SetParam", "paraMap", paraMap)
	return govPlugin.govDB.SetParam(paraMap, state)
}

func (govPlugin *GovPlugin) ListParam(state xcom.StateDB) (map[string]string, error) {
	log.Debug("call ListParam")
	paramList, err := govPlugin.govDB.ListParam(state)
	if err != nil {
		log.Error("list all parameters failed", "msg", err.Error())
		return nil, err
	}
	return paramList, nil
}

func (govPlugin *GovPlugin) GetParamValue(name string, state xcom.StateDB) (string, error) {
	log.Debug("call GetParamValue", "name", name)
	value, err := govPlugin.govDB.GetParam(name, state)
	if err != nil {
		log.Error("fina a parameter failed", "msg", err.Error())
		return "", err
	}
	return value, nil
}

func (govPlugin *GovPlugin) updateParam(proposal gov.ParamProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {
	log.Debug("call updateParam", "blockHash", blockHash, "blockNumber", blockNumber)
	if err := govPlugin.govDB.UpdateParam(proposal.ParamName, proposal.CurrentValue, proposal.NewValue, state); err != nil {
		log.Error("update parameter value failed", "msg", err.Error())
		return err
	}
	return nil
}

func inNodeList(proposer discover.NodeID, vList []discover.NodeID) bool {
	for _, v := range vList {
		if proposer == v {
			return true
		}
	}
	return false
}
