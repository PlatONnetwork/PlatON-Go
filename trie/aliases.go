package trie

import (
	"github.com/PlatONnetwork/PlatON-Go/trie/triedb/hashdb"
	"github.com/PlatONnetwork/PlatON-Go/trie/trienode"
)

// Type aliases for backward compatibility.
type (
	NodeSet       = trienode.NodeSet
	MergedNodeSet = trienode.MergedNodeSet
	CachedNode    = hashdb.CachedNode
)

var (
	NewNodeSet       = trienode.NewNodeSet
	NewMergedNodeSet = trienode.NewMergedNodeSet
	NewWithNodeSet   = trienode.NewWithNodeSet
)
