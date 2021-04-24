package dag

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Vertex struct {
	InDegree int
	OutEdges []int
	InEdges  []int
}

func NewVertex() *Vertex {
	vertex := &Vertex{
		InDegree: 0,
		OutEdges: make([]int, 0),
		InEdges:  make([]int, 0),
	}
	return vertex
}

type Dag struct {
	V          int
	Consumed   int
	VertexList []*Vertex
	NextList   []int
}

func NewDag(v int) *Dag {
	dag := &Dag{
		V: v,
	}
	for i := 0; i < v; i++ {
		dag.VertexList = append(dag.VertexList, NewVertex())
	}
	return dag
}

func (dag *Dag) AddEdge(from, to int) {
	dag.VertexList[from].OutEdges = append(dag.VertexList[from].OutEdges, to)
	dag.VertexList[to].InDegree++
	dag.VertexList[to].InEdges = append(dag.VertexList[to].InEdges, from)
}

func (dag *Dag) GetOutEdges(from int) []int {
	return dag.VertexList[from].OutEdges
}

func (dag *Dag) GetInEdges(to int) []int {
	return dag.VertexList[to].InEdges
}

func (dag *Dag) HasNext() bool {
	return dag.Consumed < dag.V
}

func (dag *Dag) Next() []int {
	if len(dag.NextList) == 0 {
		for i := 0; i < dag.V; i++ {
			if dag.VertexList[i].InDegree == 0 {
				dag.NextList = append(dag.NextList, i)
			}
		}
	} else {
		//need copy?
		preList := dag.NextList
		dag.NextList = []int{}
		for _, vtxIdx := range preList {
			for _, outVtxIdx := range dag.VertexList[vtxIdx].OutEdges {
				dag.VertexList[outVtxIdx].InDegree--
				if dag.VertexList[outVtxIdx].InDegree == 0 {
					dag.NextList = append(dag.NextList, outVtxIdx)
				}
			}
		}
	}
	dag.Consumed = dag.Consumed + len(dag.NextList)
	return dag.NextList
}

func (dag *Dag) InDegree(i int) int {
	return dag.VertexList[i].InDegree
}

func (dag *Dag) OutDegree(i int) int {
	return len(dag.VertexList[i].OutEdges)
}

func (dag *Dag) Print() (bytes.Buffer, error) {
	var buffer bytes.Buffer

	var dagCpy = &Dag{}
	if err := clone(dag, dagCpy); err != nil {
		return buffer, err
	}
	var level = 0
	for dagCpy.HasNext() {
		buffer.WriteString(fmt.Sprintf("level:%d \n    ", level))
		idxs := dagCpy.Next()
		for _, idx := range idxs {
			if level == 0 {
				buffer.WriteString(fmt.Sprintf("idx:%d ", idx))
			} else {
				buffer.WriteString(fmt.Sprintf("idx:%d%+v ", idx, dagCpy.VertexList[idx].InEdges))
			}
		}
		buffer.WriteString("\n")
		level++
	}
	return buffer, nil
}

func clone(src, dest interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(src); err != nil {
		return err
	}
	if err := dec.Decode(dest); err != nil {
		return err
	}
	return nil
}
