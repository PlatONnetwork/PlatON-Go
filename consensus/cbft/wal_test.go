package cbft

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"hash/crc32"
	"io"
	"os"
	"testing"
	"time"
)

type msgInfoa struct {
	Msg    *prepareBlock
	PeerID discover.NodeID
}

type JournalMessagev struct {
	Timestamp uint64
	Data      *msgInfoa
}

func TestWalWrite(t *testing.T) {
	wal, _ := NewWal(nil)
	var err error

	// test rotate
	//time.Sleep(6 * time.Second)

	// UpdateViewChange
	wal.UpdateViewChange(&ViewChangeMessage{
		Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065747"),
		Number: 110,
	})
	// WriteJournal
	beginTime := uint64(time.Now().UnixNano())
	countW := 0
	for i := 0; i < 3000000; i++ {
		peerId, _ := discover.HexID("b6c8c9f99bfebfa4fb174df720b9385dbd398de699ec36750af3f38f8e310d4f0b90447acbef64bdf924c4b59280f3d42bb256e6123b53e9a7e99e4c432549d6")
		if i%2 == 0 {
			viewChangeVotes := make([]*viewChangeVote, 0)
			viewChangeVotes = append(viewChangeVotes, &viewChangeVote{
				Timestamp:      uint64(time.Now().UnixNano()),
				BlockNum:       111,
				BlockHash:      common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065747"),
				ProposalIndex:  1111,
				ProposalAddr:   common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185822"),
				ValidatorIndex: 11111,
				ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185822"),
			})
			viewChangeVotes = append(viewChangeVotes, &viewChangeVote{
				Timestamp:      uint64(time.Now().UnixNano()),
				BlockNum:       222,
				BlockHash:      common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065749"),
				ProposalIndex:  2222,
				ProposalAddr:   common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185829"),
				ValidatorIndex: 22222,
				ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185829"),
			})
			err = wal.Write(&MsgInfo{
				Msg: &prepareBlock{
					Timestamp:     uint64(time.Now().UnixNano()),
					ProposalIndex: 666,
					View: &viewChange{
						Timestamp:     uint64(time.Now().UnixNano()),
						ProposalIndex: 12,
						ProposalAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185822"),
						BaseBlockNum:  10086,
						BaseBlockHash: common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
					},
					ViewChangeVotes: viewChangeVotes,
				},
				PeerID: peerId,
			})
		} else if i%3 == 0 {
			//wal.Write(&MsgInfo{
			//	Msg: &prepareBlockHash{
			//		Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065747"),
			//		Number: 13333,
			//	},
			//	PeerID: peerId,
			//})
			pvs := make([]*prepareVote, 0)
			pvs = append(pvs, &prepareVote{
				Timestamp: uint64(time.Now().UnixNano()),
				Number:    7777,
			})
			votes := make([]*prepareVotes, 0)
			votes = append(votes, &prepareVotes{
				Hash:   common.HexToHash("0x76fded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
				Number: 5678,
				Votes:  pvs,
			})
			votes = append(votes, &prepareVotes{
				Hash:   common.HexToHash("0x76fded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
				Number: 6789,
				Votes:  pvs,
			})
			err = wal.Write(&MsgInfo{
				Msg: &highestPrepareBlock{
					Votes: votes,
				},
				PeerID: peerId,
			})
		} else if i%5 == 0 {
			err = wal.Write(&MsgInfo{
				Msg: &prepareVote{
					Timestamp:      uint64(time.Now().UnixNano()),
					Hash:           common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066326"),
					Number:         16666,
					ValidatorIndex: 1,
					ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185821"),
				},
				PeerID: peerId,
			})
		} else if i%7 == 0 {
			votes := make([]*prepareVote, 0)
			votes = append(votes, &prepareVote{
				Timestamp:      uint64(time.Now().UnixNano()),
				Hash:           common.HexToHash("0x76fded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
				Number:         8877,
				ValidatorIndex: 9900,
				ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185827"),
			})
			votes = append(votes, &prepareVote{
				Timestamp:      uint64(time.Now().UnixNano()),
				Hash:           common.HexToHash("0x76fded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
				Number:         8878,
				ValidatorIndex: 9901,
				ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185828"),
			})
			err = wal.Write(&MsgInfo{
				Msg: &prepareVotes{
					Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
					Number: 7788,
					Votes:  votes,
				},
				PeerID: peerId,
			})
		} else {
			err = wal.Write(&MsgInfo{
				Msg: &viewChange{
					Timestamp:     uint64(time.Now().UnixNano()),
					ProposalIndex: 12,
					ProposalAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185822"),
					BaseBlockNum:  10086,
					BaseBlockHash: common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
				},
				PeerID: peerId,
			})
		}
		if err != nil {
			fmt.Println("write error", err)
			panic(err)
		}
		countW ++
	}
	wal.Close() // force flush
	fmt.Println("write total msg info", countW)
	endTime := uint64(time.Now().UnixNano())
	fmt.Println("write elapsed time", endTime-beginTime)
}

func TestWalLoad(t *testing.T) {
	wal, _ := NewWal(nil)
	var err error

	// LoadJournal
	beginTime := uint64(time.Now().UnixNano())
	countR := 0
	err = wal.Load(func(info *MsgInfo) {
		countR ++
		//fmt.Printf("info=%#v\n", info)
	})
	if err != nil {
		fmt.Println("load error", err)
		//panic(err)
	}
	endTime := uint64(time.Now().UnixNano())
	fmt.Println("total msg info", countR)
	fmt.Println("load elapsed time", endTime-beginTime)

}

func TestLevelDB(t *testing.T) {
	db, err := leveldb.OpenFile("C:\\Users\\jungle\\Desktop\\wal.tar\\wal\\wal_meta", nil)
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

func TestBufferRead(t *testing.T) {
	file, err := os.Open("D://data/platon/wal/wal.1")
	if err != nil {
		fmt.Println(err)
	}

	bufReader := bufio.NewReaderSize(file, 1024)
	fmt.Printf("%d \n", bufReader.Buffered())
	//bufReader.Discard(400)

	for {
		index, _ := bufReader.Peek(12)
		fmt.Printf("%d \n", bufReader.Buffered())
		crc := binary.BigEndian.Uint32(index[0:4])
		length := binary.BigEndian.Uint32(index[4:8])
		msgType := binary.BigEndian.Uint32(index[8:12])
		fmt.Println(msgType)

		pack := make([]byte, length+12)
		var (
			totalNum = 0
			readNum  = 0
		)
		for totalNum, err = 0, error(nil); err == nil && uint32(totalNum) < length+12; {
			readNum, err = bufReader.Read(pack[totalNum:])
			totalNum = totalNum + readNum
		}

		fmt.Printf("%d \n", bufReader.Buffered())
		if 0 == readNum {
			break
		}

		crcc := crc32.Checksum(pack[12:], crc32c)
		if crc != crcc {
			panic("check crc error")
		}

		var v JournalMessagev
		if err := rlp.DecodeBytes(pack[12:], &v); err == nil {
			fmt.Println("Timestamp", v.Timestamp)
			fmt.Println("Data", v.Data)
			dd := &MsgInfo{
				Msg:    v.Data.Msg,
				PeerID: v.Data.PeerID,
			}
			fmt.Println("", dd.Msg.(*prepareBlock).Timestamp)
		} else {
			panic(err)
		}

		fmt.Printf("%d \n", bufReader.Buffered())
		if err != nil && err != io.EOF {
			panic(err)
		}
	}
}
