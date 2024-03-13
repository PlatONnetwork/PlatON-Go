// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"bytes"
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

// EVMLogger is used to collect execution traces from an EVM transaction
// execution. CaptureState is called for each step of the VM with the
// current VM state.
// Note that reference types are actual VM data structures; make copies
// if you need to retain them beyond the current call.
type EVMLogger interface {
	CaptureStart(env *EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int)
	CaptureState(pc uint64, op OpCode, gas, cost uint64, scope *ScopeContext, rData []byte, depth int, err error)
	CaptureEnter(typ OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int)
	CaptureExit(output []byte, gasUsed uint64, err error)
	CaptureFault(pc uint64, op OpCode, gas, cost uint64, scope *ScopeContext, depth int, err error)
	CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error)
}

type WasmLogger struct {
	log.Logger
	root   log.Logger
	buf    *bytes.Buffer
	logger log.Logger
}

func NewWasmLogger(cfg Config, root log.Logger) *WasmLogger {
	l := &WasmLogger{
		root:   root,
		logger: root.New(),
	}

	l.buf = new(bytes.Buffer)

	level := log.LvlInfo

	if cfg.Debug {
		level = log.LvlDebug
	}
	if log.GetWasmLogLevel() >= log.LvlDebug {
		level = log.GetWasmLogLevel()
	}

	l.logger.SetHandler(log.LvlFilterHandler(level, log.StreamHandler(l.buf, log.FormatFunc(func(r *log.Record) []byte {
		return []byte(r.Msg)
	}))))

	return l
}

func (wl *WasmLogger) Flush() {
	if wl.buf.Len() != 0 {
		wl.root.Debug(wl.buf.String())
	}
	wl.buf.Reset()
}

func (wl *WasmLogger) New(ctx ...interface{}) log.Logger {
	return nil
}

// GetHandler gets the handler associated with the logger.
func (wl *WasmLogger) GetHandler() log.Handler {
	return nil
}

// SetHandler updates the logger to write records to the specified handler.
func (wl *WasmLogger) SetHandler(h log.Handler) {
}

// Log a message at the given level with context key/value pairs
func (wl *WasmLogger) Trace(msg string, ctx ...interface{}) {
	wl.logger.Trace(msg, ctx...)
}
func (wl *WasmLogger) Debug(msg string, ctx ...interface{}) {
	wl.logger.Debug(msg, ctx...)
}
func (wl *WasmLogger) Info(msg string, ctx ...interface{}) {
	wl.logger.Info(msg, ctx...)
}
func (wl *WasmLogger) Warn(msg string, ctx ...interface{}) {
	wl.logger.Warn(msg, ctx...)
}
func (wl *WasmLogger) Error(msg string, ctx ...interface{}) {
	wl.logger.Error(msg, ctx...)
}
func (wl *WasmLogger) Crit(msg string, ctx ...interface{}) {
	wl.logger.Crit(msg, ctx...)
}
