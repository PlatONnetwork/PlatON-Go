package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"math/big"
	"regexp"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/suite"
)

func TestViewChangeSuite(t *testing.T) {
	suite.Run(t, new(ViewChangeTestSuite))
}
func TestPrepareBlockSuite(t *testing.T) {
	suite.Run(t, new(PrepareBlockTestSuite))
}
func TestPrepareVoteSuite(t *testing.T) {
	suite.Run(t, new(PrepareVoteTestSuite))
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
	suit.blockOne = NewBlock(suit.view.genesisBlock.Hash(), 1)
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstCbft.state.ViewNumber()
}

// 发起viewChange
// 校验本地viewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeBuild() {
	time.Sleep((testPeriod + 200) * time.Millisecond)
	suit.Equal(1, suit.view.secondProposer().state.ViewChangeLen())
}

// 发起viewChange
// 非共识节点不会发起viewChange,校验本地viewChangeLen=0
func (suit *ViewChangeTestSuite) TestViewChangeBuildWithNotConsensus() {
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 1)
	time.Sleep((testPeriod + 200) * time.Millisecond)
	suit.Equal(0, notConsensusNodes[0].engine.state.ViewChangeLen())
}

// viewChange消息基本校验
// 1.缺少签名的viewChange消息
// 2.blockNumber与blockHash不匹配的viewChange消息
// 3.签名与验证人不一致的viewChange消息
// 4.提议人非共识节点的viewChange消息
// 5.prepareQC伪造的viewChange消息
// 6.blockNumber不为零的prepareQC为空的viewChange消息
// 7.携带的prepareQC不满足N-f
// 8.epoch大于本地的
// 9.epoch小于本地的
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

// Block与HighestQCBlock本地一致的viewChange消息
// 校验通过，ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeCheckCorrect() {
	suit.insertOneBlock()
	viewChange := mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
		suit.blockOne.Hash(), suit.blockOne.NumberU64(), suit.view.secondProposerIndex(), suit.blockOneQC.BlockQC)
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(1, suit.view.firstProposer().state.ViewChangeLen())
}

// blockNumber为零的prepareQC为空的viewChange消息
// 校验通过，ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeCheckZero() {
	viewChange := mockViewChange(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.view.secondProposer().state.ViewNumber(),
		suit.view.genesisBlock.Hash(), suit.view.genesisBlock.NumberU64(), suit.view.secondProposerIndex(), nil)
	if err := suit.view.firstProposer().OnViewChange(suit.view.secondProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.Equal(1, suit.view.firstProposer().state.ViewChangeLen())
}

// Block领先本地HighestQCBlock的viewChange消息
// 校验通过，ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeLeadHighestQCBlock() {
	block2 := NewBlock(suit.blockOne.Hash(), 2)
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

// Block落后本地HighestQCBlock的viewChange消息
// 校验通过，ViewChangeLen=1
func (suit *ViewChangeTestSuite) TestViewChangeBehindHighestQCBlock() {
	block2 := NewBlock(suit.blockOne.Hash(), 2)
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

// viewNumber小于当前viewNumber
// 校验不通过，ViewChangeLen=0
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

// viewNumber大于当前viewNumber
// 校验不通过，ViewChangeLen=0
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

// 收到已经处理过的viewChange消息
// 校验通过，ViewChangeLen不变
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

// 同一人，基于不同区块的viewChange消息
// 校验不通过，返回双viewChange的错误
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

// 非共识节点收到viewChange
// 校验通过
func (suit *ViewChangeTestSuite) TestViewChangeNotConsensus() {
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 1)
	suit.insertOneBlock()
	viewChange := mockViewChange(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.view.firstProposer().state.ViewNumber(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), suit.view.firstProposerIndex(), suit.blockOneQC.BlockQC)
	if err := notConsensusNodes[0].engine.OnViewChange(suit.view.firstProposer().NodeID().String(), viewChange); err != nil {
		suit.T().Error(err.Error())
	}
}

type PrepareBlockTestSuite struct {
	suite.Suite
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
}

