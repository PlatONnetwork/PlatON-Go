package vm

import (
	"Platon-go/common/byteutil"
	"Platon-go/common/hexutil"
	"Platon-go/rlp"
	"bytes"
	"fmt"
	"math/big"
	"testing"
)

func TestTicketPool(t *testing.T)  {
	nodeId := []byte("0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429")
	// owner := []byte("0xf216d6e4c17097a60ee2b8e5c88941cd9f07263b")

	// VoteTicket(count uint64, price *big.Int, nodeId discover.NodeID)
	var VoteTicket [][]byte
	VoteTicket = make([][]byte, 0)
	VoteTicket = append(VoteTicket, uint64ToBytes(0xf1))
	VoteTicket = append(VoteTicket, []byte("VoteTicket"))
	VoteTicket = append(VoteTicket, byteutil.Uint64ToBytes(1000))
	price, ok :=new(big.Int).SetString("14d1120d7b160000", 16)
	if !ok {
		t.Errorf("big int setstring fail")
	}
	VoteTicket = append(VoteTicket, price.Bytes())
	VoteTicket = append(VoteTicket, nodeId)
	bufVote := new(bytes.Buffer)
	err := rlp.Encode(bufVote, VoteTicket)
	if err != nil {
		fmt.Println(err)
		t.Errorf("VoteTicket encode rlp data fail")
	} else {
		fmt.Println("VoteTicket data rlp: ", hexutil.Encode(bufVote.Bytes()))
	}
}