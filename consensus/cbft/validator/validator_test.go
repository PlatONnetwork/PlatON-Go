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

package validator

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	vm2 "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress = crypto.PubkeyToAddress(testKey.PublicKey)

	//testAccount       = common.HexToAddress("0xbB6223a41b24A4394725257f2eD0F8868ddDf513")
	//testPrivateKey, _ = crypto.HexToECDSA("a765d8c78aa1ea7e535d5326b4c865e8aca446782ffa9c4cdc3c8d1e65acc302")
	//testPubKey        = "26d1fe36f5945af63e8ab09253fc02e3101397095f48d8aadab6dcbdb155419e1148e13ed4a871b109c0d861f4472607260e9477c193beb802ba74b6e73dda4e"

	testValidators = `
{
  "validateNodes": [
    {
      "index": 0,
      "host": "10.10.8.180:8000",
      "nodeID": "70f07dbcfb7348143c0c3a451710f99dfc9711fffaecc0121b24b45f03c02df79ea4bfa89b98f6119c85c2d0f642e525c7e1f6e8101db6fe43eb4244b8bfbc62"
    },
    {
      "index": 1,
      "host": "10.10.8.180:8001",
      "nodeID": "5f1f08b35a0d40765e740d70f259588becc963dd93cc19367ea3da7dff32354381e4cc40f60e15bfd75e6612e3ac02681757e9ef1b5da28990542c99af95000e"
    },
    {
      "index": 2,
      "host": "10.10.8.180:8002",
      "nodeID": "fe3e2dfe55186a636aa43193e4aff711bc8f742bff25a4bc270ba81afb10a07a700ac190f6f6f8fc35ac496350947656347ac44b5a989379213cc014ba6681ce"
    },
    {
      "index": 3,
      "host": "10.10.8.180:8003",
      "nodeID": "de133c1aa904a2050221cd45f2cae2e0c12be335bdb4b1ac28be3b4d724b36edd253fa6843dfd4d48a2a0f1f2a495e1b7ac99e4efa6d63e8bb19ace7b0caaa1c"
    }
  ]
}
`
)

func newTestNode() []params.CbftNode {
	nodes := make([]params.CbftNode, 0)

	n0, _ := discover.ParseNode("enode://e74864b27aecf5cbbfcd523da7657f126b0a5330a970c8264140704d280e6737fd8098d0ee4299706b825771f3d7017aa02f662e4e9a48e9112d93bf05fea66d@127.0.0.1:16789")
	var sec0 bls.SecretKey
	sec0.SetByCSPRNG()
	n1, _ := discover.ParseNode("enode://bf0cd4c95bc3d48cc7999bcf5b3fe6ab9974fd5dabc5253e3e7506c075d0c7a699251caa76672b144be0fc75fe34cee9aaac20753036b0dbd1cb2b3691f26965@127.0.0.1:26789")
	var sec1 bls.SecretKey
	sec1.SetByCSPRNG()
	n2, _ := discover.ParseNode("enode://84c59064dd3b2df54204c52d772cf3809bb0ad6be268843e406f473cef61dacc6d4d4546779dbfa1480deddc64016179ecefdf75d837914f69b679a71ad9711a@127.0.0.1:36789")
	var sec2 bls.SecretKey
	sec2.SetByCSPRNG()
	n3, _ := discover.ParseNode("enode://a9b7e60fa1290c1013cb862c0693d9e87113e8d4cb87d60452749acd978c9fd3a80b49ab5ce7916a5bbfe0b0a0d7e4cde201bd59acccdf97006990156bfe73a5@127.0.0.1:46789")
	var sec3 bls.SecretKey
	sec3.SetByCSPRNG()

	nodes = append(nodes, params.CbftNode{Node: *n0, BlsPubKey: *sec0.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n1, BlsPubKey: *sec1.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n2, BlsPubKey: *sec2.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n3, BlsPubKey: *sec3.GetPublicKey()})
	return nodes
}

