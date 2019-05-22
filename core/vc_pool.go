package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	//"github.com/PlatONnetwork/PlatON-Go/accounts"
	//"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/life/utils"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var VC_POOL *VCPool

const (
	TX_VC = TX_MPC + 1
)

type VCBlockChain interface {
	CurrentBlock() *types.Block
	GetBlock(hash common.Hash, number uint64) *types.Block

	SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription
}

type VCPoolConfig struct {
	VCEnable    bool          // the switch of vc compute
	NoLocals    bool          // Whether local transaction handling should be disabled
	Journal     string        // Journal of local transactions to survive node restarts
	Rejournal   time.Duration // Time interval to regenerate the local transaction journal
	GlobalQueue uint64        // Maximum number of non-executable transaction slots for all accounts
	Lifetime    time.Duration // Maximum amount of time non-executable transaction are queued

	LocalRpcPort int            // LocalRpcPort of local rpc port
	IceConf      string         // ice conf to init vm
	VcActor      common.Address // the actor of VC compute
	VcPassword   string         // the actor of VC compute
}

var DefaultVCPoolConfig = VCPoolConfig{
	Journal:     "VC_transactions.rlp",
	Rejournal:   time.Second * 4,
	GlobalQueue: 1024,
	Lifetime:    3 * time.Hour,
}

type VCPool struct {
	config      VCPoolConfig
	chainconfig *params.ChainConfig
	chain       *BlockChain

	mu        sync.RWMutex
	all       *vcLookup // All transactions to allow lookups
	queue     *vcList   // All transactions sorted by price
	work_pool *GoroutinePool
	quiteSign chan interface{}

	wg sync.WaitGroup // for shutdown sync
}

func NewVCPool(config VCPoolConfig, chainconfig *params.ChainConfig, chain *BlockChain) *VCPool {
	pool := &VCPool{
		config:      config,
		chainconfig: chainconfig,
		chain:       chain,
		all:         newVCLookup(),
		work_pool:   new(GoroutinePool),
	}

	pool.queue = newVCList(pool.all)
	pool.work_pool.Init(3, 10)

	pool.wg.Add(1)
	go pool.loop()

	//start work pool
	go pool.work_pool.Start()

	// save to global attr
	VC_POOL = pool

	return pool
}

func genCompInput(taskid string) string {
	var input [][]byte
	input = make([][]byte, 0)
	log.Debug("VC genCompInput", "taskid", taskid)
	input = append(input, utils.Int64ToBytes(2))
	input = append(input, []byte("real_compute"))
	input = append(input, utils.String2bytes(taskid))

	// var result []byte
	// input = append(input, result)
	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	if err != nil {
		fmt.Println("geninput fail.", err)
	}
	return common.Bytes2Hex(buffer.Bytes())
}

func genSetResultInput(taskid string, result []byte) string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, utils.Int64ToBytes(2))
	input = append(input, []byte("set_result"))
	input = append(input, utils.String2bytes(taskid))
	input = append(input, result)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	if err != nil {
		fmt.Println("geninput fail.", err)
		return ""
	}

	//log.Debug("VC genSetResultInput", "result", common.Bytes2Hex(buffer.Bytes()))
	return common.Bytes2Hex(buffer.Bytes())
}

//send POST
func (pool *VCPool) Post(url string, data interface{}, contentType string) (content string, err error) {
	jsonStr, _ := json.Marshal(data)
	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return content, nil
}

