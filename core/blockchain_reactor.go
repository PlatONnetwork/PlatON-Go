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

package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/handler"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

type BlockChainReactor struct {
	vh            *handler.VrfHandler
	eventMux      *event.TypeMux
	bftResultSub  *event.TypeMuxSubscription
	basePluginMap map[int]plugin.BasePlugin // xxPlugin container
	beginRule     []int                     // Order rules for xxPlugins called in BeginBlocker
	endRule       []int                     // Order rules for xxPlugins called in EndBlocker
	validatorMode string                    // mode: static, inner, ppos
	NodeId        discover.NodeID           // The nodeId of current node
	exitCh        chan chan struct{}        // Used to receive an exit signal
	exitOnce      sync.Once
	chainID       *big.Int
}

var (
	bcrOnce sync.Once
	bcr     *BlockChainReactor
)

func NewBlockChainReactor(mux *event.TypeMux, chainId *big.Int) *BlockChainReactor {
	bcrOnce.Do(func() {
		log.Info("Init BlockChainReactor ...")
		bcr = &BlockChainReactor{
			eventMux:      mux,
			basePluginMap: make(map[int]plugin.BasePlugin, 0),
			exitCh:        make(chan chan struct{}),
			chainID:       chainId,
		}
	})
	return bcr
}

func (bcr *BlockChainReactor) Start(mode string) {
	bcr.setValidatorMode(mode)
	if mode == common.PPOS_VALIDATOR_MODE {
		// Subscribe events for confirmed blocks
		bcr.bftResultSub = bcr.eventMux.Subscribe(cbfttypes.CbftResult{})
		// start the loop rutine
		go bcr.loop()
	}
}

func (bcr *BlockChainReactor) Close() {
	if bcr.validatorMode == common.PPOS_VALIDATOR_MODE {
		bcr.exitOnce.Do(func() {
			exitDone := make(chan struct{})
			bcr.exitCh <- exitDone
			<-exitDone
			close(exitDone)
		})
	}
	log.Info("blockchain_reactor closed")
}

func (bcr *BlockChainReactor) GetChainID() *big.Int {
	return bcr.chainID
}

// Getting the global bcr single instance
func GetReactorInstance() *BlockChainReactor {
	return bcr
}

func (bcr *BlockChainReactor) loop() {

	for {
		select {
		case obj := <-bcr.bftResultSub.Chan():
			if obj == nil {
				//log.Error("blockchain_reactor receive nil bftResultEvent maybe channel is closed")
				continue
			}
			cbftResult, ok := obj.Data.(cbfttypes.CbftResult)
			if !ok {
				log.Error("blockchain_reactor receive bft result type error")
				continue
			}
			bcr.commit(cbftResult.Block)
		// stop this routine
		case done := <-bcr.exitCh:
			close(bcr.exitCh)
			log.Info("blockChain reactor loop exit")
			done <- struct{}{}
			return
		}
	}

}

func (bcr *BlockChainReactor) commit(block *types.Block) error {
	if block == nil {
		log.Error("blockchain_reactor receive Cbft result error, block is nil")
		return nil
	}
	/**
	notify P2P module the nodeId of the next round validator
	*/
	if plugin, ok := bcr.basePluginMap[xcom.StakingRule]; ok {
		if err := plugin.Confirmed(bcr.NodeId, block); nil != err {
			log.Error("Failed to call Staking Confirmed", "blockNumber", block.Number(), "blockHash", block.Hash().Hex(), "err", err.Error())
		}

	}

	log.Info("Call snapshotdb commit on blockchain_reactor", "blockNumber", block.Number(), "blockHash", block.Hash())
	if err := snapshotdb.Instance().Commit(block.Hash()); nil != err {
		log.Error("Failed to call snapshotdb commit on blockchain_reactor", "blockNumber", block.Number(), "blockHash", block.Hash(), "err", err)
		return err
	}
	return nil
}

func (bcr *BlockChainReactor) OnCommit(block *types.Block) error {
	if bcr.validatorMode == common.PPOS_VALIDATOR_MODE {
		return bcr.commit(block)
	}
	return nil
}

func (bcr *BlockChainReactor) RegisterPlugin(pluginRule int, plugin plugin.BasePlugin) {
	bcr.basePluginMap[pluginRule] = plugin
}

func (bcr *BlockChainReactor) SetPluginEventMux() {
	plugin.StakingInstance().SetEventMux(bcr.eventMux)
}

func (bcr *BlockChainReactor) setValidatorMode(mode string) {
	bcr.validatorMode = mode
}

