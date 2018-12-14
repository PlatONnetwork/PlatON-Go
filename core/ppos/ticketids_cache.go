package pposm

import (
	"Platon-go/p2p/discover"
	"Platon-go/common"
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
)

func newTicketIdsCache() (*NumBlocks, error)  {

	//read from leveldb
	instance := &NumBlocks{}

	//test
	fname := ".//pb1.bin"
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	if len(in)!=0 {
		if err := proto.Unmarshal(in, instance); err != nil {
			log.Fatalln("Failed to parse address book:", err)
		}
	}
	return instance, ErrFaile
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

	return ErrFaile
}

func (nb *NumBlocks) Del(ticketIds []common.Hash) error {

	return ErrFaile
}

func (nb *NumBlocks) Get(blocknumber *big.Int, blockhash common.Hash, nodeId discover.NodeID)([]common.Hash, error) {

	return []common.Hash{}, ErrFaile
}

func (nb *NumBlocks) Hash(blockhash common.Hash) common.Hash {

	return common.Hash{}
}

func (nb *NumBlocks) TCount() *big.Int {
	return big.NewInt(0)
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

