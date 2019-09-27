package ethdb

// Iterator iterates over a database's key/value pairs in ascending key order.
//
// When it encounters an error any seek will return false and will yield no key/
// value pairs. The error can be queried by calling the Error method. Calling
// Release is still necessary.
//
// An iterator must be released after use, but it is not necessary to read an
// iterator until exhaustion. An iterator is not safe for concurrent use, but it
// is safe to use multiple iterators concurrently.
type Iterator interface {
	// Next moves the iterator to the next key/value pair. It returns whether the
	// iterator is exhausted.
	Next() bool

	// Error returns any accumulated error. Exhausting all the key/value pairs
	// is not considered to be an error.
	Error() error

	// Key returns the key of the current key/value pair, or nil if done. The caller
	// should not modify the contents of the returned slice, and its contents may
	// change on the next call to Next.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done. The
	// caller should not modify the contents of the returned slice, and its contents
	// may change on the next call to Next.
	Value() []byte

	// Release releases associated resources. Release should always succeed and can
	// be called multiple times without causing error.
	Release()
}
