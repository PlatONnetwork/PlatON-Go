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
	ErrOwnerNotOnly          = errors.New("Node ID cannot bind multiple owners")
	ErrPermissionDenied      = errors.New("Transaction from address permission denied")
	ErrFeeIllegal            = errors.New("The fee is illegal")
	ErrDepositEmpty          = errors.New("Deposit balance not zero")
	ErrWithdrawEmpty         = errors.New("No withdrawal amount")
	ErrCandidatePoolEmpty    = errors.New("Candidate Pool is null")
	ErrCandidateNotExist     = errors.New("The candidate is not exist")
	ErrCandidateAlreadyExist = errors.New("The candidate is already exist")
)

const (
	CandidateDepositEvent       = "CandidateDepositEvent"
	CandidateApplyWithdrawEvent = "CandidateApplyWithdrawEvent"
	CandidateWithdrawEvent      = "CandidateWithdrawEvent"
	SetCandidateExtraEvent      = "SetCandidateExtraEvent"
)

type candidatePoolContext interface {
	SetCandidate(state StateDB, nodeId discover.NodeID, can *types.Candidate) error
	GetCandidate(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) *types.Candidate
	GetCandidateArr(state StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) types.CandidateQueue
	WithdrawCandidate(state StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error
	GetChosens(state StateDB, flag int, blockNumber *big.Int) types.KindCanQueue
	GetChairpersons(state StateDB, blockNumber *big.Int) types.CandidateQueue
	GetDefeat(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) types.RefundQueue
	IsDefeat(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) bool
	IsChosens(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) bool
	RefundBalance(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) error
	GetOwner(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) common.Address
	SetCandidateExtra(state StateDB, nodeId discover.NodeID, extra string) error
	GetRefundInterval(blockNumber *big.Int) uint32
	MaxCount() uint32
	MaxChair() uint32
}

type CandidateContract struct {
	Contract *Contract
	Evm      *EVM
}

func (c *CandidateContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (c *CandidateContract) Run(input []byte) ([]byte, error) {
	if nil == c.Evm.CandidatePoolContext {
		log.Error("Failed to CandidateContract Run", "ErrCandidatePoolEmpty: ", ErrCandidatePoolEmpty.Error())
		return nil, ErrCandidatePoolEmpty
	}
	var command = map[string]interface{}{
		"CandidateDeposit":          c.CandidateDeposit,
		"CandidateApplyWithdraw":    c.CandidateApplyWithdraw,
		"CandidateWithdraw":         c.CandidateWithdraw,
		"SetCandidateExtra":         c.SetCandidateExtra,
		"GetCandidateWithdrawInfos": c.GetCandidateWithdrawInfos,
		"GetCandidateDetails":       c.GetCandidateDetails,
		"GetCandidateList":          c.GetCandidateList,
		"GetVerifiersList":          c.GetVerifiersList,
	}
	return execute(input, command)
}

// Candidate Application && Increase Quality Deposit
func (c *CandidateContract) CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint32, host, port, extra string) ([]byte, error) {
	deposit := c.Contract.value
	txHash := c.Evm.StateDB.TxHash()
	txIdx := c.Evm.StateDB.TxIdx()
	height := c.Evm.Context.BlockNumber
	//from := c.Contract.caller.Address()
	log.Info("Input to CandidateDeposit", "blockNumber", height.String(), "nodeId: ", nodeId.String(), " owner: ", owner.Hex(), " deposit: ", deposit,
		"  fee: ", fee, " txhash: ", txHash.Hex(), " txIdx: ", txIdx, " height: ", height, " host: ", host, " port: ", port, " extra: ", extra)
	if fee > 10000 {
		log.Error("Failed to CandidateDeposit", "blockNumber", height.String(), "ErrFeeIllegal: ", ErrFeeIllegal.Error())
		return nil, ErrFeeIllegal
	}
	if deposit.Cmp(big.NewInt(0)) < 1 {
		log.Error("Failed to CandidateDeposit", "blockNumber", height.String(), "ErrDepositEmpty: ", ErrDepositEmpty.Error())
		return nil, ErrDepositEmpty
	}
	addr := c.Evm.CandidatePoolContext.GetOwner(c.Evm.StateDB, nodeId, height)
	if common.ZeroAddr != addr {
		if ok := bytes.Equal(addr.Bytes(), owner.Bytes()); !ok {
			log.Error("Failed to CandidateDeposit==> ", "blockNumber", height.String(), "old owner", addr.Hex(), "new owner", owner, "ErrOwnerNotOnly: ", ErrOwnerNotOnly.Error())
			return nil, ErrOwnerNotOnly
		}
	}
	//var alldeposit *big.Int
	var txhash common.Hash
	var towner common.Address
	can := c.Evm.CandidatePoolContext.GetCandidate(c.Evm.StateDB, nodeId, height)
	if nil != can {
		log.Error("Failed to CandidateDeposit, the candidate is already exist", "blockNumber", height.String(), "nodeId", nodeId.String())
		return nil, ErrCandidateAlreadyExist
	}
	//if nil != can {
	//	alldeposit = new(big.Int).Add(can.Deposit, deposit)
	//	txhash = can.TxHash
	//	towner = can.TOwner
	//	log.Info("CandidateDeposit==> ", "alldeposit: ", alldeposit, " can.Deposit: ", can.Deposit, " deposit: ", deposit)
	//} else {
	//	alldeposit = deposit
	//}
	canDeposit := types.Candidate{
		//alldeposit,
		deposit,
		height,
		txIdx,
		nodeId,
		host,
		port,
		owner,
		extra,
		fee,
		txhash,
		towner,
	}
	log.Info("CandidateDeposit", "blockNumber", height.String(), "canDeposit: ", canDeposit)
	if err := c.Evm.CandidatePoolContext.SetCandidate(c.Evm.StateDB, nodeId, &canDeposit); nil != err {
		log.Error("Failed to CandidateDeposit", "blockNumber", height.String(), "SetCandidate return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "", "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateDepositEvent, string(event))
	log.Info("Result of CandidateDeposit", "blockNumber", height.String(), "json: ", string(event))
	return nil, nil
}

// Apply for a refund of the deposit
func (c *CandidateContract) CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	from := c.Contract.caller.Address()
	height := c.Evm.Context.BlockNumber
	log.Info("Input to CandidateApplyWithdraw on WithdrawCandidate", "blockNumber", height.String(), "nodeId: ", nodeId.String(), " from: ", from.Hex(), " txHash: ", txHash.Hex(), " withdraw: ", withdraw, " height: ", height)
	can := c.Evm.CandidatePoolContext.GetCandidate(c.Evm.StateDB, nodeId, height)

	if nil == can {
		log.Error("Failed to CandidateApplyWithdraw on WithdrawCandidate", "blockNumber", height.String(), "ErrCandidateNotExist: ", ErrCandidateNotExist.Error())
		return nil, ErrCandidateNotExist
	}
	if can.Deposit.Cmp(big.NewInt(0)) < 1 {
		log.Error("Failed to CandidateApplyWithdraw on WithdrawCandidate", "blockNumber", height.String(), "ErrWithdrawEmpty: ", ErrWithdrawEmpty.Error())
		return nil, ErrWithdrawEmpty
	}
	if ok := bytes.Equal(can.Owner.Bytes(), from.Bytes()); !ok {
		log.Error("Failed to CandidateApplyWithdraw on WithdrawCandidate", "blockNumber", height.String(), "ErrPermissionDenied: ", ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if withdraw.Cmp(can.Deposit) > 0 {
		withdraw = can.Deposit
	}
	if err := c.Evm.CandidatePoolContext.WithdrawCandidate(c.Evm.StateDB, nodeId, withdraw, height); nil != err {
		log.Error("Failed to CandidateApplyWithdraw on WithdrawCandidate", "blockNumber", height.String(), "WithdrawCandidate return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "", "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateApplyWithdrawEvent, string(event))
	log.Info("Result of CandidateApplyWithdraw on WithdrawCandidate", "blockNumber", height.String(), "json: ", string(event))
	return nil, nil
}

// Deposit withdrawal
func (c *CandidateContract) CandidateWithdraw(nodeId discover.NodeID) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	height := c.Evm.Context.BlockNumber
	log.Info("Input to CandidateWithdraw to RefundBalance", "nodeId: ", nodeId.String(), " height: ", height, " txHash: ", txHash.Hex())
	if err := c.Evm.CandidatePoolContext.RefundBalance(c.Evm.StateDB, nodeId, height); nil != err {
		log.Error("Failed to CandidateWithdraw to RefundBalance", "blockNumber", height.String(), "RefundBalance return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "", "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateWithdrawEvent, string(event))
	log.Info("Result of CandidateWithdraw to RefundBalance", "blockNumber", height.String(), "json: ", string(event))
	return nil, nil
}

// Set up additional information
func (c *CandidateContract) SetCandidateExtra(nodeId discover.NodeID, extra string) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	from := c.Contract.caller.Address()
	height := c.Evm.Context.BlockNumber
	log.Info("Input to SetCandidateExtra", "blockNumber", height.String(), "nodeId: ", nodeId.String(), " extra: ", extra, " from: ", from.Hex(), " txHash: ", txHash.Hex())
	owner := c.Evm.CandidatePoolContext.GetOwner(c.Evm.StateDB, nodeId, height)
	if ok := bytes.Equal(owner.Bytes(), from.Bytes()); !ok {
		log.Error("Failed to SetCandidateExtra", "blockNumber", height.String(), "ErrPermissionDenied: ", ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if err := c.Evm.CandidatePoolContext.SetCandidateExtra(c.Evm.StateDB, nodeId, extra); nil != err {
		log.Error("Failed to SetCandidateExtra", "blockNumber", height.String(), "SetCandidateExtra return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "", "success"}
	event, _ := json.Marshal(r)
	c.addLog(SetCandidateExtraEvent, string(event))
	log.Info("Result of SetCandidateExtra", "blockNumber", height.String(), "json: ", string(event))
	return nil, nil
}

// Get the refund history you have applied for
func (c *CandidateContract) GetCandidateWithdrawInfos(nodeId discover.NodeID) ([]byte, error) {
	height := c.Evm.Context.BlockNumber
	log.Info("Input to CandidateWithdrawInfos", "blockNumber", height.String(), "nodeId: ", nodeId.String())

	refunds := c.Evm.CandidatePoolContext.GetDefeat(c.Evm.StateDB, nodeId, height)
	type WithdrawInfo struct {
		Balance        *big.Int
		LockNumber     *big.Int
		LockBlockCycle uint32
	}
	r := make([]WithdrawInfo, len(refunds))
	for i, v := range refunds {
		refundBlockNumber := c.Evm.CandidatePoolContext.GetRefundInterval(height)
		log.Debug("Call CandidateWithdrawInfos", "Deposit", v.Deposit, "BlockNumber", v.BlockNumber.String(), "RefundBlockNumber", refundBlockNumber)
		r[i] = WithdrawInfo{v.Deposit, v.BlockNumber, refundBlockNumber}
	}
	data, _ := json.Marshal(r)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of CandidateWithdrawInfos", "blockNumber", height.String(), "json: ", string(data))
	return sdata, nil
}

// GetCandidateDetails returns the batch of candidate info.
func (c *CandidateContract) GetCandidateDetails(nodeIds []discover.NodeID) ([]byte, error) {

	height := c.Evm.Context.BlockNumber
	input, _ := json.Marshal(nodeIds)
	log.Info("Input to GetCandidateDetails", "blockNumber", height.String(), "length: ", len(nodeIds), " nodeIds: ", string(input))
	candidates := c.Evm.CandidatePoolContext.GetCandidateArr(c.Evm.StateDB, height, nodeIds...)
	data, _ := json.Marshal(candidates)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetCandidateDetails", "blockNumber", height.String(), "len(candidates): ", len(candidates), "json: ", string(data))
	return sdata, nil
}

// Get the current block candidate list
func (c *CandidateContract) GetCandidateList() ([]byte, error) {

	height := c.Evm.Context.BlockNumber
	candidates := c.Evm.CandidatePoolContext.GetChosens(c.Evm.StateDB, 0, height)
	data, _ := json.Marshal(candidates)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetCandidateList", "blockNumber", height.String(), "len(candidates): ", len(candidates[0])+len(candidates[1]), "json: ", string(data))
	return sdata, nil
}

// Get the current block round certifier list
func (c *CandidateContract) GetVerifiersList() ([]byte, error) {
	height := c.Evm.Context.BlockNumber
	verifiers := c.Evm.CandidatePoolContext.GetChairpersons(c.Evm.StateDB, height)
	data, _ := json.Marshal(verifiers)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetVerifiersList", "blockNumber", height.String(), "len(verifiers): ", len(verifiers), "json: ", string(data))
	return sdata, nil
}

// addLog let the result add to event.
func (c *CandidateContract) addLog(event, data string) {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); nil != err {
		log.Error("Failed to CandidateContract addlog", "rlp encode fail: ", err.Error())
	}
	c.Evm.StateDB.AddLog(&types.Log{
		Address:     common.CandidatePoolAddr,
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data:        buf.Bytes(),
		BlockNumber: c.Evm.Context.BlockNumber.Uint64(),
	})
}
