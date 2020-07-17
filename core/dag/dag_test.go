package dag

import (
	"fmt"
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

func TestDag(t *testing.T) {

	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	dag := NewDag(10)
	dag.AddEdge(0, 1)
	dag.AddEdge(0, 2)
	dag.AddEdge(3, 4)
	dag.AddEdge(3, 5)
	dag.AddEdge(1, 6)
	dag.AddEdge(2, 6)
	dag.AddEdge(4, 6)
	dag.AddEdge(5, 6)
	dag.AddEdge(6, 7)
	dag.AddEdge(7, 8)
	dag.AddEdge(7, 9)

	buff, err := dag.Print()
	if err != nil {
		fmt.Print("print DAG Graph error!", err)
	}
	fmt.Printf("DAG Graph for blockNumber:%d\n%s", 1, buff.String())

	fmt.Printf("iterate over second times")
	for dag.HasNext() {
		ids := dag.Next()
		fmt.Printf("ids:%+v", ids)
	}

}