func (bcr *BlockChainReactor) SetVRFhandler(vher *handler.VrfHandler) {
	bcr.vh = vher
}

func (bcr *BlockChainReactor) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	if bcr.validatorMode == common.PPOS_VALIDATOR_MODE && nil != privateKey {
		if nil != bcr.vh {
			bcr.vh.SetPrivateKey(privateKey)
		}
		plugin.SlashInstance().SetPrivateKey(privateKey)
		bcr.NodeId = discover.PubkeyID(&privateKey.PublicKey)
	}
}

func (bcr *BlockChainReactor) SetBeginRule(rule []int) {
	bcr.beginRule = rule
}
func (bcr *BlockChainReactor) SetEndRule(rule []int) {
	bcr.endRule = rule
}

func (bcr *BlockChainReactor) SetWorkerCoinBase(header *types.Header, nodeId discover.NodeID) {

	/**
	this things about ppos
	*/
	if bcr.validatorMode != common.PPOS_VALIDATOR_MODE {
		return
	}

	nodeIdAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to SetWorkerCoinBase: parse current nodeId is failed", "err", err)
		panic(fmt.Sprintf("parse current nodeId is failed: %s", err.Error()))
	}

	if plu, ok := bcr.basePluginMap[xcom.StakingRule]; ok {
		stake := plu.(*plugin.StakingPlugin)
		can, err := stake.GetCandidateInfo(common.ZeroHash, nodeIdAddr)
		if nil != err {
			log.Error("Failed to SetWorkerCoinBase: Query candidate info is failed", "blockNumber", header.Number,
				"nodeId", nodeId.String(), "nodeIdAddr", nodeIdAddr.Hex(), "err", err)
			return
		}
		header.Coinbase = can.BenefitAddress
		log.Info("SetWorkerCoinBase Successfully", "blockNumber", header.Number,
			"nodeId", nodeId.String(), "nodeIdAddr", nodeIdAddr.Hex(), "coinbase", header.Coinbase.String())
	}

}

