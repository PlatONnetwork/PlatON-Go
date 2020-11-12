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

package validator

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/core/state"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func newValidators(nodes []params.CbftNode, validBlockNumber uint64) *cbfttypes.Validators {
	vds := &cbfttypes.Validators{
		Nodes:            make(cbfttypes.ValidateNodeMap, len(nodes)),
		ValidBlockNumber: validBlockNumber,
	}

	for i, node := range nodes {
		pubkey, err := node.Node.ID.Pubkey()
		if err != nil {
			panic(err)
		}

		blsPubKey := node.BlsPubKey

		vds.Nodes[node.Node.ID] = &cbfttypes.ValidateNode{
			Index:     uint32(i),
			Address:   crypto.PubkeyToNodeAddress(*pubkey),
			PubKey:    pubkey,
			NodeID:    node.Node.ID,
			BlsPubKey: &blsPubKey,
		}
	}
	return vds
}

type StaticAgency struct {
	consensus.Agency

	validators *cbfttypes.Validators
}

func NewStaticAgency(nodes []params.CbftNode) consensus.Agency {
	return &StaticAgency{
		validators: newValidators(nodes, 0),
	}
}

func (d *StaticAgency) Flush(header *types.Header) error {
	return nil
}

func (d *StaticAgency) Sign(interface{}) error {
	return nil
}

func (d *StaticAgency) VerifySign(interface{}) error {
	return nil
}

func (d *StaticAgency) VerifyHeader(header *types.Header, statedb *state.StateDB) error {
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

func (d *StaticAgency) OnCommit(block *types.Block) error {
	return nil
}

type MockAgency struct {
	consensus.Agency
	validators *cbfttypes.Validators
	interval   uint64
}

func NewMockAgency(nodes []params.CbftNode, interval uint64) consensus.Agency {
	return &MockAgency{
		validators: newValidators(nodes, 0),
		interval:   interval,
	}
}

func (d *MockAgency) Flush(header *types.Header) error {
	return nil
}

func (d *MockAgency) Sign(interface{}) error {
	return nil
}

func (d *MockAgency) VerifySign(interface{}) error {
	return nil
}

func (d *MockAgency) VerifyHeader(header *types.Header, statedb *state.StateDB) error {
	return nil
}

func (d *MockAgency) GetLastNumber(blockNumber uint64) uint64 {
	if blockNumber%d.interval == 1 {
		return blockNumber + d.interval - 1
	}
	return 0
}

func (d *MockAgency) GetValidator(blockNumber uint64) (*cbfttypes.Validators, error) {
	if blockNumber > d.interval && blockNumber%d.interval == 1 {
		d.validators.ValidBlockNumber = d.validators.ValidBlockNumber + d.interval + 1
	}
	return d.validators, nil
}

func (d *MockAgency) IsCandidateNode(nodeID discover.NodeID) bool {
	return false
}

func (d *MockAgency) OnCommit(block *types.Block) error {
	return nil
}

type InnerAgency struct {
	consensus.Agency

	blocksPerNode         uint64
	defaultBlocksPerRound uint64
	offset                uint64
	blockchain            *core.BlockChain
	defaultValidators     *cbfttypes.Validators
}

func NewInnerAgency(nodes []params.CbftNode, chain *core.BlockChain, blocksPerNode, offset int) consensus.Agency {
	return &InnerAgency{
		blocksPerNode:         uint64(blocksPerNode),
		defaultBlocksPerRound: uint64(len(nodes) * blocksPerNode),
		offset:                uint64(offset),
		blockchain:            chain,
		defaultValidators:     newValidators(nodes, 1),
	}
}

func (ia *InnerAgency) Flush(header *types.Header) error {
	return nil
}

func (ia *InnerAgency) Sign(interface{}) error {
	return nil
}

func (ia *InnerAgency) VerifySign(interface{}) error {
	return nil
}

func (ia *InnerAgency) VerifyHeader(header *types.Header, stateDB *state.StateDB) error {
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
			if blockNumber%blocksPerRound == 0 {
				lastBlockNumber = blockNumber
			} else {
				baseNum := blockNumber - (blockNumber % blocksPerRound)
				lastBlockNumber = baseNum + blocksPerRound
			}
		}
	}
	//log.Debug("Get last block number", "blockNumber", blockNumber, "lastBlockNumber", lastBlockNumber)
	return lastBlockNumber
}

