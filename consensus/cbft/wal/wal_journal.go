// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package wal

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
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

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
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

var crc32c = crc32.MakeTable(crc32.Castagnoli) // The crc verifier

var (
	errNoActiveJournal = errors.New("no active journal")
	errOpenNewJournal  = errors.New("failed to open new journal file")
	errWriteJournal    = errors.New("failed to write journal")
	errLoadJournal     = errors.New("failed to load journal")
)

// Message is a combination of consensus msg.
type Message struct {
	Timestamp uint64
	Data      interface{}
}

// sortFile represents the name and index of the journal file.
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

type journal struct {
	path                 string         // Filesystem path to store the msgInfo at
	writer               *WriterWrapper // Output stream to write new msgInfo into
	fileID               uint32
	mu                   sync.Mutex
	exitCh               chan struct{}
	mockJournalLimitSize uint64
}

// listJournalFiles sort existing files based on journal index.
// it retrieves an ascending collection
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

// newJournal creates journal object
func newJournal(path string) (*journal, error) {
	journal := &journal{
		path:   path,
		exitCh: make(chan struct{}),
		fileID: 1,
	}
	if files := listJournalFiles(path); files != nil && files.Len() > 0 {
		journal.fileID = files[len(files)-1].num
	}
	// open the corresponding journal file
	newFileID, newWriter, err := journal.newJournalFile(journal.fileID)
	if err != nil {
		return nil, err
	}
	// update field fileID and writer
	journal.fileID = newFileID
	journal.writer = newWriter

	go journal.mainLoop(syncLoopDuration)

	return journal, nil
}

func (journal *journal) mainLoop(syncLoopDuration time.Duration) {
	ticker := time.NewTicker(syncLoopDuration)
	<-ticker.C // discard the initial tick

	for {
		select {
		case <-ticker.C:
			if journal.writer != nil {
				log.Trace("Rotate timer trigger")
				journal.mu.Lock()
				if err := journal.rotate(journalLimitSize); err != nil {
					log.Error("Failed to rotate cbft journal", "err", err)
				}
				journal.mu.Unlock()
			}

		case <-journal.exitCh:
			return
		}
	}
}

// CurrentJournal retrieves the current fileID and fileSeq of the cbft journal.
func (journal *journal) CurrentJournal() (uint32, uint64, error) {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	// Forced to flush
	journal.writer.Flush()
	fileSeq, err := journal.currentFileSize()
	if err != nil {
		return 0, 0, err
	}

	log.Trace("CurrentJournal", "fileID", journal.fileID, "fileSeq", fileSeq)
	return journal.fileID, fileSeq, nil
}

// Insert adds the specified message to the local disk journal.
func (journal *journal) Insert(msg *Message, sync bool) error {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	if journal.writer == nil {
		return errNoActiveJournal
	}

	buf, err := encodeJournal(msg)
	if err != nil {
		return err
	}
	//
	if err := journal.rotate(journalLimitSize); err != nil {
		log.Error("Failed to rotate cbft journal", "err", err)
		return err
	}

	n := 0
	if n, err = journal.writer.Write(buf); err != nil || n <= 0 {
		log.Error("Write data error", "err", err)
		return errWriteJournal
	}
	if sync {
		// Forced to flush
		if err = journal.writer.Flush(); err != nil {
			log.Error("Flush data error", "err", err)
			return err
		}
	}

	log.Trace("Successful to insert journal message", "n", n)
	return nil
}

// encodeJournal tries to encode journal message with rlp.
func encodeJournal(msg *Message) ([]byte, error) {
	data, err := rlp.EncodeToBytes(msg)
	if err != nil {
		log.Error("Failed to encode journal message", "err", err)
		return nil, err
	}

	crc := crc32.Checksum(data, crc32c)
	length := uint32(len(data))
	totalLength := 10 + int(length)

	pack := make([]byte, totalLength)
	binary.BigEndian.PutUint32(pack[0:4], crc)                                         // 4 byte
	binary.BigEndian.PutUint32(pack[4:8], length)                                      // 4 byte
	binary.BigEndian.PutUint16(pack[8:10], uint16(protocols.WalMessageType(msg.Data))) // 2 byte

	copy(pack[10:], data)
	return pack, nil
}

// Close flushes the journal contents to disk and closes the file.
func (journal *journal) Close() {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	if journal.writer != nil {
		log.Debug("Close journal, flush data")
		journal.writer.FlushAndClose()
		journal.writer = nil
	}
	close(journal.exitCh)
}

