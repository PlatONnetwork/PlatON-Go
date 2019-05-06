package cbft

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

type msgInfo struct {
	Id   uint64
	Name string
	Desc string
}

type JournalMessagev struct {
	Timestamp uint64
	Msg       msgInfo
}

func TestWal(t *testing.T) {
	wal, _ := NewWal(nil)

	for i := 0; i < 1000000; i++ {
		if i%2 == 0 {
			wal.Write(&msgInfo{
				Id:   uint64(i),
				Name: "platon" + strconv.Itoa(i),
				Desc: "platonDesc" + strconv.Itoa(i),
			})
		} else {
			wal.Write(&msgInfo{
				Id:   uint64(i),
				Name: "platon" + strconv.Itoa(i),
			})
		}

	}

	wal.Write(&ViewChangeMessage{
		Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065747"),
		Number: 110,
	})
}

func TestReadJournal(t *testing.T) {
	file, err := os.Open("D://data/platon/wal/wal.1")
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 44)
	ix := 0;
	for {
		//ReadAt从指定的偏移量开始读取，不会改变文件偏移量
		len, _ := file.ReadAt(buf, int64(ix))
		ix = ix + len
		if len == 0 {
			break
		}

		var v JournalMessagev
		if err := rlp.DecodeBytes(buf[8:], &v); err == nil {
			fmt.Println("Timestamp", v.Timestamp)
			fmt.Println("Msg Id", v.Msg.Id)
			fmt.Println("Msg Name", v.Msg.Name)
			fmt.Println("Msg Desc", v.Msg.Desc)
		} else {
			fmt.Println(err)
		}
	}
	file.Close()
}

func TestLevelDB(t *testing.T) {
	db, err := leveldb.OpenFile("D://data/platon/wal/wal_meta", nil)
	if err == nil {
		data, err := db.Get([]byte("view-change"), nil)
		if err == nil {
			var v ViewChangeMeta
			if err := rlp.DecodeBytes(data, &v); err == nil {
				fmt.Println(v.Number)
				fmt.Println(v.Hash.Hex())
				fmt.Println(v.FileID)
				fmt.Println(v.Seq)
			}
		}
	}
}

func TestReader(t *testing.T) {
	sr := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	buf := bufio.NewReaderSize(sr, 16)
	b := make([]byte, 20)
	//n, err := buf.Read(b)
	//fmt.Println(n)
	//fmt.Println(err)

	fmt.Println(buf.Buffered()) // 0
	s, _ := buf.Peek(5)
	s[0], s[1], s[2] = 'a', 'b', 'c'
	fmt.Printf("%d   %q\n", buf.Buffered(), s) // 16   "abcDE"

	buf.Discard(1)

	for n, err := 0, error(nil); err == nil; {
		n, err = buf.Read(b)
		fmt.Printf("%d   %q   %v\n", buf.Buffered(), b[:n], err)
	}
}

func TestReadJournal2(t *testing.T) {
	file, err := os.Open("D://data/platon/wal/wal.1")
	if err != nil {
		fmt.Println(err)
	}

	bufReader := bufio.NewReaderSize(file, 1024)
	fmt.Printf("%d \n", bufReader.Buffered())
	bufReader.Discard(400)

	for {
		index, _ := bufReader.Peek(8)
		fmt.Printf("%d \n", bufReader.Buffered())
		crc := binary.BigEndian.Uint32(index[0:4])
		length := binary.BigEndian.Uint32(index[4:8])

		pack := make([]byte, length+8)
		var (
			totalNum = 0
			readNum  = 0
		)
		for totalNum, err = 0, error(nil); err == nil && uint32(totalNum) < length+8; {
			readNum, err = bufReader.Read(pack[totalNum:])
			totalNum = totalNum + readNum
		}

		fmt.Printf("%d \n", bufReader.Buffered())
		if 0 == readNum {
			break
		}

		crcc := crc32.Checksum(pack[8:], crc32c)
		if crc != crcc {
			panic("check crc error")
		}
		var v JournalMessagev
		if err := rlp.DecodeBytes(pack[8:], &v); err == nil {
			fmt.Println("Timestamp", v.Timestamp)
			fmt.Println("Msg Id", v.Msg.Id)
			fmt.Println("Msg Name", v.Msg.Name)
			fmt.Println("Msg Desc", v.Msg.Desc)
		} else {
			panic(err)
		}

		fmt.Printf("%d \n", bufReader.Buffered())
		if err != nil && err != io.EOF {
			panic(err)
		}
	}
}
