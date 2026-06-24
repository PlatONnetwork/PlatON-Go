package snapshotdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/VictoriaMetrics/fastcache"
	lru "github.com/hashicorp/golang-lru"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/ethdb/leveldb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"

	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	defaultCapNodePercent = common.StorageSize(1) / 4
	defaultPreimageCache  = 32
	defaultRankingCache   = 8
	vrfNoncePrefix        = []byte("vn")
	archiveBlockPrefix    = []byte("abp")
	currentBlockKey       = []byte("cb")
	nonceStorageKey       = []byte("nonceStorageKey")
	ErrSnapshotArchiveDB  = errors.New("snapshot archive database not found")
)

func RankingCacheKey(root common.Hash, key []byte, ranges int) []byte {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(ranges))
	return append(append(root.Bytes(), key[:]...), buf[:]...)
}
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

type PreimageCache struct {
	db    ethdb.KeyValueReader
	cache *fastcache.Cache
}

func NewPreimageCache(db ethdb.KeyValueReader, mb int) *PreimageCache {
	return &PreimageCache{
		db:    db,
		cache: fastcache.New(mb * 1024 * 1024),
	}
}

func (pc *PreimageCache) Get(hash common.Hash) []byte {
	value := pc.cache.Get(nil, hash.Bytes())
	if len(value) != 0 {
		return value
	}
	value = rawdb.ReadPreimage(pc.db, hash)
	if len(value) != 0 {
		pc.cache.Set(hash.Bytes(), value)
	}
	return value
}

type MemDB struct {
	sync.Mutex
	*memdb.DB
	ref int
}

func NewMemDB(db *memdb.DB) *MemDB {
	return &MemDB{
		DB: db,
	}
}

func (m *MemDB) Ref() {
	m.Lock()
	defer m.Unlock()
	m.ref++
}

func (m *MemDB) Deref() {
	m.Lock()
	defer m.Unlock()
	m.ref--
}

func (m *MemDB) Release() {
	m.Lock()
	defer m.Unlock()
	m.ref--
	if m.ref == 0 {
		returnMdbPool(m)
	}
}

type RankingCache struct {
	sync.Mutex
	rankingCache *lru.Cache
}

func NewRankingCache(num int) *RankingCache {
	cache, _ := lru.NewWithEvict(num, func(key, value interface{}) {
		value.(*MemDB).Release()
	})
	return &RankingCache{rankingCache: cache}
}

func (rc *RankingCache) Add(key []byte, db *MemDB) {
	rc.Lock()
	defer rc.Unlock()
	db.Ref()
	rc.rankingCache.Add(string(key), db)
}

func (rc *RankingCache) Get(key []byte) *MemDB {
	rc.Lock()
	defer rc.Unlock()
	v, ok := rc.rankingCache.Get(string(key))
	if ok {
		db := v.(*MemDB)
		db.Ref()
		log.Debug("RankingCache hit", "key", string(key))
		return db
	}
	return nil
}

type archiveDB struct {
	db            ethdb.Database
	preimageCache *PreimageCache
	rankingCache  *RankingCache
	tracedb       *trie.Database
	triedb        *trie.Database
	trie          *trie.StateTrie
}

func OpenArchiveDB(path string, cache int, handles int) (*archiveDB, error) {

	db, err := leveldb.New(getArchiveDBPath(path), cache, handles, "", false)
	if err != nil {
		log.Error("Open archive db failed", "err", err, "path", getArchiveDBPath(path))
		return nil, err
	}
	log.Error("Open archive db succeed", "path", getArchiveDBPath(path))

	// Set read options to not fill cache for tracedb operations
	db.SetReadOptions(&opt.ReadOptions{
		DontFillCache: true,
	})

	triedb := trie.NewDatabaseWithConfig(rawdb.NewDatabase(db), &trie.Config{
		Cache:     archiveDatabaseCache,
		Preimages: true,
	})

	tracedb := trie.NewDatabaseWithConfig(rawdb.NewDatabase(db), &trie.Config{
		Cache:     archiveDatabaseCache,
		Preimages: false,
	})

	return &archiveDB{
		db:            rawdb.NewDatabase(db),
		preimageCache: NewPreimageCache(db, defaultPreimageCache),
		rankingCache:  NewRankingCache(defaultRankingCache),
		tracedb:       tracedb,
		triedb:        triedb,
		trie:          nil,
	}, nil
}

