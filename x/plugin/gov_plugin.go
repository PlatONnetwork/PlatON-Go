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

package plugin

import (
	"math"
	"math/big"
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
	chainID *big.Int
}

var govp *GovPlugin

func GovPluginInstance() *GovPlugin {
	govPluginOnce.Do(func() {
		log.Info("Init Governance plugin ...")
		govp = &GovPlugin{}
	})
	return govp
}

func (govPlugin *GovPlugin) SetChainID(chainId *big.Int) {
	govPlugin.chainID = chainId
}
func (govPlugin *GovPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

//implement BasePlugin
func (govPlugin *GovPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	var blockNumber = header.Number.Uint64()
	//log.Debug("call BeginBlock()", "blockNumber", blockNumber, "blockHash", blockHash)

	if !xutil.IsBeginOfConsensus(blockNumber) {
		return nil
	}

	if xutil.IsBeginOfEpoch(blockNumber) {
		if err := accuVerifiersAtBeginOfSettlement(blockHash, blockNumber); err != nil {
			log.Error("accumulates all distinct verifiers for voting proposal failed.", "blockNumber", blockNumber, "err", err)
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
		//log.Debug("found pre-active version proposal", "proposalID", preActiveVersionProposalID, "blockNumber", blockNumber, "blockHash", blockHash, "activeBlockNumber", versionProposal.GetActiveBlock())
		if blockNumber == versionProposal.GetActiveBlock() {
			//log.Debug("it's time to active the pre-active version proposal")
			tallyResult, err := gov.GetTallyResult(preActiveVersionProposalID, state)
			if err != nil || tallyResult == nil {
				log.Error("find pre-active version proposal tally result failed.", "blockNumber", blockNumber, "blockHash", blockHash, "preActiveVersionProposalID", preActiveVersionProposalID)
				return err
			}
			//update tally status to "active"
			tallyResult.Status = gov.Active

			if err := gov.SetTallyResult(*tallyResult, state); err != nil {
				log.Error("update version proposal tally result failed.", "blockNumber", blockNumber, "preActiveVersionProposalID", preActiveVersionProposalID)
				return err
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

			log.Info("version proposal is active", "blockNumber", blockNumber, "proposalID", versionProposal.ProposalID, "newVersion", versionProposal.NewVersion, "newVersionString", xutil.ProgramVersion2Str(versionProposal.NewVersion))
		}
	}
	return nil
}

//implement BasePlugin
func (govPlugin *GovPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	var blockNumber = header.Number.Uint64()
	//log.Debug("call EndBlock()", "blockNumber", blockNumber, "blockHash", blockHash)

	//text/version/cancel proposal's end voting block is ElectionBlock
	//param proposal's end voting block is end of Epoch
	isEndOfEpoch := false
	isElection := false
	if xutil.IsElection(blockNumber) {
		isElection = true
	} else if xutil.IsEndOfEpoch(blockNumber) {
		isEndOfEpoch = true
	} else {
		return nil
	}

	votingProposalIDs, err := gov.ListVotingProposal(blockHash)
	if err != nil {
		return err
	}
	if len(votingProposalIDs) == 0 {
		//log.Debug("there's no voting proposal", "blockNumber", blockNumber, "blockHash", blockHash)
		return nil
	}

	//iterate each voting proposal, to check if current block is proposal's end-voting block.
	for _, votingProposalID := range votingProposalIDs {
		//log.Debug("iterate each voting proposal", "proposalID", votingProposalID)
		votingProposal, err := gov.GetExistProposal(votingProposalID, state)
		//log.Debug("find voting proposal", "votingProposal", votingProposal)
		if nil != err {
			return err
		}
		if votingProposal.GetEndVotingBlock() == blockNumber {
			log.Debug("current block is end-voting block", "proposalID", votingProposal.GetProposalID(), "blockNumber", blockNumber)
			//tally the results
			if votingProposal.GetProposalType() == gov.Text && isElection {
				_, err := tallyText(votingProposal.(*gov.TextProposal), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
			} else if votingProposal.GetProposalType() == gov.Version && isElection {
				err = tallyVersion(votingProposal.(*gov.VersionProposal), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
			} else if votingProposal.GetProposalType() == gov.Cancel && isElection {
				_, err := tallyCancel(votingProposal.(*gov.CancelProposal), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
			} else if votingProposal.GetProposalType() == gov.Param && isEndOfEpoch {
				_, err := tallyParam(votingProposal.(*gov.ParamProposal), blockHash, blockNumber, state)
				if err != nil {
					return err
				}
			} else {
				log.Error("invalid proposal type", "type", votingProposal.GetProposalType())
				return gov.ProposalTypeError
			}
		}
	}
	return nil
}

/*func NewVerifiersForNextEpoch(newVerifiers []discover.NodeID, endBlockHashOfCurrentEpoch common.Hash, endBlockNumberOfCurrentEpoch uint64) error {
	if xutil.IsEndOfEpoch(endBlockNumberOfCurrentEpoch) {
		if err := accuVerifiersAtBeginOfSettlement(newVerifiers, endBlockHashOfCurrentEpoch, endBlockNumberOfCurrentEpoch); err != nil {
			log.Error("accumulates all distinct verifiers for voting proposal failed.", "err", err)
			return err
		}
	}
	return nil
}*/

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
	//log.Debug("proposal tally", "proposalID", proposal.ProposalID, "blockHash", blockHash, "blockNumber", blockNumber)

	verifierList, err := gov.ListAccuVerifier(blockHash, proposalID)
	if err != nil {
		return err
	}
	verifiersCnt := uint64(len(verifierList))

	voteList, err := gov.ListVoteValue(proposalID, blockHash)
	if err != nil {
		return err
	}

	voteCnt := uint64(len(voteList))
	yeas := voteCnt //`voteOption` can be ignored in version proposal, set voteCount to passCount as default.

	status := gov.Failed
	var supportRate uint64 = 0
	if verifiersCnt > 0 {
		supportRate = (yeas * gov.RateCoefficient) / verifiersCnt
	}

	//log.Debug("version proposal", "supportRate", supportRate, "required", Decimal(xcom.VersionProposalSupportRate()))

	if supportRate >= xcom.VersionProposal_SupportRate() {
		status = gov.PreActive

		if err := gov.AddPIPID(proposal.GetPIPID(), state); err != nil {
			log.Error("save passed PIPID failed", "proposalID", proposalID, "blockNumber", blockNumber, "blockHash", blockHash)
			return err
		}

		if err := gov.MoveVotingProposalIDToPreActive(blockHash, proposalID, proposal.NewVersion); err != nil {
			log.Error("move version proposal ID to pre-active failed", "proposalID", proposalID, "blockNumber", blockNumber, "blockHash", blockHash)
			return err
		}

		if err := gov.SetPreActiveVersion(blockHash, proposal.NewVersion); err != nil {
			log.Error("save pre-active version to state failed", "proposalID", proposalID, "blockHash", blockHash, "newVersion", proposal.NewVersion, "newVersionString", xutil.ProgramVersion2Str(proposal.NewVersion))
			return err
		}

		activeList, err := gov.GetActiveNodeList(blockHash, proposalID)
		if err != nil {
			log.Error("list active node failed", "proposalID", proposalID, "blockNumber", blockNumber, "blockHash", blockHash)
			return err
		}
		//log.Debug("call stk.ProposalPassedNotify", "proposalID", proposalID, "activeList", activeList)
		if err := stk.ProposalPassedNotify(blockHash, blockNumber, activeList, proposal.NewVersion); err != nil {
			log.Error("call stk.ProposalPassedNotify failed", "proposalID", proposalID, "blockHash", blockHash, "newVersion", proposal.NewVersion, "activeList", activeList)
			return err
		}

	} else {
		if err := gov.MoveVotingProposalIDToEnd(proposalID, blockHash); err != nil {
			log.Error("move proposalID from voting proposalID list to end list failed", "proposalID", proposalID, "blockNumber", blockNumber, "blockHash", blockHash)
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
		log.Error("save tally result failed", "blockNumber", blockNumber, "blockHash", blockHash, "proposalID", proposalID, "tallyResult", tallyResult)
		return err
	}

	// for now, do not remove these data.
	// If really want to remove these data, please confirmed with PlatON Explorer Project
	/*if err := gov.ClearVoteValue(proposalID, blockHash); err != nil {
		log.Error("clear vote value failed", "proposalID", proposalID, "blockHash", blockHash, "err", err)
		return err
	}*/

	log.Info("version proposal tally result", "blockNumber", blockNumber, "blockHash", blockHash, "proposalID", proposalID, "tallyResult", tallyResult, "verifierList", verifierList)
	return nil
}

func tallyText(tp *gov.TextProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (pass bool, err error) {
	return tally(gov.Text, tp.ProposalID, tp.PIPID, blockHash, blockNumber, state)
}

func tallyCancel(cp *gov.CancelProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (pass bool, err error) {
	if pass, err := tally(gov.Cancel, cp.ProposalID, cp.PIPID, blockHash, blockNumber, state); err != nil {
		log.Info("canceled a proposal failed", "proposalID", cp.TobeCanceled, "tobeCanceledProposalID", cp.TobeCanceled)
		return false, err
	} else if pass {
		if proposal, err := gov.GetExistProposal(cp.TobeCanceled, state); err != nil {
			return false, err
		} else if proposal.GetProposalType() != gov.Version && proposal.GetProposalType() != gov.Param {
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
				log.Error("the to be canceled proposal is not at voting stage, but the cancel proposal is passed")
				return false, err
			}
			verifierList, err := gov.ListAccuVerifier(blockHash, cp.TobeCanceled)
			if err != nil {
				return false, err
			}
			verifiersCnt := uint64(len(verifierList))

			voteList, err := gov.ListVoteValue(cp.TobeCanceled, blockHash)
			if err != nil {
				return false, err
			}

			voteCnt := uint64(len(voteList))
			yeas := voteCnt

			tallyResult.Yeas = yeas
			tallyResult.AccuVerifiers = verifiersCnt
			tallyResult.Status = gov.Canceled
			tallyResult.CanceledBy = cp.ProposalID

			if err := gov.SetTallyResult(*tallyResult, state); err != nil {
				log.Error("to cancel a proposal failed, cannot save its tally result", "blockNumber", blockNumber, "blockHash", blockHash, "proposalID", cp.ProposalID, "tallyResult", tallyResult)
				return false, err
			}

			if cp.ProposalType == gov.Version {
				if err := gov.ClearActiveNodes(blockHash, cp.TobeCanceled); err != nil {
					return false, err
				}
			}

			if err := gov.MoveVotingProposalIDToEnd(cp.TobeCanceled, blockHash); err != nil {
				return false, err
			}

			log.Info("canceled a proposal success", "blockNumber", blockNumber, "blockHash", blockHash, "proposalID", cp.TobeCanceled, "tallyResult", tallyResult)
		}
	}
	return true, nil
}

func tallyParam(pp *gov.ParamProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (pass bool, err error) {
	if pass, err := tally(gov.Param, pp.ProposalID, pp.PIPID, blockHash, blockNumber, state); err != nil {
		return false, err
	} else if pass {
		if err := gov.UpdateGovernParamValue(pp.Module, pp.Name, pp.NewValue, blockNumber+1, blockHash); err != nil {
			return false, err
		}
	}
	return true, nil
}

func tally(proposalType gov.ProposalType, proposalID common.Hash, pipID string, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (pass bool, err error) {
	//log.Debug("proposal tally", "proposalID", proposalID, "blockHash", blockHash, "blockNumber", blockNumber, "proposalID", proposalID)

	verifierList, err := gov.ListAccuVerifier(blockHash, proposalID)
	if err != nil {
		return false, err
	}

	verifiersCnt := uint64(len(verifierList))

	status := gov.Voting

	yeas, nays, abstentions, err := gov.TallyVoteValue(proposalID, blockHash)
	if err != nil {
		return false, err
	}
	var voteRate uint64 = 0
	var supportRate uint64 = 0
	if yeas+nays+abstentions > 0 && verifiersCnt > 0 {
		voteRate = (yeas + nays + abstentions) * gov.RateCoefficient / verifiersCnt
		supportRate = (yeas * gov.RateCoefficient) / (yeas + nays + abstentions)
	}

	switch proposalType {
	case gov.Text:
		//log.Debug("text proposal", "voteRate", voteRate, "required", xcom.TextProposalVoteRate(), "supportRate", supportRate, "required", Decimal(xcom.TextProposalSupportRate()))
		if voteRate > xcom.TextProposal_VoteRate() && supportRate >= xcom.TextProposal_SupportRate() {
			status = gov.Pass
		} else {
			status = gov.Failed
		}
	case gov.Cancel:
		//log.Debug("cancel proposal", "voteRate", voteRate, "required", xcom.CancelProposalVoteRate(), "supportRate", supportRate, "required", Decimal(xcom.CancelProposalSupportRate()))
		if voteRate > xcom.CancelProposal_VoteRate() && supportRate >= xcom.CancelProposal_SupportRate() {
			status = gov.Pass
		} else {
			status = gov.Failed
		}
	case gov.Param:
		//log.Debug("param proposal", "voteRate", voteRate, "required", xcom.ParamProposalVoteRate(), "supportRate", supportRate, "required", Decimal(xcom.ParamProposalSupportRate()))
		if voteRate > xcom.ParamProposal_VoteRate() && supportRate >= xcom.ParamProposal_SupportRate() {
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
	if err := gov.MoveVotingProposalIDToEnd(proposalID, blockHash); err != nil {
		log.Error("move proposalID from voting proposalID list to end list failed", "proposalID", proposalID, "blockNumber", blockNumber, "blockHash", blockHash, "err", err)
		return false, err
	}

	if status == gov.Pass {
		if err := gov.AddPIPID(pipID, state); err != nil {
			log.Error("save passed PIPID failed", "proposalID", proposalID, "blockNumber", blockNumber, "blockHash", blockHash)
			return false, err
		}
	}
	// for now, do not remove these data.
	// If really want to remove these data, please confirmed with PlatON Explorer Project
	/*if err := gov.ClearVoteValue(proposalID, blockHash); err != nil {
		log.Error("clear vote value failed", "proposalID", proposalID, "blockHash", blockHash, "err", err)
		return false, err
	}*/

	log.Debug("proposal tally result", "blockNumber", blockNumber, "blockHash", blockHash, "proposalID", proposalID, "tallyResult", tallyResult, "verifierList", verifierList)
	return status == gov.Pass, nil
}

func Decimal(value float64) int {
	return int(math.Floor(value * 1000))
}
