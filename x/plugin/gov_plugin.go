package plugin

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"sync"
)

var (
	govPluginOnce         sync.Once
	SupportRate_Threshold = 0.85
)

type GovPlugin struct {
	govDB *gov.GovDB
}

var govPlugin *GovPlugin

func GovPluginInstance() *GovPlugin {
	govPluginOnce.Do(func() {
		govPlugin = &GovPlugin{govDB: gov.GovDBInstance()}
	})
	return govPlugin
}

func (govPlugin *GovPlugin) Confirmed(block *types.Block) error {
	return nil
}

//implement BasePlugin
func (govPlugin *GovPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {

	if xutil.IsSettlementPeriod(header.Number.Uint64()) {
		verifierList, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())
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

//implement BasePlugin
func (govPlugin *GovPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	votingProposalIDs, err := govPlugin.govDB.ListVotingProposal(blockHash, state)
	if err != nil {
		log.Error("[GOV] EndBlock(): ListVerifierNodeID failed.")
		return err
	}

	for _, votingProposalID := range votingProposalIDs {
		votingProposal, err := govPlugin.govDB.GetExistProposal(votingProposalID, state)
		if nil != err {
			log.Error("[GOV] EndBlock(): get proposal by ID failed.", "proposalID", votingProposalID)
			return err
		}
		if votingProposal.GetEndVotingBlock() == header.Number.Uint64() {
			if xutil.IsSettlementPeriod(header.Number.Uint64()) {

				verifierList, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())
				if err != nil {
					log.Error("[GOV] EndBlock(): list verifier's node id by blockHash failed.", "blockHash", blockHash, "blockNumber", header.Number.Uint64())
					return err
				}

				if err := govPlugin.govDB.AccuVerifiers(blockHash, votingProposalID, verifierList); err != nil {
					log.Error("[GOV] EndBlock(): accu verifiers group by proposal ID failed.", "votingProposalID", votingProposalID, "blockHash", blockHash)
					return err
				}
			}

			accuVerifiersCnt, err := govPlugin.govDB.AccuVerifiersLength(blockHash, votingProposal.GetProposalID())
			if err != nil {
				log.Error("[GOV] EndBlock(): get accu verifiers length failed.", "blockHash", blockHash, "votingProposalID", votingProposalID)
				return err
			}

			verifierList, err := govPlugin.govDB.ListVotedVerifier(votingProposal.GetProposalID(), state)
			if err != nil {
				return err
			}

			if votingProposal.GetProposalType() == gov.Text {
				govPlugin.tallyForTextProposal(verifierList, accuVerifiersCnt, votingProposal.(gov.TextProposal), blockHash, state)
			} else if votingProposal.GetProposalType() == gov.Version {
				govPlugin.tallyForVersionProposal(verifierList, accuVerifiersCnt, votingProposal.(gov.VersionProposal), blockHash, header.Number.Uint64(), state)
			} else {
				err = errors.New("[GOV] EndBlock(): invalid Proposal Type")
				return err
			}
		}
	}
	preActiveProposalID, err := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)
	if err != nil {
		log.Error("[GOV] EndBlock(): to check if there's a preactive proposal failed.", "blockHash", blockHash)
		return err
	}
	if len(preActiveProposalID) <= 0 {
		return nil
	}
	//exsits a PreActiveProposal
	proposal, err := govPlugin.govDB.GetProposal(preActiveProposalID, state)
	if err != nil {
		return err
	}
	versionProposal, ok := proposal.(gov.VersionProposal)
	if ok {

		sub := header.Number.Uint64() - versionProposal.GetActiveBlock()

		if sub >= 0 && sub%xcom.ConsensusSize == 0 {
			validatorList, err := stk.ListCurrentValidatorID(blockHash, header.Number.Uint64())
			if err != nil {
				err := errors.New("[GOV] BeginBlock(): ListValidatorNodeID failed.")
				log.Error("[GOV] EndBlock(): list current round validators faied.", "blockHash", blockHash, "blockNumber", header.Number.Uint64())
				return err
			}
			var updatedNodes uint64 = 0

			//all active validators
			activeList, err := govPlugin.govDB.GetActiveNodeList(blockHash, preActiveProposalID)
			if err != nil {
				log.Error("[GOV] EndBlock(): list all active nodes failed.", "blockHash", blockHash, "preActiveProposalID", preActiveProposalID)
				return err
			}

			//check if all validators are active
			for _, val := range validatorList {
				if inNodeList(val, activeList) {
					if inNodeList(val, activeList) {
						updatedNodes++
					}
				}
			}
			if updatedNodes == xcom.ConsValidatorNum {
				tallyResult, err := govPlugin.govDB.GetTallyResult(preActiveProposalID, state)
				if err != nil {
					log.Error("[GOV] EndBlock(): find tally result by proposal ID failed.", "preActiveProposalID", preActiveProposalID)
					return err
				}
				//change tally status to "active"
				tallyResult.Status = gov.Active
				if err := govPlugin.govDB.SetTallyResult(*tallyResult, state); err != nil {
					log.Error("[GOV] EndBlock(): Update tally result failed.", "preActiveProposalID", preActiveProposalID)
					return err
				}
				err = govPlugin.govDB.MovePreActiveProposalIDToEnd(blockHash, preActiveProposalID, state)
				if err != nil {
					log.Error("[GOV] EndBlock(): MovePreActiveProposalIDToEnd() failed.", "blockHash", blockHash, "preActiveProposalID", preActiveProposalID)
					return err
				}

				err = govPlugin.govDB.ClearActiveNodes(blockHash, preActiveProposalID)
				if err != nil {
					log.Error("[GOV] EndBlock(): ClearActiveNodes() failed.", "blockHash", blockHash, "preActiveProposalID", preActiveProposalID)
					return err
				}
				err = govPlugin.govDB.SetActiveVersion(versionProposal.NewVersion, state)
				if err != nil {
					log.Error("[GOV] EndBlock(): SetActiveVersion() failed.", "blockHash", blockHash, "preActiveProposalID", preActiveProposalID)
					return err
				}
			}
		}
	}
	return nil
}

