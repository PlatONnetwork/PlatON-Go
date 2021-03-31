// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package cbft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"github.com/PlatONnetwork/PlatON-Go/core"
	cstate "github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	cvm "github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var (
	BaseMs = uint64(10000)
)

func init() {
	bls.Init(bls.BLS12_381)
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
	bls.Init(bls.BLS12_381)
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
		validatorPool: validator.NewValidatorPool(agency, 0, 0, nodes[0].Node.ID),
		config: ctypes.Config{
			Option: &ctypes.OptionsConfig{
				BlsPriKey: owner,
			},
		},
	}

	pb := &protocols.PrepareVote{}
	cbft.signMsgByBls(pb)
	msg, _ := pb.CannibalizeBytes()
	assert.Nil(t, cbft.validatorPool.Verify(0, 0, msg, pb.Sign()))
}

func TestPrepareBlockBls(t *testing.T) {
	bls.Init(bls.BLS12_381)
	pk, sk := GenerateKeys(1)
	owner := sk[0]
	node := params.CbftNode{
		Node:      *discover.NewNode(discover.PubkeyID(&pk[0].PublicKey), nil, 0, 0),
		BlsPubKey: *sk[0].GetPublicKey(),
	}
	agency := validator.NewStaticAgency([]params.CbftNode{node})

	cbft := &Cbft{
		validatorPool: validator.NewValidatorPool(agency, 0, 0, node.Node.ID),
		config: ctypes.Config{
			Option: &ctypes.OptionsConfig{
				BlsPriKey: owner,
			},
		},
	}

	header := &types.Header{
		Number:      big.NewInt(int64(1)),
		ParentHash:  common.BytesToHash(utils.Rand32Bytes(32)),
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		ReceiptHash: common.BytesToHash(utils.Rand32Bytes(32)),
		Root:        common.BytesToHash(utils.Rand32Bytes(32)),
		Coinbase:    common.Address{},
		GasLimit:    10000000000,
	}

	txs := make([]*types.Transaction, 0)
	receipts := make([]*types.Receipt, 0)
	for i := 0; i < 1000; i++ {
		tx := types.NewTransaction(uint64(i), common.BytesToAddress(utils.Rand32Bytes(32)), big.NewInt(9000000000), 90000, big.NewInt(11111), []byte{0x99, 0x99, 0x99, 0x99, 0x99, 0x99, 0x99, 0x99, 0x99})
		txs = append(txs, tx)
		receipt := &types.Receipt{
			Status:            types.ReceiptStatusFailed,
			CumulativeGasUsed: 1,
			Logs: []*types.Log{
				{Address: common.BytesToAddress([]byte{0x11})},
				{Address: common.BytesToAddress([]byte{0x01, 0x11})},
			},
			TxHash:          common.BytesToHash([]byte{0x11, 0x11}),
			ContractAddress: common.BytesToAddress([]byte{0x01, 0x11, 0x11}),
			GasUsed:         111111,
		}
		receipts = append(receipts, receipt)
	}

	block := types.NewBlock(header, txs, receipts)
	pb := &protocols.PrepareBlock{
		Epoch:         100,
		ViewNumber:    99,
		Block:         block,
		BlockIndex:    9,
		ProposalIndex: 24,
	}
	tstart := time.Now()
	cbft.signMsgByBls(pb)
	t.Log("sign elapsed", "time", time.Since(tstart))
	msg, _ := pb.CannibalizeBytes()
	assert.Nil(t, cbft.validatorPool.Verify(0, 0, msg, pb.Sign()))
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
			validatorPool: validator.NewValidatorPool(agency, 0, 0, nodes[0].Node.ID),
			config: ctypes.Config{
				Option: &ctypes.OptionsConfig{
					BlsPriKey: sk[i],
				},
			},
			state: state.NewViewState(BaseMs, nil),
		}

		cnode[i].state.SetHighestQCBlock(NewBlock(common.Hash{}, 1))
	}

	testPrepareQC(t, cnode)
	testViewChangeQC(t, cnode)
}