func newTestNode3() []params.CbftNode {
	nodes := make([]params.CbftNode, 0)

	n0, _ := discover.ParseNode("enode://e74864b27aecf5cbbfcd523da7657f126b0a5330a970c8264140704d280e6737fd8098d0ee4299706b825771f3d7017aa02f662e4e9a48e9112d93bf05fea66d@127.0.0.1:16789")
	var sec0 bls.SecretKey
	sec0.SetByCSPRNG()
	n1, _ := discover.ParseNode("enode://bf0cd4c95bc3d48cc7999bcf5b3fe6ab9974fd5dabc5253e3e7506c075d0c7a699251caa76672b144be0fc75fe34cee9aaac20753036b0dbd1cb2b3691f26965@127.0.0.1:26789")
	var sec1 bls.SecretKey
	sec1.SetByCSPRNG()
	n2, _ := discover.ParseNode("enode://84c59064dd3b2df54204c52d772cf3809bb0ad6be268843e406f473cef61dacc6d4d4546779dbfa1480deddc64016179ecefdf75d837914f69b679a71ad9711a@127.0.0.1:36789")
	var sec2 bls.SecretKey
	sec2.SetByCSPRNG()

	nodes = append(nodes, params.CbftNode{Node: *n0, BlsPubKey: *sec0.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n1, BlsPubKey: *sec1.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n2, BlsPubKey: *sec2.GetPublicKey()})

	return nodes
}

func TestValidators(t *testing.T) {
	bls.Init(bls.BLS12_381)
	nodes := newTestNode()

	vds := newValidators(nodes, 0)

	assert.True(t, len(nodes) == vds.Len())
	assert.Equal(t, vds.NodeID(0), nodes[0].Node.ID)

	validator, err := vds.FindNodeByID(nodes[2].Node.ID)
	assert.True(t, err == nil, "get node idex fail")
	assert.True(t, validator.Index == 2)

	pubkey, _ := nodes[1].Node.ID.Pubkey()
	addrN1 := crypto.PubkeyToNodeAddress(*pubkey)

	validator, err = vds.FindNodeByID(nodes[1].Node.ID)
	assert.True(t, err == nil, "get node index and address fail")
	assert.Equal(t, validator.Address, addrN1)
	assert.Equal(t, validator.Index, uint32(1))

	idxN1, err := vds.FindNodeByAddress(addrN1)
	assert.True(t, err == nil, "get index by address fail")
	assert.Equal(t, validator.Index, idxN1.Index)

	nl := vds.NodeList()
	assert.True(t, len(nl) == vds.Len())

	emptyNodeID := discover.NodeID{}
	validator, err = vds.FindNodeByID(emptyNodeID)
	assert.True(t, validator == nil)
	assert.True(t, err != nil)

	notFound := vds.NodeID(4)
	assert.Equal(t, notFound, emptyNodeID)

	emptyAddr := common.NodeAddress{}
	validator, err = vds.FindNodeByAddress(emptyAddr)
	assert.True(t, validator == nil)
	assert.True(t, err != nil)

	validator, err = vds.FindNodeByID(emptyNodeID)
	assert.True(t, validator == nil)
	assert.True(t, err != nil)

	node3 := newTestNode3()
	vds3 := newValidators(node3, 0)

	assert.False(t, vds.Equal(vds3))

	badNodes := make([]params.CbftNode, 0)
	badNode, _ := discover.ParseNode("enode://111164b27aecf5cbbfcd523da7657f126b0a5330a970c8264140704d280e6737fd8098d0ee4299706b825771f3d7017aa02f662e4e9a48e9112d93bf05fea66d@127.0.0.1:16789")
	var sec bls.SecretKey
	sec.SetByCSPRNG()
	badNodes = append(badNodes, params.CbftNode{Node: *badNode, BlsPubKey: *sec.GetPublicKey()})
	assert.Panics(t, func() { newValidators(badNodes, 0) })
}

