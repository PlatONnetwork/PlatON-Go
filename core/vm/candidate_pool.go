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
	ErrOwnerNotOnly     = errors.New("Node ID cannot bind multiple owners")
	ErrPermissionDenied = errors.New("Transaction from address permission denied")
	ErrFeeIllegal       = errors.New("The fee is illegal")
	ErrDepositEmpty     = errors.New("Deposit balance not zero")
	ErrWithdrawEmpty    = errors.New("No withdrawal amount")
	ErrCandidateEmpty   = errors.New("CandidatePool is null")
)

const (
	CandidateDepositEvent       = "CandidateDepositEvent"
	CandidateApplyWithdrawEvent = "CandidateApplyWithdrawEvent"
	CandidateWithdrawEvent      = "CandidateWithdrawEvent"
	SetCandidateExtraEvent      = "SetCandidateExtraEvent"
)

type candidatePool interface {
	SetCandidate(state StateDB, nodeId discover.NodeID, can *types.Candidate) error
	GetCandidate(state StateDB, nodeId discover.NodeID) (*types.Candidate, error)
	GetCandidateArr(state StateDB, nodeIds ...discover.NodeID) (types.CandidateQueue, error)
	WithdrawCandidate(state StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error
	GetChosens(state StateDB, flag int) types.CandidateQueue
	GetChairpersons(state StateDB) types.CandidateQueue
	GetDefeat(state StateDB, nodeId discover.NodeID) (types.CandidateQueue, error)
	IsDefeat(state StateDB, nodeId discover.NodeID) (bool, error)
	RefundBalance(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) error
	GetOwner(state StateDB, nodeId discover.NodeID) common.Address
	SetCandidateExtra(state StateDB, nodeId discover.NodeID, extra string) error
	GetRefundInterval() uint64
	MaxCount() uint64
}

type CandidateContract struct {
	Contract *Contract
	Evm      *EVM
}

func (c *CandidateContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (c *CandidateContract) Run(input []byte) ([]byte, error) {
	if c.Evm.CandidatePool == nil {
		log.Error("Failed to Run==> ", "ErrCandidateEmpty: ", ErrCandidateEmpty.Error())
		return nil, ErrCandidateEmpty
	}
	var command = map[string]interface{}{
		"CandidateDetails":        c.CandidateDetails,
		"CandidateApplyWithdraw":  c.CandidateApplyWithdraw,
		"CandidateDeposit":        c.CandidateDeposit,
		"CandidateList":           c.CandidateList,
		"CandidateWithdraw":       c.CandidateWithdraw,
		"SetCandidateExtra":       c.SetCandidateExtra,
		"CandidateWithdrawInfos":  c.CandidateWithdrawInfos,
		"VerifiersList":           c.VerifiersList,
		"GetBatchCandidateDetail": c.GetBatchCandidateDetail,
	}
	return execute(input, command)
}

// Candidate Application && Increase Quality Deposit
func (c *CandidateContract) CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string) ([]byte, error) {
	deposit := c.Contract.value
	txHash := c.Evm.StateDB.TxHash()
	txIdx := c.Evm.StateDB.TxIdx()
	height := c.Evm.Context.BlockNumber
	from := c.Contract.caller.Address()
	log.Info("Input to CandidateDeposit==> ", "nodeId: ", nodeId.String(), " owner: ", owner.Hex(), " deposit: ", deposit,
		"  fee: ", fee, " txhash: ", txHash.Hex(), " txIdx: ", txIdx, " height: ", height, " from: ", from.Hex(),
		" host: ", host, " port: ", port, " extra: ", extra)
	if fee > 10000 {
		log.Error("Failed to CandidateDeposit==> ", "ErrFeeIllegal: ", ErrFeeIllegal.Error())
		return nil, ErrFeeIllegal
	}
	if deposit.Cmp(big.NewInt(0)) < 1 {
		log.Error("Failed to CandidateDeposit==> ", "ErrDepositEmpty: ", ErrDepositEmpty.Error())
		return nil, ErrDepositEmpty
	}
	can, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to CandidateDeposit==> ", "GetCandidate return err: ", err.Error())
		return nil, err
	}
	var alldeposit *big.Int
	if nil != can {
		if ok := bytes.Equal(can.Owner.Bytes(), owner.Bytes()); !ok {
			log.Error("Failed to CandidateDeposit==> ", "ErrOwnerNotOnly: ", ErrOwnerNotOnly.Error())
			return nil, ErrOwnerNotOnly
		}
		alldeposit = new(big.Int).Add(can.Deposit, deposit)
		log.Info("CandidateDeposit==> ", "alldeposit: ", alldeposit, " can.Deposit: ", can.Deposit, " deposit: ", deposit)
	} else {
		alldeposit = deposit
	}
	canDeposit := types.Candidate{
		alldeposit,
		height,
		txIdx,
		nodeId,
		host,
		port,
		owner,
		from,
		extra,
		fee,
		common.Hash{},
	}
	log.Info("CandidateDeposit==> ", "canDeposit: ", canDeposit)
	if err = c.Evm.CandidatePool.SetCandidate(c.Evm.StateDB, nodeId, &canDeposit); err != nil {
		log.Error("Failed to CandidateDeposit==> ", "SetCandidate return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "", "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateDepositEvent, string(event))
	log.Info("Result of CandidateDeposit==> ", "json: ", string(event))
	return nil, nil
}

// Apply for a refund of the deposit
func (c *CandidateContract) CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	from := c.Contract.caller.Address()
	height := c.Evm.Context.BlockNumber
	log.Info("Input to CandidateApplyWithdraw==> ", "nodeId: ", nodeId.String(), " from: ", from.Hex(), " txHash: ", txHash.Hex(), " withdraw: ", withdraw, " height: ", height)
	can, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to CandidateApplyWithdraw==> ", "GetCandidate return err: ", err.Error())
		return nil, err
	}
	if can.Deposit.Cmp(big.NewInt(0)) < 1 {
		log.Error("Failed to CandidateApplyWithdraw==> ", "ErrWithdrawEmpty: ", ErrWithdrawEmpty.Error())
		return nil, ErrWithdrawEmpty
	}
	if ok := bytes.Equal(can.Owner.Bytes(), from.Bytes()); !ok {
		log.Error("Failed to CandidateApplyWithdraw==> ", "ErrPermissionDenied: ", ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if withdraw.Cmp(can.Deposit) > 0 {
		withdraw = can.Deposit
	}
	if err := c.Evm.CandidatePool.WithdrawCandidate(c.Evm.StateDB, nodeId, withdraw, height); nil != err {
		log.Error("Failed to CandidateApplyWithdraw==> ", "WithdrawCandidate return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "", "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateApplyWithdrawEvent, string(event))
	log.Info("Result of CandidateApplyWithdraw==> ", "json: ", string(event))
	return nil, nil
}

// Deposit withdrawal
func (c *CandidateContract) CandidateWithdraw(nodeId discover.NodeID) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	height := c.Evm.Context.BlockNumber
	log.Info("Input to CandidateWithdraw==> ", "nodeId: ", nodeId.String(), " height: ", height, " txHash: ", txHash.Hex())
	if err := c.Evm.CandidatePool.RefundBalance(c.Evm.StateDB, nodeId, height); nil != err {
		log.Error("Failed to CandidateWithdraw==> ", "RefundBalance return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "", "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateWithdrawEvent, string(event))
	log.Info("Result of CandidateWithdraw==> ", "json: ", string(event))
	return nil, nil
}

// Get the refund history you have applied for
func (c *CandidateContract) CandidateWithdrawInfos(nodeId discover.NodeID) ([]byte, error) {
	log.Info("Input to CandidateWithdrawInfos==> ", "nodeId: ", nodeId.String())
	infos, err := c.Evm.CandidatePool.GetDefeat(c.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to CandidateWithdrawInfos==> ", "GetDefeat return err: ", err.Error())
		return nil, err
	}
	type WithdrawInfo struct {
		Balance        *big.Int
		LockNumber     *big.Int
		LockBlockCycle uint64
	}
	type WithdrawInfos struct {
		Ret    bool
		ErrMsg string
		Infos  []WithdrawInfo
	}
	r := WithdrawInfos{true, "success", make([]WithdrawInfo, len(infos))}
	for i, v := range infos {
		r.Infos[i] = WithdrawInfo{v.Deposit, v.BlockNumber, c.Evm.CandidatePool.GetRefundInterval()}
	}
	data, _ := json.Marshal(r)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of CandidateWithdrawInfos==> ", "json: ", string(data))
	return sdata, nil
}

// Set up additional information
func (c *CandidateContract) SetCandidateExtra(nodeId discover.NodeID, extra string) ([]byte, error) {
	txHash := c.Evm.StateDB.TxHash()
	from := c.Contract.caller.Address()
	log.Info("Input to SetCandidateExtra==> ", "nodeId: ", nodeId.String(), " extra: ", extra, " from: ", from.Hex(), " txHash: ", txHash.Hex())
	owner := c.Evm.CandidatePool.GetOwner(c.Evm.StateDB, nodeId)
	if ok := bytes.Equal(owner.Bytes(), from.Bytes()); !ok {
		log.Error("Failed to SetCandidateExtra==> ", "ErrPermissionDenied: ", ErrPermissionDenied.Error())
		return nil, ErrPermissionDenied
	}
	if err := c.Evm.CandidatePool.SetCandidateExtra(c.Evm.StateDB, nodeId, extra); nil != err {
		log.Error("Failed to SetCandidateExtra==> ", "SetCandidateExtra return err: ", err.Error())
		return nil, err
	}
	r := ResultCommon{true, "", "success"}
	event, _ := json.Marshal(r)
	c.addLog(SetCandidateExtraEvent, string(event))
	log.Info("Result of SetCandidateExtra==> ", "json: ", string(event))
	return nil, nil
}

// Get candidate details
func (c *CandidateContract) CandidateDetails(nodeId discover.NodeID) ([]byte, error) {
	log.Info("Input to CandidateDetails==> ", "nodeId: ", nodeId.String())
	candidate, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("Failed to CandidateDetails==> ", "GetCandidate return err: ", err.Error())
		return nil, err
	}
	if nil == candidate {
		log.Error("Failed to CandidateDetails==> ", "The query does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(candidate)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of CandidateDetails==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetBatchCandidateDetail returns the batch of candidate info.
func (c *CandidateContract) GetBatchCandidateDetail(nodeIds []discover.NodeID) ([]byte, error) {
	input, _ := json.Marshal(nodeIds)
	log.Info("Input to GetBatchCandidateDetail==>", "length: ", len(nodeIds), " nodeIds: ", string(input))
	candidates, err := c.Evm.CandidatePool.GetCandidateArr(c.Evm.StateDB, nodeIds...)
	if nil != err {
		log.Error("Failed to GetBatchCandidateDetail==> ", "GetCandidateArr return err: ", err.Error())
		return nil, err
	}
	if 0 == len(candidates) {
		log.Error("Failed to GetBatchCandidateDetail==> ", "The query does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(candidates)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of GetBatchCandidateDetail==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// Get the current block candidate list
func (c *CandidateContract) CandidateList() ([]byte, error) {
	candidates := c.Evm.CandidatePool.GetChosens(c.Evm.StateDB, 0)
	if 0 == len(candidates) {
		log.Error("Failed to CandidateList==> ", "The query does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(candidates)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of CandidateList==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// Get the current block round certifier list
func (c *CandidateContract) VerifiersList() ([]byte, error) {
	verifiers := c.Evm.CandidatePool.GetChairpersons(c.Evm.StateDB)
	if 0 == len(verifiers) {
		log.Error("Failed to VerifiersList==> ", "The query does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(verifiers)
	sdata := DecodeResultStr(string(data))
	log.Info("Result of VerifiersList==> ", "json: ", string(data), "[]byte: ", sdata)
	return sdata, nil
}

// addLog let the result add to event.
func (c *CandidateContract) addLog(event, data string) {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); nil != err {
		log.Error("Failed to addlog==> ", "rlp encode fail: ", err.Error())
	}
	c.Evm.StateDB.AddLog(&types.Log{
		Address:     common.CandidatePoolAddr,
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data:        buf.Bytes(),
		BlockNumber: c.Evm.Context.BlockNumber.Uint64(),
	})
}
