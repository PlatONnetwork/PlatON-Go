// Package bft implements the BFT consensus engine.
package cbft

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/consensus"
	"Platon-go/core"
	"Platon-go/core/cbfttypes"
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
	"math/big"
	"sync"
	"time"
)

var (
	errSign               = errors.New("sign error")
	errUnauthorizedSigner = errors.New("unauthorized signer")
	errIllegalBlock       = errors.New("illegal block")
	errDuplicatedBlock    = errors.New("duplicated block")
	errBlockNumber        = errors.New("error block number")
	errUnknownBlock       = errors.New("unknown block")
	errFutileBlock        = errors.New("futile block")
	errGenesisBlock       = errors.New("cannot handle genesis block")
	errForNextBlock       = errors.New("cannot find a block for next")
	errListIrrBlocks      = errors.New("list irreversible blocks error")
	errMissingSignature   = errors.New("extra-data 65 byte signature suffix missing")
	extraSeal             = 65

	windowSize = uint64(20)

	blockSpeedRatio = uint64(20) //meaning 20%
)

type Cbft struct {
	config                *params.CbftConfig
	dpos                  *dpos
	rotating              *rotating
	blockSignOutCh        chan *cbfttypes.BlockSignature //a channel to send block signature
	cbftResultOutCh       chan *cbfttypes.CbftResult     //a channel to send consensus result
	highestLogicalBlockCh chan *types.Block
	closeOnce             sync.Once
	exitCh                chan chan error

	//todo:（先log,再处理）
	blockExtMap map[common.Hash]*BlockExt //store all received blocks and signs

	dataReceiveCh  chan interface{} //a channel to receive block signature
	blockChain     *core.BlockChain //the block chain
	highestLogical *BlockExt        //for next block
	irreversible   *BlockExt        //highest irreversible block

	//todo:（先log,再处理）
	signedSet      map[uint64]struct{} //all block numbers signed by local node
	lock           sync.RWMutex
	consensusCache *Cache //cache for cbft consensus

	netLatencyMap map[discover.NodeID]*list.List
}

