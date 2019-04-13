package ppos_storage

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/golang/protobuf/proto"
	"math/big"
	"errors"
	"fmt"
	"sort"
	"sync"
	"crypto/md5"
)

const (
	PREVIOUS = iota
	CURRENT
	NEXT
	IMMEDIATE
	RESERVE
)

var (
	ParamsIllegalErr 			= errors.New("Params illegal")
	TicketNotFindErr        	= errors.New("The Ticket not find")
)

var ticketCache = sync.Map{}

func PutTicket(txHash common.Hash, ticket *types.Ticket) {
	ticketCache.Store(txHash, ticket)
}

func GetTicket(txHash common.Hash) *types.Ticket {
	if value, ok := ticketCache.Load(txHash); ok {
		return value.(*types.Ticket)
	}
	return nil
}

func RemoveTicket(txHash common.Hash) {
	ticketCache.Delete(txHash)
}

type refundStorage map[discover.NodeID]types.RefundQueue

type candidate_temp struct {
	// previous witness
	pres types.CandidateQueue
	// current witness
	currs types.CandidateQueue
	// next witness
	nexts types.CandidateQueue
	//immediate
	imms types.CandidateQueue
	// reserve
	res types.CandidateQueue
	// refund
	refunds refundStorage
}

type ticketDependency struct {
	// ticket age
	//Age uint64
	// ticket count
	Num uint32
	// ticketIds
	Tinfo []*ticketInfo
}

type ticketInfo struct {
	TxHash			common.Hash
	// The number of remaining tickets
	Remaining		uint32
	Price 			*big.Int
}

func (t *ticketInfo) SubRemaining() {
	if t.Remaining > 0 {
		t.Remaining--
	}
}

//func (td *ticketDependency) AddAge(number uint64) {
//	if number > 0 {
//		td.Age += number
//	}
//}
//
//func (td *ticketDependency) SubAge(number uint64) {
//	if number > 0 && td.Age >= number {
//		td.Age -= number
//	}
//}

func (td *ticketDependency) subNum() {
	if td.Num > 0 {
		td.Num--
	}
}

type ticket_temp struct {
	// total remian  k-v
	Sq int32
	// ticketInfo  map[txHash]ticketInfo
	//Infos map[common.Hash]*types.Ticket
	// ExpireTicket  map[blockNumber]txHash
	//Ets map[string][]common.Hash
	// ticket's attachment  of node
	Dependencys map[discover.NodeID]*ticketDependency
}

type Ppos_storage struct {
	c_storage *candidate_temp
	t_storage *ticket_temp
}

func (ps *Ppos_storage) Copy() *Ppos_storage {

	//if  verifyStorageEmpty(ps) {
	//	return NewPPOS_storage()
	//}

	if nil == ps {
		return NewPPOS_storage()
	}

	ppos_storage := &Ppos_storage{
		c_storage: ps.CopyCandidateStorage(),
		t_storage: ps.CopyTicketStorage(),
	}
	return ppos_storage
}


func NewPPOS_storage () *Ppos_storage {

	cache := &Ppos_storage{
		c_storage: &candidate_temp{
			pres: 	make(types.CandidateQueue, 0),
			currs:  make(types.CandidateQueue, 0),
			nexts: 	make(types.CandidateQueue, 0),
			imms: 	make(types.CandidateQueue, 0),
			res: 	make(types.CandidateQueue, 0),
			refunds: make(refundStorage, 0),
		},

		t_storage: &ticket_temp{
			Sq: 	-1,
			Dependencys: 	make(map[discover.NodeID]*ticketDependency),
		},
	}

	/*cache := new(Ppos_storage)

	c := new(candidate_temp)
	t:= new(ticket_temp)

	c.pres = make(types.CandidateQueue, 0)
	c.currs = make(types.CandidateQueue, 0)
	c.nexts = make(types.CandidateQueue, 0)

	c.imms = make(types.CandidateQueue, 0)
	c.res = make(types.CandidateQueue, 0)

	c.refunds = make(refundStorage, 0)

	t.Sq = -1
	t.Dependencys = make(map[discover.NodeID]*ticketDependency)

	cache.c_storage = c
	cache.t_storage = t*/


	return cache
}

