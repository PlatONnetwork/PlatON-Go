package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageType(t *testing.T) {
	view := &viewChange{}

	assert.Equal(t, MessageType(view), uint64(ViewChangeMsg))
}

func TestPrepareBlock(t *testing.T) {
	//Timestamp       uint64
	//Block           *types.Block
	//ProposalIndex   uint32            `json:"proposal_index"`
	//ProposalAddr    common.Address    `json:"proposal_address"`
	//View            *viewChange       `json:"view"`
	//ViewChangeVotes []*viewChangeVote `json:"viewchange_votes"`
	pb := prepareBlock{
		Timestamp:     1,
		ProposalIndex: 2,
	}

	buf, err := rlp.EncodeToBytes(&pb)
	if err != nil {
		t.Error(err)
	}
	var pb2 prepareBlock
	err = rlp.DecodeBytes(buf, &pb2)
	if err != nil {
		t.Error(err)
	}
	t.Log(pb2.Timestamp)
}
