package core

import (
	"bytes"
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	xcommon "github.com/PlatONnetwork/PlatON-Go/x/common"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

type BlockChainReactor struct {


	eventMux      	*event.TypeMux
	bftResultSub 	*event.TypeMuxSubscription


	// xxPlugin container
	basePluginMap  	map[int]plugin.BasePlugin
	// Order rules for xxPlugins called in BeginBlocker
	beginRule		[]int
	// Order rules for xxPlugins called in EndBlocker
	endRule 		[]int
}


var (
	DecodeTxDataErr = errors.New("decode tx data is err")
)


var bcr *BlockChainReactor


func New (mux *event.TypeMux) *BlockChainReactor {
	if nil == bcr {
		bcr = &BlockChainReactor{
			eventMux: 		mux,
			basePluginMap: 	make(map[int]plugin.BasePlugin, 0),
			//beginRule:		make([]string, 0),
			//endRule: 		make([]string, 0),
		}
		// Subscribe events for confirmed blocks
		bcr.bftResultSub = bcr.eventMux.Subscribe(cbfttypes.CbftResult{})

		// start the loop rutine
		go bcr.loop()
	}
	return bcr
}

// Getting the global bcr single instance
func GetReactorInstance () *BlockChainReactor {
	return bcr
}


func (brc *BlockChainReactor) loop () {

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

			if plugin, ok := brc.basePluginMap[xcommon.StakingRule]; ok {
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


func (bcr *BlockChainReactor) RegisterPlugin (pluginRule int, plugin plugin.BasePlugin) {
	bcr.basePluginMap[pluginRule] = plugin
}
func (bcr *BlockChainReactor) SetBeginRule(rule []int) {
	bcr.beginRule = rule
}
func (bcr *BlockChainReactor) SetEndRule(rule []int) {
	bcr.endRule = rule
}


// Called before every block has not executed all txs
func (bcr *BlockChainReactor) BeginBlocker (header *types.Header, state plugin.StateDB) (bool, error) {

	for _, pluginName := range bcr.beginRule {
		if plugin, ok := bcr.basePluginMap[pluginName]; ok {
			if flag, err := plugin.BeginBlock(header, state); nil != err {
				return flag, err
			}
		}
	}
	return false, nil
}

// Called after every block had executed all txs
func (bcr *BlockChainReactor) EndBlocker (header *types.Header, state plugin.StateDB) (bool, error) {

	for _, pluginName := range bcr.endRule {
		if plugin, ok := bcr.basePluginMap[pluginName]; ok {
			if flag, err := plugin.EndBlock(header, state); nil != err {
				return flag, err
			}
		}
	}
	return false, nil
}


func (bcr *BlockChainReactor) Verify_tx (tx *types.Transaction, from common.Address) (err error) {

	//if _, ok := vm.PrecompiledContracts[from]; !ok {
	//	err = nil
	//	return
	//}

	input := tx.Data()

	var args [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &args); nil != err {
		return DecodeTxDataErr
	}

	var plugin plugin.BasePlugin
	switch from {
	case vm.StakingContractAddr:
		plugin = bcr.basePluginMap[xcommon.StakingRule]
	case vm.LockRepoContractAddr:
		plugin = bcr.basePluginMap[xcommon.LockrepoRule]
	case vm.AwardMgrContractAddr:
		plugin = bcr.basePluginMap[xcommon.AwardmgrRule]
	case vm.SlashingContractAddr:
		plugin = bcr.basePluginMap[xcommon.SlashingRule]
	default:
		return nil
	}
	err = plugin.Verify_tx_data(args)
	return
}


func (bcr *BlockChainReactor) GetPlugin(pluginLabel int) plugin.StakingPlugin {
	//return bcr.basePluginMap[pluginLabel]
	return plugin.StakingPlugin{}
}