/** candidate related func */

//func (p *Ppos_storage) CopyCandidateStorage ()  *candidate_temp {
//	start := common.NewTimer()
//	start.Begin()
//
//	temp := new(candidate_temp)
//
//	type result struct {
//		Status int
//		Data   interface{}
//	}
//	var wg sync.WaitGroup
//	wg.Add(6)
//	resCh := make(chan *result, 6)
//
//	loadQueueFunc := func(flag int) {
//		res := new(result)
//		switch flag {
//		case PREVIOUS:
//			res.Status = PREVIOUS
//			res.Data = p.c_storage.pres.DeepCopy()
//			resCh <- res
//		case CURRENT:
//			res.Status = CURRENT
//			res.Data = p.c_storage.currs.DeepCopy()
//			resCh <- res
//		case NEXT:
//			res.Status = NEXT
//			res.Data = p.c_storage.nexts.DeepCopy()
//			resCh <- res
//		case IMMEDIATE:
//			res.Status = IMMEDIATE
//			res.Data = p.c_storage.imms.DeepCopy()
//			resCh <- res
//		case RESERVE:
//			res.Status = RESERVE
//			res.Data = p.c_storage.res.DeepCopy()
//			resCh <- res
//		}
//		wg.Done()
//	}
//
//	go loadQueueFunc(PREVIOUS)
//	go loadQueueFunc(CURRENT)
//	go loadQueueFunc(NEXT)
//	go loadQueueFunc(IMMEDIATE)
//	go loadQueueFunc(RESERVE)
//
//	go func() {
//		res := new(result)
//		cache := make(refundStorage, len(p.c_storage.refunds))
//		for nodeId, queue := range p.c_storage.refunds {
//			cache[nodeId] = queue.DeepCopy()
//		}
//		res.Status = REFUND
//		res.Data = cache
//		resCh <- res
//		wg.Done()
//	}()
//	wg.Wait()
//	close(resCh)
//	for res := range resCh {
//		switch res.Status {
//		case PREVIOUS:
//			temp.pres = res.Data.(types.CandidateQueue)
//		case CURRENT:
//			temp.currs = res.Data.(types.CandidateQueue)
//		case NEXT:
//			temp.nexts = res.Data.(types.CandidateQueue)
//		case IMMEDIATE:
//			temp.imms = res.Data.(types.CandidateQueue)
//		case RESERVE:
//			temp.res = res.Data.(types.CandidateQueue)
//		case REFUND:
//			temp.refunds = res.Data.(refundStorage)
//		}
//	}
//	log.Debug("CopyCandidateStorage", "Time spent", fmt.Sprintf("%v ms", start.End()))
//	return temp
//}

func (p *Ppos_storage) CopyCandidateStorage ()  *candidate_temp {
	start := common.NewTimer()
	start.Begin()

	cache := make(refundStorage, len(p.c_storage.refunds))
	for nodeId, queue := range p.c_storage.refunds {
		cache[nodeId] = queue.DeepCopy()
	}

	temp := &candidate_temp{
		pres: 		p.c_storage.pres.DeepCopy(),
		currs: 		p.c_storage.currs.DeepCopy(),
		nexts: 		p.c_storage.nexts.DeepCopy(),

		imms: 		p.c_storage.imms.DeepCopy(),
		res: 		p.c_storage.res.DeepCopy(),

		refunds: 	cache,
	}

	/*temp := new(candidate_temp)

	temp.pres = p.c_storage.pres.DeepCopy()
	temp.currs = p.c_storage.currs.DeepCopy()
	temp.nexts = p.c_storage.nexts.DeepCopy()

	temp.imms = p.c_storage.imms.DeepCopy()
	temp.res = p.c_storage.res.DeepCopy()

	temp.refunds = cache*/


	log.Debug("CopyCandidateStorage", "Time spent", fmt.Sprintf("%v ms", start.End()))
	return temp
}