func (suit *PrepareBlockTestSuite) SetupTest() {
	suit.view = newTestView(false, testNodeNumber)
	suit.blockOne = NewBlock(suit.view.genesisBlock.Hash(), 1)
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstProposer().state.ViewNumber()
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

// prepareBlock消息基本校验
// 1.缺少签名的prepareBlock消息
// 2.签名与提议人索引不一致的prepareBlock消息
// 3.出块节点不是当前提议节点的prepareBlock消息
// 4.提议人索引与提议人不匹配的prepareBlock消息
// 5.提议人非共识节点的prepareBlock消息
// 6.epoch 太大
// 7.epoch 太小
func (suit *PrepareBlockTestSuite) TestCheckErrPrepareBlock() {
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 1)
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
		if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), testcase.data); err == nil {
			suit.T().Errorf("case %s is failed", testcase.name)
			suit.view.secondProposer().state.ResetView(suit.view.Epoch(), suit.oldViewNumber)
			// suit.T().Error(err.Error())
		}
	}
}

// 携带prepareQC和viewChangeQC的prepareBlock消息，本节点未完成viewChangeQC
// 校验通过，blockIndex为0的票数为1，viewNumber+1
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithViewChangePrepareQCAndViewChangeQC() {
	suit.insertOneBlock()
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	viewQC := mockViewQC(suit.blockOne, suit.view.allNode, suit.blockOneQC.BlockQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	suit.Equal(1, suit.view.firstProposer().state.PrepareVoteLenByIndex(0))
	suit.Equal(suit.oldViewNumber+1, suit.view.firstProposer().state.ViewNumber())
}

// 携带prepareQC和viewChangeQC的prepareBlock消息,接收消息的节点已经收到viewChangeQC
// 校验通过，blockIndex为0的票数为1，viewNumber不变
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithViewChangePrepareQCAndViewChangeQCHadViewChangQC() {
	suit.insertOneBlock()
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	viewQC := mockViewQC(suit.blockOne, suit.view.allNode, suit.blockOneQC.BlockQC)
	fmt.Println(viewQC.String())
	suit.view.firstProposer().changeView(suit.view.Epoch(), suit.oldViewNumber+1, suit.blockOne, suit.blockOneQC.BlockQC, viewQC)
	suit.view.firstProposer().state.ResetView(suit.view.Epoch(), suit.oldViewNumber+1)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	fmt.Println(prepareBlock.BlockNum())
	fmt.Println(prepareBlock.Block.ParentHash().String())
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	suit.Equal(1, suit.view.firstProposer().state.PrepareVoteLenByIndex(0))
	suit.Equal(suit.oldViewNumber+1, suit.view.firstProposer().state.ViewNumber())
}

// viewChangeQC的第一个块,接收消息的节点未收到viewChangeQC
// 1.不携带viewChangeQC的prepareBlock消息
// 2.不携带prepareQC的prepareBlock消息
// 3.不携带prepareQC和viewChangeQC的prepareBlock消息
// 4.blockIndex为1的prepareBlock消息
// 5.携带prepareQC，携带不满足2f+1的viewChangeQC
// 6.epoch 太大
// 7.epoch 太小
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithViewChangeErrFirstBlock() {
	suit.insertOneBlock()
	block2 := NewBlock(suit.blockOne.Hash(), 2)
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
		if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), testcase.data); err == nil {
			suit.T().Errorf("CASE:%s is failed", testcase.name)
			suit.view.firstProposer().state.ResetView(oldEpoch, suit.oldViewNumber+1)
			// suit.T().Error(err.Error())
		} else {
			fmt.Printf("case:%s-err:%s\n", testcase.name, err.Error())
		}
	}
}

// viewChangeQC的第一个块,接收消息的节点已经收到viewChangeQC
// 1.不携带viewChangeQC的prepareBlock消息
// 2.不携带prepareQC的prepareBlock消息
// 3.不携带prepareQC和viewChangeQC的prepareBlock消息
// 4.blockIndex为1的prepareBlock消息
// 5.携带prepareQC，携带不满足2f+1的viewChangeQC
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithViewChangeErrFirstBlockHadViewChangQC() {
	suit.insertOneBlock()
	block2 := NewBlock(suit.blockOne.Hash(), 2)
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
		if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), testcase.data); err == nil {
			suit.T().Errorf("CASE:%s is failed", testcase.name)
			suit.view.firstProposer().state.ResetView(oldEpoch, suit.oldViewNumber+1)
			// suit.T().Error(err.Error())
		} else {
			fmt.Printf("case:%s-err:%s\n", testcase.name, err.Error())
		}
	}
}

