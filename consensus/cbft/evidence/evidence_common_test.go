package evidence

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
)

type NodeData struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    common.Address
	nodeID     discover.NodeID
	index      int
}

func newBlock(blockNumber int64) *types.Block {
	header := &types.Header{
		Number:      big.NewInt(blockNumber),
		ParentHash:  common.BytesToHash(utils.Rand32Bytes(32)),
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 77),
		ReceiptHash: common.BytesToHash(utils.Rand32Bytes(32)),
		Root:        common.BytesToHash(utils.Rand32Bytes(32)),
	}
	block := types.NewBlockWithHeader(header)
	return block
}

func createAccount(n int) []*ecdsa.PrivateKey {
	var pris []*ecdsa.PrivateKey
	for i := 0; i < n; i++ {
		pri, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		pris = append(pris, pri)
	}
	return pris
}

func makePrepareBlock(epoch, viewNumber uint64, block *types.Block, blockIndex uint32, ProposalIndex uint32) *protocols.PrepareBlock {
	p := &protocols.PrepareBlock{
		Epoch:         epoch,
		ViewNumber:    viewNumber,
		Block:         block,
		BlockIndex:    blockIndex,
		ProposalIndex: ProposalIndex,
	}

	return p
}

func makePrepareVote(epoch, viewNumber uint64, blockHash common.Hash, blockNumber uint64, blockIndex uint32, validatorIndex uint32) *protocols.PrepareVote {
	p := &protocols.PrepareVote{
		Epoch:          epoch,
		ViewNumber:     viewNumber,
		BlockHash:      blockHash,
		BlockNumber:    blockNumber,
		BlockIndex:     blockIndex,
		ValidatorIndex: validatorIndex,
	}

	return p
}

func makeViewChange(epoch, viewNumber uint64, blockHash common.Hash, blockNumber uint64, validatorIndex uint32) *protocols.ViewChange {
	p := &protocols.ViewChange{
		Epoch:          epoch,
		ViewNumber:     viewNumber,
		BlockHash:      blockHash,
		BlockNumber:    blockNumber,
		ValidatorIndex: validatorIndex,
	}

	return p
}
