// Copyright 2021 The PlatON Network Authors
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
	"math/big"
	"regexp"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type PrepareBlockTestSuite struct {
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
}

func SetupPrepareBlockTestTest(period uint64) *PrepareBlockTestSuite {
	suit := new(PrepareBlockTestSuite)
	suit.view = newTestView(false, period)
	suit.blockOne = NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstProposer().state.ViewNumber()
	return suit
}
func (suit *PrepareBlockTestSuite) insertOneBlock() {
	for _, cbft := range suit.view.allCbft {
		// cbft.insertQCBlock(suit.blockOne, suit.blockOneQC.BlockQC)
		insertBlock(cbft, suit.blockOne, suit.blockOneQC.BlockQC)
	}
}

func (suit *PrepareBlockTestSuite) createEvPool(paths []string) {
	if len(paths) != len(suit.view.allNode) {
		panic("paths len err")
	}
	for i, path := range paths {
		pool, _ := evidence.NewBaseEvidencePool(path)
		suit.view.allCbft[i].evPool = pool
	}

}

func (suit *PrepareBlockTestSuite) waitVote() {
	time.Sleep(time.Millisecond * 500)
}

// prepareBlock message basic check
// 1.Missing signed prepareBlock message
// 2.a prepareBlock message whose signature is inconsistent with the proposer index
// 3.The block node is not the prepareBlock message of the current proposed node.
// 4.The proposer index matches the prepareBlock message that does not match the proposer
// 5.The prepareBlock message of the proposer non-consensus node
// 6.epoch too big
// 7.epoch too small
func TestCheckErrPrepareBlock(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	notConsensusNodes := mockNotConsensusNode(suit.view.nodeParams, 1, testPeriod)
	// suit.view.secondProposer().state.ResetView(suit.view.Epoch(), suit.oldViewNumber+1)
	testcases := []struct {
		name string
		data *protocols.PrepareBlock
		err  error
	}{
		{
			name: "Missing signed prepareBlock message",
			data: mockPrepareBlock(nil, suit.view.Epoch(), suit.oldViewNumber,
				0, suit.view.firstProposerIndex(), suit.blockOne, nil, nil),
		},
		{
			name: "The proposer index matches the prepareBlock message that does not match the proposer",
			data: mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber, 0, suit.view.firstProposerIndex()+1,
				suit.blockOne, nil, nil),
		},
		{
			name: "The suit.blockOne node is not the prepareBlock message of the current proposed node.",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber, 0,
				suit.view.secondProposerIndex(), suit.blockOne, nil, nil),
		},
		{
			name: "a prepareBlock message whose signature is inconsistent with the proposer index",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber, 0,
				suit.view.firstProposerIndex(), suit.blockOne, nil, nil),
		},
		{
			name: "The prepareBlock message of the proposer non-consensus node",
			data: mockPrepareBlock(notConsensusNodes[0].engine.config.Option.BlsPriKey,
				suit.view.Epoch(), suit.view.firstProposer().state.ViewNumber(),
				0, suit.view.firstProposerIndex(), suit.blockOne, nil, nil),
		},
		{
			name: "epoch big",
			data: mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch()+1,
				suit.oldViewNumber, 0, suit.view.firstProposerIndex(),
				suit.blockOne, nil, nil),
		},
		{
			name: "epoch small",
			data: mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch()-1,
				suit.oldViewNumber, 0, suit.view.firstProposerIndex(),
				suit.blockOne, nil, nil),
		},
	}
	for _, testcase := range testcases {
		if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().Node().ID().String(), testcase.data); err == nil {
			t.Errorf("case %s is failed", testcase.name)
			suit.view.secondProposer().state.ResetView(suit.view.Epoch(), suit.oldViewNumber)
			// t.Error(err.Error())
		}
	}
}

