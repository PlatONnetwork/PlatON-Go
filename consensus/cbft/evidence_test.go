package cbft

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/suite"
)

type EvidenceTestSuite struct {
	suite.Suite
	view          *testView
	blockOne      *types.Block
	blockOneQC    *protocols.BlockQuorumCert
	oldViewNumber uint64
	epoch         uint64
}

func (suit *EvidenceTestSuite) insertOneBlock() {
	for _, cbft := range suit.view.allCbft {
		insertBlock(cbft, suit.blockOne, suit.blockOneQC.BlockQC)
	}
}

func (suit *EvidenceTestSuite) createEvPool(paths []string) {
	if len(paths) != len(suit.view.allNode) {
		panic("paths len err")
	}
	for i, path := range paths {
		pool, _ := evidence.NewBaseEvidencePool(path)
		suit.view.allCbft[i].evPool = pool
	}

}

func (suit *EvidenceTestSuite) SetupTest() {
	suit.view = newTestView(false, testNodeNumber)
	suit.blockOne = NewBlock(suit.view.genesisBlock.Hash(), 1)
	suit.blockOneQC = mockBlockQC(suit.view.allNode, suit.blockOne, 0, nil)
	suit.oldViewNumber = suit.view.firstCbft.state.ViewNumber()
	suit.epoch = suit.view.firstCbft.state.Epoch()
}

// 双viewChange
func (suit *EvidenceTestSuite) TestViewChangeDup() {
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
	}
	ev := suit.view.firstProposer().Evidences()
	var duplicateViewChange *evidence.EvidenceData
	if err := json.Unmarshal([]byte(ev), duplicateViewChange); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.EqualValues(duplicateViewChange.DC[0].ViewA, viewChange1)
	suit.EqualValues(duplicateViewChange.DC[0].ViewB, viewChange2)
}

// 双出
func (suit *EvidenceTestSuite) TestPrepareBlockDup() {
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
	prepareBlock2 := mockPrepareBlock(suit.view.secondProposerBlsKey(), suit.view.Epoch(), suit.oldViewNumber+1, 0,
		suit.view.secondProposerIndex(), block2, qc, nil)
	if err := suit.view.firstProposer().OnPrepareBlock(suit.view.secondProposer().NodeID().String(), prepareBlock2); err == nil {
		suit.T().Error("FAIL")
	}
	ev := suit.view.firstProposer().Evidences()
	var duplicatePrepareBlock *evidence.EvidenceData
	if err := json.Unmarshal([]byte(ev), duplicatePrepareBlock); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.EqualValues(duplicatePrepareBlock.DP[0].PrepareA, prepareBlock1)
	suit.EqualValues(duplicatePrepareBlock.DP[0].PrepareB, prepareBlock2)
}

// 双签
func (suit *EvidenceTestSuite) TestPrepareVoteDup() {
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
	}
	ev := suit.view.firstProposer().Evidences()
	var duplicatePrepareVote *evidence.EvidenceData
	if err := json.Unmarshal([]byte(ev), duplicatePrepareVote); err != nil {
		suit.T().Fatal(err.Error())
	}
	suit.EqualValues(duplicatePrepareVote.DV[0].VoteA, prepareVote1)
	suit.EqualValues(duplicatePrepareVote.DV[0].VoteB, prepareVote2)
}
