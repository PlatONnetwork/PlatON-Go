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
	"testing"
	"time"
	"unsafe"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

func TestVerifyMsgTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyQCTestSuite))
}

type VerifyQCTestSuite struct {
	suite.Suite
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
	epoch         uint64
}

func (suit *VerifyQCTestSuite) SetupTest() {
	suit.view = newTestView(false, testNodeNumber)
	suit.blockOne = NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstProposer().state.ViewNumber()
	suit.epoch = suit.view.Epoch()
}

func (suit *VerifyQCTestSuite) insertOneBlock() {
	for _, cbft := range suit.view.allCbft {
		insertBlock(cbft, suit.blockOne, suit.blockOneQC.BlockQC)
	}
}

func (cbft *Cbft) mockGenerateViewChangeQuorumCert(v *protocols.ViewChange, index uint32) (*ctypes.ViewChangeQuorumCert, error) {
	// node, err := cbft.isCurrentValidator()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "local node is not validator")
	// }
	total := uint32(cbft.validatorPool.Len(cbft.state.Epoch()))
	var aggSig bls.Sign
	if err := aggSig.Deserialize(v.Sign()); err != nil {
		return nil, err
	}

	blockEpoch, blockView := uint64(0), uint64(0)
	if v.PrepareQC != nil {
		blockEpoch, blockView = v.PrepareQC.Epoch, v.PrepareQC.ViewNumber
	}
	cert := &ctypes.ViewChangeQuorumCert{
		Epoch:           v.Epoch,
		ViewNumber:      v.ViewNumber,
		BlockHash:       v.BlockHash,
		BlockNumber:     v.BlockNumber,
		BlockEpoch:      blockEpoch,
		BlockViewNumber: blockView,
		ValidatorSet:    utils.NewBitArray(total),
	}
	cert.Signature.SetBytes(aggSig.Serialize())
	cert.ValidatorSet.SetIndex(index, true)
	return cert, nil
}
func (cbft *Cbft) mockGenerateViewChangeQuorumCertWithViewNumber(qc *ctypes.QuorumCert) (*ctypes.ViewChangeQuorumCert, error) {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		return nil, errors.Wrap(err, "local node is not validator")
	}
	v := &protocols.ViewChange{
		Epoch:          qc.Epoch,
		ViewNumber:     qc.ViewNumber,
		BlockHash:      qc.BlockHash,
		BlockNumber:    qc.BlockNumber,
		ValidatorIndex: uint32(node.Index),
		PrepareQC:      qc,
	}
	if err := cbft.signMsgByBls(v); err != nil {
		return nil, errors.Wrap(err, "Sign ViewChange failed")
	}

	total := uint32(cbft.validatorPool.Len(cbft.state.Epoch()))
	var aggSig bls.Sign
	if err := aggSig.Deserialize(v.Sign()); err != nil {
		return nil, err
	}
	cert := &ctypes.ViewChangeQuorumCert{
		Epoch:           qc.Epoch,
		ViewNumber:      qc.ViewNumber,
		BlockHash:       qc.BlockHash,
		BlockNumber:     qc.BlockNumber,
		BlockEpoch:      qc.Epoch,
		BlockViewNumber: qc.ViewNumber,
		ValidatorSet:    utils.NewBitArray(total),
	}
	cert.Signature.SetBytes(aggSig.Serialize())
	cert.ValidatorSet.SetIndex(node.Index, true)
	return cert, nil
}

// Normal viewChangeQC message
// Verification pass
func (suit *VerifyQCTestSuite) TestVerifyViewChangeQC() {
	qc := mockViewQC(suit.view.genesisBlock, suit.view.allNode[0:3], nil)
	if err := suit.view.firstProposer().verifyViewChangeQC(qc); err != nil {
		suit.T().Fatal(err.Error())
	}
}

