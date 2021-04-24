// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package snapshotdb

import (
	"bytes"
	"container/heap"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

func newRankingHeap(hepNum int) *rankingHeap {
	r := new(rankingHeap)
	r.hepMaxNum = hepNum
	r.handledKey = make(map[string]struct{}, 300)
	r.heap = make(kvsMaxToMin, 0)
	return r
}

type rankingHeap struct {
	handledKey map[string]struct{}
	//max heap
	heap      kvsMaxToMin
	hepMaxNum int
}

//   the key is gt than or eq than   heap top
func (r *rankingHeap) geMaxHeap(k []byte) bool {
	if bytes.Compare(k, r.heap[0].key) > 0 {
		return true
	}
	return false
}

func (r *rankingHeap) addHandledKey(key []byte) {
	tmpKey := make([]byte, len(key))
	copy(tmpKey, key)
	r.handledKey[string(tmpKey)] = struct{}{}
}

func (r *rankingHeap) findHandledKey(key []byte) bool {
	if _, ok := r.handledKey[string(key)]; ok {
		return true
	}
	return false
}

// the heap length is less or eq than r.hepMaxNum.
// except baseDB, every block must range.
// find key, continue ,handle key add to HandledKey.
// the key must less than the top.
func (r *rankingHeap) itr2Heap(itr iterator.Iterator, baseDB, deepCopy bool) {
	unlimited := r.hepMaxNum <= 0
	if unlimited {
		for itr.Next() {
			k, v := itr.Key(), itr.Value()
			if r.findHandledKey(k) {
				continue
			}
			r.push2Heap(k, v, deepCopy)
			r.addHandledKey(k)
		}
	} else {
		for itr.Next() {
			k, v := itr.Key(), itr.Value()
			if r.findHandledKey(k) {
				continue
			}
			if r.heap.Len() >= r.hepMaxNum && bytes.Compare(k, r.heap[0].key) >= 0 {
				r.addHandledKey(k)
				break
			}
			r.push2Heap(k, v, deepCopy)
			r.addHandledKey(k)
			for r.heap.Len() > r.hepMaxNum {
				heap.Pop(&r.heap)
			}
		}
	}

	itr.Release()
}

func (r *rankingHeap) push2Heap(k, v []byte, deepCopy bool) {
	condtion := v == nil || len(v) == 0
	if !condtion {
		//if r.hepMaxNum > 0 && r.heap.Len() >= r.hepMaxNum {
		//	heap.Pop(&r.heap)
		//}
		if deepCopy {
			sk, sv := make([]byte, len(k)), make([]byte, len(v))
			copy(sk, k)
			copy(sv, v)
			heap.Push(&r.heap, kv{key: sk, value: sv})
		} else {
			heap.Push(&r.heap, kv{k, v})
		}
	}
}
