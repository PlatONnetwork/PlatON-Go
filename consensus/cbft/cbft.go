// Package bft implements the BFT consensus engine.
package cbft

import (
	"bytes"
	"container/list"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	math2 "github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
)

const (
	former int32 = iota
	current
	next
	all
)

var (
	errSign                = errors.New("sign error")
	errUnauthorizedSigner  = errors.New("unauthorized signer")
	errIllegalBlock        = errors.New("illegal block")
	lateBlock              = errors.New("block is late")
	errDuplicatedBlock     = errors.New("duplicated block")
	errBlockNumber         = errors.New("error block number")
	errUnknownBlock        = errors.New("unknown block")
	errFutileBlock         = errors.New("futile block")
	errGenesisBlock        = errors.New("cannot handle genesis block")
	errHighestLogicalBlock = errors.New("cannot find a logical block")
	errListConfirmedBlocks = errors.New("list confirmed blocks error")
	errMissingSignature    = errors.New("extra-data 65 byte signature suffix missing")
	extraSeal              = 65
	windowSize             = 10

	//periodMargin is a percentum for period margin
	periodMargin = uint64(20)

	//maxPingLatency is the time in milliseconds between Ping and Pong
	maxPingLatency = int64(5000)

	//maxAvgLatency is the time in milliseconds between two peers
	maxAvgLatency = int64(2000)
)

type Cbft struct {
	config                *params.CbftConfig
	ppos                  *ppos
	rotating              *rotating
	blockSignOutCh        chan *cbfttypes.BlockSignature //a channel to send block signature
	cbftResultOutCh       chan *cbfttypes.CbftResult     //a channel to send consensus result
	highestLogicalBlockCh chan *types.Block
	closeOnce             sync.Once
	exitCh                chan struct{}
	txPool                *core.TxPool
	blockExtMap           sync.Map         //map[common.Hash]*BlockExt //store all received blocks and signs
	dataReceiveCh         chan interface{} //a channel to receive data from miner
	blockChain            *core.BlockChain //the block chain
	highestLogical        atomic.Value     //highest block in logical path, local packages new block will base on it
	highestConfirmed      atomic.Value     //highest confirmed block in logical path
	rootIrreversible      atomic.Value     //the latest block has stored in chain
	//signedSet             map[uint64]struct{} //all block numbers signed by local node
	signedSet       sync.Map
	blockChainCache *core.BlockChainCache
	netLatencyMap   map[discover.NodeID]*list.List
	netLatencyLock  sync.RWMutex
}

func (cbft *Cbft) getRootIrreversible() *BlockExt {
	v := cbft.rootIrreversible.Load()
	if v == nil {
		return nil
	} else {
		return v.(*BlockExt)
	}
}

func (cbft *Cbft) getHighestConfirmed() *BlockExt {
	v := cbft.highestConfirmed.Load()
	if v == nil {
		return nil
	} else {
		return v.(*BlockExt)
	}
}
func (cbft *Cbft) getHighestLogical() *BlockExt {
	v := cbft.highestLogical.Load()
	if v == nil {
		return nil
	} else {
		return v.(*BlockExt)
	}
}

var cbft *Cbft

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, blockSignatureCh chan *cbfttypes.BlockSignature, cbftResultCh chan *cbfttypes.CbftResult, highestLogicalBlockCh chan *types.Block) *Cbft {

	_ppos := newPpos(config)

	cbft = &Cbft{
		config:                config,
		ppos:                  _ppos,
		rotating:              newRotating(_ppos, config.Duration),
		blockSignOutCh:        blockSignatureCh,
		cbftResultOutCh:       cbftResultCh,
		highestLogicalBlockCh: highestLogicalBlockCh,
		exitCh:                make(chan struct{}),
		//blockExtMap:         make(map[common.Hash]*BlockExt),
		//signedSet:     make(map[uint64]struct{}),
		dataReceiveCh: make(chan interface{}, 256),
		netLatencyMap: make(map[discover.NodeID]*list.List),
	}

	_ppos.ticketContext.SetChainInfo(cbft)

	flowControl = NewFlowControl()

	go cbft.dataReceiverLoop()

	return cbft
}

// BlockExt is an extension from Block
type BlockExt struct {
	block       *types.Block	`json:"-"`
	inTree      bool	`json:"-"`
	inTurn      bool	`json:"-"`
	isExecuted  bool	`json:"-"`
	isSigned    bool	`json:"-"`
	isConfirmed bool 	`json:"-"`
	Number      uint64 	`json:"number"`
	Hash        string 	`json:"hash"`
	rcvTime     int64
	signs       []*common.BlockConfirmSign `json:"-"`//all signs for block
	parent      *BlockExt
	Children    []*BlockExt	`json:"children"`

}


func (ext *BlockExt) toJson() string {
	jsons, errs := json.Marshal(ext)
	if errs != nil {
		log.Warn("can marshal BlockExt struct to json string", "number", ext.Number, "hash", ext.Hash)
		return ""
	}
	return string(jsons)
}


// New creates a BlockExt object
func NewBlockExt(block *types.Block, blockNum uint64) *BlockExt {
	return &BlockExt{
		block:  block,
		Number: blockNum,
		signs:  make([]*common.BlockConfirmSign, 0),
		Hash:	block.Hash().TerminalString(),
	}
}

func NewEmptyExt(blockNum uint64) *BlockExt {
	return &BlockExt{
		Number: blockNum,
		signs:  make([]*common.BlockConfirmSign, 0),
	}
}

var flowControl *FlowControl

// FlowControl is a rectifier for sequential blocks
type FlowControl struct {
	nodeID      discover.NodeID
	lastTime    int64
	maxInterval int64
	minInterval int64
}

func NewFlowControl() *FlowControl {
	return &FlowControl{
		nodeID:      discover.NodeID{},
		maxInterval: int64(cbft.config.Period*1000 + cbft.config.Period*1000*periodMargin/100),
		minInterval: int64(cbft.config.Period*1000 - cbft.config.Period*1000*periodMargin/100),
	}
}

// control checks if the block is received at a proper rate
func (flowControl *FlowControl) control(nodeID discover.NodeID, rcvTime int64) bool {
	passed := false
	if flowControl.nodeID == nodeID {
		differ := rcvTime - flowControl.lastTime
		if differ >= flowControl.minInterval && differ <= flowControl.maxInterval {
			passed = true
		} else {
			passed = false
		}
	} else {
		passed = true
	}
	flowControl.nodeID = nodeID
	flowControl.lastTime = rcvTime

	return passed
}

// findBlockExt finds BlockExt in cbft.blockExtMap
func (cbft *Cbft) findBlockExt(hash common.Hash) *BlockExt {
	if v, ok := cbft.blockExtMap.Load(hash); ok {
		return v.(*BlockExt)
	}
	return nil
}

//collectSign collects all signs for a block
func (cbft *Cbft) collectSign(ext *BlockExt, sign *common.BlockConfirmSign) {
	if sign != nil {
		ext.signs = append(ext.signs, sign)
		blockNumber := big.NewInt((int64(ext.Number)))
		parentNumber := new(big.Int).Sub(blockNumber, common.Big1)
		//if ext.isLinked && ext.block != nil {
		if ext.inTree { // ext.block != nil is unnecessary
			if len(ext.signs) >= cbft.getThreshold(parentNumber, ext.block.ParentHash(), blockNumber) {
				ext.isConfirmed = true
			}
		}
	}
}

// isParent checks if a block is another's parent
func (parent *BlockExt) isParent(child *types.Block) bool {
	if parent.block != nil && parent.block.NumberU64()+1 == child.NumberU64() && parent.block.Hash() == child.ParentHash() {
		return true
	}
	return false
}

// findParent finds ext's parent with non-nil block
func (cbft *Cbft) findParent(ext *BlockExt) *BlockExt {
	if ext.block == nil {
		return nil
	}
	parent := cbft.findBlockExt(ext.block.ParentHash())
	if parent != nil {
		if parent.block == nil {
			log.Warn("parent block has not received")
		} else if parent.block.NumberU64()+1 == ext.block.NumberU64() {
			return parent
		} else {
			log.Warn("data error, parent block hash is not mapping to number")
		}
	}
	return nil
}

// collectTxs collects exts's transactions
func (cbft *Cbft) collectTxs(exts []*BlockExt) types.Transactions {
	txs := make([]*types.Transaction, 0)
	for _, ext := range exts {
		copy(txs, ext.block.Transactions())
	}
	return types.Transactions(txs)
}

// findChildren finds current blockExt's all children with non-nil block
func (cbft *Cbft) findChildren(parent *BlockExt) []*BlockExt {
	if parent.block == nil {
		return nil
	}
	children := make([]*BlockExt, 0)

	f := func(k, v interface{}) bool {
		// The input and input types of this function are fixed and cannot be modified.
		// You can write your own code in the body of the function, call k, v in the map
		child := v.(*BlockExt)
		if child.block != nil && child.block.ParentHash() == parent.block.Hash() {
			if child.block.NumberU64()-1 == parent.block.NumberU64() {
				children = append(children, child)
			} else {
				log.Warn("data error, child block hash is not mapping to number")
			}
		}
		return true
	}
	cbft.blockExtMap.Range(f)

	if len(children) == 0 {
		return nil
	} else {
		return children
	}
}

// saveBlockExt saves block in memory
func (cbft *Cbft) saveBlockExt(hash common.Hash, ext *BlockExt) {
	cbft.blockExtMap.Store(hash, ext)

	length := 0
	cbft.blockExtMap.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	log.Debug("save block in memory", "hash", hash, "number", ext.Number, "totalBlocks", length)
}

// isAncestor checks if a block is another's ancestor
func (lower *BlockExt) isAncestor(higher *BlockExt) bool {

	if higher == nil || higher.block == nil || lower == nil || lower.block == nil {
		return false
	}
	generations := higher.block.NumberU64() - lower.block.NumberU64()
	if generations <= 0 {
		return false
	}

	for i := uint64(0); i < generations; i++ {
		parent := higher.parent
		if parent != nil {
			higher = parent
		} else {
			return false
		}
	}

	if lower.block.Hash() == higher.block.Hash() && lower.block.NumberU64() == higher.block.NumberU64() {
		return true
	}
	return false
}