// rotate tries to create a new journal file when the current journal file exceed the size limit.
func (journal *journal) rotate(journalLimitSize uint64) error {
	//journal.mu.Lock()
	//defer journal.mu.Unlock()

	if journal.mockJournalLimitSize > 0 {
		journalLimitSize = journal.mockJournalLimitSize
	}
	if journal.checkFileSize(journalLimitSize) {
		journalWriter := journal.writer
		if journalWriter == nil {
			return errNoActiveJournal
		}

		// Forced to flush
		journalWriter.FlushAndClose()
		journal.writer = nil

		// open another new journal file
		newFileID, newWriter, err := journal.newJournalFile(journal.fileID + 1)
		if err != nil {
			log.Error("Failed to create journal file", "fileID", journal.fileID+1, "error", err)
			return err
		}
		// update field fileID and writer
		journal.fileID = newFileID
		journal.writer = newWriter

		log.Debug("Successful to rotate journal file", "newFileID", newFileID)
	}
	return nil
}

// checkFileSize determine if the current journal file exceed the size limit.
// the limit size is configurable.
func (journal *journal) checkFileSize(journalLimitSize uint64) bool {
	fileSize, err := journal.currentFileSize()
	return err == nil && fileSize >= journalLimitSize
}

// currentFileSize retrieves the size of current journal file.
func (journal *journal) currentFileSize() (uint64, error) {
	var (
		fileInfo os.FileInfo
		err      error
	)
	if fileInfo, err = journal.writer.file.Stat(); err != nil {
		log.Error("Get the current journal file size error", "err", err)
		return 0, err
	}
	return uint64(fileInfo.Size()), nil
}

// newJournalFile create a new journal file.
// Subsequent messages will be written to the new file.
func (journal *journal) newJournalFile(fileID uint32) (uint32, *WriterWrapper, error) {
	newJournalFilePath := filepath.Join(journal.path, fmt.Sprintf("wal.%d", fileID))
	file, err := os.OpenFile(newJournalFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Error("Failed to open new journal file", "fileID", fileID, "filePath", newJournalFilePath, "err", err)
		return 0, nil, errOpenNewJournal
	}

	return fileID, NewWriterWrapper(file, writeBufferLimitSize), nil
}

// ExpireJournalFile tries to remove the expired journal file
// when a new confirm viewChange is written, the previous message will expire.
func (journal *journal) ExpireJournalFile(fileID uint32) error {
	if files := listJournalFiles(journal.path); files != nil && files.Len() > 0 {
		for _, file := range files {
			if file.num != journal.fileID && file.num < fileID {
				os.Remove(filepath.Join(journal.path, fmt.Sprintf("wal.%d", file.num)))
			}
		}
	}
	return nil
}

// LoadJournal tries to load consensus message from journal file
// search starting from the specified file and seq,will verify each message at the same time
// recovery is the callback function
func (journal *journal) LoadJournal(fromFileID uint32, fromSeq uint64, recovery recoveryConsensusMsgFn) (err error) {
	journal.mu.Lock()
	defer journal.mu.Unlock()

	if files := listJournalFiles(journal.path); files != nil && files.Len() > 0 {
		log.Debug("Begin to load journal", "fromFileID", fromFileID, "fromSeq", fromSeq)
		for _, file := range files {
			if file.num == fromFileID {
				err = journal.loadJournal(file.num, fromSeq, recovery)
			} else if file.num > fromFileID {
				err = journal.loadJournal(file.num, 0, recovery)
			}
			if err != nil {
				return err
			}
		}
	} else {
		log.Error("Failed to load journal", "fromFileID", fromFileID, "fromSeq", fromSeq)
		return errLoadJournal
	}
	return nil
}

// loadJournal is a concrete implementation to load consensus message from journal file
// Each message is loaded into the caller as a callback function
func (journal *journal) loadJournal(fileID uint32, seq uint64, recovery recoveryConsensusMsgFn) error {
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
		index, _ := bufReader.Peek(10)
		crc := binary.BigEndian.Uint32(index[0:4])      // 4 byte
		length := binary.BigEndian.Uint32(index[4:8])   // 4 byte
		msgType := binary.BigEndian.Uint16(index[8:10]) // 2 byte

		pack := make([]byte, length+10)
		var (
			totalNum int
			readNum  int
		)
		for totalNum, err = 0, error(nil); err == nil && uint32(totalNum) < length+10; {
			readNum, err = bufReader.Read(pack[totalNum:])
			totalNum = totalNum + readNum
		}

		if 0 == readNum {
			log.Debug("Load journal complete", "fileID", fileID, "fileSeq", seq)
			break
		}

		// check crc
		_crc := crc32.Checksum(pack[10:], crc32c)
		if crc != _crc {
			log.Error("Crc is invalid", "crc", crc, "_crc", _crc, "msgType", msgType)
			return errLoadJournal
		}

		// decode journal message
		if msgInfo, err := WALDecode(pack[10:], msgType); err == nil {
			if err = recovery(msgInfo); err != nil {
				return err
			}
		} else {
			log.Error("Failed to decode journal msg", "err", err)
			return errLoadJournal
		}
	}
	return nil
}