func (p *Ppos_storage) CopyTicketStorage() *ticket_temp {

	start := common.NewTimer()
	start.Begin()

	cache := make(map[discover.NodeID]*ticketDependency, len(p.t_storage.Dependencys))

	for key := range p.t_storage.Dependencys {
		temp := p.t_storage.Dependencys[key]
		tinfos := make([]*ticketInfo, len(temp.Tinfo))

		for j, tin := range temp.Tinfo {

			t := &ticketInfo{
				TxHash: 	tin.TxHash,
				Remaining: 	tin.Remaining,
				Price: 		tin.Price,
			}
			tinfos[j] = t

		}

		cache[key] = &ticketDependency{
			temp.Num,
			tinfos,
		}
	}

	ticket_cache := &ticket_temp{
		Sq: 	p.t_storage.Sq,
		Dependencys: 	cache,
	}


	/*ticket_cache := new(ticket_temp)

	ticket_cache.Sq = p.t_storage.Sq
	//ticket_cache.Infos = make(map[common.Hash]*types.Ticket)
	//ticket_cache.Ets = make(map[string][]common.Hash)
	ticket_cache.Dependencys = make(map[discover.NodeID]*ticketDependency)

	//for key := range p.t_storage.Infos {
	//	ticket := p.t_storage.Infos[key]
	//	ticket_cache.Infos[key] = ticket.DeepCopy()
	//}
	//for key := range p.t_storage.Ets {
	//	value := p.t_storage.Ets[key]
	//	list := make([]common.Hash, len(value))
	//	copy(list, value)
	//	ticket_cache.Ets[key] = list
	//}
	//for key := range p.t_storage.Dependencys {
	//	temp := p.t_storage.Dependencys[key]
	//	tids := make([]common.Hash, len(temp.Tids))
	//	copy(tids, temp.Tids)
	//	ticket_cache.Dependencys[key] = &ticketDependency{
	//		//temp.Age,
	//		temp.Num,
	//		tids,
	//	}
	//}
	for key := range p.t_storage.Dependencys {
		temp := p.t_storage.Dependencys[key]

		tinfos := make([]*ticketInfo, len(temp.Tinfo))

		for j, tin := range temp.Tinfo {

			t := &ticketInfo{
				TxHash: 	tin.TxHash,
				Remaining: 	tin.Remaining,
				Price: 		tin.Price,
			}
			tinfos[j] = t

		}
		ticket_cache.Dependencys[key] = &ticketDependency{
			//temp.Age,
			temp.Num,
			tinfos,
		}
	}*/


	log.Debug("CopyTicketStorage", "Time spent", fmt.Sprintf("%v ms", start.End()))
	return ticket_cache
}

// Get CandidateQueue
// flag:
// 0: previous witness
// 1: current witness
// 2: next witness
// 3: immediate
// 4: reserve
func (p *Ppos_storage) GetCandidateQueue(flag int) types.CandidateQueue {
	switch flag {
	case PREVIOUS:
		PrintObject("Pres queue", p.c_storage.pres)

		queueCopy := make(types.CandidateQueue, len(p.c_storage.pres))
		copy(queueCopy, p.c_storage.pres)

		return queueCopy
	case CURRENT:
		PrintObject("Curr queue", p.c_storage.currs)

		queueCopy := make(types.CandidateQueue, len(p.c_storage.currs))
		copy(queueCopy, p.c_storage.currs)

		return queueCopy
	case NEXT:
		PrintObject("Next queue", p.c_storage.nexts)

		queueCopy := make(types.CandidateQueue, len(p.c_storage.nexts))
		copy(queueCopy, p.c_storage.nexts)

		return queueCopy
	case IMMEDIATE:
		PrintObject("Imms queue", p.c_storage.imms)

		queueCopy := make(types.CandidateQueue, len(p.c_storage.imms))
		copy(queueCopy, p.c_storage.imms)

		return queueCopy
	case RESERVE:
		PrintObject("Res queue", p.c_storage.res)

		queueCopy := make(types.CandidateQueue, len(p.c_storage.res))
		copy(queueCopy, p.c_storage.res)

		return queueCopy
	default:
		return nil
	}
}

