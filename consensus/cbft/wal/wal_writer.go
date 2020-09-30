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
