package ppos_storage

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	//"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/golang/protobuf/proto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
	"strconv"
	"testing"
	"time"
	//"math/rand"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tecbot/gorocksdb"
)

func TestData(t *testing.T) {
	ldb, err := ethdb.NewLDBDatabase("E:/platon-data/platon/ppos_storage", 0, 0)
	if err!=nil {
		t.Errorf("NewLDBDatabase faile")
	}
	pposTemp := NewPPosTemp(ldb)
	t.Logf("pposTemp info, pposTemp=%+v", pposTemp)

	pposStorage := NewPPOS_storage()
	t.Logf("pposTemp info, pposStorage=%+v", pposStorage)

	pposStorage.t_storage.Sq = 51200

	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")

	//voteOwner := common.HexToAddress("0x20")
	//deposit := new(big.Int).SetUint64(1)
	blockNumber := new(big.Int).SetUint64(10)

	for i := 0; i < 51200; i++ {
		txHash := common.Hash{}
		txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))

		/*ticket := &types.Ticket{
			voteOwner,
			deposit,
			nodeId,
			blockNumber,
			2,
		}*/

		count := uint32(i) + uint32(time.Now().UnixNano())
		price := big.NewInt(int64(count))

		//pposStorage.SetExpireTicket(blockNumber, txHash)
		pposStorage.AppendTicket(nodeId, txHash, count, price)
	}

	for i := 0; i < 1; i++ {
		blockHash := common.Hash{}
		blockHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
		startTempTime := time.Now().UnixNano()
		pposTemp.SubmitPposCache2Temp(blockNumber, new(big.Int).SetUint64(1), blockHash, pposStorage)
		endTempTime := time.Now().UnixNano()
		t.Log("Testing Cache2Temp efficiency", "startTime", startTempTime, "endTime", endTempTime, "time", endTempTime/1e6-startTempTime/1e6)
		startTime := time.Now().UnixNano()
		pposTemp.Commit2DB(blockNumber, blockHash)
		endTime := time.Now().UnixNano()
		t.Log("Testing Cache2Temp efficiency", "startTime", startTime, "endTime", endTime, "time", endTime/1e6-startTime/1e6)
	}
}


//
func TestData2(t *testing.T) {
	//ldb, err := ethdb.NewLDBDatabase("E:/platon-data/platon/ppos_storage", 0, 0)
	ldb, err := ethdb.NewPPosDatabase("E:/platon-data/platon/ppos_storage")
	if err!=nil {
		t.Errorf("NewLDBDatabase faile")
	}
	pposTemp := NewPPosTemp(ldb)
	t.Logf("pposTemp info, pposTemp=%+v", pposTemp)




	for i := 0; i < 100; i++ {

		pposStorage := NewPPOS_storage()
		t.Logf("pposTemp info, pposStorage=%+v", pposStorage)

		pposStorage.t_storage.Sq = 51200

		nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")

		/*voteOwner := common.HexToAddress("0x20")


		deposit := new(big.Int).SetUint64(uint64(rand.Int63()))*/

		for i := 0; i < 51200; i++ {

			//now := time.Now().UnixNano()


			txHash := common.Hash{}
			txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
			//blockNumber := new(big.Int).SetUint64(uint64(i) + uint64(now))
			/*ticket := &types.Ticket{
				voteOwner,
				deposit,
				nodeId,
				blockNumber,
				2,
			}*/

			count := uint32(i) + uint32(time.Now().UnixNano())
			price := big.NewInt(int64(count))

			//pposStorage.SetExpireTicket(blockNumber, txHash)
			pposStorage.AppendTicket(nodeId, txHash, count, price)
		}

		blockHash := common.Hash{}
		blockHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
		startTempTime := time.Now().UnixNano()
		pposTemp.SubmitPposCache2Temp(new(big.Int).SetUint64(uint64(i)), new(big.Int).SetUint64(1), blockHash, pposStorage)
		endTempTime := time.Now().UnixNano()
		t.Log("Testing Cache2Temp efficiency", "startTime", startTempTime, "endTime", endTempTime, "time", endTempTime/1e6-startTempTime/1e6)
		startTime := time.Now().UnixNano()
		pposTemp.Commit2DB(new(big.Int).SetUint64(uint64(i)), blockHash)
		endTime := time.Now().UnixNano()
		t.Log("Testing Cache2Temp efficiency", "startTime", startTime, "endTime", endTime, "time", endTime/1e6-startTime/1e6)
	}
}


func newSqliteDB(file, tableName string, t *testing.T) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		t.Error("open database err: ", err)
		return nil, err
	}
	//defer db.Close()

	sqlStmt := fmt.Sprintf(`create table if not exists %s (id Integer primary key, key varchar(10), value blob);`, tableName)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		t.Logf("create table %q: %s\n", err, sqlStmt)
		return nil, err
	}
	return db, nil
}


