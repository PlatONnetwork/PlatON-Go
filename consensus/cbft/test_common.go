package cbft

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress = crypto.PubkeyToAddress(testKey.PublicKey)

	chainConfig      = params.TestnetChainConfig
	testTxPoolConfig = core.DefaultTxPoolConfig
)

func NewBlock(parent common.Hash, number uint64) *types.Block {
	header := &types.Header{
		Number:      big.NewInt(int64(number)),
		ParentHash:  parent,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 77),
		ReceiptHash: common.BytesToHash(hexutil.MustDecode("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")),
		Root:        common.BytesToHash(hexutil.MustDecode("0x936fc1c230a4a6d65cece52afadcf64b3f622f31faef4fa87c7c0335eaf67c2f")),
		Coinbase:    common.Address{},
		GasLimit:    10000000000,
	}
	block := types.NewBlockWithHeader(header)
	return block
}

func GenerateKeys(num int) ([]*ecdsa.PrivateKey, []*bls.SecretKey) {
	pk := make([]*ecdsa.PrivateKey, 0)
	sk := make([]*bls.SecretKey, 0)

	for i := 0; i < num; i++ {
		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		ecdsaKey, _ := crypto.GenerateKey()
		pk = append(pk, ecdsaKey)
		sk = append(sk, &blsKey)
	}
	return pk, sk
}
func GenerateCbftNode(num int) ([]*ecdsa.PrivateKey, []*bls.SecretKey, []params.CbftNode) {
	pk, sk := GenerateKeys(num)
	nodes := make([]params.CbftNode, num)
	for i := 0; i < num; i++ {

		nodes[i].Node = *discover.NewNode(discover.PubkeyID(&pk[i].PublicKey), nil, 0, 0)
		nodes[i].BlsPubKey = *sk[i].GetPublicKey()
	}
	return pk, sk, nodes
}

func CreateCBFT(pk *ecdsa.PrivateKey, sk *bls.SecretKey, period uint64, amount uint32) *Cbft {

	sysConfig := &params.CbftConfig{
		Epoch:        1,
		Period:       10,
		Amount:       10,
		InitialNodes: []params.CbftNode{},
	}

	optConfig := &ctypes.OptionsConfig{
		NodePriKey: pk,
		NodeID:     discover.PubkeyID(&pk.PublicKey),
		BlsPriKey:  sk,
	}

	ctx := node.NewServiceContext(&node.Config{DataDir: ""}, nil, new(event.TypeMux), nil)

	return New(sysConfig, optConfig, ctx.EventMux, ctx)
}

func CreateBackend(engine *Cbft, nodes []params.CbftNode) (*core.BlockChain, *core.BlockChainCache, *core.TxPool, consensus.Agency) {

	var (
		db    = ethdb.NewMemDatabase()
		gspec = core.Genesis{
			Config: chainConfig,
			Alloc:  core.GenesisAlloc{},
		}
	)
	gspec.MustCommit(db)

	chain, _ := core.NewBlockChain(db, nil, gspec.Config, engine, vm.Config{}, nil)
	cache := core.NewBlockChainCache(chain)
	txpool := core.NewTxPool(testTxPoolConfig, chainConfig, cache)

	return chain, cache, txpool, validator.NewStaticAgency(nodes)
}

func CreateValidatorBackend(engine *Cbft, nodes []params.CbftNode) (*core.BlockChain, *core.BlockChainCache, *core.TxPool, consensus.Agency) {
	var (
		db    = ethdb.NewMemDatabase()
		gspec = core.Genesis{
			Config: chainConfig,
			Alloc:  core.GenesisAlloc{},
		}
	)
	balanceBytes, _ := hexutil.Decode("0x2000000000000000000000000000000000000000000000000000000000000")
	balance := big.NewInt(0)
	gspec.Alloc[testAddress] = core.GenesisAccount{
		Code:    nil,
		Storage: nil,
		Balance: balance.SetBytes(balanceBytes),
		Nonce:   0,
	}
	gspec.MustCommit(db)

	chain, _ := core.NewBlockChain(db, nil, gspec.Config, engine, vm.Config{}, nil)
	cache := core.NewBlockChainCache(chain)
	txpool := core.NewTxPool(testTxPoolConfig, chainConfig, cache)

	return chain, cache, txpool, validator.NewInnerAgency(nodes, chain, int(engine.config.Sys.Amount), int(engine.config.Sys.Amount)*2)
}

type TestCBFT struct {
	engine *Cbft
	chain  *core.BlockChain
	cache  *core.BlockChainCache
	txpool *core.TxPool
	agency consensus.Agency
}

func (t *TestCBFT) Start() error {
	return t.engine.Start(t.chain, t.cache, t.txpool, t.agency)
}

func MockNode(pk *ecdsa.PrivateKey, sk *bls.SecretKey, nodes []params.CbftNode, period uint64, amount uint32) *TestCBFT {
	engine := CreateCBFT(pk, sk, period, amount)

	chain, cache, txpool, agency := CreateBackend(engine, nodes)
	return &TestCBFT{
		engine: engine,
		chain:  chain,
		cache:  cache,
		txpool: txpool,
		agency: agency,
	}
}

func MockValidator(pk *ecdsa.PrivateKey, sk *bls.SecretKey, nodes []params.CbftNode, period uint64, amount uint32) *TestCBFT {
	engine := CreateCBFT(pk, sk, period, amount)

	chain, cache, txpool, agency := CreateValidatorBackend(engine, nodes)
	return &TestCBFT{
		engine: engine,
		chain:  chain,
		cache:  cache,
		txpool: txpool,
		agency: agency,
	}
}
