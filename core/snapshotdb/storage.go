package snapshotdb

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func sortFds(fds []fileDesc) {
	sort.Sort(fileDescs(fds))
}

type fileDescs []fileDesc

func (f fileDescs) Len() int {
	return len(f)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (f fileDescs) Less(i, j int) bool {
	return f[i].Num < f[j].Num
}

func (f fileDescs) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type fileDesc struct {
	Type      fileType
	Num       uint64
	BlockHash common.Hash
}

func (fd fileDesc) String() string {
	switch fd.Type {
	case TypeJournal:
		return fmt.Sprintf("%010d-%s.log", fd.Num, fd.BlockHash.String())
	case TypeCurrent:
		return fmt.Sprintf("current")
	default:
		return ""
	}
}

// FileDescOk returns true if fd is a valid 'file descriptor'.
func FileDescOk(fd fileDesc) bool {
	switch fd.Type {
	case TypeJournal:
	case TypeCurrent:
	default:
		return false
	}
	return fd.Num >= 0
}

// Common error.
var (
	ErrInvalidFile = errors.New("snapshotdb/storage: invalid file for argument")
	ErrLocked      = errors.New("snapshotdb/storage: already locked")
	ErrClosed      = errors.New("snapshotdb/storage: closed")
	//errFileOpen    = errors.New("snapshotdb/storage: file still open")
	errReadOnly = errors.New("snapshotdb/storage: storage is read-only")
)

type fileLock interface {
	release() error
}

const logSizeThreshold = 1024 * 1024 // 1 MiB

// storage is the storage. A storage instance must be safe for concurrent use.
type storage interface {
	// Create creates file with the given 'file descriptor', truncate if already
	// exist and opens write-only.
	// Returns ErrClosed if the underlying storage is closed.
	Create(fd fileDesc) (io.WriteCloser, error)

	// Close closes the storage.
	// It is valid to call Close multiple times. Other methods should not be
	// called after the storage has been closed.
	Close() error

	// Open opens file with the given 'file descriptor' read-only.
	// Returns os.ErrNotExist error if the file does not exist.
	// Returns ErrClosed if the underlying storage is closed.
	Open(fd fileDesc) (Reader, error)

	// Rename renames file from oldfd to newfd.
	// Returns ErrClosed if the underlying storage is closed.
	Rename(oldfd, newfd fileDesc) error

	// Remove removes file with the given 'file descriptor'.
	// Returns ErrClosed if the underlying storage is closed.
	Remove(fd fileDesc) error

	// Append append file with the given 'file descriptor', opens write-only.
	// Returns ErrClosed if the underlying storage is closed.
	Append(fd fileDesc) (io.WriteCloser, error)

	// List returns file descriptors that match the given file types.
	// The file types may be OR'ed together.
	List(ft fileType) ([]fileDesc, error)

	// Path return path of the storage
	Path() string
}
