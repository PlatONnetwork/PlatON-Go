package pposm

import (
	"Platon-go/p2p/discover"
	"Platon-go/common"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"math/big"
)

var (
	ErrFaile = errors.New("fail")
	ErrNotfindFromblockNumber = errors.New("Not find tickets from block number")
	ErrNotfindFromblockHash = errors.New("Not find tickets from block hash")
	ErrNotfindFromnodeId = errors.New("Not find tickets from node id")
)

var ticketidsCache *NumBlocks

func NewTicketIdsCache()  *NumBlocks  {

	//read from leveldb
	ticketidsCache = &NumBlocks{}

	//read from leveldb
	fname := ".//pb1.bin"
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file:", err)
		panic("Error reading file:" + err.Error())
	}
	if len(in)!=0 {
		if err := proto.Unmarshal(in, ticketidsCache); err != nil {
			log.Fatalln("Failed to parse address book:", err)
			panic("Failed to parse address book:" + err.Error())
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

	blockNodes, ok := nb.NBlocks[blocknumber.String()]	//to -> func check(){...}
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

	//对比下list储存的时候删除快，还是map储存的时候删除快
	//for i:=

	for _, tin := range tIds {
		for i, tcache := range ticketIds.TicketId {
			if cmp:=bytes.Equal(tin.Bytes(), tcache); cmp{
				ticketIds.TicketId = append(ticketIds.TicketId[:i], ticketIds.TicketId[i+1:]...)
				break
			}
		}
	}
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
	ret := common.Hash{}
	blockNodes, ok := nb.NBlocks[blocknumber.String()]
	if !ok {
		return ret, ErrNotfindFromblockNumber
	}
	nodeTicketIds, ok := blockNodes.BNodes[blockhash.String()]
	if !ok {
		return ret, ErrNotfindFromblockHash
	}
	fmt.Println(nodeTicketIds)
	// ret = sha3(nodeTicketIds.NTickets)

	return ret, nil
}

func (nb *NumBlocks) TCount(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID)(*big.Int, error) {



	return big.NewInt(0), nil
}

func (nb *NumBlocks) Commit() error {

	out, err := proto.Marshal(nb)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}

	//write level db
	fname := ".//pb1.bin"
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
	fmt.Println("Person Marshal: ", hex.EncodeToString(out))

	return ErrFaile
}

//票池调用
/*type ticketidsCache interface {
	Init() error
	Put(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID, ticketIds []common.Hash) error
	Del(ticketIds []common.Hash) error
	Get(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID)([]common.Hash, error)
	Hash(blockhash common.Hash) common.Hash
}

//出块时调用
type tickeidsCommit interface {
	Commit() error	//cache提交到db 供Cbft区块上链时调用
}*/

