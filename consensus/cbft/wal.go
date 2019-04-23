package cbft

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"hash/crc32"
	"io"
	"os"
	"sync"
	"time"
)

//WAL: write consensus message to log file before handle it.
//Execute message from WAL before start cbft state machine
// WAL format:
//| CRC | Length |  Data |
//| ------ | ------ | ------ |
//| 4byte | 4byte |  n byte |

// WAL Meta (need fsync)
// MinLogFileID uint64
// MaxLogFileID uint64
// ViewChange (FileID, FileSequence)
var (
	metaDBName = "wal-meta"
	walPrefix  = "wal"
	MinLogKey = []byte("MinLog")
	MaxLogKey = []byte("MaxLog")
	ViewChangeKey = []byte("ViewChange")
 	crc32c = crc32.MakeTable(crc32.Castagnoli)

)

type WALMessage interface{}

type TimedWALMessage struct {
	Timestamp uint64
	Msg       WALMessage
}

type ViewChangeWALMessage struct {
	Hash   common.Hash
	Number uint64
}

type WAL interface {
	Write(msg WALMessage) error
	WriteSync(msg WALMessage) error
	LocationHeight(hash common.Hash, height uint64) (io.Reader, error)
	Start()
	Stop()
}

type baseWAL struct {
	file *WalFile
	enc  *WALEncode
}

func NewWAL(path string, options *WalOptions) (WAL, error) {
	file, err := NewWalFile(path, options)
	if err != nil {
		return nil, err
	}

	return &baseWAL{
		file: file,
		enc:  NewBaseEncode(file),
	}, nil
}

func (b *baseWAL) Write(msg WALMessage) error {
	return b.enc.Encode(msg)

}

func (b *baseWAL) WriteSync(msg WALMessage) error {
	switch m := msg.(type) {
	case *ViewChangeWALMessage:
		//CommittedWALMessage is special msg, it indicate cbft had commit block.
		//Update meta info when commit message was wrote disk
		return b.file.UpdateViewChangeMeta(m)
	}

	if err := b.Write(msg); err != nil {
		return err
	}
	return b.file.Flush()
}

//Start async write loop
//Notice: need replay wal message when wal started.
func (b *baseWAL) Start() {
	b.file.Start()
}

//Stop async write
func (b *baseWAL) Stop() {
	b.file.Stop()
}

//Used replay, location wal position
func (b *baseWAL) LocationHeight(hash common.Hash, height uint64) (io.Reader, error) {
	return nil, nil
}

type WALEncode struct {
	wr io.Writer
}

func NewBaseEncode(wr io.Writer) *WALEncode {
	return &WALEncode{
		wr: wr,
	}
}

func (w *WALEncode) Encode(msg WALMessage) error {
	tw := &TimedWALMessage{
		Timestamp:uint64(time.Now().UnixNano()),
		Msg:msg,
	}

	data, err := rlp.EncodeToBytes(tw)
	if err != nil {
		return err
	}

	crc := crc32.Checksum(data, crc32c)
	length := uint32(len(data))
	totalLength := 8 + int(length)

	buf := make([]byte, totalLength)
	binary.BigEndian.PutUint32(buf[0:4], crc)
	binary.BigEndian.PutUint32(buf[4:8], length)
	copy(buf[8:], data)
	_, err = w.wr.Write(buf)
	return err
}

type WALDecode struct {
	io.Reader
}

type WalOptions struct {
	bufferLimitSize int
	fileLimitSize   int64
	syncLoopDuration time.Duration
}

type WalFile struct {
	ID       uint64
	Path     string
	minID uint64
	writeBuf *bufio.Writer

	metaDB  *leveldb.DB
	options *WalOptions
	file    *os.File

	fileSize  int64
	end *endViewChange

	mux    sync.Mutex
	ticker *time.Ticker
	quit   chan struct{}
}

type endViewChange struct {
	Hash common.Hash
	Number uint64
	FileID uint64
	Seq uint64
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return fmt.Errorf("could not create directory %v. %v", path, err)
		}
	}
	return nil
}

type logSeq []uint64

func (w logSeq) Len() int {
	return len(w)
}

func (p logSeq) Less(i, j int) bool {
	return p[i] < p[j]
}
func (p logSeq) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

//func findWALFile(path string) []uint64 {
//	files := make(logSeq, 0)
//	info, err := ioutil.ReadDir(path)
//	if err != nil {
//		return files
//	}
//
//	reg := regexp.MustCompile("^wal.([1-9][0-9]*)$")
//	regNum := regexp.MustCompile("([1-9][0-9]*)$")
//	for _, f := range info {
//		if reg.MatchString(f.Name()) {
//			if seq, err :=strconv.ParseUint(regNum.FindString(f.Name()), 10, 64); err == nil {
//				files = append(files, seq)
//			}
//		}
//	}
//
//	sort.Sort(files)
//	return files
//}

