package xcom

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math/big"
)


// StateDB is an Plugin database for full state querying.
type StateDB interface {
	CreateAccount(common.Address)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	// todo: new func for abi of contract.
	GetAbiHash(common.Address) common.Hash
	GetAbi(common.Address) []byte
	SetAbi(common.Address, []byte)

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	// todo: hash -> bytes
	GetCommittedState(common.Address, []byte) []byte
	//GetState(common.Address, common.Hash) common.Hash
	//SetState(common.Address, common.Hash, common.Hash)
	GetState(common.Address, []byte) []byte
	SetState(common.Address, []byte, []byte)


	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)

	//ppos add
	TxHash() common.Hash
	TxIdx() uint32
}

type Evidence interface {
	Verify(ecdsa.PublicKey) error
	Equal(Evidence) bool
	//return lowest number
	BlockNumber() uint64
	Hash() []byte
	Address() common.Address
	Validate() error
}


// inner contract event data
type Result struct {
	Status bool
	Data   string
	ErrMsg string
}


// addLog let the result add to event.
func AddLog(state StateDB, blockNumber uint64, contractAddr common.Address, event, data string) error {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); nil != err {
		return err
	}
	state.AddLog(&types.Log{
		Address:     contractAddr,
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data:        buf.Bytes(),
		BlockNumber: blockNumber,
	})
	return nil
}




type SnapshotDB interface {
	Put(hash common.Hash, key, value []byte) (bool, error)
	NewBlock(blockNumber *big.Int, parentHash common.Hash, hash common.Hash) (bool, error)
	Get(hash common.Hash, key []byte) ([]byte, error)
	GetFromCommittedBlock(key []byte) ([]byte, error)
	Del(hash common.Hash, key []byte) (bool, error)
	Has(hash common.Hash, key []byte) (bool, error)
	Flush(hash common.Hash, blocknumber *big.Int) (bool, error)
	Ranking(hash common.Hash, key []byte, ranges int) iterator.Iterator
	WalkBaseDB(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error
	Commit(hash common.Hash) (bool, error)
	Clear() (bool, error)
	PutBaseDB(key, value []byte) (bool, error)
	DelBaseDB(key []byte) (bool, error)
	GetLastKVHash(blockHash common.Hash) []byte
	BaseNum() (*big.Int, error)
	Close() (bool, error)
	Compaction() (bool, error)
}