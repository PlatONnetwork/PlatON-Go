package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	state2 "github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	cvm "github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

var (
	chainConfig = params.TestnetChainConfig
	engine      = consensus.NewFaker()

	nodePriKey = crypto.HexMustToECDSA("1191dc5317d5930beb77848f416ee023921fa4452f4d783384f35352409c0ad0")
	nodeID     = crypto.PubkeyToAddress(nodePriKey.PublicKey)

	fromAccountList []*account
	toAccountList   []*account
	testTxList      types.Transactions
	accountCount    = 50
	txCount         = 200
	balance         = int64(7300000000000)
	blockGasLimit   = uint64(1000000000000000000)
	signer          = types.NewEIP155Signer(chainConfig.ChainID)

	blockchain         *BlockChain
	stateDb            *state2.StateDB
	block              *types.Block
	header             *types.Header
	presetTxExtendInfo = true
)

type account struct {
	priKey  *ecdsa.PrivateKey
	address common.Address
	nonce   uint64
}

func initTx() {
	fromAccountList = make([]*account, accountCount)
	toAccountList = make([]*account, accountCount)
	testTxList = make(types.Transactions, txCount)

	for i := 0; i < accountCount; i++ {
		fromKey, _ := crypto.GenerateKey()
		fromAccount := &account{}
		fromAccount.priKey = fromKey
		fromAccount.address = crypto.PubkeyToAddress(fromKey.PublicKey)
		fromAccount.nonce = 0
		fromAccountList[i] = fromAccount

		toKey, _ := crypto.GenerateKey()
		toAccount := &account{}
		toAccount.priKey = toKey
		toAccount.address = crypto.PubkeyToAddress(toKey.PublicKey)
		toAccount.nonce = 0
		toAccountList[i] = toAccount
	}

	for i := 0; i < txCount; i++ {
		fromAccount := fromAccountList[rand.Intn(accountCount)]
		toAccount := toAccountList[rand.Intn(accountCount)]
		var tx *types.Transaction
		if i%10 == 0 {
			tx, _ = types.SignTx(
				types.NewTransaction(
					fromAccount.nonce,
					vm.GovContractAddr,
					big.NewInt(1),
					//21000+9000+320000, // it is short.
					21000+9000+320000+21000, // it is enough.
					big.NewInt(1500000),
					hexutil.MustDecode("0xf853838207d0b842b84006463ca71944647572a3ffcf96ab229f2e607651a40d89ff3ec36759fbc62b9f72ba1c07a9a6de87f61ec0e9574ebe338914da0931f1701a8bba3ca4162c23378a89746578745049504944")),
				signer,
				fromAccount.priKey)
		} else {
			tx, _ = types.SignTx(
				types.NewTransaction(
					fromAccount.nonce,
					toAccount.address,
					big.NewInt(1),
					21000,
					big.NewInt(10),
					nil),
				signer,
				fromAccount.priKey)
		}

		from, _ := types.Sender(signer, tx)
		tx.SetFromAddr(&from)

		testTxList[i] = tx
		fromAccount.nonce++
	}
}

