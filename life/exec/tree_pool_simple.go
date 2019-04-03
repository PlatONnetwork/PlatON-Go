package exec

//
//import (
//	"fmt"
//	"math"
//	"sync"
//)
//
//type tree []int
//type TreePool struct {
//	//trees    map[int][]tree
//	//maxPage int
//	pool map[int]*sync.Pool
//}
//
//func buildTree(size int) tree {
//
//	fmt.Println("new tree.........")
//	if size < 1 || !isPowOf2(size) {
//		panic(fmt.Errorf("initTree failed,wrong Size:%d", size))
//	}
//
//	tree := make([]int, (2*size)-1)
//
//	nodeSize := size * 2
//	for i := 0; i < (2*size)-1; i++ {
//		if isPowOf2(i + 1) {
//			nodeSize /= 2
//		}
//		tree[i] = nodeSize
//	}
//	return tree
//}
//
///**
// */
//func NewTreePool(poolSize int, size int) TreePool {
//	pool := make(map[int]*sync.Pool, 0)
//	for i := 0; i < poolSize; i++ {
//		size := int(math.Pow(2, float64(i))) * DefaultPageSize
//		pool[i] = &sync.Pool{
//			New: func() interface{} {
//				return buildTree(size)
//			},
//		}
//		tree := buildTree(size)
//		pool[i].Put(tree)
//	}
//
//	return TreePool{pool}
//}
//
//func (tp *TreePool) GetTree(pages int) tree {
//	pages = fixSize(pages)
//	key := int(math.Log2(float64(pages)))
//	pool, ok := tp.pool[key]
//	if ok {
//		return pool.Get().(tree)
//	}
//
//	tp.pool[key] = &sync.Pool{
//		New: func() interface{} {
//			return buildTree(pages * DefaultPageSize)
//		},
//	}
//
//	return tp.pool[key].Get().(tree)
//}
//
//func (tp *TreePool) PutTree(tree tree) {
//
//	size := (len(tree) + 1) / 2
//	pages := size / DefaultPageSize
//	key := int(math.Log2(float64(pages)))
//	tp.pool[key].Put(tree)
//}
