package cbft

import (
	"bufio"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
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
				log: f,
				seq: seq,
			})
		}
	}
	sort.Sort(files)

	for _, i := range files {
		t.Log(i.log)
	}
}

func TestS(t *testing.T) {
	l := make([][]uint64, 0)
	d := []uint64{1, 2, 3}
	l = append(l, d)
	l = append(l, []uint64{4, 5, 6})
	b, err := rlp.EncodeToBytes(l)
	if err != nil {
		t.Error(err)
	}
	t.Log(hexutil.Encode(b))

	content, rest, _ := rlp.SplitList(b)
	t.Log(hexutil.Encode(content))
	t.Log(hexutil.Encode(rest))

	content1, rest1, _ := rlp.SplitList(content)
	t.Log(hexutil.Encode(content1))
	t.Log(hexutil.Encode(rest1))
	content2, rest2, _ := rlp.SplitList(content1)
	t.Log(hexutil.Encode(content2))
	t.Log(hexutil.Encode(rest2))
}

type viewChangeMeta struct {
	Number uint64
	Hash   common.Hash
	FileID uint32
	Seq    uint64
}

func TestLevelDB(t *testing.T) {
	db, err := leveldb.OpenFile("D://data/platon/wal/wal_meta", nil)
	if err == nil {
		//encode, err := rlp.EncodeToBytes(&viewChangeMeta{
		//	Number: 111,
		//	Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065747"),
		//	FileID: 222,
		//	Seq:    9,
		//})
		//db.Put([]byte("view-change"), encode, &opt.WriteOptions{Sync:true})

		data, err := db.Get([]byte("view-change"), nil)
		if err == nil {
			var v viewChangeMeta
			if err := rlp.DecodeBytes(data, &v); err == nil {
				fmt.Println(v.Number)
				fmt.Println(v.Hash.Hex())
				fmt.Println(v.FileID)
				fmt.Println(v.Seq)
			}
		}
	}
}

func TestJournal(t *testing.T) {
	path := "D://data/platon/wala"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			fmt.Println("could not create directory")
		}
	}

	file, err := os.OpenFile("D://data/platon/wala/wal.1", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

	if err == nil {
		buf := bufio.NewWriterSize(file, 1024)
		fmt.Println(buf.Available(), buf.Buffered()) // 1024 0	Available 返回缓存中的可以空间 / Buffered 返回缓存中未提交的数据长度

		writeWrapper := &writeWrapper{
			writer: buf,
		}
		for i := 0; i < 10000; i++ {
			go func() {
				writeWrapper.write([]byte("hello world "))
				//buf.Write([]byte("hello world "))
			}()
		}
		time.Sleep(5 * time.Second)
		buf.Flush()

		//buf.Write([]byte("hello world "))
		//fmt.Println(buf.Available(), buf.Buffered()) // 1012 12
		//input, _ := os.Open("D://data/platon/wala/wal.1")
		//contentByte, _ := ioutil.ReadAll(input)
		//fmt.Println(string(contentByte))
		//buf.Write([]byte("hello world2 "))
		//buf.Flush()
	}

	//if err := journal.load(pool.AddLocals); err != nil {
	//	log.Warn("Failed to load transaction journal", "err", err)
	//}
	//if err := pool.journal.rotate(pool.local()); err != nil {
	//	log.Warn("Failed to rotate transaction journal", "err", err)
	//}
}

type writeWrapper struct {
	writer io.Writer
	mu     sync.RWMutex
}

func (w *writeWrapper) write(b []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.writer.Write(b)
}

func TestListDir(t *testing.T) {
	// 不递归遍历指定目录下的文件列表
	files, _ := ioutil.ReadDir("D://data/platon/wala")
	for _, f := range files {
		fmt.Println(f.Name())
	}

	fmt.Println(fmt.Sprintf("wal.%d", uint32(789)))

	// 查询指定文件的大小
	fileInfo, _ := os.Stat("D://data/platon/wala/wal.3")
	fmt.Println(fileInfo.Size()) // byte
}
