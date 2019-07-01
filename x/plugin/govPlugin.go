package plugin

import (
	"errors"
	"fmt"
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
	govPluginOnce  sync.Once
	SupportRate_Threshold = 0.85

)

type GovPlugin struct {
	govDB *gov.GovDB

}

var govPlugin *GovPlugin

func GovPluginInstance() *GovPlugin {
	govPluginOnce.Do(func() {
		govPlugin = &GovPlugin{govDB : gov.GovDBInstance()}
	})
	return govPlugin
}

func (govPlugin *GovPlugin) Confirmed(block *types.Block) error {
	return nil
}

//implement BasePlugin
func (govPlugin *GovPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {
	if xutil.IsSettlementPeriod(header.Number.Uint64()) {
		verifierList, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())

		if err != nil {
			err := errors.New("[GOV] BeginBlock(): ListVerifierNodeID failed.")
			return false, err
		}

		votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
		for _, votingProposalID := range votingProposalIDs {
			ok := govPlugin.govDB.AccuVerifiers(blockHash, votingProposalID, verifierList)
			if !ok {
				err := errors.New("[GOV] BeginBlock(): add Verifiers failed.")
				return false, err
			}
		}
	}
	return true, nil
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
func (govPlugin *GovPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {

	votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
	for _, votingProposalID := range votingProposalIDs {
		votingProposal, err := govPlugin.govDB.GetProposal(votingProposalID, state)
		if nil != err {
			msg := fmt.Sprintf("[GOV] EndBlock(): Unable to get proposal: %s", votingProposalID)
			err = errors.New(msg)
			return false, err
		}
		if votingProposal.GetEndVotingBlock() == header.Number.Uint64() {
			if xutil.IsSettlementPeriod(header.Number.Uint64()) {

				verifierList, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())
				if err != nil {
					err := errors.New("[GOV] BeginBlock(): ListVerifierNodeID failed.")
					return false, err
				}

				ok := govPlugin.govDB.AccuVerifiers(blockHash, votingProposalID, verifierList)
				if !ok {
					err = errors.New("[GOV] EndBlock(): add Verifiers failed.")
					return false, err
				}
			}

			accuVerifiersCnt := uint16(govPlugin.govDB.AccuVerifiersLength(blockHash, votingProposal.GetProposalID()))

			verifierList, err := govPlugin.govDB.ListVotedVerifier(votingProposal.GetProposalID(), state)
			if err != nil {
				return false, err
			}

			if votingProposal.GetProposalType() == gov.Text {
				govPlugin.tallyForTextProposal(verifierList, accuVerifiersCnt, votingProposal.(gov.TextProposal), blockHash, state)
			} else if votingProposal.GetProposalType() == gov.Version {
				govPlugin.tallyForVersionProposal(verifierList, accuVerifiersCnt, votingProposal.(gov.VersionProposal), blockHash, header.Number.Uint64(), state)
			} else {

			}
		}
	}
	preActiveProposalID := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)
	if len(preActiveProposalID) <= 0 {
		return true, nil
	}
	//exsits a PreActiveProposal
	proposal, err := govPlugin.govDB.GetProposal(preActiveProposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] EndBlock(): Unable to get proposal: %s", preActiveProposalID)
		err = errors.New(msg)
		return false, err
	}
	versionProposal, ok := proposal.(gov.VersionProposal)
	if ok {

		sub := header.Number.Uint64() - versionProposal.GetActiveBlock()

		if sub >= 0 && sub % xcom.ConsensusSize == 0 {
			validatorList, err := stk.ListCurrentValidatorID(blockHash, header.Number.Uint64())
			if err != nil {
				err := errors.New("[GOV] BeginBlock(): ListValidatorNodeID failed.")
				return false, err
			}
			var updatedNodes uint64 = 0

			//all active validators
			activeList := govPlugin.govDB.GetActiveNodeList(blockHash, preActiveProposalID)

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
					err = errors.New("[GOV] EndBlock(): get tallyResult failed.")
					return false, err
				}
				//change tally status to "active"
				tallyResult.Status = gov.Active
				if !govPlugin.govDB.SetTallyResult(*tallyResult, state) {
					err = errors.New("[GOV] EndBlock(): Unable to set tallyResult.")
					return false, err
				}
				govPlugin.govDB.MovePreActiveProposalIDToEnd(blockHash, preActiveProposalID, state)

				govPlugin.govDB.ClearActiveNodes(blockHash, preActiveProposalID)

				//todo:
				//stk.NotifyActive(blockHash, blockNumber, proposal.NewVersion)
			}
		}
	}
	return true, nil
}

