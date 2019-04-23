package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"regexp"
	"sort"
	"strconv"
	"testing"
)

type WALLog struct {
	log string
	seq int
}

type WALLogs []WALLog

func (w WALLogs) Len() int {
	return len(w)
}

func (w WALLogs) Less(i, j int) bool {
	return w[i].seq < w[j].seq
}

func (w WALLogs) Swap(i, j int) { w[i], w[j] = w[j], w[i] }

func TestWalFile(t *testing.T) {
	reg := regexp.MustCompile("^wal.([1-9][0-9]*)$")
	regNum := regexp.MustCompile("([1-9][0-9]*)$")
	files := make(WALLogs, 0)
	for _, f := range []string{"wal.1", "wal.4555", "wal.2", "wal.4", "wal.10", "wal.8", "wal."} {
		if reg.MatchString(f) {
			seq, _ := strconv.Atoi(regNum.FindString(f))
			files = append(files, WALLog{
				log:f,
				seq:seq,
			})
		}
	}
	sort.Sort(files)

	for _, i := range files {
		t.Log(i.log)
	}
}



func TestS(t *testing.T)  {
	l := make([][]uint64, 0)
	d := []uint64{1,2,3}
	l = append(l, d)
	l = append(l, []uint64{4,5,6})
	b , err := rlp.EncodeToBytes(l)
	if err != nil {
		t.Error(err)
	}
	t.Log(hexutil.Encode(b))

	content, rest , _ := rlp.SplitList(b)
	t.Log(hexutil.Encode(content))
	t.Log(hexutil.Encode(rest))

	content1, rest1 , _ := rlp.SplitList(content)
	t.Log(hexutil.Encode(content1))
	t.Log(hexutil.Encode(rest1))
	content2, rest2 , _ := rlp.SplitList(content1)
	t.Log(hexutil.Encode(content2))
	t.Log(hexutil.Encode(rest2))

}