// Carry the prepareBlock message of prepareQC and viewChangeQC. The node does not complete viewChangeQC.
// The verification passes, the number of votes with blockIndex 0 is 1, viewNumber+1
func TestPrepareBlockOneWithViewChangePrepareQCAndViewChangeQC(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.insertOneBlock()
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	viewQC := mockViewQC(suit.blockOne, suit.view.allNode, suit.blockOneQC.BlockQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err != nil {
		t.Fatal(err.Error())
	}
	suit.waitVote()
	if 1 != suit.view.firstProposer().state.PrepareVoteLenByIndex(0) {
		t.Fatal("fail")
	}
	if suit.oldViewNumber+1 != suit.view.firstProposer().state.ViewNumber() {
		t.Fatal("fail,suit.oldViewNumber", suit.oldViewNumber+1)
	}
}

// Carry the prepareBlock message of prepareQC and viewChangeQC, and the node receiving the message has received the viewChangeQC.
// The verification passes, the number of votes with blockIndex 0 is 1, and the viewNumber is unchanged.
func TestPrepareBlockOneWithViewChangePrepareQCAndViewChangeQCHadViewChangQC(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.insertOneBlock()
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	viewQC := mockViewQC(suit.blockOne, suit.view.allNode, suit.blockOneQC.BlockQC)
	fmt.Println(viewQC.String())
	suit.view.firstProposer().changeView(suit.view.Epoch(), suit.oldViewNumber+1, suit.blockOne, suit.blockOneQC.BlockQC, viewQC)
	suit.view.firstProposer().state.ResetView(suit.view.Epoch(), suit.oldViewNumber+1)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	fmt.Println(prepareBlock.BlockNum())
	fmt.Println(prepareBlock.Block.ParentHash().String())
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err != nil {
		t.Fatal(err.Error())
	}
	suit.waitVote()
	if 1 != suit.view.firstProposer().state.PrepareVoteLenByIndex(0) {
		t.Fatal("fail")
	}
	if suit.oldViewNumber+1 != suit.view.firstProposer().state.ViewNumber() {
		t.Fatal("fail,suit.oldViewNumber", suit.oldViewNumber+1)
	}
}

// The first block of viewChangeQC, the node receiving the message did not receive viewChangeQC
// 1.Does not carry the prepareBlock message of viewChangeQC
// 2.Does not carry the prepareBlock message of prepareQC
// 3.The prepareBlock message does not carry prepareQC and viewChangeQC
// 4.prepareBlock message with blockIndex 1
// 5.Carry prepareQC and carry viewChangeQC that does not satisfy 2f+1
// 6.epoch too big
// 7.epoch too small
func TestPrepareBlockOneWithViewChangeErrFirstBlock(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.insertOneBlock()
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	viewQC := mockViewQC(suit.blockOne, suit.view.allNode, suit.blockOneQC.BlockQC)
	errViewQC := mockViewQC(suit.blockOne, suit.view.allNode[0:1], suit.blockOneQC.BlockQC)
	oldEpoch := suit.view.Epoch()
	testcases := []struct {
		name string
		data *protocols.PrepareBlock
		err  error
	}{
		{
			name: "Does not carry the prepareBlock message of viewChangeQC",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 0,
				suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, nil),
		},
		{
			name: "Does not carry the prepareBlock message of prepareQC",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 0,
				suit.view.secondProposerIndex(), block2, nil, viewQC),
		},
		{
			name: "The prepareBlock message does not carry prepareQC and viewChangeQC",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 0,
				suit.view.secondProposerIndex(), block2, nil, nil),
		},
		{
			name: "prepareBlock message with blockIndex 1",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 1,
				suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC),
		},
		{
			name: "Carry prepareQC and carry viewChangeQC that does not satisfy 2f+1",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 0,
				suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, errViewQC),
		},
		{
			name: "epoch big",
			data: mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch()+1,
				suit.oldViewNumber, 0, suit.view.firstProposerIndex(),
				suit.blockOne, nil, nil),
		},
		{
			name: "epoch small",
			data: mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch()-1,
				suit.oldViewNumber, 0, suit.view.firstProposerIndex(),
				suit.blockOne, nil, nil),
		},
	}
	for _, testcase := range testcases {
		if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), testcase.data); err == nil {
			t.Errorf("CASE:%s is failed", testcase.name)
			suit.view.firstProposer().state.ResetView(oldEpoch, suit.oldViewNumber+1)
			// t.Error(err.Error())
		} else {
			fmt.Printf("case:%s-err:%s\n", testcase.name, err.Error())
		}
	}
}

