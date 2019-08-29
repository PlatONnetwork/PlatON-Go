package node

import (
	"crypto/ecdsa"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var FakeNetEnable bool = false

var indexMock = map[int][]int{
	1:  []int{2},
	2:  []int{3},
	3:  []int{4},
	4:  []int{5},
	5:  []int{6},
	6:  []int{7},
	7:  []int{8},
	8:  []int{9},
	9:  []int{10},
	10: []int{1},
}

// MockDiscoveryNode returns to a specific network topology.
func MockDiscoveryNode(privateKey *ecdsa.PrivateKey, nodes []*discover.Node) []*discover.Node {
	selfNodeID := discover.PubkeyID(&privateKey.PublicKey)
	mockNodes := make([]*discover.Node, 0)
	ok, idxs := needAdd(selfNodeID, nodes)
	for idx, n := range nodes {
		if idxs == nil {
			break
		}
		for _, i := range idxs {
			if ok && i == (idx+1) {
				mockNodes = append(mockNodes, n)
				break
			}
		}
	}
	return mockNodes
}

// mock
func needAdd(self discover.NodeID, nodes []*discover.Node) (bool, []int) {
	selfIndex := -1
	for idx, n := range nodes {
		if n.ID.TerminalString() == self.TerminalString() {
			selfIndex = idx
			break
		}
	}
	if selfIndex == -1 {
		return false, nil
	}
	selfIndex++
	return true, indexMock[selfIndex]
}
