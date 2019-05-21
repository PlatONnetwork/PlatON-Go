package ppos_storage

import (
	"sync"
	"math/big"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/golang/protobuf/proto"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"fmt"
	"encoding/json"
	"crypto/md5"
)

const ppos_empty_indb  = "leveldb: not found"

var (
	WRITE_PPOS_ERR = errors.New("Failed to Write ppos storage into disk")

	// The key of ppos storage in disk （leveldb）
	PPOS_STORAGE_KEY = []byte("PPOS_STORAGE_KEY")
)

type numTempMap map[string]hashTempMap
type hashTempMap map[common.Hash]*Ppos_storage

// Global PPOS Dependency TEMP
type PPOS_TEMP struct {

	BlockNumber 	*big.Int
	BlockHash 		common.Hash

	db ethdb.Database

	// Record block total count
	BlockCount 	uint32

	// global data temp
	TempMap numTempMap

	lock  *sync.Mutex
}

/**
This is a Global ppos data temp
 */
var  ppos_temp *PPOS_TEMP


func NewPPosTemp(db ethdb.Database) *PPOS_TEMP {

	timer := common.NewTimer()
	timer.Begin()

	log.Info("NewPPosTemp start ...")
	if nil != ppos_temp {
		return ppos_temp
	}
	ppos_temp = new(PPOS_TEMP)

	ppos_temp.db = db

	ppos_temp.BlockCount = 0

	ntemp := make(numTempMap, 0)
	ppos_temp.TempMap = ntemp
	ppos_temp.lock = &sync.Mutex{}

	// defualt value
	ppos_temp.BlockNumber = big.NewInt(0)
	ppos_temp.BlockHash = common.Hash{}

	if data, err := db.Get(PPOS_STORAGE_KEY); nil != err {
		if ppos_empty_indb != err.Error() {
			log.Error("Failed to Call NewPPosTemp to get Global ppos temp by levelDB", "err", err)
			return ppos_temp
		}
	} else {
		log.Debug("Call NewPPosTemp to Unmarshal Global ppos temp", "pb data len", len(data))

		pb_pposTemp := new(PB_PPosTemp)
		if err := proto.Unmarshal(data, pb_pposTemp); err != nil {
			log.Error("Failed to Call NewPPosTemp to Unmarshal Global ppos temp", "err", err)
			return ppos_temp
		}else {
			/**
			build global ppos_temp
			 */

			 // TODO
			log.Debug("NewPPosTemp  loading data from disk:", "data len", len(data), "dataMD5", md5.Sum(data))

			//PrintObject("NewPPosTemp  loading data from disk:", pb_pposTemp)


			pposStorage := unmarshalPBStorage(pb_pposTemp)


			hashMap := make(map[common.Hash]*Ppos_storage, 0)
			blockHash := common.HexToHash(pb_pposTemp.BlockHash)
			hashMap[blockHash] = pposStorage
			ppos_temp.TempMap[pb_pposTemp.BlockNumber] = hashMap

			num, _ := new(big.Int).SetString(pb_pposTemp.BlockNumber, 10)

			ppos_temp.BlockNumber = num
			ppos_temp.BlockHash = blockHash

			log.Debug("Call NewPPosTemp loading into memory data", "blockNumber", pb_pposTemp.BlockNumber, "blockHash", pb_pposTemp.BlockHash)

		}
	}

	log.Debug("Call NewPPosTemp finish ...", "time long ms: ", timer.End())
	return ppos_temp
}

func GetPPosTempPtr() *PPOS_TEMP {
	return ppos_temp
}


func BuildPposCache(blockNumber *big.Int, blockHash common.Hash) *Ppos_storage {
	return ppos_temp.getPposCacheFromTemp(blockNumber, blockHash)
}


