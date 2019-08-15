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

//var dbOnce sync.Once
var govDB *GovDB

type GovDB struct {
	snapdb GovSnapshotDB
}

func GovDBInstance() *GovDB {
	//dbOnce.Do(func() {
	//	govDB = &GovDB{snapdb: NewGovSnapshotDB()}
	//})
	if govDB == nil {
		log.Info("Init GovDB ........")
		govDB = &GovDB{snapdb: NewGovSnapshotDB()}
	}
	return govDB
}

func (self *GovDB) Reset() {
	govDB = nil
	self.snapdb.reset()
}

func (self *GovDB) SetProposal(proposal Proposal, state xcom.StateDB) error {
	bytes, e := json.Marshal(proposal)
	if e != nil {
		return common.NewSysError(e.Error())
	}

	value := append(bytes, byte(proposal.GetProposalType()))
	state.SetState(vm.GovContractAddr, KeyProposal(proposal.GetProposalID()), value)

	return self.AddPIPID(proposal.GetPIPID(), state)
}

func (self *GovDB) GetProposal(proposalID common.Hash, state xcom.StateDB) (Proposal, error) {
	value := state.GetState(vm.GovContractAddr, KeyProposal(proposalID))
	if len(value) == 0 {
		return nil, nil
	}
	var p Proposal
	pData := value[0 : len(value)-1]
	pType := value[len(value)-1]
	if pType == byte(Text) {
		var proposal TextProposal
		if e := json.Unmarshal(pData, &proposal); e != nil {
			log.Error("cannot parse data to text proposal")
			return nil, common.NewSysError(e.Error())
		}
		p = proposal
	} else if pType == byte(Version) {
		var proposal VersionProposal
		if e := json.Unmarshal(pData, &proposal); e != nil {
			log.Error("cannot parse data to version proposal")
			return nil, common.NewSysError(e.Error())
		}
		p = proposal
	} else if pType == byte(Cancel) {
		var proposal CancelProposal
		if e := json.Unmarshal(pData, &proposal); e != nil {
			log.Error("cannot parse data to cancel proposal")
			return nil, common.NewSysError(e.Error())
		}
		p = proposal
	} else {
		return nil, common.NewSysError("Incorrect proposal type.")
	}

	return p, nil
}

func (self *GovDB) GetExistProposal(proposalID common.Hash, state xcom.StateDB) (Proposal, error) {
	p, err := self.GetProposal(proposalID, state)
	if err != nil {
		return nil, err
	} else if p == nil {
		log.Error("Cannot find proposal.", "proposalID", proposalID)
		return nil, common.NewSysError("Cannot find proposal.")
	} else {
		return p, nil
	}
}

// Select proposal id list from snapshot database ,then get propsal detail from statedb one by one
func (self *GovDB) GetProposalList(blockHash common.Hash, state xcom.StateDB) ([]Proposal, error) {
	proposalIds, err := self.snapdb.getAllProposalIDList(blockHash)
	if err != nil {
		return nil, common.NewSysError(err.Error())
	}
	var proposls []Proposal
	for _, proposalId := range proposalIds {
		proposal, err := self.GetExistProposal(proposalId, state)
		if err != nil {
			return nil, err
		} else {
			proposls = append(proposls, proposal)
		}
	}
	return proposls, nil
}

//Save the Vote detail
func (self *GovDB) SetVote(proposalID common.Hash, voter discover.NodeID, option VoteOption, state xcom.StateDB) error {
	voteValueList, err := self.ListVoteValue(proposalID, state)
	if err != nil {
		return common.NewSysError(err.Error())
	}
	voteValueList = append(voteValueList, VoteValue{voter, option})

	voteListBytes, _ := json.Marshal(voteValueList)

	state.SetState(vm.GovContractAddr, KeyVote(proposalID), voteListBytes)
	return nil
}

//list vote detail
func (self *GovDB) ListVoteValue(proposalID common.Hash, state xcom.StateDB) ([]VoteValue, error) {
	voteListBytes := state.GetState(vm.GovContractAddr, KeyVote(proposalID))
	if len(voteListBytes) == 0 {
		return nil, nil
	}
	var voteList []VoteValue
	if err := json.Unmarshal(voteListBytes, &voteList); err != nil {
		return nil, common.NewSysError(err.Error()) //errors.New("Unmarshal VoteValue error")
	}
	return voteList, nil
}

func (self *GovDB) ListVotedVerifier(proposalID common.Hash, state xcom.StateDB) ([]discover.NodeID, error) {
	var voterList []discover.NodeID
	valueList, err := self.ListVoteValue(proposalID, state)
	if err != nil {
		return nil, common.NewSysError(err.Error())
	}
	for _, value := range valueList {
		voterList = append(voterList, value.VoteNodeID)
	}

	return voterList, nil
}