// The first block of viewChangeQC, the node receiving the message has received viewChangeQC
// 1.Does not carry the prepareBlock message of viewChangeQC
// 2.Does not carry the prepareBlock message of prepareQC
// 3.The prepareBlock message does not carry prepareQC and viewChangeQC
// 4.prepareBlock message with blockIndex 1
// 5.Carry prepareQC and carry viewChangeQC that does not satisfy 2f+1
func TestPrepareBlockOneWithViewChangeErrFirstBlockHadViewChangQC(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.insertOneBlock()
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	viewQC := mockViewQC(suit.blockOne, suit.view.allNode, suit.blockOneQC.BlockQC)
	errViewQC := mockViewQC(suit.blockOne, suit.view.allNode[0:1], suit.blockOneQC.BlockQC)
	oldEpoch := suit.view.Epoch()
	suit.view.firstProposer().state.ResetView(suit.view.Epoch(), suit.oldViewNumber+1)
	testcases := []struct {
		name string
		data *protocols.PrepareBlock
		err  error
	}{
		{
			name: "Does not carry the prepareBlock message of viewChangeQC",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 0,
				suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, nil),
		},
		{
			name: "Does not carry the prepareBlock message of prepareQC",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 0,
				suit.view.secondProposerIndex(), block2, nil, viewQC),
		},
		{
			name: "The prepareBlock message does not carry prepareQC and viewChangeQC",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 0,
				suit.view.secondProposerIndex(), block2, nil, nil),
		},
		{
			name: "prepareBlock message with blockIndex 1",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 1,
				suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC),
		},
		{
			name: "Carry prepareQC and carry viewChangeQC that does not satisfy 2f+1",
			data: mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
				suit.oldViewNumber+1, 0,
				suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, errViewQC),
		},
	}
	for _, testcase := range testcases {
		if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), testcase.data); err == nil {
			t.Errorf("CASE:%s is failed", testcase.name)
			suit.view.firstProposer().state.ResetView(oldEpoch, suit.oldViewNumber+1)
			// t.Error(err.Error())
		} else {
			fmt.Printf("case:%s-err:%s\n", testcase.name, err.Error())
		}
	}
}

// The first block of viewChangeQC, the block node HighestQCBlock leads the prepareBlock message of the local HighestQCBlock, and the node receiving the block does not complete the viewChangeQC
// Due to backwardness, synchronization will be triggered, verification will not pass, and the error value returned is viewNumber higher then local(local:0, msg:1)
func TestPrepareBlockOneWithViewChangeFirstBlockTooHigh(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	block3 := NewBlockWithSign(block2.Hash(), 3, suit.view.allNode[0])
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	suit.view.secondProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewQC := mockViewQC(block2, suit.view.allNode, block2QC.BlockQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
		suit.oldViewNumber+1, 0, suit.view.secondProposerIndex(), block3, block2QC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		if "viewNumber higher than local(local:0, msg:1)" != err.Error() {
			t.Fatal("FAIL")
		}
	}
}

// The first block of viewChangeQC, the block node HighestQCBlock is behind the local HighestQCBlock prepareBlock message, and the node receiving the block is not completed viewChangeQC
// The verification failed, and the error value returned is viewNumber higher then local(local:0, msg:1)
func TestPrepareBlockOneWithViewChangeFirstBlockTooLow(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	suit.view.firstProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewQC := mockViewQC(block2, suit.view.allNode, block2QC.BlockQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		if "viewNumber higher than local(local:0, msg:1)" != err.Error() {
			t.Fatal("FAIL")
		}
	}
}

// The first block of viewChangeQC, the block node HighestQCBlock is behind the local HighestQCBlock prepareBlock message, the node receiving the block has completed viewChangeQC, and this time is not based on viewQC.MaxBlock
// Verification cannot pass
func TestPrepareBlockOneWithViewChangeFirstBlockTooLowHad(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	suit.view.firstProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewQC := mockViewQC(block2, suit.view.allNode, block2QC.BlockQC)
	suit.view.firstProposer().state.ResetView(suit.view.Epoch(), suit.oldViewNumber+1)
	suit.view.firstProposer().state.SetLastViewChangeQC(viewQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	}
}

