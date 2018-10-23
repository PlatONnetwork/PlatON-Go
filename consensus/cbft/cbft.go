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
	//epochLength = uint64(210000) //所有共识节点完成一轮出块的时间，21个节点，每个节点10秒出块时间，共2100000毫秒的时间。
	extraSeal = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
)

type Cbft struct {
	config           *params.CbftConfig // Consensus engine configuration parameters
	dpos             *dpos
	rotating         *rotating
	blockSignatureCh chan *cbfttypes.BlockSignature
	cbftResultCh     chan *cbfttypes.CbftResult
	closeOnce        sync.Once       // Ensures exit channel will not be closed twice.
	exitCh           chan chan error // Notification channel to exiting backend threads

	blockChain          *core.BlockChain              //区块链指针
	highestLogicalBlock *types.Block                  //区块块号最高的合理块
	masterTree          *Tree                         //主树，根节点含有最近的不可逆区块
	slaveTree           *Tree                         //副树，根节点没有实际意义，不包含区块信息
	signCacheMap        map[common.Hash]*SignCache    //签名Map
	receiptCacheMap     map[common.Hash]*ReceiptCache //块执行后的回执Map
	stateCacheMap       map[common.Hash]*StateCache   //块执行后的状态Map
	lock                sync.RWMutex                  //保护LogicalChainTree
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

//签名缓存
type SignCache struct {
	blockNum   uint64                     //区块高度
	counter    uint                       //区块签名计数器
	updateTime time.Time                  //签名计数器最新更新时间
	signs      []*common.BlockConfirmSign //签名map，key=签名
	signedByMe bool                       //本节点是否签名过
}

//收据缓存
type ReceiptCache struct {
	blockNum uint64         //区块高度
	receipts types.Receipts //执行区块后的收据
}

//state缓存
type StateCache struct {
	blockNum uint64         //区块高度
	state    *state.StateDB //执行区块后的状态
}

func (cbft *Cbft) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	cbft.config.PrivateKey = privateKey
	cbft.config.NodeID = discover.PubkeyID(&privateKey.PublicKey)
}

var cbft *Cbft

func SetBlockChain(blockChain *core.BlockChain) {
	log.Info("初始化cbft.blockChain")

	cbft.blockChain = blockChain
	cbft.dpos.SetStartTimeOfEpoch(blockChain.Genesis().Time().Int64())

	currentBlock := blockChain.CurrentBlock()

	genesisParentHash := bytes.Repeat([]byte{0x00}, 32)
	if bytes.Equal(currentBlock.ParentHash().Bytes(), genesisParentHash) && currentBlock.Number() == nil {
		currentBlock.Header().Number = big.NewInt(0)
		currentBlock.Number()
	}

	log.Info("初始化cbft.highestLogicalBlock", "hash", currentBlock.Hash().String(), "number", currentBlock.NumberU64())

	cbft.highestLogicalBlock = currentBlock

	_masterRoot := &Node{
		isLogical: true,
		block:     currentBlock,
		children:  make([]*Node, 0),
		parent:    nil,
	}
	_masterRoot.block = currentBlock

	_masterTree := &Tree{
		nodeMap: make(map[common.Hash]*Node),
		root:    _masterRoot,
	}

	_masterTree.nodeMap[currentBlock.Hash()] = _masterRoot

	cbft.masterTree = _masterTree
}

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, blockSignatureCh chan *cbfttypes.BlockSignature, cbftResultCh chan *cbfttypes.CbftResult) *Cbft {
	_dpos := newDpos(config.InitialNodes)

	_slaveRoot := &Node{
		isLogical: false,
		children:  make([]*Node, 0),
		parent:    nil,
	}
	_slaveTree := &Tree{
		nodeMap: make(map[common.Hash]*Node),
		root:    _slaveRoot,
	}

	cbft = &Cbft{
		config:           config,
		dpos:             _dpos,
		rotating:         newRotating(_dpos, config.Duration),
		blockSignatureCh: blockSignatureCh,
		cbftResultCh:     cbftResultCh,

		slaveTree:       _slaveTree,
		signCacheMap:    make(map[common.Hash]*SignCache),
		receiptCacheMap: make(map[common.Hash]*ReceiptCache),
		stateCacheMap:   make(map[common.Hash]*StateCache),
	}
	return cbft
}

