package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

type GovSnapshotDB struct {
	snapdb snapshotdb.DB
}

func NewGovSnapshotDB() *GovSnapshotDB {
	return &GovSnapshotDB{snapdb: snapshotdb.Instance()}
}

func (self *GovSnapshotDB) get(blockHash common.Hash, key []byte) ([]byte, error) {
	return self.snapdb.Get(blockHash, key)
}

func (self *GovSnapshotDB) put(blockHash common.Hash, key, value []byte) error {
	return self.snapdb.Put(blockHash, key, value)
}

func (self *GovSnapshotDB) del(blockHash common.Hash, key []byte) error {
	return self.snapdb.Del(blockHash, key)
}

func (self *GovSnapshotDB) addProposalByKey(blockHash common.Hash, key []byte, proposalId common.Hash) error {
	hashes, err := self.getProposalIDListByKey(blockHash, key)
	if err != nil {
		return err
	}

	hashes = append(hashes, proposalId)

	//重新设置,在编码
	value, err := rlp.EncodeToBytes(hashes)
	if err != nil {
		return err
	}

	return self.put(blockHash, key, value)
}

func (self *GovSnapshotDB) getVotingIDList(blockHash common.Hash) ([]common.Hash, error) {
	return self.getProposalIDListByKey(blockHash, KeyVotingProposals())
}

func (self *GovSnapshotDB) getPreActiveIDList(blockHash common.Hash) ([]common.Hash, error) {
	return self.getProposalIDListByKey(blockHash, KeyPreActiveProposals())
}

func (self *GovSnapshotDB) getEndIDList(blockHash common.Hash) ([]common.Hash, error) {
	return self.getProposalIDListByKey(blockHash, KeyEndProposals())
}

func (self *GovSnapshotDB) getProposalIDListByKey(blockHash common.Hash, key []byte) ([]common.Hash, error) {
	bytes, err := self.get(blockHash, key)
	if err != nil {
		return nil, err
	}

	var idList []common.Hash

	if err = rlp.DecodeBytes(bytes, &idList); err != nil {
		return nil, err
	}

	return idList, nil
}

func (self *GovSnapshotDB) getAllProposalIDList(blockHash common.Hash) ([]common.Hash, error) {
	var total []common.Hash

	hashes, _ := self.getVotingIDList(blockHash)
	total = append(total, hashes...)

	hashes, _ = self.getPreActiveIDList(blockHash)
	total = append(total, hashes...)

	hashes, _ = self.getEndIDList(blockHash)
	total = append(total, hashes...)
	return total, nil
}

func (self *GovSnapshotDB) addVotedVerifier(blockHash common.Hash, node discover.NodeID, proposalId common.Hash) error {

	nodes, err := self.getVotedVerifierList(blockHash, proposalId)
	if err != nil {
		return err
	}

	nodes = append(nodes, node)

	value, err := rlp.EncodeToBytes(nodes)
	if err != nil {
		return err
	}

	self.put(blockHash, keyPrefixVotedVerifiers, value)
	return nil
}

func (self *GovSnapshotDB) getVotedVerifierList(blockHash common.Hash, proposalId common.Hash) ([]discover.NodeID, error) {
	value, err := self.get(blockHash, KeyVotedVerifiers(proposalId))
	if err != nil {
		return nil, err
	}
	var nodes []discover.NodeID
	if err := rlp.DecodeBytes(value, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (self *GovSnapshotDB) addActiveNode(blockHash common.Hash, node discover.NodeID, proposalId common.Hash) error {

	nodes, err := self.getActiveNodeList(blockHash, proposalId)
	if err != nil {
		return err
	}

	nodes = append(nodes, node)

	value, err := rlp.EncodeToBytes(nodes)
	if err != nil {
		return err
	}

	self.put(blockHash, keyPrefixActiveNodes, value)
	return nil
}

func (self *GovSnapshotDB) getActiveNodeList(blockHash common.Hash, proposalId common.Hash) ([]discover.NodeID, error) {
	value, err := self.get(blockHash, KeyActiveNodes(proposalId))
	if err != nil {
		return nil, err
	}
	var nodes []discover.NodeID
	if err := rlp.DecodeBytes(value, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (self *GovSnapshotDB) deleteActiveNodeList(blockHash common.Hash, proposalId common.Hash) error {
	return self.del(blockHash, KeyActiveNodes(proposalId))
}

func (self *GovSnapshotDB) addTotalVerifiers(blockHash common.Hash, proposalId common.Hash, nodes []discover.NodeID) error {
	value, err := self.get(blockHash, KeyAccuVerifier(proposalId))
	if err != nil {
		return err
	}

	var verifiers []discover.NodeID
	if err := rlp.DecodeBytes(value, &verifiers); err != nil {
		return err
	}

	verifiers = append(verifiers, nodes...)

	return nil
}

func (self *GovSnapshotDB) getAccuVerifiersLength(blockHash common.Hash, proposalId common.Hash) (int, error) {
	value, err := self.get(blockHash, KeyAccuVerifier(proposalId))
	if err != nil {
		return 0, err
	}

	return len(value), nil
}