func (ia *InnerAgency) GetValidator(blockNumber uint64) (v *cbfttypes.Validators, err error) {
	defaultValidators := *ia.defaultValidators
	baseNumber := blockNumber
	if blockNumber == 0 {
		baseNumber = 1
	}
	defaultValidators.ValidBlockNumber = ((baseNumber-1)/ia.defaultBlocksPerRound)*ia.defaultBlocksPerRound + 1
	if blockNumber <= ia.defaultBlocksPerRound {
		return &defaultValidators, nil
	}

	// Otherwise, get validators from inner contract.
	vdsCftNum := blockNumber - ia.offset - 1
	block := ia.blockchain.GetBlockByNumber(vdsCftNum)
	if block == nil {
		log.Error("Get the block fail, use default validators", "number", vdsCftNum)
		return &defaultValidators, nil
	}
	state, err := ia.blockchain.StateAt(block.Root())
	if err != nil {
		log.Error("Get the state fail, use default validators", "number", block.Number(), "hash", block.Hash(), "error", err)
		return &defaultValidators, nil
	}
	b := state.GetState(cvm.ValidatorInnerContractAddr, []byte(vm.CurrentValidatorKey))
	if len(b) == 0 {
		return &defaultValidators, nil
	}
	var vds vm.Validators
	err = rlp.DecodeBytes(b, &vds)
	if err != nil {
		log.Error("RLP decode fail, use default validators", "number", block.Number(), "error", err)
		return &defaultValidators, nil
	}
	var validators cbfttypes.Validators
	validators.Nodes = make(cbfttypes.ValidateNodeMap, len(vds.ValidateNodes))
	for _, node := range vds.ValidateNodes {
		pubkey, _ := node.NodeID.Pubkey()
		blsPubKey := node.BlsPubKey
		validators.Nodes[node.NodeID] = &cbfttypes.ValidateNode{
			Index:     uint32(node.Index),
			Address:   node.Address,
			PubKey:    pubkey,
			NodeID:    node.NodeID,
			BlsPubKey: &blsPubKey,
		}
	}
	validators.ValidBlockNumber = vds.ValidBlockNumber
	return &validators, nil
}

func (ia *InnerAgency) IsCandidateNode(nodeID discover.NodeID) bool {
	return true
}

func (ia *InnerAgency) OnCommit(block *types.Block) error {
	return nil
}

// ValidatorPool a pool storing validators.
type ValidatorPool struct {
	agency consensus.Agency

	lock sync.RWMutex

	// Current node's public key
	nodeID discover.NodeID

	// A block number which validators switch point.
	switchPoint uint64
	lastNumber  uint64

	epoch uint64

	prevValidators    *cbfttypes.Validators // Previous validators
	currentValidators *cbfttypes.Validators // Current validators

}

// NewValidatorPool new a validator pool.
func NewValidatorPool(agency consensus.Agency, blockNumber uint64, epoch uint64, nodeID discover.NodeID) *ValidatorPool {
	pool := &ValidatorPool{
		agency: agency,
		nodeID: nodeID,
		epoch:  epoch,
	}
	// FIXME: Check `GetValidator` return error
	if agency.GetLastNumber(blockNumber) == blockNumber {
		pool.prevValidators, _ = agency.GetValidator(blockNumber)
		pool.currentValidators, _ = agency.GetValidator(NextRound(blockNumber))
		pool.lastNumber = agency.GetLastNumber(NextRound(blockNumber))
		if blockNumber != 0 {
			pool.epoch += 1
		}
	} else {
		pool.currentValidators, _ = agency.GetValidator(blockNumber)
		pool.prevValidators = pool.currentValidators
		pool.lastNumber = agency.GetLastNumber(blockNumber)
	}
	// When validator mode is `static`, the `ValidatorBlockNumber` always 0,
	// means we are using static validators. Otherwise, represent use current
	// validators validate start from `ValidatorBlockNumber` block,
	// so `ValidatorBlockNumber` - 1 is the switch point.
	if pool.currentValidators.ValidBlockNumber > 0 {
		pool.switchPoint = pool.currentValidators.ValidBlockNumber - 1
	}

	log.Debug("Update validator", "validators", pool.currentValidators.String(), "switchpoint", pool.switchPoint, "epoch", pool.epoch, "lastNumber", pool.lastNumber)
	return pool
}