func testPrepareQC(t *testing.T, cnode []*Cbft) {
	pbs := make(map[uint32]*protocols.PrepareVote)

	for i := 0; i < len(cnode); i++ {
		pb := &protocols.PrepareVote{ValidatorIndex: uint32(i)}
		assert.NotNil(t, cnode[i])
		cnode[i].signMsgByBls(pb)
		pbs[uint32(i)] = pb
	}
	qc := cnode[0].generatePrepareQC(pbs)
	fmt.Println(qc)

	assert.Nil(t, cnode[0].verifyPrepareQC(qc.BlockNumber, qc.BlockHash, qc))
	qc.ValidatorSet = nil
	assert.NotNil(t, cnode[0].verifyPrepareQC(qc.BlockNumber, qc.BlockHash, qc))

}

func testViewChangeQC(t *testing.T, cnode []*Cbft) {
	pbs := make(map[uint32]*protocols.ViewChange)

	for i := 0; i < len(cnode); i++ {
		pb := &protocols.ViewChange{BlockHash: common.BigToHash(big.NewInt(int64(i))), BlockNumber: uint64(i), ValidatorIndex: uint32(i)}
		assert.NotNil(t, cnode[i])
		cnode[i].signMsgByBls(pb)
		pbs[uint32(i)] = pb
	}
	qc := cnode[0].generateViewChangeQC(pbs)
	assert.Len(t, qc.QCs, len(cnode))
	_, _, _, _, _, num := qc.MaxBlock()
	assert.Equal(t, uint64(len(cnode)-1), num)

	assert.Nil(t, cnode[0].verifyViewChangeQC(qc))
}

func TestNode(t *testing.T) {
	pk, sk, nodes := GenerateCbftNode(4)
	node := MockNode(pk[0], sk[0], nodes, 5000, 10)
	node2 := MockNode(pk[1], sk[1], nodes, 5000, 10)
	assert.Nil(t, node.Start())
	assert.Nil(t, node2.Start())

	testSeal(t, node, node2)
	testPrepare(t, node, node2)
	testTimeout(t, node, node2)
}

func testSeal(t *testing.T, node, node2 *TestCBFT) {
	block := NewBlock(node.chain.Genesis().Hash(), 1)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	node.engine.Seal(node.cache, block, result, nil, complete)
	<-complete
	//node.engine.OnSeal(block, result, nil)
	b := <-result
	assert.NotNil(t, b)
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

func testExecuteBlock(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()
	for i := 0; i < 8; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil, complete)
		<-complete

		_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i-1), nodes[0].engine.state.MaxQCIndex())
			for j := 1; j < 4; j++ {
				msg := &protocols.PrepareVote{
					Epoch:          nodes[0].engine.state.Epoch(),
					ViewNumber:     nodes[0].engine.state.ViewNumber(),
					BlockIndex:     uint32(i),
					BlockHash:      b.Hash(),
					BlockNumber:    b.NumberU64(),
					ValidatorIndex: uint32(j),
					ParentQC:       qc,
				}
				pb := nodes[0].engine.state.PrepareBlockByIndex(uint32(i))
				assert.NotNil(t, pb)
				assert.Nil(t, nodes[j].engine.OnPrepareBlock("id", pb))
				time.Sleep(50 * time.Millisecond)
				index, finish := nodes[j].engine.state.Executing()
				assert.True(t, index == uint32(i) && finish, fmt.Sprintf("%d,%v", index, finish))
				assert.Nil(t, nodes[j].engine.signMsgByBls(msg))
				assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
				assert.Nil(t, nodes[1].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
			}
			parent = b
		}
	}
	assert.Equal(t, uint64(8), nodes[0].engine.state.HighestQCBlock().NumberU64())
	assert.Equal(t, uint64(8), nodes[1].engine.state.HighestQCBlock().NumberU64())

	//assert.Equal(t, uint64(2), nodes[0].engine.state.ViewNumber())

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
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()
	for i := 0; i < 10; i++ {
		block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil, complete)
		<-complete

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
	assert.Equal(t, uint64(1), nodes[0].engine.state.ViewNumber())
}

