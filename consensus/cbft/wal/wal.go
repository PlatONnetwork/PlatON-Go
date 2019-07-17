// Copyright 2019 The go-ethereum Authors
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

// Package wal implements the similar write-ahead logging for cbft consensus.
package wal

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	// Wal working directory
	walDir = "wal"

	// Wal database name
	metaDBName = "wal_meta"
)

var (
	viewChangeKey = []byte("view-change")
)

var (
	errCreateWalDir         = errors.New("Failed to create wal directory")
	errUpdateViewChangeMeta = errors.New("Failed to update viewChange meta")
	errGetViewChangeMeta    = errors.New("Failed to get viewChange meta")
)

type ViewChangeMessage struct {
	Hash   common.Hash
	Number uint64
}

type ViewChangeMeta struct {
	Number uint64
	Hash   common.Hash
	FileID uint32
	Seq    uint64
}

type State struct {
	Block      *protocols.PrepareBlock
	QuorumCert *types.QuorumCert
}

type ChainState struct {
	Commit *State
	Lock   *State
	QC     []*State
}

type WalMsg struct {
	Msg interface{}
}

type Wal interface {
	UpdateChainState(chainState *ChainState) error
	LoadChainState(recovery func(chainState *ChainState)) error
	Write(msg *WalMsg) error
	WriteSync(msg *WalMsg) error
	UpdateViewChange(info *ViewChangeMessage) error
	Load(add func(msg *WalMsg)) error
	Close()
}

type emptyWal struct {
}

func (w *emptyWal) UpdateChainState(chainState *ChainState) error {
	return nil
}

func (w *emptyWal) LoadChainState(recovery func(chainState *ChainState)) error {
	return nil
}

func (w *emptyWal) Write(msg *WalMsg) error {
	return nil
}

func (w *emptyWal) WriteSync(msg *WalMsg) error {
	return nil
}

func (w *emptyWal) UpdateViewChange(info *ViewChangeMessage) error {
	return nil
}

func (w *emptyWal) Load(add func(msg *WalMsg)) error {
	return nil
}

func (w *emptyWal) Close() {
}

type baseWal struct {
	path    string // Wal working directory
	metaDB  IWALDatabase
	journal *journal
}

func NewWal(ctx *node.ServiceContext, specifiedPath string) (Wal, error) {
	if ctx == nil && len(specifiedPath) == 0 {
		return &emptyWal{}, nil
	}
	var (
		originPath = specifiedPath
		metaDB     IWALDatabase
		walPath    string
		journal    *journal
		err        error
	)
	if originPath == "" {
		originPath = ctx.ResolvePath(walDir)
	}

	// Make sure the wal directory exists,If not exist create it.
	ensureWalDir := func(path string) (string, error) {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.MkdirAll(path, 0700)
			if err != nil {
				return "", errCreateWalDir
			}
		}
		return path, nil
	}
	// Open or create WAL Database
	createMetaDB := func(path, name string) (IWALDatabase, error) {
		db, err := createWalDB(filepath.Join(path, name))
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	if walPath, err = ensureWalDir(originPath); err != nil {
		log.Error("Failed to create wal directory", "err", err)
		return nil, err
	}
	if journal, err = NewJournal(walPath); err != nil {
		return nil, err
	}
	if metaDB, err = createMetaDB(walPath, metaDBName); err != nil {
		log.Error("Failed to create wal database", "err", err)
		return nil, err
	}

	wal := &baseWal{
		path:    walPath,
		metaDB:  metaDB,
		journal: journal,
	}

	return wal, nil
}

func (wal *baseWal) UpdateChainState(chainState *ChainState) error {
	return nil
}

func (wal *baseWal) LoadChainState(recovery func(chainState *ChainState)) error {
	return nil
}

// insert adds the specified MsgInfo to the local disk journal.
func (wal *baseWal) Write(msg *WalMsg) error {
	return wal.journal.Insert(&JournalMessage{
		Timestamp: uint64(time.Now().UnixNano()),
		Data:      msg,
	}, false)
}

func (wal *baseWal) WriteSync(msg *WalMsg) error {
	return wal.journal.Insert(&JournalMessage{
		Timestamp: uint64(time.Now().UnixNano()),
		Data:      msg,
	}, true)
}

func (wal *baseWal) UpdateViewChange(info *ViewChangeMessage) error {
	return wal.updateViewChangeMeta(info)
}

func (wal *baseWal) Load(add func(msg *WalMsg)) error {
	// open wal database
	data, err := wal.metaDB.Get(viewChangeKey)
	if err != nil {
		log.Warn("Failed to get viewChange meta from db,may be the first time to run platon")
		return nil
	}
	var v ViewChangeMeta
	err = rlp.DecodeBytes(data, &v)
	if err != nil {
		log.Error("Failed to decode viewChange meta")
		return errGetViewChangeMeta
	}

	return wal.journal.LoadJournal(v.FileID, v.Seq, add)
}

// Update the ViewChange Meta Data to the database.
func (wal *baseWal) updateViewChangeMeta(vc *ViewChangeMessage) error {
	fileID, seq, err := wal.journal.CurrentJournal()
	if err != nil {
		log.Error("Failed to update viewChange meta", "number", vc.Number, "hash", vc.Hash, "err", err)
		return errUpdateViewChangeMeta
	}

	viewChangeMeta := &ViewChangeMeta{
		Number: vc.Number,
		Hash:   vc.Hash,
		FileID: fileID,
		Seq:    seq,
	}
	data, err := rlp.EncodeToBytes(viewChangeMeta)
	if err != nil {
		return err
	}
	// Write the ViewChangeMeta to the WAL database
	err = wal.metaDB.Put(viewChangeKey, data, &opt.WriteOptions{Sync: true})
	if err != nil {
		return err
	}
	log.Debug("success to update viewChange meta", "number", vc.Number, "hash", vc.Hash, "fileID", fileID, "seq", seq)
	// Delete previous journal logs
	go wal.journal.ExpireJournalFile(fileID)
	return nil
}

func (wal *baseWal) Close() {
	wal.metaDB.Close()
	wal.journal.Close()
}
