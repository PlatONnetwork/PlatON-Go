// Package bft implements the BFT consensus engine.
package cbft

import (
	"Platon-go/common"
	"Platon-go/consensus"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/params"
	"Platon-go/rpc"
	"errors"
	"math/big"
	"time"
)

var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errInvalidCheckpointBeneficiary is returned if a checkpoint/epoch transition
	// block has a beneficiary set to non-zeroes.
	errInvalidCheckpointBeneficiary = errors.New("beneficiary in checkpoint block non-zero")
	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")
	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")
)
var (
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
)

type Cbft struct {
	config           *params.CbftConfig // Consensus engine configuration parameters
	dpos             *dpos
	rotating         *rotating
	blockSignatureCh chan *types.BlockSignature
	cbftResultCh     chan *types.Block
}

// New creates a concurrent BFT consensus engine
func New(config *params.CbftConfig, blockSignatureCh chan *types.BlockSignature, cbftResultCh chan *types.Block) *Cbft {
	_dpos := newDpos(config.InitialNodes)
	return &Cbft {
		dpos:              _dpos,
		rotating :         newRotating(_dpos, 10000),
		blockSignatureCh : blockSignatureCh,
		cbftResultCh :     cbftResultCh,
	}
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (b *Cbft) Author(header *types.Header) (common.Address, error) {
	// 返回出块节点对应的矿工钱包地址
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (b *Cbft) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
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
	// Ensure that the block doesn't contain any uncles which are meaningless in PoA
	if header.UncleHash != uncleHash {
		return errInvalidUncleHash
	}

	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (b *Cbft) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return nil, nil
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (b *Cbft) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return nil
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (b *Cbft) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// VerifySeal()函数基于跟Seal()完全一样的算法原理
	// 通过验证区块的某些属性(Header.Nonce，Header.MixDigest等)是否正确，来确定该区块是否已经经过Seal操作
	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (b *Cbft) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// 完成Header对象的准备
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given, and returns the final block.
func (b *Cbft) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// 生成具体的区块信息
	// 填充上Header.Root, TxHash, ReceiptHash, UncleHash等几个属性
	return nil, nil
}

// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
func (b *Cbft) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// 打包区块
	// 对传入的Block进行最终的授权
	// Seal()函数可对一个调用过Finalize()的区块进行授权或封印，并将封印过程产生的一些值赋予区块中剩余尚未赋值的成员(Header.Nonce, Header.MixDigest)。
	// Seal()成功时返回的区块全部成员齐整，可视为一个正常区块，可被广播到整个网络中，也可以被插入区块链等。
	// 所以，对于挖掘一个新区块来说，所有相关代码里Engine.Seal()是其中最重要，也是最复杂的一步
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
}

// Close implements consensus.Engine. It's a noop for clique as there is are no background threads.
func (b *Cbft) Close() error {
	return nil
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (b *Cbft) APIs(chain consensus.ChainReader) []rpc.API {
	return nil
}

func (b *Cbft) OnBlockSignature(chain consensus.ChainReader, sig *types.BlockSignature) error {
	return nil
}

func (b *Cbft) OnNewBlock(chain consensus.ChainReader, block *types.Block) error {
	return nil
}