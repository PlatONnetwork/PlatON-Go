package cbfttypes

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
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

type CbftResult struct {
	Block     *types.Block
	ExtraData []byte
	SyncState chan error
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

type UpdateValidatorEvent struct{}

type ValidateNode struct {
	Index   int            `json:"index"`
	Address common.Address `json:"-"`
	PubKey  *ecdsa.PublicKey
}

type ValidateNodeMap map[discover.NodeID]*ValidateNode

type Validators struct {
	Nodes            ValidateNodeMap `json:"validateNodes"`
	ValidBlockNumber uint64          `json:"-"`
}

func (vn *ValidateNode) String() string {
	return fmt.Sprintf("{Index:%d Address:%s}", vn.Index, vn.Address.String())
}

func (vn *ValidateNode) Verify(data, sign []byte) bool {
	recPubKey, err := crypto.Ecrecover(data, sign)
	if err != nil {
		return false
	}

	pbytes := elliptic.Marshal(vn.PubKey.Curve, vn.PubKey.X, vn.PubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return false
	}
	return true
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

func (vs *Validators) NodeIndexAddress(id discover.NodeID) (*ValidateNode, error) {
	node, ok := vs.Nodes[id]
	if ok {
		return node, nil
	}
	return nil, errors.New("not found the node")
}

func (vs *Validators) NodeID(idx int) discover.NodeID {
	for id, node := range vs.Nodes {
		if node.Index == idx {
			return id
		}
	}
	// I think never run here ^_^
	return discover.NodeID{}
}

func (vs *Validators) AddressIndex(addr common.Address) (*ValidateNode, error) {
	for _, node := range vs.Nodes {
		if bytes.Equal(node.Address[:], addr[:]) {
			return node, nil
		}
	}
	return nil, errors.New("invalid address")
}

func (vs *Validators) NodeIndex(id discover.NodeID) (*ValidateNode, error) {
	for nodeID, node := range vs.Nodes {
		if nodeID == id {
			return node, nil
		}
	}
	return nil, errors.New("not found the node")
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
