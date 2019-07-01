package cbft

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/consensus"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sort"
	"testing"
)

func path() string {
	name, err := ioutil.TempDir(os.TempDir(), "evidence")

	if err != nil {
		panic(err)
	}
	return name
}

func createAccount(n int) []*ecdsa.PrivateKey {
	var pris []*ecdsa.PrivateKey
	for i := 0; i < n; i++ {
		pri, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		pris = append(pris, pri)
	}
	return pris
}

func TestTimeOrderViewChange_Add(t *testing.T) {
	var p TimeOrderViewChange
	p.Add(nil)
	assert.Len(t, p, 1)
}

func TestNewBaseEvidencePool(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	_, err := NewBaseEvidencePool(p)
	assert.Nil(t, err)
}

func TestAdd(t *testing.T) {

	p := path()
	defer os.RemoveAll(p)
	pool, err := NewBaseEvidencePool(p)
	if err != nil {
		t.Error(err)
		return
	}
	proposer := common.BytesToAddress(Rand32Bytes(40))
	accounts := createAccount(100)
	assert.Len(t, accounts, 100)
	for i, a := range accounts {
		hash := common.BytesToHash(Rand32Bytes(32))
		timestamp := uint64(i)
		address := crypto.PubkeyToAddress(a.PublicKey)
		v := makeViewChangeVote(a, timestamp, uint64(i), hash, 0, proposer, uint32(i+1), address)
		assert.Nil(t, pool.AddViewChangeVote(v))
		p := makePrepareVote(a, timestamp, uint64(i), hash, uint32(i+1), address)
		assert.Nil(t, pool.AddPrepareVote(p))

		assert.True(t, sort.IsSorted(pool.vt[address]))
		assert.True(t, sort.IsSorted(pool.vn[address]))
		assert.True(t, sort.IsSorted(pool.pe[address]))
	}

	assert.Len(t, pool.vt, 100)
	assert.Len(t, pool.vn, 100)
	assert.Len(t, pool.pe, 100)

	pool.Clear(200, 200)
	assert.Len(t, pool.vt, 0)
	assert.Len(t, pool.vn, 0)
	assert.Len(t, pool.pe, 0)
}

func TestTimestampViewChangeVoteEvidence(t *testing.T) {

	name := path()
	defer os.RemoveAll(name)
	pool, err := NewBaseEvidencePool(name)
	if err != nil {
		t.Error(err)
		return
	}
	account := createAccount(2)
	validator := account[1]
	proposerAddr := crypto.PubkeyToAddress(account[0].PublicKey)
	validatorAddr := crypto.PubkeyToAddress(account[1].PublicKey)

	p := makeViewChangeVote(validator, 0, 1, common.BytesToHash(Rand32Bytes(32)), 0, proposerAddr, uint32(2), validatorAddr)

	assert.Nil(t, pool.AddViewChangeVote(p))

	p = makeViewChangeVote(validator, 1, 0, common.BytesToHash(Rand32Bytes(32)), 0, proposerAddr, uint32(2), validatorAddr)

	assert.IsType(t, &TimestampViewChangeVoteEvidence{}, pool.AddViewChangeVote(p))
	assert.Len(t, pool.Evidences(), 1)
	for _, e := range pool.Evidences() {
		assert.Nil(t, e.Validate())
		assert.Nil(t, e.Verify(validator.PublicKey))
	}
}

func TestDuplicateViewChangeVoteEvidence(t *testing.T) {
	name := path()
	defer os.RemoveAll(name)
	pool, err := NewBaseEvidencePool(name)
	if err != nil {
		t.Error(err)
		return
	}
	account := createAccount(2)
	validator := account[1]
	proposerAddr := crypto.PubkeyToAddress(account[0].PublicKey)
	validatorAddr := crypto.PubkeyToAddress(account[1].PublicKey)

	p := makeViewChangeVote(validator, 0, 1, common.BytesToHash(Rand32Bytes(32)), 0, proposerAddr, uint32(2), validatorAddr)

	assert.Nil(t, pool.AddViewChangeVote(p))
	p = makeViewChangeVote(validator, 1, 1, common.BytesToHash(Rand32Bytes(32)), 0, proposerAddr, uint32(2), validatorAddr)

	assert.IsType(t, &DuplicateViewChangeVoteEvidence{}, pool.AddViewChangeVote(p))
	assert.Len(t, pool.Evidences(), 1)
	for _, e := range pool.Evidences() {
		assert.Nil(t, e.Validate())
		assert.Nil(t, e.Verify(validator.PublicKey))
	}
}

