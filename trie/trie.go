// Copyright 2014 The go-ethereum Authors
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

// Package trie implements Merkle Patricia Tries.
package trie

import (
	"bytes"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// emptyState is the known hash of an empty state trie entry.
	emptyState = crypto.Keccak256Hash(nil)
	//storagePrefix = "storage-value-"
	emptyStorage = crypto.Keccak256Hash(nil)
)

// LeafCallback is a callback type invoked when a trie operation reaches a leaf
// node. It's used by state sync and commit to allow handling external references
// between account and storage tries.
type LeafCallback func(leaf []byte, parent common.Hash) error

// Trie is a Merkle Patricia Trie.
// The zero value is an empty trie with no database.
// Use New to create a trie that sits on top of a database.
//
// Trie is not safe for concurrent use.
type Trie struct {
	db           *Database
	root         node

	dag *trieDag
}

// newFlag returns the cache flag value for a newly created node.
func (t *Trie) newFlag() nodeFlag {
	dirty := true
	return nodeFlag{hash: &hashNode{}, dirty: &dirty}
}

// New creates a trie with an existing root node from db.
//
// If root is the zero hash or the sha3 hash of an empty string, the
// trie is initially empty and does not require a database. Otherwise,
// New will panic if db is nil and returns a MissingNodeError if root does
// not exist in the database. Accessing the trie loads nodes from db on demand.
func New(root common.Hash, db *Database) (*Trie, error) {
	if db == nil {
		panic("trie.New called without a database")
	}
	trie := &Trie{
		db:           db,
		dag:          newTrieDag(),
	}
	// If root is not empty, restore the node from the DB (the whole tree)
	if root != (common.Hash{}) && root != emptyRoot {
		rootnode, err := trie.resolveHash(root[:], nil)
		if err != nil {
			return nil, err
		}
		trie.root = rootnode
	}
	return trie, nil
}

// NodeIterator returns an iterator that returns nodes of the trie. Iteration starts at
// the key after the given start key.
func (t *Trie) NodeIterator(start []byte) NodeIterator {
	return newNodeIterator(t, start)
}

// Get returns the value for key stored in the trie.
// The value bytes must not be modified by the caller.
func (t *Trie) Get(key []byte) []byte {
	res, err := t.TryGet(key)
	if err != nil {
		log.Error(fmt.Sprintf("Unhandled trie error: %v", err))
	}
	return res
}

// TryGet returns the value for key stored in the trie.
// The value bytes must not be modified by the caller.
// If a node was not found in the database, a MissingNodeError is returned.
func (t *Trie) TryGet(key []byte) ([]byte, error) {
	key = keybytesToHex(key)
	value, newroot, didResolve, err := t.tryGet(t.root, key, 0)
	if err == nil && didResolve {
		t.root = newroot
	}
	return value, err
}

func (t *Trie) tryGet(origNode node, key []byte, pos int) (value []byte, newnode node, didResolve bool, err error) {
	switch n := (origNode).(type) {
	case nil:
		return nil, nil, false, nil
	case valueNode:
		return n, n, false, nil
	case *shortNode:
		if len(key)-pos < len(n.Key) || !bytes.Equal(n.Key, key[pos:pos+len(n.Key)]) {
			// key not found in trie
			return nil, n, false, nil
		}
		value, newnode, didResolve, err = t.tryGet(n.Val, key, pos+len(n.Key))
		if err == nil && didResolve {
			n = n.copy()
			n.Val = newnode
		}
		return value, n, didResolve, err
	case *fullNode:
		value, newnode, didResolve, err = t.tryGet(n.Children[key[pos]], key, pos+1)
		if err == nil && didResolve {
			n = n.copy()
			n.Children[key[pos]] = newnode
		}
		return value, n, didResolve, err
	case hashNode:
		child, err := t.resolveHash(n, key[:pos])
		if err != nil {
			return nil, n, true, err
		}
		value, newnode, _, err := t.tryGet(child, key, pos)
		return value, newnode, true, err
	default:
		panic(fmt.Sprintf("%T: invalid node: %v", origNode, origNode))
	}
}

// Update associates key with value in the trie. Subsequent calls to
// Get will return value. If value has length zero, any existing value
// is deleted from the trie and calls to Get will return nil.
//
// The value bytes must not be modified by the caller while they are
// stored in the trie.
func (t *Trie) Update(key, value []byte) {
	if err := t.TryUpdate(key, value); err != nil {
		log.Error(fmt.Sprintf("Unhandled trie error: %v", err))
		if t.dag != nil {
			t.dag.clear()
		}
	}
}

