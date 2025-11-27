package snapshotdb

import (
	"bytes"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type archiveSnapshot struct {
	trie        *trie.StateTrie
	blockNumber uint64
	kvHash      common.Hash
	vrfNonce    []byte
}

func (a *archiveSnapshot) Put(hash common.Hash, key, value []byte) error {
	a.trie.Update(key, value)
	return nil
}

func (a archiveSnapshot) NewBlock(blockNumber *big.Int, parentHash common.Hash, hash common.Hash) error {
	//TODO implement me
	panic("implement me")
}

func (a archiveSnapshot) Get(hash common.Hash, key []byte) ([]byte, error) {
	if bytes.HasPrefix(key, nonceStorageKey) {
		return a.vrfNonce, nil
	}
	log.Debug("Get from archive snapshot", "hash", hash, "key", common.Bytes2Hex(key))
	v, err := a.trie.TryGet(key)
	log.Debug("Try get from archive snapshot", "key", common.Bytes2Hex(key), "v", common.Bytes2Hex(v), "err", err)
	return v, err
}

func (a archiveSnapshot) GetFromCommittedBlock(key []byte) ([]byte, error) {
	return a.trie.Get(key), nil
}

func (a *archiveSnapshot) Del(hash common.Hash, key []byte) error {
	a.trie.Delete(key)
	return nil
}

func (a archiveSnapshot) Has(hash common.Hash, key []byte) (bool, error) {
	return len(a.trie.Get(key)) != 0, nil
}

func (a archiveSnapshot) Flush(hash common.Hash, blocknumber *big.Int) error {
	return nil
}

func (a archiveSnapshot) Ranking(hash common.Hash, key []byte, ranges int) iterator.Iterator {
	//TODO implement me
	panic("implement me")
}

func (a archiveSnapshot) WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
	panic("unsupported")
}

func (a archiveSnapshot) WalkDB(num uint64, f func(baseBlock uint64, iter iterator.Iterator, blocks []rlp.RawValue) error) error {
	panic("unsupported")

}

func (a archiveSnapshot) Commit(hash common.Hash) error {
	//do nothing
	return nil
}

func (a archiveSnapshot) Clear() error {
	//do nothing
	return nil
}

func (a archiveSnapshot) PutBaseDB(key, value []byte) error {
	panic("unsupported")
}

func (a archiveSnapshot) GetBaseDB(key []byte) ([]byte, error) {
	panic("unsupported")
}

func (a archiveSnapshot) DelBaseDB(key []byte) error {
	panic("unsupported")
}

func (a archiveSnapshot) WriteBaseDB(kvs [][2][]byte) error {
	panic("unsupported")
}

func (a archiveSnapshot) WriteBaseDBWithBlock(current *types.Header, blocks []BlockData) error {
	panic("unsupported")
}

func (a archiveSnapshot) SetCurrent(highestHash common.Hash, base, height big.Int) error {
	panic("unsupported")
}

func (a archiveSnapshot) GetCurrent() *current {
	panic("unsupported")
}

func (a *archiveSnapshot) GetLastKVHash(blockHash common.Hash) []byte {
	return a.kvHash.Bytes()
}

func (a *archiveSnapshot) BaseNum() (*big.Int, error) {
	panic("unsupported")
}

func (a archiveSnapshot) Close() error {
	return nil
}

func (a archiveSnapshot) Compaction() error {
	return nil
}

func (a archiveSnapshot) SetEmpty() error {
	panic("unsupported")
}

func (a archiveSnapshot) RevertToSnapshot(hash common.Hash, revid int) {
	//do nothind
}

func (a archiveSnapshot) Snapshot(hash common.Hash) int {
	return 0
}
