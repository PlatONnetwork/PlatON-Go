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
)

var (
	WRITE_PPOS_ERR = errors.New("Failed to Write ppos storage into disk")

	// The key of ppos storage in disk （leveldb）
	PPOS_STORAGE_KEY = []byte("PPOS_STORAGE_KEY")
)

type numTempMap map[string]hashTempMap
type hashTempMap map[common.Hash]*Ppos_storage

// Global PPOS Dependency TEMP
type PPOS_TEMP struct {
	// Record block total count
	BlockCount 	uint32

	// global data temp
	TempMap numTempMap

	lock  *sync.Mutex
}


// Get ppos storage cache by same block
func (temp *PPOS_TEMP) GetPposCacheFromTemp(blockNumber *big.Int, blockHash common.Hash) *Ppos_storage {

	ppos_storage := GetPPOS_storage()

	notGenesisBlock := blockNumber.Cmp(big.NewInt(0)) > 0

	if nil == temp && notGenesisBlock {
		log.Warn("Warn Call GetPposCacheByNumAndHash of PPOS_TEMP, the Global PPOS_TEMP instance is nil !!!!!!!!!!!!!!!", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex())
		return ppos_storage
	}

	if !notGenesisBlock || (common.Hash{}) == blockHash {
		return ppos_storage
	}

	var storage *Ppos_storage

	temp.lock.Lock()
	if hashTemp, ok := temp.TempMap[blockNumber.String()]; !ok {
		log.Warn("Warn Call GetPposCacheByNumAndHash of PPOS_TEMP, the PPOS storage cache is empty by blockNumber !!!!! Direct short-circuit", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex())
		temp.lock.Unlock()
		return ppos_storage
	}else {

		if pposStorage, ok := hashTemp[blockHash]; !ok {
			log.Warn("Warn Call GetPposCacheByNumAndHash of PPOS_TEMP, the PPOS storage cache is empty by blockHash !!!!! Direct short-circuit", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex())
			temp.lock.Unlock()
			return ppos_storage
		}else {
			storage = pposStorage.Copy()
			temp.lock.Unlock()
		}
	}
	return storage
}

// Set ppos storage cache by same block
func (temp *PPOS_TEMP) SubmitPposCache2Temp(blockNumber, blockInterval *big.Int, blockHash common.Hash, storage *Ppos_storage)  {
	log.Info("Call SubmitPposCache2Temp of PPOS_TEMP", "blockNumber", blockNumber.String(), "blockHash", blockHash.Hex(),
		"blockInterval", blockInterval, "Before SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount)

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
			"blockHash", blockHash.Hex(), "blockInterval", blockInterval, " Before SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount)
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
			 "blockHash", blockHash.Hex(), "blockInterval", blockInterval, " Before SubmitPposCache2Temp, THEN PPOS_TEMP len ", len(temp.TempMap), "Block Count", temp.BlockCount)
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
		temp.lock.Unlock()
		return
	 }else if hasHash && !empty {
		 originHashTemp[blockHash] = storage
		 temp.TempMap[blockNumber.String()] = originHashTemp

		 temp.deleteAnyTemp(blockNumber, blockInterval, blockHash)
		 temp.lock.Unlock()
		 return
	 }
}

func (temp *PPOS_TEMP) Commit2DB(db ethdb.Database, blockNumber *big.Int, blockHash common.Hash) error {
	timer := new(common.Timer)
	timer.Begin()


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
			log.Error("Failed to Commit2DB", "proto err", err)
			return err
		}else {
			if err := db.Put(PPOS_STORAGE_KEY, data); err != nil {
				log.Error("Failed to Call Commit2DB:" + WRITE_PPOS_ERR.Error(), "blockNumber", blockNumber, "blockHash", blockHash, "err", err)
				return WRITE_PPOS_ERR
			}
			log.Info("Call Commit2DB, run time long", "ms: ", timer.End())
			log.Info("Call Commit2DB, write ppos storage data to disk", "blockNumber", blockNumber, "blockHash", blockHash, "data len", len(data))
		}
	}
	return nil
}


