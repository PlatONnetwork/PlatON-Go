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
	"regexp"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type PrepareVoteTestSuite struct {
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
	epoch         uint64
}

func SetupPrepareVoteTest(period uint64) *PrepareVoteTestSuite {
	suit := new(PrepareVoteTestSuite)
	suit.view = newTestView(false, period)
	suit.blockOne = NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstProposer().state.ViewNumber()
	suit.epoch = suit.view.Epoch()
	return suit
}
func (suit *PrepareVoteTestSuite) insertOneBlock() {
	for _, cbft := range suit.view.allCbft {
		insertBlock(cbft, suit.blockOne, suit.blockOneQC.BlockQC)
	}
}
func (suit *PrepareVoteTestSuite) createEvPool(paths []string) {
	if len(paths) != len(suit.view.allNode) {
		panic("paths len err")
	}
	for i, path := range paths {
		pool, _ := evidence.NewBaseEvidencePool(path)
		suit.view.allCbft[i].evPool = pool
	}
}

func (suit *PrepareVoteTestSuite) waitVote() {
	time.Sleep(time.Millisecond * 500)
}

// Construct prepareVote message
// Receive the block and generate the corresponding vote
// Check block height is consistent with block hash
func TestBuildPrepareVote(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
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
	vote := suit.view.firstProposer().state.AllPrepareVoteByIndex(1)[0]
	if vote.BlockNum() != uint64(12) {
		t.Fatal("vote blocknum not compare")
	}
	if block12.Hash() != vote.BlockHash {
		t.Fatal("block12 not compare vote hash")
	}
	if prepareBlock.BlockIndex != vote.BlockIndex {
		t.Fatal("prepareBlock BlockIndex not compare vote.BlockIndex ")
	}
}

// prepareVote message basic check
// 1.Unsigned
// 2.The signature is inconsistent with the validatorIndex
// 3.The signature is not the verification node
// 4.epoch too big
// 5.epoch too small
func TestCheckErrPrepareVote(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	_, notConsensusKey := GenerateKeys(1)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().Node().ID().String(), prepareBlock); err != nil {
		t.Fatal("FAIL")
	}
	testcases := []struct {
		name string
		data *protocols.PrepareVote
		err  error
	}{
		{
			name: "Missing signature prepareVote message",
			data: mockPrepareVote(nil, suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil),
		},
		{
			name: "a prepareVote message whose signature is inconsistent with the certifier index",
			data: mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.secondProposerIndex(), suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil),
		},
		{
			name: "Authenticator non-consensus node prepareVote message",
			data: mockPrepareVote(notConsensusKey[0], suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil),
		},
		{
			name: "epoch big",
			data: mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch+1, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil),
		},
		{
			name: "epoch small",
			data: mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch-1, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil),
		},
	}
	for _, testcase := range testcases {
		if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), testcase.data); err == nil {
			t.Errorf("case %s is failed", testcase.name)
		} else {
			fmt.Println(err.Error())
		}
	}
}

// ParentVote message that does not carry ParentQC with zero parent block
// Verification pass
func TestPrepareVoteWithParentIsZeroButNotParentQC(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	epoch := suit.view.Epoch()
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	// if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), prepareBlock); err != nil {
	// 	t.Fatal(err.Error())
	// }
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(),
		suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err != nil {
		t.Fatal(err.Error())
	}
	if suit.view.firstProposerIndex() != suit.view.secondProposer().state.AllPrepareVoteByIndex(0)[0].ValidatorIndex {
		t.Error("TestPrepareVoteWithParentIsZeroButNotParentQC fail")
	}
}

