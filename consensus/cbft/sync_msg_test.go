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

	"github.com/stretchr/testify/suite"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

func TestSyncMsgTestSuite(t *testing.T) {
	suite.Run(t, new(SyncMsgTestSuite))
}

type SyncMsgTestSuite struct {
	suite.Suite
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
	epoch         uint64
	msgCh         chan *ctypes.MsgPackage
}

func (suit *SyncMsgTestSuite) SetupTest() {
	suit.view = newTestView(false, testNodeNumber)
	suit.blockOne = NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstProposer().state.ViewNumber()
	suit.epoch = suit.view.Epoch()
	msgCh := make(chan *ctypes.MsgPackage, 10240)
	suit.msgCh = msgCh
	f := func(msg *ctypes.MsgPackage) {
		select {
		case suit.msgCh <- msg:
		default:
			suit.T().Fatal("fail")
		}
	}
	for _, cbft := range suit.view.allCbft {
		cbft.network.SetSendQueueHook(f)
	}
}

func (suit *SyncMsgTestSuite) insertOneBlock(pb *protocols.PrepareBlock) {
	for _, cbft := range suit.view.allCbft {
		cbft.state.AddPrepareBlock(pb)
		cbft.findExecutableBlock()
	}
}
func (suit *SyncMsgTestSuite) insertOneQCBlock() {
	for _, cbft := range suit.view.allCbft {
		insertBlock(cbft, suit.blockOne, suit.blockOneQC.BlockQC)
	}
}

// normal
func (suit *SyncMsgTestSuite) TestSyncPrepareBlock() {
	pb := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.insertOneBlock(pb)
	prepareBlock := &protocols.GetPrepareBlock{
		Epoch:      suit.epoch,
		ViewNumber: suit.oldViewNumber,
		BlockIndex: 0,
	}
	cleanCh(suit.msgCh)
	if err := suit.view.firstProposer().OnGetPrepareBlock("", prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case <-suit.msgCh:
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// Epoch behind
// Epoch leading
// viewNumber leading
// Behind viewNumber
// blockIndex does not exist
func (suit *SyncMsgTestSuite) TestSyncPrepareBlockErrData() {
	pb := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.insertOneBlock(pb)
	testcases := []struct {
		name string
		data *protocols.GetPrepareBlock
	}{
		{name: "Epoch behind", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch - 1,
			ViewNumber: suit.oldViewNumber,
			BlockIndex: 0,
		}},
		{name: "Epoch leading", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch + 1,
			ViewNumber: suit.oldViewNumber,
			BlockIndex: 0,
		}},
		{name: "viewNumber leading", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch,
			ViewNumber: suit.oldViewNumber + 1,
			BlockIndex: 0,
		}},
		{name: "Behind viewNumber", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch,
			ViewNumber: math.MaxUint32,
			BlockIndex: 0,
		}},
		{name: "blockIndex does not exist", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch,
			ViewNumber: suit.oldViewNumber,
			BlockIndex: 1,
		}},
	}
	for _, testcase := range testcases {
		cleanCh(suit.msgCh)
		if err := suit.view.firstProposer().OnGetPrepareBlock("", testcase.data); err != nil {
			suit.T().Errorf("case-%s is failed,reson:%s", testcase.name, err.Error())
		}
		select {
		case <-suit.msgCh:
		case <-time.After(time.Millisecond * 10):
		}

	}
}

