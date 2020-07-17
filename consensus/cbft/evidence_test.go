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
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/suite"
)

func TestEvidenceSuite(t *testing.T) {
	suite.Run(t, new(EvidenceTestSuite))
}

type EvidenceTestSuite struct {
	suite.Suite
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
	epoch         uint64
}

func (suit *EvidenceTestSuite) insertOneBlock() {
	for _, cbft := range suit.view.allCbft {
		insertBlock(cbft, suit.blockOne, suit.blockOneQC.BlockQC)
	}
}

func (suit *EvidenceTestSuite) createEvPool(paths []string) {
	if len(paths) != len(suit.view.allNode) {
		panic("paths len err")
	}
	for i, path := range paths {
		pool, _ := evidence.NewBaseEvidencePool(path)
		suit.view.allCbft[i].evPool = pool
	}

}

func (suit *EvidenceTestSuite) SetupTest() {
	suit.view = newTestView(false, testNodeNumber)
	suit.blockOne = NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstCbft.state.ViewNumber()
	suit.epoch = suit.view.firstCbft.state.Epoch()
}

// Double view change
func (suit *EvidenceTestSuite) TestViewChangeDuplicate() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.insertOneBlock()
	viewChange1 := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber,
		suit.blockOne.Hash(), suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC)
	viewChange2 := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber,
		suit.view.genesisBlock.Hash(), suit.view.genesisBlock.NumberU64(), suit.view.secondProposerIndex(), nil)
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange1); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange2); err == nil {
		suit.T().Fatal("fail")
	}
	ev := suit.view.firstProposer().Evidences()
	var duplicateViewChange *evidence.EvidenceData
	if err := json.Unmarshal([]byte(ev), &duplicateViewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	if len(duplicateViewChange.DC) == 0 {
		suit.T().Fatal("Evidences is empty")
	}
	suit.EqualValues(suit.view.secondProposer().config.Option.NodeID, duplicateViewChange.DC[0].ViewA.ValidateNode.NodeID)
	suit.EqualValues(suit.view.secondProposer().config.Option.BlsPriKey.GetPublicKey().Serialize(), duplicateViewChange.DC[0].ViewA.ValidateNode.BlsPubKey.Serialize())
	suit.EqualValues(suit.view.secondProposer().config.Option.NodeID, duplicateViewChange.DC[0].ViewB.ValidateNode.NodeID)
	suit.EqualValues(suit.view.secondProposer().config.Option.BlsPriKey.GetPublicKey().Serialize(), duplicateViewChange.DC[0].ViewB.ValidateNode.BlsPubKey.Serialize())
}

// viewChange view number dif
func (suit *EvidenceTestSuite) TestViewChangeDuplicateDifViewNumber() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.insertOneBlock()
	viewChange1 := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber+1,
		suit.blockOne.Hash(), suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC)
	viewChange2 := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber,
		suit.view.genesisBlock.Hash(), suit.view.genesisBlock.NumberU64(), suit.view.secondProposerIndex(), nil)
	vnode, _ := suit.view.firstProposer().validatorPool.GetValidatorByIndex(suit.epoch, suit.view.secondProposerIndex())
	if err := suit.view.firstProposer().evPool.AddViewChange(viewChange1, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().evPool.AddViewChange(viewChange2, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	ev := suit.view.firstProposer().Evidences()
	suit.Equal("{}", ev)
}

// viewChange dif epoch
func (suit *EvidenceTestSuite) TestViewChangeDuplicateDifEpoch() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.insertOneBlock()
	viewChange1 := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch+1, suit.oldViewNumber,
		suit.blockOne.Hash(), suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC)
	viewChange2 := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber,
		suit.view.genesisBlock.Hash(), suit.view.genesisBlock.NumberU64(), suit.view.secondProposerIndex(), nil)
	vnode, _ := suit.view.firstProposer().validatorPool.GetValidatorByIndex(suit.epoch, suit.view.secondProposerIndex())
	if err := suit.view.firstProposer().evPool.AddViewChange(viewChange1, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().evPool.AddViewChange(viewChange2, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	ev := suit.view.firstProposer().Evidences()
	suit.Equal("{}", ev)
}

// duplicateEvidence
func (suit *EvidenceTestSuite) TestPrepareBlockDuplicate() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block1 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	header := &types.Header{
		Number:      big.NewInt(int64(11)),
		ParentHash:  suit.view.firstProposer().state.HighestQCBlock().Hash(),
		Time:        big.NewInt(time.Now().UnixNano() + 100),
		Extra:       make([]byte, 97),
		ReceiptHash: common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Root:        common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Coinbase:    common.Address{},
		GasLimit:    100000000001,
	}
	sign, _ := suit.view.allNode[1].engine.signFn(header.SealHash().Bytes())
	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sign[:])
	block2 := types.NewBlockWithHeader(header)
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock1 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.epoch,
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block1, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal(err.Error())
	}
	// time.Sleep(time.Millisecond * 10)
	prepareBlock2 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock2); err == nil {
		suit.T().Error("FAIL")
	} else {
		fmt.Println(err.Error())
	}
	ev := suit.view.firstProposer().Evidences()
	var duplicatePrepareBlock *evidence.EvidenceData
	if err := json.Unmarshal([]byte(ev), &duplicatePrepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	if len(duplicatePrepareBlock.DP) == 0 {
		suit.T().Fatal("Evidences is empty")
	}
	suit.EqualValues(suit.view.secondProposer().config.Option.NodeID, duplicatePrepareBlock.DP[0].PrepareB.ValidateNode.NodeID)
	suit.EqualValues(suit.view.secondProposer().config.Option.BlsPriKey.GetPublicKey().Serialize(), duplicatePrepareBlock.DP[0].PrepareB.ValidateNode.BlsPubKey.Serialize())
	suit.EqualValues(suit.view.secondProposer().config.Option.NodeID, duplicatePrepareBlock.DP[0].PrepareA.ValidateNode.NodeID)
	suit.EqualValues(suit.view.secondProposer().config.Option.BlsPriKey.GetPublicKey().Serialize(), duplicatePrepareBlock.DP[0].PrepareA.ValidateNode.BlsPubKey.Serialize())
}

