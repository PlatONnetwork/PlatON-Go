// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package evidence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/stretchr/testify/assert"
)

func init() {
	bls.Init(bls.BLS12_381)
}

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

func TestAddAndClear(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	pool, err := NewBaseEvidencePool(p)
	if err != nil {
		t.Error(err)
		return
	}

	validateNodes, secretKeys := createValidateNode(25) // mock 25 nodes
	assert.Len(t, validateNodes, 25)

	epoch, viewNumber, blockNumber := uint64(1), uint64(1), int64(1)
	for i := 0; i < len(validateNodes); i++ {
		var block *types.Block
		node := validateNodes[i]
		for j := 0; j < 10; j++ { // mock seal ten block per node
			block = newBlock(blockNumber)
			pb := makePrepareBlock(epoch, viewNumber, block, uint32(j), uint32(node.Index), t, secretKeys[i])
			assert.Nil(t, pool.AddPrepareBlock(pb, node))

			pv := makePrepareVote(epoch, viewNumber, block.Hash(), block.NumberU64(), uint32(j), uint32(node.Index), t, secretKeys[i])
			assert.Nil(t, pool.AddPrepareVote(pv, node))

			blockNumber = blockNumber + 1
		}

		identity := Identity(fmt.Sprintf("%d%d%d", epoch, viewNumber, uint32(node.Index)))
		assert.True(t, sort.IsSorted(pool.pb[identity]))
		assert.True(t, sort.IsSorted(pool.pv[identity]))

		vc := makeViewChange(epoch, viewNumber, block.Hash(), block.NumberU64(), uint32(node.Index), t, secretKeys[i])
		assert.Nil(t, pool.AddViewChange(vc, node))

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
	validateNodes, secretKeys := createValidateNode(1)
	pb := makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[0].Index, t, secretKeys[0])

	assert.Nil(t, pool.AddPrepareBlock(pb, validateNodes[0]))

	pb = makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[0].Index, t, secretKeys[0])
	assert.IsType(t, &DuplicatePrepareBlockEvidence{}, pool.AddPrepareBlock(pb, validateNodes[0]))

	assert.Len(t, pool.Evidences(), 1)

	// test json
	evdata := ClassifyEvidence(pool.Evidences())
	b, _ := json.MarshalIndent(evdata, "", " ")
	t.Log(string(b))
	var ed2 EvidenceData
	assert.Nil(t, json.Unmarshal(b, &ed2))

	b2, _ := json.MarshalIndent(ed2, "", " ")
	assert.Equal(t, b, b2)

	// test NewEvidence
	dp := evdata.DP[0]
	dpSeri, _ := json.MarshalIndent(dp, "", " ")
	pbEvidence, err := NewEvidence(DuplicatePrepareBlockType, string(dpSeri))
	assert.Nil(t, err)
	assert.Nil(t, pbEvidence.Validate())
}

func TestDuplicatePrepareVoteEvidence(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	pool, err := NewBaseEvidencePool(p)
	if err != nil {
		t.Error(err)
		return
	}

	validateNodes, secretKeys := createValidateNode(1)
	block := newBlock(1)
	pv := makePrepareVote(1, 1, block.Hash(), block.NumberU64(), 1, validateNodes[0].Index, t, secretKeys[0])
	assert.Nil(t, pool.AddPrepareVote(pv, validateNodes[0]))

	block = newBlock(1)
	pv = makePrepareVote(1, 1, block.Hash(), block.NumberU64(), 1, validateNodes[0].Index, t, secretKeys[0])
	assert.IsType(t, &DuplicatePrepareVoteEvidence{}, pool.AddPrepareVote(pv, validateNodes[0]))

	assert.Len(t, pool.Evidences(), 1)

	// test json
	evdata := ClassifyEvidence(pool.Evidences())
	b, _ := json.MarshalIndent(evdata, "", " ")
	t.Log(string(b))
	var ed2 EvidenceData
	assert.Nil(t, json.Unmarshal(b, &ed2))

	b2, _ := json.MarshalIndent(ed2, "", " ")
	assert.Equal(t, b, b2)
}

