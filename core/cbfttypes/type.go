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

package cbfttypes

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"math"
	"math/big"
	"sort"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// Block's Signature info
type BlockSignature struct {
	SignHash  common.Hash // Signature hash，header[0:32]
	Hash      common.Hash // Block hash，header[:]
	Number    *big.Int
	Signature *common.BlockConfirmSign
}

func (bs *BlockSignature) Copy() *BlockSignature {
	sign := *bs.Signature
	return &BlockSignature{
		SignHash:  bs.SignHash,
		Hash:      bs.Hash,
		Number:    new(big.Int).Set(bs.Number),
		Signature: &sign,
	}
}

type UpdateChainStateFn func(qcState, lockState, commitState *protocols.State)

type CbftResult struct {
	Block              *types.Block
	ExtraData          []byte
	SyncState          chan error
	ChainStateUpdateCB func()
}

type ProducerState struct {
	count int
	miner common.Address
}

func (ps *ProducerState) Add(miner common.Address) {
	if ps.miner == miner {
		ps.count++
	} else {
		ps.miner = miner
		ps.count = 1
	}
}

func (ps *ProducerState) Get() (common.Address, int) {
	return ps.miner, ps.count
}

func (ps *ProducerState) Validate(period int) bool {
	return ps.count < period
}

type AddValidatorEvent struct {
	NodeID discover.NodeID
}

type RemoveValidatorEvent struct {
	NodeID discover.NodeID
}

type UpdateValidatorEvent struct{}

type ValidateNode struct {
	Index     uint32             `json:"index"`
	Address   common.NodeAddress `json:"address"`
	PubKey    *ecdsa.PublicKey   `json:"-"`
	NodeID    discover.NodeID    `json:"nodeID"`
	BlsPubKey *bls.PublicKey     `json:"blsPubKey"`
}

type ValidateNodeMap map[discover.NodeID]*ValidateNode

type SortedValidatorNode []*ValidateNode

func (sv SortedValidatorNode) Len() int           { return len(sv) }
func (sv SortedValidatorNode) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv SortedValidatorNode) Less(i, j int) bool { return sv[i].Index < sv[j].Index }

type Validators struct {
	Nodes            ValidateNodeMap `json:"validateNodes"`
	ValidBlockNumber uint64          `json:"validateBlockNumber"`

	sortedNodes SortedValidatorNode
}

func (vn *ValidateNode) String() string {
	b, _ := json.Marshal(vn)
	return string(b)
}

func (vn *ValidateNode) Verify(data, sign []byte) error {
	var sig bls.Sign
	err := sig.Deserialize(sign)
	if err != nil {
		return err
	}

	if !sig.Verify(vn.BlsPubKey, string(data)) {
		return errors.New(fmt.Sprintf("bls verifies signature fail, data:%s, sign:%s, pubkey:%s", hexutil.Encode(data), hexutil.Encode(sign), hexutil.Encode(vn.BlsPubKey.Serialize())))
	}
	return nil
}

func (vnm ValidateNodeMap) String() string {
	s := ""
	for k, v := range vnm {
		s = s + fmt.Sprintf("{%s:%s},", k, v)
	}
	return s
}

func (vs *Validators) String() string {
	b, _ := json.Marshal(vs)
	return string(b)
}

func (vs *Validators) NodeList() []discover.NodeID {
	nodeList := make([]discover.NodeID, 0)
	for id, _ := range vs.Nodes {
		nodeList = append(nodeList, id)
	}
	return nodeList
}

func (vs *Validators) NodeListByIndexes(indexes []uint32) ([]*ValidateNode, error) {
	if len(vs.sortedNodes) == 0 {
		vs.sort()
	}
	l := make([]*ValidateNode, 0)
	for _, index := range indexes {
		if int(index) >= len(vs.sortedNodes) {
			return nil, errors.New("invalid index")
		}
		l = append(l, vs.sortedNodes[int(index)])
	}
	return l, nil
}

func (vs *Validators) NodeListByBitArray(vSet *utils.BitArray) ([]*ValidateNode, error) {
	if len(vs.sortedNodes) == 0 {
		vs.sort()
	}
	l := make([]*ValidateNode, 0)

	for index := uint32(0); index < vSet.Size(); index++ {
		if vSet.GetIndex(index) {
			if int(index) >= len(vs.sortedNodes) {
				return nil, errors.New("invalid index")
			}
			l = append(l, vs.sortedNodes[int(index)])
		}
	}
	return l, nil
}

func (vs *Validators) FindNodeByID(id discover.NodeID) (*ValidateNode, error) {
	node, ok := vs.Nodes[id]
	if ok {
		return node, nil
	}
	return nil, errors.New("not found the node")
}

func (vs *Validators) FindNodeByIndex(index int) (*ValidateNode, error) {
	if len(vs.sortedNodes) == 0 {
		vs.sort()
	}
	if index >= len(vs.sortedNodes) {
		return nil, errors.New("not found the specified validator")
	} else {
		return vs.sortedNodes[index], nil
	}
}

func (vs *Validators) FindNodeByAddress(addr common.NodeAddress) (*ValidateNode, error) {
	for _, node := range vs.Nodes {
		if bytes.Equal(node.Address[:], addr[:]) {
			return node, nil
		}
	}
	return nil, errors.New("invalid address")
}

func (vs *Validators) NodeID(idx int) discover.NodeID {
	if len(vs.sortedNodes) == 0 {
		vs.sort()
	}
	if idx >= vs.sortedNodes.Len() {
		return discover.NodeID{}
	}
	return vs.sortedNodes[idx].NodeID
}

func (vs *Validators) Index(nodeID discover.NodeID) (uint32, error) {
	if node, ok := vs.Nodes[nodeID]; ok {
		return node.Index, nil
	}
	return math.MaxUint32, errors.New("not found the specified validator")
}

func (vs *Validators) Len() int {
	return len(vs.Nodes)
}

func (vs *Validators) Equal(rsh *Validators) bool {
	if vs.Len() != rsh.Len() {
		return false
	}

	equal := true
	for k, v := range vs.Nodes {
		if vv, ok := rsh.Nodes[k]; !ok || vv.Index != v.Index {
			equal = false
			break
		}
	}
	return equal
}

func (vs *Validators) sort() {
	for _, node := range vs.Nodes {
		vs.sortedNodes = append(vs.sortedNodes, node)
	}
	sort.Sort(vs.sortedNodes)
}
