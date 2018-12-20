package ticketcache

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/crypto"
	"Platon-go/ethdb"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/big"
	"testing"
)

func Test_All(t *testing.T)  {
	ldb, err := ethdb.NewLDBDatabase("./data/platon/chaindata", 0, 0)
	if err!=nil {
		t.Errorf("NewLDBDatabase faile")
	}
	defer ldb.Close()

	number := big.NewInt(0)
	bkhash := common.HexToHash("0xe58643cfb17fc50b8579784e134d36bb4021ed70f9baac84190d9260dc005a10")
	tc := NewTicketIdsCache(ldb)

	mapCache := tc.GetNodeTicketsMap(number, bkhash)
	fmt.Println("change befor: ", mapCache)
	mapCache["nodeid1"] = []common.Hash{common.HexToHash("0x8959adf8343c5d256feb9524b027815a18e32a379de8ecc84b515b07be405c12"), common.HexToHash("0xc3b9798e430fe0b647cc91d57141799c2455fff66c485dd1d7f9595363bdb0d0")}
	fmt.Println("change after: ", mapCache)

	numberchild := big.NewInt(1)
	bkhashchild := common.HexToHash("0xf148e650a37cc218268d8cb91972814b63f1f02fc87f28c24e0ed15d3ad4aca5")
	tc.Submit2Cache(numberchild, bkhashchild, mapCache)

	chash, err:= tc.Hash(number, bkhash)
	fmt.Println("hash: ", chash.Hex())
	tc.Commit(ldb)

}

func Test_Hash(t *testing.T)  {

	m := make(map[int]string)
	m[1] = "0x012250607c8f5124eb10b1b49a52bc057490ca6f8e49a0279f295d0d0c50764a"
	m[2] = "0x01ca495a1eac21ce385f85314f8ef2cc38288a4c0244c620b20b43dc14b1ba62"
	m[3] = "0x020eebe9551a6a4b734452ea102bf16904310402ee01b8f597aeef89a9a59383"
	m[4] = "0x02ee4f5e24a6202d99e193f0515807c30b0a0ad1a4497180f684838cbe3bb221"
	m[5] = "0x039804c2b19b8c72a2d0972e2c7a712133ffac6462c6f6dc55160618098b8f69"

	//11111111111111111111111111111111111111111111111111
	nodeTicketIds1 := &NodeTicketIds{}
	nodeTicketIds1.NTickets = make(map[string]*TicketIds)

	for i:=1; i<6; i++  {
		TicketIds := &TicketIds{}
		TicketIds.TicketId = make([][]byte, 0)
		TicketIds.TicketId = append(TicketIds.TicketId, [][]byte{[]byte("1"), []byte("2"), []byte("3")}...)
		nodeTicketIds1.NTickets[m[i]] = TicketIds
	}
	fmt.Println(nodeTicketIds1)
	for k, v := range nodeTicketIds1.NTickets {
		fmt.Printf("1111111111111111k: %s, v: %v\n", k, v)
	}

	out1, err1 := proto.Marshal(getSortStruct(nodeTicketIds1.NTickets))
	if err1 != nil {
		t.Errorf("proto marshal faile")
	}
	fmt.Println("node1 hex: ", hexutil.Encode(out1), "\n hash: ", crypto.Keccak256Hash(out1).Hex())

	//22222222222222222222222222222222222222222222222222
	nodeTicketIds2 := &NodeTicketIds{}
	nodeTicketIds2.NTickets = make(map[string]*TicketIds)
	for i:=5; i>0; i--  {
		TicketIds := &TicketIds{}
		TicketIds.TicketId = make([][]byte, 0)
		TicketIds.TicketId = append(TicketIds.TicketId, [][]byte{[]byte("1"), []byte("2"), []byte("3")}...)
		nodeTicketIds2.NTickets[m[i]] = TicketIds
	}
	fmt.Println(nodeTicketIds2)
	for k, v := range nodeTicketIds2.NTickets {
		fmt.Printf("2222222222222222k: %s, v: %v\n", k, v)
	}

	out2, err2 := proto.Marshal(getSortStruct(nodeTicketIds2.NTickets))
	if err2 != nil {
		t.Errorf("proto marshal faile")
	}
	fmt.Println("node2 hex: ", hexutil.Encode(out2), "\n hash: ", crypto.Keccak256Hash(out2).Hex())

}