// findHighest finds the highest block from current start; If there are multiple highest blockExts, returns the one that singed by self; if none of blocks signed by self, returns the one that has most signs
func (cbft *Cbft) findHighest(current *BlockExt) *BlockExt {
	highest := current
	for _, child := range current.Children {
		current := cbft.findHighest(child)
		if current.block.NumberU64() > highest.block.NumberU64() || (current.block.NumberU64() == highest.block.NumberU64() && (current.isSigned || len(current.signs) > len(highest.signs))) {
			highest = current
		}
	}
	return highest
}

// findHighestLogical finds a logical path and return the highest block.
// the precondition is cur is a logical block, so, findHighestLogical will return cur if the path only has one block.
func (cbft *Cbft) findHighestLogical(cur *BlockExt) *BlockExt {
	lastClosestConfirmed := cbft.findLastClosestConfirmedIncludingSelf(cur)
	if lastClosestConfirmed == nil {
		return cbft.findHighest(cur)
	} else {
		return cbft.findHighest(lastClosestConfirmed)
	}
}

// findLastClosestConfirmedIncludingSelf return the last found block by call findClosestConfirmedExcludingSelf in a circular manner
func (cbft *Cbft) findLastClosestConfirmedIncludingSelf(cur *BlockExt) *BlockExt {
	log.Debug("findLastClosestConfirmedIncludingSelf", "cur.number", cur.Number, "cur.hash", cur.Hash, "cur.isConfirmed", cur.isConfirmed)
	var lastClosestConfirmed *BlockExt
	for {
		lastClosestConfirmed = cbft.findClosestConfirmedExcludingSelf(cur)
		if lastClosestConfirmed == nil || lastClosestConfirmed.block.Hash() == cur.block.Hash() {
			break
		} else {
			//fmt.Printf("lastClosestConfirmed, number=%d, rcvTime=%d\r\n", lastClosestConfirmed.Number, lastClosestConfirmed.rcvTime)
			cur = lastClosestConfirmed
		}
	}
	if lastClosestConfirmed != nil {
		return lastClosestConfirmed
	} else if cur.isConfirmed {
		return cur
	} else {
		log.Debug("findLastClosestConfirmedIncludingSelf return nil", "cur.number", cur.Number, "cur.hash", cur.Hash, "cur.isConfirmed", cur.isConfirmed)
		return nil
	}
}

// findClosestConfirmedIncludingSelf returns the closest confirmed block in current's descendant (including current itself).
// return nil if there's no confirmed in current's descendant.
func (cbft *Cbft) findClosestConfirmedIncludingSelf(current *BlockExt) *BlockExt {
	closest := current
	if current.inTree && current.isExecuted && !current.isConfirmed {
		closest = nil
	}
	for _, child := range current.Children {
		temp := cbft.findClosestConfirmedIncludingSelf(child)
		if closest == nil || (temp != nil && temp.inTree && temp.isExecuted && temp.isConfirmed && temp.Number < closest.Number) {
			closest = temp
		}
	}
	return closest
}

// findClosestConfirmedExcludingSelf returns the closest confirmed block in current's descendant (excluding current itself).
// return nil if there's no confirmed in current's descendant.
func (cbft *Cbft) findClosestConfirmedExcludingSelf(current *BlockExt) *BlockExt {
	var closest *BlockExt
	for _, child := range current.Children {
		if child != nil && child.inTree && child.isExecuted && child.isConfirmed {
			return child
		} else {
			temp := cbft.findClosestConfirmedIncludingSelf(child)
			if closest == nil || (temp != nil && temp.Number < closest.Number) {
				closest = temp
			}
		}
	}
	return closest
}

// signLogicalAndDescendant signs logical block go along the logical path from current block, and will not sign the block if there's another same number block has been signed.
func (cbft *Cbft) signLogicalAndDescendant(current *BlockExt) {
	log.Debug("sign logical block and its descendant", "hash", current.block.Hash(), "number", current.block.NumberU64())
	highestLogical := cbft.findHighestLogical(current)

	logicalBlocks := cbft.backTrackBlocks(highestLogical, current, true)

	//var highestConfirmed *BlockExt
	for _, logical := range logicalBlocks {
		if logical.inTurn && !logical.isSigned {
			//if _, signed := cbft.signedSet[logical.block.NumberU64()]; !signed {
			if _, signed := cbft.signedSet.Load(logical.block.NumberU64()); !signed {
				cbft.sign(logical)
				log.Debug("reset TxPool after block signed", "hash", logical.block.Hash(), "number", logical.Number)
				cbft.txPool.Reset(logical.block)
			}
		}
	}
}

func (cbft *Cbft) signLogicalAndDescendantMock(current *BlockExt) {
	highestLogical := cbft.findHighestLogical(current)
	logicalBlocks := cbft.backTrackBlocks(highestLogical, current, true)

	for _, logical := range logicalBlocks {
		if logical.inTurn && !logical.isSigned {
			if _, signed := cbft.signedSet.Load(logical.block.NumberU64()); !signed {
				cbft.signMock(logical)
			}
		}
	}
}

// executeBlockAndDescendant executes the block's transactions and its descendant
func (cbft *Cbft) executeBlockAndDescendant(current *BlockExt, parent *BlockExt) error {
	if !current.isExecuted {
		if err := cbft.execute(current, parent); err != nil {
			current.inTree = false
			current.isExecuted = false
			//remove bad block from tree and map
			cbft.removeBadBlock(current)
			//log.Error("execute block error", "hash", current.block.Hash(), "number", current.block.NumberU64())
			return err
		} else {
			current.inTree = true
			current.isExecuted = true
		}
	}

	for _, child := range current.Children {
		if err := cbft.executeBlockAndDescendant(child, current); err != nil {
			//remove bad block from tree and map
			cbft.removeBadBlock(child)
			return err
		}
	}
	return nil
}


func (cbft *Cbft) executeBlockAndDescendantMock(current *BlockExt, parent *BlockExt) error {
	if !current.isExecuted {
		current.inTree = true
		current.isExecuted = true
	}

	for _, child := range current.Children {
		if err := cbft.executeBlockAndDescendantMock(child, current); err != nil {
			//remove bad block from tree and map
			cbft.removeBadBlock(child)
			return err
		}
	}
	return nil
}

// sign signs a block
func (cbft *Cbft) sign(ext *BlockExt) {
	sealHash := ext.block.Header().SealHash()
	if signature, err := cbft.signFn(sealHash.Bytes()); err == nil {
		log.Debug("Sign block ", "hash", ext.block.Hash(), "number", ext.block.NumberU64(), "sealHash", sealHash, "signature", hexutil.Encode(signature[:8]))

		sign := common.NewBlockConfirmSign(signature)
		ext.isSigned = true

		cbft.collectSign(ext, sign)

		//save this block number
		//cbft.signedSet[ext.block.NumberU64()] = struct{}{}
		cbft.signedSet.Store(ext.block.NumberU64(), struct{}{})

		blockHash := ext.block.Hash()

		//send the BlockSignature to channel
		blockSign := &cbfttypes.BlockSignature{
			SignHash:   sealHash,
			Hash:       blockHash,
			Number:     ext.block.Number(),
			Signature:  sign,
			ParentHash: ext.block.ParentHash(),
		}
		cbft.blockSignOutCh <- blockSign
	} else {
		panic("sign block fatal error")
	}
}

func (cbft *Cbft) signMock(ext *BlockExt) {
	sealHash := ext.block.Header().SealHash()
	if signature, err := cbft.signFn(sealHash.Bytes()); err == nil {
		sign := common.NewBlockConfirmSign(signature)
		ext.isSigned = true

		cbft.collectSign(ext, sign)

		//save this block number
		//cbft.signedSet[ext.block.NumberU64()] = struct{}{}
		cbft.signedSet.Store(ext.block.NumberU64(), struct{}{})

	} else {
		panic("sign block fatal error")
	}
}

// execute executes the block's transactions based on its parent
// if success then save the receipts and state to consensusCache
func (cbft *Cbft) execute(ext *BlockExt, parent *BlockExt) error {
	log.Debug("execute block", "hash", ext.block.Hash(), "number", ext.block.NumberU64(), "ParentHash", parent.block.Hash())
	state, err := cbft.blockChainCache.MakeStateDB(parent.block)
	if err != nil {
		log.Error("execute block error, cannot make state based on parent", "hash", ext.block.Hash(), "Number", ext.block.NumberU64(), "ParentHash", parent.block.Hash(), "err", err)
		return errors.New("execute block error")
	}

	//to execute
	blockInterval := new(big.Int).Sub(ext.block.Number(), cbft.blockChain.CurrentBlock().Number())
	receipts, err := cbft.blockChain.ProcessDirectly(ext.block, state, parent.block, blockInterval)

	if err == nil {
		//save the receipts and state to consensusCache
		stateIsNil := state == nil
		log.Debug("execute block success", "hash", ext.block.Hash(), "number", ext.block.NumberU64(), "ParentHash", parent.block.Hash(), "lenReceipts", len(receipts), "stateIsNil", stateIsNil, "root", ext.block.Root())
		sealHash := ext.block.Header().SealHash()
		cbft.blockChainCache.WriteReceipts(sealHash, receipts, ext.block.NumberU64())
		cbft.blockChainCache.WriteStateDB(sealHash, state, ext.block.NumberU64())
		//cbft.blockChainCache.MarkBlockHash(ext.block.Hash())
	} else {
		log.Error("execute block error", "hash", ext.block.Hash(), "number", ext.block.NumberU64(), "ParentHash", parent.block.Hash(), "err", err)
		return errors.New("execute block error")
	}
	return nil
}

