package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/PlatONnetwork/PlatON-Go/common"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto/vrf"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

type BlockChainReactor struct {
	privateKey *ecdsa.PrivateKey

	eventMux     *event.TypeMux
	bftResultSub *event.TypeMuxSubscription

	// xxPlugin container
	basePluginMap map[int]plugin.BasePlugin
	// Order rules for xxPlugins called in BeginBlocker
	beginRule []int
	// Order rules for xxPlugins called in EndBlocker
	endRule []int
}

var bcr *BlockChainReactor

func NewBlockChainReactor(pri *ecdsa.PrivateKey, mux *event.TypeMux) *BlockChainReactor {
	if nil == bcr {
		bcr = &BlockChainReactor{
			eventMux:      mux,
			basePluginMap: make(map[int]plugin.BasePlugin, 0),
		}
		// Subscribe events for confirmed blocks
		bcr.bftResultSub = bcr.eventMux.Subscribe(cbfttypes.CbftResult{})

		// start the loop rutine
		go bcr.loop()
	}
	return bcr
}

// Getting the global bcr single instance
//func GetReactorInstance () *BlockChainReactor {
//	return bcr
//}

func (brc *BlockChainReactor) loop() {

	for {
		select {
		case obj := <-bcr.bftResultSub.Chan():
			if obj == nil {
				log.Error("BlockChainReactor receive nil bftResultEvent maybe channel is closed")
				continue
			}
			cbftResult, ok := obj.Data.(cbfttypes.CbftResult)
			if !ok {
				log.Error("receive bft result type error")
				continue
			}
			block := cbftResult.Block
			// Short circuit when receiving empty result.
			if block == nil {
				log.Error("Cbft result error, block is nil")
				continue
			}

			/**
			TODO flush the seed and the package ratio
			*/

			if plugin, ok := brc.basePluginMap[xcom.StakingRule]; ok {
				if err := plugin.Confirmed(block); nil != err {
					log.Error("Failed to call Staking Confirmed", "blockNumber", block.Number(), "blockHash", block.Hash().Hex(), "err", err.Error())
				}

			}

			/*// TODO Slashing
			if plugin, ok := brc.basePluginMap[common.SlashingRule]; ok {
				if err := plugin.Confirmed(block); nil != err {
					log.Error("Failed to call Staking Confirmed", "blockNumber", block.Number(), "blockHash", block.Hash().Hex(), "err", err.Error())
				}

			}
			}*/


		default:
			return

		}
	}

}

func (bcr *BlockChainReactor) RegisterPlugin(pluginRule int, plugin plugin.BasePlugin) {
	bcr.basePluginMap[pluginRule] = plugin
}
func (bcr *BlockChainReactor) SetBeginRule(rule []int) {
	bcr.beginRule = rule
}
func (bcr *BlockChainReactor) SetEndRule(rule []int) {
	bcr.endRule = rule
}

// Called before every block has not executed all txs
func (bcr *BlockChainReactor) BeginBlocker(header *types.Header, state xcom.StateDB) (bool, error) {

	blockHash := common.ZeroHash

	// store the sign in  header.Extra[32:97]
	if len(header.Extra[32:]) >= 65 {
		if bytes.Equal(header.Extra[32:97], make([]byte, 65)) {
			// Generate vrf proof

		} else {
			blockHash = header.Hash()
		}
	}

	// todo maybe vrf

	for _, pluginName := range bcr.beginRule {
		if plugin, ok := bcr.basePluginMap[pluginName]; ok {
			if flag, err := plugin.BeginBlock(blockHash, header, state); nil != err {
				return flag, err
			}
		}
	}
	return false, nil
}

// Called after every block had executed all txs
func (bcr *BlockChainReactor) EndBlocker(header *types.Header, state xcom.StateDB) (bool, error) {

	blockHash := common.ZeroHash

	// store the sign in  header.Extra[32:97]
	if len(header.Extra[32:]) == 65 && !bytes.Equal(header.Extra[32:97], make([]byte, 65)) {
		blockHash = header.Hash()
	}

	// todo maybe vrf

	for _, pluginName := range bcr.endRule {
		if plugin, ok := bcr.basePluginMap[pluginName]; ok {
			if flag, err := plugin.EndBlock(blockHash, header, state); nil != err {
				return flag, err
			}
		}
	}
	return false, nil
}

func (bcr *BlockChainReactor) Verify_tx(tx *types.Transaction, from common.Address) (err error) {

	if _, ok := vm.PlatONPrecompiledContracts[from]; !ok {
		return nil
	}

	input := tx.Data()

	var contract vm.PlatONPrecompiledContract
	switch from {
	case cvm.StakingContractAddr:
		contract = vm.PlatONPrecompiledContracts[cvm.StakingContractAddr]
	case cvm.RestrictingContractAddr:
		// TODO
	case cvm.AwardMgrContractAddr:
		// TODO
	case cvm.SlashingContractAddr:
		// TODO
	default:
		return nil
	}
	_, _, err = plugin.Verify_tx_data(input, contract.FnSigns())
	return
}

func (bcr *BlockChainReactor) Sign(msg interface{}) error {
	return nil
}

func (bcr *BlockChainReactor) VerifySign(msg interface{}) error {
	return nil
}

func (bcr *BlockChainReactor) GetLastNumber(blockNumber uint64) uint64 {
	return 0
}

func (brc *BlockChainReactor) GetValidator(blockNumber uint64) (*cbfttypes.Validators, error) {
	return nil, nil
}

func (bcr *BlockChainReactor) IsCandidateNode(nodeID discover.NodeID) bool {
	return false
}

func GenerateNonce(proof []byte) ([]byte, error) {
	log.Debug("Generate proof based on input","proof", hex.EncodeToString(proof))
	return vrf.Prove(bcr.privateKey, vrf.ProofToHash(proof))
}

func VerifyVrf(header *types.Header) error {
	// Verify VRF Proof
	log.Debug("Verification block vrf prove", "blockNumber", header.Number.Uint64(), "hash", header.Hash().TerminalString(), "sealHash", header.SealHash().TerminalString(), "prove", hex.EncodeToString(header.Nonce[:]))
	/*sign := header.Extra[32:common.ExtraSeal]
	pk, err := crypto.SigToPub(header.SealHash().Bytes(), sign)
	if nil != err {
		return err
	}
	pext := cbft.blockExtMap.findBlockByHash(block.ParentHash())
	if pext == nil {
		pext = cbft.blockChain.GetBlockByHash(block.ParentHash())
	}
	if pext != nil {
		parentNonce := vrf.ProofToHash(pext.Nonce())
		if value, err := vrf.Verify(pk, block.Nonce(), parentNonce); nil != err {
			cbft.log.Error("Vrf proves verification failure", "blockNumber", block.NumberU64(), "proof", hex.EncodeToString(block.Nonce()), "input", hex.EncodeToString(parentNonce), "err", err)
			return err
		} else if !value {
			cbft.log.Error("Vrf proves verification failure", "blockNumber", block.NumberU64(), "proof", hex.EncodeToString(block.Nonce()), "input", hex.EncodeToString(parentNonce))
			return errInvalidVrfProve
		}
		cbft.log.Info("Vrf proves successful verification", "blockNumber", block.NumberU64(), "proof", hex.EncodeToString(block.Nonce()), "input", hex.EncodeToString(parentNonce))
	} else {
		cbft.log.Error("Vrf proves verification failure, Cannot find parent block", "blockNumber", block.NumberU64(), "hash", block.Hash().TerminalString(), "parentHash", block.ParentHash().TerminalString())
		return errNotFoundViewBlock
	}*/
	return nil
}