func (pool *VCPool) real_compute(tx *types.TransactionWrap) error {
	pool.mu.Lock()
	signer := types.MakeSigner(pool.chainconfig, big.NewInt(int64(tx.Bn)))
	caller, _, err := signer.SignatureAndSender(tx.Transaction)
	if err != nil {
		log.Warn("Get sig fail", "hash", tx.Hash())
		return err
	}

	bc := pool.chain
	header := pool.chain.CurrentHeader()

	state, err := bc.State()
	if err != nil {
		return err
	}

	gp := new(GasPool).AddGas(math.MaxUint64)

	input := genCompInput(tx.TaskId)
	msg := types.NewMessage(caller, tx.To(), 0, new(big.Int).SetInt64(0),
		tx.Gas(), tx.GasPrice(), common.Hex2Bytes(input), false)
	context := NewEVMContext(msg, header, bc, nil)
	evm := vm.NewEVM(context, state, bc.chainConfig, bc.vmConfig)
	log.Debug("start evm call real_compute")
	ret, _, vmerr, err := ApplyMessage(evm, msg, gp)
	if vmerr || err != nil {
		log.Error("ApplyMessage real_compute wrong ", "vmerr", vmerr, "error", err)
		return err
	}
	pool.mu.Unlock()

	//fmt.Println(ret)
	//TODO rm 64 bytes messy code
	if len(ret) < 64 {
		log.Error("ApplyMessage real_compute return error ")
		return fmt.Errorf("return len error")
	}

	res := ret[64:len(ret)]
	//fmt.Println(string(res))
	//var a accounts.Account
	fmt.Println("unlock: ", pool.config.VcActor.Hex())
	//fmt.Println("password ", pool.config.VcPassword)
	data := genSetResultInput(tx.TaskId, bytes.TrimLeft(res, "\x00"))
	//a.Address = pool.config.VcActor
	//ks := keystore.NewKeyStore(filepath.Join("./build/bin/data/", "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)
	//ks.Unlock(a, pool.config.VcPassword)

	vc_data := make(map[string]interface{})
	vc_data["jsonrpc"] = "2.0"
	vc_data["method"] = "eth_sendTransaction"
	params := make([]map[string]interface{}, 1)
	param := make(map[string]interface{})
	param["from"] = pool.config.VcActor.Hex()
	param["to"] = (*(tx.To())).Hex()
	param["gas"] = "0x166709"
	param["gasPrice"] = "0x8250de00"
	param["value"] = "0x0"
	param["data"] = "0x" + data
	params[0] = param
	vc_data["params"] = params
	vc_data["id"] = 1

	url := "http://127.0.0.1:" + strconv.FormatUint(uint64(pool.config.LocalRpcPort), 10)
	format := "application/json"
	postres, err := pool.Post(url, vc_data, format)
	if err != nil {
		return err
	}
	fmt.Println(strings.ToLower(string(postres)))
	//log.Debug("start evm call over")
	return nil
}

func (pool *VCPool) loop() {

	defer pool.wg.Done()

	var prevQueued int

	report := time.NewTicker(statsReportInterval)
	defer report.Stop()

	pop := time.NewTicker(time.Second * 1)

	// Keep waiting for and reacting to the various events
	for {
		select {
		case <-report.C:
			pool.mu.RLock()
			_, queued := pool.stats()
			pool.mu.RUnlock()
			if queued != prevQueued {
				log.Debug("VC transaction pool status report", "queued", queued)
				prevQueued = queued
			}
			//debug.FreeOSMemory()

		case <-pool.quiteSign:
			return

		case <-pop.C:
			if pool.queue.items.Len() > 0 {

				pool.mu.Lock()
				tx := pool.queue.Pop()
				bn := pool.chain.CurrentBlock().Number().Int64()

				// The latest block is smaller than the current store,
				// indicating that a fork has occurred and the transaction is removed.
				//
				// If the block to which the transaction belongs is not separated
				// from the latest block by 20 confirmation blocks, no processing is performed.
				log.Debug("start evm ------------------------------------------------call ", "bn", bn)
				if bn < int64(tx.Bn) || (bn-int64(tx.Bn)) >= MinBlockConfirms {
					pool.all.Remove(tx.Hash())
				} else {
					pool.queue.Put(tx)
				}
				pool.mu.Unlock()

				if (bn - int64(tx.Bn)) >= MinBlockConfirms {
					pool.work_pool.AddTask(func() error {
						return pool.real_compute(tx)
					})
				}

			}
		}
	}
}

func (pool *VCPool) LoadActor() error {
	absPath, err := filepath.Abs(DEFAULT_ACTOR_FILE_NAME)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil
	}
	res, _ := ioutil.ReadFile(absPath)
	pool.config.VcActor = common.BytesToAddress(res)
	fmt.Println(pool.config.VcActor.Hex())
	return nil
}

// Stop terminates the VC transaction pool.
func (pool *VCPool) Stop() {
	pool.quiteSign <- true
	pool.work_pool.Stop()
	pool.wg.Wait()

	log.Info("Transaction pool stopped")
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (pool *VCPool) SubscribeNewTxsEvent(ch chan<- NewTxsEvent) event.Subscription {
	//return pool.scope.Track(pool.txFeed.Subscribe(ch))
	return nil
}

func (pool *VCPool) AddLocals(txs []*types.TransactionWrap) []error {
	return pool.addTxs(txs)
}

func (pool *VCPool) addTxs(txs []*types.TransactionWrap) []error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.addTxsLocked(txs)
}

func (pool *VCPool) addTxsLocked(txs []*types.TransactionWrap) []error {
	errs := make([]error, len(txs))
	for i, tx := range txs {
		var replace bool
		if replace, errs[i] = pool.add(tx); errs[i] == nil && !replace {
			log.Warn("load txs from rlp fail", "hash", tx.Hash())
		}
	}
	return errs
}

func (pool *VCPool) Stats() (int, int) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.stats()
}

func (pool *VCPool) stats() (int, int) {
	return 0, len(pool.queue.all.txs)
}

