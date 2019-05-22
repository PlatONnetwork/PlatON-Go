package core

import (
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"container/heap"
)

type mpcHeap []*types.TransactionWrap

func (h mpcHeap) Len() int { return len(h) }

func (h mpcHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h mpcHeap) Less(i, j int) bool {
	return h[i].Bn < h[j].Bn
}

func (h *mpcHeap) Push(x interface{}) {
	*h = append(*h, x.(*types.TransactionWrap))
}

func (h *mpcHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type mpcList struct {
	all 	*mpcLookup
	items	*mpcHeap
	stales	int
}

func newMpcList(all *mpcLookup) *mpcList {
	return &mpcList{
		all : all,
		items: new(mpcHeap),
	}
}

func (l *mpcList) Put(tx *types.TransactionWrap) {
	heap.Push(l.items, tx)
}

func (l *mpcList) Pop() *types.TransactionWrap {
	return heap.Pop(l.items).(*types.TransactionWrap)
}

