// Set CandidateQueue
func (p *Ppos_storage) SetCandidateQueue(queue types.CandidateQueue, flag int) {
	switch flag {
	case PREVIOUS:
		p.c_storage.pres = queue
	case CURRENT:
		p.c_storage.currs = queue
	case NEXT:
		p.c_storage.nexts = queue
	case IMMEDIATE:
		p.c_storage.imms = queue
	case RESERVE:
		p.c_storage.res = queue
	}
}

// Delete CandidateQueue
func (p *Ppos_storage) DelCandidateQueue(flag int)  {
	switch flag {
	case PREVIOUS:
		p.c_storage.pres = make(types.CandidateQueue, 0)
	case CURRENT:
		p.c_storage.currs = make(types.CandidateQueue, 0)
	case NEXT:
		p.c_storage.nexts = make(types.CandidateQueue, 0)
	case IMMEDIATE:
		p.c_storage.imms = make(types.CandidateQueue, 0)
	case RESERVE:
		p.c_storage.res = make(types.CandidateQueue, 0)
	}
}

// Get Refund
func (p *Ppos_storage) GetRefunds(nodeId discover.NodeID) types.RefundQueue {
	if queue, ok := p.c_storage.refunds[nodeId]; ok {
		queueCopy := make(types.RefundQueue, len(queue))
		copy(queueCopy, queue)
		return queueCopy
	} else {
		return make(types.RefundQueue, 0)
	}
}



// Set Refund
func (p *Ppos_storage) SetRefund(nodeId discover.NodeID, refund *types.CandidateRefund) {

	if queue, ok := p.c_storage.refunds[nodeId]; ok {
		queue = append(queue, refund)
		p.c_storage.refunds[nodeId] = queue
	} else {
		queue = make(types.RefundQueue, 0)
		//queue[0] = refund
		queue = append(queue, refund)
		p.c_storage.refunds[nodeId] = queue
	}
}

func (p *Ppos_storage) SetRefunds(nodeId discover.NodeID, refundArr types.RefundQueue) {
	if len(refundArr) == 0 {
		delete(p.c_storage.refunds, nodeId)
		return
	}
	p.c_storage.refunds[nodeId] = refundArr
}

func (p *Ppos_storage) AppendRefunds(nodeId discover.NodeID, refundArr types.RefundQueue) {
	if len(refundArr) == 0 {
		return
	}
	if queue, ok := p.c_storage.refunds[nodeId]; ok {
		queue = append(queue, refundArr ...)
		p.c_storage.refunds[nodeId] = queue
	} else {
		p.c_storage.refunds[nodeId] = refundArr
	}
}

// Delete RefundArr
func (p *Ppos_storage) DelRefunds(nodeId discover.NodeID) {
	delete(p.c_storage.refunds, nodeId)
}

/** ticket related func */

// Get total remian
func (p *Ppos_storage) GetTotalRemian() int32 {
	return p.t_storage.Sq
}

// Set total remain
func (p *Ppos_storage) SetTotalRemain(count int32) {
	p.t_storage.Sq = count
}

// Get TicketInfo
/*func (p *Ppos_storage) GetTicketInfo(txHash common.Hash) *types.Ticket {
	ticket, ok := p.t_storage.Infos[txHash]
	if ok {
		return ticket
	}
	return nil
}*/