func (pool *VCPool) local() types.TransactionWraps {
	VCTxs := make([]*types.TransactionWrap, 0)
	for _, v := range pool.queue.all.txs {
		VCTxs = append(VCTxs, v)
	}
	return VCTxs
}

func (pool *VCPool) validateTx(tx *types.TransactionWrap) (err error) {
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
		txType   int
		funcName string
	)
	if v, ok := iRlpList[0].([]byte); ok {
		txType = int(common.BytesToInt64(v))
	}

	if txType != TX_VC {
		return fmt.Errorf("Invalid tx type")
	}
	if v, ok := iRlpList[1].([]byte); ok {
		funcName = string(v)
	}
	tx.FuncName = funcName
	return nil
}

func (pool *VCPool) InjectTxs(block *types.Block, receipts types.Receipts, bc *BlockChain, state *state.StateDB) {
	if !pool.config.VCEnable {
		log.Info("Wow ~ VC Disable...")
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
					if tid, err := vc_verifyStartCalcLogs(receipt.Logs); err == nil {
						taskId = tid
						isSave = true
						break
					}
				}
			}
		}
		if isSave {
			wrap := &types.TransactionWrap{
				Transaction: tx,
				Bn:          block.NumberU64(),
				TaskId:      taskId,
			}

			// basic validatel
			if err := pool.validateTx(wrap); err != nil {
				log.Trace("God ~ Discarding invalid VC transaction", "hash", wrap.Hash(), "err", err)
				return
			}
			// actor validate
			// if err := pool.validateActor(wrap, bc, state); err != nil {
			// 	log.Trace("God ~ Discarding the actor not belong to current VC contract.", "hash", wrap.Hash(), "err", err)
			// 	return
			log.Debug("Wow ~ VC add pool--------------------------------------...")
			// }
			pool.add(wrap)
		}
	}
}

func (pool *VCPool) add(tx *types.TransactionWrap) (bool, error) {

	hash := tx.Hash()
	if pool.all.Get(hash) != nil {
		log.Trace("God ~ Discarding already known transaction", "hash", hash)
		return false, fmt.Errorf("known VC transaction: %x", hash)
	}

	// If the transaction fails basic validation, discard it
	if err := pool.validateTx(tx); err != nil {
		log.Trace("God ~ Discarding invalid VC transaction", "hash", hash, "err", err)
		return false, err
	}

	replace, err := pool.enqueueTx(hash, tx)
	if err != nil {
		return false, err
	}

	pool.journalTx(tx)

	log.Debug("Pooled new future VC transaction", "hash", hash, "to", tx.To())
	return replace, nil
}

func (pool *VCPool) enqueueTx(hash common.Hash, tx *types.TransactionWrap) (bool, error) {
	if pool.all.Get(hash) == nil {
		pool.all.Add(tx)
		pool.queue.items.Push(tx)
	}
	return true, nil
}

// wirte in to file
func (pool *VCPool) journalTx(tx *types.TransactionWrap) {

}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *VCPool) addTx(tx *types.TransactionWrap) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Try to inject the transaction and update any state
	_, err := pool.add(tx)
	if err != nil {
		return err
	}
	return nil
}

type vcLookup struct {
	txs  map[common.Hash]*types.TransactionWrap
	lock sync.RWMutex
}

// newTxLookup returns a new VCLookup structure.
func newVCLookup() *vcLookup {
	return &vcLookup{
		txs: make(map[common.Hash]*types.TransactionWrap),
	}
}

// Range calls f on each key and value present in the map.
func (t *vcLookup) Range(f func(hash common.Hash, tx *types.TransactionWrap) bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for key, value := range t.txs {
		if !f(key, value) {
			break
		}
	}
}

// Get returns a transaction if it exists in the lookup, or nil if not found.
func (t *vcLookup) Get(hash common.Hash) *types.TransactionWrap {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.txs[hash]
}

// Count returns the current number of items in the lookup.
func (t *vcLookup) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return len(t.txs)
}

// Add adds a transaction to the lookup.
func (t *vcLookup) Add(tx *types.TransactionWrap) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.txs[tx.Hash()] = tx
}

// Remove removes a transaction from the lookup.
func (t *vcLookup) Remove(hash common.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.txs, hash)
}

// start_calc_event -> sha3("start_calc_event")
func vc_verifyStartCalcLogs(logs []*types.Log) (string, error) {
	topic := crypto.Keccak256([]byte("start_compute_event"))
	for _, log := range logs {
		if len(log.Topics) == 0 {
			return "", fmt.Errorf("Reason: %v", "No topic found")
		}
		//log.Trace(("start evm call")
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
					code   uint64
					taskId string
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
