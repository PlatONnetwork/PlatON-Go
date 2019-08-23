package node

import (
	"crypto/ecdsa"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var FakeNetEnable bool = false

var indexMock = map[int][]int{
	1:  []int{2, 3, 4},
	2:  []int{5, 6, 7},
	3:  []int{8, 9, 10},
	4:  []int{11, 12, 13},
	5:  []int{14, 15, 16},
	6:  []int{17, 18, 19},
	7:  []int{},
	8:  []int{20, 21, 22},
	9:  []int{},
	10: []int{23, 24, 25},
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
