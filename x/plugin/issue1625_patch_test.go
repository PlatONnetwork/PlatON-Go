package plugin

import (
	"sort"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"
)

//排序顺序
//3.根据委托节点的分红比例从小到大排序，如果委托比例相同，根据节点id从小到大排序

func TestIssue1625AccountDelInfos_Sort(t *testing.T) {
	dels := make(issue1625AccountDelInfos, 0)
	//1.委托的节点已经完全退出,并且委托的时间靠后
	dels = append(dels, &issue1625AccountDelInfo{del: &staking.Delegation{DelegateEpoch: 1}, canAddr: common.NodeAddress(common.BytesToAddress([]byte{1}))})
	dels = append(dels, &issue1625AccountDelInfo{del: &staking.Delegation{DelegateEpoch: 2}, canAddr: common.NodeAddress(common.BytesToAddress([]byte{2}))})

	//2.委托的节点处于解质押状态,并且委托的时间靠后
	dels = append(dels, &issue1625AccountDelInfo{
		del: &staking.Delegation{DelegateEpoch: 2},
		candidate: &staking.Candidate{
			&staking.CandidateBase{
				NodeId: [discover.NodeIDBits / 8]byte{13},
			},
			&staking.CandidateMutable{
				Status: staking.Invalided | staking.Withdrew,
			},
		},
		canAddr: common.NodeAddress(common.BytesToAddress([]byte{3})),
	})
	//2.委托的节点处于解质押状态,并且委托的时间靠后
	dels = append(dels, &issue1625AccountDelInfo{
		del: &staking.Delegation{DelegateEpoch: 1},
		candidate: &staking.Candidate{
			&staking.CandidateBase{
				NodeId: [discover.NodeIDBits / 8]byte{11},
			},
			&staking.CandidateMutable{
				Status: staking.Invalided | staking.Withdrew,
			},
		},
		canAddr: common.NodeAddress(common.BytesToAddress([]byte{4})),
	})

	//3.根据委托节点的分红比例从小到大排序，如果委托比例相同，根据节点id从小到大排序
	dels = append(dels, &issue1625AccountDelInfo{
		del: &staking.Delegation{DelegateEpoch: 2},
		candidate: &staking.Candidate{
			&staking.CandidateBase{
				NodeId: [discover.NodeIDBits / 8]byte{2},
			},
			&staking.CandidateMutable{
				RewardPer: 10,
			},
		},
		canAddr: common.NodeAddress(common.BytesToAddress([]byte{5})),
	})
	dels = append(dels, &issue1625AccountDelInfo{
		del: &staking.Delegation{DelegateEpoch: 2},
		candidate: &staking.Candidate{
			&staking.CandidateBase{
				NodeId: [discover.NodeIDBits / 8]byte{1},
			},
			&staking.CandidateMutable{
				RewardPer: 10,
			},
		},
		canAddr: common.NodeAddress(common.BytesToAddress([]byte{6})),
	})

	dels = append(dels, &issue1625AccountDelInfo{
		del: &staking.Delegation{DelegateEpoch: 1},
		candidate: &staking.Candidate{
			&staking.CandidateBase{
				NodeId: [discover.NodeIDBits / 8]byte{3},
			},
			&staking.CandidateMutable{
				RewardPer: 15,
			},
		},
		canAddr: common.NodeAddress(common.BytesToAddress([]byte{7})),
	})

	sort.Sort(dels)
	order := []int{2, 1, 3, 4, 6, 5, 7}
	for i, del := range dels {
		if order[i] != int(del.canAddr.Big().Uint64()) {
			t.Error("sort fail,order seems wrong")
		}
	}
}