func TestStaticAgency(t *testing.T) {
	bls.Init(bls.BLS12_381)
	nodes := newTestNode()
	vds := newValidators(nodes, 0)

	agency := NewStaticAgency(nodes)
	validators, err := agency.GetValidator(0)
	assert.True(t, err == nil)
	assert.Equal(t, *vds, *validators)
	assert.True(t, agency.Sign(nil) == nil)
	assert.True(t, agency.VerifySign(nil) == nil)
	assert.True(t, agency.GetLastNumber(0) == 0)
}

func genesisBlockForTesting(db ethdb.Database, addr common.Address, balance *big.Int) (*types.Block, *params.ChainConfig) {
	buf, err := ioutil.ReadFile("../../../eth/downloader/testdata/platon.json")
	if err != nil {
		return nil, nil
	}

	var gen core.Genesis
	if err := gen.UnmarshalJSON(buf); err != nil {
		return nil, nil
	}

	gen.Alloc[addr] = core.GenesisAccount{
		Code:    nil,
		Storage: nil,
		Balance: balance,
		Nonce:   0,
	}

	block, _ := gen.Commit(db, snapshotdb.Instance())
	return block, gen.Config

}

func TestInnerAgency(t *testing.T) {
	bls.Init(bls.BLS12_381)
	testdb := rawdb.NewMemoryDatabase()
	balanceBytes, _ := hexutil.Decode("0x2000000000000000000000000000000000000000000000000000000000000")
	balance := big.NewInt(0)
	genesis, chainConfig := genesisBlockForTesting(testdb, testAddress, balance.SetBytes(balanceBytes))

	var vmVds vm.Validators
	err := json.Unmarshal([]byte(testValidators), &vmVds)
	if err != nil {
		panic(err)
	}
	vmVdsBuf, err := json.Marshal(vmVds)
	if err != nil {
		panic(err)
	}

	Uint64ToBytes := func(val uint64) []byte {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, val)
		return buf[:]
	}

	blockchain := core.GenerateBlockChain(chainConfig, genesis, new(consensus.BftMock), testdb, 200, func(i int, block *core.BlockGen) {
		block.SetCoinbase(common.Address{1})

		if i == 50 {
			param := [][]byte{
				common.Int64ToBytes(2000),
				[]byte("UpdateValidators"),
				vmVdsBuf,
			}
			data, err := rlp.EncodeToBytes(param)
			if err != nil {
				panic(err)
			}
			signer := types.NewEIP155Signer(chainConfig.ChainID)
			tx, _ := types.SignTx(
				types.NewTransaction(
					block.TxNonce(testAddress),
					vm2.ValidatorInnerContractAddr,
					big.NewInt(1000),
					3000*3000,
					big.NewInt(3000),
					data),
				signer,
				testKey)
			block.AddTx(tx)
		}

		if i == 59 {
			param := [][]byte{
				common.Int64ToBytes(2003),
				[]byte("SwitchValidators"),
				Uint64ToBytes(uint64(81)),
			}
			data, err := rlp.EncodeToBytes(param)
			if err != nil {
				panic(err)
			}
			signer := types.NewEIP155Signer(chainConfig.ChainID)
			tx, _ := types.SignTx(
				types.NewTransaction(
					block.TxNonce(testAddress),
					vm2.ValidatorInnerContractAddr,
					big.NewInt(1000),
					3000*3000,
					big.NewInt(3000),
					data),
				signer,
				testKey)
			block.AddTx(tx)
		}
	})

	nodes := newTestNode()
	vds := newValidators(nodes, 1)

	agency := NewInnerAgency(nodes, blockchain, 10, 20)

	assert.True(t, agency.GetLastNumber(0) == 40)
	assert.True(t, agency.GetLastNumber(80) == 80, fmt.Sprintf("%d", agency.GetLastNumber(80)))
	assert.True(t, agency.GetLastNumber(81) == 120)
	assert.True(t, agency.GetLastNumber(110) == 120)

	validators, err := agency.GetValidator(0)
	assert.True(t, err == nil)
	assert.Equal(t, *vds, *validators)
	assert.True(t, blockchain.Genesis() != nil)

	newVds, err := agency.GetValidator(81)
	assert.True(t, err == nil)
	assert.True(t, newVds.Len() == 4)
	assert.True(t, newVds.ValidBlockNumber == 81)

	id3 := newVds.NodeID(3)
	assert.Equal(t, id3, vmVds.ValidateNodes[3].NodeID)
	assert.True(t, agency.GetLastNumber(81) == 120)

	assert.True(t, agency.Sign(nil) == nil)
	assert.True(t, agency.VerifySign(nil) == nil)
	assert.True(t, newVds.String() != "")
	assert.False(t, newVds.Equal(validators))

	defaultVds, _ := agency.GetValidator(60)
	assert.True(t, defaultVds.Equal(validators))

	assert.True(t, agency.GetLastNumber(120) == 120)
	assert.True(t, agency.GetLastNumber(121) == 160)
	assert.True(t, agency.GetLastNumber(125) == 160)
	assert.True(t, agency.GetLastNumber(159) == 160)
	assert.True(t, agency.GetLastNumber(160) == 160)
	assert.True(t, agency.GetLastNumber(200) == 200)
	assert.True(t, agency.GetLastNumber(201) == 240)

}

