package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"math/big"
)

type BasePlugin interface {
	BeginBlock(header *types.Header, state StateDB) (bool, error)
	EndBlock(header *types.Header, state StateDB) (bool, error)
	Confirmed(block *types.Block) error
	Verify_tx_data(source [][]byte) error
}


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
