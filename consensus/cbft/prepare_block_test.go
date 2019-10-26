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
	"github.com/stretchr/testify/suite"
)

func TestPrepareBlockSuite(t *testing.T) {
	suite.Run(t, new(PrepareBlockTestSuite))
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
	suit.blockOne = NewBlockWithSign(suit.view.genesisBlock.Hash(), 1, suit.view.allNode[0])
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
	block3 := NewBlockWithSign(block2.Hash(), 3, suit.view.allNode[0])
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
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
	fmt.Println(common.Bytes2Hex(block2.Extra()))
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
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[0])
	_, qc := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(), suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.waitVote()
	suit.Equal(1, suit.view.firstProposer().state.PrepareVoteLenByIndex(0))
}

// 第一个块，携带错误的qc
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithNumberIsOne() {
	block1 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 1, suit.view.allNode[0])
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 4)
	errQC := mockErrBlockQC(notConsensusNodes, block1, 0, nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 0,
		suit.view.firstProposerIndex(), block1, errQC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("fail")
	}
}

// 非第一个块,携带错误的qc
func (suit *PrepareBlockTestSuite) TestPrepareBlockWithBlockIndexNotIsZero() {
	suit.insertOneBlock()
	block1 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 2, suit.view.allNode[0])
	notConsensusNodes := mockNotConsensusNode(false, suit.view.nodeParams, 4)
	errQC := mockErrBlockQC(notConsensusNodes, block1, 0, nil)
	prepareBlock := mockPrepareBlock(suit.view.firstProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber, 1,
		suit.view.firstProposerIndex(), block1, errQC.BlockQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("fail")
	}
}

// 最后一个块确认的第一个块,不携带prepareQC的prepareBlock消息
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockOneWithLastBlockQCNotQC() {
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
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
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
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
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(block11.Hash(), 12, suit.view.allNode[1])
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
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	block11QC := mockBlockQC(suit.view.allNode, block11, 0, oldQC)
	suit.view.firstProposer().insertQCBlock(block11, block11QC.BlockQC)
	block12 := NewBlockWithSign(suit.blockOne.Hash(), 12, suit.view.allNode[1])
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
	suit.Equal(1, suit.view.firstProposer().state.PrepareVoteLenByIndex(1))
}

// 上一个块未prepareQC的prepareBlock
// 校验通过，pengding中对于区块存在票数
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithParentNotQC() {
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
	_, oldQC := suit.view.firstProposer().blockTree.FindBlockAndQC(suit.view.firstProposer().state.HighestQCBlock().Hash(),
		suit.view.firstProposer().state.HighestQCBlock().NumberU64())
	prepareBlock11 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block11, oldQC, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock11); err != nil {
		suit.T().Fatal("FAIL")
	}
	suit.waitVote()
	block12 := NewBlockWithSign(block11.Hash(), 12, suit.view.allNode[1])
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
	suit.view.setBlockQC(10, suit.view.allNode[0])
	block11 := NewBlockWithSign(suit.view.firstProposer().state.HighestQCBlock().Hash(), 11, suit.view.allNode[1])
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 1, suit.view.allNode[0])
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
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// 区块哈希不连续块高连续的prepareBlock
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithBlockHashDiscontinuous() {
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
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock); err == nil {
		suit.T().Fatal("FAIL")
	} else {
		fmt.Println(err.Error())
	}
}

// BlockIndex与实际区块index不匹配的prepareBlock
// 校验不通过
func (suit *PrepareBlockTestSuite) TestPrepareBlockNotOneWithBlockIndexErr() {
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 3, suit.view.allNode[0])
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
	block2 := NewBlockWithSign(suit.blockOne.Hash(), 2, suit.view.allNode[0])
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
