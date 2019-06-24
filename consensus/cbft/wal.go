package cbft

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"path/filepath"
	"time"
)

const (
	// WAL working directory
	walDir = "wal"

	// WAL Database name
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

type Wal interface {
	Write(info *MsgInfo) error
	WriteSync(info *MsgInfo) error
	UpdateViewChange(info *ViewChangeMessage) error
	Load(add func(info *MsgInfo)) error
	Close()
}

type emptyWal struct {
}

func (w *emptyWal) Write(info *MsgInfo) error {
	return nil
}

func (w *emptyWal) WriteSync(info *MsgInfo) error {
	return nil
}

func (w *emptyWal) UpdateViewChange(info *ViewChangeMessage) error {
	return nil
}
func (w *emptyWal) Load(add func(info *MsgInfo)) error {
	return nil
}
func (w *emptyWal) Close() {
}

type baseWal struct {
	path    string // WAL working directory
	metaDB  IWALDatabase
	journal *journal
}

func NewWal(ctx *node.ServiceContext, specifiedPath string) (Wal, error) {
	var (
		originPath  = specifiedPath
		metaDB  IWALDatabase
		walPath string
		journal *journal
		err     error
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

// insert adds the specified MsgInfo to the local disk journal.
func (wal *baseWal) Write(info *MsgInfo) error {
	return wal.journal.Insert(&JournalMessage{
		Timestamp: uint64(time.Now().UnixNano()),
		Data:      info,
	}, false)
}

func (wal *baseWal) WriteSync(info *MsgInfo) error {
	return wal.journal.Insert(&JournalMessage{
		Timestamp: uint64(time.Now().UnixNano()),
		Data:      info,
	}, true)
}

func (wal *baseWal) UpdateViewChange(info *ViewChangeMessage) error {
	return wal.updateViewChangeMeta(info)
}

func (wal *baseWal) Load(add func(info *MsgInfo)) error {
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
