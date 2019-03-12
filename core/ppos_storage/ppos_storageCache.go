package ppos_storage

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
	"sync"
)

const (
	PREVIOUS = iota
	CURRENT
	NEXT
	IMMEDIATE
	RESERVE
	REFUND
)

var (
	ParamsIllegalErr = errors.New("Params illegal")
)

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
	Age *big.Int
	// ticket count
	Num uint64
	// ticketIds
	Tids []common.Hash
}

func (td *ticketDependency) AddAge(number *big.Int) {
	if nil != td.Age {
		td.Age.Add(td.Age, number)
	}
}

func (td *ticketDependency) SubAge(number *big.Int) {
	if nil != td.Age {
		if td.Age.Cmp(number) >= 0 && number.Uint64() > 0 {
			td.Age.Sub(td.Age, number)
		}
	}
}

func (td *ticketDependency) SubNum() {
	if td.Num > 0 {
		td.Num--
	}
}

type ticket_temp struct {
	// total remian  k-v
	Sq int
	// ticketInfo  map[txHash]ticketInfo
	Infos map[common.Hash]*types.Ticket
	// ExpireTicket  map[blockNumber]txHash
	Ets map[string][]common.Hash
	// ticket's attachment  of node
	Dependencys map[discover.NodeID]*ticketDependency
}

type Ppos_storage struct {
	c_storage *candidate_temp
	t_storage *ticket_temp
}

func (ps *Ppos_storage) Copy() *Ppos_storage {
	ppos_storage := &Ppos_storage{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		if nil != ps.c_storage {
			ppos_storage.c_storage = ps.CopyCandidateStorage()
		}
		wg.Done()
	}()

	go func() {

		// TICKET RELATED
		wg.Done()
	}()
	wg.Wait()
	return ppos_storage
}


func GetPPOS_storage () *Ppos_storage {

	cache := new(Ppos_storage)

	can_cache := new(candidate_temp)
	ticket_cache := new(ticket_temp)

	ticket_cache.Sq = -1
	ticket_cache.Infos = make(map[common.Hash]*types.Ticket)
	ticket_cache.Ets = make(map[string][]common.Hash)
	ticket_cache.Dependencys = make(map[discover.NodeID]*ticketDependency)

	queue := make(types.CandidateQueue, 0)
	refund := make(refundStorage, 0)
	can_cache.pres = queue
	can_cache.currs = queue
	can_cache.nexts = queue
	can_cache.imms = queue
	can_cache.res = queue
	can_cache.refunds = refund

	cache.c_storage = can_cache
	cache.t_storage = ticket_cache
	return cache
}

/** candidate related func */


func (p *Ppos_storage) CopyCandidateStorage ()  *candidate_temp {

	temp := new(candidate_temp)

	type result struct {
		Status int
		Data   interface{}
	}
	var wg sync.WaitGroup
	wg.Add(6)
	resCh := make(chan *result, 5)

	loadQueueFunc := func(flag int) {
		res := new(result)
		switch flag {
		case PREVIOUS:
			res.Status = PREVIOUS
			res.Data = p.c_storage.pres.DeepCopy()
		case CURRENT:
			res.Status = CURRENT
			res.Data = p.c_storage.currs.DeepCopy()
		case NEXT:
			res.Status = NEXT
			res.Data = p.c_storage.nexts.DeepCopy()
		case IMMEDIATE:
			res.Status = IMMEDIATE
			res.Data = p.c_storage.imms.DeepCopy()
		case RESERVE:
			res.Status = RESERVE
			res.Data = p.c_storage.res.DeepCopy()
		}
		resCh <- res
		wg.Done()
	}

	go loadQueueFunc(PREVIOUS)
	go loadQueueFunc(CURRENT)
	go loadQueueFunc(NEXT)
	go loadQueueFunc(IMMEDIATE)
	go loadQueueFunc(RESERVE)

	go func() {
		res := new(result)
		cache := make(refundStorage, len(p.c_storage.refunds))
		for nodeId, queue := range p.c_storage.refunds {
			temp.refunds[nodeId] = queue.DeepCopy()
		}
		res.Status = REFUND
		res.Data = cache
		resCh <- res
		wg.Done()
	}()
	wg.Wait()
	close(resCh)
	for res := range resCh {
		switch res.Status {
		case PREVIOUS:
			temp.pres = res.Data.(types.CandidateQueue)
		case CURRENT:
			temp.currs = res.Data.(types.CandidateQueue)
		case NEXT:
			temp.nexts = res.Data.(types.CandidateQueue)
		case IMMEDIATE:
			temp.imms = res.Data.(types.CandidateQueue)
		case RESERVE:
			temp.res = res.Data.(types.CandidateQueue)
		case REFUND:
			temp.refunds = res.Data.(refundStorage)
		}
	}
	return temp
}

