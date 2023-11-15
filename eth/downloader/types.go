// Copyright 2015 The go-ethereum Authors
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

package downloader

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

// dataPack is a data message returned by a peer for some query.
type dataPack interface {
	PeerId() string
	Items() int
	Stats() string
}

// headerPack is a batch of block headers returned by a peer.
type headerPack struct {
	peerID  string
	headers []*types.Header
}

func (p *headerPack) PeerId() string { return p.peerID }
func (p *headerPack) Items() int     { return len(p.headers) }
func (p *headerPack) Stats() string  { return fmt.Sprintf("%d", len(p.headers)) }

// pposStoragePack is a batch of ppos storage returned by a peer.
type pposStoragePack struct {
	peerID string
	kvs    [][2][]byte
	last   bool
	kvNum  uint64
	base   bool

	blocks    []snapshotdb.BlockData
	baseBlock uint64
}

func (p *pposStoragePack) PeerId() string { return p.peerID }
func (p *pposStoragePack) Items() int     { return len(p.kvs) + len(p.blocks) }
func (p *pposStoragePack) Stats() string  { return fmt.Sprintf("%d", len(p.kvs)+len(p.blocks)) }
func (p *pposStoragePack) KVs() [][2][]byte {
	var kv [][2][]byte
	for _, value := range p.kvs {
		kv = append(kv, value)
	}
	return kv
}

// pposStoragePack is a batch of ppos storage returned by a peer.
type pposInfoPack struct {
	peerID string
	latest *types.Header
	pivot  *types.Header
}

func (p *pposInfoPack) PeerId() string { return p.peerID }
func (p *pposInfoPack) Items() int     { return 1 }
func (p *pposInfoPack) Stats() string  { return fmt.Sprint(1) }
