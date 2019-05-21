package pposm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos_storage"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"sort"
	"strings"
	"sync"
)

var (
	TicketPoolNilErr      = errors.New("Ticket Insufficient quantity")
	TicketPoolOverflowErr = errors.New("Number of ticket pool overflow")
	EncodeTicketErr       = errors.New("Encode Ticket error")
	EncodePoolNumberErr   = errors.New("Encode SurplusQuantity error")
	DecodeTicketErr       = errors.New("Decode Ticket error")
	DecodePoolNumberErr   = errors.New("Decode SurplusQuantity error")
	RecordExpireTicketErr = errors.New("Record Expire Ticket error")
	CandidateNotFindErr   = errors.New("The Candidate not find")
	CandidateNilTicketErr = errors.New("This candidate has no ticket")
	TicketPoolBalanceErr  = errors.New("TicketPool not sufficient funds")
	TicketNotFindErr      = errors.New("The Ticket not find")
	HandleExpireTicketErr = errors.New("Failure to deal with expired tickets")
	GetCandidateAttachErr = errors.New("Get CandidateAttach error")
	SetCandidateAttachErr = errors.New("Update CandidateAttach error")
	VoteTicketErr         = errors.New("Voting failed")
)

type TicketPool struct {
	// Ticket price
	TicketPrice *big.Int
	// Maximum number of ticket pool
	MaxCount uint32
	// Reach expired quantity
	ExpireBlockNumber uint32
	lock              *sync.Mutex
}

//var ticketPool *TicketPool

// initialize the global ticket pool object
func NewTicketPool(configs *params.PposConfig) *TicketPool {
	//if nil != ticketPool {
	//	return ticketPool
	//}
	log.Debug("Build a New TicketPool Info ...")
	if "" == strings.TrimSpace(configs.TicketConfig.TicketPrice) {
		configs.TicketConfig.TicketPrice = "100000000000000000000"
	}
	var ticketPrice *big.Int
	if price, ok := new(big.Int).SetString(configs.TicketConfig.TicketPrice, 10); !ok {
		ticketPrice, _ = new(big.Int).SetString("100000000000000000000", 10)
	} else {
		ticketPrice = price
	}

	ticketPool := &TicketPool{
		TicketPrice:       ticketPrice,
		MaxCount:          configs.TicketConfig.MaxCount,
		ExpireBlockNumber: configs.TicketConfig.ExpireBlockNumber,
		lock:              &sync.Mutex{},
	}
	return ticketPool
}

func (t *TicketPool) VoteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint32, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) (uint32, error) {
	log.Debug("Call Voting", "statedb addr", fmt.Sprintf("%p", stateDB))
	log.Info("Start Voting,VoteTicket", "owner", owner.Hex(), "voteNumber", voteNumber, "price", deposit.Uint64(), "nodeId", nodeId.String(), "blockNumber", blockNumber.Uint64())
	successCount, err := t.voteTicket(stateDB, owner, voteNumber, deposit, nodeId, blockNumber)
	if nil != err {
		log.Error("Voting failed", "nodeId", nodeId.String(), "voteNumber", voteNumber, "successNum", successCount, "err", err)
		return successCount, err
	}
	// Voting completed, candidates reordered
	log.Debug("Successfully voted to start updating the list of candidates,VoteTicket", "successNum", successCount)
	if err := cContext.UpdateElectedQueue(stateDB, blockNumber, nodeId); nil != err {
		log.Error("Failed to Update candidate when voteTicket success", "err", err)
	}
	if successCount > 0 {
		t := &types.Ticket{
			owner,
			deposit,
			nodeId,
			blockNumber,
			0,
		}
		ppos_storage.PutTicket(stateDB.TxHash(), t)
	}
	log.Debug("Successful vote, candidate list updated successfully,VoteTicket", "successNum", successCount)
	return successCount, nil
}

