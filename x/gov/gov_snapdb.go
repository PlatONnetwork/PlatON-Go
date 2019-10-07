package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

func get(blockHash common.Hash, key []byte) ([]byte, error) {
	return snapshotdb.Instance().Get(blockHash, key)
}

func put(blockHash common.Hash, key []byte, value interface{}) error {
	bytes, err := rlp.EncodeToBytes(value)
	if err != nil {
		return err
	}
	return snapshotdb.Instance().Put(blockHash, key, bytes)
}

func del(blockHash common.Hash, key []byte) error {
	return snapshotdb.Instance().Del(blockHash, key)
}

func addProposalByKey(blockHash common.Hash, key []byte, proposalId common.Hash) error {
	proposalIDList, err := getProposalIDListByKey(blockHash, key)
	if err != nil {
		return err
	}

	for _, pID := range proposalIDList {
		if pID == proposalId {
			return nil
		}
	}
	proposalIDList = append(proposalIDList, proposalId)
	return put(blockHash, key, proposalIDList)
}

func getVotingIDList(blockHash common.Hash) ([]common.Hash, error) {
	return getProposalIDListByKey(blockHash, KeyVotingProposals())
}

func getPreActiveProposalID(blockHash common.Hash) (common.Hash, error) {
	//return self.getProposalIDListByKey(blockHash, KeyPreActiveProposals())
	bytes, err := get(blockHash, KeyPreActiveProposal())

	if err != nil && err != snapshotdb.ErrNotFound {
		return common.Hash{}, err
	}

	var proposalID common.Hash
	if len(bytes) > 0 {
		if err = rlp.DecodeBytes(bytes, &proposalID); err != nil {
			return common.Hash{}, err
		}
	}
	return proposalID, nil

}

func getEndIDList(blockHash common.Hash) ([]common.Hash, error) {
	return getProposalIDListByKey(blockHash, KeyEndProposals())
}

func getProposalIDListByKey(blockHash common.Hash, key []byte) ([]common.Hash, error) {
	bytes, err := get(blockHash, key)
	if err != nil && err != snapshotdb.ErrNotFound {
		return nil, err
	}
	var idList []common.Hash
	if len(bytes) > 0 {
		if err = rlp.DecodeBytes(bytes, &idList); err != nil {
			return nil, err
		}
	}
	return idList, nil
}

func getAllProposalIDList(blockHash common.Hash) ([]common.Hash, error) {
	var total []common.Hash

	proposalIDList, err := getVotingIDList(blockHash)
	if err != nil {
		log.Error("list voting proposal IDs failed", "blockHash", blockHash)
		return nil, err
	} else if len(proposalIDList) > 0 {
		total = append(total, proposalIDList...)
	}

	proposalID, err := getPreActiveProposalID(blockHash)
	if err != nil {
		log.Error("list pre-active proposal IDs failed", "blockHash", blockHash)
		return nil, err
	} else if proposalID != common.ZeroHash {
		total = append(total, proposalID)
	}
	proposalIDList, err = getEndIDList(blockHash)
	if err != nil {
		log.Error("list end proposal IDs failed", "blockHash", blockHash)
		return nil, err
	} else if len(proposalIDList) > 0 {
		total = append(total, proposalIDList...)
	}

	return total, nil
}

func addActiveNode(blockHash common.Hash, node discover.NodeID, proposalId common.Hash) error {
	nodes, err := getActiveNodeList(blockHash, proposalId)
	if err != nil && err != snapshotdb.ErrNotFound {
		return err
	}

	//distinct the nodeID
	if xutil.InNodeIDList(node, nodes) {
		return nil
	} else {
		nodes = append(nodes, node)
		return put(blockHash, KeyActiveNodes(proposalId), nodes)
	}
}

func getActiveNodeList(blockHash common.Hash, proposalId common.Hash) ([]discover.NodeID, error) {
	value, err := get(blockHash, KeyActiveNodes(proposalId))
	if err != nil && err != snapshotdb.ErrNotFound {
		return nil, err
	}
	var nodes []discover.NodeID
	if len(value) > 0 {
		if err := rlp.DecodeBytes(value, &nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func deleteActiveNodeList(blockHash common.Hash, proposalId common.Hash) error {
	return del(blockHash, KeyActiveNodes(proposalId))
}

func addAccuVerifiers(blockHash common.Hash, proposalId common.Hash, nodes []discover.NodeID) error {
	value, err := get(blockHash, KeyAccuVerifier(proposalId))
	if err != nil && err != snapshotdb.ErrNotFound {
		return err
	}
	var accuVerifiers []discover.NodeID

	if value != nil {
		if err := rlp.DecodeBytes(value, &accuVerifiers); err != nil {
			return err
		}
	}

	existMap := make(map[discover.NodeID]struct{}, len(accuVerifiers))
	for _, nodeID := range accuVerifiers {
		existMap[nodeID] = struct{}{}
	}

	for _, nodeID := range nodes {
		if _, ok := existMap[nodeID]; !ok {
			accuVerifiers = append(accuVerifiers, nodeID)
		}
		/*
			if !xutil.InNodeIDList(nodeID, accuVerifiers) {
				accuVerifiers = append(accuVerifiers, nodeID)
			}
		*/
	}
	log.Debug("accumulated verifiers", "proposalID", proposalId, "total", len(accuVerifiers))
	return put(blockHash, KeyAccuVerifier(proposalId), accuVerifiers)
}

func getAccuVerifiers(blockHash common.Hash, proposalId common.Hash) ([]discover.NodeID, error) {
	value, err := get(blockHash, KeyAccuVerifier(proposalId))
	if err != nil && err != snapshotdb.ErrNotFound {
		return nil, err
	}

	if len(value) > 0 {
		var verifiers []discover.NodeID
		if err := rlp.DecodeBytes(value, &verifiers); err != nil {
			return nil, err
		} else {
			return verifiers, nil
		}
	}
	return nil, nil
}
