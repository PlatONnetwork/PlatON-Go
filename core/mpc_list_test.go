package core

import (
	"math/rand"
	"testing"

	"Platon-go/core/types"
	"Platon-go/crypto"
)

func TestMPCListAdd(t *testing.T) {
	// Generate a list of transactions to insert
	key, _ := crypto.GenerateKey()

	txs := make(types.TransactionWraps, 1024)
	for i := 0; i < len(txs); i++ {
		txs[i] = mpcTransaction("a2c4d041f7f88c8be5ea8bac94c0a28178b47bae1dfc01100a26b01de04dd368", uint64(i), 0, key)
	}
	// Insert the transactions in a random order
	all := newMpcLookup()
	list := newMpcList(all)
	for _, v := range rand.Perm(len(txs)) {
		all.Add(txs[v])
		list.Put(txs[v])
	}
	// Verify internal state
	if len(list.all.txs) != len(txs) {
		t.Errorf("mpc transaction count mismatch: have %d, want %d", len(list.all.txs), len(txs))
	}
	if list.items.Len() != len(txs) {
		t.Errorf("mpc transaction count mismatch: have %d, want %d", list.items.Len(), len(txs))
	}
	for i, tx := range txs {
		if list.all.txs[tx.Hash()] != tx {
			t.Errorf("item %d: transaction mismatch: have %v, want %v", i, list.all.txs[tx.Hash()], tx)
		}
	}

}