// backTrackBlocks return blocks from start to end, these blocks are in a same tree branch.
// The result is sorted by block number from lower to higher.
func (cbft *Cbft) backTrackBlocks(start *BlockExt, end *BlockExt, includeEnd bool) []*BlockExt {
	log.Trace("back track blocks", "startHash", start.block.Hash(), "startParentHash", end.block.ParentHash(), "endHash", start.block.Hash())

	result := make([]*BlockExt, 0)

	if start.block.Hash() == end.block.Hash() && includeEnd {
		result = append(result, start)
	} else if start.block.NumberU64() > end.block.NumberU64() {
		found := false
		result = append(result, start)

		for {
			parent := start.parent
			if parent == nil {
				break
			} else if parent.block.Hash() == end.block.Hash() && parent.block.NumberU64() == end.block.NumberU64() {
				//log.Debug("ending of back track block ")
				if includeEnd {
					result = append(result, parent)
				}
				found = true
				break
			} else {
				//log.Debug("found new block", "hash", parent.block.Hash(), "ParentHash", parent.block.ParentHash(), "number", parent.block.NumberU64())
				result = append(result, parent)
				start = parent
			}
		}

		if found {
			//sorted by block number from lower to higher
			if len(result) > 1 {
				reverse(result)
			}
		} else {
			result = nil
		}
	}
	return result
}

func reverse(s []*BlockExt) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// SetPrivateKey sets local's private key by the backend.go
func (cbft *Cbft) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	cbft.config.PrivateKey = privateKey
	cbft.config.NodeID = discover.PubkeyID(&privateKey.PublicKey)
}

func SetBlockChainCache(blockChainCache *core.BlockChainCache) {
	cbft.blockChainCache = blockChainCache
}

// setHighestLogical sets highest logical block and send it to the highestLogicalBlockCh
func (cbft *Cbft) setHighestLogical(highestLogical *BlockExt) {
	cbft.highestLogical.Store(highestLogical)
	cbft.highestLogicalBlockCh <- highestLogical.block
}

// SetBackend sets blockChain and txPool into cbft
func SetBackend(blockChain *core.BlockChain, txPool *core.TxPool) {
	log.Debug("call SetBackend()")
	cbft.blockChain = blockChain
	cbft.ppos.SetStartTimeOfEpoch(blockChain.Genesis().Time().Int64() / 1000)
	cbft.ppos.ticketContext.SetChainConfig(blockChain.Config())

	currentBlock := blockChain.CurrentBlock()

	genesisParentHash := bytes.Repeat([]byte{0x00}, 32)
	if bytes.Equal(currentBlock.ParentHash().Bytes(), genesisParentHash) && currentBlock.Number() == nil {
		currentBlock.Header().Number = big.NewInt(0)
	}

	log.Debug("init cbft.highestLogicalBlock", "hash", currentBlock.Hash(), "number", currentBlock.NumberU64())

	current := NewBlockExt(currentBlock, currentBlock.NumberU64())
	current.inTree = true
	current.isExecuted = true
	current.isSigned = true
	current.isConfirmed = true
	current.Number = currentBlock.NumberU64()

	cbft.saveBlockExt(currentBlock.Hash(), current)

	cbft.highestConfirmed.Store(current)

	//cbft.highestLogical = current
	cbft.setHighestLogical(current)

	cbft.rootIrreversible.Store(current)

	cbft.txPool = txPool
}

func SetPposOption(blockChain *core.BlockChain) {
	cbft.ppos.setPPOS_Temp()
	cbft.ppos.SetCandidateContextOption(blockChain, cbft.config.InitialNodes)
}

// BlockSynchronisation reset the cbft env, such as cbft.highestLogical, cbft.highestConfirmed.
// This function is invoked after that local has synced new blocks from other node.
func (cbft *Cbft) OnBlockSynced() {
	log.Debug("call OnBlockSynced()", "cbft.dataReceiveCh.len", len(cbft.dataReceiveCh))
	cbft.dataReceiveCh <- &cbfttypes.BlockSynced{}
}

func (cbft *Cbft) blockSynced() {
	currentBlock := cbft.blockChain.CurrentBlock()
	log.Debug("=== call blockSynced() ===",
		"highestLogicalHash", cbft.getHighestLogical().block.Hash(),
		"highestLogicalNumber", cbft.getHighestLogical().Number,
		"highestConfirmedHash", cbft.getHighestConfirmed().block.Hash(),
		"highestConfirmedNumber", cbft.getHighestConfirmed().Number,
		"rootIrreversibleHash", cbft.getRootIrreversible().block.Hash(),
		"rootIrreversibleNumber", cbft.getRootIrreversible().Number,
		"number", currentBlock.Number(),
		"hash", currentBlock.Hash())

	if currentBlock.NumberU64() > cbft.getRootIrreversible().Number {
		log.Debug("chain has a higher irreversible block", "hash", currentBlock.Hash(), "number", currentBlock.NumberU64())
		newRoot := cbft.findBlockExt(currentBlock.Hash())
		if newRoot == nil || newRoot.block == nil {
			log.Debug("higher irreversible block is not existing in memory", "newTree", newRoot.toJson())
			//the block synced from other peer is a new block in local peer
			//remove all blocks referenced in old tree after being cut off
			cbft.cleanByTailoredTree(cbft.getRootIrreversible())

			newRoot = NewBlockExt(currentBlock, currentBlock.NumberU64())
			newRoot.inTree = true
			newRoot.isExecuted = true
			newRoot.isSigned = true
			newRoot.isConfirmed = true
			newRoot.Number = currentBlock.NumberU64()
			newRoot.parent = nil
			cbft.saveBlockExt(newRoot.block.Hash(), newRoot)

			//only need to build up child
			cbft.buildChildNode(newRoot)

			if len(newRoot.Children) > 0{
				//the new root's children should re-execute base on new state
				for _, child := range newRoot.Children {
					if err := cbft.executeBlockAndDescendant(child, newRoot); err != nil {
						log.Error("execute the block error", "err", err)
						break
					}
				}
				//there are some redundancy code for newRoot, but these codes are necessary for other logical blocks
				cbft.signLogicalAndDescendant(newRoot)
			}

		} else if newRoot.block != nil {
			log.Debug("higher irreversible block is existing in memory","newTree", newRoot.toJson())
			//the block synced from other peer exists in local peer
			newRoot.isExecuted = true
			newRoot.isSigned = true
			newRoot.isConfirmed = true
			newRoot.Number = currentBlock.NumberU64()

			if newRoot.inTree == false {
				newRoot.inTree = true

				cbft.setDescendantInTree(newRoot)

				if len(newRoot.Children) > 0{
					//the new root's children should re-execute base on new state
					for _, child := range newRoot.Children {
						if err := cbft.executeBlockAndDescendant(child, newRoot); err != nil {
							log.Error("execute the block error", "err", err)
							break
						}
					}
					//there are some redundancy code for newRoot, but these codes are necessary for other logical blocks
					cbft.signLogicalAndDescendant(newRoot)
				}
			}else{
				//cut off old tree from new root,
				tailorTree(newRoot)
			}

			//remove all blocks referenced in old tree after been cut off
			cbft.cleanByTailoredTree(cbft.getRootIrreversible())
		}

		//remove all other blocks those their numbers are too low
		cbft.cleanByNumber(newRoot.Number)

		log.Debug("the cleared new tree in memory", "json", newRoot.toJson())

		//reset the new root irreversible
		cbft.rootIrreversible.Store(newRoot)

		log.Debug("reset the new root irreversible by synced", "hash", newRoot.block.Hash(), "number", newRoot.block.NumberU64())

		//reset highest confirmed block
		highestConfirmed :=cbft.findLastClosestConfirmedIncludingSelf(newRoot)
		cbft.highestConfirmed.Store(highestConfirmed)
		if cbft.getHighestConfirmed() != nil {
			log.Debug("reset the highestConfirmed by synced successful", "hash", highestConfirmed.block.Hash(), "number", highestConfirmed.block.NumberU64())

			//reset logical path
			highestLogical := cbft.findHighest(highestConfirmed)
			cbft.setHighestLogical(highestLogical)
			log.Debug("reset the highestLogical by synced", "hash", highestLogical.block.Hash(), "number", highestLogical.block.NumberU64(), "newRoot.hash", newRoot.block.Hash(), "newRoot.number", newRoot.block.NumberU64() )

			if !cbft.flushReadyBlock() {
				//remove all other blocks those their numbers are too low
				cbft.cleanByNumber(cbft.getRootIrreversible().Number)
			}

			log.Debug("reset TxPool after block synced", "hash", currentBlock.Hash(), "number", currentBlock.NumberU64())
			cbft.txPool.Reset(currentBlock)

		} else {
			log.Debug("reset the highestConfirmed by synced failure because findLastClosestConfirmedIncludingSelf() returned nil", "newRoot.hash", newRoot.block.Hash(), "newRoot.number", newRoot.block.NumberU64())
			//cbft.highestConfirmed.Store(newRoot)
		}
	}

	log.Debug("=== end of blockSynced() ===",
		"highestLogicalHash", cbft.getHighestLogical().block.Hash(),
		"highestLogicalNumber", cbft.getHighestLogical().Number,
		"highestConfirmedHash", cbft.getHighestConfirmed().block.Hash(),
		"highestConfirmedNumber", cbft.getHighestConfirmed().Number,
		"rootIrreversibleHash", cbft.getRootIrreversible().block.Hash(),
		"rootIrreversibleNumber", cbft.getRootIrreversible().Number,
		"number", currentBlock.Number(),
		"hash", currentBlock.Hash())
}

// dataReceiverLoop is the main loop that handle the data from worker, or eth protocol's handler
// the new blocks packed by local in worker will be handled here; the other blocks and signs received by P2P will be handled here.
func (cbft *Cbft) dataReceiverLoop() {
	for {
		select {
		case v := <-cbft.dataReceiveCh:
			sign, ok := v.(*cbfttypes.BlockSignature)
			if ok {
				err := cbft.signReceiver(sign)
				if err != nil {
					log.Error("Error", "msg", err)
				}
			} else {
				blockExt, ok := v.(*BlockExt)
				if ok {
					err := cbft.blockReceiver(blockExt)
					if err != nil {
						log.Error("Error", "msg", err)
					}
				} else {
					_, ok := v.(*cbfttypes.BlockSynced)
					if ok {
						cbft.blockSynced()
					} else {
						log.Error("Received wrong data type")
					}
				}
			}
		case <-cbft.exitCh:
			log.Debug("consensus engine exit")
			return
		}
	}
}

