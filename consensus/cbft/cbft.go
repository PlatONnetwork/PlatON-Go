// Package bft implements the BFT consensus engine.
package cbft

import (
	"Platon-go/common"
	"Platon-go/consensus"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/crypto"
	"Platon-go/crypto/sha3"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"Platon-go/rlp"
	"Platon-go/rpc"
	"errors"
	"math/big"
	"sync"
	"time"
)

const (
	inmemorySignatures = 4096
)

var (
	errUnauthorizedSigner = errors.New("unauthorized signer")
	errOverdueBlock       = errors.New("overdue block")
	errBlockNumber        = errors.New("error block number")
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")
	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")
)
var (
	epochLength = uint64(210000)
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
)

type Cbft struct {
	config           *params.CbftConfig // Consensus engine configuration parameters
	dpos             *dpos
	rotating         *rotating
	blockSignatureCh chan *types.BlockSignature
	cbftResultCh     chan *types.CbftResult
	closeOnce        sync.Once       // Ensures exit channel will not be closed twice.
	exitCh           chan chan error // Notification channel to exiting backend threads

	//blockNumGenerator uint64   		//块高生成器
	highestLogicalBlock *types.Block //区块块号最高的合理块
	masterTree          *Tree
	slaveTree           *Tree
	signCounterMap      map[common.Hash]*SignCounter //签名计数器Map
	lock                sync.RWMutex                 //保护LogicalChainTree
}

type Tree struct {
	nodeMap map[common.Hash]*Node
	root    *Node
}
type Node struct {
	block     *types.Block
	isLogical bool
	children  []*Node
	parent    *Node
}

type CauseType uint

const (
	RcvBlock CauseType = 1 << iota //收到新区块
	RcvSign                        //收到新签名
)

//签名计数器
type SignCounter struct {
	blockNum   uint64                //区块高度
	counter    uint                  //区块签名计数器
	updateTime time.Time             //签名计数器最新更新时间
	signs      map[[65]byte]struct{} //签名map，key=签名
	signedByMe bool                  //本节点是否签名过
}

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, blockSignatureCh chan *types.BlockSignature, cbftResultCh chan *types.CbftResult) *Cbft {
	_dpos := newDpos(config.InitialNodes)

	conf := *config
	if conf.Epoch == 0 {
		conf.Epoch = epochLength
	}

	_masterRoot := &Node{
		isLogical: true,
		children:  make([]*Node, 0),
		parent:    nil,
	}
	_masterTree := &Tree{
		nodeMap: make(map[common.Hash]*Node),
		root:    _masterRoot,
	}

	_slaveRoot := &Node{
		isLogical: false,
		children:  make([]*Node, 0),
		parent:    nil,
	}
	_slaveTree := &Tree{
		nodeMap: make(map[common.Hash]*Node),
		root:    _slaveRoot,
	}

	return &Cbft{
		config:           &conf,
		dpos:             _dpos,
		rotating:         newRotating(_dpos, 10000),
		blockSignatureCh: blockSignatureCh,
		cbftResultCh:     cbftResultCh,

		masterTree:     _masterTree,
		slaveTree:      _slaveTree,
		signCounterMap: make(map[common.Hash]*SignCounter),
	}
}

func (cbft *Cbft) ShouldSeal() (bool, error) {
	return cbft.inTurn(), nil
}

func (cbft *Cbft) ConsensusNodes() ([]discover.Node, error) {
	return cbft.dpos.primaryNodeList, nil
}

func (cbft *Cbft) CheckConsensusNode(nodeID discover.NodeID) (bool, error) {
	return cbft.dpos.NodeIndex(nodeID) >= 0, nil
}

