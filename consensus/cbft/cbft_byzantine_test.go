// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/executor"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/rules"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/stretchr/testify/assert"
)

const (
	fetchPrepare        = "previous index block not exists"
	noBaseMaxBlock      = "prepareBlock is not based on viewChangeQC maxBlock"
	errorViewChangeQC   = "verify viewchange qc failed"
	MismatchedPrepareQC = "verify prepare qc failed,not the corresponding qc"
	missingViewChangeQC = "prepareBlock need ViewChangeQC"
	dupBlockHash        = "has duplicated blockHash"
	errorSignature      = "bls verifies signature fail"
	enableVerifyEpoch   = "enable verify epoch"
)

func MockNodes(t *testing.T, num int) []*TestCBFT {
	pk, sk, nodes := GenerateCbftNode(num)
	engines := make([]*TestCBFT, 0)

	for i := 0; i < num; i++ {
		e := MockNode(pk[i], sk[i], nodes, 10000, 10)
		assert.Nil(t, e.Start())
		engines = append(engines, e)
	}
	return engines
}

func ReachBlock(t *testing.T, nodes []*TestCBFT, reach int) {
	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()
	for i := 0; i < reach; i++ {
		block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil, complete)
		<-complete

		_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i-1), nodes[0].engine.state.MaxQCIndex())
			for j := 1; j < len(nodes); j++ {
				msg := &protocols.PrepareVote{
					Epoch:          nodes[0].engine.state.Epoch(),
					ViewNumber:     nodes[0].engine.state.ViewNumber(),
					BlockIndex:     uint32(i),
					BlockHash:      b.Hash(),
					BlockNumber:    b.NumberU64(),
					ValidatorIndex: uint32(j),
					ParentQC:       qc,
				}
				assert.Nil(t, nodes[j].engine.signMsgByBls(msg))
				assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
			}
			parent = b
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func ChangeView(t *testing.T, nodes []*TestCBFT, prepareQC *ctypes.QuorumCert) *ctypes.ViewChangeQC {
	for i := 0; i < len(nodes); i++ {
		epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()
		viewchange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      prepareQC.BlockHash,
			BlockNumber:    prepareQC.BlockNumber,
			ValidatorIndex: uint32(i),
			PrepareQC:      prepareQC,
		}
		assert.Nil(t, nodes[i].engine.signMsgByBls(viewchange))
		assert.Nil(t, nodes[0].engine.OnViewChanges("id", &protocols.ViewChanges{
			VCs: []*protocols.ViewChange{
				viewchange,
			},
		}))
	}
	assert.NotNil(t, nodes[0].engine.state.LastViewChangeQC())
	assert.Equal(t, uint64(1), nodes[0].engine.state.ViewNumber())
	return nodes[0].engine.state.LastViewChangeQC()
}

func FakeViewChangeQC(t *testing.T, node *TestCBFT, epoch, viewNumber uint64, nodeIndex uint32, prepareQC *ctypes.QuorumCert) *ctypes.ViewChangeQC {
	viewChanges := make(map[uint32]*protocols.ViewChange)
	v := &protocols.ViewChange{
		Epoch:          epoch,
		ViewNumber:     viewNumber,
		BlockHash:      prepareQC.BlockHash,
		BlockNumber:    prepareQC.BlockNumber,
		ValidatorIndex: nodeIndex,
		PrepareQC:      prepareQC,
	}
	assert.Nil(t, node.engine.signMsgByBls(v))
	viewChanges[nodeIndex] = v
	viewChangeQC := node.engine.generateViewChangeQC(viewChanges)
	viewChangeQC.QCs = append(append(viewChangeQC.QCs, viewChangeQC.QCs[0]))
	return viewChangeQC
}