// Get ppos storage cache by same block
func (temp *PPOS_TEMP) getPposCacheFromTemp(blockNumber *big.Int, blockHash common.Hash) *Ppos_storage {

	ppos_storage := NewPPOS_storage()

	notGenesisBlock := blockNumber.Cmp(big.NewInt(0)) > 0

	if nil == temp && notGenesisBlock {
		log.Warn("Warn Call getPposCacheFromTemp of PPOS_TEMP, the Global PPOS_TEMP instance is nil !!!!!!!!!!!!!!!", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex())
		return ppos_storage
	}

	if !notGenesisBlock || (common.Hash{}) == blockHash {
		return ppos_storage
	}

	var storage *Ppos_storage

	temp.lock.Lock()
	if hashTemp, ok := temp.TempMap[blockNumber.String()]; !ok {
		log.Warn("Warn Call getPposCacheFromTemp of PPOS_TEMP, the PPOS storage cache is empty by blockNumber !!!!! Direct short-circuit", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex())
		temp.lock.Unlock()
		return ppos_storage
	}else {

		if pposStorage, ok := hashTemp[blockHash]; !ok {
			log.Warn("Warn Call getPposCacheFromTemp of PPOS_TEMP, the PPOS storage cache is empty by blockHash !!!!! Direct short-circuit", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex())
			temp.lock.Unlock()
			return ppos_storage
		}else {
			start := common.NewTimer()
			start.Begin()
			storage = pposStorage.Copy()
			log.Debug("Call getPposCacheFromTemp of PPOS_TEMP, Copy ppos_storage FINISH !!!!!!", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "Time spent", fmt.Sprintf("%v ms", start.End()))
		}
	}
	temp.lock.Unlock()
	return storage
}