// viewChangeQC的第一个块,出块节点HighestQCBlock领先本地HighestQCBlock的prepareBlock消息,收到块的节点未完成viewChangeQC
// 由于落后，会触发同步,校验无法通过，返回错误值为 viewNumber higher then local(local:0, msg:1)
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithViewChangeFirstBlockTooHigh() {
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	block3 := NewBlock(block2.Hash(), 3)
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	suit.view.secondProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewQC := mockViewQC(block2, suit.view.allNode, block2QC.BlockQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
		suit.oldViewNumber+1, 0, suit.view.secondProposerIndex(), block3, block2QC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		suit.EqualValues("viewNumber higher than local(local:0, msg:1)", err.Error())
	}
}

// viewChangeQC的第一个块,出块节点HighestQCBlock落后本地HighestQCBlock的prepareBlock消息,收到块的节点未完成viewChangeQC
// 校验无法通过，返回错误值为 viewNumber higher then local(local:0, msg:1)
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithViewChangeFirstBlockTooLow() {
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	suit.view.firstProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewQC := mockViewQC(block2, suit.view.allNode, block2QC.BlockQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		suit.EqualValues("viewNumber higher than local(local:0, msg:1)", err.Error())
	}
}

// viewChangeQC的第一个块,出块节点HighestQCBlock落后本地HighestQCBlock的prepareBlock消息,收到块的节点已经完成viewChangeQC,此时没有基于viewQC.MaxBlock
// 校验无法通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithViewChangeFirstBlockTooLowHad() {
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	suit.view.firstProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewQC := mockViewQC(block2, suit.view.allNode, block2QC.BlockQC)
	suit.view.firstProposer().state.ResetView(suit.view.Epoch(), suit.oldViewNumber+1)
	suit.view.firstProposer().state.SetLastViewChangeQC(viewQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	}
}

// viewChangeQC的第一个块,没有基于viewQC.MaxBlock,收到块的节点未完成viewChangeQC
// 校验无法通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithViewChangeFirstBlockNotWithMaxBlock() {
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	block2QC := mockBlockQC(suit.view.allNode, block2, 1, suit.blockOneQC.BlockQC)
	suit.insertOneBlock()
	// suit.view.firstProposer().insertQCBlock(block2, block2QC.BlockQC)
	viewQC := mockViewQC(block2, suit.view.allNode, block2QC.BlockQC)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, suit.blockOneQC.BlockQC, viewQC)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 最后一个块确认的第一个块,blockNumber相同的hash不同prepareBlock消息
// 第一个块通过，第二个校验不通过，对应的PrepareVoteLen=1，返回双出error
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithDifHash() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	suit.view.setBlockQC(10)
	block1 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	header := &types.Header{
		Number:      big.NewInt(int64(11)),
		ParentHash:  suit.view.firstProposer().state.HighestQCBlock().Hash(),
		Time:        big.NewInt(time.Now().UnixNano() + 100),
		Extra:       make([]byte, 77),
		ReceiptHash: common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Root:        common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Coinbase:    common.Address{},
		GasLimit:    100000000001,
	}
	block2 := types.NewBlockWithHeader(header)
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock1 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block1, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	prepareBlock2 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, qc, nil)
	fmt.Println(block1.Hash().String(), block2.Hash().String())
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock2); err == nil {
		suit.T().Error("FAIL")
	} else {
		reg := regexp.MustCompile(`DuplicatePrepareBlockEvidence`)
		if len(reg.FindAllString(err.Error(), -1)) == 0 {
			suit.T().Fatal(err.Error())
		}
	}
	suit.waitVote()
	suit.Equal(1, suit.view.firstProposer().state.PrepareVoteLenByIndex(0))
}

// 最后一个块确认的第一个块,携带prepareQC的prepareBlock消息
// 校验通过，PrepareVoteLenByIndex(0)=1
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithLastBlockQC() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(), suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	suit.Equal(1, suit.view.firstProposer().state.PrepareVoteLenByIndex(0))
}

