package cbft

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/stretchr/testify/assert"
)

func TestNewViewChange(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)

	priA := validators.neighbors[0]
	//addrA := crypto.PubkeyToAddress(*priA.publicKey)

	gen := backend.chain.Genesis()

	var blocks []*types.Block
	blocks = append(blocks, gen)
	for i := uint64(1); i < 10; i++ {
		blocks = append(blocks, createBlock(priA.privateKey, blocks[i-1].Hash(), blocks[i-1].NumberU64()+1))
		//t.Log(blocks[i].NumberU64(), blocks[i].Hash().TerminalString(), blocks[i].ParentHash().TerminalString())
	}

	node := nodeIndexNow(validators, engine.startTimeOfEpoch)
	viewChange := makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address, nil)

	//create view base of genesis
	err := engine.OnViewChange(node.nodeID, viewChange)
	assert.Nil(t, err)

	err = engine.OnViewChange(node.nodeID, viewChange)
	assert.Equal(t, errDuplicationConsensusMsg, err)

	//wait switch next validator
	time.Sleep(time.Second * time.Duration(engine.config.Duration+1))

	//last validator's timestamp doesn't match current timestamp
	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address, nil)

	err = engine.OnViewChange(node.nodeID, viewChange)
	assert.Equal(t, errRecvViewTimeout, err)

	//create new viewchange base on current validator
	node = nodeIndexNow(validators, engine.startTimeOfEpoch)
	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address, nil)

	//newest viewchange is satisfied
	err = engine.OnViewChange(node.nodeID, viewChange)
	assert.Nil(t, err)

	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6-engine.config.Duration*1e3), 0, gen.Hash(), uint32(node.index), node.address, nil)

	//newest viewchange is satisfied
	err = engine.OnViewChange(node.nodeID, viewChange)
	assert.Equal(t, errTimestamp, err)

	//set localHighestPrepareVoteNum
	engine.localHighestPrepareVoteNum = backend.chain.CurrentBlock().NumberU64() + 1
	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6)+1, 0, gen.Hash(), uint32(node.index), node.address, nil)
	err = engine.OnViewChange(node.nodeID, viewChange)
	assert.Equal(t, errViewChangeBaseTooLow, err)

	//reset localHighestPrepareVoteNum
	engine.localHighestPrepareVoteNum = backend.chain.CurrentBlock().NumberU64()
	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6)+1, 0, common.BytesToHash(Rand32Bytes(32)), uint32(node.index), node.address, nil)
	err = engine.OnViewChange(node.nodeID, viewChange)
	assert.Equal(t, errViewChangeForked, err)

	//set higher baseblock
	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6)+1, 1, common.BytesToHash(Rand32Bytes(32)), uint32(node.index), node.address, nil)
	err = engine.OnViewChange(node.nodeID, viewChange)
	assert.Equal(t, errNotFoundViewBlock, err)
	time.Sleep(time.Second * time.Duration(engine.config.Period*3))

	assert.Nil(t, engine.viewChange)
	assert.Empty(t, engine.viewChangeVotes)
}

func TestCbft_OnSendViewChange(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 4)

	engine.OnSendViewChange()

	assert.NotNil(t, engine.viewChange)
	time.Sleep(time.Second * time.Duration(engine.config.Period*3))
	assert.Nil(t, engine.viewChange)
}

func TestCbft_ShouldSeal(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, validators := randomCBFT(path, 4)

	seal, err := engine.ShouldSeal(100)
	assert.False(t, seal)
	assert.Equal(t, errTwoThirdViewchangeVotes, err)

	time.Sleep(time.Second * time.Duration(engine.config.Period) * 3)
	assert.Nil(t, engine.viewChange)

	seal, err = engine.ShouldSeal(100)
	assert.False(t, seal)
	assert.Equal(t, errTwoThirdViewchangeVotes, err)

	for _, v := range validators.neighbors {
		engine.viewChangeVotes[v.address] = nil
	}

	seal, err = engine.ShouldSeal(100)
	assert.True(t, seal)
	assert.Nil(t, err)
}

func TestCbft_GetBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)

	priA := validators.neighbors[0]
	gen := backend.chain.Genesis()

	var blocks []*types.Block
	blocks = append(blocks, gen)
	for i := uint64(1); i < 10; i++ {
		block := createBlock(priA.privateKey, blocks[i-1].Hash(), blocks[i-1].NumberU64()+1)
		blocks = append(blocks, block)
		ext := NewBlockExtByPeer(block, block.NumberU64(), len(validators.Nodes()))
		engine.blockExtMap.Add(block.Hash(), block.NumberU64(), ext)
		//t.Log(blocks[i].NumberU64(), blocks[i].Hash().TerminalString(), blocks[i].ParentHash().TerminalString())
	}

	res := engine.GetBlock(blocks[3].Hash(), blocks[3].NumberU64())
	assert.Equal(t, blocks[3].Hash(), res.Hash())
	assert.Equal(t, blocks[3].NumberU64(), res.NumberU64())

	assert.True(t, engine.HasBlock(blocks[0].Hash(), blocks[0].NumberU64()))
	assert.False(t, engine.HasBlock(common.BytesToHash(Rand32Bytes(32)), res.NumberU64()))
	assert.NotNil(t, engine.GetBlockByHash(blocks[3].Hash()))
}

