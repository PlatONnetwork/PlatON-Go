package pposm

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/crypto"
	"Platon-go/p2p/discover"
	"Platon-go/log"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	_ "github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"math/big"
)

var (
	ErrFaile = errors.New("fail")
	ErrNotfindFromblockNumber = errors.New("Not find tickets from block number")
	ErrNotfindFromblockHash = errors.New("Not find tickets from block hash")
	ErrNotfindFromnodeId = errors.New("Not find tickets from node id")
	ErrProbufMarshal = errors.New("protocol buffer Marshal faile")
)

var ticketidsCache *NumBlocks

func NewTicketIdsCache()  *NumBlocks  {

	//read from leveldb
	ticketidsCache = &NumBlocks{}

	//read from leveldb
	fname := ".//pb1.bin"
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		//log.Fatalln("Error reading file:", err)
		panic("Error reading file:" + err.Error())
	}
	if len(in)!=0 {
		if err := proto.Unmarshal(in, ticketidsCache); err != nil {
			//log.Fatalln("Failed to parse address book:", err)
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

func (nb *NumBlocks) TCount(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID)(*big.Int, error) {



	return big.NewInt(0), nil
}

func (nb *NumBlocks) Commit() error {
	out, err := proto.Marshal(nb)
	if err != nil {
		log.Error("Failed to Marshal :", nb, " err: ", err)
		return ErrProbufMarshal
	}
	log.Info("Marshal out: ", hexutil.Encode(out))

	//write level db
	//filedb := ".//testData"
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		//log.Fatalln("Failed to write address book:", err)
	}
	fmt.Println("Person Marshal: ", hex.EncodeToString(out))


	//leveldb.OpenFile()

	return ErrFaile
}