// 最后一个块确认的第一个块,不携带prepareQC的prepareBlock消息
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithLastBlockQCNotQC() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 最后一个块确认的第一个块,blockIndex为1的prepareBlock消息
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithLastBlockQCBlockIndexIsOne() {
	// oldEpoch := suit.view.Epoch()
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block11, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 最后一个块确认的第一个块,出块节点HighestQCBlock领先本地HighestQCBlock的prepareBlock消息
// 校验不通过，触发同步
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithLastBlockQCLead() {
	otherNode := suit.view.thirdProposer()
	suit.view.setBlockQC(9)
	block10 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 10)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block10QC := mockBlockQC(suit.view.allNode, block10, 0,
		oldQC)
	block11 := NewBlock(block10.Hash(), 11)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(),
		suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, block10QC.BlockQC, nil)
	if err := otherNode.OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// viewNumber小于当前viewNumber
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithViewNumberTooLow() {
	// oldEpoch := suit.view.Epoch()
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	block12 := NewBlock(block11.Hash(), 12)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// viewNumber大于当前viewNumber
// 校验不通过，触发同步
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithViewNumberTooHigh() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	block12 := NewBlock(suit.blockOne.Hash(), 12)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+2, 1,
		suit.view.secondProposerIndex(), block12, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 上一个块prepareQC的prepareBlock
// 校验通过，票数+1
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithParentQC() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlock(block11.Hash(), 12)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	suit.Equal(1, suit.view.firstProposer().state.PrepareVoteLenByIndex(1))
}

// 上一个块未prepareQC的prepareBlock
// 校验通过，pengding中对于区块存在票数
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithParentNotQC() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock11 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, oldQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock11); err != nil {
		suit.T().Fatal("FAIL")
	}
	suit.waitVote()
	block12 := NewBlock(block11.Hash(), 12)
	prepareBlock12 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock12); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	suit.Equal(uint64(12), suit.view.firstProposer().state.PendingPrepareVote().Votes[0].BlockNum())

}

// 块数超出一轮限制
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithAmountTooMany() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 10,
		suit.view.secondProposerIndex(), block11, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 相同块高，不同哈希的prepareBlock
// 第二个校验不通过，双出error
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithBlockNumberRepeat() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal(err.Error())
	}
	block2 := NewBlock(suit.blockOne.Hash(), 1)
	prepareBlock2 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), prepareBlock2); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		reg := regexp.MustCompile(`DuplicatePrepareBlockEvidence`)
		if len(reg.FindAllString(err.Error(), -1)) == 0 {
			suit.T().Fatal(err.Error())
		}
	}
}

// 块高不连续的prepareBlock,区块hash连续
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithBlockNumberDiscontinuous() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlock(block11.Hash(), 13)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 区块哈希不连续块高连续的prepareBlock
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithBlockHashDiscontinuous() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlock(suit.view.genesisBlock.Hash(), 12)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// BlockIndex与实际区块index不匹配的prepareBlock
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithBlockIndexErr() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlock(block11.Hash(), 12)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 2,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}

}

// 本地区块存在相同 BlockIndex 区块，但BlockHash，BlockNumber 不相等
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithBlockIndexRepeat() {
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal(err.Error())
	}
	block2 := NewBlock(suit.blockOne.Hash(), 3)
	prepareBlock2 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock2); err == nil {
		suit.T().Fatal("FAIL")
	}
}

// 数据正确的重复的prepareBlock消息
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockDup() {
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal(err.Error())
	}
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 本地不存在BlockIndex对应的前一个索引区块的prepareBlock消息
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockBlockIndexTooHigh() {
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 数据正确的超时的prepareBlock消息
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithTimeout() {
	time.Sleep((testPeriod + 200) * time.Millisecond)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 非共识节点收到合法prepareBlock消息
// 校验通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithNotConsensus() {
	notConsensus := mockNotConsensusNode(false, suit.view.nodeParams, 1)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := notConsensus[0].engine.OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal(err.Error())
	}
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
	suit.blockOne = NewBlock(suit.view.genesisBlock.Hash(), 1)
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

// 构造prepareVote消息
// 收到区块，生成对应的投票
// 校验块高与块哈希一致
func (suit *PrepareVoteTestSuite) TestBuildPrepareVote() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlock(block11.Hash(), 12)
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, block11QC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	vote := suit.view.firstProposer().state.AllPrepareVoteByIndex(1)[0]
	suit.Equal(uint64(12), vote.BlockNum())
	suit.Equal(block12.Hash().String(), vote.BlockHash.String())
	suit.Equal(prepareBlock.BlockIndex, vote.BlockIndex)

}

