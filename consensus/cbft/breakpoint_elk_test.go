package cbft

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElk_Log(t *testing.T) {

	assert.NotNil(t, elkBP.PrepareBP())
	assert.NotNil(t, elkBP.ViewChangeBP())
	assert.NotNil(t, elkBP.SyncBlockBP())
	assert.NotNil(t, elkBP.InternalBP())

	ctx := context.WithValue(context.TODO(), "peer", randomID())

	// elkPrepareBP
	elkp := elkBP.PrepareBP()

	elkp.ReceiveBlock(ctx, makeFakePrepareBlock(), nil)
	elkp.ReceiveVote(ctx, makeFakePrepareVote(), nil)
	elkp.AcceptBlock(ctx, makeFakePrepareBlock(), nil)
	elkp.DiscardBlock(ctx, makeFakePrepareBlock(), nil)
	elkp.AcceptVote(ctx, makeFakePrepareVote(), nil)
	elkp.DiscardVote(ctx, makeFakePrepareVote(), nil)
	elkp.InvalidBlock(ctx, makeFakePrepareBlock(), nil,nil)
	elkp.InvalidVote(ctx, makeFakePrepareVote(), nil,nil)
	elkp.InvalidViewChangeVote(ctx, makeFakePrepareBlock(), nil,nil)

	elkview := elkBP.ViewChangeBP()
	elkview.ReceiveViewChange(ctx, makeFakeViewChange(), nil)
	elkview.ReceiveViewChangeVote(ctx, makeFakeViewChangeVote(), nil)
	elkview.InvalidViewChange(ctx, makeFakeViewChange(), nil,nil)
	elkview.InvalidViewChangeVote(ctx, makeFakeViewChangeVote(), nil,nil)
	elkview.InvalidViewChangeBlock(ctx, makeFakeViewChange(), nil)

	//elksync := elkBP.SyncBlockBP()

}