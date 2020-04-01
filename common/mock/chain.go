package mock

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type Chain struct {
	Genesis *types.Block
	//	chain    []common.Hash
	//	headerm  map[common.Hash]*types.Header
	//	blockm   map[common.Hash]*types.Block
	//	receiptm map[common.Hash][]*types.Receipt
	StateDB          *MockStateDB
	SnapDB           snapshotdb.DB
	h                []*types.Header
	timeGenerate     func(*big.Int) *big.Int
	coinBaseGenerate func() common.Address
}

//notic AddBlock will not append snapshotdb
func (c *Chain) AddBlock() {
	header := generateHeader(new(big.Int).Add(c.h[len(c.h)-1].Number, common.Big1), c.h[len(c.h)-1].Hash(), c.timeGenerate(c.CurrentHeader().Time), c.coinBaseGenerate())
	c.h = append(c.h, header)
}

func (c *Chain) AddBlockWithTxHash(txHash common.Hash) {
	c.AddBlock()
	c.StateDB.Prepare(txHash, c.CurrentHeader().Hash(), 1)
}

func (c *Chain) SetHeaderTimeGenerate(f func(*big.Int) *big.Int) {
	c.timeGenerate = f
}

func (c *Chain) SetCoinbaseGenerate(f func() common.Address) {
	c.coinBaseGenerate = f
}

func (c *Chain) AddBlockWithTxHashAndCommit(txHash common.Hash, miner bool, f func(hash common.Hash, header *types.Header, sdb snapshotdb.DB) error) error {
	c.AddBlockWithTxHash(txHash)
	return c.commitWithSnapshotDB(miner, f)
}

func (c *Chain) commitWithSnapshotDB(miner bool, f func(hash common.Hash, header *types.Header, sdb snapshotdb.DB) error) error {
	var useHash common.Hash
	if miner {
		useHash = common.ZeroHash
	} else {
		useHash = c.CurrentHeader().Hash()
	}
	if err := c.SnapDB.NewBlock(c.CurrentHeader().Number, c.CurrentHeader().ParentHash, useHash); err != nil {
		return err
	}
	if err := f(useHash, c.CurrentHeader(), c.SnapDB); err != nil {
		return err
	}
	if miner {
		if err := c.SnapDB.Flush(c.CurrentHeader().Hash(), c.CurrentHeader().Number); err != nil {
			return err
		}
	}
	if err := c.SnapDB.Commit(c.CurrentHeader().Hash()); err != nil {
		return err
	}
	return nil
}

func (c *Chain) AddBlockWithSnapDB(miner bool, f func(hash common.Hash, header *types.Header, sdb snapshotdb.DB) error) error {
	c.AddBlock()
	return c.commitWithSnapshotDB(miner, f)
}

func (c *Chain) CurrentHeader() *types.Header {
	return c.h[len(c.h)-1]
}

func (c *Chain) CurrentForkHeader() *types.Header {
	newhead := new(types.Header)
	newhead.Number = c.h[len(c.h)-1].Number
	newhead.ParentHash = c.h[len(c.h)-1].ParentHash
	newhead.GasUsed = rand.Uint64()
	return newhead
}

func (c *Chain) GetHeaderByHash(hash common.Hash) *types.Header {
	for i := len(c.h) - 1; i >= 0; i-- {
		if c.h[i].Hash() == hash {
			return c.h[i]
		}
	}
	return nil
}

func (c *Chain) GetHeaderByNumber(number uint64) *types.Header {
	for i := len(c.h) - 1; i >= 0; i-- {
		if c.h[i].Number.Uint64() == number {
			return c.h[i]
		}
	}
	return nil
}

func NewChain() *Chain {
	c := new(Chain)

	c.timeGenerate = func(b *big.Int) *big.Int {
		return new(big.Int).SetInt64(time.Now().UnixNano() / 1e6)
	}
	c.coinBaseGenerate = func() common.Address {
		privateKey, err := crypto.GenerateKey()
		if nil != err {
			panic(err)
		}
		addr := crypto.PubkeyToAddress(privateKey.PublicKey)
		return addr
	}
	header := generateHeader(big.NewInt(0), common.ZeroHash, c.timeGenerate(nil), c.coinBaseGenerate())
	block := new(types.Block).WithSeal(header)

	c.Genesis = block
	c.h = make([]*types.Header, 0)
	c.h = append(c.h, header)

	c.StateDB = NewMockStateDB()
	c.SnapDB = snapshotdb.Instance()
	return c
}

func NewMockStateDB() *MockStateDB {
	db := new(MockStateDB)
	db.State = make(map[common.Address]map[string][]byte)
	db.Balance = make(map[common.Address]*big.Int)
	db.logs = make(map[common.Hash][]*types.Log)
	return db
}

