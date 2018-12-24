package core

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/mpc"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var MPC_POOL *MPCPool

const (
	MinBlockConfirms int64 = 20
	DEFAULT_ACTOR_FILE_NAME = ".actor"
)

const (
	TX_TRANSFER  = iota
	TX_CONTRACT_DEPLOY
	TX_CONTRACT_CALL
	TX_VOTE
	TX_PERMISSION
	TX_MPC
)

type mpcBlockChain interface {
	CurrentBlock() *types.Block
	GetBlock(hash common.Hash, number uint64) *types.Block

	SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription
}

type MPCPoolConfig struct {
	MPCEnable       bool 				// the switch of mpc compute

	NoLocals  		bool             	// Whether local transaction handling should be disabled
	Journal   		string           	// Journal of local transactions to survive node restarts
	Rejournal 		time.Duration    	// Time interval to regenerate the local transaction journal
	GlobalQueue  	uint64 				// Maximum number of non-executable transaction slots for all accounts
	Lifetime 		time.Duration 		// Maximum amount of time non-executable transaction are queued

	LocalRpcPort	int					// LocalRpcPort of local rpc port
	IceConf			string				// ice conf to init vm
	MpcActor		common.Address		// the actor of mpc compute
}

var DefaultMPCPoolConfig = MPCPoolConfig{
	Journal:   		"mpc_transactions.rlp",
	Rejournal: 		time.Second * 4,
	GlobalQueue:  	1024,
	Lifetime: 		3 * time.Hour,
}

type MPCPool struct {
	config       	MPCPoolConfig
	chainconfig  	*params.ChainConfig
	chain        	mpcBlockChain

	mu           	sync.RWMutex
	journal 		*mpcJournal  			// Journal of local mpc transaction to back up to disk
	all     		*mpcLookup              // All transactions to allow lookups
	queue  			*mpcList                // All transactions sorted by price

	quiteSign 		chan interface{}

	wg sync.WaitGroup 						// for shutdown sync
}

func NewMPCPool(config MPCPoolConfig, chainconfig *params.ChainConfig, chain blockChain) *MPCPool {
	pool := &MPCPool{
		config:      	config,
		chainconfig: 	chainconfig,
		chain:       	chain,
		all:       		newMpcLookup(),
	}

	pool.queue = newMpcList(pool.all)
	if config.Journal != "" {
		pool.journal = newMPCJournal(config.Journal)
		if err := pool.journal.load(pool.AddLocals); err != nil {
			log.Warn("Failed to load mpc transaction journal", "err", err)
		}
		if err := pool.journal.rotate(pool.local()); err != nil {
			log.Warn("Failed to rotate mpc transaction journal", "err", err)
		}
	}

	pool.wg.Add(1)
	go pool.loop()

	// save to global attr
	MPC_POOL = pool

	// init mvm
	url := "http://127.0.0.1:" + strconv.FormatUint(uint64(pool.config.LocalRpcPort), 10)
	mpc.InitVM(pool.config.IceConf, url)

	return pool
}

func (pool *MPCPool) loop() {

	defer pool.wg.Done()

	var prevQueued int

	report := time.NewTicker(statsReportInterval)
	defer report.Stop()

	journal := time.NewTicker(pool.config.Rejournal)
	defer journal.Stop()

	pop := time.NewTicker(time.Second * 1)

	// Keep waiting for and reacting to the various events
	for {
		select {
		case <- report.C:
			pool.mu.RLock()
			_, queued := pool.stats()
			pool.mu.RUnlock()
			if queued != prevQueued {
				log.Debug("Mpc transaction pool status report","queued", queued)
				prevQueued = queued
			}
			//debug.FreeOSMemory()

		// Handle local mpc transaction journal rotation
		case <-journal.C:
			if pool.journal != nil {
				pool.mu.Lock()
				if err := pool.journal.rotate(pool.local()); err != nil {
					log.Warn("Failed to rotate local tx journal", "err", err)
				}
				pool.mu.Unlock()
			}

		case <- pool.quiteSign:
			return

		case <- pop.C:
			if pool.queue.items.Len() != 0 {
				pool.mu.Lock()
				tx := pool.queue.Pop()
				bn := pool.chain.CurrentBlock().Number().Int64()
				if bn < int64(tx.Bn) || (bn - int64(tx.Bn)) >= MinBlockConfirms {
					pool.all.Remove(tx.Hash())
				} else {
					pool.queue.Put(tx)
				}
				pool.mu.Unlock()

				if (bn - int64(tx.Bn)) >= MinBlockConfirms {
					log.Trace("Wow ~ Received mpc transaction", "hash", tx.Hash()," bn", tx.Bn)
					singer := types.MakeSigner(pool.chainconfig, big.NewInt(int64(tx.Bn)))
					caller, pubKey, err := singer.SignatureAndSender(tx.Transaction)
					if err != nil {
						log.Warn("Get sig fail", "hash", tx.Hash())
					}
					log.Info("Recover pubKey success", "pubKey", common.Bytes2Hex(pubKey))
					mpc.ExecuteMPCTx(mpc.MPCParams{
						TaskId: 	tx.TaskId,
						Pubkey: 	common.Bytes2Hex(pubKey),
						From : 		caller,
						IRAddr :	*tx.To(),
						Method :	tx.FuncName,
						Extra :		"",
					})
				}
			}
		}
	}
}