// Insufficient viewChangeQC message
// Verification cannot pass
func (suit *VerifyQCTestSuite) TestVerifyViewChangeQCErrNum() {
	qc := mockViewQC(suit.view.genesisBlock, suit.view.allNode[0:2], nil)
	if err := suit.view.firstProposer().verifyViewChangeQC(qc); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// False a sufficient number of viewChangeQC messages
// Verification cannot pass
func (suit *VerifyQCTestSuite) TestVerifyViewChangeQCErrNodeNum() {
	suit.insertOneBlock()
	view := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber, suit.blockOne.Hash(), suit.blockOne.NumberU64(), 0, suit.blockOneQC.BlockQC)
	vs := &ctypes.ViewChangeQC{}
	for i := 0; i <= 4; i++ {
		qc, err := suit.view.secondProposer().generateViewChangeQuorumCert(view)
		if err != nil {
			panic(err.Error())
		}
		vs.QCs = append(vs.QCs, qc)
	}
	if err := suit.view.firstProposer().verifyViewChangeQC(vs); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// viewChangeQC composed of viewChange with different viewNumber
// Verification cannot pass
func (suit *VerifyQCTestSuite) TestVerifyChangeQCViewChangeDifViewNumber() {
	suit.insertOneBlock()
	blocks := make([]*types.Block, 0)
	blocks = append(blocks, suit.blockOne)
	for i := uint64(1); i <= 3; i++ {
		blocks = append(blocks, NewBlock(blocks[len(blocks)-1].Hash(), i+1))
	}
	qcs := make([]*protocols.BlockQuorumCert, 0)
	qcs = append(qcs, suit.blockOneQC)
	for i, b := range blocks[1:] {
		qcs = append(qcs, mockBlockQCWithViewNumber(suit.view.allNode, b, 0, qcs[len(qcs)-1].BlockQC, uint64(i+1)))
	}
	vs := &ctypes.ViewChangeQC{}
	for i, qc := range qcs {
		cert, err := suit.view.allCbft[i].mockGenerateViewChangeQuorumCertWithViewNumber(qc.BlockQC)
		if err != nil {
			panic(err.Error())
		}
		vs.QCs = append(vs.QCs, cert)
	}
	if err := suit.view.firstProposer().verifyViewChangeQC(vs); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// False a sufficient number of viewChangeQC messages generated by non-consensus nodes
// Verification cannot pass
func (suit *VerifyQCTestSuite) TestVerifyViewChangeQCErrData() {
	suit.insertOneBlock()
	nodes := mockNotConsensusNode(false, suit.view.nodeParams, 4)
	view := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber, suit.blockOne.Hash(), suit.blockOne.NumberU64(), 0, suit.blockOneQC.BlockQC)
	vs := &ctypes.ViewChangeQC{}
	for i := 0; i < 4; i++ {
		cert, err := nodes[i].engine.mockGenerateViewChangeQuorumCert(view, uint32(i))
		if err != nil {
			panic(err.Error())
		}
		vs.QCs = append(vs.QCs, cert)
	}
	if err := suit.view.firstProposer().verifyViewChangeQC(vs); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// Forge a prepareQC message that is much larger than the number of consensus nodes
// Execution time is less than or equal to 10ms
func (suit *VerifyQCTestSuite) TestSyncViewChangeQCTooBig() {
	suit.insertOneBlock()
	view := mockViewChange(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber, suit.blockOne.Hash(), suit.blockOne.NumberU64(), 0, suit.blockOneQC.BlockQC)
	vs := &ctypes.ViewChangeQC{}
	var qcs [1000]*ctypes.ViewChangeQuorumCert
	for i := 0; i < len(qcs); i++ {
		qc, err := suit.view.secondProposer().generateViewChangeQuorumCert(view)
		if err != nil {
			panic(err.Error())
		}
		qcs[i] = qc
		vs.QCs = append(vs.QCs, qc)
	}
	fmt.Println(unsafe.Sizeof(qcs))
	start := time.Now()
	if err := suit.view.firstProposer().verifyViewChangeQC(vs); err != nil {
		fmt.Println(err.Error())
	}
	end := time.Since(start)
	if end > time.Millisecond*10 {
		suit.T().Fatal("Execution time is too long")
	}
}

// Forge a prepareQC message that is much larger than the number of consensus nodes
// Execution time is less than or equal to 10ms
func (suit *VerifyQCTestSuite) TestPrepareQCTooBig() {
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		0, suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	votes := make(map[uint32]*protocols.PrepareVote)
	for i := uint32(0); i <= 1000; i++ {
		votes[i] = prepareVote
	}
	qc := suit.view.firstProposer().generateErrPrepareQC(votes)
	fmt.Println(unsafe.Sizeof(*qc))
	start := time.Now()
	if err := suit.view.secondProposer().verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc); err != nil {
		fmt.Println(err.Error())
	}
	end := time.Since(start)
	if end > time.Millisecond*10 {
		suit.T().Fatal("Execution time is too long")
	}
}

// Normal prepareQC
// Verification pass
func (suit *VerifyQCTestSuite) TestVerifyPrepareQC() {
	qc := mockBlockQC(suit.view.allNode[0:3], suit.blockOne, 0, nil)
	if err := suit.view.firstProposer().verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc.BlockQC); err != nil {
		suit.T().Fatal(err.Error())
	}
}

// Insufficient prepareVote generated by a small number of prepareQC
// Verification failed
func (suit *VerifyQCTestSuite) TestVerifyPrepareQCErrNum() {
	qc := mockBlockQC(suit.view.allNode[0:2], suit.blockOne, 0, nil)
	if err := suit.view.firstProposer().verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc.BlockQC); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// False a sufficient number of prepareQC messages
// Verification cannot pass
func (suit *VerifyQCTestSuite) TestPrepareVoteErr() {
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		0, suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	votes := make(map[uint32]*protocols.PrepareVote)
	votes[0] = prepareVote
	votes[1] = prepareVote
	votes[2] = prepareVote
	qc := suit.view.firstProposer().generateErrPrepareQC(votes)
	if err := suit.view.secondProposer().verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// prepareQC generated using prepareVote of different blocks
// Verification failed
func (suit *VerifyQCTestSuite) TestPrepareVoteErr2() {
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		0, suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	prepareVote2 := mockPrepareVote(suit.view.secondProposerBlsKey(), suit.epoch, suit.oldViewNumber, 1,
		1, block2.Hash(),
		block2.NumberU64(), nil)
	block3 := NewBlock(block2.Hash(), 3)
	prepareVote3 := mockPrepareVote(suit.view.thirdProposer().config.Option.BlsPriKey, suit.epoch, suit.oldViewNumber, 2,
		2, block3.Hash(),
		block3.NumberU64(), nil)
	votes := make(map[uint32]*protocols.PrepareVote)
	votes[0] = prepareVote
	votes[1] = prepareVote2
	votes[2] = prepareVote3
	qc := suit.view.firstProposer().generatePrepareQC(votes)
	if err := suit.view.secondProposer().verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// Forging a sufficient number of prepareQC messages generated by non-consensus nodes
// Verification cannot pass
func (suit *VerifyQCTestSuite) TestVerifyPrepareQCErrData() {
	suit.insertOneBlock()
	nodes := mockNotConsensusNode(false, suit.view.nodeParams, 4)
	qc := mockBlockQCWithNotConsensus(nodes[0:3], suit.blockOne, 0, nil)
	if err := nodes[0].engine.verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc.BlockQC); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}