// nil is allowed
func (govPlugin *GovPlugin) GetPreActiveVersion(state xcom.StateDB) uint32 {
	return govPlugin.govDB.GetPreActiveVersion(state)
}

// should not be a nil value
func (govPlugin *GovPlugin) GetActiveVersion(state xcom.StateDB) uint32 {
	return govPlugin.govDB.GetActiveVersion(state)
}

// submit a proposal
func (govPlugin *GovPlugin) Submit(curBlockNum uint64, from common.Address, proposal gov.Proposal, blockHash common.Hash, state xcom.StateDB) error {

	hex := common.Bytes2Hex(blockHash.Bytes())
	log.Debug("check sender", "blockHash", hex, "blockNumber", curBlockNum )

	//param check
	if err := proposal.Verify(curBlockNum, state); err != nil {
		log.Error("verify proposal parameters failed", "from", from)
		return common.NewBizError(err.Error())
	}

	//check caller and proposer
	if !govPlugin.checkVerifier(from, proposal.GetProposer(), blockHash, curBlockNum) {
		return common.NewBizError("[GOV] Submit(): Tx sender is not a verifier.")
	}

	//handle version proposal
	_, isVP := proposal.(gov.VersionProposal)
	if isVP {
		//another versionProposal in voting, exit.
		vp, err := govPlugin.findVotingVersionProposal(blockHash, state)
		if err != nil {
			log.Error("[GOV] Submit(): to find if there's a voting version proposal failed", "blockHash", blockHash)
			return err
		} else if vp != nil {
			log.Error("[GOV] Submit(): existing a voting version proposal.", "votingProposalID", vp.GetProposalID())
			return err
		}
		//another VersionProposal in Pre-active process，exit
		proposalID, err := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)
		if err != nil {
			log.Error("[GOV] Submit(): to check if there's a pre-active version proposal failed.", "blockHash", blockHash)
			return err
		}
		if proposalID != common.ZeroHash {
			return common.NewBizError("existing a pre-active version proposal")
		}
	}

	//handle storage
	if err := govPlugin.govDB.SetProposal(proposal, state); err != nil {
		log.Error("[GOV] Submit(): save proposal failed", "proposalID", proposal.GetProposalID())
		return err
	}
	if err := govPlugin.govDB.AddVotingProposalID(blockHash, proposal.GetProposalID(), state); err != nil {
		log.Error("[GOV] Submit(): add proposal ID to voting proposal ID list failed", "proposalID", proposal.GetProposalID())
		return err
	}
	return nil
}