func (p *Ppos_storage) GetTicketInfo(txHash common.Hash) *ticketInfo {
	for _, obj := range p.t_storage.Dependencys {
		for index := range obj.Tinfo {
			 tinfo := obj.Tinfo[index]
			if tinfo.TxHash == txHash {
				return tinfo
			}
		}
	}
	return nil
}

//Set TicketInfo
//func (p *Ppos_storage) SetTicketInfo(txHash common.Hash, ) {
//	p.t_storage.Infos[txHash] = ticket
//}
//
//func (p *Ppos_storage) removeTicketInfo(txHash common.Hash) {
//	delete(p.t_storage.Infos, txHash)
//}

//GetTiketArr
//func (p *Ppos_storage) GetTicketArr(txHashs ...common.Hash) []*types.Ticket {
//	tickets := make([]*types.Ticket, 0)
//	if len(txHashs) > 0 {
//		for index := range txHashs {
//			if ticket := p.GetTicketInfo(txHashs[index]); ticket != nil {
//				newTicket := *ticket
//				tickets = append(tickets, &newTicket)
//			}
//		}
//	}
//	return tickets
//}

//Get ExpireTicket
//func (p *Ppos_storage) GetExpireTicket(blockNumber *big.Int) []common.Hash {
//	ids, ok := p.t_storage.Ets[blockNumber.String()]
//	if ok {
//		return ids
//	}
//	return nil
//}
//
//// Set ExpireTicket
//func (p *Ppos_storage) SetExpireTicket(blockNumber *big.Int, txHash common.Hash) {
//	ids, ok := p.t_storage.Ets[blockNumber.String()]
//	if !ok {
//		ids = make([]common.Hash, 0)
//	}
//	ids = append(ids, txHash)
//	p.t_storage.Ets[blockNumber.String()] = ids
//}
//
//func (p *Ppos_storage) RemoveExpireTicket(blockNumber *big.Int, txHash common.Hash) {
//	ids, ok := p.t_storage.Ets[blockNumber.String()]
//	if ok {
//		ids = removeTicketId(txHash, ids)
//		if len(ids) == 0 {
//			delete(p.t_storage.Ets, blockNumber.String())
//		} else {
//			p.t_storage.Ets[blockNumber.String()] = ids
//		}
//	}
//}

//Get ticket dependency
func (p *Ppos_storage) GetTicketDependency(nodeId discover.NodeID) *ticketDependency {
	value, ok := p.t_storage.Dependencys[nodeId]
	if ok {
		return value
	}
	return nil
}

// Set ticket dependency
func (p *Ppos_storage) SetTicketDependency(nodeId discover.NodeID, ependency *ticketDependency) {
	p.t_storage.Dependencys[nodeId] = ependency
}

func (p *Ppos_storage) RemoveTicketDependency(nodeId discover.NodeID) {
	delete(p.t_storage.Dependencys, nodeId)
}

/*func (p *Ppos_storage) GetCandidateTxHashs(nodeId discover.NodeID) []common.Hash {
	value, ok := p.t_storage.Dependencys[nodeId]
	if ok {
		return value.Tids
	}
	return nil
}*/

func (p *Ppos_storage) GetCandidateTxHashs(nodeId discover.NodeID) []common.Hash {
	value, ok := p.t_storage.Dependencys[nodeId]
	if ok {
		tids := make([]common.Hash, 0)
		for index := range value.Tinfo {
			tids = append(tids, value.Tinfo[index].TxHash)
		}
		return tids
	}
	return nil
}

/*func (p *Ppos_storage) AppendTicket(nodeId discover.NodeID, txHash common.Hash, ticket *types.Ticket) error {
	p.SetTicketInfo(txHash, ticket)
	value := p.GetTicketDependency(nodeId)
	if nil == value {
		value = new(ticketDependency)
		value.Tids = make([]common.Hash, 0)
	}
	value.Num += ticket.Remaining
	value.Tids = append(value.Tids, txHash)
	p.SetTicketDependency(nodeId, value)
	return nil
}*/

