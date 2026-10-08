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
	"testing"
	"time"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
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
	suit.view = newTestView(false, 10000)
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
		ValidatorIndex: node.Index,
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

// mockViewQCSignedByOne builds a ViewChangeQC in which every sub-certificate is
// signed by one and the same validator, each over a different candidate block
// hash. It emulates a maliciously crafted certificate.
func mockViewQCSignedByOne(cbft *Cbft, epoch, viewNumber uint64, entries int) *ctypes.ViewChangeQC {
	index, err := cbft.validatorPool.GetIndexByNodeID(epoch, cbft.Node().ID())
	if err != nil {
		panic(err.Error())
	}
	total := uint32(cbft.validatorPool.Len(epoch))

	qc := &ctypes.ViewChangeQC{QCs: make([]*ctypes.ViewChangeQuorumCert, 0, entries)}
	for i := 0; i < entries; i++ {
		vc := mockViewChange(cbft.config.Option.BlsPriKey, epoch, viewNumber,
			common.BytesToHash([]byte(fmt.Sprintf("viewchange-forged-block-%d", i))),
			uint64(i+1), index, nil)

		var aggSig bls.Sign
		if err := aggSig.Deserialize(vc.Sign()); err != nil {
			panic(err.Error())
		}
		cert := &ctypes.ViewChangeQuorumCert{
			Epoch:        vc.Epoch,
			ViewNumber:   vc.ViewNumber,
			BlockHash:    vc.BlockHash,
			BlockNumber:  vc.BlockNumber,
			ValidatorSet: utils.NewBitArray(total),
		}
		cert.Signature.SetBytes(aggSig.Serialize())
		cert.ValidatorSet.SetIndex(index, true)
		qc.QCs = append(qc.QCs, cert)
	}
	return qc
}

// A ViewChangeQC whose sub-certificates are all signed by a single validator, one
// per candidate block hash. The summed signature count reaches the quorum
// threshold on its own, but only one distinct validator ever signed, so the
// certificate must not pass verification.
func (suit *VerifyQCTestSuite) TestVerifyViewChangeQCSignedByOneValidator() {
	attacker := suit.view.secondProposer()
	entries := attacker.validatorPool.Len(suit.epoch)
	threshold := attacker.threshold(entries)

	qc := mockViewQCSignedByOne(attacker, suit.epoch, suit.oldViewNumber, entries)

	// The summed count clears the threshold, only the distinct count must not.
	if got := qc.Len(); got < threshold {
		suit.T().Fatalf("precondition failed, expect summed signature count >= %d, got %d", threshold, got)
	}
	if got := countDistinctViewChangeSigners(qc); got != 1 {
		suit.T().Fatalf("precondition failed, expect 1 distinct validator, got %d", got)
	}

	if err := suit.view.firstProposer().verifyViewChangeQC(qc); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// A genuine ViewChangeQC in which the signers are spread over two candidate block
// hashes. Three distinct validators out of four reach the threshold, so this has
// to keep passing: a view change may legitimately be backed by validators that
// voted for different blocks.
func (suit *VerifyQCTestSuite) TestVerifyViewChangeQCDistinctSignersAcrossBlocks() {
	suit.insertOneBlock()

	nodes := suit.view.allNode[0:3]
	total := nodes[0].engine.validatorPool.Len(suit.epoch)

	// the second candidate block, voted for by the middle node only
	altBlock := NewBlock(suit.blockOne.ParentHash(), suit.blockOne.NumberU64()+1)

	viewChanges := make(map[uint32]*protocols.ViewChange)
	for i, node := range nodes {
		block := suit.blockOne
		if i == 1 {
			block = altBlock
		}
		index, err := node.engine.validatorPool.GetIndexByNodeID(suit.epoch, node.engine.Node().ID())
		if err != nil {
			panic(err.Error())
		}
		viewChanges[index] = mockViewChange(node.engine.config.Option.BlsPriKey, suit.epoch,
			suit.oldViewNumber, block.Hash(), block.NumberU64(), index, suit.blockOneQC.BlockQC)
	}

	qc := genViewChangeQC(uint32(total), viewChanges)

	if len(qc.QCs) != 2 {
		suit.T().Fatalf("precondition failed, expect 2 sub-certificates, got %d", len(qc.QCs))
	}
	if got := countDistinctViewChangeSigners(qc); got != 3 {
		suit.T().Fatalf("precondition failed, expect 3 distinct validators, got %d", got)
	}

	if err := suit.view.firstProposer().verifyViewChangeQC(qc); err != nil {
		suit.T().Fatal(err.Error())
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
	nodes := mockNotConsensusNode(suit.view.nodeParams, 4, testPeriod)
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
	nodes := mockNotConsensusNode(suit.view.nodeParams, 4, testPeriod)
	qc := mockBlockQCWithNotConsensus(nodes[0:3], suit.blockOne, 0, nil)
	if err := nodes[0].engine.verifyPrepareQC(suit.blockOne.NumberU64(), suit.blockOne.Hash(), qc.BlockQC); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}