// normal
func (suit *SyncMsgTestSuite) TestOnGetBlockQuorumCert() {
	suit.insertOneQCBlock()
	getQC := &protocols.GetBlockQuorumCert{
		BlockHash:   suit.blockOne.Hash(),
		BlockNumber: suit.blockOne.NumberU64(),
	}
	cleanCh(suit.msgCh)
	if err := suit.view.firstProposer().OnGetBlockQuorumCert("", getQC); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case <-suit.msgCh:
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// wrong
func (suit *SyncMsgTestSuite) TestOnGetBlockQuorumCertErr() {
	fmt.Println(suit.view.genesisBlock.Root().String())
	suit.insertOneQCBlock()
	getQC := &protocols.GetBlockQuorumCert{
		BlockHash:   common.Hash{},
		BlockNumber: math.MaxUint64,
	}
	if err := suit.view.firstProposer().OnGetBlockQuorumCert("", getQC); err != nil {
		suit.T().Fatal(err.Error())
	}
}

// normal
func (suit *SyncMsgTestSuite) TestOnBlockQuorumCert() {
	pb := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.insertOneBlock(pb)
	time.Sleep(time.Millisecond * 30)
	if err := suit.view.secondProposer().OnBlockQuorumCert("", suit.blockOneQC); err != nil {
		suit.T().Fatal(err.Error())
	}
	if _, findQC := suit.view.secondProposer().blockTree.FindBlockAndQC(suit.blockOne.Hash(), 1); findQC == nil {
		suit.T().Fatal("fail")
	}
}

// normal Locally existing
func (suit *SyncMsgTestSuite) TestOnBlockQuorumCertExists() {
	suit.insertOneQCBlock()
	time.Sleep(time.Millisecond * 20)
	if err := suit.view.secondProposer().OnBlockQuorumCert("", suit.blockOneQC); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// normal Local does not exist in this block
func (suit *SyncMsgTestSuite) TestOnBlockQuorumCertBlockNotExists() {
	if err := suit.view.secondProposer().OnBlockQuorumCert("", suit.blockOneQC); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// wrong
// Epoch behind
// Epoch leading
// viewNumber leading
// Behind viewNumber
// Insufficient signature
func (suit *SyncMsgTestSuite) TestOnBlockQuorumCertErr() {
	pb := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.insertOneBlock(pb)
	time.Sleep(time.Millisecond * 20)
	testcases := []struct {
		name string
		data *protocols.BlockQuorumCert
	}{
		{
			name: "viewnumber big",
			data: mockBlockQCWithViewNumber(suit.view.allNode, suit.blockOne, 0, nil, 1),
		},
		{
			name: "viewnumber small",
			data: mockBlockQCWithViewNumber(suit.view.allNode, suit.blockOne, 0, nil, math.MaxUint64),
		},
		{
			name: "epoch big",
			data: mockBlockQCWithEpoch(suit.view.allNode, suit.blockOne, 0, nil, 2),
		},
		{
			name: "epoch small",
			data: mockBlockQCWithEpoch(suit.view.allNode, suit.blockOne, 0, nil, 0),
		},
		{
			name: "qc small number of signature",
			data: mockBlockQC(suit.view.allNode[:2], suit.blockOne, 0, nil),
		},
	}
	for _, testcase := range testcases {
		if err := suit.view.firstProposer().OnBlockQuorumCert("", testcase.data); err == nil {
			suit.T().Errorf("case:%s is failed", testcase.name)
		} else {
			fmt.Println(err.Error())
		}
	}
}

// Want a block
func (suit *SyncMsgTestSuite) TestOnGetQCBlockListWith1() {
	suit.view.setBlockQC(5, suit.view.allNode[0])
	lockBlock := suit.view.firstProposer().state.HighestLockBlock()
	getBlockList := &protocols.GetQCBlockList{
		BlockHash:   lockBlock.Hash(),
		BlockNumber: lockBlock.NumberU64(),
	}
	cleanCh(suit.msgCh)
	if err := suit.view.secondProposer().OnGetQCBlockList("", getBlockList); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case <-suit.msgCh:
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// Want two blocks
func (suit *SyncMsgTestSuite) TestOnGetQCBlockListWith2() {
	suit.view.setBlockQC(5, suit.view.allNode[0])
	commitBlock := suit.view.firstProposer().state.HighestCommitBlock()
	getBlockList := &protocols.GetQCBlockList{
		BlockHash:   commitBlock.Hash(),
		BlockNumber: commitBlock.NumberU64(),
	}
	if err := suit.view.secondProposer().OnGetQCBlockList("", getBlockList); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case <-suit.msgCh:
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// Want three blocks
func (suit *SyncMsgTestSuite) TestOnGetQCBlockListWith3() {
	suit.view.setBlockQC(3, suit.view.allNode[0])
	getBlockList := &protocols.GetQCBlockList{
		BlockHash:   suit.view.genesisBlock.Hash(),
		BlockNumber: suit.view.genesisBlock.NumberU64(),
	}
	if err := suit.view.secondProposer().OnGetQCBlockList("", getBlockList); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case <-suit.msgCh:
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// Want four blocks
func (suit *SyncMsgTestSuite) TestOnGetQCBlockListTooLow() {
	suit.view.setBlockQC(5, suit.view.allNode[0])
	getBlockList := &protocols.GetQCBlockList{
		BlockHash:   suit.view.genesisBlock.Hash(),
		BlockNumber: suit.view.genesisBlock.NumberU64(),
	}

	if err := suit.view.secondProposer().OnGetQCBlockList("", getBlockList); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// Want 0 blocks
func (suit *SyncMsgTestSuite) TestOnGetQCBlockListEqual() {
	suit.view.setBlockQC(5, suit.view.allNode[0])
	getBlockList := &protocols.GetQCBlockList{
		BlockHash:   suit.view.firstProposer().state.HighestQCBlock().Hash(),
		BlockNumber: suit.view.firstProposer().state.HighestQCBlock().NumberU64(),
	}
	if err := suit.view.secondProposer().OnGetQCBlockList("", getBlockList); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// Number and hash does not match
func (suit *SyncMsgTestSuite) TestOnGetQCBlockListDifNumber() {
	suit.view.setBlockQC(5, suit.view.allNode[0])
	getBlockList := &protocols.GetQCBlockList{
		BlockHash:   suit.view.firstProposer().state.HighestQCBlock().Hash(),
		BlockNumber: suit.view.firstProposer().state.HighestQCBlock().NumberU64() + 1,
	}
	if err := suit.view.secondProposer().OnGetQCBlockList("", getBlockList); err != nil {
		suit.T().Fatal(err.Error())
	}
}

// normal
func (suit *SyncMsgTestSuite) TestOnGetPrepareVote() {
	votes := make([]*protocols.PrepareVote, 0)
	for _, node := range suit.view.allCbft {
		index, err := node.validatorPool.GetIndexByNodeID(suit.epoch, node.config.Option.NodeID)
		if err != nil {
			panic(err.Error())
		}
		vote := mockPrepareVote(node.config.Option.BlsPriKey, suit.epoch, suit.oldViewNumber,
			0, index, suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil)
		votes = append(votes, vote)
		suit.view.firstProposer().state.AddPrepareVote(index, vote)
	}
	unKnownSet := utils.NewBitArray(uint32(len(suit.view.allCbft)))
	for i := uint32(0); i < unKnownSet.Size(); i++ {
		unKnownSet.SetIndex(i, true)
	}
	getPrepareVote := &protocols.GetPrepareVote{
		Epoch:      suit.epoch,
		ViewNumber: suit.oldViewNumber,
		BlockIndex: 0,
		UnKnownSet: unKnownSet,
	}
	cleanCh(suit.msgCh)
	suit.view.firstProposer().OnGetPrepareVote("", getPrepareVote)
	select {
	case <-suit.msgCh:
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// normal
func (suit *SyncMsgTestSuite) TestOnPrepareVotes() {
	pb := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.firstProposer().state.AddPrepareBlock(pb)
	votes := make([]*protocols.PrepareVote, 0)
	for _, node := range suit.view.allCbft {
		index, err := node.validatorPool.GetIndexByNodeID(suit.epoch, node.config.Option.NodeID)
		if err != nil {
			panic(err.Error())
		}
		vote := mockPrepareVote(node.config.Option.BlsPriKey, suit.epoch, suit.oldViewNumber,
			0, index, suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil)
		votes = append(votes, vote)
	}
	vs := &protocols.PrepareVotes{
		Epoch:      suit.epoch,
		ViewNumber: suit.oldViewNumber,
		BlockIndex: 0,
	}
	vs.Votes = append(vs.Votes, votes...)
	if err := suit.view.firstProposer().OnPrepareVotes("", vs); err != nil {
		suit.T().Fatal(err.Error())
	}
}

// Repeated
func (suit *SyncMsgTestSuite) TestOnPrepareVotesDup() {
	pb := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.firstProposer().state.AddPrepareBlock(pb)
	votes := make([]*protocols.PrepareVote, 0)
	for _, node := range suit.view.allCbft {
		index, err := node.validatorPool.GetIndexByNodeID(suit.epoch, node.config.Option.NodeID)
		if err != nil {
			panic(err.Error())
		}
		vote := mockPrepareVote(node.config.Option.BlsPriKey, suit.epoch, suit.oldViewNumber,
			0, index, suit.blockOne.Hash(), suit.blockOne.NumberU64(), nil)
		votes = append(votes, vote)
	}
	vs := &protocols.PrepareVotes{
		Epoch:      suit.epoch,
		ViewNumber: suit.oldViewNumber,
		BlockIndex: 0,
	}
	vs.Votes = append(vs.Votes, votes...)
	vs.Votes = append(vs.Votes, votes...)
	if err := suit.view.firstProposer().OnPrepareVotes("", vs); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

func cleanCh(ch chan *ctypes.MsgPackage) {
	for i := 0; i < len(ch); i++ {
		<-ch
	}
}
