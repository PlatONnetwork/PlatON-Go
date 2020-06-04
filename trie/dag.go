package trie

import (
	"container/list"
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

const invalidID = math.MaxUint64

type vertex struct {
	inDegree uint32
	outEdge  []uint64
}

type dag struct {
	vtxs     map[uint64]*vertex
	topLevel *list.List

	lock sync.Mutex
	cv   *sync.Cond

	totalVertexs  uint32
	totalConsumed uint32
}

func newDag() *dag {
	dag := &dag{
		vtxs:          make(map[uint64]*vertex),
		topLevel:      list.New(),
		totalConsumed: 0,
	}
	dag.cv = sync.NewCond(&dag.lock)

	return dag
}

func (d *dag) addVertex(id uint64) {
	if _, ok := d.vtxs[id]; !ok {
		d.vtxs[id] = &vertex{
			inDegree: 0,
			outEdge:  make([]uint64, 0),
		}
		d.totalVertexs++
	}
}

func (d *dag) delVertex(id uint64) {
	if _, ok := d.vtxs[id]; ok {
		d.totalVertexs--
		delete(d.vtxs, id)
	}
}

func (d *dag) addEdge(from, to uint64) {
	if _, ok := d.vtxs[from]; !ok {
		d.vtxs[from] = &vertex{
			inDegree: 0,
			outEdge:  make([]uint64, 0),
		}
		d.totalVertexs++
	}
	vtx := d.vtxs[from]
	found := false
	for _, t := range vtx.outEdge {
		if t == to {
			found = true
			break
		}
	}
	if !found {
		vtx.outEdge = append(vtx.outEdge, to)
	}
}

func (d *dag) generate() {
	for id, v := range d.vtxs {
		for _, pid := range v.outEdge {
			if d.vtxs[pid] == nil {
				panic(fmt.Sprintf("%d: out found parent id %d", id, pid))
			}
			d.vtxs[pid].inDegree++
		}
	}

	for k, v := range d.vtxs {
		if v.inDegree == 0 {
			d.topLevel.PushBack(k)
		}
	}
}

func (d *dag) waitPop() uint64 {
	if d.hasFinished() {
		return invalidID
	}

	d.cv.L.Lock()
	defer d.cv.L.Unlock()
	for d.topLevel.Len() == 0 && !d.hasFinished() {
		//log.Error("Wait Pop", "dag", fmt.Sprintf("%p", d), "topLevel", d.topLevel.Len(), "consumed", d.totalConsumed, "vtxs", d.totalVertexs, "degreeGt", d.degreeGt(), "cv", d.cv)
		d.cv.Wait()
	}

	if d.hasFinished() || d.topLevel.Len() == 0 {
		return invalidID
	}

	el := d.topLevel.Front()
	id := el.Value.(uint64)
	d.topLevel.Remove(el)
	return id
}

func (d *dag) hasFinished() bool {
	return atomic.LoadUint32(&d.totalConsumed) >= d.totalVertexs
}

func (d *dag) consume(id uint64) uint64 {
	var (
		producedNum        = 0
		nextID      uint64 = invalidID
		degree      uint32 = 0
	)

	for _, k := range d.vtxs[id].outEdge {
		vtx := d.vtxs[k]
		degree = atomic.AddUint32(&vtx.inDegree, ^uint32(0))
		if degree == 0 {
			producedNum += 1
			if producedNum == 1 {
				nextID = k
			} else {
				d.cv.L.Lock()
				d.topLevel.PushBack(k)
				d.cv.Signal()
				d.cv.L.Unlock()
			}
		}
	}

	if atomic.AddUint32(&d.totalConsumed, 1) == d.totalVertexs {
		d.cv.L.Lock()
		d.cv.Broadcast()
		d.cv.L.Unlock()
		log.Trace("Consume done", "consumed", d.totalConsumed, "vtxs", d.totalVertexs)
	}
	return nextID
}

func (d *dag) clear() {
	d.vtxs = make(map[uint64]*vertex)
	d.topLevel = list.New()
	d.totalConsumed = 0
	d.totalVertexs = 0
}

func (d *dag) reset() {
	for _, v := range d.vtxs {
		v.inDegree = 0
	}
	d.topLevel = list.New()
	d.totalConsumed = 0
}
