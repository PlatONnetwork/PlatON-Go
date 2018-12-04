// Package bft implements the BFT consensus engine.
package cbft

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/consensus"
	"Platon-go/core"
	"Platon-go/core/cbfttypes"
	"Platon-go/core/ppos"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/crypto"
	"Platon-go/crypto/sha3"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"Platon-go/rlp"
	"Platon-go/rpc"
	"bytes"
	"container/list"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"sync"
	"time"
)

var (
	errSign                = errors.New("sign error")
	errUnauthorizedSigner  = errors.New("unauthorized signer")
	errIllegalBlock        = errors.New("illegal block")
	errDuplicatedBlock     = errors.New("duplicated block")
	errBlockNumber         = errors.New("error block number")
	errUnknownBlock        = errors.New("unknown block")
	errFutileBlock         = errors.New("futile block")
	errGenesisBlock        = errors.New("cannot handle genesis block")
	errHighestLogicalBlock = errors.New("cannot find a logical block")
	errListConfirmedBlocks = errors.New("list confirmed blocks error")
	errMissingSignature    = errors.New("extra-data 65 byte signature suffix missing")
	extraSeal              = 65
	windowSize             = uint64(20)
	periodMargin           = uint64(20) //meaning 20%

	//一次Ping/Pong的最大网络延迟，毫秒
	maxPingLatency = int64(5000)

	//最大平均网络延迟，毫秒
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
	exitCh                chan chan error

	txPool *core.TxPool

	//todo:（先log,再处理）
	blockExtMap map[common.Hash]*BlockExt //store all received blocks and signs

	dataReceiveCh    chan interface{} //a channel to receive block signature
	blockChain       *core.BlockChain //the block chain
	highestLogical   *BlockExt        //for next block
	highestConfirmed *BlockExt        //highest highestConfirmed block

	//todo:（先log,再处理）
	signedSet map[uint64]struct{} //all block numbers signed by local node
	lock      sync.RWMutex
	// modify by platon remove consensusCache
	//consensusCache *Cache //cache for cbft consensus

	netLatencyMap map[discover.NodeID]*list.List
}

var cbft *Cbft

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, blockSignatureCh chan *cbfttypes.BlockSignature, cbftResultCh chan *cbfttypes.CbftResult, highestLogicalBlockCh chan *types.Block) *Cbft {

	pposm.PrintObject("获取ppos config：", *config)
	_ppos := newPpos(/*config.InitialNodes, */config)

	cbft = &Cbft{
		config:                config,
		ppos:                  _ppos,
		rotating:              newRotating(_ppos, config.Duration),
		blockSignOutCh:        blockSignatureCh,
		cbftResultOutCh:       cbftResultCh,
		highestLogicalBlockCh: highestLogicalBlockCh,

		blockExtMap:   make(map[common.Hash]*BlockExt),
		signedSet:     make(map[uint64]struct{}),
		dataReceiveCh: make(chan interface{}, 250),
		netLatencyMap: make(map[discover.NodeID]*list.List),
	}

	flowControl = NewFlowControl()

	//启动协程处理新收到的块/签名
	go cbft.dataReceiverGoroutine()

	return cbft
}

//the extension for Block
type BlockExt struct {
	block       *types.Block
	isLinked    bool
	isSigned    bool
	isStored    bool
	isConfirmed bool
	number      uint64
	signs       []*common.BlockConfirmSign //all signs for block
	// modify by platon remove consensusCache
	Receipts types.Receipts
	State    *state.StateDB
}

// New creates a BlockExt object
func NewBlockExt(block *types.Block, blockNum uint64) *BlockExt {
	return &BlockExt{
		block:  block,
		number: blockNum,
		signs:  make([]*common.BlockConfirmSign, 0),
	}
}

var flowControl *FlowControl

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

//收到新块时，要判断新块是否是出块人在出块窗口内按合理节奏出块的
func (flowControl *FlowControl) control(nodeID discover.NodeID, curTime int64) bool {
	passed := false
	if flowControl.nodeID == nodeID {
		differ := curTime - flowControl.lastTime
		if differ >= flowControl.minInterval && differ <= flowControl.maxInterval {
			passed = true
		} else {
			passed = false
		}
	} else {
		passed = true
	}
	flowControl.nodeID = nodeID
	flowControl.lastTime = curTime

	return passed
}

//find BlockExt in cbft.blockExtMap
func (cbft *Cbft) findBlockExt(hash common.Hash) *BlockExt {
	if v, ok := cbft.blockExtMap[hash]; ok {
		return v
	}
	return nil
}

//collect all signs for block
// todo: to filter the duplicated sign
func (cbft *Cbft) collectSign(ext *BlockExt, sign *common.BlockConfirmSign) {
	if sign != nil {
		ext.signs = append(ext.signs, sign)
		if len(ext.signs) >= cbft.getThreshold(ext.block.Number().Sub(ext.block.Number(), common.Big1), ext.block.ParentHash(), ext.block.Number()) {
			ext.isConfirmed = true
		}
	}
}

func (parent *BlockExt) isParent(child *types.Block) bool {
	if parent.block != nil && parent.block.NumberU64()+1 == child.NumberU64() && parent.block.Hash() == child.ParentHash() {
		return true
	}
	return false
}

