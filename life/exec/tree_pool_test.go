package exec

import (
	"math"
	"testing"
)

var (
	poolSize  = 4
	cacheSize = 4
)

func TestNewTreePool(t *testing.T) {
	treePool := NewTreePool(poolSize, cacheSize)
	for k, v := range treePool.pool {
		expect := int(math.Pow(2, float64(k)) * DefaultPageSize)
		if get := v.Get().(tree)[0]; get != expect {
			t.Fatalf("new pool error ,expect tree[0]=%d,get=%d", expect, get)
		}
	}
	for k, v := range treePool.trees {
		expect := int(math.Pow(2, float64(k)) * DefaultPageSize)
		for _, tree := range v {
			if get := tree[0]; get != expect {
				t.Fatalf("new pool error ,expect tree[0]=%d,get=%d", expect, get)
			}
		}
	}
}

func TestTreePool_GetPutTree(t *testing.T) {

	treePool := NewTreePool(poolSize, cacheSize)
	pages := fixSize(5)
	key := int(math.Log2(float64(pages)))
	getTree := treePool.GetTree(pages)
	if len(treePool.trees[key]) != cacheSize-1 {
		t.Fatalf("get Tree error ,expect trees size =%d,get=%d", cacheSize-1, len(treePool.trees[key]))
	}
	treePool.PutTree(getTree)
	if len(treePool.trees[key]) != cacheSize {
		t.Fatalf("get Tree error ,expect trees size =%d,get=%d", cacheSize, len(treePool.trees[key]))
	}

}

func TestTreePool_GetPutTree_Large(t *testing.T) {

	treePool := NewTreePool(1, 1)
	pages := fixSize(10)
	//key := int(math.Log2(float64(pages)))

	//tr := buildTree(DefaultPageSize * fixSize(10))
	//
	//treePool.pool[key] = &sync.Pool{
	//	New: func() interface{} {
	//		return buildTree(pages * DefaultPageSize)
	//	},
	//}
	//treePool.pool[key].Put(tr)

	getTree := treePool.GetTree(pages)
	if getTree[0] != pages*DefaultPageSize {
		t.Fatalf("get large Tree error ,expect tree[0] =%d,get=%d", pages*DefaultPageSize, getTree[0])
	}

}

func TestTreePool_Reset_Large(t *testing.T) {

	treePool := NewTreePool(poolSize, cacheSize)
	pages := fixSize(10)
	key := int(math.Log2(float64(pages)))

	getTree := treePool.GetTree(pages)
	if getTree[0] != pages*DefaultPageSize {
		t.Fatalf("get large Tree error ,expect tree[0] =%d,get=%d", pages*DefaultPageSize, getTree[0])
	}

	getTree[0] = 12
	treePool.PutTree(getTree)
	expect := int(math.Pow(2, float64(key)) * DefaultPageSize)
	for _, v := range treePool.trees[key] {
		if v[0] != expect {
			t.Fatalf("reset error ,expect tree[0]=%d,get=%d", expect, v[0])
		}
	}

}
func TestTreePool_Reset(t *testing.T) {

	treePool := NewTreePool(poolSize, cacheSize)
	pages := fixSize(5)
	key := int(math.Log2(float64(pages)))

	getTree := treePool.GetTree(pages)
	if getTree[0] != pages*DefaultPageSize {
		t.Fatalf("get large Tree error ,expect tree[0] =%d,get=%d", pages*DefaultPageSize, getTree[0])
	}

	getTree[0] = 12
	if len(treePool.trees[key]) != cacheSize-1 {
		t.Fatalf("get Tree error ,expect trees size =%d,get=%d", cacheSize-1, len(treePool.trees[key]))
	}

	treePool.PutTree(getTree)
	if len(treePool.trees[key]) != cacheSize {
		t.Fatalf("get Tree error ,expect trees size =%d,get=%d", cacheSize, len(treePool.trees[key]))
	}

	expect := int(math.Pow(2, float64(key)) * DefaultPageSize)
	for _, v := range treePool.trees[key] {
		if v[0] != expect {
			t.Fatalf("reset error ,expect tree[0]=%d,get=%d", expect, v[0])
		}
	}

}
