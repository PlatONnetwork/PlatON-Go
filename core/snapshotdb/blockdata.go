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

package snapshotdb

import (
	"fmt"
	"io"
	"math/big"
	"sort"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/syndtr/goleveldb/leveldb/memdb"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

type blockData struct {
	BlockHash  common.Hash
	ParentHash common.Hash
	Number     *big.Int
	data       *memdb.DB
	readOnly   bool
	kvHash     common.Hash

	//only use for not commit block
	journal        []journalEntry // Current changes tracked by the journal
	validRevisions []revision
	nextRevisionId int
}

type revision struct {
	id           int
	journalIndex int
}

func (b *blockData) DecodeRLP(s *rlp.Stream) error {
	jk := new(blockWal)
	if err := s.Decode(jk); err != nil {
		return err
	}
	b.Number = jk.BlockNumber
	b.ParentHash = jk.ParentHash
	b.BlockHash = jk.BlockHash
	b.kvHash = jk.KvHash
	b.data = memdb.New(DefaultComparer, 100)
	b.readOnly = true
	for _, kv := range jk.Data {
		b.data.Put(kv.Key, kv.Value)
	}
	return nil
}

func (b *blockData) EncodeRLP(w io.Writer) error {
	jk := new(blockWal)
	jk.BlockHash = b.BlockHash
	jk.ParentHash = b.ParentHash
	jk.BlockNumber = new(big.Int).Set(b.Number)
	jk.KvHash = b.kvHash
	jk.Data = make([]journalData, 0)
	if b.data.Size() != 0 {
		itr := b.data.NewIterator(nil)
		defer itr.Release()
		for itr.Next() {
			key, val := common.CopyBytes(itr.Key()), common.CopyBytes(itr.Value())
			jData := journalData{
				Key:   key,
				Value: val,
			}
			jk.Data = append(jk.Data, jData)
		}
	}
	return rlp.Encode(w, jk)
}

func (b *blockData) BlockKey() []byte {
	return EncodeWalKey(b.Number)
}

func (b *blockData) BlockVal() []byte {
	val, err := rlp.EncodeToBytes(b)
	if err != nil {
		panic("Encode BlockVal to byte fail:" + err.Error())
	}
	return val
}

func (b *blockData) cleanJournal() {
	b.journal = nil
	b.validRevisions = nil
	b.nextRevisionId = 0
}

func (b *blockData) revert(snapshot int) {
	for i := len(b.journal) - 1; i >= snapshot; i-- {
		// Undo the changes made by the operation
		en := b.journal[i]
		if en.oldValNotExist {
			b.data.Delete(en.key)
		} else {
			b.data.Put(en.key, en.oldVal)
		}
		b.kvHash = en.oldkvHash
	}
	b.journal = b.journal[:snapshot]
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (b *blockData) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(b.validRevisions), func(i int) bool {
		return b.validRevisions[i].id >= revid
	})
	if idx == len(b.validRevisions) || b.validRevisions[idx].id != revid {
		panic(fmt.Errorf("snapshotdb , revision id %v cannot be reverted", revid))
	}
	snapshot := b.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	b.revert(snapshot)
	b.validRevisions = b.validRevisions[:idx]
}

// Snapshot returns an identifier for the current revision of the state.
func (b *blockData) Snapshot() int {
	id := b.nextRevisionId
	b.nextRevisionId++
	b.validRevisions = append(b.validRevisions, revision{id, len(b.journal)})
	return id
}

func (b *blockData) Write(key, val []byte) error {
	var entry journalEntry = journalEntry{
		key:       key,
		newVal:    val,
		oldkvHash: b.kvHash,
	}
	if v, err := b.data.Get(key); err != nil {
		entry.oldVal = nil
		entry.oldValNotExist = true
	} else {
		entry.oldVal = v[:]
		entry.oldValNotExist = false
	}
	if err := b.data.Put(key, val); err != nil {
		return err
	}
	b.kvHash = generateKVHash(key, val, b.kvHash)

	// append inserts a new modification entry to the end of the change journal.
	b.journal = append(b.journal, entry)
	return nil
}

type journalEntry struct {
	key            []byte
	newVal         []byte
	oldValNotExist bool
	oldVal         []byte
	oldkvHash      common.Hash
}

type unCommitBlocks struct {
	blocks map[common.Hash]*blockData
	sync.RWMutex
}

func (u *unCommitBlocks) Get(key common.Hash) *blockData {
	u.RLock()
	block, ok := u.blocks[key]
	u.RUnlock()
	if !ok {
		return nil
	}
	return block
}

func (u *unCommitBlocks) Set(key common.Hash, block *blockData) {
	u.Lock()
	u.blocks[key] = block
	u.Unlock()
}