func TestData3(t *testing.T) {


	db, err := newSqliteDB("E:/platon-data/platon/ppos_storage.db", "storage", t)

	if nil != err {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}

	start := common.NewTimer()
	start.Begin()


	//pposStorage := NewPPOS_storage()
	//
	//pposStorage.t_storage.Sq = 51200
	//
	//nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
	//
	//voteOwner := common.HexToAddress("0x20")

	for i := 0; i < 20; i++ {

		pposStorage := NewPPOS_storage()

		pposStorage.t_storage.Sq = 51200

		nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")

		//voteOwner := common.HexToAddress("0x20")
		//
		//
		//deposit := new(big.Int).SetUint64(uint64(rand.Int63()))

		for j := 0; j < 51200; j++ {

			txHash := common.Hash{}
			txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
			//blockNumber := new(big.Int).SetUint64(uint64(j))
			/*ticket := &types.Ticket{
				voteOwner,
				deposit,
				nodeId,
				blockNumber,
				2,
			}*/


			count := uint32(i) + uint32(time.Now().UnixNano())
			price := big.NewInt(int64(count))
			//pposStorage.SetExpireTicket(blockNumber, txHash)
			pposStorage.AppendTicket(nodeId, txHash, count, price)
		}


		blockHash := common.Hash{}
		blockHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))

		blockNumber := new(big.Int).SetUint64(uint64(i))

		if pposTemp := buildPBStorage(blockNumber, blockHash, pposStorage); nil == pposTemp {
			t.Log("Call Commit2DB FINISH !!!! , PPOS storage is Empty, do not write disk AND direct short-circuit ...")
		}else{
			// write ppos_storage into disk with protobuf
			if data, err := proto.Marshal(pposTemp); nil != err {
				t.Log("Failed to Commit2DB", "proto err", err, "Time spent", fmt.Sprintf("%v ms", start.End()))
			}else {

				t.Log("Call Commit2DB, write ppos storage data to disk", "blockNumber", blockNumber, "blockHash", blockHash, "data len", len(data), "Time spent", fmt.Sprintf("%v ms", start.End()))
				// replace into student( _id , name ,age ) VALUES ( 1,'zz7zz7zz',25)
				txsql := fmt.Sprintf("replace into %s (id, key, value) values(1, ?, ?)", "storage")
				args := []interface{}{
					PPOS_STORAGE_KEY,
					data,
					//[]byte{},
				}



				res, err := tx.Exec(txsql, args...)
				if err != nil {
					t.Error("tx.Prepare err: ", err)
					tx.Rollback()
					return
				}

				id, _ := res.LastInsertId()
				t.Log("Call Commit2DB, write ppos storage data to disk", "blockNumber", blockNumber, "blockHash", blockHash, "data len", len(data), "lastRow", id, "Time spent", fmt.Sprintf("%v ms", start.End()))
			}
		}
	}
	tx.Commit()
}

func TestData4(t *testing.T) {

	opts := gorocksdb.NewDefaultOptions()
	opts.SetCompactionStyle(gorocksdb.FIFOCompactionStyle)
	opts.SetWriteBufferSize(1024)
	db, err := gorocksdb.OpenDb(opts, "E:/platon-data/platon/ppos_storage")
	if nil != err {
		t.Error("Failed to Open rocksDB", "err", err)
		return
	}

	var (
		writeOptions = gorocksdb.NewDefaultWriteOptions()
	)



	start := common.NewTimer()
	start.Begin()


	for i := 0; i < 100; i++ {

		pposStorage := NewPPOS_storage()

		pposStorage.t_storage.Sq = 51200

		nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")

		//voteOwner := common.HexToAddress("0x20")
		//
		//
		//deposit := new(big.Int).SetUint64(uint64(rand.Int63()))

		for j := 0; j < 51200; j++ {

			txHash := common.Hash{}
			txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
			//blockNumber := new(big.Int).SetUint64(uint64(j))
			/*ticket := &types.Ticket{
				voteOwner,
				deposit,
				nodeId,
				blockNumber,
				2,
			}*/



			count := uint32(i) + uint32(time.Now().UnixNano())
			price := big.NewInt(int64(count))
			//pposStorage.SetExpireTicket(blockNumber, txHash)
			pposStorage.AppendTicket(nodeId, txHash, count, price)
		}


		blockHash := common.Hash{}
		blockHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))

		blockNumber := new(big.Int).SetUint64(uint64(i))

		if pposTemp := buildPBStorage(blockNumber, blockHash, pposStorage); nil == pposTemp {
			t.Log("Call Commit2DB FINISH !!!! , PPOS storage is Empty, do not write disk AND direct short-circuit ...")
		}else{
			// write ppos_storage into disk with protobuf
			// write ppos_storage into disk with protobuf
			if data, err := proto.Marshal(pposTemp); nil != err {
				t.Log("Failed to Commit2DB", "proto err", err, "Time spent", fmt.Sprintf("%v ms", start.End()))
			}else {
				if err := db.Put(writeOptions, []byte(PPOS_STORAGE_KEY), data); nil != err {
					t.Log("Failed to Commit2DB", "write disk err", err, "Time spent", fmt.Sprintf("%v ms", start.End()))
					return
				}
				t.Log("Call Commit2DB, write ppos storage data to disk", "blockNumber", blockNumber, "blockHash", blockHash, "data len", len(data), "Time spent", fmt.Sprintf("%v ms", start.End()))
			}
		}

	}

}