func buildPBStorage(blockNumber *big.Int, blockHash common.Hash, ps *Ppos_storage) *PPOSTemp {
	ppos_temp := new(PPOSTemp)
	ppos_temp.BlockNumber = blockNumber.String()
	ppos_temp.BlockHash = blockHash.Hex()

	var empty int = 0  // 0: empty 1: no
	var wg sync.WaitGroup

	// candidate related
	if nil != ps.c_storage {

		canTemp := new(CandidateTemp)


		wg.Add(6)
		// previous witness
		go func() {
			if queue := buildPBcanqueue(ps.c_storage.pres); len(queue) != 0 {
				canTemp.Pres = queue
				empty |= 1
			}
			wg.Done()
		}()
		// current witness
		go func() {
			if queue := buildPBcanqueue(ps.c_storage.currs); len(queue) != 0 {
				canTemp.Currs = queue
				empty |= 1
			}
			wg.Done()
		}()
		// next witness
		go func() {
			if queue := buildPBcanqueue(ps.c_storage.nexts); len(queue) != 0 {
				canTemp.Nexts = queue
				empty |= 1
			}
			wg.Done()
		}()
		// immediate
		go func() {
			if queue := buildPBcanqueue(ps.c_storage.imms); len(queue) != 0 {
				canTemp.Imms = queue
				empty |= 1
			}
			wg.Done()
		}()
		// reserve
		go func() {
			if queue := buildPBcanqueue(ps.c_storage.res); len(queue) != 0 {
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

	// ticket related
	if nil != ps.t_storage {
		tickTemp := new(TicketTemp)

		// SQ
		if ps.t_storage.Sq != -1 {
			empty |= 1
		}

		wg.Add(3)

		// ticketInfos
		go func() {
			if ticketMap := buildPBticketMap(ps.t_storage.Infos); len(ticketMap) != 0 {
				tickTemp.Infos = ticketMap
				empty |= 1
			}
			wg.Done()
		}()

		// ExpireTicket
		go func() {
			if ets := buildPBexpireTicket(ps.t_storage.Ets); len(ets) != 0 {
				tickTemp.Ets = ets
				empty |= 1
			}
			wg.Done()
		}()

		// ticket's attachment  of node
		go func() {
			if dependency := buildPBdependencys(ps.t_storage.Dependencys); len(dependency) != 0 {
				tickTemp.Dependencys = dependency
				empty |= 1
			}
			wg.Done()
		}()

		wg.Wait()
		ppos_temp.TickTmp = tickTemp
	}

	if empty == 0 {
		return nil
	}
	return ppos_temp
}


func buildPBcanqueue (canQqueue types.CandidateQueue) []*CandidateInfo {
	if len(canQqueue) == 0 {
		return nil
	}

	pbQueue := make([]*CandidateInfo, len(canQqueue))
	for _, can := range canQqueue {
		canInfo := &CandidateInfo{
			Deposit: 		can.Deposit.String(),
			BlockNumber:	can.BlockNumber.String(),
			TxIndex:		can.TxIndex,
			CandidateId:	can.CandidateId.Bytes(),
			Host:			can.Host,
			Port:			can.Port,
			Owner:			can.Owner.Bytes(),
			Extra:			can.Extra,
			TicketId: 		can.TicketId.Bytes(),
		}
		pbQueue = append(pbQueue, canInfo)
	}
	return pbQueue
}


func buildPBrefunds(refunds refundStorage) map[string]*RefundArr {
	if len(refunds) == 0 {
		return nil
	}

	refundMap := make(map[string]*RefundArr, len(refunds))

	for nodeId, rs := range refunds {

		if len(rs) == 0 {
			continue
		}
		defeats := make([]*Refund, len(rs))
		for _, refund := range rs {
			refundInfo := &Refund{
				Deposit:     	refund.Deposit.String(),
				BlockNumber:	refund.BlockNumber.String(),
				Owner:			refund.Owner.Bytes(),
			}
			defeats = append(defeats, refundInfo)
		}

		refundArr := &RefundArr{
			Defeats: defeats,
		}
		refundMap[nodeId.String()] = refundArr
	}
	return refundMap
}


func buildPBticketMap(tickets map[common.Hash]*types.Ticket) map[string]*TicketInfo {
	if len(tickets) == 0 {
		return nil
	}

	pb_ticketMap := make(map[string]*TicketInfo, len(tickets))

	for tid, tinfo := range tickets {
		ticketInfo := &TicketInfo{
			Owner:       tinfo.Owner.Bytes(),
			Deposit:     tinfo.Deposit.String(),
			CandidateId: tinfo.CandidateId.Bytes(),
			BlockNumber: tinfo.BlockNumber.String(),
			Remaining:   tinfo.Remaining,
		}
		pb_ticketMap[tid.String()] = ticketInfo
	}
	return pb_ticketMap
}


func buildPBexpireTicket(ets map[string][]common.Hash) map[string]*TxHashArr  {
	if len(ets) == 0 {
		return nil
	}

	pb_ets := make(map[string]*TxHashArr, len(ets))

	for blockNumber, ticketIdArr := range ets {

		if len(ticketIdArr) == 0 {
			continue
		}

		txHashs := make([][]byte, len(ticketIdArr))

		for _, tid := range ticketIdArr {
			txHashs = append(txHashs, tid.Bytes())
		}

		txHashArr := new(TxHashArr)
		txHashArr.TxHashs = txHashs
		pb_ets[blockNumber] = txHashArr
	}
	return pb_ets
}


func buildPBdependencys(dependencys map[discover.NodeID]*ticketDependency) map[string]*TicketDependency {
	if len(dependencys) == 0 {
		return nil
	}

	pb_dependency := make(map[string]*TicketDependency, len(dependencys))

	for nodeId, dependency := range dependencys {

		tidArr := make([][]byte, len(dependency.Tids))

		for _, ticketId := range dependency.Tids {
			tidArr = append(tidArr, ticketId.Bytes())
		}

		depenInfo := &TicketDependency{
			Age:  dependency.Age,
			Num:  dependency.Num,
			Tids: tidArr,
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
		if tickStorage.Sq == -1 && len(tickStorage.Infos) == 0 &&
			len(tickStorage.Ets) == 0 && len(tickStorage.Dependencys) == 0 {
			tickEmpty = true
		}
	}
	if canEmpty && tickEmpty {
		return true
	}
	return false
}