// buildIntoTree inserts current BlockExt to the tree structure
func (cbft *Cbft) buildIntoTree(current *BlockExt) {
	parent := cbft.findParent(current)
	if parent != nil {
		//catch up with parent
		parent.Children = append(parent.Children, current)
		current.parent = parent
		current.inTree = parent.inTree
	} else {
		log.Warn("cannot find parent block", "hash", current.block.Hash(), "number", current.block.NumberU64())
	}

	cbft.buildChildNode(current)
}

func (cbft *Cbft) buildChildNode(current *BlockExt) {
	children := cbft.findChildren(current)
	if len(children) > 0 {
		current.Children = children
		for _, child := range children {
			//child should catch up with current
			child.parent = current
		}
		cbft.setDescendantInTree(current)
	}
}

func (cbft *Cbft) setDescendantInTree(child *BlockExt) {
	log.Debug("set descendant inTree attribute", "hash", child.block.Hash(), "number", child.Number)
	for _, grandchild := range child.Children {
		grandchild.inTree = child.inTree
		cbft.setDescendantInTree(grandchild)
	}
}

// removeBadBlock removes bad block executed error from the tree structure and cbft.blockExtMap.
func (cbft *Cbft) removeBadBlock(badBlock *BlockExt) {
	tailorTree(badBlock)
	cbft.removeByTailored(badBlock)
	/*for _, child := range badBlock.children {
		child.parent = nil
	}
	cbft.blockExtMap.Delete(badBlock.block.Hash())*/
	//delete(cbft.blockExtMap, badBlock.block.Hash())
}

func (cbft *Cbft) removeByTailored(badBlock *BlockExt) {
	if len(badBlock.Children) > 0 {
		for _, child := range badBlock.Children {
			cbft.removeByTailored(child)
			cbft.blockExtMap.Delete(child.block.Hash())
		}
	} else {
		cbft.blockExtMap.Delete(badBlock.block.Hash())
	}
}

// signReceiver handles the received block signature
func (cbft *Cbft) signReceiver(sig *cbfttypes.BlockSignature) error {
	log.Debug("=== call signReceiver() ===",
		"hash", sig.Hash,
		"number", sig.Number.Uint64(),
		"highestLogicalHash", cbft.getHighestLogical().block.Hash(),
		"highestLogicalNumber", cbft.getHighestLogical().Number,
		"highestConfirmedHash", cbft.getHighestConfirmed().block.Hash(),
		"highestConfirmedNumber", cbft.getHighestConfirmed().Number,
		"rootIrreversibleHash", cbft.getRootIrreversible().block.Hash(),
		"rootIrreversibleNumber", cbft.getRootIrreversible().Number)

	if sig.Number.Uint64() <= cbft.getRootIrreversible().Number {
		log.Warn("block sign is too late")
		return nil
	}

	current := cbft.findBlockExt(sig.Hash)
	if current == nil {
		log.Warn("have not received the corresponding block")
		//the block is nil
		current = NewEmptyExt(sig.Number.Uint64())
		current.inTree = false
		current.isExecuted = false
		current.isSigned = false
		current.isConfirmed = false

		cbft.saveBlockExt(sig.Hash, current)
	}

	cbft.collectSign(current, sig.Signature)

	var hashLog interface{}
	if current.block != nil {
		hashLog = current.block.Hash()
	} else {
		hashLog = "sign received before block,"
	}

	log.Debug("count signatures",

		"hash", hashLog,
		"number", current.Number,
		"signCount", len(current.signs),
		"inTree", current.inTree,
		"isExecuted", current.isExecuted,
		"isConfirmed", current.isConfirmed,
		"isSigned", current.isSigned)

	if current.inTree && current.isConfirmed && cbft.getRootIrreversible().isAncestor(current){
		//the current is new highestConfirmed on the same logical path
		if current.Number > cbft.getHighestConfirmed().Number && cbft.getHighestConfirmed().isAncestor(current) {
			cbft.highestConfirmed.Store(current)
			newHighestLogical := cbft.findHighestLogical(current)
			cbft.setHighestLogical(newHighestLogical)
		} else if current.Number < cbft.getHighestConfirmed().Number && !current.isAncestor(cbft.getHighestConfirmed()) {
			//only this case may cause a new fork
			cbft.checkFork(current)
		}
		cbft.flushReadyBlock()
	}

	log.Debug("=== end of signReceiver()  ===",
		"hash", hashLog,
		"number", current.Number,
		"highestLogicalHash", cbft.getHighestLogical().block.Hash(),
		"highestLogicalNumber", cbft.getHighestLogical().Number,
		"highestConfirmedHash", cbft.getHighestConfirmed().block.Hash(),
		"highestConfirmedNumber", cbft.getHighestConfirmed().Number,
		"rootIrreversibleHash", cbft.getRootIrreversible().block.Hash(),
		"rootIrreversibleNumber", cbft.getRootIrreversible().Number)
	return nil
}

//blockReceiver handles the new block
func (cbft *Cbft) blockReceiver(tmp *BlockExt) error {
	block := tmp.block
	rcvTime := tmp.rcvTime
	log.Debug("=== call blockReceiver() ===",
		"hash", block.Hash(),
		"number", block.NumberU64(),
		"parentHash", block.ParentHash(),
		"ReceiptHash", block.ReceiptHash(),
		"highestLogicalHash", cbft.getHighestLogical().block.Hash(),
		"highestLogicalNumber", cbft.getHighestLogical().Number,
		"highestConfirmedHash", cbft.getHighestConfirmed().block.Hash(),
		"highestConfirmedNumber", cbft.getHighestConfirmed().Number,
		"rootIrreversibleHash", cbft.getRootIrreversible().block.Hash(),
		"rootIrreversibleNumber", cbft.getRootIrreversible().Number)

	consensusNodes := cbft.ConsensusNodes(new(big.Int).Sub(block.Number(), common.Big1), block.ParentHash(), block.Number())
	if consensusNodes != nil && len(consensusNodes) == 1 && cbft.config.NodeID==consensusNodes[0] {
		log.Debug("single node Mode")
		cbft.flushReadyBlock()
		return nil
	}

	if block.NumberU64() <= 0 {
		return errGenesisBlock
	}

	if block.NumberU64() <= cbft.getRootIrreversible().Number {
		return lateBlock
	}

	//recover the producer's NodeID
	producerID, sign, err := ecrecover(block.Header())
	if err != nil {
		return err
	}

	// curTime := toMilliseconds(time.Now())

	// TODO
	//isLegal := cbft.isLegal(rcvTime, producerID)
	blockNumber := block.Number()
	parentNumber := new(big.Int).Sub(blockNumber, common.Big1)
	isLegal := cbft.isLegal(rcvTime, parentNumber, block.ParentHash(), blockNumber, producerID)
	log.Debug("check if block is legal",
		"result", isLegal,
		"hash", block.Hash(),
		"number", block.NumberU64(),
		"parentHash", block.ParentHash(),
		"rcvTime", rcvTime,
		"producerID", producerID)
	if !isLegal {
		return errIllegalBlock
	}
	//to check if there's a existing blockExt for received block
	//sometime we'll receive the block's sign before the block self.
	blockExt := cbft.findBlockExt(block.Hash())
	if blockExt == nil {
		blockExt = tmp
		cbft.saveBlockExt(blockExt.block.Hash(), blockExt)
	} else if blockExt.block == nil {
		//received its sign before.
		blockExt.block = block
		blockExt.rcvTime = rcvTime
	} else {
		return errDuplicatedBlock
	}

	//make tree node
	cbft.buildIntoTree(blockExt)

	//collect the block's sign of producer
	cbft.collectSign(blockExt, common.NewBlockConfirmSign(sign))

	log.Debug("count signatures",

		"hash", blockExt.block.Hash(),
		"number", blockExt.Number,
		"signCount", len(blockExt.signs),
		"inTree", blockExt.inTree,
		"isExecuted", blockExt.isExecuted,
		"isConfirmed", blockExt.isConfirmed,
		"isSigned", blockExt.isSigned)

	if blockExt.inTree {
		if err := cbft.executeBlockAndDescendant(blockExt, blockExt.parent); err != nil {
			return err
		}

		// TODO
		//inTurn := cbft.inTurnVerify(blockExt.rcvTime, producerID)
		blockNumber := block.Number()
		parentNumber := new(big.Int).Sub(blockNumber, common.Big1)
		inTurn := cbft.inTurnVerify(parentNumber, block.ParentHash(), blockNumber, blockExt.rcvTime, producerID)
		if !inTurn {
			log.Warn("not in turn",
				"RoutineID", common.CurrentGoRoutineID(),
				"hash", block.Hash(),
				"number", block.NumberU64(),
				"parentHash", block.ParentHash(),
				"curTime", rcvTime,
				"producerID", producerID)
		}

		blockExt.inTurn = inTurn

		//flowControl := flowControl.control(producerID, curTime)
		flowControl := true
		highestConfirmedIsAncestor := cbft.getHighestConfirmed().isAncestor(blockExt)

		isLogical := inTurn && flowControl && highestConfirmedIsAncestor

		log.Debug("check if block is logical", "result", isLogical, "hash", blockExt.block.Hash(), "number", blockExt.Number, "inTurn", inTurn, "flowControl", flowControl, "highestConfirmedIsAncestor", highestConfirmedIsAncestor)

		if isLogical {
			cbft.signLogicalAndDescendant(blockExt)

			newHighestConfirmed := cbft.findLastClosestConfirmedIncludingSelf(cbft.getHighestConfirmed())
			if newHighestConfirmed != nil {
				cbft.highestConfirmed.Store(newHighestConfirmed)
			}

			//rearrange logical path
			newHighestLogical := cbft.findHighestLogical(cbft.getHighestConfirmed())
			if newHighestLogical != nil {
				cbft.setHighestLogical(newHighestLogical)
			}

		} else {
			closestConfirmed := cbft.findClosestConfirmedIncludingSelf(blockExt)

			//if closestConfirmed != nil && closestConfirmed.number < cbft.highestConfirmed.number && !closestConfirmed.isAncestor(cbft.highestConfirmed){
			if closestConfirmed != nil && closestConfirmed.Number < cbft.getHighestConfirmed().Number {
				//only this case may cause a new fork
				cbft.checkFork(closestConfirmed)
			}
		}

		cbft.flushReadyBlock()
	}
	log.Debug("=== end of blockReceiver() ===",
		"hash", block.Hash(),
		"number", block.NumberU64(),
		"parentHash", block.ParentHash(),
		"ReceiptHash", block.ReceiptHash(),
		"highestLogicalHash", cbft.getHighestLogical().block.Hash(),
		"highestLogicalNumber", cbft.getHighestLogical().Number,
		"highestConfirmedHash", cbft.getHighestConfirmed().block.Hash(),
		"highestConfirmedNumber", cbft.getHighestConfirmed().Number,
		"rootIrreversibleHash", cbft.getRootIrreversible().block.Hash(),
		"rootIrreversibleNumber", cbft.getRootIrreversible().Number)
	return nil
}