// NewArchiveDBWithDB creates a new archiveDB instance with the provided database
func NewArchiveDBWithDB(db ethdb.Database) (*archiveDB, error) {
	return &archiveDB{
		db:            db,
		preimageCache: NewPreimageCache(db, defaultPreimageCache),
		triedb:        trie.NewDatabaseWithConfig(db, &trie.Config{Cache: 16, Preimages: true}),
		tracedb:       trie.NewDatabaseWithConfig(db, &trie.Config{Cache: 16, Preimages: true}),
		trie:          nil,
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
				a.trie.MustUpdate(common.CopyBytes(iter.Key()), common.CopyBytes(iter.Value()))
			}
			root, set, err := a.trie.Commit(false)
			if err != nil {
				return err
			}
			log.Info("walk total:", "total", total, "size", size)

			nodes := trie.NewWithNodeSet(set)
			a.triedb.Update(root, types.EmptyRootHash, 0, nodes)
			a.triedb.Commit(root, false, false)
			a.SetCurrentBlock(num.Uint64())
			batch := a.db.NewBatch()
			a.SetArchiveBlock(batch, num.Uint64(), &ArchiveBlock{Number: num.Uint64(), Root: root, KvHash: common.Hash{}})
			batch.Write()
			log.Info("Init archive db", "num", num, "root", root)
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
	currentBlock, err := a.CurrentBlock()
	if err != nil || currentBlock == nil {
		panic(fmt.Sprintf("get current failed:%v, block:%v", err, currentBlock))
	}
	if *currentBlock+1 < block.Number.Uint64() {
		panic(fmt.Sprintf("blockdata too far away from current block, current:%d, block:%d", *currentBlock, block.Number.Uint64()))
	}

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
				a.trie.MustUpdate(common.CopyBytes(itr.Key()), common.CopyBytes(itr.Value()))
			} else {
				a.trie.MustDelete(common.CopyBytes(itr.Key()))
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
	a.triedb.Update(root, oldRoot, 0, nodes)
	a.triedb.Commit(root, false, true)
	batch.Write()
	a.triedb.IncrVersion()
	a.triedb.ReferenceVersion(root)
	a.triedb.Dereference(oldRoot)

	size, _ := a.triedb.Size()
	limit := common.StorageSize(archiveTrieOversizeThreshold) * 1024 * 1024
	oversize := size > limit
	if oversize {
		log.Info("Trie oversize, need to reset", "size", size, "threshold", limit)
		a.triedb.CapNode(limit * defaultCapNodePercent)
		a.triedb.ResetUseless()
	}
	insert, deletes := set.Size()
	log.Info("Commit archive snapshot", "block", block.Number, "update", total, "treeinsert", insert, "treedelete", deletes)
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
	if a == nil {
		return nil, ErrSnapshotArchiveDB
	}
	log.Debug("Get archive snapshotdb", "blockNumber", blockNumber)
	block, err := a.GetArchiveBlock(blockNumber)
	if err != nil {
		return nil, err
	}
	snapTree, err := trie.NewStateTrie(trie.TrieID(block.Root), a.tracedb)
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
		dbreader:      a.db,
		preimageCache: a.preimageCache,
		rankingCache:  a.rankingCache,
		trie:          snapTree,
		blockNumber:   blockNumber,
		kvHash:        block.KvHash,
		vrfNonce:      vrfNonce,
	}, nil
}
