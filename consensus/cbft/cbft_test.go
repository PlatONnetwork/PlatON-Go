package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
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
	time.Sleep(time.Second * time.Duration(engine.config.Period*2))
	assert.Nil(t, engine.viewChange)
}

func TestCbft_ShouldSeal(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, validators := randomCBFT(path, 4)

	seal, err := engine.ShouldSeal(100)
	assert.False(t, seal)
	assert.Equal(t, errTwoThirdViewchangeVotes, err)

	time.Sleep(time.Second * time.Duration(engine.config.Period) * 2)
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

	close(engine.handler.sendQueue)
	sendQueue := make(chan *MsgPackage, 20)
	engine.handler.sendQueue = sendQueue

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

	wait := sync.WaitGroup{}
	wait.Add(1)
	var msg *MsgPackage
	go func() {
		select {
		case msg = <-engine.handler.sendQueue:
			wait.Done()
		}
	}()
	time.Sleep(2 * time.Second)
	assert.Nil(t, engine.OnGetHighestPrepareBlock(v.validator(2).nodeID, &getHighestPrepareBlock{2}))
	//wait.Wait()
	t.Log("send success")

	//t.Log(len(sendQueue))

	t.Log(msg)

}

func TestCbft_OnHighestPrepareBlock(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)

	engine, _, v := randomCBFT(path, 4)

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
	testCases := []struct{
		pid 		string
		pga 		bool
		msgHash 	common.Hash
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
	pvote= makePrepareVote(node.privateKey, uint64(time.Now().UnixNano()/1e6), 0, gen.Hash(), uint32(node.index), node.address)
	pvote.ValidatorAddr = common.BytesToAddress([]byte("fake address"))
	err  = engine.OnPrepareVote(node.nodeID, pvote, false)
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