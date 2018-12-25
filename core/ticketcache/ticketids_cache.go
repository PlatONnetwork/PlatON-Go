package ticketcache

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/crypto"
	"Platon-go/ethdb"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"errors"
	"github.com/golang/protobuf/proto"
	"math/big"
	"sort"
	"sync"
	"time"
)

type TicketCache map[discover.NodeID][]common.Hash

type Timer struct {
	start time.Time
}

func (t *Timer) Begin()  {
	t.start = time.Now()
	//fmt.Println("Begin=> ", "now: ", time.Now().Nanosecond(), " tCalc: ", t.tCalc)
}

func (t *Timer) End() float64 {
	//fmt.Println("End=> ", "now: ", time.Now().Nanosecond(), " tCalc: ", t.tCalc)
	tns := time.Since(t.start).Nanoseconds()
	tms := float64(tns)/float64(1e6)
	return tms

}

var (
	//error def
	ErrNotfindFromNodeId = errors.New("Not find tickets from node id")
	ErrProbufMarshal = errors.New("protocol buffer Marshal faile")
	ErrLeveldbPut = errors.New("level db put faile")

	//const def
	ticketPoolCacheKey = []byte("ticketPoolCache")
)

var ticketidsCache *NumBlocks

func NewTicketIdsCache(db ethdb.Database)  *NumBlocks {
	/*
		Put 购票交易新增选票
		Del 节点掉榜，选票过期，选票被选中
	*/
	//logInfo("NewTicketIdsCache==> Init ticketidsCache call NewTicketIdsCache func")
	timer := Timer{}
	timer.Begin()
	if nil != ticketidsCache {
		return ticketidsCache
	}
	ticketidsCache = &NumBlocks{}
	ticketidsCache.NBlocks = make(map[string]*BlockNodes)
	cache, err := db.Get(ticketPoolCacheKey)
	if err == nil {
		log.Info("NewTicketIdsCache==> ", "CacheHex: ", hexutil.Encode(cache))
		if err := proto.Unmarshal(cache, ticketidsCache); err != nil {
			log.Error("NewTicketIdsCache==> protocol buffer Unmarshal faile")
		}
	}
	log.Info("NewTicketIdsCache==> ", "ms: ", timer.End())
	return ticketidsCache
}

func GetNodeTicketsCacheMap(blocknumber *big.Int, blockhash common.Hash) (ret TicketCache) {
	log.Info("GetNodeTicketsCacheMap==> ", "blocknumber: ", blocknumber, " blockhash: ", blockhash.Hex())
	if ticketidsCache!=nil {
		ret = ticketidsCache.GetNodeTicketsMap(blocknumber, blockhash)
	}else {
		log.Error("GetNodeTicketsCacheMap==> ticketidsCache instance is nil!")
	}
	return
}

func GetTicketidsCachePtr() *NumBlocks {
	return ticketidsCache
}

////////////////////////////
func (nb *NumBlocks) Hash(cache TicketCache) (common.Hash, error) {
	timer := Timer{}
	timer.Begin()
	out, err := proto.Marshal(cache.GetSortStruct())
	log.Info("Hash==> ", "lenOut: ", len(out), " hexOut: ", hexutil.Encode(out))
	if err != nil {
		log.Error("Hash==> ", "ErrProbufMarshal: ", ErrProbufMarshal.Error())
		return common.Hash{}, ErrProbufMarshal
	}
	ret := crypto.Keccak256Hash(out)
	log.Info("Hash==> ", "run time ",  " ms: ", timer.End())
	return ret, nil
}

func (nb *NumBlocks) GetNodeTicketsMap(blocknumber *big.Int, blockhash common.Hash) TicketCache{
	log.Info("GetNodeTicketsMap==> ", "blocknumber: ", blocknumber, " blockhash: ", blockhash.Hex())
	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		blockNodes = &BlockNodes{}
		blockNodes.BNodes = make(map[string]*NodeTicketIds)
		nb.NBlocks[blocknumber.String()] = blockNodes
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if !ok {
		nodeTicketIds = &NodeTicketIds{}
		nodeTicketIds.NTickets = make(map[string]*TicketIds)
		blockNodes.BNodes[blockhash.String()] = nodeTicketIds
	}
	//go thread
	type result struct {
		key discover.NodeID
		tids []common.Hash
	}
	resCh := make(chan *result, len(nodeTicketIds.NTickets))
	var wg sync.WaitGroup
	wg.Add(len(nodeTicketIds.NTickets))
	for k, v := range nodeTicketIds.NTickets {
		nid, err := discover.HexID(k)
		if err == nil {
			go func  (nodeid discover.NodeID, tidslice *TicketIds){
				tids := make([]common.Hash, 0, len(tidslice.TicketId))
				for _, tid := range tidslice.TicketId {
					tids = append(tids, common.BytesToHash(tid))
				}
				res := new(result)
				res.key = nodeid
				res.tids = tids
				resCh <- res
				wg.Done()
			}(nid, v)
		} else {
			log.Error("GetNodeTicketsMap==> discover.HexID error ", "hex: ", k)
		}
	}
	wg.Wait()
	close(resCh)
	out := NewTicketCache()
	for res := range resCh {
		out[res.key] = res.tids
	}
	return out
}