func testValidatorSwitch(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockValidator(pk[i], sk[i], cbftnodes, 1000000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	switchNode, switchCbftNode := func() (*TestCBFT, params.CbftNode) {
		pk, sk, ns := GenerateCbftNode(1)
		node := MockValidator(pk[0], sk[0], cbftnodes, 1000000, 10)
		assert.Nil(t, node.Start())
		return node, ns[0]
	}()

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)

	parent := nodes[0].chain.Genesis()
	for i := 0; i < 4; i++ {
		for j := 0; j < 10; j++ {
			block := NewBlock(parent.Hash(), parent.NumberU64()+1)
			header := block.Header()
			assert.Nil(t, nodes[i].engine.Prepare(nodes[i].chain, header))
			if i == 1 && j == 0 {
				tx, receipt, statedb := newUpdateValidatorTx(t, parent, header, cbftnodes, switchCbftNode, nodes[i])
				block, _ = nodes[i].engine.Finalize(nodes[i].chain, header, statedb, []*types.Transaction{tx}, []*types.Receipt{receipt})
				sealHash := block.Header().SealHash()
				nodes[i].cache.WriteStateDB(sealHash, statedb, block.NumberU64())
				nodes[i].cache.WriteReceipts(sealHash, []*types.Receipt{receipt}, block.NumberU64())
			} else {
				statedb, _ := nodes[i].cache.MakeStateDB(parent)
				block, _ = nodes[i].engine.Finalize(nodes[i].chain, header, statedb, []*types.Transaction{}, []*types.Receipt{})
				nodes[i].cache.WriteStateDB(block.Header().SealHash(), statedb, block.NumberU64())
			}

			nodes[i].engine.Seal(nodes[i].chain, block, result, nil, complete)
			<-complete

			_, qc := nodes[i].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
			select {
			case b := <-result:
				assert.NotNil(t, b)
				assert.Equal(t, uint32(j-1), nodes[i].engine.state.MaxQCIndex())

				for k := 0; k < 4; k++ {
					if k == i {
						continue
					}
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
					assert.Nil(t, nodes[i].engine.OnPrepareVote(fmt.Sprintf("%d", i), msg))
				}
				parent = b
				for ii := 0; ii < 4; ii++ {
					if ii == i {
						continue
					}

					qcBlock := nodes[i].engine.state.HighestQCBlock()
					_, qqc := nodes[i].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
					assert.NotNil(t, qqc)
					p := nodes[ii].engine.state.HighestQCBlock()
					assert.Nil(t, nodes[ii].engine.blockCacheWriter.Execute(qcBlock, p), fmt.Sprintf("execute block error, parent: %d block: %d", p.NumberU64(), qcBlock.NumberU64()))
					assert.Nil(t, nodes[ii].engine.OnInsertQCBlock([]*types.Block{qcBlock}, []*ctypes.QuorumCert{qqc}))
				}

				{
					qcBlock := nodes[i].engine.state.HighestQCBlock()
					_, qqc := nodes[i].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
					assert.NotNil(t, qqc)
					p := switchNode.engine.state.HighestQCBlock()
					assert.Nil(t, switchNode.engine.blockCacheWriter.Execute(qcBlock, p), fmt.Sprintf("execute block error, parent: %d block: %d", p.NumberU64(), qcBlock.NumberU64()))
					assert.Nil(t, switchNode.engine.OnInsertQCBlock([]*types.Block{qcBlock}, []*ctypes.QuorumCert{qqc}))
				}
			}
		}
	}

	assert.False(t, switchNode.engine.IsConsensusNode())

	nodes = append(nodes, nodes[:3]...)
	nodes = append(nodes, switchNode)

	block := NewBlock(parent.Hash(), parent.NumberU64()+1)
	header := block.Header()
	assert.Nil(t, nodes[0].engine.Prepare(nodes[0].chain, header))
	statedb, _ := nodes[0].cache.MakeStateDB(parent)
	block, _ = nodes[0].engine.Finalize(nodes[0].chain, header, statedb, []*types.Transaction{}, []*types.Receipt{})
	nodes[0].cache.WriteStateDB(block.Header().SealHash(), statedb, block.NumberU64())
	nodes[0].engine.Seal(nodes[0].chain, block, result, nil, complete)
	<-complete
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
				ValidatorIndex: uint32(i),
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number: %d", block.NumberU64()))
		}
		parent = b
	}
	qcBlock0, _ := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())

	{
		qcBlock := nodes[0].engine.state.HighestQCBlock()
		_, qqc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
		assert.NotNil(t, qqc)
		p := nodes[3].engine.state.HighestQCBlock()
		assert.Nil(t, nodes[3].engine.blockCacheWriter.Execute(qcBlock, p), fmt.Sprintf("execute block error, parent: %d block: %d", p.NumberU64(), qcBlock.NumberU64()))
		assert.Nil(t, nodes[3].engine.OnInsertQCBlock([]*types.Block{qcBlock}, []*ctypes.QuorumCert{qqc}))
	}

	qcBlock1, _ := nodes[3].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	assert.NotNil(t, qcBlock0)
	assert.Equal(t, qcBlock0.NumberU64(), block.NumberU64())
	assert.Equal(t, qcBlock0.NumberU64(), qcBlock1.NumberU64())
	assert.Equal(t, qcBlock0.Hash(), qcBlock1.Hash())
	assert.True(t, nodes[3].engine.IsConsensusNode())
}

