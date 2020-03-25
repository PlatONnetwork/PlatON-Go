package trie

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/cespare/xxhash"
	"github.com/panjf2000/ants/v2"
	"github.com/petermattis/goid"
)

var fullNodeSuffix = []byte("fullnode")

// dagNode
type dagNode struct {
	collapsed node
	cached    node

	pid uint64
	idx int
}

// trieDag
type trieDag struct {
	nodes map[uint64]*dagNode
	dag   *dag

	lock sync.Mutex

	cachegen   uint16
	cachelimit uint16
}

func newTrieDag() *trieDag {
	return &trieDag{
		nodes:      make(map[uint64]*dagNode),
		dag:        newDag(),
		cachegen:   0,
		cachelimit: 0,
	}
}

func (td *trieDag) addVertexAndEdge(pprefix, prefix []byte, n node) {
	td.lock.Lock()
	defer td.lock.Unlock()
	td.internalAddVertexAndEdge(pprefix, prefix, n, true)
}

func (td *trieDag) internalAddVertexAndEdge(pprefix, prefix []byte, n node, recursive bool) {
	var pid uint64
	if len(pprefix) > 0 {
		pid = xxhash.Sum64(pprefix)
	}

	cachedHash := func(n node) (node, bool) {
		if hash, _ := n.cache(); len(hash) != 0 {
			return hash, true
		}
		return n, false
	}

	switch nc := n.(type) {
	case *shortNode:
		collapsed, cached := nc.copy(), nc.copy()
		collapsed.Key = hexToCompact(nc.Key)
		cached.Key = common.CopyBytes(nc.Key)

		hash, has := cachedHash(nc.Val)
		if has {
			hash, _ = hash.(hashNode)
			collapsed.Val = hash
		}

		id := xxhash.Sum64(append(prefix, nc.Key...))
		td.nodes[id] = &dagNode{
			collapsed: collapsed,
			cached:    cached,
			pid:       pid,
		}
		if len(prefix) > 0 {
			td.nodes[id].idx = int(prefix[len(prefix)-1])
		}
		td.dag.addVertex(id)

		if pid > 0 {
			td.dag.addEdge(id, pid)
		}
		//fmt.Printf("add short -> pprefix: %x, prefix: %x\n", pprefix, append(prefix, nc.Key...))
		//fmt.Printf("add short -> id: %d, pid: %d\n", id, pid)

	case *fullNode:
		collapsed, cached := nc.copy(), nc.copy()
		cached.Children[16] = nc.Children[16]

		dagNode := &dagNode{
			collapsed: collapsed,
			cached:    cached,
			pid:       pid,
		}
		if len(prefix) > 0 {
			dagNode.idx = int(prefix[len(prefix)-1])
		}

		id := xxhash.Sum64(append(prefix, fullNodeSuffix...))
		td.nodes[id] = dagNode
		td.dag.addVertex(id)
		if pid > 0 {
			td.dag.addEdge(id, pid)
		}
		//fmt.Printf("add full -> pprefix: %x, prefix: %x\n", pprefix, append(prefix, fullNodeSuffix...))
		//fmt.Printf("add full -> id: %d, pid: %d\n", id, pid)

		if recursive {
			for i := 0; i < 16; i++ {
				if nc.Children[i] != nil {
					cn := nc.Children[i]
					//_, dirty := cn.cache()
					//if !dirty {
					td.internalAddVertexAndEdge(append(prefix, fullNodeSuffix...), append(prefix, byte(i)), cn, false)
					//}
				}
			}
		}
	}
}

func (td *trieDag) delVertexAndEdge(key []byte) {
	id := xxhash.Sum64(key)
	td.delVertexAndEdgeByID(id)
}

func (td *trieDag) delVertexAndEdgeByID(id uint64) {
	td.lock.Lock()
	defer td.lock.Unlock()
	//td.dag.delEdge(id)
	td.dag.delVertex(id)
	delete(td.nodes, id)
	//fmt.Printf("del: %d\n", id)
}

