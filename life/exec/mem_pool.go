package exec

import (
	"math"
	"sync"
)

var (
	EmptyPage = make([]byte, DefaultPageSize)
)

func init() {
	for i := 0; i < len(EmptyPage); i++ {
		EmptyPage[i] = 0
	}
}

type MemBlock struct {
	FreeMem [][]byte
	memPool *sync.Pool
	size    int
	pages   int
}

type MemPool struct {
	sync.Mutex
	memBlock []*MemBlock
	largeMem map[int]*sync.Pool
}

func NewMemBlock(size, pages int) *MemBlock {
	block := &MemBlock{size: size, pages: pages}

	block.FreeMem = make([][]byte, 0, size)
	for i := 0; i < size; i++ {
		block.FreeMem = append(block.FreeMem, make([]byte, DefaultPageSize*pages))
	}
	block.memPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, DefaultPageSize*pages)
		},
	}
	return block
}

func (mb *MemBlock) Get() []byte {

	length := len(mb.FreeMem)
	if length == 0 {
		return mb.memPool.Get().([]byte)
	}
	m := mb.FreeMem[length-1]
	mb.FreeMem = mb.FreeMem[0 : length-1]
	return m
}

func (mb *MemBlock) Put(mem []byte) {
	if len(mem) != mb.pages*DefaultPageSize {
		panic("add wrong mem")
	}

	if len(mb.FreeMem) == mb.size {
		mb.memPool.Put(mem)
		return
	}

	mb.FreeMem = append(mb.FreeMem, mem)
}

func NewMemPool(count int, size int) *MemPool {
	pool := &MemPool{}
	pool.memBlock = make([]*MemBlock, 0, count)
	for i := 0; i < count; i++ {
		pool.memBlock = append(pool.memBlock, NewMemBlock(size, DefaultMemoryPages+int(math.Pow(2, float64(i)))))
	}
	pool.largeMem = make(map[int]*sync.Pool)
	return pool
}

func (mp *MemPool) Get(pages int) []byte {
	mp.Lock()
	defer mp.Unlock()
	if pages <= 0 {
		return nil
	}
	var mem []byte
	pages = fixSize(pages - DefaultMemoryPages)

	pos := int(math.Log2(float64(pages)))
	if pos >= len(mp.memBlock) {

		pool, ok := mp.largeMem[pages]
		if !ok {
			pool = &sync.Pool{
				New: func() interface{} {
					return make([]byte, DefaultPageSize*(pages+DefaultMemoryPages))
				},
			}
			mp.largeMem[pages] = pool
		}
		mem = pool.Get().([]byte)
	} else {
		mem = mp.memBlock[pos].Get()
	}

	memset(mem)

	return mem
}

func memset(mem []byte) {
	pages := len(mem) / DefaultPageSize
	for i := 0; i < pages; i++ {
		copy(mem[i*DefaultPageSize:(i+1)*DefaultPageSize], EmptyPage)
	}
}

func (mp *MemPool) Put(mem []byte) {
	mp.Lock()
	defer mp.Unlock()
	pages := len(mem) / DefaultPageSize

	pages = fixSize(pages - DefaultMemoryPages)
	pos := int(math.Log2(float64(pages)))
	if pos >= len(mp.memBlock) {
		pool, ok := mp.largeMem[pages]
		if !ok {
			return
		}
		pool.Put(mem)
	} else {
		mp.memBlock[pos].Put(mem)
	}
}
