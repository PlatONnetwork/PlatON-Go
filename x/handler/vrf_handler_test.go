// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package handler

import (
	"crypto/ecdsa"
	"math/big"
	"strconv"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/vrf"
	"github.com/stretchr/testify/assert"
)

var chain *mock.Chain

func initHandler() *ecdsa.PrivateKey {
	vh = &VrfHandler{
		db:           snapshotdb.Instance(),
		genesisNonce: hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23"),
	}
	//	NewVrfHandler(hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23"))
	pri, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	vh.SetPrivateKey(pri)

	chain = mock.NewChain()

	return pri
}

func TestVrfHandler_StorageLoad(t *testing.T) {
	initHandler()
	defer func() {
		vh.db.Clear()
	}()

	gov.InitGenesisGovernParam(common.ZeroHash, vh.db, 2048)

	blockNumber := new(big.Int).SetUint64(1)
	phash := common.BytesToHash([]byte("h"))
	hash := common.ZeroHash
	for i := 0; i < int(xcom.MaxValidators())+10; i++ {
		if err := vh.db.NewBlock(blockNumber, phash, common.ZeroHash); nil != err {
			t.Fatal(err)
		}
		pi, err := vh.GenerateNonce(blockNumber, phash)
		if nil != err {
			t.Fatal(err)
		}
		if err := vh.Storage(blockNumber, phash, common.ZeroHash, vrf.ProofToHash(pi)); nil != err {
			t.Fatal(err)
		}
		hash = common.BytesToHash([]byte(strconv.Itoa(i)))
		phash = hash
		if err := vh.db.Flush(hash, blockNumber); nil != err {
			t.Fatal(err)
		}
		blockNumber.Add(blockNumber, common.Big1)
	}
	if value, err := vh.Load(phash); nil != err {
		t.Fatal(err)
	} else {
		assert.Equal(t, len(value), int(xcom.MaxValidators()))
	}
}

func TestVrfHandler_Verify(t *testing.T) {
	sk := initHandler()
	defer func() {
		vh.db.Clear()
	}()
	blockNumber := new(big.Int).SetUint64(1)
	hash := common.BytesToHash([]byte("h1"))
	if value, err := vh.GenerateNonce(blockNumber, common.Hash{}); nil != err {
		t.Fatal(err)
	} else {
		if err := vh.VerifyVrf(&sk.PublicKey, blockNumber, hash, common.ZeroHash, value); nil != err {
			t.Fatal(err)
		}
		pri, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		vh.SetPrivateKey(pri)
		nonce, err := vh.GenerateNonce(blockNumber, common.Hash{})
		if nil != err {
			t.Fatal(err)
		}
		err = vh.VerifyVrf(&sk.PublicKey, blockNumber, hash, common.ZeroHash, nonce)
		assert.Equal(t, ErrInvalidVrfProve, err)
	}
}