// Reset reset validator pool.
func (vp *ValidatorPool) Reset(blockNumber uint64, epoch uint64) {
	if vp.agency.GetLastNumber(blockNumber) == blockNumber {
		vp.prevValidators, _ = vp.agency.GetValidator(blockNumber)
		vp.currentValidators, _ = vp.agency.GetValidator(NextRound(blockNumber))
		vp.lastNumber = vp.agency.GetLastNumber(NextRound(blockNumber))
		vp.epoch = epoch + 1
	} else {
		vp.currentValidators, _ = vp.agency.GetValidator(blockNumber)
		vp.prevValidators = vp.currentValidators
		vp.lastNumber = vp.agency.GetLastNumber(blockNumber)
		vp.epoch = epoch
	}
	if vp.currentValidators.ValidBlockNumber > 0 {
		vp.switchPoint = vp.currentValidators.ValidBlockNumber - 1
	}
	log.Debug("Update validator", "validators", vp.currentValidators.String(), "switchpoint", vp.switchPoint, "epoch", vp.epoch, "lastNumber", vp.lastNumber)
}

// ShouldSwitch check if should switch validators at the moment.
func (vp *ValidatorPool) ShouldSwitch(blockNumber uint64) bool {
	if blockNumber == 0 {
		return false
	}
	if blockNumber == vp.switchPoint {
		return true
	}
	return blockNumber == vp.lastNumber
}

// EqualSwitchPoint returns boolean which representment the switch point
// equal the inputs number.
func (vp *ValidatorPool) EqualSwitchPoint(number uint64) bool {
	return vp.switchPoint > 0 && vp.switchPoint == number
}

func (vp *ValidatorPool) EnableVerifyEpoch(epoch uint64) error {
	if epoch+1 == vp.epoch || epoch == vp.epoch {
		return nil
	}
	return fmt.Errorf("enable verify epoch:%d,%d, request:%d", vp.epoch-1, vp.epoch, epoch)
}

func (vp *ValidatorPool) MockSwitchPoint(number uint64) {
	vp.switchPoint = 0
	vp.lastNumber = number
}

// Update switch validators.
func (vp *ValidatorPool) Update(blockNumber uint64, epoch uint64, eventMux *event.TypeMux) error {
	vp.lock.Lock()
	defer vp.lock.Unlock()

	// Only updated once
	if blockNumber <= vp.switchPoint {
		log.Debug("Already update validator before", "blockNumber", blockNumber, "switchPoint", vp.switchPoint)
		return errors.New("already updated before")
	}

	nds, err := vp.agency.GetValidator(NextRound(blockNumber))
	if err != nil {
		log.Error("Get validator error", "blockNumber", blockNumber, "err", err)
		return err
	}
	vp.prevValidators = vp.currentValidators
	vp.currentValidators = nds
	vp.switchPoint = nds.ValidBlockNumber - 1
	vp.lastNumber = vp.agency.GetLastNumber(NextRound(blockNumber))
	vp.epoch = epoch
	log.Info("Update validator", "validators", nds.String(), "switchpoint", vp.switchPoint, "epoch", vp.epoch, "lastNumber", vp.lastNumber)

	isValidatorBefore := vp.isValidator(epoch-1, vp.nodeID)

	isValidatorAfter := vp.isValidator(epoch, vp.nodeID)

	if isValidatorBefore {
		// If we are still a consensus node, that adding
		// new validators as consensus peer, and removing
		// validators. Added as consensus peersis because
		// we need to keep connect with other validators
		// in the consensus stages. Also we are not needed
		// to keep connect with old validators.
		if isValidatorAfter {
			for _, nodeID := range vp.currentValidators.NodeList() {
				if node, _ := vp.prevValidators.FindNodeByID(nodeID); node == nil {
					eventMux.Post(cbfttypes.AddValidatorEvent{NodeID: nodeID})
					log.Trace("Post AddValidatorEvent", "nodeID", nodeID.String())
				}
			}

			for _, nodeID := range vp.prevValidators.NodeList() {
				if node, _ := vp.currentValidators.FindNodeByID(nodeID); node == nil {
					eventMux.Post(cbfttypes.RemoveValidatorEvent{NodeID: nodeID})
					log.Trace("Post RemoveValidatorEvent", "nodeID", nodeID.String())
				}
			}
		} else {
			for _, nodeID := range vp.prevValidators.NodeList() {
				eventMux.Post(cbfttypes.RemoveValidatorEvent{NodeID: nodeID})
				log.Trace("Post RemoveValidatorEvent", "nodeID", nodeID.String())
			}
		}
	} else {
		// We are become a consensus node, that adding all
		// validators as consensus peer except us. Added as
		// consensus peers is because we need to keep connecting
		// with other validators in the consensus stages.
		if isValidatorAfter {
			for _, nodeID := range vp.currentValidators.NodeList() {
				eventMux.Post(cbfttypes.AddValidatorEvent{NodeID: nodeID})
				log.Trace("Post AddValidatorEvent", "nodeID", nodeID.String())
			}
		}
	}

	return nil
}

