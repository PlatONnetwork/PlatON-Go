package vm

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"strconv"
)

var (
	ErrTicketPrice     = errors.New("Ticket Price is illegal")
	ErrIllegalDeposit  = errors.New("Deposit balance not match or too low")
	ErrTicketPoolEmpty = errors.New("Ticket Pool is null")
)

const (
	VoteTicketEvent = "VoteTicketEvent"
)

type ticketPoolContext interface {
	VoteTicket(stateDB StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error)
	GetTicket(stateDB StateDB, ticketId common.Hash) (*types.Ticket, error)
	GetTicketList(stateDB StateDB, ticketIds []common.Hash) ([]*types.Ticket, error)
	GetCandidateTicketIds(stateDB StateDB, nodeId discover.NodeID) ([]common.Hash, error)
	GetCandidatesTicketIds(stateDB StateDB, nodeIds []discover.NodeID) (map[discover.NodeID][]common.Hash, error)
	GetCandidatesTicketCount(stateDB StateDB, nodeIds []discover.NodeID) (map[discover.NodeID]int, error)
	GetCandidateEpoch(stateDB StateDB, nodeId discover.NodeID) (uint64, error)
	GetPoolNumber(stateDB StateDB) (uint64, error)
	GetTicketPrice(stateDB StateDB) (*big.Int, error)
}

type TicketContract struct {
	Contract *Contract
	Evm      *EVM
}

