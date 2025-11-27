package snapshotdb

import (
	"bytes"
	"encoding/binary"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/ethdb/leveldb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"
	"github.com/syndtr/goleveldb/leveldb/iterator"

	"math/big"

	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	vrfNoncePrefix     = []byte("vn")
	archiveBlockPrefix = []byte("abp")
	currentBlockKey    = []byte("cb")
	nonceStorageKey    = []byte("nonceStorageKey")
)

func VrfNonceKey(number uint64) []byte {
	var buf [10]byte
	copy(buf[:], vrfNoncePrefix)
	binary.BigEndian.PutUint64(buf[2:], number)
	return buf[:]
}
func ArchiveBlockKey(number uint64) []byte {
	var buf [11]byte
	copy(buf[:], archiveBlockPrefix)
	binary.BigEndian.PutUint64(buf[3:], number)
	return buf[:]
}

type VRFNonce struct {
	MaxValidatorNum uint32
	Nonce           []byte
}
type ArchiveBlock struct {
	Number uint64
	Root   common.Hash
	KvHash common.Hash
}

type archiveDB struct {
	db     ethdb.Database
	triedb *trie.Database
	trie   *trie.StateTrie
}

func OpenArchiveDB(path string, cache int, handles int) (*archiveDB, error) {

	db, err := leveldb.New(getArchiveDBPath(path), cache, handles, "", false)
	if err != nil {
		log.Error("Open archiveDB db failed", "err", err, "path", getArchiveDBPath(path))
		return nil, err
	}
	log.Error("Open archiveDB db succeed", "path", getArchiveDBPath(path))

	triedb := trie.NewDatabase(rawdb.NewDatabase(db))
	return &archiveDB{
		db:     rawdb.NewDatabase(db),
		triedb: triedb,
		trie:   nil,
	}, nil
}

// NewArchiveDBWithDB creates a new archiveDB instance with the provided database
func NewArchiveDBWithDB(db ethdb.Database) (*archiveDB, error) {
	triedb := trie.NewDatabase(db)
	return &archiveDB{
		db:     db,
		triedb: triedb,
		trie:   nil,
	}, nil
}

func (a *archiveDB) init(walk func(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error) error {
	num, err := a.CurrentBlock()
	if err != nil {
		return err
	}
	if num == nil {
		a.trie, _ = trie.NewStateTrie(trie.TrieID(common.Hash{}), a.triedb)
		walk(nil, func(num *big.Int, iter iterator.Iterator) error {
			total := 0
			size := 0
			for iter.Next() {
				total += 1
				size += len(iter.Key()) + len(iter.Value())
				a.trie.Update(iter.Key(), iter.Value())
			}
			root, set, err := a.trie.Commit(false)
			if err != nil {
				return err
			}
			log.Info("walk total:", "total", total, "size", size)

			nodes := trie.NewWithNodeSet(set)
			a.triedb.Update(nodes)
			a.triedb.Commit(root, false, false)
			a.SetCurrentBlock(num.Uint64())
			batch := a.db.NewBatch()
			a.SetArchiveBlock(batch, num.Uint64(), &ArchiveBlock{Number: num.Uint64(), Root: root, KvHash: common.Hash{}})
			batch.Write()
			log.Info("Init archiveDB db", "num", num, "root", root)
			return nil
		})
	} else {
		block, err := a.GetArchiveBlock(*num)
		if err != nil {
			return err
		}
		a.trie, _ = trie.NewStateTrie(trie.TrieID(block.Root), a.triedb)
		log.Info("Get archive block", "num", num, "root", block.Root)
	}
	return nil
}
func (a *archiveDB) CommitBlock(block *BlockData) error {
	itr := block.data.NewIterator(nil)
	defer itr.Release()
	oldRoot := a.trie.Hash()
	batch := a.db.NewBatch()
	total := 0
	for itr.Next() {
		if bytes.HasPrefix(itr.Key(), nonceStorageKey) {
			nonces := make([][]byte, 0)
			if err := rlp.DecodeBytes(itr.Value(), &nonces); nil != err {
				return err
			}
			a.SetVrfNonce(batch, block.Number.Uint64(), &VRFNonce{MaxValidatorNum: uint32(len(nonces)), Nonce: nonces[len(nonces)-1]})
		} else {
			if itr.Value() != nil {
				a.trie.Update(itr.Key(), itr.Value())
			} else {
				a.trie.Delete(itr.Key())
			}
		}
		total++
	}
	root, set, err := a.trie.Commit(false)
	if err != nil {
		return err
	}
	a.SetArchiveBlock(batch, block.Number.Uint64(), &ArchiveBlock{
		Number: block.Number.Uint64(),
		Root:   root,
		KvHash: block.kvHash,
	})
	a.SetCurrentBlock(block.Number.Uint64())
	nodes := trie.NewWithNodeSet(set)
	a.triedb.Update(nodes)
	a.triedb.Commit(root, false, false)
	batch.Write()
	a.triedb.IncrVersion()
	a.triedb.Reference(root, common.Hash{})
	a.triedb.DereferenceDB(oldRoot)
	insert, deletes := set.Size()
	log.Info("Commit archiveDB snapshot", "block", block.Number, "update", total, "treeinsert", insert, "treedelete", deletes)
	return nil
}
func (a *archiveDB) CurrentBlock() (*uint64, error) {
	val, _ := a.db.Get(currentBlockKey)
	if val == nil {
		return nil, nil
	}
	res := binary.BigEndian.Uint64(val)
	return &res, nil
}
func (a *archiveDB) SetCurrentBlock(block uint64) error {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], block)
	return a.db.Put(currentBlockKey, buf[:])
}

