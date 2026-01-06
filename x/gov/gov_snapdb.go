// Copyright 2021 The PlatON Network Authors
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

package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

func (gd *GovDB) get(blockHash common.Hash, key []byte) ([]byte, error) {
	return gd.db.Get(blockHash, key)
}

func (gd *GovDB) put(blockHash common.Hash, key []byte, value interface{}) error {
	bytes, err := rlp.EncodeToBytes(value)
	if err != nil {
		return err
	}
	return gd.db.Put(blockHash, key, bytes)
}

func (gd *GovDB) del(blockHash common.Hash, key []byte) error {
	return gd.db.Del(blockHash, key)
}

func (gd *GovDB) addProposalByKey(blockHash common.Hash, key []byte, proposalId common.Hash) error {
	proposalIDList, err := gd.getProposalIDListByKey(blockHash, key)
	if err != nil {
		return err
	}

	for _, pID := range proposalIDList {
		if pID == proposalId {
			return nil
		}
	}
	proposalIDList = append(proposalIDList, proposalId)
	return gd.put(blockHash, key, proposalIDList)
}

func (gd *GovDB) getVotingIDList(blockHash common.Hash) ([]common.Hash, error) {
	return gd.getProposalIDListByKey(blockHash, KeyVotingProposals())
}

// Set pre-active version
func (gd *GovDB) setPreActiveVersion(blockHash common.Hash, preActiveVersion uint32) error {
	return gd.put(blockHash, KeyPreActiveVersion(), preActiveVersion)
}

// Get pre-active version
func (gs *GovDB) getPreActiveVersion(blockHash common.Hash) uint32 {
	bytes, err := gs.get(blockHash, KeyPreActiveVersion())
	if snapshotdb.NonDbNotFoundErr(err) {
		return uint32(0)
	}

	var activeVersion uint32
	if len(bytes) > 0 {
		if err = rlp.DecodeBytes(bytes, &activeVersion); err != nil {
			return uint32(0)
		}
	}
	return activeVersion
}

func (gd *GovDB) delPreActiveVersion(blockHash common.Hash) error {
	return gd.del(blockHash, KeyPreActiveVersion())
}

