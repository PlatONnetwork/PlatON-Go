package cbft

import (
	"errors"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"bytes"
)

type Node struct {
	Index   int
	Address common.Address
}

type Validators struct {
	nodes            map[discover.NodeID]*Node
	startTimeOfEpoch uint64 // second
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

// Validator
type Validator interface {
	Sign(msg interface{}) error
	VerifySign(msg interface{}) error
	GetValidator() (*Validators, error)
}

type DefaultValidator struct {
	Validator

	validators Validators
}

func NewDefaultValidator(nodes []discover.Node, startTime uint64) Validator {
	validator := &DefaultValidator{
		validators: Validators{
			nodes:            make(map[discover.NodeID]*Node, len(nodes)),
			startTimeOfEpoch: startTime,
		},
	}
	for i, node := range nodes {
		pubkey, err := node.ID.Pubkey()
		if err != nil {
			panic(err)
		}

		validator.validators.nodes[node.ID] = &Node{
			Index:   i,
			Address: crypto.PubkeyToAddress(*pubkey),
		}
	}

	return validator
}

func (v *DefaultValidator) Sign(interface{}) error {
	return nil
}

func (v *DefaultValidator) VerifySign(msg interface{}) error {
	return nil
}

func (v *DefaultValidator) GetValidator() (*Validators, error) {
	return &v.validators, nil
}
