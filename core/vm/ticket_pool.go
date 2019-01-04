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
)

var (
	ErrIllegalDeposit    = errors.New("Deposit balance not match")
	ErrCandidateNotExist = errors.New("Voted candidate not exist")
	ErrTicketPoolEmpty   = errors.New("Ticket Pool is null")
)

const (
	VoteTicketEvent = "VoteTicketEvent"
)

type ticketPool interface {
	VoteTicket(stateDB StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error)
	GetTicket(stateDB StateDB, ticketId common.Hash) (*types.Ticket, error)
	GetTicketList(stateDB StateDB, ticketIds []common.Hash) ([]*types.Ticket, error)
	GetCandidateTicketIds(stateDB StateDB, nodeId discover.NodeID) ([]common.Hash, error)
	GetCandidatesTicketIds(stateDB StateDB, nodeIds []discover.NodeID) (map[discover.NodeID][]common.Hash, error)
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
	if nil == t.Evm.TicketPool {
		log.Error("Failed to Run==> ", "ErrTicketPoolEmpty: ", ErrTicketPoolEmpty.Error())
		return nil, ErrTicketPoolEmpty
	}
	var command = map[string]interface{}{
		"VoteTicket":                 t.VoteTicket,
		"GetTicketDetail":            t.GetTicketDetail,
		"GetBatchTicketDetail":       t.GetBatchTicketDetail,
		"GetCandidateTicketIds":      t.GetCandidateTicketIds,
		"GetBatchCandidateTicketIds": t.GetBatchCandidateTicketIds,
		"GetCandidateEpoch":          t.GetCandidateEpoch,
		"GetPoolRemainder":           t.GetPoolRemainder,
		"GetTicketPrice":             t.GetTicketPrice,
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
	can, err := t.Evm.CandidatePool.GetCandidate(t.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to VoteTicket==> ", "GetCandidate return err: ", err.Error())
		return nil, err
	}
	if nil == can {
		log.Error("Failed to VoteTicket==> ", "ErrCandidateNotExist: ", ErrCandidateNotExist.Error())
		return nil, ErrCandidateNotExist
	}
	totalPrice := new(big.Int).Mul(new(big.Int).SetUint64(count), price)
	if totalPrice.Cmp(value) != 0 || totalPrice.Cmp(big.NewInt(0)) != 1 {
		log.Error("Failed to VoteTicket==> ", "ErrIllegalDeposit: ", ErrIllegalDeposit.Error())
		return nil, ErrIllegalDeposit
	}
	ticketIds, err := t.Evm.TicketPool.VoteTicket(t.Evm.StateDB, from, count, price, nodeId, blockNumber)
	if nil == ticketIds {
		log.Error("Failed to VoteTicket==> ", "VoteTicket return err(nil == ticketIds): ", err.Error())
		return nil, err
	}
	data := len(ticketIds)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of VoteTicket==> ", "len(successTicketIds): ", data, " []byte: ", sdata)
	if nil != err {
		log.Error("Failed to VoteTicket==> ", "VoteTicket return err: ", err.Error())
		r := ResultCommon{true, string(data), err.Error()}
		event, _ := json.Marshal(r)
		t.addLog(VoteTicketEvent, string(event))
		return sdata, nil
	}
	r := ResultCommon{true, string(data), "success"}
	event, _ := json.Marshal(r)
	t.addLog(VoteTicketEvent, string(event))
	return sdata, nil
}

// GetTicketDetail returns the ticket info.
func (t *TicketContract) GetTicketDetail(ticketId common.Hash) ([]byte, error) {
	log.Info("Input to GetTicketDetail==> ", "ticketId: ", ticketId.Hex())
	ticket, err := t.Evm.TicketPool.GetTicket(t.Evm.StateDB, ticketId)
	if nil != err {
		log.Error("Failed to GetTicketDetail==> ", "GetTicket return err: ", err.Error())
		return nil, err
	}
	if (common.Hash{}) == ticket.TicketId {
		log.Error("Failed to GetTicketDetail==> ", "The query does not exist")
		return nil, nil
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
	tickets, err := t.Evm.TicketPool.GetTicketList(t.Evm.StateDB, ticketIds)
	if nil != err {
		if 0 == len(tickets) {
			log.Error("Failed to GetBatchTicketDetail==> ", "GetTicketList return err: ", err.Error())
			return nil, err
		}
		data, _ := json.Marshal(tickets)
		sdata := DecodeResultStr(string(data))
		log.Info("Result of GetBatchTicketDetail==> ", "json: ", string(data), "[]byte: ", sdata)
		return sdata, err
	}
	data, _ := json.Marshal(tickets)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetBatchTicketDetail==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetCandidateTicketIds returns the list of ticketId for the candidate.
func (t *TicketContract) GetCandidateTicketIds(nodeId discover.NodeID) ([]byte, error) {
	log.Info("Input to GetCandidateTicketIds==> ", " nodeId: ", nodeId.String())
	candidateTicketIds, err := t.Evm.TicketPool.GetCandidateTicketIds(t.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to GetCandidateTicketIds==> ", "GetCandidateTicketIds return err: ", err.Error())
		return nil, err
	}
	data, _ := json.Marshal(candidateTicketIds)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetCandidateTicketIds==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetBatchCandidateTicketIds returns the batch of candidate's ticketIds.
func (t *TicketContract) GetBatchCandidateTicketIds(nodeIds []discover.NodeID) ([]byte, error) {
	input, _ := json.Marshal(nodeIds)
	log.Info("Input to GetBatchCandidateTicketIds==> ", "length: ", len(nodeIds), "nodeIds: ", string(input))
	candidatesTicketIds, err := t.Evm.TicketPool.GetCandidatesTicketIds(t.Evm.StateDB, nodeIds)
	if nil != err {
		if 0 == len(candidatesTicketIds) {
			log.Error("Failed to GetBatchCandidateTicketIds==> ", "GetCandidatesTicketIds return err: ", err.Error())
			return nil, err
		}
		data, _ := json.Marshal(candidatesTicketIds)
		sdata := DecodeResultStr(string(data))
		log.Info("Result of GetBatchCandidateTicketIds==> ", "json: ", string(data), "[]byte: ", sdata)
		return sdata, err
	}
	data, _ := json.Marshal(candidatesTicketIds)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetBatchCandidateTicketIds==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetEpoch returns the current ticket age for the candidate.
func (t *TicketContract) GetCandidateEpoch(nodeId discover.NodeID) ([]byte, error) {
	log.Info("Input to GetCandidateEpoch==> ", " nodeId: ", nodeId.String())
	epoch, err := t.Evm.TicketPool.GetCandidateEpoch(t.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to GetCandidateEpoch==> ", "GetCandidateEpoch return err: ", err.Error())
		return nil, err
	}
	data, _ := json.Marshal(epoch)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetCandidateEpoch==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetPoolRemainder returns the amount of remaining tickets in the ticket pool.
func (t *TicketContract) GetPoolRemainder() ([]byte, error) {
	remainder, err := t.Evm.TicketPool.GetPoolNumber(t.Evm.StateDB)
	if nil != err {
		log.Error("Failed to GetPoolRemainder==> ", "GetPoolNumber return err: ", err.Error())
		return nil, err
	}
	data, _ := json.Marshal(remainder)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetPoolRemainder==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetTicketPrice returns the current ticket price for the ticket pool.
func (t *TicketContract) GetTicketPrice() ([]byte, error) {
	price, err := t.Evm.TicketPool.GetTicketPrice(t.Evm.StateDB)
	if nil != err {
		log.Error("Failed to GetTicketPrice==> ", "GetTicketPrice return err: ", err.Error())
		return nil, err
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
