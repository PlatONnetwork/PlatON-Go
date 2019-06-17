package core

import (
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

type BlockChainReactor struct {


	eventMux      	*event.TypeMux
	bftResultSub 	*event.TypeMuxSubscription


	// xxPlugin container
	basePluginMap  	map[string]BasePlugin
	// Order rules for xxPlugins called in BeginBlocker
	beginRule		[]string
	// Order rules for xxPlugins called in EndBlocker
	endRule 		[]string
}


var bcr *BlockChainReactor


func New (mux *event.TypeMux) *BlockChainReactor {
	if nil == bcr {
		bcr = &BlockChainReactor{
			eventMux: 		mux,
			basePluginMap: 	make(map[string]BasePlugin, 0),
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
func GetInstance () *BlockChainReactor {
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

			default:
				return

		}
	}

}


func (bcr *BlockChainReactor) RegisterPlugin (pluginName string, plugin BasePlugin) {
	bcr.basePluginMap[pluginName] = plugin
}
func (bcr *BlockChainReactor) SetBeginRule(rule []string) {
	bcr.beginRule = rule
}
func (bcr *BlockChainReactor) SetEndRule(rule []string) {
	bcr.endRule = rule
}


// Called before every block has not executed all txs
func (bcr *BlockChainReactor) BeginBlocker (header *types.Header, state *state.StateDB) (bool, error) {

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
func (bcr *BlockChainReactor) EndBlocker (header *types.Header, state *state.StateDB) (bool, error) {

	for _, pluginName := range bcr.endRule {
		if plugin, ok := bcr.basePluginMap[pluginName]; ok {
			if flag, err := plugin.EndBlock(header, state); nil != err {
				return flag, err
			}
		}
	}
	return false, nil
}


