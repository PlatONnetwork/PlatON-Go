package core

import (
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type Context interface {
	SetState(state *state.StateDB)
	GetState() *state.StateDB
	SetHeader(header *types.Header)
	GetHeader() *types.Header
	SetBlockHash(blockHash common.Hash)
	GetBlockHash() common.Hash
	SetGasPool(gp *GasPool)
	GetGasPool() *GasPool
	CumulateBlockGasUsed(txGasUsed uint64)
	GetBlockGasUsed() uint64
	GetBlockGasUsedHolder() *uint64
	SetTxList(txList []*types.Transaction)
	GetTxList() []*types.Transaction
	GetTx(idx int) *types.Transaction
	GetPackedTxList() []*types.Transaction
	AddPackedTx(tx *types.Transaction)
	SetResult(idx int, result *Result)
	GetResults() []*Result
	GetReceipts() types.Receipts
	AddReceipt(receipt *types.Receipt)
	SetTxDag(txDag *TxDag)
	GetTxDag() *TxDag
	SetPoppedAddress(poppedAddress *common.Address)
	GetPoppedAddresses() map[*common.Address]struct{}
	//AddLogs([]*types.Log)
	GetLogs() []*types.Log
	GetEarnings() *big.Int
	AddEarnings(*big.Int)
	SetTimeout(isTimeout bool)
	IsTimeout() bool
}

type Result struct {
	fromStateObject *state.ParallelStateObject
	toStateObject   *state.ParallelStateObject
	receipt         *types.Receipt
	minerEarnings   *big.Int
	err             error
}

type PackBlockContext struct {
	state           *state.StateDB
	header          *types.Header
	blockHash       common.Hash
	gp              *GasPool
	txList          []*types.Transaction
	packedTxList    []*types.Transaction
	resultList      []*Result
	receipts        types.Receipts
	logs            []*types.Log
	txDag           *TxDag
	poppedAddresses map[*common.Address]struct{}

	blockGasUsedHolder *uint64
	earnings           *big.Int
	startTime          time.Time
	blockDeadline      time.Time
	timeout            *common.AtomicBool
}

func NewPackBlockContext(state *state.StateDB, header *types.Header, blockHash common.Hash, gp *GasPool, startTime, blockDeadline time.Time) *PackBlockContext {
	ctx := &PackBlockContext{}
	ctx.state = state
	ctx.header = header
	ctx.blockHash = blockHash
	ctx.gp = gp
	ctx.startTime = startTime
	ctx.blockDeadline = blockDeadline
	ctx.poppedAddresses = make(map[*common.Address]struct{})
	ctx.earnings = big.NewInt(0)
	ctx.blockGasUsedHolder = &header.GasUsed
	ctx.timeout = &common.AtomicBool{}
	return ctx
}

func (ctx *PackBlockContext) SetState(state *state.StateDB) {
	ctx.state = state
}
func (ctx *PackBlockContext) GetState() *state.StateDB {
	return ctx.state
}
func (ctx *PackBlockContext) SetHeader(header *types.Header) {
	ctx.header = header
}
func (ctx *PackBlockContext) GetHeader() *types.Header {
	return ctx.header
}
func (ctx *PackBlockContext) SetBlockHash(blockHash common.Hash) {
	ctx.blockHash = blockHash
}
func (ctx *PackBlockContext) GetBlockHash() common.Hash {
	return ctx.blockHash
}
func (ctx *PackBlockContext) SetGasPool(gp *GasPool) {
	ctx.gp = gp
}
func (ctx *PackBlockContext) GetGasPool() *GasPool {
	return ctx.gp
}

func (ctx *PackBlockContext) SetTxList(txs []*types.Transaction) {
	ctx.txList = txs
	ctx.resultList = make([]*Result, len(txs))
}
func (ctx *PackBlockContext) GetTxList() []*types.Transaction {
	return ctx.txList
}
func (ctx *PackBlockContext) GetTx(idx int) *types.Transaction {
	if len(ctx.txList) > idx {
		return ctx.txList[idx]
	} else {
		return nil
	}
}

func (ctx *PackBlockContext) GetPackedTxList() []*types.Transaction {
	return ctx.packedTxList
}
func (ctx *PackBlockContext) AddPackedTx(tx *types.Transaction) {
	ctx.packedTxList = append(ctx.packedTxList, tx)
}

func (ctx *PackBlockContext) SetResult(idx int, result *Result) {
	if idx > len(ctx.resultList)-1 {
		return
	}
	ctx.resultList[idx] = result
}
func (ctx *PackBlockContext) GetResults() []*Result {
	return ctx.resultList
}

func (ctx *PackBlockContext) GetReceipts() types.Receipts {
	return ctx.receipts
}

func (ctx *PackBlockContext) AddReceipt(receipt *types.Receipt) {
	ctx.receipts = append(ctx.receipts, receipt)
}

func (ctx *PackBlockContext) SetTxDag(txDag *TxDag) {
	ctx.txDag = txDag
}
func (ctx *PackBlockContext) GetTxDag() *TxDag {
	return ctx.txDag
}

func (ctx *PackBlockContext) SetPoppedAddress(poppedAddress *common.Address) {
	ctx.poppedAddresses[poppedAddress] = struct{}{}
}
func (ctx *PackBlockContext) GetPoppedAddresses() map[*common.Address]struct{} {
	return ctx.poppedAddresses
}

func (ctx *PackBlockContext) GetLogs() []*types.Log {
	return ctx.state.Logs()
}

func (ctx *PackBlockContext) CumulateBlockGasUsed(txGasUsed uint64) {
	*ctx.blockGasUsedHolder += txGasUsed
}
func (ctx *PackBlockContext) GetBlockGasUsed() uint64 {
	return *ctx.blockGasUsedHolder
}

func (ctx *PackBlockContext) GetBlockGasUsedHolder() *uint64 {
	return ctx.blockGasUsedHolder
}

func (ctx *PackBlockContext) GetEarnings() *big.Int {
	return ctx.earnings
}

func (ctx *PackBlockContext) AddEarnings(earning *big.Int) {
	ctx.earnings = new(big.Int).Add(ctx.earnings, earning)
}

func (ctx *PackBlockContext) SetStartTime(startTime time.Time) {
	ctx.startTime = startTime
}
func (ctx *PackBlockContext) GetStartTime() time.Time {
	return ctx.startTime
}

func (ctx *PackBlockContext) SetBlockDeadline(blockDeadline time.Time) {
	ctx.blockDeadline = blockDeadline
}

func (ctx *PackBlockContext) GetBlockDeadline() time.Time {
	return ctx.blockDeadline
}

func (ctx *PackBlockContext) SetTimeout(isTimeout bool) {
	ctx.timeout.Set(isTimeout)
}

func (ctx *PackBlockContext) IsTimeout() bool {
	return ctx.timeout.IsSet()
}

type VerifyBlockContext struct {
	state              *state.StateDB
	header             *types.Header
	blockHash          common.Hash
	gp                 *GasPool
	receipts           types.Receipts
	logs               []*types.Log
	txList             []*types.Transaction
	packedTxList       []*types.Transaction
	resultList         []*Result
	txDag              *TxDag
	poppedAddresses    map[*common.Address]struct{}
	blockGasUsedHolder *uint64
	earnings           *big.Int
	verifyResult       bool
	startTime          time.Time
}

func NewVerifyBlockContext(state *state.StateDB, header *types.Header, blockHash common.Hash, gp *GasPool, blockGasUsed *uint64, startTime time.Time) *VerifyBlockContext {
	ctx := &VerifyBlockContext{}
	ctx.state = state
	ctx.header = header
	ctx.blockHash = blockHash
	ctx.gp = gp
	ctx.poppedAddresses = make(map[*common.Address]struct{})
	ctx.earnings = big.NewInt(0)
	ctx.blockGasUsedHolder = blockGasUsed
	ctx.startTime = startTime
	return ctx
}

func (ctx *VerifyBlockContext) SetState(state *state.StateDB) {
	ctx.state = state
}
func (ctx *VerifyBlockContext) GetState() *state.StateDB {
	return ctx.state
}
func (ctx *VerifyBlockContext) SetHeader(header *types.Header) {
	ctx.header = header
}
func (ctx *VerifyBlockContext) GetHeader() *types.Header {
	return ctx.header
}
func (ctx *VerifyBlockContext) SetBlockHash(blockHash common.Hash) {
	ctx.blockHash = blockHash
}
func (ctx *VerifyBlockContext) GetBlockHash() common.Hash {
	return ctx.blockHash
}
func (ctx *VerifyBlockContext) SetGasPool(gp *GasPool) {
	ctx.gp = gp
}
func (ctx *VerifyBlockContext) GetGasPool() *GasPool {
	return ctx.gp
}

func (ctx *VerifyBlockContext) SetTxList(txs []*types.Transaction) {
	ctx.txList = txs
	ctx.resultList = make([]*Result, len(txs))
}
func (ctx *VerifyBlockContext) GetTxList() []*types.Transaction {
	return ctx.txList
}

func (ctx *VerifyBlockContext) GetTx(idx int) *types.Transaction {
	if len(ctx.txList) > idx {
		return ctx.txList[idx]
	} else {
		return nil
	}
}

func (ctx *VerifyBlockContext) GetPackedTxList() []*types.Transaction {
	return ctx.packedTxList
}
func (ctx *VerifyBlockContext) AddPackedTx(tx *types.Transaction) {
	ctx.packedTxList = append(ctx.packedTxList, tx)
}

func (ctx *VerifyBlockContext) SetResult(idx int, result *Result) {
	if idx > len(ctx.resultList)-1 {
		return
	}
	ctx.resultList[idx] = result
}
func (ctx *VerifyBlockContext) GetResults() []*Result {
	return ctx.resultList
}

func (ctx *VerifyBlockContext) GetReceipts() types.Receipts {
	return ctx.receipts
}

func (ctx *VerifyBlockContext) AddReceipt(receipt *types.Receipt) {
	ctx.receipts = append(ctx.receipts, receipt)
}

func (ctx *VerifyBlockContext) SetTxDag(txDag *TxDag) {
	ctx.txDag = txDag
}
func (ctx *VerifyBlockContext) GetTxDag() *TxDag {
	return ctx.txDag
}

func (ctx *VerifyBlockContext) SetPoppedAddress(poppedAddress *common.Address) {
	ctx.poppedAddresses[poppedAddress] = struct{}{}
}
func (ctx *VerifyBlockContext) GetPoppedAddresses() map[*common.Address]struct{} {
	return ctx.poppedAddresses
}

func (ctx *VerifyBlockContext) GetLogs() []*types.Log {
	return ctx.state.Logs()
}

func (ctx *VerifyBlockContext) CumulateBlockGasUsed(txGasUsed uint64) {
	*ctx.blockGasUsedHolder += txGasUsed
}
func (ctx *VerifyBlockContext) GetBlockGasUsed() uint64 {
	return *ctx.blockGasUsedHolder
}

func (ctx *VerifyBlockContext) GetBlockGasUsedHolder() *uint64 {
	return ctx.blockGasUsedHolder
}

func (ctx *VerifyBlockContext) GetEarnings() *big.Int {
	return ctx.earnings
}

func (ctx *VerifyBlockContext) AddEarnings(earning *big.Int) {
	ctx.earnings = new(big.Int).Add(ctx.earnings, earning)
}

func (ctx *VerifyBlockContext) SetVerifyResult(verifyResult bool) {
	ctx.verifyResult = verifyResult
}
func (ctx *VerifyBlockContext) GetVerifyResult() bool {
	return ctx.verifyResult
}

func (ctx *VerifyBlockContext) SetTimeout(isTimeout bool) {
	return
}

func (ctx *VerifyBlockContext) IsTimeout() bool {
	return false
}