func (a *archiveDB) SetArchiveBlock(batch ethdb.Batch, blockNumber uint64, block *ArchiveBlock) error {
	value, err := rlp.EncodeToBytes(block)
	if err != nil {
		return err
	}
	batch.Put(ArchiveBlockKey(blockNumber), value)
	return nil
}
func (a *archiveDB) GetArchiveBlock(blockNumber uint64) (*ArchiveBlock, error) {
	val, err := a.db.Get(ArchiveBlockKey(blockNumber))
	if err != nil {
		return nil, err
	}
	var block ArchiveBlock
	err = rlp.DecodeBytes(val, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (a *archiveDB) SetVrfNonce(batch ethdb.Batch, blockNumber uint64, nonce *VRFNonce) error {
	value, err := rlp.EncodeToBytes(nonce)
	if err != nil {
		return err
	}
	batch.Put(VrfNonceKey(blockNumber), value)
	return nil
}
func (a *archiveDB) GetVrfNonce(blockNumber uint64) (*VRFNonce, error) {
	val, err := a.db.Get(VrfNonceKey(blockNumber))
	if err != nil {
		return nil, err
	}
	var vn VRFNonce
	err = rlp.DecodeBytes(val, &vn)
	if err != nil {
		return nil, err
	}
	return &vn, nil
}

func (a *archiveDB) GetVrfNonces(blockNumber uint64) ([][]byte, error) {
	vrfNonce, err := a.GetVrfNonce(blockNumber)
	if err != nil {
		return nil, err
	}
	nonces := make([][]byte, 0, vrfNonce.MaxValidatorNum)

	for i := blockNumber - uint64(vrfNonce.MaxValidatorNum) + 1; i < blockNumber; i++ {
		nonce, err := a.GetVrfNonce(i)
		if err != nil {
			return nil, err
		}
		nonces = append(nonces, nonce.Nonce)
	}
	nonces = append(nonces, vrfNonce.Nonce)
	return nonces, nil
}
func (a *archiveDB) SnapshotDB(blockNumber uint64) (DB, error) {
	log.Debug("Get archive snapshotdb", "blockNumber", blockNumber)
	block, err := a.GetArchiveBlock(blockNumber)
	if err != nil {
		return nil, err
	}
	snapTree, err := trie.NewStateTrie(trie.TrieID(block.Root), a.triedb)
	if err != nil {
		return nil, err
	}
	nonces, err := a.GetVrfNonces(blockNumber)
	if err != nil {
		return nil, err
	}
	vrfNonce, err := rlp.EncodeToBytes(nonces)
	if err != nil {
		return nil, err
	}
	return &archiveSnapshot{
		trie:        snapTree,
		blockNumber: blockNumber,
		kvHash:      block.KvHash,
		vrfNonce:    vrfNonce,
	}, nil
}
