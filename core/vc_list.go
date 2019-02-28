package core

import (
	"container/heap"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type vcHeap []*types.TransactionWrap

func (h vcHeap) Len() int { return len(h) }

func (h vcHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h vcHeap) Less(i, j int) bool {
	return h[i].Bn < h[j].Bn
}

func (h *vcHeap) Push(x interface{}) {
	*h = append(*h, x.(*types.TransactionWrap))
}

func (h *vcHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type vcList struct {
	all    *vcLookup
	items  *vcHeap
	stales int
}

func newVCList(all *vcLookup) *vcList {
	return &vcList{
		all:   all,
		items: new(vcHeap),
	}
}

func (l *vcList) Put(tx *types.TransactionWrap) {
	heap.Push(l.items, tx)
}

func (l *vcList) Pop() *types.TransactionWrap {
	return heap.Pop(l.items).(*types.TransactionWrap)
}