func (t *TicketPool) voteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint32, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) (uint32, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	// check ticket pool count
	surplusQuantity := t.GetPoolNumber(stateDB)
	log.Debug("Execute voteTicket", "surplusQuantity", surplusQuantity, "voteNumber", voteNumber, "blockNumber", blockNumber.Uint64())
	if surplusQuantity == 0 {
		log.Error("Ticket Insufficient quantity")
		return 0, TicketPoolNilErr
	}
	if surplusQuantity < voteNumber {
		voteNumber = surplusQuantity
	}
	log.Debug("Start circular voting", "nodeId", nodeId.String(), "voteNumber", voteNumber)

	ticketId := stateDB.TxHash()
	//ticket := &types.Ticket{
	//	owner,
	//	deposit,
	//	nodeId,
	//	blockNumber,
	//	voteNumber,
	//}
	//t.recordExpireTicket(stateDB, blockNumber, ticketId)
	//log.Debug("Record the success of the ticket to expire, and start reducing the number of tickets", "blockNumber", blockNumber.Uint64(), "surplusQuantity", surplusQuantity, "ticketId", ticketId.Hex())
	t.setPoolNumber(stateDB, surplusQuantity-voteNumber)
	stateDB.GetPPOSCache().AppendTicket(nodeId, ticketId, voteNumber, deposit)
	log.Debug("Voting SUCCUESS !!!!!!  Reduce the remaining amount of the ticket pool successfully", "surplusQuantity", t.GetPoolNumber(stateDB), "nodeId", nodeId.String(), "blockNumber", blockNumber.Uint64(), "ticketId", ticketId.Hex())
	return voteNumber, nil
}

func (t *TicketPool) calcExpireBlockNumber(stateDB vm.StateDB, blockNumber *big.Int) (*big.Int, bool) {
	num := new(big.Int).SetUint64(0)
	if blockNumber.Cmp(new(big.Int).SetUint64(uint64(t.ExpireBlockNumber))) >= 0 {
		num.Sub(blockNumber, new(big.Int).SetUint64(uint64(t.ExpireBlockNumber)))
		return num, true
	}
	return num, false
}

/*func (t *TicketPool) GetExpireTicketIds(stateDB vm.StateDB, blockNumber *big.Int) []common.Hash {
	log.Debug("Call GetExpireTicketIds", "statedb addr", fmt.Sprintf("%p", stateDB))
	return stateDB.GetPPOSCache().GetExpireTicket(blockNumber)
}*/

func (t *TicketPool) GetExpireTicketIds(stateDB vm.StateDB, blockNumber *big.Int) []common.Hash {
	log.Debug("Call GetExpireTicketIds", "statedb addr", fmt.Sprintf("%p", stateDB))
	start := common.NewTimer()
	start.Begin()
	body := tContext.GetBody(blockNumber.Uint64())
	txs := make([]common.Hash, 0)
	for _, tx := range body.Transactions {
		if *tx.To() == common.TicketPoolAddr {
			txs = append(txs, tx.Hash())
		}
	}
	log.Debug("GetExpireTicketIds Time", "Time spent", fmt.Sprintf("%v ms", start.End()))
	return txs
}

// In the current block,
// the ticket id is placed in the value slice with the block height as the key to find the expired ticket.
func (t *TicketPool) recordExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketId common.Hash) {
	//stateDB.GetPPOSCache().SetExpireTicket(blockNumber, ticketId)
}

func (t *TicketPool) removeExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketId common.Hash) {
	//stateDB.GetPPOSCache().RemoveExpireTicket(blockNumber, ticketId)
}

func (t *TicketPool) handleExpireTicket(stateDB vm.StateDB, expireBlockNumber *big.Int, currentBlockNumber *big.Int) ([]discover.NodeID, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	ticketIdList := t.GetExpireTicketIds(stateDB, expireBlockNumber)
	if len(ticketIdList) == 0 {
		return nil, nil
	}
	log.Info("Pending ticket to be processed", "amount", len(ticketIdList), "expireBlockNumber", expireBlockNumber.Uint64(), "currentBlockNumber", currentBlockNumber.Uint64())
	candidateAttachMap := make(map[discover.NodeID]bool)
	changeNodeIdList := make([]discover.NodeID, 0)
	for _, ticketId := range ticketIdList {
		ticket := t.GetTicket(stateDB, ticketId)
		if ticket == nil {
			continue
		}
		_, ok := candidateAttachMap[ticket.CandidateId]
		if !ok {
			candidateAttachMap[ticket.CandidateId] = true
			changeNodeIdList = append(changeNodeIdList, ticket.CandidateId)
		}
		if _, err := t.releaseTxTicket(stateDB, ticket.CandidateId, ticketId, currentBlockNumber); nil != err {
			return changeNodeIdList, err
		}
	}
	return changeNodeIdList, nil
}

