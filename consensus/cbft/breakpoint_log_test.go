package cbft

import (
	"context"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"os"
	"testing"
)

func TestLogPrepareBP_ReceiveVote(t *testing.T) {
	log.Root().SetHandler(log.StdoutHandler)
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 1)
	ctx := context.WithValue(context.TODO(), "peer", "xxxxxxxx")
	pvote := makeFakePrepareVote()
	logBP.PrepareBP().ReceiveVote(ctx, pvote, engine)
	logBP.PrepareBP().AcceptVote(ctx, pvote, engine)
	logBP.PrepareBP().CacheVote(ctx, pvote, engine)
	logBP.PrepareBP().DiscardVote(ctx, pvote, engine)
	logBP.PrepareBP().SendPrepareVote(ctx, pvote, engine)
	logBP.PrepareBP().InvalidVote(ctx, pvote, fmt.Errorf("fail:%v",""), engine)
	logBP.PrepareBP().TwoThirdVotes(ctx, pvote, engine)
	t.Log("done")
}

func BenchmarkLogInternalBP_ExecuteBlock(b *testing.B) {
	l := logBP
	for i := 0; i < 10000; i++ {
		l.InternalBP().ExecuteBlock(nil, common.BytesToHash(Rand32Bytes(32)), 40, 100, 40)
	}
}