// The parent block is non-zero and does not carry the parentVC prepareVote message.
// Verification failed
func TestPrepareVoteWithParentIsNotZeroButNotParentQC(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	suit.insertOneBlock()
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err != nil {
		t.Fatal(err.Error())
	}
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 1, suit.view.firstProposerIndex(), block2.Hash(),
		block2.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The parent block is non-zero but the blockIndex is 0. The parentVote message does not carry ParentQC.
// Verification failed
func TestPrepareVoteWithParentIsNotZeroAndBlockIndexNotParentQC(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	n, h := suit.view.firstProposer().HighestQCBlockBn()
	block1 := NewBlockWithSign(h, n+1, suit.view.allNode[1])
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(h, n)
	fmt.Println(qc.String())
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block1, qc, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock); err != nil {
		t.Fatal(err.Error())
	}
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 0, suit.view.firstProposerIndex(), block1.Hash(),
		block1.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// blockNumber=1, qc forged message
// Verification failed
func TestPrepareVoteWithBlockNumberIsOneAndErrParentQC(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	block1 := NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	notConsensusNodes := mockNotConsensusNode(suit.view.nodeParams, 4, testPeriod)
	errQC := mockErrBlockQC(notConsensusNodes, block1, 0, nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), block1, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock("", prepareBlock); err != nil {
		t.Fatal(err.Error())
	}
	suit.waitVote()
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), block1.Hash(), 1, errQC.BlockQC)
	if err := suit.view.secondProposer().OnPrepareVote("", prepareVote); err == nil {
		t.Fatal("fail")
	}
}

// Received the prepareVote message received by prepareBlock first
// Verification failed
func TestPrepareVoteWithNotPrepareBlock(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The prepareVote message that the block exceeds the limit of one round of the block number
// Verification failed
func TestPrepareVoteWithExceedLimit(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 21, suit.view.allNode[1])
	n, h := suit.view.firstProposer().HighestQCBlockBn()
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(h, n)
	prepareVote := mockPrepareVote(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 10,
		suit.view.secondProposerIndex(), block11.Hash(), block11.NumberU64(), oldQC)
	if err := suit.view.firstProposer().OnPrepareVote(suit.view.secondProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// Received duplicate prepareVote messages
// Legal prepareVote message
// The first Verification pass, the second prompt vote already exists, the total number of votes is 1
func TestPrepareVoteWithRepeat(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err != nil {
		t.Fatal(err.Error())
	}
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		if err.Error() != "prepare vote has exist(blockIndex:0, validatorIndex:0)" {
			t.Fatal("FAIL")
		}
	}
	if 1 != suit.view.secondProposer().state.PrepareVoteLenByIndex(0) {
		t.Fatal("FAIL")
	}
}

// duplicate sign
// Return duplicate sign error
func TestPrepareVoteDu(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	paths := createPaths(len(suit.view.allCbft), t)
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
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote1); err != nil {
		t.Fatal(err.Error())
	}
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote2); err == nil {
		t.Fatal("FAIL")
	} else {
		reg := regexp.MustCompile(`DuplicatePrepareVoteEvidence`)
		if len(reg.FindAllString(err.Error(), -1)) == 0 {
			t.Fatal(err.Error())
		}
	}
}

// viewNumber is less than the current viewNumber
// Verification failed
func TestPrepareVoteWithViewNumberTooLow(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	suit.view.secondProposer().state.ResetView(suit.epoch, suit.oldViewNumber+1)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		if err.Error() != "viewNumber too low(local:1, msg:0)" {
			t.Error("fail")
		}
	}
}

// viewNumber is greater than the current viewNumber
// Verification failed，Trigger synchronization
func TestPrepareVoteWithViewNumberTooHigh(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		if err.Error() != "viewNumber higher than local(local:0, msg:1)" {
			t.Error("fail")
		}
	}
}

// Vote's parent block did not reach prepareQC on this node
// Verification pass
func TestPrepareVoteWithParentIsNotParentQC(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	qc := mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().Node().ID().String(), prepareBlock1); err != nil {
		t.Fatal(err.Error())
	}
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	prepareBlock2 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock2); err != nil {
		t.Fatal("FAIL")
	}
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 1, suit.view.firstProposerIndex(), block2.Hash(),
		block2.NumberU64(), qc.BlockQC)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err != nil {
		t.Fatal(err.Error())
	}
}

// Vote's parent block did not reach prepareQC on the sending node (not legal)
// Verification failed
func TestPrepareVoteWithParentErrParentQC(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	qc := mockBlockQC(suit.view.allNode[0:1], suit.blockOne, 0, nil)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock1); err != nil {
		t.Fatal("FAIL")
	}
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	prepareBlock2 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock2); err != nil {
		t.Fatal("FAIL")
	}
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2.Hash(),
		block2.NumberU64(), qc.BlockQC)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// When the prepareQC is reached, there is a sub-block prepareVote
