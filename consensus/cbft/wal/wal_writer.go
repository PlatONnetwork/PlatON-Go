package wal

import (
	"bufio"
	"os"
)

type WriterWrapper struct {
	file   *os.File
	writer *bufio.Writer
}

func NewWriterWrapper(file *os.File, bufferLimitSize int) *WriterWrapper {
	bufWriter := bufio.NewWriterSize(file, bufferLimitSize)
	w := &WriterWrapper{
		file:   file,
		writer: bufWriter,
	}
	return w
}

func (w *WriterWrapper) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *WriterWrapper) Flush() (err error) {
	err = w.writer.Flush()
	if err != nil {
		return err
	}
	return w.file.Sync()
}

func (w *WriterWrapper) Close() (err error) {
	return w.file.Close()
}

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
