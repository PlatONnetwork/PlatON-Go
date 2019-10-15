package gov

import (
	"encoding/json"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	ValueDelimiter = []byte(":")
)

func SetProposal(proposal Proposal, state xcom.StateDB) error {
	bytes, e := json.Marshal(proposal)
	if e != nil {
		return e
	}

	value := append(bytes, byte(proposal.GetProposalType()))
	state.SetState(vm.GovContractAddr, KeyProposal(proposal.GetProposalID()), value)

	return AddPIPID(proposal.GetPIPID(), state)
}

func GetProposal(proposalID common.Hash, state xcom.StateDB) (Proposal, error) {
	value := state.GetState(vm.GovContractAddr, KeyProposal(proposalID))
	if len(value) == 0 {
		return nil, nil
	}
	pData := value[0 : len(value)-1]
	pType := value[len(value)-1]
	if pType == byte(Text) {
		var proposal TextProposal
		if e := json.Unmarshal(pData, &proposal); e != nil {
			log.Error("cannot parse data to text proposal")
			return nil, e
		}
		return &proposal, nil
	} else if pType == byte(Version) {
		var proposal VersionProposal
		if e := json.Unmarshal(pData, &proposal); e != nil {
			log.Error("cannot parse data to version proposal")
			return nil, e
		}
		return &proposal, nil
	} else if pType == byte(Cancel) {
		var proposal CancelProposal
		if e := json.Unmarshal(pData, &proposal); e != nil {
			log.Error("cannot parse data to cancel proposal")
			return nil, e
		}
		return &proposal, nil
	} else {
		return nil, common.InternalError.Wrap("Incorrect proposal type.")
	}
}

// Select proposal id list from snapshot database ,then get proposal detail from statedb one by one
func GetProposalList(blockHash common.Hash, state xcom.StateDB) ([]Proposal, error) {
	proposalIds, err := getAllProposalIDList(blockHash)
	if err != nil {
		return nil, err
	}
	var proposls []Proposal
	for _, proposalId := range proposalIds {
		proposal, err := GetExistProposal(proposalId, state)
		if err != nil {
			return nil, err
		} else {
			proposls = append(proposls, proposal)
		}
	}
	return proposls, nil
}

//Add the Vote detail
func AddVoteValue(proposalID common.Hash, voter discover.NodeID, option VoteOption, state xcom.StateDB) error {
	voteValueList, err := ListVoteValue(proposalID, state)
	if err != nil {
		return err
	}
	voteValueList = append(voteValueList, VoteValue{voter, option})
	return UpdateVoteValue(proposalID, voteValueList, state)
}

//list vote detail
func ListVoteValue(proposalID common.Hash, state xcom.StateDB) ([]VoteValue, error) {
	voteListBytes := state.GetState(vm.GovContractAddr, KeyVote(proposalID))
	if len(voteListBytes) == 0 {
		return nil, nil
	}
	var voteList []VoteValue
	if err := json.Unmarshal(voteListBytes, &voteList); err != nil {
		return nil, err
	}
	return voteList, nil
}

func UpdateVoteValue(proposalID common.Hash, voteValueList []VoteValue, state xcom.StateDB) error {
	voteListBytes, _ := json.Marshal(voteValueList)
	state.SetState(vm.GovContractAddr, KeyVote(proposalID), voteListBytes)
	return nil
}

// TallyVoteValue statistics vote option for a proposal
func TallyVoteValue(proposalID common.Hash, state xcom.StateDB) (yeas, nays, abstentions uint16, e error) {
	yes := uint16(0)
	no := uint16(0)
	abst := uint16(0)

	voteList, err := ListVoteValue(proposalID, state)
	if err == nil {
		for _, v := range voteList {
			if v.VoteOption == Yes {
				yes++
			}
			if v.VoteOption == No {
				no++
			}
			if v.VoteOption == Abstention {
				abst++
			}
		}
	}
	return yes, no, abst, err
}

/*
func ListVotedVerifier(proposalID common.Hash, state xcom.StateDB) ([]discover.NodeID, error) {
	var voterList []discover.NodeID
	valueList, err := ListVoteValue(proposalID, state)
	if err != nil {
		return nil, err
	}
	for _, value := range valueList {
		voterList = append(voterList, value.VoteNodeID)
	}

	return voterList, nil
}
*/

