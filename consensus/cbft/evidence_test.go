package cbft

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
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

func TestNewEvidencePool(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	_, err := NewEvidencePool(p)
	assert.Nil(t, err)
}

//type viewChangeVote struct {
//	Timestamp      uint64                  `json:"timestamp"`
//	BlockNum       uint64                  `json:"block_number"`
//	BlockHash      common.Hash             `json:"block_hash"`
//	ProposalIndex  uint32                  `json:"proposal_index"`
//	ProposalAddr   common.Address          `json:"proposal_address"`
//	ValidatorIndex uint32                  `json:"validator_index"`
//	ValidatorAddr  common.Address          `json:"-"`
//	Signature      common.BlockConfirmSign `json:"-"`
//	Extra          []byte
//}

//type prepareVote struct {
//	Timestamp      uint64
//	Hash           common.Hash
//	Number         uint64
//	ValidatorIndex uint32
//	ValidatorAddr  common.Address
//	Signature      common.BlockConfirmSign
//	Extra          []byte
//}

func TestAdd(t *testing.T) {

	p := path()
	defer os.RemoveAll(p)
	pool, err := NewEvidencePool(p)
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
		v := &viewChangeVote{
			Timestamp:      timestamp,
			BlockNum:       uint64(i),
			BlockHash:      hash,
			ProposalIndex:  0,
			ProposalAddr:   proposer,
			ValidatorIndex: uint32(i + 1),
			ValidatorAddr:  address,
		}
		cb, _ := v.CannibalizeBytes()
		sign, _ := crypto.Sign(cb, a)
		v.Signature.SetBytes(sign)

		assert.Nil(t, pool.AddViewChangeVote(v))

		p := &prepareVote{
			Timestamp:      timestamp,
			Number:         uint64(i),
			Hash:           hash,
			ValidatorIndex: uint32(i + 1),
			ValidatorAddr:  address,
		}

		cb, _ = p.CannibalizeBytes()
		sign, _ = crypto.Sign(cb, a)
		p.Signature.SetBytes(sign)

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
	pool, err := NewEvidencePool(name)
	if err != nil {
		t.Error(err)
		return
	}
	account := createAccount(2)
	validator := account[1]
	proposerAddr := crypto.PubkeyToAddress(account[0].PublicKey)
	validatorAddr := crypto.PubkeyToAddress(account[1].PublicKey)

	p := &viewChangeVote{
		Timestamp:      0,
		BlockNum:       1,
		BlockHash:      common.BytesToHash(Rand32Bytes(32)),
		ProposalIndex:  0,
		ProposalAddr:   proposerAddr,
		ValidatorIndex: uint32(2),
		ValidatorAddr:  validatorAddr,
	}

	cb, _ := p.CannibalizeBytes()
	sign, _ := crypto.Sign(cb, validator)
	p.Signature.SetBytes(sign)

	assert.Nil(t, pool.AddViewChangeVote(p))
	p = &viewChangeVote{
		Timestamp:      1,
		BlockNum:       0,
		BlockHash:      common.BytesToHash(Rand32Bytes(32)),
		ProposalIndex:  0,
		ProposalAddr:   proposerAddr,
		ValidatorIndex: uint32(2),
		ValidatorAddr:  validatorAddr,
	}

	cb, _ = p.CannibalizeBytes()
	sign, _ = crypto.Sign(cb, validator)
	p.Signature.SetBytes(sign)
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
	pool, err := NewEvidencePool(name)
	if err != nil {
		t.Error(err)
		return
	}
	account := createAccount(2)
	validator := account[1]
	proposerAddr := crypto.PubkeyToAddress(account[0].PublicKey)
	validatorAddr := crypto.PubkeyToAddress(account[1].PublicKey)

	p := &viewChangeVote{
		Timestamp:      0,
		BlockNum:       1,
		BlockHash:      common.BytesToHash(Rand32Bytes(32)),
		ProposalIndex:  0,
		ProposalAddr:   proposerAddr,
		ValidatorIndex: uint32(2),
		ValidatorAddr:  validatorAddr,
	}

	cb, _ := p.CannibalizeBytes()
	sign, _ := crypto.Sign(cb, validator)
	p.Signature.SetBytes(sign)

	assert.Nil(t, pool.AddViewChangeVote(p))
	p = &viewChangeVote{
		Timestamp:      1,
		BlockNum:       1,
		BlockHash:      common.BytesToHash(Rand32Bytes(32)),
		ProposalIndex:  0,
		ProposalAddr:   proposerAddr,
		ValidatorIndex: uint32(2),
		ValidatorAddr:  validatorAddr,
	}

	cb, _ = p.CannibalizeBytes()
	sign, _ = crypto.Sign(cb, validator)
	p.Signature.SetBytes(sign)
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
	pool, err := NewEvidencePool(name)
	if err != nil {
		t.Error(err)
		return
	}
	account := createAccount(1)[0]

	address := crypto.PubkeyToAddress(account.PublicKey)

	p := &prepareVote{
		Timestamp:      0,
		Number:         uint64(1),
		Hash:           common.BytesToHash(Rand32Bytes(32)),
		ValidatorIndex: uint32(1),
		ValidatorAddr:  address,
	}

	cb, _ := p.CannibalizeBytes()
	sign, _ := crypto.Sign(cb, account)
	p.Signature.SetBytes(sign)

	assert.Nil(t, pool.AddPrepareVote(p))
	p = &prepareVote{
		Timestamp:      0,
		Number:         uint64(1),
		Hash:           common.BytesToHash(Rand32Bytes(32)),
		ValidatorIndex: uint32(1),
		ValidatorAddr:  address,
	}

	cb, _ = p.CannibalizeBytes()
	sign, _ = crypto.Sign(cb, account)
	p.Signature.SetBytes(sign)
	assert.IsType(t, &DuplicatePrepareVoteEvidence{}, pool.AddPrepareVote(p))
	assert.Len(t, pool.Evidences(), 1)
	for _, e := range pool.Evidences() {
		assert.Nil(t, e.Validate())
		assert.Nil(t, e.Verify(account.PublicKey))
	}

}