func TestDuplicatePrepareVoteEvidence(t *testing.T) {

	name := path()
	defer os.RemoveAll(name)
	pool, err := NewBaseEvidencePool(name)
	if err != nil {
		t.Error(err)
		return
	}
	account := createAccount(1)[0]

	address := crypto.PubkeyToAddress(account.PublicKey)

	p := makePrepareVote(account, 0, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), address)

	assert.Nil(t, pool.AddPrepareVote(p))
	p = makePrepareVote(account, 0, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), address)

	assert.IsType(t, &DuplicatePrepareVoteEvidence{}, pool.AddPrepareVote(p))
	assert.Len(t, pool.Evidences(), 1)
	for _, e := range pool.Evidences() {
		assert.Nil(t, e.Validate())
		assert.Nil(t, e.Verify(account.PublicKey))
	}

}

type prepareVoteData struct {
	voteA *prepareVote
	voteB *prepareVote
	valid bool
}

func TestDuplicatePrepareVoteEvidence_Validate(t *testing.T) {
	priA := createAccount(1)[0]
	addrA := crypto.PubkeyToAddress(priA.PublicKey)
	priB := createAccount(1)[0]
	addrB := crypto.PubkeyToAddress(priB.PublicKey)

	voteA := makePrepareVote(priA, 0, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), addrA)
	invalidVoteB := makePrepareVote(priA, 1, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), addrA)
	invalidVoteB.Signature.SetBytes(Rand32Bytes(32))
	testCases := []prepareVoteData{
		{voteA, makePrepareVote(priA, 0, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), addrA), true},
		{voteA, voteA, false},
		{voteA, makePrepareVote(priA, 1, uint64(2), common.BytesToHash(Rand32Bytes(32)), uint32(1), addrA), false},
		{voteA, makePrepareVote(priA, 1, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(2), addrA), false},
		{voteA, invalidVoteB, false},
		{voteA, makePrepareVote(priB, 1, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(2), addrB), false},
	}

	for i, p := range testCases {
		d := DuplicatePrepareVoteEvidence{
			p.voteA,
			p.voteB,
		}
		if p.valid {
			assert.Nil(t, d.Validate(), fmt.Sprintf("testcase:%d error", i))
		} else {
			assert.NotNil(t, d.Validate(), fmt.Sprintf("testcase:%d error", i))
		}
	}
}

func TestDuplicatePrepareVoteEvidence_Equal(t *testing.T) {
	priA := createAccount(1)[0]
	addrA := crypto.PubkeyToAddress(priA.PublicKey)
	voteA := makePrepareVote(priA, 0, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), addrA)
	voteB := makePrepareVote(priA, 0, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), addrA)
	voteC := makePrepareVote(priA, 0, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), addrA)
	p1 := &DuplicatePrepareVoteEvidence{
		VoteA: voteA,
		VoteB: voteB,
	}

	p2 := &DuplicatePrepareVoteEvidence{
		VoteA: voteA,
		VoteB: voteC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}

