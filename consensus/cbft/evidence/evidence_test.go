package evidence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/stretchr/testify/assert"
)

func path() string {
	name, err := ioutil.TempDir(os.TempDir(), "evidence")

	if err != nil {
		panic(err)
	}
	return name
}

func TestNewBaseEvidencePool(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	_, err := NewEvidencePool(nil, p)
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

	accounts := createAccount(25) // mock 25 nodes
	assert.Len(t, accounts, 25)

	epoch, viewNumber, blockNumber := uint64(1), uint64(1), int64(1)
	for i, _ := range accounts {
		var block *types.Block
		for j := 0; j < 10; j++ { // mock seal ten block per node
			block = newBlock(blockNumber)
			pb := makePrepareBlock(epoch, viewNumber, block, uint32(j), uint32(i))
			assert.Nil(t, pool.AddPrepareBlock(pb))

			pv := makePrepareVote(epoch, viewNumber, block.Hash(), block.NumberU64(), uint32(j), uint32(i))
			assert.Nil(t, pool.AddPrepareVote(pv))

			blockNumber = blockNumber + 1
		}

		identity := Identity(fmt.Sprintf("%d%d%d", epoch, viewNumber, uint32(i)))
		assert.True(t, sort.IsSorted(pool.pb[identity]))
		assert.True(t, sort.IsSorted(pool.pv[identity]))

		vc := makeViewChange(epoch, viewNumber, block.Hash(), block.NumberU64(), uint32(i))
		assert.Nil(t, pool.AddViewChange(vc))

		viewNumber = viewNumber + 1
	}

	assert.Equal(t, pool.pb.Size(), 250)
	assert.Equal(t, pool.pv.Size(), 250)
	assert.Equal(t, pool.vc.Size(), 25)

	pool.Clear(epoch, 15)
	assert.Equal(t, pool.pb.Size(), 100)
	assert.Equal(t, pool.pv.Size(), 100)
	assert.Equal(t, pool.vc.Size(), 10)
}

func TestDuplicatePrepareBlockEvidence(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	pool, err := NewBaseEvidencePool(p)
	if err != nil {
		t.Error(err)
		return
	}

	pb := makePrepareBlock(1, 1, newBlock(1), 1, 1)
	assert.Nil(t, pool.AddPrepareBlock(pb))

	pb = makePrepareBlock(1, 1, newBlock(1), 1, 1)
	assert.IsType(t, &DuplicatePrepareBlockEvidence{}, pool.AddPrepareBlock(pb))

	assert.Len(t, pool.Evidences(), 1)
}

func TestDuplicatePrepareVoteEvidence(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	pool, err := NewBaseEvidencePool(p)
	if err != nil {
		t.Error(err)
		return
	}

	block := newBlock(1)
	pv := makePrepareVote(1, 1, block.Hash(), block.NumberU64(), 1, 1)
	assert.Nil(t, pool.AddPrepareVote(pv))

	block = newBlock(1)
	pv = makePrepareVote(1, 1, block.Hash(), block.NumberU64(), 1, 1)
	assert.IsType(t, &DuplicatePrepareVoteEvidence{}, pool.AddPrepareVote(pv))

	assert.Len(t, pool.Evidences(), 1)
}

func TestDuplicateViewChangeEvidence(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	pool, err := NewBaseEvidencePool(p)
	if err != nil {
		t.Error(err)
		return
	}

	block := newBlock(1)
	vc := makeViewChange(1, 1, block.Hash(), block.NumberU64(), 1)
	assert.Nil(t, pool.AddViewChange(vc))

	block = newBlock(1)
	vc = makeViewChange(1, 1, block.Hash(), block.NumberU64(), 1)
	assert.IsType(t, &DuplicateViewChangeEvidence{}, pool.AddViewChange(vc))

	assert.Len(t, pool.Evidences(), 1)
}

func TestJson(t *testing.T) {
	pb := makePrepareBlock(1, 1, newBlock(1), 1, 1)
	block1 := newBlock(1)
	pv := makePrepareVote(1, 1, block1.Hash(), block1.NumberU64(), 1, 1)
	block2 := newBlock(1)
	vc := makeViewChange(1, 1, block2.Hash(), block2.NumberU64(), 1)

	evs := []consensus.Evidence{
		&DuplicatePrepareBlockEvidence{
			PrepareA: pb,
			PrepareB: pb,
		},
		&DuplicatePrepareVoteEvidence{
			VoteA: pv,
			VoteB: pv,
		},
		&DuplicateViewChangeEvidence{
			ViewA: vc,
			ViewB: vc,
		},
	}
	ed := ClassifyEvidence(evs)
	b, _ := json.MarshalIndent(ed, "", " ")
	t.Log(string(b))
	var ed2 EvidenceData
	assert.Nil(t, json.Unmarshal(b, &ed2))

	b2, _ := json.MarshalIndent(ed2, "", " ")
	assert.Equal(t, b, b2)

	// test UnmarshalEvidence
	p := path()
	defer os.RemoveAll(p)
	pool, err := NewBaseEvidencePool(p)
	if err != nil {
		t.Error(err)
		return
	}

	evidences, err := pool.UnmarshalEvidence(string(b2))
	assert.Nil(t, err)
	assert.Equal(t, 3, evidences.Len())
}

func TestDuplicatePrepareBlockEvidence_Equal(t *testing.T) {
	pbA := makePrepareBlock(1, 1, newBlock(1), 1, 1)
	pbB := makePrepareBlock(1, 1, newBlock(1), 1, 1)
	pbC := makePrepareBlock(1, 1, newBlock(1), 1, 1)
	p1 := &DuplicatePrepareBlockEvidence{
		PrepareA: pbA,
		PrepareB: pbB,
	}

	p2 := &DuplicatePrepareBlockEvidence{
		PrepareA: pbA,
		PrepareB: pbC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}

func TestDuplicatePrepareVoteEvidence_Equal(t *testing.T) {
	pvA := makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, 1)
	pvB := makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, 1)
	pvC := makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, 1)
	p1 := &DuplicatePrepareVoteEvidence{
		VoteA: pvA,
		VoteB: pvB,
	}

	p2 := &DuplicatePrepareVoteEvidence{
		VoteA: pvA,
		VoteB: pvC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}

func TestDuplicateViewChangeEvidence_Equal(t *testing.T) {
	vcA := makeViewChange(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1)
	vcB := makeViewChange(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1)
	vcC := makeViewChange(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1)
	p1 := &DuplicateViewChangeEvidence{
		ViewA: vcA,
		ViewB: vcB,
	}

	p2 := &DuplicateViewChangeEvidence{
		ViewA: vcA,
		ViewB: vcC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}
