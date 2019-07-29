package cbft

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	cvm "github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
)

func init() {
	bls.Init(bls.CurveFp254BNb)
}
func TestThreshold(t *testing.T) {
	f := &Cbft{}
	assert.Equal(t, 1, f.threshold(1))
	assert.Equal(t, 2, f.threshold(2))
	assert.Equal(t, 3, f.threshold(3))
	assert.Equal(t, 3, f.threshold(4))
	assert.Equal(t, 4, f.threshold(5))
	assert.Equal(t, 5, f.threshold(6))
	assert.Equal(t, 5, f.threshold(7))

}

func TestBls(t *testing.T) {
	bls.Init(bls.CurveFp254BNb)
	num := 4
	pk, sk := GenerateKeys(num)
	owner := sk[0]
	nodes := make([]params.CbftNode, num)
	for i := 0; i < num; i++ {

		nodes[i].Node = *discover.NewNode(discover.PubkeyID(&pk[i].PublicKey), nil, 0, 0)
		nodes[i].BlsPubKey = *sk[i].GetPublicKey()
	}

	agency := validator.NewStaticAgency(nodes)

	cbft := &Cbft{
		validatorPool: validator.NewValidatorPool(agency, 0, nodes[0].Node.ID),
		config: ctypes.Config{
			Option: &ctypes.OptionsConfig{
				BlsPriKey: owner,
			},
		},
	}

	pb := &protocols.PrepareVote{}
	cbft.signMsgByBls(pb)
	msg, _ := pb.CannibalizeBytes()
	assert.True(t, cbft.validatorPool.Verify(0, 0, msg, pb.Sign()))
}
func TestAgg(t *testing.T) {
	num := 4
	pk, sk := GenerateKeys(num)
	nodes := make([]params.CbftNode, num)
	for i := 0; i < num; i++ {
		nodes[i].Node = *discover.NewNode(discover.PubkeyID(&pk[i].PublicKey), nil, 0, 0)
		nodes[i].BlsPubKey = *sk[i].GetPublicKey()
	}

	agency := validator.NewStaticAgency(nodes[0:num])

	cnode := make([]*Cbft, num)

	for i := 0; i < num; i++ {
		cnode[i] = &Cbft{
			validatorPool: validator.NewValidatorPool(agency, 0, nodes[0].Node.ID),
			config: ctypes.Config{
				Option: &ctypes.OptionsConfig{
					BlsPriKey: sk[i],
				},
			},
			state: state.NewViewState(),
		}

		cnode[i].state.SetHighestQCBlock(NewBlock(common.Hash{}, 1))
	}

	testPrepareQC(t, cnode)
	testViewChangeQC(t, cnode)
}

func testPrepareQC(t *testing.T, cnode []*Cbft) {
	pbs := make(map[uint32]*protocols.PrepareVote, 0)

	for i := 0; i < len(cnode); i++ {
		pb := &protocols.PrepareVote{ValidatorIndex: uint32(i)}
		assert.NotNil(t, cnode[i])
		cnode[i].signMsgByBls(pb)
		pbs[uint32(i)] = pb
	}
	qc := cnode[0].generatePrepareQC(pbs)

	assert.Nil(t, cnode[0].verifyPrepareQC(qc))
	qc.ValidatorSet = nil
	assert.NotNil(t, cnode[0].verifyPrepareQC(qc))

}
func testViewChangeQC(t *testing.T, cnode []*Cbft) {
	pbs := make(map[uint32]*protocols.ViewChange, 0)

	for i := 0; i < len(cnode); i++ {
		pb := &protocols.ViewChange{BlockHash: common.BigToHash(big.NewInt(int64(i))), BlockNumber: uint64(i), ValidatorIndex: uint32(i)}
		assert.NotNil(t, cnode[i])
		cnode[i].signMsgByBls(pb)
		pbs[uint32(i)] = pb
	}
	qc := cnode[0].generateViewChangeQC(pbs)
	assert.Len(t, qc.QCs, len(cnode))
	_, _, _, num := qc.MaxBlock()
	assert.Equal(t, uint64(len(cnode)-1), num)

	assert.Nil(t, cnode[0].verifyViewChangeQC(qc))
}

