package cbft

import (
	"errors"
	"fmt"
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
	//
	walDir = "wal"

	//
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

type WALMessage interface{}

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

type Wal struct {
	path        string // WAL working directory
	metaDB      IWALDatabase
	cbftJournal *cbftJournal
}

func NewWal(ctx *node.ServiceContext) (*Wal, error) {
	var (
		//TODO
		//originPath  = ctx.ResolvePath(walDir)
		originPath  = "D://data/platon/wal"
		metaDB      IWALDatabase
		walPath     string
		cbftJournal *cbftJournal
		err         error
	)

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
	if cbftJournal, err = NewCbftJournal(walPath); err != nil {
		return nil, err
	}
	if metaDB, err = createMetaDB(walPath, metaDBName); err != nil {
		log.Error("Failed to create wal database", "err", err)
		return nil, err
	}

	wal := &Wal{
		path:        walPath,
		metaDB:      metaDB,
		cbftJournal: cbftJournal,
	}

	return wal, nil
}

// insert adds the specified transaction to the local disk journal.
func (wal *Wal) Write(msg WALMessage) error {
	switch m := msg.(type) {
	case *ViewChangeMessage:
		// CommittedWALMessage is special msg, it indicate cbft had commit block.
		// Update meta info when commit message was wrote disk
		return wal.updateViewChangeMeta(m)

	default:
		return wal.cbftJournal.insert(&JournalMessage{
			Timestamp: uint64(time.Now().UnixNano()),
			Msg:       m,
		})
	}
}

func (wal *Wal) Load(add func(msg WALMessage)) error {
	// 打开leveldb
	data, err := wal.metaDB.Get(viewChangeKey)
	if err != nil {
		log.Error("Failed to get viewChange meta from db")
		return errGetViewChangeMeta
	}
	var v ViewChangeMeta
	err = rlp.DecodeBytes(data, &v)
	if err != nil {
		log.Error("Failed to decode viewChange meta")
		return errGetViewChangeMeta
	}

	return wal.cbftJournal.LoadJournal(v.FileID, v.Seq, add)
}

// Update the ViewChange Meta Data to the database.
func (wal *Wal) updateViewChangeMeta(vc *ViewChangeMessage) error {
	fileID, seq, err := wal.cbftJournal.CurrentJournal()
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
	// 写入leveldb数据库
	err = wal.metaDB.Put(viewChangeKey, data, &opt.WriteOptions{Sync: true})
	if err != nil {
		return err
	}
	// 删除多余日志
	return wal.cbftJournal.ExpireJournalFile(fileID)
}

func (wal *Wal) close() {
	wal.metaDB.Close()
	wal.cbftJournal.close()
}
