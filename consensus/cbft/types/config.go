package types

import "github.com/PlatONnetwork/PlatON-Go/params"

type OptionsConfig struct {
	WalMode bool

	PeerMsgQueueSize uint64
	EvidenceDir      string
	MaxPingLatency   int64 // maxPingLatency is the time in milliseconds between Ping and Pong
	MaxAvgLatency    int64 //maxAvgLatency is the time in milliseconds between two peers
}

type Config struct {
	Sys    *params.CbftConfig
	Option *OptionsConfig
}