// Get ticket list
func (t *TicketPool) GetTicketList(stateDB vm.StateDB, ticketIds []common.Hash) []*types.Ticket {
	log.Debug("Call GetTickList", "statedb addr", fmt.Sprintf("%p", stateDB))
	var tickets []*types.Ticket
	for _, ticketId := range ticketIds {
		ticket := t.GetTicket(stateDB, ticketId)
		if ticket == nil {
			log.Error("find this ticket fail", "ticketId", ticketId.Hex())
			continue
		}
		tickets = append(tickets, ticket)
	}
	return tickets
}

// Get ticket details based on TicketId
func (t *TicketPool) GetTicket(stateDB vm.StateDB, txHash common.Hash) *types.Ticket {
	log.Debug("Call GetTicket", "statedb addr", fmt.Sprintf("%p", stateDB))

	start := common.NewTimer()
	start.Begin()

	if value := ppos_storage.GetTicket(txHash); nil != value {
		return value
	}

	startTx := common.NewTimer()
	startTx.Begin()
	tx, _, blockNumber,_ := tContext.FindTransaction(txHash)
	log.Debug("GetTicket Time Tx", "Time spent", fmt.Sprintf("%v ms", startTx.End()))
	if nil != tx {
		startDecode := common.NewTimer()
		startDecode.Begin()
		var source [][]byte
		if err := rlp.Decode(bytes.NewReader(tx.Data()), &source); nil != err {
			log.Error("Failed to GetTicket", "txHash", txHash.Hex(), "err", err.Error())
			return nil
		}
		if byteutil.BytesToString(source[1]) != "VoteTicket" {
			return nil
		}
		log.Debug("GetTicket Time Decode", "Time spent", fmt.Sprintf("%v ms", startDecode.End()))
		ticket := new(types.Ticket)
		startSigner := common.NewTimer()
		startSigner.Begin()
		signer := types.NewEIP155Signer(tContext.chainConfig.ChainID)
		if addr, err := signer.Sender(tx); nil != err {
			log.Error("Failed to GetTicket, get tx owner is empty !!!!", "tx", tx.Hash().Hex(), "err", err)
			return nil
		} else {
			ticket.Owner = addr
		}
		log.Debug("GetTicket Time startSigner", "Time spent", fmt.Sprintf("%v ms", startSigner.End()))
		//startGetHeader := common.NewTimer()
		//startGetHeader.Begin()
		//block := tContext.GetHeader(blockHash, blockNumber)
		//log.Debug("GetTicket Time startGetHeader", "Time spent", fmt.Sprintf("%v ms", startGetHeader.End()))
		//startGetNewStateDB := common.NewTimer()
		//startGetNewStateDB.Begin()
		//if oldState, err := tContext.GetNewStateDB(block.Root, new(big.Int).SetUint64(blockNumber), blockHash); nil != err {
		//	return nil
		//} else {
		//	ticket.Deposit = t.GetTicketPrice(oldState)
		//}
		//log.Debug("GetTicket Time startGetNewStateDB", "Time spent", fmt.Sprintf("%v ms", startGetNewStateDB.End()))
		ticket.Deposit = t.GetTicketPrice(stateDB)
		ticket.CandidateId = byteutil.BytesToNodeId(source[4])
		ticket.BlockNumber = new(big.Int).SetUint64(blockNumber)
		log.Debug("GetTicket Time", "Time spent", fmt.Sprintf("%v ms", start.End()))
		return ticket
	}else {
		log.Error("Failed to GetTicket, the tx is empty", "txHash", txHash.Hex())
	}
	return nil
}

