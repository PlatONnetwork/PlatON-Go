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

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/suite"
)

func TestViewChangeSuite(t *testing.T) {
	suite.Run(t, new(ViewChangeTestSuite))
}

type ViewChangeTestSuite struct {
	suite.Suite
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
}

func (suit *ViewChangeTestSuite) insertOneBlock() {
	for _, cbft := range suit.view.allCbft {
		insertBlock(cbft, suit.blockOne, suit.blockOneQC.BlockQC)
	}
}

func (suit *ViewChangeTestSuite) createEvPool(paths []string) {
	if len(paths) != len(suit.view.allNode) {
		panic("paths len err")
	}
	for i, path := range paths {
		pool, _ := evidence.NewBaseEvidencePool(path)
		suit.view.allCbft[i].evPool = pool
	}

}

func (suit *ViewChangeTestSuite) SetupTest() {
	suit.view = newTestView(false, testNodeNumber)
	suit.blockOne = NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstCbft.state.ViewNumber()
}

// Initiate viewChange
// Verify local viewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeBuild() {
	time.Sleep((testPeriod + 200) * time.Millisecond)
	suit.Equal(1, suit.view.secondProposer().state.ViewChangeLen())
}

// Initiate viewChange
// Non-consensus nodes do not Initiate viewChange, check local viewChangeLen=0
func (suit *ViewChangeTestSuite) TestViewChangeBuildWithNotConsensus() {
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 1)
	time.Sleep((testPeriod + 200) * time.Millisecond)
	suit.Equal(0, notConsensusNodes[0].engine.state.ViewChangeLen())
}

// viewChange message basic check
// 1.Missing signed viewChange message
// 2.blockChange message that blockNumber does not match blockHash
// 3.Signature inconsistency with the certifier viewChange message
// 4.Proposal non-consensus node viewChange message
// 5.prepareQC forged viewChange message
// 6.prepareQC with blockNumber not zero is an empty viewChange message
// 7.Carrying prepareQC does not satisfy N-f
// 8.Epoch is larger than local
// 9.Epoch is smaller than local
func (suit *ViewChangeTestSuite) TestViewChangeCheckErr() {
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 4)
	errQC := mockErrBlockQC(notConsensusNodes, suit.view.genesisBlock, 0, nil)
	notEmpty := mockBlockQC(suit.view.allNode[0:1], suit.view.genesisBlock, 0, nil)
	suit.insertOneBlock()
	testcases := []struct {
		name string
		data *protocols.ViewChange
		err  error
	}{
		{
			name: "Missing signed viewChange message",
			data: mockViewChange(nil, suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(), suit.blockOne.Hash(),
				suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC),
		},
		{
			name: "blockChange message that blockNumber does not match blockHash",
			data: mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
				suit.blockOne.Hash(), suit.blockOne.NumberU64()+1, suit.view.secondProposerIndex(),
				suit.blockOneQC.BlockQC),
		},
		{
			name: "Signature inconsistency with the certifier viewChange message",
			data: mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
				suit.blockOne.Hash(), suit.blockOne.NumberU64(), suit.view.secondProposerIndex()+1,
				suit.blockOneQC.BlockQC),
		},
		{
			name: "Proposal non-consensus node viewChange message",
			data: mockViewChange(notConsensusNodes[0].engine.config.Option.BlsPriKey, suit.view.Epoch(),
				suit.view.secondProposer().state.ViewNumber(), suit.blockOne.Hash(),
				suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC),
		},
		{
			name: "prepareQC forged viewChange message",
			data: mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
				suit.blockOne.Hash(), suit.blockOne.NumberU64(), suit.view.secondProposerIndex(),
				errQC.BlockQC),
		},
		{
			name: "prepareQC is empty viewChange message",
			data: mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
				suit.blockOne.Hash(), suit.blockOne.NumberU64(),
				suit.view.secondProposerIndex(), nil),
		},
		{
			name: "prepareQC is not 2f+1 viewChange message",
			data: mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
				suit.blockOne.Hash(), suit.blockOne.NumberU64(),
				suit.view.secondProposerIndex(), notEmpty.BlockQC),
		},
		{
			name: "epoch big",
			data: mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch()+1, suit.view.secondProposer().state.ViewNumber(),
				suit.blockOne.Hash(), suit.blockOne.NumberU64(),
				suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC),
		},
		{
			name: "epoch small",
			data: mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch()-1, suit.view.secondProposer().state.ViewNumber(),
				suit.blockOne.Hash(), suit.blockOne.NumberU64(),
				suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC),
		},
	}
	for _, testcase := range testcases {
		if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), testcase.data); err == nil {
			suit.T().Errorf("CASE:%s is failefd", testcase.name)
		} else {
			fmt.Println(err.Error())
		}
	}
}

// Block and HighestQCBlock local consistent viewChange message
// Verification passed, ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeCheckCorrect() {
	suit.insertOneBlock()
	viewChange := mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
		suit.blockOne.Hash(), suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC)
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(1, suit.view.firstProposer().state.ViewChangeLen())
}