// NewBlock returns a bad block for testing.
func NewBadBlock(parent common.Hash, number uint64, node *TestCBFT) *types.Block {
	header := &types.Header{
		Number:      big.NewInt(int64(number)),
		ParentHash:  parent,
		Time:        big.NewInt(time.Now().UnixNano() / 1e6),
		Extra:       make([]byte, 97),
		ReceiptHash: common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b422")),
		Root:        common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Coinbase:    common.Address{},
		GasLimit:    10000000000,
	}

	sign, _ := node.engine.signFn(header.SealHash().Bytes())
	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sign[:])

	block := types.NewBlockWithHeader(header)
	return block
}

func newPrepareBlock(epoch, viewNumber uint64, parentHash common.Hash, blockNumber uint64, blockIndex uint32, nodeIndex uint32, parentQC *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC, secretKeys *bls.SecretKey, badBlock bool, node *TestCBFT, t *testing.T) *protocols.PrepareBlock {
	block := NewBlockWithSign(parentHash, blockNumber, node)
	if badBlock {
		block = NewBadBlock(parentHash, blockNumber, node)
	}
	p := &protocols.PrepareBlock{
		Epoch:         epoch,
		ViewNumber:    viewNumber,
		Block:         block,
		BlockIndex:    blockIndex,
		ProposalIndex: nodeIndex,
		PrepareQC:     parentQC,
		ViewChangeQC:  viewChangeQC,
	}

	// bls sign
	buf, err := p.CannibalizeBytes()
	if err != nil {
		t.Fatalf("%s", "prepareBlock cannibalizeBytes error")
	}
	p.Signature.SetBytes(secretKeys.Sign(string(buf)).Serialize())
	return p
}

func newPrepareVote(epoch, viewNumber uint64, blockHash common.Hash, blockNumber uint64, blockIndex uint32, nodeIndex uint32, parentQC *ctypes.QuorumCert, secretKeys *bls.SecretKey, t *testing.T) *protocols.PrepareVote {
	p := &protocols.PrepareVote{
		Epoch:          epoch,
		ViewNumber:     viewNumber,
		BlockHash:      blockHash,
		BlockNumber:    blockNumber,
		BlockIndex:     blockIndex,
		ValidatorIndex: nodeIndex,
		ParentQC:       parentQC,
	}

	// bls sign
	buf, err := p.CannibalizeBytes()
	if err != nil {
		t.Fatalf("%s", "prepareVote cannibalizeBytes error")
	}
	p.Signature.SetBytes(secretKeys.Sign(string(buf)).Serialize())
	return p
}

func newViewChange(epoch, viewNumber uint64, blockHash common.Hash, blockNumber uint64, blockIndex uint32, nodeIndex uint32, parentQC *ctypes.QuorumCert, secretKeys *bls.SecretKey, t *testing.T) *protocols.ViewChange {
	v := &protocols.ViewChange{
		Epoch:          epoch,
		ViewNumber:     viewNumber,
		BlockHash:      blockHash,
		BlockNumber:    blockNumber,
		ValidatorIndex: nodeIndex,
		PrepareQC:      parentQC,
	}

	// bls sign
	buf, err := v.CannibalizeBytes()
	if err != nil {
		t.Fatalf("%s", "viewChange cannibalizeBytes error")
	}
	v.Signature.SetBytes(secretKeys.Sign(string(buf)).Serialize())
	return v
}

func TestPB01(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// validator fake future index prepare
	proposalIndex := uint32(1)
	fakeHash := common.BytesToHash(utils.Rand32Bytes(32))
	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), fakeHash, qc.BlockNumber+3, qc.BlockIndex+3, proposalIndex, nil, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	_, ok := err.(rules.SafetyError)
	assert.True(t, ok)
	if ok {
		assert.True(t, strings.HasPrefix(err.Error(), fetchPrepare))
	}
}

