package trie

import (
	"container/list"
	"math"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/panjf2000/ants/v2"
)

const (
	invalidId uint64 = math.MaxUint64
)

//var fullNodeSuffix = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
var fullNodeSuffix = []byte("fullnode")

type Vertex struct {
	inDegree uint32
	outEdge  []uint64
}

type DAG struct {
	vtxs     map[uint64]*Vertex
	topLevel *list.List

	lock sync.Mutex
	cv   *sync.Cond

	totalVertexs  uint32
	totalConsumed uint32
}

func NewDAG() *DAG {
	dag := &DAG{
		vtxs:          make(map[uint64]*Vertex),
		topLevel:      list.New(),
		totalConsumed: 0,
	}
	dag.cv = sync.NewCond(&dag.lock)

	return dag
}

func (d *DAG) addVertex(id uint64) {
	if _, ok := d.vtxs[id]; !ok {
		d.totalVertexs++
	}
	d.vtxs[id] = &Vertex{
		inDegree: 0,
		outEdge:  make([]uint64, 0),
	}
}

func (d *DAG) addEdge(from, to uint64) {
	if _, ok := d.vtxs[from]; !ok {
		d.vtxs[from] = &Vertex{
			inDegree: 0,
			outEdge:  make([]uint64, 0),
		}
	}
	vtx := d.vtxs[from]
	vtx.outEdge = append(vtx.outEdge, to)

	if _, ok := d.vtxs[to]; !ok {
		d.vtxs[to] = &Vertex{
			inDegree: 0,
			outEdge:  make([]uint64, 0),
		}
	}
	d.vtxs[to].inDegree += 1
}

func (d *DAG) generate() {
	for k, v := range d.vtxs {
		if v.inDegree == 0 {
			d.topLevel.PushBack(k)
		}
	}
}

func (d *DAG) waitPop() uint64 {
	if d.hasFinished() {
		return invalidId
	}

	d.cv.L.Lock()
	defer d.cv.L.Unlock()
	if d.topLevel.Len() == 0 && !d.hasFinished() {
		d.cv.Wait()
	}

	if d.hasFinished() || d.topLevel.Len() == 0 {
		return invalidId
	}

	el := d.topLevel.Front()
	id := el.Value.(uint64)
	d.topLevel.Remove(el)
	return id
}

func (d *DAG) hasFinished() bool {
	return d.totalConsumed >= d.totalVertexs
}

func (d *DAG) consume(id uint64) uint64 {
	producedNum := 0
	var nextID uint64 = invalidId
	var degree uint32 = 0

	for _, k := range d.vtxs[id].outEdge {
		vtx := d.vtxs[k]
		degree = atomic.AddUint32(&vtx.inDegree, ^uint32(0))
		if degree == 0 {
			producedNum += 1
			if producedNum == 1 {
				nextID = k
			} else {
				d.lock.Lock()
				d.topLevel.PushBack(k)
				d.lock.Unlock()
			}
		}
	}

	if atomic.AddUint32(&d.totalConsumed, 1) == d.totalVertexs {
		d.cv.L.Lock()
		d.cv.Broadcast()
		d.cv.L.Unlock()
	}

	return nextID
}

func (d *DAG) clear() {
	d.vtxs = make(map[uint64]*Vertex)
	d.topLevel = list.New()
	d.totalConsumed = 0
	d.totalVertexs = 0
}

type DAGNode struct {
	collapsed node
	cached    node
	pid       uint64
	idx       int
}

// TrieDAG
type TrieDAG struct {
	nodes map[uint64]*DAGNode

	dag *DAG

	cachegen   uint16
	cachelimit uint16
}

func NewTrieDAG(cachegen, cachgelimit uint16) *TrieDAG {
	return &TrieDAG{
		nodes:      make(map[uint64]*DAGNode),
		dag:        NewDAG(),
		cachegen:   cachegen,
		cachelimit: cachgelimit,
	}
}

func (td *TrieDAG) init(root node) {
	hashed, _ := td.hasChildren(root, 0, nil)
	if hashed != nil {
		td.nodes[0] = &DAGNode{
			collapsed: hashed,
			cached:    root,
			pid:       0,
		}
		td.dag.addVertex(0)
	}
	td.dag.generate()
}