func newUpdateValidatorTx(t *testing.T, parent *types.Block, header *types.Header, nodes []params.CbftNode, switchNode params.CbftNode, mineNode *TestCBFT) (*types.Transaction, *types.Receipt, *cstate.StateDB) {
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

	gp := new(core.GasPool).AddGas(10000000000)

	statedb, err := mineNode.cache.MakeStateDB(parent)
	assert.Nil(t, err)
	statedb.Prepare(tx.Hash(), common.Hash{}, 1)
	receipt, _, err := core.ApplyTransaction(chainConfig, mineNode.chain, gp, statedb, header, tx, &header.GasUsed, cvm.Config{})
	assert.Nil(t, err)
	return tx, receipt, statedb
}

func TestCalc(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(1)
	node := MockNode(pk[0], sk[0], cbftnodes, 5000, 10)
	assert.Nil(t, node.Start())

	now := time.Now()
	interval := 500 * time.Millisecond
	blockTime := node.engine.CalcBlockDeadline(now)
	assert.Equal(t, blockTime, now.Add(interval-200*time.Millisecond-150*time.Millisecond))

	nextBlockTime := node.engine.CalcNextBlockTime(now)
	assert.Equal(t, nextBlockTime, now.Add(200*time.Millisecond+150*time.Millisecond))

	time.Sleep(4600 * time.Millisecond)
	old := now
	now = time.Now()
	blockTime = node.engine.CalcBlockDeadline(now)
	assert.Equal(t, blockTime, node.engine.state.Deadline())

	nextBlockTime = node.engine.CalcNextBlockTime(old)
	assert.Equal(t, nextBlockTime, old.Add(500*time.Millisecond))
}

func TestShouldSeal(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(5)
	node := MockNode(pk[0], sk[0], cbftnodes[:4], 3000, 10)
	should, _ := node.engine.ShouldSeal(time.Now())
	assert.False(t, should)

	assert.Nil(t, node.Start())

	should, err := node.engine.ShouldSeal(time.Now())
	assert.Nil(t, err)
	assert.True(t, should)

	node5 := MockNode(pk[4], sk[4], cbftnodes[:4], 3000, 10)
	assert.Nil(t, node5.Start())

	should, err = node5.engine.ShouldSeal(time.Now())
	assert.Equal(t, err.Error(), "current node not a validator")
	assert.False(t, should)

	node1 := MockNode(pk[1], sk[1], cbftnodes[:4], 3000, 10)
	assert.Nil(t, node1.Start())

	should, err = node1.engine.ShouldSeal(time.Now())
	assert.Equal(t, err.Error(), "current node not the proposer")
	assert.False(t, should)

	parent := node.cache.Genesis()
	for i := 0; i < 10; i++ {
		pb := &protocols.PrepareBlock{
			Epoch:         0,
			ViewNumber:    0,
			Block:         NewBlock(parent.Hash(), parent.NumberU64()),
			BlockIndex:    uint32(i),
			ProposalIndex: 0,
		}
		parent = pb.Block
		node.engine.state.AddPrepareBlock(pb)
	}
	should, err = node.engine.ShouldSeal(time.Now())
	assert.Equal(t, err.Error(), "produce block over limit")
	assert.False(t, should)

	time.Sleep(4 * time.Second)
	should, err = node.engine.ShouldSeal(time.Now())
	assert.NotNil(t, err.Error())
	assert.False(t, should)
}