func TestNode(t *testing.T) {
	pk, sk, nodes := GenerateCbftNode(4)
	node := MockNode(pk[0], sk[0], nodes, 10000, 10)
	node2 := MockNode(pk[1], sk[1], nodes, 10000, 10)
	assert.Nil(t, node.Start())
	assert.Nil(t, node2.Start())

	testSeal(t, node, node2)
	testPrepare(t, node, node2)
	testTimeout(t, node, node2)
}

func testSeal(t *testing.T, node, node2 *TestCBFT) {
	block := NewBlock(node.chain.Genesis().Hash(), 1)

	result := make(chan *types.Block, 1)

	node.engine.OnSeal(block, result, nil)

	select {
	case b := <-result:
		assert.NotNil(t, b)
	}
}

func testPrepare(t *testing.T, node, node2 *TestCBFT) {
	pb := node.engine.state.PrepareBlockByIndex(0)
	assert.NotNil(t, pb)
	assert.Nil(t, node2.engine.OnPrepareBlock("id", pb))
	pb2 := node2.engine.state.PrepareBlockByIndex(0)
	assert.NotNil(t, pb2)
	_, err := node.engine.verifyConsensusMsg(pb)
	assert.Nil(t, err)
	_, err = node2.engine.verifyConsensusMsg(pb)
	assert.Nil(t, err)
}
func testTimeout(t *testing.T, node, node2 *TestCBFT) {
	time.Sleep(10 * time.Second)
	pb := node.engine.state.PrepareBlockByIndex(0)
	assert.Len(t, node.engine.state.AllViewChange(), 1)
	assert.NotNil(t, node2.engine.OnPrepareBlock(node.engine.config.Option.NodeID.TerminalString(), pb))
	assert.Nil(t, node.engine.OnViewChange(node.engine.config.Option.NodeID.TerminalString(), node.engine.state.AllViewChange()[0]))
	assert.Nil(t, node2.engine.OnViewChange(node.engine.config.Option.NodeID.TerminalString(), node.engine.state.AllViewChange()[0]))
}

func TestChangeView(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	result := make(chan *types.Block, 1)

	parent := nodes[0].chain.Genesis()
	for i := 0; i < 10; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)

		_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i-1), nodes[0].engine.state.MaxQCIndex())
			for j := 1; j < 3; j++ {
				msg := &protocols.PrepareVote{
					Epoch:          nodes[0].engine.state.Epoch(),
					ViewNumber:     nodes[0].engine.state.ViewNumber(),
					BlockIndex:     uint32(i),
					BlockHash:      b.Hash(),
					BlockNumber:    b.NumberU64(),
					ValidatorIndex: uint32(j),
					ParentQC:       qc,
				}
				assert.Nil(t, nodes[j].engine.signMsgByBls(msg))
				assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
			}
			parent = b
		}
	}
	assert.Equal(t, uint64(2), nodes[0].engine.state.ViewNumber())

}