// Verification pass, sending sub-block votes
func TestPrepareVoteWithParentQCHasChild(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(block11.Hash(), 12, suit.view.allNode[1])
	prepareBlock12 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock12); err != nil {
		t.Fatal(err.Error())
	}
	block12QC := mockBlockQC(suit.view.allNode, block12, 1, block11QC.BlockQC)
	block13 := NewBlockWithSign(block12.Hash(), 13, suit.view.allNode[1])
	prepareBlock13 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 2,
		suit.view.secondProposerIndex(), block13, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().Node().ID().String(), prepareBlock13); err != nil {
		t.Fatal(err.Error())
	}
	suit.waitVote()
	if block13.NumberU64() != suit.view.firstProposer().state.PendingPrepareVote().Votes[0].BlockNum() {
		t.Error("fail")
	}
	suit.view.firstProposer().insertQCBlock(block12, block12QC.BlockQC)
	suit.view.firstProposer().trySendPrepareVote()
	if 0 != suit.view.firstProposer().state.PendingPrepareVote().Len() {
		t.Error("fail")
	}
}

// When the prepareQC is reached, there is no sub-block prepareVote
// Verify commit and lock
func TestPrepareVoteWithParentQCNotHasChild(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	suit.view.setBlockQC(5, suit.view.allNode[0])
	block6 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 6, suit.view.allNode[0])
	qc := mockBlockQC(suit.view.allNode, block6, 0, nil)
	suit.view.secondProposer().insertQCBlock(block6, qc.BlockQC)
	commitNumber, _ := suit.view.secondProposer().HighestCommitBlockBn()
	lockNumber, _ := suit.view.secondProposer().HighestLockBlockBn()
	if 4 != commitNumber && 5 != lockNumber {
		t.Error("fail")
	}
}

// Data valid timeout prepareVote message
// Verification failed
func TestPrepareVoteWithTimeout(t *testing.T) {
	suit := SetupPrepareVoteTest(3000)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	time.Sleep(time.Millisecond * 3000)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err == nil {
		t.Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The data just meets the 2f+1 prepareQC message
// Verification pass
func TestPrepareVote2fAndOne(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	qc := mockBlockQC(suit.view.allNode[0:3], suit.blockOne, 0, nil)
	if err := suit.view.secondProposer().verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc.BlockQC); err != nil {
		t.Fatal(err.Error())
	}
}

func (cbft *Cbft) generateErrPrepareQC(votes map[uint32]*protocols.PrepareVote) *ctypes.QuorumCert {
	if len(votes) == 0 {
		return nil
	}
	vote := votes[0]
	// Validator set prepareQC is the same as highestQC
	total := cbft.validatorPool.Len(cbft.state.Epoch())

	vSet := utils.NewBitArray(uint32(total))
	vSet.SetIndex(vote.NodeIndex(), true)

	var aggSig bls.Sign
	if err := aggSig.Deserialize(vote.Sign()); err != nil {
		return nil
	}
	qc := &ctypes.QuorumCert{
		Epoch:        vote.Epoch,
		ViewNumber:   vote.ViewNumber,
		BlockHash:    vote.BlockHash,
		BlockNumber:  vote.BlockNumber,
		BlockIndex:   vote.BlockIndex,
		ValidatorSet: utils.NewBitArray(vSet.Size()),
	}
	for i, p := range votes {
		if i != 0 {
			var sig bls.Sign
			err := sig.Deserialize(p.Sign())
			if err != nil {
				return nil
			}
			aggSig.Add(&sig)
			vSet.SetIndex(i, true)
		}
	}
	qc.Signature.SetBytes(aggSig.Serialize())
	qc.ValidatorSet.Update(vSet)
	return qc
}

// Non-consensus nodes receive prepareVote
// Verification pass
func TestPrepareVoteOfNotConsensus(t *testing.T) {
	suit := SetupPrepareVoteTest(10000)
	notConsensus := mockNotConsensusNode(suit.view.nodeParams, 1, testPeriod)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	notConsensus[0].engine.state.AddPrepareBlock(prepareBlock)
	if err := notConsensus[0].engine.OnPrepareVote(suit.view.firstProposer().Node().ID().String(), prepareVote); err != nil {
		t.Error(err.Error())
	}
}