// forked returns the blocks forked from original branch
// original[0] == newFork[0] == cbft.rootIrreversible, len(origPath) > len(newPath)
func (cbft *Cbft) forked(origPath []*BlockExt, newPath []*BlockExt) (oldTress, newTress []*BlockExt) {
	for i := 0; i < len(newPath); i++ {
		if newPath[i].block.Hash() != origPath[i].block.Hash() {
			return origPath[i:], newPath[i:]
		}
	}
	return nil, nil
}

func extraBlocks(exts []*BlockExt) []*types.Block {
	blocks := make([]*types.Block, len(exts))
	for idx, ext := range exts {
		blocks[idx] = ext.block
	}
	return blocks
}

// checkFork checks if the logical path is changed cause the newConfirmed, if changed, this is a new fork.
func (cbft *Cbft) checkFork(newConfirmed *BlockExt) {
	newHighestConfirmed := cbft.findLastClosestConfirmedIncludingSelf(newConfirmed)
	if newHighestConfirmed != nil && cbft.getRootIrreversible().isAncestor(newHighestConfirmed) && newHighestConfirmed.block.Hash() != cbft.getHighestConfirmed().block.Hash() {
		//forked
		newHighestLogical := cbft.findHighestLogical(newHighestConfirmed)
		newPath := cbft.backTrackBlocks(newHighestLogical, cbft.getRootIrreversible(), true)

		origPath := cbft.backTrackBlocks(cbft.getHighestLogical(), cbft.getRootIrreversible(), true)

		oldTress, newTress := cbft.forked(origPath, newPath)

		if len(newTress) > 0 {

			//fork
			log.Warn("the block chain in memory forked", "newHighestConfirmedHash", newHighestConfirmed.block.Hash(), "newHighestConfirmedNumber", newHighestConfirmed.Number)

			cbft.txPool.ForkedReset(extraBlocks(oldTress), extraBlocks(newTress))

			//forkFrom to lower block
			cbft.highestConfirmed.Store(newHighestConfirmed)
			cbft.setHighestLogical(newHighestLogical)
		}
	}
}

// flushReadyBlock finds ready blocks and flush them to chain
func (cbft *Cbft) flushReadyBlock() bool {
	log.Debug("check if there's any block ready to flush to chain", "highestConfirmedNumber", cbft.getHighestConfirmed().Number, "rootIrreversibleNumber", cbft.getRootIrreversible().Number)

	fallCount := int(cbft.getHighestConfirmed().Number - cbft.getRootIrreversible().Number)
	var newRoot *BlockExt
	if fallCount == 1 && cbft.getRootIrreversible().isParent(cbft.getHighestConfirmed().block) {
		cbft.storeBlocks([]*BlockExt{cbft.getHighestConfirmed()})
		newRoot = cbft.getHighestConfirmed()
	} else if fallCount > windowSize {
		//find the completed path from root to highest logical
		logicalBlocks := cbft.backTrackBlocks(cbft.getHighestConfirmed(), cbft.getRootIrreversible(), false)
		total := len(logicalBlocks)
		toFlushs := logicalBlocks[:total-windowSize]

		logicalBlocks = logicalBlocks[total-windowSize:]

		for _, confirmed := range logicalBlocks {
			if confirmed.isConfirmed {
				toFlushs = append(toFlushs, confirmed)
			} else {
				break
			}
		}

		cbft.storeBlocks(toFlushs)

		for _, confirmed := range toFlushs {
			log.Debug("blocks should be flushed to chain  ", "hash", confirmed.block.Hash(), "number", confirmed.Number)
		}

		newRoot = toFlushs[len(toFlushs)-1]
	}
	if newRoot != nil {
		// blocks[0] == cbft.rootIrreversible
		oldRoot := cbft.getRootIrreversible()
		log.Debug("blockExt tree reorged, root info", "origHash", oldRoot.block.Hash(), "origNumber", oldRoot.Number, "newHash", newRoot.block.Hash(), "newNumber", newRoot.Number)
		//cut off old tree from new root,
		tailorTree(newRoot)

		//set the new root as cbft.rootIrreversible
		cbft.rootIrreversible.Store(newRoot)

		//remove all blocks referenced in old tree after being cut off
		cbft.cleanByTailoredTree(oldRoot)

		//remove all other blocks those their numbers are too low
		cbft.cleanByNumber(cbft.getRootIrreversible().Number)
		return true
	}
	return false

	/*if exceededCount := cbft.highestConfirmed.number - cbft.rootIrreversible.number; exceededCount > 0 {
		//find the completed path from root to highest logical
		logicalBlocks := cbft.backTrackBlocks(cbft.highestConfirmed, cbft.rootIrreversible, false)

		total := len(logicalBlocks)

		var newRoot *BlockExt

		if total > 20 {
			forced := logicalBlocks[:total-20]
			log.Warn("force to flush blocks to chain", "blockCount", len(forced))

			cbft.storeBlocks(forced)

			newRoot = forced[len(forced)-1]
			logicalBlocks = logicalBlocks[total-20:]
		}

		count := 0
		for _, confirmed := range logicalBlocks {
			if confirmed.isConfirmed {
				newRoot = confirmed
				log.Debug("find confirmed block that can be flushed to chain  ", "hash", newRoot.block.Hash(), "number", newRoot.number)
				count++
			} else {
				break
			}
		}
		if count > 0 {
			cbft.storeBlocks(logicalBlocks[:count])
		}
		if newRoot != nil {
			// blocks[0] == cbft.rootIrreversible
			oldRoot := cbft.rootIrreversible
			log.Debug("oldRoot", "hash", oldRoot.block.Hash(), "number", oldRoot.number)
			log.Debug("newRoot", "hash", newRoot.block.Hash(), "number", newRoot.number)
			//cut off old tree from new root,
			tailorTree(newRoot)

			//set the new root as cbft.rootIrreversible
			cbft.rootIrreversible = newRoot

			//remove all blocks referenced in old tree after being cut off
			cbft.cleanByTailoredTree(oldRoot)

			//remove all other blocks those their numbers are too low
			cbft.cleanByNumber(cbft.rootIrreversible.number)
		}
	}*/
}

// tailorTree tailors the old tree from new root
func tailorTree(newRoot *BlockExt) {
	if newRoot.parent != nil && newRoot.parent.Children != nil {
		for i := 0; i < len(newRoot.parent.Children); i++ {
			//remove newRoot from its parent's children list
			if newRoot.parent.Children[i].block.Hash() == newRoot.block.Hash() {
				newRoot.parent.Children = append(newRoot.parent.Children[:i], newRoot.parent.Children[i+1:]...)
				break
			}
		}
		newRoot.parent = nil
	}
}

// cleanByTailoredTree removes all blocks in the tree which has been tailored.
func (cbft *Cbft) cleanByTailoredTree(root *BlockExt) {
	log.Trace("call cleanByTailoredTree()", "rootHash", root.block.Hash(), "rootNumber", root.block.NumberU64())
	if len(root.Children) > 0 {
		for _, child := range root.Children {
			cbft.cleanByTailoredTree(child)
			log.Debug("remove block in memory", "hash", root.block.Hash(), "number", root.block.NumberU64())
			cbft.blockExtMap.Delete(root.block.Hash())
			//delete(cbft.blockExtMap, root.block.Hash())
			//delete(cbft.signedSet, root.block.NumberU64())
			cbft.signedSet.Delete(root.block.NumberU64())
		}
	} else {
		log.Debug("remove block in memory", "hash", root.block.Hash(), "number", root.block.NumberU64())
		cbft.blockExtMap.Delete(root.block.Hash())
		//delete(cbft.blockExtMap, root.block.Hash())
	}
}

// cleanByNumber removes all blocks lower than upperLimit in BlockExtMap.
func (cbft *Cbft) cleanByNumber(upperLimit uint64) {
	log.Trace("call cleanByNumber()", "upperLimit", upperLimit)

	f := func(k, v interface{}) bool {
		// The input and input types of this function are fixed and cannot be modified.
		// You can write your own code in the body of the function, call k, v in the map
		hash := k.(common.Hash)
		ext := v.(*BlockExt)
		if ext.Number < upperLimit {
			log.Debug("remove block in memory", "hash", hash, "number", ext.Number)
			cbft.blockExtMap.Delete(hash)
		}
		return true
	}
	cbft.blockExtMap.Range(f)

	/*for number, _ := range cbft.signedSet {
		if number < upperLimit {
			delete(cbft.signedSet, number)
		}
	}*/
	f2 := func(k, v interface{}) bool {
		number := k.(uint64)
		if number < upperLimit {
			cbft.signedSet.Delete(number)
		}
		return true
	}

	cbft.signedSet.Range(f2)
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	log.Trace("call Author()", "hash", header.Hash(), "number", header.Number.Uint64())
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	log.Trace("call VerifyHeader()", "hash", header.Hash(), "number", header.Number.Uint64(), "seal", seal)

	if header.Number == nil {
		return errUnknownBlock
	}

	if len(header.Extra) < extraSeal {
		return errMissingSignature
	}
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (cbft *Cbft) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	log.Trace("call VerifyHeaders()", "Headers count", len(headers))

	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for _, header := range headers {
			err := cbft.VerifyHeader(chain, header, false)

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (cbft *Cbft) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return nil
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (cbft *Cbft) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	log.Trace("call VerifySeal()", "hash", header.Hash(), "number", header.Number.String())

	return cbft.verifySeal(chain, header, nil)
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (cbft *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	log.Debug("call Prepare()", "hash", header.Hash(), "number", header.Number.Uint64())

	if cbft.getHighestLogical().block == nil || header.ParentHash != cbft.getHighestLogical().block.Hash() || header.Number.Uint64()-1 != cbft.getHighestLogical().block.NumberU64() {
		return consensus.ErrUnknownAncestor
	}

	// header.Extra[0:32] to store block's version info etc. and right pad with 0x00;
	// header.Extra[32:97] to store block's sign of producer, the length of sign is 65.
	if len(header.Extra) < 32 {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, 32-len(header.Extra))...)
	}
	header.Extra = header.Extra[:32]

	//init header.Extra[32: 32+65]
	header.Extra = append(header.Extra, make([]byte, consensus.ExtraSeal)...)
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given, and returns the final block.
func (cbft *Cbft) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	log.Debug("call Finalize()", "RoutineID", common.CurrentGoRoutineID(), "hash", header.Hash(), "number", header.Number.Uint64(), "txs", len(txs), "receipts", len(receipts), " extra: ", hexutil.Encode(header.Extra))
	cbft.accumulateRewards(chain.Config(), state, header)
	cbft.IncreaseRewardPool(state, header.Number)

	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	return types.NewBlock(header, txs, nil, receipts), nil
}