//to find current blockExt's parent with non-nil block
func (ext *BlockExt) findParent() *BlockExt {
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

//find
func (ext *BlockExt) findStoredAncestor() *BlockExt {
	if ext.block == nil {
		return nil
	}
	parent := cbft.findBlockExt(ext.block.ParentHash())
	if parent != nil && parent.block != nil && parent.block.NumberU64()+1 == ext.block.NumberU64() {
		if parent.isStored {
			return parent
		} else {
			return parent.findStoredAncestor()
		}
	} else {
		log.Warn("cannot find stored ancestor")
		return nil
	}
}

func (cbft *Cbft) isFork(newFork []*BlockExt, oldFork []*BlockExt) (bool, *BlockExt) {
	for i := 1; i < len(newFork); i++ {
		if oldFork[i].isConfirmed {
			return false, oldFork[i]
		}
	}
	return true, nil
}

func (cbft *Cbft) collectTxs(exts []*BlockExt) types.Transactions {
	txs := make([]*types.Transaction, 0)
	for _, ext := range exts {
		copy(txs, ext.block.Transactions())
	}
	return types.Transactions(txs)
}

//to find current blockExt's children with non-nil block
func (ext *BlockExt) findChildren() []*BlockExt {
	if ext.block == nil {
		return nil
	}
	children := make([]*BlockExt, 0)

	for _, child := range cbft.blockExtMap {
		if child.block != nil && child.block.ParentHash() == ext.block.Hash() {
			if child.block.NumberU64()-1 == ext.block.NumberU64() {
				children = append(children, child)
			} else {
				log.Warn("data error, child block hash is not mapping to number")
			}
		}
	}

	if len(children) == 0 {
		return nil
	} else {
		return children
	}
}

func (cbft *Cbft) saveBlock(hash common.Hash, ext *BlockExt) {
	cbft.blockExtMap[hash] = ext
	log.Debug("total blocks in memory", "totalBlocks", len(cbft.blockExtMap))
}

func (lower *BlockExt) isAncestor(higher *BlockExt) bool {
	if higher.block == nil || lower.block == nil {
		return false
	}
	generations := higher.block.NumberU64() - lower.block.NumberU64()

	for i := uint64(0); i < generations; i++ {
		parent := higher.findParent()
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

//to find a new highest confirmed blockExt from ext; If there are multiple highest confirmed blockExts, return the first.
func (cbft *Cbft) findNewHighestConfirmed(ext *BlockExt) *BlockExt {
	log.Debug("find first, highest confirmed block ")
	found := ext
	if !found.isConfirmed {
		found = nil
	}
	//each child has non-nil block
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			current := cbft.findNewHighestConfirmed(child)
			if current != nil && current.isConfirmed && (found == nil || current.block.NumberU64() > found.block.NumberU64()) {
				found = current
			}
		}
	}
	return found
}

//to find the highest block from ext; If there are multiple highest blockExts, return the one that has most signs
func (cbft *Cbft) findHighest(ext *BlockExt) *BlockExt {
	log.Debug("find highest block has most signs")
	highest := ext
	//each child has non-nil block
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			current := cbft.findHighest(child)
			if current.block.NumberU64() > highest.block.NumberU64() || (current.block.NumberU64() == highest.block.NumberU64() && len(current.signs) > len(highest.signs)) {
				highest = current
			}
		}
	}
	return highest
}

//to find the highest signed block from ext; If there are multiple highest blockExts, return the one that has most signs
func (cbft *Cbft) findHighestSigned(ext *BlockExt) *BlockExt {
	log.Debug("find highest block has most signs including mine's ")
	highest := ext

	if !highest.isSigned {
		highest = nil
	}
	//each child has non-nil block
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			current := cbft.findHighestSigned(child)
			if current != nil && current.isSigned {
				if highest == nil {
					highest = current
				} else if current.block.NumberU64() > highest.block.NumberU64() || (current.block.NumberU64() == highest.block.NumberU64() && len(current.signs) > len(highest.signs)) {
					highest = current
				}
			}
		}
	}
	return highest
}

func (cbft *Cbft) handleBlockAndDescendant(ext *BlockExt, parent *BlockExt, signIfPossible bool) error {
	log.Debug("handle block recursively", "hash", ext.block.Hash(), "number", ext.block.NumberU64())

	err := cbft.executeBlockAndDescendant(ext, parent)

	if err != nil {
		return nil
	}

	if ext.findChildren() == nil {
		if signIfPossible {
			if _, signed := cbft.signedSet[ext.block.NumberU64()]; !signed {
				cbft.sign(ext)
			}
		}
	} else {
		highest := cbft.findHighest(ext)
		logicalExts := cbft.backTrackBlocks(highest, ext, false)
		for _, logical := range logicalExts {
			if _, signed := cbft.signedSet[logical.block.NumberU64()]; !signed {
				cbft.sign(logical)
			}
		}
	}
	return nil
}

