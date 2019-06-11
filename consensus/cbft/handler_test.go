package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var id = randomID()

func TestNewHandler(t *testing.T) {
	cbft := &Cbft{
		log: log.New(),
	}
	router := NewRouter(cbft, nil)
	cbft.router = router
	handler := NewHandler(cbft)

	pos := handler.Protocols()
	assert.Equal(t, 1, len(pos))

	ps := handler.PeerSet()
	assert.Equal(t, 0, len(ps.peers))

	p2pPeer := p2p.NewPeer(id, "pid", nil)
	fake := &fakeMessageRW{}
	peer := newPeer(p2pPeer, fake)
	ps.Register(peer)

	handler.Send(id, &prepareBlockHash{})
	assert.Equal(t, 1, len(handler.sendQueue))

	handler.SendAllConsensusPeer(&cbftStatusData{})
	assert.Equal(t,2, len(handler.sendQueue))

	handler.SendBroadcast(&cbftStatusData{})
	assert.Equal(t, 3,len(handler.sendQueue))

	handler.SendPartBroadcast(&cbftStatusData{})
	assert.Equal(t, 4, len(handler.sendQueue))

	handler.Start()
	time.Sleep(1 * time.Second)
	close(handler.sendQueue)
}

func TestBaseHandler(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 1)

	handler := NewHandler(engine)
	p2pPeer := p2p.NewPeer(id, "pid", nil)
	fake := &fakeMessageRW{}
	//peer := newPeer(p2pPeer, fake)

	err := handler.handler(p2pPeer, fake)
	assert.NotNil(t, err)
}

func TestHandlerMsg(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 1)

	handler := NewHandler(engine)
	p2pPeer := p2p.NewPeer(id, "pid", nil)
	fake := &fakeMessageRW{}
	peer := newPeer(p2pPeer, fake)

	// test
	testCases := []struct{
		code uint64
		msg Message
		err error
	}{
		{CBFTStatusMsg, &cbftStatusData{}, nil},
		{CBFTStatusMsg, &cbftStatusData{}, fmt.Errorf("error")},
		{GetPrepareBlockMsg, &getPrepareBlock{}, nil},
		{PrepareBlockMsg, makeFakePrepareBlock(), nil},
		{PrepareVoteMsg, makeFakePrepareVote(), nil},
		{PrepareVoteMsg, makeFakePrepareVote(), fmt.Errorf("error")},
		{ViewChangeMsg, makeFakeViewChange(), nil},
		{ViewChangeMsg, makeFakeViewChange(), fmt.Errorf("error")},
		{ViewChangeVoteMsg, &viewChangeVote{}, nil},
		{ConfirmedPrepareBlockMsg, makeFakeConfirmedPrepareBlock(), nil},
		{GetPrepareVoteMsg, makeFakeGetPrepareVote(), nil},
		{GetPrepareVoteMsg, makeFakeGetPrepareVote(), fmt.Errorf("error")},
		{PrepareVotesMsg, &prepareVotes{}, nil},
		{GetHighestPrepareBlockMsg, &getHighestPrepareBlock{}, nil},
		{HighestPrepareBlockMsg, &highestPrepareBlock{}, nil},
		{PrepareBlockHashMsg, &prepareBlockHash{}, nil},
	}

	for _, v := range testCases {
		size, r, err := rlp.EncodeToReader(v.msg)
		if err != nil {
			t.Errorf("encode fail")
		}
		msg := p2p.Msg{Code: v.code, Size: uint32(size), Payload: r}
		fake.update(msg, v.err)
		err = handler.handleMsg(peer)
		if v.code == CBFTStatusMsg {
			assert.NotNil(t, err)
		}
	}
}

type fakeMessageRW struct {
	msg p2p.Msg
	err error
}

func (rw *fakeMessageRW) update(msg p2p.Msg, err error) {
	rw.msg = msg
	rw.err = err
}

func (rw *fakeMessageRW) ReadMsg() (p2p.Msg, error) {
	return rw.msg, rw.err
}

func (rw *fakeMessageRW) WriteMsg(msg p2p.Msg) error {
	return fmt.Errorf("fake error")
}