func TestPB03(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "evidence")
	defer os.RemoveAll(tempDir)

	nodes := MockNodes(t, 2)
	nodes[0].engine.evPool, _ = evidence.NewBaseEvidencePool(tempDir)
	ReachBlock(t, nodes, 5)
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	// seal duplicate prepareBlock
	proposalIndex := uint32(0)
	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), lockBlock.Hash(), lockBlock.NumberU64()+1, lockQC.BlockIndex+1, proposalIndex, nil, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	_, ok := err.(*evidence.DuplicatePrepareBlockEvidence)
	assert.True(t, ok)
	if ok {
		evds := nodes[0].engine.evPool.Evidences()
		assert.Equal(t, 1, evds.Len())
		_, ok = evds[0].(evidence.DuplicatePrepareBlockEvidence)
		if ok {
			assert.Equal(t, lockBlock.NumberU64()+1, evds[0].BlockNumber())
			assert.Equal(t, discover.PubkeyID(&nodes[0].engine.config.Option.NodePriKey.PublicKey), evds[0].NodeID())
			assert.Nil(t, evds[0].Validate())
		}
	}
}

func TestPB04(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	viewChangeQC := ChangeView(t, nodes, qc)

	// base lock block seal first index prepare
	proposalIndex := uint32(1)
	blockIndex := uint32(0)
	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), lockBlock.Hash(), lockBlock.NumberU64()+1, blockIndex, proposalIndex, lockQC, viewChangeQC, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	_, ok := err.(authFailedError)
	assert.True(t, ok)
	if ok {
		assert.True(t, strings.HasPrefix(err.Error(), noBaseMaxBlock))
	}
}

func TestPB05(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	viewChangeQC := ChangeView(t, nodes, qc)
	viewChangeQC.QCs[0].BlockNumber = lockBlock.NumberU64()
	viewChangeQC.QCs[0].BlockHash = lockBlock.Hash()

	// base lock block seal first index prepare
	proposalIndex := uint32(1)
	blockIndex := uint32(0)
	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), lockBlock.Hash(), lockBlock.NumberU64()+1, blockIndex, proposalIndex, lockQC, viewChangeQC, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	_, ok := err.(authFailedError)
	assert.True(t, ok)
	if ok {
		assert.True(t, strings.HasPrefix(err.Error(), errorViewChangeQC))
	}
}

func TestPB06(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	viewChangeQC := ChangeView(t, nodes, qc)

	// base lock block seal first index prepare
	proposalIndex := uint32(1)
	blockIndex := uint32(0)
	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), qcBlock.Hash(), qcBlock.NumberU64()+1, blockIndex, proposalIndex, lockQC, viewChangeQC, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	_, ok := err.(authFailedError)
	assert.True(t, ok)
	if ok {
		assert.True(t, strings.HasPrefix(err.Error(), MismatchedPrepareQC))
	}
}

func TestPB07(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 10)
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	// base lock block seal first index prepare
	proposalIndex := uint32(1)
	blockIndex := uint32(0)
	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), lockBlock.Hash(), lockBlock.NumberU64()+1, blockIndex, proposalIndex, lockQC, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	_, ok := err.(authFailedError)
	assert.True(t, ok)
	if ok {
		assert.True(t, strings.HasPrefix(err.Error(), missingViewChangeQC))
	}
}

func TestPB08(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// base qc block seal first index prepare
	proposalIndex := uint32(1)
	blockIndex := uint32(0)
	viewChangeQC := FakeViewChangeQC(t, nodes[1], nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()+1, proposalIndex, qc)

	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()+1, qcBlock.Hash(), qcBlock.NumberU64()+1, blockIndex, proposalIndex, qc, viewChangeQC, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), dupBlockHash))
}

func TestPB09(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// base qc block seal first index prepare
	proposalIndex := uint32(1)
	blockIndex := uint32(0)

	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()+1, qcBlock.Hash(), qcBlock.NumberU64()+1, blockIndex, proposalIndex, qc, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	err1, ok := err.(rules.SafetyError)
	assert.True(t, ok)
	if ok {
		assert.True(t, err1.Fetch())
	}
}