func TestDuplicateViewChangeEvidence(t *testing.T) {
	p := path()
	defer os.RemoveAll(p)
	pool, err := NewBaseEvidencePool(p)
	if err != nil {
		t.Error(err)
		return
	}

	validateNodes, secretKeys := createValidateNode(1)
	block := newBlock(1)
	vc := makeViewChange(1, 1, block.Hash(), block.NumberU64(), validateNodes[0].Index, t, secretKeys[0])
	assert.Nil(t, pool.AddViewChange(vc, validateNodes[0]))

	block = newBlock(1)
	vc = makeViewChange(1, 1, block.Hash(), block.NumberU64(), validateNodes[0].Index, t, secretKeys[0])
	assert.IsType(t, &DuplicateViewChangeEvidence{}, pool.AddViewChange(vc, validateNodes[0]))

	assert.Len(t, pool.Evidences(), 1)

	// test json
	evdata := ClassifyEvidence(pool.Evidences())
	b, _ := json.MarshalIndent(evdata, "", " ")
	t.Log(string(b))
	var ed2 EvidenceData
	assert.Nil(t, json.Unmarshal(b, &ed2))

	b2, _ := json.MarshalIndent(ed2, "", " ")
	assert.Equal(t, b, b2)
}

func TestJson(t *testing.T) {
	validateNodes, secretKeys := createValidateNode(1)

	pb := makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[0].Index, t, secretKeys[0])
	evidencePrepare, _ := NewEvidencePrepare(pb, validateNodes[0])

	block1 := newBlock(1)
	pv := makePrepareVote(1, 1, block1.Hash(), block1.NumberU64(), 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceVote, _ := NewEvidenceVote(pv, validateNodes[0])

	block2 := newBlock(1)
	vc := makeViewChange(1, 1, block2.Hash(), block2.NumberU64(), validateNodes[0].Index, t, secretKeys[0])
	evidenceView, _ := NewEvidenceView(vc, validateNodes[0])

	evs := []consensus.Evidence{
		&DuplicatePrepareBlockEvidence{
			PrepareA: evidencePrepare,
			PrepareB: evidencePrepare,
		},
		&DuplicatePrepareVoteEvidence{
			VoteA: evidenceVote,
			VoteB: evidenceVote,
		},
		&DuplicateViewChangeEvidence{
			ViewA: evidenceView,
			ViewB: evidenceView,
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

	evidences, err := NewEvidences(string(b2))
	assert.Nil(t, err)
	assert.Equal(t, 3, evidences.Len())

	testNewEvidence(t, evs)
}

func testNewEvidence(t *testing.T, evs []consensus.Evidence) {
	for _, e := range evs {
		b, _ := json.MarshalIndent(e, "", " ")
		switch e.(type) {
		case *DuplicatePrepareBlockEvidence:
			prepare, err := NewEvidence(DuplicatePrepareBlockType, string(b))
			assert.Nil(t, err)
			assert.NotNil(t, prepare)

		case *DuplicatePrepareVoteEvidence:
			vote, err := NewEvidence(DuplicatePrepareVoteType, string(b))
			assert.Nil(t, err)
			assert.NotNil(t, vote)

		case *DuplicateViewChangeEvidence:
			view, err := NewEvidence(DuplicateViewChangeType, string(b))
			assert.Nil(t, err)
			assert.NotNil(t, view)
		}
	}
}

func TestDuplicatePrepareBlockEvidence_Equal(t *testing.T) {
	validateNodes, secretKeys := createValidateNode(1)

	pbA := makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[0].Index, t, secretKeys[0])
	evidencePrepareA, _ := NewEvidencePrepare(pbA, validateNodes[0])

	pbB := makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[0].Index, t, secretKeys[0])
	evidencePrepareB, _ := NewEvidencePrepare(pbB, validateNodes[0])

	pbC := makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[0].Index, t, secretKeys[0])
	evidencePrepareC, _ := NewEvidencePrepare(pbC, validateNodes[0])

	p1 := &DuplicatePrepareBlockEvidence{
		PrepareA: evidencePrepareA,
		PrepareB: evidencePrepareB,
	}

	p2 := &DuplicatePrepareBlockEvidence{
		PrepareA: evidencePrepareA,
		PrepareB: evidencePrepareC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}

