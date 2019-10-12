package plugin

import (
	"errors"
	"math"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

var (
	govPluginOnce sync.Once
)

type GovPlugin struct {
}

var govp *GovPlugin

func GovPluginInstance() *GovPlugin {
	govPluginOnce.Do(func() {
		log.Info("Init Governance plugin ...")
		govp = &GovPlugin{}
	})
	return govp
}

func (govPlugin *GovPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

//implement BasePlugin
func (govPlugin *GovPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	var blockNumber = header.Number.Uint64()
	log.Debug("call BeginBlock()", "blockNumber", blockNumber, "blockHash", blockHash)

	if !xutil.IsBeginOfConsensus(blockNumber) {
		return nil
	}

	if xutil.IsBeginOfEpoch(blockNumber) {
		if err := accuVerifiersAtBeginOfSettlement(blockHash, blockNumber); err != nil {
			log.Error("accumulates all distinct verifiers for voting proposal failed.", "err", err)
			return err
		}
	}

	//check if there's a pre-active version proposal that can be activated
	preActiveVersionProposalID, err := gov.GetPreActiveProposalID(blockHash)
	if err != nil {
		log.Error("check if there's a pre-active version proposal failed.", "blockNumber", blockNumber, "blockHash", blockHash)
		return err
	}
	if preActiveVersionProposalID == common.ZeroHash {
		return nil
	}

	//handle a PreActiveProposal
	preActiveVersionProposal, err := gov.GetExistProposal(preActiveVersionProposalID, state)
	if err != nil {
		return err
	}
	versionProposal, isVersionProposal := preActiveVersionProposal.(*gov.VersionProposal)

	if isVersionProposal {
		log.Debug("found pre-active version proposal", "proposalID", preActiveVersionProposalID, "blockNumber", blockNumber, "blockHash", blockHash, "activeBlockNumber", versionProposal.GetActiveBlock())
		if blockNumber >= versionProposal.GetActiveBlock() && (blockNumber-versionProposal.GetActiveBlock())%xutil.ConsensusSize() == 0 {
			currentValidatorList, err := stk.ListCurrentValidatorID(blockHash, blockNumber)
			if err != nil {
				log.Error("list current round validators failed.", "blockHash", blockHash, "blockNumber", blockNumber)
				return err
			}
			var updatedNodes int = 0
			var totalValidators int = len(currentValidatorList)

			//all active validators (including node that has either voted or declared)
			activeList, err := gov.GetActiveNodeList(blockHash, preActiveVersionProposalID)
			if err != nil {
				log.Error("list all active nodes failed.", "blockNumber", blockNumber, "blockHash", blockHash, "preActiveVersionProposalID", preActiveVersionProposalID)
				return err
			}

			activeMap := make(map[discover.NodeID]struct{}, len(activeList))
			for _, activeNode := range activeList {
				activeMap[activeNode] = struct{}{}
			}

			//check if all validators are active
			for _, validator := range currentValidatorList {
				if _, isActive := activeMap[validator]; isActive {
					updatedNodes++
				}
			}

			log.Debug("check active criteria", "blockNumber", blockNumber, "blockHash", blockHash, "pre-active nodes", updatedNodes, "total validators", totalValidators, "activeList", activeList, "currentValidator", currentValidatorList)
			if updatedNodes == totalValidators {
				log.Debug("the pre-active version proposal has passed")
				tallyResult, err := gov.GetTallyResult(preActiveVersionProposalID, state)
				if err != nil || tallyResult == nil {
					log.Error("find pre-active version proposal tally result failed.", "blockNumber", blockNumber, "blockHash", blockHash, "preActiveVersionProposalID", preActiveVersionProposalID)
					return err
				}
				//change tally status to "active"
				tallyResult.Status = gov.Active

				if err := gov.SetTallyResult(*tallyResult, state); err != nil {
					log.Error("update version proposal tally result failed.", "preActiveVersionProposalID", preActiveVersionProposalID)
					return err
				}

				if versionProposal.GetActiveBlock() != blockNumber {
					versionProposal.ActiveBlock = blockNumber
					if err := gov.SetProposal(versionProposal, state); err != nil {
						log.Error("update activeBlock of version proposal failed.", "preActiveVersionProposalID", preActiveVersionProposalID, "blockNumber", blockNumber, "blockHash", blockHash)
					}
				}

				if err = gov.MovePreActiveProposalIDToEnd(blockHash, preActiveVersionProposalID); err != nil {
					log.Error("move version proposal ID to EndProposalID list failed.", "blockNumber", blockNumber, "blockHash", blockHash, "preActiveVersionProposalID", preActiveVersionProposalID)
					return err
				}

				if err = gov.ClearActiveNodes(blockHash, preActiveVersionProposalID); err != nil {
					log.Error("clear version proposal active nodes failed.", "blockNumber", blockNumber, "blockHash", blockHash, "preActiveVersionProposalID", preActiveVersionProposalID)
					return err
				}

				if err = gov.AddActiveVersion(versionProposal.NewVersion, blockNumber, state); err != nil {
					log.Error("save active version to stateDB failed.", "blockNumber", blockNumber, "blockHash", blockHash, "preActiveProposalID", preActiveVersionProposalID)
					return err
				}
				log.Debug("version proposal active.", "proposalID", versionProposal.ProposalID, "newVersion", versionProposal.NewVersion, "newVersionString", xutil.ProgramVersion2Str(versionProposal.NewVersion))
			}
		}
	}
	return nil
}

//implement BasePlugin

func (govPlugin *GovPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	var blockNumber = header.Number.Uint64()
	log.Debug("call EndBlock()", "blockNumber", blockNumber, "blockHash", blockHash)

	//the endVotingBlock must be consensus Election block
	if !xutil.IsElection(blockNumber) {
		return nil
	}

	votingProposalIDs, err := gov.ListVotingProposal(blockHash)
	if err != nil {
		return err
	}
	if len(votingProposalIDs) == 0 {
		log.Debug("there's no voting proposal", "blockNumber", blockNumber, "blockHash", blockHash)
		return nil
	}

	//iterate each voting proposal, to check if current block is proposal's end-voting block.
	for _, votingProposalID := range votingProposalIDs {
		log.Debug("iterate each voting proposal", "proposalID", votingProposalID)
		votingProposal, err := gov.GetExistProposal(votingProposalID, state)
		log.Debug("find voting proposal", "votingProposal", votingProposal)
		if nil != err {
			return err
		}
		if votingProposal.GetEndVotingBlock() == blockNumber {
			log.Debug("current block is end-voting block", "proposalID", votingProposal.GetProposalID(), "blockNumber", blockNumber)
			//tally the results
			if votingProposal.GetProposalType() == gov.Text {
				_, err := tallyText(votingProposal.GetProposalID(), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
			} else if votingProposal.GetProposalType() == gov.Version {
				err = tallyVersion(votingProposal.(*gov.VersionProposal), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
			} else if votingProposal.GetProposalType() == gov.Cancel {
				_, err := tallyCancel(votingProposal.(*gov.CancelProposal), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
			} else {
				log.Error("invalid proposal type", "type", votingProposal.GetProposalType())
				err = errors.New("invalid proposal type")
				return err
			}
		}
	}
	return nil
}

// According to the proposal's rules, the submit block maybe is the begin block of a settlement, even then, it's ok, gov.AccuVerifiers will remove the duplicated verifiers.
func accuVerifiersAtBeginOfSettlement(blockHash common.Hash, blockNumber uint64) error {
	votingProposalIDs, err := gov.ListVotingProposal(blockHash)
	if err != nil {
		return err
	}
	if len(votingProposalIDs) == 0 {
		log.Debug("there's no voting proposal", "blockNumber", blockNumber, "blockHash", blockHash)
		return nil
	}

	verifierList, err := stk.ListVerifierNodeID(blockHash, blockNumber)
	if err != nil {
		return err
	}
	log.Debug("get verifier nodes from staking", "verifierCount", len(verifierList))

	//note: if the proposal's submit block == blockNumber, it's ok, gov.AccuVerifiers will remove the duplicated verifiers
	for _, votingProposalID := range votingProposalIDs {
		if err := gov.AccuVerifiers(blockHash, votingProposalID, verifierList); err != nil {
			return err
		}
	}
	return nil
}

// tally a version proposal
func tallyVersion(proposal *gov.VersionProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {
	proposalID := proposal.ProposalID
	log.Debug("proposal tally", "proposalID", proposal.ProposalID, "blockHash", blockHash, "blockNumber", blockNumber)

	verifierList, err := gov.ListAccuVerifier(blockHash, proposalID)
	if err != nil {
		return err
	}
	verifiersCnt := uint16(len(verifierList))

	voteList, err := gov.ListVoteValue(proposalID, state)
	if err != nil {
		return err
	}

	voteCnt := uint16(len(voteList))
	yeas := voteCnt //`voteOption` can be ignored in version proposal, set voteCount to passCount as default.

	status := gov.Failed
	supportRate := float64(yeas) / float64(verifiersCnt)

	log.Debug("version proposal", "supportRate", supportRate, "required", Decimal(xcom.VersionProposal_SupportRate()))

	if Decimal(supportRate) >= Decimal(xcom.VersionProposal_SupportRate()) {
		status = gov.PreActive

		if err := gov.MoveVotingProposalIDToPreActive(blockHash, proposalID); err != nil {
			log.Error("move version proposal ID to pre-active failed", "proposalID", proposalID, "blockNumber", blockNumber, "blockHash", blockHash)
			return err
		}

		if err := gov.SetPreActiveVersion(proposal.NewVersion, state); err != nil {
			log.Error("save pre-active version to state failed", "proposalID", proposalID, "blockHash", blockHash, "newVersion", proposal.NewVersion, "newVersionString", xutil.ProgramVersion2Str(proposal.NewVersion))
			return err
		}

		activeList, err := gov.GetActiveNodeList(blockHash, proposalID)
		if err != nil {
			log.Error("list active node failed", "proposalID", proposalID, "blockNumber", blockNumber, "blockHash", blockHash)
			return err
		}
		log.Debug("call stk.ProposalPassedNotify", "proposalID", proposalID, "activeList", activeList)
		if err := stk.ProposalPassedNotify(blockHash, blockNumber, activeList, proposal.NewVersion); err != nil {
			log.Error("call stk.ProposalPassedNotify failed", "proposalID", proposalID, "blockHash", blockHash, "newVersion", proposal.NewVersion, "activeList", activeList)
			return err
		}

	} else {
		if err := gov.MoveVotingProposalIDToEnd(blockHash, proposalID); err != nil {
			log.Error("move proposalID from voting proposalID list to end list failed", "proposalID", proposalID, "blockHash", blockHash)
			return err
		}
		if err := gov.ClearActiveNodes(blockHash, proposalID); err != nil {
			return err
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

	if err := gov.SetTallyResult(*tallyResult, state); err != nil {
		log.Error("save tally result failed", "proposalID", proposalID, "tallyResult", tallyResult)
		return err
	}

	log.Debug("proposal tally result", "proposalID", proposalID, "tallyResult", tallyResult, "verifierList", verifierList)
	return nil
}

func tallyText(proposalID common.Hash, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (pass bool, err error) {
	return tally(gov.Text, proposalID, blockHash, blockNumber, state)
}

func tallyCancel(cp *gov.CancelProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (pass bool, err error) {
	if pass, err := tally(gov.Cancel, cp.ProposalID, blockHash, blockNumber, state); err != nil {
		return false, err
	} else if pass {
		if proposal, err := gov.GetExistProposal(cp.TobeCanceled, state); err != nil {
			return false, err
		} else if proposal.GetProposalType() != gov.Version {
			return false, gov.TobeCanceledProposalTypeError
		}
		if votingProposalIDList, err := gov.ListVotingProposalID(blockHash); err != nil {
			return false, err
		} else if !xutil.InHashList(cp.TobeCanceled, votingProposalIDList) {
			return false, gov.TobeCanceledProposalNotAtVoting
		}

		if tallyResult, err := gov.GetTallyResult(cp.TobeCanceled, state); err != nil {
			return false, err
		} else {
			if tallyResult == nil {
				tallyResult = &gov.TallyResult{
					ProposalID:    cp.TobeCanceled,
					Yeas:          0,
					Nays:          0,
					Abstentions:   0,
					AccuVerifiers: 0,
				}
			} else if tallyResult.Status != gov.Voting {
				log.Error("the version proposal to be canceled is not at voting stage, but the cancel proposal is passed")
				return false, err
			}
			verifierList, err := gov.ListAccuVerifier(blockHash, cp.TobeCanceled)
			if err != nil {
				return false, err
			}
			verifiersCnt := uint16(len(verifierList))

			voteList, err := gov.ListVoteValue(cp.TobeCanceled, state)
			if err != nil {
				return false, err
			}

			voteCnt := uint16(len(voteList))
			yeas := voteCnt

			tallyResult.Yeas = yeas
			tallyResult.AccuVerifiers = verifiersCnt
			tallyResult.Status = gov.Canceled
			tallyResult.CanceledBy = cp.ProposalID

			if err := gov.SetTallyResult(*tallyResult, state); err != nil {
				log.Error("version proposal is canceled failed, cannot save its tally result", "tallyResult", tallyResult)
				return false, err
			}

			if err := gov.ClearActiveNodes(blockHash, cp.TobeCanceled); err != nil {
				return false, err
			}

			if err := gov.MoveVotingProposalIDToEnd(blockHash, cp.TobeCanceled); err != nil {
				return false, err
			}

			log.Debug("version proposal is canceled success by others", "proposalID", cp.TobeCanceled, "tallyResult", tallyResult)
		}
	}
	return true, nil
}

func tally(proposalType gov.ProposalType, proposalID common.Hash, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (pass bool, err error) {
	log.Debug("proposal tally", "proposalID", proposalID, "blockHash", blockHash, "blockNumber", blockNumber, "proposalID", proposalID)

	verifierList, err := gov.ListAccuVerifier(blockHash, proposalID)
	if err != nil {
		return false, err
	}

	verifiersCnt := uint16(len(verifierList))

	status := gov.Voting

	yeas, nays, abstentions, err := gov.TallyVoteValue(proposalID, state)
	if err != nil {
		return false, err
	}

	voteRate := Decimal(float64(yeas+nays+abstentions) / float64(verifiersCnt))
	supportRate := Decimal(float64(yeas) / float64(yeas+nays+abstentions))

	switch proposalType {
	case gov.Text:
		log.Debug("text proposal", "voteRate", voteRate, "required", xcom.TextProposal_VoteRate(), "supportRate", supportRate, "required", Decimal(xcom.TextProposal_SupportRate()))
		if voteRate > Decimal(xcom.TextProposal_VoteRate()) && supportRate >= Decimal(xcom.TextProposal_SupportRate()) {
			status = gov.Pass
		} else {
			status = gov.Failed
		}
	case gov.Cancel:
		log.Debug("cancel proposal", "voteRate", voteRate, "required", xcom.CancelProposal_VoteRate(), "supportRate", supportRate, "required", Decimal(xcom.CancelProposal_SupportRate()))
		if voteRate > Decimal(xcom.CancelProposal_VoteRate()) && supportRate >= Decimal(xcom.CancelProposal_SupportRate()) {
			status = gov.Pass
		} else {
			status = gov.Failed
		}
	}
	tallyResult := &gov.TallyResult{
		ProposalID:    proposalID,
		Yeas:          yeas,
		Nays:          nays,
		Abstentions:   abstentions,
		AccuVerifiers: verifiersCnt,
		Status:        status,
	}
	if err := gov.SetTallyResult(*tallyResult, state); err != nil {
		log.Error("save tally result failed", "tallyResult", tallyResult)
		return false, err
	}
	//gov.MoveVotingProposalIDToEnd(blockHash, proposalID, state)
	if err := gov.MoveVotingProposalIDToEnd(blockHash, proposalID); err != nil {
		log.Error("move proposalID from voting proposalID list to end list failed", "blockHash", blockHash, "proposalID", proposalID)
		return false, err
	}
	if err := gov.ClearActiveNodes(blockHash, proposalID); err != nil {
		return false, err
	}
	log.Debug("proposal tally result", "proposalID", proposalID, "tallyResult", tallyResult, "verifierList", verifierList)
	return status == gov.Pass, nil
}

func Decimal(value float64) int {
	return int(math.Floor(value * 1000))
}