func (cbft *Cbft) executeBlockAndDescendant(ext *BlockExt, parent *BlockExt) error {
	log.Debug("execute block recursively", "hash", ext.block.Hash(), "number", ext.block.NumberU64())
	if ext.isLinked == false {
		err := cbft.execute(ext, parent)
		if err != nil {
			return err
		}
		ext.isLinked = true
	}
	//each child has non-nil block
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			return cbft.executeBlockAndDescendant(child, ext)
		}
	}
	return nil
}

//to sign a block
func (cbft *Cbft) sign(ext *BlockExt) {
	//签名
	sealHash := sealHash(ext.block.Header())
	signature, err := cbft.signFn(sealHash.Bytes())
	if err == nil {
		log.Debug("Sign block ", "Hash", ext.block.Hash(), "number", ext.block.NumberU64(), "sealHash", sealHash, "signature", hexutil.Encode(signature))

		sign := common.NewBlockConfirmSign(signature)
		ext.isSigned = true

		cbft.collectSign(ext, sign)

		//save this block number
		cbft.signedSet[ext.block.NumberU64()] = struct{}{}

		blockHash := ext.block.Hash()

		//send the BlockSignature to channel
		blockSign := &cbfttypes.BlockSignature{
			SignHash:  sealHash,
			Hash:      blockHash,
			Number:    ext.block.Number(),
			Signature: sign,
			ParentHash: ext.block.ParentHash(),
		}
		cbft.blockSignOutCh <- blockSign
	} else {
		panic("sign a block error")
	}
}

//execute the block based on its parent
// if success then set this block's level with Ledge, and save the receipts and state to consensusCache
func (cbft *Cbft) execute(ext *BlockExt, parent *BlockExt) error {
	// modify by platon remove consensusCache
	//state, err := cbft.consensusCache.MakeStateDB(parent.block)
	state, err := cbft.blockChain.StateAt(parent.block.Root())
	if err != nil {
		log.Error("execute block error, cannot make state based on parent")
		return err
	}

	//to execute
	receipts, err := cbft.blockChain.ProcessDirectly(ext.block, state, parent.block)
	if err == nil {
		// modify by platon remove consensusCache
		////save the receipts and state to consensusCache
		//cbft.consensusCache.WriteReceipts(ext.block.Hash(), receipts, ext.block.NumberU64())
		//cbft.consensusCache.WriteStateDB(ext.block.Root(), state, ext.block.NumberU64())
		ext.Receipts = receipts
		ext.State = state
	} else {
		log.Error("execute a block error", err)
	}
	return err
}

func reverse(s []*BlockExt) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (cbft *Cbft) backTrackBlocks(start *BlockExt, end *BlockExt, includeEnd bool) []*BlockExt {
	log.Debug("back track blocks", "startHash", start.block.Hash(), "startParentHash", end.block.ParentHash(), "endHash", start.block.Hash())

	found := false
	logicalExts := make([]*BlockExt, 1)
	logicalExts[0] = start

	for {
		parent := start.findParent()
		if parent == nil {
			break
		} else if parent.block.Hash() == end.block.Hash() && parent.block.NumberU64() == end.block.NumberU64() {
			log.Debug("ending of back track block ")
			if includeEnd {
				logicalExts = append(logicalExts, parent)
			}
			found = true
			break
		} else {
			log.Debug("found new block", "Hash", parent.block.Hash(), "ParentHash", parent.block.ParentHash(), "number", parent.block.NumberU64())
			logicalExts = append(logicalExts, parent)
			start = parent
		}
	}

	if found {
		//sorted by block number from lower to higher
		if len(logicalExts) > 1 {
			reverse(logicalExts)
		}
		return logicalExts
	} else {
		return nil
	}
}

//gather all blocks from the current highest confirmed one (excluded) to the new confirmed one
//the result is sorted by block number from lower to higher
func (cbft *Cbft) backTrackTillConfirmed(newConfirmed *BlockExt) []*BlockExt {

	log.Debug("found new block to store", "Hash", newConfirmed.block.Hash(), "ParentHash", newConfirmed.block.ParentHash(), "number", newConfirmed.block.NumberU64())

	existMap := make(map[common.Hash]struct{})

	foundExts := make([]*BlockExt, 1)
	foundExts[0] = newConfirmed

	existMap[newConfirmed.block.Hash()] = struct{}{}
	foundRoot := false
	for {
		parent := newConfirmed.findParent()
		if parent == nil {
			break
		}

		foundExts = append(foundExts, parent)

		if parent.isStored {
			foundRoot = true
			break
		} else {
			log.Debug("found new block to store", "Hash", parent.block.Hash(), "ParentHash", parent.block.ParentHash(), "number", parent.block.NumberU64())
			if _, exist := existMap[parent.block.Hash()]; exist {
				log.Error("get into a loop when finding new block to store")
				return nil
			}

			newConfirmed = parent
		}
	}

	if !foundRoot {
		log.Error("cannot lead to a stored block")
		return nil
	}

	//sorted by block number from lower to higher
	if len(foundExts) > 1 {
		reverse(foundExts)
	}
	return foundExts
}

func (cbft *Cbft) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	cbft.config.PrivateKey = privateKey
	cbft.config.NodeID = discover.PubkeyID(&privateKey.PublicKey)
}

