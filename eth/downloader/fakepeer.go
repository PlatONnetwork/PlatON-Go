// Copyright 2017 The go-ethereum Authors
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
	"errors"
	"math/big"

	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
)

// FakePeer is a mock downloader peer that operates on a local database instance
// instead of being an actual live node. It's useful for testing and to implement
// sync commands from an existing local database.
type FakePeer struct {
	id         string
	db         ethdb.Database
	hc         *core.HeaderChain
	dl         *Downloader
	snapshotDB snapshotdb.DB
}

// NewFakePeer creates a new mock downloader peer with the given data sources.
func NewFakePeer(id string, db ethdb.Database, sdb snapshotdb.DB, hc *core.HeaderChain, dl *Downloader) *FakePeer {
	return &FakePeer{id: id, db: db, hc: hc, dl: dl, snapshotDB: sdb}
}

// Head implements downloader.Peer, returning the current head hash and number
// of the best known header.
func (p *FakePeer) Head() (common.Hash, *big.Int) {
	header := p.hc.CurrentHeader()
	return header.Hash(), header.Number
}

// RequestHeadersByHash implements downloader.Peer, returning a batch of headers
// defined by the origin hash and the associated query parameters.
func (p *FakePeer) RequestHeadersByHash(hash common.Hash, amount int, skip int, reverse bool) error {
	var (
		headers []*types.Header
		unknown bool
	)
	for !unknown && len(headers) < amount {
		origin := p.hc.GetHeaderByHash(hash)
		if origin == nil {
			break
		}
		number := origin.Number.Uint64()
		headers = append(headers, origin)
		if reverse {
			for i := 0; i <= skip; i++ {
				if header := p.hc.GetHeader(hash, number); header != nil {
					hash = header.ParentHash
					number--
				} else {
					unknown = true
					break
				}
			}
		} else {
			var (
				current = origin.Number.Uint64()
				next    = current + uint64(skip) + 1
			)
			if header := p.hc.GetHeaderByNumber(next); header != nil {
				if p.hc.GetBlockHashesFromHash(header.Hash(), uint64(skip+1))[skip] == hash {
					hash = header.Hash()
				} else {
					unknown = true
				}
			} else {
				unknown = true
			}
		}
	}
	p.dl.DeliverHeaders(p.id, headers)
	return nil
}

// RequestHeadersByNumber implements downloader.Peer, returning a batch of headers
// defined by the origin number and the associated query parameters.
func (p *FakePeer) RequestHeadersByNumber(number uint64, amount int, skip int, reverse bool) error {
	var (
		headers []*types.Header
		unknown bool
	)
	base, _ := p.snapshotDB.BaseNum()
	for !unknown && len(headers) < amount {
		if number > base.Uint64() {
			break
		}
		origin := p.hc.GetHeaderByNumber(number)
		if origin == nil {
			break
		}
		if reverse {
			if number >= uint64(skip+1) {
				number -= uint64(skip + 1)
			} else {
				unknown = true
			}
		} else {
			number += uint64(skip + 1)
		}
		headers = append(headers, origin)
	}
	p.dl.DeliverHeaders(p.id, headers)
	return nil
}

// RequestBodies implements downloader.Peer, returning a batch of block bodies
// corresponding to the specified block hashes.
func (p *FakePeer) RequestBodies(hashes []common.Hash) error {
	var (
		txs [][]*types.Transaction
		ex  [][]byte
	)
	for _, hash := range hashes {
		block := rawdb.ReadBlock(p.db, hash, *p.hc.GetBlockNumber(hash))

		txs = append(txs, block.Transactions())
		ex = append(ex, block.ExtraData())
	}
	p.dl.DeliverBodies(p.id, txs, ex)
	return nil
}

// RequestReceipts implements downloader.Peer, returning a batch of transaction
// receipts corresponding to the specified block hashes.
func (p *FakePeer) RequestReceipts(hashes []common.Hash) error {
	var receipts [][]*types.Receipt
	for _, hash := range hashes {
		receipts = append(receipts, rawdb.ReadRawReceipts(p.db, hash, *p.hc.GetBlockNumber(hash)))
	}
	p.dl.DeliverReceipts(p.id, receipts)
	return nil
}

// RequestNodeData implements downloader.Peer, returning a batch of state trie
// nodes corresponding to the specified trie hashes.
func (p *FakePeer) RequestNodeData(hashes []common.Hash) error {
	var data [][]byte
	for _, hash := range hashes {
		if entry, err := p.db.Get(hash.Bytes()); err == nil {
			data = append(data, entry)
		} else {
			secureKey := make([]byte, 11+32)
			var secureKeyPrefix = []byte("secure-key-")
			secureKey = append(secureKey[:0], secureKeyPrefix...)
			secureKey = append(secureKey, hash[:]...)
			if v, err := p.db.Get(secureKey); err != nil {
				return err
			} else {
				data = append(data, v)
			}
		}
	}
	log.Debug("RequestNodeData", "DeliverReceipts", len(data), "len", len(hashes))
	p.dl.DeliverNodeData(p.id, data)
	return nil
}

func (p *FakePeer) RequestPPOSStorage() error {
	f := func(num *big.Int, iter iterator.Iterator) error {
		var (
			count int
			KVNum uint64
		)
		KVs := make([]PPOSStorageKV, 0)
		if num == nil {
			return errors.New("num should not be nil")
		}
		Pivot := p.hc.GetHeaderByNumber(num.Uint64())
		Latest := p.hc.CurrentHeader()
		if err := p.dl.DeliverPposInfo(p.id, Latest, Pivot); err != nil {
			log.Error("[GetPPOSStorageMsg]send last ppos meassage fail", "error", err)
			return err
		}
		for iter.Next() {
			k, v := make([]byte, len(iter.Key())), make([]byte, len(iter.Value()))
			copy(k, iter.Key())
			copy(v, iter.Value())
			kv := [2][]byte{
				k,
				v,
			}
			KVs = append(KVs, kv)
			KVNum++
			count++
			if count >= PPOSStorageKVSizeFetch {
				if err := p.dl.DeliverPposStorage(p.id, KVs, false, KVNum); err != nil {
					log.Error("[GetPPOSStorageMsg]send ppos meassage fail", "error", err, "kvnum", KVNum)
					return err
				}
				count = 0
				KVs = make([]PPOSStorageKV, 0)
			}
		}
		if err := p.dl.DeliverPposStorage(p.id, KVs, true, KVNum); err != nil {
			log.Error("[GetPPOSStorageMsg]send last ppos meassage fail", "error", err)
			return err
		}
		return nil
	}

	if err := p.snapshotDB.WalkBaseDB(nil, f); err != nil {
		log.Error("[GetPPOSStorageMsg]send  ppos storage fail", "error", err)
		return err
	}
	return nil
}
func (p *FakePeer) RequestOriginAndPivotByCurrent(m uint64) error {
	oHead := p.hc.GetHeaderByNumber(m)
	pivot, err := p.snapshotDB.BaseNum()
	if err != nil {
		return errors.New("GetOriginAndPivot get snapshotdb baseNum fail")
	}
	if pivot == nil {
		return errors.New("[GetOriginAndPivot] pivot should not be nil")
	}
	pHead := p.hc.GetHeaderByNumber(pivot.Uint64())

	data := make([]*types.Header, 0)
	data = append(data, oHead, pHead)
	if err := p.dl.DeliverOriginAndPivot(p.id, data); err != nil {
		log.Error("[GetOriginAndPivotMsg]send data meassage fail", "error", err)
		return err
	}
	return nil
}