func TestPB10(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// base qc block seal first index prepare
	proposalIndex := uint32(1)
	blockIndex := uint32(0)

	p := newPrepareBlock(nodes[0].engine.state.Epoch()+2, nodes[0].engine.state.ViewNumber(), qcBlock.Hash(), qcBlock.NumberU64()+1, blockIndex, proposalIndex, qc, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	err := nodes[0].engine.OnPrepareBlock("id", p)

	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), enableVerifyEpoch))
	}
}

func TestPB11(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	viewChangeQC := ChangeView(t, nodes, qc)

	// base qc block seal bad block
	proposalIndex := uint32(1)
	blockIndex := uint32(0)
	p := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), qcBlock.Hash(), qcBlock.NumberU64()+1, blockIndex, proposalIndex, qc, viewChangeQC, nodes[proposalIndex].engine.config.Option.BlsPriKey, true, nodes[0], t)

	result := make(chan interface{}, 1)
	nodes[0].engine.executeStatusHook = func(status *executor.BlockExecuteStatus) {
		assert.NotNil(t, status.Err)
		assert.Equal(t, p.Block.NumberU64(), status.Number)
		assert.Equal(t, p.Block.Hash(), status.Hash)
		result <- struct{}{}
	}

	err := nodes[0].engine.OnPrepareBlock("id", p)
	assert.Nil(t, err)
	if err != nil {
		result <- struct{}{}
	}
	<-result
}

func TestVT01(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// validator fake future index vote
	proposalIndex := uint32(1)
	fakeHash := common.BytesToHash(utils.Rand32Bytes(32))
	p := newPrepareVote(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), fakeHash, qc.BlockNumber+3, qc.BlockIndex+3, proposalIndex, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, t)
	err := nodes[0].engine.OnPrepareVote("id", p)

	err1, ok := err.(rules.SafetyError)
	assert.True(t, ok)
	if ok {
		assert.True(t, err1.FetchPrepare())
	}
}

func TestVT02(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "evidence")
	defer os.RemoveAll(tempDir)

	nodes := MockNodes(t, 2)
	nodes[0].engine.evPool, _ = evidence.NewBaseEvidencePool(tempDir)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	// validator duplicate prepareVote
	validatorIndex := uint32(1)
	fakeHash := common.BytesToHash(utils.Rand32Bytes(32))
	p := newPrepareVote(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), fakeHash, qc.BlockNumber, qc.BlockIndex, validatorIndex, lockQC, nodes[validatorIndex].engine.config.Option.BlsPriKey, t)
	err := nodes[0].engine.OnPrepareVote("id", p)

	_, ok := err.(*evidence.DuplicatePrepareVoteEvidence)
	assert.True(t, ok)
	if ok {
		evds := nodes[0].engine.evPool.Evidences()
		assert.Equal(t, 1, evds.Len())
		_, ok = evds[0].(evidence.DuplicatePrepareVoteEvidence)
		if ok {
			assert.Equal(t, qcBlock.NumberU64()+1, evds[0].BlockNumber())
			assert.Equal(t, discover.PubkeyID(&nodes[0].engine.config.Option.NodePriKey.PublicKey), evds[0].NodeID())
			assert.Nil(t, evds[0].Validate())
		}
	}
}

func TestVT03(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// validator fake next view vote
	proposalIndex := uint32(1)
	fakeHash := common.BytesToHash(utils.Rand32Bytes(32))
	p := newPrepareVote(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()+1, fakeHash, qc.BlockNumber+1, 0, proposalIndex, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, t)
	err := nodes[0].engine.OnPrepareVote("id", p)

	err1, ok := err.(rules.SafetyError)
	assert.True(t, ok)
	if ok {
		assert.True(t, err1.Fetch())
	}
}

