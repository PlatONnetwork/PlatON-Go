package miner

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"path/filepath"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/naoina/toml"
)

func TestInnerSignerMigrationAndExport(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "inner-admin")
	if err := crypto.SaveECDSA(path, key); err != nil {
		t.Fatal(err)
	}
	config := &Config{InnerAdminKeyFile: path, GasPrice: big.NewInt(1)}
	chain := &params.ChainConfig{Cbft: &params.CbftConfig{ValidatorMode: "inner"}}
	if config.ValidateInnerSigner(chain, false) == nil {
		t.Fatal("missing signer accepted")
	}
	if err := config.LoadInnerSigner(path, chain); err != nil {
		t.Fatal(err)
	}
	if config.ValidateInnerSigner(chain, false) != nil || config.ValidateInnerSigner(chain, true) == nil {
		t.Fatal("unexpected validation across admin activation")
	}
	chain.Cbft.InnerAdmin = crypto.PubkeyToAddress(key.PublicKey)
	if err := config.ValidateInnerSigner(chain, true); err != nil {
		t.Fatal(err)
	}
	out, err := toml.Marshal(config)
	if err != nil || bytes.Contains(out, []byte(hex.EncodeToString(crypto.FromECDSA(key)))) || !bytes.Contains(out, []byte("inner-admin")) {
		t.Fatal("expected only the key file reference in exported config", err)
	}
	chain.Cbft.InnerAdmin = common.Address{1}
	if err := config.ValidateInnerSigner(chain, true); err == nil {
		t.Fatal("cached signer does not match the configured admin")
	}
	if err := config.LoadInnerSigner(path, chain); err == nil {
		t.Fatal("non-admin signer accepted")
	}
}

func TestInnerTxIncluded(t *testing.T) {
	want := common.HexToHash("0x01")
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	if !innerTxIncluded([]*types.Transaction{tx}, tx.Hash()) {
		t.Fatal("included tx not found")
	}
	if innerTxIncluded([]*types.Transaction{tx}, want) {
		t.Fatal("unrelated hash reported as included")
	}
	if innerTxIncluded(nil, tx.Hash()) {
		t.Fatal("empty list reported a hit")
	}
}