func TestCbft_IsSignedBySelf(t *testing.T) {

	path := path()
	defer os.RemoveAll(path)

	engine, _, validators := randomCBFT(path, 4)
	sealHash := common.BytesToHash(Rand32Bytes(32))

	sign, _ := crypto.Sign(sealHash[:], validators.owner.privateKey)

	assert.True(t, engine.IsSignedBySelf(sealHash, sign))

	assert.False(t, engine.IsSignedBySelf(sealHash, sign[:20]))
}

func TestCbft_Seal(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)

	node := validators.owner
	genesis := backend.chain.Genesis()
	block := createBlock(node.privateKey, genesis.Hash(), 1)

	viewChange := makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, genesis.Hash(), uint32(node.index), node.address, nil)

	engine.viewChange = viewChange

	sealResultCh := make(chan *types.Block, 1)
	stopCh := make(chan struct{}, 1)
	assert.Nil(t, engine.Seal(backend.chain, block, sealResultCh, stopCh))

}

func TestCbft_NextBaseBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, _ := randomCBFT(path, 4)

	block := engine.NextBaseBlock()
	genesis := backend.chain.Genesis()
	assert.Equal(t, block.Hash(), genesis.Hash())
	assert.Equal(t, block.NumberU64(), genesis.NumberU64())
}

func TestCbft_OnPrepareBlockHash(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, validators := randomCBFT(path, 4)

	assert.Nil(t, engine.OnPrepareBlockHash(validators.neighbors[0].nodeID,
		&prepareBlockHash{
			Hash:   common.BytesToHash(Rand32Bytes(32)),
			Number: 20,
		}))
}

//create insert chain
func TestCbft_InsertChain(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)

	// create view by validator 1
	view := makeViewChange(v.validator(1).privateKey, uint64(engine.config.Duration*1000*1+100), 0, engine.blockChain.Genesis().Hash(), 1, v.validator(1).address, nil)

	blocks := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view, 2)

	// create view by validator 2
	view2 := makeViewChange(v.validator(2).privateKey, uint64(engine.config.Duration*1000*2+100), 2, blocks[1].block.Hash(), 2, v.validator(2).address, blocks[1].prepareVotes.Votes())

	blocks2 := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view2, 2)

	states := make([]chan error, 0)
	for _, b := range blocks {
		syncState := make(chan error, 1)
		engine.InsertChain(b.block, syncState)
		states = append(states, syncState)
	}

	for _, b := range blocks2 {
		syncState := make(chan error, 1)
		engine.InsertChain(b.block, syncState)
		states = append(states, syncState)
	}

	for i, s := range states {
		if err := <-s; err != nil {
			t.Error(fmt.Sprintf("%d th block error", i))
		}
	}
}

func TestCbft_OnGetHighestPrepareBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)

	mockHandler := NewMockHandler()
	engine.handler = mockHandler

	//sendQueue <- &MsgPackage{
	//	v.validator(2).nodeID.String(),
	//	&getHighestPrepareBlock{2},
	//	1,
	//}

	view := makeViewChange(v.validator(1).privateKey, uint64(engine.config.Duration*1000*1+100), 0, engine.blockChain.Genesis().Hash(), 1, v.validator(1).address, nil)

	blocks := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view, 2)

	for _, block := range blocks {
		engine.blockExtMap.Add(block.block.Hash(), block.number, block)
	}

	assert.Nil(t, engine.OnGetHighestPrepareBlock(v.validator(2).nodeID, &getHighestPrepareBlock{2}))

	assert.Len(t, mockHandler.sendQueue, 1)
}

func TestCbft_OnHighestPrepareBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)

	badBlocks := createEmptyBlocks(v.validator(1).privateKey, engine.blockChain.Genesis().Hash(), engine.blockChain.Genesis().NumberU64(), int(maxBlockDist+1))

	assert.Error(t, engine.OnHighestPrepareBlock(v.validator(1).nodeID, &highestPrepareBlock{CommitedBlock: badBlocks}))

	view := makeViewChange(v.validator(1).privateKey, uint64(engine.config.Duration*1000*1+100), 0, engine.blockChain.Genesis().Hash(), 1, v.validator(1).address, nil)

	blocksExt := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view, 9)
	var blocks []*types.Block
	for _, ext := range blocksExt[0:3] {
		blocks = append(blocks, ext.block)
	}
	var prepare []*prepareVotes
	for _, ext := range blocksExt[3:9] {
		prepare = append(prepare, &prepareVotes{
			Hash:   ext.block.Hash(),
			Number: ext.block.NumberU64(),
			Votes:  ext.prepareVotes.Votes(),
		})
	}
	hp := &highestPrepareBlock{
		CommitedBlock: blocks,
		Votes:         prepare,
	}
	engine.OnHighestPrepareBlock(v.validator(1).nodeID, hp)
}

