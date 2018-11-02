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
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync"
	"time"
)

var (
	errSign               = errors.New("sign error")
	errUnauthorizedSigner = errors.New("unauthorized signer")
	errOverdueBlock       = errors.New("overdue block")
	errDuplicatedBlock    = errors.New("duplicated block")
	errBlockNumber        = errors.New("error block number")
	errUnknownBlock       = errors.New("unknown block")

	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")
	extraSeal           = 65
)

type Cbft struct {
	config          *params.CbftConfig
	dpos            *dpos
	rotating        *rotating
	blockSignOutCh  chan *cbfttypes.BlockSignature //a channel to send block signature
	cbftResultOutCh chan *cbfttypes.CbftResult     //a channel to send consensus result
	closeOnce       sync.Once
	exitCh          chan chan error

	blockExtMap sync.Map //store all received blocks and signs

	signReceiveCh          chan *cbfttypes.BlockSignature //a channel to receive block signature
	blockReceiveCh         chan *types.Block              //a channel to receive block
	blockChain             *core.BlockChain               //the block chain
	highestLogicalBlockExt *BlockExt                      //highest logical block
	irreversibleBlockExt   *BlockExt                      //highest irreversible block
	signedSet              sync.Map                       //all block numbers signed by local node
	lock                   sync.Mutex
	consensusCache         *Cache //cache for cbft consensus
}

var cbft *Cbft

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, blockSignatureCh chan *cbfttypes.BlockSignature, cbftResultCh chan *cbfttypes.CbftResult) *Cbft {
	_dpos := newDpos(config.InitialNodes)

	cbft = &Cbft{
		config:          config,
		dpos:            _dpos,
		rotating:        newRotating(_dpos, config.Duration),
		blockSignOutCh:  blockSignatureCh,
		cbftResultOutCh: cbftResultCh,

		signReceiveCh:  make(chan *cbfttypes.BlockSignature, 240),
		blockReceiveCh: make(chan *types.Block, 24),
	}

	//启动协程处理新收到的签名
	go cbft.signReceiverGoroutine()

	//启动协程处理新收到的块
	go cbft.blockReceiverGoroutine()

	return cbft
}

//each block has differ level in it's lifecycle
type Level int

const (
	Discrete Level = iota //neither should execute nor be signed
	Legal                 //should execute
	Logical               //should be singed
)

//the extension for Block
type BlockExt struct {
	block           *types.Block
	level           Level
	isIrreversible  bool
	signsUpdateTime time.Time
	signs           []*common.BlockConfirmSign //all signs for block
}

// New creates a BlockExt object
func NewBlockExt(block *types.Block) *BlockExt {
	return &BlockExt{
		block: block,
		signs: make([]*common.BlockConfirmSign, 0),
	}
}