// prepare block view number dif
func (suit *EvidenceTestSuite) TestPrepareBlockDuplicateDifViewNumber() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block1 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	header := &types.Header{
		Number:      big.NewInt(int64(11)),
		ParentHash:  suit.view.firstProposer().state.HighestQCBlock().Hash(),
		Time:        big.NewInt(time.Now().UnixNano() + 100),
		Extra:       make([]byte, 97),
		ReceiptHash: common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Root:        common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Coinbase:    common.Address{},
		GasLimit:    100000000001,
	}
	sign, _ := suit.view.allNode[1].engine.signFn(header.SealHash().Bytes())
	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sign[:])
	block2 := types.NewBlockWithHeader(header)
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock1 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.epoch,
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block1, qc, nil)
	prepareBlock2 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber+2, 0,
		suit.view.secondProposerIndex(), block2, qc, nil)
	vnode, _ := suit.view.firstProposer().validatorPool.GetValidatorByIndex(suit.epoch, suit.view.secondProposerIndex())
	if err := suit.view.firstProposer().evPool.AddPrepareBlock(prepareBlock1, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().evPool.AddPrepareBlock(prepareBlock2, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	ev := suit.view.firstProposer().Evidences()
	suit.Equal("{}", ev)
}

// prepare block epoch dif
func (suit *EvidenceTestSuite) TestPrepareBlockDuplicateDifEpoch() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block1 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	header := &types.Header{
		Number:      big.NewInt(int64(11)),
		ParentHash:  suit.view.firstProposer().state.HighestQCBlock().Hash(),
		Time:        big.NewInt(time.Now().UnixNano() + 100),
		Extra:       make([]byte, 97),
		ReceiptHash: common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Root:        common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Coinbase:    common.Address{},
		GasLimit:    100000000001,
	}
	sign, _ := suit.view.allNode[1].engine.signFn(header.SealHash().Bytes())
	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sign[:])
	block2 := types.NewBlockWithHeader(header)
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock1 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.epoch+1,
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block1, qc, nil)
	prepareBlock2 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, qc, nil)
	vnode, _ := suit.view.firstProposer().validatorPool.GetValidatorByIndex(suit.epoch, suit.view.secondProposerIndex())
	if err := suit.view.firstProposer().evPool.AddPrepareBlock(prepareBlock1, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().evPool.AddPrepareBlock(prepareBlock2, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	ev := suit.view.firstProposer().Evidences()
	suit.Equal("{}", ev)
}

// duplicate sign
func (suit *EvidenceTestSuite) TestPrepareVoteDuplicate() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	block1 := NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	prepareVote1 := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	prepareVote2 := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block1.Hash(),
		block1.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote1); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote2); err == nil {
		suit.T().Fatal("FAIL")
	}
	ev := suit.view.secondProposer().Evidences()
	var duplicatePrepareVote *evidence.EvidenceData
	if err := json.Unmarshal([]byte(ev), &duplicatePrepareVote); err != nil {
		suit.T().Fatal(err.Error())
	}
	if len(duplicatePrepareVote.DV) == 0 {
		suit.T().Fatal("Evidences is empty")
	}
	suit.EqualValues(suit.view.firstProposer().config.Option.NodeID, duplicatePrepareVote.DV[0].VoteB.ValidateNode.NodeID)
	suit.EqualValues(suit.view.firstProposer().config.Option.BlsPriKey.GetPublicKey().Serialize(), duplicatePrepareVote.DV[0].VoteB.ValidateNode.BlsPubKey.Serialize())
	suit.EqualValues(suit.view.firstProposer().config.Option.NodeID, duplicatePrepareVote.DV[0].VoteB.ValidateNode.NodeID)
	suit.EqualValues(suit.view.firstProposer().config.Option.BlsPriKey.GetPublicKey().Serialize(), duplicatePrepareVote.DV[0].VoteB.ValidateNode.BlsPubKey.Serialize())
}

// prepare vote view number dif
func (suit *EvidenceTestSuite) TestPrepareVoteDuplicateDifViewNumber() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	block1 := NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	prepareVote1 := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 0,
		suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	prepareVote2 := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block1.Hash(),
		block1.NumberU64(), nil)
	vnode, _ := suit.view.firstProposer().validatorPool.GetValidatorByIndex(suit.epoch, suit.view.secondProposerIndex())
	if err := suit.view.firstProposer().evPool.AddPrepareVote(prepareVote1, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().evPool.AddPrepareVote(prepareVote2, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	ev := suit.view.firstProposer().Evidences()
	suit.Equal("{}", ev)
}

// prepare vote epoch dif
func (suit *EvidenceTestSuite) TestPrepareVoteDuplicateDifEpoch() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	block1 := NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	prepareVote1 := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch+1, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	prepareVote2 := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block1.Hash(),
		block1.NumberU64(), nil)
	vnode, _ := suit.view.firstProposer().validatorPool.GetValidatorByIndex(suit.epoch, suit.view.secondProposerIndex())
	if err := suit.view.firstProposer().evPool.AddPrepareVote(prepareVote1, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().evPool.AddPrepareVote(prepareVote2, vnode); err != nil {
		suit.T().Fatal(err.Error())
	}
	ev := suit.view.firstProposer().Evidences()
	suit.Equal("{}", ev)
}