// modify by platon remove consensusCache
//func SetConsensusCache(cache *Cache) {
//	cbft.consensusCache = cache
//}

func setHighestLogical(highestLogical *BlockExt) {
	cbft.highestLogical = highestLogical
	cbft.highestLogicalBlockCh <- highestLogical.block
}

func SetBackend(blockChain *core.BlockChain, txPool *core.TxPool) {
	log.Debug("init backend for CBFT")

	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	cbft.blockChain = blockChain
	cbft.ppos.SetStartTimeOfEpoch(blockChain.Genesis().Time().Int64())

	currentBlock := blockChain.CurrentBlock()

	genesisParentHash := bytes.Repeat([]byte{0x00}, 32)
	if bytes.Equal(currentBlock.ParentHash().Bytes(), genesisParentHash) && currentBlock.Number() == nil {
		currentBlock.Header().Number = big.NewInt(0)
	}

	log.Debug("init cbft.highestLogicalBlock", "Hash", currentBlock.Hash(), "number", currentBlock.NumberU64())

	confirmedBlock := NewBlockExt(currentBlock, currentBlock.NumberU64())
	confirmedBlock.isLinked = true
	confirmedBlock.isStored = true
	confirmedBlock.isConfirmed = true
	confirmedBlock.number = currentBlock.NumberU64()

	cbft.saveBlock(currentBlock.Hash(), confirmedBlock)

	cbft.highestConfirmed = confirmedBlock
	//cbft.highestLogical = confirmedBlock
	setHighestLogical(confirmedBlock)

	txPool = txPool
}

func SetDopsOption(blockChain *core.BlockChain) {
	cbft.ppos.SetCandidatePool(blockChain, cbft.config.InitialNodes)
}

func BlockSynchronisation() {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Debug("sync blocks finished.")

	currentBlock := cbft.blockChain.CurrentBlock()

	if currentBlock.NumberU64() > cbft.highestConfirmed.block.NumberU64() {
		log.Debug("found higher highestConfirmed block")

		confirmedBlock := NewBlockExt(currentBlock, currentBlock.NumberU64())
		confirmedBlock.isLinked = true
		confirmedBlock.isStored = true
		confirmedBlock.isConfirmed = true
		confirmedBlock.number = currentBlock.NumberU64()

		cbft.slideWindow(confirmedBlock)

		cbft.saveBlock(currentBlock.Hash(), confirmedBlock)

		cbft.highestConfirmed = confirmedBlock

		highestLogical := cbft.findHighestSigned(confirmedBlock)
		if highestLogical == nil {
			highestLogical = cbft.findHighest(confirmedBlock)
		}

		if highestLogical == nil {
			log.Warn("cannot find a logical block")
			return
		}

		setHighestLogical(highestLogical)

		children := confirmedBlock.findChildren()
		for _, child := range children {
			err := cbft.handleBlockAndDescendant(child, confirmedBlock, true)
			log.Error("block sync error", "err", err)
		}
	}
}

func (cbft *Cbft) dataReceiverGoroutine() {
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
				block, ok := v.(*types.Block)
				if ok {
					err := cbft.blockReceiver(block)
					if err != nil {
						log.Error("Error", "msg", err)
					}
				} else {
					log.Error("Received wrong data type")
				}
			}
		}
	}
}

//newConfirmed不能从blockExtMap中删除，但是可以从signedSet删除
func (cbft *Cbft) slideWindow(newConfirmed *BlockExt) {
	for hash, ext := range cbft.blockExtMap {
		if ext.number <= newConfirmed.block.NumberU64()-windowSize {
			if ext.block == nil {
				log.Debug("delete blockExt(only signs) from blockExtMap", "Hash", hash)
				delete(cbft.blockExtMap, hash)
			} else if ext.block.Hash() != newConfirmed.block.Hash() {
				log.Debug("delete blockExt from blockExtMap", "Hash", hash, "number", ext.block.NumberU64())
				delete(cbft.blockExtMap, hash)
			}
		}
	}

	for number, _ := range cbft.signedSet {
		if number <= cbft.highestConfirmed.block.NumberU64()-windowSize {
			log.Debug("delete number from signedSet", "number", number)
			delete(cbft.signedSet, number)
		}
	}

	log.Debug("remaining blocks in memory", "remainingBlocks", len(cbft.blockExtMap))
}