func newTestInnerAgency(nodes []params.CbftNode) consensus.Agency {
	testdb := rawdb.NewMemoryDatabase()
	balanceBytes, _ := hexutil.Decode("0x2000000000000000000000000000000000000000000000000000000000000")
	balance := big.NewInt(0)
	genesis, chainConfig := genesisBlockForTesting(testdb, testAddress, balance.SetBytes(balanceBytes))

	var vmVds vm.Validators
	err := json.Unmarshal([]byte(testValidators), &vmVds)
	if err != nil {
		panic(err)
	}
	vmVdsBuf, err := json.Marshal(vmVds)
	if err != nil {
		panic(err)
	}

	Uint64ToBytes := func(val uint64) []byte {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, val)
		return buf[:]
	}

	blockchain := core.GenerateBlockChain(chainConfig, genesis, new(consensus.BftMock), testdb, 80, func(i int, block *core.BlockGen) {
		block.SetCoinbase(common.Address{1})

		if i == 50 {
			param := [][]byte{
				common.Int64ToBytes(2000),
				[]byte("UpdateValidators"),
				vmVdsBuf,
			}
			data, err := rlp.EncodeToBytes(param)
			if err != nil {
				panic(err)
			}
			signer := types.NewEIP155Signer(chainConfig.ChainID)
			tx, _ := types.SignTx(
				types.NewTransaction(
					block.TxNonce(testAddress),
					vm2.ValidatorInnerContractAddr,
					big.NewInt(1000),
					3000*3000,
					big.NewInt(3000),
					data),
				signer,
				testKey)
			block.AddTx(tx)
		}

		if i == 59 {
			param := [][]byte{
				common.Int64ToBytes(2003),
				[]byte("SwitchValidators"),
				Uint64ToBytes(uint64(81)),
			}
			data, err := rlp.EncodeToBytes(param)
			if err != nil {
				panic(err)
			}
			signer := types.NewEIP155Signer(chainConfig.ChainID)
			tx, _ := types.SignTx(
				types.NewTransaction(
					block.TxNonce(testAddress),
					vm2.ValidatorInnerContractAddr,
					big.NewInt(1000),
					3000*3000,
					big.NewInt(3000),
					data),
				signer,
				testKey)
			block.AddTx(tx)
		}
	})

	agency := NewInnerAgency(nodes, blockchain, 10, 20)
	return agency
}

