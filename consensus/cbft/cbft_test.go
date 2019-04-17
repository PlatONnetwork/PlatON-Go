package cbft

//
//import (
//	"crypto/ecdsa"
//	"github.com/PlatONnetwork/PlatON-Go/common"
//	"github.com/PlatONnetwork/PlatON-Go/core"
//	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
//	"github.com/PlatONnetwork/PlatON-Go/core/types"
//	"github.com/PlatONnetwork/PlatON-Go/crypto"
//	"github.com/PlatONnetwork/PlatON-Go/event"
//	"github.com/PlatONnetwork/PlatON-Go/p2p"
//	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
//	"github.com/PlatONnetwork/PlatON-Go/params"
//	"github.com/go-errors/errors"
//	"github.com/stretchr/testify/assert"
//	"math/big"
//	"testing"
//	"time"
//)
//
//type testCbftConfig struct {
//	c     *params.CbftConfig
//	nodes map[discover.NodeID]*ecdsa.PrivateKey
//}
//
//type testCbft struct {
//	testCbftConfig
//	cbft *Cbft
//}
//
//func (t *testCbft) init(c *testCbftConfig) {
//	in := make([]discover.NodeID, 0)
//
//	for k, _ := range c.nodes {
//		in = append(in, k)
//	}
//	t.cbft = &Cbft{}
//	t.cbft.config = c.c
//	t.nodes = c.nodes
//	t.cbft.dpos = newDpos(in)
//}
//
//func (t testCbftConfig) signs(msg []byte) []common.BlockConfirmSign {
//	signs := make([]common.BlockConfirmSign, 0)
//	for _, v := range t.nodes {
//		s, _ := crypto.Sign(msg, v)
//		signs = append(signs, *common.NewBlockConfirmSign(s))
//	}
//	return signs
//}
//
//type nullConn struct {
//	p2p.MsgReadWriter
//}
//
//func (n nullConn) ReadMsg() (p2p.Msg, error) {
//	return p2p.Msg{}, nil
//}
//
//func (n nullConn) WriteMsg(p2p.Msg) error {
//	return nil
//}
//
//func newNodeId() (*ecdsa.PrivateKey, discover.NodeID) {
//	key, err := crypto.GenerateKey()
//	if err != nil {
//		panic("new key failed")
//	}
//	nodeId := discover.PubkeyID(&key.PublicKey)
//	return key, nodeId
//}
//
//func newCbftConfig(nodes int) *testCbftConfig {
//	config := &testCbftConfig{
//		c: &params.CbftConfig{
//			Period:           0,
//			Epoch:            0,
//			MaxLatency:       0,
//			LegalCoefficient: 0,
//			Duration:         0,
//		},
//	}
//
//	config.nodes = make(map[discover.NodeID]*ecdsa.PrivateKey)
//
//	for i := 0; i < nodes; i++ {
//		k, n := newNodeId()
//		config.nodes[n] = k
//		config.c.NodeID = n
//		config.c.PrivateKey = k
//	}
//	return config
//}
//
//func newCbft() *testCbft {
//	tc := newCbftConfig(4)
//
//	t := &testCbft{}
//	t.init(tc)
//
//	t.cbft.clearViewChange()
//	t.cbft.clearPending()
//	t.cbft.handler = NewHandler(t.cbft)
//	t.cbft.maxVotedBlockNum = big.NewInt(0)
//	return t
//}
//func TestViewChangeConn(t *testing.T) {
//	tc := newCbftConfig(4)
//	in := make([]discover.NodeID, 0)
//	for k, _ := range tc.nodes {
//		in = append(in, k)
//	}
//
//	ntc := newCbft()
//	ntc.cbft.handler = NewHandler(ntc.cbft)
//	ntc.cbft.maxVotedBlockNum = big.NewInt(0)
//	hash := common.HexToHash("0x59f18af5b772a730d0fc3ba105d9343b0c9202d0a48d7f87c149d70aeb7d2ae2")
//
//	closer, rw, _, _ := p2p.NewMockPeerNodeID(in[0], ntc.cbft.Protocols())
//	defer closer()
//
//	v := &viewChange{
//		Timestamp:              big.NewInt(1),
//		IrreversibleBlockNum:   big.NewInt(1),
//		IrreversibleBlockHash:  hash,
//		IrreversibleBlockSigns: tc.signs(hash[:]),
//	}
//	s, err := ntc.cbft.signFn(hash[:])
//	if err != nil {
//		t.Error(err)
//	}
//	v.Sign.Set(s)
//	if err := p2p.Send(rw, 16+ViewChangeMsg, v); err != nil {
//		t.Error(err)
//	}
//	time.Sleep(time.Second)
//	assert.Equal(t, ntc.cbft.viewChange.Timestamp, v.Timestamp, "timestamp is not equal")
//	assert.Equal(t, ntc.cbft.viewChange.IrreversibleBlockNum, v.IrreversibleBlockNum, "block number is not equal")
//}
//
//func TestViewChange(t *testing.T) {
//	ntc := newCbft()
//	in := ntc.cbft.dpos.primaryNodeList
//	p := p2p.NewPeer(in[0], "", nil)
//	ntc.cbft.handler.peers.Register(newPeer(p, nullConn{}))
//	hash := common.HexToHash("0x59f18af5b772a730d0fc3ba105d9343b0c9202d0a48d7f87c149d70aeb7d2ae2")
//
//	v := &viewChange{
//		Timestamp:              big.NewInt(1),
//		IrreversibleBlockNum:   big.NewInt(1),
//		IrreversibleBlockHash:  hash,
//		IrreversibleBlockSigns: ntc.signs(hash[:]),
//	}
//	if s, err := ntc.cbft.signFn(hash[:]); err != nil {
//		t.Error(err)
//	} else {
//		v.Sign.Set(s)
//	}
//
//	if err := ntc.cbft.OnViewChange(in[0], v); err != nil {
//		t.Error(err)
//	}
//
//	assert.Equal(t, ntc.cbft.viewChange.Timestamp, v.Timestamp, "timestamp is not equal")
//	assert.Equal(t, ntc.cbft.viewChange.IrreversibleBlockNum, v.IrreversibleBlockNum, "block number is not equal")
//
//	ntc.cbft.viewChange = &viewChange{
//		Timestamp:              big.NewInt(2),
//		IrreversibleBlockNum:   big.NewInt(2),
//		IrreversibleBlockHash:  hash,
//		IrreversibleBlockSigns: nil,
//	}
//
//	v2 := &viewChange{
//		Timestamp:              big.NewInt(1),
//		IrreversibleBlockNum:   big.NewInt(1),
//		IrreversibleBlockHash:  hash,
//		IrreversibleBlockSigns: ntc.signs(hash[:]),
//	}
//	if s, err := ntc.cbft.signFn(hash[:]); err != nil {
//		t.Error(err)
//	} else {
//		v2.Sign.Set(s)
//	}
//
//	assert.Equal(t, ntc.cbft.OnViewChange(in[0], v2), errTimestamp, "")
//}
//
//func TestViewChangeResp(t *testing.T) {
//	ntc := newCbft()
//	in := ntc.cbft.dpos.primaryNodeList
//	p := p2p.NewPeer(in[0], "", nil)
//	ntc.cbft.handler.peers.Register(newPeer(p, nullConn{}))
//	hash := common.HexToHash("0x59f18af5b772a730d0fc3ba105d9343b0c9202d0a48d7f87c149d70aeb7d2ae2")
//
//	v := &viewChange{
//		Timestamp:              big.NewInt(1),
//		IrreversibleBlockNum:   big.NewInt(1),
//		IrreversibleBlockHash:  hash,
//		IrreversibleBlockSigns: ntc.signs(hash[:]),
//	}
//
//	if s, err := ntc.cbft.signFn(hash[:]); err != nil {
//		t.Error(err)
//	} else {
//		v.Sign.Set(s)
//	}
//	ntc.cbft.viewChange = v
//
//	clearFunc := func() {
//		ntc.cbft.clearViewChange()
//		ntc.cbft.viewChange = v
//	}
//
//	testCase := func(resp *viewChangeResp, expect error) error {
//		if s, err := crypto.Sign(resp.BlockHash[:], ntc.nodes[resp.ID]); err != nil {
//			t.Error(err)
//		} else {
//			resp.Sign.Set(s)
//		}
//
//		if err := ntc.cbft.OnViewChangeResp(resp.ID, resp); err != expect {
//			t.Error(err)
//			return errors.New("expect error")
//		} else {
//			return err
//		}
//	}
//
//	resp := &viewChangeResp{
//		ID:        in[0],
//		Timestamp: v.Timestamp,
//		BlockNum:  v.IrreversibleBlockNum,
//		BlockHash: v.IrreversibleBlockHash,
//	}
//
//	assert.Equal(t, testCase(resp, nil), nil, "")
//
//	clearFunc()
//	resp = &viewChangeResp{
//		ID:        in[0],
//		Timestamp: big.NewInt(2),
//		BlockNum:  v.IrreversibleBlockNum,
//		BlockHash: v.IrreversibleBlockHash,
//	}
//
//	assert.Equal(t, testCase(resp, errTimestampNotEqual), errTimestampNotEqual, "")
//
//	clearFunc()
//	resp = &viewChangeResp{
//		ID:        in[0],
//		Timestamp: v.Timestamp,
//		BlockNum:  v.IrreversibleBlockNum,
//		BlockHash: common.HexToHash("0x59f18af5b772a730d0fc3ba105d9343b0c9202d0a48d7f87c149d70aeb7d2ae3"),
//	}
//
//	assert.Equal(t, testCase(resp, errBlockHashNotEqual), errBlockHashNotEqual, "")
//}
//
//func TestSendViewChange(t *testing.T) {
//	ntc := newCbft()
//
//	block := types.NewBlock(&types.Header{
//		GasLimit: 0,
//		Number:   big.NewInt(1),
//	}, nil, nil)
//	ntc.cbft.rootIrreversible.Store(NewBlockExt(block, block.NumberU64()))
//	if err := ntc.cbft.sendViewChange(); err != nil {
//		t.Error(err)
//		return
//	}
//	assert.True(t, ntc.cbft.hadSendViewChange())
//}
//
//func TestSendViewChangeTimeout(t *testing.T) {
//	ntc := newCbft()
//
//	block := types.NewBlock(&types.Header{
//		GasLimit: 0,
//		Number:   big.NewInt(1),
//	}, nil, nil)
//	ntc.cbft.rootIrreversible.Store(NewBlockExt(block, block.NumberU64()))
//	if err := ntc.cbft.sendViewChange(); err != nil {
//		t.Error(err)
//		return
//	}
//	time.Sleep(time.Second * 2)
//	assert.True(t, !ntc.cbft.hadSendViewChange())
//}
//
//func TestBroadcastPrepareBlockPending(t *testing.T) {
//	ntc := newCbft()
//	in := ntc.cbft.dpos.primaryNodeList
//	p := p2p.NewPeer(in[0], "", nil)
//	ntc.cbft.handler.peers.Register(newPeer(p, nullConn{}))
//	hash := common.HexToHash("0x59f18af5b772a730d0fc3ba105d9343b0c9202d0a48d7f87c149d70aeb7d2ae2")
//
//	v := &viewChange{
//		Timestamp:              big.NewInt(1),
//		IrreversibleBlockNum:   big.NewInt(1),
//		IrreversibleBlockHash:  hash,
//		IrreversibleBlockSigns: ntc.signs(hash[:]),
//	}
//
//	if s, err := ntc.cbft.signFn(hash[:]); err != nil {
//		t.Error(err)
//	} else {
//		v.Sign.Set(s)
//	}
//	ntc.cbft.viewChange = v
//
//	go ntc.cbft.broadcastPrepare()
//	ntc.cbft.eventMux = new(event.TypeMux)
//	ntc.cbft.prepareMinedBlockSub = ntc.cbft.eventMux.Subscribe(core.PrepareMinedBlockEvent{})
//	if err := ntc.cbft.eventMux.Post(core.PrepareMinedBlockEvent{Block: types.NewBlock(&types.Header{
//		GasLimit: 0,
//		Number:   big.NewInt(2),
//	}, nil, nil),
//	}); err != nil {
//		t.Error(err)
//	}
//
//	if err := ntc.cbft.eventMux.Post(core.PrepareMinedBlockEvent{Block: types.NewBlock(&types.Header{
//		GasLimit: 0,
//		Number:   big.NewInt(3),
//	}, nil, nil),
//	}); err != nil {
//		t.Error(err)
//	}
//
//	time.Sleep(time.Second)
//	ntc.cbft.mux.Lock()
//	defer ntc.cbft.mux.Unlock()
//
//	assert.Equal(t, 2, len(ntc.cbft.pendingBlocks))
//}
//
//func TestBroadcastBlockSigns(t *testing.T) {
//	ntc := newCbft()
//	in := ntc.cbft.dpos.primaryNodeList
//	p := p2p.NewPeer(in[0], "", nil)
//	ntc.cbft.handler.peers.Register(newPeer(p, nullConn{}))
//	hash := common.HexToHash("0x59f18af5b772a730d0fc3ba105d9343b0c9202d0a48d7f87c149d70aeb7d2ae2")
//
//	v := &viewChange{
//		Timestamp:              big.NewInt(1),
//		IrreversibleBlockNum:   big.NewInt(1),
//		IrreversibleBlockHash:  hash,
//		IrreversibleBlockSigns: ntc.signs(hash[:]),
//	}
//
//	if s, err := ntc.cbft.signFn(hash[:]); err != nil {
//		t.Error(err)
//	} else {
//		v.Sign.Set(s)
//	}
//	ntc.cbft.viewChange = v
//
//	go ntc.cbft.broadcastSigns()
//	ntc.cbft.eventMux = new(event.TypeMux)
//	ntc.cbft.blockSignatureSub = ntc.cbft.eventMux.Subscribe(core.BlockSignatureEvent{})
//
//	sign2 := hash
//	sign2[2] = 2
//	if err := ntc.cbft.eventMux.Post(core.BlockSignatureEvent{BlockSignature: &cbfttypes.BlockSignature{
//		SignHash:  sign2,
//		Hash:      sign2,
//		Number:    big.NewInt(2),
//		Signature: new(common.BlockConfirmSign),
//	}}); err != nil {
//		t.Error(err)
//	}
//
//	sign3 := hash
//	sign3[3] = 3
//	if err := ntc.cbft.eventMux.Post(core.BlockSignatureEvent{BlockSignature: &cbfttypes.BlockSignature{
//		SignHash:  sign3,
//		Hash:      sign3,
//		Number:    big.NewInt(3),
//		Signature: new(common.BlockConfirmSign),
//	}}); err != nil {
//		t.Error(err)
//	}
//
//	time.Sleep(time.Second)
//	ntc.cbft.mux.Lock()
//	defer ntc.cbft.mux.Unlock()
//
//	assert.Equal(t, 2, len(ntc.cbft.pendingVotes))
//}
