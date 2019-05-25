package cbft

import (
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
)

var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress = crypto.PubkeyToAddress(testKey.PublicKey)
)

func newTestNode() []discover.Node {
	nodes := make([]discover.Node, 0)

	n0, _ := discover.ParseNode("enode://e74864b27aecf5cbbfcd523da7657f126b0a5330a970c8264140704d280e6737fd8098d0ee4299706b825771f3d7017aa02f662e4e9a48e9112d93bf05fea66d@127.0.0.1:16789")
	n1, _ := discover.ParseNode("enode://bf0cd4c95bc3d48cc7999bcf5b3fe6ab9974fd5dabc5253e3e7506c075d0c7a699251caa76672b144be0fc75fe34cee9aaac20753036b0dbd1cb2b3691f26965@127.0.0.1:26789")
	n2, _ := discover.ParseNode("enode://84c59064dd3b2df54204c52d772cf3809bb0ad6be268843e406f473cef61dacc6d4d4546779dbfa1480deddc64016179ecefdf75d837914f69b679a71ad9711a@127.0.0.1:36789")
	n3, _ := discover.ParseNode("enode://a9b7e60fa1290c1013cb862c0693d9e87113e8d4cb87d60452749acd978c9fd3a80b49ab5ce7916a5bbfe0b0a0d7e4cde201bd59acccdf97006990156bfe73a5@127.0.0.1:46789")

	nodes = append(nodes, *n0)
	nodes = append(nodes, *n1)
	nodes = append(nodes, *n2)
	nodes = append(nodes, *n3)
	return nodes
}

func TestValidators(t *testing.T) {
	nodes := newTestNode()

	vds := newValidators(nodes, 0)

	t.Log(vds)

	assert.True(t, len(nodes) == vds.Len())
	assert.Equal(t, vds.NodeID(0), nodes[0].ID)

	validator, err := vds.NodeIndex(nodes[2].ID)
	assert.True(t, err == nil, "get node idex fail")
	assert.True(t, validator.Index == 2)

	pubkey, err := nodes[1].ID.Pubkey()
	addrN1 := crypto.PubkeyToAddress(*pubkey)

	validator, err = vds.NodeIndexAddress(nodes[1].ID)
	assert.True(t, err == nil, "get node index and address fail")
	assert.Equal(t, validator.Address, addrN1)
	assert.Equal(t, validator.Index, 1)

	idxN1, err := vds.AddressIndex(addrN1)
	assert.True(t, err == nil, "get index by address fail")
	assert.Equal(t, validator.Index, idxN1.Index)

	nl := vds.NodeList()
	assert.True(t, len(nl) == vds.Len())

	emptyNodeID := discover.NodeID{}
	validator, err = vds.NodeIndexAddress(emptyNodeID)
	assert.True(t, validator == nil)
	assert.True(t, err != nil)

	notFound := vds.NodeID(4)
	assert.Equal(t, notFound, emptyNodeID)

	emptyAddr := common.Address{}
	validator, err = vds.AddressIndex(emptyAddr)
	assert.True(t, validator == nil)
	assert.True(t, err != nil)

	validator, err = vds.NodeIndex(emptyNodeID)
	assert.True(t, validator == nil)
	assert.True(t, err != nil)
}

func TestStaticAgency(t *testing.T) {
	nodes := newTestNode()
	vds := newValidators(nodes, 0)

	agency := NewStaticAgency(nodes)
	validators, err := agency.GetValidator(0)
	assert.True(t, err == nil)
	assert.Equal(t, *vds, *validators)
}

func genesisBlockForTesting(db ethdb.Database, addr common.Address, balance *big.Int) (*types.Block, *params.ChainConfig) {
	buf, err := ioutil.ReadFile("../../eth/downloader/testdata/platon.json")
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

	block, _ := gen.Commit(db)
	return block, gen.Config

}

func TestInnerAgency(t *testing.T) {
	testdb := ethdb.NewMemDatabase()
	balanceBytes, _ := hexutil.Decode("0x2000000000000000000000000000000000000000000000000000000000000")
	balance := big.NewInt(0)
	genesis, chainConfig := genesisBlockForTesting(testdb, testAddress, balance.SetBytes(balanceBytes))
	vmConfig := vm.Config{}

	blocks, receipts := core.GenerateChain(chainConfig, genesis, new(consensus.BftMock), testdb, 60, func(i int, block *core.BlockGen) {
		block.SetCoinbase(common.Address{1})

		gas := big.NewInt(0)
		gas = gas.SetBytes(hexutil.MustDecode("0x99988888"))
		gasPrice := big.NewInt(0)
		gasPrice = gasPrice.SetBytes(hexutil.MustDecode("0x8250"))

		if i%3 == 0 {
			signer := types.MakeSigner(chainConfig, block.Number())
			tx, err := types.SignTx(types.NewTransaction(block.TxNonce(testAddress), common.HexToAddress("0x0384d39b9cbf9bab2a3b41692d426ad57e41c54c"), big.NewInt(1000), gas.Uint64(), gasPrice, hexutil.MustDecode("0xd3880000000000000002857072696e7483616263")), signer, testKey)
			if err != nil {
				panic(err)
			}
			block.AddTx(tx)
		}
	})
	assert.True(t, len(blocks) == 60)
	assert.True(t, len(receipts) == 60)

	for i, block := range blocks {
		rawdb.WriteBlock(testdb, block)
		rawdb.WriteReceipts(testdb, block.Hash(), block.NumberU64(), receipts[i])
		rawdb.WriteCanonicalHash(testdb, block.Hash(), block.NumberU64())
		rawdb.WriteHeadBlockHash(testdb, block.Hash())
		rawdb.WriteHeadHeaderHash(testdb, block.Hash())
	}

	blockchain, _ := core.NewBlockChain(testdb, nil, chainConfig, new(consensus.BftMock), vmConfig, func(*types.Block) bool {
		return true
	})

	nodes := newTestNode()
	vds := newValidators(nodes, 0)

	agency := NewInnerAgency(nodes, blockchain, 10, 20)

	assert.True(t, agency.GetLastNumber(0) == 40)

	validators, err := agency.GetValidator(0)
	assert.True(t, err == nil)
	assert.Equal(t, *vds, *validators)
	assert.True(t, blockchain.Genesis() != nil)
	block := blockchain.GetBlockByNumber(20)
	assert.True(t, block != nil)
	state, err := blockchain.StateAt(block.Root())

	cvds := &vm.Validators{
		ValidateNodes:    make(vm.NodeList, 0),
		ValidBlockNumber: 41,
	}
	for k, v := range validators.Nodes {
		cvds.ValidateNodes = append(cvds.ValidateNodes, &vm.ValidateNode{
			Index:   uint(v.Index),
			Address: v.Address,
			NodeID:  k,
		})
	}

	b, _ := rlp.EncodeToBytes(cvds)
	state.SetState(vm.ValidatorInnerContractAddr, []byte(vm.CurrentValidatorKey), b)
	state.Commit(false)

	newVds, err := agency.GetValidator(40)
	assert.True(t, err == nil)
	assert.True(t, newVds.Len() == 4)
	assert.True(t, newVds.ValidBlockNumber == 41)

	id3 := newVds.NodeID(3)
	assert.Equal(t, id3, nodes[3].ID)
	assert.True(t, agency.GetLastNumber(41) == 80)
}
