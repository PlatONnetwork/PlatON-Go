// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package les

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/light"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

// lesCommons contains fields needed by both server and client.
type lesCommons struct {
	config                       *eth.Config
	iConfig                      *light.IndexerConfig
	chainDb                      ethdb.Database
	protocolManager              *ProtocolManager
	chtIndexer, bloomTrieIndexer *core.ChainIndexer
}

// NodeInfo represents a short summary of the Ethereum sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network    uint64                   `json:"network"`    // Ethereum network ID (1=Frontier, 2=Morden, Ropsten=3, Rinkeby=4)
	Genesis    common.Hash              `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig      `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash              `json:"head"`       // SHA3 hash of the host's best owned block
	CHT        params.TrustedCheckpoint `json:"cht"`        // Trused CHT checkpoint for fast catchup
}

// makeProtocols creates protocol descriptors for the given LES versions.
func (c *lesCommons) makeProtocols(versions []uint) []p2p.Protocol {
	protos := make([]p2p.Protocol, len(versions))
	for i, version := range versions {
		version := version
		protos[i] = p2p.Protocol{
			Name:     "les",
			Version:  version,
			Length:   ProtocolLengths[version],
			NodeInfo: c.nodeInfo,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				return c.protocolManager.runPeer(version, p, rw)
			},
			PeerInfo: func(id enode.ID) interface{} {
				if p := c.protocolManager.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		}
	}
	return protos
}

// nodeInfo retrieves some protocol metadata about the running host node.
func (c *lesCommons) nodeInfo() interface{} {
	var cht params.TrustedCheckpoint
	sections, _, _ := c.chtIndexer.Sections()
	sections2, _, _ := c.bloomTrieIndexer.Sections()

	if sections2 < sections {
		sections = sections2
	}
	if sections > 0 {
		sectionIndex := sections - 1
		sectionHead := c.bloomTrieIndexer.SectionHead(sectionIndex)
		cht = params.TrustedCheckpoint{
			SectionIndex: sectionIndex,
			SectionHead:  sectionHead,
			CHTRoot:      light.GetChtRoot(c.chainDb, sectionIndex, sectionHead),
			BloomRoot:    light.GetBloomTrieRoot(c.chainDb, sectionIndex, sectionHead),
		}
	}

	chain := c.protocolManager.blockchain
	return &NodeInfo{
		Network:    c.config.NetworkId,
		Genesis:    chain.Genesis().Hash(),
		Config:     chain.Config(),
		Head:       chain.CurrentHeader().Hash(),
		CHT:        cht,
	}
}