func (t *TicketPool) GetTicketRemainByTxHash (stateDB vm.StateDB, txHash common.Hash) uint32 {
	return stateDB.GetPPOSCache().GetTicketRemainByTxHash(txHash)
}

//func (t *TicketPool) setTicket(stateDB vm.StateDB, ticketId common.Hash, ticket *types.Ticket) {
//	stateDB.GetPPOSCache().SetTicketInfo(ticketId, ticket)
//}

func (t *TicketPool) DropReturnTicket(stateDB vm.StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Debug("Call DropReturnTicket", "statedb addr", fmt.Sprintf("%p", stateDB))
	log.Info("Start processing tickets for the drop list on DropReturnTicket", "candidateNum", len(nodeIds), "blockNumber", blockNumber.Uint64())
	for _, nodeId := range nodeIds {
		if nodeId == (discover.NodeID{}) {
			continue
		}
		candidateTicketIds := t.GetCandidateTicketIds(stateDB, nodeId)
		if len(candidateTicketIds) == 0 {
			continue
		}
		//epoch := t.GetCandidateEpoch(stateDB, nodeId)
		ticketCount := t.GetCandidateTicketCount(stateDB, nodeId)
		surplusQuantity := t.GetPoolNumber(stateDB)
		log.Debug("Start reducing the number of tickets on DropReturnTicket", "surplusQuantity", surplusQuantity, "candidateTicketIds", ticketCount)
		t.setPoolNumber(stateDB, surplusQuantity+ticketCount)
		log.Debug("Start processing each invalid ticket on DropReturnTicket", "nodeId", nodeId.String(), "ticketSize", ticketCount)
		for _, ticketId := range candidateTicketIds {
			ticket := t.GetTicket(stateDB, ticketId)
			if ticket == nil {
				continue
			}
			if tinfo, err := stateDB.GetPPOSCache().RemoveTicket(nodeId, ticketId); err != nil {
				return err
			} else {
				ticket.Remaining = tinfo.Remaining
				ticket.Deposit = tinfo.Price
			}
			log.Debug("Start transfer on DropReturnTicket", "nodeId", nodeId.String(), "ticketId", ticketId.Hex(), "deposit", ticket.Deposit, "remaining", ticket.Remaining)
			if err := transfer(stateDB, common.TicketPoolAddr, ticket.Owner, ticket.TotalDeposit()); nil != err {
				return err
			}
			//t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
		}
		log.Debug("Delete candidate ticket collection on DropReturnTicket", "nodeId", nodeId.String(), "ticketSize", ticketCount)
		stateDB.GetPPOSCache().RemoveTicketDependency(nodeId)
	}
	log.Debug("End processing the list on DropReturnTicket")
	return nil
}

func (t *TicketPool) ReturnTicket(stateDB vm.StateDB, nodeId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) error {
	log.Debug("Call ReturnTicket", "statedb addr", fmt.Sprintf("%p", stateDB))
	log.Info("Release the selected ticket on ReturnTicket", "nodeId", nodeId.String(), "ticketId", ticketId.Hex(), "blockNumber", blockNumber.Uint64())
	t.lock.Lock()
	defer t.lock.Unlock()
	if ticketId == (common.Hash{}) {
		return TicketNotFindErr
	}
	if nodeId == (discover.NodeID{}) {
		return CandidateNotFindErr
	}
	_, err := t.releaseTicket(stateDB, nodeId, ticketId, blockNumber)
	if nil != err {
		return err
	}
	return nil
}

