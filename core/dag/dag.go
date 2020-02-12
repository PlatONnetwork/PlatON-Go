package dag

type Vertex struct {
	inDegree int
	outEdges []int
}

func NewVertex() *Vertex {
	vertex := &Vertex{
		inDegree: 0,
		outEdges: make([]int, 0),
	}
	return vertex
}

type Dag struct {
	v          int
	consumed   int
	vertexList []*Vertex
	nextList   []int
}

func NewDag(v int) *Dag {
	dag := &Dag{
		v: v,
	}
	for i := 0; i < v; i++ {
		dag.vertexList = append(dag.vertexList, NewVertex())
	}
	return dag
}

func (dag *Dag) AddEdge(from, to int) {
	dag.vertexList[from].outEdges = append(dag.vertexList[from].outEdges, to)
	dag.vertexList[to].inDegree++
}

func (dag *Dag) GetEdges(from int) []int {
	return dag.vertexList[from].outEdges
}

func (dag *Dag) HasNext() bool {
	return dag.consumed < dag.v
}

func (dag *Dag) Next() []int {
	if len(dag.nextList) == 0 {
		for i := 0; i < dag.v; i++ {
			if dag.vertexList[i].inDegree == 0 {
				dag.nextList = append(dag.nextList, i)
			}
		}
	} else {
		//need copy?
		preList := dag.nextList
		dag.nextList = []int{}
		for _, vtxIdx := range preList {
			for _, outVtxIdx := range dag.vertexList[vtxIdx].outEdges {
				dag.vertexList[outVtxIdx].inDegree--
				if dag.vertexList[outVtxIdx].inDegree == 0 {
					dag.nextList = append(dag.nextList, outVtxIdx)
				}
			}
		}
	}
	dag.consumed = dag.consumed + len(dag.nextList)
	return dag.nextList
}

func (dag *Dag) InDegree(i int) int {
	return dag.vertexList[i].inDegree
}

func (dag *Dag) OutDegree(i int) int {
	return len(dag.vertexList[i].outEdges)
}
