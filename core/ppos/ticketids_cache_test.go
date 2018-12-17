package pposm

import (
	_ "Platon-go/common/hexutil"
	_ "github.com/syndtr/goleveldb/leveldb"
	"testing"
)

func Test(t *testing.T) {

	//filename := ".//testDb"
	//db, err := leveldb.OpenFile(filename, nil)
	//fmt.Println("open err: ", err.Error())
	//if err==nil {
	//	cache, err :=db.Get(ticketPoolCacheKey, nil)
	//	fmt.Println("get cache: ", hexutil.Encode(cache), " err: ", err.Error())
	//	//err = db.Put(ticketPoolCacheKey, []byte("aabbccdd"), nil);
	//	//fmt.Println("set cache err: ", err.Error())
	//}


	//instance, _ := newTicketIdsCache()
	//blocknumber := big.NewInt(600)
	//blockhash := common.HexToHash("0x5678901234567890123456789012345678901234567890123456789012345678")
	//nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	//tIds := make([]common.Hash, 0)
	//for i:=0; i<10 ; i++ {
	//	tIds = append(tIds, common.HexToHash("0xff789012345678901234567890123456789012345678901234567890123456ff"))
	//}
	//instance.Put(blocknumber, blockhash, nodeId, tIds)
	//instance.Commit()
}
