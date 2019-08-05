package cbfttypes

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sort"

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

type UpdateChainState func(qc *types.Block, lock *types.Block, commit *types.Block)

type ChainStateResult struct {
	QC       *types.Block
	Lock     *types.Block
	Commit   *types.Block
	Callback UpdateChainState
}

type CbftResult struct {
	Block            *types.Block
	ExtraData        []byte
	SyncState        chan error
	ChainStateResult *ChainStateResult
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
	Index     int            `json:"index"`
	Address   common.Address `json:"-"`
	PubKey    *ecdsa.PublicKey
	NodeID    discover.NodeID
	BlsPubKey *bls.PublicKey
}

type ValidateNodeMap map[discover.NodeID]*ValidateNode

type SortedValidatorNode []*ValidateNode

func (sv SortedValidatorNode) Len() int           { return len(sv) }
func (sv SortedValidatorNode) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv SortedValidatorNode) Less(i, j int) bool { return sv[i].Index < sv[j].Index }

type Validators struct {
	Nodes            ValidateNodeMap `json:"validateNodes"`
	ValidBlockNumber uint64          `json:"-"`

	sortedNodes SortedValidatorNode
}

func (vn *ValidateNode) String() string {
	return fmt.Sprintf("{Index:%d Address:%s BlsPubKey:%s}", vn.Index, vn.Address.String(), fmt.Sprintf("%x", vn.BlsPubKey.Serialize()))
}

func (vn *ValidateNode) Verify(data, sign []byte) bool {
	var sig bls.Sign
	err := sig.Deserialize(sign)
	if err != nil {
		return false
	}

	return sig.Verify(vn.BlsPubKey, string(data))
}

func (vnm ValidateNodeMap) String() string {
	s := ""
	for k, v := range vnm {
		s = s + fmt.Sprintf("{%s:%s},", k, v)
	}
	return s
}

func (vs *Validators) String() string {
	return fmt.Sprintf("{Nodes:[%s] ValidBlockNumber:%d}", vs.Nodes, vs.ValidBlockNumber)
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

func (vs *Validators) FindNodeByAddress(addr common.Address) (*ValidateNode, error) {
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

func (vs *Validators) Index(nodeID discover.NodeID) (int, error) {
	if node, ok := vs.Nodes[nodeID]; ok {
		return node.Index, nil
	}
	return -1, errors.New("not found the specified validator")
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