// GetValidatorByNodeID get the validator by node id.
func (vp *ValidatorPool) GetValidatorByNodeID(epoch uint64, nodeID discover.NodeID) (*cbfttypes.ValidateNode, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()
	return vp.getValidatorByNodeID(epoch, nodeID)
}

func (vp *ValidatorPool) getValidatorByNodeID(epoch uint64, nodeID discover.NodeID) (*cbfttypes.ValidateNode, error) {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.FindNodeByID(nodeID)
	}
	return vp.currentValidators.FindNodeByID(nodeID)
}

// GetValidatorByAddr get the validator by address.
func (vp *ValidatorPool) GetValidatorByAddr(epoch uint64, addr common.NodeAddress) (*cbfttypes.ValidateNode, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.getValidatorByAddr(epoch, addr)
}

func (vp *ValidatorPool) getValidatorByAddr(epoch uint64, addr common.NodeAddress) (*cbfttypes.ValidateNode, error) {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.FindNodeByAddress(addr)
	}
	return vp.currentValidators.FindNodeByAddress(addr)
}

// GetValidatorByIndex get the validator by index.
func (vp *ValidatorPool) GetValidatorByIndex(epoch uint64, index uint32) (*cbfttypes.ValidateNode, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.getValidatorByIndex(epoch, index)
}

func (vp *ValidatorPool) getValidatorByIndex(epoch uint64, index uint32) (*cbfttypes.ValidateNode, error) {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.FindNodeByIndex(int(index))
	}
	return vp.currentValidators.FindNodeByIndex(int(index))
}

// GetNodeIDByIndex get the node id by index.
func (vp *ValidatorPool) GetNodeIDByIndex(epoch uint64, index int) discover.NodeID {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.getNodeIDByIndex(epoch, index)
}

func (vp *ValidatorPool) getNodeIDByIndex(epoch uint64, index int) discover.NodeID {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.NodeID(index)
	}
	return vp.currentValidators.NodeID(index)
}

// GetIndexByNodeID get the index by node id.
func (vp *ValidatorPool) GetIndexByNodeID(epoch uint64, nodeID discover.NodeID) (uint32, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.getIndexByNodeID(epoch, nodeID)
}

func (vp *ValidatorPool) getIndexByNodeID(epoch uint64, nodeID discover.NodeID) (uint32, error) {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.Index(nodeID)
	}
	return vp.currentValidators.Index(nodeID)
}

// ValidatorList get the validator list.
func (vp *ValidatorPool) ValidatorList(epoch uint64) []discover.NodeID {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.validatorList(epoch)
}

func (vp *ValidatorPool) validatorList(epoch uint64) []discover.NodeID {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.NodeList()
	}
	return vp.currentValidators.NodeList()
}