func (t *TicketPool) releaseTicket(stateDB vm.StateDB, candidateId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	log.Debug("Start executing releaseTicket", "nodeId", candidateId.String(), "ticketId", ticketId.Hex(), "blockNumber", blockNumber.Uint64())
	ticket := t.GetTicket(stateDB, ticketId)
	if ticket == nil {
		return nil, TicketNotFindErr
	}
	log.Debug("releaseTicket,Start Update", "nodeId", candidateId.String(), "ticketId", ticketId.Hex())
	if tinfo, err := stateDB.GetPPOSCache().SubTicket(candidateId, ticketId); err != nil {
		return ticket, err
	} else {
		ticket.Remaining = tinfo.Remaining
		ticket.Deposit = tinfo.Price
	}
	/*if ticket.Remaining == 0 {
		// Remove from pending expire tickets
		log.Debug("releaseTicket, Ticket has been used, deleted from waiting for expiration", "ticketId", ticketId.Hex(), "candidateId", candidateId.String(), "blockNumber", blockNumber.Uint64())
		t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
	}*/
	log.Debug("releaseTicket, end update", "nodeId", candidateId.String())
	surplusQuantity := t.GetPoolNumber(stateDB)
	log.Debug("releaseTicket, start to update the ticket pool", "surplusQuantity", surplusQuantity)
	if err := t.addPoolNumber(stateDB); err != nil {
		return ticket, err
	}
	surplusQuantity = t.GetPoolNumber(stateDB)
	log.Debug("releaseTicket, end the update ticket pool", "surplusQuantity", surplusQuantity)
	//epoch := t.GetCandidateEpoch(stateDB, candidateId)
	//log.Debug("releaseTicket, start updating the total epoch of candidates", "nodeId", candidateId.String(), "totalEpoch", epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	//if err := t.subCandidateEpoch(stateDB, candidateId, ticket.CalcEpoch(blockNumber)); nil != err {
	//	return ticket, err
	//}
	//epoch = t.GetCandidateEpoch(stateDB, candidateId)
	//log.Debug("releaseTicket, the end of the update candidate total epoch", "nodeId", candidateId.String(), "totalEpoch", epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	return ticket, transfer(stateDB, common.TicketPoolAddr, ticket.Owner, ticket.Deposit)
}

func (t *TicketPool) releaseTxTicket(stateDB vm.StateDB, candidateId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	log.Debug("Start executing releaseTxTicket", "nodeId", candidateId.String(), "ticketId", ticketId.Hex(), "blockNumber", blockNumber.Uint64())
	ticket := t.GetTicket(stateDB, ticketId)
	if ticket == nil {
		return nil, TicketNotFindErr
	}
	log.Debug("releaseTxTicket,Start Update", "nodeId", candidateId.String(), "ticketId", ticketId.Hex())
	if tinfo, err := stateDB.GetPPOSCache().RemoveTicket(candidateId, ticketId); err != nil && err != ppos_storage.TicketNotFindErr {
		return ticket, err
	} else {
		if tinfo == nil {
			log.Warn("releaseTxTicket, Not find TicketId", "ticketId", ticketId.Hex(), "nodeId", candidateId.String(), "blockNumber", blockNumber.Uint64())
			return nil, nil
		}
		ticket.Remaining = tinfo.Remaining
		ticket.Deposit = tinfo.Price
	}
	log.Debug("releaseTxTicket, end update", "nodeId", candidateId.String())
	//t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
	surplusQuantity := t.GetPoolNumber(stateDB)
	log.Debug("releaseTxTicket, start to update the ticket pool", "surplusQuantity", surplusQuantity)
	t.setPoolNumber(stateDB, surplusQuantity + ticket.Remaining)
	surplusQuantity = t.GetPoolNumber(stateDB)
	log.Debug("releaseTxTicket, end the update ticket pool", "surplusQuantity", surplusQuantity)
	//epoch := t.GetCandidateEpoch(stateDB, candidateId)
	//log.Debug("releaseTicket, start updating the total epoch of candidates", "nodeId", candidateId.String(), "totalEpoch", epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	//if err := t.subCandidateEpoch(stateDB, candidateId, ticket.TotalEpoch(blockNumber)); nil != err {
	//	return ticket, err
	//}
	//epoch = t.GetCandidateEpoch(stateDB, candidateId)
	//log.Debug("releaseTicket, the end of the update candidate total epoch", "nodeId", candidateId.String(), "totalEpoch", epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	return ticket, transfer(stateDB, common.TicketPoolAddr, ticket.Owner, ticket.TotalDeposit())
}

