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
	"regexp"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/suite"
)

func TestPrepareVoteSuite(t *testing.T) {
	suite.Run(t, new(PrepareVoteTestSuite))
}

type PrepareVoteTestSuite struct {
	suite.Suite
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
	epoch         uint64
}

func (suit *PrepareVoteTestSuite) SetupTest() {
	suit.view = newTestView(false, testNodeNumber)
	suit.blockOne = NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstProposer().state.ViewNumber()
	suit.epoch = suit.view.Epoch()
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
func (suit *PrepareVoteTestSuite) TestBuildPrepareVote() {
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
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	vote := suit.view.firstProposer().state.AllPrepareVoteByIndex(1)[0]
	suit.Equal(uint64(12), vote.BlockNum())
	suit.Equal(block12.Hash().String(), vote.BlockHash.String())
	suit.Equal(prepareBlock.BlockIndex, vote.BlockIndex)

}

// prepareVote message basic check
// 1.Unsigned
// 2.The signature is inconsistent with the validatorIndex
// 3.The signature is not the verification node
// 4.epoch too big
// 5.epoch too small
func (suit *PrepareVoteTestSuite) TestCheckErrPrepareVote() {
	_, notConsensusKey := GenerateKeys(1)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal("FAIL")
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
		if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), testcase.data); err == nil {
			suit.T().Errorf("case %s is failed", testcase.name)
		} else {
			fmt.Println(err.Error())
		}
	}
}

// ParentVote message that does not carry ParentQC with zero parent block
// Verification pass
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentIsZeroButNotParentQC() {
	epoch := suit.view.Epoch()
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	// if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), prepareBlock); err != nil {
	// 	suit.T().Fatal(err.Error())
	// }
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(),
		suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(suit.view.firstProposerIndex(), suit.view.secondProposer().state.AllPrepareVoteByIndex(0)[0].ValidatorIndex)
}

// The parent block is non-zero and does not carry the parentVC prepareVote message.
// Verification failed
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentIsNotZeroButNotParentQC() {
	suit.insertOneBlock()
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 1, suit.view.firstProposerIndex(), block2.Hash(),
		block2.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The parent block is non-zero but the blockIndex is 0. The parentVote message does not carry ParentQC.
// Verification failed
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentIsNotZeroAndBlockIndexNotParentQC() {
	suit.view.setBlockQC(10, suit.view.allNode[0])
	n, h := suit.view.firstProposer().HighestQCBlockBn()
	block1 := NewBlockWithSign(h, n+1, suit.view.allNode[1])
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(h, n)
	fmt.Println(qc.String())
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block1, qc, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 0, suit.view.firstProposerIndex(), block1.Hash(),
		block1.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// blockNumber=1, qc forged message
// Verification failed
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithBlockNumberIsOneAndErrParentQC() {
	block1 := NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 4)
	errQC := mockErrBlockQC(notConsensusNodes, block1, 0, nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), block1, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock("", prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), block1.Hash(), 1, errQC.BlockQC)
	if err := suit.view.secondProposer().OnPrepareVote("", prepareVote); err == nil {
		suit.T().Fatal("fail")
	}
}

// Received the prepareVote message received by prepareBlock first
// Verification failed
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithNotPrepareBlock() {
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The prepareVote message that the block exceeds the limit of one round of the block number
// Verification failed
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithExceedLimit() {
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 21, suit.view.allNode[1])
	n, h := suit.view.firstProposer().HighestQCBlockBn()
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(h, n)
	prepareVote := mockPrepareVote(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 10,
		suit.view.secondProposerIndex(), block11.Hash(), block11.NumberU64(), oldQC)
	if err := suit.view.firstProposer().OnPrepareVote(suit.view.secondProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// Received duplicate prepareVote messages
// Legal prepareVote message
// The first Verification pass, the second prompt vote already exists, the total number of votes is 1
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithRepeat() {
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		suit.Equal(err.Error(), "prepare vote has exist(blockIndex:0, validatorIndex:0)")
	}
	suit.Equal(1, suit.view.secondProposer().state.PrepareVoteLenByIndex(0))
}

// duplicate sign
// Return duplicate sign error
func (suit *PrepareVoteTestSuite) TestPrepareVoteDu() {
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
	} else {
		reg := regexp.MustCompile(`DuplicatePrepareVoteEvidence`)
		if len(reg.FindAllString(err.Error(), -1)) == 0 {
			suit.T().Fatal(err.Error())
		}

	}
}