func TestCbft_OnGetPrepareVote2(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)

	mockHandler := NewMockHandler()
	engine.handler = mockHandler

	// create view by validator 1
	view := makeViewChange(v.validator(1).privateKey, uint64(engine.config.Duration*1000*1+100), 0, engine.blockChain.Genesis().Hash(), 1, v.validator(1).address, nil)

	blocks := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view, 2)

	// create view by validator 2
	view2 := makeViewChange(v.validator(2).privateKey, uint64(engine.config.Duration*1000*2+100), 2, blocks[1].block.Hash(), 2, v.validator(2).address, blocks[1].prepareVotes.Votes())

	blocks2 := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view2, 2)

	states := make([]chan error, 0)
	for _, b := range blocks {
		syncState := make(chan error, 1)
		engine.InsertChain(b.block, syncState)
		states = append(states, syncState)
	}

	for _, b := range blocks2 {
		engine.blockExtMap.Add(b.block.Hash(), b.number, b)
	}

	for i, s := range states {
		if err := <-s; err != nil {
			t.Error(fmt.Sprintf("%d th block error", i))
		}
	}

	pv1 := &getPrepareVote{
		Hash:     common.BytesToHash(Rand32Bytes(32)),
		Number:   1000,
		VoteBits: NewBitArray(4),
	}

	assert.Nil(t, engine.OnGetPrepareVote(v.validator(1).nodeID, pv1))

	assert.Len(t, mockHandler.sendQueue, 0)

	pv1 = &getPrepareVote{
		Hash:     blocks[1].block.Hash(),
		Number:   blocks[1].number,
		VoteBits: NewBitArray(4),
	}

	assert.Nil(t, engine.OnGetPrepareVote(v.validator(1).nodeID, pv1))

	assert.Len(t, mockHandler.sendQueue, 1)

	mockHandler.clear()

	pv1 = &getPrepareVote{
		Hash:     blocks[1].block.Hash(),
		Number:   blocks[1].number,
		VoteBits: NewBitArray(4),
	}

	assert.Nil(t, engine.OnGetPrepareVote(v.validator(1).nodeID, pv1))

	assert.Len(t, mockHandler.sendQueue, 1)

}