func (nb *NumBlocks) Submit2Cache(blocknumber *big.Int, blockhash common.Hash, in map[discover.NodeID][]common.Hash) {
	log.Info("Submit2Cache==> ", "blocknumber: ", blocknumber, " blockhash: ", blockhash.Hex())
	blockNodes, ok := nb.NBlocks[blocknumber.String()];
	if !ok {
		blockNodes = &BlockNodes{}
		blockNodes.BNodes = make(map[string]*NodeTicketIds)
	}
	//The same block hash data will be overwritten
	nodeTicketIds := &NodeTicketIds{}
	nodeTicketIds.NTickets = make(map[string]*TicketIds)
	//go thread
	type result struct {
		key discover.NodeID
		value *TicketIds
	}
	resCh := make(chan *result, len(in))
	var wg sync.WaitGroup
	wg.Add(len(in))
	for k, v := range in {
		go func  (key discover.NodeID, val []common.Hash){
			tIds := &TicketIds{}
			for _, va := range val {
				tIds.TicketId = append(tIds.TicketId, va.Bytes())
			}
			res := new(result)
			res.key = key
			res.value = tIds
			resCh <- res
			wg.Done()
		}(k, v)
	}
	wg.Wait()
	close(resCh)
	for res := range resCh {
		nodeTicketIds.NTickets[res.key.String()] = res.value
	}
	blockNodes.BNodes[blockhash.String()] = nodeTicketIds
	nb.NBlocks[blocknumber.String()] = blockNodes
}

func (nb *NumBlocks) Commit(db ethdb.Database) error {
	timer := Timer{}
	timer.Begin()
	out, err := proto.Marshal(nb)
	if err != nil {
		log.Error("Commit==> ","ErrProbufMarshal: ", err.Error())
		return ErrProbufMarshal
	}
	//logInfo("Marshal out: ", hexutil.Encode(out))
	log.Info("Commit==> ", "outlen: ", len(out), " outhex: ", hexutil.Encode(out))
	if err := db.Put(ticketPoolCacheKey, out); err != nil  {
		log.Error("level db put faile: ", err.Error())
		return ErrLeveldbPut
	}
	log.Info("Commit==> run time ", "ms: ", timer.End())
	return nil
}

///////////////////////////////////

func NewTicketCache() TicketCache  {
	return make(map[discover.NodeID][]common.Hash)
}

func (tc TicketCache) AppendTicketCache (nodeid discover.NodeID, tids []common.Hash){
	value, ok := tc[nodeid]
	if !ok {
		value = make([]common.Hash, 0)
	}
	for _, id := range tids {
		value = append(value, id)
	}
	tc[nodeid] = value
}

func (tc TicketCache) GetTicketCache(nodeid discover.NodeID) ([]common.Hash, error) {
	tids, ok := tc[nodeid]
	if !ok {
		return nil, ErrNotfindFromNodeId
	}
	ret := make([]common.Hash, 0, len(tids))
	copy(ret, tids)
	return ret, nil
}

func (tc TicketCache) RemoveTicketCache(nodeid discover.NodeID, tids []common.Hash) error {
	cache, ok := tc[nodeid]
	if !ok {
		return ErrNotfindFromNodeId
	}
	mapTIds := make(map[common.Hash]common.Hash)
	for _, id := range tids {
		mapTIds[id] = id
	}
	for i := 0; i < len(cache); i++ {
		if _, ok := mapTIds[cache[i]]; ok {
			cache = append(cache[:i], cache[i+1:]...)
			i = i - 1
		}
	}
	tc[nodeid] = cache
	return nil
}

func (tc TicketCache) TCount(nodeid discover.NodeID) uint64 {
	count := uint64(len(tc[nodeid]))
	return count
}

func (tc TicketCache) TicketCaceheSnapshot() TicketCache {
	ret := NewTicketCache()
	for nodeid, tids := range tc {
		ret[nodeid] = make([]common.Hash, 0, len(tids))
		copy(ret[nodeid], tids)
	}
	return ret
}

func (tc TicketCache) GetSortStruct() *SortCalcHash {
	sc := &SortCalcHash{}
	sc.Nodeids = make([]string, 0, len(tc))
	sc.Tids = make([]*TicketIds, 0, len(tc))
	for k := range tc {
		sc.Nodeids = append(sc.Nodeids, k.String())
	}
	sort.Strings(sc.Nodeids)
	for _, k := range sc.Nodeids {
		nodeid , err := discover.HexID(k)
		if err == nil {
			tids := &TicketIds{}
			tids.TicketId = make([][]byte, 0, len(tc[nodeid]))
			for _, tid := range tc[nodeid] {
				tids.TicketId = append(tids.TicketId, tid.Bytes())
			}
			sc.Tids = append(sc.Tids, tids)
		} else {
			log.Error("GetSortStruct==> discover.HexID error ", "hex: ", k)
		}
	}
	return sc
}