//1.1 如果新确认块高 > 当前不可逆块高，并且当前不可逆块高 是 新确认块祖先，则新确认块追溯到已经入链块之间的块，都需要入链
//1.2 如果新确认块高 > 当前不可逆块高，并且当前不可逆块高 不是 新确认块祖先，则无需处理
//2.1 如果新确认块高 = 当前不可逆块， 无需处理
//3.1 如果新确认块高 < 当前不可逆块，则可能分叉;如果允许分叉，则新确认块追溯到已经入链块之间的块，都需要入链；如果不允许分叉，则无需处理
func (cbft *Cbft) handleNewConfirmed(newConfirmed *BlockExt) error {
	if newConfirmed.block.NumberU64() > cbft.highestConfirmed.block.NumberU64() {
		//新确认块的高大于当前不可逆块高
		if cbft.highestConfirmed.isAncestor(newConfirmed) {
			//当前最高确认块是新确认块的祖先
			log.Debug("consensus success, new confirmed block is higher, and it is a descendant of the highest confirmed")

			blocksToStore := cbft.backTrackTillConfirmed(newConfirmed)

			if blocksToStore == nil {
				return errListConfirmedBlocks
			}

			return cbft.handleNewConfirmedContinue(newConfirmed, blocksToStore[1:])
		} else {
			//当前最高确认区块不是新确认块的祖先
			log.Warn("consensus error, new confirmed block is higher, but it is not a descendant of the highest confirmed")
			return nil
		}
	} else if newConfirmed.block.NumberU64() == cbft.highestConfirmed.block.NumberU64() {
		//新确认块的高等于当前最高确认区块
		log.Warn("consensus error, new confirmed block's number is same as the highest confirmed's")
		return nil
	} else {
		//新确认块的高小于当前最高确认区块，则可能分叉
		newFork := cbft.backTrackTillConfirmed(newConfirmed)
		if len(newFork) <= 1 {
			log.Error("new fork error")
			return nil
		}

		//newForm[0] is a confirmed and stored block
		oldFork := cbft.backTrackBlocks(cbft.highestConfirmed, newFork[0], true)
		if len(newFork) >= len(oldFork) {
			log.Error("new fork error")
			return nil
		}

		if isFork, cause := cbft.isFork(newFork, oldFork); !isFork {
			log.Warn("consensus success, but cannot fork to new confirmed", "causeHash", cause.block.Hash(), "causeNumber", cause.block.NumberU64())
			return nil
		} else {
			err := cbft.handleNewConfirmedContinue(newConfirmed, newFork[1:])

			if err == nil {
				// modify by platon remove consensusCache
				//state, err := cbft.consensusCache.MakeStateDB(newConfirmed.block)
				_, err := cbft.blockChain.StateAt(newConfirmed.block.Root())
				if err == nil {
					cbft.ppos.UpdateNodeList(cbft.blockChain, newConfirmed.block.Number(), newConfirmed.block.Hash())
				} else {
					log.Error("consensus success, but updateNodeList error", "err", err)
					return nil
				}

				//准备写入新分叉支链的交易
				txsInNewFork := cbft.collectTxs(newFork[1:])
				//已经写入被分叉支链的交易，需要恢复
				txsInOldFork := cbft.collectTxs(oldFork[1:])

				differ := types.TxDifference(txsInNewFork, txsInOldFork)

				log.Debug("consensus success, to recover the transactions from old fork", "txsCount", len(differ))

				//cbft.txPool.RecoverTxs(differ)
				return nil
			} else {
				log.Error("consensus success, but fork to new confirmed error", "err", err)
				return nil
			}
		}
	}
}

func (cbft *Cbft) handleNewConfirmedContinue(newConfirmed *BlockExt, blocksToStore []*BlockExt) error {
	log.Debug("store blocks to chain", "toStoreCount", len(blocksToStore))

	cbft.storeBlocks(blocksToStore)

	cbft.highestConfirmed = newConfirmed
	highestLogical := cbft.findHighestSigned(newConfirmed)
	if highestLogical == nil {
		highestLogical = cbft.findHighest(newConfirmed)
	}
	if highestLogical == nil {
		return errHighestLogicalBlock
	}

	setHighestLogical(highestLogical)

	log.Debug("to free memory")
	cbft.slideWindow(newConfirmed)
	return nil
}

//handle the received block signature
func (cbft *Cbft) signReceiver(sig *cbfttypes.BlockSignature) error {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Debug("=== begin to handle new signature ===", "Hash", sig.Hash, "number", sig.Number.Uint64())

	if sig.Number.Uint64() <= cbft.highestConfirmed.number {
		log.Warn("block sign is too late")
		return nil
	}

	ext := cbft.findBlockExt(sig.Hash)
	if ext == nil {
		log.Debug("have not received the corresponding block")

		//the block is nil
		ext = NewBlockExt(nil, sig.Number.Uint64())
		ext.isLinked = false

		cbft.saveBlock(sig.Hash, ext)
	} else if ext.isStored {
		//收到已经确认块的签名，直接扔掉
		log.Debug("received a highestConfirmed block's signature, just discard it")
		return nil
	}

	cbft.collectSign(ext, sig.Signature)

	log.Debug("count signatures", "signCount", len(ext.signs))

	if ext.isConfirmed && ext.isLinked {
		return cbft.handleNewConfirmed(ext)
	}

	log.Debug("=== end to handle new signature ===", "Hash", sig.Hash, "number", sig.Number.Uint64())

	return nil
}

