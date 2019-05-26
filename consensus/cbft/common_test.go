package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)


var (
	chainConfig = params.TestnetChainConfig
	testTxPoolConfig = core.DefaultTxPoolConfig
)

func init()  {
	chainConfig.Cbft.Period = 1
	chainConfig.Cbft.Epoch = 1
	chainConfig.Cbft.MaxLatency = 1
	chainConfig.Cbft.Duration = 10
	testTxPoolConfig.Journal = ""

}
func CreateCBFT(path string) *Cbft {
	ctx := node.NewServiceContext(&node.Config{DataDir:path}, nil, new(event.TypeMux), nil)

	return New(chainConfig.Cbft, ctx.EventMux, ctx)
}

type testBackend struct {
	db ethdb.Database
	chain *core.BlockChain
	txpool *core.TxPool
	cache *core.BlockChainCache
	cbft *Cbft
	worker *mockWorker
}

type mockWorker struct {
	mux *event.TypeMux
}

func NewMockWorker (mux *event.TypeMux) *mockWorker {
	m := &mockWorker{mux:mux}
	go m.loop()
	return m
}

func (m *mockWorker) loop()  {
	sub :=m.mux.Subscribe(cbfttypes.CbftResult{})
	for {
		select {
			case   <-sub.Chan():

		}
	}
}


func CreateBackend(engine *Cbft) * testBackend{
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

	engine.SetBlockChainCache(cache)
	txpool := core.NewTxPool(testTxPoolConfig, chainConfig, core.NewTxPoolBlockChain(cache))

	engine.Start(chain, txpool,NewStaticAgency(nil) )

	return &testBackend{
		db:db,
		chain:chain,
		cache:cache,
		txpool:txpool,
		worker:NewMockWorker(engine.eventMux),
	}
}



func TestNewCBFT(t *testing.T)  {
	path := path()
	defer os.RemoveAll(path)

	engine :=CreateCBFT(path)
	assert.NotNil(t, engine)
	backend := CreateBackend(engine)
	assert.NotNil(t, backend)
	assert.NotNil(t, backend.cache)

}