func (pool *MPCPool) LoadActor() error {
	absPath, err := filepath.Abs(DEFAULT_ACTOR_FILE_NAME)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil
	}
	res, _ := ioutil.ReadFile(absPath)
	pool.config.MpcActor = common.BytesToAddress(res)
	fmt.Println(pool.config.MpcActor.Hex())
	return nil
}

// Stop terminates the mpc transaction pool.
func (pool *MPCPool) Stop() {
	pool.quiteSign <- true
	pool.wg.Wait()

	if pool.journal != nil {
		pool.journal.close()
	}
	log.Info("Transaction pool stopped")
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (pool *MPCPool) SubscribeNewTxsEvent(ch chan<- NewTxsEvent) event.Subscription {
	//return pool.scope.Track(pool.txFeed.Subscribe(ch))
	return nil
}

func (pool *MPCPool) AddLocals(txs []*types.TransactionWrap) []error {
	return pool.addTxs(txs)
}

func (pool *MPCPool) addTxs(txs []*types.TransactionWrap) []error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.addTxsLocked(txs)
}

func (pool *MPCPool) addTxsLocked(txs []*types.TransactionWrap) []error {
	errs := make([]error, len(txs))
	for i, tx := range txs {
		var replace bool
		if replace, errs[i] = pool.add(tx); errs[i] == nil && !replace {
			log.Warn("load txs from rlp fail", "hash", tx.Hash())
		}
	}
	return errs
}

func (pool *MPCPool) Stats() (int, int) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.stats()
}

func (pool *MPCPool) stats() (int, int) {
	return 0, len(pool.queue.all.txs)
}

func (pool *MPCPool) local() types.TransactionWraps {
	mpcTxs := make([]*types.TransactionWrap, 0)
	for _, v := range pool.queue.all.txs {
		mpcTxs = append(mpcTxs, v)
	}
	return mpcTxs
}

func (pool *MPCPool) validateActor(wrapTx *types.TransactionWrap, bc *BlockChain, state *state.StateDB) (err error) {

	header := bc.CurrentHeader()

	singer := types.MakeSigner(pool.chainconfig, big.NewInt(int64(wrapTx.Bn)))
	from, err := singer.Sender(wrapTx.Transaction)
	if err != nil {
		return fmt.Errorf("Retrive caller from transaction get failed : %v", err.Error())
	}

	input := "da880000000000000009906765745f7061727469636970616e7473"

	msg := types.NewMessage(from, wrapTx.To(), 0, new(big.Int).SetInt64(0), wrapTx.Gas(), wrapTx.GasPrice(), common.Hex2Bytes(input), false)
	context := NewEVMContext(msg, header, bc, nil)
	vm := vm.NewEVM(context, state, bc.chainConfig, bc.vmConfig)
	gp := new(GasPool).AddGas(math.MaxUint64)

	ret, _, _, err := ApplyMessage(vm, msg, gp)
	if err != nil {
		return fmt.Errorf("get call error: %v", err.Error())
	}
	//fmt.Println(strings.ToLower(string(ret)))

	if !strings.Contains(strings.ToLower(string(ret)), strings.ToLower(pool.config.MpcActor.String())) {
		return fmt.Errorf("%v", "Invalid caller")
	}
	return nil
}

func (pool *MPCPool) validateTx(tx *types.TransactionWrap) (err error) {
	input := tx.Data()
	if input == nil || len(input) <= 1 {
		return fmt.Errorf("Invalid input")
	}
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(input), &ptr)
	if err != nil {
		return err
	}
	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return fmt.Errorf("Invalid rlp encoded")
	}
	iRlpList := rlpList.([]interface{})
	if len(iRlpList) < 2 {
		return fmt.Errorf("Invalid input. ele must greater than 2")
	}
	var (
		txType 		int
		funcName 	string
	)
	if v, ok := iRlpList[0].([]byte); ok {
		txType = int(common.BytesToInt64(v))
	}
	if txType != TX_MPC {
		return fmt.Errorf("Invalid tx type")
	}
	if v, ok := iRlpList[1].([]byte); ok {
		funcName = string(v)
	}
	tx.FuncName = funcName
	return nil
}

