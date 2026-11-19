package gov

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func TestVersionActivationParentSnapshot(t *testing.T) {
	db, err := snapshotdb.OpenWithStorage(storage.NewMemStorage(), 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	state := mock.NewMockStateDB()
	committed, parent, sibling, missing := common.Hash{1}, common.Hash{2}, common.Hash{3}, common.Hash{4}
	if err := db.SetCurrent(committed, *big.NewInt(10), *big.NewInt(10)); err != nil {
		t.Fatal(err)
	}
	for _, hash := range []common.Hash{parent, sibling} {
		if err := db.NewBlock(big.NewInt(11), committed, hash); err != nil {
			t.Fatal(err)
		}
	}
	check := func(hash common.Hash, height int64, wantActive, wantCertain bool) {
		t.Helper()
		active, certain := VersionActivationAt(db, hash, big.NewInt(height), state)
		if active != wantActive || certain != wantCertain {
			t.Fatalf("parent %x height %d: got (%t, %t), want (%t, %t)", hash, height, active, certain, wantActive, wantCertain)
		}
	}
	check(committed, 11, false, true)
	check(missing, 12, false, false)
	proposal := &VersionProposal{ProposalID: common.Hash{5}, ProposalType: Version, NewVersion: params.FORKVERSION_1_6_0, ActiveBlock: 12}
	id, err := rlp.EncodeToBytes(proposal.ProposalID)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Put(parent, KeyPreActiveProposal(), id); err != nil {
		t.Fatal(err)
	}
	check(parent, 12, false, false)
	if err := NewGovDB(db).SetProposal(proposal, state); err != nil {
		t.Fatal(err)
	}
	check(parent, 12, false, false)
	gdb := NewGovDB(db)
	if err := gdb.SetPreActiveVersion(parent, params.FORKVERSION_1_5_0); err != nil {
		t.Fatal(err)
	}
	check(parent, 12, false, false)
	if err := gdb.SetPreActiveVersion(parent, proposal.NewVersion); err != nil {
		t.Fatal(err)
	}
	check(parent, 12, false, false)
	if err := gdb.SetTallyResult(TallyResult{ProposalID: proposal.ProposalID, Status: PreActive}, state); err != nil {
		t.Fatal(err)
	}
	check(parent, 12, true, true)
	check(parent, 13, false, true)
	check(sibling, 12, false, true)
	if err := db.PutBaseDB(KeyPreActiveProposal(), id); err != nil {
		t.Fatal(err)
	}
	if err := db.PutBaseDB(KeyPreActiveVersion(), common.MustRlpEncode(proposal.NewVersion)); err != nil {
		t.Fatal(err)
	}
	check(sibling, 12, true, true)
	check(missing, 12, false, false)
	if err := db.Del(parent, KeyPreActiveProposal()); err != nil {
		t.Fatal(err)
	}
	if err := db.Del(parent, KeyPreActiveVersion()); err != nil {
		t.Fatal(err)
	}
	check(parent, 12, false, true)
	if err := db.Put(parent, KeyPreActiveProposal(), []byte{1}); err != nil {
		t.Fatal(err)
	}
	if err := gdb.SetPreActiveVersion(parent, proposal.NewVersion); err != nil {
		t.Fatal(err)
	}
	check(parent, 12, false, false)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	check(parent, 12, false, false)
}

func TestHawkingActiveConfigForkWithoutState(t *testing.T) {
	cfg := &params.ChainConfig{
		PauliBlock:   big.NewInt(0),
		HawkingBlock: big.NewInt(0),
	}
	if !HawkingActive(cfg, big.NewInt(1), nil) {
		t.Fatal("HawkingBlock must enable rules without governance state")
	}
	before := &params.ChainConfig{PauliBlock: big.NewInt(0)}
	if HawkingActive(before, big.NewInt(1), nil) {
		t.Fatal("missing HawkingBlock and state must stay inactive")
	}
}