func TestDuplicateViewChangeVoteEvidence_Equal(t *testing.T) {
	priA := createAccount(1)[0]
	addrA := crypto.PubkeyToAddress(priA.PublicKey)
	voteA := makeViewChangeVote(priA, 0, 1, common.BytesToHash([]byte{1}), 0, addrA, uint32(2), addrA)
	voteB := makeViewChangeVote(priA, 0, 1, common.BytesToHash([]byte{3}), 0, addrA, uint32(2), addrA)
	voteC := makeViewChangeVote(priA, 0, 2, common.BytesToHash([]byte{1}), 0, addrA, uint32(2), addrA)
	p1 := &DuplicateViewChangeVoteEvidence{
		VoteA: voteA,
		VoteB: voteB,
	}

	p2 := &DuplicateViewChangeVoteEvidence{
		VoteA: voteA,
		VoteB: voteC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}

func TestTimestampViewChangeVoteEvidence_Equal(t *testing.T) {
	priA := createAccount(1)[0]
	addrA := crypto.PubkeyToAddress(priA.PublicKey)
	voteA := makeViewChangeVote(priA, 0, 1, common.BytesToHash([]byte{1}), 0, addrA, uint32(2), addrA)
	voteB := makeViewChangeVote(priA, 0, 1, common.BytesToHash([]byte{3}), 0, addrA, uint32(2), addrA)
	voteC := makeViewChangeVote(priA, 0, 2, common.BytesToHash([]byte{1}), 0, addrA, uint32(2), addrA)
	p1 := &TimestampViewChangeVoteEvidence{
		VoteA: voteA,
		VoteB: voteB,
	}

	p2 := &TimestampViewChangeVoteEvidence{
		VoteA: voteA,
		VoteB: voteC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}

type viewChangeVoteData struct {
	voteA *viewChangeVote
	voteB *viewChangeVote
	valid bool
}

func TestDuplicateViewChangeVoteEvidence_Validate(t *testing.T) {
	priA := createAccount(1)[0]
	addrA := crypto.PubkeyToAddress(priA.PublicKey)
	priB := createAccount(1)[0]
	addrB := crypto.PubkeyToAddress(priB.PublicKey)

	voteA := makeViewChangeVote(priA, 0, 1, common.BytesToHash([]byte{1}), 0, addrA, uint32(2), addrA)

	invalidVoteB := makeViewChangeVote(priA, 0, 1, common.BytesToHash(Rand32Bytes(32)), 0, addrA, uint32(2), addrA)
	invalidVoteB.Signature.SetBytes(Rand32Bytes(32))

	testCases := []viewChangeVoteData{
		{voteA, makeViewChangeVote(priA, 0, 1, common.BytesToHash([]byte{4}), 0, addrA, uint32(2), addrA), true},
		{voteA, voteA, false},
		{voteA, makeViewChangeVote(priA, 0, 2, common.BytesToHash(Rand32Bytes(32)), 0, addrA, uint32(2), addrA), false},
		{voteA, makeViewChangeVote(priA, 0, 1, common.BytesToHash(Rand32Bytes(32)), 0, addrA, uint32(3), addrA), false},
		{voteA, invalidVoteB, false},
		{voteA, makeViewChangeVote(priA, 0, 1, common.BytesToHash(Rand32Bytes(32)), 0, addrA, uint32(2), addrB), false},
	}

	for i, p := range testCases {
		d := DuplicateViewChangeVoteEvidence{
			p.voteA,
			p.voteB,
		}
		if p.valid {
			assert.Nil(t, d.Validate(), fmt.Sprintf("testcase:%d error", i))
		} else {
			assert.NotNil(t, d.Validate(), fmt.Sprintf("testcase:%d error", i))
		}
	}
}

func TestTimestampViewChangeVoteEvidence_Validate(t *testing.T) {
	priA := createAccount(1)[0]
	addrA := crypto.PubkeyToAddress(priA.PublicKey)
	priB := createAccount(1)[0]
	addrB := crypto.PubkeyToAddress(priB.PublicKey)

	voteA := makeViewChangeVote(priA, 0, 5, common.BytesToHash([]byte{1}), 0, addrA, uint32(2), addrA)

	invalidVoteB := makeViewChangeVote(priA, 0, 1, common.BytesToHash(Rand32Bytes(32)), 0, addrA, uint32(2), addrA)
	invalidVoteB.Signature.SetBytes(Rand32Bytes(32))

	testCases := []viewChangeVoteData{
		{voteA, makeViewChangeVote(priA, 1, 4, common.BytesToHash([]byte{4}), 0, addrA, uint32(2), addrA), true},
		{voteA, voteA, false},
		{voteA, makeViewChangeVote(priA, 1, 7, common.BytesToHash(Rand32Bytes(32)), 0, addrA, uint32(2), addrA), false},
		{voteA, makeViewChangeVote(priA, 0, 5, common.BytesToHash(Rand32Bytes(32)), 0, addrA, uint32(3), addrA), false},
		{voteA, invalidVoteB, false},
		{voteA, makeViewChangeVote(priA, 0, 5, common.BytesToHash(Rand32Bytes(32)), 0, addrA, uint32(2), addrB), false},
	}

	for i, p := range testCases {
		d := TimestampViewChangeVoteEvidence{
			p.voteA,
			p.voteB,
		}
		if p.valid {
			assert.Nil(t, d.Validate(), fmt.Sprintf("testcase:%d error", i))
		} else {
			assert.NotNil(t, d.Validate(), fmt.Sprintf("testcase:%d error", i))
		}
	}
}

func TestJson(t *testing.T) {
	priA := createAccount(1)[0]
	addrA := crypto.PubkeyToAddress(priA.PublicKey)

	voteA := makeViewChangeVote(priA, 0, 5, common.BytesToHash([]byte{1}), 0, addrA, uint32(2), addrA)
	voteB := makePrepareVote(priA, 0, uint64(1), common.BytesToHash(Rand32Bytes(32)), uint32(1), addrA)

	evs := []consensus.Evidence{
		&DuplicateViewChangeVoteEvidence{
			VoteB: voteA,
			VoteA: voteA,
		},
		&TimestampViewChangeVoteEvidence{
			VoteB: voteA,
			VoteA: voteA,
		},
		&DuplicatePrepareVoteEvidence{
			VoteA: voteB,
			VoteB: voteB,
		},
	}
	eds := ClassifyEvidence(evs)
	b, _ := json.MarshalIndent(eds, "", "  ")
	t.Log(string(b))
	var eds2 EvidenceData
	assert.Nil(t, json.Unmarshal(b, &eds2))

	b2, _ := json.MarshalIndent(eds2, "", "  ")
	assert.Equal(t, b, b2)

	_, err := NewEvidences(string(b))
	assert.Nil(t, err)

}
