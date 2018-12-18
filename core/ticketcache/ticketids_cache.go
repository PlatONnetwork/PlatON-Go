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
)

var (
	//error def
	ErrNotfindFromblockNumber = errors.New("Not find tickets from block number")
	ErrNotfindFromblockHash = errors.New("Not find tickets from block hash")
	ErrNotfindFromnodeId = errors.New("Not find tickets from node id")
	ErrProbufMarshal = errors.New("protocol buffer Marshal faile")
	ErrLeveldbPut = errors.New("level db put faile")
	ErrExistFromblockHash = errors.New("nodeId->tickets map is exist")

	//const def
	ticketPoolCacheKey = []byte("ticketPoolCache")
)

var ticketidsCache *NumBlocks

func GetNodeTicketsCacheMap(blocknumber *big.Int, blockhash common.Hash) (ret map[string][]common.Hash) {
	if ticketidsCache!=nil {
		var err error
		ret, err = ticketidsCache.GetNodeTicketsMap(blocknumber, blockhash)
		if err!=nil {
			log.Error("GetNodeTicketsMap err: ", err.Error())
		}
	}else {
		log.Error("ticketidsCache==nil!")
	}
	return
}

func NewTicketIdsCache(db ethdb.Database)  *NumBlocks {
	/*
		Put 购票交易新增选票
		Del 节点掉榜，选票过期，选票被选中
	*/
	ticketidsCache = &NumBlocks{}
	cache, err := db.Get(ticketPoolCacheKey)
	if err == nil {
		if err := proto.Unmarshal(cache, ticketidsCache); err != nil {
			log.Error("protocol buffer Unmarshal faile hex: ", hexutil.Encode(cache))
		}
	}
	return ticketidsCache
}

func (nb *NumBlocks) Put(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID, tIds []common.Hash) error  {

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
	ticketIds, ok := nodeTicketIds.NTickets[nodeId.String()]
	if !ok {
		ticketIds = &TicketIds{}
		ticketIds.TicketId = make([][]byte, 0)
		nodeTicketIds.NTickets[nodeId.String()] = ticketIds
	}
	for _, v := range tIds{
		ticketIds.TicketId = append(ticketIds.TicketId, v.Bytes())
	}
	return nil
}

func (nb *NumBlocks) Del(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID, tIds []common.Hash) error {

	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		return ErrNotfindFromblockNumber
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if !ok {
		return ErrNotfindFromblockHash
	}
	ticketIds, ok := nodeTicketIds.NTickets[nodeId.String()]
	if !ok {
		return ErrNotfindFromnodeId
	}
	//1
	mapTIds := make(map[string]common.Hash)
	for _, tin := range tIds {
		mapTIds[tin.String()] = tin
	}
	for i:=0; i<len(ticketIds.TicketId); i++ {
		if _, ok := mapTIds[hexutil.Encode(ticketIds.TicketId[i])]; ok {
			ticketIds.TicketId = append(ticketIds.TicketId[:i], ticketIds.TicketId[i+1:]...)
			i = i-1
		}
	}
	//2
	/*for _, tin := range tIds {
		for i, tcache := range ticketIds.TicketId {
			if cmp:=bytes.Equal(tin.Bytes(), tcache); cmp{
				ticketIds.TicketId = append(ticketIds.TicketId[:i], ticketIds.TicketId[i+1:]...)
				break
			}
		}
	}*/
	return nil
}

func (nb *NumBlocks) Get(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID)([]common.Hash, error) {

	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		return nil, ErrNotfindFromblockNumber
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if !ok {
		return nil, ErrNotfindFromblockHash
	}
	ticketIds, ok := nodeTicketIds.NTickets[nodeId.String()]
	if !ok {
		return nil, ErrNotfindFromnodeId
	}
	ret := make([]common.Hash, 0)
	for _, v := range ticketIds.TicketId{
		ret = append(ret, common.BytesToHash(v))
	}
	return ret, nil
}

func (nb *NumBlocks) Hash(blocknumber *big.Int, blockhash common.Hash) (common.Hash, error) {

	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		return common.Hash{}, ErrNotfindFromblockNumber
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if !ok {
		return common.Hash{}, ErrNotfindFromblockHash
	}
	out, err := proto.Marshal(nodeTicketIds)
	if err != nil {
		return common.Hash{}, ErrProbufMarshal
	}

	return crypto.Keccak256Hash(out), nil
}

func (nb *NumBlocks) TCount(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID)(uint64, error) {

	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		return 0, ErrNotfindFromblockNumber
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if !ok {
		return 0, ErrNotfindFromblockHash
	}
	ticketIds, ok := nodeTicketIds.NTickets[nodeId.String()]
	if !ok {
		return 0, ErrNotfindFromnodeId
	}
	return uint64(len(ticketIds.TicketId)), nil
}

func (nb *NumBlocks) GetNodeTicketsMap(blocknumber *big.Int, blockhash common.Hash) (map[string][]common.Hash, error){

	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		return nil, ErrNotfindFromblockNumber
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if !ok {
		return nil, ErrNotfindFromblockHash
	}
	out := make(map[string][]common.Hash)
	for k, v := range nodeTicketIds.NTickets{
		tids := make([]common.Hash, len(v.TicketId))
		for _, t := range v.TicketId {
			tid := common.Hash{}
			tid.SetBytes(t)
			tids = append(tids, tid)
		}
		out[k] = tids
	}
	return out, nil
}

func (nb *NumBlocks) Submit2Cache(blocknumber *big.Int, blockhash common.Hash, in map[string][]common.Hash) error  {

	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		blockNodes = &BlockNodes{}
		blockNodes.BNodes = make(map[string]*NodeTicketIds)
		nb.NBlocks[blocknumber.String()] = blockNodes
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if ok {
		return ErrExistFromblockHash
	}
	nodeTicketIds = &NodeTicketIds{}
	nodeTicketIds.NTickets = make(map[string]*TicketIds)
	blockNodes.BNodes[blockhash.String()] = nodeTicketIds
	for k, v := range in {
		tIds := &TicketIds{}
		for _, va := range v {
			tIds.TicketId = append(tIds.TicketId, va.Bytes())
		}
		nodeTicketIds.NTickets[k] = tIds
	}
	return nil
}

func (nb *NumBlocks) Commit(db ethdb.Database) error {

	out, err := proto.Marshal(nb)
	if err != nil {
		log.Error("Protocol buffer failed to marshal :", nb, " err: ", err.Error())
		return ErrProbufMarshal
	}
	log.Info("Marshal out: ", hexutil.Encode(out))
	if err := db.Put(ticketPoolCacheKey, out); err != nil  {
		log.Error("level db put faile: ", err.Error())
		return ErrLeveldbPut
	}
	return nil
}



