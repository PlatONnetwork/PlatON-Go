package vrfstatistics

import (
	"bytes"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
)

var (
	KeySumOfRandomNum = []byte("asum")
	KeyRandomNumTxs   = []byte("btxs")
	Prefix            = "randomNum-"
	Tool              *Statistics
)

func init() {
	Tool = new(Statistics)
	Tool.blocks = make(map[uint64]*RandomNumRequest)
}

func EncodeKeyRandomNumTxs(block uint64) []byte {
	return append(KeyRandomNumTxs, common.Uint64ToBytes(block)...)
}

func DecodeKeyRandomNumTxs(key []byte) uint64 {
	return common.BytesToUint64(key[len(KeyRandomNumTxs):])
}

type Statistics struct {
	blocks map[uint64]*RandomNumRequest
	mu     sync.RWMutex
}

func (c *Statistics) AddRequest(block uint64, seedNum uint64, txhash common.Hash, sender common.Address) {
	c.mu.Lock()
	defer c.mu.Unlock()
	request, ok := c.blocks[block]
	if ok {
		if _, ok := request.txHashs[txhash]; ok {
			return
		}
		request.seedNums += seedNum
		request.txs = append(request.txs, TxInfo{txhash, sender})
		request.txHashs[txhash] = struct{}{}
		return
	}
	data := new(RandomNumRequest)
	data.seedNums = seedNum
	data.txs = make([]TxInfo, 0)
	data.txs = append(data.txs, TxInfo{txhash, sender})
	data.txHashs = make(map[common.Hash]struct{})
	data.txHashs[txhash] = struct{}{}
	c.blocks[block] = data
}

func (c *Statistics) Save(block uint64, database ethdb.KeyValueStore) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	info, ok := c.blocks[block]
	if !ok {
		return nil
	}
	val, err := database.Get(KeySumOfRandomNum)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}

	batch := database.NewBatch()
	if err := batch.Put(KeySumOfRandomNum, common.Uint64ToBytes(info.seedNums+common.BytesToUint64(val))); err != nil {
		return err
	}
	txdatas, err := rlp.EncodeToBytes(info.txs)
	if err != nil {
		return err
	}

	if err := batch.Put(EncodeKeyRandomNumTxs(block), txdatas); err != nil {
		return err
	}

	if err := batch.Write(); err != nil {
		return err
	}
	delete(c.blocks, block)
	return nil
}

func (c *Statistics) SumOfRandomNum(database ethdb.KeyValueReader) (uint64, error) {
	val, err := database.Get(KeySumOfRandomNum)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}
	return common.BytesToUint64(val), nil
}

func (c *Statistics) GetRandomNumberTxs(from, to uint64, database ethdb.KeyValueStore) (map[uint64][]TxInfo, error) {
	iter := database.NewIteratorWithStart(EncodeKeyRandomNumTxs(from))
	defer iter.Release()
	result := make(map[uint64][]TxInfo)
	for iter.Next() {
		key, value := iter.Key(), iter.Value()
		if !bytes.HasPrefix(key, KeyRandomNumTxs) {
			break
		}
		block := DecodeKeyRandomNumTxs(key)
		if block > to {
			break
		}
		var txinfos []TxInfo
		if err := rlp.DecodeBytes(value, &txinfos); err != nil {
			return nil, err
		}
		result[block] = txinfos
	}
	return result, nil
}

type RandomNumRequest struct {
	seedNums uint64
	txs      []TxInfo
	txHashs  map[common.Hash]struct{}
}

type TxInfo struct {
	TxHash common.Hash    `json:"txhash"`
	Sender common.Address `json:"sender"`
}
