package wal

import (
	"bufio"
	"os"
)

// WriterWrapper wrap the write file.
// it contains the actual file handle and memory bufio
type WriterWrapper struct {
	file   *os.File
	writer *bufio.Writer
}

// NewWriterWrapper creates a new buffer writer to write data to file.
func NewWriterWrapper(file *os.File, bufferLimitSize int) *WriterWrapper {
	bufWriter := bufio.NewWriterSize(file, bufferLimitSize)
	w := &WriterWrapper{
		file:   file,
		writer: bufWriter,
	}
	return w
}

// Write tries to record data to memory bufio.
func (w *WriterWrapper) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

// Flush tries to flush the data from memory bufio to actual file.
func (w *WriterWrapper) Flush() (err error) {
	err = w.writer.Flush()
	if err != nil {
		return err
	}
	return w.file.Sync()
}

// Close tries to actual file handle.
func (w *WriterWrapper) Close() (err error) {
	return w.file.Close()
}

// FlushAndClose successive invoke function Flush and Close
func (w *WriterWrapper) FlushAndClose() (err error) {
	err = w.writer.Flush()
	if err != nil {
		return err
	}
	err = w.file.Sync()
	if err != nil {
		return err
	}
	return w.Close()
}
