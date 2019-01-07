package exec

import (
	"testing"
)

func TestMemory_Malloc(t *testing.T) {
	size := 32
	m := &Memory{
		Memory: make([]byte, size),
		Size:   size,
		Start:  0,
		tree:   make([]int, (2*size)-1),
	}
	initTree(m.tree, size)

	expect := []int{0, 4, 8, 16}

	for i := 0; i < 4; i++ {
		pos := m.Malloc((i + 1) * 2)
		if expect[i] != pos {
			t.Fatalf("malloc error,expect %d,get %d", expect[i], pos)
		}
	}
}

func TestMemory_Free(t *testing.T) {
	size := 32
	m := &Memory{
		Memory: make([]byte, size+100),
		Size:   size,
		Start:  100,
		tree:   make([]int, (2*size)-1),
	}
	initTree(m.tree, size)

	for i := 0; i < 4; i++ {
		pos := m.Malloc((i + 1) * 2)

		if pos != 100 {
			t.Fatalf("malloc error,expect 100,get %d", pos)
		}

		e := m.Free(pos)
		if e != nil {
			t.Fatalf("free error")
		}
	}
}

func TestMemory_Free2(t *testing.T) {
	size := 32
	m := &Memory{
		Memory: make([]byte, size+100),
		Size:   size,
		Start:  100,
		tree:   make([]int, (2*size)-1),
	}
	initTree(m.tree, size)

	originTree := m.tree
	offsets := make([]int, 0)

	for i := 0; i < 4; i++ {
		pos := m.Malloc((i + 1) * 2)
		offsets = append(offsets, pos)
	}

	for _, offset := range offsets {
		e := m.Free(offset)
		if e != nil {
			t.Fatalf("free error")
		}
	}

	if !Compare(originTree, m.tree) {
		t.Fatalf("expect tree %v, get %v", originTree, m.tree)
	}
}

func Compare(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
