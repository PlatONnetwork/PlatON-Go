package network

import (
	ctpyes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
)

// SetSendQueueHook
func (h *EngineManager) SetSendQueueHook(f func(*ctpyes.MsgPackage)) {
	h.sendQueueHook = f
}