func TestCbft_OnConfirmedPrepareBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)

	mockHandler := NewMockHandler()
	engine.handler = mockHandler

	// create view by validator 1
	view := makeViewChange(v.validator(1).privateKey, uint64(engine.config.Duration*1000*1+100), 0, engine.blockChain.Genesis().Hash(), 1, v.validator(1).address, nil)

	blocks := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view, 2)

	states := make([]chan error, 0)
	for _, b := range blocks {
		syncState := make(chan error, 1)
		engine.InsertChain(b.block, syncState)
		states = append(states, syncState)
	}

	view2 := makeViewChange(v.validator(2).privateKey, uint64(engine.config.Duration*1000*2+100), 2, blocks[1].block.Hash(), 2, v.validator(2).address, blocks[1].prepareVotes.Votes())

	blocks2 := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view2, 2)

	engine.blockExtMap.Add(blocks2[0].block.Hash(), blocks2[0].number, blocks2[0])

	copyBits := blocks2[1].prepareVotes.voteBits.copy()
	blocks2[1].prepareVotes = NewPrepareVoteSet(blocks2[1].prepareVotes.voteBits.Size())
	engine.blockExtMap.Add(blocks2[1].block.Hash(), blocks2[1].number, blocks2[1])
	c := &confirmedPrepareBlock{
		Hash:     blocks2[1].block.Hash(),
		Number:   blocks2[1].number,
		VoteBits: copyBits,
	}

	assert.Nil(t, engine.OnConfirmedPrepareBlock(v.validator(1).nodeID, c))

	assert.Len(t, mockHandler.sendQueue, 1)
	mockHandler.clear()
	c = &confirmedPrepareBlock{
		Hash:     blocks2[0].block.Hash(),
		Number:   blocks2[0].number,
		VoteBits: blocks2[0].prepareVotes.voteBits,
	}

	assert.Nil(t, engine.OnConfirmedPrepareBlock(v.validator(1).nodeID, c))
	assert.Len(t, mockHandler.sendQueue, 1)

}
func TestCbft_OnViewChangeVote(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)
	timestamp := uint64(common.Millis(time.Now()))

	view := makeViewChange(v.validator(0).privateKey, timestamp, 0, engine.blockChain.Genesis().Hash(), 0, v.validator(0).address, nil)
	engine.viewChange = view

	vote0 := makeViewChangeVote(v.validator(1).privateKey, timestamp, 0, engine.blockChain.Genesis().Hash(), 0, v.validator(0).address, 1, v.validator(1).address)
	assert.Nil(t, engine.OnViewChangeVote(v.validator(1).nodeID, vote0))

	// less 2f+1 votes
	assert.False(t, engine.agreeViewChange())

	// inconsistent block number
	voteErrBlockNum := makeViewChangeVote(v.validator(2).privateKey, timestamp, 1, engine.blockChain.Genesis().Hash(), 0, v.validator(0).address, 2, v.validator(2).address)
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(2).nodeID, voteErrBlockNum))
	assert.False(t, engine.agreeViewChange())

	// inconsistent timestamp
	voteErrTs := makeViewChangeVote(v.validator(1).privateKey, 1000, 0, engine.blockChain.Genesis().Hash(), 0, v.validator(0).address, 1, v.validator(1).address)
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(1).nodeID, voteErrTs))
	assert.False(t, engine.agreeViewChange())

	// inconsistent block number
	voteErrBlockHash := makeViewChangeVote(v.validator(2).privateKey, timestamp, 0, common.Hash{}, 0, v.validator(0).address, 2, v.validator(2).address)
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(2).nodeID, voteErrBlockHash))
	assert.False(t, engine.agreeViewChange())

	// empty parameters
	voteEmpty := &viewChangeVote{}
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(2).nodeID, voteEmpty))
	assert.False(t, engine.agreeViewChange())

	// parameter error
	voteErr := makeViewChangeVote(v.validator(2).privateKey, timestamp, 0, common.Hash{}, 0, common.Address{}, 100, common.Address{})
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(2).nodeID, voteErr))
	assert.False(t, engine.agreeViewChange())

	vote1 := makeViewChangeVote(v.validator(2).privateKey, timestamp, 0, engine.blockChain.Genesis().Hash(), 0, v.validator(0).address, 2, v.validator(2).address)
	assert.Nil(t, engine.OnViewChangeVote(v.validator(2).nodeID, vote1))

	// fulfil 2f+1 votes
	assert.True(t, engine.agreeViewChange())

	vote2 := makeViewChangeVote(v.validator(3).privateKey, timestamp, 0, engine.blockChain.Genesis().Hash(), 0, v.validator(0).address, 3, v.validator(3).address)
	assert.Nil(t, engine.OnViewChangeVote(v.validator(3).nodeID, vote2))

	assert.True(t, engine.agreeViewChange())

	blocksExt := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view, 20)
	for _, ext := range blocksExt {
		ch := make(chan error, 1)
		engine.InsertChain(ext.block, ch)
		<-ch
	}

	timestamp = uint64(common.Millis(time.Now()))
	block := blocksExt[19].block
	engine.clearViewChange()
	view0 := makeViewChange(v.validator(0).privateKey, timestamp, block.NumberU64(), block.Hash(), 0, v.validator(0).address, nil)
	engine.viewChange = view0

	errTs := blocksExt[9].timestamp - 100
	voteErrTs1 := makeViewChangeVote(v.validator(1).privateKey, errTs, block.NumberU64(), block.Hash(), 0, v.validator(0).address, 1, v.validator(1).address)
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(1).nodeID, voteErrTs1))

	voteErrTs2 := makeViewChangeVote(v.validator(2).privateKey, errTs, block.NumberU64(), block.Hash(), 0, v.validator(0).address, 2, v.validator(2).address)
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(2).nodeID, voteErrTs2))

	assert.False(t, engine.agreeViewChange())

	block1 := blocksExt[8].block
	voteErrBlockNum1 := makeViewChangeVote(v.validator(1).privateKey, timestamp, block1.NumberU64(), block.Hash(), 0, v.validator(0).address, 1, v.validator(1).address)
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(1).nodeID, voteErrBlockNum1))

	voteErrBlockNum2 := makeViewChangeVote(v.validator(2).privateKey, timestamp, block1.NumberU64(), block.Hash(), 0, v.validator(0).address, 2, v.validator(2).address)
	assert.NotNil(t, engine.OnViewChangeVote(v.validator(2).nodeID, voteErrBlockNum2))

	assert.False(t, engine.agreeViewChange())

	voteOk1 := makeViewChangeVote(v.validator(1).privateKey, timestamp, block.NumberU64(), block.Hash(), 0, v.validator(0).address, 1, v.validator(1).address)
	assert.Nil(t, engine.OnViewChangeVote(v.validator(1).nodeID, voteOk1))

	voteOk2 := makeViewChangeVote(v.validator(2).privateKey, timestamp, block.NumberU64(), block.Hash(), 0, v.validator(0).address, 2, v.validator(2).address)
	assert.Nil(t, engine.OnViewChangeVote(v.validator(2).nodeID, voteOk2))

	assert.True(t, engine.agreeViewChange())
	assert.True(t, engine.HasTwoThirdsMajorityViewChangeVotes())
}