func (p *Ppos_storage) CopyTicketStorage() *ticket_temp {
	ticket_cache := new(ticket_temp)

	ticket_cache.Sq = p.t_storage.Sq
	ticket_cache.Infos = make(map[common.Hash]*types.Ticket)
	ticket_cache.Ets = make(map[string][]common.Hash)
	ticket_cache.Dependencys = make(map[discover.NodeID]*ticketDependency)

	for key := range p.t_storage.Infos {
		ticket := p.t_storage.Infos[key]
		ticket_cache.Infos[key] = ticket.DeepCopy()
	}
	for key := range p.t_storage.Ets {
		ticket_cache.Ets[key] = p.t_storage.Ets[key][:]
	}
	for key := range p.t_storage.Dependencys {
		temp := *p.t_storage.Dependencys[key]
		ticket_cache.Dependencys[key] = &temp
	}
	return ticket_cache
}

// Get CandidateQueue
// flag:
// 0: previous witness
// 1: current witness
// 2: next witness
// 3: immediate
// 4: reserve
func (p *Ppos_storage) GetCandidateQueue(flag int) (types.CandidateQueue, error) {
	switch flag {
	case PREVIOUS:
		return p.c_storage.pres, nil
	case CURRENT:
		return p.c_storage.currs, nil
	case NEXT:
		return p.c_storage.nexts, nil
	case IMMEDIATE:
		return p.c_storage.imms, nil
	case RESERVE:
		return p.c_storage.res, nil
	default:
		return nil, ParamsIllegalErr
	}
}

// Set CandidateQueue
func (p *Ppos_storage) SetCandidateQueue(queue types.CandidateQueue, flag int) error {
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
	default:
		return ParamsIllegalErr

	}
	return nil
}

// Delete CandidateQueue
func (p *Ppos_storage) DelCandidateQueue(flag int) error {
	switch flag {
	case PREVIOUS:
		p.c_storage.pres = nil
	case CURRENT:
		p.c_storage.currs = nil
	case NEXT:
		p.c_storage.nexts = nil
	case IMMEDIATE:
		p.c_storage.imms = nil
	case RESERVE:
		p.c_storage.res = nil
	default:
		return ParamsIllegalErr

	}
	return nil
}

// Get Refund
func (p *Ppos_storage) GetRefund(nodeId discover.NodeID) types.RefundQueue {
	if queue, ok := p.c_storage.refunds[nodeId]; ok {
		return queue
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
		queue = make(types.RefundQueue, 1)
		queue[0] = refund
		p.c_storage.refunds[nodeId] = queue
	}
}

// Delete Refund
func (p *Ppos_storage) DelRefund(nodeId discover.NodeID) {
	delete(p.c_storage.refunds, nodeId)
}

/** ticket related func */

// Get total remian
func (p *Ppos_storage) GetTotalRemian() int {
	return p.t_storage.Sq
}

// Set total remain
func (p *Ppos_storage) SetTotalRemain(count int) error {
	p.t_storage.Sq = count
	return nil
}

// Get TicketInfo
func (p *Ppos_storage) GetTicketInfo(txHash common.Hash) (*types.Ticket, error) {
	ticket, ok := p.t_storage.Infos[txHash]
	if ok {
		return ticket, nil
	}
	return nil, nil
}

//Set TicketInfo
func (p *Ppos_storage) SetTicketInfo(txHash common.Hash, ticket *types.Ticket) error {
	p.t_storage.Infos[txHash] = ticket
	return nil
}

func (p *Ppos_storage) removeTicketInfo(txHash common.Hash) {
	delete(p.t_storage.Infos, txHash)
}

//GetTiketArr
func (p *Ppos_storage) GetTicketArr(txHashs ...common.Hash) ([]*types.Ticket, error) {
	tickets := make([]*types.Ticket, 0)
	if len(txHashs) > 0 {
		for index := range txHashs {
			if ticket, err := p.GetTicketInfo(txHashs[index]); nil != err && ticket != nil {
				newTicket := *ticket
				tickets = append(tickets, &newTicket)
			}
		}
	}
	return tickets, nil
}

//Get ExpireTicket
func (p *Ppos_storage) GetExpireTicket(blockNumber *big.Int) ([]common.Hash, error) {
	ids, ok := p.t_storage.Ets[blockNumber.String()]
	if ok {
		return ids, nil
	}
	return nil, nil
}