var cbft *Cbft

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, blockSignatureCh chan *cbfttypes.BlockSignature, cbftResultCh chan *cbfttypes.CbftResult, highestLogicalBlockCh chan *types.Block) *Cbft {
	_dpos := newDpos(config.InitialNodes)

	cbft = &Cbft{
		config:                config,
		dpos:                  _dpos,
		rotating:              newRotating(_dpos, config.Duration),
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
	block    *types.Block
	isLinked bool
	isSigned bool
	isStored bool
	number   uint64
	signs    []*common.BlockConfirmSign //all signs for block
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
	nodeID    discover.NodeID
	lastTime  int64
	maxOffset int64
	minOffset int64
}

func NewFlowControl() *FlowControl {
	return &FlowControl{
		nodeID:    discover.NodeID{},
		maxOffset: int64(cbft.config.Period*1000 + cbft.config.Period*1000*blockSpeedRatio/100),
		minOffset: int64(cbft.config.Period*1000 - cbft.config.Period*1000*blockSpeedRatio/100),
	}
}

//收到新块时，要判断新块是否是出块人在出块窗口内按合理节奏出块的
func (flowControl *FlowControl) control(nodeID discover.NodeID, curTime int64) bool {
	passed := false
	if flowControl.nodeID == nodeID {
		differ := curTime - flowControl.lastTime
		if differ >= flowControl.minOffset && differ <= flowControl.maxOffset {
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
func (ext *BlockExt) collectSign(sign *common.BlockConfirmSign) int {
	if sign != nil {
		ext.signs = append(ext.signs, sign)
		return len(ext.signs)
	}
	return 0
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
	log.Debug("Blocks in memory", "counts", len(cbft.blockExtMap))
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

//to find a new irreversible blockExt from ext; If there are multiple irreversible blockExts, return the first.
func (cbft *Cbft) findHighestNewIrreversible(ext *BlockExt) *BlockExt {
	log.Info("recursion in findHighestNewIrreversible()")
	irr := ext
	if len(irr.signs) < cbft.getThreshold() {
		irr = nil
	}
	//each child has non-nil block
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			current := cbft.findHighestNewIrreversible(child)
			if current != nil && len(current.signs) >= cbft.getThreshold() && (irr == nil || current.block.NumberU64() > irr.block.NumberU64()) {
				irr = current
			}
		}
	}
	return irr
}

//to find the highest block from ext; If there are multiple highest blockExts, return the one that has most signs
func (cbft *Cbft) findHighest(ext *BlockExt) *BlockExt {
	log.Info("recursion in findHighest()")
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
	log.Info("recursion in findHighestSigned()")
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

func (cbft *Cbft) handleBlockAndDescendant(ext *BlockExt, parent *BlockExt, signIfPossible bool) {
	log.Info("handle block recursively", "hash", ext.block.Hash(), "number", ext.block.NumberU64())

	cbft.executeBlockAndDescendant(ext, parent)

	if ext.findChildren() == nil {
		if signIfPossible {
			if _, signed := cbft.signedSet[ext.block.NumberU64()]; !signed {
				cbft.sign(ext)
			}
		}
	} else {
		highest := cbft.findHighest(ext)
		logicalExts := cbft.backTrackLogicals(ext, highest)
		for _, logical := range logicalExts {
			if _, signed := cbft.signedSet[logical.block.NumberU64()]; !signed {
				cbft.sign(logical)
			}
		}
	}
}

func (cbft *Cbft) executeBlockAndDescendant(ext *BlockExt, parent *BlockExt) {
	log.Info("handle block recursively", "hash", ext.block.Hash(), "number", ext.block.NumberU64())
	if ext.isLinked == false {
		cbft.execute(ext, parent)
		ext.isLinked = true
	}
	//each child has non-nil block
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			cbft.executeBlockAndDescendant(child, ext)
		}
	}
}

//to sign a block
func (cbft *Cbft) sign(ext *BlockExt) {
	//签名
	sealHash := sealHash(ext.block.Header())
	signature, err := cbft.signFn(sealHash.Bytes())
	if err == nil {
		log.Info("Sign block ", "Hash", ext.block.Hash(), "number", ext.block.NumberU64(), "sealHash", sealHash, "signature", hexutil.Encode(signature))

		sign := common.NewBlockConfirmSign(signature)
		ext.collectSign(sign)
		ext.isSigned = true

		//save this block number
		cbft.signedSet[ext.block.NumberU64()] = struct{}{}

		blockHash := ext.block.Hash()

		//send the BlockSignature to channel
		blockSign := &cbfttypes.BlockSignature{
			SignHash:  sealHash,
			Hash:      blockHash,
			Number:    ext.block.Number(),
			Signature: sign,
		}
		cbft.blockSignOutCh <- blockSign
	} else {
		panic("can't sign a block")
	}
}

//execute the block based on its parent
// if success then set this block's level with Ledge, and save the receipts and state to consensusCache
func (cbft *Cbft) execute(ext *BlockExt, parent *BlockExt) {
	state, err := cbft.consensusCache.MakeStateDB(parent.block)
	if err != nil {
		log.Error("execute block error, cannot make state based on parent")
		return
	}

	//to execute
	receipts, err := cbft.blockChain.ProcessDirectly(ext.block, state)
	if err == nil {
		//save the receipts and state to consensusCache
		cbft.consensusCache.WriteReceipts(ext.block.Hash(), receipts, ext.block.NumberU64())
		cbft.consensusCache.WriteStateDB(ext.block.Root(), state, ext.block.NumberU64())
	} else {
		log.Error("execute a block error", err)
	}
}

//exclude start
func (cbft *Cbft) backTrackLogicals(start *BlockExt, end *BlockExt) []*BlockExt {

	log.Info("backTrackLogicals", "start", start.block.Hash(), "startParentHash", start.block.ParentHash(), "end", end.block.Hash())

	logicalExts := make([]*BlockExt, 1)
	logicalExts[0] = end

	for {
		parent := end.findParent()
		if parent == nil {
			break
		} else if end.block.ParentHash() == start.block.Hash() && end.block.NumberU64() == start.block.NumberU64()+1 {
			log.Debug("ending of back track logicals ")
			break
		} else {
			log.Debug("Found new logical block", "Hash", parent.block.Hash(), "ParentHash", parent.block.ParentHash(), "number", parent.block.NumberU64())
			logicalExts = append(logicalExts, parent)
		}
	}

	//sorted by block number from lower to higher
	if len(logicalExts) > 1 {
		for i, j := 0, len(logicalExts)-1; i < j; i, j = i+1, j-1 {
			logicalExts[i], logicalExts[j] = logicalExts[j], logicalExts[i]
		}
	}
	return logicalExts
}

//gather all irreversible blocks from the root one (excluded) to the highest one
//the result is sorted by block number from lower to higher
func (cbft *Cbft) backTrackIrreversibles(newIrr *BlockExt) []*BlockExt {

	log.Info("Found new irreversible block", "Hash", newIrr.block.Hash(), "ParentHash", newIrr.block.ParentHash(), "number", newIrr.block.NumberU64())

	existMap := make(map[common.Hash]struct{})

	IrrExts := make([]*BlockExt, 1)
	IrrExts[0] = newIrr

	existMap[newIrr.block.Hash()] = struct{}{}
	findRootIrr := false
	for {
		parent := newIrr.findParent()
		if parent == nil {
			break
		}
		if parent.isStored {
			findRootIrr = true
			break
		} else {
			log.Info("Found new irreversible block", "Hash", parent.block.Hash(), "ParentHash", parent.block.ParentHash(), "number", parent.block.NumberU64())
			if _, exist := existMap[parent.block.Hash()]; exist {
				log.Error("New irreversible block get into a loop")
				return nil
			}
			IrrExts = append(IrrExts, parent)
		}

		newIrr = parent
	}

	if !findRootIrr {
		log.Error("cannot lead to a irreversible block")
		return nil
	}

	//sorted by block number from lower to higher
	if len(IrrExts) > 1 {
		for i, j := 0, len(IrrExts)-1; i < j; i, j = i+1, j-1 {
			IrrExts[i], IrrExts[j] = IrrExts[j], IrrExts[i]
		}
	}
	return IrrExts
}

func (cbft *Cbft) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	cbft.config.PrivateKey = privateKey
	cbft.config.NodeID = discover.PubkeyID(&privateKey.PublicKey)
}

func SetConsensusCache(cache *Cache) {
	cbft.consensusCache = cache
}

func setHighestLogical(highestLogical *BlockExt) {
	cbft.highestLogical = highestLogical
	cbft.highestLogicalBlockCh <- highestLogical.block
}

func SetBlockChain(blockChain *core.BlockChain) {
	log.Info("init cbft.blockChain")

	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	cbft.blockChain = blockChain
	cbft.dpos.SetStartTimeOfEpoch(blockChain.Genesis().Time().Int64())

	currentBlock := blockChain.CurrentBlock()

	genesisParentHash := bytes.Repeat([]byte{0x00}, 32)
	if bytes.Equal(currentBlock.ParentHash().Bytes(), genesisParentHash) && currentBlock.Number() == nil {
		currentBlock.Header().Number = big.NewInt(0)
	}

	log.Info("init cbft.highestLogicalBlock", "Hash", currentBlock.Hash(), "number", currentBlock.NumberU64())

	irrBlock := NewBlockExt(currentBlock, currentBlock.NumberU64())
	irrBlock.isLinked = true
	irrBlock.isStored = true
	irrBlock.number = currentBlock.NumberU64()

	cbft.saveBlock(currentBlock.Hash(), irrBlock)

	cbft.irreversible = irrBlock
	//cbft.highestLogical = irrBlock
	setHighestLogical(irrBlock)

}

func BlockSynchronisation() {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Info("sync blocks finished.")

	currentBlock := cbft.blockChain.CurrentBlock()

	if currentBlock.NumberU64() > cbft.irreversible.block.NumberU64() {
		log.Info("found higher irreversible block")

		irrBlock := NewBlockExt(currentBlock, currentBlock.NumberU64())
		irrBlock.isLinked = true
		irrBlock.isStored = true
		irrBlock.number = currentBlock.NumberU64()

		cbft.slideWindow(irrBlock)

		cbft.saveBlock(currentBlock.Hash(), irrBlock)

		cbft.irreversible = irrBlock

		highestLogical := cbft.findHighestSigned(irrBlock)
		if cbft.highestLogical == nil {
			highestLogical = cbft.findHighest(irrBlock)
		}
		setHighestLogical(highestLogical)

		children := irrBlock.findChildren()
		for _, child := range children {
			cbft.handleBlockAndDescendant(child, irrBlock, true)
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

func (cbft *Cbft) slideWindow(newIrr *BlockExt) {
	for hash, ext := range cbft.blockExtMap {
		if ext.number <= cbft.irreversible.block.NumberU64()-windowSize {
			if ext.block == nil || ext.block.Hash() != cbft.irreversible.block.Hash() {
				log.Info("to delete hash from blockExtMap", "Hash", ext.block.Hash(), "number", ext.block.NumberU64())
				delete(cbft.blockExtMap, hash)
			}
		}
	}

	for number, _ := range cbft.signedSet {
		if number <= cbft.irreversible.block.NumberU64()-windowSize {
			log.Info("to delete number from signedSet", "number", number)
			delete(cbft.signedSet, number)
		}
	}
}

//to handle the new irreversible block
// (this block has executed and signed, all it's descendant has executed, and the logical path has reset, the block on the path has signed if necessary)
//gather the blocks from original irreversible block to this new one (excluding the original one), and save these blocks to chain
//clean cbft.blockExtMap, cbft.signedSet
//reset the cbft.irreversibleBlockExt
//1.1 如果新确认块高 > 当前不可逆块高，并且当前不可逆块高 是 新确认块祖先，则新确认块追溯到已经入链块之间的块，都需要入链
//1.2 如果新确认块高 > 当前不可逆块高，并且当前不可逆块高 不是 新确认块祖先，则无需处理
//2.1 如果新确认块高 = 当前不可逆块， 无需处理
//3.1 如果新确认块高 < 当前不可逆块， 则发送分叉，新确认块追溯到已经入链块之间的块，都需要入链
func (cbft *Cbft) handleNewIrreversible(newIrr *BlockExt) error {
	needSlideWindow := false
	if newIrr.block.NumberU64() > cbft.irreversible.block.NumberU64() {
		//新确认块的高大于当前不可逆块高
		if cbft.irreversible.isAncestor(newIrr) {
			//当前不可逆块是新确认块的祖先
			log.Info("consensus success, the block is higher and in the same branch")

			needSlideWindow = true
			//at last, handle the new irreversible
			//cbft.handleNewIrreversible(newIrr)
		} else {
			//当前不可逆块不是新确认块的祖先
			log.Warn("consensus useless, the block is higher and in b branch")
			return nil
		}
	} else if newIrr.block.NumberU64() == cbft.irreversible.block.NumberU64() {
		//新确认块的高等于当前不可逆块高
		log.Warn("consensus error, the block with same number are confirmed")
		return nil
	} else {
		//新确认块的高小于当前不可逆块高，并且是另一条分支（相同分支的小于当前不可逆块高的块，都已经入链了）
		log.Info("consensus success, the block is lower and in another branch")
		//cbft.handleNewIrreversible(newIrr)
	}

	irrs := cbft.backTrackIrreversibles(newIrr)

	if irrs == nil {
		return errListIrrBlocks
	}

	log.Info("found irreversible blocks", "count", len(irrs))

	cbft.storeIrreversibles(irrs)

	cbft.irreversible = newIrr
	highestLogical := cbft.findHighestSigned(newIrr)
	if highestLogical == nil {
		highestLogical = cbft.findHighest(newIrr)
	}
	if highestLogical == nil {
		return errForNextBlock
	}

	setHighestLogical(highestLogical)

	if needSlideWindow {
		cbft.slideWindow(newIrr)
	}
	return nil
}

//handle the received block signature
func (cbft *Cbft) signReceiver(sig *cbfttypes.BlockSignature) error {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Info("=== begin to handle new signature ===", "Hash", sig.Hash, "number", sig.Number.Uint64())

	if sig.Number.Uint64() <= cbft.irreversible.number {
		log.Warn("block sign is too late")
		return nil
	}

	ext := cbft.findBlockExt(sig.Hash)
	if ext == nil {
		log.Info("have not received the corresponding block")

		//the block is nil
		ext = NewBlockExt(nil, sig.Number.Uint64())
		ext.isLinked = false

		cbft.saveBlock(sig.Hash, ext)
	} else if ext.isStored {
		//收到已经确认块的签名，直接扔掉
		log.Info("received a irreversible block's signature, just discard it")
		return nil
	}

	signCount := ext.collectSign(sig.Signature)

	log.Info("count signatures", "Count", signCount)

	if signCount >= cbft.getThreshold() && ext.isLinked {
		return cbft.handleNewIrreversible(ext)
	}

	log.Info("=== end to handle new signature ===", "Hash", sig.Hash, "number", sig.Number.Uint64())

	return nil
}

//handle the received block
func (cbft *Cbft) blockReceiver(block *types.Block) error {

	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Info("=== begin to handle block ===", "Hash", block.Hash(), "number", block.Number().Uint64(), "ParentHash", block.ParentHash())

	if block.NumberU64() <= cbft.irreversible.block.NumberU64() {
		log.Warn("Received block is lower than the irreversible block")
		return nil
	}

	if block.NumberU64() <= 0 {
		return errGenesisBlock
	}

	if block.NumberU64() <= cbft.irreversible.number {
		return errBlockNumber
	}
	//recover the producer's NodeID
	producerNodeID, sign, err := ecrecover(block.Header())
	if err != nil {
		return err
	}

	curTime := toMilliseconds(time.Now())

	keepIt := cbft.shouldKeepIt(curTime, producerNodeID)
	log.Info("check if block should be kept", "result", keepIt, "producerNodeID", hex.EncodeToString(producerNodeID.Bytes()[:8]))
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
	log.Info("collect this block's sign")
	ext.collectSign(common.NewBlockConfirmSign(sign))

	parent := ext.findParent()
	if parent != nil && parent.isLinked {
		inTurn := cbft.inTurnVerify(curTime, producerNodeID)
		log.Info("check if block is in turn", "result", inTurn, "producerNodeID", hex.EncodeToString(producerNodeID.Bytes()[:8]))

		passed := flowControl.control(producerNodeID, curTime)
		log.Info("check if block is allowed by flow control", "result", passed, "producerNodeID", hex.EncodeToString(producerNodeID.Bytes()[:8]))

		signIfPossible := inTurn && passed && cbft.irreversible.isAncestor(ext)

		cbft.handleBlockAndDescendant(ext, parent, signIfPossible)

		newIrr := cbft.findHighestNewIrreversible(ext)
		if newIrr != nil {
			//处理新不可逆块
			log.Info("found higher new irreversible")
			return cbft.handleNewIrreversible(newIrr)
		}
	} else {
		log.Info("cannot find block's parent, just keep it")
	}

	log.Info("=== end to handle block ===", "Hash", block.Hash(), "number", block.Number().Uint64())
	return nil
}

func (cbft *Cbft) ShouldSeal() (bool, error) {
	return cbft.inTurn(), nil
}

func (cbft *Cbft) ConsensusNodes() ([]discover.NodeID, error) {
	log.Info("call ConsensusNodes()")
	return cbft.dpos.primaryNodeList, nil
}

func (cbft *Cbft) CheckConsensusNode(nodeID discover.NodeID) (bool, error) {
	log.Info("call CheckConsensusNode()", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]))
	return cbft.dpos.NodeIndex(nodeID) >= 0, nil
}

func (cbft *Cbft) IsConsensusNode() (bool, error) {
	log.Info("call IsConsensusNode()")
	return cbft.dpos.NodeIndex(cbft.config.NodeID) >= 0, nil
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	log.Info("call Author()", "Hash", header.Hash(), "number", header.Number.Uint64())

	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	log.Info("call VerifyHeader()", "Hash", header.Hash(), "number", header.Number.Uint64(), "seal", seal)

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
	log.Info("call VerifyHeaders()", "Headers count", len(headers))

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
	log.Info("call VerifySeal()", "Hash", header.Hash(), "number", header.Number.String())

	return cbft.verifySeal(chain, header, nil)
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (b *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	log.Info("call Prepare()", "Hash", header.Hash(), "number", header.Number.Uint64())

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
	log.Info("call Finalize()", "Hash", header.Hash(), "number", header.Number.Uint64(), "txs", len(txs), "receipts", len(receipts))
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	return types.NewBlock(header, txs, nil, receipts), nil
}

//to sign the block, and store the sign to header.Extra[32:], send the sign to chanel to broadcast to other consensus nodes
func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) error {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Info("call Seal()", "number", block.NumberU64(), "parentHash", block.ParentHash())

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
	curExt.collectSign(common.NewBlockConfirmSign(sign))

	//save the block to cbft.blockExtMap
	cbft.saveBlock(sealedBlock.Hash(), curExt)

	log.Info("seal complete", "Hash", sealedBlock.Hash(), "number", block.NumberU64())

	if len(cbft.dpos.primaryNodeList) == 1 {
		//only one consensus node, so, each block is irreversible. (lock is needless)
		return cbft.handleNewIrreversible(curExt)
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
	log.Info("call CalcDifficulty()", "time", time, "parentHash", parent.Hash(), "parentNumber", parent.Number.Uint64())
	return big.NewInt(2)
}

// SealHash returns the hash of a block prior to it being sealed.
func (b *Cbft) SealHash(header *types.Header) common.Hash {
	log.Info("call SealHash()", "Hash", header.Hash(), "number", header.Number.Uint64())
	return sealHash(header)
}

// Close implements consensus.Engine. It's a noop for clique as there is are no background threads.
func (cbft *Cbft) Close() error {
	log.Info("call Close() ...")

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
	log.Info("call APIs() ... ")

	return []rpc.API{{
		Namespace: "cbft",
		Version:   "1.0",
		Service:   &API{chain: chain, cbft: cbft},
		Public:    false,
	}}
}

//receive the new block signature
func (cbft *Cbft) OnBlockSignature(chain consensus.ChainReader, nodeID discover.NodeID, rcvSign *cbfttypes.BlockSignature) error {
	log.Info("Received a new signature", "Hash", rcvSign.Hash, "number", rcvSign.Number, "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]), "signHash", rcvSign.SignHash)

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
	log.Info("Received a new block, put into chanel", "Hash", rcvBlock.Hash(), "number", rcvBlock.NumberU64(), "ParentHash", rcvBlock.ParentHash())

	cbft.dataReceiveCh <- rcvBlock
	return nil
}

//receive the new block
//netLatency：当前节点和nodeID直接的网络延迟
func (cbft *Cbft) OnPong(nodeID discover.NodeID, netLatency int64) {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Info("Received a report for net latency", "Hash", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]), "netLatency", netLatency)

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
}

func (cbft *Cbft) avgLatency(nodeID discover.NodeID) int64 {
	cbft.lock.RLock()
	defer cbft.lock.RUnlock()

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

	log.Info("call HighestLogicalBlock() ...")

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
func (cbft *Cbft) storeIrreversibles(exts []*BlockExt) {
	for _, ext := range exts {
		cbftResult := &cbfttypes.CbftResult{
			Block:             ext.block,
			BlockConfirmSigns: ext.signs,
		}
		ext.isStored = true
		log.Info("send to channel", "Hash", ext.block.Hash(), "number", ext.block.NumberU64(), "signCount", len(ext.signs))
		cbft.cbftResultOutCh <- cbftResult
	}
}

//to check if it's my turn to produce blocks
func (cbft *Cbft) inTurn() bool {
	curTime := toMilliseconds(time.Now())

	nodeIdx := cbft.dpos.NodeIndex(cbft.config.NodeID)
	start := cbft.dpos.StartTimeOfEpoch() * 1000
	if nodeIdx >= 0 {
		durationMilliseconds := cbft.config.Duration * 1000
		totalDuration := durationMilliseconds * int64(len(cbft.dpos.primaryNodeList))

		startTime := nodeIdx * (durationMilliseconds)

		curTime := (curTime - start) % totalDuration

		endTime := (nodeIdx + 1) * durationMilliseconds

		log.Info("calTurn", "nodeIdx", nodeIdx, "startTime", startTime, "endTime", endTime, "curTime", curTime, "curTime", curTime, "startTimeOfEpoch", start)

		if curTime > startTime && curTime < endTime {
			return true
		}
	}
	return false
}

//time in milliseconds
func (cbft *Cbft) inTurnVerify(curTime int64, nodeID discover.NodeID) bool {

	nodeIdx := cbft.dpos.NodeIndex(nodeID)
	start := cbft.dpos.StartTimeOfEpoch() * 1000
	durationPerNode := cbft.config.Duration * 1000
	durationPerTurn := durationPerNode * int64(len(cbft.dpos.primaryNodeList))

	avgLatency := cbft.avgLatency(nodeID)
	blockTime := curTime - avgLatency

	rounds := (curTime - start) / durationPerTurn

	startTimeOfCurTurn := start + durationPerTurn*rounds

	startTime := startTimeOfCurTurn + durationPerNode*nodeIdx - avgLatency/3
	endTime := startTimeOfCurTurn + durationPerNode*(nodeIdx+1) - avgLatency*2/3

	log.Info("inTurnVerify", "start", start, "startTime", startTime, "endTime", endTime, "blockTime", blockTime, "curTime", curTime, "avgLatency", avgLatency)

	if blockTime > startTime && blockTime < endTime {
		return true
	} else {
		return false
	}
}

//time in milliseconds
func (cbft *Cbft) shouldKeepIt(curTime int64, nodeID discover.NodeID) bool {

	offset := 1000 * (cbft.config.Duration/2 - 1)

	nodeIdx := cbft.dpos.NodeIndex(nodeID)
	start := cbft.dpos.StartTimeOfEpoch() * 1000
	durationPerNode := cbft.config.Duration * 1000
	durationPerTurn := durationPerNode * int64(len(cbft.dpos.primaryNodeList))

	avgLatency := cbft.avgLatency(nodeID)
	blockTime := curTime - avgLatency

	rounds := (curTime - start) / durationPerTurn

	startTimeOfCurTurn := start + durationPerTurn*rounds

	startTime := startTimeOfCurTurn + durationPerNode*nodeIdx - offset
	endTime := startTimeOfCurTurn + durationPerNode*(nodeIdx+1) + offset

	log.Info("shouldKeepIt", "start", start, "startTime", startTime, "endTime", endTime, "blockTime", blockTime, "curTime", curTime, "avgLatency", avgLatency)

	if blockTime > startTime && blockTime < endTime {
		return true
	} else {
		return false
	}
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

func (cbft *Cbft) getThreshold() int {
	trunc := len(cbft.dpos.primaryNodeList) * 2 / 3
	return int(trunc + 1)
}

func toMilliseconds(t time.Time) int64 {
	return t.UnixNano() / 1e6
}