func TestCbft_CreateGenesis(t *testing.T) {
	var db = rawdb.NewMemoryDatabase()
	_, block := CreateGenesis(db)
	fmt.Println(block.Root().Hex())
}

func TestInsertChain(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)

	parent := nodes[0].chain.Genesis()
	hasQCBlock := make([]*types.Block, 0)
	for i := 0; i < 10; i++ {
		block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil, complete)
		<-complete

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
			block, mqc := nodes[0].engine.blockTree.FindBlockAndQC(b.Hash(), b.NumberU64())

			assert.NotNil(t, block)
			qcBytes, _ := ctypes.EncodeExtra(cbftVersion, mqc)

			hasQCBlock = append(hasQCBlock, block.WithBody(nil, qcBytes))
			parent = b
		}
	}
	assert.Equal(t, uint64(1), nodes[0].engine.state.ViewNumber())

	for _, b := range hasQCBlock {
		assert.Nil(t, nodes[1].engine.InsertChain(b))
	}
}

func TestViewChangeCannibalizeBytes(t *testing.T) {
	v := &protocols.ViewChange{
		Epoch:          0,
		ViewNumber:     1,
		BlockHash:      common.HexToHash("0x8ad5dee5aee35b5231ccc19eb2152eb06226031fce7a6ead4f934dc488a1be4c"),
		BlockNumber:    7,
		ValidatorIndex: 2,
		PrepareQC: &ctypes.QuorumCert{
			Epoch:      0,
			ViewNumber: 0,
		},
	}
	vq := ctypes.ViewChangeQuorumCert{
		Epoch:           0,
		ViewNumber:      1,
		BlockHash:       common.HexToHash("0x8ad5dee5aee35b5231ccc19eb2152eb06226031fce7a6ead4f934dc488a1be4c"),
		BlockNumber:     7,
		BlockEpoch:      0,
		BlockViewNumber: 0,
	}

	vc, err := v.CannibalizeBytes()
	assert.Nil(t, err)
	vqc, err := vq.CannibalizeBytes()
	assert.Nil(t, err)

	assert.True(t, bytes.Equal(vc, vqc))

	vq.BlockViewNumber = 1
	vqc, err = vq.CannibalizeBytes()
	assert.Nil(t, err)

	assert.False(t, bytes.Equal(vc, vqc))
}

func Test_StatMessage(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(1)
	node := MockNode(pk[0], sk[0], cbftnodes, 5000, 10)
	err := node.engine.statMessage(nil)
	assert.NotNil(t, err)
	for i := 0; i < 500; i++ {
		pid1 := fmt.Sprintf("id%d", i)
		pid2 := fmt.Sprintf("id%d", i)
		msgInfo1 := ctypes.NewMsgInfo(&protocols.PrepareBlockHash{
			BlockHash:   common.BytesToHash([]byte("1")),
			BlockNumber: uint64(i),
		}, pid1)
		msgInfo2 := ctypes.NewMsgInfo(&protocols.PrepareBlockHash{
			BlockHash:   common.BytesToHash([]byte("1")),
			BlockNumber: uint64(i),
		}, pid2)
		node.engine.statMessage(msgInfo1)
		node.engine.statMessage(msgInfo1)
		node.engine.statMessage(msgInfo2)
	}
	assert.Equal(t, 199, len(node.engine.statQueues))
}