//nil is allowed
func (govPlugin *GovPlugin) GetPreActiveVersion(state xcom.StateDB) uint32 {
	return govPlugin.govDB.GetPreActiveVersion(state)
}

//should not be a nil value
func (govPlugin *GovPlugin) GetActiveVersion(state xcom.StateDB) uint32 {
	return govPlugin.govDB.GetActiveVersion(state)
}

func (govPlugin *GovPlugin) Submit(curBlockNum uint64, from common.Address, proposal gov.Proposal, blockHash common.Hash, state xcom.StateDB) (bool, error) {

	//param check
	success, err := proposal.Verify(curBlockNum, state)
	if !success || err != nil {
		err = errors.New("[GOV] Submit(): param error.")
		return false, err
	}

	//check caller and proposer
	if !govPlugin.checkVerifier(from, proposal.GetProposer(), blockHash, curBlockNum) {
		return false, errors.New("[GOV] Submit(): Tx sender is not a verifier.")
	}

	//handle version proposal
	_, ok := proposal.(gov.VersionProposal)
	if ok {
		//another versionProposal in voting, exit.
		votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
		for _, votingProposalID := range votingProposalIDs {
			votingProposal, err := govPlugin.govDB.GetProposal(votingProposalID, state)
			if err != nil {
				msg := fmt.Sprintf("[GOV] Submit(): Unable to get proposal: %s", votingProposalID)
				err = errors.New(msg)
				return false, err
			}
			_, ok := votingProposal.(gov.VersionProposal)
			if ok {
				var err error = errors.New("[GOV] Submit(): existing a voting version proposal.")
				return false, err
			}
		}
		//another VersionProposal in Pre-active process，exit
		if len(govPlugin.govDB.GetPreActiveProposalID(blockHash, state)) > 0 {
			err = errors.New("[GOV] Submit(): existing a pre-active version proposal.")
			return false, err
		}
	}

	//handle storage
	err = govPlugin.govDB.SetProposal(proposal, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Submit(): Unable to set proposal: %s", proposal.GetProposalID())
		err = errors.New(msg)
		return false, err
	}
	//if !ok {
	//	err = errors.New("[GOV] Submit(): set proposal failed.")
	//	return false, err
	//}
	ok = govPlugin.govDB.AddVotingProposalID(blockHash, proposal.GetProposalID(), state)
	if !ok {
		err = errors.New("[GOV] Submit(): add VotingProposalID failed.")
		return false, err
	}
	return true, nil
}