func (p *Ppos_storage) AppendTicket(nodeId discover.NodeID, txHash common.Hash, count uint32, price *big.Int) error {
	value := p.GetTicketDependency(nodeId)
	if nil == value {
		value = new(ticketDependency)
		value.Tinfo = make([]*ticketInfo, 0)
	}
	value.Num += count
	tinfo := &ticketInfo{
		txHash,
		count,
		price,
	}
	value.Tinfo = append(value.Tinfo, tinfo)
	p.SetTicketDependency(nodeId, value)
	return nil
}


/*func (p *Ppos_storage) SubTicket(nodeId discover.NodeID, txHash common.Hash) error {
	value := p.GetTicketDependency(nodeId)
	if nil != value {
		ticket := p.GetTicketInfo(txHash)
		if ticket == nil {
			return TicketNotFindErr
		}
		ticket.SubRemaining()
		value.subNum()
		if ticket.Remaining == 0 {
			p.removeTicketInfo(txHash)
			if list := removeTicketId(txHash, value.Tids); len(list) > 0 {
				value.Tids = list
			} else {
				value.Tids = make([]common.Hash, 0)
			}
		} else {
			p.SetTicketInfo(txHash, ticket)
		}
	}
	return nil
}*/

func (p *Ppos_storage) SubTicket(nodeId discover.NodeID, txHash common.Hash) (*ticketInfo, error) {
	value := p.GetTicketDependency(nodeId)
	if nil != value {
		ticket := p.GetTicketInfo(txHash)
		if ticket == nil {
			return nil, TicketNotFindErr
		}
		ticket.SubRemaining()
		value.subNum()
		if ticket.Remaining == 0 {
			if list := removeTinfo(txHash, value.Tinfo); len(list) > 0 {
				value.Tinfo = list
			} else {
				value.Tinfo = make([]*ticketInfo, 0)
			}
		}
		return ticket, nil
	}
	return nil, nil
}

/*func (p *Ppos_storage) RemoveTicket(nodeId discover.NodeID, txHash common.Hash) error {
	ticket := p.GetTicketInfo(txHash)
	if ticket == nil {
		return TicketNotFindErr
	}
	value := p.GetTicketDependency(nodeId)
	if nil != value {
		value.Num -= ticket.Remaining
		if list := removeTicketId(txHash, value.Tids); len(list) > 0 {
			value.Tids = list
		} else {
			value.Tids = make([]common.Hash, 0)
		}
	}
	p.removeTicketInfo(txHash)
	return nil
}*/

func (p *Ppos_storage) RemoveTicket(nodeId discover.NodeID, txHash common.Hash) (*ticketInfo, error) {
	ticket := p.GetTicketInfo(txHash)
	if ticket == nil {
		return nil, TicketNotFindErr
	}
	value := p.GetTicketDependency(nodeId)
	if nil != value {
		value.Num -= ticket.Remaining
		if list := removeTinfo(txHash, value.Tinfo); len(list) > 0 {
			value.Tinfo = list
		} else {
			value.Tinfo = make([]*ticketInfo, 0)
		}
		if value.Num == 0 {
			p.RemoveTicketDependency(nodeId)
		}
	}
	return ticket, nil
}

func (p *Ppos_storage) GetCandidateTicketCount(nodeId discover.NodeID) uint32 {
	if value := p.GetTicketDependency(nodeId); value != nil {
		log.Debug("Gets the ticket count of node", "nodeId", nodeId.String(), "tcount", value.Num)
		return value.Num
	}
	log.Debug("Gets the ticket count of node", "nodeId", nodeId.String(), "tcount", 0)
	return 0
}

func (p *Ppos_storage) GetCandidateTicketAge(nodeId discover.NodeID) uint64 {
	/*if value := p.GetTicketDependency(nodeId); value != nil {
		return value.Age
	}*/
	return 0
}

