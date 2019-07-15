package xcom

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/vrf"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func initHandler() *ecdsa.PrivateKey {
	NewVrfHandler(hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23"))
	pri, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	vh.SetPrivateKey(pri)
	SetEconomicModel(&DefaultConfig)
	return pri
}

func TestVrfHandler_StorageLoad(t *testing.T) {
	initHandler()
	defer func() {
		vh.db.Clear()
	}()
	blockNumber := new(big.Int).SetUint64(1)
	hash := common.BytesToHash([]byte("h1"))
	if err := vh.db.NewBlock(blockNumber, common.ZeroHash, common.ZeroHash); nil != err {
		t.Error(err)
	}
	pi, err := vh.GenerateNonce(blockNumber, common.ZeroHash)
	if nil != err {
		t.Error(err)
	}
	if err := vh.Storage(new(big.Int).SetUint64(1), common.ZeroHash, common.ZeroHash, vrf.ProofToHash(pi)); nil != err {
		t.Error(err)
	}
	if err := vh.db.Flush(hash, blockNumber); nil != err {
		t.Error(err)
	}
	if err := vh.db.NewBlock(new(big.Int).SetUint64(2), hash, common.ZeroHash); nil != err {
		t.Error(err)
	}
	pi, err = vh.GenerateNonce(new(big.Int).SetUint64(2), common.ZeroHash)
	if nil != err {
		t.Error(err)
	}
	if err := vh.Storage(new(big.Int).SetUint64(2), hash, common.ZeroHash, vrf.ProofToHash(pi)); nil != err {
		t.Error(err)
	}
	if _, err := vh.Load(common.Hash{}); nil != err {
		t.Error(err)
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
		t.Error(err)
	} else {
		if err := vh.VerifyVrf(&sk.PublicKey, blockNumber, hash, common.ZeroHash, value); nil != err {
			t.Error(err)
		}
		pri, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		vh.SetPrivateKey(pri)
		nonce, err := vh.GenerateNonce(blockNumber, common.Hash{})
		if nil != err {
			t.Error(err)
		}
		err = vh.VerifyVrf(&sk.PublicKey, blockNumber, hash, common.ZeroHash, nonce)
		assert.Equal(t, ErrInvalidVrfProve, err)
	}
}