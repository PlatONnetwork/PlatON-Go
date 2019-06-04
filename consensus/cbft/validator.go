package cbft

import (
	"bytes"
	"crypto/elliptic"
	"errors"

	"fmt"

	"crypto/ecdsa"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

type ValidateNode struct {
	Index   int            `json:"index"`
	Address common.Address `json:"-"`
	PubKey  *ecdsa.PublicKey
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

type ValidateNodeMap map[discover.NodeID]*ValidateNode

func (vnm ValidateNodeMap) String() string {
	s := ""
	for k, v := range vnm {
		s = s + fmt.Sprintf("{%s:%s},", k, v)
	}
	return s
}

type Validators struct {
	Nodes            ValidateNodeMap `json:"validateNodes"`
	ValidBlockNumber uint64          `json:"-"`
}

func newValidators(nodes []discover.Node, validBlockNumber uint64) *Validators {
	vds := &Validators{
		Nodes:            make(ValidateNodeMap, len(nodes)),
		ValidBlockNumber: validBlockNumber,
	}

	for i, node := range nodes {
		pubkey, err := node.ID.Pubkey()
		if err != nil {
			panic(err)
		}

		vds.Nodes[node.ID] = &ValidateNode{
			Index:   i,
			Address: crypto.PubkeyToAddress(*pubkey),
			PubKey:  pubkey,
		}
	}
	return vds
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

// Agency
type Agency interface {
	Sign(msg interface{}) error
	VerifySign(msg interface{}) error
	GetLastNumber(blockNumber uint64) uint64
	GetValidator(blockNumber uint64) (*Validators, error)
	IsCandidateNode(nodeID discover.NodeID) bool
}

type StaticAgency struct {
	Agency

	validators *Validators
}

func NewStaticAgency(nodes []discover.Node) Agency {
	return &StaticAgency{
		validators: newValidators(nodes, 0),
	}
}

func (d *StaticAgency) Sign(interface{}) error {
	return nil
}

func (d *StaticAgency) VerifySign(interface{}) error {
	return nil
}

func (d *StaticAgency) GetLastNumber(blockNumber uint64) uint64 {
	return 0
}

func (d *StaticAgency) GetValidator(uint64) (*Validators, error) {
	return d.validators, nil
}

func (d *StaticAgency) IsCandidateNode(nodeID discover.NodeID) bool {
	return false
}

type InnerAgency struct {
	Agency

	blocksPerNode         uint64
	defaultBlocksPerRound uint64
	offset                uint64
	blockchain            *core.BlockChain
	defaultValidators     *Validators
}

func NewInnerAgency(nodes []discover.Node, chain *core.BlockChain, blocksPerNode, offset int) Agency {
	return &InnerAgency{
		blocksPerNode:         uint64(blocksPerNode),
		defaultBlocksPerRound: uint64(len(nodes) * blocksPerNode),
		offset:                uint64(offset),
		blockchain:            chain,
		defaultValidators:     newValidators(nodes, 0),
	}
}

func (ia *InnerAgency) Sign(interface{}) error {
	return nil
}

func (ia *InnerAgency) VerifySign(interface{}) error {
	return nil
}

func (ia *InnerAgency) GetLastNumber(blockNumber uint64) uint64 {
	var lastBlockNumber uint64
	if blockNumber <= ia.defaultBlocksPerRound {
		lastBlockNumber = ia.defaultBlocksPerRound
	} else {
		vds, err := ia.GetValidator(blockNumber)
		if err != nil {
			log.Error("Get validator fail", "blockNumber", blockNumber)
			return 0
		}

		if vds.ValidBlockNumber == 0 && blockNumber%ia.defaultBlocksPerRound == 0 {
			return blockNumber
		}

		// lastNumber = vds.ValidBlockNumber + ia.blocksPerNode * vds.Len() - 1
		lastBlockNumber = vds.ValidBlockNumber + ia.blocksPerNode*uint64(vds.Len()) - 1

		// May be `CurrentValidators ` had not updated, so we need to calcuate `lastBlockNumber`
		// via `blockNumber`.
		if lastBlockNumber < blockNumber {
			blocksPerRound := ia.blocksPerNode * uint64(vds.Len())
			baseNum := blockNumber - (blockNumber % blocksPerRound)
			lastBlockNumber = baseNum + blocksPerRound
		}
	}
	log.Debug("Get last block number", "blockNumber", blockNumber, "lastBlockNumber", lastBlockNumber)
	return lastBlockNumber
}

func (ia *InnerAgency) GetValidator(blockNumber uint64) (v *Validators, err error) {
	//var lastBlockNumber uint64
	/*
		defer func() {
			log.Trace("Get validator",
				"lastBlockNumber", lastBlockNumber,
				"blocksPerNode", ia.blocksPerNode,
				"blockNumber", blockNumber,
				"validators", v,
				"error", err)
		}()*/

	if blockNumber <= ia.defaultBlocksPerRound {
		return ia.defaultValidators, nil
	}

	// Otherwise, get validators from inner contract.
	vdsCftNum := blockNumber - ia.offset - 1
	block := ia.blockchain.GetBlockByNumber(vdsCftNum)
	if block == nil {
		log.Error("Get the block fail, use default validators", "number", vdsCftNum)
		return ia.defaultValidators, nil
	}
	state, err := ia.blockchain.StateAt(block.Root())
	if err != nil {
		log.Error("Get the state fail, use default validators", "number", block.Number(), "hash", block.Hash(), "error", err)
		return ia.defaultValidators, nil
	}
	b := state.GetState(vm.ValidatorInnerContractAddr, []byte(vm.CurrentValidatorKey))
	var vds vm.Validators
	err = rlp.DecodeBytes(b, &vds)
	if err != nil {
		log.Error("RLP decode fail, use default validators", "number", block.Number(), "error", err)
		return ia.defaultValidators, nil
	}
	var validators Validators
	validators.Nodes = make(ValidateNodeMap, len(vds.ValidateNodes))
	for _, node := range vds.ValidateNodes {
		pubkey, _ := node.NodeID.Pubkey()
		validators.Nodes[node.NodeID] = &ValidateNode{
			Index:   int(node.Index),
			Address: node.Address,
			PubKey:  pubkey,
		}
	}
	validators.ValidBlockNumber = vds.ValidBlockNumber
	//lastBlockNumber = vds.ValidBlockNumber + ia.blocksPerNode*uint64(validators.Len()) - 1
	return &validators, nil
}

func (ia *InnerAgency) IsCandidateNode(nodeID discover.NodeID) bool {
	return true
}
