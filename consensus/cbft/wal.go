package cbft

import (
	"bufio"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"io"
	"os"
	"sync"
	"time"
)

//WAL: write consensus message to log file before handle it.
//Execute message from WAL before start cbft state machine
// WAL format:
//| CRC | Length | Type | Data |
//| ------ | ------ | ------ |------ |
//| 4byte | 4byte | 1byte | n byte |

// WAL Meta (need fsync)
// MinLogFileID uint64
// MaxLogFileID uint64
// CommitedBlockNumber (FileID, FileSequence)
var (
	metaDBName = "walmeta"
)

type WALMessage interface{}

type WAL interface {
	Write(msg WALMessage)
	WriteSync(msg WALMessage)
	Start()
	Stop()
}

type baseWAL struct {
	file *WalFile
	enc  WALEncode
}

func NewWAL(path string, options WalOptions) (WAL, error) {
	return nil, nil
}

func (b *baseWAL) Write(msg WALMessage) {

}

func (b *baseWAL) WriteSync(msg WALMessage) {

}
func (b *baseWAL) Start() {

}
func (b *baseWAL) Stop() {

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
	return nil
}

type WALDecode struct {
	io.Reader
}

type WalOptions struct {
	bufferLimitSize int
	fileLimitSize   uint32
}

type WalFile struct {
	ID       uint64
	Path     string
	writeBuf *bufio.Writer

	metaDB  *leveldb.DB
	options WalOptions
	file    *os.File

	fileSize  uint32
	endHeight uint32
	endFileID uint32

	mux    sync.Mutex
	ticker *time.Ticker
	quit   chan struct{}
}

func NewWalFile() (*WalFile, error) {
	return nil, nil
}

func (w *WalFile) fileName(id uint64) string {
	return fmt.Sprintf("%s/wal.%d", w.Path, id)
}

func (w *WalFile) Write(p []byte) (n int, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	return 0, nil
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
	go w.writeLoop()
}

func (w *WalFile) Stop() {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.writeBuf.Flush()
	w.file.Sync()

	close(w.quit)
}
