package exec

import (
	"github.com/PlatONnetwork/PlatON-Go/log"
	"fmt"
	"math"
	"sync"
)

type tree []int
type TreePool struct {
	trees     map[int][]tree
	pool      map[int]*sync.Pool
	cacheSize int
}

func initTree(t tree, size int) {
	nodeSize := size * 2
	for i := 0; i < 2*size-1; i++ {
		if isPowOf2(i + 1) {
			nodeSize /= 2
		}
		t[i] = nodeSize
	}
}

func buildTree(size int) tree {

	if size < 1 || !isPowOf2(size) {
		panic(fmt.Errorf("build tree failed,wrong Size:%d", size))
	}
	tree := make([]int, size*2-1)
	initTree(tree, size)

	return tree
}

func NewTreePool(poolSize int, cacheSize int) TreePool {

	trees := make(map[int][]tree, 0)
	pool := make(map[int]*sync.Pool, 0)

	for i := 0; i < poolSize; i++ {
		size := int(math.Pow(2, float64(i))) * DefaultPageSize
		treeC := buildTree(size)
		for j := 0; j < cacheSize; j++ {
			t := make([]int, (2*size)-1)
			copy(t, treeC)
			trees[i] = append(trees[i], t)
		}
		pool[i] = &sync.Pool{
			New: func() interface{} {
				return buildTree(size)
			},
		}
	}

	return TreePool{trees, pool, cacheSize}
}

func (tp *TreePool) GetTree(pages int) tree {
	pages = fixSize(pages)
	key := int(math.Log2(float64(pages)))
	treeArr, ok := tp.trees[key]
	if !ok {
		tp.pool[key] = &sync.Pool{
			New: func() interface{} {
				return buildTree(pages * DefaultPageSize)
			},
		}
		return tp.pool[key].Get().(tree)
	}
	if len(treeArr) > 0 {
		tree := treeArr[0]
		//删除
		tp.trees[key] = append(treeArr[:0], treeArr[1:]...)
		return tree
	} else {
		return tp.pool[key].Get().(tree)
	}
}

func (tp *TreePool) PutTree(tree []int) {
	size := (len(tree) + 1) / 2
	pages := size / DefaultPageSize
	key := int(math.Log2(float64(pages)))

	//判断是否需要reset
	if tree[0] != size {
		//需要reset
		log.Debug("reset memory tree...")
		reset(tree, tp.trees[key], size)
	}
	treeArr, ok := tp.trees[key]
	if !ok || len(treeArr) >= tp.cacheSize {
		tp.pool[key].Put(tree)
		return
	}
	tp.trees[key] = append(treeArr, tree)
}

func reset(t tree, trees []tree, size int) {
	if trees != nil {
		for _, treeC := range trees {
			if treeC != nil {
				copy(t, treeC)
				return
			}
		}
	}
	initTree(t, size)
}