func (pool *MPCPool) InjectTxs(block *types.Block, receipts types.Receipts, bc *BlockChain, state *state.StateDB) {
	if !pool.config.MPCEnable {
		log.Info("Wow ~ MPC Disable...")
		return
	}
	for _, tx := range block.Transactions() {
		isSave := false
		var taskId string
		for _, receipt := range receipts {
			if len(receipt.Logs) == 0 {
				continue
			}
			if bytes.Equal(receipt.TxHash.Bytes(), tx.Hash().Bytes()) { // fail
				if receipt.Status == 1 {
					// valid logs : error = 1 -> success
					if tid, err := verifyStartCalcLogs(receipt.Logs); err == nil {
						taskId = tid
						isSave = true
						break
					}
				}
			}
		}
		if isSave {
			wrap := &types.TransactionWrap {
				Transaction : tx,
				Bn: block.NumberU64(),
				TaskId: taskId,
			}
			// basic validate
			if err := pool.validateTx(wrap); err != nil {
				log.Trace("God ~ Discarding invalid mpc transaction", "hash", wrap.Hash(), "err", err)
				return
			}
			// actor validate
			if err := pool.validateActor(wrap, bc, state); err != nil {
				log.Debug("God ~ Discarding the actor not belong to current mpc contract.", "hash", wrap.Hash(), "err", err.Error())
				return
			}
			pool.add(wrap)
		}
	}
}

func (pool *MPCPool) add(tx *types.TransactionWrap) (bool, error) {

	hash := tx.Hash()
	if pool.all.Get(hash) != nil {
		log.Trace("God ~ Discarding already known transaction", "hash", hash)
		return false, fmt.Errorf("known mpc transaction: %x", hash)
	}

	// If the transaction fails basic validation, discard it
	if err := pool.validateTx(tx); err != nil {
		log.Trace("God ~ Discarding invalid mpc transaction", "hash", hash, "err", err)
		return false, err
	}

	replace, err := pool.enqueueTx(hash, tx)
	if err != nil {
		return false, err
	}

	pool.journalTx(tx)

	log.Trace("Pooled new future mpc transaction", "hash", hash, "to", tx.To())
	return replace, nil
}

func (pool *MPCPool) enqueueTx(hash common.Hash, tx *types.TransactionWrap) (bool, error) {
	if pool.all.Get(hash) == nil {
		pool.all.Add(tx)
		pool.queue.items.Push(tx)
	}
	return true, nil
}

// 写入文件
func (pool *MPCPool) journalTx(tx *types.TransactionWrap) {
	if pool.journal == nil {
		return
	}
	if err := pool.journal.insert(tx); err != nil {
		log.Warn("Failed to journal local transaction", "err", err)
	}
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *MPCPool) addTx(tx *types.TransactionWrap) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Try to inject the transaction and update any state
	_, err := pool.add(tx)
	if err != nil {
		return err
	}
	return nil
}

type mpcLookup struct {
	txs  		map[common.Hash]*types.TransactionWrap
	lock 		sync.RWMutex
}

// newTxLookup returns a new mpcLookup structure.
func newMpcLookup() *mpcLookup {
	return &mpcLookup{
		txs: make(map[common.Hash]*types.TransactionWrap),
	}
}

// Range calls f on each key and value present in the map.
func (t *mpcLookup) Range(f func(hash common.Hash, tx *types.TransactionWrap) bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for key, value := range t.txs {
		if !f(key, value) {
			break
		}
	}
}

// Get returns a transaction if it exists in the lookup, or nil if not found.
func (t *mpcLookup) Get(hash common.Hash) *types.TransactionWrap {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.txs[hash]
}

// Count returns the current number of items in the lookup.
func (t *mpcLookup) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return len(t.txs)
}

// Add adds a transaction to the lookup.
func (t *mpcLookup) Add(tx *types.TransactionWrap) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.txs[tx.Hash()] = tx
}

// Remove removes a transaction from the lookup.
func (t *mpcLookup) Remove(hash common.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.txs, hash)
}

// start_calc_event -> sha3("start_calc_event")
func verifyStartCalcLogs(logs []*types.Log) (string, error) {
	topic := crypto.Keccak256([]byte("start_calc_event"))
	for _, log := range logs {
		if len(log.Topics) == 0 {
			return "", fmt.Errorf("Reason: %v", "No topic found")
		}
		for _, top := range log.Topics {
			if bytes.EqualFold(topic, top.Bytes()) {
				// found valid log
				ptr := new(interface{})
				err := rlp.Decode(bytes.NewReader(log.Data), &ptr)
				if err != nil {
					return "", fmt.Errorf("Decode data of log got err: %v", err.Error())
				}
				rlpList := reflect.ValueOf(ptr).Elem().Interface()
				if _, ok := rlpList.([]interface{}); !ok {
					return "", fmt.Errorf("Reason: %v", "Invalid RLPList format")
				}
				iRlpList := rlpList.([]interface{})
				// [0] -> code [1] -> taskId
				var (
					code 	uint64
					taskId 	string
				)
				if v, ok := iRlpList[0].([]byte); ok {
					code = uint64(common.BytesToInt64(common.PaddingLeft(v, 8)))
				}
				if v, ok := iRlpList[1].([]byte); ok {
					taskId = string(v)
				}
				if code == 1 {
					return taskId, nil
				}
			}
		}
	}
	return "", fmt.Errorf("Invalid logs for event on topic : { %v }", topic)
}