// TryUpdate associates key with value in the trie. Subsequent calls to
// Get will return value. If value has length zero, any existing value
// is deleted from the trie and calls to Get will return nil.
//
// The value bytes must not be modified by the caller while they are
// stored in the trie.
//
// If a node was not found in the database, a MissingNodeError is returned.
func (t *Trie) TryUpdate(key, value []byte) error {
	k := keybytesToHex(key)
	if len(value) != 0 {
		_, n, err := t.insert(t.root, nil, nil, k, valueNode(value))
		if err != nil {
			return err
		}
		t.root = n
	} else {
		_, n, err := t.delete(t.root, nil, k)
		if err != nil {
			return err
		}
		if t.dag != nil {
			t.dag.delVertexAndEdgeByNode(nil, t.root)
			t.dag.addVertexAndEdge(nil, nil, n)
		}
		t.root = n
	}
	return nil
}

func (t *Trie) insert(n node, fprefix, prefix, key []byte, value node) (bool, node, error) {
	if len(key) == 0 {
		if v, ok := n.(valueNode); ok {
			return !bytes.Equal(v, value.(valueNode)), value, nil
		}
		if t.dag != nil {
			//fmt.Printf("239: del vtx -> prefix: %x\n", prefix)
			t.dag.delVertexAndEdgeByNode(prefix, value)
			t.dag.addVertexAndEdge(fprefix, prefix, value)
		}
		return true, value, nil
	}
	switch n := n.(type) {
	case *shortNode:
		matchlen := prefixLen(key, n.Key)
		// If the whole key matches, keep this short node as is
		// and only update the value.
		if matchlen == len(n.Key) {
			dirty, nn, err := t.insert(n.Val, append(prefix, key[:matchlen]...), append(prefix, key[:matchlen]...), key[matchlen:], value)
			if !dirty || err != nil {
				return false, n, err
			}
			rn := &shortNode{n.Key, nn, t.newFlag()}
			if t.dag != nil {
				//fmt.Printf("257: del vtx -> prefix: %x\n", append(prefix, n.Key...))
				t.dag.delVertexAndEdge(append(prefix, n.Key...))
				t.dag.addVertexAndEdge(fprefix, prefix, rn)
			}
			return true, rn, nil
		}
		// Otherwise branch out at the index where they differ.
		branch := &fullNode{flags: t.newFlag()}
		pprefix := common.CopyBytes(prefix)
		if matchlen > 0 {
			pprefix = append(pprefix, key[:matchlen]...)
		}
		pprefix = append(pprefix, fullNodeSuffix...)
		if t.dag != nil {
			//fmt.Printf("281: del vtx -> prefix: %x\n", append(prefix, n.Key...))
			t.dag.delVertexAndEdge(append(prefix, n.Key...))
		}
		var err error
		_, branch.Children[n.Key[matchlen]], err = t.insert(nil, pprefix, append(prefix, n.Key[:matchlen+1]...), n.Key[matchlen+1:], n.Val)
		if err != nil {
			return false, nil, err
		}
		_, branch.Children[key[matchlen]], err = t.insert(nil, pprefix, append(prefix, key[:matchlen+1]...), key[matchlen+1:], value)
		if err != nil {
			return false, nil, err
		}

		// Replace this shortNode with the branch if it occurs at index 0.
		if matchlen == 0 {
			if t.dag != nil {
				t.dag.addVertexAndEdge(fprefix, prefix, branch)
			}
			return true, branch, nil
		}
		if t.dag != nil {
			t.dag.addVertexAndEdge(append(prefix, key[:matchlen]...), append(prefix, key[:matchlen]...), branch)
		}
		// Otherwise, replace it with a short node leading up to the branch.
		nn := &shortNode{key[:matchlen], branch, t.newFlag()}
		if t.dag != nil {
			t.dag.addVertexAndEdge(fprefix, prefix, nn)
		}
		return true, nn, nil

	case *fullNode:
		dirty, nn, err := t.insert(n.Children[key[0]], append(prefix, fullNodeSuffix...), append(prefix, key[0]), key[1:], value)
		if !dirty || err != nil {
			return false, n, err
		}
		if t.dag != nil {
			//fmt.Printf("302: del vtx -> prefix: %x\n", append(prefix, fullNodeSuffix...))
			t.dag.delVertexAndEdge(append(prefix, fullNodeSuffix...))
		}
		n = n.copy()
		n.flags = t.newFlag()
		n.Children[key[0]] = nn
		if t.dag != nil {
			t.dag.addVertexAndEdge(fprefix, prefix, n)
		}
		return true, n, nil

	case nil:
		if t.dag != nil {
			//fmt.Printf("320: del vtx -> prefix: %x\n", append(prefix, key...))
			t.dag.delVertexAndEdge(append(prefix, key...))
		}
		nn := &shortNode{key, value, t.newFlag()}
		if t.dag != nil {
			t.dag.addVertexAndEdge(fprefix, prefix, nn)
		}
		return true, nn, nil

	case hashNode:
		// We've hit a part of the trie that isn't loaded yet. Load
		// the node and insert into it. This leaves all child nodes on
		// the path to the value in the trie.
		rn, err := t.resolveHash(n, prefix)
		if err != nil {
			return false, nil, err
		}
		dirty, nn, err := t.insert(rn, fprefix, prefix, key, value)
		if !dirty || err != nil {
			return false, rn, err
		}

		t.dag.addVertexAndEdge(fprefix, prefix, nn)
		return true, nn, nil

	default:
		panic(fmt.Sprintf("%T: invalid node: %v", n, n))
	}
}

