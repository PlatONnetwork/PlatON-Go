// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

//go:generate go run ../../rlp/rlpgen -type StateAccount -out gen_account_rlp.go

// StateAccount is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type StateAccount struct {
	Nonce            uint64
	Balance          *big.Int
	Root             common.Hash // merkle root of the storage trie
	CodeHash         []byte
	StorageKeyPrefix []byte // A prefix added to the `key` to ensure that data between different accounts are not shared
}

func (self *StateAccount) Empty() bool {
	if self.Nonce != 0 {
		return false
	}
	if self.Balance.Cmp(common.Big0) != 0 {
		return false
	}
	if self.Root != common.ZeroHash {
		return false
	}
	if len(self.CodeHash) != 0 {
		return false
	}
	if len(self.StorageKeyPrefix) != 0 {
		return false
	}
	return true
}