//func TestData(t *testing.T) {
//	ldb, err := ethdb.NewLDBDatabase("E:/platon-data/platon/ppos_storage", 0, 0)
//	if err!=nil {
//		t.Errorf("NewLDBDatabase faile")
//	}
//	pposTemp := NewPPosTemp(ldb)
//	t.Logf("pposTemp info, pposTemp=%+v", pposTemp)
//
//	pposStorage := NewPPOS_storage()
//	t.Logf("pposTemp info, pposStorage=%+v", pposStorage)
//
//	pposStorage.t_storage.Sq = 51200
//
//	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
//
//	voteOwner := common.HexToAddress("0x20")
//	deposit := new(big.Int).SetUint64(1)
//	blockNumber := new(big.Int).SetUint64(10)
//
//	for i := 0; i < 51200; i++ {
//		txHash := common.Hash{}
//		txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
//
//		ticket := &types.Ticket{
//			voteOwner,
//			deposit,
//			nodeId,
//			blockNumber,
//			2,
//		}
//
//		pposStorage.SetExpireTicket(blockNumber, txHash)
//		pposStorage.AppendTicket(nodeId, txHash, ticket)
//	}
//
//	for i := 0; i < 1; i++ {
//		blockHash := common.Hash{}
//		blockHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
//		startTempTime := time.Now().UnixNano()
//		pposTemp.SubmitPposCache2Temp(blockNumber, new(big.Int).SetUint64(1), blockHash, pposStorage)
//		endTempTime := time.Now().UnixNano()
//		t.Log("Test Cache2Temp efficiency", "startTime", startTempTime, "endTime", endTempTime, "time", endTempTime/1e6-startTempTime/1e6)
//		startTime := time.Now().UnixNano()
//		pposTemp.Commit2DB(ldb, blockNumber, blockHash)
//		endTime := time.Now().UnixNano()
//		t.Log("Test Commit2DB efficiency", "startTime", startTime, "endTime", endTime, "time", endTime/1e6-startTime/1e6)
//	}
//}



//func TestData2(t *testing.T) {
//	//ldb, err := ethdb.NewLDBDatabase("E:/platon-data/platon/ppos_storage", 0, 0)
//	ldb, err := ethdb.NewPPosDatabase("E:/platon-data/platon/ppos_storage")
//	if err!=nil {
//		t.Errorf("NewLDBDatabase faile")
//	}
//	pposTemp := NewPPosTemp(ldb)
//	t.Logf("pposTemp info, pposTemp=%+v", pposTemp)
//
//	pposStorage := NewPPOS_storage()
//	t.Logf("pposTemp info, pposStorage=%+v", pposStorage)
//
//	pposStorage.t_storage.Sq = 51200
//
//	nodeId := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")
//
//	voteOwner := common.HexToAddress("0x20")
//
//
//
//	for i := 0; i < 10; i++ {
//
//
//		deposit := new(big.Int).SetUint64(uint64(rand.Int63()))
//
//
//		for i := 0; i < 51200; i++ {
//			txHash := common.Hash{}
//			txHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
//			blockNumber := new(big.Int).SetUint64(uint64(i))
//			ticket := &types.Ticket{
//				voteOwner,
//				deposit,
//				nodeId,
//				blockNumber,
//				2,
//			}
//
//			pposStorage.SetExpireTicket(blockNumber, txHash)
//			pposStorage.AppendTicket(nodeId, txHash, ticket)
//		}
//
//
//
//
//		blockHash := common.Hash{}
//		blockHash.SetBytes(crypto.Keccak256([]byte(strconv.Itoa(time.Now().Nanosecond() + i))))
//		startTempTime := time.Now().UnixNano()
//		pposTemp.SubmitPposCache2Temp(new(big.Int).SetUint64(uint64(i)), new(big.Int).SetUint64(1), blockHash, pposStorage)
//		endTempTime := time.Now().UnixNano()
//		t.Log("Test Cache2Temp efficiency", "startTime", startTempTime, "endTime", endTempTime, "time", endTempTime/1e6-startTempTime/1e6)
//		startTime := time.Now().UnixNano()
//		pposTemp.Commit2DB(ldb, new(big.Int).SetUint64(uint64(i)), blockHash)
//		endTime := time.Now().UnixNano()
//		t.Log("Test Commit2DB efficiency", "startTime", startTime, "endTime", endTime, "time", endTime/1e6-startTime/1e6)
//	}
//}
