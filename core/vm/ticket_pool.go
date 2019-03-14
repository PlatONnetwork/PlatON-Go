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
	VoteTicket(stateDB StateDB, owner common.Address, voteNumber uint32, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) (uint32, error)
	GetCandidatesTicketCount(stateDB StateDB, nodeIds []discover.NodeID) map[discover.NodeID]uint32
	GetBatchTicketRemaining(stateDB StateDB, ticketIds []common.Hash) map[common.Hash]uint32
	GetCandidateEpoch(stateDB StateDB, nodeId discover.NodeID) uint64
	GetPoolNumber(stateDB StateDB) uint32
	GetTicketPrice(stateDB StateDB) *big.Int
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
		"VoteTicket":               t.VoteTicket,
		"GetCandidateTicketCount":  t.GetCandidateTicketCount,
		"GetTicketCountByTicketId": t.GetTicketCountByTicketId,
		"GetCandidateEpoch":        t.GetCandidateEpoch,
		"GetPoolRemainder":         t.GetPoolRemainder,
		"GetTicketPrice":           t.GetTicketPrice,
	}
	return execute(input, command)
}

// VoteTicket let a account buy tickets and vote to the chosen candidate.
func (t *TicketContract) VoteTicket(count uint32, price *big.Int, nodeId discover.NodeID) ([]byte, error) {
	value := t.Contract.value
	txHash := t.Evm.StateDB.TxHash()
	txIdx := t.Evm.StateDB.TxIdx()
	blockNumber := t.Evm.Context.BlockNumber
	from := t.Contract.caller.Address()
	log.Info("Input to VoteTicket==>", " nodeId: ", nodeId.String(), " owner: ", from.Hex(), " txhash: ", txHash.Hex(),
		" txIdx: ", txIdx, " blockNumber: ", blockNumber, " value: ", value, " count: ", count, " price: ", price)
	ticketPrice := t.Evm.TicketPoolContext.GetTicketPrice(t.Evm.StateDB)
	if price.Cmp(ticketPrice) < 0 {
		log.Error("Failed to VoteTicket==> ", "ErrTicketPrice: ", ErrTicketPrice.Error())
		return nil, ErrTicketPrice
	}
	totalPrice := new(big.Int).Mul(new(big.Int).SetUint64(uint64(count)), price)
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
	successCount, err := t.Evm.TicketPoolContext.VoteTicket(t.Evm.StateDB, from, count, price, nodeId, blockNumber)
	if 0 == successCount {
		log.Error("Failed to VoteTicket==> ", "VoteTicket return err(0 == len(ticketIds)): ", err.Error())
		return nil, err
	}
	// return the extra money
	if successCount < count {
		failCount := count - successCount
		backBalance := new(big.Int).Mul(new(big.Int).SetUint64(uint64(failCount)), price)
		t.Evm.StateDB.AddBalance(from, backBalance)
		t.Evm.StateDB.SubBalance(common.TicketPoolAddr, backBalance)
	}
	sdata := DecodeResultStr(strconv.FormatUint(uint64(successCount), 10))
	log.Info("Result of VoteTicket==> ", "successCount: ", successCount, " []byte: ", sdata)
	if nil != err {
		log.Warn("Failed to VoteTicket==> ", "VoteTicket return err: ", err.Error())
		r := ResultCommon{true, strconv.FormatUint(uint64(successCount), 10), err.Error()}
		event, _ := json.Marshal(r)
		t.addLog(VoteTicketEvent, string(event))
		return sdata, nil
	}
	r := ResultCommon{true, strconv.FormatUint(uint64(successCount), 10), "success"}
	event, _ := json.Marshal(r)
	t.addLog(VoteTicketEvent, string(event))
	return sdata, nil
}

// GetCandidateTicketCount returns the number of candidate's ticket.
func (t *TicketContract) GetCandidateTicketCount(nodeIds []discover.NodeID) ([]byte, error) {
	input, _ := json.Marshal(nodeIds)
	log.Info("Input to GetCandidateTicketCount==> ", "length: ", len(nodeIds), "nodeIds: ", string(input))
	candidatesTicketCount := t.Evm.TicketPoolContext.GetCandidatesTicketCount(t.Evm.StateDB, nodeIds)
	data, _ := json.Marshal(candidatesTicketCount)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetCandidateTicketCount==> ", "len(candidatesTicketCount): ", len(candidatesTicketCount), "json: ", string(data))
	return sdata, nil
}

// GetTicketCountByTicketId returns the number of transaction's ticket.
func (t *TicketContract) GetTicketCountByTicketId(ticketIds []common.Hash) ([]byte, error) {
	input, _ := json.Marshal(ticketIds)
	log.Info("Input to GetTicketCountByTicketId==> ", "length: ", len(ticketIds), "ticketIds: ", string(input))
	ticketsRemaining := t.Evm.TicketPoolContext.GetBatchTicketRemaining(t.Evm.StateDB, ticketIds)
	data, _ := json.Marshal(ticketsRemaining)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetTicketCountByTicketId==> ", "len(ticketsRemaining): ", len(ticketsRemaining), "json: ", string(data))
	return sdata, nil
}

// GetEpoch returns the current ticket age for the candidate.
func (t *TicketContract) GetCandidateEpoch(nodeId discover.NodeID) ([]byte, error) {
	log.Info("Input to GetCandidateEpoch==> ", " nodeId: ", nodeId.String())
	epoch := t.Evm.TicketPoolContext.GetCandidateEpoch(t.Evm.StateDB, nodeId)
	data, _ := json.Marshal(epoch)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetCandidateEpoch==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetPoolRemainder returns the amount of remaining tickets in the ticket pool.
func (t *TicketContract) GetPoolRemainder() ([]byte, error) {
	remainder := t.Evm.TicketPoolContext.GetPoolNumber(t.Evm.StateDB)
	data, _ := json.Marshal(remainder)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetPoolRemainder==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetTicketPrice returns the current ticket price for the ticket pool.
func (t *TicketContract) GetTicketPrice() ([]byte, error) {
	price := t.Evm.TicketPoolContext.GetTicketPrice(t.Evm.StateDB)
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