// Delete removes any existing value for key from the trie.
func (t *Trie) Delete(key []byte) {
	if err := t.TryDelete(key); err != nil {
		log.Error(fmt.Sprintf("Unhandled trie error: %v", err))
	}
}

// TryDelete removes any existing value for key from the trie.
// If a node was not found in the database, a MissingNodeError is returned.
func (t *Trie) TryDelete(key []byte) error {
	k := keybytesToHex(key)
	_, n, err := t.delete(t.root, nil, k)
	if err != nil {
		return err
	}
	t.dag.delVertexAndEdgeByNode(nil, t.root)
	t.dag.addVertexAndEdge(nil, nil, n)
	t.root = n
	return nil
}

// delete returns the new root of the trie with key deleted.
// It reduces the trie to minimal form by simplifying
// nodes on the way up after deleting recursively.
func (t *Trie) delete(n node, prefix, key []byte) (bool, node, error) {
	switch n := n.(type) {
	case *shortNode:
		matchlen := prefixLen(key, n.Key)
		if matchlen < len(n.Key) {
			return false, n, nil // don't replace n on mismatch
		}
		if matchlen == len(key) {
			if t.dag != nil {
				//fmt.Printf("382: del vtx -> prefix: %x\n", append(prefix, key...))
				t.dag.delVertexAndEdge(append(prefix, key...))
			}
			return true, nil, nil // remove n entirely for whole matches
		}
		// The key is longer than n.Key. Remove the remaining suffix
		// from the subtrie. Child can never be nil here since the
		// subtrie must contain at least two other values with keys
		// longer than n.Key.
		dirty, child, err := t.delete(n.Val, append(prefix, key[:len(n.Key)]...), key[len(n.Key):])
		if !dirty || err != nil {
			return false, n, err
		}
		if t.dag != nil {
			//fmt.Printf("397: del vtx -> prefix: %x\n", append(prefix, n.Key...))
			t.dag.delVertexAndEdge(append(prefix, n.Key...))
		}
		switch child := child.(type) {
		case *shortNode:
			// Deleting from the subtrie reduced it to another
			// short node. Merge the nodes to avoid creating a
			// shortNode{..., shortNode{...}}. Use concat (which
			// always creates a new slice) instead of append to
			// avoid modifying n.Key since it might be shared with
			// other nodes.
			if t.dag != nil {
				//fmt.Printf("405: del vtx -> prefix: %x\n", append(prefix, concat(n.Key, child.Key...)...))
				t.dag.delVertexAndEdgeByNode(append(prefix, concat(n.Key, child.Key...)...), child.Val)
				//fmt.Printf("407: del vtx -> prefix: %x\n", append(prefix, concat(n.Key, child.Key...)...))
				t.dag.delVertexAndEdgeByNode(append(prefix, n.Key...), child)
				//fmt.Printf("409: add vtx -> prefix: %x\n", append(prefix, concat(n.Key, child.Key...)...))
				t.dag.addVertexAndEdge(append(prefix, concat(n.Key, child.Key...)...), append(prefix, concat(n.Key, child.Key...)...), child.Val)
			}
			return true, &shortNode{concat(n.Key, child.Key...), child.Val, t.newFlag()}, nil
		default:
			if t.dag != nil {
				//fmt.Printf("414: dev vtx -> prefix: %x\n", append(prefix, n.Key...))
				t.dag.delVertexAndEdgeByNode(append(prefix, n.Key...), child)
				//fmt.Printf("417: add vtx -> prefix: %x\n", append(prefix, n.Key...))
				t.dag.addVertexAndEdge(append(prefix, n.Key...), append(prefix, n.Key...), child)
			}
			return true, &shortNode{n.Key, child, t.newFlag()}, nil
		}

	case *fullNode:
		dirty, nn, err := t.delete(n.Children[key[0]], append(prefix, key[0]), key[1:])
		if !dirty || err != nil {
			return false, n, err
		}
		n = n.copy()
		n.flags = t.newFlag()
		n.Children[key[0]] = nn

		// Check how many non-nil entries are left after deleting and
		// reduce the full node to a short node if only one entry is
		// left. Since n must've contained at least two children
		// before deletion (otherwise it would not be a full node) n
		// can never be reduced to nil.
		//
		// When the loop is done, pos contains the index of the single
		// value that is left in n or -2 if n contains at least two
		// values.
		pos := -1
		for i, cld := range &n.Children {
			if cld != nil {
				if pos == -1 {
					pos = i
				} else {
					pos = -2
					break
				}
			}
		}
		if pos >= 0 {
			if t.dag != nil {
				//fmt.Printf("452: del vtx -> prefix: %x\n", append(prefix, fullNodeSuffix...))
				t.dag.delVertexAndEdge(append(prefix, fullNodeSuffix...))
			}
			if pos != 16 {
				// If the remaining entry is a short node, it replaces
				// n and its key gets the missing nibble tacked to the
				// front. This avoids creating an invalid
				// shortNode{..., shortNode{...}}.  Since the entry
				// might not be loaded yet, resolve it just for this
				// check.
				cnode, err := t.resolve(n.Children[pos], prefix)
				if err != nil {
					return false, nil, err
				}
				if cnode, ok := cnode.(*shortNode); ok {
					k := append([]byte{byte(pos)}, cnode.Key...)
					if t.dag != nil {
						//fmt.Printf("469: del vtx -> prefix: %x\n", append(prefix, byte(pos)))
						t.dag.delVertexAndEdgeByNode(append(prefix, byte(pos)), cnode)
						//fmt.Printf("473: add vtx -> prefix: %x\n", append(prefix, k...))
						t.dag.addVertexAndEdge(append(prefix, k...), append(prefix, k...), cnode.Val)
					}
					return true, &shortNode{k, cnode.Val, t.newFlag()}, nil
				}
			}
			// Otherwise, n is replaced by a one-nibble short node
			// containing the child.
			if t.dag != nil {
				//fmt.Printf("479: del vtx -> prefix: %x\n", append(prefix, byte(pos)))
				t.dag.delVertexAndEdgeByNode(append(prefix, byte(pos)), n.Children[pos])
				//fmt.Printf("484: add vtx -> prefix: %x\n", append(prefix, byte(pos)))
				t.dag.addVertexAndEdge(append(prefix, byte(pos)), append(prefix, byte(pos)), n.Children[pos])
			}
			return true, &shortNode{[]byte{byte(pos)}, n.Children[pos], t.newFlag()}, nil
		}
		// n still contains at least two values and cannot be reduced.
		if t.dag != nil {
			//fmt.Printf("491: add vtx -> prefix: %x\n", append(prefix, key[0]))
			t.dag.addVertexAndEdge(append(prefix, fullNodeSuffix...), append(prefix, key[0]), nn)
		}
		return true, n, nil

	case valueNode:
		return true, nil, nil

	case nil:
		return false, nil, nil

	case hashNode:
		// We've hit a part of the trie that isn't loaded yet. Load
		// the node and delete from it. This leaves all child nodes on
		// the path to the value in the trie.
		rn, err := t.resolveHash(n, prefix)
		if err != nil {
			return false, nil, err
		}
		dirty, nn, err := t.delete(rn, prefix, key)
		if !dirty || err != nil {
			return false, rn, err
		}
		return true, nn, nil

	default:
		panic(fmt.Sprintf("%T: invalid node: %v (%v)", n, n, key))
	}
}