func (govPlugin *GovPlugin) Vote(from common.Address, vote gov.Vote, blockHash common.Hash, curBlockNum uint64, state xcom.StateDB) (bool, error) {

	if len(vote.ProposalID) == 0 || len(vote.VoteNodeID) == 0 || vote.VoteOption == 0 {
		return false, errors.New("[GOV] Vote(): empty parameter detected.")
	}

	proposal, err := govPlugin.govDB.GetProposal(vote.ProposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Vote(): Unable to get proposal: %s", vote.ProposalID)
		err = errors.New(msg)
		return false, err
	}

	//check caller and voter
	if !govPlugin.checkVerifier(from, proposal.GetProposer(), blockHash, curBlockNum) {
		return false, errors.New("[GOV] Vote(): Tx sender is not a verifier.")
	}

	//voteOption range check
	if vote.VoteOption <= gov.Yes && vote.VoteOption >= gov.Abstention {
		err = errors.New("[GOV] Vote(): VoteOption invalid.")
		return false, err
	}

	//check if vote.proposalID is in voting
	isVoting := func(proposalID common.Hash, votingProposalList []common.Hash) bool {
		for _, votingProposal := range votingProposalList {
			if proposalID == votingProposal {
				return true
			}
		}
		return false
	}
	votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
	if !isVoting(vote.ProposalID, votingProposalIDs) {
		err = errors.New("[GOV] Vote(): vote.proposalID is not in voting.")
		return false, err
	}

	//handle storage
	err = govPlugin.govDB.SetVote(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, state)
	if err != nil {
		return false, err
	}
	if !govPlugin.govDB.AddVotedVerifier(vote.ProposalID, vote.ProposalID, vote.VoteNodeID) {
		err = errors.New("[GOV] Vote(): Add VotedVerifier failed.")
		return false, err
	}
	return true, nil
}

//Verifier or candidate can declare his version
func (govPlugin *GovPlugin) DeclareVersion(from common.Address, declaredNodeID discover.NodeID, version uint32, blockHash common.Hash, curBlockNum uint64, state xcom.StateDB) (bool, error) {

	//check caller is a Verifier or Candidate
	if !govPlugin.checkVerifier(from, declaredNodeID, blockHash, curBlockNum) {
		return false, errors.New("[GOV] Vote(): Tx sender is not a verifier.")
	}

	if !govPlugin.checkCandidate(from, declaredNodeID, blockHash, curBlockNum) {
		var err = errors.New("[GOV] DeclareVersion(): Tx sender is not a candidate.")
		return false, err
	}

	activeVersion := uint32(govPlugin.govDB.GetActiveVersion(state))
	if activeVersion <= 0 {
		err := errors.New("[GOV] DeclareVersion(): get active version failed.")
		return false, err
	}

	vp, err := govPlugin.findVotingVersionProposal(blockHash, state)
	if err != nil {
		return false, err
	}

	//there is a voting version proposal
	if vp != nil {
		if version>>8 == activeVersion>>8 && version>>8 == vp.GetNewVersion()>>8 {
			govPlugin.govDB.AddActiveNode(blockHash, vp.ProposalID, declaredNodeID)
		}else{
			return false, errors.New(fmt.Sprintf("[GOV] DeclareVersion(): invalid declared version: %s", version))
		}
	}else {
		if version>>8 == activeVersion>>8 {
			//TODO inform staking
		}else{
			return false, errors.New(fmt.Sprintf("[GOV] DeclareVersion(): invalid declared version: %s", version))
		}
	}

	return true, nil
}

func (govPlugin *GovPlugin) GetProposal(proposalID common.Hash, state xcom.StateDB) gov.Proposal {
	proposal, err := govPlugin.govDB.GetProposal(proposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Submit(): Unable to set proposal: %s", proposalID)
		err = errors.New(msg)
		return nil
	}
	if proposal != nil {
		return proposal
	}
	log.Error("[GOV] GetProposal(): Unable to get proposal.")
	return nil
}

func (govPlugin *GovPlugin) GetTallyResult(proposalID common.Hash, state xcom.StateDB) *gov.TallyResult {
	tallyResult, err := govPlugin.govDB.GetTallyResult(proposalID, state)
	if err != nil {
		log.Error("[GOV] GetTallyResult(): Unable to get tallyResult from db.")
	}
	if nil != tallyResult {
		return tallyResult
	}
	log.Error("[GOV] GetTallyResult(): Unable to get tallyResult.")
	return nil
}