func TestCBFT_OnPrepareVote(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, backend, validators := randomCBFT(path, 4)
	priA := validators.neighbors[0]
	gen := backend.chain.Genesis()

	// define block
	var blocks []*types.Block
	blocks = append(blocks, gen)
	for i := uint64(1); i < 10; i++ {
		blocks = append(blocks, createBlock(priA.privateKey, blocks[i-1].Hash(), blocks[i-1].NumberU64()+1))
	}
	node := nodeIndexNow(validators, engine.startTimeOfEpoch)
	pvote := makePrepareVote(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address)

	var err error
	testCases := []struct {
		pid     string
		pga     bool
		msgHash common.Hash
	}{
		{pid: "peer id 01", pga: true, msgHash: common.BytesToHash([]byte("Invalid hash"))},
		{pid: "peer id 02", pga: true, msgHash: pvote.MsgHash()},
		{pid: node.nodeID.TerminalString(), pga: true, msgHash: pvote.MsgHash()},
	}

	for _, v := range testCases {
		handler := makeHandler(engine, v.pid, v.msgHash)
		engine.handler = handler
		err = engine.OnPrepareVote(node.nodeID, pvote, v.pga)
		assert.Nil(t, err)
	}
	// -------------------------------------------------------------------

	// ext == nil
	pvote = makePrepareVote(node.privateKey, uint64(time.Now().UnixNano()/1e6), 1111, common.BytesToHash([]byte("Invalid block hash")), uint32(node.index), node.address)
	viewChange := makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address, nil)
	engine.setViewChange(viewChange)
	err = engine.OnPrepareVote(node.nodeID, pvote, false)
	assert.Nil(t, err)

	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address, nil)
	engine.setViewChange(viewChange)
	err = engine.OnPrepareVote(node.nodeID, pvote, false)
	assert.Nil(t, err)

	// verify the sign of validator
	pvote = makePrepareVote(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address)
	pvote.ValidatorAddr = common.BytesToAddress([]byte("fake address"))
	err = engine.OnPrepareVote(node.nodeID, pvote, false)

	assert.NotNil(t, err)

	// whether to accept: cache
	pvote = makePrepareVote(node.privateKey, uint64(time.Now().UnixNano()/1e6), 3, gen.Hash(), uint32(node.index), node.address)
	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 1, gen.Hash(), uint32(node.index), node.address, nil)
	engine.setViewChange(viewChange)
	err = engine.OnPrepareVote(node.nodeID, pvote, false)
	assert.Nil(t, err)

	// lastViewChangeVotes is nil
	engine.viewChangeVotes = nil
	err = engine.OnPrepareVote(node.nodeID, pvote, false)
	assert.Nil(t, err)

	// whether to accept: accetp
	pvote = makePrepareVote(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address)
	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 1, gen.Hash(), uint32(node.index), node.address, nil)
	engine.setViewChange(viewChange)
	err = engine.OnPrepareVote(node.nodeID, pvote, false)
	assert.Nil(t, err)

	// whether to accept: discard
	higConfirmed := engine.getHighestConfirmed()
	higConfirmed.number = 3
	engine.highestConfirmed.Store(higConfirmed)
	err = engine.OnPrepareVote(node.nodeID, pvote, false)
	assert.Nil(t, err)
}

func TestCbft_OnGetPrepareVote(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, backend, validators := randomCBFT(path, 1)
	gen := backend.chain.Genesis()
	node := nodeIndexNow(validators, engine.startTimeOfEpoch)
	gpv := makeGetPrepareVote(0, gen.Hash())

	err := engine.OnGetPrepareVote(node.nodeID, gpv)
	assert.Nil(t, err)

	engine.blockExtMap = new(BlockExtMap)
	err = engine.OnGetPrepareVote(node.nodeID, gpv)
	assert.Nil(t, err)
}

func TestCbft_OnPrepareVotes(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, backend, validators := randomCBFT(path, 1)
	gen := backend.chain.Genesis()
	node := nodeIndexNow(validators, engine.startTimeOfEpoch)

	pvs := makePrepareVotes(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address)

	err := engine.OnPrepareVotes(node.nodeID, pvs)
	assert.Nil(t, err)

	pvs.Votes[0].ValidatorAddr = common.BytesToAddress([]byte("fake address"))
	err = engine.OnPrepareVotes(node.nodeID, pvs)
	assert.NotNil(t, err)
}