func (td *trieDag) delVertexAndEdgeByNode(prefix []byte, n node) {
	var id uint64
	switch nc := n.(type) {
	case *shortNode:
		id = xxhash.Sum64(append(prefix, nc.Key...))
	case *fullNode:
		id = xxhash.Sum64(append(prefix, fullNodeSuffix...))
	}
	td.delVertexAndEdgeByID(id)
}

func (td *trieDag) replaceEdge(old, new []byte) {
	opid := xxhash.Sum64(old)
	npid := xxhash.Sum64(new)

	for id, vtx := range td.dag.vtxs {
		for _, pid := range vtx.outEdge {
			if opid == pid {
				vtx.outEdge = make([]uint64, 0)
				vtx.outEdge = append(vtx.outEdge, npid)
				td.nodes[id].pid = npid
			}
		}
	}
}

func (td *trieDag) reset() {
	td.lock.Lock()
	defer td.lock.Unlock()
	td.dag.reset()
}

func (td *trieDag) clear() {
	td.lock.Lock()
	defer td.lock.Unlock()

	td.dag.clear()
	td.nodes = make(map[uint64]*dagNode)
}

func (td *trieDag) hash(db *Database, force bool, onleaf LeafCallback) (node, node, error) {
	td.lock.Lock()
	defer td.lock.Unlock()

	td.dag.generate()

	var wg sync.WaitGroup
	var errDone common.AtomicBool
	var e atomic.Value // error
	var resHash node = hashNode{}
	var newRoot node
	numCPU := runtime.NumCPU()

	cachedHash := func(n, c node) (node, node, bool) {
		if hash, dirty := c.cache(); len(hash) != 0 {
			if db == nil {
				return hash, c, true
			}

			if c.canUnload(td.cachegen, td.cachelimit) {
				cacheUnloadCounter.Inc(1)
				return hash, hash, true
			}
			if !dirty {
				return hash, c, true
			}
		}
		return n, c, false
	}

	process := func() {
		log.Debug("Do hash", "me", fmt.Sprintf("%p", td), "routineID", goid.Get(), "dag", fmt.Sprintf("%p", td.dag), "nodes", len(td.nodes), "topLevel", td.dag.topLevel.Len(), "consumed", td.dag.totalConsumed, "vtxs", td.dag.totalVertexs, "cv", td.dag.cv)
		hasher := newHasher(td.cachegen, td.cachelimit, onleaf)

		id := td.dag.waitPop()
		if id == invalidID {
			returnHasherToPool(hasher)
			wg.Done()
			return
		}

		var hashed node
		var cached node
		var err error
		var hasCache bool
		for id != invalidID {
			n := td.nodes[id]

			tmpForce := false
			if n.pid == 0 {
				tmpForce = force
			}

			hashed, cached, hasCache = cachedHash(n.collapsed, n.cached)
			if !hasCache {
				hashed, err = hasher.store(n.collapsed, db, tmpForce)
				if err != nil {
					e.Store(err)
					errDone.Set(true)
					break
				}
				cached = n.cached
			}

			if n.pid > 0 {
				p := td.nodes[n.pid]
				switch ptype := p.collapsed.(type) {
				case *shortNode:
					ptype.Val = hashed
				case *fullNode:
					ptype.Children[n.idx] = hashed
				}

				if _, ok := cached.(hashNode); ok {
					switch nc := p.cached.(type) {
					case *shortNode:
						nc.Val = cached
					case *fullNode:
						nc.Children[n.idx] = cached
					}
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

			if id == invalidID && !td.dag.hasFinished() {
				id = td.dag.waitPop()
			}
		}
		returnHasherToPool(hasher)
		wg.Done()
		log.Error("Work done", "me", fmt.Sprintf("%p", td), "routineID", goid.Get(), "consumed", td.dag.totalConsumed, "vtxs", td.dag.totalVertexs)
	}

	wg.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		_ = ants.Submit(process)
	}

	wg.Wait()
	td.dag.reset()

	if e.Load() != nil && e.Load().(error) != nil {
		return hashNode{}, nil, e.Load().(error)
	}
	return resHash, newRoot, nil
}

func (td *trieDag) init(root node) {
	td.lock.Lock()
	defer td.lock.Unlock()
	td.dag.generate()
}