// vote for a proposal
func (govPlugin *GovPlugin) Vote(from common.Address, vote gov.Vote, blockHash common.Hash, curBlockNum uint64, state xcom.StateDB) error {
	if len(vote.ProposalID) == 0 || len(vote.VoteNodeID) == 0 || vote.VoteOption == 0 {
		return common.NewBizError("Empty parameter detected.")
	}

	proposal, err := govPlugin.govDB.GetProposal(vote.ProposalID, state)
	if err != nil {
		log.Error("[GOV] Vote(): cannot find proposal by ID", "proposalID", vote.ProposalID)
		return err
	} else if proposal == nil {
		log.Error("[GOV] Vote(): incorrect proposal ID.", "proposalID", vote.ProposalID)
		return common.NewBizError("Incorrect proposal ID.")
	}

	//check caller and voter
	if !govPlugin.checkVerifier(from, proposal.GetProposer(), blockHash, curBlockNum) {
		return common.NewBizError("The sender is not a verifier.")
	}

	//voteOption range check
	if !(vote.VoteOption >= gov.Yes && vote.VoteOption <= gov.Abstention) {
		return common.NewBizError("The vote option is error.")
	}

	//check if vote.proposalID is in voting
	vp, err := govPlugin.findVotingVersionProposal(blockHash, state)
	if err != nil {
		log.Error("[GOV] Submit(): to find if there's a voting version proposal failed", "blockHash", blockHash)
		return err
	} else if vp != nil {
		log.Error("[GOV] Submit(): existing a voting version proposal.", "votingProposalID", vp.GetProposalID())
		return err
	}
	if vp.GetProposalID() != vote.ProposalID {
		log.Error("[GOV] Vote(): proposal is not voting", "proposalID", vote.ProposalID)
		return common.NewBizError("Proposal is not voting.")
	}

	//handle storage
	if err := govPlugin.govDB.SetVote(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, state); err != nil {
		log.Error("[GOV] Vote(): save vote failed", "proposalID", vote.ProposalID)
		return err
	}
	if err := govPlugin.govDB.AddVotedVerifier(blockHash, vote.ProposalID, vote.VoteNodeID); err != nil {
		log.Error("[GOV] Vote(): Add nodeID to voted verifier list failed", "proposalID", vote.ProposalID, "voteNodeID", vote.VoteNodeID)
		return err
	}
	return nil
}

// node declares it's version
func (govPlugin *GovPlugin) DeclareVersion(from common.Address, declaredNodeID discover.NodeID, version uint32, blockHash common.Hash, curBlockNum uint64, state xcom.StateDB) error {

	//check caller is a Verifier or Candidate
	isVerifier := govPlugin.checkVerifier(from, declaredNodeID, blockHash, curBlockNum)
	isCandidate := govPlugin.checkCandidate(from, declaredNodeID, blockHash, curBlockNum)
	if !(isVerifier || isCandidate) {
		return common.NewBizError("The sender is not a verifier or candidate.")
	}

	activeVersion := uint32(govPlugin.govDB.GetActiveVersion(state))
	if activeVersion <= 0 {
		return common.NewBizError("wrong active version.")
	}

	votingVP, err := govPlugin.findVotingVersionProposal(blockHash, state)
	if err != nil {
		log.Error("[GOV] Vote(): to find if there's a voting version proposal failed", "blockHash", blockHash)
		return err
	}

	//there is a voting version proposal
	if votingVP != nil {
		if version>>8 == activeVersion>>8 {
			//the declared version is the current active version, notify staking immediately
			stk.DeclarePromoteNotify(blockHash, curBlockNum, declaredNodeID, version)
		} else if version>>8 == votingVP.GetNewVersion()>>8 {
			//the declared version is the next version, will notify staking when the proposal is passed
			govPlugin.govDB.AddActiveNode(blockHash, votingVP.ProposalID, declaredNodeID)
		} else {
			log.Error("[GOV] DeclareVersion(): declared version invalid.", "version", version)
			return common.NewBizError("declared version invalid.")
		}
	} else {
		if version>>8 == activeVersion>>8 {
			//the declared version is the current active version, notify staking immediately
			stk.DeclarePromoteNotify(blockHash, curBlockNum, declaredNodeID, version)
		} else {
			log.Error("[GOV] DeclareVersion(): declared version invalid.", "version", version)
			return common.NewBizError("declared version invalid.")
		}
	}

	return nil
}

