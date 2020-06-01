package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	fromAccountList      []*account
	toAccountList        []*account
	contractAccountList  []*account
	testTxList           types.Transactions
	accountCount         = 100
	contractAccountCount = 100
	txCount              = 20
	balance              = int64(730000000000000000)
	blockGasLimit        = uint64(500000000000)

	gasPrice = big.NewInt(15000)
	signer   = types.NewEIP155Signer(chainConfig.ChainID)

	blockchain *BlockChain
	stateDb    *state2.StateDB
	block      *types.Block
	header     *types.Header
)

type account struct {
	priKey  *ecdsa.PrivateKey
	address common.Address
	nonce   uint64
}

func initAccount() {
	contractAccountList = make([]*account, contractAccountCount)
	for i := 0; i < contractAccountCount; i++ {
		contractAccountKey, _ := crypto.GenerateKey()
		contractAccount := &account{}
		contractAccount.priKey = contractAccountKey
		contractAccount.address = crypto.PubkeyToAddress(contractAccountKey.PublicKey)
		contractAccount.nonce = 0
		contractAccountList[i] = contractAccount
	}

	fromAccountList = make([]*account, accountCount)
	toAccountList = make([]*account, accountCount)

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
}

func initTx() {
	testTxList = make(types.Transactions, txCount)
	for i := 0; i < txCount; i++ {
		fromAccount := fromAccountList[rand.Intn(accountCount)]
		//toAccount := toAccountList[rand.Intn(accountCount)]
		toAccount := fromAccountList[rand.Intn(accountCount)] //it is possible from=to
		//contractAccount := contractAccountList[rand.Intn(contractAccountCount)]
		contractAccount := contractAccountList[1]
		txType := rand.Intn(4)
		var tx *types.Transaction
		if txType == 0 || txType == 1 {
			tx, _ = types.SignTx(
				types.NewTransaction(
					fromAccount.nonce,
					toAccount.address,
					//contractAccount.address,
					big.NewInt(1),
					//21000+9000+320000, // it is short.
					//21000+9000+320000+21000, // it is enough.
					uint64(rand.Intn(200000000)),
					gasPrice,
					hexutil.MustDecode("0xf853838207d0b842b84006463ca71944647572a3ffcf96ab229f2e607651a40d89ff3ec36759fbc62b9f72ba1c07a9a6de87f61ec0e9574ebe338914da0931f1701a8bba3ca4162c23378a89746578745049504944")),
				signer,
				fromAccount.priKey)

		} else if txType == 2 {
			tx, _ = types.SignTx(
				types.NewTransaction(
					fromAccount.nonce,
					//toAccount.address,
					contractAccount.address,
					big.NewInt(0),
					//21000+9000+320000, // it is short.
					//21000+9000+320000+21000, // it is enough.
					uint64(rand.Intn(200000000)),
					gasPrice,
					hexutil.MustDecode("0xf853838207d0b842b84006463ca71944647572a3ffcf96ab229f2e607651a40d89ff3ec36759fbc62b9f72ba1c07a9a6de87f61ec0e9574ebe338914da0931f1701a8bba3ca4162c23378a89746578745049504944")),
				signer,
				fromAccount.priKey)
		} else if txType == 3 {
			tx, _ = types.SignTx(
				types.NewTransaction(
					fromAccount.nonce,
					//toAccount.address,
					//vm.GovContractAddr,
					contractAccount.address,
					big.NewInt(0),
					//21000+9000+320000, // it is short.
					//21000+9000+320000+21000, // it is enough.
					uint64(rand.Intn(200000000)),
					big.NewInt(1500000),
					hexutil.MustDecode("0xf853838207d0b842b84006463ca71944647572a3ffcf96ab229f2e607651a40d89ff3ec36759fbc62b9f72ba1c07a9a6de87f61ec0e9574ebe338914da0931f1701a8bba3ca4162c23378a89746578745049504944")),
				signer,
				fromAccount.priKey)
		}
		types.Sender(signer, tx)

		testTxList[i] = tx
		fromAccount.nonce++
	}

	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

}