func (p *Ppos_storage) SetCandidateTicketAge(nodeId discover.NodeID, age uint64) {
	/*if value := p.GetTicketDependency(nodeId); value != nil {
		value.Age = age
	}*/
}

func (p *Ppos_storage) GetTicketRemainByTxHash(txHash common.Hash) uint32 {
	//PrintObject("Call GetTicketRemainByTxHash", p.t_storage.Dependencys)
	//log.Debug("Call GetTicketRemainByTxHash", "ticketId", txHash.Hex())
	for _, depen := range p.t_storage.Dependencys {
		for _, field := range depen.Tinfo {
			if txHash == field.TxHash {
				return field.Remaining
			}
		}
	}
	return 0
}


func removeTicketId(hash common.Hash, hashs []common.Hash) []common.Hash {
	for index, tempHash := range hashs {
		if tempHash == hash {
			if len(hashs) == 1 {
				return nil
			}
			start := hashs[:index]
			end := hashs[index+1:]
			return append(start, end...)
		}
	}
	return hashs
}

func removeTinfo(hash common.Hash, tinfos []*ticketInfo) []*ticketInfo {
	for index, info := range tinfos {
		if info.TxHash == hash {
			if len(tinfos) == 1 {
				return nil
			}
			start := tinfos[:index]
			end := tinfos[index+1:]
			return append(start, end...)
		}
	}
	return tinfos
}