func (cbft *Cbft) IsConsensusNode() (bool, error) {
	return cbft.dpos.NodeIndex(cbft.config.NodeID) >= 0, nil
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (cbft *Cbft) Author(header *types.Header) (common.Address, error) {
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
	if header.Number == nil {
		return errUnknownBlock
	}
	// Don't waste time checking blocks from the future
	if header.Time.Cmp(big.NewInt(time.Now().Unix())) > 0 {
		return consensus.ErrFutureBlock
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
	return cbft.verifySeal(chain, header, nil)
	//return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (b *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// 完成Header对象的准备
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Difficulty = nil
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given, and returns the final block.
func (cbft *Cbft) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// 生成具体的区块信息
	// 填充上Header.Root, TxHash, ReceiptHash, UncleHash等几个属性
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}

// 完成对区块的签名成功，并设置到header.Extra中，然后把区块发送到sealResultCh通道中（然后会被组播到其它共识节点）
func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) error {
	header := block.Header()
	number := header.Number.Uint64()

	if number == 0 {
		// 不密封创世块
		return errUnknownBlock
	}

	if cbft.config.Period == 0 && len(block.Transactions()) == 0 {
		// 不支持0-period的链，不支持空块密封
		log.Info("Sealing paused, waiting for transactions")
		return nil
	}

	//不是合法共识节点
	if ok, _ := cbft.dpos.IsConsensusNode(); !ok {
		return errUnauthorizedSigner
	}

	//不在出块的时间窗口内
	if !cbft.inTurn() {
		log.Info("Not my turn")
		return nil
	}

	//todo:
	//检验区块难度

	// 核心工作：开始签名。注意，delay的不是签名，而是结果的返回
	sighash, err := cbft.signFn(sigHash(header).Bytes())
	if err != nil {
		return err
	}
	//将签名结果替换区块头的Extra字段（专门支持记录额外信息的）
	copy(header.Extra[len(header.Extra)-extraSeal:], sighash)

	go func() {
		select {
		case <-stopCh: //如果先收到stop（客户端RPC发出)，则直接返回
			return
		case sealResultCh <- block.WithSeal(header): //有接受才能发送数据，去执行区块
		default: //如果没有接收数据，则走default
			log.Warn("Sealing result is not read by miner", "sealhash", cbft.SealHash(header))
		}
	}()
	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func (b *Cbft) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (b *Cbft) SealHash(header *types.Header) common.Hash {
	return consensus.SigHash(header)
	//return sigHash(header)
}

// Close implements consensus.Engine. It's a noop for clique as there is are no background threads.
func (cbft *Cbft) Close() error {
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
	return []rpc.API{{
		Namespace: "cbft",
		Version:   "1.0",
		Service:   &API{chain: chain, cbft: cbft},
		Public:    false,
	}}
}

//收到新的区块签名
func (cbft *Cbft) OnBlockSignature(chain consensus.ChainReader, sig *types.BlockSignature) {
	signCounter := cbft.addSignCounter(sig.Hash, sig.Number.Uint64(), sig.Signature, false)
	if signCounter >= 15 {
		node, exists := cbft.masterTree.nodeMap[sig.Hash]
		if exists {
			//如果这个hash对应的区块，已经在masterTree中，
			//这个区块所在的节点，将成为masterTree的新的根节点
			//当这个节点是isLogical==false时，需要重新设定合理节点路径（重新设定的合理节点，需要补签名码？）

			if !node.isLogical {
				tempNode = nil
				cbft.findHighestNode(node)
				highestNode := copyPointer(tempNode)

				//设置最高合理区块
				cbft.highestLogicalBlock = highestNode.block

				//重新设定合理节点路径（重新设定的合理节点，如果没有签名过，则需要补签名,广播签名）
				for highestNode.parent.block.Hash() != node.block.Hash() {
					//设定合理节点
					highestNode.isLogical = true

					signCounter, exists := cbft.signCounterMap[highestNode.block.Hash()]
					if exists {
						if !signCounter.signedByMe {
							//如果没有签名过，则需要补签名,广播签名
							cbft.signNode(highestNode)
						}
					} else {
						log.Warn("cannot find SignCounter for block:", highestNode)
					}
					highestNode = highestNode.parent
				}
			}
			//那么这个区块所在节点到根的所有区块，都可以写入链了
			cbft.storeConfirmed(node, RcvSign)
		}
	}
}

//收到新的区块
func (cbft *Cbft) OnNewBlock(chain consensus.ChainReader, rcvBlock *types.Block) error {
	rcvHeader := rcvBlock.Header()
	rcvNumber := rcvHeader.Number.Uint64()
	if rcvNumber <= 0 {
		return nil
	}

	//从签名恢复出出块人地址
	nodeID, rcvSign, err := ecrecover(rcvHeader)
	if err != nil {
		return err
	}
	//收到的新块中，包含着出块人的一个签名，所以签名数量+1,
	cbft.addSignCounter(rcvBlock.Hash(), rcvNumber, rcvSign, false)

	//检查块是否在出块人的时间窗口内生成的
	//时间合法性计算，不合法返回error
	if cbft.isOverdue(rcvHeader.Time.Int64(), nodeID) {
		return errOverdueBlock
	}
	masterParent, hasMasterParent, err := queryParent(cbft.masterTree.root, rcvHeader)
	if err != nil {
		return err
	}

	//可以加入masterTree
	if hasMasterParent {
		//新块缺省被认为：不是合理块
		//合理时间窗口内出的块，则此时暂时可认为新块：是合理块
		isLogical := true
		if masterParent.isLogical {
			for _, child := range masterParent.children {
				if child.isLogical {
					//如果父块是合理块，而且父块已经有合理子块，则新块被认为：不是合理块
					isLogical = false
					break
				}
			}
		}

		//用新块构建masterTree节点,暂时此节点的父节点=nil
		node := &Node{
			block:     rcvBlock,
			isLogical: isLogical,
			children:  make([]*Node, 0),
			parent:    nil,
		}

		//从slave树中，嫁接可能的子树到node上，根就是node节点
		cbft.graftingFromSlaveTree(node)

		if node.isLogical {
			//如果这棵树的根是合理的，则从slave里嫁接过来的树，有一条支也是合理的，需要补签名
			tempNode = nil
			cbft.findHighestNode(node)
			highestNode := copyPointer(tempNode)

			//设置最高合理区块
			cbft.highestLogicalBlock = highestNode.block

			//设置本节点出块块高
			//cbft.blockNumGenerator = highestNode.block.Number().Uint64()

			//设置一条合理节点路径
			for highestNode.parent != nil {
				highestNode.isLogical = true
				highestNode = highestNode.parent
			}
		}
		//执行子树中的区块，如果区块是合理的，还需要签名并广播
		cbft.recursionESOnNewBlock(node)

		//正式接入masterTree
		node.parent = masterParent

		//查找子树node是否有可以写入链的块
		tempNode = nil
		cbft.findConfirmedAndHighestNode(node)
		if tempNode != nil {
			newRoot := copyPointer(tempNode)
			cbft.storeConfirmed(newRoot, RcvBlock)
		}

	} else {
		//其它情况，把块放入slave树中，不需要执行，也不需要签名
		slaveParent, hasSlaveParent, err := queryParent(cbft.slaveTree.root, rcvHeader)
		if err != nil {
			return err
		}

		if !hasSlaveParent {
			slaveParent = cbft.slaveTree.root
		}
		node := &Node{
			block:     rcvBlock,
			isLogical: false,
			children:  make([]*Node, 0),
			parent:    slaveParent,
		}
		slaveParent.children = append(slaveParent.children, node)
	}

	return nil
}

func (cbft *Cbft) signNode(node *Node) {
	//签名

	sign, err := cbft.signFn(sigHash(node.block.Header()).Bytes())
	if err == nil {
		//块签名计数器+1
		cbft.addSignCounter(node.block.Hash(), node.block.Number().Uint64(), sign, true)
		//广播签名
		blockSign := &types.BlockSignature{
			Hash:      node.block.Hash(),
			Number:    node.block.Number(),
			Signature: sign,
		}
		cbft.blockSignatureCh <- blockSign

	} else {
		log.Warn("sign the received block error", err)
	}
}

//E:Execute
//S:Sign
// 执行这棵子树的所有节点，如果节点是isLogical=true，还需要签名并广播签名
func (cbft *Cbft) recursionESOnNewBlock(node *Node) {
	//todo:执行

	if node.isLogical {
		//签名
		cbft.signNode(node)
	}
	for _, child := range node.children {
		cbft.recursionESOnNewBlock(child)
	}
}

func copyPointer(node *Node) *Node {
	address := *node
	return &address
}

//保存确认块
func (cbft *Cbft) storeConfirmed(newRoot *Node, cause CauseType) {

	cbft.lock.Lock()

	confirmedBlocks := make([]*types.Block, 1)
	confirmedBlocks = append(confirmedBlocks, newRoot.block)

	for newRoot.parent != nil {
		newRoot = newRoot.parent
		confirmedBlocks = append(confirmedBlocks, newRoot.block)
	}
	//去掉原来的root
	confirmedBlocks = confirmedBlocks[:len(confirmedBlocks)-1]

	//反转slice，按顺序把区块写入链
	if len(confirmedBlocks) > 1 {
		for i, j := 0, len(confirmedBlocks)-1; i < j; i, j = i+1, j-1 {
			confirmedBlocks[i], confirmedBlocks[j] = confirmedBlocks[j], confirmedBlocks[i]
		}
	}

	//todo:考虑cbftResultCh改成[]types.Block
	for _, block := range confirmedBlocks {
		cbft.cbftResultCh <- block
	}
	//saveBlocks(confirmedBlocks)

	//把node作为新的root
	newRoot.parent.children = nil
	newRoot.parent = nil
	cbft.masterTree.root = newRoot

	//重置cbft.masterTree.nodeMap
	cbft.masterTree.nodeMap = map[common.Hash]*Node{}
	cbft.resetNodeMap(cbft.masterTree.root)

	//清理slaveTree
	cbft.cleanSlaveTree()

	//清理signCounter
	cbft.cleanSignCounter()
	cbft.lock.Unlock()
}

//查询树中块高最高节点; 相同块高，取签名数多的节点
func (cbft *Cbft) findHighestNode(subTree *Node) {
	for _, node := range subTree.children {
		signCounter := cbft.getSignCounter(node.block.Hash())
		//找到一个更高的块
		if tempNode == nil || node.block.Number().Uint64() > tempNode.block.Number().Uint64() {
			tempNode = node
		} else if node.block.Number().Uint64() == tempNode.block.Number().Uint64() {
			if signCounter > cbft.getSignCounter(tempNode.block.Hash()) {
				tempNode = node
			}
		}
		cbft.findHighestNode(node)
	}
}

//查询树中，签名数>=15并且块高最高的节点; 相同块高，则取签名数多的节点
var tempNode *Node = nil

func (cbft *Cbft) findConfirmedAndHighestNode(subTree *Node) {
	for _, node := range subTree.children {
		signCounter := cbft.getSignCounter(node.block.Hash())
		if signCounter >= 15 {
			//找到一个更高的确认块
			if tempNode == nil || node.block.Number().Uint64() > tempNode.block.Number().Uint64() {
				tempNode = node
			} else if node.block.Number().Uint64() == tempNode.block.Number().Uint64() {
				if signCounter > cbft.getSignCounter(tempNode.block.Hash()) {
					tempNode = node
				}
			}
		}
		cbft.findConfirmedAndHighestNode(node)
	}
}

//在slaveTree中查某个节点（此节点是masterTree中的某个节点）的子树（必定是masterTree的直接子树），并把此子树嫁接到masterTree中的某个节点上
func (cbft *Cbft) graftingFromSlaveTree(parent *Node) {
	slaveRoot := cbft.slaveTree.root
	if slaveRoot != nil && len(slaveRoot.children) > 0 {
		for idx, sonChild := range slaveRoot.children {
			//找到子树
			if parent.block.Hash() == sonChild.block.ParentHash() && parent.block.Number().Uint64()+1 == sonChild.block.Number().Uint64() {
				//在slaveTree中删除此子树
				slaveRoot.children = append(slaveRoot.children[:idx], slaveRoot.children[idx:]...)
				//子树从新指定父节点
				sonChild.parent = parent
				//父节点中加入此子树
				parent.children = append(parent.children, sonChild)
				return
			}
		}
	}
}

//重置cbft.masterTree.nodeMap
//当有新的确认区块产生后，有可能需要重置cbft.masterTree.nodeMap
func (cbft *Cbft) resetNodeMap(node *Node) {
	if node != nil && len(node.children) > 0 {
		for _, child := range cbft.masterTree.root.children {
			cbft.masterTree.nodeMap[child.block.Hash()] = child
			cbft.resetNodeMap(child)
		}
	}
}

//清除cbft.slaveTree,把块高 <= cbft.blockNumGenerator的节点清除掉；如果清除掉的节点还有子树，则把子树接到cbft.slaveTree根节点上
func (cbft *Cbft) cleanSlaveTree() {
	root := cbft.masterTree.root
	if root != nil && len(root.children) > 0 {
		for idx, sonChild := range root.children {
			if sonChild.block.Number().Uint64() <= cbft.highestLogicalBlock.Number().Uint64() {
				//从root删除儿子
				root.children = append(root.children[:idx], root.children[idx:]...)
				//在root里加入孙子(提升孙子作为儿子）
				root.children = append(root.children, sonChild.children...)
				for _, grandChild := range sonChild.children {
					//孙子节点指向root
					grandChild.parent = root
				}
				//删除儿子节点
				sonChild = nil
			}
		}
	}
}

//清理signCounter，清理块高低于masterTree.root块高的签名计数器数据
func (cbft *Cbft) cleanSignCounter() {
	root := cbft.masterTree.root
	rootBlockNum := root.block.Number().Uint64()

	keysDeleted := make([]common.Hash, 0)
	for hash, signCounter := range cbft.signCounterMap {
		if signCounter.blockNum <= rootBlockNum {
			keysDeleted = append(keysDeleted, hash)
		}
	}
	for _, key := range keysDeleted {
		delete(cbft.signCounterMap, key)
	}
}

//签名计数器，保存块收到的签名以及总数；自己出的块也需要增加签名计数器
func (cbft *Cbft) addSignCounter(blockHash common.Hash, blockNum uint64, sign []byte, signedByMe bool) uint {
	signCounter, exists := cbft.signCounterMap[blockHash]

	var signArray [65]byte
	copy(signArray[:], sign[:])

	if exists {
		signCounter.counter = signCounter.counter + 1
		cbft.signCounterMap[blockHash].signs[signArray] = struct{}{}
		signCounter.updateTime = time.Now()
		return signCounter.counter
	} else {
		cbft.signCounterMap[blockHash] = &SignCounter{
			blockNum:   blockNum,
			counter:    1,
			updateTime: time.Now(),
			signs:      make(map[[65]byte]struct{}),
			signedByMe: signedByMe,
		}
		cbft.signCounterMap[blockHash].signs[signArray] = struct{}{}
		return 1
	}
}

//查询签名计数器
func (cbft *Cbft) getSignCounter(blockHash common.Hash) uint {
	signCounter, exists := cbft.signCounterMap[blockHash]
	if exists {
		return signCounter.counter
	} else {
		return 0
	}
}

//获取区块的所有签名
func (cbft *Cbft) getSigns(blockHash common.Hash) map[[65]byte]struct{} {
	signCounter, exists := cbft.signCounterMap[blockHash]
	if exists {
		return signCounter.signs
	} else {
		return nil
	}
}

//查询root开始的数中，
func queryParent(root *Node, rcvHeader *types.Header) (*Node, bool, error) {
	if root.children != nil && len(root.children) > 0 {
		for _, node := range root.children {
			if node.block.Hash() == rcvHeader.ParentHash {
				if node.block.Number().Uint64()+1 == rcvHeader.Number.Uint64() {
					return node, true, nil
				} else {
					return nil, false, errBlockNumber
				}
			} else {
				return queryParent(node, rcvHeader)
			}
		}
	}
	return nil, false, nil
}

//是否签过当前块高的其他区块
func hasSameBlockNumInMaster(root *Node, header *types.Header) bool {
	if root.children != nil && len(root.children) > 0 {
		for _, node := range root.children {
			if node.isLogical {
				if node.block.Number().Uint64() == header.Number.Uint64() {
					return true
				} else {
					return hasSameBlockNumInMaster(node, header)
				}
			}
		}
	}
	return false
}

//出块时间窗口期与出块节点匹配
func (cbft *Cbft) inTurn() bool {
	singerIdx := cbft.dpos.NodeIndex(cbft.config.NodeID)
	if singerIdx >= 0 {
		value1 := int64(singerIdx*10*1000 - int(cbft.config.MaxLatency/3))

		value2 := (time.Now().Unix() - cbft.dpos.StartTimeOfEpoch()) * 1000 % int64(epochLength)

		value3 := int64((singerIdx+1)*10*1000 - int(cbft.config.MaxLatency*2/3))

		if value2 > value1 && value3 > value2 {
			return true
		}
	}
	return false
}

//收到新的区块后，检查新区块的时间合法性
func (cbft *Cbft) isOverdue(blockTimeInSecond int64, nodeID discover.NodeID) bool {
	singerIdx := cbft.dpos.NodeIndex(nodeID)

	round := time.Now().Unix() - cbft.dpos.StartTimeOfEpoch() - int64(10*(singerIdx+1))/210

	deadline := cbft.dpos.StartTimeOfEpoch() + 210*round + int64(10*(singerIdx+1))

	deadline = deadline + int64(cbft.config.MaxLatency*cbft.config.LegalCoefficient)

	if deadline < time.Now().Unix() {
		return true
	}
	return false
}

func ecrecover(header *types.Header) (discover.NodeID, []byte, error) {
	// Retrieve the signature from the header extra-data

	var nodeID discover.NodeID
	if len(header.Extra) < extraSeal {
		return nodeID, []byte{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(sigHash(header).Bytes(), signature)
	if err != nil {
		return nodeID, []byte{}, err
	}

	//转成discover.NodeID
	nodeID, err = discover.BytesID(pubkey)
	if err != nil {
		return nodeID, []byte{}, err
	}
	return nodeID, signature, nil
}

func sigHash(header *types.Header) (hash common.Hash) {
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
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
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

func (cbft *Cbft) signFn(headerHash []byte) (signature []byte, err error) {
	return crypto.Sign(headerHash, cbft.config.PrivateKey)
}