func concat(s1 []byte, s2 ...byte) []byte {
	r := make([]byte, len(s1)+len(s2))
	copy(r, s1)
	copy(r[len(s1):], s2)
	return r
}

func (t *Trie) resolve(n node, prefix []byte) (node, error) {
	if n, ok := n.(hashNode); ok {
		return t.resolveHash(n, prefix)
	}
	return n, nil
}

func (t *Trie) resolveHash(n hashNode, prefix []byte) (node, error) {

	hash := common.BytesToHash(n)
	if node := t.db.node(hash); node != nil {
		return node, nil
	}
	return nil, &MissingNodeError{NodeHash: hash, Path: prefix}
}

// Root returns the root hash of the trie.
// Deprecated: use Hash instead.
func (t *Trie) Root() []byte { return t.Hash().Bytes() }

// Hash returns the root hash of the trie. It does not write to the
// database and can be used even if the trie doesn't have one.
func (t *Trie) Hash() common.Hash {
	hash, cached, _ := t.hashRoot(nil, nil)
	t.root = cached
	return common.BytesToHash(hash.(hashNode))
}

func (t *Trie) ParallelHash() common.Hash {
	hash, cached, err := t.parallelHashRoot(nil, nil)
	if err == nil {
		t.root = cached
	}
	return common.BytesToHash(hash.(hashNode))
}