//handle the received block
func (cbft *Cbft) blockReceiver(block *types.Block) error {

	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Debug("=== begin to handle block ===", "Hash", block.Hash(), "number", block.Number().Uint64(), "ParentHash", block.ParentHash())

	if block.NumberU64() <= cbft.highestConfirmed.block.NumberU64() {
		log.Warn("Received block is lower than the highestConfirmed block")
		return nil
	}

	if block.NumberU64() <= 0 {
		return errGenesisBlock
	}

	if block.NumberU64() <= cbft.highestConfirmed.number {
		return errBlockNumber
	}
	//recover the producer's NodeID
	producerNodeID, sign, err := ecrecover(block.Header())
	if err != nil {
		return err
	}

	curTime := toMilliseconds(time.Now())

	keepIt := cbft.shouldKeepIt(block.Number().Sub(block.Number(),common.Big1), block.ParentHash(), block.Number(), curTime, producerNodeID)
	log.Debug("check if block should be kept", "result", keepIt, "producerNodeID", hex.EncodeToString(producerNodeID.Bytes()[:8]))
	if !keepIt {
		return errIllegalBlock
	}

	//to check if there's a existing blockExt for received block
	//sometime we'll receive the block's sign before the block self.
	ext := cbft.findBlockExt(block.Hash())
	if ext == nil {
		ext = NewBlockExt(block, block.NumberU64())
		//default
		ext.isLinked = false

		cbft.saveBlock(block.Hash(), ext)

	} else if ext.block == nil {
		//received its sign before.
		ext.block = block
	} else {
		return errDuplicatedBlock
	}

	//collect the block's sign of producer
	log.Debug("collect this block's sign")

	cbft.collectSign(ext, common.NewBlockConfirmSign(sign))

	parent := ext.findParent()
	if parent != nil && parent.isLinked {
		//inTurn := cbft.inTurnVerify(curTime, producerNodeID)
		inTurn := cbft.inTurnVerify(block.Number().Sub(block.Number(),common.Big1), block.ParentHash(), block.Number(), curTime, producerNodeID)
		log.Debug("check if block is in turn", "result", inTurn, "producerNodeID", hex.EncodeToString(producerNodeID.Bytes()[:8]))

		passed := flowControl.control(producerNodeID, curTime)
		log.Debug("check if block is allowed by flow control", "result", passed, "producerNodeID", hex.EncodeToString(producerNodeID.Bytes()[:8]))

		signIfPossible := inTurn && passed && cbft.highestConfirmed.isAncestor(ext)

		err := cbft.handleBlockAndDescendant(ext, parent, signIfPossible)

		if err != nil {
			return err
		}

		newConfirmed := cbft.findNewHighestConfirmed(ext)
		if newConfirmed != nil {
			//处理新不可逆块
			log.Debug("found new highest confirmed block")
			return cbft.handleNewConfirmed(newConfirmed)
		}
	} else {
		log.Debug("cannot find block's parent, just keep it")
	}

	log.Debug("=== end to handle block ===", "Hash", block.Hash(), "number", block.Number().Uint64())
	return nil
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	log.Debug("call Author()", "Hash", header.Hash(), "number", header.Number.Uint64())

	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	log.Debug("call VerifyHeader()", "Hash", header.Hash(), "number", header.Number.Uint64(), "seal", seal)

	//todo:每秒一个交易，校验块高/父区块
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
	log.Debug("call VerifyHeaders()", "Headers count", len(headers))

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
	log.Debug("call VerifySeal()", "Hash", header.Hash(), "number", header.Number.String())

	return cbft.verifySeal(chain, header, nil)
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (b *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	log.Debug("call Prepare()", "Hash", header.Hash(), "number", header.Number.Uint64())

	cbft.lock.RLock()
	defer cbft.lock.RUnlock()

	//检查父区块
	if cbft.highestLogical.block == nil || header.ParentHash != cbft.highestLogical.block.Hash() || header.Number.Uint64()-1 != cbft.highestLogical.block.NumberU64() {
		return consensus.ErrUnknownAncestor
	}

	header.Difficulty = big.NewInt(2)

	//header.Extra[0:31] to store block's version info etc. and right pad with 0x00;
	//header.Extra[32:] to store block's sign of producer, the length of sign is 65.
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
	log.Debug("call Finalize()", "Hash", header.Hash(), "number", header.Number.Uint64(), "txs", len(txs), "receipts", len(receipts))
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	return types.NewBlock(header, txs, nil, receipts), nil
}

//to sign the block, and store the sign to header.Extra[32:], send the sign to chanel to broadcast to other consensus nodes
func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) error {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Debug("call Seal()", "number", block.NumberU64(), "parentHash", block.ParentHash())

	header := block.Header()
	number := block.NumberU64()

	if number == 0 {
		return errUnknownBlock
	}

	if !cbft.highestLogical.isParent(block) {
		log.Error("Futile block cause highest logical block changed", "parentHash", block.ParentHash())
		return errFutileBlock
	}

	// sign the seal hash
	sign, err := cbft.signFn(sealHash(header).Bytes())
	if err != nil {
		return err
	}

	//store the sign in  header.Extra[32:]
	copy(header.Extra[len(header.Extra)-extraSeal:], sign[:])

	sealedBlock := block.WithSeal(header)

	curExt := NewBlockExt(sealedBlock, sealedBlock.NumberU64())

	//this block is produced by local node, so need not execute in cbft.
	curExt.isLinked = true

	//collect the sign
	cbft.collectSign(curExt, common.NewBlockConfirmSign(sign))

	//save the block to cbft.blockExtMap
	cbft.saveBlock(sealedBlock.Hash(), curExt)

	log.Debug("seal complete", "Hash", sealedBlock.Hash(), "number", block.NumberU64())

	consensusNodes := cbft.ConsensusNodes(block.Number().Sub(block.Number(), common.Big1), block.ParentHash(), block.Number())
	if consensusNodes != nil && len(consensusNodes) == 1 {
		//only one consensus node, so, each block is highestConfirmed. (lock is needless)
		return cbft.handleNewConfirmed(curExt)
	}

	//reset cbft.highestLogicalBlockExt cause this block is produced by myself
	setHighestLogical(curExt)

	go func() {
		select {
		case <-stopCh:
			return
		case sealResultCh <- sealedBlock:
		default:
			log.Warn("Sealing result is not ready by miner", "sealHash", cbft.SealHash(header))
		}
	}()
	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func (b *Cbft) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	log.Debug("call CalcDifficulty()", "time", time, "parentHash", parent.Hash(), "parentNumber", parent.Number.Uint64())
	return big.NewInt(2)
}

