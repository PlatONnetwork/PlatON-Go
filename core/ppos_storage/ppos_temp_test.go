package ppos_storage

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func TestData(t *testing.T) {
	ldb, err := ethdb.NewLDBDatabase("./data/platon/chaindata", 0, 0)
	if err!=nil {
		t.Errorf("NewLDBDatabase faile")
	}
	pposTemp := NewPPosTemp(ldb)
	t.Logf("pposTemp info, pposTemp=%+v", pposTemp)

	pposStorage := NewPPOS_storage()
	t.Logf("pposTemp info, pposStorage=%+v", pposStorage)

	pposStorage.t_storage.Sq = 51200

	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")

	voteOwner := common.HexToAddress("0x20")
	deposit := new(big.Int).SetUint64(1)
	blockNumber := new(big.Int).SetUint64(10)

	for i := 0; i < 51200; i++ {
		txHash := common.Hash{}
		txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))

		ticket := &types.Ticket{
			voteOwner,
			deposit,
			nodeId,
			blockNumber,
			2,
		}

		pposStorage.SetExpireTicket(blockNumber, txHash)
		pposStorage.AppendTicket(nodeId, txHash, ticket)
	}

	for i := 0; i < 1; i++ {
		blockHash := common.Hash{}
		blockHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
		startTempTime := time.Now().UnixNano()
		pposTemp.SubmitPposCache2Temp(blockNumber, new(big.Int).SetUint64(1), blockHash, pposStorage)
		endTempTime := time.Now().UnixNano()
		t.Log("测试Cache2Temp效率", "startTime", startTempTime, "endTime", endTempTime, "time", endTempTime/1e6-startTempTime/1e6)
		startTime := time.Now().UnixNano()
		pposTemp.Commit2DB(ldb, blockNumber, blockHash)
		endTime := time.Now().UnixNano()
		t.Log("测试Commit2DB效率", "startTime", startTime, "endTime", endTime, "time", endTime/1e6-startTime/1e6)
	}
}