// prepareVote消息基本校验
// 1.没有签名的
// 2.签名与validatorIndex不一致的
// 3.签名不是验证节点的
// 4.epoch 太大
// 5.epoch 太小
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

// 父区块为零的不携带ParentQC的prepareVote消息
// 校验通过
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

// 父区块非零的不携带ParentQC的prepareVote消息
// 校验不通过
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentIsNotZeroButNotParentQC() {
	suit.insertOneBlock()
	block2 := NewBlock(suit.blockOne.Hash(), 2)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block2, suit.blockOneQC.BlockQC, nil)
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

// 未收到prepareBlock先收到的prepareVote消息
// 校验不通过
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithNotPrepareBlock() {
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 区块超出一轮出块数限制的prepareVote消息
// 校验不通过
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithExceedLimit() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 21)
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

// 收到重复的prepareVote消息
// 合法的prepareVote消息
// 第一次校验通过，第二次提示投票已经存在，票总数为1
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

// 双签
// 返回双签错误
func (suit *PrepareVoteTestSuite) TestPrepareVoteDu() {
	paths := createPaths(len(suit.view.allCbft))
	defer removePaths(paths)
	suit.createEvPool(paths)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	suit.view.secondProposer().state.AddPrepareBlock(prepareBlock)
	block1 := NewBlock(suit.view.genesisBlock.Hash(), 1)
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

// viewNumber小于当前viewNumber
// 校验不通过
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

// viewNumber大于当前viewNumber
// 校验不通过，触发同步
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithViewNumberTooHigh() {
	prepareVote := mockPrepareVote(suit.view.firstProposerBlsKey(), suit.epoch, suit.oldViewNumber+1, 0, suit.view.firstProposerIndex(), suit.blockOne.Hash(),
		suit.blockOne.NumberU64(), nil)
	if err := suit.view.secondProposer().OnPrepareVote(suit.view.firstProposer().NodeID().String(), prepareVote); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		suit.Equal(err.Error(), "viewNumber higher than local(local:0, msg:1)")
	}
}

// Vote的父区块在本节点未达成prepareQC
// 校验通过
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentIsNotParentQC() {
	qc := mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.firstProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal(err.Error())
	}
	block2 := NewBlock(suit.blockOne.Hash(), 2)
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

// Vote的父区块在发送节点未达成prepareQC（不合法）
// 校验不通过
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentErrParentQC() {
	qc := mockBlockQC(suit.view.allNode[0:1], suit.blockOne, 0, nil)
	prepareBlock1 := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), suit.blockOne, nil, nil)
	if err := suit.view.secondProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock1); err != nil {
		suit.T().Fatal("FAIL")
	}
	block2 := NewBlock(suit.blockOne.Hash(), 2)
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

// 达成prepareQC时，存在子区块prepareVote
// 校验通过，发送子区块的投票
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentQCHasChild() {
	suit.view.setBlockQC(10)
	block11 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11)
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	// suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	insertBlock(suit.view.firstProposer(), block11, block11QC.BlockQC)
	block12 := NewBlock(block11.Hash(), 12)
	prepareBlock12 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 1,
		suit.view.secondProposerIndex(), block12, nil, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock12); err != nil {
		suit.T().Fatal(err.Error())
	}
	block12QC := mockBlockQC(suit.view.allNode, block12, 1, block11QC.BlockQC)
	block13 := NewBlock(block12.Hash(), 13)
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

// 达成prepareQC时，不存在子区块prepareVote
// 校验commit和lock
func (suit *PrepareVoteTestSuite) TestPrepareVoteWithParentQCNotHasChild() {
	suit.view.setBlockQC(5)
	block6 := NewBlock(suit.view.firstProposer().state.HighestQCBlock().Hash(), 6)
	qc := mockBlockQC(suit.view.allNode, block6, 0, nil)
	suit.view.secondProposer().insertQCBlock(block6, qc.BlockQC)
	commitNumber, _ := suit.view.secondProposer().HighestCommitBlockBn()
	lockNumber, _ := suit.view.secondProposer().HighestLockBlockBn()
	suit.Equal(uint64(4), commitNumber)
	suit.Equal(uint64(5), lockNumber)
}

// 数据合法的超时的prepareVote消息
// 校验不通过
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

// 数据刚好满足2f+1的prepareQC消息
// 校验通过
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

// 非共识节点收到prepareVote
// 校验通过
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