// client query a specified proposal
func (govPlugin *GovPlugin) GetProposal(proposalID common.Hash, state xcom.StateDB) (gov.Proposal, error) {
	proposal, err := govPlugin.govDB.GetProposal(proposalID, state)
	if err != nil {
		log.Error("[GOV] GetProposal(): get proposal by ID failed", "proposalID", proposalID, "msg", err.Error())
		return nil, err
	}
	if proposal == nil {
		return nil, common.NewBizError("Incorrect proposal ID.")
	}
	return proposal, nil
}

// query a specified proposal's tally result
func (govPlugin *GovPlugin) GetTallyResult(proposalID common.Hash, state xcom.StateDB) (*gov.TallyResult, error) {
	tallyResult, err := govPlugin.govDB.GetTallyResult(proposalID, state)
	if err != nil {
		log.Error("[GOV] GetTallyResult(): Unable to get tallyResult.", "proposalID", proposalID, "msg", err.Error())
		return nil, err
	}
	if nil == tallyResult {
		return nil, common.NewBizError("Unable to get tallyResult.")
	}

	return tallyResult, nil
}

// query proposal list
func (govPlugin *GovPlugin) ListProposal(blockHash common.Hash, state xcom.StateDB) ([]gov.Proposal, error) {
	var proposalIDs []common.Hash
	var proposals []gov.Proposal

	votingProposals, err := govPlugin.govDB.ListVotingProposal(blockHash, state)
	if err != nil {
		log.Error("[GOV] ListProposal(): list voting proposals failed.", "blockHash", blockHash)
		return nil, err
	}
	endProposals, err := govPlugin.govDB.ListEndProposalID(blockHash, state)
	if err != nil {
		log.Error("[GOV] ListProposal(): list end proposals failed.", "blockHash", blockHash)
		return nil, err
	}

	preActiveProposals, err := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)
	if err != nil {
		log.Error("[GOV] ListProposal(): find pre-active proposal failed.", "blockHash", blockHash)
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
			log.Error("[GOV] ListProposal(): find proposal failed.", "proposalID", proposalID)
			return nil, err
		}
		proposals = append(proposals, proposal)
	}
	return proposals, nil
}