// viewNumber is less than the current viewNumber
// Verification failed
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithViewNumberTooLow() {
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	suit.view.secondProposer().state.ResetView(suit.epoch, suit.oldViewNumber+1)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		suit.Equal(err.Error(), "viewNumber too low(local:1, msg:0)")
	}
}

// viewNumber is greater than the current viewNumber
// Verification failed，Trigger synchronization
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithViewNumberTooHigh() {
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		suit.Equal(err.Error(), "viewNumber higher than local(local:0, msg:1)")
	}
}

// Vote's parent block did not reach prepareQC on this node
// Verification pass
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentIsNotParentQC() {
	qc := mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal(err.Error())
	}
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	prepareBlock2 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock2); err != nil {
		suit.T().Fatal("FAIL")
	}
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 1, suit.view.firstProposerIndex(), block2.Hash(),
		block2.NumberU64(), qc.BlockQC)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err != nil {
		suit.T().Fatal(err.Error())
	}
}

// Vote's parent block did not reach prepareQC on the sending node (not legal)
// Verification failed
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentErrParentQC() {
	qc := mockBlockQC(suit.view.allNode[0:1], suit.blockOne, 0, nil)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal("FAIL")
	}
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	prepareBlock2 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock2); err != nil {
		suit.T().Fatal("FAIL")
	}
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2.Hash(),
		block2.NumberU64(), qc.BlockQC)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// When the prepareQC is reached, there is a sub-block prepareVote
// Verification pass, sending sub-block votes
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentQCHasChild() {
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
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock12); err != nil {
		suit.T().Fatal(err.Error())
	}
	block12QC := mockBlockQC(suit.view.allNode, block12, 1, block11QC.BlockQC)
	block13 := NewBlockWithSign(block12.Hash(), 13, suit.view.allNode[1])
	prepareBlock13 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 2,
		suit.view.secondProposerIndex(), block13, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock13); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	suit.Equal(block13.NumberU64(), suit.view.firstProposer().state.PendingPrepareVote().Votes[0].BlockNum())
	suit.view.firstProposer().insertQCBlock(block12, block12QC.BlockQC)
	suit.view.firstProposer().trySendPrepareVote()
	suit.Equal(0, suit.view.firstProposer().state.PendingPrepareVote().Len())
}

// When the prepareQC is reached, there is no sub-block prepareVote
// Verify commit and lock
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentQCNotHasChild() {
	suit.view.setBlockQC(5, suit.view.allNode[0])
	block6 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 6, suit.view.allNode[0])
	qc := mockBlockQC(suit.view.allNode, block6, 0, nil)
	suit.view.secondProposer().insertQCBlock(block6, qc.BlockQC)
	commitNumber, _ := suit.view.secondProposer().HighestCommitBlockBn()
	lockNumber, _ := suit.view.secondProposer().HighestLockBlockBn()
	suit.Equal(uint64(4), commitNumber)
	suit.Equal(uint64(5), lockNumber)
}

// Data valid timeout prepareVote message
// Verification failed
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithTimeout() {
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	time.Sleep(time.Millisecond * testPeriod)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// The data just meets the 2f+1 prepareQC message
// Verification pass
func (suit *PrepareVoteTestSuite) TestPrepareVote2fAndOne() {
	qc := mockBlockQC(suit.view.allNode[0:3], suit.blockOne, 0, nil)
	if err := suit.view.secondProposer().verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc.BlockQC); err != nil {
		suit.T().Fatal(err.Error())
	}
}

func (cbft *Cbft) generateErrPrepareQC(votes map[uint32]*protocols.PrepareVote) *ctypes.QuorumCert {
	if len(votes) == 0 {
		return nil
	}
	var vote *protocols.PrepareVote
	vote = votes[0]
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
func (suit *PrepareVoteTestSuite) TestPrepareVoteOfNotConsensus() {
	notConsensus := mockNotConsensusNode(false, suit.view.nodeParams, 1)
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	notConsensus[0].engine.state.AddPrepareBlock(prepareBlock)
	if err := notConsensus[0].engine.OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err != nil {
		suit.T().Error(err.Error())
	}
}