func (cbft *Cbft) ShouldSeal() (bool, error) {
	return cbft.inTurn(), nil
}

func (cbft *Cbft) ConsensusNodes() ([]discover.Node, error) {
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
	log.Info("call Author(), parameter: ", "headerHash", header.Hash().String(), "headerNumber", header.Number)

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
	log.Info("call VerifyHeader(), parameter: ", "headerHash", header.Hash().String(), "headerNumber", header.Number, "seal", seal)

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
	log.Info("call VerifySeal(), parameter: ", "headerHash", header.Hash().String(), "headerNumber", header.Number.String())

	return cbft.verifySeal(chain, header, nil)
	//return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (b *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	log.Info("call Prepare(), parameter: ", "headerHash", header.Hash().String(), "headerNumber", header.Number.String())

	//检查父区块
	if header.ParentHash != cbft.highestLogicalBlock.Hash() || header.Number.Uint64()-1 != cbft.highestLogicalBlock.NumberU64() {
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
	log.Info("call Finalize(), parameter: ", "headerHash", header.Hash().String(), "headerNumber", header.Number.String(), "state", state, "txs", txs, "uncles", uncles, "receipts", receipts)

	// 生成具体的区块信息
	// 填充上Header.Root, TxHash, ReceiptHash, UncleHash等几个属性
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	return types.NewBlock(header, txs, nil, receipts), nil
}

// 完成对区块的签名成功，并设置到header.Extra中，然后把区块发送到sealResultCh通道中（然后会被组播到其它共识节点）
func (cbft *Cbft) Seal(chain consensus.ChainReader, block *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) error {
	log.Info("call Seal(), parameter", "block", block)

	header := block.Header()
	number := header.Number.Uint64()

	if number == 0 {
		// 不密封创世块
		return errUnknownBlock
	}

	// 核心工作：开始签名。注意，delay的不是签名，而是结果的返回
	sign, err := cbft.signFn(signHash(header).Bytes())
	if err != nil {
		return err
	}
	//将签名结果替换区块头的Extra字段（专门支持记录额外信息的）
	copy(header.Extra[len(header.Extra)-extraSeal:], sign[:])

	newBlock := block.WithSeal(header)

	parentNode := cbft.masterTree.nodeMap[newBlock.ParentHash()]
	if parentNode == nil {
		log.Error("找不到父节点", "blockHash", newBlock.Hash().String(), "parentHash", newBlock.ParentHash().String())
		return errUnknownBlock
	}

	//把新节点加入masterTree
	cbft.addBlockToMasterTree(parentNode, newBlock)

	//把当前新块作为最高区块
	cbft.highestLogicalBlock = newBlock

	go func() {
		select {
		case <-stopCh: //如果先收到stop（客户端RPC发出)，则直接返回
			return
		case sealResultCh <- newBlock: //发送给p2p，把区块广播到其它节点
		default: //如果没有接收数据，则走default
			log.Warn("Sealing result is not read by miner", "sealhash", cbft.SealHash(header).String())
		}
	}()
	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func (b *Cbft) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	log.Info("call CalcDifficulty(), parameter", "time", time, "parentHash", parent.Hash().String(), "parentNumber", parent.Number.String())

	return big.NewInt(2)
}

// SealHash returns the hash of a block prior to it being sealed.
func (b *Cbft) SealHash(header *types.Header) common.Hash {
	log.Info("call SealHash(), parameter", "headerHash", header.Hash().String(), "headerNumber", header.Number.String())

	//return consensus.SigHash(header)
	return signHash(header)
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
func (cbft *Cbft) OnBlockSignature(chain consensus.ChainReader, nodeID discover.NodeID, sig *cbfttypes.BlockSignature) error {
	log.Info("收到新的区块签名==>>, parameter", "nodeID", nodeID.String(), "blockHash", sig.Hash, "signHash", sig.SignHash, "sigNUmber", sig.Number, "sig", sig.Signature.String())

	ok, err := verifySign(nodeID, sig.SignHash, sig.Signature[:])
	if err != nil {
		return err
	}

	if !ok {
		log.Error("unauthorized signer")
		return errUnauthorizedSigner
	}

	signCounter := cbft.addSign(sig.Hash, sig.Number.Uint64(), sig.Signature, false)
	log.Info("区块签名数量", "hash", sig.Hash, "signCounter", signCounter, "共识成功阈值", cbft.getThreshold())

	if signCounter >= cbft.getThreshold() {
		//区块收到的签名数量>=15，可以入链了
		log.Info("收到的签名数量>2f+1", "signCounter", signCounter)
		node, exists := cbft.masterTree.nodeMap[sig.Hash]
		if exists {
			//如果这个hash对应的区块，已经在masterTree中，
			//这个区块所在的节点，将成为masterTree的新的根节点
			//当这个节点是isLogical==false时，需要重新设定合理节点路径（重新设定的合理节点，需要补签名码？）
			log.Info("签名对应的区块已经在masterTree中")
			if !node.isLogical {
				log.Info("签名对应的区块不是合理节点，需要先转成合理节点")

				tempNode = nil

				cbft.findHighestNode(node)

				//tempNode至少是node
				highestNode := tempNode

				//设置最高合理区块
				cbft.highestLogicalBlock = highestNode.block

				//重新设定合理节点路径（重新设定的合理节点，如果没有签名过，则需要补签名,广播签名）
				// highestNode.parent != nil 可以保证 highestNode.parent.block != nil。有节点，就有block
				for highestNode.parent != nil && highestNode.parent.block.Hash() != node.block.Hash() {
					//设定合理节点
					highestNode.isLogical = true

					signCounter, exists := cbft.signCacheMap[highestNode.block.Hash()]
					if exists {
						if !signCounter.signedByMe {
							//如果没有签名过，则需要补签名,广播签名
							cbft.signNode(highestNode)
						}
					} else {
						log.Warn("cannot find SignCounter for block:", "highestNode", highestNode)
					}
					highestNode = highestNode.parent
				}
			}
			//那么这个区块所在节点到根的所有区块，都可以写入链了
			cbft.storeConfirmed(node, RcvSign)
		}
	}
	return nil
}

//收到新的区块
func (cbft *Cbft) OnNewBlock(chain consensus.ChainReader, rcvBlock *types.Block) error {
	log.Info("收到新的区块==>>, parameter", "rcvBlockHash", rcvBlock.Hash().String(), "rcvBlockNumber", rcvBlock.Header().Number, "ParentHash", rcvBlock.ParentHash().String(), "rcvBlockExtra", hexutil.Encode(rcvBlock.Header().Extra))

	rcvHeader := rcvBlock.Header()
	rcvNumber := rcvHeader.Number.Uint64()
	if rcvNumber <= 0 {
		return nil
	}

	//从签名恢复出出块人地址·
	nodeID, rcvSign, err := ecrecover(rcvHeader)
	if err != nil {
		return err
	}
	log.Info("收到的新块，出块方：", "nodeID", nodeID.String())

	//收到的新块中，包含着出块人的一个签名，所以签名数量+1,
	sign := common.NewBlockConfirmSign(rcvSign)
	cbft.addSign(rcvBlock.Hash(), rcvNumber, sign, false)

	//检查块是否在出块人的时间窗口内生成的
	//时间合法性计算，不合法返回error
	log.Info("检查块是否在出块人的时间窗口内生成的")
	if cbft.isOverdue(rcvHeader.Time.Int64(), nodeID) {
		return errOverdueBlock
	}

	log.Info("查询新块是否能接上cbft.masterTree")

	masterParent, hasMasterParent, err := queryParent(cbft.masterTree.root, rcvHeader)
	if err != nil {
		return err
	}

	//可以加入masterTree
	if hasMasterParent {
		//新块缺省被认为：不是合理块
		//合理时间窗口内出的块，则此时暂时可认为新块：是合理块
		log.Info("新块可以加入masterTree")
		isLogical := true
		if masterParent.isLogical {
			log.Info("父块是合理块")
			for _, child := range masterParent.children {
				if child.isLogical {
					//如果父块是合理块，而且父块已经有合理子块，则新块被认为：不是合理块
					log.Info("如果父块是合理块，而且父块已经有合理子块，则新块被认为：不是合理块")
					isLogical = false
					break
				}
			}
			log.Info("新块是否是合理块", "isLogical", isLogical)
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
			log.Info("新块是合理块")
			tempNode = nil
			cbft.findHighestNode(node)

			//tempNode至少是node
			highestNode := tempNode

			log.Info("新块开始的临时树中，最高节点是：", "highestNode", highestNode)

			//设置最高合理区块
			log.Info("设置最高合理区块：", "highestNode.block", highestNode.block)
			cbft.highestLogicalBlock = highestNode.block

			//设置本节点出块块高
			//cbft.blockNumGenerator = highestNode.block.Number().Uint64()

			//设置一条合理节点路径 （注意退出循环的条件，这也是创建node时，设置paretn=nil的原因）
			for highestNode.parent != nil {
				highestNode.isLogical = true
				highestNode = highestNode.parent
			}
		}
		//正式接入masterTree
		//需要先接入masterTree，这样，子树的root才能找到parent（执行的时候需要）
		log.Info("新块开始的临时树，正式接入masterTree：")
		cbft.addNodeToMasterTree(masterParent, node)

		//执行子树中的区块，如果区块是合理的，还需要签名并广播
		cbft.recursionESOnNewBlock(node)

		//查找新接入的子树，是否有可以写入链的块
		log.Info("查找新接入的子树，是否有可以写入链的块")
		tempNode = nil
		cbft.findConfirmedAndHighestNode(node)
		if tempNode != nil {

			newRoot := tempNode

			log.Info("新接入的子树，有可以写入链的块", "blockHash", newRoot.block.Hash().String())

			cbft.storeConfirmed(newRoot, RcvBlock)
		}

	} else {

		log.Info("新块只能加入slaveTree")

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
		cbft.masterTree.nodeMap[node.block.Hash()] = node
	}

	return nil
}

func (cbft *Cbft) HighestLogicalBlock() *types.Block {

	log.Info("call HighestLogicalBlock() ...")

	return cbft.highestLogicalBlock
}

func (cbft *Cbft) processNode(node *Node) {
	//执行
	receipts, state, err := cbft.blockChain.ProcessDirectly(node.block, node.parent.block)
	if err == nil {
		if receipts == nil {
			receipts = types.Receipts{}
		}
		receiptsCache := &ReceiptCache{
			blockNum: node.block.NumberU64(),
			receipts: receipts,
		}
		cbft.receiptCacheMap[node.block.Hash()] = receiptsCache

		stateCache := &StateCache{
			blockNum: node.block.NumberU64(),
			state:    state,
		}
		cbft.stateCacheMap[node.block.Hash()] = stateCache

	} else {
		log.Warn("process block error", err)
	}
}

func (cbft *Cbft) signNode(node *Node) uint {
	//签名
	signHash := signHash(node.block.Header())
	signature, err := cbft.signFn(signHash.Bytes())
	if err == nil {
		log.Info("区块签名值", "signHash", signHash.String(), "sign", hexutil.Encode(signature))

		blockHash := node.block.Hash()

		//块签名计数器+1
		sign := common.NewBlockConfirmSign(signature)
		//addSign()的key必须是blok.hash()
		signCounter := cbft.addSign(blockHash, node.block.Number().Uint64(), sign, true)
		//广播签名
		blockSign := &cbfttypes.BlockSignature{
			SignHash:  signHash,
			Hash:      blockHash,
			Number:    node.block.Number(),
			Signature: sign,
		}
		cbft.blockSignatureCh <- blockSign

		return signCounter
	} else {
		panic("can't sign for block")
	}
}

//E:Execute
//S:Sign
// 执行这棵子树的所有节点，如果节点是isLogical=true，还需要签名并广播签名
func (cbft *Cbft) recursionESOnNewBlock(node *Node) {

	//执行
	log.Info("区块执行", "blockHash", node.block.Hash().String())
	cbft.processNode(node)

	if node.isLogical {
		//签名
		log.Info("区块是合理块，需要签名")
		signCounter := cbft.signNode(node)

		if signCounter > cbft.getThreshold() {
			cbft.storeConfirmed(node, RcvBlock)
		}
	}

	log.Info("子块数量：", "lenChildren", len(node.children))
	for _, child := range node.children {
		cbft.recursionESOnNewBlock(child)
	}
}

func (cbft *Cbft) storeConfirmed(newRoot *Node, cause CauseType) {

	log.Info("把masterTree中，从root到node的节点都入链，并重构masterTree")

	cbft.lock.Lock()
	defer cbft.lock.Unlock()

	//保存最初的输入节点（新树的根节点）
	tempNode := newRoot

	confirmedBlocks := make([]*types.Block, 1)
	confirmedBlocks[0] = newRoot.block

	for tempNode.parent != nil {
		tempNode = tempNode.parent
		confirmedBlocks = append(confirmedBlocks, tempNode.block)
	}
	//去掉原来的root
	confirmedBlocks = confirmedBlocks[:len(confirmedBlocks)-1]

	//反转slice，按顺序把区块写入链
	if len(confirmedBlocks) > 1 {
		for i, j := 0, len(confirmedBlocks)-1; i < j; i, j = i+1, j-1 {
			confirmedBlocks[i], confirmedBlocks[j] = confirmedBlocks[j], confirmedBlocks[i]
		}
	}

	//todo:考虑cbftResultCh改成[]types.CbftResult
	for _, block := range confirmedBlocks {

		cbftResult := &cbfttypes.CbftResult{
			Block: block,
			//Receipts:          cbft.receiptCacheMap[block.Hash()].receipts,
			//State:             cbft.stateCacheMap[block.Hash()].state,
			BlockConfirmSigns: cbft.signCacheMap[block.Hash()].signs,
		}

		receiptCache := cbft.receiptCacheMap[block.Hash()]
		if receiptCache != nil {
			cbftResult.Receipts = receiptCache.receipts
		}

		stateCache := cbft.stateCacheMap[block.Hash()]
		if stateCache != nil {
			cbftResult.State = stateCache.state
		}

		//把需要保存的数据，发往通道：cbftResultCh
		cbft.cbftResultCh <- cbftResult
	}

	//把node作为新的root
	newRoot.parent.children = nil
	newRoot.parent = nil
	cbft.masterTree.root = newRoot

	//重置cbft.masterTree.nodeMap
	cbft.masterTree.nodeMap = map[common.Hash]*Node{}
	cbft.resetNodeMap(cbft.masterTree.root)

	//清理slaveTree
	cbft.cleanSlaveTree()

	//清理signCacheMap
	cbft.cleanSignCacheMap()

	//清理receiptCacheMap
	cbft.cleanReceiptCacheMap()

	//清理stateCacheMap
	cbft.cleanStateCacheMap()
}

//查询树中块高最高节点; 相同块高，取签名数多的节点
func (cbft *Cbft) findHighestNode(subTree *Node) {
	if subTree.children == nil || len(subTree.children) == 0 {
		tempNode = subTree
	}
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
	signCounter := cbft.getSignCounter(subTree.block.Hash())
	if signCounter >= cbft.getThreshold() {
		//找到一个更高的确认块
		if tempNode == nil || subTree.block.Number().Uint64() > tempNode.block.Number().Uint64() {
			tempNode = subTree
		} else if subTree.block.Number().Uint64() == tempNode.block.Number().Uint64() {
			if signCounter > cbft.getSignCounter(tempNode.block.Hash()) {
				tempNode = subTree
			}
		}
	}
	for _, node := range subTree.children {
		signCounter := cbft.getSignCounter(node.block.Hash())
		if signCounter >= cbft.getThreshold() {
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
	cbft.masterTree.nodeMap[node.block.Hash()] = node
	if node != nil && len(node.children) > 0 {
		for _, child := range cbft.masterTree.root.children {
			//cbft.masterTree.nodeMap[child.block.Hash()] = child
			cbft.resetNodeMap(child)
		}
	}
}

// 清除cbft.slaveTree,把块高 <= cbft.masterTree.root的节点清除掉
// 处理时，是从cbft.slaveTree.root的儿子层开始的，每次循环只处理儿子节点。
// 当儿子节点<=highLimiting时，把此儿子从cbft.slaveTree.root根上删除，把此儿子的儿子（即cbft.slaveTree.root的孙子）提升为cbft.slaveTree.root的儿子
// 当cbft.slaveTree.root根的所有儿子都不满足<= highLimiting 时，则退出循环
func (cbft *Cbft) cleanSlaveTree() {
	//masterTree根节点区块的块高
	highLimiting := cbft.masterTree.root.block.NumberU64()
	root := cbft.slaveTree.root
	if root != nil && len(root.children) > 0 {
		//退出循环处理标识
		exit := false
		for !exit {
			exit = true
			for idx, sonChild := range root.children {
				if sonChild.block.NumberU64() <= highLimiting {
					exit = false
					//从root删除儿子
					root.children = append(root.children[:idx], root.children[idx+1:]...)
					//在root里加入孙子(提升所有孙子作为儿子）
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
}

//清理signCacheMap，清理块高低于masterTree.root块高的签名计数器数据
func (cbft *Cbft) cleanSignCacheMap() {
	root := cbft.masterTree.root
	rootBlockNum := root.block.Number().Uint64()

	keysDeleted := make([]common.Hash, 0)
	for hash, signCache := range cbft.signCacheMap {
		if signCache.blockNum <= rootBlockNum {
			keysDeleted = append(keysDeleted, hash)
		}
	}
	for _, key := range keysDeleted {
		delete(cbft.signCacheMap, key)
	}
}

//清理receiptCacheMap，清理块高低于masterTree.root块高的签名计数器数据
func (cbft *Cbft) cleanReceiptCacheMap() {
	root := cbft.masterTree.root
	rootBlockNum := root.block.Number().Uint64()

	keysDeleted := make([]common.Hash, 0)
	for hash, receiptCache := range cbft.receiptCacheMap {
		if receiptCache.blockNum <= rootBlockNum {
			keysDeleted = append(keysDeleted, hash)
		}
	}
	for _, key := range keysDeleted {
		delete(cbft.receiptCacheMap, key)
	}
}

//清理stateCacheMap，清理块高低于masterTree.root块高的签名计数器数据
func (cbft *Cbft) cleanStateCacheMap() {
	root := cbft.masterTree.root
	rootBlockNum := root.block.Number().Uint64()

	keysDeleted := make([]common.Hash, 0)
	for hash, stateCache := range cbft.stateCacheMap {
		if stateCache.blockNum <= rootBlockNum {
			keysDeleted = append(keysDeleted, hash)
		}
	}
	for _, key := range keysDeleted {
		delete(cbft.stateCacheMap, key)
	}
}

// 保存区块的签名（自己出的块也需要保存签名）
// 返回区块的签名总数
func (cbft *Cbft) addSign(blockHash common.Hash, blockNum uint64, sign *common.BlockConfirmSign, signedByMe bool) uint {
	signCache, exists := cbft.signCacheMap[blockHash]

	if !exists {
		signCache = &SignCache{
			blockNum:   blockNum,
			counter:    0,
			signs:      make([]*common.BlockConfirmSign, 0),
			signedByMe: signedByMe,
		}
		cbft.signCacheMap[blockHash] = signCache
	}

	signCache.counter = signCache.counter + 1
	log.Info("增加签名的计数器", "hash", blockHash, "counter", signCache.counter)

	signCache.signs = append(signCache.signs, sign)

	signCache.updateTime = time.Now()
	return signCache.counter
}

//查询区块的签名总数
func (cbft *Cbft) getSignCounter(blockHash common.Hash) uint {
	signCounter, exists := cbft.signCacheMap[blockHash]
	if exists {
		return signCounter.counter
	} else {
		return 0
	}
}

//获取区块的所有签名
func (cbft *Cbft) getSigns(blockHash common.Hash) []*common.BlockConfirmSign {
	signCounter, exists := cbft.signCacheMap[blockHash]
	if exists {
		return signCounter.signs
	} else {
		return nil
	}
}

//本节点出的区块，都是本节点认为的合理区块
/*func (cbft *Cbft) addMasterTreeNode(newBlock *types.Block) error {
	parentNode := cbft.masterTree.nodeMap[newBlock.ParentHash()]
	if parentNode == nil {
		log.Error("找不到父节点", "blockHash", newBlock.Hash().String(), "parentHash", newBlock.ParentHash().String())
		return errUnknownBlock
	}
	newNode := &Node{
		block:     newBlock,
		isLogical: true, //合理区块
		children:  make([]*Node, 0),
		parent:    parentNode,
	}

	parentNode.children = append(parentNode.children, newNode)

	cbft.masterTree.nodeMap[newBlock.Hash()] = newNode
	return nil
}
*/

func (cbft *Cbft) addBlockToMasterTree(parent *Node, newBlock *types.Block) {
	newNode := &Node{
		block:     newBlock,
		isLogical: true, //合理区块
		children:  make([]*Node, 0),
		parent:    parent,
	}
	parent.children = append(parent.children, newNode)
	cbft.masterTree.nodeMap[newBlock.Hash()] = newNode
}

func (cbft *Cbft) addNodeToMasterTree(parent *Node, child *Node) {
	child.parent = parent
	parent.children = append(parent.children, child)
	cbft.masterTree.nodeMap[child.block.Hash()] = child
}

//查询root开始的树中，是否有父节点
func queryParent(root *Node, rcvHeader *types.Header) (*Node, bool, error) {

	//slaveTree树的root，直接返回false
	if root.block == nil && root.parent == nil {
		return nil, false, nil
	}

	log.Info("查询root开始的树中，是否有父节点", "rootHash", root.block.Hash().String(), "rootNumber", root.block.Number().Uint64(), "newHash", rcvHeader.Hash().String(), "newNumber", rcvHeader.Number)

	if root.block.Hash() == rcvHeader.ParentHash && root.block.Number().Uint64()+1 == rcvHeader.Number.Uint64() {
		return root, true, nil
	} else {
		if root.children != nil && len(root.children) > 0 {
			for _, child := range root.children {
				return queryParent(child, rcvHeader)

				/*				log.Info("查询root开始的树中，是否有父节点", "childHash", child.block.Hash(), "childNumber", child.block.Number().Uint64(), "newHash", rcvHeader.Hash(), "newNumber", rcvHeader.Number)

								if child.block.Hash() == rcvHeader.ParentHash {
									if child.block.Number().Uint64()+1 == rcvHeader.Number.Uint64() {
										return child, true, nil
									} else {
										return nil, false, errBlockNumber
									}
								} else {
									return queryParent(child, rcvHeader)
								}*/
			}
		}
		return nil, false, nil
	}
}

//出块时间窗口期与出块节点匹配
func (cbft *Cbft) inTurn() bool {
	singerIdx := cbft.dpos.NodeIndex(cbft.config.NodeID)
	if singerIdx >= 0 {
		durationMilliseconds := cbft.config.Duration * 1000
		totalDuration := durationMilliseconds * int64(len(cbft.dpos.primaryNodeList))

		value1 := singerIdx*(durationMilliseconds) - int64(cbft.config.MaxLatency/3)

		value2 := (time.Now().Unix()*1000 - cbft.dpos.StartTimeOfEpoch()) % totalDuration

		value3 := (singerIdx+1)*durationMilliseconds - int64(cbft.config.MaxLatency*2/3)

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
	rounds := (time.Now().Unix() - cbft.dpos.StartTimeOfEpoch()) / totalDuration

	//nodeID的最晚出块时间
	deadline := cbft.dpos.StartTimeOfEpoch() + totalDuration*rounds + durationMilliseconds*(singerIdx+1)

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
	signHash := signHash(header)
	log.Info("收到新块", "signHash", signHash.String())

	pubkey, err := crypto.Ecrecover(signHash.Bytes(), signature)
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

func verifySign(expectedNodeID discover.NodeID, signHash common.Hash, signature []byte) (bool, error) {

	log.Info("验证签名", "signHash", signHash.String(), "signature", hexutil.Encode(signature), "expectedNodeID", hexutil.Encode(expectedNodeID.Bytes()))

	pubkey, err := crypto.SigToPub(signHash.Bytes(), signature)

	if err != nil {
		return false, err
	}

	nodeID := discover.PubkeyID(pubkey)
	//比较两个[]byte
	log.Info("从签名恢复出的NodeID", "nodeID", nodeID.String())
	if bytes.Equal(nodeID.Bytes(), expectedNodeID.Bytes()) {
		return true, nil
	}
	return false, nil
}

func signHash(header *types.Header) (hash common.Hash) {
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
	return cbft.highestLogicalBlock
}

func (cbft *Cbft) getThreshold() uint {
	trunc := len(cbft.dpos.primaryNodeList) * 2 / 3
	remainder := len(cbft.dpos.primaryNodeList) * 2 % 3

	if remainder == 0 {
		return uint(trunc)
	} else {
		return uint(trunc + 1)
	}
}