// tally for a text proposal
func (govPlugin *GovPlugin) tallyForTextProposal(votedVerifierList []discover.NodeID, accuCnt uint16, proposal gov.TextProposal, blockHash common.Hash, state xcom.StateDB) error {

	proposalID := proposal.ProposalID
	verifiersCnt, err := govPlugin.govDB.AccuVerifiersLength(blockHash, proposalID)
	if err != nil {
		log.Error("count accu verifiers failed", "proposalID", proposalID, "blockHash", blockHash)
		return err
	}

	status := gov.Voting
	yeas := uint16(0)
	nays := uint16(0)
	abstentions := uint16(0)

	voteList, err := govPlugin.govDB.ListVoteValue(proposal.ProposalID, state)
	if err != nil {
		log.Error("[GOV] tallyForTextProposal(): list vote value failed.", "blockHash", blockHash)
		return err
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
	supportRate := float64(yeas) * 100 / float64(accuCnt)

	if supportRate >= SupportRate_Threshold {
		status = gov.Pass
	} else {
		status = gov.Failed
	}

	tallyResult := &gov.TallyResult{
		ProposalID:    proposal.ProposalID,
		Yeas:          yeas,
		Nays:          nays,
		Abstentions:   abstentions,
		AccuVerifiers: verifiersCnt,
		Status:        status,
	}

	govPlugin.govDB.MoveVotingProposalIDToEnd(blockHash, proposal.ProposalID, state)

	if err := govPlugin.govDB.SetTallyResult(*tallyResult, state); err != nil {
		log.Error("Save tally result failed", "tallyResult", tallyResult)
		return err
	}
	return nil
}

// tally for a version proposal
func (govPlugin *GovPlugin) tallyForVersionProposal(votedVerifierList []discover.NodeID, accuCnt uint16, proposal gov.VersionProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {

	proposalID := proposal.ProposalID
	verifiersCnt, err := govPlugin.govDB.AccuVerifiersLength(blockHash, proposalID)
	if err != nil {
		log.Error("count accu verifiers failed", "proposalID", proposalID, "blockHash", blockHash)
		return err
	}

	voteList, err := govPlugin.govDB.ListVoteValue(proposalID, state)
	if err != nil {
		log.Error("list vote value failed", "proposalID", proposalID)
		return err
	}

	voteCnt := uint16(len(voteList))
	yeas := voteCnt //`voteOption` can be ignored in version proposal, set voteCount to passCount as default.

	status := gov.Voting
	supportRate := float64(yeas) * 100 / float64(verifiersCnt)
	if supportRate > SupportRate_Threshold {
		status = gov.PreActive

		activeList, err := govPlugin.govDB.GetActiveNodeList(blockHash, proposalID)
		if err != nil {
			log.Error("list active nodes failed", "blockHash", blockHash, "proposalID", proposalID)
			return err
		}
		govPlugin.govDB.MoveVotingProposalIDToPreActive(blockHash, proposalID, state)
		//todo: handle error
		stk.ProposalPassedNotify(blockHash, blockNumber, activeList, proposal.NewVersion)
	} else {
		status = gov.Failed
		govPlugin.govDB.MoveVotingProposalIDToEnd(blockHash, proposalID, state)
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
		log.Error("Save tally result failed", "tallyResult", tallyResult)
		return err
	}
	return nil
}

func (govPlugin *GovPlugin) TestTally(votedVerifierList []discover.NodeID, accuCnt uint16, proposal gov.Proposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {
	vp, ok := proposal.(gov.VersionProposal)
	if ok {
		return govPlugin.tallyForVersionProposal(votedVerifierList, accuCnt, vp, blockHash, blockNumber, state)
	}
	tp, ok := proposal.(gov.TextProposal)
	if ok {
		return govPlugin.tallyForTextProposal(votedVerifierList, accuCnt, tp, blockHash, state)
	}
	err := errors.New("[GOV] TestTally(): proposal type error")
	return err
}

// check if the node a verifier, and the caller address is same as the staking address
func (govPlugin *GovPlugin) checkVerifier(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64) bool {
	verifierList, err := stk.GetVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if err != nil {
		log.Error("list verifiers failed", "blockHash", blockHash)
		return false
	}

	for _, verifier := range verifierList {
		if verifier!= nil && verifier.NodeId == nodeID {
			if verifier.StakingAddress == from {
				return true
			} else {
				log.Warn("Verifier should send the tx by staking address")
				return false
			}
		}
	}
	return false
}

// check if the node a candidate, and the caller address is same as the staking address
func (govPlugin *GovPlugin) checkCandidate(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64) bool {
	candidateList, err := stk.GetCandidateList(blockHash, QueryStartNotIrr)
	if err != nil {
		log.Error("list candidates failed", "blockHash", blockHash)
		return false
	}

	for _, candidate := range candidateList {
		if candidate.NodeId == nodeID {
			if candidate.StakingAddress == from {
				return true
			} else {
				log.Warn("Candidate should send the tx by staking address")
				return false
			}
		}
	}
	return false
}

// find a version proposal at the voting stage
func (govPlugin *GovPlugin) findVotingVersionProposal(blockHash common.Hash, state xcom.StateDB) (*gov.VersionProposal, error) {
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