func (govPlugin *GovPlugin) ListProposal(blockHash common.Hash, state xcom.StateDB) []gov.Proposal {
	var proposalIDs []common.Hash
	var proposals []gov.Proposal

	votingProposals := govPlugin.govDB.ListVotingProposal(blockHash, state)
	endProposals := govPlugin.govDB.ListEndProposalID(blockHash, state)
	preActiveProposals := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)

	proposalIDs = append(proposalIDs, votingProposals...)
	proposalIDs = append(proposalIDs, endProposals...)
	proposalIDs = append(proposalIDs, preActiveProposals)

	for _, proposalID := range proposalIDs {
		proposal, err := govPlugin.govDB.GetProposal(proposalID, state)
		if err != nil {
			msg := fmt.Sprintf("[GOV] ListProposal(): Unable to get proposal: %s", proposalID)
			err = errors.New(msg)
			return nil
		}
		proposals = append(proposals, proposal)
	}
	return proposals
}

func (govPlugin *GovPlugin) tallyForTextProposal(votedVerifierList []discover.NodeID, accuCnt uint16, proposal gov.TextProposal, blockHash common.Hash, state xcom.StateDB) error {

	status := gov.Voting
	yeas := uint16(0)

	supportRate := float64(yeas) / float64(accuCnt)

	if supportRate >= SupportRate_Threshold {
		status = gov.Pass
	} else {
		status = gov.Failed
	}

	tallyResult := gov.TallyResult{
		ProposalID:    proposal.ProposalID,
		Yeas:          yeas,
		AccuVerifiers: accuCnt,
		Status:        status,
	}

	govPlugin.govDB.MoveVotingProposalIDToEnd(blockHash, proposal.ProposalID, state)

	if govPlugin.govDB.SetTallyResult(tallyResult, state) {
		return nil
	} else {
		log.Error("[GOV] tally(): Unable to save text proposal tally result.")
		return errors.New("save text proposal tally result error")
	}
}

func (govPlugin *GovPlugin) tallyForVersionProposal(votedVerifierList []discover.NodeID, accuCnt uint16, proposal gov.VersionProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {

	proposalID := proposal.ProposalID
	verifiersCnt := uint16(govPlugin.govDB.AccuVerifiersLength(blockHash, proposalID))


	voterList, err := govPlugin.govDB.ListVoteValue(proposalID, state)
	if err!= nil {
		return err
	}

	voteCnt := uint16(len(voterList))
	yeas := voteCnt //`voteOption` can be ignored in version proposal, set voteCount to passCount as default.

	status := gov.Voting
	supportRate := float64(yeas) * 100 / float64(verifiersCnt)
	if supportRate > SupportRate_Threshold {
		status = gov.Pass
	} else {
		status = gov.Failed
		govPlugin.govDB.MoveVotingProposalIDToEnd(blockHash, proposal.ProposalID, state)
	}

	tallyResult := &gov.TallyResult{
		ProposalID:    proposalID,
		Yeas:          yeas,
		Nays:          0x0,
		Abstentions:   0x0,
		AccuVerifiers: verifiersCnt,
		Status:        status,
	}

	if govPlugin.govDB.SetTallyResult(*tallyResult, state) {
		return nil
	} else {
		log.Error("[GOV] tally(): Unable to save version proposal tally result.")
		return errors.New("save version proposal tally result error")
	}
}


func (govPlugin *GovPlugin) checkVerifier(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64) bool {
	stk.GetVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	return true
}

func (govPlugin *GovPlugin) checkCandidate(from common.Address, nodeID discover.NodeID, blockHash common.Hash, blockNumber uint64) bool {
	stk.GetCandidateList(blockHash, QueryStartNotIrr)
	return true
}

func  (govPlugin *GovPlugin) findVotingVersionProposal(blockHash common.Hash, state xcom.StateDB) (*gov.VersionProposal, error) {
	idList := govPlugin.govDB.ListVotingProposal(blockHash, state)
	for _, proposalID := range idList {
		p, err := govPlugin.govDB.GetProposal(proposalID, state)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, errors.New(fmt.Sprintf("cannot find specified proposal: %s", proposalID))
		}
		if p.GetProposalType() == gov.Version {
			vp := p.(gov.VersionProposal)
			return &vp, nil
		}
	}
	return nil, nil
}