func (t *TicketPool) Notify(stateDB vm.StateDB, blockNumber *big.Int) error {
	log.Debug("Call Notify", "statedb addr", fmt.Sprintf("%p", stateDB))
	// Check expired tickets

	//ticket := t.GetTicket(stateDB, common.HexToHash("0xafdd2a272c9af265410369bba200960229e6c90e044d8241cbcd6abf8a1706f8"))


	expireBlockNumber, ok := t.calcExpireBlockNumber(stateDB, blockNumber)
	log.Debug("Check expired tickets on Notify", "isOk", ok, "expireBlockNumber", expireBlockNumber.Uint64())
	if ok {
		if nodeIdList, err := t.handleExpireTicket(stateDB, expireBlockNumber, blockNumber); nil != err {
			log.Error("OutBlockNotice method handleExpireTicket error", "blockNumber", blockNumber.Uint64(), "err", err)
			return HandleExpireTicketErr
		} else {
			// Notify the candidate to update the list information after processing the expired ticket
			log.Debug("After processing the expired ticket, start updating the candidate list on Notify", "blockNumber", blockNumber.Uint64(), "nodeIdList", len(nodeIdList))
			if len(nodeIdList) > 0 {
				if err := cContext.UpdateElectedQueue(stateDB, blockNumber, nodeIdList...); nil != err {
					log.Error("Failed to Update candidate when handleExpireTicket success on Notify", "err", err)
				}
			}
		}
	}
	// Increase the total number of epoch for each candidate
	/*log.Debug("Increase the total number of epoch for each candidate on Notify", "blockNumber", blockNumber.Uint64())
	if err := t.calcCandidateEpoch(stateDB, blockNumber); nil != err {
		return err
	}*/
	return nil
}

/*func (t *TicketPool) calcCandidateEpoch(stateDB vm.StateDB, blockNumber *big.Int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	candidateList := cContext.GetCandidatePendArr(stateDB, 0)
	for _, candidate := range candidateList {
		epoch := t.GetCandidateEpoch(stateDB, candidate.CandidateId)
		// Get the total number of votes, increase the total epoch
		ticketCount := stateDB.GetPPOSCache().GetCandidateTicketCount(candidate.CandidateId)
		log.Debug("increase the total epoch", "candidateId", candidate.CandidateId.String(), "blockNumber", blockNumber.Uint64(), "ticketCount", ticketCount, "epoch", epoch)
		if ticketCount > 0 {
			t.addCandidateEpoch(stateDB, candidate.CandidateId, uint64(ticketCount))
			epoch = t.GetCandidateEpoch(stateDB, candidate.CandidateId)
			log.Debug("increase the total epoch success", "candidateId", candidate.CandidateId.String(), "blockNumber", blockNumber.Uint64(), "ticketCount", ticketCount, "epoch", epoch)
		}
	}
	return nil
}*/

// Simple version of the lucky ticket algorithm
// According to the previous block Hash,
// find the first ticket Id which is larger than the Hash. If not found, the last ticket Id is taken.
func (t *TicketPool) SelectionLuckyTicket(stateDB vm.StateDB, nodeId discover.NodeID, blockHash common.Hash) (common.Hash, error) {
	log.Debug("Call SelectionLuckyTicket", "statedb addr", fmt.Sprintf("%p", stateDB))
	candidateTicketIds := t.GetCandidateTicketIds(stateDB, nodeId)
	log.Debug("Start picking lucky tickets on SelectionLuckyTicket", "nodeId", nodeId.String(), "blockHash", blockHash.Hex(), "candidateTicketIds", len(candidateTicketIds))
	luckyTicketId := common.Hash{}
	if len(candidateTicketIds) == 0 {
		return luckyTicketId, nil
	}
	if len(candidateTicketIds) == 1 {
		return candidateTicketIds[0], nil
	}
	decList := make([]float64, 0)
	decMap := make(map[float64]common.Hash, 0)
	for _, ticketId := range candidateTicketIds {
		decNumber := hexutil.HexDec(ticketId.Hex()[2:])
		decList = append(decList, decNumber)
		decMap[decNumber] = ticketId
	}
	sort.Float64s(decList)
	index := findFirstMatch(decList, hexutil.HexDec(blockHash.Hex()[2:]))
	log.Debug("Pick out a lucky ticket on SelectionLuckyTicket", "index", index)
	luckyTicketId = decMap[decList[index]]
	log.Debug("End the selection of lucky tickets on SelectionLuckyTicket", "nodeId", nodeId.String(), "blockHash", blockHash.Hex(), "luckyTicketId", luckyTicketId.Hex(), "candidateTicketIds", len(candidateTicketIds))
	return luckyTicketId, nil
}