// Set ExpireTicket
func (p *Ppos_storage) SetExpireTicket(blockNumber *big.Int, txHash common.Hash) error {
	ids, ok := p.t_storage.Ets[blockNumber.String()]
	if !ok {
		ids = make([]common.Hash, 0)
	}
	ids = append(ids, txHash)
	p.t_storage.Ets[blockNumber.String()] = ids
	return nil
}

func (p *Ppos_storage) RemoveExpireTicket(blockNumber *big.Int, txHash common.Hash) error {
	ids, ok := p.t_storage.Ets[blockNumber.String()]
	if ok {
		ids = removeTicketId(txHash, ids)
		if ids == nil {
			delete(p.t_storage.Ets, blockNumber.String())
		} else {
			p.t_storage.Ets[blockNumber.String()] = ids
		}
	}
	return nil
}

//Get ticket dependency
func (p *Ppos_storage) GetTicketDependency(nodeId discover.NodeID) (*ticketDependency, error) {
	value, ok := p.t_storage.Dependencys[nodeId]
	if ok {
		return value, nil
	}
	return nil, nil
}

// Set ticket dependency
func (p *Ppos_storage) SetTicketDependency(nodeId discover.NodeID, ependency *ticketDependency) error {
	p.t_storage.Dependencys[nodeId] = ependency
	return nil
}

func (p *Ppos_storage) RemoveTicketDependency(nodeId discover.NodeID) error {
	delete(p.t_storage.Dependencys, nodeId)
	return nil
}

func (p *Ppos_storage) GetCandidateTxHashs(nodeId discover.NodeID) ([]common.Hash, error) {
	value, ok := p.t_storage.Dependencys[nodeId]
	if ok {
		return value.Tids, nil
	}
	return nil, nil
}

func (p *Ppos_storage) AppendTicket(nodeId discover.NodeID, txHash common.Hash, ticket *types.Ticket) error {
	if err := p.SetTicketInfo(txHash, ticket); nil != err {
		return err
	}
	if value, err := p.GetTicketDependency(nodeId); nil != err {
		return err
	} else {
		if nil == value {
			value = new(ticketDependency)
			value.Tids = make([]common.Hash, 0)
		}
		value.Num += ticket.Remaining
		value.Age = new(big.Int)
		value.Tids = append(value.Tids, txHash)
		p.SetTicketDependency(nodeId, value)
	}
	return nil
}

func (p *Ppos_storage) SubTicket(nodeId discover.NodeID, txHash common.Hash) error {
	if value, err := p.GetTicketDependency(nodeId); nil != err {
		return err
	} else {
		if nil != value {
			ticket, terr := p.GetTicketInfo(txHash)
			if nil != terr || ticket == nil {
				return terr
			}
			ticket.SubRemaining()
			value.SubNum()
			if ticket.Remaining == 0 {
				p.removeTicketInfo(txHash)
				for index := range value.Tids {
					if value.Tids[index] == txHash {
						start := value.Tids[:index]
						end := value.Tids[index+1:]
						value.Tids = append(start, end...)
						break
					}
				}
			} else {
				if err := p.SetTicketInfo(txHash, ticket); nil != err {
					return err
				}
			}
		}
	}
	return nil
}

func (p *Ppos_storage) RemoveTicket(nodeId discover.NodeID, txHash common.Hash) error {
	ticket, terr := p.GetTicketInfo(txHash)
	if nil != terr || ticket == nil {
		return terr
	}
	if value, err := p.GetTicketDependency(nodeId); nil != err {
		return err
	} else {
		if nil != value {
			value.Num -= ticket.Remaining
			value.Tids = removeTicketId(txHash, value.tIds)
		}
		p.removeTicketInfo(txHash)
	}
	return nil
}

func (p *Ppos_storage) GetCandidateTicketCount(nodeId discover.NodeID) uint64 {
	if value, err := p.GetTicketDependency(nodeId); nil == err && value != nil {
		return value.Num
	}
	return 0
}

func (p *Ppos_storage) GetCandidateTicketAge(nodeId discover.NodeID) (*big.Int, error) {
	if value, err := p.GetTicketDependency(nodeId); nil != err {
		return nil, err
	} else {
		if value != nil {
			return value.Age, nil
		}
	}
	return nil, nil
}

func (p *Ppos_storage) SetCandidateTicketAge(nodeId discover.NodeID, age *big.Int) error {
	if value, err := p.GetTicketDependency(nodeId); nil != err {
		return err
	} else {
		if value != nil {
			value.Age = age
		}
	}
	return nil
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