func TestVT05(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// proposal seal next two prepare,p1 and p2
	proposalIndex := uint32(0)
	p1 := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), qcBlock.Hash(), qcBlock.NumberU64()+1, qc.BlockIndex+1, proposalIndex, nil, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	nodes[0].engine.OnPrepareBlock("id", p1)
	p2 := newPrepareBlock(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), p1.Block.Hash(), p1.Block.NumberU64()+1, p1.BlockIndex+1, proposalIndex, nil, nil, nodes[proposalIndex].engine.config.Option.BlsPriKey, false, nodes[0], t)
	nodes[0].engine.OnPrepareBlock("id", p2)

	// validator fake p1QC and vote p2
	validaterIndex := uint32(1)
	p := newPrepareVote(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), p2.Block.Hash(), p2.Block.NumberU64(), p2.BlockIndex, validaterIndex, qc, nodes[validaterIndex].engine.config.Option.BlsPriKey, t)
	err := nodes[0].engine.OnPrepareVote("id", p)

	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), MismatchedPrepareQC))
	}
}

func TestVC01(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	// validator base qcBlock and lockQC send viewChange
	validaterIndex := uint32(1)
	v := newViewChange(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), qcBlock.Hash(), qcBlock.NumberU64(), qc.BlockIndex, validaterIndex, lockQC, nodes[validaterIndex].engine.config.Option.BlsPriKey, t)
	err := nodes[0].engine.OnViewChange("id", v)

	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), MismatchedPrepareQC))
	}
}

func TestVC02(t *testing.T) {
	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// fake validator base qcBlock and qc send viewChange
	fakeNode := MockNodes(t, 1)
	fakeNodeIndex := uint32(1)
	v := newViewChange(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), qcBlock.Hash(), qcBlock.NumberU64(), qc.BlockIndex, fakeNodeIndex, qc, fakeNode[0].engine.config.Option.BlsPriKey, t)
	err := nodes[0].engine.OnViewChange("id", v)

	assert.NotNil(t, err)
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), errorSignature))
	}
}

func TestVC03(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "evidence")
	defer os.RemoveAll(tempDir)

	nodes := MockNodes(t, 2)
	nodes[0].engine.evPool, _ = evidence.NewBaseEvidencePool(tempDir)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	// validator duplicate viewChange
	validatorIndex := uint32(1)
	v1 := newViewChange(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), qcBlock.Hash(), qcBlock.NumberU64(), qc.BlockIndex, validatorIndex, qc, nodes[validatorIndex].engine.config.Option.BlsPriKey, t)
	nodes[0].engine.OnViewChange("id", v1)
	v2 := newViewChange(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), lockBlock.Hash(), lockBlock.NumberU64(), lockQC.BlockIndex, validatorIndex, lockQC, nodes[validatorIndex].engine.config.Option.BlsPriKey, t)
	err := nodes[0].engine.OnViewChange("id", v2)

	_, ok := err.(*evidence.DuplicateViewChangeEvidence)
	assert.True(t, ok)
	if ok {
		evds := nodes[0].engine.evPool.Evidences()
		assert.Equal(t, 1, evds.Len())
		_, ok = evds[0].(evidence.DuplicateViewChangeEvidence)
		if ok {
			assert.Equal(t, qcBlock.NumberU64()+1, evds[0].BlockNumber())
			assert.Equal(t, discover.PubkeyID(&nodes[0].engine.config.Option.NodePriKey.PublicKey), evds[0].NodeID())
			assert.Nil(t, evds[0].Validate())
		}
	}
}

func TestVC04(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "evidence")
	defer os.RemoveAll(tempDir)

	nodes := MockNodes(t, 2)
	ReachBlock(t, nodes, 5)
	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	// validator fake next viewChange
	validatorIndex := uint32(1)
	v := newViewChange(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()+1, qcBlock.Hash(), qcBlock.NumberU64(), qc.BlockIndex, validatorIndex, qc, nodes[validatorIndex].engine.config.Option.BlsPriKey, t)
	err := nodes[0].engine.OnViewChange("id", v)

	err1, ok := err.(rules.SafetyError)
	assert.True(t, ok)
	if ok {
		assert.True(t, err1.Fetch())
	}
}
