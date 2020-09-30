// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

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
