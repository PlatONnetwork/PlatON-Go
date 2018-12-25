package ticketcache

import (
	"Platon-go/common"
	"Platon-go/common/byteutil"
	"Platon-go/crypto"
	"Platon-go/ethdb"
	"Platon-go/p2p/discover"
	"fmt"
	"math/big"
	"testing"
)

const (
	blockCount = 20
	nodeCount = 200
	ticketCount = 51200
)

func getBlockMaxData() (TicketCache, error) {
	//every nodeid has 256 ticket total has 200 nodeid
	ret := NewTicketCache()
	for n:=0; n<nodeCount; n++ {
		nodeid := make([]byte, 0, 64)
		nodeid = append(nodeid, crypto.Keccak256Hash([]byte("nodeid"), byteutil.IntToBytes(n)).Bytes()...)
		nodeid = append(nodeid, crypto.Keccak256Hash([]byte("nodeid"), byteutil.IntToBytes(n*10)).Bytes()...)
		NodeId, err := discover.BytesID(nodeid)
		if err!=nil {
			return ret, err
		}
		tids := make([]common.Hash, 0)
		for i:=0; i<ticketCount/nodeCount ; i++ {
			tids = append(tids, crypto.Keccak256Hash([]byte("tid"), byteutil.IntToBytes(i)))
		}
		ret[NodeId] = tids
	}
	return ret, nil
}

func Test_Timer (t *testing.T)  {
	timer := Timer{}
	for i:=0; i<1000; i++  {
		timer.Begin()
		//time.Sleep()
		fmt.Println("i: ", i, " t: ", timer.End())
	}
}

func Test_GenerateData (t *testing.T)  {
	for i:=0; i<blockCount; i++  {
		_,err := getBlockMaxData()
		if err!=nil {
			fmt.Println("getMaxtickets faile err: ", err.Error())
			t.Errorf("getMaxtickets faile")
		}
	}
}

func Test_New(t *testing.T)  {
	ldb, err := ethdb.NewLDBDatabase("./data/platon/chaindata", 0, 0)
	if err!=nil {
		t.Errorf("NewLDBDatabase faile")
	}
	timer := Timer{}
	timer.Begin()
	NewTicketIdsCache(ldb)
	fmt.Printf("NewTicketIdsCache time [ms=%.3f]\n", timer.End())
	ldb.Close()
}

func Test_Submit2Cache(t *testing.T)  {
	ldb, err := ethdb.NewLDBDatabase("./data/platon/chaindata", 0, 0)
	if err!=nil {
		t.Errorf("NewLDBDatabase faile")
	}
	tc := NewTicketIdsCache(ldb)
	for i:=0; i<blockCount; i++  {
		number := big.NewInt(int64(i))
		bkhash := crypto.Keccak256Hash(byteutil.IntToBytes(i))
		mapCache ,err := getBlockMaxData()
		if err!=nil {
			fmt.Println("getMaxtickets faile err: ", err.Error())
			t.Errorf("getMaxtickets faile")
		}
		timer := Timer{}
		timer.Begin()
		fmt.Println("msg==> ", "begin: ", timer.End())
		tc.Submit2Cache(number, bkhash, mapCache)
		fmt.Printf("run submit time [index=%d][ms=%.3f]\n", i, timer.End())
	}
	ldb.Close()
}

func Test_Write(t *testing.T)  {
	ldb, err := ethdb.NewLDBDatabase("./data/platon/chaindata", 0, 0)
	if err!=nil {
		t.Errorf("NewLDBDatabase faile")
	}
	timer := Timer{}
	tc := NewTicketIdsCache(ldb)
	for i:=0; i<blockCount; i++  {
		number := big.NewInt(int64(i))
		bkhash := crypto.Keccak256Hash(byteutil.IntToBytes(i))
		mapCache ,err := getBlockMaxData()
		if err!=nil {
			fmt.Println("getMaxtickets faile err: ", err.Error())
			t.Errorf("getMaxtickets faile")
		}

		////Submit2Cache
		timer.Begin()
		tc.Submit2Cache(number, bkhash, mapCache)
		fmt.Printf("run submit time [index=%d][ms=%.3f]\n", i, timer.End())

		//Hash
		timer.Begin()
		chash, err:= tc.Hash(mapCache)
		fmt.Printf("run hash time [index=%d][ms=%.3f][hash=%s]\n", i, timer.End(), chash.Hex())
	}

	//Commit
	timer.Begin()
	tc.Commit(ldb)
	fmt.Printf("run Commit time [ms=%.3f]\n", timer.End())
	ldb.Close()
}

func Test_Read(t *testing.T)  {
	ldb, err := ethdb.NewLDBDatabase("./data/platon/chaindata", 0, 0)
	if err!=nil {
		t.Errorf("NewLDBDatabase faile")
	}
	timer := Timer{}
	tcCopy := NewTicketIdsCache(ldb)
	for i:=0; i<blockCount; i++  {
		number := big.NewInt(int64(i))
		bkhash := crypto.Keccak256Hash(byteutil.IntToBytes(i))

		//==>run GetNodeTicketsMap time
		timer.Begin()
		tcCopy.GetNodeTicketsMap(number, bkhash)
		fmt.Printf("run getNodeTicketsMap time [index=%d] [ms=%.3f]\n", i, timer.End())
	}
	ldb.Close()
}

