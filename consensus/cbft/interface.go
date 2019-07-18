package cbft

import "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"

// Define the function of the router that other packages need to access.
type Router interface {
	// Broadcast forwards the message to the router for distribution.
	gossip(m *types.MsgPackage)

	// Send message to a known peerId. Determine if the peerId has established
	// a connection before sending.
	sendMessage(m *types.MsgPackage)
}
