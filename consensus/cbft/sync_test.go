package cbft

//
//import (
//	"fmt"
//	"github.com/PlatONnetwork/PlatON-Go/common"
//	"github.com/PlatONnetwork/PlatON-Go/core/types"
//	"github.com/PlatONnetwork/PlatON-Go/p2p"
//	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
//	"math/big"
//	"testing"
//	"time"
//)
//
//type aa struct {
//	hash common.Hash
//}
//
//func (a aa) String() string {
//	return a.hash.String()
//}
//func TestNewFlowControl(t *testing.T) {
//	a := &aa{
//		hash: common.HexToHash("0x59f18af5b772a730d0fc3ba105d9343b0c9202d0a48d7f87c149d70aeb7d2ae2"),
//	}
//	fmt.Println(a.String())
//	b := true
//	hash := common.HexToHash("0x59f18af5b772a730d0fc3ba105d9343b0c9202d0a48d7f87c149d70aeb7d2ae2")
//	fmt.Println(fmt.Sprintf("%s, %v", hash.TerminalString(), b))
//}
//func TestSync(t *testing.T) {
//	//sync := NewSyncing()
//
//	closer, rw, _, errc := p2p.NewMockPeer(nil)
//
//	defer closer()
//
//	sba := &signBitArray{
//		BlockHash: common.HexToHash("0x1"),
//		BlockNum:  big.NewInt(2),
//		SignBits:  NewBitArray(2),
//	}
//	if err := p2p.Send(rw, 16+2, sba); err != nil {
//		t.Error(err)
//	}
//	//Send(rw, baseProtocolLength+3, []uint{2})
//	//Send(rw, baseProtocolLength+4, []uint{3})
//
//	select {
//	case err := <-errc:
//		if err != nil {
//			t.Errorf("peer returned error: %v", err)
//		}
//	case <-time.After(2 * time.Second):
//		t.Errorf("receive timeout")
//	}
//}
//
//func initCbftSync(tc *testCbft) {
//	tc.cbft.producerBlocks = NewProducerBlocks(tc.cbft.dpos.primaryNodeList[0], big.NewInt(1))
//	block := types.NewBlock(&types.Header{
//		GasLimit: 0,
//		Number:   big.NewInt(1),
//	}, nil, nil)
//
//	tc.cbft.producerBlocks.AddBlock(block)
//	tc.cbft.producerBlocks.AddBlock(types.NewBlock(&types.Header{
//		GasLimit: 0,
//		Number:   big.NewInt(2),
//	}, nil, nil))
//	tc.cbft.producerBlocks.AddBlock(types.NewBlock(&types.Header{
//		GasLimit: 0,
//		Number:   big.NewInt(3),
//	}, nil, nil))
//	tc.cbft.rootIrreversible.Store(NewBlockExt(block, block.NumberU64()))
//}
//func TestSyncPeer(t *testing.T) {
//	tc := newCbft()
//	initCbftSync(tc)
//
//	nodeId1 := discover.PubkeyID(&tc.nodes[tc.cbft.dpos.primaryNodeList[0]].PublicKey)
//
//	nodeId2 := discover.PubkeyID(&tc.nodes[tc.cbft.dpos.primaryNodeList[1]].PublicKey)
//	cbft := tc.cbft
//	cbft.handler = NewHandler(cbft)
//	_, rw1, p2, rw2 := p2p.NewPeerByNodeID(nodeId1, nodeId2, cbft.Protocols())
//	p2p.Send(rw1, ConsensusStateMsg, &consensusState{IrreversibleBlockNum: big.NewInt(0), MemMaxBlockNum: big.NewInt(1)})
//
//	time.Sleep(time.Second * 3)
//}
