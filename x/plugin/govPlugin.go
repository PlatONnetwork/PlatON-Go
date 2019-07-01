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
	"math/big"
	"sync"
)

var SupportRate_Threshold = 0.85

type GovPlugin struct {
	govDB *gov.GovDB
	once  sync.Once
}

var govPlugin *GovPlugin

func GovPluginInstance() *StakingPlugin {
	if nil == govPlugin {
		govPlugin = &GovPlugin{
			govDB: gov.NewGovDB(),
		}
	}
	return stk
}

func (govPlugin *GovPlugin) Confirmed(block *types.Block) error {
	return nil
}

//implement BasePlugin
func (govPlugin *GovPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {

	//TODO: Staking Plugin
	//if is end of Settle Cycle
	//TODO: if stk.IsEndofSettleCycle(state) {
	if true {
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
		if votingProposal.GetEndVotingBlock().Uint64() == header.Number.Uint64() {
			//TODO: if !stk.isEndofSettleCycle(blockHash) {
			if true {

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

			verifierList, err := govPlugin.govDB.ListVotedVerifier(blockHash, votingProposal.GetProposalID())
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
		//TODO: condition 2 should be: if header blockNum corresponds to is a new consensus round.
		if versionProposal.GetActiveBlock().Cmp(header.Number) >= 0 && true {
			//TODO: Fix, should be validator
			validatorList, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())
			if err != nil {
				err := errors.New("[GOV] BeginBlock(): ListValidatorNodeID failed.")
				return false, err
			}
			var updatedNodes uint8 = 0

			declareList := govPlugin.govDB.GetActiveNodeList(blockHash, preActiveProposalID)

			//voteList := govPlugin.govDB.ListVote(preActiveProposalID, state)
			var voterList []discover.NodeID
			//for _, vote := range voteList {
			//	var voter discover.NodeID //vote.voter
			//	voterList = append(voterList, voter)
			//}

			//check if validator has declared his version, or has voted for a version
			for _, val := range validatorList {
				if inNodeList(val, declareList) {
					if inNodeList(val, declareList) || inNodeList(val, voterList) {
						updatedNodes++
					}
				}
			}
			if updatedNodes == 25 {
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

				//stk.NotifyActive(blockHash, blockNumber, proposal.NewVersion)
			} else {
				//TODO inform staking of un-upgraded validators
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

func (govPlugin *GovPlugin) Submit(curBlockNum *big.Int, from common.Address, proposal gov.Proposal, blockHash common.Hash, state xcom.StateDB) (bool, error) {

	//TODO check if Address corresponds to NodeID
	/*if !stk.isAdressCorrespondingToNodeID(from, proposal.GetProposer()) {
		var err error = errors.New("[GOV] Submit(): tx sender is not the declare proposer.")
		return false, err
	}*/

	//param check
	success, err := proposal.Verify(curBlockNum, state)
	if !success || err != nil {
		err = errors.New("[GOV] Submit(): param error.")
		return false, err
	}

	//check if proposer is Verifier
	success, err = stk.IsCurrVerifier(blockHash, proposal.GetProposer())
	if !success || err != nil {
		err = errors.New("[GOV] Submit(): proposer is not verifier.")
		return false, err
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

func (govPlugin *GovPlugin) Vote(from common.Address, vote gov.Vote, blockHash common.Hash, state xcom.StateDB) (bool, error) {

	if len(vote.ProposalID) == 0 || len(vote.VoteNodeID) == 0 || vote.VoteOption == 0 {
		err := errors.New("[GOV] Vote(): empty parameter detected.")
		return false, err
	}
	//TODO check if Address corresponds to NodeID
	/*proposal, err := govPlugin.govDB.GetProposal(vote.ProposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Vote(): Unable to get proposal: %s", vote.ProposalID)
		err = errors.New(msg)
		return false, err
	}

	if !stk.isAdressCorrespondingToNodeID(from, proposal.GetProposer()) {
		var err error = errors.New("[GOV] Vote(): tx sender is not the declare proposer.")
		return false, err
	}*/

	success, err := stk.IsCurrVerifier(blockHash, vote.VoteNodeID)

	if !success || err != nil {
		err = errors.New("[GOV] Vote(): proposer is not verifier.")
		return false, err
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
	if !govPlugin.govDB.SetVote(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, state) {
		err = errors.New("[GOV] Vote(): Set vote failed.")
		return false, err
	}
	if !govPlugin.govDB.AddVotedVerifier(vote.ProposalID, vote.ProposalID, vote.VoteNodeID) {
		err = errors.New("[GOV] Vote(): Add VotedVerifier failed.")
		return false, err
	}
	return true, nil
}

//Verifier or candidate can declare his version
func (govPlugin *GovPlugin) DeclareVersion(from common.Address, declaredNodeID *discover.NodeID, version uint32, blockHash common.Hash, state xcom.StateDB) (bool, error) {

	////check address
	//if !plugin.StakingInstance(nil).isAdressCorrespondingToNodeID(from, declaredNodeID) {
	//	var err error = errors.New("[GOV] Vote(): tx sender is not the declare proposer.")
	//	return false, err
	//}

	//check voteNodeID: Verifier or Candidate
	isCurrVerifier, err := stk.IsCurrVerifier(blockHash, *declaredNodeID)
	if err != nil {
		var err = errors.New("[GOV] DeclareVersion(): run isCurrVerifier() failed.")
		return false, err
	}

	//TODO: Replace to stk.IsCandidate()
	var isCandidate bool

	if !(isCurrVerifier || isCandidate) {
		var err = errors.New("[GOV] DeclareVersion(): declared Node is not a verifier or candidate.")
		return false, err
	}

	activeVersion := uint32(govPlugin.govDB.GetActiveVersion(state))
	if activeVersion <= 0 {
		err := errors.New("[GOV] DeclareVersion(): get active version failed.")
		return false, err
	}

	if version>>8 == activeVersion>>8 {
		//TODO inform staking
	}

	//in voting process, and is a versionProposal
	votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
	for _, votingProposalID := range votingProposalIDs {
		votingProposal, err := govPlugin.govDB.GetProposal(votingProposalID, state)
		if err != nil {
			msg := fmt.Sprintf("[GOV] DeclareVersion(): Unable to get proposal: %s", votingProposalID)
			err = errors.New(msg)
			return false, err
		}
		if votingProposal == nil {
			continue
		}
		versionProposal, ok := votingProposal.(gov.VersionProposal)
		if !ok {
			continue
		}
		if versionProposal.GetNewVersion()>>8 == version>>8 {
			proposer := versionProposal.GetProposer()
			//store vote to ActiveNode
			if !govPlugin.govDB.AddActiveNode(blockHash, votingProposalID, proposer) {
				var err error = errors.New("[GOV] DeclareVersion(): add active node failed.")
				return false, err
			}
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
	//votedCnt := uint16(len(votedVerifierList))
	status := gov.Voting

	yeas := uint16(0)
	//nays := uint16(0)
	//abstentions := uint16(0)

	//voteList := govPlugin.govDB.ListVote(proposal.ProposalID, state)
	//for _, v := range voteList {
	//	//TODO
	//	if v.option == gov.Yes {
	//	yeas++
	//	}
	//	if v.option == gov.No {
	//	nays++
	//	}
	//	if v.option == gov.Abstention {
	//	abstentions++
	//	}
	//}
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
	voteCnt := uint16(len(govPlugin.govDB.ListVote(proposalID, state)))
	yeas := voteCnt //`voteOption` can be ignored in version proposal, set voteCount to passCount as default.

	status := gov.Voting
	supportRate := float64(yeas) * 100 / float64(verifiersCnt)
	//TODO: define SupportRateThreshold
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