func TestDuplicatePrepareVoteEvidence_Equal(t *testing.T) {
	validateNodes, secretKeys := createValidateNode(1)

	pvA := makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceVoteA, _ := NewEvidenceVote(pvA, validateNodes[0])

	pvB := makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceVoteB, _ := NewEvidenceVote(pvB, validateNodes[0])

	pvC := makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceVoteC, _ := NewEvidenceVote(pvC, validateNodes[0])

	p1 := &DuplicatePrepareVoteEvidence{
		VoteA: evidenceVoteA,
		VoteB: evidenceVoteB,
	}

	p2 := &DuplicatePrepareVoteEvidence{
		VoteA: evidenceVoteA,
		VoteB: evidenceVoteC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}

func TestDuplicateViewChangeEvidence_Equal(t *testing.T) {
	validateNodes, secretKeys := createValidateNode(1)

	vcA := makeViewChange(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceViewA, _ := NewEvidenceView(vcA, validateNodes[0])

	vcB := makeViewChange(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceViewB, _ := NewEvidenceView(vcB, validateNodes[0])

	vcC := makeViewChange(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceViewC, _ := NewEvidenceView(vcC, validateNodes[0])

	p1 := &DuplicateViewChangeEvidence{
		ViewA: evidenceViewA,
		ViewB: evidenceViewB,
	}

	p2 := &DuplicateViewChangeEvidence{
		ViewA: evidenceViewA,
		ViewB: evidenceViewC,
	}

	assert.True(t, p1.Equal(p1))
	assert.False(t, p1.Equal(p2))
}

func TestDuplicatePrepareBlockEvidence_Validate(t *testing.T) {
	validateNodes, secretKeys := createValidateNode(2)

	pbA := makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[0].Index, t, secretKeys[0])
	evidencePrepareA, _ := NewEvidencePrepare(pbA, validateNodes[0])

	pbB := makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[0].Index, t, secretKeys[0])
	evidencePrepareB, _ := NewEvidencePrepare(pbB, validateNodes[0])

	d := &DuplicatePrepareBlockEvidence{
		PrepareA: evidencePrepareA,
		PrepareB: evidencePrepareB,
	}
	assert.Nil(t, d.Validate())

	pbB = makePrepareBlock(1, 1, newBlock(1), 1, validateNodes[1].Index, t, secretKeys[1])
	evidencePrepareB, _ = NewEvidencePrepare(pbB, validateNodes[1])
	d = &DuplicatePrepareBlockEvidence{
		PrepareA: evidencePrepareA,
		PrepareB: evidencePrepareB,
	}
	assert.NotNil(t, d.Validate())
}

func TestDuplicatePrepareVoteEvidence_Validate(t *testing.T) {
	validateNodes, secretKeys := createValidateNode(2)

	pvA := makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceVoteA, _ := NewEvidenceVote(pvA, validateNodes[0])

	pvB := makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceVoteB, _ := NewEvidenceVote(pvB, validateNodes[0])

	d := &DuplicatePrepareVoteEvidence{
		VoteA: evidenceVoteA,
		VoteB: evidenceVoteB,
	}
	assert.Nil(t, d.Validate())

	pvB = makePrepareVote(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, 1, validateNodes[1].Index, t, secretKeys[1])
	evidenceVoteB, _ = NewEvidenceVote(pvB, validateNodes[1])
	d = &DuplicatePrepareVoteEvidence{
		VoteA: evidenceVoteA,
		VoteB: evidenceVoteB,
	}
	assert.NotNil(t, d.Validate())
}

func TestDuplicateViewChangeEvidence_Validate(t *testing.T) {
	validateNodes, secretKeys := createValidateNode(2)

	hash := common.BytesToHash(utils.Rand32Bytes(32))
	vcA := makeViewChange(1, 1, hash, 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceViewA, _ := NewEvidenceView(vcA, validateNodes[0])

	vcB := makeViewChange(1, 1, hash, 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceViewB, _ := NewEvidenceView(vcB, validateNodes[0])

	d := &DuplicateViewChangeEvidence{
		ViewA: evidenceViewA,
		ViewB: evidenceViewB,
	}
	assert.NotNil(t, d.Validate())

	// different validater
	vcB = makeViewChange(1, 1, hash, 1, validateNodes[1].Index, t, secretKeys[1])
	evidenceViewB, _ = NewEvidenceView(vcB, validateNodes[1])

	d = &DuplicateViewChangeEvidence{
		ViewA: evidenceViewA,
		ViewB: evidenceViewB,
	}
	assert.NotNil(t, d.Validate())

	// different number
	vcB = makeViewChange(1, 1, hash, 2, validateNodes[0].Index, t, secretKeys[0])
	evidenceViewB, _ = NewEvidenceView(vcB, validateNodes[0])

	d = &DuplicateViewChangeEvidence{
		ViewA: evidenceViewA,
		ViewB: evidenceViewB,
	}
	assert.Nil(t, d.Validate())

	// different hash
	vcB = makeViewChange(1, 1, common.BytesToHash(utils.Rand32Bytes(32)), 1, validateNodes[0].Index, t, secretKeys[0])
	evidenceViewB, _ = NewEvidenceView(vcB, validateNodes[0])

	d = &DuplicateViewChangeEvidence{
		ViewA: evidenceViewA,
		ViewB: evidenceViewB,
	}
	assert.Nil(t, d.Validate())
}