// Commit writes all nodes to the trie's memory database, tracking the internal
// and external (for account tries) references.
func (t *Trie) Commit(onleaf LeafCallback) (root common.Hash, err error) {
	if t.db == nil {
		panic("commit called on trie with nil database")
	}
	hash, cached, err := t.hashRoot(t.db, onleaf)
	if err != nil {
		return common.Hash{}, err
	}
	t.root = cached
	return common.BytesToHash(hash.(hashNode)), nil
}

func (t *Trie) ParallelCommit(onleaf LeafCallback) (root common.Hash, err error) {
	if t.db == nil {
		panic("commit called on trie with nil database")
	}

	hash, cached, err := t.parallelHashRoot(t.db, onleaf)
	if err != nil {
		return common.Hash{}, err
	}
	t.root = cached

	// clear dag
	t.dag.clear()
	return common.BytesToHash(hash.(hashNode)), nil
}

func (t *Trie) hashRoot(db *Database, onleaf LeafCallback) (node, node, error) {
	if t.root == nil {
		return hashNode(emptyRoot.Bytes()), nil, nil
	}
	h := newHasher(onleaf)
	defer returnHasherToPool(h)
	return h.hash(t.root, db, true)
}

func (t *Trie) parallelHashRoot(db *Database, onleaf LeafCallback) (node, node, error) {
	if t.root == nil {
		return hashNode(emptyRoot.Bytes()), nil, nil
	}
	if len(t.dag.nodes) > 0 {
		//t.dag.init(t.root)
		return t.dag.hash(db, true, onleaf)
	} else {
		return t.hashRoot(db, onleaf)
	}
}

func (t *Trie) DeepCopyTrie() *Trie {
	cpyRoot := t.root
	switch n := t.root.(type) {
	case *shortNode:
		cpyRoot = n.copy()
	case *fullNode:
		cpyRoot = n.copy()
	}
	t.copyNode(cpyRoot)
	return &Trie{
		db:           t.db,
		root:         cpyRoot,
		//dag:          t.dag.DeepCopy(),
		dag: newTrieDag(),
	}
}

func (t *Trie) copyNode(n node) {

	//hash, dirty := n.cache()
	switch n := n.(type) {
	case *shortNode:
		if _, ok := n.Val.(valueNode); !ok {
			if _, dirty := n.Val.cache(); !dirty {
				if hash, _ := n.Val.cache(); len(hash) != 0 {
					n.Val = hash
				}
			} else {
				switch child := n.Val.(type) {
				case *shortNode:
					n.Val = child.copy()
				case *fullNode:
					n.Val = child.copy()
				}
				t.copyNode(n.Val)
			}
		}

	case *fullNode:
		for i := 0; i < len(n.Children); i++ {
			if n.Children[i] != nil {
				if _, ok := n.Children[i].(valueNode); !ok {

					if _, dirty := n.Children[i].cache(); !dirty {
						if hash, _ := n.Children[i].cache(); len(hash) != 0 {
							n.Children[i] = hash
						}
					} else {
						switch child := n.Children[i].(type) {
						case *shortNode:
							n.Children[i] = child.copy()
						case *fullNode:
							n.Children[i] = child.copy()
						}
						t.copyNode(n.Children[i])
					}
				}
			}
		}
	}
}