func TestValidatorPool(t *testing.T) {
	bls.Init(bls.BLS12_381)
	nodes := newTestNode()
	agency := newTestInnerAgency(nodes)

	validatorPool := NewValidatorPool(agency, 0, 0, nodes[0].Node.ID)
	assert.False(t, validatorPool.ShouldSwitch(0))
	assert.True(t, validatorPool.ShouldSwitch(40))

	node, err := validatorPool.GetValidatorByNodeID(0, nodes[0].Node.ID)
	assert.Nil(t, err)
	assert.Equal(t, node.NodeID, nodes[0].Node.ID)

	_, err = validatorPool.GetValidatorByNodeID(0, discover.NodeID{})
	assert.Equal(t, err, errors.New("not found the node"))

	node, err = validatorPool.GetValidatorByIndex(0, 1)
	assert.Nil(t, err)
	assert.Equal(t, node.NodeID, nodes[1].Node.ID)

	_, err = validatorPool.GetValidatorByIndex(0, 5)
	assert.Equal(t, err, errors.New("not found the specified validator"))

	vds := newValidators(nodes, 0)
	node0, _ := vds.FindNodeByIndex(2)
	node, err = validatorPool.GetValidatorByAddr(0, node0.Address)
	assert.Nil(t, err)
	assert.Equal(t, node.Address, node0.Address)

	_, err = validatorPool.GetValidatorByAddr(0, common.NodeAddress{})
	assert.Equal(t, err, errors.New("invalid address"))

	nodeID := validatorPool.GetNodeIDByIndex(0, 4)
	assert.Equal(t, nodeID, discover.NodeID{})

	nodeID = validatorPool.GetNodeIDByIndex(0, 0)
	assert.Equal(t, nodeID, nodes[0].Node.ID)

	index, err := validatorPool.GetIndexByNodeID(0, nodeID)
	assert.Nil(t, err)
	assert.Equal(t, index, uint32(0))

	index, err = validatorPool.GetIndexByNodeID(0, discover.NodeID{})
	assert.Equal(t, err, errors.New("not found the specified validator"))
	assert.Equal(t, index, uint32(0xffffffff))

	nl := validatorPool.ValidatorList(0)
	assert.True(t, len(nl) == len(nodes))

	assert.True(t, validatorPool.IsValidator(0, nodes[0].Node.ID))
	assert.True(t, validatorPool.Len(0) == len(nodes))
	assert.True(t, validatorPool.IsCandidateNode(discover.NodeID{}))

	eventMux := &event.TypeMux{}

	validatorPool.Update(80, 1, eventMux)
	assert.True(t, validatorPool.IsValidator(0, nodes[0].Node.ID))
	assert.False(t, validatorPool.IsValidator(1, nodes[0].Node.ID))
}

