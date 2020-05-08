// Copyright 2018-2019 The PlatON Network Authors
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
	NodePriKey *ecdsa.PrivateKey
	NodeID     discover.NodeID
	BlsPriKey  *bls.SecretKey
	WalMode    bool

	PeerMsgQueueSize  uint64
	EvidenceDir       string
	MaxPingLatency    int64 // maxPingLatency is the time in milliseconds between Ping and Pong
	MaxQueuesLimit    int64 // The maximum value that a single node can send a message.
	BlacklistDeadline int64 // Blacklist expiration time. unit: minute.

	Period uint64
	Amount uint32
}

type Config struct {
	Sys    *params.CbftConfig
	Option *OptionsConfig
}