func (gd *GovDB) getPreActiveProposalID(blockHash common.Hash) (common.Hash, error) {
	//return self.getProposalIDListByKey(blockHash, KeyPreActiveProposals())
	bytes, err := gd.get(blockHash, KeyPreActiveProposal())

	if snapshotdb.NonDbNotFoundErr(err) {
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

func (gd *GovDB) getEndIDList(blockHash common.Hash) ([]common.Hash, error) {
	return gd.getProposalIDListByKey(blockHash, KeyEndProposals())
}

func (gd *GovDB) getProposalIDListByKey(blockHash common.Hash, key []byte) ([]common.Hash, error) {
	bytes, err := gd.get(blockHash, key)
	if snapshotdb.NonDbNotFoundErr(err) {
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

func (gd *GovDB) getAllProposalIDList(blockHash common.Hash) ([]common.Hash, error) {
	var total []common.Hash

	proposalIDList, err := gd.getVotingIDList(blockHash)
	if err != nil {
		log.Error("list voting proposal IDs failed", "blockHash", blockHash)
		return nil, err
	} else if len(proposalIDList) > 0 {
		total = append(total, proposalIDList...)
	}

	proposalID, err := gd.getPreActiveProposalID(blockHash)
	if err != nil {
		log.Error("list pre-active proposal IDs failed", "blockHash", blockHash)
		return nil, err
	} else if proposalID != common.ZeroHash {
		total = append(total, proposalID)
	}
	proposalIDList, err = gd.getEndIDList(blockHash)
	if err != nil {
		log.Error("list end proposal IDs failed", "blockHash", blockHash)
		return nil, err
	} else if len(proposalIDList) > 0 {
		total = append(total, proposalIDList...)
	}

	return total, nil
}

func (gd *GovDB) addActiveNode(blockHash common.Hash, node enode.IDv0, proposalId common.Hash) error {
	nodes, err := gd.getActiveNodeList(blockHash, proposalId)
	if snapshotdb.NonDbNotFoundErr(err) {
		return err
	}

	//distinct the nodeID
	if xutil.InNodeIDList(node, nodes) {
		return nil
	} else {
		nodes = append(nodes, node)
		return gd.put(blockHash, KeyActiveNodes(proposalId), nodes)
	}
}

func (gd *GovDB) getActiveNodeList(blockHash common.Hash, proposalId common.Hash) ([]enode.IDv0, error) {
	value, err := gd.get(blockHash, KeyActiveNodes(proposalId))
	if snapshotdb.NonDbNotFoundErr(err) {
		return nil, err
	}
	var nodes []enode.IDv0
	if len(value) > 0 {
		if err := rlp.DecodeBytes(value, &nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (gd *GovDB) deleteActiveNodeList(blockHash common.Hash, proposalId common.Hash) error {
	return gd.del(blockHash, KeyActiveNodes(proposalId))
}

func (gd *GovDB) addAccuVerifiers(blockHash common.Hash, proposalId common.Hash, nodes []enode.IDv0) error {
	value, err := gd.get(blockHash, KeyAccuVerifier(proposalId))
	if snapshotdb.NonDbNotFoundErr(err) {
		return err
	}
	var accuVerifiers []enode.IDv0

	if value != nil {
		if err := rlp.DecodeBytes(value, &accuVerifiers); err != nil {
			return err
		}
	}

	existMap := make(map[enode.IDv0]struct{}, len(accuVerifiers))
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
	return gd.put(blockHash, KeyAccuVerifier(proposalId), accuVerifiers)
}

func (gd *GovDB) getAccuVerifiers(blockHash common.Hash, proposalId common.Hash) ([]enode.IDv0, error) {
	value, err := gd.get(blockHash, KeyAccuVerifier(proposalId))
	if snapshotdb.NonDbNotFoundErr(err) {
		return nil, err
	}

	if len(value) > 0 {
		var verifiers []enode.IDv0
		if err := rlp.DecodeBytes(value, &verifiers); err != nil {
			return nil, err
		} else {
			return verifiers, nil
		}
	}
	return nil, nil
}

func (gd *GovDB) delAccuVerifiers(blockHash common.Hash, proposalId common.Hash) error {
	return gd.del(blockHash, KeyAccuVerifier(proposalId))
}

func (gd *GovDB) addGovernParam(module, name, desc string, paramValue *ParamValue, blockHash common.Hash) error {
	itemList, err := gd.listGovernParamItem("", blockHash)
	if err != nil {
		return nil
	}
	itemList = append(itemList, &ParamItem{module, name, desc})
	if err := gd.put(blockHash, keyPrefixParamItems, itemList); err != nil {
		return err
	}

	if err := gd.put(blockHash, KeyParamValue(module, name), paramValue); err != nil {
		return err
	}
	return nil
}

func (gd *GovDB) findGovernParamValue(module, name string, blockHash common.Hash) (*ParamValue, error) {
	value, err := gd.get(blockHash, KeyParamValue(module, name))
	if snapshotdb.NonDbNotFoundErr(err) {
		return nil, err
	}

	if len(value) > 0 {
		var paramValue ParamValue
		if err := rlp.DecodeBytes(value, &paramValue); err != nil {
			return nil, err
		} else {
			return &paramValue, nil
		}
	}
	return nil, nil
}

func (gd *GovDB) updateGovernParamValue(module, name, newValue string, activeBlock uint64, blockHash common.Hash) error {
	value, err := gd.get(blockHash, KeyParamValue(module, name))
	if snapshotdb.NonDbNotFoundErr(err) {
		return err
	}
	if len(value) > 0 {
		var paramValue ParamValue
		if err := rlp.DecodeBytes(value, &paramValue); err != nil {
			return err
		}
		paramValue.StaleValue = paramValue.Value
		paramValue.Value = newValue
		paramValue.ActiveBlock = activeBlock

		if err := gd.put(blockHash, KeyParamValue(module, name), paramValue); err != nil {
			return err
		}
		return nil
	}
	return UnsupportedGovernParam
}

func (gd *GovDB) listGovernParam(module string, blockHash common.Hash) ([]*GovernParam, error) {
	itemList, err := gd.listGovernParamItem(module, blockHash)
	if err != nil {
		return nil, err
	}
	var paraList []*GovernParam
	for _, item := range itemList {
		if value, err := gd.findGovernParamValue(item.Module, item.Name, blockHash); err != nil {
			return nil, err
		} else {
			param := &GovernParam{item, value, nil}
			paraList = append(paraList, param)
		}
	}
	return paraList, nil
}

func (gd *GovDB) listGovernParamItem(module string, blockHash common.Hash) ([]*ParamItem, error) {
	itemBytes, err := gd.get(blockHash, KeyParamItems())
	if snapshotdb.NonDbNotFoundErr(err) {
		return nil, err
	}

	if len(itemBytes) > 0 {
		var itemList []*ParamItem
		if err := rlp.DecodeBytes(itemBytes, &itemList); err != nil {
			return nil, err
		}
		if len(module) == 0 {
			return itemList, nil
		} else {
			idx := 0
			for _, item := range itemList {
				if item.Module == module {
					itemList[idx] = item
					idx++
				}
			}
			return itemList[:idx], nil
		}
	}
	return nil, nil
}