// SealHash returns the hash of a block prior to it being sealed.
func (b *Cbft) SealHash(header *types.Header) common.Hash {
	log.Debug("call SealHash()", "Hash", header.Hash(), "number", header.Number.Uint64())
	return sealHash(header)
}

// Close implements consensus.Engine. It's a noop for clique as there is are no background threads.
func (cbft *Cbft) Close() error {
	log.Debug("call Close() ...")

	var err error
	cbft.closeOnce.Do(func() {
		// Short circuit if the exit channel is not allocated.
		if cbft.exitCh == nil {
			return
		}
		errc := make(chan error)
		cbft.exitCh <- errc
		err = <-errc
		close(cbft.exitCh)
	})
	return err
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (cbft *Cbft) APIs(chain consensus.ChainReader) []rpc.API {
	log.Debug("call APIs() ... ")

	return []rpc.API{{
		Namespace: "cbft",
		Version:   "1.0",
		Service:   &API{chain: chain, cbft: cbft},
		Public:    false,
	}}
}

//receive the new block signature
func (cbft *Cbft) OnBlockSignature(chain consensus.ChainReader, nodeID discover.NodeID, rcvSign *cbfttypes.BlockSignature) error {
	log.Debug("Receive a new signature", "Hash", rcvSign.Hash, "number", rcvSign.Number, "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]), "signHash", rcvSign.SignHash)

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

//receive the new block
func (cbft *Cbft) OnNewBlock(chain consensus.ChainReader, rcvBlock *types.Block) error {
	log.Debug("Receive a new block, put into chanel", "Hash", rcvBlock.Hash(), "number", rcvBlock.NumberU64(), "ParentHash", rcvBlock.ParentHash())

	cbft.dataReceiveCh <- rcvBlock
	return nil
}

//receive the new block
//netLatency：当前节点和nodeID直接的网络延迟
func (cbft *Cbft) OnPong(nodeID discover.NodeID, netLatency int64) error {
	log.Debug("Receive a net latency report", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]), "netLatency", netLatency)
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

func (cbft *Cbft) HighestLogicalBlock() *types.Block {
	cbft.lock.RLock()
	defer cbft.lock.RUnlock()

	log.Debug("call HighestLogicalBlock() ...")

	return cbft.highestLogical.block
}

func IsSignedBySelf(sealHash common.Hash, signature []byte) bool {
	ok, err := verifySign(cbft.config.NodeID, sealHash, signature)
	if err != nil {
		log.Error("verify sign error", "errors", err)
		return false
	}
	return ok
}

// send blockExt to channel
func (cbft *Cbft) storeBlocks(blocksToStore []*BlockExt) {
	for _, ext := range blocksToStore {
		cbftResult := &cbfttypes.CbftResult{
			Block:             ext.block,
			BlockConfirmSigns: ext.signs,
		}
		ext.isStored = true
		log.Debug("send to channel", "Hash", ext.block.Hash(), "number", ext.block.NumberU64(), "signCount", len(ext.signs))
		cbft.cbftResultOutCh <- cbftResult
	}
}

//to check if it's my turn to produce blocks
//func (cbft *Cbft) inTurn() bool {
//	curTime := toMilliseconds(time.Now())
//	inturn := cbft.calTurn(0, curTime, cbft.config.NodeID)
//	log.Debug("inTurn", "result", inturn)
//	return inturn
//
//}

func (cbft *Cbft) inTurn(parentNumber *big.Int, parentHash common.Hash, commitNumber *big.Int) bool {
	curTime := toMilliseconds(time.Now())
	inturn := cbft.calTurn(parentNumber, parentHash, commitNumber, curTime, cbft.config.NodeID)
	log.Debug("inTurn", "result", inturn)
	return inturn
}

//time in milliseconds
/*func (cbft *Cbft) inTurnVerify(curTime int64, nodeID discover.NodeID) bool {
	latency := cbft.avgLatency(nodeID)
	if latency >= maxAvgLatency {
		log.Debug("inTurnVerify, return false cause of net latency", "result", false, "latency", latency)
		return false
	}
	inTurnVerify := cbft.calTurn(curTime-latency, nodeID)
	log.Debug("inTurnVerify", "result", inTurnVerify, "latency", latency)
	return inTurnVerify
}
*/
func (cbft *Cbft) inTurnVerify(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, curTime int64, nodeID discover.NodeID) bool {
	latency := cbft.avgLatency(nodeID)
	if latency >= maxAvgLatency {
		log.Debug("inTurnVerify, return false cause of net latency", "result", false, "latency", latency)
		return false
	}
	inTurnVerify := cbft.calTurn(parentNumber, parentHash, blockNumber, curTime-latency, nodeID)
	log.Debug("inTurnVerify", "result", inTurnVerify, "latency", latency)
	return inTurnVerify
}

//time in milliseconds
func (cbft *Cbft) shouldKeepIt(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, curTime int64, nodeID discover.NodeID) bool {
	offset := 1000 * (cbft.config.Duration/2 - 1)
	keepIt := cbft.calTurn(parentNumber, parentHash, blockNumber, curTime-offset, nodeID)
	if !keepIt {
		keepIt = cbft.calTurn(parentNumber, parentHash, blockNumber, curTime+offset, nodeID)
	}
	log.Debug("shouldKeepIt", "result", keepIt, "offset", offset)
	return keepIt
}

func (cbft *Cbft) calTurn(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, curTime int64, nodeID discover.NodeID) bool {
	nodeIdx := cbft.ppos.BlockProducerIndex(parentNumber, parentHash, blockNumber, nodeID)
	startEpoch := cbft.ppos.StartTimeOfEpoch() * 1000

	if nodeIdx >= 0 {
		durationPerNode := cbft.config.Duration * 1000

		consensusNodes := cbft.ConsensusNodes(parentNumber, parentHash, blockNumber)
		if consensusNodes == nil || len(consensusNodes) <= 0 {
			return false
		}

		durationPerTurn := durationPerNode * int64(len(consensusNodes))

		min := nodeIdx * (durationPerNode)

		value := (curTime - startEpoch) % durationPerTurn

		max := (nodeIdx + 1) * durationPerNode

		log.Debug("calTurn", "idx", nodeIdx, "min", min, "value", value, "max", max, "curTime", curTime, "startEpoch", startEpoch)

		if value > min && value < max {
			return true
		}
	}
	return false
}

// producer's signature = header.Extra[32:]
// public key can be recovered from signature, the length of public key is 65,
// the length of NodeID is 64, nodeID = publicKey[1:]
func ecrecover(header *types.Header) (discover.NodeID, []byte, error) {
	var nodeID discover.NodeID
	if len(header.Extra) < extraSeal {
		return nodeID, []byte{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]
	sealHash := sealHash(header)

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

// seal hash, only include from byte[0] to byte[32] of header.Extra
func sealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[0:32],
		header.MixDigest,
		header.Nonce,
	})
	hasher.Sum(hash[:0])

	return hash
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

func (cbft *Cbft) getThreshold(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) int {
	consensusNodes := cbft.ConsensusNodes(parentNumber, parentHash, blockNumber)
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
	return cbft.inTurn(parentNumber, parentHash, commitNumber)
}

//func (cbft *Cbft) FormerNodeID() []discover.NodeID {
//	return cbft.ppos.getFormerNodeID()
//}
//
//func (cbft *Cbft) CurrentNodeID() []discover.NodeID {
//	return cbft.ppos.getCurrentNodeID()
//}
//
//func (cbft *Cbft) NextNodeID() []discover.NodeID {
//	return cbft.ppos.getNextNodeID()
//}

func (cbft *Cbft) FormerNodes(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []*discover.Node {
	return cbft.ppos.getFormerNodes(parentNumber, parentHash, blockNumber)
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

// wether nodeID in current or next
func (cbft *Cbft) CheckFutureConsensusNode(nodeID discover.NodeID) (bool, error) {
	log.Debug("call CheckFutureConsensusNode()", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]))
	return cbft.ppos.NodeIndexInFuture(nodeID) >= 0, nil
}

// wether nodeID in former or current or next
//func (cbft *Cbft) IsConsensusNode() (bool, error) {
//	log.Debug("call IsConsensusNode()")
//	return cbft.ppos.AnyIndex(cbft.config.NodeID) >= 0, nil
//}

func (cbft *Cbft) Election(state *state.StateDB, blockNumber *big.Int) ([]*discover.Node, error) {
	return cbft.ppos.Election(state, blockNumber)
}

func (cbft *Cbft) Switch(state *state.StateDB) bool {
	return cbft.ppos.Switch(state)
}

func (cbft *Cbft) GetWitness(state *state.StateDB, flag int) ([]*discover.Node, error) {
	return cbft.ppos.GetWitness(state, flag)
}

func (cbft *Cbft) GetOwnNodeID() discover.NodeID {
	return cbft.config.NodeID
}

func (cbft *Cbft) SetNodeCache(state *state.StateDB, parentNumber, currentNumber *big.Int, parentHash, currentHash common.Hash) error {
	log.Info("cbft SetNodeCache", "parentNumber", parentNumber, "parentHash", parentHash, "currentNumber", currentNumber, "currentHash", currentHash)
	return cbft.ppos.SetNodeCache(state, parentNumber, currentNumber, parentHash, currentHash)
}


