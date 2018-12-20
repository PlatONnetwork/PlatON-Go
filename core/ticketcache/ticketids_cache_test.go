package ticketcache

import (
	"Platon-go/common"
	"Platon-go/ethdb"
	"fmt"
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
