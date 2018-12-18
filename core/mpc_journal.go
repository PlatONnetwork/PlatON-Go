package core

import (
	"io"
	"os"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// devNull is a WriteCloser that just discards anything written into it. Its
// goal is to allow the transaction journal to write into a fake journal when
// loading transactions on startup without printing warnings due to no file
// being read for write.
type devMPCNull struct{}

func (*devMPCNull) Write(p []byte) (n int, err error) { return len(p), nil }
func (*devMPCNull) Close() error                      { return nil }

// mpcJournal is a rotating log of transactions with the aim of storing locally
// created transactions to allow non-executed ones to survive node restarts.
type mpcJournal struct {
	path   string         // Filesystem path to store the mpc transactions at
	writer io.WriteCloser // Output stream to write new mpc transactions into
}

// newTxJournal creates a new mpc transaction journal to
func newMPCJournal(path string) *mpcJournal {
	return &mpcJournal{
		path: path,
	}
}

// load parses a mpc transaction journal dump from disk, loading its contents into
// the specified pool.
func (journal *mpcJournal) load(add func([]*types.TransactionWrap) []error) error {

	// skip the parsing if the journal file doesn't exists at all
	if _, err := os.Stat(journal.path); os.IsNotExist(err) {
		return nil
	}

	// Open the journal for loading any past mpc transactions
	input, err := os.Open(journal.path)
	if err != nil {
		return err
	}
	defer func() { input.Close() }()

	// Temporaily discard any journal additions(dont't double add on load)
	journal.writer = new(devMPCNull)
	defer func() { journal.writer = nil }()

	// Inject all transactions from the journal into the pool
	stream := rlp.NewStream(input, 0)
	total, dropped := 0, 0

	// Create a method to load a limited batch of mpc transactions and bump
	// the appropriate progress counters. Then use this method to load all the
	// journaled mpc transactions in small-ish batches.
	loadBatch := func(txs types.TransactionWraps) {
		for _, err := range add(txs) {
			if err != nil {
				log.Debug("Failed to add journaled mpc transaction", "err", err)
				dropped++
			}
		}
	}

	var (
		failure error
		batch   types.TransactionWraps
	)
	for {
		// Parse the next mpc transaction and terminate on error
		tx := new(types.TransactionWrap)
		if err = stream.Decode(tx); err != nil {
			if err != io.EOF {
				failure = err
			}
			if batch.Len() > 0 {
				loadBatch(batch)
			}
			break
		}

		// New mpc transaction parsed, queue up for later, import if threshold is reached
		total++

		if batch = append(batch, tx); batch.Len() > 1024 {
			loadBatch(batch)
			batch = batch[:0]
		}
	}
	log.Info("Loaded local mpc transaction journal", "mpc transactions", total, "dropped", dropped)

	return failure
}

// insert adds the specified mpc transaction to the local disk journal.
func (journal *mpcJournal) insert(tx *types.TransactionWrap) error {
	if journal.writer == nil {
		return errNoActiveJournal
	}
	if err := rlp.Encode(journal.writer, tx); err != nil {
		return err
	}
	return nil
}

// rotate regenerates the mpc transaction journal based on the current contents of
// the mpc transaction pool.
func (journal *mpcJournal) rotate(all types.TransactionWraps) error {
	// Close the current journal (if any is open)
	if journal.writer != nil {
		if err := journal.writer.Close(); err != nil {
			return err
		}
		journal.writer = nil
	}

	// Generate a new journal with the contents of the current mpc pool
	replacement, err := os.OpenFile(journal.path + ".new", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	journaled := 0
	for _, tx := range all {
		//fmt.Println("write to rlp file, hash : ", tx.Hash().Hex())
		if err = rlp.Encode(replacement, tx); err != nil {
			replacement.Close()
			return err
		}
		journaled++
	}
	replacement.Close()

	// Replace the live journal with the newly generated one
	if err = os.Rename(journal.path + ".new", journal.path); err != nil {
		return err
	}
	sink, err := os.OpenFile(journal.path, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	journal.writer = sink
	log.Info("Regenerated local mpc transaction journal", "mpc transactions", journaled)
	return nil
}

// close flushes the mpc transaction journal contents to disk and closes the file.
func (journal *mpcJournal) close() error {
	var err error
	if journal.writer != nil {
		err = journal.writer.Close()
		journal.writer = nil
	}
	return err
}