func initPrecompiledTx() {
	fromAccountList = make([]*account, accountCount)
	toAccountList = make([]*account, accountCount)
	testTxList = make(types.Transactions, txCount)

	for i := 0; i < accountCount; i++ {
		fromKey, _ := crypto.GenerateKey()
		fromAccount := &account{}
		fromAccount.priKey = fromKey
		fromAccount.address = crypto.PubkeyToAddress(fromKey.PublicKey)
		fromAccount.nonce = 0
		fromAccountList[i] = fromAccount

		toKey, _ := crypto.GenerateKey()
		toAccount := &account{}
		toAccount.priKey = toKey
		toAccount.address = crypto.PubkeyToAddress(toKey.PublicKey)
		toAccount.nonce = 0
		toAccountList[i] = toAccount
	}

	testTxList = make(types.Transactions, txCount)

	fromAccount := fromAccountList[rand.Intn(accountCount)]

	tx, _ := types.SignTx(
		types.NewTransaction(
			fromAccount.nonce,
			vm.GovContractAddr,
			big.NewInt(1),
			21000+9000+320000, // it is short.
			big.NewInt(1500000),
			hexutil.MustDecode("0xf853838207d0b842b84006463ca71944647572a3ffcf96ab229f2e607651a40d89ff3ec36759fbc62b9f72ba1c07a9a6de87f61ec0e9574ebe338914da0931f1701a8bba3ca4162c23378a89746578745049504944")),
		signer,
		fromAccount.priKey)
	fromAccount.nonce++
	testTxList[0] = tx

	fromAccount2 := fromAccountList[rand.Intn(accountCount)]
	toAccount := toAccountList[rand.Intn(accountCount)]
	tx2, _ := types.SignTx(
		types.NewTransaction(
			fromAccount2.nonce,
			toAccount.address,
			big.NewInt(1),
			21000,
			big.NewInt(10),
			nil),
		signer,
		fromAccount2.priKey)
	fromAccount2.nonce++
	testTxList[1] = tx2
}
func initChain() {
	db := ethdb.NewMemDatabase()
	stateDb, _ = state2.New(common.Hash{}, state2.NewDatabase(db))

	stateDb.SetBalance(nodeID, big.NewInt(0))
	for i := 0; i < accountCount; i++ {
		stateDb.SetBalance(fromAccountList[i].address, big.NewInt(balance))
		stateDb.SetBalance(toAccountList[i].address, big.NewInt(balance))
	}
	stateDb.Finalise(false)

	gspec := &Genesis{
		Config: chainConfig,
		//Alloc:  GenesisAlloc{address: {Balance: funds}},
	}
	gspec.MustCommit(db)

	blockchain, _ = NewBlockChain(db, nil, gspec.Config, engine, cvm.Config{}, nil)

	parent := blockchain.Genesis()
	block, header = NewBlock(parent.Hash(), parent.NumberU64()+1)
	header.Coinbase = nodeID
}

func NewBlock(parentHash common.Hash, number uint64) (*types.Block, *types.Header) {
	header := &types.Header{
		ParentHash: parentHash,
		Number:     big.NewInt(int64(number)),
		GasLimit:   blockGasLimit,
		Time:       big.NewInt(time.Now().UnixNano()),
		Coinbase:   nodeID,
	}

	Prepare(header)

	extraData := makeExtraData()
	copy(header.Extra[:len(extraData)], extraData)

	block := types.NewBlockWithHeader(header)
	return block, header
}

func makeExtraData() []byte {
	// create default extradata
	extra, _ := rlp.EncodeToBytes([]interface{}{
		uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
		"platon",
		runtime.Version(),
		runtime.GOOS,
	})
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}

	log.Debug("Prepare header extra", "data", hex.EncodeToString(extra))
	return extra
}

func Prepare(header *types.Header) error {
	if len(header.Extra) < 32 {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, 32-len(header.Extra))...)
	}
	header.Extra = header.Extra[:32]

	//init header.Extra[32: 32+65]
	header.Extra = append(header.Extra, make([]byte, consensus.ExtraSeal)...)
	return nil
}

func Finalize(chain consensus.ChainReader, header *types.Header, state *state2.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	header.Root = state.IntermediateRoot(true)
	return types.NewBlock(header, txs, receipts), nil
}

func signFn(m []byte) ([]byte, error) {
	return crypto.Sign(m, nodePriKey)
}

func Seal(chain consensus.ChainReader, block *types.Block) (sealedBlock *types.Block, err error) {
	header := block.Header()
	if block.NumberU64() == 0 {
		return nil, errors.New("unknown block")
	}

	sign, err := signFn(header.SealHash().Bytes())
	if err != nil {
		return nil, err
	}

	copy(header.Extra[len(header.Extra)-consensus.ExtraSeal:], sign[:])

	sealBlock := block.WithSeal(header)
	return sealBlock, nil
}