func (vp *ValidatorPool) Validators(epoch uint64) *cbfttypes.Validators {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators
	}
	return vp.currentValidators
}

// VerifyHeader verify block's header.
func (vp *ValidatorPool) VerifyHeader(header *types.Header) error {
	_, err := crypto.Ecrecover(header.SealHash().Bytes(), header.Signature())
	if err != nil {
		return err
	}
	// todo: need confirmed.
	return vp.agency.VerifyHeader(header, nil)
}

// IsValidator check if the node is validator.
func (vp *ValidatorPool) IsValidator(epoch uint64, nodeID discover.NodeID) bool {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.isValidator(epoch, nodeID)
}

func (vp *ValidatorPool) isValidator(epoch uint64, nodeID discover.NodeID) bool {
	_, err := vp.getValidatorByNodeID(epoch, nodeID)
	return err == nil
}

// IsCandidateNode check if the node is candidate node.
func (vp *ValidatorPool) IsCandidateNode(nodeID discover.NodeID) bool {
	return vp.agency.IsCandidateNode(nodeID)
}

// Len return number of validators.
func (vp *ValidatorPool) Len(epoch uint64) int {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.Len()
	}
	return vp.currentValidators.Len()
}

// Verify verifies signature using the specified validator's bls public key.
func (vp *ValidatorPool) Verify(epoch uint64, validatorIndex uint32, msg, signature []byte) error {
	validator, err := vp.GetValidatorByIndex(epoch, validatorIndex)
	if err != nil {
		return err
	}

	return validator.Verify(msg, signature)
}

// VerifyAggSig verifies aggregation signature using the specified validators' public keys.
func (vp *ValidatorPool) VerifyAggSig(epoch uint64, validatorIndexes []uint32, msg, signature []byte) bool {
	vp.lock.RLock()
	validators := vp.currentValidators
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		validators = vp.prevValidators
	}

	nodeList, err := validators.NodeListByIndexes(validatorIndexes)
	if err != nil {
		vp.lock.RUnlock()
		return false
	}
	vp.lock.RUnlock()

	var pub bls.PublicKey
	for _, node := range nodeList {
		pub.Add(node.BlsPubKey)
	}

	var sig bls.Sign
	err = sig.Deserialize(signature)
	if err != nil {
		return false
	}
	return sig.Verify(&pub, string(msg))
}

func (vp *ValidatorPool) VerifyAggSigByBA(epoch uint64, vSet *utils.BitArray, msg, signature []byte) error {
	vp.lock.RLock()
	validators := vp.currentValidators
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		validators = vp.prevValidators
	}

	nodeList, err := validators.NodeListByBitArray(vSet)
	if err != nil || len(nodeList) == 0 {
		vp.lock.RUnlock()
		return fmt.Errorf("not found validators: %v", err)
	}
	vp.lock.RUnlock()

	var pub bls.PublicKey
	pub.Deserialize(nodeList[0].BlsPubKey.Serialize())
	for i := 1; i < len(nodeList); i++ {
		pub.Add(nodeList[i].BlsPubKey)
	}

	var sig bls.Sign
	err = sig.Deserialize(signature)
	if err != nil {
		return err
	}
	if !sig.Verify(&pub, string(msg)) {
		log.Error("Verify signature fail", "epoch", epoch, "vSet", vSet.String(), "msg", hex.EncodeToString(msg), "signature", hex.EncodeToString(signature), "nodeList", nodeList, "validators", validators.String())
		return errors.New("bls verifies signature fail")
	}
	return nil
}

func (vp *ValidatorPool) epochToBlockNumber(epoch uint64) uint64 {
	if epoch > vp.epoch {
		panic(fmt.Sprintf("get unknown epoch, current:%d, request:%d", vp.epoch, epoch))
	}
	if epoch+1 == vp.epoch {
		return vp.switchPoint
	}
	return vp.switchPoint + 1
}

func (vp *ValidatorPool) Flush(header *types.Header) error {
	return vp.agency.Flush(header)
}

func (vp *ValidatorPool) Commit(block *types.Block) error {
	return vp.agency.OnCommit(block)
}

func NextRound(blockNumber uint64) uint64 {
	return blockNumber + 1
}
