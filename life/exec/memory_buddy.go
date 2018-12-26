package exec

import (
	"github.com/PlatONnetwork/PlatON-Go/log"
	"fmt"
)

type Memory struct {
	Memory []byte
	Start  int //start position for malloc
	Size   int //memory size for malloc
	tree   []int
}

func (m *Memory) Malloc(size int) int {
	if size <= 0 {
		panic(fmt.Errorf("wrong Size=%d", size))
	} else {
		size = fixSize(size)
	}
	if size > m.tree[0] {
		panic(fmt.Errorf("malloc Size=%d exceed avalable memory Size", size))
	}

	/*
		find the suitable nodeSize
	*/
	index := 0
	nodeSize := 0
	for nodeSize = m.Size; nodeSize != size; nodeSize /= 2 {
		if m.tree[left(index)] >= size {
			index = left(index)
		} else {
			index = right(index)
		}
	}
	m.tree[index] = 0
	//计算节点对应的地址
	offset := (index+1)*nodeSize - m.Size

	//向上修改收到影响的父节点size大小
	for index > 0 {
		index = parent(index)
		m.tree[index] = max(m.tree[left(index)], m.tree[right(index)])
	}

	return offset + m.Start
}

func (m *Memory) Free(offset int) error {
	if offset == 0 {
		log.Debug("free offset = 0...")
		return nil
	}
	offset = offset - m.Start
	if offset < 0 || offset >= m.Size {
		panic(fmt.Errorf("error offset=%d", offset))
	}

	//最下层节点
	nodeSize := 1
	//offset对应得节点索引
	index := offset + m.Size - 1
	//从最后的节点开始一直往上找到size为0的节点，即当初分配块所适配的大小和位置
	for ; m.tree[index] != 0; index = parent(index) {
		nodeSize *= 2
		if index == 0 {
			return nil
		}
	}

	//恢复节点
	m.tree[index] = nodeSize

	//清除节点对应的内存数据
	clear(offset+m.Start, offset+m.Start+nodeSize, m.Memory)

	//向上遍历恢复影响的节点
	var leftNode int
	var rightNode int
	for index = parent(index); index >= 0; index = parent(index) {
		nodeSize *= 2
		leftNode = m.tree[left(index)]
		rightNode = m.tree[right(index)]
		if leftNode+rightNode == nodeSize {
			m.tree[index] = nodeSize
		} else {
			m.tree[index] = max(leftNode, rightNode)
		}
	}

	return nil
}

func clear(start, end int, mem []byte) {
	for i := start; i < end; i++ {
		mem[i] = 0
	}
}

/**
计算当前节点的index计算左叶子结点的index
*/
func left(index int) int {
	return index*2 + 1
}

/**
计算当前节点的index计算右叶子结点的index
*/
func right(index int) int {
	return index*2 + 2
}

/**
计算当前节点的index计算左叶子结点的index
*/
func parent(index int) int {
	return ((index)+1)/2 - 1
}

func max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

/**
判断是否是2的次幂数
*/
func isPowOf2(n int) bool {
	if n <= 0 {
		return false
	}
	return n&(n-1) == 0
}

/*
  获取大于size的最小2的次幂数
*/
func fixSize(size int) int {

	result := 1
	for result < size {
		result = result << 1
	}
	return result
}
