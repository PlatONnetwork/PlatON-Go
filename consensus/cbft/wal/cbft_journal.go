package cbft

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// The limit size of a single journal file
	journalLimitSize = 100 * 1024 * 1024

	// A new Writer whose buffer has at least the specified size
	writeBufferLimitSize = 16 * 1024

	// A new Reader whose buffer has at least the specified size
	readBufferLimitSize = 16 * 1024

	// The setting of rotate timer ticker
	syncLoopDuration = 5 * time.Second
)

var crc32c = crc32.MakeTable(crc32.Castagnoli)

// errNoActiveJournal is returned if a transaction is attempted to be inserted
// into the journal, but no such file is currently open.
var (
	errNoActiveJournal = errors.New("no active journal")
	errOpenNewJournal  = errors.New("Failed to open new journal file")
	errLoadJournal  = errors.New("Failed to load journal")
)

//
type JournalMessage struct {
	Timestamp uint64
	Msg       interface{}
}

type sortFile struct {
	name string
	num  uint32
}

type sortFiles []sortFile

func (s sortFiles) Len() int {
	return len(s)
}

func (s sortFiles) Less(i, j int) bool {
	return s[i].num < s[j].num
}

func (s sortFiles) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// txJournal is a rotating log of transactions with the aim of storing locally
// created transactions to allow non-executed ones to survive node restarts.
type cbftJournal struct {
	path   string         // Filesystem path to store the transactions at
	writer *WriterWrapper // Output stream to write new transactions into
	fileID uint32
	mu     sync.Mutex
	exitCh chan struct{}
}

// 扫描wal目录,遍历所有journal文件并按文件下标排序
func listJournalFiles(path string) sortFiles {
	files, err := ioutil.ReadDir(path)

	if err == nil && len(files) > 0 {
		var s []string
		for _, f := range files {
			s = append(s, f.Name())
		}
		log.Trace("The list of wal directory", "directory", path, "files", strings.Join(s, ","))
		reg := regexp.MustCompile("^wal.([1-9][0-9]*)$")
		regNum := regexp.MustCompile("([1-9][0-9]*)$")
		fs := make(sortFiles, 0)

		for _, f := range s {
			if reg.MatchString(f) {
				n, _ := strconv.Atoi(regNum.FindString(f))
				fs = append(fs, sortFile{
					name: f,
					num:  uint32(n),
				})
			}
		}
		sort.Sort(fs)
		return fs
	}
	return nil
}

// newTxJournal creates a new transaction journal to
func NewCbftJournal(path string) (*cbftJournal, error) {
	journal := &cbftJournal{
		path:   path,
		exitCh: make(chan struct{}),
		fileID: 1,
	}
	if files := listJournalFiles(path); files != nil && files.Len() > 0 {
		journal.fileID = files[len(files)-1].num
	}
	// 2 重新开辟一个journal文件
	newFileID, newWriter, err := journal.newJournalFile(journal.fileID)
	if err != nil {
		return nil, err
	}
	// 3 更新fileID，writer
	journal.fileID = newFileID
	journal.writer = newWriter

	go journal.mainLoop(syncLoopDuration)

	return journal, nil
}

func (journal *cbftJournal) mainLoop(syncLoopDuration time.Duration) {
	ticker := time.NewTicker(syncLoopDuration)
	<-ticker.C // discard the initial tick

	for {
		select {
		case <-ticker.C:
			if journal.writer != nil {
				log.Trace("Rotate timer trigger")
				if err := journal.rotate(journalLimitSize); err != nil {
					log.Error("Failed to rotate cbft journal", "err", err)
				}
			}

		case <-journal.exitCh:
			return
		}
	}
}

// currentJournal retrieves the current fileID and fileSeq of the cbft journal.
func (journal *cbftJournal) CurrentJournal() (uint32, uint64, error) {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	// 刷一次磁盘，更新journal信息
	journal.writer.Flush()
	fileSeq, err := journal.currentFileSize()
	if err != nil {
		return 0, 0, err
	}

	log.Trace("currentJournal", "fileID", journal.fileID, "fileSeq", fileSeq)
	return journal.fileID, fileSeq, nil
}

// insert adds the specified consensus message to the local disk journal.
func (journal *cbftJournal) insert(msg *JournalMessage) error {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	if journal.writer == nil {
		return errNoActiveJournal
	}

	buf, err := encodeJournal(msg)
	if err != nil {
		return err
	}
	n, err := journal.writer.Write(buf)
	if err == nil && n > 0 {
		log.Trace("Successful to insert journal message", "n", n)
		return nil
	}
	return err
}

