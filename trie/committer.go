// Copyright 2020 The go-ethereum Authors
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

package trie

import (
	"fmt"

	"golang.org/x/crypto/sha3"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// leaf represents a trie leaf node
type leaf struct {
	blob   []byte      // raw blob of leaf
	parent common.Hash // the hash of parent node
}

// committer is the tool used for the trie Commit operation. The committer will
// capture all dirty nodes during the commit process and keep them cached in
// insertion order.
type committer struct {
	sha         crypto.KeccakState
	tmp         sliceBuffer
	encbuf      rlp.EncoderBuffer
	nodes       *NodeSet
	collectLeaf bool
}

// newCommitter creates a new committer or picks one from the pool.
func newCommitter(nodeset *NodeSet, collectLeaf bool) *committer {
	return &committer{
		tmp:         make(sliceBuffer, 0, 550), // cap is as large as a full fullNode.
		sha:         sha3.NewLegacyKeccak256().(crypto.KeccakState),
		encbuf:      rlp.NewEncoderBuffer(nil),
		nodes:       nodeset,
		collectLeaf: collectLeaf,
	}
}

// commit collapses a node down into a hash node and inserts it into the database
func (c *committer) commit(path []byte, n node, force bool) (node, node, error) {
	// If we're not storing the node, just hashing, use available cached data
	if hash, dirty := n.cache(); len(hash) != 0 {
		if !dirty {
			switch n.(type) {
			case *fullNode, *shortNode:
				return hash, hash, nil
			default:
				return hash, n, nil
			}
		}
	}
	// Trie not processed yet or needs storage, walk the children
	collapsed, cached, err := c.commitChildren(path, n)
	if err != nil {
		return hashNode{}, n, err
	}
	hashed, err := c.store(path, collapsed, force)
	if err != nil {
		return hashNode{}, n, err
	}
	// Cache the hash of the node for later reuse and remove
	// the dirty flag in commit mode. It's fine to assign these values directly
	// without copying the node first because hashChildren copies it.
	cachedHash, _ := hashed.(hashNode)

	/*
		// Mark the node as deleted if it's present in database previously.
		// It's equivalent as deletion from database's perspective.
		if prev := c.tracer.getPrev(path); len(prev) != 0 {
			c.nodes.markDeleted(path, prev)
		}
	*/

	switch cn := cached.(type) {
	case *shortNode:
		*cn.flags.hash = cachedHash
		*cn.flags.dirty = false
	case *fullNode:
		*cn.flags.hash = cachedHash
		*cn.flags.dirty = false
	}
	return hashed, cached, nil
}

func (c *committer) commitChildren(path []byte, original node) (node, node, error) {
	var err error

	switch n := original.(type) {
	case *shortNode:
		// Hash the short node's child, caching the newly hashed subtree
		collapsed, cached := n.copy(), n.copy()
		collapsed.Key = hexToCompact(n.Key)
		cached.Key = common.CopyBytes(n.Key)

		if _, ok := n.Val.(valueNode); !ok {
			collapsed.Val, cached.Val, err = c.commit(append(path, n.Key...), n.Val, false)
			if err != nil {
				return original, original, err
			}
		}
		return collapsed, cached, nil

	case *fullNode:
		// Hash the full node's children, caching the newly hashed subtrees
		collapsed, cached := n.copy(), n.copy()

		for i := 0; i < 16; i++ {
			if n.Children[i] != nil {
				collapsed.Children[i], cached.Children[i], err = c.commit(append(path, byte(i)), n.Children[i], false)
				if err != nil {
					return original, original, err
				}
			}
		}
		cached.Children[16] = n.Children[16]
		return collapsed, cached, nil

	default:
		// Value and hash nodes don't have children so they're left as were
		return n, original, nil
	}
}

// store hashes the node n and if we have a storage layer specified, it writes
// the key/value pair to it and tracks any node->child references as well as any
// node->external trie references.
func (c *committer) store(path []byte, n node, force bool) (node, error) {
	// Don't store hashes or empty nodes.
	if _, isHash := n.(hashNode); n == nil || isHash {
		return n, nil
	}
	// Generate the RLP encoding of the node
	n.encode(c.encbuf)
	enc := c.encodedBytes()
	if len(enc) < 32 && !force {
		// The node is embedded in its parent, in other words, this node
		// will not be stored in the database independently, mark it as
		// deleted only if the node was existent in database before.
		if _, ok := c.nodes.accessList[string(path)]; ok {
			c.nodes.markDeleted(path)
		}
		return n, nil // Nodes smaller than 32 bytes are stored inside their parent
	}

	// Larger nodes are replaced by their hash and stored in the database.
	hash, _ := n.cache()
	if len(hash) == 0 {
		hash = c.makeHashNode(enc)
	}

	// We have the hash already, estimate the RLP encoding-size of the node.
	// The size is used for mem tracking, does not need to be exact
	var (
		nhash = common.BytesToHash(hash)
		mnode = &memoryNode{
			hash: nhash,
			node: nodeToBytes(n),
		}
	)
	// Collect the dirty node to nodeset for return.
	c.nodes.markUpdated(path, mnode)

	// Collect the corresponding leaf node if it's required. We don't check
	// full node since it's impossible to store value in fullNode. The key
	// length of leaves should be exactly same.
	if c.collectLeaf {
		switch n := n.(type) {
		case *shortNode:
			if val, ok := n.Val.(valueNode); ok {
				c.nodes.addLeaf(&leaf{blob: val, parent: nhash})
			}
		case *fullNode:
			for i := 0; i < 16; i++ {
				if val, ok := n.Children[i].(valueNode); ok {
					c.nodes.addLeaf(&leaf{blob: val, parent: nhash})
				}
			}
		}
	}
	/*
		// We are pooling the trie nodes into an intermediate memory cache
		//hash2 := common.BytesToHash(hash)
		db.lock.Lock()
		db.insert(hash2, estimateSize(n), n)
		db.insertFreshNode(hash2)
		db.lock.Unlock()

		// Track external references from account->storage trie
		if c.onleaf != nil {
			switch n := n.(type) {
			case *shortNode:
				if child, ok := n.Val.(valueNode); ok {
					c.onleaf(nil, nil, child, hash2, nil)
				}
			case *fullNode:
				for i := 0; i < 16; i++ {
					if child, ok := n.Children[i].(valueNode); ok {
						c.onleaf(nil, nil, child, hash2, nil)
					}
				}
			}
		}*/

	return hash, nil
}

// encodedBytes returns the result of the last encoding operation on h.encbuf.
// This also resets the encoder buffer.
//
// All node encoding must be done like this:
//
//	node.encode(h.encbuf)
//	enc := h.encodedBytes()
//
// This convention exists because node.encode can only be inlined/escape-analyzed when
// called on a concrete receiver type.
func (c *committer) encodedBytes() []byte {
	c.tmp = c.encbuf.AppendToBytes(c.tmp[:0])
	c.encbuf.Reset(nil)
	return c.tmp
}

// makeHashNode hashes the provided data
func (c *committer) makeHashNode(data []byte) hashNode {
	n := make(hashNode, c.sha.Size())
	c.sha.Reset()
	c.sha.Write(data)
	c.sha.Read(n)
	return n
}

// mptResolver the children resolver in merkle-patricia-tree.
type mptResolver struct{}

// ForEach implements childResolver, decodes the provided node and
// traverses the children inside.
func (resolver mptResolver) forEach(node []byte, onChild func(common.Hash)) {
	forGatherChildren(mustDecodeNodeUnsafe(nil, node), onChild)
}

// forGatherChildren traverses the node hierarchy and invokes the callback
// for all the hashnode children.
func forGatherChildren(n node, onChild func(hash common.Hash)) {
	switch n := n.(type) {
	case *shortNode:
		forGatherChildren(n.Val, onChild)
	case *fullNode:
		for i := 0; i < 16; i++ {
			forGatherChildren(n.Children[i], onChild)
		}
	case hashNode:
		onChild(common.BytesToHash(n))
	case valueNode, nil:
	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}
