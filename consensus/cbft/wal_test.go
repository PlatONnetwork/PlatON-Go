package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	viewChangeNumber = uint64(100)
	viewChangeHash   = common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065747")
	times            = 3000000
	tempDir          string
	wal              Wal
)

func TestMain(m *testing.M) {
	fmt.Println("begin test wal")
	tempDir, _ = ioutil.TempDir("", "wal")
	m.Run()
	os.RemoveAll(tempDir)
	fmt.Println("end test wal")
}

func getWal() Wal {
	if wal == nil {
		wal, _ = NewWal(nil, tempDir)
	}
	return wal
}

func TestWalUpdateViewChange(t *testing.T) {
	// UpdateViewChange
	getWal().UpdateViewChange(&ViewChangeMessage{
		Hash:   viewChangeHash,
		Number: viewChangeNumber,
	})
}

func TestWalWrite(t *testing.T) {
	var err error
	// WriteJournal
	beginTime := uint64(time.Now().UnixNano())
	count := 0
	header := &types.Header{
		Number: big.NewInt(1),
	}
	block := types.NewBlock(header, nil, nil)
	for i := 0; i < times; i++ {
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
			err = getWal().Write(&MsgInfo{
				Msg: &prepareBlock{
					Timestamp:     uint64(time.Now().UnixNano()),
					Block:         block,
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
			//getWal().Write(&MsgInfo{
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
			err = getWal().Write(&MsgInfo{
				Msg: &highestPrepareBlock{
					Votes: votes,
				},
				PeerID: peerId,
			})
		} else if i%5 == 0 {
			err = getWal().Write(&MsgInfo{
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
			err = getWal().Write(&MsgInfo{
				Msg: &prepareVotes{
					Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
					Number: 7788,
					Votes:  votes,
				},
				PeerID: peerId,
			})
		} else {
			err = getWal().Write(&MsgInfo{
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
			t.Errorf("%s", "write error")
		}
		count ++
	}
	getWal().Close() // force flush
	wal = nil
	fmt.Println("write total msg info", count)
	if count != times {
		t.Errorf("%s", "write error")
	}
	endTime := uint64(time.Now().UnixNano())
	fmt.Println("write elapsed time", endTime-beginTime)
}

func TestWalLoad(t *testing.T) {
	var err error
	// LoadJournal
	beginTime := uint64(time.Now().UnixNano())
	count := 0
	err = getWal().Load(func(info *MsgInfo) {
		count ++
	})
	if err != nil {
		fmt.Println("load error", err)
		t.Errorf("%s", "load error")
	}
	getWal().Close() // force flush
	wal = nil
	fmt.Println("total msg info", count)
	if count != times {
		t.Errorf("%s", "load error")
	}
	endTime := uint64(time.Now().UnixNano())
	fmt.Println("load elapsed time", endTime-beginTime)

}

func TestLevelDB(t *testing.T) {
	path := filepath.Join(tempDir, "wal_meta")
	if db, err := leveldb.OpenFile(path, nil); err != nil {
		t.Errorf("%s", "TestLevelDB error")
	} else {
		data, err := db.Get([]byte("view-change"), nil)
		if err == nil {
			var v ViewChangeMeta
			if err := rlp.DecodeBytes(data, &v); err == nil {
				fmt.Println(v.Number)
				fmt.Println(v.Hash.Hex())
				fmt.Println(v.FileID)
				fmt.Println(v.Seq)
				db.Close()
				if v.Number != 100 || v.Hash.Hex() != viewChangeHash.Hex() {
					t.Errorf("%s", "TestLevelDB error")
				}
			}
		} else {
			db.Close()
			t.Errorf("%s", "TestLevelDB error")
		}
	}
}