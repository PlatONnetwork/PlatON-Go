package xcom

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"math/big"
	"testing"
)

var pk *ecdsa.PublicKey

func initHandler() {
	NewVrfHandler(snapshotdb.Instance(), hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23"))
	if pk == nil {
		pri, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		pk = &pri.PublicKey
		vh.SetPrivateKey(pri)
	}
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
	if err := vh.Storage(new(big.Int).SetUint64(1), common.ZeroHash, common.ZeroHash, hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23")); nil != err {
		t.Error(err)
	}
	if err := vh.db.Flush(hash, blockNumber); nil != err {
		t.Error(err)
	}
	if err := vh.db.NewBlock(new(big.Int).SetUint64(2), hash, common.ZeroHash); nil != err {
		t.Error(err)
	}
	if err := vh.Storage(new(big.Int).SetUint64(2), hash, common.ZeroHash, hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd33")); nil != err {
		t.Error(err)
	}
	if _, err := vh.Load(common.Hash{}); nil != err {
		t.Error(err)
	}
}

func TestVrfHandler_Verify(t *testing.T) {
	initHandler()
	defer func() {
		vh.db.Clear()
	}()
	blockNumber := new(big.Int).SetUint64(1)
	hash := common.BytesToHash([]byte("h1"))
	if value, err := vh.GenerateNonce(blockNumber, common.Hash{}); nil != err {
		t.Error(err)
	} else {
		if err := vh.VerifyVrf(pk, blockNumber, hash, common.ZeroHash, value); nil != err {
			t.Error(err)
		}
	}
}