func (self *GovDB) SetTallyResult(tallyResult TallyResult, state xcom.StateDB) error {
	value, err := json.Marshal(tallyResult)
	if err != nil {
		return common.NewSysError(err.Error())
	}
	state.SetState(vm.GovContractAddr, KeyTallyResult(tallyResult.ProposalID), value)
	return nil
}

func (self *GovDB) GetTallyResult(proposalID common.Hash, state xcom.StateDB) (*TallyResult, error) {
	value := state.GetState(vm.GovContractAddr, KeyTallyResult(proposalID))

	if len(value) == 0 {
		return nil, nil
	}

	var tallyResult TallyResult
	if err := json.Unmarshal(value, &tallyResult); err != nil {
		return nil, common.NewSysError(err.Error())
	}
	return &tallyResult, nil

}

// Set pre-active version
func (self *GovDB) SetPreActiveVersion(preActiveVersion uint32, state xcom.StateDB) error {
	state.SetState(vm.GovContractAddr, KeyPreActiveVersion(), common.Uint32ToBytes(preActiveVersion))
	return nil
}

// Get pre-active version
func (self *GovDB) GetPreActiveVersion(state xcom.StateDB) uint32 {
	value := state.GetState(vm.GovContractAddr, KeyPreActiveVersion())
	return common.BytesToUint32(value)
}