// Set ppos storage cache by same block
func (temp *PPOS_TEMP) SubmitPposCache2Temp(blockNumber, blockInterval *big.Int, blockHash common.Hash, storage *Ppos_storage)  {
	log.Info("Call SubmitPposCache2Temp of PPOS_TEMP", "blockNumber", blockNumber.String(), "blockHash", blockHash.Hex(),
		"blockInterval", blockInterval, "Before SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount)

	start := common.NewTimer()
	start.Begin()

	empty := verifyStorageEmpty(storage)


	temp.lock.Lock()
	/**
	first condition blockNumber
	There are four kinds.
	1a、origin data is empty by blockNumber AND input data is not empty; 	then we will Add （set） input data into global temp
	2a、origin data is empty by blockNumber AND input data is empty;     	then we direct short-circuit
	3a、origin data is not empty by blockNumber AND input data is not empty; then we will update （set） global temp by indata
	4a、origin data is not empty by blockNumber AND input data is empty; 	then we will delete global temp data
	 */

	originHashTemp, hasNum := temp.TempMap[blockNumber.String()]
	// match 1a
	if !hasNum && !empty {
		originHashTemp = make(hashTempMap, 1)
		originHashTemp[blockHash] = storage
		temp.TempMap[blockNumber.String()] = originHashTemp

		temp.BlockCount += 1

		temp.deleteAnyTemp(blockNumber, blockInterval, blockHash)
		temp.lock.Unlock()
		return
	}else if  !hasNum && empty { // match 2a
		log.Debug("Call SubmitPposCache2Temp of PPOS_TEMP， origin ppos_storage and input ppos_storage is empty by BlockNumber !!!! Direct short-circuit", "blockNumber", blockNumber.Uint64(),
			"blockHash", blockHash.Hex(), "blockInterval", blockInterval, " Before SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount, "Time spent", fmt.Sprintf("%v ms", start.End()))
		temp.lock.Unlock()
		return
	}

	_, hasHash := originHashTemp[blockHash]

	/**
	second condition blockHash
	There are four kinds, too.
	1b、origin data is empty by blockHash AND input data is not empty; 		then we will Add （set） indata into global temp
	2b、origin data is empty by blockHash AND input data is empty;     		then we direct short-circuit
	3b、origin data is not empty by blockHash AND input data is not empty; 	then we will update （set） global temp by indata
	4b、origin data is not empty by blockHash AND input data is empty; 		then we will delete global temp data
	 */

	 // match 1b
	 if !hasHash && !empty {
		 originHashTemp[blockHash] = storage
		 temp.TempMap[blockNumber.String()] = originHashTemp

		 temp.BlockCount += 1

		 temp.deleteAnyTemp(blockNumber, blockInterval, blockHash)
		 temp.lock.Unlock()
		 return
	 }else if  !hasHash && empty { // match 2b
		 log.Debug("Call SubmitPposCache2Temp of PPOS_TEMP， origin ppos_storage and input ppos_storage is empty by BlockHash !!!! Direct short-circuit", "blockNumber", blockNumber.Uint64(),
			 "blockHash", blockHash.Hex(), "blockInterval", blockInterval, " Before SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount, "Time spent", fmt.Sprintf("%v ms", start.End()))
		 temp.lock.Unlock()
		 return
	 }

	 // now remain 3a 4a 3b 4b
	 if hasHash && empty { // delete origin data
		delete(originHashTemp, blockHash)
		temp.BlockCount -= 1
		if len(originHashTemp) == 0 {
			delete(temp.TempMap, blockNumber.String())
		}

		temp.deleteAnyTemp(blockNumber, blockInterval, blockHash)
	 }else if hasHash && !empty {
		 originHashTemp[blockHash] = storage
		 temp.TempMap[blockNumber.String()] = originHashTemp

		 temp.deleteAnyTemp(blockNumber, blockInterval, blockHash)
	 }
	temp.lock.Unlock()
	log.Debug("Call SubmitPposCache2Temp of PPOS_TEMP，SUCCESS !!!!!!", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
		"blockInterval", blockInterval, " Before SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount, "Time spent", fmt.Sprintf("%v ms", start.End()))

}

func (temp *PPOS_TEMP) Commit2DB(blockNumber *big.Int, blockHash common.Hash) error {
	start := common.NewTimer()
	start.Begin()


	var ps *Ppos_storage
	temp.lock.Lock()
	if hashMap, ok := temp.TempMap[blockNumber.String()]; !ok {
		temp.lock.Unlock()
		return nil
	}else {
		if ppos_storage, ok := hashMap[blockHash]; !ok {
			temp.lock.Unlock()
			return nil
		}else {
			ps = ppos_storage
		}
	}
	temp.lock.Unlock()




	if pposTemp := buildPBStorage(blockNumber, blockHash, ps); nil == pposTemp {
		log.Debug("Call Commit2DB FINISH !!!! , PPOS storage is Empty, do not write disk AND direct short-circuit ...")
		return nil
	}else{
		// write ppos_storage into disk with protobuf
		if data, err := proto.Marshal(pposTemp); nil != err {
			log.Error("Failed to Commit2DB", "proto err", err, "Time spent", fmt.Sprintf("%v ms", start.End()))
			return err
		}else {
			if len(data) != 0 {

				// TODO

				if err := temp.db.Put(PPOS_STORAGE_KEY, data); err != nil {
					log.Error("Failed to Call Commit2DB:" + WRITE_PPOS_ERR.Error(), "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "data len", len(data), "Time spent", fmt.Sprintf("%v ms", start.End()), "err", err)
					return WRITE_PPOS_ERR
				}

				temp.BlockNumber = blockNumber
				temp.BlockHash = blockHash

			}
			log.Info("Call Commit2DB, write ppos storage data to disk", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "data len", len(data), "dataMD5", md5.Sum(data), "Time spent", fmt.Sprintf("%v ms", start.End()))
		}
	}
	return nil
}

// Gets ppos_storag pb from db
func  (temp *PPOS_TEMP) GetPPosStorageProto() (common.Hash, []byte, error) {
	start := common.NewTimer()
	start.Begin()
	if data, err := temp.db.Get(PPOS_STORAGE_KEY); nil != err {
		if ppos_empty_indb == err.Error() {
			log.Debug("Call GetPPosStorageProto, ppos storage is empty in disk ...")
			return common.Hash{}, nil, nil
		}else {
			log.Warn("Failed to Call GetPPosStorageProto to get Global ppos temp by levelDB", "err", err)
			return common.Hash{}, nil, err
		}
	} else {
		log.Debug("Call GetPPosStorageProto to Unmarshal Global ppos temp", "pb data len", len(data))

		pb_pposTemp := new(PB_PPosTemp)
		if err := proto.Unmarshal(data, pb_pposTemp); err != nil {
			log.Error("Failed to Call GetPPosStorageProto to Unmarshal Global ppos temp", "err", err)
			return common.Hash{}, nil, err
		}else {
			// TODO

			//PrintObject("GetPPosStorageProto resolve the data of PB:", pb_pposTemp)
			curr_Num, _ := new(big.Int).SetString(pb_pposTemp.BlockNumber, 10)
			if curr_Num.Cmp(big.NewInt(common.BaseElection - 1)) < 0 {
				return common.Hash{}, nil, nil
			}

			log.Debug("Call GetPPosStorageProto FINISH !!!!", "blockNumber", pb_pposTemp.BlockNumber, "blockHash", pb_pposTemp.BlockHash, "data len", len(data), "dataMD5", md5.Sum(data), "Time spent", fmt.Sprintf("%v ms", start.End()))
			return common.HexToHash(pb_pposTemp.BlockHash), data, nil
		}
	}
}

// Flush data into db
func (temp *PPOS_TEMP) PushPPosStorageProto(data []byte)  error {
	if len(data) == 0 {
		return nil
	}
	start := common.NewTimer()
	start.Begin()
	pb_pposTemp := new(PB_PPosTemp)
	if err := proto.Unmarshal(data, pb_pposTemp); err != nil {
		log.Error("Failed to Call PushPPosStorageProto to Unmarshal Global ppos temp", "err", err)
		return err
	}else {
		/**
		build global ppos_temp
		 */
		 // TODO

		log.Debug("PushPPosStorageProto input params:", "data len", len(data), "dataMD5", md5.Sum(data))

		//PrintObject("PushPPosStorageProto resolve the data of PB, will flush disk:", pb_pposTemp)

		var hashMap map[common.Hash]*Ppos_storage

		ppos_temp.lock.Lock()

		if hashData, ok := ppos_temp.TempMap[pb_pposTemp.BlockNumber]; ok {
			hashMap = hashData
		}else {
			hashMap = make(map[common.Hash]*Ppos_storage, 1)
		}

		pposStorage := unmarshalPBStorage(pb_pposTemp)



		blockHash := common.HexToHash(pb_pposTemp.BlockHash)
		hashMap[blockHash] = pposStorage
		ppos_temp.TempMap[pb_pposTemp.BlockNumber] = hashMap



		ppos_temp.lock.Unlock()


		// flush data into disk
		if len(data) != 0 {
			log.Debug("Call PushPPosStorageProto flush data into disk start ...", "blockNumber", pb_pposTemp.BlockNumber, "blockHash", pb_pposTemp.BlockHash, "data len", len(data))
			if err := temp.db.Put(PPOS_STORAGE_KEY, data); err != nil {
				log.Error("Failed to Call PushPPosStorageProto:" + WRITE_PPOS_ERR.Error(), "blockNumber", pb_pposTemp.BlockNumber, "blockHash", pb_pposTemp.BlockHash, "data len", len(data), "Time spent", fmt.Sprintf("%v ms", start.End()), "err", err)
				return WRITE_PPOS_ERR
			}
		}

		num, _ := new(big.Int).SetString(pb_pposTemp.BlockNumber, 10)

		temp.BlockNumber = num
		temp.BlockHash = blockHash
	}



	log.Debug("Call PushPPosStorageProto FINISH !!!!", "blockNumber", pb_pposTemp.BlockNumber, "blockHash", pb_pposTemp.BlockHash, "data len", len(data), "Time spent", fmt.Sprintf("%v ms", start.End()))
	return nil
}


func buildPBStorage(blockNumber *big.Int, blockHash common.Hash, ps *Ppos_storage) *PB_PPosTemp {
	ppos_temp := new(PB_PPosTemp)
	ppos_temp.BlockNumber = blockNumber.String()
	ppos_temp.BlockHash = blockHash.Hex()

	var empty int = 0  // 0: empty 1: no
	var wg sync.WaitGroup

	/**
	candidate related
	*/
	if nil != ps.c_storage {

		canTemp := new(CandidateTemp)


		wg.Add(6)
		// previous witness
		go func() {
			if queue := buildPBcanqueue("buildPBStorage pres", ps.c_storage.pres); len(queue) != 0 {
				canTemp.Pres = queue
				empty |= 1
			}
			wg.Done()
		}()
		// current witness
		go func() {
			if queue := buildPBcanqueue("buildPBStorage currs", ps.c_storage.currs); len(queue) != 0 {
				canTemp.Currs = queue
				empty |= 1
			}
			wg.Done()
		}()
		// next witness
		go func() {
			if queue := buildPBcanqueue("buildPBStorage nexts", ps.c_storage.nexts); len(queue) != 0 {
				canTemp.Nexts = queue
				empty |= 1
			}
			wg.Done()
		}()
		// immediate
		go func() {
			if queue := buildPBcanqueue("buildPBStorage imms", ps.c_storage.imms); len(queue) != 0 {
				canTemp.Imms = queue
				empty |= 1
			}
			wg.Done()
		}()
		// reserve
		go func() {
			if queue := buildPBcanqueue("buildPBStorage res", ps.c_storage.res); len(queue) != 0 {
				canTemp.Res = queue
				empty |= 1
			}
			wg.Done()
		}()

		// refund
		go func() {
			if refundMap := buildPBrefunds(ps.c_storage.refunds); len(refundMap) != 0 {
				canTemp.Refunds = refundMap
				empty |= 1
			}
			wg.Done()
		}()

		wg.Wait()
		ppos_temp.CanTmp = canTemp
	}

	/**
	ticket related
	*/
	if nil != ps.t_storage {
		tickTemp := new(TicketTemp)

		tickTemp.Sq = ps.t_storage.Sq

		// SQ
		if ps.t_storage.Sq != -1 {
			empty |= 1
		}

		//wg.Add(3)

		// ticketInfos
		/*go func() {
			if ticketMap := buildPBticketMap(ps.t_storage.Infos); len(ticketMap) != 0 {
				tickTemp.Infos = ticketMap
				empty |= 1
			}
			wg.Done()
		}()*/

		//// ExpireTicket
		//go func() {
		//	if ets := buildPBexpireTicket(ps.t_storage.Ets); len(ets) != 0 {
		//		tickTemp.Ets = ets
		//		empty |= 1
		//	}
		//	wg.Done()
		//}()

		// ticket's attachment  of node
		//go func() {
			if dependency := buildPBdependencys(ps.t_storage.Dependencys); len(dependency) != 0 {
				tickTemp.Dependencys = dependency
				empty |= 1
			}
			//wg.Done()
		//}()

		//wg.Wait()
		ppos_temp.TickTmp = tickTemp
	}

	if empty == 0 {
		return nil
	}
	return ppos_temp
}


func unmarshalPBStorage(pb_temp *PB_PPosTemp) *Ppos_storage {

	ppos_storage := new(Ppos_storage)

	/**
	candidate related
	*/
	canGlobalTemp := pb_temp.CanTmp
	if nil !=  canGlobalTemp {

		canTemp := new(candidate_temp)

		buildQueueFunc := func(arr []*CandidateInfo) types.CandidateQueue {
			if len(arr) == 0 {
				return nil
			}
			queue := make(types.CandidateQueue, len(arr))
			for i, can := range arr {
				deposit, _ := new(big.Int).SetString(can.Deposit, 10)
				num, _ := new(big.Int).SetString(can.BlockNumber, 10)
				canInfo := &types.Candidate{
					Deposit: 		deposit,
					BlockNumber:	num,
					TxIndex:		can.TxIndex,
					CandidateId:	discover.MustHexID(can.CandidateId),
					Host:        	can.Host,
					Port:         	can.Port,
					Owner:  		common.HexToAddress(can.Owner),
					Extra:  		can.Extra,
					Fee:  			can.Fee,
					TxHash: 		common.HexToHash(can.TxHash),
					TOwner: 		common.HexToAddress(can.TOwner),
				}
				queue[i] = canInfo
			}
			return queue
		}

		// previous witness
		canTemp.pres = buildQueueFunc(canGlobalTemp.Pres)
		// current witness
		canTemp.currs = buildQueueFunc(canGlobalTemp.Currs)
		// next witness
		canTemp.nexts = buildQueueFunc(canGlobalTemp.Nexts)
		// immediate
		canTemp.imms = buildQueueFunc(canGlobalTemp.Imms)
		// reserve
		canTemp.res = buildQueueFunc(canGlobalTemp.Res)
		// refund
		/*if len(canGlobalTemp.Refunds) != 0 {


		}*/



		defeatMap := make(refundStorage, len(canGlobalTemp.Refunds))

		for nodeId, refundArr := range canGlobalTemp.Refunds {

			if len(refundArr.Defeats) == 0 {
				continue
			}

			defeatArr := make(types.RefundQueue, len(refundArr.Defeats))
			for i, defeat := range refundArr.Defeats {

				deposit, _ := new(big.Int).SetString(defeat.Deposit, 10)
				num, _ := new(big.Int).SetString(defeat.BlockNumber, 10)

				refund := &types.CandidateRefund{
					Deposit:  		deposit,
					BlockNumber: 	num,
					Owner: 			common.HexToAddress(defeat.Owner),
				}
				defeatArr[i] = refund
			}
			defeatMap[discover.MustHexID(nodeId)] = defeatArr
		}

		canTemp.refunds = defeatMap
		ppos_storage.c_storage = canTemp
	}


	/**
	ticket related
	 */
	tickGlobalTemp := pb_temp.TickTmp
	if nil != tickGlobalTemp {

		tickTemp := new(ticket_temp)

		// SQ
		tickTemp.Sq = tickGlobalTemp.Sq

		// ticketInfo map
		/*if len(tickGlobalTemp.Infos) != 0 {

			infoMap := make(map[common.Hash]*types.Ticket, len(tickGlobalTemp.Infos))
			for tid, tinfo := range tickGlobalTemp.Infos {
				deposit, _ := new(big.Int).SetString(tinfo.Deposit, 10)
				num, _ := new(big.Int).SetString(tinfo.BlockNumber, 10)
				ticketInfo := &types.Ticket{
					Owner: 			common.BytesToAddress(tinfo.Owner),
					Deposit:		deposit,
					CandidateId: 	discover.MustBytesID(tinfo.CandidateId),
					BlockNumber: 	num,
					Remaining:		tinfo.Remaining,
				}

				infoMap[common.HexToHash(tid)] = ticketInfo
			}
			tickTemp.Infos = infoMap
		}*/


		//// ExpireTicket map
		//if len(tickGlobalTemp.Ets) != 0 {
		//	ets := make(map[string][]common.Hash, len(tickGlobalTemp.Ets))
		//
		//	for blockNum, ticketIdArr := range tickGlobalTemp.Ets {
		//
		//		if len(ticketIdArr.TxHashs) == 0 {
		//			continue
		//		}
		//
		//		ticketIds := make([]common.Hash, len(ticketIdArr.TxHashs))
		//		for i, ticketId := range ticketIdArr.TxHashs {
		//			ticketIds[i] = common.BytesToHash(ticketId)
		//		}
		//		ets[blockNum] = ticketIds
		//	}
		//
		//	tickTemp.Ets = ets
		//}

		// ticket's attachment  of node
		/*if len(tickGlobalTemp.Dependencys) != 0 {


		}*/



		dependencyMap := make(map[discover.NodeID]*ticketDependency, len(tickGlobalTemp.Dependencys))

		for nodeIdStr, pb_dependency := range tickGlobalTemp.Dependencys {

			dependencyInfo := new(ticketDependency)
			//dependencyInfo.Age = pb_dependency.Age
			dependencyInfo.Num = pb_dependency.Num


			/*tidArr := make([]common.Hash, len(pb_dependency.Tids))

			for i, ticketId := range pb_dependency.Tids {
				tidArr[i] = common.BytesToHash(ticketId)
			}

			dependencyInfo.Tids = tidArr*/



			fieldArr := make([]*ticketInfo, len(pb_dependency.Tinfo))
			for j, field := range pb_dependency.Tinfo {

				price, _ := new(big.Int).SetString(field.Price, 10)
				f := &ticketInfo{
					TxHash: 	common.HexToHash(field.TxHash),
					Remaining: 	field.Remaining,
					Price: 		price,
				}

				fieldArr[j] =  f
			}
			dependencyInfo.Tinfo = fieldArr

			dependencyMap[discover.MustHexID(nodeIdStr)] = dependencyInfo
		}

		tickTemp.Dependencys = dependencyMap
		ppos_storage.t_storage = tickTemp
	}

	return ppos_storage
}

func buildPBcanqueue (title string, canQqueue types.CandidateQueue) []*CandidateInfo {

	PrintObject(title + " ,buildPBcanqueue:", canQqueue)

	pbQueue := make([]*CandidateInfo, len(canQqueue))
	if len(canQqueue) == 0 {
		return pbQueue
	}

	for i, can := range canQqueue {
		canInfo := &CandidateInfo{
			Deposit: 		can.Deposit.String(),
			BlockNumber:	can.BlockNumber.String(),
			TxIndex:		can.TxIndex,
			CandidateId:	can.CandidateId.String(),
			Host:			can.Host,
			Port:			can.Port,
			Owner:			can.Owner.String(),
			Extra:			can.Extra,
			Fee: 			can.Fee,
			TxHash: 		can.TxHash.String(),
			TOwner: 		can.TOwner.String(),
		}
		pbQueue[i] = canInfo
	}
	return pbQueue
}


func buildPBrefunds(refunds refundStorage) map[string]*RefundArr {

	PrintObject("buildPBrefunds", refunds)

	if len(refunds) == 0 {
		return nil
	}

	refundMap := make(map[string]*RefundArr, len(refunds))

	for nodeId, rs := range refunds {

		if len(rs) == 0 {
			continue
		}
		defeats := make([]*Refund, len(rs))
		for i, refund := range rs {
			refundInfo := &Refund{
				Deposit:     	refund.Deposit.String(),
				BlockNumber:	refund.BlockNumber.String(),
				Owner:			refund.Owner.String(),
			}
			defeats[i] = refundInfo
		}

		refundArr := &RefundArr{
			Defeats: defeats,
		}
		refundMap[nodeId.String()] = refundArr
	}
	return refundMap
}


//func buildPBticketMap(tickets map[common.Hash]*types.Ticket) map[string]*TicketInfo {
//	if len(tickets) == 0 {
//		return nil
//	}
//
//	pb_ticketMap := make(map[string]*TicketInfo, len(tickets))
//
//	for tid, tinfo := range tickets {
//		ticketInfo := &TicketInfo{
//			Owner:       tinfo.Owner.Bytes(),
//			Deposit:     tinfo.Deposit.String(),
//			CandidateId: tinfo.CandidateId.Bytes(),
//			BlockNumber: tinfo.BlockNumber.String(),
//			Remaining:   tinfo.Remaining,
//		}
//		pb_ticketMap[tid.String()] = ticketInfo
//	}
//	return pb_ticketMap
//}


//func buildPBexpireTicket(ets map[string][]common.Hash) map[string]*TxHashArr  {
//	if len(ets) == 0 {
//		return nil
//	}
//
//	pb_ets := make(map[string]*TxHashArr, len(ets))
//
//	for blockNumber, ticketIdArr := range ets {
//
//		if len(ticketIdArr) == 0 {
//			continue
//		}
//
//		txHashs := make([][]byte, len(ticketIdArr))
//
//		for i, tid := range ticketIdArr {
//			txHashs[i] = tid.Bytes()
//		}
//
//		txHashArr := new(TxHashArr)
//		txHashArr.TxHashs = txHashs
//		pb_ets[blockNumber] = txHashArr
//	}
//	return pb_ets
//}


func buildPBdependencys(dependencys map[discover.NodeID]*ticketDependency) map[string]*TicketDependency {
	if len(dependencys) == 0 {
		return nil
	}

	pb_dependency := make(map[string]*TicketDependency, len(dependencys))

	for nodeId, dependency := range dependencys {

		/*tidArr := make([][]byte, len(dependency.Tids))

		for i, ticketId := range dependency.Tids {
			tidArr[i] = ticketId.Bytes()
		}*/


		if dependency.Num == 0 && len(dependency.Tinfo) == 0 {
			continue
		}

		fieldArr := make([]*Field, len(dependency.Tinfo))


		for i, field := range dependency.Tinfo {

			f := &Field{
				TxHash:		field.TxHash.String(),
				Remaining: 	field.Remaining,
				Price: 		field.Price.String(),
			}
			fieldArr[i] = f
		}

		depenInfo := &TicketDependency{
			//Age:  dependency.Age,
			Num:  dependency.Num,
			//Tids: tidArr,
			Tinfo: 	fieldArr,
		}

		pb_dependency[nodeId.String()] = depenInfo
	}
	return pb_dependency
}

func (temp *PPOS_TEMP) deleteAnyTemp (blockNumber, blockInterval *big.Int, blockHash common.Hash) {

	// delete font any data
	if big.NewInt(0).Cmp(blockInterval) > 0 {
		log.Error("WARN WARN WARN !!! Call SubmitPposCache2Temp of PPOS_TEMP FINISH !!!!!! blockInterval is NEGATIVE NUMBER", "blockNumber", blockNumber.String(),
			"blockHash", blockHash.Hex(), "blockInterval", blockInterval, "After SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount)
		return
	}

	// blockInterval is the difference of block height between
	// the highest block in memory and the highest block in the chain
	interval := new(big.Int).Add(blockInterval, big.NewInt(30))

	// del old cache
	// blocknumber: current memory block
	target := new(big.Int).Sub(blockNumber, interval)
	for number := range temp.TempMap {
		if currentNum, ok := new(big.Int).SetString(number, 0); ok {
			if currentNum.Cmp(target) < 0 {

				hashMap, ok := temp.TempMap[currentNum.String()]

				// delete current number related ppos data
				delete(temp.TempMap, number)
				// decr block count
				if ok {
					temp.BlockCount -= uint32(len(hashMap))
				}
			}
		}
	}
	log.Debug("Call SubmitPposCache2Temp of PPOS_TEMP FINISH !!!!!!", "blockNumber", blockNumber.String(), "blockHash", blockHash.Hex(),
		"blockInterval", blockInterval, "After SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount)

}

func verifyStorageEmpty(storage *Ppos_storage) bool {
	if nil == storage.c_storage && nil == storage.t_storage {
		return true
	}

	var canEmpty, tickEmpty bool

	canStorage := storage.c_storage
	if nil != canStorage {
		if len(canStorage.pres) == 0 && len(canStorage.currs) == 0 && len(canStorage.nexts) == 0 &&
			len(canStorage.imms) == 0 && len(canStorage.res) == 0 && len(canStorage.refunds) == 0 {
			canEmpty = true
		}
	}

	tickStorage := storage.t_storage
	if nil != tickStorage {
		if tickStorage.Sq == -1 && /*len(tickStorage.Infos) == 0 &&
			len(tickStorage.Ets) == 0 &&*/ len(tickStorage.Dependencys) == 0 {
			tickEmpty = true
		}
	}
	if canEmpty && tickEmpty {
		return true
	}
	return false
}



func PrintObject(s string, obj interface{}) {
	objs, _ := json.Marshal(obj)
	log.Debug(s, "==", string(objs))
}