// prepareQC with blockNumber zero is an empty viewChange message
// Verification passed, ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeCheckZero() {
	viewChange := mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
		suit.view.genesisBlock.Hash(), suit.view.genesisBlock.NumberU64(), suit.view.secondProposerIndex(), nil)
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(1, suit.view.firstProposer().state.ViewChangeLen())
}

// prepareQC with blockNumber zero is not empty viewChange message
// Verification passed, ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeCheckZeroPrepareQCNotNil() {
	suit.view.setBlockQC(9, suit.view.allNode[0])
	_, h := suit.view.firstProposer().HighestQCBlockBn()
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 4)
	block := NewBlockWithSign(h, 12, suit.view.allNode[1])
	errQC := mockErrBlockQC(notConsensusNodes, block, 8, nil)
	errViewChange := mockViewChange(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
		h, 0, suit.view.firstProposerIndex(), errQC.BlockQC)
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), errViewChange); err == nil {
		suit.T().Fatal("fail")
	}
}

// Block leads the local HighestQCBlock viewChange message
// Verification passed, ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeLeadHighestQCBlock() {
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	suit.view.firstProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewChange := mockViewChange(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.view.firstProposer().state.ViewNumber(), block2.Hash(),
		block2.NumberU64(), suit.view.firstProposerIndex(), block2QC.BlockQC)
	if err := suit.view.secondProposer().OnViewChange(suit.view.firstProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(1, suit.view.secondProposer().state.ViewChangeLen())
}

// Block behind the local HighestQCBlock viewChange message
// Verification passed, ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeBehindHighestQCBlock() {
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	suit.view.secondProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewChange := mockViewChange(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.view.firstProposer().state.ViewNumber(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), suit.view.firstProposerIndex(), suit.blockOneQC.BlockQC)
	if err := suit.view.secondProposer().OnViewChange(suit.view.firstProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(1, suit.view.secondProposer().state.ViewChangeLen())
}

// viewNumber is less than the current viewNumber
// The verification fails, ViewChangeLen=0
func (suit *ViewChangeTestSuite) TestViewChangeViewNumberBehind() {
	suit.insertOneBlock()
	suit.view.secondProposer().state.ResetView(1, 2)
	viewChange := mockViewChange(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.view.firstProposer().state.ViewNumber(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), suit.view.firstProposerIndex(), suit.blockOneQC.BlockQC)
	if err := suit.view.secondProposer().OnViewChange(suit.view.firstProposer().NodeID().String(), viewChange); err == nil {
		suit.T().Fatal("FAIL")
	} else if err.Error() != "viewNumber too low(local:2, msg:0)" {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(0, suit.view.secondProposer().state.ViewChangeLen())
}

// viewNumber is greater than the current viewNumber
// The verification fails, ViewChangeLen=0
func (suit *ViewChangeTestSuite) TestViewChangeViewNumberLead() {
	suit.insertOneBlock()
	viewChange := mockViewChange(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.view.firstProposer().state.ViewNumber()+1, suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), suit.view.firstProposerIndex(), suit.blockOneQC.BlockQC)
	if err := suit.view.secondProposer().OnViewChange(suit.view.firstProposer().NodeID().String(), viewChange); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		suit.EqualValues("viewNumber higher than local(local:0, msg:1)", err.Error())
	}
	suit.Equal(0, suit.view.secondProposer().state.ViewChangeLen())
}

// Received a viewChange message that has been processed
// Verification passed, ViewChangeLen unchanged
func (suit *ViewChangeTestSuite) TestCheckCorrectViewChangeRepeat() {
	suit.insertOneBlock()
	viewChange := mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC)

	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(1, suit.view.firstProposer().state.ViewChangeLen())
}

// The same person, viewChange message based on different blocks
// The verification fails, and the error of returning double viewChange is returned.
func (suit *ViewChangeTestSuite) TestViewChangeRepeatWithDifBlock() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.insertOneBlock()
	viewChange1 := mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
		suit.blockOne.Hash(), suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC)
	viewChange2 := mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
		suit.view.genesisBlock.Hash(), suit.view.genesisBlock.NumberU64(), suit.view.secondProposerIndex(), nil)
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange1); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange2); err == nil {
		suit.T().Fatal("fail")
	} else {
		reg := regexp.MustCompile(`DuplicateViewChangeEvidence`)
		if len(reg.FindAllString(err.Error(), -1)) == 0 {
			suit.T().Fatal(err.Error())
		}
	}
	suit.Equal(1, suit.view.firstProposer().state.ViewChangeLen())
}

// Non-consensus nodes receive viewChange
// Verification pass
func (suit *ViewChangeTestSuite) TestViewChangeNotConsensus() {
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 1)
	suit.insertOneBlock()
	viewChange := mockViewChange(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.view.firstProposer().state.ViewNumber(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), suit.view.firstProposerIndex(), suit.blockOneQC.BlockQC)
	if err := notConsensusNodes[0].engine.OnViewChange(suit.view.firstProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Error(err.Error())
	}
}