func (p *Ppos_storage) CalculateHash(blockNumber *big.Int, blockHash common.Hash) (common.Hash, error) {
	log.Debug("Call CalculateHash start ...", "blockNumber", blockNumber, "blockHash", blockHash.Hex())
	start := common.NewTimer()
	start.Begin()

	if verifyStorageEmpty(p) {
		return common.Hash{}, nil
	}

	// declare can refund func
	RefundIdQueueFunc := func(refundMap refundStorage) ([]string, []*RefundArr) {

		PrintObject("RefundIdQueueFunc, Refunds", refundMap)

		if len(refundMap) == 0 {
			return nil, nil
		}

		nodeIdStrArr := make([]string, len(refundMap))

		tempMap := make(map[string]discover.NodeID, len(refundMap))

		var i int = 0

		for nodeId := range refundMap {

			nodeIdStr := nodeId.String()

			nodeIdStrArr[i]= nodeIdStr

			tempMap[nodeIdStr] = nodeId
			
			i ++
		}

		sort.Strings(nodeIdStrArr)

		refundArrQueue := make([]*RefundArr, 0)

		for _, nodeIdStr := range nodeIdStrArr {

			nodeId := tempMap[nodeIdStr]

			rs := refundMap[nodeId]

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

			refundArrQueue = append(refundArrQueue, refundArr)
		}
		return nodeIdStrArr, refundArrQueue
	}

	// declare can dependency func
	DependencyFunc := func(dependencys map[discover.NodeID]*ticketDependency) ([]string, []*TicketDependency) {

		PrintObject("DependencyFunc, dependencys", dependencys)

		if len(dependencys) == 0 {
			return nil, nil
		}

		nodeIdStrArr := make([]string, len(dependencys))

		tempMap := make(map[string]discover.NodeID, len(dependencys))


		var i int = 0
		for nodeId := range dependencys {


			nodeIdStr := nodeId.String()

			nodeIdStrArr[i] = nodeIdStr

			tempMap[nodeIdStr] = nodeId

			i++
		}

		sort.Strings(nodeIdStrArr)

		dependencyArr := make([]*TicketDependency, 0)


		for _, nodeIdStr := range nodeIdStrArr {
			nodeId := tempMap[nodeIdStr]


			depen := dependencys[nodeId]

			if depen.Num == 0 && len(depen.Tinfo) == 0 {
				continue
			}

			fieldArr := make([]*Field, len(depen.Tinfo))


			for i, field := range depen.Tinfo {

				f := &Field{
					TxHash:		field.TxHash.String(),
					Remaining: 	field.Remaining,
					Price: 		field.Price.String(),
				}
				fieldArr[i] = f
			}

			depenInfo := &TicketDependency{
				//Age:  dependency.Age,
				Num:  depen.Num,
				//Tids: tidArr,
				Tinfo: 	fieldArr,
			}

			dependencyArr = append(dependencyArr, depenInfo)

		}

		return nodeIdStrArr, dependencyArr
	}


	sortTemp := new(SortTemp)

	var wg sync.WaitGroup
	wg.Add(7)

	resqueue := make([][]*CandidateInfo, 5)

	/**
	calculate can dependency Hash
	*/
	go func() {
		resqueue[0] = buildPBcanqueue("DependencyFunc, pres", p.c_storage.pres)
		wg.Done()
	}()

	go func() {
		resqueue[1] = buildPBcanqueue("DependencyFunc, currs", p.c_storage.currs)
		wg.Done()
	}()

	go func() {
		resqueue[2] = buildPBcanqueue("DependencyFunc, nexts", p.c_storage.nexts)
		wg.Done()
	}()

	go func() {
		resqueue[3] = buildPBcanqueue("DependencyFunc, imms", p.c_storage.imms)
		wg.Done()
	}()

	go func() {
		resqueue[4] = buildPBcanqueue("DependencyFunc, res", p.c_storage.res)
		wg.Done()
	}()

	go func() {
		refundNodeIdArr, refundArr := RefundIdQueueFunc(p.c_storage.refunds)
		sortTemp.ReIds = refundNodeIdArr
		sortTemp.Refunds = refundArr
		wg.Done()
	}()

	// calculate tick dependency Hash
	go func() {
		dependencyNodeIdArr, dependencyArr := DependencyFunc(p.t_storage.Dependencys)
		sortTemp.NodeIds = dependencyNodeIdArr
		sortTemp.Deps = dependencyArr
		wg.Done()
	}()

	wg.Wait()

	// assemble data
	sortTemp.Sq = p.t_storage.Sq

	for _, canArr := range resqueue {
		if len(canArr) != 0 {
			sortTemp.Cans = append(sortTemp.Cans, canArr...)
		}
	}

	//PrintObject("Call CalculateHash build SortTemp: blockNumber:" + blockNumber.String() + ",blockHash:" + blockHash.Hex() + ", sortTemp", sortTemp)

	log.Debug("Call CalculateHash build SortTemp success ...","blockNumber", blockNumber, "blockHash", blockHash.Hex(), "Build data Time spent", fmt.Sprintf("%v ms", start.End()))

	data, err := proto.Marshal(sortTemp)
	if err != nil {
		log.Error("Failed to Call CalculateHash, protobuf is failed", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),  "err", err)
		return common.Hash{}, err
	}
	log.Debug("Call CalculateHash protobuf success ...", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),  "Made protobuf Time spent", fmt.Sprintf("%v ms", start.End()))
	ret := crypto.Keccak256Hash(data)
	log.Debug("Call CalculateHash finish ...", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "proto out len", len(data), "Hash", string(ret[:]),"md5",  md5.Sum(data),  "ppos storage Hash", ret.Hex(), "Total Time spent", fmt.Sprintf("%v ms", start.End()))

	/*PrintObject("Call CalculateHash Data Cans", sortTemp.Cans)
	PrintObject("Call CalculateHash Data ReIds", sortTemp.ReIds)
	PrintObject("Call CalculateHash Data Refunds", sortTemp.Refunds)
	PrintObject("Call CalculateHash Data Sq", sortTemp.Sq)
	PrintObject("Call CalculateHash Data NodeIds", sortTemp.NodeIds)
	PrintObject("Call CalculateHash Data Deps", sortTemp.Deps)*/

	return ret, nil

}