//find BlockExt in cbft.blockExtMap
func (cbft *Cbft) findBlockExt(hash common.Hash) *BlockExt {
	if v, ok := cbft.blockExtMap.Load(hash); ok {
		return v.(*BlockExt)
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

//to check if current blockExt is the next logical block
func (ext *BlockExt) isNextLogical() bool {
	if ext.block != nil && ext.block.ParentHash() == cbft.highestLogicalBlockExt.block.Hash() && ext.block.NumberU64() == cbft.highestLogicalBlockExt.block.NumberU64()+1 {
		return true
	}
	return false
}

//to check if current blockExt has a logical child
func (ext *BlockExt) hasLogicalChild() bool {
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			if child.level == Logical {
				return true
			}
		}
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
		if parent.block != nil && parent.block.NumberU64()+1 == ext.block.NumberU64() {
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

	f := func(k, v interface{}) bool {
		child, _ := v.(*BlockExt)
		if child != nil && child.block != nil && child.block.ParentHash() == ext.block.Hash() {
			if child.block.NumberU64()-1 == ext.block.NumberU64() {
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

func (cbft *Cbft) saveBlock(ext *BlockExt) {
	cbft.blockExtMap.Store(ext.block.Hash(), ext)
}

func (cbft *Cbft) saveEmptyBlock(hash common.Hash, ext *BlockExt) {
	cbft.blockExtMap.Store(hash, ext)
}

//reset logical path from start (the starter is a logical block already)
//sign the on-path block and set the last blockExt as the highest logical blockExt
func (cbft *Cbft) resetLogicalPath(start *BlockExt) {
	log.Info("reset the path of logical blocks", "start blockHash", start.block.Hash())

	highest := cbft.findHighest(start)

	logicals := cbft.collectLogicals(start, highest)

	cbft.resetHighestLogical(highest)

	cbft.signLogicals(logicals)
}

//to check two blocks are the same or not
func (ext *BlockExt) isSame(other *BlockExt) bool {
	if ext.block.Hash() == other.block.Hash() {
		return true
	}
	return false
}

//reset the highest logical block
func (cbft *Cbft) resetHighestLogical(highest *BlockExt) {
	log.Info("reset the highest logical block", "hash", highest.block.Hash(), "number", highest.block.NumberU64())
	cbft.highestLogicalBlockExt = highest

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

//to find the blockExt that has most signs; If there are multiple blockExts, return the highest and first one.
func (cbft *Cbft) findMaxSigns(ext *BlockExt) *BlockExt {
	log.Info("recursion in findMaxSigns()")
	max := ext
	//each child has non-nil block
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			current := cbft.findMaxSigns(child)
			if len(current.signs) > len(max.signs) || ((len(current.signs) == len(max.signs)) && current.block.NumberU64() > max.block.NumberU64()) {
				max = current
			}
		}
	}
	return max
}

func (cbft *Cbft) execDescendant(ext *BlockExt, parent *BlockExt) {
	log.Info("recursion in findMaxSigns()")

	log.Info("execute the block recursively", "hash", ext.block.Hash(), "number", ext.block.NumberU64())
	if ext.level == Discrete {
		cbft.execute(ext, parent)
	}
	//each child has non-nil block
	children := ext.findChildren()
	if children != nil {
		for _, child := range children {
			cbft.execDescendant(child, ext)
		}
	}
}

//to sign the blocks excluding those are already logical and whose block numbers are signed.
func (cbft *Cbft) signLogicals(exts []*BlockExt) {
	for _, ext := range exts {
		if ext.level == Logical {
			continue
		}
		if _, signed := cbft.signedSet.Load(ext.block.NumberU64()); !signed {
			cbft.sign(ext)
		}
		ext.level = Logical
	}
}

//to sign a block
func (cbft *Cbft) sign(ext *BlockExt) {
	//签名
	sealHash := sealHash(ext.block.Header())
	signature, err := cbft.signFn(sealHash.Bytes())
	if err == nil {
		log.Info("Sign block ", "hash", ext.block.Hash(), "number", ext.block.NumberU64(), "sealHash", sealHash, "signature", hexutil.Encode(signature))

		sign := common.NewBlockConfirmSign(signature)
		ext.collectSign(sign)

		//save this block number
		cbft.signedSet.Store(ext.block.NumberU64(), struct{}{})

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
		//set this block's level with Ledge
		ext.level = Legal

		//save the receipts and state to consensusCache
		cbft.consensusCache.WriteReceipts(ext.block.Hash(), receipts, ext.block.NumberU64())
		cbft.consensusCache.WriteStateDB(ext.block.Root(), state, ext.block.NumberU64())

	} else {
		log.Error("execute a block error", err)
	}
}

//gather all irreversible blocks from the root one (excluded) to the highest one
//the result is sorted by block number from lower to higher
func (cbft *Cbft) listIrreversibles(newIrr *BlockExt) []*BlockExt {

	log.Info("Found new irreversible block", "Hash", newIrr.block.Hash(), "ParentHash", newIrr.block.ParentHash(), "Number", newIrr.block.NumberU64())

	existMap := make(map[common.Hash]struct{})

	exts := make([]*BlockExt, 1)
	exts[0] = newIrr

	existMap[newIrr.block.Hash()] = struct{}{}

	findRootIrr := false
	for {
		parent := newIrr.findParent()
		if parent == nil {
			break
		}
		if parent.isIrreversible {
			findRootIrr = true
			break
		} else {
			log.Info("Found new irreversible block", "Hash", parent.block.Hash(), "ParentHash", parent.block.ParentHash(), "Number", parent.block.NumberU64())
			if _, exist := existMap[newIrr.block.Hash()]; exist {
				log.Error("Irreversible blocks get into a loop")
				return nil
			}
			exts = append(exts, parent)
		}
	}

	if !findRootIrr {
		log.Error("cannot lead to root irreversible block")
		return nil
	}

	//sorted by block number from lower to higher
	if len(exts) > 1 {
		for i, j := 0, len(exts)-1; i < j; i, j = i+1, j-1 {
			exts[i], exts[j] = exts[j], exts[i]
		}
	}
	return exts
}

//收集从end到start（start已经是合理块）的合理块列表，按块高排序
func (cbft *Cbft) collectLogicals(start *BlockExt, end *BlockExt) []*BlockExt {
	exts := make([]*BlockExt, 1)
	exts[0] = end

	if start.isSame(end) {
		return exts
	}

	findStart := false
	for {
		//父节点满足ext.Block!=nil
		parent := end.findParent()
		if parent == nil {
			break
		}
		exts = append(exts, parent)
		if parent.isSame(start) {
			findStart = true
			break
		}
	}

	if !findStart {
		panic("cannot link to logical block")
	}

	//按块高排好序
	if len(exts) > 1 {
		for i, j := 0, len(exts)-1; i < j; i, j = i+1, j-1 {
			exts[i], exts[j] = exts[j], exts[i]
		}
	}
	return exts
}

func (cbft *Cbft) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	cbft.config.PrivateKey = privateKey
	cbft.config.NodeID = discover.PubkeyID(&privateKey.PublicKey)
}

func SetConsensusCache(cache *Cache) {
	cbft.consensusCache = cache
}

func SetBlockChain(blockChain *core.BlockChain) {
	log.Info("初始化cbft.blockChain")

	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	cbft.blockChain = blockChain
	cbft.dpos.SetStartTimeOfEpoch(blockChain.Genesis().Time().Int64())

	currentBlock := blockChain.CurrentBlock()

	genesisParentHash := bytes.Repeat([]byte{0x00}, 32)
	if bytes.Equal(currentBlock.ParentHash().Bytes(), genesisParentHash) && currentBlock.Number() == nil {
		currentBlock.Header().Number = big.NewInt(0)
		currentBlock.Number()
	}

	log.Info("初始化cbft.highestLogicalBlock", "hash", currentBlock.Hash(), "number", currentBlock.NumberU64())

	blockExt := NewBlockExt(currentBlock)
	blockExt.level = Logical
	blockExt.isIrreversible = true

	//最高合理块
	cbft.resetHighestLogical(blockExt)

	//当前不可逆块
	cbft.irreversibleBlockExt = blockExt

	cbft.saveBlock(blockExt)
}

func BlockSynchronisation() {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Info("区块同步处理")
	currentBlock := cbft.blockChain.CurrentBlock()
	//如果链上块高>内存中不可逆块高
	if currentBlock.NumberU64() > cbft.irreversibleBlockExt.block.NumberU64() {
		log.Info("区块同步处理， 发现更高的不可逆区块")
		//准备新的不可逆块
		newIrr := NewBlockExt(currentBlock)
		newIrr.level = Logical
		newIrr.isIrreversible = true

		//保存
		cbft.saveBlock(newIrr)

		//当前不可逆块
		cbft.irreversibleBlockExt = newIrr

		//执行新块和所有后代块
		cbft.execDescendant(newIrr, cbft.highestLogicalBlockExt)

		//重置从ext开始的合理链路径（如果链路上节点需要签名则签名，并设置最高合理节点）
		cbft.resetLogicalPath(newIrr)
		//清理blockExtMap/signedSet
		cbft.cleanCtx()

	}
}

func (cbft *Cbft) signReceiverGoroutine() {
	for {
		select {
		case sign := <-cbft.signReceiveCh:
			cbft.signReceiver(sign)
		}
	}
}

func (cbft *Cbft) cleanCtx() {
	f1 := func(k, v interface{}) bool {
		ext, _ := v.(*BlockExt)
		if ext.block.NumberU64() < cbft.irreversibleBlockExt.block.NumberU64() {
			cbft.blockExtMap.Delete(ext.block.Hash())
		}
		return true
	}
	cbft.blockExtMap.Range(f1)

	f2 := func(k, v interface{}) bool {
		number, _ := k.(uint64)
		if number < cbft.irreversibleBlockExt.block.NumberU64() {
			cbft.blockExtMap.Delete(number)
		}
		return true
	}

	cbft.signedSet.Range(f2)
}

func (cbft *Cbft) blockReceiverGoroutine() {
	for {
		select {
		case block := <-cbft.blockReceiveCh:
			cbft.blockReceiver(block)
		}
	}
}

//处理新的不可逆块
//把新的不可逆块到原不可逆块路径上的区块，写入链。并设置新的不可逆块。
//此新不可逆块，已经被处理（执行/签名），所有后代也已经处理（执行），并且新的合理链路已经设置好，新合理链上的块也已经处理（执行/签名）
//从ext开始查看是否有可入链块
func (cbft *Cbft) handleNewIrreversible(newIrr *BlockExt) {
	log.Info("处理新的不可逆块")
	//收集可入链的块，不包括原来的不可逆块，并按块高排好序
	irrs := cbft.listIrreversibles(newIrr)

	if irrs == nil {
		log.Error("gather all irreversible blocks error")
		return
	}

	log.Info("处理新的不可逆块", "可以入链的区块数量", len(irrs))

	cbft.writeChain(irrs)

	//清理blockExtMap/signedSet
	cbft.cleanCtx()

	//设置新节点为不可逆块
	newIrr.isIrreversible = true
	newIrr.level = Logical
}

//收到签名处理
func (cbft *Cbft) signReceiver(sig *cbfttypes.BlockSignature) {
	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	log.Info("收到新签名==>>", "blockHash", sig.Hash)

	ext := cbft.findBlockExt(sig.Hash)
	if ext == nil {
		log.Info("新签名对应的区块还没收到")
		//没有的话，new一个BlockExt
		ext = NewBlockExt(nil)
		ext.level = Discrete
		cbft.saveEmptyBlock(sig.Hash, ext)
	} else if ext.isIrreversible {
		log.Info("新签名对应的区块已经是不可逆块，直接丢弃此签名")
		return
	}

	//收集签名
	signCount := ext.collectSign(sig.Signature)

	log.Info("区块签名数量", "signCount", signCount)

	if signCount >= cbft.getThreshold() && ext.level != Discrete {
		//区块收到的签名数量>=15，可以入链了
		log.Info("区块签名数量>2f+1，成为最新的不可逆区块")
		//处理新不可逆块
		cbft.handleNewIrreversible(ext)

		if ext.level == Legal {
			//重置从ext开始的合理链路径（如果链路上节点需要签名则签名，并设置最高合理节点）
			cbft.resetLogicalPath(ext)
		}
	}
}

//收到区块处理
func (cbft *Cbft) blockReceiver(block *types.Block) error {

	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	if block.NumberU64() <= 0 {
		return nil
	}
	//从签名恢复出出块人地址·
	producerNodeID, sign, err := ecrecover(block.Header())
	if err != nil {
		return err
	}
	log.Info("处理收到的新块", "producerNodeID", producerNodeID.String())

	//检查块是否在出块人的时间窗口内生成的
	//时间合法性计算，不合法返回error
	log.Info("检查块是否在出块人的时间窗口内生成的")
	if cbft.isOverdue(block.Time().Int64(), producerNodeID) {
		log.Error("迟到的不合法区块，直接丢弃")
		return errOverdueBlock
	}

	//查看是否先收到过的签名
	ext := cbft.findBlockExt(block.Hash())
	if ext == nil {
		//没有的话，new一个BlockExt
		ext = NewBlockExt(block)
		//设置缺省level，并保存
		ext.level = Discrete
		cbft.saveBlock(ext)

	} else if ext.block == nil {
		//先收到签名，然后再收到区块
		ext.block = block
	} else {
		return errDuplicatedBlock
	}

	//收集出块人的签名
	ext.collectSign(common.NewBlockConfirmSign(sign))

	if ext.isNextLogical() { //是下一个最高合理节点
		//执行新块和所有后代块
		log.Info("是下一个最高合理节点")
		cbft.execDescendant(ext, cbft.highestLogicalBlockExt)

		//重置从ext开始的合理链路径（如果链路上节点需要签名则签名，并设置最高合理节点）
		cbft.resetLogicalPath(ext)

	} else {
		parent := ext.findParent()
		if parent != nil {
			if parent.level == Discrete {
				log.Info("收到的新块的父区块是孤区块，暂时不处理")
			} else if parent.level == Legal || (parent.level == Logical && parent.hasLogicalChild()) {
				// 执行所有后代
				cbft.execDescendant(ext, parent)
			}
		} else {
			log.Info("收到的新块是孤区块，暂时不处理")
		}
	}

	//新块已经处理（执行/签名），后代也已经处理（执行），并且可能的新合理路径也设置好（未签名的也已经签名）
	//从ext开始查看是否有可入链块
	if ext.level != Discrete {
		log.Info("从新块开始，查看是否有可入链区块")
		newIrr := cbft.findHighestNewIrreversible(ext)

		if newIrr != nil {
			//处理新不可逆块
			cbft.handleNewIrreversible(newIrr)
		}
	}
	return nil
}

func (cbft *Cbft) ShouldSeal() (bool, error) {
	return cbft.inTurn(), nil
}

func (cbft *Cbft) ConsensusNodes() ([]discover.NodeID, error) {
	log.Info("call ConsensusNodes() ...")
	return cbft.dpos.primaryNodeList, nil
}

func (cbft *Cbft) CheckConsensusNode(nodeID discover.NodeID) (bool, error) {
	log.Info("call CheckConsensusNode(), parameter: ", "nodeID", nodeID.String())
	return cbft.dpos.NodeIndex(nodeID) >= 0, nil
}

func (cbft *Cbft) IsConsensusNode() (bool, error) {
	log.Info("call IsConsensusNode() ...")
	return cbft.dpos.NodeIndex(cbft.config.NodeID) >= 0, nil
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
	log.Info("call Author(), parameter: ", "headerHash", header.Hash(), "headerNumber", header.Number)

	// 返回出块节点对应的矿工钱包地址
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
// 从区块头的 extraData 字段中取出出块点的签名，从签名推导出公钥信息，然后从公钥算出出块节点的地址
// 如果该出块节点是共识节点之一，并且此区块头是在此出块节点的出块时间窗口内打包生成的，则认为该区块为合法区块
// chain: 	当前的链
// header: 	需要验证的区块头
// seal:	是否要验证封印（出块签名）
func (cbft *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	log.Info("call VerifyHeader(), parameter: ", "headerHash", header.Hash(), "headerNumber", header.Number, "seal", seal)

	//todo:每秒一个交易，校验块高/父区块
	if header.Number == nil {
		return errUnknownBlock
	}
	// Don't waste time checking blocks from the future
	//if header.Time.Cmp(big.NewInt(time.Now().Unix())) > 0 {
	//return consensus.ErrFutureBlock
	//}

	if len(header.Extra) < extraSeal {
		return errMissingSignature
	}
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (cbft *Cbft) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	log.Info("call VerifyHeaders(), parameter: ", "lenHeaders", len(headers))

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
// 校验(别的结点广播过来的)区块信息
// 主要是对区块的出块节点，以及区块难度值的确认
func (cbft *Cbft) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	log.Info("call VerifySeal(), parameter: ", "headerHash", header.Hash(), "headerNumber", header.Number.String())

	return cbft.verifySeal(chain, header, nil)
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (b *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	log.Info("call Prepare(), parameter: ", "headerHash", header.Hash(), "headerNumber", header.Number.String())

	//检查父区块
	if header.ParentHash != cbft.highestLogicalBlockExt.block.Hash() || header.Number.Uint64()-1 != cbft.highestLogicalBlockExt.block.NumberU64() {
		return consensus.ErrUnknownAncestor
	}

	header.Difficulty = big.NewInt(2)

	//Extra中，有32个字节存放版本等信息，占用27个字节（后补0到32个字节），后65个字节存放出块人的签名）
	if len(header.Extra) < 32 {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, 32-len(header.Extra))...)
	}
	header.Extra = header.Extra[:32]

	header.Extra = append(header.Extra, make([]byte, consensus.ExtraSeal)...)
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given, and returns the final block.
func (cbft *Cbft) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	log.Info("call Finalize(), parameter: ", "headerHash", header.Hash(), "headerNumber", header.Number.String(), "txs", len(txs), "receipts", len(receipts))

	// 生成具体的区块信息
	// 填充上Header.Root, TxHash, ReceiptHash, UncleHash等几个属性
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	return types.NewBlock(header, txs, nil, receipts), nil
}

// 完成对区块的签名成功，并设置到header.Extra中，然后把区块发送到sealResultCh通道中（然后会被组播到其它共识节点）
func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) error {
	log.Info("call Seal(), parameter", "number", block.NumberU64(), "parentHash", block.ParentHash())

	header := block.Header()
	number := block.NumberU64()

	if number == 0 {
		// 不密封创世块
		return errUnknownBlock
	}

	parent := cbft.findBlockExt(header.ParentHash)
	if parent == nil {
		log.Error("找不到父节点", "parentHash", header.ParentHash)
		return errUnknownBlock
	}

	//todo:检查ext.block 和 入参 block

	// 开始签名
	sign, err := cbft.signFn(sealHash(header).Bytes())
	if err != nil {
		return err
	}

	//将签名结果替换区块头的Extra字段（专门支持记录额外信息的）
	copy(header.Extra[len(header.Extra)-extraSeal:], sign[:])

	//得到完成签名后的区块
	sealedBlock := block.WithSeal(header)

	//ext中保存的是完成签名后的block。header.extra中必须有签名
	ext := NewBlockExt(sealedBlock)

	// 标识为合理块，并且执行（在挖矿过程中执行）/签名过。
	// 这样，本地节点出的块，就不会在共识引擎中执行，这样，相应的BlockExt中就没有此区块的执行回执receipts和状态state
	ext.level = Logical

	//收集新区块的签名
	ext.collectSign(common.NewBlockConfirmSign(sign))

	//保存(blockExtMap.key必须是块经过Seal后的hash)
	cbft.saveBlock(ext)

	log.Info("签名完成", "number", block.NumberU64(), "blockHash", sealedBlock.Hash(), "parentHash", block.ParentHash())

	if len(cbft.dpos.primaryNodeList) == 1 {
		//单个节点，直接出块
		cbft.handleNewIrreversible(ext)

		cbft.highestLogicalBlockExt = ext

		return nil
	}

	cbft.lock.Lock()
	//把新块作为最高区块
	cbft.highestLogicalBlockExt = ext
	cbft.lock.Unlock()

	go func() {
		select {
		case <-stopCh: //如果先收到stop（客户端RPC发出)，则直接返回
			return
		case sealResultCh <- sealedBlock: //发送给p2p，把区块广播到其它节点
		default: //如果没有接收数据，则走default
			log.Warn("Sealing result is not read by miner", "sealHash", cbft.SealHash(header))
		}
	}()
	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func (b *Cbft) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	log.Info("call CalcDifficulty(), parameter", "time", time, "parentHash", parent.Hash(), "parentNumber", parent.Number.String())

	return big.NewInt(2)
}

// SealHash returns the hash of a block prior to it being sealed.
func (b *Cbft) SealHash(header *types.Header) common.Hash {
	log.Info("call SealHash(), parameter", "headerHash", header.Hash(), "headerNumber", header.Number.String())

	//return consensus.SigHash(header)
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

//收到新的区块签名
//需要验证签名是否时nodeID签名的
func (cbft *Cbft) OnBlockSignature(chain consensus.ChainReader, nodeID discover.NodeID, rcvSign *cbfttypes.BlockSignature) error {
	log.Info("收到新的区块签名==>>, parameter", "nodeID", nodeID.String(), "blockHash", rcvSign.Hash, "Number", rcvSign.Number, "signHash", rcvSign.SignHash, "sign", rcvSign.Signature.String())

	ok, err := verifySign(nodeID, rcvSign.SignHash, rcvSign.Signature[:])
	if err != nil {
		return err
	}

	if !ok {
		log.Error("unauthorized signer")
		return errUnauthorizedSigner
	}

	cbft.signReceiveCh <- rcvSign

	return nil
}

//收到新的区块
func (cbft *Cbft) OnNewBlock(chain consensus.ChainReader, rcvBlock *types.Block) error {
	log.Info("Received a new block, put into chanel", "Hash", rcvBlock.Hash(), "Number", rcvBlock.NumberU64(), "ParentHash", rcvBlock.ParentHash(), "headerExtra", hexutil.Encode(rcvBlock.Header().Extra))

	cbft.blockReceiveCh <- rcvBlock
	return nil
}

func (cbft *Cbft) HighestLogicalBlock() *types.Block {

	log.Info("call HighestLogicalBlock() ...")

	return cbft.highestLogicalBlockExt.block
}

func IsSignedBySelf(sealHash common.Hash, signature []byte) bool {
	ok, err := verifySign(cbft.config.NodeID, sealHash, signature)
	if err != nil {
		return false
	}
	return ok
}

//把区块写入出块channel
func (cbft *Cbft) writeChain(exts []*BlockExt) {

	for _, ext := range exts {

		log.Info("区块准备入链", "Hash", ext.block.Hash(), "Number", ext.block.NumberU64(), "signCount", len(ext.signs))

		cbftResult := &cbfttypes.CbftResult{
			Block: ext.block,
			/*Receipts:          ext.receipts,
			State:             ext.state,*/
			BlockConfirmSigns: ext.signs,
		}

		//把需要保存的数据，发往通道：cbftResultCh
		cbft.cbftResultOutCh <- cbftResult
	}
}

//出块时间窗口期与出块节点匹配
func (cbft *Cbft) inTurn() bool {
	singerIdx := cbft.dpos.NodeIndex(cbft.config.NodeID)
	/*if singerIdx >= 0 {
		durationMilliseconds := cbft.config.Duration * 1000
		totalDuration := durationMilliseconds * int64(len(cbft.dpos.primaryNodeList))

		value1 := singerIdx*(durationMilliseconds) - int64(cbft.config.MaxLatency/3)

		value2 := (time.Now().Unix()*1000 - cbft.dpos.StartTimeOfEpoch()) % totalDuration

		value3 := (singerIdx+1)*durationMilliseconds - int64(cbft.config.MaxLatency*2/3)

		if value2 > value1 && value3 > value2 {
			return true
		}
	}
	return false*/

	idxInturn := -1
	for idx := 0; idx < len(cbft.dpos.primaryNodeList); idx++ {
		inturn := calTurn(int64(idx))
		log.Info("节点出块轮值情况", "inturn", inturn, "nodeID", hexutil.Encode(cbft.dpos.primaryNodeList[idx][:]))

		if inturn {
			idxInturn = idx
		}
	}

	return singerIdx == int64(idxInturn)
}

func calTurn(idx int64) bool {
	if idx >= 0 {
		durationMilliseconds := cbft.config.Duration * 1000
		totalDuration := durationMilliseconds * int64(len(cbft.dpos.primaryNodeList))

		value1 := idx*(durationMilliseconds) - int64(cbft.config.MaxLatency/3)

		value2 := (time.Now().Unix() - cbft.dpos.StartTimeOfEpoch()) * 1000 % totalDuration

		value3 := (idx+1)*durationMilliseconds - int64(cbft.config.MaxLatency*2/3)

		//log.Info("计算参数", "now", now, "idx", idx, "durationMilliseconds", durationMilliseconds, "totalDuration", totalDuration, "MaxLatency", cbft.config.MaxLatency, "StartTimeOfEpoch", cbft.dpos.StartTimeOfEpoch())

		if value2 > value1 && value3 > value2 {
			return true
		}
	}
	return false
}

//收到新的区块后，检查新区块的时间合法性
func (cbft *Cbft) isOverdue(blockTimeInSecond int64, nodeID discover.NodeID) bool {
	singerIdx := cbft.dpos.NodeIndex(nodeID)

	durationMilliseconds := cbft.config.Duration * 1000

	totalDuration := durationMilliseconds * int64(len(cbft.dpos.primaryNodeList))

	//从StartTimeOfEpoch开始到now的完整轮数
	rounds := (time.Now().Unix() - cbft.dpos.StartTimeOfEpoch()) * 1000 / totalDuration

	//nodeID的最晚出块时间
	deadline := cbft.dpos.StartTimeOfEpoch()*1000 + totalDuration*rounds + durationMilliseconds*(singerIdx+1)

	//nodeID加上合适的延迟后的最晚出块时间
	deadline = deadline + int64(float64(cbft.config.MaxLatency)*cbft.config.LegalCoefficient)

	if deadline < time.Now().Unix() {
		//出块时间+延迟后，仍然小于当前时间（即收到区块的时间），则认为是超时的废区块，直接丢弃
		return true
	}
	return false
}

//NodeID是64字节，而publicKey是65字节，publicKey后64字节才是NodeID
func ecrecover(header *types.Header) (discover.NodeID, []byte, error) {
	// Retrieve the signature from the header extra-data

	//NodeID是64字节，而publicKey是65字节，publicKey后64字节才是NodeID
	var nodeID discover.NodeID
	if len(header.Extra) < extraSeal {
		return nodeID, []byte{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]
	log.Info("收到新块", "sign", hexutil.Encode(signature))
	sealHash := sealHash(header)
	log.Info("收到新块", "sealHash", sealHash)

	pubkey, err := crypto.Ecrecover(sealHash.Bytes(), signature)
	if err != nil {
		return nodeID, []byte{}, err
	}

	//转成discover.NodeID
	nodeID, err = discover.BytesID(pubkey[1:])
	if err != nil {
		return nodeID, []byte{}, err
	}
	return nodeID, signature, nil
}

// verify sign, check the sign is from the right node.
func verifySign(expectedNodeID discover.NodeID, sealHash common.Hash, signature []byte) (bool, error) {
	log.Info("verify sign", "sealHash", sealHash, "signature", hexutil.Encode(signature), "expectedNodeID", hexutil.Encode(expectedNodeID.Bytes()))

	pubkey, err := crypto.SigToPub(sealHash.Bytes(), signature)

	if err != nil {
		log.Error("verify sign error", "errors", err)
		return false, err
	}

	nodeID := discover.PubkeyID(pubkey)
	//比较两个[]byte
	log.Info("从签名恢复出的NodeID", "nodeID", nodeID.String())
	if bytes.Equal(nodeID.Bytes(), expectedNodeID.Bytes()) {
		log.Warn("the node id discover from signature is different from the node id from which the sign is received.")
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
		//header.Extra[0:len(header.Extra)-65],
		header.Extra[0:32],
		header.MixDigest,
		header.Nonce,
	})
	hasher.Sum(hash[:0])

	return hash
}

// VerifySeal()函数基于跟Seal()完全一样的算法原理
// 通过验证区块的某些属性(Header.Nonce，Header.MixDigest等)是否正确，来确定该区块是否已经经过Seal操作
func (cbft *Cbft) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	return nil
}

//返回当前时间UTC时区的毫秒数
func nowMillisecond() int64 {
	return time.Now().UnixNano() % 1e6 / 1e3
}

func (cbft *Cbft) signFn(headerHash []byte) (sign []byte, err error) {
	log.Info("signFN", "headerHash", hexutil.Encode(headerHash))

	nodeID := discover.PubkeyID(&cbft.config.PrivateKey.PublicKey)

	log.Info("signFN", "nodeIDFromPubKey", nodeID.String())

	return crypto.Sign(headerHash, cbft.config.PrivateKey)
}

//取最高区块
func (cbft *Cbft) getHighestLogicalBlock() *types.Block {
	log.Info("获取最高合理区块", "hash", cbft.highestLogicalBlockExt.block.Hash(), "number", cbft.highestLogicalBlockExt.block.NumberU64())
	return cbft.highestLogicalBlockExt.block
}

func (cbft *Cbft) getThreshold() int {
	trunc := len(cbft.dpos.primaryNodeList) * 2 / 3
	remainder := len(cbft.dpos.primaryNodeList) * 2 % 3

	if remainder == 0 {
		return int(trunc)
	} else {
		return int(trunc + 1)
	}
}
