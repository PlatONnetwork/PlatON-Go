package cbft

import (
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func newValidators(nodes []discover.Node, validBlockNumber uint64) *cbfttypes.Validators {
	vds := &cbfttypes.Validators{
		Nodes:            make(cbfttypes.ValidateNodeMap, len(nodes)),
		ValidBlockNumber: validBlockNumber,
	}

	for i, node := range nodes {
		pubkey, err := node.ID.Pubkey()
		if err != nil {
			panic(err)
		}

		vds.Nodes[node.ID] = &cbfttypes.ValidateNode{
			Index:   i,
			Address: crypto.PubkeyToAddress(*pubkey),
			PubKey:  pubkey,
		}
	}
	return vds
}

// Agency
type Agency interface {
	Sign(msg interface{}) error
	VerifySign(msg interface{}) error
	Flush(header *types.Header) error
	VerifyHeader(header *types.Header) error
	GetLastNumber(blockNumber uint64) uint64
	GetValidator(blockNumber uint64) (*cbfttypes.Validators, error)
	IsCandidateNode(nodeID discover.NodeID) bool
}

type StaticAgency struct {
	Agency

	validators *cbfttypes.Validators
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

func (d *StaticAgency) Flush(header *types.Header) error {
	return nil
}

func (d *StaticAgency) VerifyHeader(*types.Header) error {
	return nil
}

func (d *StaticAgency) GetLastNumber(blockNumber uint64) uint64 {
	return 0
}

func (d *StaticAgency) GetValidator(uint64) (*cbfttypes.Validators, error) {
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
	defaultValidators     *cbfttypes.Validators
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

func (ia *InnerAgency) Flush(header *types.Header) error {
	return nil
}

func (ia *InnerAgency) VerifyHeader(*types.Header) error {
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

func (ia *InnerAgency) GetValidator(blockNumber uint64) (v *cbfttypes.Validators, err error) {
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
	b := state.GetState(cvm.ValidatorInnerContractAddr, []byte(vm.CurrentValidatorKey))
	var vds vm.Validators
	err = rlp.DecodeBytes(b, &vds)
	if err != nil {
		log.Error("RLP decode fail, use default validators", "number", block.Number(), "error", err)
		return ia.defaultValidators, nil
	}
	var validators cbfttypes.Validators
	validators.Nodes = make(cbfttypes.ValidateNodeMap, len(vds.ValidateNodes))
	for _, node := range vds.ValidateNodes {
		pubkey, _ := node.NodeID.Pubkey()
		validators.Nodes[node.NodeID] = &cbfttypes.ValidateNode{
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
