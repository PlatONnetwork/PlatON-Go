package evidence

import (
	"encoding/json"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"

	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	prepare := buildPrepareBlock()
	vote := buildPrepareVote()
	view := buildViewChange()

	evs := []consensus.Evidence{
		&DuplicatePrepareBlockEvidence{
			PrepareA: prepare,
			PrepareB: prepare,
		},
		&DuplicatePrepareVoteEvidence{
			VoteA: vote,
			VoteB: vote,
		},
		&DuplicateViewChangeEvidence{
			ViewA: view,
			ViewB: view,
		},
	}
	ed := ClassifyEvidence(evs)
	b, _ := json.MarshalIndent(ed, "", " ")
	t.Log(string(b))
	var ed2 EvidenceData
	assert.Nil(t, json.Unmarshal(b, &ed2))

	b2, _ := json.MarshalIndent(ed2, "", " ")
	assert.Equal(t, b, b2)
}