// The first block of viewChangeQC, not based on viewQC.MaxBlock, the node that received the block did not complete viewChangeQC
// Verification cannot pass
func TestPrepareBlockOneWithViewChangeFirstBlockNotWithMaxBlock(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	// suit.view.firstProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewQC := mockViewQC(block2, suit.view.allNode, block2QC.BlockQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The first block confirmed by the last block, the hash with the same blockNumber is different from the prepareBlock message.
// The first block passes, the second check fails, and the corresponding PrepareVoteLen=1 returns double evidence error.
func TestPrepareBlockOneWithDifHash(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	paths := createPaths(len(suit.view.allCbft), t)
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block1 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	header := &types.Header{
		Number:      big.NewInt(int64(11)),
		ParentHash:  suit.view.firstProposer().state.HighestQCBlock().Hash(),
		Time:        uint64(time.Now().UnixNano() + 100),
		Extra:       make([]byte, 97),
		ReceiptHash: common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Root:        common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Coinbase:    common.Address{},
		GasLimit:    100000000001,
	}
	sign, _ := suit.view.allNode[1].engine.signFn(header.SealHash().Bytes())
	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sign[:])
	block2 := types.NewBlockWithHeader(header)
	fmt.Println(common.Bytes2Hex(block2.Extra()))
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock1 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block1, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock1); err != nil {
		t.Fatal(err.Error())
	}
	suit.waitVote()
	prepareBlock2 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, qc, nil)
	fmt.Println(block1.Hash().String(), block2.Hash().String())
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock2); err == nil {
		t.Error("FAIL")
	} else {
		reg := regexp.MustCompile(`DuplicatePrepareBlockEvidence`)
		if len(reg.FindAllString(err.Error(), -1)) == 0 {
			t.Fatal(err.Error())
		}
	}
	suit.waitVote()
	if 1 != suit.view.firstProposer().state.PrepareVoteLenByIndex(0) {
		t.Fatal("fail")
	}
}

// The first block confirmed by the last block, carrying the prepareBlock message of prepareQC
// Verification pass，PrepareVoteLenByIndex(0)=1
func TestPrepareBlockOneWithLastBlockQC(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[0])
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(), suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err != nil {
		t.Fatal(err.Error())
	}
	suit.waitVote()
	if 1 != suit.view.firstProposer().state.PrepareVoteLenByIndex(0) {
		t.Fatal("fail")
	}
}

// The first block, carrying the wrong qc
func TestPrepareBlockOneWithNumberIsOne(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	block1 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 1, suit.view.allNode[0])
	notConsensusNodes := mockNotConsensusNode(suit.view.nodeParams, 4, testPeriod)
	errQC := mockErrBlockQC(notConsensusNodes, block1, 0, nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block1, errQC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("fail")
	}
}

// Non-first block, carrying the wrong qc
func TestPrepareBlockWithBlockIndexNotIsZero(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.insertOneBlock()
	block1 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 2, suit.view.allNode[0])
	notConsensusNodes := mockNotConsensusNode(suit.view.nodeParams, 4, testPeriod)
	errQC := mockErrBlockQC(notConsensusNodes, block1, 0, nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block1, errQC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("fail")
	}
}