func initChain() {
	db := ethdb.NewMemDatabase()
	stateDb, _ = state2.New(common.Hash{}, state2.NewDatabase(db))

	stateDb.SetBalance(nodeID, big.NewInt(0))
	for i := 0; i < accountCount; i++ {
		stateDb.SetBalance(fromAccountList[i].address, big.NewInt(balance))
		stateDb.SetBalance(toAccountList[i].address, big.NewInt(balance))
	}

	for i := 0; i < contractAccountCount; i++ {
		stateDb.SetCode(contractAccountList[i].address, hexutil.MustDecode("0xf87303843b9aca008347e7c494e80cbe05d8b7de0b8b2e436deda5ea6a70e4bf90808ecd888f9af6c4c62d90d8830186a081eca09610dd9c17164e5675e593c1b7b59aa865a2f120ecc0287538cf18ba05d76a78a07f49582e1850d7cff2bad12b5acd333cd8d50f25eccc04fd896c93b97b75f66a"))
	}

	stateDb.Finalise(false)

	gspec := &Genesis{
		Config: chainConfig,
		//Alloc:  GenesisAlloc{address: {Balance: funds}},
	}
	gspec.MustCommit(db)

	vmConfig := cvm.Config{
		ConsoleOutput: true,
		WasmType:      cvm.Wagon,
	}

	blockchain, _ = NewBlockChain(db, nil, gspec.Config, engine, vmConfig, nil)

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

/*func TestMain(m *testing.M) {
	initAccount()
	initChain()
	initTx()

	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	exitCode := m.Run()
	os.Exit(exitCode)
}*/

func TestParallel_PackParallel_VerifyParallel(t *testing.T) {
	initAccount()
	initChain()
	initTx()

	blockchain.SetProcessor(NewParallelStateProcessor(chainConfig, blockchain, engine))
	parallelMode(t)
}

func TestParallel_PackParallel_VerifySerial(t *testing.T) {
	initAccount()
	initChain()
	initTx()

	blockchain.SetProcessor(NewStateProcessor(chainConfig, blockchain, engine))
	parallelMode(t)
	t.Log("ok")
}

func parallelMode(t testing.TB) {
	initState := stateDb.Copy()
	NewExecutor(chainConfig, blockchain, blockchain.vmConfig)

	start := time.Now()
	gp := new(GasPool).AddGas(header.GasLimit)
	ctx := NewParallelContext(stateDb, header, common.Hash{}, gp, true, nil)
	ctx.SetBlockDeadline(time.Now().Add(200 * time.Second))
	ctx.SetBlockGasUsedHolder(&header.GasUsed)
	ctx.SetTxList(testTxList)

	if err := GetExecutor().ExecuteTransactions(ctx); err != nil {
		t.Fatal("pack txs err", "err", err)
	}
	end := time.Now()
	executeTxsCost := end.Sub(start).Nanoseconds()
	t.Logf("Executed txs cost(parallel mode, including to make DAG graph): %d Nanoseconds.\n", executeTxsCost)

	finalizedBlock, err := Finalize(blockchain, header, stateDb, ctx.packedTxList, ctx.receipts)
	if err != nil {
		t.Fatal("Finalize block failed", "err", err)
	}

	t.Log("======================")
	t.Logf("Finalize block cost(parallel mode): %d Nanoseconds.\n", time.Now().Sub(end).Nanoseconds())
	t.Log("======================")

	if sealedBlock, err := Seal(blockchain, finalizedBlock); err != nil {
		t.Fatal("Seal block failed", "err", err)
	} else {
		//fmt.Println(fmt.Sprintf("total txs=%d", len(sealedBlock.Transactions())))
		if _, err := blockchain.ProcessDirectly(sealedBlock, initState, blockchain.Genesis()); err != nil {
			t.Fatal("ProcessDirectly block error", "err", err)
		}
	}
}

func TestParallel_PackSerial_VerifyParallel(t *testing.T) {
	initAccount()
	initChain()
	initTx()

	blockchain.SetProcessor(NewParallelStateProcessor(chainConfig, blockchain, engine))
	serialMode(t)
}

func TestParallel_PackSerial_VerifySerial(t *testing.T) {
	initAccount()
	initChain()
	initTx()

	blockchain.SetProcessor(NewStateProcessor(chainConfig, blockchain, engine))
	serialMode(t)
}

func serialMode(t testing.TB) {
	initState := stateDb.Copy()
	NewExecutor(chainConfig, blockchain, blockchain.vmConfig)

	t.Logf("begin to executed txs cost(serial mode), blockGasLimit=%d, blockGasUsed=%d \n", header.GasLimit, header.GasUsed)

	gp := new(GasPool).AddGas(header.GasLimit)
	start := time.Now()
	txs := types.Transactions{}
	var receipts = types.Receipts{}
	for idx, tx := range testTxList {
		stateDb.Prepare(tx.Hash(), common.Hash{}, idx)
		t.Logf("idx=%d, txHash=%s, gasPool=%d, txGasLimit=%d", idx, tx.Hash().Hex(), gp.Gas(), tx.Gas())
		preUsed := header.GasUsed
		receipt, _, err := ApplyTransaction(chainConfig, blockchain, gp, stateDb, header, tx, &header.GasUsed, blockchain.vmConfig)
		used := header.GasUsed - preUsed
		t.Logf("txHash=%s, gasPool=%d, txGasLimit=%d, gasUsed=%d", tx.Hash().Hex(), gp.Gas(), tx.Gas(), used)

		if err != nil {
			t.Logf("apply tx error, err:%v", err)
			continue
		}
		receipts = append(receipts, receipt)
		txs = append(txs, tx)
	}
	end := time.Now()
	executeTxsCost := end.Sub(start).Nanoseconds()
	t.Logf("Executed txs cost(serial mode): %d Nanoseconds, blockGasUsed: %d, txCount: %d \n", executeTxsCost, header.GasUsed, len(txs))

	finalizedBlock, err := Finalize(blockchain, header, stateDb, txs, receipts)

	if err != nil {
		t.Fatal("Finalize block failed", "err", err)
	}
	t.Logf("Finalize block cost(parallel mode): %d Nanoseconds.\n", time.Now().Sub(end).Nanoseconds())

	if sealedBlock, err := Seal(blockchain, finalizedBlock); err != nil {
		t.Fatal("Seal block failed", "err", err)
	} else {
		//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
		t.Logf("verify block, blockGasLimit=%d", sealedBlock.GasLimit())
		if _, err := blockchain.ProcessDirectly(sealedBlock, initState, blockchain.Genesis()); err != nil {
			t.Fatal("ProcessDirectly block error", "err", err)
		}
	}
}

func TestParallel_EstimateTransferIntrinsicGas(t *testing.T) {
	gas, err := EstimateTransferIntrinsicGas([]byte("test"))
	if err != nil {
		t.Fatal("error", err)
	} else {
		t.Log("gas", gas)
	}

	gas2, err := EstimateTransferIntrinsicGas([]byte("test"))
	if err != nil {
		t.Fatal("error1", err)
	} else {
		t.Log("gas2", gas)
	}

	assert.Equal(t, gas, gas2)

}