type MockStateDB struct {
	Balance      map[common.Address]*big.Int
	State        map[common.Address]map[string][]byte
	thash, bhash common.Hash
	txIndex      int
	logs         map[common.Hash][]*types.Log
}

func (s *MockStateDB) Prepare(thash, bhash common.Hash, ti int) {
	s.thash = thash
	s.bhash = bhash
	s.txIndex = ti
}

func (s *MockStateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	return common.ZeroHash
}

func (s *MockStateDB) SubBalance(adr common.Address, amount *big.Int) {
	if balance, ok := s.Balance[adr]; ok {
		balance.Sub(balance, amount)
	}
}

func (s *MockStateDB) AddBalance(adr common.Address, amount *big.Int) {
	if balance, ok := s.Balance[adr]; ok {
		balance.Add(balance, amount)
	} else {
		s.Balance[adr] = amount
	}
}

func (s *MockStateDB) GetBalance(adr common.Address) *big.Int {
	if balance, ok := s.Balance[adr]; ok {
		return balance
	} else {
		return big.NewInt(0)
	}
}

func (s *MockStateDB) GetState(adr common.Address, key []byte) []byte {
	return s.State[adr][string(key)]
}

func (s *MockStateDB) SetState(adr common.Address, key, val []byte) {
	if len(val) == 0 {
		delete(s.State[adr], string(key))
	} else {
		if stateVal, ok := s.State[adr]; ok {
			stateVal[string(key)] = val
		} else {
			stateVal := make(map[string][]byte)
			stateVal[string(key)] = val
			s.State[adr] = stateVal
		}
	}
}

func (s *MockStateDB) CreateAccount(common.Address) {
	return
}

func (s *MockStateDB) GetNonce(common.Address) uint64 {
	return 0
}
func (s *MockStateDB) SetNonce(common.Address, uint64) {
	return
}

func (s *MockStateDB) GetCodeHash(common.Address) common.Hash {
	return common.ZeroHash
}
func (s *MockStateDB) GetCode(common.Address) []byte {
	return nil
}
func (s *MockStateDB) SetCode(common.Address, []byte) {
	return
}
func (s *MockStateDB) GetCodeSize(common.Address) int {
	return 0
}

func (s *MockStateDB) GetAbiHash(common.Address) common.Hash {
	return common.ZeroHash
}
func (s *MockStateDB) GetAbi(common.Address) []byte {
	return nil
}
func (s *MockStateDB) SetAbi(common.Address, []byte) {
	return
}

func (s *MockStateDB) AddRefund(uint64) {
	return
}
func (s *MockStateDB) SubRefund(uint64) {
	return
}
func (s *MockStateDB) GetRefund() uint64 {
	return 0
}

func (s *MockStateDB) GetCommittedState(common.Address, []byte) []byte {
	return nil
}

//GetState(common.Address, common.Hash) common.Hash
//SetState(common.Address, common.Hash, common.Hash)

func (s *MockStateDB) Suicide(common.Address) bool {
	return true
}
func (s *MockStateDB) HasSuicided(common.Address) bool {
	return true
}

// Exist reports whether the given account exists in state.
// Notably this should also return true for suicided accounts.
func (s *MockStateDB) Exist(common.Address) bool {
	return true
}

// Empty returns whether the given account is empty. Empty
// is defined according to EIP161 (balance = nonce = code = 0).
func (s *MockStateDB) Empty(common.Address) bool {
	return true
}

func (s *MockStateDB) RevertToSnapshot(int) {
	return
}
func (s *MockStateDB) Snapshot() int {
	return 0
}

func (s *MockStateDB) AddLog(l *types.Log) {
	if _, ok := s.logs[s.thash]; !ok {
		s.logs[s.thash] = make([]*types.Log, 0)
	}
	s.logs[s.thash] = append(s.logs[s.thash], l)
	return
}

func (s *MockStateDB) GetLog() []*types.Log {
	return s.logs[s.thash]
}

func (s *MockStateDB) AddPreimage(common.Hash, []byte) {
	return
}

func (s *MockStateDB) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {
	return
}

func (s *MockStateDB) TxHash() common.Hash {
	return s.thash
}
func (s *MockStateDB) TxIdx() uint32 {
	return uint32(s.txIndex)
}

func generateHeader(num *big.Int, parentHash common.Hash, htime *big.Int, coninbase common.Address) *types.Header {
	h := new(types.Header)
	h.Number = num
	h.ParentHash = parentHash
	h.Coinbase = coninbase
	h.Time = htime
	return h
}
