package vm

import (
	"Platon-go/common"
	"Platon-go/core/types"
	"Platon-go/crypto"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"Platon-go/rlp"
	"bytes"
	"encoding/json"
	"errors"
	"math/big"
)

// error def
var (
	ErrIllegalDeposit = errors.New("Deposit balance not match")
	ErrCandidateNotExist = errors.New("Voted candidate not exist")
	ErrTicketPoolEmpty = errors.New("Ticket Pool is null")
)

const (
	VoteTicketEvent = "VoteTicketEvent"
)

type ticketPool interface {
	VoteTicket(stateDB StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error)
	GetTicket(stateDB StateDB, ticketId common.Hash) (*types.Ticket, error)
	GetCandidateTicketIds(stateDB StateDB, nodeId discover.NodeID) ([]common.Hash, error)
	GetPoolNumber(stateDB StateDB) (uint64, error)
}

type ticketContract struct {
	contract *Contract
	evm *EVM
}

func (t *ticketContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (t *ticketContract) Run(input []byte) ([]byte, error) {
	if nil == t.evm.TicketPool {
		log.Error("Run==> ", "ErrTicketPoolEmpty: ", ErrTicketPoolEmpty.Error())
		return nil, ErrTicketPoolEmpty
	}
	var command = map[string] interface{}{
		"VoteTicket" : t.VoteTicket,
		"GetTicketDetail" : t.GetTicketDetail,
		"GetCandidateTicketIds" : t.GetCandidateTicketIds,
		"GetEpoch" : t.GetEpoch,
		"GetPoolRemainder" : t.GetPoolRemainder,
		"GetTicketPrice": t.GetTicketPrice,
	}
	return execute(input, command)
}

// VoteTicket let a account buy tickets and vote to the chosen candidate.
func (t *ticketContract) VoteTicket(count uint64, price *big.Int, nodeId discover.NodeID) ([]byte, error) {
	// input params
	deposit := t.contract.value
	txHash := t.evm.StateDB.TxHash()
	txIdx := t.evm.StateDB.TxIdx()
	height := t.evm.Context.BlockNumber
	from := t.contract.caller.Address()
	log.Info("VoteTicket==>", " nodeId: ", nodeId.String(), " owner: ", from.Hex(), " txhash: ", txHash.Hex(),
		" txIdx: ", txIdx, " height: ", height, " deposit: ", deposit, " count: ", count, " price: ", price)
	can, err := t.evm.CandidatePool.GetCandidate(t.evm.StateDB, nodeId)
	if nil != err {
		log.Error("VoteTicket==> ","GetCandidate occured error", err.Error())
		return nil, err
	}
	if nil == can {
		return nil, ErrCandidateNotExist
	}
	totalPrice := new(big.Int).Mul(new(big.Int).SetUint64(count), price)
	if totalPrice != deposit || deposit.Cmp(big.NewInt(0)) != 1 || totalPrice.Cmp(big.NewInt(0)) != 1 {
		return nil, ErrIllegalDeposit
	}

	// return ([]common.hash, error)
	arr, err := t.evm.TicketPool.VoteTicket(t.evm.StateDB, from, count, deposit, nodeId, height)
	if nil != err {
		log.Error("VoteTicket==> ","voteTicket occured error", err.Error())
		return nil, err
	}

	data, _ := json.Marshal(arr)
	t.addLog(VoteTicketEvent, string(data))
	sdata := DecodeResultStr(string(data))
	log.Info("VoteTicket==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetTicketDetail returns the ticket info.
func (t *ticketContract) GetTicketDetail(ticketId common.Hash) ([]byte, error) {
	// input params
	log.Info("GetTicketDetail==>", " nodeId: ", ticketId.Hex())
	ticket, err := t.evm.TicketPool.GetTicket(t.evm.StateDB, ticketId)
	if nil != err {
		log.Error("GetTicketDetail==> ","GetTicketDetail() occured error: ", err.Error())
		return nil, err
	}
	if nil == ticket {
		log.Error("GetTicketDetail==> The ticket for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(ticket)
	sdata := DecodeResultStr(string(data))
	log.Info("GetTicketDetail==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetCandidateTicketIds returns the list of ticketId for the candidate.
func (t *ticketContract) GetCandidateTicketIds(nodeId discover.NodeID, blockNumber *big.Int) ([]byte, error) {
	// input params
	log.Info("GetCandidateTicketIds==>", " nodeId: ", nodeId.String(), " blockNumber: ", blockNumber)
	return nil, nil
}

// GetEpoch returns the current ticket age for the candidate.
func (t *ticketContract) GetEpoch(nodeId discover.NodeID) ([]byte, error) {
	// input params
	log.Info("GetEpoch==>", " nodeId: ", nodeId.String())
	return nil, nil
}

// GetPoolRemainder returns the amount of remaining tikcets in the ticket pool.
func (t *ticketContract) GetPoolRemainder() ([]byte, error) {
	return nil, nil
}

// GetTicketPrice returns the current ticket price for the ticket pool.
func (t *ticketContract) GetTicketPrice() ([]byte, error) {
	return nil, nil
}

// transaction add event
func (t *ticketContract) addLog(event, data string) {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); nil != err {
		log.Error("addlog==> ","rlp encode fail: ", err.Error())
	}
	t.evm.StateDB.AddLog(&types.Log{
		Address:common.TicketPoolAddr,
		Topics: []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data: buf.Bytes(),
		BlockNumber: t.evm.Context.BlockNumber.Uint64(),
	})
}