// Called before every block has not executed all txs
func (bcr *BlockChainReactor) BeginBlocker(header *types.Header, state xcom.StateDB) error {

	/**
	this things about ppos
	*/
	if bcr.validatorMode != common.PPOS_VALIDATOR_MODE {
		return nil
	}

	blockHash := common.ZeroHash

	// store the sign in  header.Extra[32:97]
	if xutil.IsWorker(header.Extra) {
		// Generate vrf proof
		if value, err := bcr.vh.GenerateNonce(header.Number, header.ParentHash); nil != err {
			return err
		} else {
			header.Nonce = types.EncodeNonce(value)
		}
	} else {
		blockHash = header.Hash()
		// Verify vrf proof
		sign := header.Extra[32:97]
		sealHash := header.SealHash().Bytes()
		pk, err := crypto.SigToPub(sealHash, sign)
		if nil != err {
			return err
		}
		if err := bcr.vh.VerifyVrf(pk, header.Number, header.ParentHash, blockHash, header.Nonce.Bytes()); nil != err {
			return err
		}
	}

	log.Debug("Call snapshotDB newBlock on blockchain_reactor", "blockNumber", header.Number.Uint64(),
		"hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(header.ParentHash.Bytes()))
	if err := snapshotdb.Instance().NewBlock(header.Number, header.ParentHash, blockHash); nil != err {
		log.Error("Failed to call snapshotDB newBlock on blockchain_reactor", "blockNumber",
			header.Number.Uint64(), "hash", hex.EncodeToString(blockHash.Bytes()), "parentHash",
			hex.EncodeToString(header.ParentHash.Bytes()), "err", err)
		return err
	}

	for _, pluginRule := range bcr.beginRule {
		if plugin, ok := bcr.basePluginMap[pluginRule]; ok {
			if err := plugin.BeginBlock(blockHash, header, state); nil != err {
				return err
			}
		}
	}

	// This must not be deleted
	root := state.IntermediateRoot(true)
	log.Debug("BeginBlock StateDB root, end", "blockHash", header.Hash().Hex(), "blockNumber",
		header.Number.Uint64(), "root", root.Hex(), "pointer", fmt.Sprintf("%p", state))

	return nil
}

// Called after every block had executed all txs
func (bcr *BlockChainReactor) EndBlocker(header *types.Header, state xcom.StateDB) error {

	/**
	this things about ppos
	*/
	if bcr.validatorMode != common.PPOS_VALIDATOR_MODE {
		return nil
	}

	blockHash := common.ZeroHash

	if !xutil.IsWorker(header.Extra) {
		blockHash = header.Hash()
	}

	// Store the previous vrf random number
	if err := bcr.vh.Storage(header.Number, header.ParentHash, blockHash, header.Nonce.Bytes()); nil != err {
		log.Error("blockchain_reactor Storage proof failed", "blockNumber", header.Number.Uint64(),
			"blockHash", hex.EncodeToString(blockHash.Bytes()), "err", err)
		return err
	}

	for _, pluginRule := range bcr.endRule {
		if plugin, ok := bcr.basePluginMap[pluginRule]; ok {
			if err := plugin.EndBlock(blockHash, header, state); nil != err {
				return err
			}
		}
	}

	// storage the ppos k-v Hash
	pposHash := snapshotdb.Instance().GetLastKVHash(blockHash)

	if len(pposHash) != 0 && !bytes.Equal(pposHash, make([]byte, len(pposHash))) {
		// store hash about ppos
		state.SetState(cvm.StakingContractAddr, staking.GetPPOSHASHKey(), pposHash)
		log.Debug("Store ppos hash", "blockHash", blockHash.Hex(), "blockNumber", header.Number.Uint64(),
			"pposHash", hex.EncodeToString(pposHash))
	}

	// This must not be deleted
	root := state.IntermediateRoot(true)
	log.Debug("EndBlock StateDB root, end", "blockHash", blockHash.Hex(), "blockNumber",
		header.Number.Uint64(), "root", root.Hex(), "pointer", fmt.Sprintf("%p", state))

	return nil
}

func (bcr *BlockChainReactor) VerifyTx(tx *types.Transaction, to common.Address) error {

	if !vm.IsPlatONPrecompiledContract(to) {
		return nil
	}

	input := tx.Data()
	if len(input) == 0 {
		return nil
	}

	var contract vm.PlatONPrecompiledContract
	switch to {
	case cvm.StakingContractAddr:
		c := vm.PlatONPrecompiledContracts[cvm.StakingContractAddr]
		contract = c.(vm.PlatONPrecompiledContract)
	case cvm.RestrictingContractAddr:
		c := vm.PlatONPrecompiledContracts[cvm.RestrictingContractAddr]
		contract = c.(vm.PlatONPrecompiledContract)
	case cvm.GovContractAddr:
		c := vm.PlatONPrecompiledContracts[cvm.GovContractAddr]
		contract = c.(vm.PlatONPrecompiledContract)
	case cvm.SlashingContractAddr:
		c := vm.PlatONPrecompiledContracts[cvm.SlashingContractAddr]
		contract = c.(vm.PlatONPrecompiledContract)
	default:
		// pass if the contract is validatorInnerContract
		return nil
	}
	// verify the ppos contract tx.data
	if contract != nil {
		if fcode, _, _, err := plugin.VerifyTxData(input, contract.FnSigns()); nil != err {
			return err
		} else {
			return contract.CheckGasPrice(tx.GasPrice(), fcode)
		}
	} else {
		log.Warn("Cannot find an appropriate PlatONPrecompiledContract!")
		return nil
	}
}

func (bcr *BlockChainReactor) Sign(msg interface{}) error {
	return nil
}

func (bcr *BlockChainReactor) VerifySign(msg interface{}) error {
	return nil
}

func (bcr *BlockChainReactor) VerifyHeader(header *types.Header, stateDB *state.StateDB) error {
	return nil
}

func (bcr *BlockChainReactor) GetLastNumber(blockNumber uint64) uint64 {
	return plugin.StakingInstance().GetLastNumber(blockNumber)
}

func (bcr *BlockChainReactor) GetValidator(blockNumber uint64) (*cbfttypes.Validators, error) {
	return plugin.StakingInstance().GetValidator(blockNumber)
}

func (bcr *BlockChainReactor) IsCandidateNode(nodeID discover.NodeID) bool {
	return plugin.StakingInstance().IsCandidateNode(nodeID)
}

func (bcr *BlockChainReactor) Flush(header *types.Header) error {
	log.Debug("Call snapshotdb flush on blockchain_reactor", "blockNumber", header.Number.Uint64(), "hash", hex.EncodeToString(header.Hash().Bytes()))
	if err := snapshotdb.Instance().Flush(header.Hash(), header.Number); nil != err {
		log.Error("Failed to call snapshotdb flush on blockchain_reactor", "blockNumber", header.Number.Uint64(), "hash", hex.EncodeToString(header.Hash().Bytes()), "err", err)
		return err
	}
	return nil
}
