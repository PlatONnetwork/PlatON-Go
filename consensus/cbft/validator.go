package cbft

import (
	"errors"
	"math/big"

	"bytes"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// Valiadator event
type UpdateValidatorEvent struct{}
type StopConsensusEvent struct{}

type ValidateNode struct {
	Index   int
	Address common.Address
}

type ValidateNodeMap map[discover.NodeID]*ValidateNode

type Validators struct {
	nodes            ValidateNodeMap
	startTimeOfEpoch uint64 // second
}

func newValidators(nodes []discover.Node, startTime uint64) *Validators {
	vds := &Validators{
		nodes:            make(ValidateNodeMap, len(nodes)),
		startTimeOfEpoch: startTime,
	}

	for i, node := range nodes {
		pubkey, err := node.ID.Pubkey()
		if err != nil {
			panic(err)
		}

		vds.nodes[node.ID] = &ValidateNode{
			Index:   i,
			Address: crypto.PubkeyToAddress(*pubkey),
		}
	}
	return vds
}

func (vs *Validators) NodeList() []discover.NodeID {
	nodeList := make([]discover.NodeID, len(vs.nodes))
	for id, _ := range vs.nodes {
		nodeList = append(nodeList, id)
	}
	return nodeList
}

func (vs *Validators) NodeIndexAddress(id discover.NodeID) (int, common.Address, error) {
	node, ok := vs.nodes[id]
	if ok {
		return node.Index, node.Address, nil
	}
	return -1, common.Address{}, errors.New("not found the node")
}

func (vs *Validators) NodeID(idx int) discover.NodeID {
	for id, node := range vs.nodes {
		if node.Index == idx {
			return id
		}
	}
	// I think never run here ^_^
	return discover.NodeID{}
}

func (vs *Validators) AddressIndex(addr common.Address) (int, error) {
	for _, node := range vs.nodes {
		if bytes.Equal(node.Address[:], addr[:]) {
			return node.Index, nil
		}
	}
	return -1, errors.New("invalid address")
}

func (vs *Validators) NodeIndex(id discover.NodeID) (int, error) {
	for nodeID, node := range vs.nodes {
		if nodeID == id {
			return node.Index, nil
		}
	}
	return -1, errors.New("not found the node")
}

func (vs *Validators) StartTimeOfEpoch() int64 {
	return int64(vs.startTimeOfEpoch)
}

func (vs *Validators) Len() int {
	return len(vs.nodes)
}

// Agency
type Agency interface {
	Sign(msg interface{}) error
	VerifySign(msg interface{}) error
	GetValidator(blockNumber *big.Int) (*Validators, error)
}

type StaticAgency struct {
	Agency

	validators *Validators
}

func NewStaticAgency(nodes []discover.Node, startTime uint64) Agency {
	agency := &StaticAgency{
		validators: newValidators(nodes, startTime),
	}
	return agency
}

func (d *StaticAgency) Sign(interface{}) error {
	return nil
}

func (d *StaticAgency) VerifySign(interface{}) error {
	return nil
}

func (d *StaticAgency) GetValidator(*big.Int) (*Validators, error) {
	return d.validators, nil
}

type InnerAgency struct {
	Agency

	defaultValidators *Validators
}

func NewInnerAgency(nodes []discover.Node, startTime uint64) Agency {
	agency := &InnerAgency{
		defaultValidators: newValidators(nodes, startTime),
	}

	return agency
}

func (ia *InnerAgency) Sign(interface{}) error {
	return nil
}

func (ia *InnerAgency) VerifySign(interface{}) error {
	return nil
}

func (ia *InnerAgency) GetValidator(blockNumber *big.Int) (*Validators, error) {
	if blockNumber.Cmp(big.NewInt(0)) == 0 {
		return ia.defaultValidators, nil
	}

	// Otherwise, get validators from inner contract.
	// TODO: Get validator from inner contract.
	return nil, nil
}