func (t *TicketContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (t *TicketContract) Run(input []byte) ([]byte, error) {
	if nil == t.Evm.TicketPoolContext {
		log.Error("Failed to Run==> ", "ErrTicketPoolEmpty: ", ErrTicketPoolEmpty.Error())
		return nil, ErrTicketPoolEmpty
	}
	var command = map[string]interface{}{
		"VoteTicket":                   t.VoteTicket,
		"GetTicketDetail":              t.GetTicketDetail,
		"GetBatchTicketDetail":         t.GetBatchTicketDetail,
		"GetCandidateTicketIds":        t.GetCandidateTicketIds,
		"GetBatchCandidateTicketIds":   t.GetBatchCandidateTicketIds,
		"GetBatchCandidateTicketCount": t.GetBatchCandidateTicketCount,
		"GetCandidateEpoch":            t.GetCandidateEpoch,
		"GetPoolRemainder":             t.GetPoolRemainder,
		"GetTicketPrice":               t.GetTicketPrice,
	}
	return execute(input, command)
}

// VoteTicket let a account buy tickets and vote to the chosen candidate.
func (t *TicketContract) VoteTicket(count uint64, price *big.Int, nodeId discover.NodeID) ([]byte, error) {
	value := t.Contract.value
	txHash := t.Evm.StateDB.TxHash()
	txIdx := t.Evm.StateDB.TxIdx()
	blockNumber := t.Evm.Context.BlockNumber
	from := t.Contract.caller.Address()
	log.Info("Input to VoteTicket==>", " nodeId: ", nodeId.String(), " owner: ", from.Hex(), " txhash: ", txHash.Hex(),
		" txIdx: ", txIdx, " blockNumber: ", blockNumber, " value: ", value, " count: ", count, " price: ", price)
	if ticketPrice, err := t.Evm.TicketPoolContext.GetTicketPrice(t.Evm.StateDB); nil == err {
		if price.Cmp(ticketPrice) < 0 {
			log.Error("Failed to VoteTicket==> ", "ErrTicketPrice: ", ErrTicketPrice.Error())
			return nil, ErrTicketPrice
		}
	} else {
		log.Error("Failed to VoteTicket==> ", "GetTicketPrice return err: ", err.Error())
		return nil, err
	}
	totalPrice := new(big.Int).Mul(new(big.Int).SetUint64(count), price)
	if value.Cmp(totalPrice) < 0 || totalPrice.Cmp(big.NewInt(0)) != 1 {
		log.Error("Failed to VoteTicket==> ", "ErrIllegalDeposit: ", ErrIllegalDeposit.Error())
		return nil, ErrIllegalDeposit
	}
	can, err := t.Evm.CandidatePoolContext.GetCandidate(t.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to VoteTicket==> ", "GetCandidate return err: ", err.Error())
		return nil, err
	}
	if nil == can {
		log.Error("Failed to VoteTicket==> ", "ErrCandidateNotExist: ", ErrCandidateNotExist.Error())
		return nil, ErrCandidateNotExist
	}
	ticketIds, err := t.Evm.TicketPoolContext.VoteTicket(t.Evm.StateDB, from, count, price, nodeId, blockNumber)
	data := len(ticketIds)
	if 0 == data {
		log.Error("Failed to VoteTicket==> ", "VoteTicket return err(0 == len(ticketIds)): ", err.Error())
		return nil, err
	}
	// return the extra money
	if uint64(data) < count {
		failNum := count - uint64(data)
		backBalance := new(big.Int).Mul(new(big.Int).SetUint64(failNum), price)
		t.Evm.StateDB.AddBalance(from, backBalance)
		t.Evm.StateDB.SubBalance(common.TicketPoolAddr, backBalance)
	}
	sdata := DecodeResultStr(strconv.Itoa(data))
	log.Info("Result of VoteTicket==> ", "len(successTicketIds): ", strconv.Itoa(data), " []byte: ", sdata)
	if nil != err {
		log.Warn("Failed to VoteTicket==> ", "VoteTicket return err: ", err.Error())
		r := ResultCommon{true, strconv.Itoa(data), err.Error()}
		event, _ := json.Marshal(r)
		t.addLog(VoteTicketEvent, string(event))
		return sdata, nil
	}
	r := ResultCommon{true, strconv.Itoa(data), "success"}
	event, _ := json.Marshal(r)
	t.addLog(VoteTicketEvent, string(event))
	return sdata, nil
}

// GetTicketDetail returns the ticket info.
func (t *TicketContract) GetTicketDetail(ticketId common.Hash) ([]byte, error) {
	log.Info("Input to GetTicketDetail==> ", "ticketId: ", ticketId.Hex())
	ticket, err := t.Evm.TicketPoolContext.GetTicket(t.Evm.StateDB, ticketId)
	if nil != err {
		log.Error("Failed to GetTicketDetail==> ", "GetTicket return err: ", err.Error())
		ticket := types.Ticket{}
		data, _ := json.Marshal(ticket)
		sdata := DecodeResultStr(string(data))
		return sdata, err
	}
	data, _ := json.Marshal(ticket)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetTicketDetail==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetBatchTicketDetail returns the batch of ticket info.
func (t *TicketContract) GetBatchTicketDetail(ticketIds []common.Hash) ([]byte, error) {
	input, _ := json.Marshal(ticketIds)
	log.Info("Input to GetBatchTicketDetail==>", "length: ", len(ticketIds), "ticketIds: ", string(input))
	tickets, _ := t.Evm.TicketPoolContext.GetTicketList(t.Evm.StateDB, ticketIds)
	if 0 == len(tickets) {
		log.Warn("Failed to GetBatchTicketDetail==> The query does not exist")
		tickets := make([]types.Ticket, 0)
		data, _ := json.Marshal(tickets)
		sdata := DecodeResultStr(string(data))
		return sdata, nil
	}
	data, _ := json.Marshal(tickets)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetBatchTicketDetail==> ", "len(tickets): ", len(tickets), "json: ", string(data))
	return sdata, nil
}

// GetCandidateTicketIds returns the list of ticketId for the candidate.
func (t *TicketContract) GetCandidateTicketIds(nodeId discover.NodeID) ([]byte, error) {
	log.Info("Input to GetCandidateTicketIds==> ", " nodeId: ", nodeId.String())
	candidateTicketIds, err := t.Evm.TicketPoolContext.GetCandidateTicketIds(t.Evm.StateDB, nodeId)
	if nil != err {
		log.Warn("Failed to GetCandidateTicketIds==> ", "GetCandidateTicketIds return err: ", err.Error())
		candidateTicketIds := make([]common.Hash, 0)
		data, _ := json.Marshal(candidateTicketIds)
		sdata := DecodeResultStr(string(data))
		return sdata, err
	}
	data, _ := json.Marshal(candidateTicketIds)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetCandidateTicketIds==> ", "len(candidateTicketIds): ", len(candidateTicketIds), "json: ", string(data))
	return sdata, nil
}

// GetBatchCandidateTicketIds returns the batch of candidate's ticketIds.
func (t *TicketContract) GetBatchCandidateTicketIds(nodeIds []discover.NodeID) ([]byte, error) {
	input, _ := json.Marshal(nodeIds)
	log.Info("Input to GetBatchCandidateTicketIds==> ", "length: ", len(nodeIds), "nodeIds: ", string(input))
	candidatesTicketIds, _ := t.Evm.TicketPoolContext.GetCandidatesTicketIds(t.Evm.StateDB, nodeIds)
	if 0 == len(candidatesTicketIds) {
		log.Warn("Failed to GetBatchCandidateTicketIds==> The query does not exist")
		candidatesTicketIds := make(map[discover.NodeID][]common.Hash, 0)
		data, _ := json.Marshal(candidatesTicketIds)
		sdata := DecodeResultStr(string(data))
		return sdata, nil
	}
	data, _ := json.Marshal(candidatesTicketIds)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetBatchCandidateTicketIds==> ", "len(candidatesTicketIds): ", len(candidatesTicketIds), "json: ", string(data))
	return sdata, nil
}

// GetBatchCandidateTicketCount returns the number of candidate's ticket.
func (t *TicketContract) GetBatchCandidateTicketCount(nodeIds []discover.NodeID) ([]byte, error) {
	input, _ := json.Marshal(nodeIds)
	log.Info("Input to GetBatchCandidateTicketCount==> ", "length: ", len(nodeIds), "nodeIds: ", string(input))
	candidatesTicketCount, _ := t.Evm.TicketPoolContext.GetCandidatesTicketCount(t.Evm.StateDB, nodeIds)
	if 0 == len(candidatesTicketCount) {
		log.Warn("Failed to GetBatchCandidateTicketCount==> The query does not exist")
		candidatesTicketCount := make(map[discover.NodeID]int, 0)
		data, _ := json.Marshal(candidatesTicketCount)
		sdata := DecodeResultStr(string(data))
		return sdata, nil
	}
	data, _ := json.Marshal(candidatesTicketCount)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetBatchCandidateTicketCount==> ", "len(candidatesTicketCount): ", len(candidatesTicketCount), "json: ", string(data))
	return sdata, nil
}

// GetEpoch returns the current ticket age for the candidate.
func (t *TicketContract) GetCandidateEpoch(nodeId discover.NodeID) ([]byte, error) {
	log.Info("Input to GetCandidateEpoch==> ", " nodeId: ", nodeId.String())
	epoch, err := t.Evm.TicketPoolContext.GetCandidateEpoch(t.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to GetCandidateEpoch==> ", "GetCandidateEpoch return err: ", err.Error())
		data, _ := json.Marshal(epoch)
		sdata := DecodeResultStr(string(data))
		return sdata, err
	}
	data, _ := json.Marshal(epoch)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetCandidateEpoch==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetPoolRemainder returns the amount of remaining tickets in the ticket pool.
func (t *TicketContract) GetPoolRemainder() ([]byte, error) {
	remainder, err := t.Evm.TicketPoolContext.GetPoolNumber(t.Evm.StateDB)
	if nil != err {
		log.Error("Failed to GetPoolRemainder==> ", "GetPoolNumber return err: ", err.Error())
		data, _ := json.Marshal(remainder)
		sdata := DecodeResultStr(string(data))
		return sdata, err
	}
	data, _ := json.Marshal(remainder)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetPoolRemainder==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetTicketPrice returns the current ticket price for the ticket pool.
func (t *TicketContract) GetTicketPrice() ([]byte, error) {
	price, err := t.Evm.TicketPoolContext.GetTicketPrice(t.Evm.StateDB)
	if nil != err {
		log.Error("Failed to GetTicketPrice==> ", "GetTicketPrice return err: ", err.Error())
		data, _ := json.Marshal(price)
		sdata := DecodeResultStr(string(data))
		return sdata, err
	}
	data, _ := json.Marshal(price)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetTicketPrice==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// addLog let the result add to event.
func (t *TicketContract) addLog(event, data string) {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); nil != err {
		log.Error("Failed to addlog==> ", "rlp encode fail: ", err.Error())
	}
	t.Evm.StateDB.AddLog(&types.Log{
		Address:     common.TicketPoolAddr,
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data:        buf.Bytes(),
		BlockNumber: t.Evm.Context.BlockNumber.Uint64(),
	})
}
