package snapshotdb

import (
	"bytes"
	"container/heap"
	"math/big"
	"sync"

	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"
)

var (
	mdbPool = sync.Pool{
		New: func() interface{} { return NewMemDB(memdb.New(DefaultComparer, 256)) },
	}
)

func newMdb() *MemDB {
	return mdbPool.Get().(*MemDB)
}

func returnMdbPool(h *MemDB) {
	h.Reset()
	mdbPool.Put(h)
}

// trieKeyValueIterator adapts trie.NodeIterator to iterator.Iterator interface
// It only returns leaf nodes that match the given prefix
type trieKeyValueIterator struct {
	preimageCache *PreimageCache
	nodeIt        trie.NodeIterator
	prefix        []byte
	key           []byte
	value         []byte
	valid         bool
	err           error
}

func (t *trieKeyValueIterator) First() bool {
	panic("unsupported")
}

func (t *trieKeyValueIterator) Last() bool {
	panic("unsupported")
}

func (t *trieKeyValueIterator) Seek(key []byte) bool {
	panic("unsupported")

}

func (t *trieKeyValueIterator) Prev() bool {
	panic("unsupported")
}

func (t *trieKeyValueIterator) SetReleaser(releaser util.Releaser) {
	panic("unsupported")
}

func (t *trieKeyValueIterator) Valid() bool {
	panic("unsupported")
}

func (t *trieKeyValueIterator) Next() bool {
	for t.nodeIt.Next(true) {
		if t.nodeIt.Leaf() {
			key := t.nodeIt.LeafKey()
			origin := t.preimageCache.Get(common.BytesToHash(key))
			// Filter by prefix
			if bytes.HasPrefix(origin, t.prefix) {
				t.key = origin
				t.value = t.nodeIt.LeafBlob()
				t.valid = true
				return true
			}
		}
	}
	t.valid = false
	t.err = t.nodeIt.Error()
	return false
}

func (t *trieKeyValueIterator) Key() []byte {
	if !t.valid {
		return nil
	}
	return t.key
}

func (t *trieKeyValueIterator) Value() []byte {
	if !t.valid {
		return nil
	}
	return t.value
}

func (t *trieKeyValueIterator) Error() error {
	return t.err
}

func (t *trieKeyValueIterator) Release() {
	t.key = nil
	t.value = nil
	t.valid = false
	t.err = nil
	t.prefix = nil
	t.nodeIt = nil
}

type mdbIterator struct {
	iterator.Iterator
	mdb *MemDB
}

func (m *mdbIterator) Release() {
	m.Iterator.Release()
	m.mdb.Release()
}

type archiveSnapshot struct {
	dbreader      ethdb.KeyValueReader
	preimageCache *PreimageCache
	rankingCache  *RankingCache
	trie          *trie.StateTrie
	blockNumber   uint64
	kvHash        common.Hash
	vrfNonce      []byte
}

func (a *archiveSnapshot) Put(hash common.Hash, key, value []byte) error {
	a.trie.Update(key, value)
	return nil
}

func (a archiveSnapshot) NewBlock(blockNumber *big.Int, parentHash common.Hash, hash common.Hash) error {
	//TODO implement me
	panic("implement me")
}

func (a *archiveSnapshot) Get(hash common.Hash, key []byte) ([]byte, error) {
	if bytes.HasPrefix(key, nonceStorageKey) {
		return a.vrfNonce, nil
	}
	log.Debug("Get from archive snapshot", "hash", hash, "key", common.Bytes2Hex(key))
	v, err := a.trie.TryGetStorage(common.Address{}, key)
	if len(v) == 0 && err == nil {
		return nil, ErrNotFound
	}
	log.Debug("Try get from archive snapshot", "key", common.Bytes2Hex(key), "v", common.Bytes2Hex(v), "err", err)
	return v, err
}

func (a *archiveSnapshot) GetFromCommittedBlock(key []byte) ([]byte, error) {
	v, err := a.trie.TryGetStorage(common.Address{}, key)
	if len(v) == 0 && err == nil {
		return nil, ErrNotFound
	}
	return v, err
}

func (a *archiveSnapshot) Del(hash common.Hash, key []byte) error {
	a.trie.Delete(key)
	return nil
}

func (a *archiveSnapshot) Has(hash common.Hash, key []byte) (bool, error) {
	v, err := a.trie.TryGetStorage(common.Address{}, key)
	if len(v) == 0 && err == nil {
		return true, ErrNotFound
	}
	return len(v) != 0, err
}

func (a archiveSnapshot) Flush(hash common.Hash, blocknumber *big.Int) error {
	return nil
}

func (a *archiveSnapshot) Ranking(hash common.Hash, key []byte, ranges int) iterator.Iterator {
	rangingKey := RankingCacheKey(a.trie.Hash(), key, ranges)
	if db := a.rankingCache.Get(rangingKey); db != nil {
		return &mdbIterator{
			db.NewIterator(nil),
			db,
		}
	}

	log.Debug("Archive ranking", "key", hexutil.Encode(key))
	// Create a trie iterator to traverse all nodes
	nodeIt := a.trie.NodeIterator(nil)

	// Create ranking heap for sorting and limiting results
	rankingHeap := newRankingHeap(ranges)
	// Create a custom iterator that converts trie nodes to key-value pairs
	trieIter := &trieKeyValueIterator{
		preimageCache: a.preimageCache,
		nodeIt:        nodeIt,
		prefix:        key,
	}

	// Add trie key-value pairs to heap
	rankingHeap.itr2Heap(trieIter, true, false)
	trieIter.Release()
	defer rankingHeap.release()
	// Create memdb to store sorted results
	mdb := newMdb()
	for rankingHeap.heap.Len() > 0 {
		kv := heap.Pop(&rankingHeap.heap).(kv)
		if err := mdb.Put(kv.key, kv.value); err != nil {
			returnMdbPool(mdb)
			return iterator.NewEmptyIterator(err)
		}
	}

	a.rankingCache.Add(rangingKey, mdb)
	mdb.Ref()
	return &mdbIterator{
		mdb.NewIterator(nil),
		mdb,
	}
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
