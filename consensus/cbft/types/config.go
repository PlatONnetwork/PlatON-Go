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

package types

import (
	"crypto/ecdsa"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

type OptionsConfig struct {
	NodePriKey *ecdsa.PrivateKey `json:"-"`
	NodeID     discover.NodeID   `json:"nodeID"`
	BlsPriKey  *bls.SecretKey    `json:"-"`
	WalMode    bool              `json:"walMode"`

	PeerMsgQueueSize  uint64 `json:"peerMsgQueueSize"`
	EvidenceDir       string `json:"evidenceDir"`
	MaxPingLatency    int64  `json:"maxPingLatency"`    // maxPingLatency is the time in milliseconds between Ping and Pong
	MaxQueuesLimit    int64  `json:"maxQueuesLimit"`    // The maximum value that a single node can send a message.
	BlacklistDeadline int64  `json:"blacklistDeadline"` // Blacklist expiration time. unit: minute.

	Period uint64 `json:"period"`
	Amount uint32 `json:"amount"`
}

type Config struct {
	Sys    *params.CbftConfig `json:"sys"`
	Option *OptionsConfig     `json:"option"`
}