// to sign the block, and store the sign to header.Extra[32:97], send the sign to chanel to broadcast to other consensus nodes
func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) error {
	log.Debug("call Seal()", "number", block.NumberU64(), "parentHash", block.ParentHash())
	/*cbft.lock.Lock()
	defer cbft.lock.Unlock()*/

	header := block.Header()
	number := block.NumberU64()

	if number == 0 {
		return errUnknownBlock
	}

	if !cbft.getHighestLogical().isParent(block) {
		log.Error("Futile block cause highest logical block changed", "parentHash", block.ParentHash())
		return errFutileBlock
	}

	// sign the seal hash
	sign, err := cbft.signFn(header.SealHash().Bytes())
	if err != nil {
		return err
	}

	// store the sign into the header.Extra[32:97]
	copy(header.Extra[32:97], sign[:])

	sealedBlock := block.WithSeal(header)

	current := NewBlockExt(sealedBlock, sealedBlock.NumberU64())

	//this block is produced by local node, so need not execute in cbft.
	current.inTree = true
	current.isExecuted = true
	current.isSigned = true

	//save the block to cbft.blockExtMap
	cbft.saveBlockExt(sealedBlock.Hash(), current)

	//collect the sign
	cbft.collectSign(current, common.NewBlockConfirmSign(sign))

	//log this signed block's number
	//cbft.signedSet[sealedBlock.NumberU64()] = struct{}{}
	cbft.signedSet.Store(sealedBlock.NumberU64(), struct{}{})

	//build tree node
	cbft.buildIntoTree(current)

	log.Debug("seal complete", "hash", sealedBlock.Hash(), "number", block.NumberU64())

	// SetNodeCache
	blockNumber := current.block.Number()
	parentNumber := new(big.Int).Sub(blockNumber, common.Big1)
	sealhash := cbft.SealHash(current.block.Header())
	state := cbft.blockChainCache.ReadStateDB(sealhash)
	log.Debug("setNodeCache", "parentNumber", parentNumber, "parentHash", current.block.ParentHash(), "blockNumber", blockNumber, "blockHash", current.block.Hash())
	if state != nil {

		genesis := cbft.blockChain.Genesis()

		cbft.ppos.SetNodeCache(state, genesis.Number(), parentNumber, blockNumber, genesis.Hash(), block.ParentHash(), current.block.Hash())
		blockInterval := new(big.Int).Sub(current.block.Number(), cbft.blockChain.CurrentBlock().Number())
		cbft.ppos.Submit2Cache(state, blockNumber, blockInterval, current.block.Hash())
		root := state.IntermediateRoot(cbft.blockChain.Config().IsEIP158(current.block.Number()))
		log.Debug("Consensus packaged, Finally, after Submit2Cache", "blockNumber", current.block.NumberU64(), "blockHash", current.block.Hash().Hex(), "block.root", current.block.Root().Hex(), "Realt-time state.root", root.Hex())
	} else {
		log.Error("setNodeCache error")
	}

	consensusNodes := cbft.ConsensusNodes(parentNumber, current.block.ParentHash(), blockNumber)

	if consensusNodes != nil && len(consensusNodes) == 1 && cbft.config.NodeID==consensusNodes[0]{
		log.Debug("single node Mode")
		//only one consensus node, so, each block is highestConfirmed. (lock is needless)
		current.rcvTime = toMilliseconds(time.Now())
		current.inTree = true
		current.isExecuted = true
		current.isSigned = true
		current.isConfirmed = true

		cbft.setHighestLogical(current)
		cbft.highestConfirmed.Store(current)

		cbft.dataReceiveCh <- current

		log.Debug("reset TxPool after block sealed", "hash", current.block.Hash(), "number", current.Number)
		cbft.txPool.Reset(current.block)
		return nil
	}

	//reset cbft.highestLogicalBlockExt cause this block is produced by myself
	cbft.setHighestLogical(current)

	go func() {
		select {
		case <-stopCh:
			return
		case sealResultCh <- sealedBlock:
		default:
			log.Warn("Sealing result is not ready by miner", "sealHash", header.SealHash())
		}
	}()

	log.Debug("reset TxPool after block sealed", "hash", current.block.Hash(), "number", current.Number)
	cbft.txPool.Reset(current.block)
	return nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (b *Cbft) SealHash(header *types.Header) common.Hash {
	log.Debug("call SealHash()", "hash", header.Hash(), "number", header.Number.Uint64())
	return header.SealHash()
}

// Close implements consensus.Engine. It's a noop for cbft as there is are no background threads.
func (cbft *Cbft) Close() error {
	log.Trace("call Close()")
	cbft.closeOnce.Do(func() {
		// Short circuit if the exit channel is not allocated.
		if cbft.exitCh == nil {
			return
		}
		cbft.exitCh <- struct{}{}
		close(cbft.exitCh)
	})
	return nil
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (cbft *Cbft) APIs(chain consensus.ChainReader) []rpc.API {
	log.Trace("call APIs()")

	return []rpc.API{{
		Namespace: "cbft",
		Version:   "1.0",
		Service:   &API{chain: chain, cbft: cbft},
		Public:    false,
	}}
}

// OnBlockSignature is called by by protocol handler when it received a new block signature by P2P.
func (cbft *Cbft) OnBlockSignature(chain consensus.ChainReader, nodeID discover.NodeID, rcvSign *cbfttypes.BlockSignature) error {
	log.Debug("call OnBlockSignature()", "hash", rcvSign.Hash, "number", rcvSign.Number, "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]), "signHash", rcvSign.SignHash, "cbft.dataReceiveCh.len", len(cbft.dataReceiveCh))
	ok, err := verifySign(nodeID, rcvSign.SignHash, rcvSign.Signature[:])
	if err != nil {
		log.Error("verify sign error", "errors", err)
		return err
	}

	if !ok {
		log.Error("unauthorized signer")
		return errUnauthorizedSigner
	}
	cbft.dataReceiveCh <- rcvSign
	return nil
}

// OnNewBlock is called by protocol handler when it received a new block by P2P.
func (cbft *Cbft) OnNewBlock(chain consensus.ChainReader, rcvBlock *types.Block) error {
	log.Debug("call OnNewBlock()", "hash", rcvBlock.Hash(), "number", rcvBlock.NumberU64(), "ParentHash", rcvBlock.ParentHash(), "cbft.dataReceiveCh.len", len(cbft.dataReceiveCh))
	tmp := NewBlockExt(rcvBlock, rcvBlock.NumberU64())
	tmp.rcvTime = toMilliseconds(time.Now())
	tmp.inTree = false
	tmp.isExecuted = false
	tmp.isSigned = false
	tmp.isConfirmed = false

	cbft.dataReceiveCh <- tmp
	return nil
}

// OnPong is called by protocol handler when it received a new Pong message by P2P.
func (cbft *Cbft) OnPong(nodeID discover.NodeID, netLatency int64) error {
	log.Trace("call OnPong()", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]), "netLatency", netLatency)

	cbft.netLatencyLock.Lock()
	defer cbft.netLatencyLock.Unlock()

	if netLatency >= maxPingLatency {
		return nil
	}

	latencyList, exist := cbft.netLatencyMap[nodeID]
	if !exist {
		cbft.netLatencyMap[nodeID] = list.New()
		cbft.netLatencyMap[nodeID].PushBack(netLatency)
	} else {
		if latencyList.Len() > 5 {
			e := latencyList.Front()
			cbft.netLatencyMap[nodeID].Remove(e)
		}
		cbft.netLatencyMap[nodeID].PushBack(netLatency)
	}
	return nil
}

// RemovePeer remove the net latency info from netLatencyMap.
func (cbft *Cbft) RemovePeer(nodeID discover.NodeID){
	log.Trace("call RemovePeer()", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]))
	cbft.netLatencyLock.Lock()
	defer cbft.netLatencyLock.Unlock()
	delete(cbft.netLatencyMap, nodeID)
}

// avgLatency statistics the net latency between local and other peers.
func (cbft *Cbft) avgLatency(nodeID discover.NodeID) int64 {
	if latencyList, exist := cbft.netLatencyMap[nodeID]; exist {
		sum := int64(0)
		counts := int64(0)
		for e := latencyList.Front(); e != nil; e = e.Next() {
			if latency, ok := e.Value.(int64); ok {
				counts++
				sum += latency
			}
		}
		if counts > 0 {
			return sum / counts
		}
	}
	return cbft.config.MaxLatency
}

// HighestLogicalBlock returns the cbft.highestLogical.block.
func (cbft *Cbft) HighestLogicalBlock() *types.Block {
	log.Debug("call HighestLogicalBlock() ...")
	if cbft.getHighestLogical() == nil {
		return nil
	} else {
		return cbft.getHighestLogical().block
	}
}