func GetVotedVerifierMap(proposalID common.Hash, state xcom.StateDB) (map[discover.NodeID]struct{}, error) {
	valueList, err := ListVoteValue(proposalID, state)
	if err != nil {
		return nil, err
	}

	votedMap := make(map[discover.NodeID]struct{}, len(valueList))
	for _, value := range valueList {
		votedMap[value.VoteNodeID] = struct{}{}
	}
	return votedMap, nil
}

func SetTallyResult(tallyResult TallyResult, state xcom.StateDB) error {
	value, err := json.Marshal(tallyResult)
	if err != nil {
		return err
	}
	state.SetState(vm.GovContractAddr, KeyTallyResult(tallyResult.ProposalID), value)
	return nil
}

func GetTallyResult(proposalID common.Hash, state xcom.StateDB) (*TallyResult, error) {
	proposal, err := GetProposal(proposalID, state)
	if err != nil {
		return nil, err
	} else if proposal == nil {
		return nil, ProposalNotFound
	}

	value := state.GetState(vm.GovContractAddr, KeyTallyResult(proposalID))

	if len(value) == 0 {
		return nil, nil
	}

	var tallyResult TallyResult
	if err := json.Unmarshal(value, &tallyResult); err != nil {
		return nil, err
	}
	return &tallyResult, nil

}

// Set pre-active version
func SetPreActiveVersion(preActiveVersion uint32, state xcom.StateDB) error {
	state.SetState(vm.GovContractAddr, KeyPreActiveVersion(), common.Uint32ToBytes(preActiveVersion))
	return nil
}

// Get pre-active version
func GetPreActiveVersion(state xcom.StateDB) uint32 {
	value := state.GetState(vm.GovContractAddr, KeyPreActiveVersion())
	return common.BytesToUint32(value)
}

// Set active version record
func AddActiveVersion(activeVersion uint32, activeBlock uint64, state xcom.StateDB) error {
	avList, err := ListActiveVersion(state)
	if err != nil {
		return err
	}
	curAv := ActiveVersionValue{ActiveVersion: activeVersion, ActiveBlock: activeBlock}
	//Insert the object into the head of the list
	avList = append([]ActiveVersionValue{curAv}, avList...)

	avListBytes, _ := json.Marshal(avList)
	state.SetState(vm.GovContractAddr, KeyActiveVersions(), avListBytes)
	return nil
}

// Get voting proposal
func ListVotingProposal(blockHash common.Hash) ([]common.Hash, error) {
	value, err := getVotingIDList(blockHash)
	if err != nil {
		log.Error("List voting proposal ID error")
		return nil, err
	}
	return value, nil
}