func encodeJournal(msg *JournalMessage) ([]byte, error) {
	data, err := rlp.EncodeToBytes(msg)
	if err != nil {
		log.Error("Failed to encode journal message", "err", err)
		return nil, err
	}

	crc := crc32.Checksum(data, crc32c)
	length := uint32(len(data))
	totalLength := 8 + int(length)

	buf := make([]byte, totalLength)
	binary.BigEndian.PutUint32(buf[0:4], crc)
	binary.BigEndian.PutUint32(buf[4:8], length)
	copy(buf[8:], data)
	return buf, nil
}

// close flushes the transaction journal contents to disk and closes the file.
func (journal *cbftJournal) close() {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	if journal.writer != nil {
		// 将缓存数据输入磁盘，关闭writer
		journal.writer.FlushAndClose()
		journal.writer = nil
	}
	close(journal.exitCh)
}

func (journal *cbftJournal) rotate(journalLimitSize uint64) error {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	if journal.checkFileSize(journalLimitSize) {
		journalWriter := journal.writer
		if journalWriter == nil {
			return errNoActiveJournal
		}

		// 1 将缓存数据输入磁盘
		journalWriter.FlushAndClose()
		journal.writer = nil

		// 2 重新开辟一个journal文件
		newFileID, newWriter, err := journal.newJournalFile(journal.fileID + 1)
		if err != nil {
			return err
		}
		// 3 更新fileID，writer
		journal.fileID = newFileID
		journal.writer = newWriter

		log.Debug("Successful to rotate journal file", "newFileID", newFileID)
	}
	return nil
}

func (journal *cbftJournal) checkFileSize(journalLimitSize uint64) bool {
	fileSize, err := journal.currentFileSize()
	return err == nil && fileSize >= journalLimitSize
}

func (journal *cbftJournal) currentFileSize() (uint64, error) {
	currentFile := filepath.Join(journal.path, fmt.Sprintf("wal.%d", journal.fileID))
	if fileInfo, err := os.Stat(currentFile); err != nil {
		log.Error("Get the current journal file size error", "err", err)
		return 0, err
	} else {
		return uint64(fileInfo.Size()), nil
	}
}

func (journal *cbftJournal) newJournalFile(fileID uint32) (uint32, *WriterWrapper, error) {
	newJournalFilePath := filepath.Join(journal.path, fmt.Sprintf("wal.%d", fileID))
	file, err := os.OpenFile(newJournalFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Error("Failed to open new journal file", "fileID", fileID, "filePath", newJournalFilePath, "err", err)
		return 0, nil, errOpenNewJournal
	}

	return fileID, NewWriterWrapper(file, writeBufferLimitSize), nil
}

func (journal *cbftJournal) ExpireJournalFile(fileID uint32) error {
	if files := listJournalFiles(journal.path); files != nil && files.Len() > 0 {
		for _, file := range files {
			if file.num != journal.fileID && file.num <  fileID {
				os.Remove(filepath.Join(journal.path, fmt.Sprintf("wal.%d", file.num)))
			}
		}
	}
	return nil
}

func (journal *cbftJournal) LoadJournal(fromFileID uint32, fromSeq uint64, add func(j *JournalMessage)) error {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	if files := listJournalFiles(journal.path); files != nil && files.Len() > 0 {
		for i:= files.Len(); i >= 0; i-- {
			file := files[i]

			if file.num == fromFileID {
				journal.loadJournal(file.num, fromSeq, add)
			} else if file.num > fromFileID {
				journal.loadJournal(file.num, 0, add)
			}
		}
	} else {
		log.Error("Failed to load journal", "fromFileID", fromFileID, "fromSeq", fromSeq)
		return errLoadJournal
	}
	return nil
}

func (journal *cbftJournal) loadJournal(fileID uint32, seq uint64, add func(j *JournalMessage) error) error {
	file, err := os.Open(filepath.Join(journal.path, fmt.Sprintf("wal.%d", fileID)))
	if err != nil {
		return err
	}
	defer file.Close()

	bufReader := bufio.NewReaderSize(file, readBufferLimitSize)
	if seq > 0 {
		bufReader.Discard(int(seq))
	}

	for {
		index, _ := bufReader.Peek(8)
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

		if 0 == readNum {
			break
		}

		_crc := crc32.Checksum(index[8:], crc32c)
		if crc != _crc {
			log.Error("crc is invalid", "crc", crc, "_crc", _crc)
			return errLoadJournal
		}
		var j JournalMessage
		if err := rlp.DecodeBytes(pack[8:], &j); err == nil {
			fmt.Println("Timestamp", j.Timestamp)
			// TODO
		} else {
			log.Error("Failed to decode journal msg", "err", err)
			return errLoadJournal
		}
	}
	return nil
}