// HighestConfirmedBlock returns the cbft.highestConfirmed.block.
func (cbft *Cbft) HighestConfirmedBlock() *types.Block {
	log.Debug("call HighestConfirmedBlock() ...")
	if cbft.getHighestConfirmed() == nil {
		return nil
	} else {
		return cbft.getHighestConfirmed().block
	}
}

// GetBlock returns the block in blockExtMap.
func (cbft *Cbft) GetBlock(hash common.Hash, number uint64) *types.Block {
	/*if ext := cbft.findBlockExt(hash); ext!=nil {
		if ext.block != nil && ext.number == number {
			return ext.block
		}
	}*/

	if ext, ok := cbft.blockExtMap.Load(hash); ok {
		if ext.(*BlockExt).block != nil && ext.(*BlockExt).Number == number {
			return ext.(*BlockExt).block
		}
	}
	return nil
}

// IsSignedBySelf returns if the block is signed by local.
func IsSignedBySelf(sealHash common.Hash, signature []byte) bool {
	ok, err := verifySign(cbft.config.NodeID, sealHash, signature)
	if err != nil {
		log.Error("verify sign error", "errors", err)
		return false
	}
	return ok
}

// storeBlocks sends the blocks to cbft.cbftResultOutCh, the receiver will write them into chain
func (cbft *Cbft) storeBlocks(blocksToStore []*BlockExt) {
	for _, ext := range blocksToStore {
		cbftResult := &cbfttypes.CbftResult{
			Block:             ext.block,
			BlockConfirmSigns: ext.signs,
		}
		log.Debug("send consensus result to worker", "hash", ext.block.Hash(), "number", ext.block.NumberU64(), "signCount", len(ext.signs))
		cbft.cbftResultOutCh <- cbftResult
	}
}

// inTurn return if it is local's turn to package new block.
//func (cbft *Cbft) inTurn() bool {
//	curTime := toMilliseconds(time.Now())
//	inturn := cbft.calTurn(curTime, cbft.config.NodeID)
//	log.Debug("inTurn", "result", inturn)
//	return inturn
//}

func (cbft *Cbft) inTurn(parentNumber *big.Int, parentHash common.Hash, commitNumber *big.Int) bool {
	curTime := toMilliseconds(time.Now())
	inturn := cbft.calTurn(curTime-300, parentNumber, parentHash, commitNumber, cbft.config.NodeID, current)
	if inturn {
		inturn = cbft.calTurn(curTime+600, parentNumber, parentHash, commitNumber, cbft.config.NodeID, current)
	}
	log.Debug("inTurn", "result", inturn)
	return inturn
}

// inTurnVerify verifies the time is in the time-window of the nodeID to package new block.
//func (cbft *Cbft) inTurnVerify(curTime int64, nodeID discover.NodeID) bool {
//	latency := cbft.avgLatency(nodeID)
//	if latency >= maxAvgLatency {
//		log.Debug("inTurnVerify, return false cause of net latency", "result", false, "latency", latency)
//		return false
//	}
//	inTurnVerify := cbft.calTurn(curTime-latency, nodeID)
//	log.Debug("inTurnVerify", "result", inTurnVerify, "latency", latency)
//	return inTurnVerify
//}
func (cbft *Cbft) inTurnVerify(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, rcvTime int64, nodeID discover.NodeID) bool {
	latency := cbft.avgLatency(nodeID)
	if latency >= maxAvgLatency {
		log.Warn("check if peer's turn to commit block", "result", false, "peerID", nodeID, "high latency ", latency)
		return false
	}
	inTurnVerify := cbft.calTurn(rcvTime-latency, parentNumber, parentHash, blockNumber, nodeID, all)
	log.Debug("check if peer's turn to commit block", "result", inTurnVerify, "peerID", nodeID, "latency", latency)
	return inTurnVerify
}

//shouldKeepIt verifies the time is legal to package new block for the nodeID.
//func (cbft *Cbft) isLegal(rcvTime int64, producerID discover.NodeID) bool {
//	offset := 1000 * (cbft.config.Duration/2 - 1)
//	isLegal := cbft.calTurn(curTime-offset, producerID)
//	if !isLegal {
//		isLegal = cbft.calTurn(curTime+offset, producerID)
//	}
//	return isLegal
//}
func (cbft *Cbft) isLegal(rcvTime int64, parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, producerID discover.NodeID) bool {
	offset := 1000 * (cbft.config.Duration/2 - 1)
	isLegal := cbft.calTurn(rcvTime-offset, parentNumber, parentHash, blockNumber, producerID, all)
	if !isLegal {
		isLegal = cbft.calTurn(rcvTime+offset, parentNumber, parentHash, blockNumber, producerID, all)
	}
	return isLegal
}

//func (cbft *Cbft) calTurn(curTime int64, nodeID discover.NodeID) bool {
//	nodeIdx := cbft.dpos.NodeIndex(nodeID)
//	startEpoch := cbft.dpos.StartTimeOfEpoch() * 1000
//
//	if nodeIdx >= 0 {
//		durationPerNode := cbft.config.Duration * 1000
//		durationPerTurn := durationPerNode * int64(len(cbft.dpos.primaryNodeList))
//
//		min := nodeIdx * (durationPerNode)
//
//		value := (curTime - startEpoch) % durationPerTurn
//
//		max := (nodeIdx + 1) * durationPerNode
//
//		log.Debug("calTurn", "idx", nodeIdx, "min", min, "value", value, "max", max, "curTime", curTime, "startEpoch", startEpoch)
//
//		if value > min && value < max {
//			return true
//		}
//	}
//	return false
//}

func (cbft *Cbft) calTurn(timePoint int64, parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, nodeID discover.NodeID, round int32) bool {

	nodeIdx, consensusNodes := cbft.ppos.BlockProducerIndex(parentNumber, parentHash, blockNumber, nodeID, round)
	startEpoch := cbft.ppos.StartTimeOfEpoch() * 1000

	if nodeIdx >= 0 {
		durationPerNode := cbft.config.Duration * 1000

		//consensusNodes := cbft.ConsensusNodes(parentNumber, parentHash, blockNumber)
		if consensusNodes == nil || len(consensusNodes) <= 0 {
			log.Error("there is no consensus node",  "number", blockNumber)
			return false
		} else if len(consensusNodes) == 1 {
			return true
		}

		durationPerTurn := durationPerNode * int64(len(consensusNodes))

		min := nodeIdx * (durationPerNode)

		value := (timePoint - startEpoch) % durationPerTurn

		max := (nodeIdx + 1) * durationPerNode

		if value > min && value < max {
			return true
		}
	}else{
		log.Debug("local is not a consensus node", "localNode", nodeID.String(), "number", blockNumber)
		for idx, nid := range  consensusNodes{
			log.Debug("consensus node list", "idx", idx, "nodeID", nid.String())
		}
	}
	return false
}

// producer's signature = header.Extra[32:97]
// public key can be recovered from signature, the length of public key is 65,
// the length of NodeID is 64, nodeID = publicKey[1:]
func ecrecover(header *types.Header) (discover.NodeID, []byte, error) {
	var nodeID discover.NodeID
	if len(header.Extra) < extraSeal {
		return nodeID, []byte{}, errMissingSignature
	}
	signature := header.Extra[32:97]
	sealHash := header.SealHash()

	pubkey, err := crypto.Ecrecover(sealHash.Bytes(), signature)
	if err != nil {
		return nodeID, []byte{}, err
	}

	nodeID, err = discover.BytesID(pubkey[1:])
	if err != nil {
		return nodeID, []byte{}, err
	}
	return nodeID, signature, nil
}

// verify sign, check the sign is from the right node.
func verifySign(expectedNodeID discover.NodeID, sealHash common.Hash, signature []byte) (bool, error) {
	pubkey, err := crypto.SigToPub(sealHash.Bytes(), signature)

	if err != nil {
		return false, err
	}

	nodeID := discover.PubkeyID(pubkey)
	if bytes.Equal(nodeID.Bytes(), expectedNodeID.Bytes()) {
		return true, nil
	}
	return false, nil
}

func (cbft *Cbft) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	return nil
}

func (cbft *Cbft) signFn(headerHash []byte) (sign []byte, err error) {
	return crypto.Sign(headerHash, cbft.config.PrivateKey)
}

//func (cbft *Cbft) getThreshold() int {
//	trunc := len(cbft.dpos.primaryNodeList) * 2 / 3
//	return int(trunc + 1)
//}


func (cbft *Cbft) getConnectedConsensusNodes(consensusNodes []discover.NodeID) int {
	cbft.netLatencyLock.RLock()
	defer cbft.netLatencyLock.RUnlock()

	connected :=0
	for _,nodeId :=range consensusNodes {
		if _, exist := cbft.netLatencyMap[nodeId]; exist {
			connected++
		}
	}
	return connected
}


func (cbft *Cbft) getThreshold(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) int {
	consensusNodes := cbft.ConsensusNodes(parentNumber, parentHash, blockNumber)
	if consensusNodes != nil {
		total := len(consensusNodes)
		log.Debug("consensus node quantity", "total", total)
		trunc := total * 2 / 3
		return int(trunc + 1)
	}
	return math.MaxInt16
}

func (cbft *Cbft) calculateThreshold(consensusNodes []discover.NodeID) int {
	if consensusNodes != nil {
		trunc := len(consensusNodes) * 2 / 3
		return int(trunc + 1)
	}
	return math.MaxInt16
}


