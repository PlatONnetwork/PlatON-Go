package snapshotdb

import "github.com/PlatONnetwork/PlatON-Go/common"

// Putter wraps the database write operation supported by both batches and regular databases.
type Putter interface {
	Put(hash common.Hash, key []byte, value []byte) error
}

// Deleter wraps the database delete operation supported by both batches and regular databases.
type Deleter interface {
	Delete(hash common.Hash, key []byte) error
}

type Writer interface {
	Putter
	Deleter
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	Putter
	Deleter
	Get(hash common.Hash, key []byte) ([]byte, error)
	Has(hash common.Hash, key []byte) (bool, error)
}
