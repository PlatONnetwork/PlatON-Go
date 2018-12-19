package core

import (
	"github.com/PlatONnetwork/PlatON-Go/log"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

var mpcTestTxPoolConfig MPCPoolConfig

func init() {
	mpcTestTxPoolConfig = DefaultMPCPoolConfig
	mpcTestTxPoolConfig.Journal = "mpc.rlp"
}

type testMpcBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testMpcBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{
		GasLimit: bc.gasLimit,
	}, nil, nil, nil)
}

func (bc *testMpcBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testMpcBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testMpcBlockChain) SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

func mpcTransaction(taskId string, nonce uint64, gaslimit uint64, key *ecdsa.PrivateKey) *types.TransactionWrap {
	 tx := mpcPricedTransaction(nonce, gaslimit, big.NewInt(1), key)
	return &types.TransactionWrap{
		Transaction : tx,
		Bn: nonce,
		TaskId: taskId,
	}
}

func mpcPricedTransaction(nonce uint64, gaslimit uint64, gasprice *big.Int, key *ecdsa.PrivateKey) *types.Transaction {
	tx, _ := types.SignTx(types.NewTransaction(nonce, common.Address{}, big.NewInt(100), gaslimit, gasprice, nil), types.HomesteadSigner{}, key)
	return tx
}

// Configuring the MPCPool transaction pool.
func setupMpcPool() (*MPCPool, *ecdsa.PrivateKey) {
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	blockchain := &testMpcBlockChain{statedb, 1000000, new(event.Feed)}

	key, _ := crypto.GenerateKey()
	pool := NewMPCPool(mpcTestTxPoolConfig, params.TestChainConfig, blockchain)

	return pool, key
}

func TestMpcTransactionQueue(t *testing.T) {

	handler := log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stdout, log.FormatFunc(func(record *log.Record) []byte {
		return []byte(record.Msg + "\n")
	})))
	log.Root().SetHandler(handler)
	log.Debug("debug.........")
	t.Parallel()

	pool, key := setupMpcPool()
	defer pool.Stop()
	tx := mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd368",0, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd311",1, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd362",2, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd366",6, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd363",3, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd365",5, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd364",4, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd368",8, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd367",7, 100, key)
	pool.addTx(tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd369",9, 100, key)
	pool.addTx(tx)

	/*go func() {
		count := uint64(3)
		for {
			tx = transaction(count, 100, key)
			pool.addTx(tx)
			fmt.Println("gen tx , hash :", tx.Hash().Hex())
			//count++
			time.Sleep(time.Second * 3)
		}
	}()*/

}

func TestMpcTransactionJournaling(t *testing.T) {
	testMpcTransactionJournaling(t, false)
}

func TestMpcTransactionJournalingNoLocals(t *testing.T) {
	testMpcTransactionJournaling(t, true)
}

func testMpcTransactionJournaling(t *testing.T, nolocals bool) {

	t.Parallel()

	// Create a temporary file for the journal
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("failed to create temporary journal: %v", err)
	}
	journal := file.Name()
	defer os.Remove(journal)

	// Clean up the temporary file, we only need the path for now
	file.Close()
	os.Remove(journal)

	// Create the original pool to inject transaction into the journal
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	blockchain := &testMpcBlockChain{statedb, 1000000, new(event.Feed)}

	config := mpcTestTxPoolConfig
	config.NoLocals = nolocals
	config.Journal = journal
	config.Rejournal = time.Second

	pool := NewMPCPool(config, params.TestChainConfig, blockchain)

	key, _ := crypto.GenerateKey()
	tx := mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd360",0, 100, key)
	pool.enqueueTx(tx.Hash(), tx)

	tx = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd361",1, 100, key)
	pool.enqueueTx(tx.Hash(), tx)

	pending, queued := pool.Stats()
	if pending != 4 {
		t.Fatalf("pending transactions mismatched: have %d, want %d", pending, 4)
	}
	if queued != 0 {
		t.Fatalf("queued transactions mismatched: have %d, want %d", queued, 0)
	}
}