func TestCbft_OnPong(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 2)

	assert.Equal(t, engine.avgLatency(v.validator(1).nodeID), engine.config.MaxLatency)

	assert.Nil(t, engine.OnPong(v.validator(1).nodeID, 5001))
	assert.Nil(t, engine.OnPong(v.validator(1).nodeID, 100))
	assert.Equal(t, engine.avgLatency(v.validator(1).nodeID), int64(100))

	assert.Nil(t, engine.OnPong(v.validator(1).nodeID, 65))
	assert.Equal(t, engine.avgLatency(v.validator(1).nodeID), int64(82))

	assert.Nil(t, engine.OnPong(v.validator(1).nodeID, 10))
	assert.Nil(t, engine.OnPong(v.validator(1).nodeID, 10))
	assert.Nil(t, engine.OnPong(v.validator(1).nodeID, 10))
	assert.Nil(t, engine.OnPong(v.validator(1).nodeID, 10))
	assert.Nil(t, engine.OnPong(v.validator(1).nodeID, 10))
	assert.Equal(t, engine.avgLatency(v.validator(1).nodeID), int64(19))
}

func TestCbft_HighestConfirmedBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, _ := randomCBFT(path, 1)
	assert.Equal(t, engine.HighestConfirmedBlock(), engine.blockChain.Genesis())
}

func TestCbft_inTurnVerify(t *testing.T) {
	path0 := path()
	defer os.RemoveAll(path0)

	engine, _, v := randomCBFT(path0, 4)
	engine.OnPong(v.validator(0).nodeID, 2000)
	assert.False(t, engine.inTurnVerify(0, v.validator(0).nodeID))

	engine.OnPong(v.validator(1).nodeID, 100)
	assert.False(t, engine.inTurnVerify(100, v.validator(1).nodeID))
	assert.True(t, engine.inTurnVerify(10200, v.validator(1).nodeID))

	path1 := path()
	defer os.RemoveAll(path1)
	engine1, _, v1 := randomCBFT(path1, 1)
	engine.OnPong(v1.validator(0).nodeID, 100)
	assert.True(t, engine1.inTurnVerify(100, v1.validator(0).nodeID))
}

func TestCbft_update(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, _ := randomCBFT(path, 1)

	<-time.After(1000 * time.Millisecond)

	assert.Nil(t, engine.eventMux.Post(downloader.StartEvent{}))
	<-time.After(500 * time.Millisecond)
	assert.False(t, engine.isRunning())

	assert.Nil(t, engine.eventMux.Post(downloader.DoneEvent{}))
	<-time.After(500 * time.Millisecond)
	assert.True(t, engine.isRunning())
}

func TestCbft_Consensus(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 1)
	assert.False(t, engine.CheckConsensusNode(common.Address{}))
	assert.True(t, engine.CheckConsensusNode(v.validator(0).address))
	assert.True(t, engine.IsConsensusNode())
}

func TestCbft_VerifyHeader(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, _ := randomCBFT(path, 1)

	header := &types.Header{}
	assert.Equal(t, engine.VerifyHeader(engine.blockChain, header, false), errUnknownBlock)

	header.Number = big.NewInt(1)
	assert.Equal(t, engine.VerifyHeader(engine.blockChain, header, false), errMissingSignature)

	header.Extra = make([]byte, 65)
	assert.Nil(t, engine.VerifyHeader(engine.blockChain, header, false))
}

func TestCbft_CalcBlockDeadline(t *testing.T) {
	path1 := path()
	defer os.RemoveAll(path1)

	now := common.Millis(time.Now())
	engine1, _, _ := randomCBFT(path1, 1)
	tdl, err := engine1.CalcBlockDeadline(now)
	assert.Nil(t, err)
	assert.Equal(t, common.Millis(tdl), now+1000)

	path2 := path()
	defer os.RemoveAll(path2)

	engine2, _, _ := randomCBFT(path2, 4)
	ts := 0
	tdl, err = engine2.CalcBlockDeadline(int64(ts))
	assert.Nil(t, err)
	assert.Equal(t, tdl, common.MillisToTime(int64(ts)).Add(1000*time.Millisecond))

	engine2.config.BlockInterval = 100
	ts = 9900
	tdl, err = engine2.CalcBlockDeadline(int64(ts))
	assert.Nil(t, err)
	assert.Equal(t, tdl, common.MillisToTime(int64(ts)).Add(50*time.Millisecond))
}

