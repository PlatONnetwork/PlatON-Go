package cbft

import (
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
)

func (cbft *Cbft) newQCState(number uint64, hash common.Hash) error {
	// 调用函数获取QC PrepareBlock和QuorumCert
	// TODO
	return nil
}

func (cbft *Cbft) addQCState(number uint64, hash common.Hash) error {
	// 调用函数获取QC PrepareBlock和QuorumCert
	// TODO
	return nil
}

func (cbft *Cbft) confirmViewChange(view *protocols.ViewChange, master bool) error {
	// TODO
	return nil
}

func (cbft *Cbft) sendViewChange(view *protocols.ViewChange) error {
	// TODO
	return nil
}

func (cbft *Cbft) sendPrepareBlock(pb *protocols.PrepareBlock) error {
	// TODO
	return nil
}

func (cbft *Cbft) sendPrepareVote(pb *protocols.PrepareBlock, pv *protocols.PrepareVote) error {
	// TODO
	return nil
}

func (cbft *Cbft) recoveryMsg(msg interface{}) {
	cbft.log.Debug("Recover journal message from wal", "msgType", reflect.TypeOf(msg))

	switch m := msg.(type) {
	case *protocols.ConfirmedViewChange:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "confirmedViewChange", m.ViewChange.String())

	case *protocols.SendViewChange:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendViewChange", m.ViewChange.String())

	case *protocols.SendPrepareBlock:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareBlock", m.PrepareBlock.String())

	case *protocols.SendPrepareVote:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareVote", m.PrepareVote.String())
	}
}

func (cbft *Cbft) recoveryChainState(chainState *protocols.ChainState) {
	cbft.log.Debug("Recover chain state from wal", "chainState", chainState.String())

	// TODO
}