// Set active version record
func (self *GovDB) AddActiveVersion(activeVersion uint32, activeBlock uint64, state xcom.StateDB) error {
	avList, err := self.ListActiveVersion(state)
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

func (self *GovDB) ListActiveVersion(state xcom.StateDB) ([]ActiveVersionValue, error) {
	avListBytes := state.GetState(vm.GovContractAddr, KeyActiveVersions())
	if len(avListBytes) == 0 {
		return nil, nil
	}
	var avList []ActiveVersionValue
	if err := json.Unmarshal(avListBytes, &avList); err != nil {
		return nil, common.NewSysError(err.Error())
	}
	return avList, nil
}

// Get current active version record
func (self *GovDB) GetCurrentActiveVersion(state xcom.StateDB) uint32 {
	avList, err := self.ListActiveVersion(state)
	if err != nil {
		log.Error("Cannot find active version list")
		return 0
	}
	return avList[0].ActiveVersion
}

// Get voting proposal
func (self *GovDB) ListVotingProposal(blockHash common.Hash) ([]common.Hash, error) {
	value, err := govDB.snapdb.getVotingIDList(blockHash)
	if err != nil {
		log.Error("List voting proposal ID error")
		return nil, common.NewSysError(err.Error())
	}
	return value, nil
}

func (self *GovDB) ListEndProposalID(blockHash common.Hash) ([]common.Hash, error) {
	value, err := govDB.snapdb.getEndIDList(blockHash)
	if err != nil {
		return nil, common.NewSysError(err.Error())
	}

	return value, nil
}

func (self *GovDB) GetPreActiveProposalID(blockHash common.Hash) (common.Hash, error) {
	value, err := govDB.snapdb.getPreActiveProposalID(blockHash)
	if err != nil {
		//log.Error("Get pre-active proposal ID error")
		return common.ZeroHash, common.NewSysError(err.Error())
	}
	return value, nil
}

func (self *GovDB) AddVotingProposalID(blockHash common.Hash, proposalID common.Hash) error {
	if err := govDB.snapdb.addProposalByKey(blockHash, KeyVotingProposals(), proposalID); err != nil {
		//log.Error("add voting proposal to snapshot db error:%s", err)
		return common.NewSysError(err.Error())
	}

	return nil
}

func (self *GovDB) MoveVotingProposalIDToPreActive(blockHash common.Hash, proposalID common.Hash) error {

	voting, err := self.snapdb.getVotingIDList(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}
	voting = remove(voting, proposalID)

	/*pre, err := self.snapdb.getPreActiveProposalID(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	pre = append(pre, proposalID)*/

	err = self.snapdb.put(blockHash, KeyVotingProposals(), voting)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	err = self.snapdb.put(blockHash, KeyPreActiveProposal(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
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

func (self *GovDB) MoveVotingProposalIDToEnd(blockHash common.Hash, proposalID common.Hash, state xcom.StateDB) error {

	voting, err := self.snapdb.getVotingIDList(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}
	voting = remove(voting, proposalID)
	err = self.snapdb.put(blockHash, KeyVotingProposals(), voting)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	err = self.snapdb.addProposalByKey(blockHash, KeyEndProposals(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	return nil
}

func (self *GovDB) MovePreActiveProposalIDToEnd(blockHash common.Hash, proposalID common.Hash, state xcom.StateDB) error {
	//only one proposalID in PreActiveProposalIDList, so, just set it empty.
	err := self.snapdb.put(blockHash, KeyPreActiveProposal(), common.Hash{})
	if err != nil {
		return common.NewSysError(err.Error())
	}

	err = self.snapdb.addProposalByKey(blockHash, KeyEndProposals(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	return nil
}

// Add the node that has made a new version declare or vote during voting period
func (self *GovDB) AddActiveNode(blockHash common.Hash, proposalID common.Hash, nodeID discover.NodeID) error {
	if err := self.snapdb.addActiveNode(blockHash, nodeID, proposalID); err != nil {
		log.Error("add active node to snapshot db failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return common.NewSysError(err.Error())
	}
	return nil
}

// Get the node list that have made a new version declare or vote during voting period
func (self *GovDB) GetActiveNodeList(blockHash common.Hash, proposalID common.Hash) ([]discover.NodeID, error) {
	nodes, err := self.snapdb.getActiveNodeList(blockHash, proposalID)
	if err != nil {
		log.Error("get active nodes from snapshot db failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return nil, common.NewSysError(err.Error())
	}
	return nodes, nil
}

// Clear the version declaration records after upgrade
func (self *GovDB) ClearActiveNodes(blockHash common.Hash, proposalID common.Hash) error {
	err := self.snapdb.deleteActiveNodeList(blockHash, proposalID)
	if err != nil {
		log.Error("clear active nodes in snapshot db failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return common.NewSysError(err.Error())
	}
	return nil
}

// All verifiers who can vote accumulatively in the settlement cycle
func (self *GovDB) AccuVerifiers(blockHash common.Hash, proposalID common.Hash, verifierList []discover.NodeID) error {
	if err := self.snapdb.addAccuVerifiers(blockHash, proposalID, verifierList); err != nil {
		log.Error("save accumulated verifiers to snapshot db failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return common.NewSysError(err.Error())
	}
	return nil
}

// Get the total number of all voting verifiers
func (self *GovDB) AccuVerifiersLength(blockHash common.Hash, proposalID common.Hash) (uint16, error) {
	if l, err := self.snapdb.getAccuVerifiersLength(blockHash, proposalID); err != nil {
		log.Error("get accumulated verifiers count failed", "blockHash", blockHash.String(), "proposalID", proposalID, "error", err)
		return 0, common.NewSysError(err.Error())
	} else {
		return l, nil
	}
}

func (self *GovDB) AddPIPID(pipID string, state xcom.StateDB) error {
	pipIDList, err := self.ListPIPID(state)
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

func (self *GovDB) ListPIPID(state xcom.StateDB) ([]string, error) {
	pipIDListBytes := state.GetState(vm.GovContractAddr, KeyPIPIDs())
	if len(pipIDListBytes) > 0 {
		var pipIDList []string
		if err := json.Unmarshal(pipIDListBytes, &pipIDList); err != nil {
			return nil, common.NewSysError(err.Error())
		}
		return pipIDList, nil
	} else {
		return nil, nil
	}
}

// find a version proposal at voting stage
func (self *GovDB) FindVotingVersionProposal(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (*VersionProposal, error) {
	log.Debug("call findVotingVersionProposal", "blockHash", blockHash, "blockNumber", blockNumber)
	idList, err := self.ListVotingProposal(blockHash)
	if err != nil {
		log.Error("find voting version proposal failed", "blockHash", blockHash)
		return nil, err
	}
	for _, proposalID := range idList {
		p, err := self.GetExistProposal(proposalID, state)
		if err != nil {
			return nil, err
		}
		if p.GetProposalType() == Version {
			vp := p.(VersionProposal)
			return &vp, nil
		}
	}
	return nil, nil
}

// find a cancel proposal at voting stage
func (self *GovDB) FindVotingCancelProposal(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) (*CancelProposal, error) {
	log.Debug("call findVotingCancelProposal", "blockHash", blockHash, "blockNumber", blockNumber)
	idList, err := self.ListVotingProposal(blockHash)
	if err != nil {
		log.Error("find voting proposal failed", "blockHash", blockHash)
		return nil, err
	}
	for _, proposalID := range idList {
		p, err := self.GetExistProposal(proposalID, state)
		if err != nil {
			return nil, err
		}
		if p.GetProposalType() == Cancel {
			vp := p.(CancelProposal)
			return &vp, nil
		}
	}
	return nil, nil
}
