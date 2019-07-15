package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"log"
	"sync"
)

type Cbft struct {
	config     *params.CbftConfig
	eventMux   *event.TypeMux
	closeOnce  sync.Once
	exitCh     chan struct{}
	txPool     *core.TxPool
	blockChain *core.BlockChain //the block chain
	peerMsgCh  chan *MsgInfo
	syncMsgCh  chan *MsgInfo
	evPool     EvidencePool
	log        log.Logger

	//Control the current view state
	state viewState

	//Block executor, the block responsible for executing the current view
	executor blockExecutor

	//Verification security rules for proposed blocks and viewchange
	safetyRules safetyRules

	//Determine when to allow voting
	voteRules voteRules

	//Store blocks that are not committed
	blockTree blockTree
}
