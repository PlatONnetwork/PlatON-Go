package ticketcache

import (
	"Platon-go/common"
	"Platon-go/crypto"
	"Platon-go/ethdb"
	"Platon-go/log"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/big"
	"sort"
	"sync"
)

var (
	//error def
	ErrNotfindFromblockNumber = errors.New("Not find tickets from block number")
	ErrNotfindFromblockHash = errors.New("Not find tickets from block hash")
	ErrProbufMarshal = errors.New("protocol buffer Marshal faile")
	ErrLeveldbPut = errors.New("level db put faile")

	//const def
	ticketPoolCacheKey = []byte("ticketPoolCache")
)

var ticketidsCache *NumBlocks

func logError(msg string, ctx ...interface{})  {
	args := []interface{}{msg}
	args = append(args, ctx...)
	fmt.Println(args...)
	//log.Error(msg, ctx...)
}

func logInfo(msg string, ctx ...interface{})  {
	args := []interface{}{msg}
	args = append(args, ctx...)
	fmt.Println(args...)
	//log.logInfo(msg, ctx...)
}

func GetNodeTicketsCacheMap(blocknumber *big.Int, blockhash common.Hash) (ret map[string][]common.Hash) {
	logInfo("GetNodeTicketsCacheMap==> ", blocknumber, "  ", blockhash.Hex())
	if ticketidsCache!=nil {
		ret = ticketidsCache.GetNodeTicketsMap(blocknumber, blockhash)
	}else {
		logError("ticketidsCache instance is nil!")
		log.Error("ticketidsCache is Empty no GetNodeTicketsCacheMap ... ")
	}
	return
}

func NewTicketIdsCache(db ethdb.Database)  *NumBlocks {
	/*
		Put 购票交易新增选票
		Del 节点掉榜，选票过期，选票被选中
	*/
	logInfo("NewTicketIdsCache==> in")
	log.Info("Init ticketidsCache call NewTicketIdsCache func...")
	if nil != ticketidsCache {
		return ticketidsCache
	}
	ticketidsCache = &NumBlocks{}
	ticketidsCache.NBlocks = make(map[string]*BlockNodes)
	cache, err := db.Get(ticketPoolCacheKey)
	if err == nil {
		//logInfo("NewTicketIdsCache==> Get db cache hex: ", hexutil.Encode(cache))
		if err := proto.Unmarshal(cache, ticketidsCache); err != nil {
			//logError("protocol buffer Unmarshal faile hex: ", hexutil.Encode(cache))
			logError("protocol buffer Unmarshal faile hex: &&&&&&&&&&&&&&&")
		}
	}
	return ticketidsCache
}

func (nb *NumBlocks) Hash(blocknumber *big.Int, blockhash common.Hash) (common.Hash, error) {

	logInfo("Hash==> ", blocknumber, "  ", blockhash.Hex())
	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		logError(ErrNotfindFromblockNumber.Error())
		return common.Hash{}, ErrNotfindFromblockNumber
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if !ok {
		logError(ErrNotfindFromblockHash.Error())
		return common.Hash{}, ErrNotfindFromblockHash
	}
	out, err := proto.Marshal(getSortStruct(nodeTicketIds.NTickets))
	if err != nil {
		logError(ErrProbufMarshal.Error())
		return common.Hash{}, ErrProbufMarshal
	}
	ret := crypto.Keccak256Hash(out)
	logInfo("Hash==> output: ", ret.Hex())
	return ret, nil
}

func (nb *NumBlocks) GetNodeTicketsMap(blocknumber *big.Int, blockhash common.Hash) map[string][]common.Hash{

	//logInfo("GetNodeTicketsMap==> ", blocknumber, "  ", blockhash.Hex())
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
	out := make(map[string][]common.Hash)
	for k, v := range nodeTicketIds.NTickets{
		tids := make([]common.Hash, 0)
		for _, t := range v.TicketId {
			tid := common.Hash{}
			tid.SetBytes(t)
			tids = append(tids, tid)
		}
		out[k] = tids
	}
	return out
}

func (nb *NumBlocks) Submit2Cache(blocknumber *big.Int, blockhash common.Hash, in map[string][]common.Hash) {

	//logInfo("Submit2Cache==> ", blocknumber, "  ", blockhash.Hex())
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
		key string
		value *TicketIds
	}
	resCh := make(chan *result, len(in))
	var wg sync.WaitGroup
	wg.Add(len(in))
	for k, v := range in {
		go func  (key string, val []common.Hash){
			tIds := &TicketIds{}
			for _, va := range v {
				tIds.TicketId = append(tIds.TicketId, va.Bytes())
			}
			res := new(result)
			res.key = k
			res.value = tIds
			resCh <- res
			wg.Done()
		}(k, v)
	}
	wg.Wait()
	close(resCh)
	for res := range resCh {
		nodeTicketIds.NTickets[res.key] = res.value
	}

	//not thread
	//for k, v := range in {
	//	tIds := &TicketIds{}
	//	for _, va := range v {
	//		tIds.TicketId = append(tIds.TicketId, va.Bytes())
	//	}
	//	nodeTicketIds.NTickets[k] = tIds
	//}

	blockNodes.BNodes[blockhash.String()] = nodeTicketIds
	nb.NBlocks[blocknumber.String()] = blockNodes
}

func (nb *NumBlocks) Commit(db ethdb.Database) error {

	logInfo("Commit==> in")
	out, err := proto.Marshal(nb)
	if err != nil {
		logError("Protocol buffer failed to marshal :", nb, " err: ", err.Error())
		return ErrProbufMarshal
	}
	//logInfo("Marshal out: ", hexutil.Encode(out))
	logInfo("Marshal out len: ", len(out))
	if err := db.Put(ticketPoolCacheKey, out); err != nil  {
		logError("level db put faile: ", err.Error())
		return ErrLeveldbPut
	}
	return nil
}

func GetTicketidsCachePtr() *NumBlocks {
	return ticketidsCache
}

func getSortStruct(NTickets map[string]*TicketIds) *SortCalcHash {
	sc := &SortCalcHash{}
	sc.Nodeids = make([]string, 0, len(NTickets))
	sc.Tids = make([]*TicketIds, 0, len(NTickets))
	for k := range NTickets {
		sc.Nodeids = append(sc.Nodeids, k)
	}
	sort.Strings(sc.Nodeids)
	for _, k := range sc.Nodeids {
		sc.Tids = append(sc.Tids, NTickets[k])
	}
	return sc
}