func TestValidatorSwitch(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockValidator(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	switchNode, switchCbftNode := func() (*TestCBFT, params.CbftNode) {
		pk, sk, ns := GenerateCbftNode(1)
		node := MockValidator(pk[0], sk[0], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())
		return node, ns[0]
	}()

	result := make(chan *types.Block, 1)

	parent := nodes[0].chain.Genesis()
	for i := 0; i < 4; i++ {
		for j := 0; j < 10; j++ {
			block := NewBlock(parent.Hash(), parent.NumberU64()+1)
			if i == 1 && j == 0 {
				newUpdateValidatorTx(t, parent, block, cbftnodes, switchCbftNode, nodes[i])
			}
			assert.True(t, nodes[i].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
			nodes[i].engine.OnSeal(block, result, nil)

			_, qc := nodes[i].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
			switchNode.engine.InsertChain(parent)
			select {
			case b := <-result:
				assert.NotNil(t, b)
				assert.Equal(t, uint32(j-1), nodes[i].engine.state.MaxQCIndex())
				for k := 0; k < 4; k++ {
					if k == i {
						continue
					}

					msgPrepare := &protocols.PrepareBlock{
						Epoch:         nodes[i].engine.state.Epoch(),
						ViewNumber:    nodes[i].engine.state.ViewNumber(),
						Block:         block,
						BlockIndex:    uint32(j),
						ProposalIndex: uint32(i),
					}
					if j == 0 {
						msgPrepare.PrepareQC = qc
					}
					assert.Nil(t, nodes[i].engine.signMsgByBls(msgPrepare))
					assert.Nil(t, nodes[k].engine.OnPrepareBlock(fmt.Sprintf("%d", i), msgPrepare))

					msg := &protocols.PrepareVote{
						Epoch:          nodes[i].engine.state.Epoch(),
						ViewNumber:     nodes[i].engine.state.ViewNumber(),
						BlockIndex:     uint32(j),
						BlockHash:      b.Hash(),
						BlockNumber:    b.NumberU64(),
						ValidatorIndex: uint32(k),
						ParentQC:       qc,
					}
					assert.Nil(t, nodes[k].engine.signMsgByBls(msg))
					for kk := 0; kk < 4; kk++ {
						if kk == k {
							continue
						}
						assert.Nil(t, nodes[kk].engine.OnPrepareVote(fmt.Sprintf("%d", i), msg), fmt.Sprintf("number: %d", b.NumberU64()))
					}
				}
				parent = b
			}
		}
	}
	switchNode.engine.InsertChain(parent)

	assert.False(t, switchNode.engine.IsConsensusNode())

	nodes = append(nodes, nodes[:3]...)
	nodes = append(nodes, switchNode)

	block := NewBlock(parent.Hash(), parent.NumberU64()+1)
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil)
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	select {
	case b := <-result:
		assert.NotNil(t, b)
		//assert.Equal(t, , nodes[0].engine.state.MaxQCIndex())
		for i := 1; i < 4; i++ {
			if i == 2 {
				continue
			}
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     0,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: uint32(0),
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number: %d", block.NumberU64()))
			assert.Nil(t, nodes[3].engine.OnPrepareVote("id", msg), fmt.Sprintf("number: %d", block.NumberU64()))
		}
	}
	qcBlock0, _ := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	qcBlock1, _ := nodes[3].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	assert.NotNil(t, qcBlock0)
	assert.NotNil(t, qcBlock1)
	assert.Equal(t, qcBlock0.NumberU64(), qcBlock1.NumberU64())
	assert.Equal(t, qcBlock1.Hash(), qcBlock1.Hash())
}

func newUpdateValidatorTx(t *testing.T, parent, block *types.Block, nodes []params.CbftNode, switchNode params.CbftNode, mineNode *TestCBFT) *types.Block {
	type Vd struct {
		Index     uint            `json:"index"`
		NodeID    discover.NodeID `json:"nodeID"`
		BlsPubKey bls.PublicKey   `json:"blsPubKey"`
	}
	type VdList struct {
		NodeList []*Vd `json:"validateNode"`
	}

	vdl := &VdList{
		NodeList: make([]*Vd, 0),
	}

	for i := 0; i < 3; i++ {
		vdl.NodeList = append(vdl.NodeList, &Vd{
			Index:     uint(i),
			NodeID:    nodes[i].Node.ID,
			BlsPubKey: nodes[i].BlsPubKey,
		})
	}
	vdl.NodeList = append(vdl.NodeList, &Vd{
		Index:     3,
		NodeID:    switchNode.Node.ID,
		BlsPubKey: switchNode.BlsPubKey,
	})

	buf, _ := json.Marshal(vdl)

	param := [][]byte{
		common.Int64ToBytes(2000),
		[]byte("UpdateValidators"),
		buf,
	}
	data, err := rlp.EncodeToBytes(param)
	assert.Nil(t, err)
	signer := types.NewEIP155Signer(chainConfig.ChainID)
	tx, err := types.SignTx(
		types.NewTransaction(
			0,
			vm.ValidatorInnerContractAddr,
			big.NewInt(1000),
			3000*3000,
			big.NewInt(3000),
			data),
		signer,
		testKey)
	assert.Nil(t, err)

	statedb, err := mineNode.chain.StateAt(parent.Root())
	assert.Nil(t, err)
	statedb.Prepare(tx.Hash(), common.Hash{}, 1)
	receipt, _, err := core.ApplyTransaction(chainConfig, nil, &block.Header().Coinbase, nil, statedb, block.Header(), tx, &block.Header().GasUsed, cvm.Config{})
	block, err = mineNode.engine.Finalize(mineNode.chain, block.Header(), statedb, []*types.Transaction{tx}, []*types.Receipt{receipt})
	assert.Nil(t, err)
	return block
}