func TestCbft_CalcNextBlockTime(t *testing.T) {
	path1 := path()
	defer os.RemoveAll(path1)

	engine1, _, _ := randomCBFT(path1, 1)
	now := common.Millis(time.Now())
	tdl, err := engine1.CalcNextBlockTime(now)
	assert.Nil(t, err)
	assert.Equal(t, common.Millis(tdl), now+1000)

	path2 := path()
	defer os.RemoveAll(path2)

	engine2, _, _ := randomCBFT(path2, 4)
	ts := 0
	tdl, err = engine2.CalcNextBlockTime(int64(ts))
	assert.Nil(t, err)
	assert.Equal(t, tdl, common.MillisToTime(int64(ts)).Add(1000*time.Millisecond))

	ts = 9000
	tdl, err = engine2.CalcNextBlockTime(int64(ts))
	assert.Nil(t, err)
	assert.Equal(t, tdl, common.MillisToTime(int64(ts)).Add(31000*time.Millisecond))

	ts = 10000
	tdl, err = engine2.CalcNextBlockTime(int64(ts))
	assert.Nil(t, err)
	assert.Equal(t, tdl, common.MillisToTime(int64(ts)).Add(30000*time.Millisecond))
}

func TestCbft_updateValidator(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)
	assert.True(t, engine.IsConsensusNode())
	engine.updateValidator()
	assert.True(t, engine.IsConsensusNode())
	assert.Equal(t, v.validator(0).nodeID, engine.getValidators().NodeID(0))

	tv := createTestValidator(createAccount(4))
	newAgency := NewStaticAgency(tv.Nodes())
	oldAgency := engine.agency
	engine.agency = newAgency
	engine.updateValidator()
	assert.False(t, engine.IsConsensusNode())

	engine.agency = oldAgency
	engine.updateValidator()
	assert.True(t, engine.IsConsensusNode())
}

func TestCbft_OnFastSyncCommitHead(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)

	assert.Equal(t, engine.CurrentBlock(), engine.blockChain.Genesis())

	assert.Nil(t, <-engine.FastSyncCommitHead())
	assert.Equal(t, engine.getHighestConfirmed().block, engine.blockChain.Genesis())
	assert.Equal(t, engine.getHighestLogical().block, engine.blockChain.Genesis())
	assert.Equal(t, engine.getRootIrreversible().block, engine.blockChain.Genesis())

	ch := make(chan error, 1)
	engine.OnFastSyncCommitHead(ch)
	err := <-ch
	assert.Nil(t, err)
	assert.Equal(t, engine.getHighestConfirmed().block, engine.blockChain.Genesis())
	assert.Equal(t, engine.getHighestLogical().block, engine.blockChain.Genesis())
	assert.Equal(t, engine.getRootIrreversible().block, engine.blockChain.Genesis())

	timestamp := uint64(common.Millis(time.Now()))
	view := makeViewChange(v.validator(0).privateKey, timestamp, 0, engine.blockChain.Genesis().Hash(), 0, v.validator(0).address, nil)
	engine.viewChange = view

	blocksExt := makeConfirmedBlock(v, engine.blockChain.Genesis().Root(), view, 10)
	for _, ext := range blocksExt {
		statedb, err := engine.blockChain.StateAt(ext.block.Root())
		assert.Nil(t, err)
		_, err = engine.blockChain.WriteBlockWithState(ext.block, nil, statedb)
		assert.Nil(t, err)
	}
	engine.OnFastSyncCommitHead(ch)
	err = <-ch
	assert.Nil(t, err)

	block := blocksExt[9].block
	assert.Equal(t, engine.getHighestLogical().block, block)
	assert.Equal(t, engine.getHighestConfirmed().block, block)
	assert.Equal(t, engine.getRootIrreversible().block, block)
}

