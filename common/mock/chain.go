package mock

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type chain struct {
	Genesis  *types.Block
	chain    []common.Hash
	headerm  map[common.Hash]*types.Header
	blockm   map[common.Hash]*types.Block
	receiptm map[common.Hash][]*types.Receipt
	StateDB  *MockStateDB
}

func NewChain(genesis *types.Block) *chain {
	c := new(chain)
	c.Genesis = genesis
	db := new(MockStateDB)
	db.State = make(map[common.Address]map[string][]byte)
	db.Balance = make(map[common.Address]*big.Int)
	c.StateDB = db
	return c
}

type MockStateDB struct {
	Balance      map[common.Address]*big.Int
	State        map[common.Address]map[string][]byte
	thash, bhash common.Hash
	txIndex      int
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

func (s *MockStateDB) AddLog(*types.Log) {
	return
}
func (s *MockStateDB) AddPreimage(common.Hash, []byte) {
	return
}

func (s *MockStateDB) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {
	return
}

func (s *MockStateDB) TxHash() common.Hash {
	return common.ZeroHash
}
func (s *MockStateDB) TxIdx() uint32 {
	return 0
}