func TestMain(m *testing.M) {
	initTx()
	//initPrecompiledTx()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestParallel_PackParallel_VerifyParallel(t *testing.T) {
	//initTx()
	initChain()
	blockchain.SetProcessor(NewParallelStateProcessor(chainConfig, blockchain, engine))
	parallelMode(t)
}

func TestParallel_PackParallel_VerifySerial(t *testing.T) {
	//initTx()
	initChain()
	blockchain.SetProcessor(NewStateProcessor(chainConfig, blockchain, engine))
	parallelMode(t)
}

func parallelMode(t testing.TB) {
	initState := stateDb.Copy()

	NewExecutor(chainConfig, blockchain, cvm.Config{})

	start := time.Now()
	gp := new(GasPool).AddGas(1000000000000000000)
	ctx := NewPackBlockContext(stateDb, header, common.Hash{}, gp, time.Now().Add(200*time.Second))
	ctx.SetTxList(testTxList)

	if _, err := GetExecutor().PackBlockTxs(ctx); err != nil {
		t.Fatal("pack txs err", "err", err)
	}
	end := time.Now()
	executeTxsCost := end.Sub(start).Milliseconds()
	t.Logf("Executed txs cost(parallel mode, including to make DAG graph): %d milliseconds.\n", executeTxsCost)

	finalizedBlock, err := Finalize(blockchain, header, stateDb, ctx.packedTxList, ctx.receipts)
	if err != nil {
		t.Fatal("Finalize block failed", "err", err)
	}
	t.Logf("Finalize block cost(parallel mode): %d milliseconds.\n", time.Now().Sub(end).Milliseconds())

	if sealedBlock, err := Seal(blockchain, finalizedBlock); err != nil {
		t.Fatal("Seal block failed", "err", err)
	} else {
		if _, err := blockchain.ProcessDirectly(sealedBlock, initState, blockchain.Genesis()); err != nil {
			t.Fatal("ProcessDirectly block error", "err", err)
		}
	}
}

func TestParallel_PackSerial_VerifyParallel(t *testing.T) {
	//initTx()
	initChain()
	blockchain.SetProcessor(NewParallelStateProcessor(chainConfig, blockchain, engine))
	serialMode(t)
}

func TestParallel_PackSerial_VerifySerial(t *testing.T) {
	//initTx()
	initChain()
	blockchain.SetProcessor(NewStateProcessor(chainConfig, blockchain, engine))
	serialMode(t)
}

func serialMode(t testing.TB) {
	initState := stateDb.Copy()
	gp := new(GasPool).AddGas(1000000000000000000)
	start := time.Now()
	var receipts = types.Receipts{}
	for idx, tx := range testTxList {
		stateDb.Prepare(tx.Hash(), common.Hash{}, idx)
		receipt, _, err := ApplyTransaction(chainConfig, blockchain, gp, stateDb, header, tx, &header.GasUsed, cvm.Config{})

		if err != nil {
			t.Fatalf("apply tx error, err:%v", err)
		}
		receipts = append(receipts, receipt)
	}

	end := time.Now()
	executeTxsCost := end.Sub(start).Milliseconds()
	t.Logf("Executed txs cost(serial mode): %d milliseconds.\n", executeTxsCost)

	finalizedBlock, err := Finalize(blockchain, header, stateDb, testTxList, receipts)

	if err != nil {
		t.Fatal("Finalize block failed", "err", err)
	}
	t.Logf("Finalize block cost(parallel mode): %d milliseconds.\n", time.Now().Sub(end).Milliseconds())

	if sealedBlock, err := Seal(blockchain, finalizedBlock); err != nil {
		t.Fatal("Seal block failed", "err", err)
	} else {
		if _, err := blockchain.ProcessDirectly(sealedBlock, initState, blockchain.Genesis()); err != nil {
			t.Fatal("ProcessDirectly block error", "err", err)
		}
	}
}
