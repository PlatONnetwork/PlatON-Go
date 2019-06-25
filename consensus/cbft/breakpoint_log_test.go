package cbft

import (
	"context"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestLogPrepareBP_ReceiveVote(t *testing.T) {
	p, _ := ioutil.TempDir(os.TempDir(), "test")
	logBP, _ = NewLogBP(p)
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
	logBP.PrepareBP().InvalidVote(ctx, pvote, fmt.Errorf("fail:%v", ""), engine)
	logBP.PrepareBP().TwoThirdVotes(ctx, pvote, engine)
	t.Log("done")
}

func TestJsonLog(b *testing.T) {
	p, _ := ioutil.TempDir(os.TempDir(), "test")
	logBP, _ = NewLogBP(p)
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 1)
	ctx := context.WithValue(context.TODO(), "peer", "xxxxxxxx")
	pvote := makeFakePrepareVote()

	start := time.Now()
	n := 100000
	for i := 0; i < n; i++ {
		logBP.PrepareBP().ReceiveVote(ctx, pvote, engine)
		logBP.PrepareBP().AcceptVote(ctx, pvote, engine)
		logBP.PrepareBP().CacheVote(ctx, pvote, engine)
		logBP.PrepareBP().DiscardVote(ctx, pvote, engine)
		logBP.PrepareBP().SendPrepareVote(ctx, pvote, engine)
		logBP.PrepareBP().InvalidVote(ctx, pvote, fmt.Errorf("fail:%v", ""), engine)
		logBP.PrepareBP().TwoThirdVotes(ctx, pvote, engine)
		logBP.InternalBP().ExecuteBlock(nil, common.BytesToHash(Rand32Bytes(32)), 40, 100, 40)
	}
	b.Log(fmt.Sprintf("tps:%v/ns", time.Since(start).Nanoseconds()/int64(n)/8))
	logBP.Close()
}