func (t *TicketPool) addPoolNumber(stateDB vm.StateDB) error {
	surplusQuantity := t.GetPoolNumber(stateDB)
	if surplusQuantity == t.MaxCount {
		return TicketPoolOverflowErr
	}
	surplusQuantity++
	t.setPoolNumber(stateDB, surplusQuantity)
	return nil
}

func (t *TicketPool) subPoolNumber(stateDB vm.StateDB) error {
	surplusQuantity := t.GetPoolNumber(stateDB)
	if surplusQuantity == 0 {
		return TicketPoolNilErr
	}
	surplusQuantity--
	t.setPoolNumber(stateDB, surplusQuantity)
	return nil
}

func (t *TicketPool) setPoolNumber(stateDB vm.StateDB, surplusQuantity uint32) {
	stateDB.GetPPOSCache().SetTotalRemain(int32(surplusQuantity))
}

func (t *TicketPool) GetPoolNumber(stateDB vm.StateDB) uint32 {
	if val := stateDB.GetPPOSCache().GetTotalRemian(); val >= 0 {
		return uint32(val)
	} else {
		// Default initialization values
		return t.MaxCount
	}
}

/*func (t *TicketPool) subCandidateEpoch(stateDB vm.StateDB, nodeId discover.NodeID, epoch uint64) error {
	dependency := stateDB.GetPPOSCache().GetTicketDependency(nodeId)
	if nil == dependency {
		return CandidateNotFindErr
	}
	dependency.SubAge(epoch)
	return nil
}*/

/*func (t *TicketPool) addCandidateEpoch(stateDB vm.StateDB, nodeId discover.NodeID, epoch uint64) error {
	dependency := stateDB.GetPPOSCache().GetTicketDependency(nodeId)
	if nil == dependency {
		return CandidateNotFindErr
	}
	dependency.AddAge(epoch)
	return nil
}*/

// Get the remaining number of ticket
func (t *TicketPool) GetTicketRemaining(stateDB vm.StateDB, ticketId common.Hash) uint32 {
	return t.GetTicketRemainByTxHash(stateDB, ticketId)
	/*if nil == ticket {
		return 0
	}
	return ticket.Remaining*/
}

// Get the batch remaining number of ticket
func (t *TicketPool) GetBatchTicketRemaining(stateDB vm.StateDB, ticketIds []common.Hash) map[common.Hash]uint32 {
	ticketsRemaining := make(map[common.Hash]uint32, len(ticketIds))
	var wg sync.WaitGroup
	wg.Add(len(ticketIds))

	type result struct {
		id 		common.Hash
		count 	uint32
	}
	resCh := make(chan *result, len(ticketIds))

	for _, ticketId := range ticketIds {
		/*remaining := t.GetTicketRemaining(stateDB, ticketId)
		ticketsRemaining[ticketId] = remaining*/
		go func(txHash common.Hash) {
			res := new(result)
			res.id = txHash
			res.count = t.GetTicketRemaining(stateDB, txHash)
			resCh <- res
			wg.Done()
		}(ticketId)
	}
	wg.Wait()
	close(resCh)
	for res := range resCh {
		ticketsRemaining[res.id] = res.count
	}
	return ticketsRemaining
}

func (t *TicketPool) GetCandidateTicketIds(stateDB vm.StateDB, nodeId discover.NodeID) []common.Hash {
	log.Debug("Call GetCandidaieTicketIds", "statedb addr", fmt.Sprintf("%p", stateDB))
	return stateDB.GetPPOSCache().GetCandidateTxHashs(nodeId)
}

