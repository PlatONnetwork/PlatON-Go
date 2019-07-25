package types

import (
	"crypto/ecdsa"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

type OptionsConfig struct {
	NodePriKey *ecdsa.PrivateKey
	NodeID     discover.NodeID
	BlsPriKey  *bls.SecretKey
	WalMode    bool

	PeerMsgQueueSize uint64
	EvidenceDir      string
	MaxPingLatency   int64 // maxPingLatency is the time in milliseconds between Ping and Pong
	MaxQueuesLimit   int64 // The maximum value that a single node can send a message.
}

type Config struct {
	Sys    *params.CbftConfig
	Option *OptionsConfig
}