func ListEndProposalID(blockHash common.Hash) ([]common.Hash, error) {
	value, err := getEndIDList(blockHash)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func GetPreActiveProposalID(blockHash common.Hash) (common.Hash, error) {
	value, err := getPreActiveProposalID(blockHash)
	if err != nil {
		//log.Error("Get pre-active proposal ID error")
		return common.ZeroHash, err
	}
	return value, nil
}

func AddVotingProposalID(blockHash common.Hash, proposalID common.Hash) error {
	if err := addProposalByKey(blockHash, KeyVotingProposals(), proposalID); err != nil {
		//log.Error("add voting proposal to snapshot db error:%s", err)
		return err
	}

	return nil
}

func MoveVotingProposalIDToPreActive(blockHash common.Hash, proposalID common.Hash) error {
	voting, err := getVotingIDList(blockHash)
	if err != nil {
		return err
	}
	voting = remove(voting, proposalID)

	/*pre, err := self.snapdb.getPreActiveProposalID(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	pre = append(pre, proposalID)*/

	err = put(blockHash, KeyVotingProposals(), voting)
	if err != nil {
		return err
	}

	err = put(blockHash, KeyPreActiveProposal(), proposalID)
	if err != nil {
		return err
	}

	return nil
}

func remove(list []common.Hash, item common.Hash) []common.Hash {
	if len(list) == 0 {
		return list
	}

	for i, id := range list {
		if id == item {
			if len(list) > 1 {
				list = append(list[:i], list[i+1:]...)
			} else {
				list = []common.Hash{}
			}
		}
	}
	return list
}

func MoveVotingProposalIDToEnd(blockHash common.Hash, proposalID common.Hash) error {

	voting, err := getVotingIDList(blockHash)
	if err != nil {
		return err
	}
	voting = remove(voting, proposalID)
	err = put(blockHash, KeyVotingProposals(), voting)
	if err != nil {
		return err
	}

	err = addProposalByKey(blockHash, KeyEndProposals(), proposalID)
	if err != nil {
		return err
	}

	return nil
}

func MovePreActiveProposalIDToEnd(blockHash common.Hash, proposalID common.Hash) error {
	//only one proposalID in PreActiveProposalIDList, so, just set it empty.
	err := put(blockHash, KeyPreActiveProposal(), common.Hash{})
	if err != nil {
		return err
	}

	err = addProposalByKey(blockHash, KeyEndProposals(), proposalID)
	if err != nil {
		return err
	}

	return nil
}

// Add the node that has made a new version declare or vote during voting period
func AddActiveNode(blockHash common.Hash, proposalID common.Hash, nodeID discover.NodeID) error {
	if err := addActiveNode(blockHash, nodeID, proposalID); err != nil {
		log.Error("add active node to snapshot db failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return err
	}
	return nil
}

// Get the node list that have made a new version declare or vote during voting period
func GetActiveNodeList(blockHash common.Hash, proposalID common.Hash) ([]discover.NodeID, error) {
	nodes, err := getActiveNodeList(blockHash, proposalID)
	if err != nil {
		log.Error("get active nodes from snapshot db failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return nil, err
	}
	return nodes, nil
}

// Clear the version declaration records after upgrade
func ClearActiveNodes(blockHash common.Hash, proposalID common.Hash) error {
	err := deleteActiveNodeList(blockHash, proposalID)
	if err != nil {
		log.Error("clear active nodes in snapshot db failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return err
	}
	return nil
}

// AccuVerifiers accumulates all distinct verifiers those can vote this proposal ID
func AccuVerifiers(blockHash common.Hash, proposalID common.Hash, verifierList []discover.NodeID) error {
	if err := addAccuVerifiers(blockHash, proposalID, verifierList); err != nil {
		log.Error("accumulates verifiers to snapshot db failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return err
	}
	return nil
}

// Get the total number of all voting verifiers
func ListAccuVerifier(blockHash common.Hash, proposalID common.Hash) ([]discover.NodeID, error) {
	if l, err := getAccuVerifiers(blockHash, proposalID); err != nil {
		log.Error("list accumulated verifiers failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return nil, err
	} else {
		return l, nil
	}
}

func AddPIPID(pipID string, state xcom.StateDB) error {
	pipIDList, err := ListPIPID(state)
	if err != nil {
		return err
	}

	if pipIDList == nil || len(pipIDList) == 0 {
		pipIDList = []string{pipID}
	} else {
		pipIDList = append(pipIDList, pipID)
	}

	pipIDListBytes, _ := json.Marshal(pipIDList)
	state.SetState(vm.GovContractAddr, KeyPIPIDs(), pipIDListBytes)
	return nil
}

func ListPIPID(state xcom.StateDB) ([]string, error) {
	pipIDListBytes := state.GetState(vm.GovContractAddr, KeyPIPIDs())
	if len(pipIDListBytes) > 0 {
		var pipIDList []string
		if err := json.Unmarshal(pipIDListBytes, &pipIDList); err != nil {
			return nil, err
		}
		return pipIDList, nil
	} else {
		return nil, nil
	}
}

func GetExistProposal(proposalID common.Hash, state xcom.StateDB) (Proposal, error) {
	p, err := GetProposal(proposalID, state)
	if err != nil {
		return nil, err
	} else if p == nil {
		//log.Error("Cannot find proposal.", "proposalID", proposalID)
		return nil, ProposalNotFound
	} else {
		return p, nil
	}
}

// find a version proposal at voting stage
func FindVotingVersionProposal(blockHash common.Hash, state xcom.StateDB) (*VersionProposal, error) {
	log.Debug("call findVotingVersionProposal", "blockHash", blockHash)
	idList, err := ListVotingProposal(blockHash)
	if err != nil {
		log.Error("find voting version proposal failed", "blockHash", blockHash)
		return nil, err
	}
	for _, proposalID := range idList {
		p, err := GetExistProposal(proposalID, state)
		if err != nil {
			return nil, err
		}
		if p.GetProposalType() == Version {
			vp := p.(*VersionProposal)
			return vp, nil
		}
	}
	return nil, nil
}

func ListActiveVersion(state xcom.StateDB) ([]ActiveVersionValue, error) {
	avListBytes := state.GetState(vm.GovContractAddr, KeyActiveVersions())
	if len(avListBytes) == 0 {
		return nil, nil
	}
	var avList []ActiveVersionValue
	if err := json.Unmarshal(avListBytes, &avList); err != nil {
		return nil, err
	}
	return avList, nil
}