func toMilliseconds(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func (cbft *Cbft) ShouldSeal(parentNumber *big.Int, parentHash common.Hash, commitNumber *big.Int) bool {
	log.Trace("call ShouldSeal()")

	//consensusNodes := cbft.ConsensusNodes(parentNumber, parentHash, commitNumber)
	//if consensusNodes != nil && len(consensusNodes) == 1 {
	//	return true
	//}

	inturn := cbft.inTurn(parentNumber, parentHash, commitNumber)

	if inturn {
		consensusNodes := cbft.ConsensusNodes(parentNumber, parentHash, commitNumber)
		threshold := cbft.calculateThreshold(consensusNodes)
		connected := cbft.getConnectedConsensusNodes(consensusNodes)

		log.Debug("connected consensus peers(including self)", "count", connected+1, "threshold", threshold)
		if connected < threshold-1 {
			inturn = false
		}
	}
	log.Debug("end of ShouldSeal()", "result", inturn, "number", commitNumber)
	return inturn
}

func (cbft *Cbft) CurrentNodes(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []*discover.Node {
	return cbft.ppos.getCurrentNodes(parentNumber, parentHash, blockNumber)
}

func (cbft *Cbft) IsCurrentNode(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) bool {
	currentNodes := cbft.ppos.getCurrentNodes(parentNumber, parentHash, blockNumber)
	nodeID := cbft.GetOwnNodeID()
	for _, n := range currentNodes {
		if nodeID == n.ID {
			return true
		}
	}
	return false
}

func (cbft *Cbft) ConsensusNodes(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []discover.NodeID {
	return cbft.ppos.consensusNodes(parentNumber, parentHash, blockNumber)
}

// wether nodeID in former or current or next
//func (cbft *Cbft) CheckConsensusNode(nodeID discover.NodeID) (bool, error) {
//	log.Debug("call CheckConsensusNode()", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]))
//	return cbft.ppos.AnyIndex(nodeID) >= 0, nil
//}

// wether nodeID in former or current or next
//func (cbft *Cbft) IsConsensusNode() (bool, error) {
//	log.Debug("call IsConsensusNode()")
//	return cbft.ppos.AnyIndex(cbft.config.NodeID) >= 0, nil
//}

func (cbft *Cbft) Election(state *state.StateDB, parentHash common.Hash, blockNumber *big.Int) ([]*discover.Node, error) {
	return cbft.ppos.Election(state, parentHash, blockNumber)
}

func (cbft *Cbft) Switch(state *state.StateDB, blockNumber *big.Int) bool {
	return cbft.ppos.Switch(state, blockNumber)
}

func (cbft *Cbft) GetWitness(state *state.StateDB, flag int, blockNumber *big.Int) ([]*discover.Node, error) {
	return cbft.ppos.GetWitness(state, flag, blockNumber)
}

func (cbft *Cbft) GetOwnNodeID() discover.NodeID {
	return cbft.config.NodeID
}

func (cbft *Cbft) SetNodeCache(state *state.StateDB, parentNumber, currentNumber *big.Int, parentHash, currentHash common.Hash) error {
	log.Info("cbft SetNodeCache", "parentNumber", parentNumber, "parentHash", parentHash, "currentNumber", currentNumber, "currentHash", currentHash)
	genesis := cbft.blockChain.Genesis()
	return cbft.ppos.SetNodeCache(state, genesis.Number(), parentNumber, currentNumber, genesis.Hash(), parentHash, currentHash)
}

func (cbft *Cbft) Notify(state vm.StateDB, blockNumber *big.Int) error {
	return cbft.ppos.Notify(state, blockNumber)
}

func (cbft *Cbft) StoreHash(state *state.StateDB, blockNumber *big.Int, blockHash common.Hash) {
	cbft.ppos.StoreHash(state, blockNumber, blockHash)
}

func (cbft *Cbft) Submit2Cache(state *state.StateDB, currBlocknumber *big.Int, blockInterval *big.Int, currBlockhash common.Hash) {
	cbft.ppos.Submit2Cache(state, currBlocknumber, blockInterval, currBlockhash)
}

// AccumulateRewards for lucky tickets
// Adjust rewards every 3600*24*365 blocks
func (cbft *Cbft) accumulateRewards(config *params.ChainConfig, state *state.StateDB, header *types.Header) {
	if len(header.Extra) < 64 {
		log.Error("Failed to Call accumulateRewards, header.Extra < 64", "blockNumber", header.Number, "blockHash", header.Hash(), "len(header.Extra):", len(header.Extra), "extra", hexutil.Encode(header.Extra))
	}


	appendEtraFunc := func(txHash common.Hash, isPackageNode bool) {
		if isPackageNode {
			var buffer bytes.Buffer
			buffer.Write(header.Extra)
			buffer.Write(txHash.Bytes())
			header.Extra = buffer.Bytes()
			log.Info("Call accumulateRewards, When After Sets the Maybe lucky ticket into header.Extra", "len(header.Extra):", len(header.Extra), "extra", hexutil.Encode(header.Extra))
		}
	}


	var nodeId discover.NodeID
	var err error
	log.Debug("Call accumulateRewards block header", "len(header.Extra):", len(header.Extra), "extra", hexutil.Encode(header.Extra))

	packageNodeFlag := bytes.Equal(header.Extra[32:97], make([]byte, 65))

	if packageNodeFlag {
		log.Warn("Call accumulateRewards block header extra[32:97] is empty!", "blockNumber", header.Number, "blockHash", header.Hash())
		nodeId = cbft.config.NodeID
	} else {
		if nodeId, _, err = ecrecover(header); err != nil {

			appendEtraFunc((common.Hash{}), packageNodeFlag)

			log.Error("Failed to Call accumulateRewards, ecrecover faile", " err: ", err, "blockNumber", header.Number, "blockHash", header.Hash())
			return
		} else {
			log.Info("Success ecrecover", "blockNumber", header.Number, "blockHash", header.Hash(), " nodeId: ", nodeId.String())
		}
	}

	//log.Info("Call accumulateRewards", "nodeid: ", nodeId.String())
	var can *types.Candidate
	if big.NewInt(0).Cmp(new(big.Int).Rem(header.Number, big.NewInt(common.BaseSwitchWitness))) == 0 {
		can = cbft.ppos.GetWitnessCandidate(state, nodeId, -1, header.Number)
	} else {
		can = cbft.ppos.GetWitnessCandidate(state, nodeId, 0, header.Number)
	}
	if err != nil {

		appendEtraFunc((common.Hash{}), packageNodeFlag)

		log.Error("Failed to Call accumulateRewards, GetCandidate faile ", " err: ", err.Error(),
			"blockNumber", header.Number, "blockHash", header.Hash(), "nodeId", nodeId.String())
		return
	}
	if can == nil {

		appendEtraFunc((common.Hash{}), packageNodeFlag)

		log.Warn("Call accumulateRewards, Witness's can is Empty !!!!!!!!!!!!!!!!!!!!!!!!",
			"blockNumber", header.Number, "blockHash", header.Hash(), "nodeId", nodeId.String())
		return
	}

	// store the lucky ticket into the header.Extra[97:129]
	appendEtraFunc(can.TxHash, packageNodeFlag)

	//Calculate current block rewards
	var blockReward *big.Int
	preYearNumber := new(big.Int).Sub(header.Number, common.YearBlocks)
	yearReward := new(big.Int).Set(common.FirstYearReward)
	if preYearNumber.Cmp(common.YearBlocks) > 0 { // otherwise is 0 year and 1 year block reward
		yearReward = new(big.Int).Sub(GetAmount(header.Number), GetAmount(preYearNumber))
	}
	blockReward = new(big.Int).Div(yearReward, common.YearBlocks)

	nodeReward := blockReward
	//log.Info("Call accumulateRewards, GetTicket ", "TicketId: ", can.TicketId.Hex())
	if can.TOwner != (common.Address{}) {
		nodeReward = new(big.Int).Div(new(big.Int).Mul(blockReward, new(big.Int).SetUint64(uint64(can.Fee))), common.FeeBase)
		ticketReward := new(big.Int).Sub(blockReward, nodeReward)

		//log.Info("Call accumulateRewards, Rewards detail", "blockReward: ", blockReward, "nodeReward: ", nodeReward, "ticketReward: ", ticketReward)
		state.AddBalance(can.TOwner, ticketReward)
		log.Info("Ticket accumulateRewards", "txHash", can.TxHash.Hex(), "ticketOwner", can.TOwner.Hex(), "ticketReward", ticketReward, "ticketOwnerBalance", state.GetBalance(can.TOwner))
	}
	state.AddBalance(header.Coinbase, nodeReward)
	state.SubBalance(common.RewardPoolAddr, blockReward)

	log.Info("Call accumulateRewards SUCCESS !! ", "blockNumber", header.Number, "blockHash", header.Hash(),
		"nodeId", nodeId.String(), "ticketId", can.TxHash.Hex(), " yearReward: ", yearReward, " blockReward:", blockReward,
		" nodeReward: ", nodeReward, " RewardPoolAddr address: ", common.RewardPoolAddr.Hex(),
		" balance: ", state.GetBalance(common.RewardPoolAddr), " Fee: ", can.Fee, " Coinbase address: ", header.Coinbase.Hex(),
		" balance: ", state.GetBalance(header.Coinbase))
}

func (cbft *Cbft) IncreaseRewardPool(state *state.StateDB, number *big.Int) {
	//add balance to reward pool
	if new(big.Int).Rem(number, common.YearBlocks).Cmp(big.NewInt(0)) == 0 {
		num := GetAmount(number)
		log.Info("Call IncreaseRewardPool SUCCESS !! ", "addr", common.RewardPoolAddr.Hex(), "num", num.String())
		state.AddBalance(common.RewardPoolAddr, num)
	}
}

func GetAmount(number *big.Int) *big.Int {
	cycle := new(big.Int).Div(number, common.YearBlocks)
	rate := math2.BigPow(common.Rate.Int64(), cycle.Int64())
	base := math2.BigPow(common.Base.Int64(), cycle.Int64())
	//fmt.Println("number: ", number, " cycle: ", cycle, " rate: ", rate, " base: ", base)
	yearAmount := new(big.Int).Mul(common.InitAmount, rate)
	ret := new(big.Int).Div(yearAmount, base)
	return ret
}

func (cbft *Cbft) FindTransaction(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return cbft.blockChain.GetTransactionByHash(txHash)
}

func (cbft *Cbft) GetHeader(blockHash common.Hash, blockNumber uint64) *types.Header {
	return cbft.blockChain.GetHeader(blockHash, blockNumber)
}

func (cbft *Cbft) GetBody(blockNumber uint64) *types.Body {
	return cbft.blockChain.GetBodyByNumber(blockNumber)
}

func (cbft *Cbft) GetNewStateDB(root common.Hash, blockNumber *big.Int, blockHash common.Hash) (*state.StateDB, error) {
	return cbft.blockChain.GetNewStateDB(root, blockNumber, blockHash)
}