// The first block confirmed by the last block, Does not carry the prepareBlock message of prepareQC
// Verification failed
func TestPrepareBlockOneWithLastBlockQCNotQC(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The first block confirmed by the last block, prepareBlock message with blockIndex 1
// Verification failed
func TestPrepareBlockOneWithLastBlockQCBlockIndexIsOne(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	// oldEpoch := suit.view.Epoch()
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block11, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The first block confirmed by the last block, the block node HighestQCBlock leads the prepareBlock message of the local HighestQCBlock
// Verification failed，Trigger synchronization
func TestPrepareBlockOneWithLastBlockQCLead(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	otherNode := suit.view.thirdProposer()
	suit.view.setBlockQC(9, suit.view.allNode[0])
	block10 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 10, suit.view.allNode[0])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block10QC := mockBlockQC(suit.view.allNode, block10, 0,
		oldQC)
	block11 := NewBlockWithSign(block10.Hash(), 11, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, block10QC.BlockQC, nil)
	if err := otherNode.OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// viewNumber is less than the current viewNumber
// Verification failed
func TestPrepareBlockNotOneWithViewNumberTooLow(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	// oldEpoch := suit.view.Epoch()
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(block11.Hash(), 12, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// viewNumber is greater than the current viewNumber
// Verification failed，Trigger synchronization
func TestPrepareBlockNotOneWithViewNumberTooHigh(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(suit.blockOne.Hash(), 12, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+2, 1,
		suit.view.secondProposerIndex(), block12, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// PreviousBlock of the previous block prepareQC
// Verification pass, number of votes +1
func TestPrepareBlockNotOneWithParentQC(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(block11.Hash(), 12, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err != nil {
		t.Fatal(err.Error())
	}
	suit.waitVote()
	if 1 != suit.view.firstProposer().state.PrepareVoteLenByIndex(1) {
		t.Fatal("fail")
	}
}

// Previous block without prepareQC prepareBlock
// Verification pass，There are votes for the block in pengding
func TestPrepareBlockNotOneWithParentNotQC(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock11 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, oldQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock11); err != nil {
		t.Fatal("FAIL")
	}
	suit.waitVote()
	block12 := NewBlockWithSign(block11.Hash(), 12, suit.view.allNode[1])
	prepareBlock12 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock12); err != nil {
		t.Fatal(err.Error())
	}
	suit.waitVote()
	if 12 != suit.view.firstProposer().state.PendingPrepareVote().Votes[0].BlockNum() {
		t.Fatal("fail")
	}

}

// The number of blocks exceeds the limit of one round
// Verification failed
func TestPrepareBlockNotOneWithAmountTooMany(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 10,
		suit.view.secondProposerIndex(), block11, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// Same block height, different hash of prepareBlock
// The second Verification failed, double evidence error
func TestPrepareBlockNotOneWithBlockNumberRepeat(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	paths := createPaths(len(suit.view.allCbft), t)
	defer removePaths(paths)
	suit.createEvPool(paths)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().Node().ID().String(), prepareBlock1); err != nil {
		t.Fatal(err.Error())
	}
	block2 := NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	prepareBlock2 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().Node().ID().String(), prepareBlock2); err == nil {
		t.Fatal("FAIL")
	} else {
		reg := regexp.MustCompile(`DuplicatePrepareBlockEvidence`)
		if len(reg.FindAllString(err.Error(), -1)) == 0 {
			t.Fatal(err.Error())
		}
	}
}

// Block high discontinuous prepareBlock, block hash continuous
// Verification failed
func TestPrepareBlockNotOneWithBlockNumberDiscontinuous(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(block11.Hash(), 13, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// Block hash discontinuous block high continuous prepareBlock
// Verification failed
func TestPrepareBlockNotOneWithBlockHashDiscontinuous(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	fmt.Println(suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(suit.view.genesisBlock.Hash(), 12, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The prepareBlock whose BlockIndex does not match the actual block index
// Verification failed
func TestPrepareBlockNotOneWithBlockIndexErr(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(block11.Hash(), 12, suit.view.allNode[1])
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 2,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}

}

// The same BlockIndex block exists in this local block, but BlockHash, BlockNumber are not equal
// Verification failed
func TestPrepareBlockNotOneWithBlockIndexRepeat(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock1); err != nil {
		t.Fatal(err.Error())
	}
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 3, suit.view.allNode[0])
	prepareBlock2 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock2); err == nil {
		t.Fatal("FAIL")
	}
}

// Data correctly repeating prepareBlock message
// Verification failed
func TestPrepareBlockDup(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock1); err != nil {
		t.Fatal(err.Error())
	}
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock1); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// There is no prepareBlock message of the previous index block corresponding to the BlockIndex.
// Verification failed
func TestPrepareBlockBlockIndexTooHigh(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock1); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The correct timeout of the data is the prepareBlock message
// Verification failed
func TestPrepareBlockNotOneWithTimeout(t *testing.T) {
	suit := SetupPrepareBlockTestTest(3000)
	time.Sleep((3000 + 200) * time.Millisecond)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock1); err == nil {
		t.Fatal("FAIL")
	} else {
		t.Log(err)
	}
}

// Non-consensus node receives a valid prepareBlock message
// Verification pass
func TestPrepareBlockNotOneWithNotConsensus(t *testing.T) {
	suit := SetupPrepareBlockTestTest(10000)
	notConsensus := mockNotConsensusNode(suit.view.nodeParams, 1, testPeriod)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := notConsensus[0].engine.OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock1); err != nil {
		t.Fatal(err.Error())
	}
}