func TestCbft_OnNewPrepareBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)
	node := nodeIndexNow(validators, engine.startTimeOfEpoch)
	gen := backend.chain.Genesis()
	block := createBlock(node.privateKey, gen.Hash(), gen.NumberU64()+1)
	propagation := true

	// test Cache prepareBlock
	p := makePrepareBlock(block, node, nil, nil)
	assert.Nil(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation))

	viewChange, _ := engine.newViewChange() // build viewChange

	// test errFutileBlock
	t.Log(viewChange.BaseBlockNum, viewChange.BaseBlockHash.Hex(), viewChange.ProposalIndex, viewChange.ProposalAddr, viewChange.Timestamp)
	assert.EqualError(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation), errFutileBlock.Error())

	viewChangeVotes := buildViewChangeVote(viewChange, validators.neighbors) // build viewChangeVotes

	// test VerifyHeader
	header := &types.Header{Number: big.NewInt(int64(gen.NumberU64() + 1)), ParentHash: gen.Hash()}
	sign, _ := crypto.Sign(header.SealHash().Bytes(), node.privateKey)
	header.Extra = make([]byte, 32)
	copy(header.Extra, sign[0:32])
	block = types.NewBlockWithHeader(header)
	p = makePrepareBlock(block, node, viewChange, viewChangeVotes)
	assert.EqualError(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation), errMissingSignature.Error())

	// test errInvalidatorCandidateAddress
	header = &types.Header{Number: big.NewInt(int64(gen.NumberU64() + 1)), ParentHash: gen.Hash()}
	sign, _ = crypto.Sign(header.SealHash().Bytes(), node.privateKey)
	header.Extra = make([]byte, 32+65)
	copy(header.Extra, sign)
	block = types.NewBlockWithHeader(header)
	p = makePrepareBlock(block, node, viewChange, viewChangeVotes)
	p.ProposalAddr = common.HexToAddress("0x27f7e1d4b9caab9d5b13803cff6da714c51de34e")
	assert.EqualError(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation), errInvalidatorCandidateAddress.Error())

	// test errInvalidatorCandidateAddress
	p.ProposalAddr = node.address
	pri, _ := crypto.GenerateKey()
	engine.config.NodeID = discover.PubkeyID(&pri.PublicKey)
	assert.EqualError(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation), errInvalidatorCandidateAddress.Error())

	// test errTwoThirdViewchangeVotes
	engine.config.NodeID = validators.owner.nodeID
	p.ViewChangeVotes = p.ViewChangeVotes[0:1]
	assert.EqualError(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation), errTwoThirdViewchangeVotes.Error())

	// test errInvalidViewChangeVote
	p = makePrepareBlock(block, node, viewChange, viewChangeVotes)
	p.ViewChangeVotes[2] = forgeViewChangeVote(viewChange)
	assert.EqualError(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation), errInvalidViewChangeVotes.Error())

	// test Discard prepareBlock
	viewChangeVotes = buildViewChangeVote(viewChange, validators.neighbors)
	p = makePrepareBlock(block, node, viewChange, viewChangeVotes)
	assert.Nil(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation))

	// test Accept prepareBlock
	viewChange = makeViewChange(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address, nil)
	viewChangeVotes = buildViewChangeVote(viewChange, validators.AllNodes())
	p = makePrepareBlock(block, node, viewChange, viewChangeVotes)
	assert.Nil(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation))

	// test exists block
	assert.Nil(t, engine.OnNewPrepareBlock(node.nodeID, p, propagation))
}

func TestCbft_AddPrepareBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)
	owner := validators.owner
	gen := backend.chain.Genesis()
	block := createBlock(owner.privateKey, gen.Hash(), gen.NumberU64()+1)
	viewChange, _ := engine.newViewChange()
	t.Log(viewChange.BaseBlockNum, viewChange.BaseBlockHash.Hex(), viewChange.ProposalIndex, viewChange.ProposalAddr, viewChange.Timestamp)
	engine.AddPrepareBlock(block)
}

func TestCbft_OnGetPrepareBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)
	owner := validators.owner
	gen := backend.chain.Genesis()
	block := createBlock(owner.privateKey, gen.Hash(), gen.NumberU64()+1)
	ext := NewBlockExtByPeer(block, block.NumberU64(), len(validators.Nodes()))
	engine.blockExtMap.Add(block.Hash(), block.NumberU64(), ext)
	assert.Nil(t, engine.OnGetPrepareBlock(validators.neighbors[0].nodeID, &getPrepareBlock{Hash: block.Hash(), Number: block.NumberU64()}))
}

func TestCbft_AddJournal(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, backend, validators := randomCBFT(path, 4)
	owner := validators.owner
	gen := backend.chain.Genesis()
	block := createBlock(owner.privateKey, gen.Hash(), gen.NumberU64()+1)

	getPrepareBlock := &getPrepareBlock{Hash: block.Hash(), Number: block.NumberU64()}
	engine.AddJournal(&MsgInfo{getPrepareBlock, validators.neighbors[0].nodeID})
}

func TestCbft_Prepare(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, _ := randomCBFT(path, 1)
	highestBlock := engine.getHighestLogical()

	header := &types.Header{
		Number: big.NewInt(1),
		Extra:  make([]byte, 60),
	}

	ext := &BlockExt{
		block: nil,
	}
	engine.highestLogical.Store(ext)
	assert.NotNil(t, engine.Prepare(engine.blockChain, header))

	engine.highestLogical.Store(highestBlock)
	assert.Nil(t, engine.Prepare(engine.blockChain, header))
	assert.True(t, len(header.Extra) == 97)
}

func TestCbft_VerifySeal(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, _ := randomCBFT(path, 1)

	header := &types.Header{
		Number: big.NewInt(0),
	}
	assert.NotNil(t, engine.VerifySeal(engine.blockChain, header))

	header.Number = big.NewInt(1)
	assert.Nil(t, engine.VerifySeal(engine.blockChain, header))
}

func TestCbft_VerifyHeaders(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, _ := randomCBFT(path, 1)

	header := &types.Header{
		Number: big.NewInt(1),
	}
	_, results := engine.VerifyHeaders(engine.blockChain, []*types.Header{header}, []bool{false})
	assert.NotNil(t, <-results)

	header.Extra = make([]byte, 65)
	_, results = engine.VerifyHeaders(engine.blockChain, []*types.Header{header}, []bool{false})
	assert.Nil(t, <-results)
}
