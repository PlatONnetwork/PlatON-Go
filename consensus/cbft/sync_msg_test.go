package cbft

import (
	"fmt"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/suite"
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
	msgCh := make(chan *ctypes.MsgPackage, 100)
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

// 正常的
func (suit *SyncMsgTestSuite) TestSyncPrepareBlock() {
	pb := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.insertOneBlock(pb)
	prepareBlock := &protocols.GetPrepareBlock{
		Epoch:      suit.epoch,
		ViewNumber: suit.oldViewNumber,
		BlockIndex: 0,
	}
	if err := suit.view.firstProposer().OnGetPrepareBlock("", prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case m := <-suit.msgCh:
		suit.EqualValues(pb.String(), m.Message().String())
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// epoch落后的
// epoch领先的
// viewNumber领先的
// viewNumber落后的
// blockIndex不存在的
func (suit *SyncMsgTestSuite) TestSyncPrepareBlockErrData() {
	pb := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.insertOneBlock(pb)
	testcases := []struct {
		name string
		data *protocols.GetPrepareBlock
	}{
		{name: "epoch落后的", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch - 1,
			ViewNumber: suit.oldViewNumber,
			BlockIndex: 0,
		}},
		{name: "epoch领先的", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch + 1,
			ViewNumber: suit.oldViewNumber,
			BlockIndex: 0,
		}},
		{name: "viewNumber领先的", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch,
			ViewNumber: suit.oldViewNumber + 1,
			BlockIndex: 0,
		}},
		{name: "viewNumber落后的", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch,
			ViewNumber: math.MaxUint32,
			BlockIndex: 0,
		}},
		{name: "blockIndex不存在的", data: &protocols.GetPrepareBlock{
			Epoch:      suit.epoch,
			ViewNumber: suit.oldViewNumber,
			BlockIndex: 1,
		}},
	}
	for _, testcase := range testcases {
		select {
		case <-suit.msgCh:
		case <-time.After(time.Millisecond * 10):
		}
		if err := suit.view.firstProposer().OnGetPrepareBlock("", testcase.data); err != nil {
			suit.T().Errorf("case-%s is failed,reson:%s", testcase.name, err.Error())
		}
		select {
		case m := <-suit.msgCh:
			if v, ok := m.Message().(*protocols.PrepareBlock); ok {
				suit.T().Error(v.String())
			}
		case <-time.After(time.Millisecond * 10):
		}

	}
}

// 正常的
func (suit *SyncMsgTestSuite) TestOnGetBlockQuorumCert() {
	suit.insertOneQCBlock()
	getQC := &protocols.GetBlockQuorumCert{
		BlockHash:   suit.blockOne.Hash(),
		BlockNumber: suit.blockOne.NumberU64(),
	}
	if err := suit.view.firstProposer().OnGetBlockQuorumCert("", getQC); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case m := <-suit.msgCh:
		suit.EqualValues(suit.blockOneQC.String(), m.Message().String())
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// 错误的
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
	select {
	case <-suit.msgCh:
		suit.T().Fatal("fail")
	case <-time.After(time.Millisecond * 10):
	}
}

// 正常的
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

// 正常的 本地已经存在的
func (suit *SyncMsgTestSuite) TestOnBlockQuorumCertExists() {
	suit.insertOneQCBlock()
	time.Sleep(time.Millisecond * 20)
	if err := suit.view.secondProposer().OnBlockQuorumCert("", suit.blockOneQC); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 正常的 本地不存在这个块的
func (suit *SyncMsgTestSuite) TestOnBlockQuorumCertBlockNotExists() {
	if err := suit.view.secondProposer().OnBlockQuorumCert("", suit.blockOneQC); err == nil {
		suit.T().Fatal("fail")
	} else {
		fmt.Println(err.Error())
	}
}

// 错误的
// epoch落后的
// epoch领先的
// viewNumber领先的
// viewNumber落后的
// 签名不足的
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

// 要一个块
func (suit *SyncMsgTestSuite) TestOnGetQCBlockListWith1() {
	suit.view.setBlockQC(5, suit.view.allNode[0])
	lockBlock := suit.view.firstProposer().state.HighestLockBlock()
	getBlockList := &protocols.GetQCBlockList{
		BlockHash:   lockBlock.Hash(),
		BlockNumber: lockBlock.NumberU64(),
	}
	if err := suit.view.secondProposer().OnGetQCBlockList("", getBlockList); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case m := <-suit.msgCh:
		if msg, ok := m.Message().(*protocols.QCBlockList); ok {
			suit.Equal(1, len(msg.QC))
			suit.Equal(1, len(msg.Blocks))
		}
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// 要两个块
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
	case m := <-suit.msgCh:
		if msg, ok := m.Message().(*protocols.QCBlockList); ok {
			suit.Equal(2, len(msg.QC))
			suit.Equal(2, len(msg.Blocks))
		}
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// 要三个块
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
	case m := <-suit.msgCh:
		if msg, ok := m.Message().(*protocols.QCBlockList); ok {
			suit.Equal(3, len(msg.QC))
			suit.Equal(3, len(msg.Blocks))
		}
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// 要四个块
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

// 要0个块
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

// number与hash不匹配的
func (suit *SyncMsgTestSuite) TestOnGetQCBlockListDifNumber() {
	suit.view.setBlockQC(5, suit.view.allNode[0])
	getBlockList := &protocols.GetQCBlockList{
		BlockHash:   suit.view.firstProposer().state.HighestQCBlock().Hash(),
		BlockNumber: suit.view.firstProposer().state.HighestQCBlock().NumberU64() + 1,
	}
	if err := suit.view.secondProposer().OnGetQCBlockList("", getBlockList); err != nil {
		suit.T().Fatal(err.Error())
	}
	select {
	case <-suit.msgCh:
		suit.T().Fatal("fail")
	case <-time.After(time.Millisecond * 10):
	}
}

// 正常的
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
	suit.view.firstProposer().OnGetPrepareVote("", getPrepareVote)
	select {
	case m := <-suit.msgCh:
		if msg, ok := m.Message().(*protocols.PrepareVotes); ok {
			suit.Equal(3, len(msg.Votes))
			// sort.Slice(votes, func(i, j int) bool {
			// 	return votes[i].ValidatorIndex < votes[j].ValidatorIndex
			// })
			// sort.Slice(msg.Votes, func(i, j int) bool {
			// 	return msg.Votes[i].ValidatorIndex < msg.Votes[j].ValidatorIndex
			// })
			// suit.EqualValues(votes[:3], msg.Votes)
		}
	case <-time.After(time.Millisecond * 10):
		suit.T().Fatal("timeout")
	}
}

// 正常的
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

// 重复的
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