func (t *TicketPool) GetCandidatesTicketIds(stateDB vm.StateDB, nodeIds []discover.NodeID) map[discover.NodeID][]common.Hash {
	log.Debug("Call GetCandidateArrTicketIds", "statedb addr", fmt.Sprintf("%p", stateDB))
	result := make(map[discover.NodeID][]common.Hash)
	if nodeIds != nil {
		for _, nodeId := range nodeIds {
			ticketIds := t.GetCandidateTicketIds(stateDB, nodeId)
			if nil == ticketIds {
				continue
			}
			result[nodeId] = ticketIds
		}
	}
	return result
}

func (t *TicketPool) GetCandidateTicketCount(stateDB vm.StateDB, nodeId discover.NodeID) uint32 {
	return stateDB.GetPPOSCache().GetCandidateTicketCount(nodeId)
}

func (t *TicketPool) GetCandidatesTicketCount(stateDB vm.StateDB, nodeIds []discover.NodeID) map[discover.NodeID]uint32 {
	log.Debug("Call GetCandidatesTicketCount", "statedb addr", fmt.Sprintf("%p", stateDB))
	result := make(map[discover.NodeID]uint32)
	if nil != nodeIds {
		for _, nodeId := range nodeIds {
			result[nodeId] = stateDB.GetPPOSCache().GetCandidateTicketCount(nodeId)
		}
	}
	return result
}

func (t *TicketPool) setCandidateEpoch(stateDB vm.StateDB, nodeId discover.NodeID, age uint64) {
	stateDB.GetPPOSCache().SetCandidateTicketAge(nodeId, age)
}

func (t *TicketPool) GetCandidateEpoch(stateDB vm.StateDB, nodeId discover.NodeID) uint64 {
	log.Debug("Call GetCandidateEpoch", "statedb addr", fmt.Sprintf("%p", stateDB))
	return stateDB.GetPPOSCache().GetCandidateTicketAge(nodeId)
}

func (t *TicketPool) GetTicketPrice(stateDB vm.StateDB) *big.Int {
	return t.TicketPrice
}

// Save the hash value of the current state of the ticket pool
func (t *TicketPool) CommitHash(stateDB vm.StateDB, blockNumber *big.Int, blockHash common.Hash) error {
	//hash := common.Hash{}
	if hash, err := stateDB.GetPPOSCache().CalculateHash(blockNumber, blockHash); nil != err {
		return err
	}else {
		setTicketPoolState(stateDB, addCommonPrefix(TicketPoolHashKey), hash.Bytes())
		return nil
	}
}

//func GetTicketPtr() *TicketPool {
//	return ticketPool
//}

func checkBalance(stateDB vm.StateDB, addr common.Address, amount *big.Int) bool {
	if stateDB.GetBalance(addr).Cmp(amount) >= 0 {
		return true
	}
	return false
}

func transfer(stateDB vm.StateDB, from common.Address, to common.Address, amount *big.Int) error {
	if !checkBalance(stateDB, from, amount) {
		log.Error("TicketPool not sufficient funds", "from", from.Hex(), "to", to.Hex(), "money", amount.Uint64())
		return TicketPoolBalanceErr
	}
	stateDB.SubBalance(from, amount)
	stateDB.AddBalance(to, amount)
	return nil
}

func getTicketPoolState(stateDB vm.StateDB, key []byte, result interface{}) error {
	return getState(common.TicketPoolAddr, stateDB, key, result)
}

func getState(addr common.Address, stateDB vm.StateDB, key []byte, result interface{}) error {
	if val := stateDB.GetState(addr, key); len(val) > 0 {
		if err := rlp.DecodeBytes(val, result); nil != err {
			log.Error("Decode Data error", "key", string(key), "err", err)
			return err
		}
	}
	return nil
}

func setTicketPoolState(stateDB vm.StateDB, key []byte, val []byte) {
	stateDB.SetState(common.TicketPoolAddr, key, val)
}

func addCommonPrefix(key []byte) []byte {
	return append(common.TicketPoolAddr.Bytes(), key...)
}

func findFirstMatch(list []float64, key float64) int {
	left := 0
	right := len(list) - 1
	for left <= right {
		mid := (left + right) / 2
		if list[mid] >= key {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	// If no match is found, the last subscript is returned by default.
	if left >= len(list) {
		return len(list) - 1
	}
	return left
}
