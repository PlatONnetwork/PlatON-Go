package handler

import (
	"crypto/ecdsa"
	"math/big"
	"strconv"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/vrf"
	"github.com/stretchr/testify/assert"
)

func initHandler() *ecdsa.PrivateKey {
	NewVrfHandler(hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23"))
	pri, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	vh.SetPrivateKey(pri)
	return pri
}

func TestVrfHandler_StorageLoad(t *testing.T) {
	initHandler()
	defer func() {
		vh.db.Clear()
	}()
	blockNumber := new(big.Int).SetUint64(1)
	phash := common.BytesToHash([]byte("h"))
	hash := common.ZeroHash
	for i := 0; i < int(xcom.EpochValidatorNum())+10; i++ {
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
		assert.Equal(t, len(value), int(xcom.EpochValidatorNum()))
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