func NewWalFile(path string, options *WalOptions) (*WalFile, error) {
	walFile := &WalFile{
		ID:0,
		Path:path,
		minID:0,
		options:options,
		fileSize:0,
		end : nil,
		quit:make(chan struct{}),
	}

	ldbName := fmt.Sprintf("%s/%s", path, metaDBName)

	initEmptyWalFile := func() error{
		if err := ensureDir(path); err != nil {
			return err
		}
		if file, err := os.OpenFile(fmt.Sprintf("%s/wal.0", path ),  os.O_WRONLY , 0700); err == nil {
			walFile.file = file
		} else {
			return fmt.Errorf("create wal log failed %v. %v", ldbName, err)
		}

		if db, err := leveldb.OpenFile(ldbName, nil); err != nil {
			walFile.metaDB = db
		} else {
			return fmt.Errorf("create db failed %v. %v", ldbName, err)
		}

		walFile.writeBuf = bufio.NewWriterSize(walFile.file, options.bufferLimitSize)
		return nil
	}

	clearLog := func() error{
		return os.RemoveAll(path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := initEmptyWalFile(); err != nil {
			return nil,  fmt.Errorf("init wal log failed %v.%v", path, err)
		}
		return walFile, nil
	}


	if _, err := os.Stat(ldbName); os.IsNotExist(err) {
		if err = clearLog(); err != nil {
			return nil, fmt.Errorf("clear wal log failed %v.%v", path, err)
		}
		if err = initEmptyWalFile(); err != nil {
			return nil, fmt.Errorf("init wal log failed %v.%v", path, err)
		}
		return walFile, nil
	}




	ldb, err := leveldb.OpenFile(ldbName, nil)
	if err != nil {
		return nil, fmt.Errorf("create db failed %v. %v", ldbName, err)
	}

	//get all meta from database, return error and note clear wal log if get failed
	if min, err := ldb.Get(MinLogKey, nil); err == nil {
		walFile.minID = binary.BigEndian.Uint64(min)
	} else {
		return nil, fmt.Errorf("meta lost, please clear all %v", path)
	}

	if max, err := ldb.Get(MaxLogKey, nil); err == nil {
		walFile.ID = binary.BigEndian.Uint64(max)
		//make sure had wal log
		if file, err := os.OpenFile(fmt.Sprintf("%s/wal.0", path),  os.O_WRONLY , 0700); err == nil {
			walFile.file = file
			walFile.writeBuf = bufio.NewWriterSize(walFile.file, options.bufferLimitSize)
			stat , err := file.Stat()
			if err != nil {
				return nil, err
			}
			walFile.fileSize = stat.Size()
		} else {
			return nil, fmt.Errorf("create wal log failed %v. %v", ldbName, err)
		}
	} else {
		return nil, fmt.Errorf("meta lost, please clear all %v", path)
	}

	value, err := ldb.Get(ViewChangeKey, nil)
	var end endViewChange
	if err := rlp.DecodeBytes(value, &end); err != nil {
		return nil,  fmt.Errorf("meta destroy, please clear all %v", path)
	}

	return walFile, nil
}


func (w *WalFile) fileName(id uint64) string {
	return fmt.Sprintf("%s/wal.%d", w.Path, id)
}

func (w *WalFile) Write(p []byte) (n int, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	return 0, nil
}

func (w *WalFile) Flush() error {
	return  nil
}

func (w *WalFile) WriteSync(p []byte) (n int, err error) {
	return 0, nil
}

func (w *WalFile) writeLoop() {
	for {
		select {
		case <-w.ticker.C:
			w.checkFileSize()
		case <-w.quit:
			return
		}
	}
}
func (w *WalFile) UpdateViewChangeMeta(c *ViewChangeWALMessage) error {
	data , err := rlp.EncodeToBytes(&endViewChange{Hash:c.Hash, Number:c.Number})
	if err != nil {
		return err
	}
	return w.metaDB.Put(ViewChangeKey, data, &opt.WriteOptions{Sync:true})
}

func (w *WalFile) checkFileSize() {
	w.mux.Lock()
	defer w.mux.Unlock()
	if w.fileSize > w.options.fileLimitSize {
		w.rotateFile()
	}
}

func (w *WalFile) rotateFile() {

}

func (w *WalFile) Start() {
	w.ticker = time.NewTicker(w.options.syncLoopDuration)
	go w.writeLoop()
}

func (w *WalFile) Stop() {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.writeBuf.Flush()
	w.file.Sync()

	close(w.quit)
}