func TestValidatorPoolVerify(t *testing.T) {
	bls.Init(bls.BLS12_381)

	nodes := make([]params.CbftNode, 0)

	n0, _ := discover.ParseNode("enode://e74864b27aecf5cbbfcd523da7657f126b0a5330a970c8264140704d280e6737fd8098d0ee4299706b825771f3d7017aa02f662e4e9a48e9112d93bf05fea66d@127.0.0.1:16789")
	var sec0 bls.SecretKey
	sec0.SetByCSPRNG()
	n1, _ := discover.ParseNode("enode://bf0cd4c95bc3d48cc7999bcf5b3fe6ab9974fd5dabc5253e3e7506c075d0c7a699251caa76672b144be0fc75fe34cee9aaac20753036b0dbd1cb2b3691f26965@127.0.0.1:26789")
	var sec1 bls.SecretKey
	sec1.SetByCSPRNG()
	n2, _ := discover.ParseNode("enode://84c59064dd3b2df54204c52d772cf3809bb0ad6be268843e406f473cef61dacc6d4d4546779dbfa1480deddc64016179ecefdf75d837914f69b679a71ad9711a@127.0.0.1:36789")
	var sec2 bls.SecretKey
	sec2.SetByCSPRNG()
	n3, _ := discover.ParseNode("enode://a9b7e60fa1290c1013cb862c0693d9e87113e8d4cb87d60452749acd978c9fd3a80b49ab5ce7916a5bbfe0b0a0d7e4cde201bd59acccdf97006990156bfe73a5@127.0.0.1:46789")
	var sec3 bls.SecretKey
	sec3.SetByCSPRNG()

	priKey, _ := crypto.GenerateKey()
	nodeStr := fmt.Sprintf("enode://%s@127.0.0.1:6666", hex.EncodeToString(crypto.FromECDSAPub(&priKey.PublicKey)[1:]))
	n4, _ := discover.ParseNode(nodeStr)
	var sec4 bls.SecretKey
	sec4.SetByCSPRNG()

	nodes = append(nodes, params.CbftNode{Node: *n0, BlsPubKey: *sec0.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n1, BlsPubKey: *sec1.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n2, BlsPubKey: *sec2.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n3, BlsPubKey: *sec3.GetPublicKey()})
	nodes = append(nodes, params.CbftNode{Node: *n4, BlsPubKey: *sec4.GetPublicKey()})

	agency := newTestInnerAgency(nodes)
	vp := NewValidatorPool(agency, 0, 0, nodes[0].Node.ID)

	m := "test sig"

	sig0 := sec0.Sign(m)
	assert.Nil(t, vp.Verify(0, 0, []byte(m), sig0.Serialize()))
	assert.NotNil(t, vp.Verify(0, 1, []byte(m), sig0.Serialize()))

	sig1 := sec1.Sign(m)
	sig0.Add(sig1)

	sig2 := sec2.Sign(m)
	sig0.Add(sig2)

	sig3 := sec3.Sign(m)
	sig0.Add(sig3)

	assert.True(t, vp.VerifyAggSig(0, []uint32{0, 1, 2, 3}, []byte(m), sig0.Serialize()))
	assert.True(t, vp.VerifyAggSig(0, []uint32{1, 0, 3, 2}, []byte(m), sig0.Serialize()))
	assert.False(t, vp.VerifyAggSig(0, []uint32{0, 1, 2}, []byte(m), sig0.Serialize()))

	header := types.Header{
		ParentHash: common.Hash{},
		Root:       common.Hash{},
		Number:     big.NewInt(0),
		Time:       big.NewInt(time.Now().Unix()),
		Extra:      make([]byte, 32+consensus.ExtraSeal),
	}
	sig, err := crypto.Sign(header.SealHash().Bytes(), priKey)
	assert.Nil(t, err)
	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sig[:])

	assert.Nil(t, vp.VerifyHeader(&header))

	priKey1, _ := crypto.GenerateKey()
	sigWrong, _ := crypto.Sign(header.SealHash().Bytes(), priKey1)
	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sigWrong[:])
}

type mockAgency struct {
	consensus.Agency

	lastNumber uint64
}

func newMockAgency(lastNumber uint64) *mockAgency {
	return &mockAgency{
		lastNumber: lastNumber,
	}
}

func (m *mockAgency) Flush(*types.Header) error { return nil }

func (m *mockAgency) Sign(interface{}) error { return nil }

func (m *mockAgency) VerifySign(interface{}) error { return nil }

func (m *mockAgency) VerifyHeader(*types.Header, *state.StateDB) error { return nil }

func (m *mockAgency) GetLastNumber(blockNumber uint64) uint64 { return m.lastNumber }

func (m *mockAgency) GetValidator(blockNumber uint64) (*cbfttypes.Validators, error) {
	return &cbfttypes.Validators{
		ValidBlockNumber: blockNumber,
	}, nil
}

func (m *mockAgency) IsCandidateNode(discover.NodeID) bool { return false }

func (m *mockAgency) OnCommit(block *types.Block) error { return nil }

func TestValidatorPoolReset(t *testing.T) {
	agency := newMockAgency(100)
	vp := NewValidatorPool(agency, 0, 0, discover.NodeID{})

	vp.Reset(100, 10)
	assert.Equal(t, vp.switchPoint, uint64(100))
	assert.Equal(t, vp.epoch, uint64(11))

	agency.lastNumber = 200
	vp.Reset(150, 15)
	assert.Equal(t, vp.epoch, uint64(15))
	assert.Equal(t, vp.switchPoint, uint64(149))
}
