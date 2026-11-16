package miner

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

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