func (td *TrieDAG) hasChildren(origin node, pid uint64, prefix []byte) (node, node) {
	if origin == nil {
		return nil, nil
	}
	if hash, dirty := origin.cache(); len(hash) != 0 {
		if origin.canUnload(td.cachegen, td.cachelimit) {
			cacheUnloadCounter.Inc(1)
			return hash, hash
		}
		if !dirty {
			return hash, origin
		}
	}

	switch n := origin.(type) {
	case *shortNode:
		collapsed, cached := n.copy(), n.copy()
		collapsed.Key = hexToCompact(n.Key)
		cached.Key = common.CopyBytes(n.Key)

		//id := xxhash.Sum64(append(prefix, n.Key...))
		idx := 0
		if len(prefix) > 0 {
			idx = int(prefix[len(prefix)-1])
		}
		id := (pid * 100) + 1 + uint64(idx)
		td.nodes[id] = &DAGNode{
			collapsed: collapsed,
			cached:    cached,
			pid:       pid,
			idx:       idx,
		}
		//if len(prefix) > 0 {
		//	td.nodes[id].idx = int(prefix[len(prefix)-1])
		//}

		td.dag.addVertex(id)

		if pid > 0 {
			td.dag.addEdge(id, pid)
		}

		if _, ok := n.Val.(valueNode); !ok {
			hash, org := td.hasChildren(n.Val, id, append(prefix, n.Key...))
			if hash != nil && org != nil {
				collapsed.Val = hash
				cached.Val = org
			}
		}

	case *fullNode:
		collapsed, cached := n.copy(), n.copy()

		dagNode := &DAGNode{
			collapsed: collapsed,
			cached:    cached,
			pid:       pid,
		}
		if len(prefix) > 0 {
			dagNode.idx = int(prefix[len(prefix)-1])
		}

		//id := xxhash.Sum64(append(prefix, fullNodeSuffix...))
		id := (pid * 100) + 1 + uint64(dagNode.idx)
		td.nodes[id] = dagNode
		td.dag.addVertex(id)
		if pid > 0 {
			td.dag.addEdge(id, pid)
		}

		for i := 0; i < 16; i++ {
			if n.Children[i] != nil {
				hash, org := td.hasChildren(n.Children[i], id, append(prefix, byte(i)))
				if hash != nil && org != nil {
					collapsed.Children[i] = hash
					cached.Children[i] = org
				}
			}
		}
		cached.Children[16] = n.Children[16]
	default:
		// Value and hash nodes don't have children so they're left as were
	}
	return nil, nil
}

func (td *TrieDAG) hash(db *Database, force bool, onleaf LeafCallback) (node, node, error) {
	var wg sync.WaitGroup
	var errDone common.AtomicBool
	var e atomic.Value // error
	var resHash node = hashNode{}
	var newRoot node
	numCPU := runtime.NumCPU()

	process := func() {
		hasher := newHasher(td.cachegen, td.cachelimit, onleaf)

		id := td.dag.waitPop()
		if id == invalidId {
			returnHasherToPool(hasher)
			wg.Done()
			return
		}

		for id != invalidId {
			n := td.nodes[id]

			tmpForce := false
			if n.pid == 0 {
				tmpForce = force
			}
			hashed, err := hasher.store(n.collapsed, db, tmpForce)
			if err != nil {
				e.Store(err)
				errDone.Set(true)
				break
			}

			if n.pid > 0 {
				p := td.nodes[n.pid]
				switch ptype := p.collapsed.(type) {
				case *shortNode:
					ptype.Val = hashed
				case *fullNode:
					ptype.Children[n.idx] = hashed
				}
			}

			cachedHash, _ := hashed.(hashNode)
			switch cn := n.cached.(type) {
			case *shortNode:
				*cn.flags.hash = cachedHash
				if db != nil {
					*cn.flags.dirty = false
				}
			case *fullNode:
				*cn.flags.hash = cachedHash
				if db != nil {
					*cn.flags.dirty = false
				}
			}

			id = td.dag.consume(id)
			if n.pid == 0 {
				resHash = hashed
				newRoot = n.cached
				break
			}

			if errDone.IsSet() {
				break
			}

			if id == invalidId && !td.dag.hasFinished() {
				id = td.dag.waitPop()
			}
		}
		returnHasherToPool(hasher)
		wg.Done()
	}

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		_ = ants.Submit(process)
	}

	wg.Wait()

	if e.Load() != nil && e.Load().(error) != nil {
		return hashNode{}, nil, e.Load().(error)
	}
	return resHash, newRoot, nil
}
